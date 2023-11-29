<script setup lang="ts">
import { VueDraggable } from 'vue-draggable-plus'
import { NBackTop, NButton, NButtonGroup, NDropdown, NModal, NSkeleton, NSpin, useDialog, useMessage } from 'naive-ui'
import { nextTick, onMounted, ref, watch } from 'vue'
import { AppIcon, EditItem, Setting } from './components'
import { Clock, SearchBox } from '@/components/deskModule'
import { SvgIcon } from '@/components/common'
import { deletes, getListByGroupId, saveSort } from '@/api/panel/itemIcon'
import { getList as getGroupList } from '@/api/panel/itemIconGroup'

import { getInfo } from '@/api/system/user'
import { usePanelState, useUserStore } from '@/store'
import { PanelPanelConfigStyleEnum, PanelStateNetworkModeEnum } from '@/enums'
import { setTitle } from '@/utils/cmn'

interface StateDragAppSort {
  status: boolean
}
interface ItemGroup extends Panel.ItemIconGroup {
  items?: Panel.ItemInfo[]
}

const stateDragAppSort = ref<StateDragAppSort>({
  status: false,
})

const ms = useMessage()
const dialog = useDialog()
const panelState = usePanelState()
const userStore = useUserStore()

const scrollContainerRef = ref<HTMLElement | undefined>(undefined)

const editItemInfoShow = ref<boolean>(false)
const editItemInfoData = ref<Panel.ItemInfo | null>(null)
const windowShow = ref<boolean>(false)
const windowSrc = ref<string>('')
const windowTitle = ref<string>('')

const windowIframeRef = ref(null)
const windowIframeIsLoad = ref<boolean>(false)

const dropdownMenuX = ref(0)
const dropdownMenuY = ref(0)
const dropdownShow = ref(false)
const currentRightSelectItem = ref<Panel.ItemInfo | null>(null)

const settingModalShow = ref(false)

const items = ref<ItemGroup[]>([])
const filterItems = ref<ItemGroup[]>([])

