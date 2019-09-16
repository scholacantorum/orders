<!--
ChooseEvent is the dialog for choosing what event we're serving.
-->

<template lang="pug">
#event
  #event-heading Choose Event
  b-spinner#event-spinner(v-if="!events")
  template(v-else)
    .event-tile(v-for="event in events" :key="event.id" @click="onClick(event)" v-text="`${event.start.substr(0, 10)} ${event.name}`")
</template>

<script>
export default {
  data: () => ({ events: null }),
  async mounted() {
    try {
      this.events = (await this.$axios.get('/api/event?future=1&freeEntries=1', {
        headers: { 'Auth': this.$store.state.auth },
      })).data
    } catch (err) {
      console.error('Error getting event list', err)
      window.alert(`Server error: ${err.toString()}`)
    }
  },
  methods: {
    async onClick(event) {
      try {
        const products = (await this.$axios.get('/api/prices', {
          headers: { 'Auth': this.$store.state.auth },
          params: { event: event.id, source: 'inperson' },
        })).data.products.filter(p => !p.message)
        if (!products.length) {
          window.alert('No tickets are on sale for that event.')
          return
        }
        this.$store.commit('event', { event, products })
      } catch (err) {
        console.error('Error getting event products', err)
        window.alert(`Server error: ${err.toString()}`)
      }
    },
  },
}
</script>

<style lang="stylus">
#event
  display flex
  flex-direction column
#event-heading
  padding 0.75rem
  background-color #017efa
  color #fff
  text-align center
  font-size 1.5rem
#event-spinner
  align-self center
  margin-top 1.5rem
.event-tile
  padding 0.75rem
  border-bottom 1px solid #ccc
  font-size 1.25rem
</style>
