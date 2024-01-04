<script setup lang="ts">
// -------------------
// 系统监控图标临时使用（与 AppIcon/index.vue 一致）
// 如果确定这种方案将 AppIcon/index.vue 封装成通用组件
// -------------------

import { computed } from 'vue'
import { MonitorType } from '../typings'
import type { CardStyle } from '../typings'
import CardCPU from './CPU.vue'
import Memory from './Memory.vue'
import Disk from './Disk.vue'
import { SvgIcon } from '@/components/common'
import { PanelPanelConfigStyleEnum } from '@/enums'

interface Prop {
  size?: number // 默认70
  extendParam?: { [key: string]: [value: any] }
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
    <!-- 详情图标 -->
    <div
      v-if="cardTypeStyle === PanelPanelConfigStyleEnum.info"
      class="w-full rounded-2xl transition-all duration-200 hover:shadow-[0_0_20px_10px_rgba(0,0,0,0.2)] flex"
      :style="{ backgroundColor: cardStyle.background }"
    >
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
          :progress-style="extendParam?.progressStyle as any"
          :refresh-interval="refreshInterval"
        />
        <Memory
          v-else-if="monitorType === MonitorType.memory"
          :card-type-style="PanelPanelConfigStyleEnum.info"
          :progress-style="extendParam?.progressStyle as any"
          :refresh-interval="refreshInterval"
        />
        <Disk
          v-else-if="monitorType === MonitorType.disk"
          :card-type-style="PanelPanelConfigStyleEnum.info"
          :progress-style="extendParam?.progressStyle as any"
          :path="extendParam?.path as any"
          :refresh-interval="refreshInterval"
        />
      </div>
    </div>

    <!-- 极简图标（APP） -->
    <div
      v-if="cardTypeStyle === PanelPanelConfigStyleEnum.icon"
    >
      <div
        class="overflow-hidden rounded-2xl sunpanel w-[70px] h-[70px] mx-auto rounded-2xl transition-all duration-200 hover:shadow-[0_0_20px_10px_rgba(0,0,0,0.2)]"
        :style="{ backgroundColor: cardStyle.background }"
      >
        <CardCPU
          v-if="monitorType === MonitorType.cpu"
          :card-type-style="PanelPanelConfigStyleEnum.icon"
          :progress-style="extendParam?.progressStyle as any"
          :refresh-interval="refreshInterval"
        />
        <Memory
          v-else-if="monitorType === MonitorType.memory"
          :card-type-style="PanelPanelConfigStyleEnum.icon"
          :progress-style="extendParam?.progressStyle as any"
          :refresh-interval="refreshInterval"
        />
        <Disk
          v-else-if="monitorType === MonitorType.disk"
          :card-type-style="PanelPanelConfigStyleEnum.icon"
          :progress-style="extendParam?.progressStyle as any"
          :path="extendParam?.path as any"
          :refresh-interval="refreshInterval"
        />
      </div>
      <div
        v-if="!iconTextIconHideTitle"
        class="text-center app-icon-text-shadow cursor-pointer mt-[2px]"
        :style="{ color: iconTextColor }"
      >
        <span
          v-if="monitorType === MonitorType.cpu"
        >
          CPU
        </span>

        <span
          v-else-if="monitorType === MonitorType.memory"
        >
          RAM
        </span>

        <span
          v-else-if="monitorType === MonitorType.disk"
        >
          {{ extendParam?.path }}
        </span>
      </div>
    </div>
  </div>
</template>
