<!--
SalesPayment executes the card payment flow for the order given to it.
-->

<template>
  <CreateOrder
    v-if="shouldCreateOrder"
    :order="order"
    :onCreated="onOrderCreated"
    :onFailure="onCancel"
  />
  <CollectPayment
    v-else-if="shouldCollectPayment"
    :order="myOrder"
    :onCancel="myOnCancel"
    :onCollected="onIntent"
    :onManualEntry="onManualEntry"
  />
  <ProcessPayment
    v-else-if="shouldProcessPayment"
    :order="myOrder"
    :onFailure="myOnCancel"
    :onProcessed="onIntent"
  />
  <CapturePayment
    v-else-if="shouldCapturePayment"
    :order="myOrder"
    :onCaptured="onCaptured"
    :onFailure="onCaptureFailed"
  />
  <ManualEntry
    v-else-if="shouldDoManualEntry"
    :order="order"
    :onCancel="onCancel"
    :onEntered="onManualToken"
  />
  <ManualOrder
    v-else-if="shouldDoManualOrder"
    :order="order"
    :token="token"
    :onFailure="onManualOrderFailure"
    :onSuccess="onManualOrderSuccess"
  />
</template>

<script>
import StripeTerminal from 'react-native-stripe-terminal'
import CapturePayment from './SalesPaymentCardCapturePayment'
import CollectPayment from './SalesPaymentCardCollectPayment'
import CreateOrder from './SalesPaymentCardCreateOrder'
import ManualEntry from './SalesPaymentCardManualEntry'
import ManualOrder from './SalesPaymentCardManualOrder'
import ProcessPayment from './SalesPaymentCardProcessPayment'
import backend from './backend'
import reader from './reader'

export default {
  components: { CapturePayment, CollectPayment, CreateOrder, ManualEntry, ManualOrder, ProcessPayment },
  props: {
    order: Object,
    onCancel: Function,
    onPaid: Function,
  },
  data: () => ({
    intent: null,
    myOrder: null,
    manual: !reader.isConnected,
    token: null,
  }),
  computed: {
    shouldCreateOrder() { return !this.manual && !this.myOrder },
    shouldCollectPayment() {
      return !this.manual && this.intent.status === StripeTerminal.PaymentIntentStatusRequiresPaymentMethod
    },
    shouldProcessPayment() {
      return !this.manual && this.intent.status === StripeTerminal.PaymentIntentStatusRequiresConfirmation
    },
    shouldCapturePayment() {
      return !this.manual && this.intent.status === StripeTerminal.PaymentIntentStatusRequiresCapture
    },
    shouldDoManualEntry() { return this.manual && !this.token },
    shouldDoManualOrder() { return this.manual && this.token },
  },
  methods: {
    myOnCancel() {
      backend.cancelOrder(this.myOrder.id)
      this.onCancel()
    },
    onCaptureFailed() { this.onPaid(this.myOrder, false) },
    onCaptured(order) { this.onPaid(order, true) },
    onIntent(intent) { this.intent = intent },
    onManualEntry() {
      backend.cancelOrder(this.myOrder.id)
      this.manual = true
    },
    onManualOrderFailure() { this.token = null },
    onManualOrderSuccess(order) { this.onPaid(order, true) },
    onManualToken(token) { this.token = token },
    onOrderCreated({ order, intent }) {
      this.myOrder = order
      this.intent = intent
    },
  }
}
</script>
