<script setup lang="ts">
// -------------------
// 系统监控图标临时使用（与 AppIcon/index.vue 一致）
// 如果确定这种方案将 AppIcon/index.vue 封装成通用组件
// -------------------

import { computed, ref } from 'vue'
import { MonitorType } from '../typings'
import type { CardStyle } from '../typings'
import CardCPU from './CPU.vue'
import Memory from './Memory.vue'
import Disk from './Disk.vue'
import { ItemCard, SvgIcon } from '@/components/common'
import { PanelPanelConfigStyleEnum } from '@/enums'

interface Prop {
  size?: number // 默认70
  extendParam?: any
  iconTextColor?: string
  iconTextIconHideTitle: boolean
  cardTypeStyle: PanelPanelConfigStyleEnum
  monitorType: string
  cardStyle: CardStyle
}

const props = withDefaults(defineProps<Prop>(), {
  size: 70,
})

const defaultBackground = '#2a2a2a6b'
const iconText = ref('自定义')
const refreshInterval = 5
const svgStyle = {
  width: '35px',
  height: '35px',
}

const calculateLuminance = (color: string) => {
  const hex = color.replace(/^#/, '')
  const r = parseInt(hex.substring(0, 2), 16)
  const g = parseInt(hex.substring(2, 4), 16)
  const b = parseInt(hex.substring(4, 6), 16)
  return (0.299 * r + 0.587 * g + 0.114 * b) / 255
}

const textColor = computed(() => {
  const luminance = calculateLuminance(props.cardStyle.background || defaultBackground)
  return luminance > 0.5 ? 'black' : 'white'
})
</script>

<template>
  <div class="w-full">
    <ItemCard
      icon-text-icon-hide-title
      :card-type-style="cardTypeStyle"
      :icon-text="iconText"
      class="hover:shadow-[0_0_20px_10px_rgba(0,0,0,0.2)]"
    >
      <template #info>
        <!-- 图标 -->
        <div class="w-[60px] h-[70px]">
          <div class="w-[60px] h-full flex items-center justify-center text-white">
            <SvgIcon v-if="monitorType === MonitorType.cpu" icon="solar-cpu-bold" :style="svgStyle" />
            <SvgIcon v-if="monitorType === MonitorType.memory" icon="material-symbols-memory-alt-rounded" :style="svgStyle" />
            <SvgIcon v-if="monitorType === MonitorType.disk" icon="clarity-hard-disk-solid" :style="svgStyle" />
          </div>
        </div>

        <!-- 如果为纯白色，将自动根据背景的明暗计算字体的黑白色 -->
        <div
          class=" w-full text-white flex items-center"
          :style="{ color: (iconTextColor === '#ffffff') ? textColor : iconTextColor, maxWidth: 'calc(100% - 80px)' }"
        >
          <CardCPU
            v-if="monitorType === MonitorType.cpu"
            :card-type-style="PanelPanelConfigStyleEnum.info"
            :progress-style="extendParam?.progressStyle"
            :refresh-interval="refreshInterval"
          />
          <Memory
            v-else-if="monitorType === MonitorType.memory"
            :card-type-style="PanelPanelConfigStyleEnum.info"
            :progress-style="extendParam?.progressStyle"
            :refresh-interval="refreshInterval"
          />
          <Disk
            v-else-if="monitorType === MonitorType.disk"
            :card-type-style="PanelPanelConfigStyleEnum.info"
            :progress-style="extendParam?.progressStyle"
            :path="extendParam?.path"
            :refresh-interval="refreshInterval"
          />
        </div>
      </template>
      <template #small>
        <CardCPU
          v-if="monitorType === MonitorType.cpu"
          :card-type-style="PanelPanelConfigStyleEnum.icon"
          :progress-style="extendParam?.progressStyle"
          :refresh-interval="refreshInterval"
        />
        <Memory
          v-else-if="monitorType === MonitorType.memory"
          :card-type-style="PanelPanelConfigStyleEnum.icon"
          :progress-style="extendParam?.progressStyle"
          :refresh-interval="refreshInterval"
        />
        <Disk
          v-else-if="monitorType === MonitorType.disk"
          :card-type-style="PanelPanelConfigStyleEnum.icon"
          :progress-style="extendParam?.progressStyle"
          :path="extendParam?.path"
          :refresh-interval="refreshInterval"
        />
      </template>
    </ItemCard>
  </div>
</template>
