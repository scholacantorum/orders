<!--
TicketQuantity displays a ticket product name and gets the quantity for it.
-->

<template lang="pug">
.tqty
  .tqty-name(v-text="product.name")
  b-button.tqty-button(variant="primary" @click="onDown") –
  .tqty-qty(v-text="value || '0'")
  b-button.tqty-button(variant="primary" @click="onUp") +
  div(:class="priceClass" v-text="priceFmt")
</template>

<script>
export default {
  props: {
    product: Object,
    value: Number,
  },
  computed: {
    priceClass() { return this.value ? 'tqty-amount' : 'tqty-price' },
    priceFmt() {
      if (this.value) return `$${this.value * this.product.price / 100}`
      return `$${this.product.price / 100}`
    },
  },
  methods: {
    onDown() {
      if (this.value) this.$emit('change', { product: this.product, value: this.value - 1 })
    },
    onUp() {
      this.$emit('change', { product: this.product, value: (this.value || 0) + 1 })
    },
  },
}
</script>

<style lang="stylus">
.tqty
  display flex
  align-items center
.tqty-name
  flex auto
  margin-left 0.5rem
  font-size 1.5rem
.tqty-button
  flex none
  width 50px
  height 50px
  font-size 2rem
  line-height 2rem
.tqty-qty
  flex none
  width 50px
  text-align center
  font-weight bold
  font-size 2.5rem
.tqty-amount
  margin-right 0.5rem
  width 4rem
  text-align right
  font-size 1.5rem
.tqty-price
  @extend .tqty-amount
  color #888
</style>
