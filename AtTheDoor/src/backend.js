// This file contains the code that communicates with the back end server.

import Stripe from 'tipsi-stripe'
import store from './store'
import { Alert } from 'react-native'

const backend = {

  // login logs the user in with the specified username and password, in test
  // mode or not as specified.  It returns a Promise resolving to null if
  // successful or an error string if the login is invalid. It throws an error
  // message string if the login attempt failed due to a network or server
  // error.
  async login(username, password, testmode, allow) {
    let error
    const baseURL = `https://orders${testmode ? '-test' : ''}.scholacantorum.org/api`
    const url = `${baseURL}/login`
    const body = new URLSearchParams()
    body.append('username', username)
    body.append('password', password)
    const resp = await fetch(url, {
      method: 'POST',
      headers: { 'Content-Type': 'application/x-www-form-urlencoded' },
      body: body.toString(),
    }).catch((err) => {
      console.error('Error logging in', err)
      throw err.toString()
    })
    if (resp.status === 401)
      return 'Login incorrect'
    if (!resp.ok) {
      console.error('Error logging in', resp.status)
      throw `Server error ${resp.status}`
    }
    const result = await resp.json()
    if (!result.privScanTickets)
      return 'Not authorized to use this app'
    if ((allow.card || allow.cash) && !result.privInPersonSales)
      return 'Not authorized to sell tickets'
    if (allow.willcall && !result.privViewOrders)
      return 'Not authorized to view will call list'
    Stripe.setOptions({ publishableKey: result.stripePublicKey })
    store.commit('login', { auth: result.token, allow, baseURL, username })
    return null
  },

  // logout terminates the user session.  It doesn't actually communicate with the
  // server; it just forgets the auth token.
  logout() {
    store.commit('logout')
  },

  // eventlist gets the list of future events.  It throws an error message string
  // if the call fails.
  async eventlist() {
    if (!store.state.auth) throw 'Not logged in'
    const resp = await fetch(`${store.state.baseURL}/event?future=1&freeEntries=1`, {
      headers: { 'Auth': store.state.auth },
    }).catch((err) => {
      console.error('Error getting event list', err)
      throw `Server error: ${err.toString()}`
    })
    if (resp.status === 401) {
      store.commit('logout')
      throw 'Login session expired'
    }
    if (!resp.ok) {
      console.error('Error getting event list', resp.status)
      throw `Server error: ${resp.status}`
    }
    return await resp.json()
  },

  // eventProducts gets the list of ticket products, and their prices, for the
  // specified event.  It throws an error message string if the call fails.
  async eventProducts(event) {
    if (!store.state.auth) throw 'Not logged in'
    const resp = await fetch(`${store.state.baseURL}/prices?event=${encodeURIComponent(event.id)}`, {
      headers: { 'Auth': store.state.auth },
    }).catch((err) => {
      console.error('Error getting event products', err)
      throw err.toString()
    })
    if (resp.status === 401) {
      store.commit('logout')
      throw 'Login session expired'
    }
    if (!resp.ok) {
      console.error('Error getting event products', resp.status)
      throw `Server error: ${resp.status}`
    }
    return await resp.json()
  },

  // placeOrder places an order.  It throws an error message string if the call
  // fails.
  async placeOrder(order) {
    if (!store.state.auth) throw 'Not logged in'
    const resp = await fetch(`${store.state.baseURL}/order`, {
      method: 'POST',
      headers: { 'Auth': store.state.auth },
      body: JSON.stringify(order),
    }).catch((err) => {
      console.error('Error placing order', err)
      throw err.toString()
    })
    let result
    switch (resp.status) {
      case 200:
        result = await resp.json()
        if (result.error) {
          throw result.error
        }
        return result
      case 400:
        throw (await resp.text())
      case 401:
        store.commit('logout')
        throw 'Login session expired'
      default:
        console.error('Error placing order', resp.status)
        throw resp.status
    }
  },

  // cancelOrder cancels an order.  Since it is used only in error handling
  // paths, it never fails.
  async cancelOrder(orderID) {
    try {
      const resp = await fetch(`${store.state.baseURL}/order/${orderID}`, {
        method: 'DELETE',
        headers: { 'Auth': store.state.auth },
      })
      if (resp.status >= 400)
        console.error('Error canceling order', orderID, resp.status)
    } catch (err) {
      console.error('Error canceling order', orderID, err)
    }
  },

  // captureOrderPayment captures the payment of an order.  It throws an error
  // message string if the call fails.
  async captureOrderPayment(orderID) {
    if (!store.state.auth) throw 'Not logged in'
    const resp = await fetch(`${store.state.baseURL}/order/${orderID}/capturePayment`, {
      method: 'POST',
      headers: { 'Auth': store.state.auth },
    }).catch((err) => {
      console.error('Error capturing payment', err)
      throw `Server error: ${err.toString()}`
    })
    if (resp.status === 401) {
      store.commit('logout')
      throw 'Login session expired'
    }
    if (!resp.ok) {
      console.error('Error capturing payment', resp.status)
      throw `Server error: ${resp.status}`
    }
    return await resp.json()
  },

  // sendOrderReceipt sends an email receipt for an order.  It throws an error
  // message string if the call fails.
  async sendEmailReceipt(orderID, email) {
    if (!store.state.auth) throw 'Not logged in'
    const resp = await fetch(`${store.state.baseURL}/order/${orderID}/sendReceipt?email=${encodeURIComponent(email)}`, {
      method: 'POST',
      headers: { 'Auth': store.state.auth },
    }).catch((err) => {
      console.error('Error sending receipt', err)
      throw err.toString()
    })
    if (resp.status === 401) {
      store.commit('logout')
      throw 'Login session expired'
    }
    if (!resp.ok) {
      console.error('Error sending receipt', resp.status)
      throw resp.status
    }
    return null
  },

  // connectionToken gets a Stripe Terminal connection token.
  async connectionToken() {
    if (!store.state.auth) throw 'Not logged in'
    const resp = await fetch(`${store.state.baseURL}/stripe/connectTerminal`, {
      headers: { 'Auth': store.state.auth },
    }).catch((err) => {
      console.error('Error getting connection token', err)
      throw `Server error: ${err.toString()}`
    })
    if (resp.status === 401) {
      store.commit('logout')
      throw 'Login session expired'
    }
    if (!resp.ok) {
      console.error('Error getting connection token', resp.status)
      throw `Server error: ${resp.status}`
    }
    return await resp.json()
  },

  // willCallList retrieves the Will Call list for an event, i.e., a sorted list
  // of customer names and order numbers who have tickets that have been or
  // could be used at the specified event.
  async willCallList(eventID) {
    const resp = await fetch(`${store.state.baseURL}/event/${eventID}/orders`, {
      headers: { 'Auth': store.state.auth },
    }).catch((err) => {
      console.error('Error getting will call list', err)
      throw err.toString()
    })
    if (resp.status === 401) {
      store.commit('logout')
      throw 'Login session expired'
    }
    if (!resp.ok) {
      console.error('Error getting will call list', resp.status)
      throw `Server error: ${resp.status}`
    }
    return await resp.json()
  },

  // fetchTicketUsage retrieves the ticket usage counts for the specified event
  // and order.
  async fetchTicketUsage(eventID, orderID) {
    const resp = await fetch(`${store.state.baseURL}/event/${eventID}/ticket/${orderID}`, {
      headers: { 'Auth': store.state.auth },
    }).catch((err) => {
      console.error('Error fetching ticket usage', err)
      throw err.toString()
    })
    if (resp.status === 401) {
      store.commit('logout')
      throw 'Login session expired'
    }
    if (!resp.ok) {
      console.error('Error fetching ticket usage', resp.status)
      throw `Server error: ${resp.status}`
    }
    return await resp.json()
  },


  // useTickets updates the ticket usage counts for the specified event and
  // order.
  async useTickets(eventID, order) {
    const body = new URLSearchParams()
    body.append('scan', order.scan)
    order.classes.forEach(cl => {
      body.append('class', cl.name)
      body.append('used', cl.used)
    })
    const resp = await fetch(`${store.state.baseURL}/event/${eventID}/ticket/${order.id}`, {
      method: 'POST',
      body: body.toString(),
      headers: {
        'Auth': store.state.auth,
        'Content-Type': 'application/x-www-form-urlencoded',
      },
    }).catch((err) => {
      console.error('Error updating ticket usage', err)
      throw err.toString()
    })
    if (resp.status === 401) {
      store.commit('logout')
      throw 'Login session expired'
    }
    if (!resp.ok) {
      console.error('Error updating ticket usage', resp.status)
      throw `Server error: ${resp.status}`
    }
    return null
  },

}
export default backend
