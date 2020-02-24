<!--
GalaRegisterGuest handles the registration form contents for a single guest.
-->

<template lang="pug">
b-form-group(:label="`Guest #${number+1}`" label-class="font-weight-bold")
  b-form-group(label="Name" :label-for="`guest-${number}-name`" label-cols-sm="auto" label-class="gala-guest-label")
    b-form-input(:id="`guest-${number}-name`" trim v-model="name")
  b-form-group(label="Email" :label-for="`guest-${number}-email`" label-cols-sm="auto" label-class="gala-guest-label" :state="emailError ? false : null" :invalid-feedback="emailError")
    b-form-input(:id="`guest-${number}-email`" lazy trim v-model="email" :state="emailError ? false : null")
  b-form-group(label="Entree" :label-for="`guest-${number}-entree`" label-cols-sm="auto" label-class="gala-guest-label")
    b-form-select(:id="`guest-${number}-entree`" :options="entreeOptions" v-model="entree")
</template>

<script>
export default {
  props: {
    guest: Object,
    number: Number,
    entreeOptions: Array,
  },
  model: {
    prop: 'guest',
    event: 'input',
  },
  data: () => ({
    name: '',
    email: '',
    emailError: null,
    entree: '',
  }),
  watch: {
    guest() {
      this.name = this.guest.name
      this.email = this.guest.email
      this.entree = this.guest.entree
      this.validate()
    },
    name: 'emit',
    email() {
      this.validate()
      this.emit()
    },
    entree: 'emit',
  },
  methods: {
    emit() {
      this.$emit('input', { name: this.name, email: this.email, entree: this.entree, valid: !this.emailError })
    },
    validate() {
      if (this.email && !this.email.match(/^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$/))
        this.emailError = 'This is not a valid email address.'
      else
        this.emailError = null
    },
  }
}
</script>

<style lang="stylus">
.gala-guest-label
  width 5rem
</style>
