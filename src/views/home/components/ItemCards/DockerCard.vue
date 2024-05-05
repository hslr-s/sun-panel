<script setup lang="ts">
import { computed, nextTick, onMounted, onUnmounted, ref } from 'vue'
import type { DropdownOption } from 'naive-ui'
import { NDropdown, NEllipsis, NModal, NSkeleton, NSpin } from 'naive-ui'
import type { PanelPanelConfigStyleEnum } from '@/enums'
import { ItemCard, ItemIcon } from '@/components/common'
import { t } from '@/locales'
import { PanelStateNetworkModeEnum } from '@/enums'
import { DockerContainerInfoState } from '@/enums/docker'
import { useAuthStore, usePanelState } from '@/store'
import { execSwitchAction as execSwitchActionApi, getContainerStatus as getContainerStatusApi } from '@/api/docker/container'

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
let intervalId: any
const showBadgeStatus = ref<DockerContainerInfoState[]>([])

const windowShow = ref<boolean>(false)
const windowSrc = ref<string>('')
const windowTitle = ref<string>('')
const windowIframeRef = ref(null)
const windowIframeIsLoad = ref<boolean>(false)

const dropdownMenuX = ref(0)
const dropdownMenuY = ref(0)
const dropdownShow = ref(false)
const badgeStatusColor = ref('')
const loading = ref(false)

const containerState = ref<Docker.Container.Stats>({
  cpuPercent: 0,
  memory: 0,
  cache: 0,
  ioRead: 0,
  ioWrite: 0,
  networkRX: 0,
  networkTX: 0,
  runStatus: DockerContainerInfoState.exited,
})
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

  // console.log('项目呗点击', containerState.value.runStatus, props.itemInfo.expandParam.openPageAndStart)
  if (containerState.value.runStatus === DockerContainerInfoState.exited && props.itemInfo.expandParam.openPageAndStart === true) {
    execSwitchAction(true, () => {
      openPage(props.itemInfo.openMethod, jumpUrl, props.itemInfo?.title)

      // 容易被拦截暂时取消延时
      // const delaySecond = props.itemInfo.expandParam.openPageAndStartDelaySecond
      // if (delaySecond !== 0) {
      //   // console.log('延时打开页面')
      //   setTimeout(() => {
      //     openPage(props.itemInfo.openMethod, jumpUrl, props.itemInfo?.title)
      //   }, delaySecond * 1000)
      // }
      // else {
      //   // console.log('不延时打开页面')
      //   openPage(props.itemInfo.openMethod, jumpUrl, props.itemInfo?.title)
      // }
    })
  }
  else {
    openPage(props.itemInfo.openMethod, jumpUrl, props.itemInfo?.title)
  }
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

    case 'startContainer':
      execSwitchAction(true)
      break

    case 'stopContainer':
      execSwitchAction(false)
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
        label: t('iconItem.dockerAction.startContainer'),
        key: 'startContainer',
      },
      {
        label: t('iconItem.dockerAction.stopContainer'),
        key: 'stopContainer',
      },
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

async function execSwitchAction(action: boolean, startSuccessCallback?: () => void) {
  loading.value = true
  const id = props.itemInfo.expandParam.containerId
  await execSwitchActionApi({
    id, action,
  }).then(() => {
    if (action) {
      containerState.value.runStatus = DockerContainerInfoState.running
      console.log(containerState.value.runStatus)
      if (startSuccessCallback)
        startSuccessCallback()
    }
    else {
      containerState.value.runStatus = DockerContainerInfoState.exited
    }

    // 成功
  }).catch(() => {
    // 失败
  })
  loading.value = false
}

async function getContainerStatus(showLoading?: boolean) {
  if (showLoading)
    loading.value = true
  await getContainerStatusApi<Docker.Container.Stats>(props.itemInfo.expandParam.containerId).then(({ data }) => {
    containerState.value = data

    // 设置徽章的颜色
    switch (containerState.value.runStatus) {
      case DockerContainerInfoState.running:
        badgeStatusColor.value = '#18A058'
        break
      case DockerContainerInfoState.restarting:
        badgeStatusColor.value = '#F0A020FF'
        break

      default:
        badgeStatusColor.value = '#7B7B7B'
        break
    }
  }).catch(() => {

  })
  if (showLoading)
    loading.value = false
}

function roundNumberToTwoDecimalPlaces(num: number): number {
  return Number(num.toFixed(2))
}

