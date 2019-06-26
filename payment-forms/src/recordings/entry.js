import '@babel/polyfill'
import 'mutationobserver-shim'
import Vue from 'vue'
import '@/plugins/axios'
import '@/plugins/bootstrap-vue'
import 'bootstrap/dist/css/bootstrap.min.css'
import Recordings from './Recordings'

Vue.config.productionTip = false

const params = new URLSearchParams(location.search)
new Vue({
  render: h => h(Recordings, {
    props: {
      productIDs: params.getAll('p'),
      auth: params.get('auth'),
    }
  }),
}).$mount('#app')
