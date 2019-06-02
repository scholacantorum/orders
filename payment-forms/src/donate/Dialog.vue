<!--
Dialog is the dialog box that opens when someone clicks "Buy Tickets".
-->

<template lang="pug">
b-modal(ref="modal" title="Donation" no-close-on-backdrop hide-footer @shown="onShown" @hide="onHide")
  Confirmation(v-if="orderID" :key="seq" :orderID="orderID" @close="onClose")
  OrderForm(v-else :key="seq" ref="form"
    :ordersURL="ordersURL" :stripeKey="stripeKey"
    @success="onOrderSuccess" @cancel="onClose" @submitting="onSubmitting"
  )
</template>

<script>
import Confirmation from './Confirmation'
import OrderForm from './OrderForm'

export default {
  components: { Confirmation, OrderForm },
  props: {
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
    onHide(evt) {
      if (this.submitting) evt.preventDefault()
    },
    onOrderSuccess(orderID) {
      this.orderID = orderID
      this.submitting = false
    },
    onShown() { this.$refs.form.setAutoFocus() },
    onSubmitting(submitting) { this.submitting = submitting },
    show() {
      this.seq++
      this.$refs.modal.show()
    },
  },
}
</script>
