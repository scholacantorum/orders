<!--
SalesQuantities displays the list of available ticket products and gets the
desired quantity of each.  It also gets the basic payment method (cash, check,
or card), since the buttons that select the payment method are the ones that
terminate the quantity selection.
-->

<template>
  <View :style="{ marginTop: 24 }">
    <TicketQuantity
      v-for="product in $store.state.products"
      :key="product.id"
      :product="product"
      :value="quantities[product.id]"
      :onChange="value => onChange(product, value)"
    />
    <View :style="{ flexDirection: 'row', justifyContent: 'flex-end', paddingTop: 6 }">
      <Text :style="{ fontSize: 24, paddingRight: 6 }">TOTAL ${{ salesTotal / 100 }}</Text>
    </View>
    <View :style="{ flexDirection: 'row', justifyContent: 'space-evenly', marginTop: 24 }">
      <Button
        :bstyle="{ minWidth: 80, width: '20%', maxWidth: 200, fontSize: 20 }"
        :disabled="!qtyTotal"
        title="Cash"
        :onPress="onCash"
      />
      <Button
        :bstyle="{ minWidth: 80, width: '20%', maxWidth: 200, fontSize: 20 }"
        :disabled="!salesTotal"
        title="Check"
        :onPress="onCheck"
      />
      <Button
        :bstyle="{ minWidth: 80, width: '20%', maxWidth: 200, fontSize: 20 }"
        :disabled="!salesTotal"
        title="Card"
        :onPress="onCard"
      />
      <Button
        :bstyle="{ minWidth: 80, width: '20%', maxWidth: 200, fontSize: 20 }"
        secondary
        title="Cancel"
        :onPress="onCancel"
      />
    </View>
  </View>
</template>

<script>
import Button from './Button'
import TicketQuantity from './TicketQuantity'

export default {
  components: { Button, TicketQuantity },
  props: {
    onCancel: Function,
    onOrder: Function,
  },
  data: () => ({
    quantities: {},
  }),
  computed: {
    qtyTotal() {
      return this.$store.state.products.reduce((accum, product) => accum + (this.quantities[product.id] || 0), 0)
    },
    salesTotal() {
      return this.$store.state.products.reduce((accum, product) => accum + product.price * (this.quantities[product.id] || 0), 0)
    },
  },
  methods: {
    onCard() { this.sendOrder('card-present') },
    onCash() { this.sendOrder('other', 'Cash') },
    onChange(product, value) {
      this.$set(this.quantities, product.id, value)
    },
    onCheck() { this.sendOrder('other', 'Check') },
    sendOrder(type, method) {
      this.onOrder({
        source: 'inperson',
        payments: [{
          type, method,
          amount: this.salesTotal,
        }],
        lines: this.$store.state.products.filter(p => this.quantities[p.id]).map(p => ({
          product: p.id,
          quantity: this.quantities[p.id],
          used: this.quantities[p.id],
          usedAt: this.$store.state.event.id,
          price: p.price,
        }))
      })
    },
  }
}
</script>
