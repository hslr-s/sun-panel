<script setup lang="ts">
import { nextTick, ref } from 'vue'
import { NDropdown } from 'naive-ui'
import type { DropdownOption } from 'naive-ui'
import StandardIcon from './common/StandardIcon.vue'
import { getItemUrl } from './common'
import type { PanelPanelConfigStyleEnum } from '@/enums'
import { t } from '@/locales'
import { PanelStateNetworkModeEnum } from '@/enums'
import { useAuthStore, usePanelState } from '@/store'
import { IframeLayer } from '@/components/common'

import { VisitMode } from '@/enums/auth'
interface Prop {
  itemInfo: Panel.ItemInfo
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
// const dialog = useDialog()

const windowShow = ref<boolean>(false)
const windowSrc = ref<string>('')
const windowTitle = ref<string>('')
const windowIframeIsLoad = ref<boolean>(false)

const dropdownMenuX = ref(0)
const dropdownMenuY = ref(0)
const dropdownShow = ref(false)

// const ms = useMessage()
// const defaultBackground = defaultItemCardBackgroundColor

// const calculateLuminance = (color: string) => {
//   const hex = color.replace(/^#/, '')
//   const r = parseInt(hex.substring(0, 2), 16)
//   const g = parseInt(hex.substring(2, 4), 16)
//   const b = parseInt(hex.substring(4, 6), 16)
//   return (0.299 * r + 0.587 * g + 0.114 * b) / 255
// }

function openPage(openMethod: number, url: string, title?: string) {
  switch (openMethod) {
    case 1:
      if (window.self !== window.top)
        window.parent.location.href = url
      else
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

// function handWindowIframeIdLoad(payload: Event) {
//   windowIframeIsLoad.value = false
// }

// 点击项目
async function handleItemClick() {
  if (props.sortStatus) {
    emit('edit')
    return
  }

  const jumpUrl = await getItemUrl(props.itemInfo, panelState.panelConfig.netModeChangeAutoMaxPingMs)

  openPage(props.itemInfo.openMethod, jumpUrl, props.itemInfo?.title)
}

// 右键
function handleContextMenu(event: MouseEvent) {
  event.preventDefault()

  if (props.sortStatus)
    return

  // currentRightSelectItem.value = item
  dropdownShow.value = false
  nextTick().then(() => {
    dropdownShow.value = true
    dropdownMenuX.value = event.clientX
    dropdownMenuY.value = event.clientY
  })
}

async function handleMousedowClick(event: MouseEvent) {
  // 鼠标中间键点击
  if (event.button === 1 && !props.sortStatus) {
    event.preventDefault()
    window.open(await getItemUrl(props.itemInfo, panelState.panelConfig.netModeChangeAutoMaxPingMs))
  }
}

// // 获取项目地址，识别当前内外网环境返回地址，废弃-采用 import { getItemUrl } from './common' 2024年12月1日记录
// function getItemUrl() {
//   let jumpUrl = ''
//   if (props.itemInfo)
//     jumpUrl = (panelState.networkMode === PanelStateNetworkModeEnum.lan ? props.itemInfo.lanUrl : props.itemInfo.url) as string

//   if (props.itemInfo?.lanUrl === '')
//     jumpUrl = props.itemInfo.url

//   return jumpUrl
// }

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
  const dropdownMenuOptions: DropdownOption[] = [
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
    dropdownMenuOptions.push(
      {
        type: 'divider',
        key: 'd1',
      },
      {
        label: t('common.edit'),
        key: 'edit',
      },
      {
        label: t('common.delete'),
        key: 'delete',
      },
    )
  }

  return dropdownMenuOptions
}

function onClickoutside() {
  // message.info('clickoutside')
  dropdownShow.value = false
}
</script>

<template>
  <div class="app-icon w-full" @click="handleItemClick" @contextmenu="(e) => handleContextMenu(e)" @mousedown="handleMousedowClick">
    <StandardIcon
      :item-info="itemInfo"
      :card-type-style="cardStyle"
      :background-color="backgroundColor"
      :icon-text-icon-hide-title="iconTextIconHideTitle"
      :icon-text-color="iconTextColor"
      :icon-text="itemInfo?.title"
      :card-style="cardStyle"
      :sort-status="sortStatus"
      :icon-text-info-hide-description="iconTextInfoHideDescription"
      :title="itemInfo.description"
    />

    <!-- 右键菜单 -->
    <NDropdown
      placement="bottom-start" trigger="manual" :x="dropdownMenuX" :y="dropdownMenuY"
      :options="getDropdownMenuOptions()" :show="dropdownShow" :on-clickoutside="onClickoutside" @select="handleRightMenuSelect"
    />

    <!-- 弹窗 -->
    <IframeLayer
      v-model:show="windowShow"
      :title="windowTitle"
      :src="windowSrc"
      move
      resize
    />

    <!-- 旧弹窗 -->
    <!-- <NModal
      v-model:show="windowShow" :mask-closable="false" preset="card"
      style="max-width: 1200px;height: 650px;border-radius: 1rem;" :bordered="true" size="small" role="dialog"
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
    </NModal> -->
  </div>
</template>
