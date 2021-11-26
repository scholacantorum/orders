import '@babel/polyfill'
import 'mutationobserver-shim'
import Vue from 'vue'
import '@/plugins/axios'
import '@/plugins/bootstrap-vue'
import BuyTickets from './BuyTickets'

Vue.config.productionTip = false

window.addEventListener('load', () => {
  const elms = document.getElementsByClassName('buy-tickets')
  while (elms.length) {
    const elm = elms[0]
    const productIDs = elm.getAttribute('data-products').split(' ')
    const title = elm.getAttribute('data-title')
    let buttonLabel = elm.getAttribute('data-button-label')
    if (!buttonLabel) buttonLabel = title.includes('Subscription') ? 'Buy Subscriptions' : 'Buy Tickets'
    new Vue({
      render(h) {
        return h(BuyTickets, {
          props: {
            ordersURL: elm.getAttribute('data-orders-url'),
            productIDs,
            stripeKey: elm.getAttribute('data-stripe-key'),
            title,
            buttonLabel,
          }
        })
      }
    }).$mount(elm)
    // The mounted Vue component does not have the buy-tickets class, so it is
    // effectively shifted from elms.
  }
})
