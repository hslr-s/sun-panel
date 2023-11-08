<script setup lang="ts">
import { NAvatar, NImage } from 'naive-ui'
import { computed, withDefaults } from 'vue'
import { SvgIcon } from '@/components/common'

interface Prop {
  itemIcon?: Panel.ItemIcon | null
  size?: number // 默认70
}

const props = withDefaults(defineProps<Prop>(), { size: 70 })
const defaultStyle = { width: `${props.size}px`, height: `${props.size}px` }
const iconExt = computed(() => {
  return props.itemIcon?.src?.split('.').pop()
})
</script>

<template>
  <div class="overflow-hidden rounded-2xl" :style="defaultStyle">
    <slot>
      <div v-if="itemIcon">
        <div v-if="itemIcon?.itemType === 1">
          <NAvatar :size="props.size" :style="{ backgroundColor: itemIcon?.bgColor }">
            {{ itemIcon.text }}
          </NAvatar>
        </div>

        <div v-else-if="itemIcon?.itemType === 2">
          <div v-if="iconExt === 'svg'" :style="defaultStyle" class="flex justify-center items-center">
            <img :src="itemIcon?.src" class="w-[35px] h-[35px]">
            <!-- <object :data="itemIcon?.src" type="image/svg+xml" class="w-[35px] h-[35px]" style="fill: rgb(255, 255, 255) !important;" /> -->
          </div>
          <NImage v-else :style="defaultStyle" :src="itemIcon?.src" preview-disabled />
        </div>

        <div v-else-if="itemIcon?.itemType === 3">
          <NAvatar :size="props.size" :style="{ backgroundColor: itemIcon?.bgColor }">
            <SvgIcon style="font-size: 35px;" :icon="itemIcon.text" />
          </NAvatar>
        </div>
      </div>

      <div v-else>
        <NAvatar :size="props.size" />
      </div>
    </slot>
  </div>
</template>
