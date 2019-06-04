<!--
EventChooser displays a list of upcoming events and allows the user to choose
which one they are taking tickets for.
-->

<template lang="pug">
#events
  LogoWide
  #events-head(v-if="events.length") Which event are you taking tickets for?
  .event(v-for="event in events" :key="event.id" @click="$emit('event', event)")
    | {{ event.start.substr(0, 10) }} {{ event.name }}
</template>

<script>
import LogoWide from './LogoWide'

export default {
  components: { LogoWide },
  props: {
    auth: String,
  },
  data: () => ({ events: [] }),
  async mounted() {
    const resp = await this.$axios.get(`/api/event?future=1&freeEntries=1`, {
      headers: { auth: this.auth },
    }).catch(err => {
      console.log(err)
      this.$emit('fatal', 'Server error')
      return null
    })
    if (!resp) return
    if (resp.status !== 200) {
      console.log(resp.statusText)
      this.$emit('fatal', 'Server error')
      return
    }
    this.events = resp.data
  },
}
</script>

<style lang="stylus">
#events
  overflow auto
  width 100%
  height 100%
#events-head
  margin 12px
  font-weight bold
  font-size 20px
.event
  display flex
  align-items center
  padding 12px
  width 100%
  height 50px
  border-bottom 1px solid #ccc
</style>
