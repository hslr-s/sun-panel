<script setup lang="ts">
import { onMounted, onUnmounted, ref } from 'vue'
import GenericProgress from '../components/GenericProgress/index.vue'
import { correctionNumber, correctionNumberByCardStyle } from './common'
import { getCpuState } from '@/api/system/systemMonitor'
import type { PanelPanelConfigStyleEnum } from '@/enums'

interface Prop {
  cardTypeStyle: PanelPanelConfigStyleEnum
  refreshInterval: number
  textColor: string
  progressColor: string
  progressRailColor: string
}

const props = defineProps<Prop>()
let timer: NodeJS.Timer
const cpuState = ref<SystemMonitor.CPUInfo | null>(null)

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
  //@ts-ignore
  clearInterval(timer)
})
</script>

<template>
  <GenericProgress
    :progress-color="progressColor"
    :progress-rail-color="progressRailColor"
    :progress-height="5"
    :percentage="correctionNumberByCardStyle(cpuState?.usages[0] || 0, cardTypeStyle)"
    :card-type-style="cardTypeStyle"
    :info-card-right-text="`${correctionNumber(cpuState?.usages[0] || 0)}%`"
    info-card-left-text="CPU"
    :text-color="textColor"
    style="width: 100%;"
  />
</template>
