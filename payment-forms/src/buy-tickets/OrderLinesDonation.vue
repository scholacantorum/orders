<!--
OrderLinesDonation displays the row of the order form that asks for a donation.
-->

<template lang="pug">
#buy-tickets-donation-row
  #buy-tickets-donation-label Add a donation to&nbsp;your&nbsp;order?
  #buy-tickets-donation-amount-cell
    | $
    |
    b-form-input#buy-tickets-donation-amount(
      :value="value || ''" :disabled="disabled"
      type="number" placeholder="0" min="0"
      @input="$emit('input', Math.max(parseInt($event) || 0), 0)"
      @focus="$event.target.select()"
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
  display flex
  justify-content space-between
  align-items center
  width 100%
#buy-tickets-donation-label
  flex 1 0 0
  margin-right 0.5em
  line-height 1
#buy-tickets-donation-amount-cell
  flex none
  margin-right calc(-0.75em - 1px) // padding and border of input box
  text-align right
  white-space nowrap
  font-weight bold
#buy-tickets-donation-amount
  display inline
  width 4em
  text-align right
  font-weight bold
  -moz-appearance textfield
  &::-webkit-outer-spin-button, &::-webkit-inner-spin-button
    margin 0
    -webkit-appearance none
</style>
