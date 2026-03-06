<script setup lang="ts">
import { computed, useAttrs } from 'vue'
import { NModal } from 'naive-ui'
const props = defineProps<{
  title?: string
  show: boolean
  size?: 'medium' | 'small' | 'large' | 'huge' | undefined
}>()

const emit = defineEmits<Emit>()
interface Emit {
  (e: 'update:show', show: boolean): void
//   (e: 'done', item: Panel.Info): void// 创建完成
}

const attrs = useAttrs()

const bindAttrs = computed<{ class: string; style: string }>(() => ({
  class: (attrs.class as string) || '',
  style: (attrs.style as string) || '',
}))

// 更新值父组件传来的值
const showModal = computed({
  get: () => props.show,
  set: (show: boolean) => {
    emit('update:show', show)
  },
})
</script>

<template>
  <NModal v-model:show="showModal" preset="card" :size="size" v-bind="bindAttrs" style="border-radius: 1rem;" :style="$parent" :title="title">
    <template #cover>
      <slot name="cover" />
    </template>
    <template #header>
      <slot name="header" />
    </template>
    <template #header-extra>
      <slot name="header-extra" />
    </template>
    <template #footer>
      <slot name="footer" />
    </template>
    <template #action>
      <slot name="action" />
    </template>
    <slot />
  </NModal>
</template>
