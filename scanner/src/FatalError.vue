<!--
FatalError displays a fatal error message.
-->

<template lang="pug">
#fatal-top
  LogoWide
  #fatal
    div(v-text="error")
    b-btn(@click="onReload" variant="primary") Reload
</template>

<script>
import LogoWide from './LogoWide'

export default {
  components: { LogoWide },
  props: {
    error: String,
  },
  methods: {
    onReload() { location.reload() },
  },
}
</script>

<style lang="stylus">
#fatal-top
  display flex
  flex-direction column
  height 100%
#fatal
  display flex
  flex auto
  flex-direction column
  justify-content space-around
  align-items center
  background-color red
  color white
  font-size 48px
  line-height 1.2
</style>
