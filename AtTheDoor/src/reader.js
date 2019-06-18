// This file contains the code that interacts with the card reader.

import { Alert } from 'react-native'
import StripeTerminal from 'react-native-stripe-terminal'
import backend from './backend'
import store from './store'

StripeTerminal.initialize({ fetchConnectionToken: backend.connectionToken })

const readerAPI = {
  discovering: false,
  isConnected: false,

  // connectReader tells this module to start trying to connect with the card
  // reader (if we aren't already trying or connected).  It is called after we
  // successfully log in (and therefore can get connection tokens).
  async connectReader() {
    const status = await StripeTerminal.getConnectionStatus()
    if (status === StripeTerminal.ConnectionStatusConnected) {
      console.log('Reader was already connected')
      store.commit('reader', {})
      this.isConnected = true
      return
    }
    store.commit('reader', { state: { color: '#fc6', message: 'Connecting to card reader' } })
    if (status === StripeTerminal.ConnectionStatusConnecting) {
      console.log('Reader connection already in progress')
      this.discovering = true
      return
    }
    try {
      console.log('Discovering readers...')
      this.discovering = true
      await StripeTerminal.discoverReaders(
        StripeTerminal.DeviceTypeChipper2X,
        StripeTerminal.DiscoveryMethodBluetoothProximity,
        false
      )
    } catch (err) {
      this.discovering = false
      console.error('Error discovering readers', err)
      Alert.alert('Card Reader Error', err.error)
      store.commit('reader', { state: { color: '#f33', message: 'Card reader not connected', press: this.connectReader } })
      return
    }
    this.discovering = false
  },

  async readersDiscovered(readers) {
    if (!readers || !readers.length) return
    const reader = readers[0]
    console.log('Discovered reader', reader)
    try {
      await StripeTerminal.connectReader(reader.serialNumber)
    } catch (err) {
      console.error('Error connecting reader', err)
      Alert.alert('Card Reader Error', err.error)
      return
    }
    console.log('Connected to reader', reader)
    this.isConnected = true
    let update
    try {
      update = await StripeTerminal.checkForUpdate()
    } catch (err) {
      console.warn('Error checking for reader updates', err)
    }
    if (update)
      Alert.alert('Card Reader Update', 'A software update is available for the connected card reader.')
    store.commit('reader', {})
    // TODO We haven't implemented code to actually apply the reader software
    // update, because there's no way to test it.  If and when Stripe actually
    // issues an update, presumably they're also enhance their simulator to
    // allow testing the update process.
  },

  didReportUnexpectedReaderDisconnect() {
    Alert.alert('Card Reader Error', 'The card reader has disconnected unexpectedly.  The software will try to reconnect.')
    this.isConnected = false
    this.discoverReaders()
  },

  didReportLowBatteryWarning() {
    Alert.alert('Card Reader Warning', 'The card reader battery is low.  Recharge it soon.')
    store.commit('reader', { state: { color: '#fc6', message: 'Card reader battery low.' } })
  },

  // disconnectReader disconnects from the card reader, if we are connected, and
  // cancels any outstanding attempt to connect to it.  It is called when we log
  // out.
  async disconnectReader() {
    if (this.discovering) {
      try {
        await StripeTerminal.abortDiscoverReaders()
        console.log('Aborted discovering readers.')
      } catch (err) {
        console.error('Error aborting discover readers: ', err.error)
      }
    }
    if (await StripeTerminal.getConnectionStatus() !== StripeTerminal.ConnectionStatusNotConnected) {
      try {
        await StripeTerminal.disconnectReader()
        console.log('Disconnected reader.')
      } catch (err) {
        console.error('Error disconnecting reader: ', err.error)
      }
    }
    this.isConnected = false
  },

  // retrievePaymentIntent retrieves a PaymentIntent from Stripe given its
  // secret token.
  async retrievePaymentIntent(secret) {
    try {
      const intent = await StripeTerminal.retrievePaymentIntent(secret)
      console.log('Retrieved payment intent', intent)
      return intent
    } catch (err) {
      console.error('Error creating payment intent', err)
      throw err.error
    }
  },

  addReaderMessageListener(f) {
    StripeTerminal.addDidRequestReaderInputListener(f)
    StripeTerminal.addDidRequestReaderDisplayMessageListener(f)
  },
  removeReaderMessageListener(f) {
    StripeTerminal.removeDidRequestReaderInputListener(f)
    StripeTerminal.removeDidRequestReaderDisplayMessageListener(f)
  },

  async collectPaymentMethod() {
    try {
      const intent = await StripeTerminal.collectPaymentMethod()
      console.log('Collected payment method', intent)
      return intent
    } catch (err) {
      if (err.code !== 2020) console.error('Error collecting payment method', err)
      throw err // whole object, not string, so we can look at code
    }
  },

  async abortCollectPaymentMethod() {
    try {
      await StripeTerminal.abortCreatePayment()
      console.log('Aborted collection of payment method')
    } catch (err) {
      // log, not error, and no throw:  not a fatal problem
      console.log('Error aborting collection of payment method', err)
    }
  },

  async processPayment() {
    try {
      const intent = await StripeTerminal.processPayment()
      console.log('Processed payment', intent)
      return intent
    } catch (err) {
      console.error('Error processing payment', err)
      throw err
    }
  },

}

StripeTerminal.addReadersDiscoveredListener(readerAPI.readersDiscovered.bind(readerAPI))
StripeTerminal.addDidReportUnexpectedReaderDisconnectListener(readerAPI.didReportUnexpectedReaderDisconnect.bind(readerAPI))
StripeTerminal.addDidReportLowBatteryWarningListener(readerAPI.didReportLowBatteryWarning.bind(readerAPI))

export default readerAPI
