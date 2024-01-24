<script setup lang="ts">
import { onBeforeUnmount, onMounted, ref } from 'vue'
import { t } from '@/locales'

const props = defineProps<{
  hideSecond?: boolean
}>()

interface CurrentDate {
  time: string
  date: string
  week: string
}

const currentDate = ref<CurrentDate>({
  time: '--:--',
  date: '------',
  week: '--',
})

function updateCurrentDate() {
  const now = new Date()
  const hours = String(now.getHours()).padStart(2, '0')
  const minutes = String(now.getMinutes()).padStart(2, '0')

  if (!props.hideSecond) {
    const seconds = String(now.getSeconds()).padStart(2, '0')
    currentDate.value.time = `${hours}:${minutes}:${seconds}`
  }
  else {
    currentDate.value.time = `${hours}:${minutes}`
  }

  // 获取当前的日期
  const day = now.getDate()
  const month = now.getMonth() + 1 // 月份从0开始，所以要加1
  // const year = now.getFullYear()

  const daysOfWeek = [
    t('deskModule.clock.sun'),
    t('deskModule.clock.mon'),
    t('deskModule.clock.tue'),
    t('deskModule.clock.wed'),
    t('deskModule.clock.thu'),
    t('deskModule.clock.fri'),
    t('deskModule.clock.sat'),
  ]
  // const daysOfWeek = ['Sunday', 'Monday', 'Tuesday', 'Wednesday', 'Thursday', 'Friday', 'Saturday']
  currentDate.value.week = daysOfWeek[now.getDay()]
  currentDate.value.date = `${month}-${day}`
}

const updateClock = () => {
  updateCurrentDate()
}

const intervalId = setInterval(updateClock, 1000)

onMounted(() => {
  updateClock()
  updateCurrentDate()
})

onBeforeUnmount(() => {
  clearInterval(intervalId)
})
</script>

<template>
  <div class="clock w-full text-center">
    <span class="clock-time text-2xl sm:text-2xl md:text-3xl font-[600]">
      {{ currentDate.time }}
    </span>
    <div class="hidden sm:hidden md:block">
      <span class="clock-date mr-1">
        {{ currentDate.date }}
      </span>
      <span class="clock-week">
        {{ currentDate.week }}
      </span>
    </div>
  </div>
</template>
