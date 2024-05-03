<script setup lang="ts">
import { VueDraggable } from 'vue-draggable-plus'
import { NBackTop, NButton, NButtonGroup, useDialog, useMessage } from 'naive-ui'
import { onMounted, ref } from 'vue'
import { AppStarter, EditItem, ItemCards } from './components'
import { Clock, SearchBox, SystemMonitor } from '@/components/deskModule'
import { SvgIcon } from '@/components/common'
import { deletes, getListByGroupId, saveSort } from '@/api/panel/itemIcon'
import { getList as getGroupList } from '@/api/panel/itemIconGroup'
import { apiRespErrMsg, setTitle, updateLocalUserInfo } from '@/utils/cmn'

import { useAuthStore, usePanelState } from '@/store'
import { PanelCardType, PanelPanelConfigStyleEnum, PanelStateNetworkModeEnum } from '@/enums'
import { VisitMode } from '@/enums/auth'
import { router } from '@/router'
import { t } from '@/locales'

interface ItemGroup extends Panel.ItemGroupInfo {
  sortStatus?: boolean
  hoverStatus: boolean
  items?: Panel.ItemInfo[]
}

const ms = useMessage()
const dialog = useDialog()
const panelState = usePanelState()
const authStore = useAuthStore()

const scrollContainerRef = ref<HTMLElement | undefined>(undefined)

const editItemInfoShow = ref<boolean>(false)
const editItemInfoData = ref<Panel.ItemInfo | null>(null)

// const dropdownShow = ref(false)
const currentAddItenIconGroupId = ref<number | undefined>()

const settingModalShow = ref(false)

const items = ref<ItemGroup[]>([])
const filterItems = ref<ItemGroup[]>([])

function getList() {
  // 获取组数据
  getGroupList<Common.ListResponse<ItemGroup[]>>().then(({ code, data, msg }) => {
    if (code === 0)
      items.value = data.list
    for (let i = 0; i < data.list.length; i++) {
      const element = data.list[i]
      if (element.id)
        updateItemIconGroupByNet(i, element.id)
    }
    filterItems.value = items.value
    // console.log(items)
  }).catch((error) => {
    apiRespErrMsg(error)
  })
}

// 从后端获取组下面的图标
function updateItemIconGroupByNet(itemIconGroupIndex: number, itemIconGroupId: number) {
  getListByGroupId<Common.ListResponse<Panel.ItemInfo[]>>(itemIconGroupId).then((res) => {
    if (res.code === 0)
      items.value[itemIconGroupIndex].items = res.data.list
  }).catch((error) => {
    apiRespErrMsg(error)
  })
}

function handleEditSuccess(item: Panel.ItemInfo) {
  getList()
}

function handleChangeNetwork(mode: PanelStateNetworkModeEnum) {
  panelState.setNetworkMode(mode)
  if (mode === PanelStateNetworkModeEnum.lan)
    ms.success(t('panelHome.changeToLanModelSuccess'))

  else
    ms.success(t('panelHome.changeToWanModelSuccess'))
}

// 删除项目
function deleteItem(itemInfo: Panel.ItemInfo) {
  dialog.warning({
    title: t('common.warning'),
    content: t('common.deleteConfirmByName', { name: itemInfo.title }),
    positiveText: t('common.confirm'),
    negativeText: t('common.cancel'),
    onPositiveClick: () => {
      deletes([itemInfo.id as number]).then(({ code, msg }) => {
        if (code === 0) {
          ms.success(t('common.deleteSuccess'))
          getList()
        }
        else {
          ms.error(`${t('common.deleteFail')}:${msg}`)
        }
      }).catch((error) => {
        apiRespErrMsg(error)
      })
    },
  })
}

// 结束拖拽
// function handleEndDrag(event: any, itemIconGroup: Panel.ItemGroupInfo) {
//   // console.log(event)
//   // console.log(items.value)
// }

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
        ms.success(t('common.saveSuccess'))
        itemGroup.sortStatus = false
      }
      else {
        ms.error(`${t('common.saveFail')}:${msg}`)
      }
    }).catch((error) => {
      apiRespErrMsg(error)
    })
  }
}

onMounted(() => {
  // 更新用户信息
  updateLocalUserInfo()
  getList()

  // 更新同步云端配置
  panelState.updatePanelConfigByCloud()

  // 设置标题
  if (panelState.panelConfig.logoText)
    setTitle(panelState.panelConfig.logoText)
})

