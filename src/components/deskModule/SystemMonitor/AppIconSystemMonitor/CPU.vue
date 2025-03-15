<script setup lang="ts">
import { onMounted, onUnmounted, ref } from 'vue'
import { NGrid, NGridItem, NPopover, NProgress } from 'naive-ui'
import GenericProgress from '../components/GenericProgress/index.vue'
import { allCPUUsages, correctionNumber, correctionNumberByCardStyle } from './common'
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
let timer: any
const cpuState = ref<SystemMonitor.CPUInfo | null>(null)

async function getData() {
  try {
    const { data, code } = await getCpuState<SystemMonitor.CPUInfo>()
    if (code === 0)
      cpuState.value = data
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
        :percentage="correctionNumberByCardStyle(allCPUUsages(cpuState), cardTypeStyle)"
        :card-type-style="cardTypeStyle" :info-card-right-text="`${correctionNumber(allCPUUsages(cpuState))}%`"
        info-card-left-text="CPU" :text-color="textColor" style="width: 100%;" class="cursor-pointer"
      />
    </template>

    <div class="mt-2">
      <div class="overflow-auto max-h-[200px] px-2">
        <NGrid x-gap="12" y-gap="15" cols="4 400:5 600:6" responsive="screen">
          <NGridItem v-for="item, index in cpuState?.usages" :key="index">
            <div class="flex flex-col justify-center items-center">
              <NProgress type="dashboard" gap-position="bottom" :percentage="item" style="width: 65px;">
                <span class="text-xs"> {{ correctionNumber(item, 0) }}%</span>
              </NProgress>
              <span class="text-xs"># {{ index + 1 }}</span>
            </div>
          </NGridItem>
        </NGrid>
      </div>
    </div>
  </NPopover>
</template>
