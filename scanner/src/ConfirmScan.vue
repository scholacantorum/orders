<!--
ConfirmScan displays the result of scanning a ticket.  If the result is an
error, it displays it.  If the result is a success, it confirms the count (if
needed) and then marks the tickets used.  It emits 'done' when finished.
-->

<template lang="pug">
#confirm
  #confirm-order(v-if="odata")
    #confirm-oname(v-if="odata.name" v-text="odata.name")
    #confirm-oid(v-text="`Order number ${odata.id}`")
  #confirm-error(v-if="error" @click="$emit('done')" v-text="error")
  template(v-else)
    #confirm-counts
      .confirm-class(
        v-for="cls, cname in odata.classes"
        :key="cname"
        :class="cname ? 'confirm-restricted' : null"
      )
        .confirm-cname(v-text="cname || 'General Admission'")
        CountChoice(v-model="cls.val" :max="cls.max" :zero="showZero")
    #confirm-buttons
      b-button#confirm-cancel(@click="$emit('done')") Cancel
      b-button#confirm-save(variant="success" :disabled="!canSave" @click="onSave") Use Tickets
</template>

<script>
import CountChoice from './CountChoice'

export default {
  components: { CountChoice },
  props: {
    event: String,
    ticket: String,
  },
  data: () => ({ error: null, odata: null }),
  watch: {
    ticket: {
      immediate: true,
      handler() {
        switch (this.ticket) {
          case 'non-schola':
            this.error = 'Not a Schola order'
            break
          case 'door':
            // TODO
            break
          case 'free':
            // TODO
            break
          default:
            this.fetchOrderData()
        }
      }
    },
  },
  computed: {
    canSave() {
      for (let cl in this.odata.classes) {
        if (this.odata.classes[cl].val) return true
      }
      return false
    },
    showZero() {
      let count = 0
      if (!this.odata.classes) return count
      // eslint-disable-next-line
      for (let _ in this.odata.classes) count++
      return count > 1
    },
  },
  methods: {
    async fetchOrderData() {
      const resp = await this.$axios.get(`/api/event/${this.event}/ticket/${this.ticket}`).catch(err => {
        console.log(err)
        this.error = 'Network failure'
        return null
      })
      if (!resp) return
      if (resp.status !== 200) {
        console.log(resp.statuText)
        this.error = 'Software error'
        return
      }
      if (resp.data.id) this.odata = resp.data
      if (resp.data.error) this.error = resp.data.error
    },
    async onSave() {
      const body = {}
      for (let cl in this.odata.classes) {
        if (this.odata.classes[cl].val)
          body[cl] = this.odata.classes[cl].val
      }
      const resp = await this.$axios.post(`/api/event/${this.event}/ticket/${this.ticket}`, JSON.stringify(body)).catch(err => {
        console.log(err)
        this.error = 'Network failure'
        return null
      })
      if (!resp) return
      switch (resp.status) {
        case 200:
          this.error = resp.data.error
          break
        case 204:
          this.$emit('done')
          break
        default:
          console.log(resp.status)
          this.error = 'Software error'
      }
    },
  },
}
</script>

<style lang="stylus">
#confirm
  display flex
  flex-direction column
  width 100%
  height 100%
#confirm-order
  display flex
  flex none
  flex-direction column
  justify-content center
  align-items center
  width 100%
  height 50px
  text-align center
#confirm-oname
  font-weight bold
  font-size 20px
  line-height 24px
#confirm-oid
  color #888
  font-size 16px
  line-height 16px
#confirm-error
  display flex
  flex 1 1 auto
  justify-content center
  align-items center
  background-color red
  color white
  text-align center
  font-size 48px
  line-height 1.2
#confirm-counts
  flex auto
.confirm-class
  padding 12px 6px 6px
.confirm-restricted
  background-color gold
#confirm-buttons
  display flex
  flex none
  justify-content space-between
  align-items center
  padding 12px
  width 100%
  height 74px
#confirm-cancel, #confirm-save
  height 50px
</style>
