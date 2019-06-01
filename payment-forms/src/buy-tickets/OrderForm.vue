<!--
OrderForm is the ticket order form displayed in the dialog box.
-->

<template lang="pug">
b-form(novalidate @submit.prevent="onSubmit")
  b-container.px-0(fluid)
    b-form-group.mb-0(
      label="Ticket Quantities" label-sr-only
      :state="quantityState" invalid-feedback="How many tickets do you want?"
    )
      table#buy-tickets-form-quantities
        QtyRow(v-for="ol in order.lines" ref="qty" :key="ol.product"
          v-model="ol.quantity" :name="ol.productName" :message="ol.message" :price="ol.price" :valid="quantityState"
        )
    DonationRow(v-model="donation")
    TotalRow(:total="total")
    #buy-tickets-form-divider
    b-form-group.mb-1(
      label="Your name" label-sr-only
      :state="nameState" invalid-feedback="Please enter your name."
    )
      b-form-input(v-model.trim="order.name" placeholder="Your name")
    b-form-group.mb-1(
      label="Email address" label-sr-only
      :state="emailState"
      :invalid-feedback="order.email ? 'This is not a valid email address.' : 'Please enter your email address.'"
    )
      b-form-input(
        v-model.trim="order.email" type="email" placeholder="Email address"
        @focus="emailFocused=true" @blur="emailFocused=false"
      )
    b-form-group.mb-1(
      label="Card number" label-sr-only
      :state="cardError ? false : null" :invalid-feedback="cardError"
    )
      #buy-tickets-form-card.form-control(ref="card")
    #buy-tickets-form-footer
      #buy-tickets-form-message(v-if="message" v-text="message")
      #buy-tickets-form-buttons
        b-btn(type="button" variant="secondary" :disabled="submitting") Cancel
        b-btn#buy-tickets-form-pay-now(v-if="submitting" type="submit" variant="primary" disabled)
          b-spinner.mr-1(small)
          | Paying...
        b-btn#buy-tickets-form-pay-now(v-else type="submit" variant="primary") Pay Now
</template>

<script>
import DonationRow from './DonationRow'
import QtyRow from './QtyRow'
import TotalRow from './TotalRow'

let stripe

export default {
  components: { DonationRow, QtyRow, TotalRow },
  props: {
    products: Array,
    title: String,
  },
  data: () => ({
    card: null,
    cardChange: null,
    cardFocus: false,
    donation: 0,
    elements: null,
    emailFocused: false,
    message: null,
    order: {
      name: '', email: '', lines: [],
    },
    submitted: false,
    submitting: false,
  }),
  watch: {
    products: { immediate: true, handler: 'createOrder' },
  },
  mounted() {
    // eslint-disable-next-line
    if (!stripe) stripe = Stripe(process.env.VUE_APP_STRIPE_KEY);
    this.elements = stripe.elements();
    this.card = this.elements.create('card', { style: { base: { fontSize: '16px', lineHeight: 1.5 } } });
    this.$nextTick(() => {
      this.card.mount(this.$refs.card)
      this.card.on('change', this.onCardChange)
      this.card.on('focus', () => { this.cardFocus = true })
      this.card.on('blur', () => { this.cardFocus = false })
    })
  },
  computed: {
    cardError() {
      if (this.cardChange && this.cardChange.error) return this.cardChange.error.message
      if (!this.submitted) return null
      if (!this.cardChange || this.cardChange.empty) return 'Please enter your card number.'
      if (!this.cardFocus && !this.cardChange.complete) return 'This card number is incomplete.'
      // Incomplete entries in one of the card fields are reflected in
      // cardChange.error on blur.  This last if catches cases where one of the
      // card fields is left blank.
      return null
    },
    emailState() {
      if (!this.emailFocused && this.order.email &&
        !this.order.email.match(/^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$/))
        return false
      if (this.submitted && !this.order.email) return false
      return null
    },
    nameState() { return this.submitted && this.order.name == '' ? false : null },
    quantityState() {
      if (!this.submitted) return null
      if (this.order.lines.some(ol => ol.quantity)) return null
      return false
    },
    total() {
      return this.order.lines.reduce((t, ol) => t + ol.quantity * (ol.price || 0), 0) + this.donation * 100
    },
  },
  methods: {
    createOrder() {
      if (!this.products) return
      let validProductCount = 0
      let lastValidProductLine
      this.order.lines = []
      this.products.forEach(p => {
        this.order.lines.push({
          product: p.id,
          productName: p.name,
          message: p.message,
          price: p.price,
          quantity: 0,
        })
        if (!p.message) {
          validProductCount++
          lastValidProductLine = this.order.lines[this.order.lines.length - 1]
        }
      })
      if (validProductCount === 1) lastValidProductLine.quantity = 1
    },
    onCardChange(evt) { this.cardChange = evt },
    async onSubmit() {
      this.submitted = true
      if (this.quantityState !== null || this.nameState !== null || this.emailState !== null || this.cardError) return
      this.submitting = true
      const { paymentMethod, error } = await stripe.createPaymentMethod('card', this.card, {
        billing_details: { name: this.order.name, email: this.order.email }
      })
      if (error) {
        console.error(error)
        this.submitting = false
        if (error.type === 'card_error' || error.type === 'validation_error') this.message = error.message
        else this.message = `We’re sorry, but we're unable to process payment cards at the moment.  Please try again later, or call our office at (650) 254-1700 to order by phone.`
        return
      }
      const result = await this.$axios.post(`${process.env.VUE_APP_ORDERS_URL}/api/order`, JSON.stringify({
        source: 'public',
        name: this.order.name,
        email: this.order.email,
        lines: this.order.lines.filter(ol => ol.quantity && !ol.message).map(ol => ({
          product: ol.product,
          quantity: ol.quantity,
          price: ol.price,
        })),
        payments: [{ type: 'card', method: paymentMethod.id, amount: this.total }],
      })).catch(err => {
        return err
      })
      this.submitting = false
      if (result && result.data && result.data.id)
        this.$emit('success', result.data.id)
      else if (result && result.data && result.data.error)
        this.message = result.data.error
      else {
        console.error(result)
        this.message = `We’re sorry, but we're unable to process payment cards at the moment.  Please try again later, or call our office at (650) 254-1700 to order by phone.`
      }
    },
    setAutoFocus() { this.$refs.qty[0].focus() },
  },
}
</script>

<style lang="stylus">
#buy-tickets-form-quantities
  width 100%
#buy-tickets-form-divider
  margin-top 16px
#buy-tickets-form-card
  padding 6px 12px
#buy-tickets-form-footer
  margin-top 16px
#buy-tickets-form-message
  color red
#buy-tickets-form-buttons
  margin-top 6px
  text-align right
  button
    margin-left 8px
#buy-tickets-form-pay-now
  // fixed width so the size doesn't change when the label does
  width 110px
</style>
