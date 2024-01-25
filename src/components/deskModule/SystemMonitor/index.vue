<script setup lang="ts">
import { nextTick, onMounted, ref } from 'vue'
import { VueDraggable } from 'vue-draggable-plus'
import { NButton, NDropdown, useDialog, useMessage } from 'naive-ui'
import AppIconSystemMonitor from './AppIconSystemMonitor/index.vue'
import { type CardStyle, type MonitorData, MonitorType } from './typings'
import Edit from './Edit/index.vue'
import { deleteByIndex, getAll, saveAll } from './common'
import { usePanelState } from '@/store'
import { PanelPanelConfigStyleEnum } from '@/enums'
import { SvgIcon } from '@/components/common'
import { t } from '@/locales'

interface MonitorGroup extends Panel.ItemIconGroup {
  sortStatus?: boolean
  hoverStatus: boolean
  items?: Panel.ItemInfo[]
}

const props = defineProps<{
  allowEdit?: boolean
  showTitle?: boolean
}>()
const panelState = usePanelState()

const dialog = useDialog()
const ms = useMessage()

const dropdownMenuX = ref(0)
const dropdownMenuY = ref(0)
const dropdownShow = ref(false)
const currentRightSelectIndex = ref<number | null>(null)

const monitorGroup = ref<MonitorGroup>({
  hoverStatus: false,
  sortStatus: false,
})

const editShowStatus = ref<boolean>(false)
const editData = ref<MonitorData | null>(null)
const editIndex = ref<number | null>(null)

function handleAddItem() {
  editShowStatus.value = true
  editData.value = null
  editIndex.value = null
}

function handleSetSortStatus(sortStatus: boolean) {
  monitorGroup.value.sortStatus = sortStatus

  // 并未保存排序重新更新数据
  if (!sortStatus)
    getData()
}

function handleSetHoverStatus(hoverStatus: boolean) {
  monitorGroup.value.hoverStatus = hoverStatus
}

const cardStyle: CardStyle = {
  background: '#2a2a2a6b',
}

const monitorDatas = ref<MonitorData[]>([])

function handleClick(index: number, item: MonitorData) {
  if (!props.allowEdit)
    return
  editShowStatus.value = true
  editData.value = item
  editIndex.value = index
}

async function getData() {
  monitorDatas.value = await getAll()

  if (monitorDatas.value.length === 0) {
    // 防止空 - 默认数据
    monitorDatas.value.push(
      {
        extendParam: {
          backgroundColor: '#2a2a2a6b',
          color: '#fff',
          progressColor: '#fff',
          progressRailColor: '#CFCFCFA8',
        },
        monitorType: MonitorType.cpu,
      },
    )

    // 生成并保存
    saveAll(monitorDatas.value)
  }
}

onMounted(() => {
  getData()
})

function handleSaveDone() {
  getData()
}

async function handleSaveSort() {
  const { code } = await saveAll(monitorDatas.value)
  if (code === 0)
    monitorGroup.value.sortStatus = false
}

function handleContextMenu(e: MouseEvent, index: number | null, item: MonitorData) {
  if (index !== null) {
    e.preventDefault()
    currentRightSelectIndex.value = index
  }

  nextTick().then(() => {
    dropdownShow.value = true
    dropdownMenuX.value = e.clientX
    dropdownMenuY.value = e.clientY
  })
}

function getDropdownMenuOptions() {
  const dropdownMenuOptions = [
    {
      label: t('common.delete'),
      key: 'delete',
    },
  ]

  return dropdownMenuOptions
}

function onClickoutside() {
  // message.info('clickoutside')
  dropdownShow.value = false
}

async function deleteOneByIndex(index: number) {
  const res = await deleteByIndex(index)
  if (res)
    getData()
}

function handleRightMenuSelect(key: string | number) {
  dropdownShow.value = false

  switch (key) {
    case 'delete':
      dialog.warning({
        title: t('common.warning'),
        content: t('common.deleteConfirm'),
        positiveText: t('common.confirm'),
        negativeText: t('common.cancel'),
        onPositiveClick: () => {
          if (monitorDatas.value.length <= 1) {
            ms.warning(t('common.leastOne'))
            return
          }
          if (currentRightSelectIndex.value !== null)
            deleteOneByIndex(currentRightSelectIndex.value)
        },
      })

      break
    default:
      break
  }
}
</script>

