<!--
Ticket:Quantity displays a ticket product name and gets the quantity for it.
-->

<template>
  <View>
    <View :style="{ flexDirection: 'row', alignItems: 'center', marginBottom: 12 }">
      <Text :style="{ fontSize: 24, flex: 1, marginLeft: 6 }">{{ product.name }}</Text>
      <Button
        :bstyle="{ minWidth: 36, width: '5%', paddingTop: 6, paddingBottom: 6 }"
        :disabled="disabled"
        title="–"
        :onPress="onSellDown"
      />
      <Text
        :style="{ fontSize: 36, fontWeight: 'bold', minWidth: 36, width: '5%', textAlign: 'center' }"
      >{{ sell || '0' }}</Text>
      <Button
        :bstyle="{ minWidth: 36, width: '5%', paddingTop: 6, paddingBottom: 6 }"
        :disabled="disabled"
        title="+"
        :onPress="onSellUp"
      />
      <Text :style="priceStyle">{{ priceFmt }}</Text>
    </View>
    <View v-if="showUse" :style="{ flexDirection: 'row', alignItems: 'center', marginBottom: 12 }">
      <Text :style="{ fontSize: 24, flex: 1, marginRight: 6, textAlign: 'right' }">and use</Text>
      <Button
        :bstyle="{ minWidth: 36, width: '5%', paddingTop: 6, paddingBottom: 6 }"
        :disabled="disabled"
        title="–"
        :onPress="onUseDown"
      />
      <Text
        :style="{ fontSize: 36, fontWeight: 'bold', minWidth: 36, width: '5%', textAlign: 'center' }"
      >{{ use || '0' }}</Text>
      <Button
        :bstyle="{ minWidth: 36, width: '5%', paddingTop: 6, paddingBottom: 6 }"
        :disabled="disabled"
        title="+"
        :onPress="onUseUp"
      />
      <Text :style="priceStyle">&nbsp;</Text>
    </View>
  </View>
</template>

<script>
import Button from './Button'

export default {
  components: { Button },
  props: {
    disabled: Boolean,
    product: Object,
    count: Number,
    sell: Number,
    use: Number,
    onChange: Function,
  },
  computed: {
    priceStyle() {
      return {
        fontSize: 24,
        marginRight: 6,
        width: 60,
        textAlign: 'right',
        color: this.sell ? '#000' : '#888',
      }
    },
    priceFmt() { return '$' + (this.value ? this.value * this.product.price : this.product.price) / 100 },
    showUse() { return (this.sell && this.count > 1) || this.sell > 1 },
  },
  methods: {
    onSellDown() {
      if (!this.sell) return
      const sell = this.sell - 1
      const use = Math.min(sell * this.count, this.use)
      this.onChange(this.product, sell, use)
    },
    onSellUp() {
      this.onChange(this.product, (this.sell || 0) + 1, (this.use || 0) + 1)
    },
    onUseDown() {
      if (!this.use) return
      this.onChange(this.product, this.sell, this.use - 1)
    },
    onUseUp() {
      if (this.use >= this.count * this.sell) return
      this.onChange(this.product, this.sell, (this.use || 0) + 1)
    },
  }
}
</script>
