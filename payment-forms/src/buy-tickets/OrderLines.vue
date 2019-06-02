<!--
OrderLines is the part of the ticket order form that determines what is being
bought and for how much.  Its properties are:
  - disabled: disables all input fields
  - products: the array of product information returned by the server for the
    product(s) being sold
  - submitted: a Boolean indicating whether the user has attempted to submit
    the form (after which validation errors are more blatant)
It emits:
  - lines: the list of lines of the order, if everything is valid; or null, if
    the order is invalid (i.e., all quantities are zero)
-->

<template lang="pug">
b-container.px-0(fluid)
  b-form-group.mb-0(
    label="Ticket Quantities" label-sr-only
    :state="state" invalid-feedback="How many tickets do you want?"
  )
    table#buy-tickets-form-quantities
      OrderLinesQty(
        v-for="ol in lines" ref="qty" :key="ol.product"
        v-model="ol.quantity" :name="ol.productName" :message="ol.message"
        :price="ol.price" :state="state" :disabled="disabled"
      )
  OrderLinesDonation(v-model="donation" :disabled="disabled")
  OrderLinesTotal(:total="total")
</template>

<script>
import OrderLinesDonation from './OrderLinesDonation'
import OrderLinesQty from './OrderLinesQty'
import OrderLinesTotal from './OrderLinesTotal'

export default {
  components: { OrderLinesDonation, OrderLinesQty, OrderLinesTotal },
  props: {
    disabled: Boolean,
    products: Array,
    submitted: Boolean,
  },
  data: () => ({
    donation: 0,
    lines: [],
  }),
  watch: {
    donation: 'emitLines',
    hasOne: 'emitLines',
    products: { immediate: true, handler: 'fillLines' },
  },
  computed: {
    hasOne() { return !!this.lines.some(ol => ol.quantity) },
    state() {
      if (!this.submitted || this.hasOne) return null
      return false
    },
    total() {
      return this.lines.reduce((t, ol) => t + ol.quantity * (ol.price || 0), 0) + this.donation * 100
    },
  },
  methods: {
    emitLines() {
      if (!this.hasOne) {
        this.$emit('lines', null)
      } else if (!this.donation) {
        this.$emit('lines', this.lines)
      } else {
        this.$emit('lines', [...this.lines, {
          product: 'donation', price: this.donation * 100, quantity: 1,
        }])
      }
    },
    fillLines() {
      if (!this.products) return
      let validProductCount = 0
      let lastValidProductLine
      this.lines = []
      this.products.forEach(p => {
        this.lines.push({
          product: p.id,
          productName: p.name,
          message: p.message,
          price: p.price,
          quantity: 0,
        })
        if (!p.message) {
          validProductCount++
          lastValidProductLine = this.lines[this.lines.length - 1]
        }
      })
      if (validProductCount === 1) lastValidProductLine.quantity = 1
    },
    focus() { this.$refs.qty[0].focus() },
  },
}
</script>

<style lang="stylus">
#buy-tickets-form-quantities
  width 100%
</style>
