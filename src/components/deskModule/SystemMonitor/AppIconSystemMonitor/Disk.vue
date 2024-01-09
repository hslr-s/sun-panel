<script setup lang="ts">
import { onMounted, onUnmounted, ref } from 'vue'
import GenericProgress from '../components/GenericProgress/index.vue'
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

function correctionNumber(v: number, keepNum = 2): number {
  return v === 0 ? 0 : Number(v.toFixed(keepNum))
}

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
  <div class="w-full">
    <GenericProgress
      :progress-color="progressColor"
      :progress-rail-color="progressRailColor"
      :progress-height="5"
      :percentage="correctionNumber(diskState?.usedPercent || 0)"
      :card-type-style="cardTypeStyle"
      :info-card-right-text="`${formatdiskSize(formatdiskToByte(diskState?.total || 0))}/${formatdiskSize(formatdiskToByte(diskState?.free || 0))}`"
      :info-card-left-text="diskState?.mountpoint"
      :text-color="textColor"
    />
  </div>
</template>
