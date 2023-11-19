<script setup lang='ts'>
import { onMounted, ref } from 'vue'
import { NSpin } from 'naive-ui'
import { fetchChatConfig } from '@/api'
import { getAboutDescription as getAboutDescriptionApi } from '@/api/openness'

interface ConfigState {
  timeoutMs?: number
  reverseProxy?: string
  apiModel?: string
  socksProxy?: string
  httpsProxy?: string
  usage?: string
}

const loading = ref(false)

const config = ref<ConfigState>()

const content = ref<string>()

async function fetchConfig() {
  try {
    loading.value = true
    const { data } = await fetchChatConfig<ConfigState>()
    config.value = data
  }
  finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchConfig()
  getAboutDescription()
})

async function getAboutDescription() {
  const { data } = await getAboutDescriptionApi<string>()
  content.value = data
}
</script>

<template>
  <NSpin :show="loading">
    <div class="p-4 space-y-4">
      <h2 class="text-xl font-bold">
        {{ $t('common.appName') }}
      </h2>
      <div>
        <span v-html="content" />
      </div>
    </div>
  </NSpin>
</template>
