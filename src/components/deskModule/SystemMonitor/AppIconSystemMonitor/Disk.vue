<script setup lang="ts">
import { onMounted, onUnmounted, ref } from 'vue'
import { NProgress } from 'naive-ui'
import type { ProgressStyle } from '../typings'
import { PanelPanelConfigStyleEnum } from '@/enums'
import { bytesToSize } from '@/utils/cmn'
import { getDiskStateByPath } from '@/api/system/systemMonitor'

interface Prop {
  cardTypeStyle: PanelPanelConfigStyleEnum
  progressStyle: ProgressStyle
  refreshInterval: number
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
    <div v-if="cardTypeStyle === PanelPanelConfigStyleEnum.info">
      <div class="mb-1">
        <span>
          {{ diskState?.mountpoint }}
        </span>
        <span class="float-right">
          {{ formatdiskSize(formatdiskToByte(diskState?.total || 0)) }}/{{ formatdiskSize(formatdiskToByte(diskState?.free || 0)) }}
        </span>
      </div>
      <NProgress
        :color="progressStyle.color"
        :rail-color="progressStyle.railColor"
        :height="progressStyle.height"
        type="line"
        :percentage="diskState?.usedPercent"
        :show-indicator="false"
        :stroke-width="15"
        style="width: 150px;"
      />
    </div>
    <div v-else>
      <div class="flex justify-center h-full w-full mt-3">
        <NProgress
          :color="progressStyle.color"
          :rail-color="progressStyle.railColor"
          type="dashboard"
          :percentage="correctionNumber(diskState?.usedPercent || 0)" :stroke-width="15"
          style="width: 50px;"
        >
          <div class="text-white" style="font-size: 10px;">
            {{ correctionNumber(diskState?.usedPercent || 0, 1) }}%
          </div>
        </NProgress>
      </div>
    </div>
  </div>
</template>
