<!--
SellTickets displays the ticket sales sequence.
-->

<template lang="pug">
SInput(v-if="!order" @order="onOrder" @cancel="$emit('done')")
Receipt(v-else-if="done" :order="order" @done="$emit('done')")
PaymentCash(v-else-if="!order.id && order.payments[0].type === 'other' && order.payments[0].method === 'Cash'" :order="order" @paid="$emit('done')" @cancel="$emit('done')")
PaymentCheck(v-else-if="!order.id && order.payments[0].type === 'other'" :order="order" @paid="$emit('done')" @cancel="$emit('done')")
PaymentCard(v-else="!order.id" :order="order" @paid="onPaid" @cancel="$emit('done')")
</template>

<script>
import PaymentCard from './SalesPaymentCard'
import PaymentCash from './SalesPaymentCash'
import PaymentCheck from './SalesPaymentCheck'
import SInput from './SalesInput'
import Receipt from './SalesReceipt'

export default {
  components: { PaymentCard, PaymentCash, PaymentCheck, Receipt, SInput },
  data: () => ({ order: null, done: false }),
  methods: {
    onOrder(order) { this.order = order },
    onPaid(order) {
      this.done = true
      // If we already have an email address on the order, then a receipt
      // already got sent and we shouldn't ask.  Otherwise we will.
      if (this.order.email) this.$emit('done')
      else this.order = order
    },
  },
}
</script>
