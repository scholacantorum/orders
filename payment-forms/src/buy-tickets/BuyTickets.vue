<!--
BuyTickets displays a Buy Tickets button if one of the target products is
orderable, and an optional message otherwise.
-->

<template lang="pug">
div(v-if="message" v-text="message")
div(v-else-if="products")
  Dialog(ref="dialog" :products="products" :title="title")
  b-btn(variant="primary" @click="onBuyTickets") Buy Tickets
</template>

<script>
import Dialog from './Dialog'

export default {
  props: {
    productIDs: Array,
    title: String,
  },
  components: { Dialog },
  data: () => ({ coupon: '', couponMatch: true, message: null, products: null }),
  async mounted() {
    const params = new (URLSearchParams)
    this.productIDs.forEach(p => { params.append("p", p) })
    let result
    try {
      result = (await this.$axios({
        method: 'GET',
        url: `${process.env.VUE_APP_ORDERS_URL}/api/prices`,
        params,
      })).data
    } catch (err) {
      console.error(err)
      return
    }
    if (typeof result === 'string') this.message = result
    if (typeof result === 'object') {
      this.products = result.products
      this.couponMatch = result.coupon
    }
  },
  methods: {
    onBuyTickets() {
      this.$refs.dialog.show()
    },
  },
}
</script>

<style lang="stylus"></style>
