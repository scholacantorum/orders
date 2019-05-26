<!--
TicketScanner displays the page that scans a ticket bar code, or allows other
methods of entering ticket usage.  It emits a 'ticket' event whenever a ticket
is entered.  The value of the 'ticket' event may be an order ID number (manually
entered), an order token (scanned), 'door' for an at-the-door sale, or 'free'
for a free entry.
-->

<template lang="pug">
#scanner
  QrcodeStream#stream(:style="streamStyle" :track="false" @init="onStreamInit" @decode="onDecode")
  form#orderidform(@submit.prevent="onSubmit")
    label(for="orderid" style="align-self:center;margin:0") Order number:
    b-form-input#orderid(v-model="orderid" type="number" min="1" step="any" size="5")
    b-button(type="submit") Submit
  #buttons
    b-button(@click="onDoor") Door Sale
    b-button(@click="onFree") Free Entry
</template>

<script>
import { QrcodeStream } from 'vue-qrcode-reader'

export default {
  components: { QrcodeStream },
  data: () => ({ orderid: null }),
  computed: {
    streamStyle() {
      // On Android at least, the size of the stream is dictated solely by its
      // width and the standard 4:3 aspect ratio.  The height affects layout of
      // items further down the page, but it doesn't actually affed the height
      // of the displayed stream.  To make everything fit on the page, we need
      // to have the width be no more than 3/4 of the available vertical space,
      // and of course, no more than the available horizontal space.
      let width = (window.innerHeight - 150) * 3 / 4
      if (width > window.innerWidth) width = window.innerWidth
      let height = width * 4 / 3
      return { width: `${width}px`, height: `${height}px` }
    },
  },
  methods: {
    onDecode(text) {
      if (text.startsWith('https://orders.scholacantorum.org/ticket/'))
        this.$emit('ticket', text.substr(42))
      else
        this.$emit('ticket', 'non-schola')
    },
    onDoor() {
      this.$emit('ticket', 'door')
    },
    onFree() {
      this.$emit('ticket', 'free')
    },
    async onStreamInit(promise) {
      await promise.catch(err => {
        if (err.name === 'NotAllowedError') {
          window.alert('ERROR: you need to grant camera access permisson')
        } else if (err.name === 'NotFoundError') {
          window.alert('ERROR: no camera on this device')
        } else if (err.name === 'NotSupportedError') {
          window.alert('ERROR: secure context required (HTTPS, localhost)')
        } else if (err.name === 'NotReadableError') {
          window.alert('ERROR: is the camera already in use?')
        } else if (err.name === 'OverconstrainedError') {
          window.alert('ERROR: installed cameras are not suitable')
        } else if (err.name === 'StreamApiNotSupportedError') {
          window.alert('ERROR: Stream API is not supported in this browser')
        } else {
          window.alert('ERROR: unable to start QR code scanner')
        }
      })
    },
    onSubmit() {
      if (this.orderid > 0) this.$emit('ticket', this.orderid)
    },
  }
}
</script>

<style lang="stylus">
#scanner
  width 100%
  height 100%
#orderidform
  display flex
  justify-content center
  box-sizing border-box
  padding 6px 12px
  width 100%
  height 50px
  border-bottom 1px solid #ccc
#orderid
  margin 0 24px
  width 5em
#buttons
  display flex
  justify-content space-around
  align-items center
  width 100%
  height 50px
  button
    width 35%
    height 38px
</style>
