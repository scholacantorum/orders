<!--
CountChoice displays a horizontal, wrapping list of buttons with which to choose
a count of tickets sold.  It emits 'input' with the new count when the count is
changed (which means it can be used with v-model).
-->

<template lang="pug">
div
  b-button.count-choice(
    v-for="n in choices"
    :key="n"
    :disabled="n === 0 && !zero"
    :variant="(n === 0 && !zero) ? null : n === value ? 'primary' : 'outline-primary'"
    @click="$emit('input', n)"
    v-text="n"
  )
  b-button.count-choice(
    v-for="n in usedChoices"
    :key="n"
    :disabled="true"
    v-text="n"
  )
</template>

<script>
export default {
  props: {
    value: Number,
    zero: Boolean,
    max: Number,
    used: Number,
  },
  computed: {
    choices() {
      const list = []
      for (let n = 0; n <= this.max; n++) list.push(n)
      return list
    },
    usedChoices() {
      const list = []
      for (let n = this.max + 1; n <= this.max + this.used; n++) list.push(n)
      return list
    },
  },
}
</script>

<style lang="stylus">
.count-choice
  margin 0 6px 6px 0
  width 50px
  height 50px
</style>
