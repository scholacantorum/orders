import Vue from 'vue-native-core'
import Vuex from 'vuex'

Vue.use(Vuex)
const store = new Vuex.Store({
  state: {
    auth: null,
    baseURL: null,
    event: null,
    products: null,
    readerState: { color: '#f33', message: 'Card reader not connected.' },
    username: null,
  },
  mutations: {
    login(state, { auth, baseURL, username }) {
      state.auth = auth
      state.baseURL = baseURL
      state.username = username
    },
    logout(state) {
      state.auth = state.baseURL = state.username = null
    },
    event(state, { event, products }) {
      state.event = event
      state.products = products
    },
    reader(st, { state }) {
      st.readerState = state
    },
  }
})

Vue.prototype.$store = store
export default store

// TEMPORARY FOR TESTING
// import Stripe from 'tipsi-stripe'
// Stripe.setOptions({ publishableKey: 'pk_test_QPwvhWbGaakWn7DGcco8J5Nd' })
// store.commit('login', { auth: '9999-9999-9999', baseURL: 'http://localhost:8100/api', username: 'sroth' })
