<!--
SalesQuantities displays the list of available ticket products and gets the
desired quantity of each.  It also gets the basic payment method (cash, check,
or card), since the buttons that select the payment method are the ones that
terminate the quantity selection.
-->

<template lang="pug">
#quantities
  TicketQuantity(v-for="product in $store.state.products" :key="product.id" :product="product" :value="quantities[product.id]" @change="onChange")
  #quantities-total(v-text="`TOTAL $${salesTotal / 100}`")
  #quantities-buttons
    b-button.quantities-button(v-if="$store.state.allow.cash" :disabled="!qtyTotal" variant="primary" @click="onCash") Cash
    b-button.quantities-button(v-if="$store.state.allow.cash" :disabled="!salesTotal" variant="primary" @click="onCheck") Check
    b-button.quantities-button(v-if="$store.state.allow.card" :disabled="!salesTotal" variant="primary" @click="onCard") Card
    b-button.quantities-button(@click="$emit('done')") Cancel
</template>

<script>
import TicketQuantity from './TicketQuantity'

export default {
  components: { TicketQuantity },
  data: () => ({ quantities: {} }),
  computed: {
    qtyTotal() {
      return this.$store.state.products.reduce((accum, product) => accum + (this.quantities[product.id] || 0), 0)
    },
    salesTotal() {
      return this.$store.state.products.reduce((accum, product) => accum + product.price * (this.quantities[product.id] || 0), 0)
    },
  },
  methods: {
    onCard() { this.sendOrder('card') },
    onCash() { this.sendOrder('other', 'Cash') },
    onCheck() { this.sendOrder('other', 'Check') },
    onChange({ product, value }) {
      this.$set(this.quantities, product.id, value)
    },
    sendOrder(type, method) {
      this.$emit('order', {
        source: 'inperson',
        payments: [{
          type, method,
          amount: this.salesTotal,
        }],
        lines: this.$store.state.products.filter(p => this.quantities[p.id]).map(p => ({
          product: p.id,
          quantity: this.quantities[p.id],
          used: this.quantities[p.id],
          usedAt: this.$store.state.event.id,
          price: p.price,
        }))
      })
    },
  },
}
</script>

<style lang="stylus">
#quantities
  display flex
  flex-direction column
  margin-top 1.5rem
#quantities-total
  align-self flex-end
  padding 0.5rem
  font-size 1.5rem
#quantities-buttons
  display flex
  justify-content space-evenly
  margin-top 1.5rem
.quantities-button
  padding 6px
  min-width 5rem
  width calc(25vw - 1rem)
  font-size 1.25rem
</style>
