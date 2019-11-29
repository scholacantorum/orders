<!--
OrderForm is the recording order form displayed in the dialog box.
-->

<template lang="pug">
b-form(novalidate @submit.prevent="onSubmit")
  OrderLines(:products="products")
  OrderPayment(ref="pmt"
    :send="onSend" :stripeKey="stripeKey" :total="total" :user="user"
    @cancel="onCancel" @submitted="onSubmitted" @submitting="onSubmitting"
  )
</template>

<script>
import OrderLines from './OrderLines'
import OrderPayment from './OrderPayment'

export default {
  components: { OrderLines, OrderPayment },
  props: {
    auth: String,
    products: Array,
    stripeKey: String,
    user: Object,
  },
  data: () => ({
    submitted: false,
    submitting: false,
  }),
  computed: {
    lines() {
      return this.products.map(p => ({
        product: p.id,
        price: p.price,
        quantity: 1,
      }))
    },
    total() {
      if (!this.lines) return null
      return this.lines.reduce((t, ol) => t + ol.quantity * (ol.price || 0), 0)
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
    , async onSend({ subtype, method }) {
      const result = await this.$axios.post(
        `/payapi/order?auth=${encodeURIComponent(this.auth)}`,
        JSON.stringify({
          source: 'members',
          name: this.user.name,
          email: this.user.email,
          address: this.user.address,
          city: this.user.city,
          state: this.user.state,
          zip: this.user.zip,
          member: this.user.id,
          lines: this.lines.filter(ol => ol.quantity && !ol.message).map(ol => ({
            product: ol.product,
            quantity: ol.quantity,
            price: ol.price,
          })),
          payments: [{ type: 'card', subtype, method, amount: this.total }],
        }),
        { headers: { 'Content-Type': 'application/json' } },
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
      return `We’re sorry, but we're unable to process payment cards at the moment.  Please try again later, or call our office at (650) 254-1700 to order by phone.`
    },
    setAutoFocus() { this.$refs.pmt.focus() },
  },
}
</script>

<style lang="stylus">
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
