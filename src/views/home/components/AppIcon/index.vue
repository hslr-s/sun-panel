<script setup lang="ts">
import { onMounted } from 'vue'
import { ItemIcon } from '@/components/common'
import { PanelPanelConfigStyleEnum } from '@/enums'

interface Prop {
  itemInfo?: Panel.ItemInfo
  size?: number // 默认70
  forceBackground?: string // 强制背景色
  iconTextColor?: string
  iconTextInfoHideDescription: boolean
  style: PanelPanelConfigStyleEnum
}

const props = withDefaults(defineProps<Prop>(), {
  size: 70,
})

const defaultBackground = '#2a2a2a6b'

onMounted(() => {
  console.log(props.style)
})
</script>

<template>
  <div class="w-full">
    <!-- 详情图标 -->
    <div
      v-if="style === PanelPanelConfigStyleEnum.info"
      class="w-full rounded-2xl  transition-all duration-200 hover:shadow-[0_0_20px_10px_rgba(0,0,0,0.2)] flex"
      :style="{ background: itemInfo?.icon?.backgroundColor || defaultBackground }"
    >
      <div class="w-[70px] h-[70px]">
        <div class="w-[70px] h-full flex items-center justify-center ">
          <ItemIcon :item-icon="itemInfo?.icon" force-background="transparent" :size="50" class="overflow-hidden rounded-xl" />
        </div>
      </div>
      <div class="text-white m-[8px_8px_0_8px]" :style="{ color: iconTextColor }">
        <div class="font-semibold">
          <NEllipsis>
            {{ itemInfo?.title }}
          </NEllipsis>
        </div>
        <div v-if="!iconTextInfoHideDescription">
          <NEllipsis :line-clamp="2" class="text-xs">
            {{ itemInfo?.description }}
          </NEllipsis>
        </div>
      </div>
    </div>

    <!-- 极简图标（APP） -->
    <div v-if="style === PanelPanelConfigStyleEnum.icon">
      <div
        class="overflow-hidden rounded-2xl sunpanel w-[70px] h-[70px] mx-auto rounded-2xl transition-all duration-200 hover:shadow-[0_0_20px_10px_rgba(0,0,0,0.2)]"
        :title="itemInfo?.description"
      >
        <ItemIcon :item-icon="itemInfo?.icon" />
      </div>
      <div
        class="text-center app-icon-text-shadow cursor-pointer mt-[2px]"
        :style="{ color: iconTextColor }"
      >
        <span>{{ itemInfo?.title }}</span>
      </div>
    </div>
  </div>
</template>
