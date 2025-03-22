<script setup lang="ts">
import { onMounted, ref } from 'vue'
import axios from 'axios'
import type { UploadFileInfo } from 'naive-ui'
import {
  NButton, NButtonGroup, NCard, NDivider, NH3, NImage,
  NInput, NPopconfirm, NSwitch, NUpload, NUploadDragger, useMessage,
} from 'naive-ui'
import Custom from './Custom.vue'
import { useGlobalSitesettingStore } from './store'
import { ProAuthOverlay, ProTag, STip } from '@/components/common'
import { useAuthStore } from '@/store'
import { getSiteStting as getSiteSttingApi, saveSiteStting as saveSiteSettingApi } from '@/api/panel/globalSetting'
import { t } from '@/locales'
import { apiRespErrMsg } from '@/utils/cmn'
import { defaultSetting } from '@/store/modules/app/helper'
import { FileType } from '@/enums'

const authStore = useAuthStore()
const globalSitesettingStore = useGlobalSitesettingStore()
const ms = useMessage()
const icoThumbnailUrl = ref('')
const icoImageFile = ref<File | null>(null)
const icoIsNeedUpload = ref(false)
const defaultSiteSetting: GlobalSetting.SiteSetting = defaultSetting().siteSetting

const siteSetting = ref<GlobalSetting.SiteSetting>({ ...defaultSiteSetting })

function handleUploadBackgroundFinish({
  file,
  event,
}: {
  file: UploadFileInfo
  event?: ProgressEvent
}) {
  const res = JSON.parse((event?.target as XMLHttpRequest).response)
  siteSetting.value.loginBackgroundUrl = res.data.imageUrl
  return file
}

async function handleSaveSiteSetting() {
  let icoImage = icoThumbnailUrl.value
  if (icoIsNeedUpload.value) {
    try {
      icoImage = await handleUploadIco()
    }
    catch (error) {
      return
    }
  }

  siteSetting.value.favicoUrl = icoImage
  icoIsNeedUpload.value = false

  await saveSiteSetting()
  getSiteSetting()
}

async function saveSiteSetting() {
  try {
    await saveSiteSettingApi(siteSetting.value)
    ms.success(t('common.saveSuccess'))
  }
  catch (error) {
    ms.error(t('common.saveFail'))
    console.error(error)
  }
}

function handleFileIcoUploadChange(event: Event) {
  const target = event.target as HTMLInputElement
  // console.log(target)

  if (target.files && target.files[0]) {
    // console.log(target.files)
    const imageFile = target.files[0]
    const imageUrl = URL.createObjectURL(imageFile)
    // console.log(imageUrl)
    icoThumbnailUrl.value = imageUrl
    icoImageFile.value = imageFile
  }
}

async function handleUploadIco(): Promise<string> {
  let imgurl = ''
  if (icoImageFile.value) {
    const formData = new FormData()
    formData.append('imgfile', icoImageFile.value)
    formData.append('type', String(FileType.CardIcon))
    await axios.post('/api/panel/globalSetting/uploadFaviconImage', formData, {
      headers: {
        'Content-Type': 'multipart/form-data',
        'token': authStore.token as string,
      },
    })
      .then((response) => {
        // console.log('上传', response)
        if (response.data.code !== 0) {
          if (response.data.code === 1301) {
            ms.error(`${t('apps.globalSetting.favicon')}: ${t('apiErrorCode.1301')}`)
          }

          throw new Error(`upload error API:${response.data.msg}`)
        }
        imgurl = response.data.data.imageUrl
      })
      .catch((error) => {
        console.error('upload error', error)
        throw error
      })
  }

  return imgurl
}

function handleSelectIcoFile() {
  const fileInput = document.getElementById('fileIcoUpload') as HTMLInputElement
  fileInput.click()
  icoIsNeedUpload.value = true
}

function handleDeleteIcoFile() {
  icoImageFile.value = null
  icoThumbnailUrl.value = ''
  icoIsNeedUpload.value = false
}

function resetSiteSetting() {
  siteSetting.value = { ...defaultSiteSetting }
  icoThumbnailUrl.value = defaultSiteSetting.favicoUrl
  handleSaveSiteSetting()
}

// 获取站点设置
async function getSiteSetting() {
  try {
    const { data } = await getSiteSttingApi<GlobalSetting.SiteSetting>()
    siteSetting.value = data
    icoThumbnailUrl.value = data.favicoUrl
  }
  catch (error) {
    apiRespErrMsg(error)
  }
}

onMounted(() => {
  getSiteSetting()
})
</script>

