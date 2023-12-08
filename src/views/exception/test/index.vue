<script setup lang="ts">
import { computed, ref } from 'vue'
import { NColorPicker } from 'naive-ui'
import { SvgIcon } from '@/components/common'
// 获取背景颜色的 RGB 值（假设这里获取到了背景颜色的值）
const bgColor = ref('#000000') // 假设为蓝色

// 计算颜色的明暗度
const calculateLuminance = (color: string) => {
  const hex = color.replace(/^#/, '')
  const r = parseInt(hex.substring(0, 2), 16)
  const g = parseInt(hex.substring(2, 4), 16)
  const b = parseInt(hex.substring(4, 6), 16)
  return (0.299 * r + 0.587 * g + 0.114 * b) / 255
}

// 根据明暗度切换文本颜色
const textColor = computed(() => {
  const luminance = calculateLuminance(bgColor.value)
  return luminance > 0.5 ? 'black' : 'white'
})
</script>

<template>
  <div :style="{ backgroundColor: bgColor, color: textColor }">
    Background Color Example
  </div>
  <NColorPicker v-model:value="bgColor" />

  <div>
    <SvgIcon icon="mingcute-more-1-fill" style="color: pink;height: 50px;width: 50px;" />
  </div>
</template>

<style scoped>
/* 样式可以根据实际需求自定义 */
div {
  width: 200px;
  height: 100px;
  display: flex;
  justify-content: center;
  align-items: center;
  font-size: 18px;
}
</style>