function handleAddAppClick() {
  editItemInfoData.value = null
  editItemInfoShow.value = true
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

function handleItemClick(item: Panel.ItemInfo) {
  let jumpUrl = ''

  if (item)
    jumpUrl = (panelState.networkMode === PanelStateNetworkModeEnum.lan ? item.lanUrl : item.url) as string
  if (item.lanUrl === '')
    jumpUrl = item.url

  openPage(item.openMethod, jumpUrl, item.title)
}

function handWindowIframeIdLoad(payload: Event) {
  windowIframeIsLoad.value = false
}

function getList() {
  // 获取组数据
  getGroupList<Common.ListResponse<ItemGroup[]>>().then(({ code, data, msg }) => {
    if (code === 0)
      items.value = data.list
    for (let i = 0; i < data.list.length; i++) {
      const element = data.list[i]
      getListByGroupId<Common.ListResponse<Panel.ItemInfo[]>>(element.id).then((res) => {
        if (res.code === 0)
          items.value[i].items = res.data.list
      })
    }
		filterItems.value = items.value
    // console.log(items)
  })
}

function handleSelect(key: string | number) {
  dropdownShow.value = false
  // console.log(currentRightSelectItem, key)
  let jumpUrl = panelState.networkMode === PanelStateNetworkModeEnum.lan ? currentRightSelectItem.value?.lanUrl : currentRightSelectItem.value?.url
  if (currentRightSelectItem.value?.lanUrl === '')
    jumpUrl = currentRightSelectItem.value.url
  switch (key) {
    case 'newWindows':
      window.open(jumpUrl)
      break
    case 'openWanUrl':
      if (currentRightSelectItem.value)
        openPage(currentRightSelectItem.value?.openMethod, currentRightSelectItem.value?.url, currentRightSelectItem.value?.title)
      break
    case 'openLanUrl':
      if (currentRightSelectItem.value && currentRightSelectItem.value.lanUrl)
        openPage(currentRightSelectItem.value?.openMethod, currentRightSelectItem.value.lanUrl, currentRightSelectItem.value?.title)
      break
    case 'edit':
      // 这里有个奇怪的问题，如果不使用{...}的方式 父组件的值会同步修改 标记一下
      editItemInfoData.value = { ...currentRightSelectItem.value } as Panel.ItemInfo
      editItemInfoShow.value = true
      break
    case 'delete':
      dialog.warning({
        title: '警告',
        content: `你确定要删除图标 ${currentRightSelectItem.value?.title} ？`,
        positiveText: '确定',
        negativeText: '取消',
        onPositiveClick: () => {
          deletes([currentRightSelectItem.value?.id as number]).then(({ code, msg }) => {
            if (code === 0) {
              ms.success('已删除')
              getList()
            }
            else {
              ms.error(`删除失败：${msg}`)
            }
          })
        },
      })

      break
    default:
      break
  }
}

function handleContextMenu(e: MouseEvent, item: Panel.ItemInfo) {
  e.preventDefault()
  currentRightSelectItem.value = item
  dropdownShow.value = false
  nextTick().then(() => {
    dropdownShow.value = true
    dropdownMenuX.value = e.clientX
    dropdownMenuY.value = e.clientY
  })
}

function onClickoutside() {
  // message.info('clickoutside')
  dropdownShow.value = false
}

function handleEditSuccess(item: Panel.ItemInfo) {
  getList()
}

function handleChangeNetwork(mode: PanelStateNetworkModeEnum) {
  panelState.setNetworkMode(mode)
  if (mode === PanelStateNetworkModeEnum.lan)
    ms.success('已经切换成局域网模式(此配置仅保存在本地)')

  else
    ms.success('已经切换成互联网模式(此配置仅保存在本地)')
}

// 结束拖拽
function handleEndDrag(event: any, itemIconGroup: Panel.ItemIconGroup) {
  // console.log(event)
  // console.log(items.value)
}

function handleSaveSort(itemGroup: ItemGroup) {
  const saveItems: Common.SortItemRequest[] = []
  if (itemGroup.items) {
    for (let i = 0; i < itemGroup.items.length; i++) {
      const element = itemGroup.items[i]
      saveItems.push({
        id: element.id as number,
        sort: i + 1,
      })
    }

    saveSort({ itemIconGroupId: itemGroup.id as number, sortItems: saveItems }).then(({ code, msg }) => {
      if (code === 0) {
        //
        ms.success('保存成功')
        // sortStatus.value = false
      }
      else {
        ms.error(`保存失败:${msg}`)
      }
    })
  }
}

function getDropdownMenuOptions() {
  const dropdownMenuOptions = [
    {
      label: '新窗口打开',
      key: 'newWindows',
    },

  ]

  if (currentRightSelectItem.value?.lanUrl && panelState.networkMode === PanelStateNetworkModeEnum.wan) {
    dropdownMenuOptions.push({
      label: '打开局域网地址',
      key: 'openLanUrl',
    })
  }

  if (currentRightSelectItem.value?.lanUrl && panelState.networkMode === PanelStateNetworkModeEnum.lan) {
    dropdownMenuOptions.push({
      label: '打开互联网地址',
      key: 'openWanUrl',
    })
  }

  dropdownMenuOptions.push({
    label: '编辑',
    key: 'edit',
  }, {
    label: '删除',
    key: 'delete',
  })

  return dropdownMenuOptions
}

watch(() => stateDragAppSort.value.status, (newvalue: boolean) => {
  if (newvalue === false)
    getList()
  else
    // 开始排序咯,禁用前端搜索功能
    filterItems.value = items.value
    ms.warning('进入排序模式，记得点击保存再退出')
})

onMounted(() => {
  getList()

  // 获取用户信息
  getInfo<User.Info>().then((res) => {
    if (res.code === 0)
      userStore.updateUserInfo(res.data)
  })

  // 更新同步云端配置
  panelState.updatePanelConfigByCloud()

  // 设置标题
  if (panelState.panelConfig.logoText)
    setTitle(panelState.panelConfig.logoText)
})

// 前端搜索过滤
function itemFrontEndSearch(keyword?: string) {
  if (stateDragAppSort.value.status) {
    //排序禁用搜索
    return
  }
  keyword = keyword?.trim()
  if (keyword !== '') {
    const filteredData = ref<ItemGroup[]>([])
    for (let i = 0; i < items.value.length; i++) {
      const element = items.value[i].items?.filter((item: Panel.ItemInfo) => {
        return (
          item.title.toLowerCase().includes(keyword?.toLowerCase() ?? '')
          || item.url.toLowerCase().includes(keyword?.toLowerCase() ?? '')
          || item.description?.toLowerCase().includes(keyword?.toLowerCase() ?? '')
        )
      })
      if (element && element.length > 0){
				filteredData.value.push({ items: element })
      }
    }
    filterItems.value = filteredData.value
  } else {
    filterItems.value = items.value
  }
}

</script>

<template>
  <div class="w-full h-full sun-main ">
    <div
      class="cover" :style="{
        filter: `blur(${panelState.panelConfig.backgroundBlur}px)`,
        background: `url(${panelState.panelConfig.backgroundImageSrc}) no-repeat`,
        backgroundSize: 'cover',
        backgroundPosition: 'center',
      }"
    />
    <div class="mask" :style="{ backgroundColor: `rgba(0,0,0,${panelState.panelConfig.backgroundMaskNumber})` }" />
    <div ref="scrollContainerRef" class="absolute w-full h-full overflow-auto">
      <div class="p-2.5 max-w-[1200px] mx-auto mt-[10%]">
        <!-- 头 -->
        <div class="mx-[auto] w-[80%]">
          <div class="flex mx-[auto] items-center justify-center text-white">
            <div>
              <span class="text-2xl md:text-5xl font-bold text-shadow">
                {{ panelState.panelConfig.logoText }}
              </span>
            </div>
            <div class="text-base lg:text-2xl mx-[10px]">
              |
            </div>
            <div class="text-shadow">
              <Clock :hide-second="!panelState.panelConfig.clockShowSecond" />
            </div>
          </div>
          <div v-if="panelState.panelConfig.searchBoxShow" class="flex mt-[20px] mx-auto sm:w-full lg:w-[80%]">
            <SearchBox  @itemSearch="itemFrontEndSearch"/>
          </div>
        </div>

        <!-- 应用盒子 -->
        <div class="mt-[50px]">
          <!-- 组纵向排列 -->
          <div
            v-for="(itemGroup, itemGroupIndex) in filterItems"
            :key="itemGroupIndex"
            class="mt-[50px]"
            :class="stateDragAppSort.status ? 'shadow-2xl border shadow-[0_0_30px_10px_rgba(0,0,0,0.8)]  p-[10px] rounded-2xl' : ''"
          >
            <!-- 分组标题 -->
            <div class="text-white text-xl font-extrabold mb-[20px] ml-[10px]">
              {{ itemGroup.title }}
            </div>

            <!-- 详情图标 -->
            <div v-if="panelState.panelConfig.iconStyle === PanelPanelConfigStyleEnum.info">
              <div v-if="itemGroup.items">
                <VueDraggable
                  v-model="itemGroup.items" item-key="sort" :animation="300"
                  class="mx-auto mt-4 grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 xl:12 gap-5"
                  filter=".not-drag"
                  :disabled="!stateDragAppSort.status"
                  @end="(event) => handleEndDrag(event, itemGroup)"
                >
                  <div v-for="item, index in itemGroup.items" :key="index" :title="item.description" @contextmenu="(e) => handleContextMenu(e, item)">
                    <AppIcon
                      :class="stateDragAppSort.status ? 'cursor-move' : 'cursor-pointer'"
                      :item-info="item"
                      :icon-text-color="panelState.panelConfig.iconTextColor"
                      :icon-text-info-hide-description="panelState.panelConfig.iconTextInfoHideDescription || false"
                      :style="0"
                      @click="handleItemClick(item)"
                    />
                  </div>

                  <div v-if="itemGroup.items.length === 0" class="not-drag">
                    <AppIcon
                      :class="stateDragAppSort.status ? 'cursor-move' : 'cursor-pointer'"
                      :item-info="{ icon: { itemType: 3, text: 'subway:add' }, title: '添加图标', url: '', openMethod: 0 }"
                      :icon-text-color="panelState.panelConfig.iconTextColor"
                      :icon-text-info-hide-description="panelState.panelConfig.iconTextInfoHideDescription || false"
                      :style="0"
                      @click="handleAddAppClick"
                    />
                  </div>
                </VueDraggable>
              </div>
            </div>

            <!-- APP图标宫型盒子 -->
            <div v-if="panelState.panelConfig.iconStyle === PanelPanelConfigStyleEnum.icon">
              <div v-if="itemGroup.items">
                <VueDraggable
                  v-model="itemGroup.items" item-key="id" :animation="300"
                  class="mx-auto mt-4 grid grid-cols-4 sm:grid-cols-6 md:grid-cols-8 lg:grid-cols-10 xl:12 gap-5"

                  filter=".not-drag"
                  :disabled="!stateDragAppSort.status"
                >
                  <div v-for="item, index in itemGroup.items" :key="index" :title="item.description" @contextmenu="(e) => handleContextMenu(e, item)">
                    <AppIcon
                      :class="stateDragAppSort.status ? 'cursor-move' : 'cursor-pointer'"
                      :item-info="item"
                      :icon-text-color="panelState.panelConfig.iconTextColor"
                      :icon-text-info-hide-description="!panelState.panelConfig.iconTextInfoHideDescription"
                      :style="1"
                      @click="handleItemClick(item)"
                    />
                  </div>

                  <div v-if="itemGroup.items.length === 0" class="not-drag">
                    <AppIcon
                      :class="stateDragAppSort.status ? 'cursor-move' : 'cursor-pointer'"
                      :item-info="{ icon: { itemType: 3, text: 'subway:add' }, title: '添加图标', url: '', openMethod: 0 }"
                      :icon-text-color="panelState.panelConfig.iconTextColor"
                      :icon-text-info-hide-description="!panelState.panelConfig.iconTextInfoHideDescription"
                      :style="1"
                      @click="handleAddAppClick"
                    />
                  </div>
                </vuedraggable>
              </div>
            </div>

            <!-- 编辑栏 -->
            <div v-if="stateDragAppSort.status" class="flex mt-[10px]">
              <div>
                <NButton color="#2a2a2a6b" @click="handleSaveSort(itemGroup)">
                  <template #icon>
                    <SvgIcon class="text-white font-xl" icon="material-symbols:save" />
                  </template>
                  <div>
                    保存排序
                  </div>
                </NButton>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 右键菜单 -->
    <NDropdown
      placement="bottom-start" trigger="manual" :x="dropdownMenuX" :y="dropdownMenuY"
      :options="getDropdownMenuOptions()" :show="dropdownShow" :on-clickoutside="onClickoutside" @select="handleSelect"
    />

    <!-- 悬浮按钮 -->
    <div class="fixed-element  shadow-[0_0_10px_2px_rgba(0,0,0,0.2)]">
      <NButton v-if="stateDragAppSort.status" color="#2a2a2a6b" @click="stateDragAppSort.status = !stateDragAppSort.status">
        <template #icon>
          <SvgIcon class="text-white font-xl" icon="ri:drag-drop-line" />
        </template>
      </NButton>
      <NButtonGroup v-if="!stateDragAppSort.status" vertical>
        <NButton color="#2a2a2a6b" @click="handleAddAppClick">
          <template #icon>
            <SvgIcon class="text-white font-xl" icon="typcn:plus" />
          </template>
        </NButton>

        <NButton
          v-if="panelState.networkMode === PanelStateNetworkModeEnum.lan" color="#2a2a2a6b"
          title="当前:局域网模式，点击切换成互联网模式" @click="handleChangeNetwork(PanelStateNetworkModeEnum.wan)"
        >
          <template #icon>
            <SvgIcon class="text-white font-xl" icon="material-symbols:lan-outline" />
          </template>
        </NButton>

        <NButton
          v-if="panelState.networkMode === PanelStateNetworkModeEnum.wan" color="#2a2a2a6b"
          title="当前:互联网模式，点击切换成局域网模式" @click="handleChangeNetwork(PanelStateNetworkModeEnum.lan)"
        >
          <template #icon>
            <SvgIcon class="text-white font-xl" icon="mdi:wan" />
          </template>
        </NButton>

        <NButton color="#2a2a2a6b" title="排序模式" @click="stateDragAppSort.status = !stateDragAppSort.status">
          <template #icon>
            <SvgIcon class="text-white font-xl" icon="ri:drag-drop-line" />
          </template>
        </NButton>

        <NButton color="#2a2a2a6b" @click="settingModalShow = !settingModalShow">
          <template #icon>
            <SvgIcon class="text-white font-xl" icon="ep:setting" />
          </template>
        </NButton>
      </NButtonGroup>

      <NBackTop
        :listen-to="() => scrollContainerRef"
        :right="10"
        :bottom="10"
        style="background-color:transparent;border: none;box-shadow: none;"
      >
        <div class="shadow-[0_0_10px_2px_rgba(0,0,0,0.2)]">
          <NButton color="#2a2a2a6b">
            <template #icon>
              <SvgIcon class="text-white font-xl" icon="icon-park-outline:to-top" />
            </template>
          </NButton>
        </div>
      </NBackTop>

      <Setting v-model:visible="settingModalShow" />
    </div>

    <EditItem v-model:visible="editItemInfoShow" :item-info="editItemInfoData" @done="handleEditSuccess" />

    <!-- 弹窗 -->
    <NModal
      v-model:show="windowShow" :mask-closable="false" preset="card"
      style="max-width: 1000px;height: 600px;border-radius: 1rem;" :bordered="false" size="small" role="dialog"
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
      <div class="w-full h-full rounded-2xl overflow-hidden border">
        <NSkeleton v-if="windowIframeIsLoad" height="100%" width="100%" />
        <iframe
          v-show="!windowIframeIsLoad" id="windowIframeId" ref="windowIframeRef" :src="windowSrc"
          class="w-full h-full" frameborder="0" @load="handWindowIframeIdLoad"
        />
      </div>
    </NModal>
  </div>
</template>

<style>
body,
html {
  overflow: hidden;
  background-color: rgb(54, 54, 54);
}
</style>

<style scoped>
.mask {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
}

.sun-main {
  user-select: none;
}

.cover {
  position: absolute;
  width: 100%;
  height: 100%;
  overflow: hidden;
  /* background: url(@/assets/start_sky.jpg) no-repeat; */

  transform: scale(1.05);
}

.text-shadow {
  text-shadow: 2px 2px 50px rgb(0, 0, 0);
}

.app-icon-text-shadow {
  text-shadow: 2px 2px 5px rgb(0, 0, 0);
}

.fixed-element {
  position: fixed;
  /* 将元素固定在屏幕上 */
  right: 10px;
  /* 距离屏幕顶部的距离 */
  bottom: 50px;
  /* 距离屏幕左侧的距离 */
}
</style>
