<!--
Recordings opens a recordings purchase dialog on demand.
-->

<template lang="pug">
Dialog(
  ref="dialog"
  :auth="auth" :products="products" :stripeKey="stripeKey" :user="user"
  @close="onClose" @success="onSuccess"
)
</template>

<script>
import Dialog from './Dialog'

export default {
  props: {
    productIDs: Array,
    auth: String,
  },
  components: { Dialog },
  data: () => ({ products: null, stripeKey: null, user: null }),
  async mounted() {
    const params = new URLSearchParams()
    this.productIDs.forEach(p => { params.append("p", p) })
    params.append('auth', this.auth)
    let result
    try {
      result = (await this.$axios({ method: 'GET', url: '/api/prices', params })).data
    } catch (err) {
      console.error(err)
      window.parent.postMessage('close', '*')
      return
    }
    this.products = result.products
    this.stripeKey = result.stripePublicKey
    this.user = result.user
    this.$nextTick(() => { this.$refs.dialog.show() })
  },
  methods: {
    onClose() {
      window.parent.postMessage('close', '*')
    },
    onSuccess(orderID) {
      window.parent.postMessage(`success:${orderID}`, '*')
    }
  },
}
</script>

<style lang="stylus">
body
  background-color transparent !important
</style>
