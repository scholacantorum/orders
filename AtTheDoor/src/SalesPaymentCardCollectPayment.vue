<!--
SalesPaymentCardCollectPayment reads the payment method from the card reader,
displaying the card reader's prompts.  It also provides "Manual Entry" and
"Cancel" buttons.
-->

<template>
  <View :style="{ flex: 1, justifyContent: 'space-between' }">
    <Summary :order="order"/>
    <View :style="{ margin: 12 }">
      <Text
        :style="{ fontSize: 24, fontWeight: 'bold', color: '#f90', lineHeight: 30, textAlign: 'center' }"
      >{{ message }}</Text>
    </View>
    <View :style="{ flexDirection: 'row', justifyContent: 'space-evenly', marginBottom: 12 }">
      <Button
        :bstyle="{ width: '40%', fontSize: 20 }"
        title="Manual Entry"
        :onPress="myOnManualEntry"
      />
      <Button
        :bstyle="{ width: '40%', fontSize: 20 }"
        secondary
        title="Cancel"
        :onPress="myOnCancel"
      />
    </View>
  </View>
</template>

<script>
import { Alert } from 'react-native'
import Button from './Button'
import Summary from './SalesSummary'
import backend from './backend'
import reader from './reader'

export default {
  components: { Button, Summary },
  props: {
    order: Object,
    onCancel: Function,
    onCollected: Function,
    onManualEntry: Function,
  },
  data: () => ({ message: 'Swipe / Insert / Tap' }),
  async mounted() {
    if (!reader.isConnected) {
      this.onManualEntry()
      return
    }
    reader.addReaderMessageListener(this.onReaderMessage)
    while (true) {
      try {
        const intent = await reader.collectPaymentMethod()
        this.onCollected(intent)
        return
      } catch (err) {
        switch (err.code) {
          case 2020: // Canceled
            return
          case 2810: case 2820: case 2830: case 2840: case 2850:
            // Transient human errors.
            Alert.alert('Card Reader Error', err.error)
            break
          default:
            Alert.alert('Card Reader Error', err.error)
            this.onCancel()
            return
        }
      }
    }
  },
  beforeDestroy() {
    reader.removeReaderMessageListener(this.onReaderMessage)
  },
  methods: {
    myOnManualEntry() {
      reader.abortCollectPaymentMethod()
      this.onManualEntry()
    },
    myOnCancel() {
      reader.abortCollectPaymentMethod()
      this.onCancel()
    },
    onReaderMessage({ text }) { this.message = text },
  },
}
</script>
