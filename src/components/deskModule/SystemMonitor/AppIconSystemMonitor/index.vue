<script setup lang="ts">
import { computed } from 'vue'
import { MonitorType } from '../typings'
import type { CardStyle } from '../typings'
import GenericMonitorCard from '../components/GenericMonitorCard/index.vue'
import CardCPU from './CPU.vue'
import Memory from './Memory.vue'
import Disk from './Disk.vue'
import { SvgIcon } from '@/components/common'
import { PanelPanelConfigStyleEnum } from '@/enums'

interface Prop {
  extendParam?: any
  iconTextColor?: string
  iconTextIconHideTitle: boolean
  cardTypeStyle: PanelPanelConfigStyleEnum
  monitorType: string
  cardStyle: CardStyle
}

const props = withDefaults(defineProps<Prop>(), {})

const iconText = computed(() => {
  switch (props.monitorType) {
    case MonitorType.cpu:
      return 'CPU'
    case MonitorType.disk:
      return props.extendParam.path
    case MonitorType.memory:
      return 'RAM'
  }
  return ''
})
const refreshInterval = 5000
</script>

<template>
  <div class="w-full">
    <GenericMonitorCard
      :icon-text-icon-hide-title="false"
      :card-type-style="cardTypeStyle"
      :icon-text="iconText"
      :icon-text-color="iconTextColor"
      :background-color="extendParam.backgroundColor"
      class="hover:shadow-[0_0_20px_10px_rgba(0,0,0,0.2)]"
    >
      <template #icon>
        <!-- 图标 -->
        <div class="w-[60px] h-[70px]">
          <div class="app-icon w-[60px] h-full flex items-center justify-center text-white">
            <SvgIcon v-if="monitorType === MonitorType.cpu" icon="solar-cpu-bold" :style="{ color: extendParam.color }" style="width:35px;height:35px" />
            <SvgIcon v-if="monitorType === MonitorType.memory" icon="material-symbols-memory-alt-rounded" :style="{ color: extendParam.color }" style="width:35px;height:35px" />
            <SvgIcon v-if="monitorType === MonitorType.disk" icon="clarity-hard-disk-solid" :style="{ color: extendParam.color }" style="width:35px;height:35px" />
          </div>
        </div>
      </template>
      <template #info>
        <!-- 如果为纯白色，将自动根据背景的明暗计算字体的黑白色 -->
        <div
          class=" w-full text-white flex items-center"
        >
          <CardCPU
            v-if="monitorType === MonitorType.cpu"
            :card-type-style="PanelPanelConfigStyleEnum.info"
            :progress-color="extendParam?.progressColor"
            :progress-rail-color="extendParam?.progressRailColor"
            :text-color="extendParam?.color"
            :refresh-interval="refreshInterval"
          />
          <Memory
            v-else-if="monitorType === MonitorType.memory"
            :card-type-style="PanelPanelConfigStyleEnum.info"
            :progress-color="extendParam?.progressColor"
            :progress-rail-color="extendParam?.progressRailColor"
            :text-color="extendParam?.color"
            :refresh-interval="refreshInterval"
          />
          <Disk
            v-else-if="monitorType === MonitorType.disk"
            :card-type-style="PanelPanelConfigStyleEnum.info"
            :progress-color="extendParam?.progressColor"
            :progress-rail-color="extendParam?.progressRailColor"
            :text-color="extendParam?.color"
            :path="extendParam?.path"
            :refresh-interval="refreshInterval"
          />
        </div>
      </template>
      <template #small>
        <CardCPU
          v-if="monitorType === MonitorType.cpu"
          :card-type-style="PanelPanelConfigStyleEnum.icon"
          :progress-color="extendParam?.progressColor"
          :progress-rail-color="extendParam?.progressRailColor"
          :text-color="extendParam?.color"
          :refresh-interval="refreshInterval"
        />
        <Memory
          v-else-if="monitorType === MonitorType.memory"
          :card-type-style="PanelPanelConfigStyleEnum.icon"
          :progress-color="extendParam?.progressColor"
          :progress-rail-color="extendParam?.progressRailColor"
          :text-color="extendParam?.color"
          :refresh-interval="refreshInterval"
        />
        <Disk
          v-else-if="monitorType === MonitorType.disk"
          :card-type-style="PanelPanelConfigStyleEnum.icon"
          :progress-color="extendParam?.progressColor"
          :progress-rail-color="extendParam?.progressRailColor"
          :text-color="extendParam?.color"
          :path="extendParam?.path"
          :refresh-interval="refreshInterval"
        />
      </template>
    </GenericMonitorCard>
  </div>
</template>
