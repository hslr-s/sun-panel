<script lang="ts">
import { computed, defineComponent } from 'vue'

function compatibleName(inputString: string): string {
  // 使用正则表达式替换所有的冒号
  const resultString = inputString.replace(/:/g, '-')
  return resultString
}

export default defineComponent({
  name: 'SvgIcon',
  props: {
    icon: {
      type: String,
      required: true,
    },
    className: {
      type: String,
      default: '',
    },
  },
  setup(props) {
    const symbolId = computed(() => `#${compatibleName(props.icon)}`)
    const svgClass = computed(() => {
      if (props.className)
        return `svg-icon ${props.className}`

      return 'svg-icon'
    })

    return {
      symbolId,
      svgClass,
    }
  },
})
</script>

<template>
  <svg :class="svgClass" aria-hidden="true">
    <use class="svg-use" :href="symbolId" />
  </svg>
</template>

<style scoped>
  .svg-icon {
    width: 1em;
    height: 1em;
    fill: currentColor;
    overflow: hidden;
  }
</style>
