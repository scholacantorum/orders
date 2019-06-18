<!--
UseTickets displays the used and available tickets on an order and allows them
to be consumed (or unconsumed).
-->

<template>
  <View v-if="!order" :style="{ flex: 1, justifyContent: 'center', alignItems: 'center' }">
    <ActivityIndicator size="large"/>
  </View>
  <View v-else :style="{ flex: 1, marginBottom: 12 }">
    <View :style="{ padding: 12, backgroundColor: '#ccc', alignItems: 'center' }">
      <Text :style="{ fontSize: 24, fontWeight: 'bold' }">{{ order.name }}</Text>
      <Text :style="{ fontSize: 20 }">Order #{{ order.id }}</Text>
    </View>
    <ScrollView :style="{ flex: 1 }">
      <UseTicketClass
        v-for="tclass in order.classes"
        :key="tclass.name"
        :tclass="tclass"
        :onChange="n => onChange(tclass, n)"
      />
    </ScrollView>
    <View :style="{ flexDirection: 'row', justifyContent: 'space-evenly', margin: 12 }">
      <Button
        :bstyle="{ width: '40%', fontSize: 20 }"
        :disabled="disabled"
        title="Use Tickets"
        :onPress="onUseTickets"
      />
      <Button
        :bstyle="{ width: '40%', fontSize: 20 }"
        secondary
        title="Cancel"
        :onPress="onCancel"
      />
    </View>
  </View>
</template>

<script>
import { Alert } from 'react-native'
import Button from './Button'
import UseTicketClass from './UseTicketClass'
import backend from './backend'

export default {
  components: { Button, UseTicketClass },
  props: {
    orderID: [Number, String],
    onCancel: Function,
    onDone: Function,
  },
  data: () => ({ order: null }),
  async mounted() {
    try {
      const order = await backend.fetchTicketUsage(this.$store.state.event.id, this.orderID)
      if (order.error)
        Alert.alert('Order Error', order.error, () => { this.onCancel() })
      else
        this.order = order
    } catch (err) {
      Alert.alert('Server Error', err, () => { this.onCancel() })
    }
  },
  computed: {
    disabled() { return this.order.classes.every(cl => cl.min === cl.used) },
  },
  methods: {
    onChange(tclass, n) {
      tclass.used = n
      tclass.overflow = false
    },
    async onUseTickets() {
      try {
        await backend.useTickets(this.$store.state.event.id, this.order)
        this.onDone()
      } catch (err) {
        Alert.alert('Server Error', err, () => { this.onCancel() })
      }
    },
  },
}
</script>
