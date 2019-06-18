<!--
WillCall displays the list of customers who have tickets that could be, or were,
used at the event, and lets the user consume some of them.
-->

<template>
  <View v-if="!orders" :style="{ flex: 1, justifyContent: 'center', alignItems: 'center' }">
    <ActivityIndicator size="large"/>
  </View>
  <View v-else-if="!selected" :style="{ flex: 1 }">
    <View :style="{ width: '100%', padding: 12 }">
      <TextInput
        v-model="search"
        :autoCapitalize="'none'"
        :autoCorrect="false"
        placeholder="Search Will Call List"
        returnKeyType="search"
        :style="{ width: '100%', fontSize: 24, borderWidth: 1, borderColor: '#ccc' }"
        textContentType="none"
      />
    </View>
    <ScrollView :style="{ flex: 1 }">
      <TouchableOpacity
        v-for="order in filteredOrders"
        :key="order.id"
        :style="{ padding: 12, borderBottomWidth: 1, borderBottomColor: '#ccc' }"
        :onPress="() => onSelect(order.id)"
      >
        <View :style="{ flexDirection: 'row', alignItems: 'flex-end' }">
          <Text :style="{ fontSize: 20, paddingRight: 24 }">{{ order.name }}</Text>
          <Text :style="{ fontSize: 16, color: '#888' }">#{{ order.id }}</Text>
        </View>
      </TouchableOpacity>
    </ScrollView>
    <View :style="{ flexDirection: 'row', justifyContent: 'center', margin: 12 }">
      <Button :bstyle="{ width: '40%', fontSize: 20 }" secondary title="Cancel" :onPress="onDone"/>
    </View>
  </View>
  <UseTickets v-else :orderID="selected" :onCancel="onCancelSelect" :onDone="onDone"/>
</template>

<script>
import { Alert } from 'react-native'
import Button from './Button'
import UseTickets from './UseTickets'
import backend from './backend'

export default {
  components: { Button, UseTickets },
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
    filteredOrders() {
      if (!this.search.trim()) return this.orders
      return this.orders.filter(
        o => o.name.toLowerCase().includes(this.search.toLowerCase()) || parseInt(this.search) === o.id
      )
    },
  },
  methods: {
    onCancelSelect() { this.selected = null },
    onSelect(orderid) { this.selected = orderid },
  },
}
</script>
