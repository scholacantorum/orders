<!--
UseTicketClass displays and changes the ticket usage for a single ticket class.
-->

<template lang="pug">
.useclass(:class="classes")
  .useclass-inner
    .useclass-name(v-text="tclass.name || 'General Admission'")
    .useclass-buttons
      b-button.useclass-button(v-for="n in max" :key="n" :disabled="disabled(n)" :variant="variant(n)" @click="onClick(n)" v-text="n")
</template>

<script>
export default {
  props: {
    tclass: Object,
  },
  computed: {
    classes() {
      return this.tclass.overflow ? 'useclass-overflow' : this.tclass.name ? 'useclass-restricted' : null
    },
    max() {
      if (this.tclass.max < 1000) return this.tclass.max
      return Math.ceil((this.tclass.used + 1) / 6) * 6
    },
  },
  methods: {
    disabled(n) { return n < this.tclass.min },
    onClick(n) {
      if (n < this.tclass.min) return
      if (n === 1 && this.tclass.min === 0 && this.tclass.used === 1)
        this.$emit('change', { tclass: this.tclass, used: 0 })
      else
        this.$emit('change', { tclass: this.tclass, used: n })
    },
    variant(n) {
      if (n <= this.tclass.min) return 'dark'
      if (n <= this.tclass.used) return 'success'
      return 'outline-success'
    },
  }
}
</script>

<style lang="stylus">
.useclass
  display flex
  flex-direction column
  align-items center
  padding 0.75rem 0.375rem 0.375rem
.useclass-overflow
  background-color red
.useclass-restricted
  background-color gold
.useclass-inner
  width 336px
.useclass-name
  font-size 1.25rem
.useclass-buttons
  display flex
  flex-wrap wrap
.useclass-button
  margin 0 6px 6px 0
  padding 12px 6px
  width 50px
  height 50px
  border-width 2px
  box-shadow none !important
  font-size 20px
  &.btn-outline-success
    // override filled-in style when active or hovered
    background-color transparent !important
    color #28a745 !important
</style>
