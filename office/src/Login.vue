<!--
Login displays the login page.
-->

<template lang="pug">
form#login(@submit.prevent="onSubmit")
  #login-header Please log in.
  b-form-group(label="Username" label-for="login-username" label-cols="4")
    b-form-input#login-username(v-model="username" autocapitalize="none" autocomplete="username" autofocus trim)
  b-form-group(label="Password" label-for="login-password" label-cols="4")
    b-form-input#login-password(v-model="password" type="password" autocomplete="password")
  #login-error(v-if="error" v-text="error")
  b-button#login-button(type="submit" variant="primary" :disabled="disabled") Login
</template>

<script>
export default {
  data: () => ({
    error: null,
    password: null,
    username: null,
  }),
  computed: {
    disabled() { return !this.password || !this.username },
  },
  methods: {
    async onSubmit() {
      const body = new URLSearchParams()
      body.append('username', this.username)
      body.append('password', this.password)
      try {
        const resp = (await this.$axios.post('/api/login', body)).data
        if (!resp.privViewOrders) {
          this.error = 'Not authorized to use this app'
          return
        }
        this.$store.commit('login', {
          auth: resp.token,
          stripeKey: resp.stripePublicKey,
          username: this.username,
          privSetupOrders: resp.privSetupOrders,
          privManageOrders: resp.privManageOrders,
        })
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
  margin 1rem auto
  width 20rem
#login-header
  align-self center
  margin-bottom 1rem
  font-weight bold
  font-size 1.25rem
#login-error
  align-self center
  margin-bottom 1rem
  color red
#login-button
  align-self center
  font-size 1.25rem
</style>
