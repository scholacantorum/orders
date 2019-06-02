<!--
OrderPayment gets the payment method.  Its properties are:
  - send - function that will send the order to the server.  It is called with
    a hash containing:
      - name - the customer name
      - email - the customer email
      - subtype - the payment subtype information
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

In the mounted() method, this component determines whether the user's browser
supports the PaymentRequest API (in other words, Apple Pay, Google Pay, etc.).
It sets the "canPR" flag to true or false based on that.  No payment UI is shown
until this flag is set.

If PaymentRequest API is supported, we give the user the option of using it or
directly entering their payment info.  Both UIs are rendered so that switching
between them is quick, but only the selected one is visible.  The user's choice
is stored in the "usePR" flag, controlled by a switch in the UI.  It defaults to
true (i.e, use the PaymentRequest API).

If the PaymentRequest API is used, we will get a 'click' event from the payment
request button when it is clicked.  If the form passes validation, we update the
payment request with the correct order total and allow the browser to proceed
with its payment UI.  When that succeeds, we will get a 'paymentmethod' event
from the payment request button.  We get the customer name and email and the
payment method ID from that event, and tell our parent component to place the
order using those; we then notify the payment request button whether the order
was successfully placed.

If the PaymentRequest API is not used, we will get a 'submit' event from the
form when the user clicks the 'Pay Now' button.  If the form passes validation,
we will request a payment method from the card entry field, and if that is
successful, we will tell our parent component to place the order using it and
using the customer name and email from our form fields.
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
      b-form-input(v-model.trim="name" placeholder="Your name" autocomplete="name" :disabled="submitting")
    b-form-group.mb-1(
      label="Email address" label-sr-only
      :state="emailState"
      :invalid-feedback="email ? 'This is not a valid email address.' : 'Please enter your email address.'"
    )
      b-form-input(
        v-model.trim="email" type="email" placeholder="Email address" autocomplete="email" :disabled="submitting"
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
let stripe // handle to Stripe API, set in mounted()

export default {
  props: {
    send: Function,
    stripeKey: String,
    total: Number,
  },
  data: () => ({
    canPR: null,         // browser supports PaymentRequest API; null means still checking
    card: null,          // Stripe card element
    cardChange: null,    // most recent cardChange event payload from card element
    cardFocus: false,    // card element currently has focus
    elements: null,      // Stripe elements collection
    email: '',           // customer email address from form
    emailFocused: false, // email address input has focus
    message: null,       // error message after failed submission
    name: '',            // customer name from form
    payreq: null,        // Stripe payment request object
    prbutton: null,      // Stripe payment request button element
    submitted: false,    // true if submission has been attempted
    submitting: false,   // true if submission is in progress (disables all fields)
    usePR: false,        // true if user wants to use PaymentRequest API
  }),
  async mounted() {
    // eslint-disable-next-line
    if (!stripe) stripe = Stripe(this.stripeKey);

    // Create the Stripe card element and set it up.
    this.elements = stripe.elements();
    this.card = this.elements.create('card', { style: { base: { fontSize: '16px', lineHeight: 1.5 } } });
    this.card.on('change', this.onCardChange)
    this.card.on('focus', () => { this.cardFocus = true })
    this.card.on('blur', () => { this.cardFocus = false })
    this.$nextTick(() => {
      this.card.mount(this.$refs.card)
    })

    // Create the payment request and check whether the Payment Request API is
    // supported.
    this.payreq = stripe.paymentRequest({
      country: 'US', currency: 'usd',
      total: { label: 'Schola Cantorum Ticket Order', amount: 100, pending: true },
      requestPayerName: true, requestPayerEmail: true,
    })

    // If the Payment Request API is supported, create the payment request
    // button element and set it up.
    this.canPR = !!(await this.payreq.canMakePayment())
    if (this.canPR) {
      this.usePR = true
      this.payreq.on('paymentmethod', this.onPaymentMethod)
      this.prbutton = this.elements.create('paymentRequestButton', { paymentRequest: this.payreq })
      this.prbutton.on('click', this.onPRButtonClick)
      this.$nextTick(() => {
        this.prbutton.mount(this.$refs.prbutton)
      })
    }
  },
  computed: {
    cardError() {
      // Returns an error message describing the problem with the card input,
      // or null if no error message should be displayed.
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
      // This is an imperfect heuristic used to tailor the label of the usePR
      // switch.  On mobile devices it asks whether to use the payment info
      // saved on device; on desktops it asks whether to use the payment info
      // saved in browser.
      const userAgent = navigator.userAgent || navigator.vendor
      if (/android|ipad|iphone|ipod|windows phone/i.test(userAgent)) return 'on device'
      return 'in browser'
    },
    emailState() {
      // Returns false if the email field should show an error message, null if
      // it should not.
      if (!this.emailFocused && this.email &&
        !this.email.match(/^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$/))
        return false
      if (this.submitted && !this.email) return false
      return null
    },
    nameState() {
      // Returns false if the name field should show an error message, null if
      // it should not.
      return this.submitted && this.name == '' ? false : null
    },
  },
  watch: {
    submitted() { if (this.submitted) this.$emit('submitted') },
    submitting() { this.$emit('submitting', this.submitting) },
    usePR() {
      // When someone changes from payment request to manual entry, it's likely
      // because their payment request failed.  If so, the name, email, and card
      // fields are all showing error messages because they are empty.  We'll
      // clear those by pretending that the form was never submitted.
      this.submitted = false
    },
  },
  methods: {
    onCancel() { this.$emit('cancel') },
    onCardChange(evt) { this.cardChange = evt },
    async onPaymentMethod(evt) {
      // The user has completed the browser payment request, and we have a
      // payment method to use.  Ask our parent to place the order, and report
      // back the result.
      this.submitting = true
      this.card.update({ disabled: true })
      const error = await this.send({
        name: evt.payerName, email: evt.payerEmail,
        subtype: `API ${evt.methodName}`, method: evt.paymentMethod.id,
      })
      this.submitting = false
      this.card.update({ disabled: false })
      if (error) this.message = error
      evt.complete(error ? 'fail' : 'success')
    },
    onPRButtonClick(evt) {
      // The user has clicked the payment request button.  Validate the form,
      // and update the payment request with the correct order total before the
      // browser starts its payment request UI.
      this.submitted = true
      this.message = null
      if (this.total === null) {
        evt.preventDefault()
        return
      }
      this.payreq.update({ total: { label: 'Schola Cantorum Ticket Order', amount: this.total, pending: false } })
    },
    async submit() {
      // The user has clicked the Pay Now button, or pressed enter in one of
      // the form fields.  If they're using Payment Request API, or we don't
      // know yet, ignore it.
      if (this.canPR === null || this.usePR) return

      // Otherwise, validate the form.
      this.submitted = true
      this.message = null
      if (this.total === null || this.nameState !== null || this.emailState !== null || this.cardError) return

      // The form is valid, so ask the card element for a payment method.
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

      // We got a payment method, so ask our parent to place the order.
      const error2 = await this.send({
        name: this.name, email: this.email, subtype: 'manual', method: paymentMethod.id,
      })
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
