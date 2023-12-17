<script setup lang="ts">
import { ref } from 'vue'

const jsonData = ref<string | null>(null)

const handleFileChange = (event: Event) => {
  const input = event.target as HTMLInputElement
  const file = input.files?.[0]

  if (file) {
    const reader = new FileReader()
    reader.onload = () => {
      if (reader.result)
        jsonData.value = reader.result as string
    }
    reader.readAsText(file)
  }
}

const exportJson = () => {
  if (jsonData.value) {
    const blob = new Blob([jsonData.value], { type: 'application/json' })
    const link = document.createElement('a')
    link.href = URL.createObjectURL(blob)
    link.download = 'exported_data.json'
    link.click()
  }
}
</script>

<template>
  <div>
    <input type="file" @change="handleFileChange">
    <button @click="exportJson">
      导出JSON
    </button>

    <div v-if="jsonData">
      <h2>JSON 数据</h2>
      <pre>{{ jsonData }}</pre>
    </div>
  </div>
</template>
