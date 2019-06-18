<!--
Button is a button, styled the way I like them.
-->

<template>
  <TouchableHighlight :disabled="disabled" :style="buttonStyle" :onPress="onGuardedPress">
    <Text :style="buttonTextStyle">{{ title }}</Text>
  </TouchableHighlight>
</template>

<script>
const textStyles = [
  'color', 'fontFamily', 'fontSize', 'fontStyle', 'fontVariant', 'fontWeight', 'includeFontPadding', 'letterSpacing', 'lineHeight',
  'textAlign', 'textAlignVertical', 'textDecorationColor', 'textDecorationLine', 'textDecorationStyle', 'textShadowColor',
  'textShadowOffset', 'textShadowRadius', 'textTransform', 'writingDirection',
]

export default {
  props: {
    bstyle: Object,
    disabled: Boolean,
    onPress: Function,
    secondary: Boolean,
    title: String,
  },
  computed: {
    activeOpacity() { this.disabled ? 1 : null },
    buttonStyle() {
      const style = {
        backgroundColor: this.secondary ? (this.disabled ? '#ccc' : '#444') : (this.disabled ? '#ccf' : '#00f'),
        borderRadius: 6,
        paddingTop: 12,
        paddingBottom: 12,
        paddingLeft: 6,
        paddingRight: 6,
      }
      if (this.bstyle)
        for (let key in this.bstyle)
          if (!textStyles.includes(key))
            style[key] = this.bstyle[key]
      return style
    },
    buttonTextStyle() {
      const style = {
        fontSize: 24,
        fontWeight: 'bold',
        color: '#fff',
        textAlign: 'center',
      }
      if (this.bstyle)
        for (let key in this.bstyle)
          if (textStyles.includes(key))
            style[key] = this.bstyle[key]
      return style
    },
  },
  methods: {
    onGuardedPress() {
      if (!this.disabled) this.onPress()
    },
  },
}
</script>

<style lang="stylus"></style>
