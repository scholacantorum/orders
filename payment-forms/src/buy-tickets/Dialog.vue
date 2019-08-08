<!--
Dialog is the dialog box that opens when someone clicks "Buy Tickets".
-->

<template lang="pug">
b-modal(ref="modal" :title="title" no-close-on-backdrop hide-footer @shown="onShown" @hide="onHide")
  Confirmation(v-if="orderID" :key="seq" :orderID="orderID" @close="onClose")
  OrderForm(v-else :key="seq" ref="form"
    :coupon="coupon" :couponMatch="couponMatch" :ordersURL="ordersURL" :products="products" :stripeKey="stripeKey"
    @coupon="onCouponChange" @success="onOrderSuccess" @cancel="onClose" @submitting="onSubmitting"
  )
</template>

<script>
import Confirmation from './Confirmation'
import OrderForm from './OrderForm'

export default {
  components: { Confirmation, OrderForm },
  props: {
    coupon: String,
    couponMatch: Boolean,
    ordersURL: String,
    products: Array,
    stripeKey: String,
    title: String,
  },
  data: () => ({
    orderID: null,
    seq: 0,
    submitting: false,
  }),
  methods: {
    onClose() { this.$refs.modal.hide() },
    onCouponChange(coupon) { this.$emit('coupon', coupon) },
    onHide(evt) { if (this.submitting) evt.preventDefault() },
    onOrderSuccess(orderID) {
      this.orderID = orderID
      this.submitting = false
    },
    onShown() { this.$refs.form.setAutoFocus() },
    onSubmitting(submitting) { this.submitting = submitting },
    show() {
      this.seq++
      this.orderID = null
      this.$refs.modal.show()
    },
  },
}
</script>

<style lang="stylus"></style>
