<script setup lang="ts">
import { ref, watch } from 'vue'
import type { UploadFileInfo } from 'naive-ui'
import {
  NButton, NCard, NColorPicker, NFlex, NGrid,
  NGridItem, NInput, NInputGroup,
  NInputNumber, NPopconfirm,
  NSelect, NSlider, NSwitch, NTag, NUpload,
  NUploadDragger,
  useDialog, useMessage,
} from 'naive-ui'
import { useAuthStore, usePanelState } from '@/store'
import { defaultStatePanelConfig } from '@/store/modules/panel/helper'
import { set as setUserConfig } from '@/api/panel/userConfig'
// import { PanelPanelConfigStyleEnum } from '@/enums/panel'
import { t } from '@/locales'
import { apiRespErrMsg } from '@/utils/cmn/apiMessage'
import { FileType, PanelPanelConfigStyleEnum } from '@/enums'
import { STip } from '@/components/common'

const authStore = useAuthStore()
const panelState = usePanelState()
const ms = useMessage()
const dialog = useDialog()
const showWallpaperInput = ref(false)

const isSaveing = ref(false)

const iconTypeOptions = [
  {
    label: t('apps.baseSettings.detailIcon'),
    value: PanelPanelConfigStyleEnum.info,
  },
  {
    label: t('apps.baseSettings.smallIcon'),
    value: PanelPanelConfigStyleEnum.icon,
  },
]

const maxWidthUnitOption = [
  {
    label: 'px',
    value: 'px',
  },
  {
    label: '%',
    value: '%',
  },
]

watch(panelState.panelConfig, () => {
  if (!isSaveing.value) {
    isSaveing.value = true

    setTimeout(() => {
      panelState.recordState()// 本地记录
      isSaveing.value = false
      uploadCloud()
    }, 1000)
  }
})

function handleUploadBackgroundFinish({
  file,
  event,
}: {
  file: UploadFileInfo
  event?: ProgressEvent
}) {
  const res = JSON.parse((event?.target as XMLHttpRequest).response)
  if (res.code === 0)
    panelState.panelConfig.backgroundImageSrc = res.data.imageUrl

  else
    apiRespErrMsg(res)

  return file
}

function uploadCloud() {
  setUserConfig({ panel: panelState.panelConfig }).then((res) => {
    ms.success(t('apps.baseSettings.configSaved'))
  }).catch((error) => {
    ms.error(t('apps.baseSettings.configFailed', { message: error.msg }))
  })
}

function handleResetSearchBox() {
  panelState.panelConfig.searchBoxShow = defaultStatePanelConfig().searchBoxShow
  panelState.panelConfig.searchBoxSearchIcon = defaultStatePanelConfig().searchBoxSearchIcon
  panelState.panelConfig.searchBoxBackground = defaultStatePanelConfig().searchBoxBackground
  panelState.panelConfig.searchBoxBorderColor = defaultStatePanelConfig().searchBoxBorderColor
  panelState.panelConfig.searchBoxTextColor = defaultStatePanelConfig().searchBoxTextColor
}

function resetPanelConfig() {
  dialog.warning({
    title: t('common.warning'),
    content: t('apps.baseSettings.resetAllStyleWarnText'),
    positiveText: t('common.confirm'),
    negativeText: t('common.cancel'),
    onPositiveClick: () => {
      panelState.resetPanelConfig()
      uploadCloud()
    },
  })
}

// 开启自动识别网络将隐藏手动切换按钮
function handleNetModeChangeAutoChange(value: boolean) {
  panelState.panelConfig.netModeChangeButtonShow = !value
}

function handleNetModeChangeAutoMaxPingMsChange(value: number) {
  if (value < 100) {
    panelState.panelConfig.netModeChangeAutoMaxPingMs = 100
  }
}
</script>

