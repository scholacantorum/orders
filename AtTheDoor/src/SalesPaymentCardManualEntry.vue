<!--
SalesPaymentCardManualEntry displays the manual credit card entry form.
-->

<template>
  <Summary :order="order"/>
</template>

<script>
import { Alert } from 'react-native'
import Stripe from 'tipsi-stripe'
import Summary from './SalesSummary'

export default {
  components: { Summary },
  props: {
    order: Object,
    onCancel: Function,
    onEntered: Function,
  },
  async mounted() {
    try {
      const token = await Stripe.paymentRequestWithCardForm()
      this.onEntered(token.tokenId)
    } catch (err) {
      if (err.code !== 'cancelled') Alert.alert(err.message)
      this.onCancel()
    }
  },
}
</script>
