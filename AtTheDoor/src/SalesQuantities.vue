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
      :count="product.ticketCount"
      :product="product"
      :sell="sellqty[product.id]"
      :use="useqty[product.id]"
      :onChange="onChange"
    />
    <View :style="{ flexDirection: 'row', justifyContent: 'flex-end', paddingTop: 6 }">
      <Text :style="{ fontSize: 24, paddingRight: 6 }">TOTAL ${{ salesTotal / 100 }}</Text>
    </View>
    <View :style="{ flexDirection: 'row', justifyContent: 'space-evenly', marginTop: 24 }">
      <Button
        :bstyle="{ minWidth: 80, width: '20%', maxWidth: 200, fontSize: 20 }"
        :disabled="!sellQtyTotal"
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
    sellqty: {},
    useqty: {},
  }),
  computed: {
    sellQtyTotal() {
      return this.$store.state.products.reduce((accum, product) => accum + (this.sellqty[product.id] || 0), 0)
    },
    salesTotal() {
      return this.$store.state.products.reduce((accum, product) => accum + product.price * (this.sellqty[product.id] || 0), 0)
    },
  },
  methods: {
    onCard() { this.sendOrder('card-present') },
    onCash() { this.sendOrder('other', 'Cash') },
    onChange(product, sell, use) {
      this.$set(this.sellqty, product.id, sell)
      this.$set(this.useqty, product.id, use)
    },
    onCheck() { this.sendOrder('other', 'Check') },
    sendOrder(type, method) {
      this.onOrder({
        source: 'inperson',
        payments: [{
          type, method,
          amount: this.salesTotal,
        }],
        lines: this.$store.state.products.filter(p => this.sellqty[p.id]).map(p => ({
          product: p.id,
          quantity: this.sellqty[p.id],
          used: this.useqty[p.id],
          usedAt: this.$store.state.event.id,
          price: p.price,
        }))
      })
    },
  }
}
</script>
