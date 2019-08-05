<!--
TreeSelectAncestor displays an ancestor (i.e., non-leaf) node in a tree
selection.
-->

<template lang="pug">
div.tree-select-ancestor
  div.tree-select-item
    Icon.tree-select-expand(
      fixed-width
      :icon="isExpanded ? caretDown : caretRight"
      @click="onExpand"
    )
    input.tree-select-cb(
      ref="cb"
      :checked="isChecked"
      type="checkbox"
      @change="onChange"
    )
    div.tree-select-label(v-text="node.label")
    div.tree-select-count(v-text="node.count")
  .tree-select-ancestor-indent(v-if="isExpanded")
    template(v-for="child in node.children")
      TreeSelectAncestor(
        v-if="child.children"
        :key="child.label"
        :checked="checked"
        :count="count"
        :expanded="expanded"
        :node="child"
        :path="[...path, child.aid]"
        :value="value"
        @bulkChange="$emit('bulkChange', $event)"
        @change="$emit('change', $event)"
        @expand="$emit('expand', $event)"
      )
      TreeSelectLeaf(
        v-else
        :key="child.id"
        :flat="false"
        :node="child"
        :path="path"
        :value="value.includes(child.id)"
        @change="$emit('change', $event)"
      )
</template>

<script>
import { FontAwesomeIcon } from '@fortawesome/vue-fontawesome'
import { faCaretDown, faCaretRight } from '@fortawesome/free-solid-svg-icons'
import TreeSelectLeaf from './TreeSelectLeaf'

export default {
  name: 'TreeSelectAncestor',
  components: { Icon: FontAwesomeIcon, TreeSelectLeaf },
  props: {
    checked: Object,
    count: Object,
    expanded: Object,
    node: Object,
    path: Array,
    value: Array,
  },
  data: () => ({
    caretDown: faCaretDown,
    caretRight: faCaretRight,
  }),
  computed: {
    isChecked() { return this.checked[this.node.aid] === this.count[this.node.aid] },
    isExpanded() { return this.expanded[this.node.aid] },
    isIndeterminate() { return !this.isChecked && this.checked[this.node.aid] !== 0 },
  },
  watch: {
    isIndeterminate: {
      immediate: true,
      handler() {
        if (this.$refs.cb) this.$refs.cb.indeterminate = this.isIndeterminate
        else this.$nextTick(() => { this.$refs.cb.indeterminate = this.isIndeterminate })
      }
    }
  },
  methods: {
    onChange() { this.$emit('bulkChange', { node: this.node, path: this.path, checked: !this.isChecked }) },
    onExpand() { this.$emit('expand', this.node.aid) },
  },
}
</script>

<style lang="stylus">
.tree-select-item
  display flex
  align-items center
.tree-select-cb
  margin-right 0.3em
.tree-select-label
  flex auto
  overflow hidden
  min-width 0
  text-overflow ellipsis
  white-space nowrap
.tree-select-count
  color #888
  font-size smaller
  &::before
    margin-left 0.5em
    content '('
  &::after
    content ')'
.tree-select-ancestor-indent
  margin-left 1.2em
</style>
