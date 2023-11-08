<script setup lang='ts'>
import { onMounted, ref } from 'vue'
import { NButton, NInput, NInputGroup, NSpin } from 'naive-ui'
import { useClipboard } from '@vueuse/core'
import { getReferralCode as getReferralCodeApi } from '@/api/system/user'
// import { copyText } from '@/utils/format'

const loading = ref(false)

const referralCode = ref<string>('')
const url = ref<string>('')
const copyButtonText = ref<string>('')

async function getReferralCode() {
  try {
    loading.value = true
    const { data } = await getReferralCodeApi<User.GetReferralCodeResponse>()
    referralCode.value = data.referralCode
    url.value = `${getCurrentDomain()}/#/register?referralCode=${referralCode.value}`
  }
  finally {
    loading.value = false
  }
}

function getCurrentDomain() {
  return `${location.protocol}//${location.hostname}${location.port === '' ? '' : (`:${location.port}`)}`
}

function handleCopy() {
  const { copy } = useClipboard()
  //   copyText({ text: urlCopy })
  copy(url.value)
  copyButtonText.value = '复制完成！'
  setInterval(() => {
    copyButtonText.value = '复制链接'
  }, 3000)
}

onMounted(() => {
  getReferralCode()
  copyButtonText.value = '复制链接'
})
</script>

<template>
  <NSpin :show="loading">
    <div class="p-4 space-y-4">
      <h2 class="text-xl ">
        您的专属推荐链接
      </h2>
      <div>
        <NInputGroup>
          <NInput v-model:value="url" readonly type="text" />
          <NButton type="primary" ghost @click="handleCopy">
            {{ copyButtonText }}
          </NButton>
        </NInputGroup>
      </div>
    </div>
  </NSpin>
</template>
