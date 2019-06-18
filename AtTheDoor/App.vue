<template>
  <View :style="{ flex: 1, backgroundColor: '#fff', marginBottom }">
    <StatusBar barStyle="light-content" background-color="#0153A5"/>
    <View :style="{ backgroundColor: '#0153A5', paddingTop: paddingTop+12, paddingBottom: 12 }">
      <Text :style="{ textAlign: 'center', fontSize: 24, color: '#fff' }">Schola At The Door</Text>
    </View>
    <Login v-if="!$store.state.auth"/>
    <ChooseEvent v-else-if="!$store.state.event"/>
    <Main v-else/>
  </View>
</template>

<script>
import { Platform, Dimensions } from "react-native";
import ChooseEvent from './src/ChooseEvent'
import Login from './src/Login'
import Main from './src/Main'
import reader from './src/reader'

export default {
  components: { ChooseEvent, Login, Main },
  computed: {
    isIPhoneX() {
      if (Platform.OS !== 'ios') return false
      const width = Dimensions.get('window').width
      const height = Dimensions.get('window').height
      return width === 812 || height === 812 || width === 896 || height === 896
    },
    paddingTop() { return this.isIPhoneX ? 24 : 0 },
    marginBottom() { return this.isIPhoneX ? 34 : 0 },
  },
  watch: {
    '$store.state.auth': {
      immediate: true,
      handler(newv, oldv) {
        if (newv && !oldv) reader.connectReader()
        if (oldv && !newv) reader.disconnectReader()
      }
    },
  },
}
</script>
