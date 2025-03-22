<script setup lang="ts">
import { onMounted, onUnmounted, ref } from 'vue'
import { NPopover } from 'naive-ui'
import GenericProgress from '../components/GenericProgress/index.vue'
import { correctionNumber, correctionNumberByCardStyle } from './common'
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
let timer: any // NodeJS.Timer
const memoryState = ref<SystemMonitor.MemoryInfo | null>(null)

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
        :percentage="correctionNumberByCardStyle(memoryState?.usedPercent || 0, cardTypeStyle)"
        :card-type-style="cardTypeStyle"
        :info-card-right-text="`${formatMemorySize(memoryState?.used || 0)}/${formatMemorySize(memoryState?.total || 0)}`"
        info-card-left-text="RAM" :text-color="textColor" class="cursor-pointer"
      />
    </template>

    <div>
      <!-- <div class="my-2 text-lg">
        <span class="font-bold">
          {{ $t('deskModule.systemMonitor.diskMountPoint') }}:
        </span>

        {{ diskState?.mountpoint }}
      </div> -->

      <div class="overflow-auto max-h-[200px]">
        <div>
          <span> {{ $t('deskModule.systemMonitor.free') }}</span>
          <span class="ml-2">{{ formatMemorySize((memoryState?.total || 0) - (memoryState?.used || 0)) }}</span>
          <span class="ml-2">({{ 100 - correctionNumber(memoryState?.usedPercent || 0) }}%)</span>
        </div>
        <div>
          <span> {{ $t('deskModule.systemMonitor.used') }}</span>
          <span class="ml-2">{{ formatMemorySize(memoryState?.used || 0) }}</span>
          <span class="ml-2">({{ correctionNumber(memoryState?.usedPercent || 0) }}%)</span>
        </div>
        <div>
          <span>{{ $t('deskModule.systemMonitor.total') }}</span>
          <span class="ml-2">{{ formatMemorySize(memoryState?.total || 0) }}</span>
        </div>
      </div>
    </div>
  </NPopover>
</template>
