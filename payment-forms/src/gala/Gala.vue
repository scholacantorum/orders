<!--
Gala displays a gala registration form if registrations are open, replaced with
a thank-you when registration is complete.
-->

<template lang="pug">
GalaRegistration(v-if="product" :product="product" :ordersURL="ordersURL" :stripeKey="stripeKey")
#gala-message(v-else-if="productError" v-text="productError")
</template>

<script>
import GalaRegistration from './GalaRegistration'

export default {
  components: { GalaRegistration },
  props: {
    ordersURL: String,
    stripeKey: String,
    productID: String,
  },
  data: () => ({ product: null, productError: null }),
  async created() {
    const response = (await this.$axios.get(`${this.ordersURL}/payapi/prices?p=${this.productID}`)).data
    if (!response) return
    if (typeof response === 'string') this.productError = response
    else this.product = response.products[0]
  },
}
</script>

<style lang="stylus"></style>
