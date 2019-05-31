<!--
TotalRow displays the order total.
-->

<template lang="pug">
#buy-tickets-form-total-row
  #buy-tickets-form-total(v-text="`$${total / 100}`")
</template>

<script>
export default {
  props: {
    total: Number,
  },
}
</script>

<style lang="stylus">
#buy-tickets-form-total-row
  text-align right
#buy-tickets-form-total
  display inline-block
  margin 4px -8px 0 0
  padding 0 8px
  min-width 4.5em
  border-top 2px solid black
  text-align right
  font-weight bold
</style>
