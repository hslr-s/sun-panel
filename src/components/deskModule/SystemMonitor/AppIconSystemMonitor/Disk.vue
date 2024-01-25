<script setup lang="ts">
import { onMounted, onUnmounted, ref } from 'vue'
import GenericProgress from '../components/GenericProgress/index.vue'
import { correctionNumberByCardStyle } from './common'
import type { PanelPanelConfigStyleEnum } from '@/enums'
import { bytesToSize } from '@/utils/cmn'
import { getDiskStateByPath } from '@/api/system/systemMonitor'

interface Prop {
  cardTypeStyle: PanelPanelConfigStyleEnum
  refreshInterval: number
  textColor: string
  progressColor: string
  progressRailColor: string
  path: string
}

const props = defineProps<Prop>()
let timer: NodeJS.Timer
const diskState = ref<SystemMonitor.DiskInfo | null>(null)

function formatdiskSize(v: number): string {
  return bytesToSize(v)
}

function formatdiskToByte(v: number): number {
  return v
}

async function getData() {
  try {
    const { data, code } = await getDiskStateByPath<SystemMonitor.DiskInfo>(props.path)
    if (code === 0)
      diskState.value = data
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
  <GenericProgress
    :progress-color="progressColor"
    :progress-rail-color="progressRailColor"
    :progress-height="5"
    :percentage="correctionNumberByCardStyle(diskState?.usedPercent || 0, cardTypeStyle)"
    :card-type-style="cardTypeStyle"
    :info-card-right-text="`${formatdiskSize(formatdiskToByte(diskState?.used || 0))}/${formatdiskSize(formatdiskToByte(diskState?.free || 0))}`"
    :info-card-left-text="diskState?.mountpoint"
    :text-color="textColor"
  />
</template>
