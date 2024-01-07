<script setup lang="ts">
import { computed } from 'vue'
import { NColorPicker, NForm, NFormItem } from 'naive-ui'
import type { ProgressStyle } from '../../typings'

const props = defineProps<{
  progressStyle: ProgressStyle
}>()

const emit = defineEmits<{
  (e: 'update:progressStyle', visible: ProgressStyle): void // 定义修改父组件（prop内）的值的事件
}>()

const data = computed({
  get: () => props.progressStyle,
  set() {
    emit('update:progressStyle', data.value)
  },
})
</script>

<template>
  <div>
    <NForm ref="formRef" :model="data">
      <NFormItem path="url" label="主色">
        <!-- <NSelect :style="{ width: '100px' }" :options="urlProtocolOptions" /> -->
        <NColorPicker v-model:value="data.color" size="small" />
      </NFormItem>
      <NFormItem path="url" label="背景色">
        <!-- <NSelect :style="{ width: '100px' }" :options="urlProtocolOptions" /> -->
        <NColorPicker v-model:value="data.railColor" size="small" />
      </NFormItem>
    </NForm>
  </div>
</template>
