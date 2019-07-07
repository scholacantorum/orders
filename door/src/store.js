import Vue from 'vue'
import Vuex from 'vuex'

Vue.use(Vuex)

const store = new Vuex.Store({
  state: {
    allow: {},
    auth: null,
    event: null,
    products: null,
    stripeKey: null,
    username: null,
    admitted: 0,
    sold: 0,
    cash: 0,
    check: 0,
  },
  mutations: {
    login(state, { allow, auth, stripeKey, username }) {
      state.allow = allow
      state.auth = auth
      state.stripeKey = stripeKey
      state.username = username
    },
    logout(state) {
      state.auth = state.stripeKey = state.username = state.event = state.products = null
      state.allow = {}
      state.admitted = state.sold = state.cash = state.check = 0
    },
    event(state, { event, products }) {
      state.event = event
      state.products = products
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
  },
})
export default store

// TEMPORARY FOR TESTING
// store.commit('login', { auth: '9999-9999-9999', username: 'sroth', allow: { card: true, cash: true, willcall: true }, stripeKey: 'pk_test_QPwvhWbGaakWn7DGcco8J5Nd' })
