<!--
SellTickets displays the ticket sales sequence.
-->

<template lang="pug">
Quantities(v-if="!order" @order="onOrder" @cancel="$emit('done')")
PaymentOther(v-else-if="!order.id && order.payments[0].type === 'other'" :order="order" @paid="$emit('done')" @cancel="$emit('done')")
PaymentCard(v-else-if="!order.id" :order="order" @paid="onPaid" @cancel="$emit('done')")
Receipt(v-else :order="order" @done="$emit('done')")
</template>

<script>
import PaymentCard from './SalesPaymentCard'
import PaymentOther from './SalesPaymentOther'
import Quantities from './SalesQuantities'
import Receipt from './SalesReceipt'

export default {
  components: { PaymentCard, PaymentOther, Quantities, Receipt },
  data: () => ({ order: null }),
  methods: {
    onOrder(order) { this.order = order },
    onPaid(order) { this.order = order },
  },
}
</script>
