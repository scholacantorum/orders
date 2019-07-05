<!--
Login displays the login page.
-->

<template lang="pug">
form#login(@submit.prevent="onSubmit")
  #login-header Please log in.
  b-form-group(label="Username" label-for="login-username" label-cols="4")
    b-form-input#login-username(v-model="username" autocomplete="username" autofocus trim)
  b-form-group(label="Password" label-for="login-password" label-cols="4")
    b-form-input#login-password(v-model="password" type="password" autocomplete="password")
  b-form-group#login-allow-row(label="Allow" label-cols="4")
    b-form-checkbox-group(v-model="privs" :options="possiblePrivs" stacked switches)
  #login-error(v-if="error" v-text="error")
  b-button#login-button(type="submit" variant="primary" :disabled="disabled") Login
</template>

<script>
export default {
  data: () => ({
    error: null,
    password: null,
    possiblePrivs: [
      { value: 'card', text: 'Credit/debit card sales' },
      { value: 'cash', text: 'Cash/check sales' },
      { value: 'willcall', text: 'Will call' },
    ],
    privs: ['card', 'cash', 'willcall'],
    username: null,
  }),
  computed: {
    disabled() { return !this.password || !this.username || !this.privs.length },
  },
  methods: {
    async onSubmit() {
      const body = new URLSearchParams()
      body.append('username', this.username)
      body.append('password', this.password)
      try {
        const resp = (await this.$axios.post('/api/login', body)).data
        if (!resp.privScanTickets) {
          this.error = 'Not authorized to use this app'
          return
        }
        const privs = {}
        this.privs.forEach(p => { privs[p] = true })
        if ((privs.card || privs.cash) && !resp.privInPersonSales) {
          this.error = 'Not authorized to sell tickets'
          return
        }
        if (privs.willcall && !resp.privViewOrders) {
          this.error = 'Not authorized to view will call list'
          return
        }
        this.$store.commit('login', { auth: resp.token, allow: privs, stripeKey: resp.stripePublicKey, username: this.username })
      } catch (err) {
        if (err.response && err.response.status === 401) {
          this.error = 'Login incorrect'
          return
        }
        console.error('Error logging in', err)
        this.error = 'Server error — login failed'
        return
      }
    },
  },
}
</script>

<style lang="stylus">
#login
  display flex
  flex-direction column
  align-self center
  align-items stretch
  margin 1rem 0.5rem
  min-width 20rem
#login-header
  align-self center
  margin-bottom 1rem
  font-weight bold
  font-size 1.25rem
#login-allow-row legend
  padding-top 0
#login-error
  align-self center
  margin-bottom 1rem
  color red
#login-button
  align-self center
  font-size 1.25rem
</style>