function inShowBadgeStatus(value: DockerContainerInfoState): boolean {
  // 检查 arr 是否为数组
  if (!Array.isArray(showBadgeStatus.value))
    throw new Error('The first argument must be an array')

  return showBadgeStatus.value.includes(value)
}

function isShowBadgeStatus(value: DockerContainerInfoState): boolean {
  return inShowBadgeStatus(value) && containerState.value.runStatus === value
}

onMounted(() => {
  getContainerStatus(true)

  if (props.itemInfo.expandParam.showBadgeStatus)
    showBadgeStatus.value = props.itemInfo.expandParam.showBadgeStatus

  intervalId = setInterval(() => {
    getContainerStatus()
  }, 5000) // 5000毫秒，即5秒
})

onUnmounted(() => {
  clearInterval(intervalId)
})
</script>

<template>
  <div class="app-icon w-full" @click="handleItemClick" @contextmenu="(e) => handleContextMenu(e)">
    <ItemCard
      :card-type-style="cardStyle"
      :background-color="backgroundColor"
      :icon-text-icon-hide-title="iconTextIconHideTitle"
      :icon-text-color="iconTextColor"
      :icon-text="itemInfo?.title"
      class="hover:shadow-[0_0_10px_5px_rgba(0,0,0,0.2)] transition-all duration-300 backdrop-blur-[2px] hover:backdrop-blur-[10px]"
    >
      <template #small>
        <NSpin :show="loading">
          <div
            class="app-icon-small-icon overflow-hidden rounded-2xl w-[70px] h-[70px] mx-auto rounded-2xl relative"
          >
            <ItemIcon :item-icon="itemInfo?.icon" :color="textColor" />
            <div
              v-if="isShowBadgeStatus(containerState.runStatus as DockerContainerInfoState)"
              class="dot dot-small flex items-center justify-center border"
              :style="{ backgroundColor: badgeStatusColor }"
            >
              <!-- <SvgIconOnline icon="tabler:activity-heartbeat" /> -->
            </div>
          </div>
        </NSpin>
      </template>

      <template #info>
        <NSpin :show="loading" style="width: 100%;">
          <div class="flex">
            <!-- 图标 -->
            <div class="app-icon-info-icon w-[70px] h-[70px]">
              <div class="w-[70px] h-full flex items-center justify-center">
                <ItemIcon :item-icon="itemInfo?.icon" :size="50" class="overflow-hidden rounded-xl" :color="textColor" />
              </div>
            </div>

            <div
              v-if="isShowBadgeStatus(containerState.runStatus as DockerContainerInfoState)"
              class="dot dot-info flex items-center justify-center border"
              :style="{ backgroundColor: badgeStatusColor }"
            >
              <!-- <SvgIconOnline icon="tabler:activity-heartbeat" /> -->
            </div>

            <!-- 文字 -->
            <!-- 如果为纯白色，将自动根据背景的明暗计算字体的黑白色 -->
            <div
              class="text-white flex items-center"
              :style="{ color: (!iconTextColor || iconTextColor === '#ffffff') ? textColor : iconTextColor, maxWidth: 'calc(100% - 80px)' }"
            >
              <div class="app-icon-info-text-box w-full">
                <div class="app-icon-info-text-box-title font-semibold w-full">
                  <NEllipsis>
                    {{ itemInfo?.title }}
                  </NEllipsis>
                </div>
                <div v-if="!iconTextInfoHideDescription" class="app-icon-info-text-box-description">
                  <NEllipsis v-if="containerState.runStatus === DockerContainerInfoState.running" :line-clamp="2" class="text-xs">
                    CPU:{{ roundNumberToTwoDecimalPlaces(containerState.cpuPercent) }}%
                    RAM:{{ roundNumberToTwoDecimalPlaces(containerState.memory) }} MB
                  </NEllipsis>
                  <NEllipsis v-else-if="containerState.runStatus === DockerContainerInfoState.restarting" :line-clamp="2" class="text-xs">
                    {{ t('iconItem.showBadgeStatusRestarting') }}
                  </NEllipsis>
                  <NEllipsis v-else :line-clamp="2" class="text-xs">
                    {{ t('iconItem.showBadgeStatusExited') }}
                  </NEllipsis>
                </div>
              </div>
            </div>
          </div>
        </NSpin>
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

<style scoped>
.dot {
  position: absolute;
  color: white;
  border-radius: 50%;
}

.dot-small{
  top: 6px;
  right: 6px;
  width: 10px;
  height: 10px;
  font-size: 10px;
}

.dot-info{
  top: 8px;
  right: 8px;
  width: 10px;
  height: 10px;
}
</style>
