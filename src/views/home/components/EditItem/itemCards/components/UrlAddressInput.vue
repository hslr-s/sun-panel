<script setup lang="ts">
import { computed, ref } from 'vue'

import {
  NButton, NInput, NInputGroup, useMessage,
} from 'naive-ui'
import { verifyUrlHttp } from '../common'
import { getSiteFavicon } from '@/api/panel/itemIcon'
import { t } from '@/locales'

const props = defineProps<{
  url?: string
  iconInfo: Panel.ItemIcon | null
}>()

const emit = defineEmits<{
  (e: 'update:url', info?: string): void
  (e: 'update:iconInfo', icon?: Panel.ItemIcon | null): void
  // (e: 'updateFaviconUrl', faviconUrl?: string): void

}>()

const ms = useMessage()

// const faviconUrl = ref('')
const modelUrl = computed({
  get() {
    return props.url
  },
  set(value) {
    emit('update:url', value)
    // emit('updateFaviconUrl', faviconUrl.value)
  },
})
const modelIconInfo = computed({
  get() {
    return props.iconInfo
  },
  set(value) {
    emit('update:iconInfo', value)
  },
})

const getIconLoading = ref(false)
async function getIconByUrl(url: string) {
  getIconLoading.value = true
  try {
    const { code, data } = await getSiteFavicon<{ iconUrl: string }>(url)
    if (code === 0) {
      modelIconInfo.value = {
        itemType: 2,
        src: data.iconUrl,
      }
      // faviconUrl.value = data.iconUrl
    }
    else {
      ms.error(t('iconItem.geticonFail'))
    }
  }
  catch (error) {
    ms.error(t('iconItem.geticonFail'))
  }
  getIconLoading.value = false
}
</script>

<template>
  <div>
    <NInputGroup>
      <NInput v-model:value="modelUrl" type="text" :maxlength="1000" :placeholder="$t('iconItem.lanUrlInputPlaceholder')" />
      <NButton :disabled="!modelUrl" :loading="getIconLoading" @click="getIconByUrl(modelUrl || '')">
        {{ $t('iconItem.getIcon') }}
      </NButton>
    </NInputGroup>
    <span v-if="!verifyUrlHttp(modelUrl || '')" class="text-xs text-orange-400">
      {{ t("iconItem.urlNoHttpStartWarn") }}
    </span>
    <!-- <RoundCardModal>
      <template #title>
        <span>选择头像</span>
      </template>
    </RoundCardModal> -->
  </div>
</template>
