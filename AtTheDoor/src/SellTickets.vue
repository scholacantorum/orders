<!--
SellTickets displays the ticket sales sequence.
-->

<template>
  <Quantities v-if="!order" :onOrder="onOrder" :onCancel="onDone"/>
  <PaymentOther
    v-else-if="!order.id && order.payments[0].type === 'other'"
    :order="order"
    :onPaid="onPaid"
    :onCancel="onDone"
  />
  <PaymentCard v-else-if="!order.id" :order="order" :onPaid="onPaid" :onCancel="onDone"/>
  <Receipt v-else :order="order" :onDone="onDone"/>
</template>

<script>
import PaymentCard from './SalesPaymentCard'
import PaymentOther from './SalesPaymentOther'
import Quantities from './SalesQuantities'
import Receipt from './SalesReceipt'

export default {
  components: { PaymentCard, PaymentOther, Quantities, Receipt },
  props: {
    onDone: Function,
  },
  data: () => ({
    order: null,
    step: 'quantities',
  }),
  methods: {
    onOrder(order) { this.order = order },
    onPaid(order, emailOK) {
      this.order = order
      if (!emailOK) this.onDone()
    },
  },
}
</script>
