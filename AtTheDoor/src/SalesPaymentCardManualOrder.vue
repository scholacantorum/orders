<!--
SalesPaymentCardManualOrder places the order using the card token that was
manually entered.  It returns the revised order.
-->

<template>
  <View :style="{ flex: 1, justifyContent: 'space-between' }">
    <Summary :order="order"/>
    <View :style="{ margin: 12 }">
      <Text
        :style="{ fontSize: 24, fontWeight: 'bold', color: '#f90', lineHeight: 30, textAlign: 'center' }"
      >Processing payment...</Text>
      <ActivityIndicator size="large" :style="{ marginTop: 12 }"/>
    </View>
    <View/>
  </View>
</template>

<script>
import { Alert } from 'react-native'
import Summary from './SalesSummary'
import backend from './backend'

export default {
  components: { Summary },
  props: {
    order: Object,
    token: String,
    onSuccess: Function,
    onFailure: Function,
  },
  async mounted() {
    try {
      let order = {
        source: this.order.source,
        lines: this.order.lines,
        payments: [{ type: 'card', method: this.token, amount: this.order.payments[0].amount }],
      }
      order = await backend.placeOrder(order)
      this.onSuccess(order)
    } catch (err) {
      Alert.alert('Order Error', err, () => { this.onFailure() })
    }
  },
}
</script>
