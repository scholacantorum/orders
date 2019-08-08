<!--
OrderDiscount displays the line asking for an order discount coupon.
-->

<template lang="pug">
#buy-tickets-discount-row
  td#buy-tickets-discount-label Discount code
  b-form-input#buy-tickets-discount-code(
    :value="coupon" :disabled="disabled" :state="couponMatch === false ? false : null"
    @change="$emit('coupon', $event.toUpperCase())"
  )
</template>

<script>
export default {
  props: {
    coupon: String,
    couponMatch: Boolean,
    disabled: Boolean,
  },
}
</script>

<style lang="stylus">
#buy-tickets-discount-row
  display flex
  justify-content space-between
  align-items center
  margin-bottom 16px
  width 100%
#buy-tickets-discount-label
  flex none
#buy-tickets-discount-code
  flex none
  width 9em
</style>
