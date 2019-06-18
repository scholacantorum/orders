<!--
SalesSummary displays a summary of the order being placed.
-->

<template>
  <View
    :style="{ flexDirection: 'row', justifyContent: 'space-between', alignItems: 'flex-end', marginTop: 24 }"
  >
    <View :style="{ marginLeft: 6 }">
      <Text v-for="(line, i) in lines" :key="i" :style="{ fontSize: 16, lineHeight: 20 }">{{ line }}</Text>
    </View>
    <Text :style="{ fontSize: 24, paddingRight: 6 }">TOTAL ${{ salesTotal / 100 }}</Text>
  </View>
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
