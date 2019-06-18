<!--
ChooseEvent is the dialog for choosing what event to sell tickets for.
-->

<template>
  <View>
    <View :style="{ padding: 12, backgroundColor: '#017EFA' }">
      <Text :style="{ fontSize: 24, textAlign: 'center', color: '#fff' }">Choose Event</Text>
    </View>
    <ActivityIndicator v-if="!events" size="large" :style="{ alignSelf: 'center', marginTop: 24 }"/>
    <View v-else>
      <TouchableOpacity v-for="event in events" :key="event.id" :onPress="() => onPress(event)">
        <View :style="{ borderBottomWidth: 1, borderBottomColor: '#ccc', padding: 12 }">
          <Text :style="{ fontSize: 20 }">{{ event.start.substr(0, 10) }} {{ event.name }}</Text>
        </View>
      </TouchableOpacity>
    </View>
  </View>
</template>

<script>
import { Alert } from 'react-native'
import backend from './backend'

export default {
  data: () => ({ events: null }),
  async mounted() {
    this.events = await backend.eventlist()
  },
  methods: {
    keyExtractor(event) { return event.id },
    async onPress(event) {
      let products
      try {
        products = await backend.eventProducts(event)
      } catch (err) {
        Alert.alert('Server Error', err)
        return
      }
      products = products.products.filter(p => !p.message)
      if (!products.length) {
        Alert.alert('Event Not Available', 'No tickets are on sale for that event.')
        return
      }
      this.$store.commit('event', { event, products })
    },
  },
}
</script>

<style>
</style>
