<!--
TreeSelectLeaf displays a leaf node in a tree selection.
-->

<template lang="pug">
div.tree-select-item
  .tree-select-leaf-expand(v-if="!flat")
  input.tree-select-cb(type="checkbox" :checked="value" @change="onChange")
  div.tree-select-label(v-text="node.label || node.id")
  div.tree-select-count(v-text="node.count")
</template>

<script>
export default {
  props: {
    flat: Boolean,
    node: Object,
    path: Array,
    value: Boolean,
  },
  methods: {
    onChange() {
      this.$emit('change', { id: this.node.id, path: this.path, checked: !this.value })
    },
  }
}
</script>

<style lang="stylus">
.tree-select-leaf-expand
  flex none
  width 1.25em
</style>
