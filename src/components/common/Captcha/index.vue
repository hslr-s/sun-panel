<script setup lang='ts'>
import { NImage } from 'naive-ui'
import { ref } from 'vue'
defineProps<{
  src: string
}>()

const emit = defineEmits<{
  (event: 'click'): void
  (event: 'refresh'): void
}>()

const randCode = ref<string>('0')

function handleClick() {
  randCode.value = String(rand(100, 99999))
  emit('click')
}

function rand(min: number, max: number) {
  return Math.floor(Math.random() * (max - min)) + min
}

defineExpose({
  // 刷新验证码
  refresh() {
    handleClick()
  },
})
</script>

<template>
  <!-- <div> -->
  <NImage
    :src="`${src}?${randCode}`"
    :preview-disabled="true"
    @click="handleClick"
  />
  <!-- </div> -->
</template>
