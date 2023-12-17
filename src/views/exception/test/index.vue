<script setup lang="ts">
import { ref } from 'vue'
import { AppStarter } from '@/views/home/components'

const componentName = ref('About')
const show = ref(true)

const changeComponent = () => {
  console.log('切换组件')
  componentName.value = 'About'
}
</script>

<template>
  <div>
    <AppStarter v-model:visible="show" />
    <button @click="changeComponent">
      Change Component
    </button>
    --
    <button @click="show = true">
      打开设置
    </button>
  </div>
</template>
