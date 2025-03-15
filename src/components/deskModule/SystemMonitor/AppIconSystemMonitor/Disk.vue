<script setup lang="ts">
import { onMounted, onUnmounted, ref } from 'vue'
import { NPopover } from 'naive-ui'
import GenericProgress from '../components/GenericProgress/index.vue'
import { correctionNumber, correctionNumberByCardStyle } from './common'
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
  title: string
}

const props = defineProps<Prop>()
let timer: any // NodeJS.Timer
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
    console.error(error)
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
  <NPopover trigger="click">
    <template #trigger>
      <GenericProgress
        :progress-color="progressColor" :progress-rail-color="progressRailColor" :progress-height="5"
        :percentage="correctionNumberByCardStyle(diskState?.usedPercent || 0, cardTypeStyle)"
        :card-type-style="cardTypeStyle"
        :info-card-right-text="`${formatdiskSize(formatdiskToByte(diskState?.used || 0))}/${formatdiskSize(formatdiskToByte(diskState?.total || 0))}`"
        :info-card-left-text="title || diskState?.mountpoint" :text-color="textColor" class="cursor-pointer"
      />
    </template>

    <div>
      <div class="overflow-auto max-h-[200px]">
        <div>
          <span> {{ $t('deskModule.systemMonitor.free') }}</span>
          <span class="ml-2">{{ formatdiskSize(formatdiskToByte(diskState?.free || 0)) }}</span>
          <span class="ml-2">({{ 100 - correctionNumber(diskState?.usedPercent || 0) }}%)</span>
        </div>
        <div>
          <span> {{ $t('deskModule.systemMonitor.used') }}</span>
          <span class="ml-2">{{ formatdiskSize(formatdiskToByte(diskState?.used || 0)) }}</span>
          <span class="ml-2">({{ correctionNumber(diskState?.usedPercent || 0) }}%)</span>
        </div>
        <div>
          <span>{{ $t('deskModule.systemMonitor.total') }}</span>
          <span class="ml-2">{{ formatdiskSize(formatdiskToByte(diskState?.total || 0)) }}</span>
        </div>
      </div>
    </div>
  </NPopover>
</template>
