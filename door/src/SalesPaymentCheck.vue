<!--
SalesPaymentCheck executes the check payment flow for the order given to it.
-->

<template lang="pug">
#paycheck
  Summary(:order="order")
  .paycheck-row
    .paycheck-label Amount Due
    #paycheck-due(v-text="`$${due}`")
  .paycheck-row
    .paycheck-label Check Amount
    input#paycheck-received(v-model="received" autocomplete="off" :formatter="numFormatter" inputmode="numeric" pattern="[0-9]*")
  .paycheck-row(v-if="donation")
    .paycheck-label Donation
    #paycheck-donation(v-text="`$${donation}`")
  #paycheck-buttons
    b-button.paycheck-button(:disabled="!canConfirm" variant="primary" @click="onConfirmed" v-text="confirmed ? 'Saving...' : 'Confirm'")
    b-button.paycheck-button(:disabled="confirmed" @click="$emit('cancel')") Cancel
</template>

<script>
import Summary from './SalesSummary'

export default {
  components: { Summary },
  props: {
    order: Object,
  },
  data: () => ({ confirmed: false, received: 0 }),
  mounted() {
    this.received = this.order.payments[0].amount / 100
  },
  computed: {
    canConfirm() { return (this.received >= this.due) && !this.confirmed },
    donation() {
      if (this.received < this.due) return 0
      return this.received - this.due
    },
    due() { return this.order.payments[0].amount / 100 },
  },
  methods: {
    numFormatter(text) { return text.replace(/[^0-9]/g, '') },
    async onConfirmed() {
      this.confirmed = true
      try {
        let order = { ...this.order }
        if (this.donation) {
          order.lines = [...order.lines]
          order.lines.push({
            product: 'donation',
            quantity: 1,
            price: this.donation * 100,
          })
          order.payments = [{ ...order.payments[0] }]
          order.payments[0].amount += this.donation * 100
        }
        const revised = (await this.$axios.post('/api/order', order, {
          headers: { 'Auth': this.$store.state.auth },
        })).data
        if (revised.error) throw revised.error
        this.$store.commit('sold', {
          count: this.order.lines.reduce((accum, line) => line.quantity + accum, 0),
          amount: revised.payments[0].amount,
          method: revised.payments[0].method,
        })
        this.$emit('paid', revised)
      } catch (err) {
        this.confirmed = false
        if (err.response && err.response.status === 401) {
          this.$store.commit('logout')
          window.alert('Login session expired')
          return
        }
        console.error('Error placing order', err)
        window.alert(`Server error: ${err.toString()}`)
      }
    },
    setReceived(v) { this.received = v },
  },
}
</script>

<style lang="stylus">
#paycheck
  display flex
  flex-direction column
.paycheck-row
  display flex
  align-self center
  align-items center
  margin 0.75rem 0.75rem 0
  width 15rem
.paycheck-label
  width 10rem
  font-size 1.25rem
#paycheck-due, #paycheck-donation, #paycheck-received
  width 5rem
  text-align right
  font-size 1.25rem
#paycheck-received
  margin-left 5px
  padding 5px
  border 1px solid #ccc
#paycheck-buttons
  display flex
  justify-content space-evenly
  margin-bottom 0.75rem
.paycheck-button
  margin-top 0.75rem
  max-width 12rem
  width 40%
  font-size 1.25rem
</style>
