<script setup lang="ts">
import { onMounted, onUnmounted, ref } from 'vue'
import { NProgress } from 'naive-ui'
import type { ProgressStyle } from '../typings'
import { getMemonyState } from '@/api/system/systemMonitor'
import { PanelPanelConfigStyleEnum } from '@/enums'
import { bytesToSize } from '@/utils/cmn'

interface Prop {
  cardTypeStyle: PanelPanelConfigStyleEnum
  progressStyle: ProgressStyle
  refreshInterval: number
}

const props = defineProps<Prop>()
let timer: NodeJS.Timer
const memoryState = ref<SystemMonitor.MemoryInfo | null>(null)

function correctionNumber(v: number, keepNum = 2): number {
  return v === 0 ? 0 : Number(v.toFixed(keepNum))
}

function formatMemorySize(v: number): string {
  return bytesToSize(v)
}

async function getData() {
  try {
    const { data, code } = await getMemonyState<SystemMonitor.MemoryInfo>()
    if (code === 0)
      memoryState.value = data
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
          RAM
        </span>
        <span class="float-right">
          {{ formatMemorySize(memoryState?.total || 0) }}/{{ formatMemorySize(memoryState?.free || 0) }}
        </span>
      </div>
      <NProgress
        type="line"
        :color="progressStyle.color"
        :rail-color="progressStyle.railColor"
        :height="progressStyle.height"
        :percentage="memoryState?.usedPercent"
        :show-indicator="false"
        :stroke-width="15" style="width: 150px;"
      />
    </div>
    <div v-else>
      <div class="flex justify-center h-full w-full mt-3">
        <NProgress
          :color="progressStyle.color"
          :rail-color="progressStyle.railColor"
          type="dashboard"
          :percentage="correctionNumber(memoryState?.usedPercent || 0)" :stroke-width="15"
          style="width: 50px;"
        >
          <div class="text-white" style="font-size: 10px;">
            {{ correctionNumber(memoryState?.usedPercent || 0, 1) }}%
          </div>
        </NProgress>
      </div>
    </div>
  </div>
</template>
