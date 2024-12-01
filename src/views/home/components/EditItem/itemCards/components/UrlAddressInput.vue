<script setup lang="ts">
import { computed, ref } from 'vue'

import {
  NButton, NImage, NInput, NInputGroup, NSpin, useMessage,
} from 'naive-ui'
import { verifyUrlHttp } from '../common'
import { downloadFavicon, getSiteFaviconList } from '@/api/panel/itemIcon'
import { t } from '@/locales'
import { RoundCardModal } from '@/components/common'
import { apiRespErrMsg } from '@/utils/cmn/apiMessage'

const props = defineProps<{
  url?: string
  iconInfo: Panel.ItemIcon | null
  placeholder?: string
}>()

const emit = defineEmits<{
  (e: 'update:url', info?: string): void
  (e: 'update:iconInfo', icon?: Panel.ItemIcon | null): void
}>()

const ms = useMessage()

const selectIconModalShow = ref(false)
const faviconList = ref<string[]>([])
const downloadingFavicon = ref(false)

const modelUrl = computed({
  get() {
    return props.url
  },
  set(value) {
    emit('update:url', value)
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
  await getSiteFaviconList < Common.ListResponse<string[]>>(url).then(({ data }) => {
    // 当有多个图标的时候，弹出选择框，让用户选择
    if (data.list.length > 1) {
      faviconList.value = data.list
      selectIconModalShow.value = true
    }
    else {
      setIconInfo(data.list[0])
    }
  }).catch((res) => {
    if (res.code === -1) {
      ms.error(t('iconItem.geticonFail'))
      return
    }
    apiRespErrMsg(res)
  })

  getIconLoading.value = false
}

async function setIconInfo(url: string) {
  downloadingFavicon.value = true
  await downloadFavicon<{ iconUrl: string }>(url).then(({ data }) => {
    modelIconInfo.value = {
      itemType: 2,
      src: data.iconUrl,
    }
    selectIconModalShow.value = false
  }).catch((res) => {
    if (res.code === -1) {
      ms.error(t('iconItem.geticonFail'))
      return
    }
    apiRespErrMsg(res)
  })
  downloadingFavicon.value = false
}
</script>

<template>
  <div>
    <NInputGroup>
      <NInput v-model:value="modelUrl" type="text" :maxlength="1000" :placeholder="placeholder" />
      <NButton :disabled="!modelUrl" :loading="getIconLoading" @click="getIconByUrl(modelUrl || '')">
        {{ $t('iconItem.getIcon') }}
      </NButton>
    </NInputGroup>
    <span v-if="!verifyUrlHttp(modelUrl || '')" class="text-xs text-orange-400">
      {{ t("iconItem.urlNoHttpStartWarn") }}
    </span>
    <RoundCardModal v-model:show="selectIconModalShow" size="small" style="width: 500px;" :title="t('iconItem.selectOneIcon')">
      <NSpin :show="downloadingFavicon">
        <span
          v-for="item, index in faviconList" :key="index"
          style="margin: 5px;cursor: pointer;padding: 2px;"
          class="rounded-xl border inline-block w-[80px] h-[80px] "
          @click="setIconInfo(item)"
        >
          <NImage
            :src="item" style="width: 100%;height: 100%;border-radius: 5px;" preview-disabled
            class="hover:brightness-105 transition-all hover:scale-95"
          />
        </span>
      </NSpin>
    </RoundCardModal>
  </div>
</template>
