<script setup lang="ts">
import { onMounted, onUnmounted, ref } from 'vue'
import { NProgress } from 'naive-ui'
import { getAll } from '@/api/system/systemMonitor'
import { SvgIcon } from '@/components/common'

interface ProgressStyle {
  color: string
  railColor: string
}

let timer: NodeJS.Timer
const systemMonitorData = ref<SystemMonitor.GetAllRes | null>(null)
const progressStyle = ref<ProgressStyle>({
  color: 'white',
  railColor: 'rgba(0, 0, 0, 0.5)',
})
const svgStyle = {
  width: '25px',
  height: '25px',
}

async function getData() {
  try {
    const { data, code } = await getAll<SystemMonitor.GetAllRes>()
    if (code === 0)
      systemMonitorData.value = data
  }
  catch (error) {

  }
}

function correctionNumber(v: number, keepNum = 2): number {
  return v === 0 ? 0 : Number(v.toFixed(keepNum))
}

function formatMemoryToGB(v: number): number {
  return correctionNumber(v / 1024 / 1024 / 1024)
}

function formatdiskToGB(v: number): number {
  return correctionNumber(v / 1024)
}
onMounted(() => {
  getData()
  timer = setInterval(() => {
    getData()
  }, 5000)
})

onUnmounted(() => {
  clearInterval(timer)
})
</script>

<template>
  <div class="w-full p-5 px-5 bg-[#2a2a2a6b] rounded-2xl system-monitor flex items-center px-2 text-white overflow-auto">
    <!-- <div class="flex flex-col items-center justify-center ">
      <div>
        <NProgress type="dashboard" :percentage="correctionNumber(systemMonitorData?.cpuInfo.usages[0] || 0)" :stroke-width="15" style="width: 50px;">
          <div class="text-white text-xs">
            {{ correctionNumber(systemMonitorData?.cpuInfo.usages[0] || 0, 1) }}%
          </div>
        </NProgress>
      </div>
      <span>
        CPU
      </span>
    </div> -->
    <div class="text-xs mr-10 flex justify-center items-center">
      <div class="mr-2">
        <SvgIcon icon="solar-cpu-bold" :style="svgStyle" />
      </div>
      <div>
        <div class="mb-1">
          <span>
            CPU
          </span>
          <span class="float-right">
            {{ correctionNumber(systemMonitorData?.cpuInfo.usages[0] || 0) }}%
          </span>
        </div>
        <NProgress
          type="line"
          :color="progressStyle.color"
          :rail-color="progressStyle.railColor"
          :percentage="correctionNumber(systemMonitorData?.cpuInfo.usages[0] || 0)"
          :show-indicator="false"
          :stroke-width="15"
          style="width: 150px;"
        />
      </div>
    </div>

    <div class="text-xs mr-10 flex justify-center items-center">
      <div class="mr-2">
        <SvgIcon icon="material-symbols-memory-alt-rounded" :style="svgStyle" />
      </div>
      <div>
        <div class="mb-1">
          <span>
            内存
          </span>
          <span class="float-right">
            {{ formatMemoryToGB(systemMonitorData?.memoryInfo.total || 0) }} GB/{{ correctionNumber(formatMemoryToGB(systemMonitorData?.memoryInfo.total || 0) - formatMemoryToGB(systemMonitorData?.memoryInfo.free || 0)) }} GB
          </span>
        </div>
        <NProgress
          type="line"
          :color="progressStyle.color"
          :rail-color="progressStyle.railColor"
          :percentage="systemMonitorData?.memoryInfo.usedPercent"
          :show-indicator="false"
          :stroke-width="15" style="width: 150px;"
        />
      </div>
    </div>

    <!-- <div class="text-xs mr-2">
      <div class="mb-1">
        <span>
          网络:{{ systemMonitorData?.netIOCountersInfo[0].name }}
        </span>
      </div>
      <div>
        <span class="float-right">
          上行：{{ netIOToKB(systemMonitorData?.netIOCountersInfo[0].bytesSent || 0) }}
          下行：{{ netIOToKB(systemMonitorData?.netIOCountersInfo[0].bytesRecv || 0) }}
        </span>
      </div>
    </div> -->

    <!-- 磁盘信息 -->
    <div v-for=" item, index in systemMonitorData?.diskInfo" :key="index">
      <div class="text-xs mr-10 flex justify-center items-center">
        <div class="mr-2">
          <SvgIcon icon="clarity-hard-disk-solid" :style="svgStyle" />
        </div>
        <div>
          <div class="mb-1">
            <span>
              {{ item.mountpoint }}
            </span>
            <span class="float-right">
              {{ formatdiskToGB(item.total || 0) }} GB/{{ formatdiskToGB(item.used || 0) }} GB
            </span>
          </div>
          <NProgress
            :color="progressStyle.color"
            :rail-color="progressStyle.railColor"
            type="line"
            :percentage="item.usedPercent"
            :show-indicator="false"
            :stroke-width="15"
            style="width: 150px;"
          />
        </div>
      </div>
    </div>
  </div>
</template>
