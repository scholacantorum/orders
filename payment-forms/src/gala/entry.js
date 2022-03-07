import '@babel/polyfill'
import 'mutationobserver-shim'
import Vue from 'vue'
import '@/plugins/axios'
import '@/plugins/bootstrap-vue'
import Gala from './Gala'

Vue.config.productionTip = false

window.addEventListener('load', () => {
  const elm = document.getElementById('gala-registration')
  if (!elm) return
  const productID = elm.getAttribute('data-product')
  new Vue({
    render(h) {
      return h(Gala, {
        props: {
          ordersURL: elm.getAttribute('data-orders-url'),
          galaRegURL: elm.getAttribute('data-gala-reg-url'),
          stripeKey: elm.getAttribute('data-stripe-key'),
          productID,
        }
      })
    }
  }).$mount(elm)
})
