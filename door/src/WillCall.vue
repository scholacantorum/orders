<!--
WillCall displays the list of customers who have tickets that could be, or were,
used at the event, and lets the user consume some of them.
-->

<template lang="pug">
#willcall-spinner(v-if="!orders")
  b-spinner
#willcall(v-else-if="!selected")
  b-form-input#willcall-search(v-model="search" :autoCapitalize="none" autocomplete="off" :autoCorrect="false" placeholder="Search Will Call List")
  #willcall-orders
    div
      .willcall-order(v-for="order in filteredOrders" :key="order.id" @click="onSelect(order.id)")
        | {{ order.name }}
        span.willcall-order-num(v-text="`#${order.id}`")
  b-button#willcall-cancel(@click="$emit('done')") Cancel
UseTickets(v-else :orderID="selected" @cancel="onUseCancel" @done="onUseDone")
</template>

<script>
import UseTickets from './UseTickets'

export default {
  components: { UseTickets },
  data: () => ({ orders: null, search: '', selected: null }),
  async mounted() {
    try {
      this.orders = (await this.$axios.get(`/api/event/${this.$store.state.event.id}/orders`, {
        headers: { 'Auth': this.$store.state.auth },
      })).data
    } catch (err) {
      if (err.response && err.response.status === 401) {
        this.$store.commit('logout')
        window.alert('Login session expired')
        return
      }
      console.error('Error getting will call list', err)
      window.alert(`Server error: ${err.toString()}`)
    }
  },
  computed: {
    filteredOrders() {
      if (!this.search.trim()) return this.orders
      return this.orders.filter(
        o => o.name.toLowerCase().includes(this.search.toLowerCase()) || parseInt(this.search) === o.id
      )
    },
  },
  methods: {
    onSelect(orderID) { this.selected = orderID },
    onUseCancel() { this.selected = null },
    onUseDone() { this.$emit('done') },
  },
}
</script>

<style lang="stylus">
#willcall-spinner
  margin-top 2rem
  text-align center
#willcall
  display flex
  flex-direction column
  height 100%
#willcall-search
  flex none
  margin 0.75rem
  width calc(100% - 1.5rem)
  font-size 1.5rem
#willcall-orders
  display flex
  flex auto
  flex-direction column
  overflow-y auto
  height calc(100% - 3.75rem - 3.375rem - 14px)
.willcall-order
  padding 0.75rem
  border-bottom 1px solid #ccc
  font-size 1.25rem
.willcall-order-num
  padding-left 1.5rem
  color #888
  font-size 1rem
#willcall-cancel
  flex none
  align-self center
  margin 0.75rem
  width 8rem
  font-size 1.25rem
</style>
