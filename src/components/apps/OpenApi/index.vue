<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { NAlert, NButton, NCard, NDivider, NH4, NTag, useDialog, useMessage } from 'naive-ui'
import { t } from '@/locales'
import { useAppStore } from '@/store'
import {
  deleteToken as deleteTokenApi,
  getTokenStatusInfo as getTokenStatusInfoApi,
  refreshToken as refreshTokenApi,
} from '@/api/system/openApi'

const ms = useMessage()
const dialog = useDialog()
const appStore = useAppStore()
const openApiLink = computed(() => {
  const url = 'http://links.sun-panel.top/openApiDoc_'
  if (appStore.language === 'zh-CN')
    return `${url}zh_cn`

  return `${url}en`
})
const tokenInfo = ref<OpenApi.GetTokenStatusInfoResp>({
  status: false,
  token: '',
})

const urlPath = '/openapi/v1'

function getCurrentDomain() {
  return `${location.protocol}//${location.hostname}${location.port === '' ? '' : (`:${location.port}`)}`
}

async function getTokenStatusInfo() {
  try {
    const { data } = await getTokenStatusInfoApi<OpenApi.GetTokenStatusInfoResp>()
    tokenInfo.value = data
  }
  catch (error) {

  }
}

function handleRefreshToken() {
  dialog.warning({
    title: t('common.warning'),
    content: t('apps.openApi.updateTokenDialogText'),
    positiveText: t('common.confirm'),
    negativeText: t('common.cancel'),
    onPositiveClick: () => {
      refreshToken()
    },
  })
}

function handleDeleteToken() {
  dialog.warning({
    title: t('common.warning'),
    content: t('apps.openApi.updateTokenDialogText'),
    positiveText: t('common.confirm'),
    negativeText: t('common.cancel'),
    onPositiveClick: () => {
      deleteToken()
    },
  })
}

async function refreshToken() {
  try {
    await refreshTokenApi<OpenApi.TokenInfo>()
    getTokenStatusInfo()
  }
  catch (error) {
    ms.error(t('common.failed'))
  }
}

async function deleteToken() {
  try {
    await deleteTokenApi()
    getTokenStatusInfo()
  }
  catch (error) {
    ms.error(t('common.failed'))
  }
}

onMounted(() => {
  getTokenStatusInfo()
})
</script>

<template>
  <div class="bg-slate-200 dark:bg-zinc-900 p-2 h-full">
    <NCard style="border-radius:10px" size="small">
      <NAlert type="info" :bordered="false">
        {{ t("apps.openApi.alertText") }}<a :href="openApiLink" target="_blank" class="text-[blue]">{{
          t("apps.openApi.docTitle") }}</a>
      </NAlert>

      <div class="mt-5">
        <div>
          <NH4>
            {{ t('apps.openApi.state') }}
          </NH4>
          <NTag v-if="tokenInfo.status" type="success" :bordered="false">
            {{ t('apps.openApi.enable') }}
          </NTag>
          <NTag v-else :bordered="false">
            {{ t('apps.openApi.deactivate') }}
          </NTag>
          <NDivider />
        </div>

        <div v-if="tokenInfo.status">
          <NH4>
            {{ t('apps.openApi.apiUrl') }} <span class="text-xs text-slate-400"> {{ t('apps.openApi.autoDomainWarning')
            }}</span>
          </NH4>
          {{ getCurrentDomain() }}{{ urlPath }}
          <NDivider />
        </div>

        <div v-if="tokenInfo.status">
          <NH4>
            Token
          </NH4>
          <NTag :bordered="false">
            {{ tokenInfo.token }}
          </NTag>
          <NDivider />
        </div>
      </div>

      <div class="flex justify-center">
        <template v-if="tokenInfo.status">
          <NButton v-if="tokenInfo.status" type="info" strong secondary @click="handleRefreshToken">
            {{ t('apps.openApi.refreshToken') }}
          </NButton>

          <NButton type="error" strong secondary style="margin-left: 4px;" @click="handleDeleteToken">
            {{ t('apps.openApi.stopOpenApi') }}
          </NButton>
        </template>

        <NButton v-else type="success" @click="refreshToken">
          {{ t('apps.openApi.enable') }}
        </NButton>
      </div>
    </NCard>
  </div>
</template>
