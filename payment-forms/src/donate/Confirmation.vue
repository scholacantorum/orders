<!--
Confirmation displays the order confirmation and attempts to get the customer
connected with us.
-->

<template lang="pug">
div
  #donate-confirm
    | We have received your donation and emailed you a receipt.  We will send
    | a confirmation by postal mail for your tax records.  Thank you for
    | supporting Schola Cantorum!
  #connect-head Stay informed of Schola news!
  div
    b-spinner.connect-done.mt-1(v-show="emailSpinner" small)
    .connect-done(v-show="emailDone") ✓
    .connect-link(@click="onEmail") Subscribe to our email list
  div
    b-spinner.connect-done.mt-1(v-show="pmailSpinner" small)
    .connect-done(v-show="pmailDone") ✓
    .connect-link(@click="onPMail") Subscribe to our postal mailing list
  div
    .connect-done(v-show="facebookDone") ✓
    .connect-link(@click="onFacebook") Follow us on Facebook
  div
    .connect-done(v-show="twitterDone") ✓
    .connect-link(@click="onTwitter") Follow us on Twitter
  div(style="text-align:right")
    b-btn#connect-close(variant="primary" @click="$emit('close')") Close
</template>

<script>
export default {
  props: {
    orderID: Number,
  },
  data: () => ({
    emailDone: false,
    emailSpinner: false,
    facebookDone: false,
    pmailDone: false,
    pmailSpinner: false,
    twitterDone: false,
  }),
  methods: {
    async onEmail() {
      if (this.emailDone) return
      this.emailSpinner = true
      const resp = await this.$axios.post(`/backend/email-signup?order=${this.orderID}`).catch(err => {
        console.error(err)
        return null
      })
      this.emailSpinner = false
      if (!resp) return
      if (resp.status < 400)
        this.emailDone = true
      else
        console.error(resp.statusText)
    },
    onFacebook() {
      window.open('https://www.facebook.com/scholacantorum.org', '_blank')
      this.facebookDone = true
    },
    async onPMail() {
      if (this.pmailDone) return
      this.pmailSpinner = true
      const resp = await this.$axios.post(`/backend/mail-signup?order=${this.orderID}`).catch(err => {
        console.error(err)
        return null
      })
      this.pmailSpinner = false
      if (!resp) return
      if (resp.status < 400)
        this.pmailDone = true
      else
        console.error(resp.statusText)
    },
    onTwitter() {
      window.open('https://twitter.com/scholacantorum1', '_blank')
      this.twitterDone = true
    },
  }
}
</script>

<style lang="stylus">
#donate-confirm
  color #0153a5
  font-weight bold
  line-height 1.2
#connect-head
  margin-top 12px
  font-weight bold
.connect-done
  float left
  color #32cd32
  font-weight bold
.connect-link
  margin-left 18px
  color #0056b3
  text-decoration underline
  cursor pointer
  user-select none
#connect-close
  margin 16px 0 0 auto
</style>
