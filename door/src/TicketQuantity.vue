<!--
TicketQuantity displays a ticket product name and gets the quantity for it.
-->

<template lang="pug">
div
  .tqty
    .tqty-name-box
      .tqty-name(v-text="productName")
      .tqty-subname(v-if="productSubName" v-text="productSubName")
    b-button.tqty-button(variant="primary" @click="onSellDown") –
    .tqty-qty(v-text="sell || '0'")
    b-button.tqty-button(variant="primary" @click="onSellUp") +
    div(:class="priceClass" v-text="priceFmt")
  .tqty(v-if="showUse")
    .tqty-and-use and use
    b-button.tqty-button(variant="info" @click="onUseDown") –
    .tqty-qty(v-text="use || '0'")
    b-button.tqty-button(variant="info" @click="onUseUp") +
    div.tqty-amount &nbsp;
</template>

<script>
export default {
  props: {
    product: Object,
    count: Number,
    sell: Number,
    use: Number,
  },
  computed: {
    priceClass() { return this.sell ? 'tqty-amount' : 'tqty-price' },
    priceFmt() {
      if (this.sell) return `$${this.sell * this.product.price / 100}`
      return `$${this.product.price / 100}`
    },
    productName() { return this.product.name.split('\n')[0] },
    productSubName() {
      const parts = this.product.name.split('\n')
      return parts.length > 1 ? parts[1] : ''
    },
    showUse() { return (this.sell && this.count > 1) || this.sell > 1 },
  },
  methods: {
    onSellDown() {
      if (!this.sell) return
      const sell = this.sell - 1
      const use = Math.min(sell * this.count, this.use)
      this.$emit('change', { product: this.product, sell, use })
    },
    onSellUp() {
      this.$emit('change', { product: this.product, sell: (this.sell || 0) + 1, use: (this.use || 0) + 1 })
    },
    onUseDown() {
      if (!this.use) return
      this.$emit('change', { product: this.product, sell: this.sell, use: this.use - 1 })
    },
    onUseUp() {
      if (this.use >= this.count * this.sell) return
      this.$emit('change', { product: this.product, sell: this.sell, use: (this.use || 0) + 1 })
    },
  },
}
</script>

<style lang="stylus">
.tqty
  display flex
  align-items center
.tqty-name-box
  display flex
  flex auto
  flex-direction column
  margin-left 0.5rem
  line-height 1
.tqty-name
  font-size 1.5rem
.tqty-subname
  color #888
.tqty-and-use
  flex auto
  margin-right 0.5rem
  text-align right
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
