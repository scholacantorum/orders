<!--
CameraView displays the camera, or an error why the camera can't be displayed.
It emits 'decode' events whenever it sees a QR code.
-->

<template lang="pug">
div#camera
  #error(v-if="error" v-text="error")
  QrcodeStream(v-else :track="false" @init="onStreamInit" @decode="onDecode")
</template>

<script>
import { QrcodeStream } from 'vue-qrcode-reader'

export default {
  components: { QrcodeStream },
  data: () => ({ error: null }),
  computed: {
    streamStyle() {
      // On Android at least, the size of the stream is dictated solely by its
      // width and the standard 4:3 aspect ratio.  The height affects layout of
      // items further down the page, but it doesn't actually affect the height
      // of the displayed stream.  To make everything fit on the page, we need
      // to have the width be no more than 3/4 of the available vertical space,
      // and of course, no more than the available horizontal space.
      let width = (window.innerHeight - 100) * 3 / 4
      if (width > window.innerWidth) width = window.innerWidth
      let height = width * 4 / 3
      return { width: `${width}px`, height: `${height}px` }
    },
  },
  methods: {
    onDecode(text) { this.$emit('decode', text) },
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
  },
}
</script>

<style lang="stylus">
#camera
  flex none
  width 37.5vh
  height 50vh
#error
  display flex
  justify-content center
  align-items center
  padding 12px
  height 100%
  border 3px solid red
  color red
  text-align center
  font-weight bold
  font-size 24px
</style>
