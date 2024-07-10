<script setup lang="ts">
import { onMounted, ref } from 'vue'
import {
  NAlert, NButton, NCard, NDivider, NH3,
  NInput, NPopconfirm, useMessage,
} from 'naive-ui'
import { useGlobalSitesettingStore } from './store'
import { ProAuthOverlay, ProTag, SvgIconOnline } from '@/components/common'
import { apiRespErrMsg } from '@/utils/cmn'

import {
  getCustomJsAndCssCode as getCustomJsAndCssCodeApi,
  getCustomStylePath as getCustomStylePathApi,
  saveCustomJsAndCssCode as saveCustomJsAndCssCodeApi,
} from '@/api/panel/globalSetting'
import { t } from '@/locales'

const globalSitesettingStore = useGlobalSitesettingStore()
const ms = useMessage()
const customPath = ref('')

const defaultData: GlobalSetting.GlobalSettingCustomJsAndCss = {
  jsContent: '',
  cssContent: '',
}

const siteSetting = ref<GlobalSetting.GlobalSettingCustomJsAndCss>({ ...defaultData })

async function handleSaveSiteSetting() {
  await save()
  getCustomJsAndCssCode()
}

async function save() {
  try {
    await saveCustomJsAndCssCodeApi(siteSetting.value)
    ms.success(t('common.saveSuccess'))
  }
  catch (error) {
    apiRespErrMsg(error)
    console.log(error)
  }
}

function resetSetting() {
  siteSetting.value = { ...defaultData }
  // handleSaveSiteSetting()
}

async function getCustomJsAndCssCode() {
  try {
    const { data } = await getCustomJsAndCssCodeApi<GlobalSetting.GlobalSettingCustomJsAndCss>()
    siteSetting.value = data
  }
  catch (error) {
    apiRespErrMsg(error)
    console.log(error)
  }
}

async function getCustomPath() {
  try {
    const { data } = await getCustomStylePathApi<GlobalSetting.GlobalSettingGetCustomPathReq>()
    customPath.value = data.customPath
  }
  catch (error) {
    console.log(error)
  }
}

onMounted(() => {
  getCustomPath()
  getCustomJsAndCssCode()
})
</script>

<template>
  <NCard style="border-radius:10px" class="mt-[10px]" size="small">
    <NH3>
      {{ t("apps.globalSetting.customJSAndCss") }}
      <ProTag show-icon>
        <div class="max-w-[400px]">
          {{ t("apps.globalSetting.customJSAndCssAlert", { path: customPath }) }}
        </div>
      </ProTag>
    </NH3>

    <NAlert type="error" :bordered="false" style="border-radius: 5px;">
      <template #icon>
        <SvgIconOnline icon="gg:danger" />
      </template>
      <span class="text-sm">
        {{ t("apps.globalSetting.customDangerAlert") }}
      </span>
    </NAlert>
    <ProAuthOverlay :show="globalSitesettingStore.proIsExpired">
      <div class="text-slate-500 my-[5px] font-bold">
        JS ({{ `${customPath}/index.js` }})
      </div>
      <div class="flex items-center mt-[5px]">
        <NInput v-model:value="siteSetting.jsContent" type="textarea" clearable />
      </div>
      <NDivider />
      <div class="text-slate-500 my-[5px] font-bold">
        CSS ({{ `${customPath}/index.css` }})
      </div>
      <div class="flex items-center mt-[5px]">
        <NInput v-model:value="siteSetting.cssContent" type="textarea" clearable />
      </div>

      <div class="mt-5 flex justify-end">
        <NPopconfirm @positive-click="resetSetting">
          <template #trigger>
            <NButton size="small" type="error">
              {{ t('common.reset') }}
            </NButton>
          </template>
          {{ $t('apps.globalSetting.resetWarnText') }}
        </NPopconfirm>

        <NButton size="small" type="success" class="ml-[10px]" style="margin-left: 5px;" @click="handleSaveSiteSetting">
          {{ t('common.save') }}
        </NButton>
      </div>
    </ProAuthOverlay>
  </NCard>
</template>
