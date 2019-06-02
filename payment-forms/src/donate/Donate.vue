<!--
Donate displays a Donate Online button, replaced with a thank-you when the
donation is complete.
-->

<template lang="pug">
#donate-div
  Dialog(ref="dialog"
    :ordersURL="ordersURL" :stripeKey="stripeKey"
  )
  b-btn#donate-btn(variant="primary" @click="$refs.dialog.show()") Donate Online
  span with a credit or debit&nbsp;card
</template>

<script>
import Dialog from './Dialog'

export default {
  props: {
    ordersURL: String,
    stripeKey: String,
  },
  components: { Dialog },
}
</script>

<style lang="stylus">
#donate-div
  display flex
  align-items center
  margin-bottom 16px
  line-height 1.2
#donate-btn
  flex none
  margin-right 1ex
</style>
