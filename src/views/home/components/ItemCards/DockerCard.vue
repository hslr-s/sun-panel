<script setup lang="ts">
import { computed, nextTick, onMounted, onUnmounted, ref, watch } from 'vue'
import type { DropdownOption } from 'naive-ui'
import { NDropdown, NEllipsis, NFlex, NModal, NSkeleton, NSpin, useMessage } from 'naive-ui'
import { useDockerCardStore } from './dockerCard/store'
import type { PanelPanelConfigStyleEnum } from '@/enums'
import { ItemCard, ItemIcon, STip } from '@/components/common'
import { t } from '@/locales'
import { PanelStateNetworkModeEnum } from '@/enums'
import { DockerContainerInfoState } from '@/enums/docker'
import { useAuthStore, usePanelState } from '@/store'
import {
  execSwitchAction as execSwitchActionApi,
  getContainerIDByName as getContainerIDByNameApi,
  getContainerStatus as getContainerStatusApi,
} from '@/api/docker/container'

import { edit as itemIconEditApi } from '@/api/panel/itemIcon'

import { VisitMode } from '@/enums/auth'
import { apiRespErrMsg } from '@/utils/request/apiMessage'
import { bytesToSize } from '@/utils/cmn'
import { useDomainLinks } from '@/hooks/useDomainLinks'

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

let isRegistListenFinish = false

const panelState = usePanelState()
const authStore = useAuthStore()
const domainLinks = useDomainLinks()
// let intervalId: any
const showBadgeStatus = ref<DockerContainerInfoState[]>([])
const ms = useMessage()

const windowShow = ref<boolean>(false)
const windowSrc = ref<string>('')
const windowTitle = ref<string>('')
const windowIframeRef = ref(null)
const windowIframeIsLoad = ref<boolean>(false)

const itemCardInfo = ref(props.itemInfo)

const store = useDockerCardStore()

const dropdownMenuX = ref(0)
const dropdownMenuY = ref(0)
const dropdownShow = ref(false)
const badgeStatusColor = ref('')
const loading = ref(false)
const badgeMsg = ref('')
const dockerErrorMsg = ref('')
const lastTime = ref('')

