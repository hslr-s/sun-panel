<script setup lang="ts">
import { NButton } from 'naive-ui'
import { clearAllCache } from '@/utils/pwaClear'
// const a = ref()

function handleClear() {
  clearAllCache()
}
</script>

<template>
  <div>
    <NButton @click="handleClear">
      Clear Cache ALL
    </NButton>
  </div>
</template>
