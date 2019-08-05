<!--
Report displays the report generator page.
-->

<template lang="pug">
Split#report
  SplitArea(:size="10" :minSize="200")
    #report-spinner(v-if="!report")
      b-spinner(label="Loading...")
    ReportCriteria(v-else :stats="report" @update="onUpdate")
  SplitArea(:size="90")
    #report-message(v-if="!report")
    #report-message(v-else-if="!haveParams") For a list of orders, please provide search criteria.
    #report-message(v-else-if="!report.lines") Too many results; please narrow the search criteria.
    #report-message(v-else-if="!report.lines.length") No orders match your search criteria.
    ReportTable(v-else :lines="report.lines")
</template>

<script>
import Split from 'vue-split-panel/src/components/split'
import ReportCriteria from './ReportCriteria'
import ReportTable from './ReportTable'

export default {
  components: { ReportCriteria, ReportTable, Split, SplitArea: Split.SplitArea },
  data: () => ({
    haveParams: false,
    report: null,
  }),
  mounted() {
    this.runReport(null)
  },
  methods: {
    onUpdate(params) {
      this.haveParams = params.toString() !== ''
      this.runReport(params)
    },
    async runReport(params) {
      try {
        this.report = (await this.$axios.get('/api/report', {
          headers: { Auth: this.$store.state.auth },
          params,
        })).data
      } catch (err) {
        window.alert(err.toString())
      }
    },
  }
}
</script>

<style lang="stylus">
#report
  height 100%
#report-spinner
  margin-top 2em
  text-align center
#report-message
  margin 0.5em
  font-weight bold
</style>
