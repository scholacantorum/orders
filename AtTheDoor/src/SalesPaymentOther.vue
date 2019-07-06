<!--
SalesPaymentOther executes the cash/check payment flow for the order given to
it.
-->

<template>
  <View :style="{ flex: 1, justifyContent: 'space-between' }">
    <Summary :order="order" />
    <View :style="{ margin: 12 }">
      <Text
        :style="{ fontSize: 24, fontWeight: 'bold', color: '#f90', lineHeight: 30, textAlign: 'center' }"
      >{{ confirmMessage }}</Text>
    </View>
    <View :style="{ flexDirection: 'row', justifyContent: 'space-evenly', marginBottom: 12 }">
      <Button
        :bstyle="{ width: '40%' }"
        :disabled="confirmed"
        :title="confirmed ? 'Saving...' : 'Confirmed'"
        :onPress="onConfirm"
      />
      <Button :bstyle="{ width: '40%' }" secondary title="Cancel" :onPress="onCancel" />
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
    onCancel: Function,
    onPaid: Function,
  },
  data: () => ({ confirmed: false }),
  computed: {
    confirmMessage() {
      if (this.order.payments[0].method === 'Cash')
        return `Confirm you have received $${this.order.payments[0].amount / 100} in cash.`
      else
        return `Confirm you have received a check for $${this.order.payments[0].amount / 100}.`
    },
  },
  methods: {
    async onConfirm() {
      this.confirmed = true
      try {
        const revised = await backend.placeOrder(this.order)
        this.$store.commit('sold', {
          count: revised.lines.reduce((accum, line) => line.quantity + accum, 0),
          amount: revised.payments[0].amount,
          method: revised.payments[0].method,
        })
        this.onPaid(revised)
      } catch (err) {
        Alert.alert('Server Error', err)
        this.onCancel()
      }
    },
  },
}
</script>
