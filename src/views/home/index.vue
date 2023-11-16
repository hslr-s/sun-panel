<script setup lang="ts">
import { VueDraggable } from 'vue-draggable-plus'
import { NButton, NButtonGroup, NDropdown, NEllipsis, NModal, NSkeleton, NSpin, useDialog, useMessage } from 'naive-ui'
import { nextTick, onMounted, ref } from 'vue'
import { EditItem, Setting } from './components'
import { Clock } from '@/components/deskModule'
import { ItemIcon, SvgIcon } from '@/components/common'
import { deletes, getListByGroupId } from '@/api/panel/itemIcon'
import { getList as getGroupList } from '@/api/panel/itemIconGroup'

import { getInfo } from '@/api/system/user'
import { usePanelState, useUserStore } from '@/store'
import { PanelStateNetworkModeEnum } from '@/enums'
import { setTitle } from '@/utils/cmn'

const ms = useMessage()
const dialog = useDialog()
const panelState = usePanelState()
const userStore = useUserStore()

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

const dropdownMenuOptions = [
  {
    label: '新窗口打开',
    key: 'newWindows',
  },
  {
    label: '编辑',
    key: 'edit',
  },
  {
    label: '删除',
    key: 'delete',
  },
]

interface ItemGroup extends Panel.ItemIconGroup {
  items?: Panel.ItemInfo[]
}
const items = ref<ItemGroup[]>([])

function handleAddAppClick() {
  editItemInfoData.value = null
  editItemInfoShow.value = true
}

function handleItemClick(item: Panel.ItemInfo) {
  let jumpUrl = ''

  if (item)
    jumpUrl = (panelState.networkMode === PanelStateNetworkModeEnum.lan ? item.lanUrl : item.url) as string
  if (item.lanUrl === '')
    jumpUrl = item.url

  switch (item.openMethod) {
    case 1:
      window.location.href = jumpUrl
      break
    case 2:
      window.open(jumpUrl)
      break
    case 3:
      windowShow.value = true
      windowSrc.value = jumpUrl
      windowTitle.value = item.title
      windowIframeIsLoad.value = true
      break

    default:
      break
  }
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
    console.log(items)
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
    ms.success('已经切换成局域网模式，此时再点击已填写局域网地址的图标将跳转至局域网地址(此配置仅保存在本地)')

  else
    ms.success('已经切换成互联网模式(此配置仅保存在本地)')
}

