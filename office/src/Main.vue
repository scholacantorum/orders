<!--
Main displays the main UI, after preliminaries like logging in and choosing an
event are taken care of.
-->

<template lang="pug">
b-card#main(no-body)
  b-tabs#main-tabs(v-model="currentTab" card content-class="main-tab-content" nav-wrapper-class="main-tabs-header")
    b-tab.main-tab-pane(title="Order Search")
      OrderSearch(@report="onOrderReport")
    b-tab.main-tab-pane(v-if="report" title="Order List" active)
      OrderList(:report="report" @close="onCloseList")
    b-tab.main-tab-pane(v-if="$store.state.privSetupOrders" title="Events")
      | Events
    b-tab.main-tab-pane(v-if="$store.state.privSetupOrders" title="Products")
      | Products
</template>

<script>
import OrderList from './OrderList'
import OrderSearch from './OrderSearch'

export default {
  components: { OrderList, OrderSearch },
  data: () => ({
    report: null,
    currentTab: 0,
  }),
  methods: {
    onCloseList() {
      this.currentTab = 0
      this.report = null
    },
    onOrderReport(report) {
      if (!report && this.report) this.currentTab = 0
      this.report = report
    },
  },
}
</script>

<style lang="stylus">
#main
  height 100%
#main-tabs
  display flex
  flex-direction column
  height 100%
.main-tabs-header
  flex none
.main-tab-content
  flex 1
  overflow-y auto
.main-tab-pane
  padding 0
  height 100%
</style>