// 前端搜索过滤
function itemFrontEndSearch(keyword?: string) {
  keyword = keyword?.trim()
  if (keyword !== '' && panelState.panelConfig.searchBoxSearchIcon) {
    const filteredData = ref<ItemGroup[]>([])
    for (let i = 0; i < items.value.length; i++) {
      const element = items.value[i].items?.filter((item: Panel.ItemInfo) => {
        return (
          item.title.toLowerCase().includes(keyword?.toLowerCase() ?? '')
          || item.url.toLowerCase().includes(keyword?.toLowerCase() ?? '')
          || item.description?.toLowerCase().includes(keyword?.toLowerCase() ?? '')
        )
      })
      if (element && element.length > 0)
        filteredData.value.push({ items: element, hoverStatus: false, cardStyle: items.value[i].cardStyle, title: items.value[i].title })
    }
    filterItems.value = filteredData.value
    console.log('搜索模式开启')
  }
  else {
    console.log('搜索模式关闭')
    filterItems.value = items.value
  }
}

function handleSetHoverStatus(groupIndex: number, hoverStatus: boolean) {
  if (items.value[groupIndex])
    items.value[groupIndex].hoverStatus = hoverStatus
}

function handleSetSortStatus(groupIndex: number, sortStatus: boolean) {
  if (items.value[groupIndex])
    items.value[groupIndex].sortStatus = sortStatus

  // 并未保存排序重新更新数据
  if (!sortStatus) {
    // 单独更新组
    if (items.value[groupIndex] && items.value[groupIndex].id)
      updateItemIconGroupByNet(groupIndex, items.value[groupIndex].id as number)
  }
}

function handleEditItem(item: Panel.ItemInfo) {
  editItemInfoData.value = item
  editItemInfoShow.value = true
  currentAddItenIconGroupId.value = undefined
}

function handleAddItem(itemIconGroupId?: number) {
  editItemInfoData.value = null
  editItemInfoShow.value = true
  if (itemIconGroupId)
    currentAddItenIconGroupId.value = itemIconGroupId
}
</script>

