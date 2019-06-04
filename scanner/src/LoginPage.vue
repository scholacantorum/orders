<!--
LoginPage displays the login page and gets the user logged in.
-->

<template lang="pug">
#login
  LogoWide
  #login-head Please log in.
  form#login-form(@submit.prevent="onSubmit")
    b-form-group(inline label="Username" label-cols="4" label-for="login-username")
      b-form-input#login-username(v-model="username" autofocus)
    b-form-group(inline label="Password" label-cols="4" label-for="login-password")
      b-form-input#login-password(v-model="password" type="password")
    b-btn#login-btn(type="submit" variant="primary") Login
  #login-error(v-show="error" v-text="error")
</template>

<script>
import LogoWide from './LogoWide'

export default {
  components: { LogoWide },
  data: () => ({ username: null, password: null, error: false }),
  methods: {
    async onSubmit() {
      const body = new URLSearchParams()
      body.append('username', this.username)
      body.append('password', this.password)
      const resp = await this.$axios.post('/api/login', body).catch(err => {
        if (err.response && err.response.status === 401)
          this.error = 'Login incorrect.'
        else {
          console.error(err)
          this.error = 'Server error.'
        }
      })
      if (!resp) return
      if (!resp.data.privScanTickets) {
        this.error = 'Not authorized to use ticket scanner.'
        return
      }
      this.$emit('auth', resp.data.token)
    }
  },
}
</script>

<style lang="stylus">
#login
  display flex
  flex-direction column
  align-items center
#login-head
  margin 32px 0 16px
  font-size 24px
#login-form
  margin-bottom 16px
  width 320px
#login-btn
  display block
  margin 0 auto 16px
#login-error
  color red
</style>
