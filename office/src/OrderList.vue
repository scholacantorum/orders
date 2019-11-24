<!--
OrderList displays the order search results as a list.
-->

<template lang="pug">
#olist
  #olist-table
    b-table(
      borderless
      :fields="fields"
      :items="report.lines"
      small
      striped
    )
  #olist-footer
    #olist-stats(v-text="statsString")
    #olist-buttons
      b-button#olist-close(variant="outline-primary" @click="$emit('close')") Close
      b-button(variant="primary" @click="onExport") Export
</template>

<script>
const sourceLabels = {
  gala: 'Gala',
  inperson: 'Box Office',
  members: 'Members Site',
  office: 'Schola Office',
  public: 'Web Site',
}

function csvEncode(s) {
  if (s.indexOf(',') >= 0 || s.indexOf('"') >= 0 || s.indexOf('\n') >= 0) {
    return `"${s.replace('"', '""')}"`
  } else {
    return s
  }
}

export default {
  props: {
    report: Object,
  },
  data: () => ({
    fields: [
      { key: 'orderID', label: 'Order', sortable: true, sortdirection: 'desc' },
      {
        key: 'orderTime', label: 'Date/Time', sortable: true, sortdirection: 'desc',
        formatter: dt => dt.substr(0, 16).replace('T', ' '),
      },
      { key: 'name', label: 'Customer', sortable: true },
      { key: 'email', label: 'Email', sortable: true },
      { key: 'quantity', label: 'Qty', sortable: true, tdClass: 'report-table-right' },
      { key: 'product', label: 'Product', sortable: true },
      {
        key: 'usedAtEvent', label: 'Used At', sortable: true,
        formatter: uae => uae || '(unused)',
      },
      {
        key: 'orderSource', label: 'Source', sortable: true,
        formatter: os => sourceLabels[os] || os,
      },
      {
        key: 'paymentType', label: 'Payment', sortable: true,
        formatter: pt => pt.replace(',', ', '),
      },
      {
        key: 'amount', label: 'Amount', sortable: true, tdClass: 'olist-table-right',
        formatter: a => `$${a.toFixed(2)}`,
      },
    ],
  }),
  computed: {
    statsString() {
      return `Matched ${this.report.orderCount} ${this.report.orderCount === 1 ? 'order' : 'orders'}, ${this.report.itemCount} ${this.report.itemCount === 1 ? 'item' : 'items'}, total amount $${this.report.totalAmount.toFixed(2)}.`
    },
  },
  methods: {
    onExport() {
      let csv = 'Order,Date/Time,Customer,Email,Qty,Product,Used At,Source,Payment,Amount\r\n'
      this.report.lines.forEach(line => {
        csv += `${line.orderID},${line.orderTime.substr(0, 16).replace('T', ' ')},${csvEncode(line.name)},${csvEncode(line.email)},${line.quantity},${csvEncode(line.product)},${line.usedAtEvent ? csvEncode(line.usedAtEvent) : '(unused)'},${sourceLabels[line.orderSource] || csvEncode(line.orderSource)},${csvEncode(line.paymentType.replace(',', ', '))},${line.amount.toFixed(2)}\r\n`
      })
      const elm = document.createElement('a');
      elm.href = 'data:text/csv;charset=utf-8,' + encodeURI(csv);
      elm.target = '_blank';
      elm.download = 'orders.csv';
      elm.click();
    },
  }
}
</script>

<style lang="stylus">
#olist
  display flex
  flex-direction column
  height 100%
#olist-table
  flex auto
  overflow auto
  td
    white-space nowrap
.olist-table-right
  text-align right
#olist-footer
  display flex
  flex none
  flex-wrap wrap
  justify-content space-between
  align-items center
  padding 0.5em 1em 0.5em 0.3em // 0.3 to match padding on table cells
  border-top 1px solid rgba(0, 0, 0, 0.125) // to match bottom border of tabs
#olist-stats
  color #888
#olist-buttons
  display flex
  flex none
#olist-close
  margin-right 1em
</style>
