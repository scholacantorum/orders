/**
 * @format
 */

import { AppRegistry } from 'react-native';
import App from './App';
import { name as appName } from './app.json';

AppRegistry.registerComponent(appName, () => App);

// import Vue from "vue-native-core";
// import { VueNativeBase } from "native-base";

// registering all native-base components to the global scope of the Vue
// Vue.use(VueNativeBase);

import './src/store'
// import './src/reader'
// import './src/event'
