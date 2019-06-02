import '@babel/polyfill'
import 'mutationobserver-shim'
import Vue from 'vue'
import '@/plugins/axios'
import '@/plugins/bootstrap-vue'
import Donate from './Donate'

Vue.config.productionTip = false

window.addEventListener('load', () => {
  const elm = document.getElementById('donate')
  if (!elm) return
  new Vue({
    render(h) {
      return h(Donate, { props: {
          ordersURL: elm.getAttribute('data-orders-url'),
          stripeKey: elm.getAttribute('data-stripe-key'),
      } })
    }
  }).$mount(elm)
})
