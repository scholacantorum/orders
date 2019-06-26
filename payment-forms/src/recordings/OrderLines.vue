<!--
OrderLines is the part of the recording order form that displays what is being
bought and for how much.  Its properties are:
  - products: the array of product information returned by the server for the
    product(s) being sold
-->

<template lang="pug">
b-container.px-0.mb-3(fluid)
  b-form-group.mb-0
    table#recordings-form-quantities
      tr.recordings-qty-row(v-for="p in products" :key="p.id")
        td.recordings-qty-row-name(v-text="p.name")
        td.recordings-qty-row-amount(v-text="`$${p.price / 100}`")
  #recordings-form-total-row(v-if="products.length > 1")
    #recordings-form-total(v-text="`$${total / 100}`")
</template>

<script>
export default {
  props: {
    products: Array,
  },
  computed: {
    total() { return this.products.reduce((t, p) => t + p.price, 0) },
  },
}
</script>

<style lang="stylus">
#recordings-form-quantities
  width 100%
.recordings-qty-row-name
  padding 2px 0
  width 100%
  vertical-align middle
  line-height 1.2
.recordings-qty-row-amount
  padding 2px 0
  width 5ch
  vertical-align middle
  text-align right
  font-weight bold
#recordings-form-total-row
  text-align right
#recordings-form-total
  display inline-block
  margin 4px -8px 0 0
  padding 0 8px
  min-width 4.5em
  border-top 2px solid black
  text-align right
  font-weight bold
</style>
