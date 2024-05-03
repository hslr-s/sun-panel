<script setup lang="ts">
import { computed, nextTick, ref } from 'vue'
import { NDropdown, NEllipsis, NModal, NSkeleton, NSpin, useDialog, useMessage } from 'naive-ui'
import type { PanelPanelConfigStyleEnum } from '@/enums'
import { ItemCard, ItemIcon } from '@/components/common'
import { t } from '@/locales'
import { PanelStateNetworkModeEnum } from '@/enums'
import { useAuthStore, usePanelState } from '@/store'

import { VisitMode } from '@/enums/auth'
interface Prop {
  itemInfo: Panel.ItemCard.CardInfo
  backgroundColor?: string
  iconTextColor?: string
  iconTextInfoHideDescription: boolean
  iconTextIconHideTitle: boolean
  cardStyle: PanelPanelConfigStyleEnum
  sortStatus: boolean
}
const props = defineProps<Prop>()
const emit = defineEmits<{
  (e: 'delete'): void
  (e: 'edit'): void
}>()
const panelState = usePanelState()
const authStore = useAuthStore()
const dialog = useDialog()

const windowShow = ref<boolean>(false)
const windowSrc = ref<string>('')
const windowTitle = ref<string>('')
const windowIframeRef = ref(null)
const windowIframeIsLoad = ref<boolean>(false)

const dropdownMenuX = ref(0)
const dropdownMenuY = ref(0)
const dropdownShow = ref(false)

const ms = useMessage()
const defaultBackground = '#2a2a2a6b'

