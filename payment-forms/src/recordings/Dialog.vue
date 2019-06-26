<!--
Dialog is the dialog box that opens when someone clicks "Place Order" for
recordings.
-->

<template lang="pug">
b-modal(ref="modal" title="Concert Recordings" no-close-on-backdrop hide-footer @shown="onShown" @hide="onHide")
  OrderForm(:key="seq" ref="form"
    :auth="auth" :products="products" :stripeKey="stripeKey" :user="user"
    @success="onOrderSuccess" @cancel="onClose" @submitting="onSubmitting"
  )
</template>

<script>
import OrderForm from './OrderForm'

export default {
  components: { OrderForm },
  props: {
    auth: String,
    products: Array,
    stripeKey: String,
    user: Object,
  },
  data: () => ({
    seq: 0,
    submitting: false,
  }),
  methods: {
    onClose() { this.$refs.modal.hide() },
    onHide(evt) {
      if (this.submitting) evt.preventDefault()
      this.$emit('close')
    },
    onOrderSuccess(orderID) {
      this.submitting = false
      this.$emit('success', orderID)
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
