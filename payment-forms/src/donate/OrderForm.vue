<!--
OrderForm is the ticket order form displayed in the dialog box.
-->

<template lang="pug">
b-form(novalidate @submit.prevent="onSubmit")
  b-form-group(:state="amountState" invalid-feedback="Please enter an amount.")
    table#donate-amount-row
      tr
        td#donate-amount-label
          label.mb-0(for="donate-amount") Donation amount?
        td#donate-amount-cell
          b-input-group(prepend="$")
            b-form-input#donate-amount(ref="amount"
              :value="amount || ''" :disabled="submitting"
              type="number" placeholder="0" min="0"
              @input="amount = Math.max(parseInt($event) || 0, 0)"
            )
  OrderPayment(ref="pmt"
    :send="onSend" :stripeKey="stripeKey" :total="(amount*100) || null"
    @cancel="onCancel" @submitted="onSubmitted" @submitting="onSubmitting"
  )
</template>

<script>
import OrderPayment from './OrderPayment'

export default {
  components: { OrderPayment },
  props: {
    ordersURL: String,
    stripeKey: String,
  },
  data: () => ({
    amount: 0,
    submitted: false,
    submitting: false,
  }),
  computed: {
    amountState() {
      if (this.submitted && !this.amount) return false
      return null
    },
  },
  watch: {
    submitting() { this.$emit('submitting', this.submitting) },
  },
  methods: {
    onCancel() { this.$emit('cancel') },
    onSubmit() { this.$refs.pmt.submit() },
    onSubmitted() { this.submitted = true },
    onSubmitting(submitting) { this.submitting = submitting }
    , async onSend({ name, email, address, city, state, zip, subtype, method }) {
      const body = new FormData()
      body.append('source', 'public')
      body.append('name', name)
      body.append('email', email)
      body.append('address', address)
      body.append('city', city)
      body.append('state', state)
      body.append('zip', zip)
      body.append('line1.product', 'donation')
      body.append('line1.quantity', 1)
      body.append('line1.price', this.amount * 100)
      body.append('payment1.type', 'card')
      body.append('payment1.subtype', subtype)
      body.append('payment1.method', method)
      body.append('payment1.amount', this.amount * 100)
      const result = await this.$axios.post(`${this.ordersURL}/payapi/order`, body,
        { headers: { 'Content-Type': 'application/x-www-form-urlencoded' } },
      ).catch(err => {
        return err
      })
      if (result && result.data && result.data.id) {
        this.$emit('success', result.data.id)
        return null
      }
      if (result && result.data && result.data.error)
        return result.data.error
      console.error(result)
      return `We’re sorry, but we're unable to process payment cards at the moment.  Please try again later, or call our office at (650) 254-1700 to donate by phone.`
    },
    setAutoFocus() { this.$refs.amount.focus() },
  },
}
</script>

<style lang="stylus">
#donate-amount-row
  width 100%
#donate-amount-label
  vertical-align middle
#donate-amount-cell
  width 8em
  vertical-align middle
  white-space nowrap
  font-weight bold
#donate-amount
  text-align right
  font-weight bold
</style>
