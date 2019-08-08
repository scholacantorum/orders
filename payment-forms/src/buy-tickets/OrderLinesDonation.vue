<!--
OrderLinesDonation displays the row of the order form that asks for a donation.
-->

<template lang="pug">
table#buy-tickets-donation-row
  tr
    td#buy-tickets-donation-label Add a donation to&nbsp;your&nbsp;order?
    td#buy-tickets-donation-amount-cell
      | $
      b-form-input#buy-tickets-donation-amount(
        :value="value || ''" :disabled="disabled"
        type="number" placeholder="0" min="0"
        @input="$emit('input', Math.max(parseInt($event) || 0), 0)"
      )
</template>

<script>
export default {
  props: {
    disabled: Boolean,
    value: Number,
  },
}
</script>

<style lang="stylus">
#buy-tickets-donation-row
  width calc(100% + 4px) // line up with unboxed numbers
#buy-tickets-donation-label
  vertical-align middle
  font-size 14px
  line-height 1.2
#buy-tickets-donation-amount-cell
  width 4.5em
  vertical-align middle
  white-space nowrap
  font-weight bold
#buy-tickets-donation-amount
  display inline
  text-align right
  font-weight bold
  -moz-appearance textfield
  &::-webkit-outer-spin-button, &::-webkit-inner-spin-button
    margin 0
    -webkit-appearance none
</style>
