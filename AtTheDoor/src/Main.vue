<!--
Main displays the main UI, after preliminaries like logging in and choosing an
event are taken care of.
-->

<template>
  <View :style="{ flex: 1 }">
    <View :style="{ padding: 12, backgroundColor: '#017EFA' }">
      <Text :style="{ fontSize: 24, textAlign: 'center', color: '#fff' }">{{ eventHeader }}</Text>
    </View>
    <View
      v-if="$store.state.readerState"
      :style="{ padding: 12, backgroundColor: $store.state.readerState.color }"
    >
      <Text :style="{ fontSize: 20, textAlign: 'center' }">{{ $store.state.readerState.message }}</Text>
    </View>
    <SellTickets v-if="mode === 'sell'" :onDone="onMenu" />
    <WillCall v-else-if="mode === 'willcall'" :onDone="onMenu" />
    <ScanTicket v-else-if="mode === 'scan'" :onDone="onMenu" />
    <View v-else :style="{ alignItems: 'center' }">
      <Button
        v-if="canSellTickets"
        :bstyle="{ width: 160, marginTop: 48 }"
        title="Sell Tickets"
        :onPress="onSellTickets"
      />
      <Button
        v-if="canWillCall"
        :bstyle="{ width: 160, marginTop: 48 }"
        title="Will Call"
        :onPress="onWillCall"
      />
      <Button :bstyle="{ width: 160, marginTop: 48 }" title="Scan Ticket" :onPress="onScanTicket" />
      <Button
        :bstyle="{ width: 100, marginTop: 144, fontSize: 20 }"
        secondary
        title="Logout"
        :onPress="onLogout"
      />
      <Text :style="{ marginTop: 12, fontSize: 12, color: '#888' }">{{ stats }}</Text>
    </View>
  </View>
</template>

<script>
import backend from './backend'
import Button from './Button'
import ScanTicket from './ScanTicket'
import SellTickets from './SellTickets'
import WillCall from './WillCall'

export default {
  components: { Button, ScanTicket, SellTickets, WillCall },
  data: () => ({ mode: null }),
  computed: {
    canSellTickets() { return this.$store.state.allow.card || this.$storestate.allow.cash },
    canWillCall() { return this.$store.state.allow.willcall },
    eventHeader() { return `${this.$store.state.event.start.substr(0, 10)} ${this.$store.state.event.name}` },
    stats() {
      return `Admitted ${this.$store.state.admitted}, sold ${this.$store.state.sold}, cash $${this.$store.state.cash / 100}, check $${this.$store.state.check / 100}`
    },
  },
  methods: {
    onLogout() { backend.logout() },
    onMenu() { this.mode = null },
    onScanTicket() { this.mode = 'scan' },
    onSellTickets() { this.mode = 'sell' },
    onWillCall() { this.mode = 'willcall' },
  },
}
</script>
