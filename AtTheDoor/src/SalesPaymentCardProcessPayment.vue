<!--
SalesPaymentCardProcessPayment processes the payment with the payment method
collected from the card reader.  It returns the revised intent.
-->

<template>
  <View :style="{ flex: 1, justifyContent: 'space-between' }">
    <Summary :order="order" />
    <View :style="{ margin: 12 }">
      <Text
        :style="{ fontSize: 24, fontWeight: 'bold', color: '#f90', lineHeight: 30, textAlign: 'center' }"
      >Processing payment...</Text>
      <ActivityIndicator size="large" :style="{ marginTop: 12 }" />
    </View>
    <View />
  </View>
</template>

<script>
import { Alert } from 'react-native'
import StripeTerminal from 'react-native-stripe-terminal'
import Summary from './SalesSummary'
import reader from './reader'

export default {
  components: { Summary },
  props: {
    order: Object,
    onProcessed: Function,
    onFailure: Function,
  },
  async mounted() {
    try {
      const intent = await reader.processPayment()
      this.$store.commit('sold', {
        count: this.order.lines.reduce((accum, line) => line.quantity + accum, 0),
        amount: this.order.payments[0].amount,
        method: 'card',
      })
      this.onProcessed(intent)
    } catch (err) {
      Alert.alert('Payment Error', err.error)
      if (!err.intent || err.intent.status === StripeTerminal.PaymentIntentStatusRequiresConfirmation)
        this.onFailure()
      else // going back to RequiresPaymentMethod
        this.onProcessed(err.intent)
    }
  },
}
</script>