const containerState = ref<Docker.Container.Stats>({
  cpuPercent: 0,
  memory: 0,
  cache: 0,
  ioRead: 0,
  ioWrite: 0,
  networkRX: 0,
  networkTX: 0,
  runStatus: DockerContainerInfoState.exited,
  errorMessage: '',
})
const defaultBackground = '#2a2a2a6b'
const describeText = ref('')

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

  if (itemCardInfo.value)
    jumpUrl = (panelState.networkMode === PanelStateNetworkModeEnum.lan ? itemCardInfo.value.lanUrl : itemCardInfo.value.url) as string

  if (itemCardInfo.value?.lanUrl === '')
    jumpUrl = itemCardInfo.value.url

  if (containerState.value.runStatus === DockerContainerInfoState.exited && itemCardInfo.value.expandParam.openPageAndStart === true) {
    // 判断当前是否为公开模式，公开模式进行拦截
    if (authStore.visitMode === VisitMode.VISIT_MODE_PUBLIC) {
      ms.warning(t('common.publicModeRestricted'))
      return
    }

    execSwitchAction(true, () => {
      openPage(itemCardInfo.value.openMethod, jumpUrl, itemCardInfo.value?.title)

      // 容易被拦截暂时取消延时
      // const delaySecond = itemCardInfo.value.expandParam.openPageAndStartDelaySecond
      // if (delaySecond !== 0) {
      //   // console.log('延时打开页面')
      //   setTimeout(() => {
      //     openPage(itemCardInfo.value.openMethod, jumpUrl, itemCardInfo.value?.title)
      //   }, delaySecond * 1000)
      // }
      // else {
      //   // console.log('不延时打开页面')
      //   openPage(itemCardInfo.value.openMethod, jumpUrl, itemCardInfo.value?.title)
      // }
    })
  }
  else {
    openPage(itemCardInfo.value.openMethod, jumpUrl, itemCardInfo.value?.title)
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
  let jumpUrl = panelState.networkMode === PanelStateNetworkModeEnum.lan ? itemCardInfo.value?.lanUrl : itemCardInfo.value?.url
  if (itemCardInfo.value?.lanUrl === '')
    jumpUrl = itemCardInfo.value.url
  switch (key) {
    case 'newWindows':
      window.open(jumpUrl)
      break
    case 'openWanUrl':
      if (itemCardInfo.value)
        openPage(itemCardInfo.value?.openMethod, itemCardInfo.value?.url, itemCardInfo.value?.title)
      break
    case 'openLanUrl':
      if (itemCardInfo.value && itemCardInfo.value.lanUrl)
        openPage(itemCardInfo.value?.openMethod, itemCardInfo.value.lanUrl, itemCardInfo.value?.title)
      break

    case 'startContainer':
      execSwitchAction(true)
      break

    case 'stopContainer':
      execSwitchAction(false)
      break
    case 'edit':
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

  if (itemCardInfo.value?.lanUrl && panelState.networkMode === PanelStateNetworkModeEnum.wan) {
    dropdownMenuOptions.push({
      label: t('panelHome.openLanUrl'),
      key: 'openLanUrl',
    })
  }

  if (itemCardInfo.value?.lanUrl && panelState.networkMode === PanelStateNetworkModeEnum.lan) {
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
  const containerId = itemCardInfo.value.expandParam.containerId
  await execSwitchActionApi({
    id: containerId, action,
  }).then(async () => {
    // 成功
    if (action) {
      containerState.value.runStatus = DockerContainerInfoState.running
      // console.log(containerState.value.runStatus)
      if (startSuccessCallback)
        startSuccessCallback()
    }
    else {
      containerState.value.runStatus = DockerContainerInfoState.exited
    }
    await getContainerStatus()
    loading.value = false
  }).catch((res) => {
    // 失败
    loading.value = false
    if (res.code === -1) {
      dockerErrorMsg.value = res.msg
      return
    }

    apiRespErrMsg(res)
  })
}

// 步骤：若容器名存在，将使用容器名查询出容器id，得到容器id后将更新卡片信息
// 参数：stopCycle：终止循环
async function fixCardDockerContainerIdByContainerName(containerId: string, containerName: string) {
  // 登录用户才能操作
  if (authStore.visitMode === VisitMode.VISIT_MODE_LOGIN) {
    await getContainerIDByNameApi(containerName).then(({ data }) => {
      itemCardInfo.value.expandParam.containerId = data.id
      // 更新卡片信息
      itemIconEditApi(itemCardInfo.value).then(() => {
        registAndListenContainerStatsUpdate(itemCardInfo.value.expandParam.containerId)
        getContainerStatus(true)
      }).catch(() => {
        // 更新卡片失败
        badgeStatusColor.value = 'red'
        badgeMsg.value = t('iconItem.dockerContainerError')
      })
    }).catch(() => {
      // 容器不存在
      badgeStatusColor.value = 'red'
      badgeMsg.value = t('iconItem.dockerContainerError')
    })
  }
  else {
    badgeStatusColor.value = 'red'
    badgeMsg.value = t('iconItem.dockerContainerError')
  }
}

// 获取容器状态
async function getContainerStatus(showLoading?: boolean) {
  const containerId = itemCardInfo.value.expandParam.containerId as string
  const containerName = itemCardInfo.value.expandParam.containerName as string
  if (!containerId)
    return

  if (showLoading)
    loading.value = true
  await getContainerStatusApi<Docker.Container.Stats>(containerId).then(({ data }) => {
    lastTime.value = data.shotTime || ''
    updateStats(data)
    badgeMsg.value = ''
    registAndListenContainerStatsUpdate(containerId)
  }).catch((res) => {
    // 容器异常（尝试自动修复）例如升级导致容器id变更，此处遍历容器尝试更新
    if (containerName) {
      (async function () {
        await fixCardDockerContainerIdByContainerName(containerId, containerName)
      }())
    }

    if (res.code === -1) {
      dockerErrorMsg.value = res.msg
      return
    }

    apiRespErrMsg(res)
  })
  if (showLoading)
    loading.value = false
}

function updateStats(data: Docker.Container.Stats) {
  if (data) {
    const d1 = new Date(data.shotTime || '')
    const lt = new Date(lastTime.value || '')

    // 如果任意两个时间有一个为空或者数据时间大于上次更新时间，更新！
    if ((data.shotTime !== '' && lastTime.value !== '') && d1 < lt) {
      return // 不更新
    }

    lastTime.value = data.shotTime || ''
  }

  containerState.value = data
  describeText.value = ''

  if (data.errorMessage !== '') {
    badgeMsg.value = t('iconItem.dockerContainerError')
    dockerErrorMsg.value = data.errorMessage
  }

  if (containerState.value.runStatus) {
    describeText.value = `CPU:${roundNumberToTwoDecimalPlaces(containerState.value.cpuPercent)}% RAM:${bytesToSize(containerState.value.memory * 1024 * 1024)}`
    // 设置徽章的颜色
    switch (containerState.value.runStatus) {
      case DockerContainerInfoState.running:
        badgeStatusColor.value = '#18A058'
        break
      case DockerContainerInfoState.restarting:
        badgeStatusColor.value = '#F0A020FF'
        break
      case DockerContainerInfoState.error:
        badgeStatusColor.value = 'red'
        break

      default:
        badgeStatusColor.value = '#7B7B7B'
        break
    }
    // console.log(describeText.value, badgeStatusColor.value)
  }
}

function roundNumberToTwoDecimalPlaces(num: number): number {
  return Number(num.toFixed(2))
}

function roundNumberToOneDecimalPlaces(num: number): number {
  return Number(num.toFixed(1))
}

function inShowBadgeStatus(value: DockerContainerInfoState): boolean {
  // 检查 arr 是否为数组
  if (!Array.isArray(showBadgeStatus.value))
    throw new Error('The first argument must be an array')

  return showBadgeStatus.value.includes(value) || value === DockerContainerInfoState.error
}

function isShowBadgeStatus(value: DockerContainerInfoState): boolean {
  return inShowBadgeStatus(value) && containerState.value.runStatus === value
}

function registAndListenContainerStatsUpdate(containerId: string) {
  // 需要保持最新的数据，防止被定时的数据刷掉，可能需要加入一个时间字段来进行时间新旧对比判断是否更新
  if (!isRegistListenFinish) {
    isRegistListenFinish = true

    watch(() => store.containerRuningStats[containerId], (value) => {
      const isEmpty = Object.keys(value).length === 0

      if (!isEmpty && value.runStatus) {
        updateStats(value)
      }
    })

    store.registerContainerQueryStats(containerId)
  }
}

onMounted(async () => {
  if (itemCardInfo.value.expandParam.containerId) {
    await getContainerStatus(true)
  }

  if (itemCardInfo.value.expandParam.showBadgeStatus)
    showBadgeStatus.value = itemCardInfo.value.expandParam.showBadgeStatus
})

function formatMemory(sizeInMB: number): string {
  // 定义单位数组
  const units = ['MB', 'GB', 'TB']
  let size = sizeInMB
  let unitIndex = 0

  // 转换到合适的单位，直到小于1024
  while (size >= 1024) {
    size /= 1024
    unitIndex++
  }

  // 四舍五入到小数点后一位
  let roundedSize = Math.round(size * 10) / 10

  // 检查整数部分长度是否超过3位
  const integerPart = Math.floor(roundedSize) // 获取整数部分
  if (integerPart.toString().length >= 3) {
    // 如果整数部分超过3位，向上取整到最近的整数，并转换到下一个单位
    roundedSize = Math.ceil(roundedSize)
    unitIndex++
    roundedSize /= 1024
  }

  // 格式化为字符串，保留小数点后一位
  const formattedSize = roundedSize.toFixed(1)

  // 返回格式化后的数值和单位
  return `${formattedSize} ${units[unitIndex]}`
}

function handleDockerErrorMsgClick() {
  // console.log('异常点击', domainLinks.faqLinkWithLang('docker_card_abnormal'))
}

onUnmounted(() => {
  // 取消监听容器的状态
  if (isRegistListenFinish) {
    store.cancelContainerQueryStats(itemCardInfo.value.expandParam.containerId)
  }
})
</script>

<template>
  <div
    class="app-icon w-full"
    :title="describeText"
    @click="handleItemClick"
    @contextmenu="(e) => handleContextMenu(e)"
  >
    <ItemCard
      :card-type-style="cardStyle" :background-color="backgroundColor"
      :icon-text-icon-hide-title="iconTextIconHideTitle" :icon-text-color="iconTextColor" :icon-text="itemCardInfo?.title"
    >
      <template #small>
        <NSpin :show="loading">
          <div class="app-icon-small-icon overflow-hidden w-[70px] h-[70px] mx-auto rounded-2xl relative">
            <ItemIcon :item-icon="itemCardInfo?.icon" :color="textColor" />
            <div
              v-if="isShowBadgeStatus(containerState.runStatus as DockerContainerInfoState)"
              class="dot dot-small flex items-center justify-center border"
              :style="{ backgroundColor: badgeStatusColor }" :title="badgeMsg"
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
                <ItemIcon
                  :item-icon="itemCardInfo?.icon" :size="50" class="overflow-hidden rounded-xl"
                  :color="textColor"
                />
              </div>
            </div>

            <div
              v-if="isShowBadgeStatus(containerState.runStatus as DockerContainerInfoState)"
              class="dot dot-info flex items-center justify-center border"
              :style="{ backgroundColor: badgeStatusColor }" :title="badgeMsg"
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
                    {{ itemCardInfo?.title }}
                  </NEllipsis>
                </div>
                <div v-if="!iconTextInfoHideDescription" class="app-icon-info-text-box-description">
                  <NEllipsis
                    v-if="containerState.runStatus === DockerContainerInfoState.running" :line-clamp="2"
                    class="text-xs"
                  >
                    CPU:{{ roundNumberToOneDecimalPlaces(containerState.cpuPercent) }}%
                    RAM:{{ formatMemory(containerState.memory) }}
                  </NEllipsis>
                  <NEllipsis
                    v-else-if="containerState.runStatus === DockerContainerInfoState.restarting"
                    :line-clamp="2" class="text-xs"
                  >
                    {{ t('iconItem.showBadgeStatusRestarting') }}
                    <span class="text-[red]">
                      {{ badgeMsg }}
                    </span>
                  </NEllipsis>
                  <NFlex v-else :line-clamp="2" class="text-xs" align="center">
                    {{ t('iconItem.showBadgeStatusExited') }}
                    <span v-if="dockerErrorMsg" class="text-[red]" :title="dockerErrorMsg" @click.stop="handleDockerErrorMsgClick">
                      <STip>
                        <div class="max-w-[300px]">
                          {{ badgeMsg }}: {{ dockerErrorMsg }}
                          <a target="_blank" class="text-lime-600" :href="domainLinks.faqLinkWithLang('docker_card_abnormal')">
                            [{{ t('common.solution') }}]
                          </a>
                        </div>
                      </STip>
                    </span>
                  </NFlex>
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
      :options="getDropdownMenuOptions()" :show="dropdownShow" :on-clickoutside="onClickoutside"
      @select="handleRightMenuSelect"
    />

    <!-- 弹窗 -->
    <NModal
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
    </NModal>
  </div>
</template>

<style scoped>
.dot {
  position: absolute;
  color: white;
  border-radius: 50%;
}

.dot-small {
  top: 6px;
  right: 6px;
  width: 10px;
  height: 10px;
  font-size: 10px;
}

.dot-info {
  top: 8px;
  right: 8px;
  width: 10px;
  height: 10px;
}
</style>
