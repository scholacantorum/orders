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
    <SellTickets v-if="mode === 'sell'" :onDone="onMenu"/>
    <WillCall v-else-if="mode === 'willcall'" :onDone="onMenu"/>
    <ScanTicket v-else-if="mode === 'scan'" :onDone="onMenu"/>
    <View v-else :style="{ alignItems: 'center' }">
      <Button
        :bstyle="{ width: 160, marginTop: 48 }"
        title="Sell Tickets"
        :onPress="onSellTickets"
      />
      <Button :bstyle="{ width: 160, marginTop: 48 }" title="Will Call" :onPress="onWillCall"/>
      <Button :bstyle="{ width: 160, marginTop: 48 }" title="Scan Ticket" :onPress="onScanTicket"/>
      <Button
        :bstyle="{ width: 100, marginTop: 144, fontSize: 20 }"
        secondary
        title="Logout"
        :onPress="onLogout"
      />
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
    eventHeader() { return `${this.$store.state.event.start.substr(0, 10)} ${this.$store.state.event.name}` },
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
