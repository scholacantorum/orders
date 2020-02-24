<!--
Dialog is the dialog box that opens when someone clicks "Buy Tickets".
-->

<template lang="pug">
b-form#gala-form(novalidate @submit.prevent="onSubmit")
  hr.w-100
  a#gala-register(name="register") Registration
  div.mb-3(v-if="success") Thank you for your registration.  A receipt has been emailed to you.  We look forward to seeing you at the gala.
  template(v-else)
    b-form-group(label="Register" label-for="gala-qty" label-cols="auto" label-class="mt-1" :state="qtyError ? false : null" :invalid-feedback="qtyError")
      b-form-input#gala-qty.d-inline(type="number" number min="1" v-model="qty")
      span(v-text="qtyLabel")
      b-form-text Register 10 seats to fill a table.
    GalaRegisterGuest(v-for="(guest, i) in guests" :key="i" :number="i" :entreeOptions="entreeOptions" v-model="guests[i]")
    b-form-group(label="Any special requests?" label-for="gala-requests" label-class="font-weight-bold")
      b-form-textarea#gala-requests(placeholder="Seating preferences, dietary restrictions, etc." v-model="requests")
    b-form-group.mb-0(label="Payment Information" label-class="font-weight-bold")
    OrderPayment(ref="pmt" :send="onSend" :stripeKey="stripeKey" :total="total")
</template>

<script>
import GalaRegisterGuest from './GalaRegisterGuest'
import OrderPayment from './OrderPayment'

export default {
  components: { GalaRegisterGuest, OrderPayment },
  props: {
    ordersURL: String,
    product: Object,
    stripeKey: String,
  },
  data: () => ({
    qty: 1,
    qtyError: null,
    qtyValid: true,
    guests: [{ name: '', email: '', entree: '', valid: true }],
    requests: '',
    success: false,
  }),
  computed: {
    entreeOptions() {
      const options = [{ text: '(please select)', value: '' }]
      this.product.options.forEach(o => {
        const [text, value] = o.split('/', 2)
        options.push({ text, value })
      })
      return options
    },
    guestEmails() {
      // I don't really care about the concatenation of guest emails, but this
      // value will change whenever a guest's email is changed, which means
      // whenever the validity of the guests might have changed, so I can watch
      // it for that.
      return this.guests.reduce((acc, g) => acc + g.email, '')
    },
    qtyLabel() {
      if (typeof (this.qty) === 'number' && this.qty === 1)
        return ` seat at $${this.product.price / 100}`
      return ` seats at $${this.product.price / 100} each`
    },
    total() {
      if (!this.valid) return null
      return this.qty * this.product.price
    },
    valid() {
      return !this.qtyError && !this.guests.find(g => !g.valid)
    },
  },
  watch: {
    qty() {
      this.qtyValid = typeof (this.qty) === 'number' && this.qty >= 1
      if (!this.qtyValid) {
        this.qtyError = 'Please enter a valid quantity.'
        return
      }
      this.qtyError = null
      if (this.qty < this.guests.length)
        this.guests.splice(this.qty)
      else {
        for (let i = this.guests.length; i < this.qty; i++) {
          this.guests.push({ name: '', email: '', entree: '', valid: true })
        }
      }
    },
  },
  methods: {
    async onSend({ name, email, address, city, state, zip, subtype, method }) {
      const body = new FormData()
      body.append('source', 'public')
      body.append('name', name)
      body.append('email', email)
      body.append('address', address)
      body.append('city', city)
      body.append('state', state)
      body.append('zip', zip)
      this.guests.forEach((guest, i) => {
        const prefix = `line${i + 1}.`
        body.append(prefix + 'product', this.product.id)
        body.append(prefix + 'quantity', 1)
        body.append(prefix + 'price', this.product.price)
        body.append(prefix + 'guestName', guest.name)
        body.append(prefix + 'guestEmail', guest.email)
        body.append(prefix + 'option', guest.entree)
      })
      body.append('payment1.type', 'card')
      body.append('payment1.subtype', subtype)
      body.append('payment1.method', method)
      body.append('payment1.amount', this.total)
      const result = await this.$axios.post(`${this.ordersURL}/payapi/order`, body,
        { headers: { 'Content-Type': 'application/x-www-form-urlencoded' } },
      ).catch(err => {
        return err
      })
      if (result && result.data && result.data.id) {
        this.success = true
        return null
      }
      if (result && result.data && result.data.error)
        return result.data.error
      console.error(result)
      return `We’re sorry, but we're unable to process payment cards at the moment.  Please try again later, or call our office at (650) 254-1700 to donate by phone.`
    },
    onSubmit() { this.$refs.pmt.submit() },
  },
}
</script>

<style lang="stylus">
#gala-register
  margin-bottom 12px
  text-align center
  font-weight bold
  font-size 18px
#gala-form
  display flex
  flex-direction column
  margin 6px 12px
  @media (min-width: 624px)
    margin 6px auto
    max-width 600px
#gala-qty
  flex none
  margin-right 8px
  width 70px // minimum to fit the placeholder text
  font-size 20px
  line-height 1
</style>
