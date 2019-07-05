<!--
SalesReceipt offers to send an email receipt to the purchaser of the completed
order.
-->

<template lang="pug">
form#receipt(@submit.prevent="onSendReceipt")
  Summary(:order="order")
  b-form-input#receipt-email(v-model="email" type="email" autocomplete="off" :autoCapitalize="'none'" :autoCorrect="false" autofocus placeholder="Email address")
  #receipt-buttons
    b-button.receipt-button(type="submit" :disabled="!validEmail || sending" variant="primary" v-text="sending ? 'Sending...' : 'Send Receipt'")
    b-button.receipt-button(@click="$emit('done')") No Receipt
</template>

<script>
import Summary from './SalesSummary'

export default {
  components: { Summary },
  props: {
    order: Object,
  },
  data: () => ({ email: null, sending: false }),
  mounted() {
    this.email = this.order.email
  },
  computed: {
    validEmail() {
      return /^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$/.test(this.email)
    },
  },
  methods: {
    async onSendReceipt() {
      if (!this.validEmail || this.sending) return
      this.sending = true
      try {
        await this.$axios.post(`/api/order/${this.order.id}/sendReceipt`, null, {
          headers: { 'Auth': this.$store.state.auth },
          params: { email: this.email },
        })
        this.$emit('done')
      } catch (err) {
        this.sending = false
        if (err.response && err.response.status === 401) {
          this.$store.commit('logout')
          window.alert('Login session expired')
          return
        }
        console.error('Error sending receipt', err)
        window.alert(`Server error: ${err.toString()}`)
      }
    },
  },
}
</script>

<style lang="stylus">
#receipt
  display flex
  flex-direction column
  height 100%
#receipt-email
  margin 0.75rem
  width calc(100% - 1.5rem)
  font-size 1.25rem
#receipt-buttons
  display flex
  justify-content space-evenly
  margin-bottom 0.75rem
.receipt-button
  max-width 12rem
  width 40%
  font-size 1.25rem
</style>
