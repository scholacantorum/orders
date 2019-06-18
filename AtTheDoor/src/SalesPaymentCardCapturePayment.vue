<!--
SalesPaymentCardCapturePayment tells the server to capture the payment that we
just processed.  It returns the revised order.
-->

<template>
  <View :style="{ flex: 1, justifyContent: 'space-between' }">
    <Summary :order="order"/>
    <View :style="{ margin: 12 }">
      <Text
        :style="{ fontSize: 24, fontWeight: 'bold', color: '#f90', lineHeight: 30, textAlign: 'center' }"
      >Recording successful payment...</Text>
      <ActivityIndicator size="large" :style="{ marginTop: 12 }"/>
    </View>
    <View/>
  </View>
</template>

<script>
import Summary from './SalesSummary'
import backend from './backend'

export default {
  components: { Summary },
  props: {
    order: Object,
    onCaptured: Function,
    onFailure: Function,
  },
  async mounted() {
    try {
      const order = await backend.captureOrderPayment(this.order.id)
      this.onCaptured(order)
    } catch (error) {
      // We are not displaying an error here because the credit card did get
      // authorized, and we'll eventually notice the authorized but uncaptured
      // payment and capture it.  It's a problem, but not one the person doing
      // at-the-door sales should be distracted with.  We do return an error,
      // however, so that the caller knows not to offer an email receipt that we
      // can't generate.
      this.onFailure()
    }
  },
}
</script>
