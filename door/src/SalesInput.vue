<!--
SalesInput displays the list of available ticket products and gets the desired
purchase and usage quantity of each.  If not all purchased tickets are being
consumed, it gets the customer's name and email address.  Finally, it gets the
basic payment method for the order (cash, check, or card).  It returns the
formatted order object.
-->

<template lang="pug">
#sinput
  TicketQuantity(v-for="product in $store.state.products" :key="product.id" :count="product.ticketCount" :product="product" :sell="sellqty[product.id]" :use="useqty[product.id]" @change="onChange")
  #sinput-total(v-text="`TOTAL $${salesTotal / 100}`")
  #sinput-name-email(v-if="needNameEmail")
    b-form-group(label="Name" label-for="sinput-name" label-cols="2")
      b-form-input#sinput-name(v-model.trim="name" autocomplete="off")
    b-form-group(label="Email" label-for="sinput-email" label-cols="2")
      b-form-input#sinput-email(v-model="email" type="email" autocomplete="off")
  #sinput-buttons
    b-button.sinput-button(v-if="$store.state.allow.cash" :disabled="!canCash" variant="primary" @click="onCash") Cash
    b-button.sinput-button(v-if="$store.state.allow.cash" :disabled="!canCheck" variant="primary" @click="onCheck") Check
    b-button.sinput-button(v-if="$store.state.allow.card" :disabled="!canCard" variant="primary" @click="onCard") Card
    b-button.sinput-button(@click="$emit('cancel')") Cancel
</template>

<script>
import TicketQuantity from './TicketQuantity'

export default {
  components: { TicketQuantity },
  data: () => ({ email: '', name: '', sellqty: {}, useqty: {} }),
  computed: {
    canCard() { return this.salesTotal && this.haveNameEmail },
    canCash() { return this.sellQtyTotal && this.haveNameEmail },
    canCheck() { return this.salesTotal && this.haveNameEmail },
    haveNameEmail() {
      if (!this.needNameEmail) return true
      if (!this.name) return false
      if (!this.email.match(/^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$/)) return false
      return true
    },
    needNameEmail() {
      return this.$store.state.products.some(p => (this.sellqty[p.id] || 0) * p.ticketCount > (this.useqty[p.id] || 0))
    },
    sellQtyTotal() {
      return this.$store.state.products.reduce((accum, product) => accum + (this.sellqty[product.id] || 0), 0)
    },
    salesTotal() {
      return this.$store.state.products.reduce((accum, product) => accum + product.price * (this.sellqty[product.id] || 0), 0)
    },
  },
  methods: {
    onCard() { this.sendOrder('card', 'manual') },
    onCash() { this.sendOrder('other', 'cash') },
    onCheck() { this.sendOrder('other', 'check') },
    onChange({ product, sell, use }) {
      this.$set(this.sellqty, product.id, sell)
      this.$set(this.useqty, product.id, use)
    },
    sendOrder(type, subtype) {
      this.$emit('order', {
        source: 'inperson',
        name: this.name,
        email: this.email,
        payments: [{
          type, subtype,
          amount: this.salesTotal,
        }],
        lines: this.$store.state.products.filter(p => this.sellqty[p.id]).map(p => ({
          product: p.id,
          quantity: this.sellqty[p.id],
          used: this.useqty[p.id],
          usedAt: this.$store.state.event.id,
          price: p.price,
        }))
      })
    },
  },
}
</script>

<style lang="stylus">
#sinput
  display flex
  flex-direction column
  margin-top 1.5rem
#sinput-total
  align-self flex-end
  padding 0.5rem
  font-size 1.5rem
#sinput-name-email
  margin 0 0.5rem
  font-size 1.25rem
#sinput-buttons
  display flex
  justify-content space-evenly
  margin-top 1.5rem
.sinput-button
  padding 6px
  min-width 5rem
  width calc(25vw - 1rem)
  font-size 1.25rem
</style>
