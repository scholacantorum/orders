<!--
SellTickets displays the ticket sales sequence.
-->

<template>
  <SInput v-if="!order" :onOrder="onOrder" :onCancel="onDone" />
  <PaymentCash
    v-else-if="!order.id && order.payments[0].type === 'other' && order.payments[0].method === 'Cash'"
    :order="order"
    :onDone="onDone"
  />
  <PaymentCheck
    v-else-if="!order.id && order.payments[0].type === 'other'"
    :order="order"
    :onDone="onDone"
  />
  <PaymentCard v-else-if="!order.id" :order="order" :onPaid="onPaid" :onCancel="onDone" />
  <Receipt v-else :order="order" :onDone="onDone" />
</template>

<script>
import PaymentCard from './SalesPaymentCard'
import PaymentCash from './SalesPaymentCash'
import PaymentCheck from './SalesPaymentCheck'
import SInput from './SalesInput'
import Receipt from './SalesReceipt'

export default {
  components: { PaymentCard, PaymentCash, PaymentCheck, Receipt, SInput },
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
