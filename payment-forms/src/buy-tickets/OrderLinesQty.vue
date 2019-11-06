<!--
OrderLinesQty displays a single product and gets the quantity of it.
-->

<template lang="pug">
tr.buy-tickets-qty-row
  template(v-if="message")
    td.buy-tickets-qty-row-names
      .buy-tickets-qty-row-name(v-html="name")
      .buy-tickets-qty-row-subname(v-if="subname" v-html="subname")
    td.buy-tickets-qty-row-message(colspan="4" v-if="message" v-text="message")
  template(v-else)
    td.buy-tickets-qty-row-names(colspan="2")
      .buy-tickets-qty-row-name(v-html="name")
      .buy-tickets-qty-row-subname(v-if="subname" v-html="subname")
    td.buy-tickets-qty-row-qty-cell
      b-form-input.buy-tickets-qty-row-qty(
        ref="input" :value="value || ''" :state="state" :disabled="disabled"
        type="number" placeholder="0" min="0"
        @input="$emit('input', Math.max(parseInt($event) || 0), 0)"
        @focus="$event.target.select()"
      )
    td.buy-tickets-qty-row-price &times;${{ price / 100 }}
    td.buy-tickets-qty-row-amount(v-text="value ? `$${price * value / 100}` : ''")
</template>

<script>
export default {
  props: {
    autofocus: Boolean,
    disabled: Boolean,
    message: String,
    name: String,
    subname: String,
    price: Number,
    state: Boolean,
    value: Number,
  },
  methods: {
    focus() { this.$refs.input.focus() },
  },
}
</script>

<style lang="stylus">
.buy-tickets-qty-row-names
  padding 2px 0
  width calc(100% - 4em - 9ch)
  vertical-align middle
.buy-tickets-qty-row-name
  overflow hidden
  text-overflow ellipsis
  white-space nowrap
  line-height 1.1
.buy-tickets-qty-row-subname
  overflow hidden
  text-overflow ellipsis
  white-space nowrap
  font-size 14px
  line-height 1.1
.buy-tickets-qty-row-message
  padding 2px 0 2px 6px
  color #888
  vertical-align middle
  line-height 1.2
.buy-tickets-qty-row-qty-cell
  padding 2px 0 2px 6px
  width 4em
  vertical-align middle
.buy-tickets-qty-row-qty
  text-align right
  -moz-appearance textfield
  &::-webkit-outer-spin-button, &::-webkit-inner-spin-button
    margin 0
    -webkit-appearance none
  &.is-invalid
    // Override the red X that Bootstrap adds since there isn't room for it.
    padding-right 0.75rem
    background-image none
td.buy-tickets-qty-row-price
  padding 2px 0
  width 4ch
  color #888
  vertical-align middle
.buy-tickets-qty-row-amount
  padding 2px 0
  width 5ch
  vertical-align middle
  text-align right
  font-weight bold
</style>
