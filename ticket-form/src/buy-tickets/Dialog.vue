<!--
Dialog is the dialog box that opens when someone clicks "Buy Tickets".
-->

<template lang="pug">
b-modal(ref="modal" :title="title" no-close-on-backdrop hide-footer @shown="onShown")
  Confirmation(v-if="orderID" :orderID="orderID" @close="onClose")
  OrderForm(v-else ref="form" :products="products" @success="onOrderSuccess")
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
  data: () => ({ orderID: 1 }),
  methods: {
    onClose() { this.$refs.modal.hide() },
    onOrderSuccess(orderID) { this.orderID = orderID },
    onShown() { this.$refs.form.setAutoFocus() },
    show() { this.$refs.modal.show() },
  },
}
</script>

<style lang="stylus"></style>