<template>
  <div class="w-full h-full sun-main">
    <div
      class="cover wallpaper" :style="{
        filter: `blur(${panelState.panelConfig.backgroundBlur}px)`,
        background: `url(${panelState.panelConfig.backgroundImageSrc}) no-repeat`,
        backgroundSize: 'cover',
        backgroundPosition: 'center',
      }"
    />
    <div class="mask" :style="{ backgroundColor: `rgba(0,0,0,${panelState.panelConfig.backgroundMaskNumber})` }" />
    <div ref="scrollContainerRef" class="absolute w-full h-full overflow-auto">
      <div
        class="p-2.5 mx-auto"
        :style="{
          marginTop: `${panelState.panelConfig.marginTop}%`,
          marginBottom: `${panelState.panelConfig.marginBottom}%`,
          maxWidth: (panelState.panelConfig.maxWidth ?? '1200') + panelState.panelConfig.maxWidthUnit,
        }"
      >
        <!-- 头 -->
        <div class="mx-[auto] w-[80%]">
          <div class="flex mx-[auto] items-center justify-center text-white">
            <div class="logo">
              <span class="text-2xl md:text-6xl font-bold text-shadow">
                {{ panelState.panelConfig.logoText }}
              </span>
            </div>
            <div class="divider text-base lg:text-2xl mx-[10px]">
              |
            </div>
            <div class="text-shadow">
              <Clock :hide-second="!panelState.panelConfig.clockShowSecond" />
            </div>
          </div>
          <div v-if="panelState.panelConfig.searchBoxShow" class="flex mt-[20px] mx-auto sm:w-full lg:w-[80%]">
            <SearchBox @itemSearch="itemFrontEndSearch" />
          </div>
        </div>

        <!-- 应用盒子 -->
        <div :style="{ marginLeft: `${panelState.panelConfig.marginX}px`, marginRight: `${panelState.panelConfig.marginX}px` }">
          <!-- 系统监控状态 -->
          <div
            v-if="panelState.panelConfig.systemMonitorShow
              && ((panelState.panelConfig.systemMonitorPublicVisitModeShow && authStore.visitMode === VisitMode.VISIT_MODE_PUBLIC)
                || authStore.visitMode === VisitMode.VISIT_MODE_LOGIN)"
            class="flex mx-auto"
          >
            <SystemMonitor
              :allow-edit="authStore.visitMode === VisitMode.VISIT_MODE_LOGIN"
              :show-title="panelState.panelConfig.systemMonitorShowTitle"
            />
          </div>

          <!-- 组纵向排列 -->
          <div
            v-for="(itemGroup, itemGroupIndex) in filterItems" :key="itemGroupIndex"
            class="item-list mt-[50px]"
            :class="itemGroup.sortStatus ? 'shadow-2xl border shadow-[0_0_30px_10px_rgba(0,0,0,0.3)]  p-[10px] rounded-2xl' : ''"
            @mouseenter="handleSetHoverStatus(itemGroupIndex, true)"
            @mouseleave="handleSetHoverStatus(itemGroupIndex, false)"
          >
            <!-- 分组标题 -->
            <div class="text-white text-xl font-extrabold mb-[20px] ml-[10px] flex items-center">
              <span class="group-title text-shadow">
                {{ itemGroup.title }}
              </span>
              <div
                v-if="authStore.visitMode === VisitMode.VISIT_MODE_LOGIN"
                class="group-buttons ml-2 delay-100 transition-opacity flex"
                :class="itemGroup.hoverStatus ? 'opacity-100' : 'opacity-0'"
              >
                <span class="mr-2 cursor-pointer" :title="t('common.add')" @click="handleAddItem(itemGroup.id)">
                  <SvgIcon class="text-white font-xl" icon="typcn:plus" />
                </span>
                <span class="mr-2 cursor-pointer " :title="t('common.sort')" @click="handleSetSortStatus(itemGroupIndex, !itemGroup.sortStatus)">
                  <SvgIcon class="text-white font-xl" icon="ri:drag-drop-line" />
                </span>
              </div>
            </div>

            <!-- 编辑栏 -->
            <div v-if="itemGroup.sortStatus" class="flex my-4">
              <NButton color="#2a2a2a6b" @click="handleSaveSort(itemGroup)">
                <template #icon>
                  <SvgIcon class="text-white font-xl" icon="material-symbols:save" />
                </template>
                <div>
                  {{ $t('common.saveSort') }}
                </div>
              </NButton>
            </div>

            <!-- 详情图标 -->
            <templete v-if="itemGroup.items && (!itemGroup.cardStyle || itemGroup.cardStyle.style === PanelPanelConfigStyleEnum.info)">
              <VueDraggable
                v-model="itemGroup.items" item-key="sort" :animation="300"
                class="icon-info-box"
                filter=".not-drag"
                :disabled="!itemGroup.sortStatus"
              >
                <div v-for="item, index in itemGroup.items" :key="index" :title="item.description">
                  <ItemCards
                    :sort-status="itemGroup.sortStatus || false"
                    :class="itemGroup.sortStatus ? 'cursor-move' : 'cursor-pointer'"
                    :card-info="item"
                    :icon-text-color="itemGroup.cardStyle.textColor"
                    :icon-text-info-hide-description="itemGroup.cardStyle.textInfoHideDescription || false"
                    :icon-text-icon-hide-title="itemGroup.cardStyle.textIconHideTitle || false"
                    :card-style="PanelPanelConfigStyleEnum.info"
                    :card-type="PanelCardType.Normal"
                    @delete="() => deleteItem(item)"
                    @edit="() => handleEditItem(item)"
                  />
                </div>

                <div v-if="itemGroup.items.length === 0" class="not-drag">
                  <ItemCards
                    :sort-status="itemGroup.sortStatus || false"
                    class="cursor-pointer"
                    :card-info="{ icon: { itemType: 3, text: 'subway:add' }, title: t('common.add'), url: '', openMethod: 0, cardType: PanelCardType.Normal }"
                    :icon-text-color="itemGroup.cardStyle.textColor"
                    :icon-text-info-hide-description="itemGroup.cardStyle.textInfoHideDescription || false"
                    :icon-text-icon-hide-title="itemGroup.cardStyle.textIconHideTitle || false"
                    :card-style="PanelPanelConfigStyleEnum.info"
                    :card-type="PanelCardType.Normal"
                    @click="handleAddItem(itemGroup.id)"
                  />
                </div>
              </VueDraggable>
            </templete>

            <!-- APP图标宫型盒子 -->
            <templete v-if="itemGroup.items && itemGroup.cardStyle.style === PanelPanelConfigStyleEnum.icon">
              <VueDraggable
                v-model="itemGroup.items" item-key="sort" :animation="300"
                class="icon-small-box"
                filter=".not-drag"
                :disabled="!itemGroup.sortStatus"
              >
                <div v-for="item, index in itemGroup.items" :key="index" :title="item.description">
                  <ItemCards
                    :sort-status="itemGroup.sortStatus || false"
                    :class="itemGroup.sortStatus ? 'cursor-move' : 'cursor-pointer'"
                    :card-info="item"
                    :icon-text-color="itemGroup.cardStyle.textColor"
                    :icon-text-info-hide-description="itemGroup.cardStyle.textInfoHideDescription || false"
                    :icon-text-icon-hide-title="itemGroup.cardStyle.textIconHideTitle || false"
                    :card-style="PanelPanelConfigStyleEnum.small"
                    :card-type="PanelCardType.Normal"
                    @delete="() => deleteItem(item)"
                    @edit="() => handleEditItem(item)"
                  />
                </div>

                <div v-if="itemGroup.items.length === 0" class="not-drag">
                  <ItemCards
                    :sort-status="itemGroup.sortStatus || false"
                    class="cursor-pointer"
                    :card-info="{ icon: { itemType: 3, text: 'subway:add' }, title: t('common.add'), url: '', openMethod: 0, cardType: PanelCardType.Normal }"
                    :icon-text-color="itemGroup.cardStyle.textColor"
                    :icon-text-info-hide-description="itemGroup.cardStyle.textInfoHideDescription || false"
                    :icon-text-icon-hide-title="itemGroup.cardStyle.textIconHideTitle || false"
                    :card-style="PanelPanelConfigStyleEnum.small"
                    :card-type="PanelCardType.Normal"
                    @click="handleAddItem(itemGroup.id)"
                  />
                </div>
              </vuedraggable>
            </templete>

            <!-- 编辑栏 -->
            <div v-if="itemGroup.sortStatus" class="flex my-4">
              <NButton color="#2a2a2a6b" @click="handleSaveSort(itemGroup)">
                <template #icon>
                  <SvgIcon class="text-white font-xl" icon="material-symbols:save" />
                </template>
                <div>
                  {{ $t('common.saveSort') }}
                </div>
              </NButton>
            </div>
          </div>
        </div>
        <div class="mt-5 footer" v-html="panelState.panelConfig.footerHtml" />
      </div>
    </div>

    <!-- 悬浮按钮 -->
    <div class="fixed-element shadow-[0_0_10px_2px_rgba(0,0,0,0.2)]">
      <NButtonGroup vertical>
        <!-- 网络模式切换按钮组 -->
        <NButton
          v-if="panelState.networkMode === PanelStateNetworkModeEnum.lan && panelState.panelConfig.netModeChangeButtonShow" color="#2a2a2a6b"
          :title="t('panelHome.changeToWanModel')" @click="handleChangeNetwork(PanelStateNetworkModeEnum.wan)"
        >
          <template #icon>
            <SvgIcon class="text-white font-xl" icon="material-symbols:lan-outline-rounded" />
          </template>
        </NButton>

        <NButton
          v-if="panelState.networkMode === PanelStateNetworkModeEnum.wan && panelState.panelConfig.netModeChangeButtonShow" color="#2a2a2a6b"
          :title="t('panelHome.changeToLanModel')" @click="handleChangeNetwork(PanelStateNetworkModeEnum.lan)"
        >
          <template #icon>
            <SvgIcon class="text-white font-xl" icon="mdi:wan" />
          </template>
        </NButton>

        <NButton v-if="authStore.visitMode === VisitMode.VISIT_MODE_LOGIN" color="#2a2a2a6b" @click="settingModalShow = !settingModalShow">
          <template #icon>
            <SvgIcon class="text-white font-xl" icon="majesticons-applications" />
          </template>
        </NButton>

        <NButton v-if="authStore.visitMode === VisitMode.VISIT_MODE_PUBLIC" color="#2a2a2a6b" :title="$t('panelHome.goToLogin')" @click="router.push('/login')">
          <template #icon>
            <SvgIcon class="text-white font-xl" icon="material-symbols:account-circle" />
          </template>
        </NButton>
      </NButtonGroup>

      <AppStarter v-model:visible="settingModalShow" />
      <!-- <Setting v-model:visible="settingModalShow" /> -->
    </div>

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

    <EditItem v-model:visible="editItemInfoShow" :item-info="editItemInfoData" :item-group-id="currentAddItenIconGroupId" @done="handleEditSuccess" />
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

.icon-info-box {
  width: 100%;
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
  gap: 18px;

}

.icon-small-box {
  width: 100%;
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(75px, 1fr));
  gap: 18px;

}

@media (max-width: 500px) {
  .icon-info-box{
    grid-template-columns: repeat(auto-fill, minmax(150px, 1fr));
  }
}
</style>
