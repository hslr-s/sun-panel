<script setup lang="ts">
import { onMounted, onUnmounted, ref } from 'vue'
import GenericProgress from '../components/GenericProgress/index.vue'
import { getMemonyState } from '@/api/system/systemMonitor'
import type { PanelPanelConfigStyleEnum } from '@/enums'
import { bytesToSize } from '@/utils/cmn'

interface Prop {
  cardTypeStyle: PanelPanelConfigStyleEnum
  refreshInterval: number
  textColor: string
  progressColor: string
  progressRailColor: string
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
    <GenericProgress
      :progress-color="progressColor"
      :progress-rail-color="progressRailColor"
      :progress-height="5"
      :percentage="correctionNumber(memoryState?.usedPercent || 0)"
      :card-type-style="cardTypeStyle"
      :info-card-right-text="`${formatMemorySize(memoryState?.total || 0)}/${formatMemorySize(memoryState?.free || 0)}`"
      info-card-left-text="RAM"
      :text-color="textColor"
    />
  </div>
</template>
