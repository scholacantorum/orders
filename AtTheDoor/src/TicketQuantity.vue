<!--
Ticket:Quantity displays a ticket product name and gets the quantity for it.
-->

<template>
  <View :style="{ flexDirection: 'row', alignItems: 'center', marginBottom: 12 }">
    <Text :style="{ fontSize: 24, flex: 1, marginLeft: 6 }">{{ product.name }}</Text>
    <Button
      :bstyle="{ minWidth: 36, width: '5%', paddingTop: 6, paddingBottom: 6 }"
      :disabled="disabled"
      title="–"
      :onPress="onDown"
    />
    <Text
      :style="{ fontSize: 36, fontWeight: 'bold', minWidth: 36, width: '5%', textAlign: 'center' }"
    >{{ value || '0' }}</Text>
    <Button
      :bstyle="{ minWidth: 36, width: '5%', paddingTop: 6, paddingBottom: 6 }"
      :disabled="disabled"
      title="+"
      :onPress="onUp"
    />
    <Text :style="priceStyle">{{ priceFmt }}</Text>
  </View>
</template>

<script>
import Button from './Button'

export default {
  components: { Button },
  props: {
    disabled: Boolean,
    product: Object,
    value: Number,
    onChange: Function,
  },
  computed: {
    priceStyle() {
      return {
        fontSize: 24,
        marginRight: 6,
        width: 60,
        textAlign: 'right',
        color: this.value ? '#000' : '#888',
      }
    },
    priceFmt() { return '$' + (this.value ? this.value * this.product.price : this.product.price) / 100 },
  },
  methods: {
    onDown() {
      if (this.value) this.onChange(this.value - 1)
    },
    onUp() {
      this.onChange((this.value || 0) + 1)
    },
  }
}
</script>
