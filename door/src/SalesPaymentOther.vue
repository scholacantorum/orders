<!--
SalesPaymentOther executes the cash/check payment flow for the order given to
it.
-->

<template lang="pug">
#payother
  Summary(:order="order")
  #payother-message(v-text="confirmMessage")
  #payother-buttons
    b-button.payother-button(:disabled="confirmed" variant="primary" @click="onConfirmed" v-text="confirmed ? 'Saving...' : 'Confirmed'")
    b-button.payother-button(:disabled="confirmed" @click="$emit('cancel')") Cancel
</template>

<script>
import Summary from './SalesSummary'

export default {
  components: { Summary },
  props: {
    order: Object,
  },
  data: () => ({ confirmed: false }),
  computed: {
    confirmMessage() {
      if (this.order.payments[0].method === 'Cash')
        return `Confirm you have received $${this.order.payments[0].amount / 100}\u00A0in\u00A0cash.`
      else
        return `Confirm you have received a check for\u00A0$${this.order.payments[0].amount / 100}.`
    },
  },
  methods: {
    async onConfirmed() {
      this.confirmed = true
      try {
        const revised = (await this.$axios.post('/api/order', this.order, {
          headers: { 'Auth': this.$store.state.auth },
        })).data
        if (revised.error) throw revised.error
        this.$store.commit('sold', {
          count: revised.lines.reduce((accum, line) => line.quantity + accum, 0),
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
  },
}
</script>

<style lang="stylus">
#payother
  display flex
  flex-direction column
  justify-content space-between
  height 100%
#payother-message
  margin 0.75rem
  color #f90
  text-align center
  font-weight bold
  font-size 1.5rem
  line-height 1.25
#payother-buttons
  display flex
  justify-content space-evenly
  margin-bottom 0.75rem
.payother-button
  max-width 12rem
  width 40%
  font-size 1.25rem
</style>
