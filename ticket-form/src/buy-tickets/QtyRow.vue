<!--
QtyRow displays a single product and gets the quantity of it.
-->

<template lang="pug">
tr.buy-tickets-qty-row
  template(v-if="message")
    td.buy-tickets-qty-row-name(v-text="name")
    td(colspan="4").buy-tickets-qty-row-message(v-if="message" v-text="message")
  template(v-else)
    td.buy-tickets-qty-row-name(colspan="2" v-text="name")
    td.buy-tickets-qty-row-qty-cell
      b-form-input.buy-tickets-qty-row-qty(
        ref="input" :value="value || ''" :state="valid"
        type="number" placeholder="0" min="0"
        @input="$emit('input', Math.max(parseInt($event) || 0), 0)"
      )
    td.buy-tickets-qty-row-price &times;${{ price / 100 }}
    td.buy-tickets-qty-row-amount(v-text="value ? `$${price * value / 100}` : ''")
</template>

<script>
export default {
  props: {
    autofocus: Boolean,
    message: String,
    name: String,
    price: Number,
    value: Number,
    valid: Boolean,
  },
  methods: {
    focus() { this.$refs.input.focus() },
    onQtyChange(evt) { this.$emit('input', parseInt(evt.target.value) || 0) },
  },
}
</script>

<style lang="stylus">
.buy-tickets-qty-row-name
  padding 2px 0
  vertical-align middle
  line-height 1.2
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
