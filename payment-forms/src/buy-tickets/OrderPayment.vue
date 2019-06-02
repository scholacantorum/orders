<!--
OrderPayment gets the payment method.  Its properties are:
  - send - function that will send the order to the server.  It is called with
    a hash containing:
      - name - the customer name
      - email - the customer email
      - method - the Stripe payment method ID
    This function must return a Promise that resolves to null if the order was
    placed successfully, or an error message otherwise.
  - stripeKey - the key for accessing the Stripe API
  - total - the amount to be paid, in cents, or null if the form is not ready
It emits:
  - cancel - when the order is cancelled by the user.
  - submitted - when the user attempts to submit the form.  This enables more
    blatant validation errors.
  - submitting - sent with payload true when the code is attempting to submit
    the payment, and sent with payload false when the attempt ends
It exposes:
  - submit() - to be called when the containing form is submitted.  The default
    action of the form submit should be prevented by the caller.
-->

<template lang="pug">
#buy-tickets-form-payment(v-show="canPR !== null")
  #buy-tickets-form-use-pr-div(v-if="canPR")
    label#buy-tickets-form-use-pr-label(for="buy-tickets-form-use-pr")
      | Use payment info saved {{ deviceOrBrowser }}?
    b-form-checkbox#buy-tickets-form-use-pr(v-model="usePR" switch)
  div(v-show="!usePR")
    b-form-group.mb-1(
      label="Your name" label-sr-only
      :state="nameState" invalid-feedback="Please enter your name."
    )
      b-form-input(v-model.trim="name" placeholder="Your name" :disabled="submitting")
    b-form-group.mb-1(
      label="Email address" label-sr-only
      :state="emailState"
      :invalid-feedback="email ? 'This is not a valid email address.' : 'Please enter your email address.'"
    )
      b-form-input(
        v-model.trim="email" type="email" placeholder="Email address" :disabled="submitting"
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
      b-btn(type="button" variant="secondary" :disabled="submitting" @click="onCancel") Cancel
      #buy-tickets-form-prbutton(v-show="usePR" ref="prbutton")
      b-btn#buy-tickets-form-pay-now(v-if="!usePR && submitting" type="submit" variant="primary" disabled)
        b-spinner.mr-1(small)
        | Paying...
      b-btn#buy-tickets-form-pay-now(v-if="!usePR && !submitting" type="submit" variant="primary")
        | Pay {{ total ? `$${total/100}` : 'Now' }}
</template>

<script>
let stripe

export default {
  props: {
    send: Function,
    stripeKey: String,
    total: Number,
  },
  data: () => ({
    canPR: null,
    card: null,
    cardChange: null,
    cardFocus: false,
    elements: null,
    email: '',
    emailFocused: false,
    message: null,
    name: '',
    payreq: null,
    prbutton: null,
    submitted: false,
    submitting: false,
    usePR: false,
  }),
  async mounted() {
    // eslint-disable-next-line
    if (!stripe) stripe = Stripe(this.stripeKey);
    this.elements = stripe.elements();
    this.card = this.elements.create('card', { style: { base: { fontSize: '16px', lineHeight: 1.5 } } });
    this.card.on('change', this.onCardChange)
    this.card.on('focus', () => { this.cardFocus = true })
    this.card.on('blur', () => { this.cardFocus = false })
    this.$nextTick(() => {
      this.card.mount(this.$refs.card)
    })
    this.payreq = stripe.paymentRequest({
      country: 'US', currency: 'usd',
      total: { label: 'Schola Cantorum Ticket Order', amount: 100, pending: true },
      requestPayerName: true, requestPayerEmail: true,
    })
    if (await this.payreq.canMakePayment()) {
      this.canPR = this.usePR = true
      this.payreq.on('paymentmethod', this.onPaymentMethod)
      this.prbutton = this.elements.create('paymentRequestButton', { paymentRequest: this.payreq })
      this.prbutton.on('click', this.onPRButtonClick)
      this.$nextTick(() => {
        this.prbutton.mount(this.$refs.prbutton)
      })
    } else {
      this.canPR = false
    }
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
    deviceOrBrowser() {
      const userAgent = navigator.userAgent || navigator.vendor
      if (/android|ipad|iphone|ipod|windows phone/i.test(userAgent)) return 'on device'
      return 'in browser'
    },
    emailState() {
      if (!this.emailFocused && this.email &&
        !this.email.match(/^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$/))
        return false
      if (this.submitted && !this.email) return false
      return null
    },
    nameState() { return this.submitted && this.name == '' ? false : null },
  },
  watch: {
    submitted() { if (this.submitted) this.$emit('submitted') },
    submitting() { this.$emit('submitting', this.submitting) },
    usePR() { this.submitted = false },
  },
  methods: {
    onCancel() { this.$emit('cancel') },
    onCardChange(evt) { this.cardChange = evt },
    async onPaymentMethod(evt) {
      this.submitting = true
      this.card.update({ disabled: true })
      const error = await this.send({ name: evt.payerName, email: evt.payerEmail, method: evt.paymentMethod.id })
      this.submitting = false
      this.card.update({ disabled: false })
      if (error) this.message = error
      evt.complete(error ? 'fail' : 'success')
    },
    onPRButtonClick(evt) {
      this.submitted = true
      this.message = null
      if (this.total === null) {
        evt.preventDefault()
        return
      }
      this.payreq.update({ total: { label: 'Schola Cantorum Ticket Order', amount: this.total, pending: false } })
    },
    async submit() {
      this.submitted = true
      this.message = null
      if (this.total === null || this.nameState !== null || this.emailState !== null || this.cardError) return
      this.submitting = true
      this.card.update({ disabled: true })
      const { paymentMethod, error } = await stripe.createPaymentMethod('card', this.card, {
        billing_details: { name: this.name, email: this.email }
      })
      if (error) {
        console.error(error)
        this.submitting = false
        this.card.update({ disabled: false })
        if (error.type === 'card_error' || error.type === 'validation_error') this.message = error.message
        else this.message = `We’re sorry, but we're unable to process payment cards at the moment.  Please try again later, or call our office at (650) 254-1700 to order by phone.`
        return
      }
      const error2 = await this.send({ name: this.name, email: this.email, method: paymentMethod.id })
      this.submitting = false
      this.card.update({ disabled: false })
      if (error2) this.message = error2
    },
  },
}
</script>

<style lang="stylus">
#buy-tickets-form-use-pr-div
  display flex
  justify-content space-between
  align-items center
  margin 0 -8px 6px 0
#buy-tickets-form-use-pr-label
  margin-bottom 0
  font-size 14px
#buy-tickets-form-card
  padding 6px 12px
#buy-tickets-form-footer
  margin-top 16px
#buy-tickets-form-message
  color red
#buy-tickets-form-buttons
  display flex
  justify-content flex-end
  margin-top 6px
  button
    margin-left 8px
#buy-tickets-form-prbutton
  margin-left 8px
  min-width 110px
#buy-tickets-form-pay-now
  // fixed width so the size doesn't change when the label does
  width 110px
</style>
