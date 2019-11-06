<!--
BuyTickets displays a Buy Tickets button if one of the target products is
orderable, and an optional message otherwise.
-->

<template lang="pug">
div(v-if="message" v-text="message")
div(v-else-if="products")
  Dialog(ref="dialog"
    :coupon="coupon"
    :couponMatch="couponMatch"
    :ordersURL="ordersURL"
    :products="products"
    :stripeKey="stripeKey"
    :title="title"
    @coupon="onChangeCoupon"
  )
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
  data: () => ({ coupon: '', couponMatch: null, message: null, products: null }),
  mounted() {
    this.getPrices()
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
    async getPrices() {
      this.couponMatch = null // signal that we're retrieving prices
      const params = new (URLSearchParams)
      this.productIDs.forEach(p => { params.append("p", p) })
      params.append("coupon", this.coupon)
      let result
      try {
        result = (await this.$axios({
          method: 'GET',
          url: `${this.ordersURL}/api/prices`,
          params,
        })).data
      } catch (err) {
        console.error(err)
        this.couponMatch = this.coupon === ''
        return
      }
      this.message = (typeof result === 'string') ? result : null
      this.products = (typeof result === 'object') ? result.products : null
      this.couponMatch = (typeof result === 'object') ? result.coupon : (this.coupon === '')
      if (this.products)
        this.products.forEach(p => {
          if (p.name && p.name.indexOf("\n") > 0) {
            const parts = p.name.split("\n")
            p.name = parts[0]
            p.subname = parts[1]
          }
        })
    },
    onBuyTickets() {
      this.$refs.dialog.show()
    },
    onChangeCoupon(coupon) {
      this.coupon = coupon
      this.getPrices()
    },
  },
}
</script>

<style lang="stylus">
.buy-tickets-price
  margin-left 0.3em
</style>