<template>
  <div class="bg-slate-200 dark:bg-zinc-900 rounded-[10px] p-[8px] overflow-auto">
    <NCard style="border-radius:10px" size="small">
      <NH3>
        {{ t("apps.globalSetting.personalise") }}
      </NH3>

      <div>
        <div class="text-slate-500 mb-[5px] font-bold">
          {{ $t('apps.globalSetting.title') }}
          <ProTag />
        </div>
        <ProAuthOverlay :show="globalSitesettingStore.proIsExpired">
          <div class="flex items-center mt-[5px]">
            <NInput
              v-model:value="siteSetting.title" type="text" show-count :maxlength="20"
              :placeholder="t('common.inputPlaceholder')"
            />
          </div>
        </ProAuthOverlay>
      </div>
      <NDivider />
      <div>
        <div class="text-slate-500 mb-[5px] font-bold">
          {{ $t('apps.globalSetting.favicon') }}
          <ProTag />
        </div>

        <ProAuthOverlay :show="globalSitesettingStore.proIsExpired">
          <div>
            <NImage
              width="100" style="width: 80px; height: 80px;border: 1px solid #e9e9e9;" :src="icoThumbnailUrl"
              preview-disabled
            />
            <br>
            <NButtonGroup size="small">
              <NButton for="fileIcoUpload" @click="handleSelectIcoFile">
                {{ t('apps.globalSetting.selectIcoImage') }}
              </NButton>
              <NButton for="fileIcoUpload" :disabled="icoThumbnailUrl === ''" @click="handleDeleteIcoFile">
                {{ t('apps.globalSetting.clearIcoImage') }}
              </NButton>
            </NButtonGroup>
            <input
              id="fileIcoUpload" type="file" accept=".ico" style="display: none;"
              @change="handleFileIcoUploadChange"
            >
          </div>
        </ProAuthOverlay>
      </div>

      <NDivider />

      <div>
        <div class="text-slate-500 mb-[5px] font-bold">
          {{ $t('apps.globalSetting.isEnableLoginCaptcha') }}
          <ProTag />
        </div>
        <ProAuthOverlay :show="globalSitesettingStore.proIsExpired">
          <div class="flex items-center mt-[5px]">
            <NSwitch v-model:value="siteSetting.isEnableLoginCaptcha" />
          </div>
        </ProAuthOverlay>
      </div>

      <NDivider />

      <div>
        <div class="text-slate-500 mb-[5px] font-bold">
          {{ $t('apps.globalSetting.loginBackgroundImage') }}
        </div>
        <NUpload
          action="/api/file/uploadImg" :show-file-list="false" name="imgfile" :headers="{
            token: authStore.token as string,
          }" :data="{ type: FileType.Wallpaper }" :directory-dnd="true" @finish="handleUploadBackgroundFinish"
        >
          <NUploadDragger style="width: 100%;">
            <div
              class="h-[200px] w-full border bg-slate-100 flex justify-center items-center cursor-pointer rounded-[10px]"
              :style="{ background: `url(${siteSetting.loginBackgroundUrl}) no-repeat`, backgroundSize: 'cover' }"
            >
              <div class="text-shadow text-white">
                {{ $t('apps.baseSettings.uploadOrDragText') }}
              </div>
            </div>
          </NUploadDragger>
        </NUpload>

        <div class="flex items-center mt-[5px]">
          <span class="mr-[10px]">{{ $t('apps.globalSetting.customImageURL') }}</span>
        </div>
        <NInput v-model:value="siteSetting.loginBackgroundUrl" type="text" size="small" clearable />
      </div>

      <NDivider />

      <div>
        <div class="text-slate-500 mb-[5px] font-bold flex items-center">
          {{ $t('apps.globalSetting.isEnableItemCardCache') }}
          <STip class="ml-2">
            <div class="max-w-[300px]">
              {{ t('apps.globalSetting.isEnableItemCardCacheTip') }}
            </div>
          </STip>
        </div>

        <div class="flex items-center mt-[5px]">
          <!-- 为了适配后端此处“isDisableItemCardCache”取反，因为前端的逻辑是“isEnableItemCardCache”的取值 -->
          <NSwitch v-model:value="siteSetting.isDisableItemCardCache" :checked-value="false" :unchecked-value="true" />
        </div>
      </div>

      <div class="mt-5 flex justify-end">
        <NPopconfirm @positive-click="resetSiteSetting">
          <template #trigger>
            <NButton size="small" type="error">
              {{ $t('common.reset') }}
            </NButton>
          </template>
          {{ $t('apps.globalSetting.resetWarnText') }}
        </NPopconfirm>

        <NButton size="small" type="success" class="ml-[10px]" style="margin-left: 5px;" @click="handleSaveSiteSetting">
          {{ $t('common.save') }}
        </NButton>
      </div>
    </NCard>

    <Custom />
  </div>
</template>

<style scoped>
.text-shadow{
  text-shadow: 0px 0px 5px gray;
}
</style>
