<!--
ReportCriteria
-->

<template lang="pug">
#report-criteria
  #report-criteria-heading Report Criteria
  .report-criteria-section Customer
  input#report-criteria-customer(v-model.lazy.trim="customer" type="text" placeholder="Name or Email")
  .report-criteria-section Order Dates
  .report-criteria-date-box
    | From
    input.report-criteria-date(v-model="createdAfter" type="date" placeholder="From Date")
  .report-criteria-date-box
    | To
    input.report-criteria-date(v-model="createdBefore" type="date" placeholder="To Date")
  .report-criteria-section Products
  TreeSelect(:tree="productsTree" :value="selectedProducts" @change="onChangeProducts")
  .report-criteria-section Tickets Used At
  TreeSelect(:tree="usedAtEventsTree" :value="selectedUsedAtEvents" @change="onChangeUsedAtEvents")
  .report-criteria-section Ticket Classes
  TreeSelect(:tree="ticketClassList" :value="selectedTicketClasses" @change="onChangeTicketClasses")
  .report-criteria-section Order Sources
  TreeSelect(:tree="orderSourcesList" :value="selectedOrderSources" @change="onChangeOrderSources")
  .report-criteria-section Payment Types
  TreeSelect(:tree="paymentTypesTree" :value="selectedPaymentTypes" @change="onChangePaymentTypes")
  .report-criteria-section Coupon Codes
  TreeSelect(:tree="orderCouponsList" :value="selectedOrderCoupons" @change="onChangeOrderCoupons")
</template>

<script>
import TreeSelect from './TreeSelect'

const productTypeLabels = {
  auctionitem: 'Auction Items',
  donation: 'Donations',
  recording: 'Recordings',
  sheetmusic: 'Sheet Music',
  ticket: 'Tickets',
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
  props: {
    stats: Object,
  },
  data: () => ({
    createdAfter: '',
    createdBefore: '',
    customer: '',
    selectedOrderCoupons: [],
    selectedOrderSources: [],
    selectedPaymentTypes: [],
    selectedProducts: [],
    selectedTicketClasses: [],
    selectedUsedAtEvents: [],
    updateTimer: null,
  }),
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
  },
  watch: {
    createdAfter: 'startUpdateTimer',
    createdBefore: 'startUpdateTimer',
    customer: 'startUpdateTimer',
  },
  methods: {
    onChangeOrderCoupons(list) {
      this.selectedOrderCoupons = list
      this.startUpdateTimer()
    },
    onChangeOrderSources(list) {
      this.selectedOrderSources = list
      this.startUpdateTimer()
    },
    onChangePaymentTypes(list) {
      this.selectedPaymentTypes = list
      this.startUpdateTimer()
    },
    onChangeProducts(list) {
      this.selectedProducts = list
      this.startUpdateTimer()
    },
    onChangeTicketClasses(list) {
      this.selectedTicketClasses = list
      this.startUpdateTimer()
    },
    onChangeUsedAtEvents(list) {
      this.selectedUsedAtEvents = list
      this.startUpdateTimer()
    },
    sendUpdate() {
      this.updateTimer = null
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
      this.$emit('update', params)
    },
    sortTree(tree) {
      for (const key in tree) {
        if (tree[key].children) tree[key].children = this.sortTree(tree[key].children)
      }
      return Object.keys(tree).sort((a, b) => (tree[a].id || tree[a].label).localeCompare(tree[b].id || tree[b].label)).map(key => tree[key])
    },
    startUpdateTimer() {
      if (this.updateTimer) window.clearTimeout(this.updateTimer)
      this.updateTimer = window.setTimeout(this.sendUpdate, 1500)
    },
  },
}
</script>

<style lang="stylus">
#report-criteria
  overflow-y auto
  margin 0.5em
  width calc(100% - 1em)
  height calc(100% - 1em)
#report-criteria-heading
  color #0153A5
  font-weight bold
  font-size larger
.report-criteria-section
  overflow hidden
  margin-top 0.8em
  text-overflow ellipsis
  white-space nowrap
  font-weight bold
#report-criteria-customer
  width 100%
  line-height 1
.report-criteria-date-box
  display flex
  justify-content space-between
  align-items center
  width 100%
  line-height 1
.report-criteria-date
  margin-bottom 2px
  width calc(100% - 3em)
</style>
