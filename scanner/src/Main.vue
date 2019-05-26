<template lang="pug">
  #top
    #logo
      | Schola Cantorum Ticket Scanner
    #main
      EventChooser(v-if="!event")
      TicketScanner(v-else-if="!ticket" @ticket="onTicket")
      ConfirmScan(v-else :event="event" :ticket="ticket" @done="onDone")
</template>

<script>
import ConfirmScan from './ConfirmScan'
import TicketScanner from './TicketScanner'

const EventChooser = {}

export default {
  components: { ConfirmScan, EventChooser, TicketScanner },
  data: () => ({
    event: '2019-07-15',
    ticket: null,
  }),
  methods: {
    onDone() {
      this.ticket = null
    },
    onTicket(ticket) {
      this.ticket = ticket
    }
  },
}
</script>

<style lang="stylus">
#top
  width 100vw
  height 100vh
#logo
  display flex
  align-items center
  padding-left 55px
  width 100vw
  height 50px
  background #0157A5 url('./logo.png') no-repeat
  background-position left center
  color white
#main
  width 100vw
  height calc(100vh - 50px)
</style>
