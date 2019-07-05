<!--
Main displays the main UI, after preliminaries like logging in and choosing an
event are taken care of.
-->

<template lang="pug">
#main
  #main-header(v-text="eventHeader")
  #main-body
    SellTickets(v-if="mode === 'sell'" @done="onMenu")
    WillCall(v-else-if="mode === 'willcall'" @done="onMenu")
    ScanTickets(v-else-if="mode === 'scan'" @done="onMenu")
    #main-menu(v-else)
      b-button.main-button(v-if="canSellTickets" variant="primary" @click="onSellTickets") Sell Tickets
      b-button.main-button(v-if="canWillCall" variant="primary" @click="onWillCall") Will Call
      b-button.main-button(variant="primary" @click="onScanTickets") Scan Tickets
      b-button#main-logout(@click="onLogout") Logout
      #main-stats(v-text="stats")
</template>

<script>
import ScanTickets from './ScanTickets'
import SellTickets from './SellTickets'
import WillCall from './WillCall'

export default {
  components: { ScanTickets, SellTickets, WillCall },
  data: () => ({ mode: null }),
  computed: {
    canSellTickets() { return this.$store.state.allow.card || this.$storestate.allow.cash },
    canWillCall() { return this.$store.state.allow.willcall },
    eventHeader() {
      return `${this.$store.state.event.start.substr(0, 10)} ${this.$store.state.event.name}`
    },
    stats() {
      return `Admitted ${this.$store.state.admitted}, sold ${this.$store.state.sold}, cash $${this.$store.state.cash / 100}, check $${this.$store.state.check / 100}`
    },
  },
  methods: {
    onLogout() { this.$store.commit('logout') },
    onMenu() { this.mode = null },
    onScanTickets() { this.mode = 'scan' },
    onSellTickets() { this.mode = 'sell' },
    onWillCall() { this.mode = 'willcall' },
  },
}
</script>

<style lang="stylus">
#main
  display flex
  flex-direction column
  height 100%
#main-header
  flex none
  padding 0.75rem
  background-color #017efa
  color #fff
  text-align center
  font-size 1.5rem
#main-body
  height calc(100% - 3.75rem)
#main-menu
  display flex
  flex-direction column
  align-items center
.main-button
  margin-top 3rem
  width 12rem
  font-size 1.5rem
#main-logout
  margin-top 6rem
  width 8rem
  font-size 1.25rem
#main-stats
  margin-top 0.75rem
  color #888
  font-size 0.75rem
</style>
