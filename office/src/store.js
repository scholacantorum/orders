import Vue from 'vue'
import Vuex from 'vuex'

Vue.use(Vuex)

const store = new Vuex.Store({
  state: {
    auth: null,
    stripeKey: null,
    username: null,
  },
  mutations: {
    login(state, { auth, stripeKey, username }) {
      state.auth = auth
      state.stripeKey = stripeKey
      state.username = username
    },
    logout(state) {
      state.auth = state.stripeKey = state.username = null
    },
  },
})
export default store

// TEMPORARY FOR TESTING
// store.commit('login', { auth: '9999-9999-9999', username: 'sroth', stripeKey: 'pk_test_QPwvhWbGaakWn7DGcco8J5Nd' })
