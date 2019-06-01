<!--
Dialog is the dialog box that opens when someone clicks "Buy Tickets".
-->

<template lang="pug">
b-modal(ref="modal" :title="title" no-close-on-backdrop hide-footer @shown="onShown")
  Confirmation(v-if="orderID" :key="seq" :orderID="orderID" @close="onClose")
  OrderForm(v-else :key="seq" ref="form" :products="products" @success="onOrderSuccess")
</template>

<script>
import Confirmation from './Confirmation'
import OrderForm from './OrderForm'

export default {
  components: { Confirmation, OrderForm },
  props: {
    products: Array,
    title: String,
  },
  data: () => ({
    orderID: null,
    seq: 0,
  }),
  methods: {
    onClose() { this.$refs.modal.hide() },
    onOrderSuccess(orderID) { this.orderID = orderID },
    onShown() { this.$refs.form.setAutoFocus() },
    show() {
      this.seq++
      this.orderID = null
      this.$refs.modal.show()
    },
  },
}
</script>

<style lang="stylus"></style>