const calculateLuminance = (color: string) => {
  const hex = color.replace(/^#/, '')
  const r = parseInt(hex.substring(0, 2), 16)
  const g = parseInt(hex.substring(2, 4), 16)
  const b = parseInt(hex.substring(4, 6), 16)
  return (0.299 * r + 0.587 * g + 0.114 * b) / 255
}

function openPage(openMethod: number, url: string, title?: string) {
  switch (openMethod) {
    case 1:
      window.location.href = url
      break
    case 2:
      window.open(url)
      break
    case 3:
      windowShow.value = true
      windowSrc.value = url
      windowTitle.value = title || url
      windowIframeIsLoad.value = true
      break

    default:
      break
  }
}

function handWindowIframeIdLoad(payload: Event) {
  windowIframeIsLoad.value = false
}

const textColor = computed(() => {
  const luminance = calculateLuminance(props.backgroundColor || defaultBackground)
  return luminance > 0.5 ? 'black' : 'white'
})

// 点击项目
function handleItemClick() {
  if (props.sortStatus) {
    emit('edit')
    return
  }

  let jumpUrl = ''

  if (props.itemInfo)
    jumpUrl = (panelState.networkMode === PanelStateNetworkModeEnum.lan ? props.itemInfo.lanUrl : props.itemInfo.url) as string

  if (props.itemInfo?.lanUrl === '')
    jumpUrl = props.itemInfo.url

  openPage(props.itemInfo.openMethod, jumpUrl, props.itemInfo?.title)
}

// 右键
function handleContextMenu(e: MouseEvent) {
  e.preventDefault()

  if (props.sortStatus)
    return

  // currentRightSelectItem.value = item
  dropdownShow.value = false
  nextTick().then(() => {
    dropdownShow.value = true
    dropdownMenuX.value = e.clientX
    dropdownMenuY.value = e.clientY
  })
}

function handleRightMenuSelect(key: string | number) {
  dropdownShow.value = false
  // console.log(currentRightSelectItem, key)
  let jumpUrl = panelState.networkMode === PanelStateNetworkModeEnum.lan ? props.itemInfo?.lanUrl : props.itemInfo?.url
  if (props.itemInfo?.lanUrl === '')
    jumpUrl = props.itemInfo.url
  switch (key) {
    case 'newWindows':
      window.open(jumpUrl)
      break
    case 'openWanUrl':
      if (props.itemInfo)
        openPage(props.itemInfo?.openMethod, props.itemInfo?.url, props.itemInfo?.title)
      break
    case 'openLanUrl':
      if (props.itemInfo && props.itemInfo.lanUrl)
        openPage(props.itemInfo?.openMethod, props.itemInfo.lanUrl, props.itemInfo?.title)
      break
    case 'edit':
      // 这里有个奇怪的问题，如果不使用{...}的方式 父组件的值会同步修改 标记一下
      emit('edit')
      break
    case 'delete':
      emit('delete')
      break
    default:
      break
  }
}

function getDropdownMenuOptions() {
  const dropdownMenuOptions = [
    {
      label: t('iconItem.newWindowOpen'),
      key: 'newWindows',
    },

  ]

  if (props.itemInfo?.lanUrl && panelState.networkMode === PanelStateNetworkModeEnum.wan) {
    dropdownMenuOptions.push({
      label: t('panelHome.openLanUrl'),
      key: 'openLanUrl',
    })
  }

  if (props.itemInfo?.lanUrl && panelState.networkMode === PanelStateNetworkModeEnum.lan) {
    dropdownMenuOptions.push({
      label: t('panelHome.openWanUrl'),
      key: 'openWanUrl',
    })
  }

  if (authStore.visitMode === VisitMode.VISIT_MODE_LOGIN) {
    dropdownMenuOptions.push({
      label: t('common.edit'),
      key: 'edit',
    }, {
      label: t('common.delete'),
      key: 'delete',
    })
  }

  return dropdownMenuOptions
}

function onClickoutside() {
  // message.info('clickoutside')
  dropdownShow.value = false
}
</script>

<template>
  <div class="app-icon w-full" @click="handleItemClick" @contextmenu="(e) => handleContextMenu(e)">
    <ItemCard
      :card-type-style="cardStyle"
      :background-color="backgroundColor"
      :icon-text-icon-hide-title="iconTextIconHideTitle"
      :icon-text-color="iconTextColor"
      :icon-text="itemInfo?.title"
      class="hover:shadow-[0_0_20px_10px_rgba(0,0,0,0.2)]"
    >
      <template #small>
        <div
          class="app-icon-small-icon overflow-hidden rounded-2xl w-[70px] h-[70px] mx-auto rounded-2xl transition-all duration-200 hover:shadow-[0_0_20px_10px_rgba(0,0,0,0.2)]"
          :title="itemInfo?.description"
        >
          <ItemIcon :item-icon="itemInfo?.icon" />
        </div>
        <div
          v-if="!iconTextIconHideTitle"
          class="app-icon-small-title text-center app-icon-text-shadow cursor-pointer mt-[2px]"
          :style="{ color: iconTextColor }"
        >
          <span>{{ itemInfo?.title }}</span>
        </div>
      </template>

      <template #info>
        <!-- 图标 -->
        <div class="app-icon-info-icon w-[70px] h-[70px]">
          <div class="w-[70px] h-full flex items-center justify-center ">
            <ItemIcon :item-icon="itemInfo?.icon" :size="50" class="overflow-hidden rounded-xl" />
          </div>
        </div>

        <!-- 文字 -->
        <!-- 如果为纯白色，将自动根据背景的明暗计算字体的黑白色 -->
        <div class="text-white flex items-center" :style="{ color: (iconTextColor === '#ffffff') ? textColor : iconTextColor, maxWidth: 'calc(100% - 80px)' }">
          <div class="app-icon-info-text-box w-full">
            <div class="app-icon-info-text-box-title font-semibold w-full">
              <NEllipsis>
                {{ itemInfo?.title }}
              </NEllipsis>
            </div>
            <div v-if="!iconTextInfoHideDescription" class="app-icon-info-text-box-description">
              <NEllipsis :line-clamp="2" class="text-xs">
                {{ itemInfo?.description }}
              </NEllipsis>
            </div>
          </div>
        </div>
      </template>
    </ItemCard>

    <!-- 右键菜单 -->
    <NDropdown
      placement="bottom-start" trigger="manual" :x="dropdownMenuX" :y="dropdownMenuY"
      :options="getDropdownMenuOptions()" :show="dropdownShow" :on-clickoutside="onClickoutside" @select="handleRightMenuSelect"
    />

    <!-- 弹窗 -->
    <NModal
      v-model:show="windowShow" :mask-closable="false" preset="card"
      style="max-width: 1000px;height: 600px;border-radius: 1rem;" :bordered="true" size="small" role="dialog"
      aria-modal="true"
    >
      <template #header>
        <div class="flex items-center">
          <span class="mr-[20px]">
            {{ windowTitle }}
          </span>

          <NSpin v-if="windowIframeIsLoad" size="small" />
        </div>
      </template>
      <div class="w-full h-full rounded-2xl overflow-hidden border dark:border-zinc-700">
        <div v-if="windowIframeIsLoad" class="flex flex-col p-5">
          <NSkeleton height="50px" width="100%" class="rounded-lg" />
          <NSkeleton height="180px" width="100%" class="mt-[20px] rounded-lg" />
          <NSkeleton height="180px" width="100%" class="mt-[20px] rounded-lg" />
        </div>
        <iframe
          v-show="!windowIframeIsLoad" id="windowIframeId" ref="windowIframeRef" :src="windowSrc"
          class="w-full h-full" frameborder="0" @load="handWindowIframeIdLoad"
        />
      </div>
    </NModal>
  </div>
</template>
