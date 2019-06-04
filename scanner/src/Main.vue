<template lang="pug">
#main
  FatalError(v-if="fatal" :error="fatal")
  EventChooser(v-else-if="!event" @event="event = $event" @fatal="fatal = $event")
  TicketScanner(v-else :event="event" :freeEntry="event.freeEntries")
</template>

<script>
import EventChooser from './EventChooser'
import FatalError from './FatalError'
import TicketScanner from './TicketScanner'

export default {
  components: { EventChooser, FatalError, TicketScanner },
  data: () => ({
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
