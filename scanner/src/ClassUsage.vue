<!--
ClassUsage displays a ticket class and its current usage, and allows that usage
to be changed.  It emits 'change' with the new usage count when the usage is
changed.
-->

<template lang="pug">
.class-usage(:class="classColor")
  .class-usage-buttons
    div(v-text="tclass.name || 'General Admission'")
    b-btn.class-usage-button(v-for="n in max" :key="n" v-text="n"
      :variant="countVariant(n)" @click="onClick(n)"
    )
</template>

<script>
export default {
  props: {
    tclass: Object,
  },
  computed: {
    classColor() {
      if (this.tclass.overflow) return 'class-usage-overflow'
      if (this.tclass.name) return 'class-usage-restricted'
      return null
    },
    max() {
      if (this.tclass.max < 1000) return this.tclass.max
      return Math.ceil((this.tclass.used + 1) / 6) * 6
    },
  },
  methods: {
    countVariant(n) {
      if (n <= this.tclass.min) return 'secondary'
      return (n <= this.tclass.used) ? 'success' : `outline-success`
    },
    onClick(n) {
      if (n < this.tclass.min) return
      if (n === 1 && this.tclass.min === 0 && this.tclass.used === 1) this.$emit('change', 0)
      else this.$emit('change', n)
    },
  },
}
</script>

<style lang="stylus">
.class-usage
  padding 12px 6px 6px
.class-usage-overflow
  background-color red
.class-usage-restricted
  background-color gold
.class-usage-buttons
  margin 0 auto
  width 336px
.class-usage-button
  margin 0 6px 6px 0
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
