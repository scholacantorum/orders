import '@babel/polyfill'
import 'mutationobserver-shim'
import Vue from 'vue'
import '@/plugins/axios'
import '@/plugins/bootstrap-vue'
import '@/plugins/vuelidate'
import BuyTickets from './BuyTickets'

Vue.config.productionTip = false

window.addEventListener('load', () => {
  const elms = document.getElementsByClassName('buy-tickets')
  while (elms.length) {
    const elm = elms[0]
    const productIDs = elm.getAttribute('data-products').split(' ')
    new Vue({
      render(h) {
        return h(BuyTickets, { props: {
          title: elm.getAttribute('data-title'),
          productIDs,
        } })
      }
    }).$mount(elm)
    // The mounted Vue component does not have the buy-tickets class, so it is
    // effectively shifted from elms.
  }
})
