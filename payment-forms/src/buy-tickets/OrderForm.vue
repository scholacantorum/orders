<!--
OrderForm is the ticket order form displayed in the dialog box.
-->

<template lang="pug">
b-form(novalidate, @submit.prevent='onSubmit')
  OrderLines(
    ref='lines',
    :disabled='submitting',
    :products='products',
    :submitted='submitted',
    @lines='onLines'
  )
  OrderDiscount(
    ref='discount',
    v-if='canDiscount',
    :coupon='coupon',
    :couponMatch='couponMatch',
    :disabled='submitting',
    @coupon='onCouponChange'
  )
  OrderPayment(
    ref='pmt',
    :couponMatch='couponMatch',
    :send='onSend',
    :stripeKey='stripeKey',
    :total='total',
    @cancel='onCancel',
    @submitted='onSubmitted',
    @submitting='onSubmitting'
  )
</template>

<script>
import OrderDiscount from './OrderDiscount'
import OrderLines from './OrderLines'
import OrderPayment from './OrderPayment'

export default {
  components: { OrderDiscount, OrderLines, OrderPayment },
  props: {
    coupon: String,
    couponMatch: Boolean,
    ordersURL: String,
    products: Array,
    stripeKey: String,
  },
  data: () => ({
    lines: null,
    submitted: false,
    submitting: false,
  }),
  computed: {
    total() {
      if (!this.lines) return null
      return this.lines.reduce((t, ol) => t + ol.quantity * (ol.price || 0), 0)
    },
    canDiscount() {
      return this.lines && !!this.lines.find(ol => !!ol.price && ol.product !== 'donation')
    },
  },
  watch: {
    submitting() { this.$emit('submitting', this.submitting) },
  },
  methods: {
    onCancel() { this.$emit('cancel') },
    onCouponChange(coupon) { this.$emit('coupon', coupon) },
    onLines(lines) { this.lines = lines },
    onSubmit() { this.$refs.pmt.submit() },
    onSubmitted() { this.submitted = true },
    onSubmitting(submitting) { this.submitting = submitting }
    , async onSend({ name, email, subtype, method }) {
      const body = new FormData()
      body.append('source', 'public')
      body.append('name', name)
      body.append('email', email)
      body.append('coupon', this.coupon)
      this.lines.filter(ol => ol.quantity && !ol.message).forEach((ol, idx) => {
        const prefix = `line${idx + 1}`
        body.append(`${prefix}.product`, ol.product)
        body.append(`${prefix}.quantity`, ol.quantity)
        body.append(`${prefix}.price`, ol.price)
      })
      if (subtype) {
        body.append('payment1.type', 'card')
        body.append('payment1.subtype', subtype)
        body.append('payment1.method', method)
        body.append('payment1.amount', this.total)
      }
      const result = await this.$axios.post(`${this.ordersURL}/payapi/order`, body,
        { headers: { 'Content-Type': 'application/x-www-form-urlencoded' } },
      ).catch(err => {
        return err
      })
      if (result && result.data && result.data.id) {
        this.$emit('success', {id: result.data.id, email: result.data.email})
        return null
      }
      if (result && result.data && result.data.error)
        return result.data.error
      console.error(result)
      return `We’re sorry, but we're unable to process your order at the moment.  Please try again later, or call our office at (650) 254-1700 to order by phone.`
    },
    setAutoFocus() { this.$refs.lines.focus() },
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
