<!--
ReportTable displays the table of report results.
-->

<template lang="pug">
#report-table
  b-table(
    borderless
    :fields="fields"
    :items="lines"
    small
    striped
  )
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
        formatter: a => `$${a.toFixed(2)}`,
      },
    ],
  }),
}
</script>

<style lang="stylus">
#report-table
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
