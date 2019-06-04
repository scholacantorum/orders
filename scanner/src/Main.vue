<template lang="pug">
#main
  FatalError(v-if="fatal" :error="fatal")
  LoginPage(v-else-if="!auth" @auth="auth=$event" @fatal="fatal=$event")
  EventChooser(v-else-if="!event" :auth="auth" @event="event=$event" @fatal="fatal=$event")
  TicketScanner(v-else :auth="auth" :event="event" :freeEntry="event.freeEntries")
</template>

<script>
import EventChooser from './EventChooser'
import FatalError from './FatalError'
import LoginPage from './LoginPage'
import TicketScanner from './TicketScanner'

export default {
  components: { EventChooser, FatalError, LoginPage, TicketScanner },
  data: () => ({
    auth: null,
    event: null,
    fatal: null,
  }),
}
</script>

<style lang="stylus">
#main
  width 100vw
  height 100vh
</style>
