<!--
SalesPaymentCard handles the order flow for paying by credit card.  Only manual
entry of the card number is supported in the webapp.  (The iOS native app can
also use a card reader.)
-->

<template lang="pug">
form#paycard(@submit.prevent="onChargeCard")
  Summary(:order="order")
  #paycard-card(ref="card")
  #paycard-buttons
    b-button.paycard-button(type="submit" :disabled="processing" variant="primary" v-text="processing ? 'Charging...' : 'Charge Card'")
    b-button.paycard-button(:disabled="processing" @click="$emit('cancel')") Cancel
  #paycard-error(v-text="cardError")
</template>

<script>
import Summary from './SalesSummary'
import orderToFormData from './orderToFormData'

let stripe // handle to Stripe API, set in mounted()

export default {
  components: { Summary },
  props: {
    order: Object,
  },
  data: () => ({
    card: null,
    cardChange: null,
    cardFocus: false,
    elements: null,
    error: null,
    processing: false,
    submitted: false,
  }),
  mounted() {
    // eslint-disable-next-line
    if (!stripe) stripe = Stripe(this.$store.state.stripeKey)

    // Create the Stripe card element and set it up.
    this.elements = stripe.elements();
    this.card = this.elements.create('card', { style: { base: { fontSize: '16px', lineHeight: 1.5 } }, hidePostalCode: true });
    this.card.on('change', this.onCardChange)
    this.card.on('focus', () => { this.cardFocus = true })
    this.card.on('blur', () => { this.cardFocus = false })
    this.card.on('ready', () => { this.card.focus() })
    this.$nextTick(() => {
      this.card.mount(this.$refs.card)
    })
  },
  computed: {
    cardError() {
      if (this.error) return this.error
      if (this.cardChange && this.cardChange.error) return this.cardChange.error.message
      if (!this.submitted) return null
      if (!this.cardChange || this.cardChange.empty) return 'Please enter the payment card number.'
      if (!this.cardFocus && !this.cardChange.complete) return 'This card number is incomplete.'
      return null
    },
  },
  methods: {
    onCardChange(evt) { this.cardChange = evt },
    async onChargeCard() {
      this.submitted = true
      this.error = null
      if (this.cardError) return
      this.processing = true
      this.card.update({ disabled: true })
      const { paymentMethod, error } = await stripe.createPaymentMethod('card', this.card)
      if (error) {
        console.error(error)
        this.processing = false
        this.card.update({ disabled: false })
        if (error.type === 'card_error' || error.type === 'validation_error') this.error = error.message
        else this.error = 'Server error, unable to process payment'
        return
      }
      this.order.payments[0].method = paymentMethod.id
      try {
        const revised = (await this.$axios.post('/posapi/order', orderToFormData(this.order), {
          headers: { 'Auth': this.$store.state.auth },
        })).data
        if (revised.error) {
          this.processing = false
          console.error(revised.error)
          this.card.update({ disabled: false })
          this.error = revised.error
          return
        }
        this.$store.commit('sold', {
          count: revised.lines.reduce((accum, line) => line.quantity + accum, 0),
          amount: revised.payments[0].amount,
          method: 'card',
        })
        this.$emit('paid', revised)
      } catch (err) {
        this.processing = false
        this.card.update({ disabled: false })
        if (err.response && err.response.status === 401) {
          this.$store.commit('logout')
          window.alert('Login session expired')
          return
        }
        console.error('Error placing order', err)
        window.alert(`Server error: ${err.toString()}`)
      }
    },
  },
}
</script>

<style lang="stylus">
#paycard
  display flex
  flex-direction column
  height 100%
#paycard-card
  margin 0.75rem
  padding 0.75rem
  border 1px solid #ccc
#paycard-buttons
  display flex
  justify-content space-evenly
  margin-bottom 0.75rem
.paycard-button
  max-width 12rem
  width 40%
  font-size 1.25rem
#paycard-error
  margin 0 0.75rem
  color red
  text-align center
  font-size 1.25rem
</style>
