<!--
UseTicketClass displays and changes the ticket usage for a single ticket class.
-->

<template>
  <View :style="classStyle">
    <View :style="{ width: 336 }">
      <Text :style="{ fontSize: 20 }">{{ tclass.name || 'General Admission' }}</Text>
      <View :style="{ flexDirection: 'row', flexWrap: 'wrap' }">
        <TouchableHighlight
          v-for="n in max"
          :key="n"
          :disabled="disabled(n)"
          :style="buttonStyle(n)"
          :onPress="() => onPress(n)"
        >
          <Text :style="buttonTextStyle(n)">{{ n }}</Text>
        </TouchableHighlight>
      </View>
    </View>
  </View>
</template>

<script>
export default {
  props: {
    tclass: Object,
    onChange: Function,
  },
  computed: {
    classStyle() {
      return {
        alignItems: 'center',
        paddingTop: 12,
        paddingLeft: 6,
        paddingRight: 6,
        paddingBottom: 6,
        backgroundColor: this.tclass.overflow ? 'red' : this.tclass.name ? 'gold' : '#fff',
      }
    },
    max() {
      if (this.tclass.max < 1000) return this.tclass.max
      return Math.ceil((this.tclass.used + 1) / 6) * 6
    },
  },
  methods: {
    buttonStyle(n) {
      const style = {
        width: 50,
        height: 50,
        marginRight: 6,
        marginBottom: 6,
        borderWidth: 2,
        borderRadius: 6,
        paddingTop: 12,
        paddingBottom: 12,
        paddingLeft: 6,
        paddingRight: 6,
      }
      if (n <= this.tclass.min) style.backgroundColor = style.borderColor = '#444'
      else if (n <= this.tclass.used) style.backgroundColor = style.borderColor = '#28a745'
      else {
        style.backgroundColor = 'transparent'
        style.borderColor = '#28a745'
      }
      return style
    },
    buttonTextStyle(n) {
      const style = { fontSize: 20, textAlign: 'center', color: '#fff' }
      if (n > this.tclass.used) style.color = '#28a745'
      return style
    },
    disabled(n) { return n < this.tclass.min },
    onPress(n) {
      if (n < this.tclass.min) return
      if (n === 1 && this.tclass.min === 0 && this.tclass.used === 1) this.onChange(0)
      else this.onChange(n)
    },
  },
}
</script>
