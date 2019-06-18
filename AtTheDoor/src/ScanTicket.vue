<!--
ScanTicket scans the barcode of a ticket and can mark it as used.
-->

<template>
  <View v-if="!selected" :style="{ flex: 1 }">
    <RNCamera
      :barCodeTypes="barCodeTypes"
      :captureAudio="false"
      :style="{ flex: 1 }"
      :onBarCodeRead="onBarCodeRead"
      :onMountError="onMountError"
    />
    <View :style="{ flexDirection: 'row', justifyContent: 'center', margin: 12 }">
      <Button :bstyle="{ width: '40%', fontSize: 20 }" secondary title="Cancel" :onPress="onDone"/>
    </View>
  </View>
  <UseTickets v-else :orderID="selected" :onCancel="onDone" :onDone="onDone"/>
</template>

<script>
import { Alert } from 'react-native'
import { RNCamera } from 'react-native-camera'
import Button from './Button'
import UseTickets from './UseTickets'
import backend from './backend'

export default {
  components: { Button, RNCamera, UseTickets },
  props: {
    onDone: Function,
  },
  data: () => ({
    orders: null,
    search: '',
    selected: null,
  }),
  async mounted() {
    try {
      this.orders = await backend.willCallList(this.$store.state.event.id)
    } catch (err) {
      Alert.alert('Will Call List Error', err, () => { this.onDone() })
    }
  },
  computed: {
    barCodeTypes() { return [RNCamera.Constants.BarCodeType.qr] },
    filteredOrders() {
      if (!this.search.trim()) return this.orders
      return this.orders.filter(
        o => o.name.toLowerCase().includes(this.search.toLowerCase()) || parseInt(this.search) === o.id
      )
    },
  },
  methods: {
    onMountError(err) { Alert.alert('Camera Error', err, this.onDone) },
    onBarCodeRead(evt) {
      if (!/\/ticket\/\d\d\d\d-\d\d\d\d-\d\d\d\d$/.test(evt.data)) {
        Alert.alert('Not a Schola Cantorum order bar code')
        return
      }
      this.selected = evt.data.substr(evt.data.length - 14)
    },
  },
}
</script>
