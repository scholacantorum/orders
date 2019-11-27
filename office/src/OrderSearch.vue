<!--
OrderSearch displays the order search page.
-->

<template lang="pug">
#osearch
  #osearch-cols
    .osearch-col
      .osearch-criteria-section Order Number
      form#osearch-ordernum-row
        input#osearch-criteria-ordernum(:value="ordernum || ''" @input="onOrderNumber")
        b-button#osearch-viewedit(type="submit" variant="primary" :disabled="!ordernum" v-text="viewEditLabel")
      #osearch-spinner(v-if="!stats")
        b-spinner(label="Loading...")
      template(v-else)
        .osearch-criteria-section(style="margin-top:0.8em") Customer
        input#osearch-criteria-customer(v-model.lazy.trim="customer" type="text" placeholder="Name or Email")
        .osearch-criteria-section(style="margin-top:0.8em") Order Dates
        .osearch-criteria-date-box
          | From
          input.osearch-criteria-date(v-model="createdAfter" type="date" placeholder="From Date")
        .osearch-criteria-date-box
          | To
          input.osearch-criteria-date(v-model="createdBefore" type="date" placeholder="To Date")
    .osearch-col(v-if="stats")
      .osearch-criteria-section Products
      TreeSelect(:tree="productsTree" :value="selectedProducts" @change="onChangeProducts")
    .osearch-col(v-if="stats")
      .osearch-criteria-section Tickets Used At
      TreeSelect(:tree="usedAtEventsTree" :value="selectedUsedAtEvents" @change="onChangeUsedAtEvents")
    .osearch-col(v-if="stats")
      .osearch-criteria-section Ticket Classes
      TreeSelect(:tree="ticketClassList" :value="selectedTicketClasses" @change="onChangeTicketClasses")
    .osearch-col(v-if="stats")
      .osearch-criteria-section Order Sources
      TreeSelect(:tree="orderSourcesList" :value="selectedOrderSources" @change="onChangeOrderSources")
    .osearch-col(v-if="stats")
      .osearch-criteria-section Payment Types
      TreeSelect(:tree="paymentTypesTree" :value="selectedPaymentTypes" @change="onChangePaymentTypes")
    .osearch-col(v-if="stats")
      .osearch-criteria-section Coupon Codes
      TreeSelect(:tree="orderCouponsList" :value="selectedOrderCoupons" @change="onChangeOrderCoupons")
  #osearch-footer
    #osearch-stats(v-text="statsString")
    #osearch-buttons
      b-button#osearch-reset(v-if="stats" variant="outline-primary" @click="onReset") Reset
      b-button#osearch-search-list(
        v-if="stats" variant="primary" :disabled="searchListDisabled" @click="onSearchList" v-text="searchListLabel"
      )
</template>

<script>
import TreeSelect from './TreeSelect'

const productTypeLabels = {
  auctionitem: 'Auction Items',
  donation: 'Donations',
  recording: 'Recordings',
  sheetmusic: 'Sheet Music',
  ticket: 'Tickets',
  wardrobe: 'Wardrobe',
  other: '[other]',
}

const sourceLabels = {
  gala: 'Gala Software',
  inperson: 'Event Box Office',
  members: 'Members Web Site',
  office: 'Schola Office',
  public: 'Public Web Site',
}

