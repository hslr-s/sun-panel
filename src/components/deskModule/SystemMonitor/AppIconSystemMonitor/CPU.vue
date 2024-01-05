<script setup lang="ts">
import { onMounted, onUnmounted, ref } from 'vue'
import type { ProgressStyle } from '../typings'
import GenericProgress from '../components/GenericProgress/index.vue'
import { getCpuState } from '@/api/system/systemMonitor'
import type { PanelPanelConfigStyleEnum } from '@/enums'

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
    <GenericProgress
      :progress-style="progressStyle"
      :percentage="correctionNumber(cpuState?.usages[0] || 0)"
      :card-type-style="cardTypeStyle"
      :info-card-right-text="`${correctionNumber(cpuState?.usages[0] || 0)}%`"
      info-card-left-text="CPU"
      text-color="black"
    />
  </div>
</template>
