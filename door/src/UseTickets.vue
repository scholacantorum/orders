<!--
UseTickets displays the used and available tickets on an order and allows them
to be consumed.
-->

<template lang="pug">
#use-spinner(v-if="!order")
  b-spinner
#use(v-else)
  #use-order
    #use-order-name(v-text="order.name")
    #use-order-num(v-text="`Order #${order.id}`")
  #use-classes
    UseTicketClass(v-for="tclass in order.classes" :key="tclass.name" :tclass="tclass" @change="onClassChange")
  #use-buttons
    b-button.use-button(:disabled="disabled" variant="primary" @click="onUseTickets" v-text="buttonText")
    b-button.use-button(@click="$emit('cancel')") Cancel
</template>

<script>
import UseTicketClass from './UseTicketClass'

export default {
  components: { UseTicketClass },
  props: {
    orderID: String,
  },
  data: () => ({ order: null, using: false }),
  async mounted() {
    try {
      const order = (await this.$axios.get(`/posapi/event/${this.$store.state.event.id}/ticket/${this.orderID}`, {
        headers: { 'Auth': this.$store.state.auth },
      })).data
      if (order.error) {
        window.alert(order.error)
        this.$emit('cancel')
        return
      }
      this.order = order
    } catch (err) {
      if (err.response && err.response.status === 401) {
        this.$store.commit('logout')
        window.alert('Login session expired')
        return
      }
      console.error('Error fetching ticket usage', err)
      window.alert(`Server error: ${err.toString()}`)
    }
  },
  computed: {
    buttonText() { return this.using ? 'Using...' : 'Use Tickets' },
    disabled() { return this.using || this.order.classes.every(cl => cl.min === cl.used) }
  },
  methods: {
    onClassChange({ tclass, used }) {
      tclass.used = used
      tclass.overflow = false
    },
    async onUseTickets() {
      try {
        this.using = true
        const body = new URLSearchParams()
        this.order.classes.forEach(cl => {
          body.append('class', cl.name)
          body.append('used', cl.used)
        })
        await this.$axios.post(`/posapi/event/${this.$store.state.event.id}/ticket/${this.order.id}`, body, {
          headers: {
            'Auth': this.$store.state.auth,
          },
        })
        this.$store.commit('admitted', this.order.classes.reduce((accum, cl) => cl.used - cl.min, 0))
        this.$emit('done')
      } catch (err) {
        if (err.response && err.response.status === 401) {
          this.$store.commit('logout')
          window.alert('Login session expired')
          return
        }
        this.using = false
        console.error('Error updating ticket usage', err)
        window.alert(`Server error: ${err.toString()}`)
      }
    },
  },
}
</script>

<style lang="stylus">
#use-spinner
  margin-top 2rem
  text-align center
#use
  display flex
  flex auto
  flex-direction column
#use-order
  display flex
  flex none
  flex-direction column
  align-items center
  padding 0.75rem
  background-color #ccc
#use-order-name
  font-weight bold
  font-size 1.5rem
#use-order-num
  font-size 1.25rem
#use-classes
  display flex
  flex auto
  flex-direction column
  overflow-y auto
#use-buttons
  display flex
  flex none
  justify-content space-around
  padding 0.75rem
.use-button
  width 10rem
  font-size 1.25rem
</style>
