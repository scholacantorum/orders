<!--
ScanTickets scans the barcode of a ticket and can mark it as used.
-->

<template lang="pug">
#scan(v-if="!selected")
  #scan-error(v-if="error" v-text="error")
  QrcodeStream#scan-stream(v-else :track="false" @init="onStreamInit" @decode="onDecode")
  form#scan-form(@submit.prevent="onSubmit")
    b-form-input#scan-oid(v-model="orderID" autocomplete="off" :formatter="oidFormatter" inputmode="numeric" pattern="[0-9]*" placeholder="Order #")
    b-button#scan-button(type="submit" :variant="buttonVariant" v-text="buttonLabel")
UseTickets(v-else :orderID="selected" @cancel="onUseDone" @done="onUseDone")
</template>

<script>
import { QrcodeStream } from 'vue-qrcode-reader'
import UseTickets from './UseTickets'

export default {
  components: { QrcodeStream, UseTickets },
  data: () => ({ error: null, orderID: null, selected: null }),
  computed: {
    buttonLabel() { return this.orderID ? 'Submit' : 'Cancel' },
    buttonVariant() { return this.orderID ? 'primary' : 'secondary' },
  },
  methods: {
    oidFormatter(text) { return text.replace(/[^0-9]/g, '') },
    onDecode(text) {
      if (!/\/ticket\/\d\d\d\d-\d\d\d\d-\d\d\d\d$/.test(text)) {
        window.alert('Not a Schola Cantorum order bar code')
        return
      }
      this.selected = text.substr(text.length - 14)
    },
    async onStreamInit(promise) {
      await promise.catch(err => {
        if (err.name === 'NotAllowedError') {
          this.error = 'no permission to use camera'
        } else if (err.name === 'NotFoundError') {
          this.error = 'no camera on this device'
        } else if (err.name === 'NotSupportedError') {
          this.error = 'secure context required (HTTPS, localhost)'
        } else if (err.name === 'NotReadableError') {
          this.error = 'camera already in use'
        } else if (err.name === 'OverconstrainedError') {
          this.error = 'camera is not suitable'
        } else if (err.name === 'StreamApiNotSupportedError') {
          this.error = "browser can't use camera"
        } else {
          this.error = "can't start QR code scanner"
        }
      })
    },
    onSubmit() {
      if (this.orderID) {
        this.selected = this.orderID
        this.orderID = null
      } else {
        this.$emit('done')
      }
    },
    onUseDone() { this.selected = null },
  },
}
</script>

<style lang="stylus">
#scan
  display flex
  flex auto
  flex-direction column
#scan-error
  padding 0.75rem
  background-color red
  color white
  text-align center
  font-size 1.25rem
#scan-stream
  flex auto
  max-height calc(var(--vh, 1vh) * 100 - 3rem - 3.75rem - 3.375rem - 14px)
  // app heading, main heading, cancel button
#scan-form
  z-index 1000
  display flex
  flex none
  justify-content space-around
  padding 0.75rem
  background-color white
#scan-oid
  width 8rem
  font-size 1.25rem
#scan-button
  width 8rem
  font-size 1.25rem
</style>