export default {
  components: { TreeSelect },
  data: () => ({
    stats: null,
    createdAfter: '',
    createdBefore: '',
    customer: '',
    ordernum: 0,
    selectedOrderCoupons: [],
    selectedOrderSources: [],
    selectedPaymentTypes: [],
    selectedProducts: [],
    selectedTicketClasses: [],
    selectedUsedAtEvents: [],
    modifiedCriteria: false,
  }),
  mounted() {
    this.getStats(null)
  },
  watch: {
    customer() { this.modifiedCriteria = true },
    createdAfter() { this.modifiedCriteria = true },
    createdBefore() { this.modifiedCriteria = true },
  },
  computed: {
    orderCouponsList() {
      return this.stats.orderCoupons.map(oc => ({
        id: oc.n, label: oc.n || '(none)', count: oc.c,
      })).sort((a, b) => a.label.localeCompare(b.label))
    },
    orderSourcesList() {
      return this.stats.orderSources.map(os => ({
        id: os.os, label: sourceLabels[os.os] || os.os, count: os.c,
      })).sort((a, b) => a.label.localeCompare(b.label))
    },
    paymentTypesTree() {
      let tree = {}
      this.stats.paymentTypes.forEach(pt => {
        const path = pt.n.split(',')
        const leaf = path.pop()
        let root = tree
        let aid = ''
        path.forEach(p => {
          aid += ',' + p
          if (!root[p]) root[p] = { label: p, aid, count: 0, children: {} }
          root[p].count += pt.c
          root = root[p].children
        })
        root[pt.n] = { id: pt.n, label: leaf, count: pt.c }
      })
      tree = this.sortTree(tree)
      return tree
    },
    productsTree() {
      let tree = {}
      this.stats.products.forEach(p => {
        if (!tree[p.ptype]) tree[p.ptype] = { label: productTypeLabels[p.ptype] || p.ptype, aid: p.ptype, count: 0, children: {} }
        if (p.series) {
          if (!tree[p.ptype].children[p.series])
            tree[p.ptype].children[p.series] = { label: p.series, aid: `${p.ptype} ${p.series}`, count: 0, children: {} }
          tree[p.ptype].children[p.series].children[p.id] = { id: p.id, label: p.name, count: p.count }
          tree[p.ptype].children[p.series].count += p.count
        } else {
          tree[p.ptype].children[p.id] = { id: p.id, label: p.name, count: p.count }
        }
        tree[p.ptype].count += p.count
      })
      tree = this.sortTree(tree)
      return tree
    },
    searchListDisabled() {
      return !this.modifiedCriteria && (!this.stats.lines || !this.stats.lines.length)
    },
    searchListLabel() {
      return this.modifiedCriteria ? 'Search' : ' Show List'
    },
    statsString() {
      return `Matched ${this.stats.orderCount} ${this.stats.orderCount === 1 ? 'order' : 'orders'}, ${this.stats.itemCount} ${this.stats.itemCount === 1 ? 'item' : 'items'}, total amount $${this.stats.totalAmount.toFixed(2)}.`
    },
    ticketClassList() {
      return this.stats.ticketClasses.map(tc => ({
        id: tc.n, label: tc.n || 'General Admission', count: tc.c,
      }))
    },
    usedAtEventsTree() {
      let tree = {}
      this.stats.usedAtEvents.forEach(e => {
        if (!e.series) {
          tree[e.id] = { id: e.id, label: `${e.start ? e.start.substr(0, 10) + ' ' : ''}${e.name}`, count: e.count }
        } else {
          if (!tree[e.series]) tree[e.series] = { label: e.series, aid: e.series, count: 0, children: {} }
          tree[e.series].children[e.id] = { id: e.id, label: `${e.start.substr(0, 10)} ${e.name}`, count: e.count }
          tree[e.series].count += e.count
        }
      })
      tree = this.sortTree(tree)
      return tree
    },
    viewEditLabel() {
      return this.$store.state.privManageOrders ? 'Edit' : 'View'
    },
  },
  methods: {
    async getStats(params) {
      this.$emit('list', null)
      try {
        this.stats = (await this.$axios.get('/api/report', {
          headers: { Auth: this.$store.state.auth },
          params,
        })).data
        this.modifiedCriteria = false
      } catch (err) {
        window.alert(err.toString())
      }
    },
    onChangeOrderCoupons(list) {
      this.selectedOrderCoupons = list
      this.modifiedCriteria = true
    },
    onChangeOrderSources(list) {
      this.selectedOrderSources = list
      this.modifiedCriteria = true
    },
    onChangePaymentTypes(list) {
      this.selectedPaymentTypes = list
      this.modifiedCriteria = true
    },
    onChangeProducts(list) {
      this.selectedProducts = list
      this.modifiedCriteria = true
    },
    onChangeTicketClasses(list) {
      this.selectedTicketClasses = list
      this.modifiedCriteria = true
    },
    onChangeUsedAtEvents(list) {
      this.selectedUsedAtEvents = list
      this.modifiedCriteria = true
    },
    onOrderNumber(evt) {
      const trim = evt.target.value.replace(/[^0-9]/g, '')
      if (trim !== evt.target.value) {
        evt.target.value = trim
        return
      }
      this.ordernum = parseInt(evt.target.value, 10) || 0
    },
    onReset() {
      this.customer = ''
      this.createdAfter = ''
      this.createdBefore = ''
      this.selectedOrderCoupons = []
      this.selectedOrderSources = []
      this.selectedPaymentTypes = []
      this.selectedProducts = []
      this.selectedTicketClasses = []
      this.selectedUsedAtEvents = []
      this.onSearch()
    },
    onSearch() {
      const params = new URLSearchParams()
      if (this.customer) params.append("customer", this.customer)
      if (this.createdAfter) params.append("createdAfter", this.createdAfter + "T00:00:00")
      if (this.createdBefore) params.append("createdBefore", this.createdBefore + "T23:59:59")
      this.selectedOrderCoupons.forEach(v => { params.append("orderCoupon", v) })
      this.selectedOrderSources.forEach(v => { params.append("orderSource", v) })
      this.selectedPaymentTypes.forEach(v => { params.append("paymentType", v) })
      this.selectedProducts.forEach(v => { params.append("product", v) })
      this.selectedTicketClasses.forEach(v => { params.append("ticketClass", v) })
      this.selectedUsedAtEvents.forEach(v => { params.append("usedAtEvent", v) })
      this.getStats(params)
    },
    onSearchList() {
      if (this.modifiedCriteria) { // Search button
        this.onSearch()
      } else { // Show List button
        this.$emit('report', this.stats)
      }
    },
    sortTree(tree) {
      for (const key in tree) {
        if (tree[key].children) tree[key].children = this.sortTree(tree[key].children)
      }
      return Object.keys(tree).sort((a, b) => (tree[a].id || tree[a].label).localeCompare(tree[b].id || tree[b].label)).map(key => tree[key])
    },
  }
}
</script>

