<script setup lang="ts">
import { NAvatar, NImage } from 'naive-ui'
import { computed, ref, withDefaults } from 'vue'
import { SvgIconOnline } from '@/components/common'

interface Prop {
  itemIcon?: Panel.ItemIcon | null
  size?: number // 默认70
  forceBackground?: string // 强制背景色
}

const props = withDefaults(defineProps<Prop>(), { size: 70 })
const defaultBackground = '#2a2a2a6b'
const defaultStyle = ref({
  width: `${props.size}px`,
  height: `${props.size}px`,
})
const iconExt = computed(() => {
  return props.itemIcon?.src?.split('.').pop()
})
</script>

<template>
  <div :style="defaultStyle">
    <slot>
      <div v-if="itemIcon">
        <div v-if="itemIcon?.itemType === 1">
          <NAvatar :size="props.size" :style="{ backgroundColor: (forceBackground ?? itemIcon?.backgroundColor) || defaultBackground }">
            {{ itemIcon.text }}
          </NAvatar>
        </div>

        <div v-else-if="itemIcon?.itemType === 2">
          <div v-if="iconExt === 'svg'" :style="{ backgroundColor: (forceBackground ?? itemIcon?.backgroundColor) || defaultBackground, ...defaultStyle }" class="flex justify-center items-center">
            <img :src="itemIcon?.src" class="w-[35px] h-[35px]">
          </div>
          <NImage v-else :style="{ backgroundColor: (forceBackground ?? itemIcon?.backgroundColor) || defaultBackground, ...defaultStyle }" :src="itemIcon?.src" preview-disabled />
        </div>

        <div v-else-if="itemIcon?.itemType === 3">
          <NAvatar :size="props.size" :style="{ backgroundColor: (forceBackground ?? itemIcon?.backgroundColor) || defaultBackground }">
            <SvgIconOnline style="font-size: 35px;" :icon="itemIcon.text" />
          </NAvatar>
        </div>
      </div>

      <div v-else>
        <NAvatar :size="props.size" />
      </div>
    </slot>
  </div>
</template>
