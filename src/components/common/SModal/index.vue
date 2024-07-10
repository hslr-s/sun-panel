<script setup lang="ts">
import { computed, nextTick, useAttrs, watch } from 'vue'
import { NModal } from 'naive-ui'
import startDrag from './Drag'

const props = defineProps<{
  title?: string
  show: boolean
  size?: 'medium' | 'small' | 'large' | 'huge' | undefined
  move?: boolean
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

const modalId = `s-model-${getRandomInt(10000, 99999)}`
const modalMoveBarId = `${modalId}-move-bar`

// 更新值父组件传来的值
const showModal = computed({
  get: () => props.show,
  set: (show: boolean) => {
    emit('update:show', show)
  },
})

watch(() => showModal.value, (v: boolean) => {
  if (v && props.move) {
    nextTick(() => {
      const oBox = document.getElementById(modalId)
      const oBar = document.getElementById(modalMoveBarId)
      startDrag(oBar as HTMLElement, oBox as HTMLElement)
    })
  }
})

function getRandomInt(min: number, max: number): number {
  // 确保min小于max
  if (min > max)
    [min, max] = [max, min] // 交换min和max的值

  // 生成随机数并四舍五入到最近的整数
  return Math.floor(Math.random() * (max - min + 1)) + min
}
</script>

<template>
  <NModal
    :id="modalId"
    v-bind="bindAttrs"
    v-model:show="showModal"
    preset="card"
    :size="size"
    :style="$parent"
  >
    <template #header>
      <div :id="modalMoveBarId" class="w-full" :class="move ? 'cursor-move select-none' : ''">
        {{ title }}
        <slot name="header" />
      </div>
    </template>

    <template #cover>
      <slot name="cover" />
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
