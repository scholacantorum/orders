<!--
SalesPaymentCheck executes the check payment flow for the order given to it.
-->

<template>
  <View>
    <Summary :order="order" />
    <View
      :style="{ alignSelf: 'center', width: 240, flexDirection: 'row', alignItems: 'center', marginTop: 12, marginLeft: 12, marginRight: 12 }"
    >
      <Text :style="{ fontSize: 20, width: 160 }">Amount Due</Text>
      <Text :style="{ fontSize: 20, width: 80, textAlign: 'right' }">${{due}}</Text>
    </View>
    <View
      :style="{ alignSelf: 'center', width: 240, flexDirection: 'row', alignItems: 'center', marginTop: 12, marginLeft: 12, marginRight: 12 }"
    >
      <Text :style="{ fontSize: 20, width: 160 }">Check Amount</Text>
      <TextInput
        clearTextOnFocus
        keyboardType="numeric"
        :style="{ fontSize: 20, width: 80, padding: 5, textAlign: 'right', borderWidth: 1, borderColor: '#ccc', marginLeft: 6 }"
        :value="received ? received.toString() : ''"
        :onChangeText="setReceived"
      ></TextInput>
    </View>
    <View
      v-if="donation"
      :style="{ alignSelf: 'center', width: 240, flexDirection: 'row', alignItems: 'center', marginTop: 12, marginLeft: 12, marginRight: 12 }"
    >
      <Text :style="{ fontSize: 20, width: 160 }">Donation</Text>
      <Text :style="{ fontSize: 20, width: 80, textAlign: 'right' }">${{donation}}</Text>
    </View>
    <View
      :style="{ flexDirection: 'row', justifyContent: 'space-evenly', marginTop: 12, marginBottom: 12 }"
    >
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
  data: () => ({ confirmed: false, received: 0 }),
  mounted() {
    this.received = this.order.payments[0].amount / 100
  },
  computed: {
    canConfirm() { return (this.received >= this.due) && !this.confirmed },
    donation() {
      if (this.received < this.due) return 0
      return this.received - this.due
    },
    due() { return this.order.payments[0].amount / 100 },
  },
  methods: {
    async onConfirm() {
      this.confirmed = true
      try {
        let order = { ...this.order }
        if (this.donation) {
          order.lines = [...order.lines]
          order.lines.push({
            product: 'donation',
            quantity: 1,
            price: this.donation * 100,
          })
          order.payments = [{ ...order.payments[0] }]
          order.payments[0].amount += this.donation * 100
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
    setReceived(v) {
      if (v === '') this.received = 0
      else if (parseInt(v)) this.received = parseInt(v)
    },
  },
}
</script>
