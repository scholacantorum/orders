<template lang="pug">
#app
  #app-header
    img#app-header-logo(src="./logo.png")
    #app-header-title Schola at the Door
  #app-body
    Login(v-if="!$store.state.auth")
    ChooseEvent(v-else-if="!$store.state.event")
    Main(v-else)
</template>

<script>
import ChooseEvent from './ChooseEvent'
import Login from './Login'
import Main from './Main.vue'

export default {
  components: { ChooseEvent, Login, Main },
  mounted() {
    const vh = window.innerHeight * 0.01;
    document.documentElement.style.setProperty('--vh', `${vh}px`);
  },
}
</script>

<style lang="stylus">
body
  box-sizing border-box
  font-family Arial, Arial, Helvetica, sans-serif
*
  box-sizing inherit
#app
  display flex
  flex-direction column
  overflow hidden
  width 100vw
  height calc(var(--vh, 1vh) * 100)
#app-header
  display flex
  flex none
  justify-content center
  align-items center
  height 3rem
  background-color #0153a5
#app-header-logo
  margin-right 0.5rem
  height 1.5rem
#app-header-title
  color white
  font-weight bold
  font-size 1.5rem
#app-body
  height calc(var(--vh, 1vh) * 100 - 3rem)
</style>
