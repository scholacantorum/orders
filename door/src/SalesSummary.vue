<!--
SalesSummary
-->

<template lang="pug">
#ssummary
  #ssummary-lines
    .ssummary-line(v-for="(line, i) in lines" :key="i" v-text="line")
  #ssummary-total(v-text="`TOTAL $${salesTotal/100}`")
</template>

<script>
export default {
  props: {
    order: Object,
  },
  computed: {
    lines() {
      return this.order.lines.map(line => {
        const name = this.$store.state.products.find(p => p.id === line.product).name
        return `${line.quantity}× ${name}`
      })
    },
    salesTotal() { return this.order.payments[0].amount },
  },
}
</script>

<style lang="stylus">
#ssummary
  display flex
  justify-content space-between
  align-items flex-end
  margin-top 1.5rem
#ssummary-lines
  margin-left 0.5rem
.ssummary-line
  line-height 1.25
#ssummary-total
  padding-right 0.5rem
  font-size 1.5rem
</style>
