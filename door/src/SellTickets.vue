<!--
SellTickets displays the ticket sales sequence.
-->

<template lang="pug">
SInput(v-if="!order" @order="onOrder" @cancel="$emit('done')")
PaymentOther(v-else-if="!order.id && order.payments[0].type === 'other'" :order="order" @paid="$emit('done')" @cancel="$emit('done')")
PaymentCard(v-else-if="!order.id" :order="order" @paid="onPaid" @cancel="$emit('done')")
Receipt(v-else :order="order" @done="$emit('done')")
</template>

<script>
import PaymentCard from './SalesPaymentCard'
import PaymentOther from './SalesPaymentOther'
import SInput from './SalesInput'
import Receipt from './SalesReceipt'

export default {
  components: { PaymentCard, PaymentOther, Receipt, SInput },
  data: () => ({ order: null }),
  methods: {
    onOrder(order) { this.order = order },
    onPaid(order) {
      // If we already have an email address on the order, then a receipt
      // already got sent and we shouldn't ask.  Otherwise we will.
      if (this.order.email) this.$emit('done')
      else this.order = order
    },
  },
}
</script>
