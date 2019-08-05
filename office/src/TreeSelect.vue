<!--
TreeSelect displays a tree of options, each with a count and a checkbox.
Ancestor nodes of the tree serve only to select and deselect their leaf
descendants; they have no independent selection state.

The tree is described by the "tree" property, which is a nested array of
objects.  Leaf objects have "id", "label", and "count" properties.  ("label" is
optional and defaults to "id".)  Leaf objects are displayed with the specified
label and associated count.  They are sorted by ID, and identified by ID in the
list of selected items, so the ID must be unique throughout the tree.  Non-leaf
objects have "aid" and "label" properties.  They are displayed with the
specified label and the count that is the sum of the counts of their
descendants.  The "aid" must be unique throughout the tree; it is used to
identify which ancestor nodes have been expanded.

The "value" prop is an array of the IDs of the leaf objects that are currently
selected.  The @change event is sent when the selection changes, with the new
list as its parameter.
-->

<template lang="pug">
.tree-select
  template(v-for="node in tree")
    TreeSelectAncestor(
      v-if="node.children"
      :key="node.aid"
      :checked="checked"
      :count="count"
      :expanded="expanded"
      :node="node"
      :path="[node.aid]"
      :value="value"
      @bulkChange="onBulkChange($event)"
      @change="onChange($event)"
      @expand="onExpand($event)"
    )
    TreeSelectLeaf(
      v-else
      :key="node.id"
      :flat="flat"
      :node="node"
      :path="[]"
      :value="value.includes(node.id)"
      @change="onChange($event)"
    )
</template>

<script>
import TreeSelectAncestor from './TreeSelectAncestor'
import TreeSelectLeaf from './TreeSelectLeaf'

export default {
  components: { TreeSelectAncestor, TreeSelectLeaf },
  props: {
    tree: Array,
    value: Array,
  },
  data: () => ({
    count: {},
    checked: {},
    expanded: {},
    flat: false,
  }),
  watch: {
    tree: { immediate: true, handler: 'setup' },
  },
  methods: {
    onBulkChange({ node, path, checked }) {
      let nv = [...this.value]
      node.children.forEach(child => { nv = this.onBulkChangeInner(child, path, checked, nv) })
      this.$emit('change', nv)
    },
    onBulkChangeInner(node, path, checked, value) {
      if (node.children) {
        path = [...path, node.aid]
        node.children.forEach(child => { value = this.onBulkChangeInner(child, path, checked, value) })
      } else {
        if (checked && !value.includes(node.id)) {
          value.push(node.id)
          path.forEach(p => { this.checked[p]++ })
        } else if (!checked && value.includes(node.id)) {
          value = value.filter(v => v !== node.id)
          path.forEach(p => { this.checked[p]-- })
        }
      }
      return value
    },
    onChange({ id, checked, path }) {
      if (checked) {
        this.$emit('change', [id, ...this.value])
        path.forEach(p => { this.checked[p]++ })
      } else {
        this.$emit('change', this.value.filter(v => v !== id))
        path.forEach(p => { this.checked[p]-- })
      }
    },
    onExpand(aid) {
      this.expanded[aid] = !this.expanded[aid]
    },
    setup() {
      this.count = {}
      this.checked = {}
      this.flat = true
      this.tree.forEach(root => { this.setupInner([], root) })
    },
    setupInner(path, root) {
      if (root.children) {
        this.flat = false
        path = [...path, root.aid]
        this.$set(this.count, root.aid, 0)
        this.$set(this.checked, root.aid, 0)
        this.$set(this.expanded, root.aid, !!this.expanded[root.aid])
        root.children.forEach(child => { this.setupInner(path, child) })
      } else {
        path.forEach(p => { this.count[p]++ })
        if (this.value.includes(root.id)) path.forEach(p => { this.checked[p]++ })
      }
    },
  },
}
</script>
