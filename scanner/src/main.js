import '@babel/polyfill'
import 'mutationobserver-shim'
import Vue from 'vue'
import './plugins/axios'
import './plugins/bootstrap-vue'
import Main from './Main.vue'

Vue.config.productionTip = false

new Vue({
  render: h => h(Main),
}).$mount('#main')
