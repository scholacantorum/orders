<!--
BuyTickets displays a Buy Tickets button if one of the target products is
orderable, and an optional message otherwise.
-->

<template lang="pug">
div(v-if="message" v-text="message")
div(v-else-if="products")
  Dialog(ref="dialog" :ordersURL="ordersURL" :products="products" :stripeKey="stripeKey" :title="title")
  b-btn(variant="primary" @click="onBuyTickets" v-text="buttonLabel")
  span.buy-tickets-price(v-if="priceLabel" v-text="priceLabel")
</template>

<script>
import Dialog from './Dialog'

export default {
  props: {
    ordersURL: String,
    productIDs: Array,
    stripeKey: String,
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
        url: `${this.ordersURL}/api/prices`,
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
  computed: {
    buttonLabel() {
      return this.title.includes('Subscription') ? 'Buy Subscriptions' : 'Buy Tickets'
    },
    priceLabel() {
      if (!this.products || !this.products.length || !this.products[0].price) return null
      if (this.products.some(p => p.price !== this.products[0].price)) return null
      return `$${this.products[0].price / 100} each`
    },
  },
  methods: {
    onBuyTickets() {
      this.$refs.dialog.show()
    },
  },
}
</script>

<style lang="stylus">
.buy-tickets-price
  margin-left 0.3em
</style>
