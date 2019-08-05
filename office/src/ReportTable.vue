<!--
ReportTable displays the table of report results.
-->

<template lang="pug">
#report-table
  #report-table-box
    b-table(
      borderless
      :fields="fields"
      :items="lines"
      small
      striped
    )
  #report-table-stats(v-text="stats")
</template>

<script>
const sourceLabels = {
  gala: 'Gala',
  inperson: 'Box Office',
  members: 'Members Site',
  office: 'Schola Office',
  public: 'Web Site',
}

export default {
  props: {
    lines: Array,
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
        key: 'amount', label: 'Amount', sortable: true, tdClass: 'report-table-right',
        formatter: a => `$${Math.floor(a / 100)}`,
      },
    ],
  }),
  computed: {
    stats() {
      let lastID = 0
      let orders = 0
      let items = 0
      let amount = 0
      this.lines.forEach(l => {
        if (l.orderID !== lastID) {
          orders++
          lastID = l.orderID
        }
        items += l.quantity
        amount += l.amount
      })
      return `Matched ${orders} order${orders !== 1 ? 's' : ''}, ${items} item${items !== 1 ? 's' : ''}, total amount $${Math.floor(amount / 100)}.`
    },
  },
}
</script>

<style lang="stylus">
#report-table
  display flex
  flex-direction column
  height 100%
#report-table-box
  flex auto
  overflow auto
  td
    white-space nowrap
.report-table-right
  text-align right
#report-table-stats
  flex none
  margin 0.5em 0.5em 0.5em 0.3em // 0.3 to match padding on table cells
  font-weight bold
</style>
