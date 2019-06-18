<!--
SalesReceipt offers to send an email receipt to the purchaser of the completed
order.
-->

<template>
  <View :style="{ flex: 1, justifyContent: 'space-between' }">
    <Summary :order="order"/>
    <View :style="{ margin: 12, alignItems: 'center' }">
      <Text :style="{ fontSize: 24, fontWeight: 'bold' }">Email Receipt?</Text>
      <View :style="{ width: '100%', padding: 12 }">
        <TextInput
          v-model="email"
          :autoCapitalize="'none'"
          :autoCorrect="false"
          keyboardType="email-address"
          placeholder="Email address"
          returnKeyType="send"
          :style="{ width: '100%', fontSize: 24, textAlign: 'center', borderWidth: 1, borderColor: '#ccc' }"
          textContentType="none"
          type="email"
          :onSubmitEditing="onSendReceipt"
        />
      </View>
    </View>
    <View :style="{ flexDirection: 'row', justifyContent: 'space-evenly', marginBottom: 12 }">
      <Button
        :bstyle="{ width: '40%', fontSize: 20 }"
        :disabled="!validEmail || sending"
        :title="sending ? 'Sending...' : 'Send Receipt'"
        :onPress="onSendReceipt"
      />
      <Button
        :bstyle="{ width: '40%', fontSize: 20 }"
        secondary
        title="No Receipt"
        :onPress="onDone"
      />
    </View>
  </View>
</template>

<script>
import { Alert } from 'react-native'
import Button from './Button'
import Summary from './SalesSummary'
import backend from './backend'

export default {
  components: { Button, Summary },
  props: {
    order: Object,
    onDone: Function,
  },
  data: () => ({ email: null, sending: false }),
  mounted() {
    this.email = this.order.email
  },
  computed: {
    validEmail() {
      return /^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$/.test(this.email)
    },
  },
  methods: {
    async onSendReceipt() {
      this.sending = true
      try {
        await backend.sendEmailReceipt(this.order.id, this.email)
      } catch (err) {
        Alert.alert('Server Error', err)
      }
      this.onDone()
    },
  },
}
</script>
