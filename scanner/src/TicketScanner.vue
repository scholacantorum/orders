<!--
TicketScanner displays the page that scans a ticket bar code, or allows other
methods of entering ticket usage.  It emits a 'ticket' event whenever a ticket
is entered.  The value of the 'ticket' event may be an order ID number (manually
entered), an order token (scanned),  or 'free' for a free entry.
-->

<template lang="pug">
#scanner
  #scanner-top
    CameraView(@decode="onDecode")
    #scanner-controls
      LogoTall
      form#orderidform(@submit.prevent="onSubmit")
        #forminner
          b-form-input#orderid(ref="orderid" v-model="orderid" type="number" min="1" step="any" size="5" placeholder="Order #")
          b-button(type="submit") Submit
        b-button(v-if="freeEntry" @click.prevent="onFree" :disabled="!freeEntry") Free Entry
      #orderinfo(:class="orderInfoClass")
        #ordername(v-text="(order && order.name) ? order.name : '\u00A0'")
        #ordernum(v-text="(order && order.id) ? `Order number ${order.id}` : '\u00A0'")
  #scanner-bottom
    #scanError(v-if="scanError" v-text="scanError")
    #quantities(v-else-if="order && order.classes")
      ClassUsage(v-for="tclass in order.classes" :key="tclass.name"
        :tclass="tclass" @change="onCountChange(tclass, $event)"
      )
</template>

<script>
import CameraView from './CameraView'
import ClassUsage from './ClassUsage'
import LogoTall from './LogoTall'

export default {
  components: { CameraView, ClassUsage, LogoTall },
  props: {
    auth: String,
    event: Object,
    freeEntry: Array,
  },
  data: () => ({
    newOrder: false,
    order: null,
    orderid: null,
    scanError: null,
  }),
  computed: {
    orderInfoClass() { return this.newOrder ? 'orderinfo-new' : null },
  },
  methods: {
    async fetchTicket(token) {
      this.scanError = this.order = null
      const resp = await this.$axios.post(`/api/event/${this.event.id}/ticket/${token}`, null,
        { headers: { auth: this.auth } }
      ).catch(err => {
        if (err.response && err.response.status === 404) {
          this.scanError = 'No such order'
        } else {
          console.log(err)
          this.scanError = 'Server error'
        }
        return null
      })
      if (!resp) return
      if (resp.data.id) {
        this.order = resp.data
        this.newOrder = true
        window.setTimeout(() => { this.newOrder = false }, 100)
      }
      if (resp.data.error) this.scanError = resp.data.error
    },
    async onCountChange(tclass, count) {
      const body = new URLSearchParams()
      body.append('scan', this.order.scan)
      body.append('class', tclass.name)
      body.append('used', count)
      const resp = await this.$axios.post(
        `/api/event/${this.event.id}/ticket/${this.order.id || 'free'}`, body,
        { headers: { auth: this.auth } }
      ).catch(err => {
        console.log(err)
        this.scanError = 'Server error'
        return null
      })
      if (!resp) return
      if (resp.data.error) {
        this.scanError = resp.data.error
        return
      }
      this.$set(this.order, 'scan', resp.data.scan)
      this.$set(this.order, 'id', resp.data.id)
      tclass.used = count
      tclass.overflow = false
    },
    onDecode(text) {
      const m = text.match(/\/ticket\/(\d{4}-\d{4}-\d{4})$/)
      if (m) this.fetchTicket(m[1])
      else this.scanError = 'Not a Schola order'
    },
    onFree() {
      this.order = {
        scan: 'free',
        name: 'Free Entry',
        classes: this.freeEntry.map(fe => ({ name: fe, min: 0, max: 1000, used: 0 })),
      }
    },
    onSubmit() {
      if (this.orderid > 0) this.fetchTicket(this.orderid)
      this.orderid = null
      this.$refs.orderid.blur()
    },
  }
}
</script>

<style lang="stylus">
#scanner
  display flex
  flex-direction column
  width 100vw
  height 100vh
#scanner-top, #scanner-bottom
  flex none
  width 100vw
  height 50vh
#scanner-top
  display flex
#scanner-controls
  display flex
  flex none
  flex-direction column
  width calc(100vw - 37.5vh) // makes CameraView 4:3 aspect ratio
#orderidform
  display flex
  flex auto
  flex-direction column
  justify-content space-around
  box-sizing border-box
  padding 6px
#forminner
  display flex
  flex-direction column
#orderid
  margin-bottom 6px
#orderinfo
  background-color #fff
  text-align center
  transition background-color 0.5s
  &.orderinfo-new
    background-color #28a745
    transition none
#ordername
  font-weight bold
  font-size 20px
  line-height 24px
#orderid
  color #888
  font-size 16px
  line-height 16px
#scanError
  display flex
  justify-content center
  align-items center
  height 100%
  background-color red
  color white
  text-align center
  font-size 24px
  line-height 1.2
#quantities
  overflow-y auto
  height 100%
</style>
