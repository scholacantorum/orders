<!--
SalesPaymentCardCreateOrder creates the order on the server, gets back the
secret token for the corresponding Stripe PaymentIntent, and retrieves that
PaymentIntent from Stripe.
-->

<template>
  <View :style="{ flex: 1, justifyContent: 'space-between' }">
    <Summary :order="order"/>
    <View :style="{ margin: 12 }">
      <Text
        :style="{ fontSize: 24, fontWeight: 'bold', color: '#f90', lineHeight: 30, textAlign: 'center' }"
      >Creating order...</Text>
      <ActivityIndicator size="large" :style="{ marginTop: 12 }"/>
    </View>
    <View/>
  </View>
</template>

<script>
import { Alert } from 'react-native'
import Summary from './SalesSummary'
import backend from './backend'
import reader from './reader'

export default {
  components: { Summary },
  props: {
    order: Object,
    onCreated: Function,
    onFailure: Function,
  },
  async mounted() {
    let order, intent
    try {
      order = await backend.placeOrder(this.order)
    } catch (err) {
      Alert.alert('Server Error', err)
      return this.onFailure()
    }
    try {
      intent = await reader.retrievePaymentIntent(order.payments[0].method)
    } catch (err) {
      Alert.alert('Server Error', err)
      backend.cancelOrder(order.id)
      return this.onFailure()
    }
    this.onCreated({ order, intent })
  },
}
</script>