<style lang="stylus">
#osearch
  display flex
  flex-direction column
  height 100%
#osearch-cols
  display flex
  flex 1
  flex-wrap wrap
  overflow-y auto
  padding 0.5em 1em
#osearch-spinner
  margin-top 2em
  text-align center
.osearch-col
  margin 0 2em 2em 0
  width 240px
.osearch-criteria-section
  overflow hidden
  text-overflow ellipsis
  white-space nowrap
  font-weight bold
#osearch-ordernum-row
  display flex
  width 240px
#osearch-criteria-ordernum
  flex auto
  min-width 0
  line-height 1
#osearch-viewedit
  flex none
  margin-left 6px
  padding 0 0.75rem
  line-height 1
#osearch-criteria-customer
  width 100%
  line-height 1
.osearch-criteria-date-box
  display flex
  justify-content space-between
  align-items center
  width 100%
  line-height 1
.osearch-criteria-date
  margin-bottom 2px
  width calc(100% - 3em)
#osearch-footer
  display flex
  flex none
  flex-wrap wrap
  justify-content space-between
  align-items center
  padding 0.5em 1em
  border-top 1px solid rgba(0, 0, 0, 0.125) // to match bottom border of tabs
#osearch-stats
  color #888
#osearch-buttons
  display flex
  flex none
#osearch-reset
  margin-right 1em
#osearch-search-list
  width 100px
</style>
