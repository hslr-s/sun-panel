<script setup lang="ts">
import { onBeforeUnmount, onMounted, ref } from 'vue'

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

  const daysOfWeek = ['日', '一', '二', '三', '四', '五', '六']
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
  <div class="w-full text-center">
    <span class="text-2xl sm:text-2xl md:text-3xl font-[600]">
      {{ currentDate.time }}
    </span>
    <div class="hidden sm:hidden md:block">
      <span>
        {{ currentDate.date }}
      </span>
      <span>
        星期{{ currentDate.week }}
      </span>
    </div>
  </div>
</template>
