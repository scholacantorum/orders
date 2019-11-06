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
b-container.px-0.mb-3(fluid)
  b-form-group.mb-0(
    label="Ticket Quantities" label-sr-only
    :state="state" invalid-feedback="How many tickets do you want?"
  )
    table#buy-tickets-form-quantities
      OrderLinesQty(
        v-for="p in products" ref="qty" :key="p.id"
        v-model="quantities[p.id]" :name="p.name" :subname="p.subname" :message="p.message"
        :price="p.price" :state="state" :disabled="disabled"
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
    quantities: {},
  }),
  watch: {
    donation: 'emitLines',
    hasOne: 'emitLines',
    products: { immediate: true, handler: 'fillQuantities' },
    quantities: { deep: true, handler: 'emitLines' },
  },
  computed: {
    hasOne() {
      return this.products.some(p => !p.message && this.quantities[p.id])
    },
    state() {
      if (!this.submitted || this.hasOne) return null
      return false
    },
    total() {
      return this.products.reduce((t, p) => t + (p.message ? 0 : ((p.price * this.quantities[p.id]) || 0)), 0) + this.donation * 100
    },
  },
  methods: {
    emitLines() {
      if (!this.products || !this.hasOne) {
        this.$emit('lines', null)
        return
      }
      const lines = []
      this.products.forEach(p => {
        if (p.message) return
        lines.push({
          product: p.id,
          price: p.price,
          quantity: this.quantities[p.id] || 0,
        })
      })
      if (this.donation) {
        lines.push({
          product: 'donation', price: this.donation * 100, quantity: 1,
        })
      }
      this.$emit('lines', lines)
    },
    fillQuantities() {
      if (!this.products) return
      let validProductCount = 0
      let lastValidProduct
      let seenQty = false
      this.products.forEach(p => {
        if (!this.quantities[p.id]) this.$set(this.quantities, p.id, 0)
        if (this.quantities[p.id]) seenQty = true
        if (!p.message) {
          validProductCount++
          lastValidProduct = p.id
        }
      })
      if (!seenQty && validProductCount === 1) this.quantities[lastValidProduct] = 1
      this.emitLines()
    },
    focus() { this.$refs.qty[0].focus() },
  },
}
</script>

<style lang="stylus">
#buy-tickets-form-quantities
  width 100%
  table-layout fixed
</style>