// 结束拖拽
function handleEndDrag() {
  console.log(items.value)
}

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
    <div class="absolute w-full h-full overflow-auto">
      <div class="p-2.5 max-w-[1200px] mx-auto mt-[10%]">
        <!-- 头 -->
        <div class="mx-[auto] w-[80%]">
          <div class="flex mx-[auto] items-center justify-center text-white">
            <div>
              <span class="text-5xl font-bold text-shadow">
                {{ panelState.panelConfig.logoText }}
              </span>
            </div>
            <div class="text-2xl mx-[10px]">
              |
            </div>
            <div class="text-shadow">
              <Clock :hide-second="!panelState.panelConfig.clockShowSecond" />
            </div>
          </div>
          <!-- <div class="flex mt-[20px] mx-auto w-[80%]">
          <SearchBox />
        </div> -->
        </div>

        <!-- 应用盒子 -->
        <div class="mt-[50px]">
          <!-- 组纵向排列 -->
          <div v-for="(itemGroup, itemGroupIndex) in items" :key="itemGroupIndex" class="mt-[50px]">
            <!-- 分组标题 -->
            <div class="text-white text-xl font-extrabold mb-[20px] ml-[10px]">
              {{ itemGroup.title }}
            </div>

            <!-- 详情图标 -->
            <div v-if="panelState.panelConfig.iconStyle === 0 && itemGroup.items ">
              <VueDraggable
                v-model="itemGroup.items" item-key="sort" :animation="300"
                class="mx-auto mt-4 grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 xl:12 gap-5"
                filter=".not-drag" @end="handleEndDrag"
              >
                <div v-for="item, index in itemGroup.items" :key="index" @contextmenu="(e) => handleContextMenu(e, item)">
                  <div
                    class="w-full rounded-2xl cursor-pointer transition-all duration-200 hover:shadow-[0_0_20px_10px_rgba(0,0,0,0.2)] bg-[#2a2a2a6b] flex"
                    @click="handleItemClick(item)"
                  >
                    <div class="w-[70px]">
                      <ItemIcon :item-icon="item.icon" />
                    </div>
                    <div class="text-white m-[8px_8px_0_8px]" :style="{ color: panelState.panelConfig.iconTextColor }">
                      <div>
                        <NEllipsis>
                          {{ item.title }}
                        </NEllipsis>
                      </div>
                      <div>
                        <NEllipsis :line-clamp="2" class="text-xs">
                          {{ item.description }}
                        </NEllipsis>
                      </div>
                    </div>
                  </div>
                </div>

                <div class="not-drag">
                  <div
                    class="w-full rounded-2xl cursor-pointer transition-all duration-200 hover:shadow-[0_0_20px_10px_rgba(0,0,0,0.2)] bg-[#2a2a2a6b] flex"
                    @click="handleAddAppClick"
                  >
                    <ItemIcon :item-icon="{ itemType: 3, text: 'subway:add', bgColor: '#00000000' }" />
                    <div class="text-white m-[8px]" :style="{ color: panelState.panelConfig.iconTextColor }">
                      <div>
                        <NEllipsis>
                          添加图标
                        </NEllipsis>
                      </div>

                      <div class="text text-xs">
                        <NEllipsis>
                          新增一个新的图标
                        </NEllipsis>
                      </div>
                    </div>
                  </div>
                </div>
              </vuedraggable>
            </div>

            <!-- APP图标宫型盒子 -->
            <div v-if="panelState.panelConfig.iconStyle === 1 && itemGroup.items">
              <VueDraggable
                v-model="itemGroup.items" item-key="id" :animation="300"
                class="mx-auto mt-4 grid grid-cols-4 sm:grid-cols-6 md:grid-cols-8 lg:grid-cols-10 xl:12 gap-5"
                filter=".not-drag"
              >
                <div v-for="item, index in itemGroup.items" :key="index" @contextmenu="(e) => handleContextMenu(e, item)">
                  <div
                    class="sunpanel w-[70px] h-[70px] mx-auto rounded-2xl cursor-pointer transition-all duration-200 hover:shadow-[0_0_20px_10px_rgba(0,0,0,0.2)] bg-[#2a2a2a6b]"
                    @click="handleItemClick(item)"
                  >
                    <ItemIcon :item-icon="item.icon" />
                  </div>
                  <div
                    class="text-center app-icon-text-shadow cursor-pointer mt-[2px]"
                    :style="{ color: panelState.panelConfig.iconTextColor }" @click="handleItemClick(item)"
                  >
                    <span>{{ item.title }}</span>
                  </div>
                </div>

                <div class="not-drag">
                  <div
                    class="w-[70px] h-[70px] mx-auto rounded-2xl cursor-pointer transition-all duration-200 hover:shadow-[0_0_20px_10px_rgba(0,0,0,0.2)]"
                    @click="handleAddAppClick"
                  >
                    <ItemIcon :item-icon="{ itemType: 3, text: 'subway:add', bgColor: '#343434' }" />
                  </div>
                  <div
                    class="text-center app-icon-text-shadow cursor-pointer mt-[2px]"
                    :style="{ color: panelState.panelConfig.iconTextColor }" @click="handleAddAppClick"
                  >
                    添加图标
                  </div>
                </div>
              </vuedraggable>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 右键菜单 -->
    <NDropdown
      placement="bottom-start" trigger="manual" :x="dropdownMenuX" :y="dropdownMenuY"
      :options="dropdownMenuOptions" :show="dropdownShow" :on-clickoutside="onClickoutside" @select="handleSelect"
    />

    <!-- 悬浮按钮 -->
    <div class="fixed-element  shadow-[0_0_10px_2px_rgba(0,0,0,0.2)]">
      <NButtonGroup vertical>
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
        <NButton color="#2a2a2a6b" @click="settingModalShow = !settingModalShow">
          <template #icon>
            <SvgIcon class="text-white font-xl" icon="ep:setting" />
          </template>
        </NButton>
      </NButtonGroup>

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
  right: 30px;
  /* 距离屏幕顶部的距离 */
  bottom: 50px;
  /* 距离屏幕左侧的距离 */
}
</style>
