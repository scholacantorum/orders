import Vue from 'vue-native-core'
import Vuex from 'vuex'

Vue.use(Vuex)
const store = new Vuex.Store({
  state: {
    allow: {},
    auth: null,
    baseURL: null,
    event: null,
    products: null,
    readerState: { color: '#f33', message: 'Card reader not connected.' },
    username: null,
    admitted: 0,
    sold: 0,
    cash: 0,
    check: 0,
  },
  mutations: {
    login(state, { auth, allow, baseURL, username }) {
      state.allow = allow
      state.auth = auth
      state.baseURL = baseURL
      state.username = username
    },
    logout(state) {
      state.auth = state.baseURL = state.username = state.event = state.products = null
      state.allow = {}
      state.admitted = state.sold = state.cash = state.check = 0
    },
    event(state, { event, products }) {
      state.event = event
      state.products = products
    },
    reader(st, { state }) {
      st.readerState = state
    },
    admitted(state, count) {
      state.admitted += count
    },
    sold(state, { count, amount, method }) {
      state.sold += count
      state.admitted += count
      switch (method) {
        case 'Cash': state.cash += amount; break
        case 'Check': state.check += amount; break;
      }
    },
  }
})

Vue.prototype.$store = store
export default store

// TEMPORARY FOR TESTING
// import Stripe from 'tipsi-stripe'
// Stripe.setOptions({ publishableKey: 'pk_test_QPwvhWbGaakWn7DGcco8J5Nd' })
// store.commit('login', { auth: '9999-9999-9999', baseURL: 'http://localhost:8100/api', username: 'sroth', allow: { card: true, cash: true, willcall: true } })