<template>
  <div class="bg-slate-200 dark:bg-zinc-900 rounded-[10px] p-[8px] overflow-auto">
    <NCard style="border-radius:10px" size="small">
      <template #header>
        <div class="text-slate-500 font-bold">
          LOGO
        </div>
      </template>

      <div class="flex items-center mt-[5px]">
        <span class="mr-[10px]">{{ $t('common.show') }}</span>
        <NSwitch v-model:value="panelState.panelConfig.logoShow" size="small" />
      </div>

      <div class="items-center mt-[5px]">
        <div>
          {{ $t('apps.baseSettings.textContent') }}
        </div>
        <div class="flex items-center mt-[5px]">
          <NInput v-model:value="panelState.panelConfig.logoText" :disabled="!panelState.panelConfig.logoShow" type="text" show-count :maxlength="20" placeholder="请输入文字" />
        </div>
      </div>
    </NCard>

    <NCard style="border-radius:10px" class="mt-[10px]" size="small">
      <template #header>
        <div class="text-slate-500 font-bold">
          {{ $t('apps.baseSettings.clock') }}
        </div>
      </template>

      <div class="flex items-center mt-[5px]">
        <span class="mr-[10px]">{{ $t('common.show') }}</span>
        <NSwitch v-model:value="panelState.panelConfig.clockShow" size="small" />
      </div>

      <div class="flex items-center mt-[5px]">
        <span class="mr-[10px]">{{ $t('apps.baseSettings.clockSecondShow') }}</span>
        <NSwitch v-model:value="panelState.panelConfig.clockShowSecond" size="small" :disabled="!panelState.panelConfig.clockShow" />
      </div>
    </NCard>

    <NCard style="border-radius:10px" class="mt-[10px]" size="small">
      <template #header>
        <div class="text-slate-500 font-bold flex">
          <span>
            {{ $t('apps.baseSettings.searchBar') }}
          </span>

          <span class="ml-auto">
            <NPopconfirm
              @positive-click="handleResetSearchBox"
            >
              <template #trigger>
                <NButton size="tiny" quaternary type="warning">{{ t('common.reset') }}</NButton>
              </template>
              {{ $t('apps.baseSettings.resetWarnText') }}
            </NPopconfirm>
          </span>
        </div>
      </template>
      <div class="flex items-center mt-[5px]">
        <span class="mr-[10px]">{{ $t('common.show') }}</span>
        <NSwitch v-model:value="panelState.panelConfig.searchBoxShow" size="small" />
      </div>
      <template v-if="panelState.panelConfig.searchBoxShow">
        <div class="flex items-center mt-[5px]">
          <span class="mr-[10px]">{{ $t('apps.baseSettings.searchBarSearchItem') }}</span>
          <NSwitch v-model:value="panelState.panelConfig.searchBoxSearchIcon" size="small" />
        </div>

        <NFlex class="mt-[5px]" align="center">
          <div>
            {{ $t('common.backgroundColor') }}
          </div>
          <div class="flex items-center mt-[5px] w-[100px]">
            <NColorPicker
              v-model:value="panelState.panelConfig.searchBoxBackground"
              size="small"
              :modes="['hex']"
              :swatches="[
                '#00000000',
                '#000000',
                '#ffffff',
                '#18A058',
                '#2080F0',
                '#F0A020',
              ]"
            />
          </div>
        </NFlex>

        <NFlex class="mt-[5px]" align="center">
          <div>
            {{ $t('common.textColor') }}
          </div>
          <div class="flex items-center mt-[5px] w-[100px]">
            <NColorPicker
              v-model:value="panelState.panelConfig.searchBoxTextColor"
              size="small"
              :modes="['hex']"
              :swatches="[
                '#00000000',
                '#000000',
                '#ffffff',
                '#18A058',
                '#2080F0',
                '#F0A020',
              ]"
            />
          </div>
        </NFlex>

        <NFlex class="mt-[5px]" align="center">
          <div>
            {{ $t('apps.baseSettings.searchBarBorderColor') }}
          </div>
          <div class="flex items-center mt-[5px] w-[100px]">
            <NColorPicker
              v-model:value="panelState.panelConfig.searchBoxBorderColor"
              size="small"
              :modes="['hex']"
              :swatches="[
                '#00000000',
                '#000000',
                '#ffffff',
                '#18A058',
                '#2080F0',
                '#F0A020',
              ]"
            />
          </div>
        </NFlex>

        <!-- <NFlex class="mt-[5px]" align="center">
          <div>
            {{ $t('apps.baseSettings.searchBarPlaceholderColor') }}
          </div>
          <div class="flex items-center mt-[5px] w-[100px]">
            <NColorPicker
              v-model:value="panelState.panelConfig.searchBoxPlaceholderColor"
              :show-alpha="false"
              size="small"
              :modes="['hex']"
              :swatches="[
                '#000000',
                '#ffffff',
                '#18A058',
                '#2080F0',
                '#F0A020',
              ]"
            />
          </div>
        </NFlex> -->
      </template>
    </NCard>

    <NCard style="border-radius:10px" class="mt-[10px]" size="small">
      <template #header>
        <div class="text-slate-500 font-bold">
          {{ $t('apps.baseSettings.systemMonitorStatus') }}
        </div>
      </template>

      <div class="flex items-center mt-[5px]">
        <span class="mr-[10px]">{{ $t('common.show') }}</span>
        <NSwitch v-model:value="panelState.panelConfig.systemMonitorShow" size="small" />
      </div>
      <div v-if="panelState.panelConfig.systemMonitorShow" class="flex items-center mt-[5px]">
        <span class="mr-[10px]">{{ $t('apps.baseSettings.showTitle') }}</span>
        <NSwitch v-model:value="panelState.panelConfig.systemMonitorShowTitle" size="small" />
      </div>
      <div v-if="panelState.panelConfig.systemMonitorShow" class="flex items-center mt-[5px]">
        <span class="mr-[10px]">{{ $t('apps.baseSettings.publicVisitModeShow') }}</span>
        <NSwitch v-model:value="panelState.panelConfig.systemMonitorPublicVisitModeShow" size="small" />
      </div>

      <div class="flex items-center mt-[5px]">
        <div>
          {{ t('apps.itemGroupManage.cardStyle') }}
        </div>
        <div class="w-[200px] ml-[5px]">
          <NSelect v-model:value="panelState.panelConfig.iconStyle" size="small" :options="iconTypeOptions" />
        </div>
      </div>
    </NCard>

    <NCard style="border-radius:10px" class="mt-[10px]" size="small">
      <template #header>
        <div class="text-slate-500 font-bold">
          {{ $t('apps.baseSettings.netModeChangeButtonShow') }}
        </div>
      </template>

      <!-- <div class="flex items-center mt-[5px]">
        <span class="mr-[10px]">{{ $t('common.auto') }}</span>
        <STip :content="t('apps.baseSettings.autoNetModeLinkOnOpen')" class="mr-2" />
        <NSwitch v-model:value="panelState.panelConfig.netModeChangeAuto" />
      </div> -->

      <div class="flex items-center mt-[5px]">
        <span class="mr-[10px]">{{ $t('apps.baseSettings.netModeChangeAuto') }}</span>
        <NSwitch v-model:value="panelState.panelConfig.netModeChangeAuto" @update:value="handleNetModeChangeAutoChange" />

        <span class="ml-2">
          <STip>
            <div class=" max-w-[500px]">
              <NTag type="warning" size="small">BETA</NTag>
              {{ t('apps.baseSettings.netModeChangeAutoTip') }}
            </div>
          </STip>
        </span>
      </div>

      <div v-if="panelState.panelConfig.netModeChangeAuto" class="flex items-center mt-[5px]">
        <span class="mr-[10px]">{{ `${$t('apps.baseSettings.netModeChangeAutoMaxPingMs')} (ms)` }}</span>
        <NInputNumber
          v-model:value="panelState.panelConfig.netModeChangeAutoMaxPingMs"
          :step="100" size="small"
          :max="2000" :min="100"
          :style="{ width: '100px' }" placeholder="200"
          @update:value="handleNetModeChangeAutoMaxPingMsChange"
        />
      </div>

      <div class="flex items-center mt-[5px]">
        <span class="mr-[10px]">{{ $t('common.show') }}</span>
        <NSwitch v-model:value="panelState.panelConfig.netModeChangeButtonShow" :disabled="panelState.panelConfig.netModeChangeAuto" size="small" />
      </div>
    </NCard>

    <!-- 1.5.0取消此设置 -->
    <NCard style="border-radius:10px" class="mt-[10px]" size="small">
      <div class="text-slate-500 mb-[5px] font-bold">
        <NFlex align="center">
          {{ $t('common.icon') }}

          <STip :content="t('apps.baseSettings.iconMigrateDescribe')" class="text-cyan-600 text-xl" />
        </NFlex>
      </div>
    </NCard>
    <NCard style="border-radius:10px" class="mt-[10px]" size="small">
      <template #header>
        <div class="text-slate-500 font-bold">
          {{ $t('apps.baseSettings.wallpaper') }}
        </div>
      </template>

      <NUpload
        action="/api/file/uploadImg"
        :show-file-list="false"
        name="imgfile"
        :headers="{
          token: authStore.token as string,
        }"
        :data="{ type: FileType.Wallpaper }"
        :directory-dnd="true"
        @finish="handleUploadBackgroundFinish"
      >
        <NUploadDragger style="width: 100%;">
          <div
            class="h-[200px] w-full border bg-slate-100 flex justify-center items-center cursor-pointer rounded-[10px]"
            :style="{ background: `url(${panelState.panelConfig.backgroundImageSrc}) no-repeat`, backgroundSize: 'cover' }"
          >
            <div class="text-shadow text-white">
              {{ $t('apps.baseSettings.uploadOrDragText') }}
            </div>
          </div>
        </NUploadDragger>
      </NUpload>

      <div class="flex items-center mt-[5px]">
        <span class="mr-[10px]">{{ $t('apps.baseSettings.editImageURL') }}</span>
        <NSwitch v-model:value="showWallpaperInput" size="small" />
      </div>
      <div v-if="showWallpaperInput" class="mt-1">
        <NInput v-model:value="panelState.panelConfig.backgroundImageSrc" type="text" size="small" clearable />
      </div>

      <div class="flex items-center mt-[10px]">
        <span class="mr-[10px]">{{ $t('apps.baseSettings.vague') }}</span>
        <NSlider v-model:value="panelState.panelConfig.backgroundBlur" class="max-w-[200px]" :step="2" :max="20" />
      </div>

      <div class="flex items-center mt-[10px]">
        <span class="mr-[10px]">{{ $t('apps.baseSettings.mask') }}</span>
        <NSlider v-model:value="panelState.panelConfig.backgroundMaskNumber" class="max-w-[200px]" :step="0.1" :max="1" />
      </div>
    </NCard>

    <NCard style="border-radius:10px" class="mt-[10px]" size="small">
      <template #header>
        <div class="text-slate-500 font-bold">
          {{ $t('apps.baseSettings.contentArea') }}
        </div>
      </template>

      <NGrid cols="2">
        <NGridItem span="12 400:12">
          <div class="flex items-center mt-[10px]">
            <span class="mr-[10px]">{{ $t('apps.baseSettings.maxWidth') }}</span>
            <div class="flex">
              <NInputGroup>
                <NInputNumber v-model:value="panelState.panelConfig.maxWidth" size="small" :maxlength="10" :style="{ width: '100px' }" placeholder="1200" />
                <NSelect v-model:value="panelState.panelConfig.maxWidthUnit" :style="{ width: '80px' }" :options="maxWidthUnitOption" size="small" />
              </NInputGroup>
            </div>
          </div>
        </NGridItem>
        <NGridItem span="12 400:12">
          <div class="flex items-center mt-[10px]">
            <span class="mr-[10px]">{{ $t('apps.baseSettings.leftRightMargin') }}</span>
            <NSlider v-model:value="panelState.panelConfig.marginX" class="max-w-[200px]" :step="1" :max="100" />
          </div>
        </NGridItem>
        <NGridItem span="12 400:12">
          <div class="flex items-center mt-[10px]">
            <span class="mr-[10px]">{{ $t('apps.baseSettings.topMargin') }} (%)</span>
            <NSlider v-model:value="panelState.panelConfig.marginTop" class="max-w-[200px]" :step="1" :max="50" />
          </div>
        </NGridItem>
        <NGridItem span="12 400:6">
          <div class="flex items-center mt-[10px]">
            <span class="mr-[10px]">{{ $t('apps.baseSettings.bottomMargin') }} (%)</span>
            <NSlider v-model:value="panelState.panelConfig.marginBottom" class="max-w-[200px]" :step="1" :max="50" />
          </div>
        </NGridItem>
      </NGrid>
    </NCard>

    <NCard style="border-radius:10px" class="mt-[10px]" size="small">
      <template #header>
        <div class="text-slate-500 font-bold flex items-center">
          <span>
            {{ $t('apps.baseSettings.customFooter') }}
          </span>
          <span class="mx-2">
            <STip class="text-base text-orange-500">
              <div style="max-width: 400px">
                {{ t('apps.baseSettings.customFooterAlert') }}
              </div>
            </STip>
          </span>
        </div>
      </template>

      <NInput
        v-model:value="panelState.panelConfig.footerHtml"
        type="textarea"
        clearable
      />
    </NCard>

    <NCard style="border-radius:10px" class="mt-[10px]" size="small">
      <NButton size="small" quaternary type="error" @click="resetPanelConfig">
        {{ $t('common.reset') }}
      </NButton>

      <NButton size="small" quaternary type="success" class="ml-[10px]" @click="uploadCloud">
        {{ $t('common.save') }}
      </NButton>
    </NCard>
  </div>
</template>

<style scoped>
.text-shadow{
  text-shadow: 0px 0px 5px gray;
}
</style>
