<!--
SalesPaymentCash executes the cash payment flow for the order given to it.
-->

<template>
  <View>
    <Summary :order="order" />
    <View
      :style="{ alignSelf: 'center', width: 384, flexDirection: 'row', alignItems: 'center', marginTop: 12, marginLeft: 12, marginRight: 12 }"
    >
      <Text :style="{ fontSize: 20, width: 160 }">Amount Due</Text>
      <Text :style="{ fontSize: 20, width: 80, textAlign: 'right' }">${{due}}</Text>
    </View>
    <View
      :style="{ alignSelf: 'center', width: 384, flexDirection: 'row', alignItems: 'center', marginTop: 12, marginLeft: 12, marginRight: 12 }"
    >
      <Text :style="{ fontSize: 20, width: 160 }">Cash Received</Text>
      <TextInput
        :value="received ? received.toString() : ''"
        :style="{ fontSize: 20, width: 80, padding: 5, textAlign: 'right', borderWidth: 1, borderColor: '#ccc', marginLeft: 6 }"
        :onChangeText="setReceived"
      ></TextInput>
      <Button
        v-for="v in rounded"
        :key="v"
        :title="`$${v}`"
        :bstyle="{ fontSize: 20, marginLeft: 6, backgroundColor: '#fff', borderColor: '#00f', color: '#00f', borderWidth: 1 }"
        :onPress="() => setReceived(v)"
      />
    </View>
    <View
      :style="{ alignSelf: 'center', width: 384, flexDirection: 'row', alignItems: 'center', margin: 12 }"
    >
      <Text :style="{ fontSize: 20, width: 160 }">Change</Text>
      <Text :style="{ fontSize: 20, width: 80, textAlign: 'right', marginRight: 6 }">${{change}}</Text>
      <Switch :value="donate" :disabled="!change" :onValueChange="onDonateChange" />
      <Text :style="{ fontSize: 20 }">Donation</Text>
    </View>
    <View :style="{ flexDirection: 'row', justifyContent: 'space-evenly', marginBottom: 12 }">
      <Button
        :bstyle="{ width: '40%' }"
        :disabled="!canConfirm"
        :title="confirmed ? 'Saving...' : 'Confirm'"
        :onPress="onConfirm"
      />
      <Button :bstyle="{ width: '40%' }" secondary title="Cancel" :onPress="onDone" />
    </View>
  </View>
</template>

<script>
import { Alert } from 'react-native'
import Button from './Button'
import Summary from './SalesSummary'
import backend from './backend'

export default {
  components: { Button, Summary },
  props: {
    order: Object,
    onDone: Function,
  },
  data: () => ({ confirmed: false, donate: false, received: 0 }),
  mounted() {
    this.received = this.order.payments[0].amount / 100
  },
  computed: {
    canConfirm() { return (this.received >= this.due) && !this.confirmed },
    change() {
      if (this.received < this.due) return 0
      return this.received - this.due
    },
    due() { return this.order.payments[0].amount / 100 },
    rounded() {
      const a = this.order.payments[0].amount / 100
      const list = []
      const a5 = Math.ceil(a / 5) * 5
      if (a5 !== a) list.push(a5)
      const a20 = Math.ceil(a / 20) * 20
      if (a20 !== a5) list.push(a20)
      return list
    },
  },
  methods: {
    async onConfirm() {
      this.confirmed = true
      try {
        let order = { ...this.order }
        if (this.change && this.donate) {
          order.lines = [...order.lines]
          order.lines.push({
            product: 'donation',
            quantity: 1,
            price: this.change * 100,
          })
          order.payments = [{ ...order.payments[0] }]
          order.payments[0].amount += this.change * 100
        }
        const revised = await backend.placeOrder(order)
        this.$store.commit('sold', {
          count: this.order.lines.reduce((accum, line) => line.quantity + accum, 0),
          amount: revised.payments[0].amount,
          method: revised.payments[0].method,
        })
        this.onDone()
      } catch (err) {
        this.confirmed = false
        Alert.alert('Server error', err)
      }
    },
    onDonateChange(v) { this.donate = v },
    setReceived(v) {
      if (v === '') this.received = 0
      else if (parseInt(v)) this.received = parseInt(v)
    },
  },
}
</script>
