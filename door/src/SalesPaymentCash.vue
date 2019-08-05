<!--
SalesPaymentCash executes the cash payment flow for the order given to it.
-->

<template lang="pug">
#paycash
  Summary(:order="order")
  .paycash-row
    .paycash-label Amount Due
    #paycash-due(v-text="`$${due}`")
  .paycash-row
    .paycash-label Cash Received
    input#paycash-received(v-model="received" autocomplete="off" :formatter="numFormatter" inputmode="numeric" pattern="[0-9]*")
    b-button.paycash-round(v-for="v in rounded" :key="v" variant="outline-primary" @click="setReceived(v)" v-text="`$${v}`")
  .paycash-row
    .paycash-label Change
    #paycash-change(v-text="`$${change}`")
    b-form-checkbox#paycash-donate.ml-3(v-model="donate" :disabled="!change" size="lg" switch) Donation
  #paycash-buttons
    b-button.paycash-button(:disabled="!canConfirm" variant="primary" @click="onConfirmed" v-text="confirmed ? 'Saving...' : 'Confirm'")
    b-button.paycash-button(:disabled="confirmed" @click="$emit('cancel')") Cancel
</template>

<script>
import Summary from './SalesSummary'

export default {
  components: { Summary },
  props: {
    order: Object,
  },
  data: () => ({ confirmed: false, donate: false, received: 0 }),
  mounted() {
    this.received = this.order.payments[0].amount / 100
  },
  computed: {
    canConfirm() { return (this.received >= this.due) && !this.confirmed },
    change() {
      if (this.received < this.due) return 0
      return this.received - this.due
    },
    due() { return this.order.payments[0].amount / 100 },
    rounded() {
      const a = this.order.payments[0].amount / 100
      const list = []
      const a5 = Math.ceil(a / 5) * 5
      if (a5 !== a) list.push(a5)
      const a20 = Math.ceil(a / 20) * 20
      if (a20 !== a5) list.push(a20)
      return list
    },
  },
  methods: {
    numFormatter(text) { return text.replace(/[^0-9]/g, '') },
    async onConfirmed() {
      this.confirmed = true
      try {
        let order = { ...this.order }
        if (this.change && this.donate) {
          order.lines = [...order.lines]
          order.lines.push({
            product: 'donation',
            quantity: 1,
            price: this.change * 100,
          })
          order.payments = [{ ...order.payments[0] }]
          order.payments[0].amount += this.change * 100
        }
        const revised = (await this.$axios.post('/api/order', order, {
          headers: { 'Auth': this.$store.state.auth },
        })).data
        if (revised.error) throw revised.error
        this.$store.commit('sold', {
          count: this.order.lines.reduce((accum, line) => line.quantity + accum, 0),
          amount: this.order.payments[0].amount,
          method: this.order.payments[0].subtype,
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
#paycash
  display flex
  flex-direction column
.paycash-row
  display flex
  align-self center
  align-items center
  margin 0.75rem 0.75rem 0
  width 24rem
.paycash-label
  width 10rem
  font-size 1.25rem
#paycash-due, #paycash-change, #paycash-received
  width 5rem
  text-align right
  font-size 1.25rem
#paycash-received
  margin-left 5px
  padding 5px
  border 1px solid #ccc
.paycash-round
  margin-left 0.5rem
  width 4rem
  font-size 1.25rem
#paycash-buttons
  display flex
  justify-content space-evenly
  margin-bottom 0.75rem
.paycash-button
  max-width 12rem
  width 40%
  font-size 1.25rem
</style>
