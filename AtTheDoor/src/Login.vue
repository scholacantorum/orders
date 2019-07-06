<!--
Login displays the login page.
-->

<template>
  <View :style="{ marginTop: 24, width: '100%', maxWidth: 200, alignSelf: 'center' }">
    <Text :style="{ fontSize: 24, fontWeight: 'bold', marginBottom: 24 }">Please log in.</Text>
    <Text :style="{ width: 110, fontSize: 20 }">Username</Text>
    <TextInput
      v-model.trim="username"
      :autoCapitalize="'none'"
      autoCompleteType="username"
      :autoCorrect="false"
      :autoFocus="true"
      enablesReturnKeyAutomatically
      returnKeyType="next"
      :style="{ width: '100%', fontSize: 20, borderColor: '#ccc', borderWidth: 1, marginBottom: 12 }"
      textContentType="username"
      :onSubmitEditing="focusPassword"
    />
    <Text :style="{ width: 110, fontSize: 20 }">Password</Text>
    <TextInput
      ref="password"
      v-model="password"
      :autoCapitalize="'none'"
      autoCompleteType="password"
      :autoCorrect="false"
      clearTextOnFocus
      enableReturnKeyAutomatically
      returnKeyType="go"
      secureTextEntry
      :style="{ width: '100%', fontSize: 20, borderColor: '#ccc', borderWidth: 1, marginBottom: 12 }"
      textContentType="password"
      :onSubmitEditing="onLogin"
    />
    <View
      :style="{ flexDirection: 'row', alignItems: 'center', justifyContent: 'space-between', width: '100%', marginBottom: 24}"
    >
      <Text :style="{ fontSize: 20 }">Test Mode</Text>
      <Switch :value="testmode" :onValueChange="onTestModeChange" />
    </View>
    <View
      :style="{ flexDirection: 'row', alignItems: 'flex-start', justifyContent: 'space-between', width: '100%', marginBottom: 24}"
    >
      <Text :style="{ fontSize: 20 }">Allow</Text>
      <View>
        <View :style="{ flexDirection: 'row' }">
          <Switch :value="allowCard" :onValueChange="onAllowCardChange" />
          <Text :style="{ fontSize: 20 }">Credit/debit card sales</Text>
        </View>
        <View :style="{ flexDirection: 'row' }">
          <Switch :value="allowCash" :onValueChange="onAllowCashChange" />
          <Text :style="{ fontSize: 20 }">Cash/check sales</Text>
        </View>
        <View :style="{ flexDirection: 'row' }">
          <Switch :value="allowWillCall" :onValueChange="onAllowWillCallChange" />
          <Text :style="{ fontSize: 20 }">Will call</Text>
        </View>
      </View>
    </View>
    <Text
      v-if="error"
      :style="{ fontSize: 16, color: 'red', alignSelf: 'center', marginBottom: 12 }"
    >{{ error }}</Text>
    <Button :style="{ alignSelf: 'center' }" title="Login" :onPress="onLogin" />
  </View>
</template>

<script>
import { Alert } from 'react-native'
import backend from './backend'

export default {
  data: () => ({
    allowCard: true,
    allowCash: true,
    allowWillCall: true,
    error: null,
    password: '',
    testmode: false,
    username: '',
  }),
  methods: {
    focusPassword() {
      this.$refs.password.focus()
    },
    onAllowCardChange(v) { this.allowCard = v },
    onAllowCashChange(v) { this.allowCash = v },
    onAllowWillCallChange(v) { this.allowWillCall = v },
    onTestModeChange(v) { this.testmode = v },
    async onLogin() {
      if (!this.username || !this.password) {
        this.error = 'Please enter username and password.'
        return
      }
      try {
        this.error = await backend.login(this.username, this.password, this.testmode, { card: this.allowCard, cash: this.allowCash, willcall: this.allowWillCall })
      } catch (err) {
        Alert.alert('Server Error', err)
        this.error = null
      }
    },
  }
}
</script>
