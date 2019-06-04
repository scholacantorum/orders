<template lang="pug">
#logo-tall
  img(:src="logo")
  | Schola<br>Ticket<br>Scanner
</template>

<script>
import logo from './logo.png'

export default {
  computed: {
    logo() { return logo },
  },
}
</script>

<style lang="stylus">
#logo-tall
  display flex
  align-items center
  padding 4px 4px 4px 0
  background #0153A5
  color white
  font-size 14px
  line-height 1.2
  img
    padding-right 4px
</style>