<template>
  <div class="system-monitor w-full">
    <div
      class="mt-[50px]"
      :class="monitorGroup.sortStatus ? 'shadow-2xl border shadow-[0_0_30px_10px_rgba(0,0,0,0.3)]  p-[10px] rounded-2xl' : ''"
      @mouseenter="handleSetHoverStatus(true)"
      @mouseleave="handleSetHoverStatus(false)"
    >
      <!-- 分组标题 -->
      <div class="system-monitor-header text-white text-xl font-extrabold mb-[20px] ml-[10px] flex items-center">
        <span v-if="showTitle" class="text-shadow">
          {{ $t('deskModule.systemMonitor.systemState') }}
        </span>
        <div
          v-if="allowEdit"
          class="system-monitor-buttons ml-2 delay-100 transition-opacity flex"
          :class="monitorGroup.hoverStatus ? 'opacity-100' : 'opacity-0'"
        >
          <span class="mr-2 cursor-pointer" @click="handleAddItem()">
            <SvgIcon class="text-white font-xl" icon="typcn:plus" />
          </span>
          <span class="mr-2 cursor-pointer" @click="handleSetSortStatus(!monitorGroup.sortStatus)">
            <SvgIcon class="text-white font-xl" icon="ri:drag-drop-line" />
          </span>
        </div>
      </div>

      <!-- 详情图标 -->
      <template v-if="panelState.panelConfig.iconStyle === PanelPanelConfigStyleEnum.info">
        <VueDraggable
          v-model="monitorDatas" item-key="sort" :animation="300"
          class="icon-info-box"
          filter=".not-drag"
          :disabled="!monitorGroup.sortStatus"
        >
          <div
            v-for="item, index in monitorDatas" :key="index"
            :title="item.description"
            @click="handleClick(index, item)"
            @contextmenu="(e) => handleContextMenu(e, index, item)"
          >
            <AppIconSystemMonitor
              :extend-param="item.extendParam"
              :icon-text-icon-hide-title="false"
              :card-type-style="panelState.panelConfig.iconStyle"
              :monitor-type="item.monitorType"
              :card-style="cardStyle"
              :icon-text-color="panelState.panelConfig.iconTextColor"
            />
          </div>
        </VueDraggable>
      </template>

      <!-- APP图标宫型盒子 -->
      <template v-if="panelState.panelConfig.iconStyle === PanelPanelConfigStyleEnum.icon">
        <VueDraggable
          v-model="monitorDatas" item-key="sort" :animation="300"
          class="icon-small-box"
          filter=".not-drag"
          :disabled="!monitorGroup.sortStatus"
        >
          <div
            v-for="item, index in monitorDatas" :key="index"
            :title="item.description"
            @click="handleClick(index, item)"
            @contextmenu="(e) => handleContextMenu(e, index, item)"
          >
            <AppIconSystemMonitor
              :extend-param="item.extendParam"
              :icon-text-icon-hide-title="false"
              :card-type-style="panelState.panelConfig.iconStyle"
              :monitor-type="item.monitorType"
              :card-style="cardStyle"
              :icon-text-color="panelState.panelConfig.iconTextColor"
            />
          </div>
        </vuedraggable>
      </template>

      <!-- 编辑栏 -->
      <template v-if="monitorGroup.sortStatus && allowEdit">
        <div class="system-monitor-edit-bar flex mt-[10px]">
          <NButton color="#2a2a2a6b" @click="handleSaveSort()">
            <template #icon>
              <SvgIcon class="text-white font-xl" icon="material-symbols:save" />
            </template>
            <div>
              {{ $t('common.saveSort') }}
            </div>
          </NButton>
        </div>
      </template>
    </div>

    <Edit v-model:visible="editShowStatus" :monitor-data="editData" :index="editIndex" @done="handleSaveDone" />

    <NDropdown
      placement="bottom-start" trigger="manual" :x="dropdownMenuX" :y="dropdownMenuY"
      :options="getDropdownMenuOptions()" :show="dropdownShow" :on-clickoutside="onClickoutside" @select="handleRightMenuSelect"
    />
  </div>
</template>

<style scoped>
.text-shadow {
  text-shadow: 2px 2px 50px rgb(0, 0, 0);
}

.app-icon-text-shadow {
  text-shadow: 2px 2px 5px rgb(0, 0, 0);
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
