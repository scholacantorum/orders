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
      <Switch v-model="testmode"/>
    </View>
    <Text
      v-if="error"
      :style="{ fontSize: 16, color: 'red', alignSelf: 'center', marginBottom: 12 }"
    >{{ error }}</Text>
    <Button :style="{ alignSelf: 'center' }" title="Login" :onPress="onLogin"/>
  </View>
</template>

<script>
import { Alert } from 'react-native'
import backend from './backend'

export default {
  data: () => ({
    error: null,
    password: '',
    testmode: true,
    username: '',
  }),
  methods: {
    focusPassword() {
      this.$refs.password.focus()
    },
    async onLogin() {
      if (!this.username || !this.password) {
        this.error = 'Please enter username and password.'
        return
      }
      try {
        this.error = await backend.login(this.username, this.password, this.testmode)
      } catch (err) {
        Alert.alert('Server Error', err)
        this.error = null
      }
    },
  }
}
</script>
