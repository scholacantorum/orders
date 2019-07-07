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
    <View v-if="needNameEmail" :style="{ marginLeft: 6, marginRight: 6 }">
      <View :style="{ flexDirection: 'row' }">
        <Text :style="{ fontSize: 20, flex: 1 }">Name</Text>
        <TextInput
          v-model="name"
          :autoCorrect="false"
          :style="{ flex: 5, fontSize: 20, borderWidth: 1, borderColor: '#ccc' }"
          textContentType="none"
        />
      </View>
      <View :style="{ flexDirection: 'row' }">
        <Text :style="{ fontSize: 20, flex: 1 }">Email</Text>
        <TextInput
          v-model="email"
          :autoCapitalize="'none'"
          :autoCorrect="false"
          keyboardType="email-address"
          :style="{ flex: 5, fontSize: 20, borderWidth: 1, borderColor: '#ccc' }"
          textContentType="none"
          type="email"
        />
      </View>
    </View>
    <View :style="{ flexDirection: 'row', justifyContent: 'space-evenly', marginTop: 24 }">
      <Button
        :bstyle="{ minWidth: 80, width: '20%', maxWidth: 200, fontSize: 20 }"
        :disabled="!canCash"
        title="Cash"
        :onPress="onCash"
      />
      <Button
        :bstyle="{ minWidth: 80, width: '20%', maxWidth: 200, fontSize: 20 }"
        :disabled="!canCheck"
        title="Check"
        :onPress="onCheck"
      />
      <Button
        :bstyle="{ minWidth: 80, width: '20%', maxWidth: 200, fontSize: 20 }"
        :disabled="!canCard"
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
    email: '',
    name: '',
    sellqty: {},
    useqty: {},
  }),
  computed: {
    canCard() { return this.salesTotal && this.haveNameEmail },
    canCash() { return this.sellQtyTotal && this.haveNameEmail },
    canCheck() { return this.salesTotal && this.haveNameEmail },
    haveNameEmail() {
      if (!this.needNameEmail) return true
      if (!this.name) return false
      if (!this.email.match(/^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$/)) return false
      return true
    },
    needNameEmail() {
      return this.$store.state.products.some(p => (this.sellqty[p.id] || 0) > (this.useqty[p.id] || 0))
    },
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
        name: this.name,
        email: this.email,
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
