<script setup lang="ts">
import { onMounted, onUnmounted, ref } from 'vue'
import { NProgress } from 'naive-ui'
import type { ProgressStyle } from '../typings'
import { getCpuState } from '@/api/system/systemMonitor'
import { PanelPanelConfigStyleEnum } from '@/enums'

interface Prop {
  cardTypeStyle: PanelPanelConfigStyleEnum
  progressStyle: ProgressStyle
  refreshInterval: number
}

const props = defineProps<Prop>()
let timer: NodeJS.Timer
const cpuState = ref<SystemMonitor.CPUInfo | null>(null)

function correctionNumber(v: number, keepNum = 2): number {
  return v === 0 ? 0 : Number(v.toFixed(keepNum))
}

async function getData() {
  try {
    const { data, code } = await getCpuState<SystemMonitor.CPUInfo>()
    if (code === 0)
      cpuState.value = data
  }
  catch (error) {

  }
}

onMounted(() => {
  getData()
  timer = setInterval(() => {
    getData()
  }, (!props.refreshInterval || props.refreshInterval <= 2000) ? 2000 : props.refreshInterval)
})

onUnmounted(() => {
  clearInterval(timer)
})
</script>

<template>
  <div class="w-full">
    <div v-if="cardTypeStyle === PanelPanelConfigStyleEnum.info">
      <div class="mb-1">
        <span>
          CPU
        </span>
        <span class="float-right">
          {{ correctionNumber(cpuState?.usages[0] || 0) }}%
        </span>
      </div>
      <NProgress
        type="line"
        :color="progressStyle.color"
        :rail-color="progressStyle.railColor"
        :height="progressStyle.height"
        :percentage="correctionNumber(cpuState?.usages[0] || 0)"
        :show-indicator="false"
        :stroke-width="15"
        style="max-width: 135px;"
      />
    </div>
    <div v-else>
      <div class="flex justify-center h-full w-full mt-3">
        <NProgress
          :color="progressStyle.color"
          :rail-color="progressStyle.railColor"
          type="dashboard"
          :percentage="correctionNumber(cpuState?.usages[0] || 0)" :stroke-width="15"
          style="width: 50px;"
        >
          <div class="text-white" style="font-size: 10px;">
            {{ correctionNumber(cpuState?.usages[0] || 0, 1) }}%
          </div>
        </NProgress>
      </div>
    </div>
  </div>
</template>
