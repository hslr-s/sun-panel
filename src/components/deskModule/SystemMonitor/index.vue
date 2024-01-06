<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { VueDraggable } from 'vue-draggable-plus'
import { NButton } from 'naive-ui'
import AppIconSystemMonitor from './AppIconSystemMonitor/index.vue'
import type { CardStyle, MonitorData } from './typings'
import Edit from './Edit/index.vue'
import { getAll } from './common'
import { useAuthStore, usePanelState } from '@/store'
import { PanelPanelConfigStyleEnum } from '@/enums'
import { VisitMode } from '@/enums/auth'
import { SvgIcon } from '@/components/common'

interface MonitorGroup extends Panel.ItemIconGroup {
  sortStatus?: boolean
  hoverStatus: boolean
  items?: Panel.ItemInfo[]
}
const panelState = usePanelState()
const authStore = useAuthStore()
// const systemMonitorData = ref<SystemMonitor.GetAllRes | null>(null)
const monitorGroup = ref<MonitorGroup>({
  hoverStatus: false,
  sortStatus: false,
})
// const progressStyle = ref<ProgressStyle>({
//   color: 'white',
//   railColor: 'rgba(0, 0, 0, 0.2)',
//   height: 5,
// })

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

  // // 并未保存排序重新更新数据
  // if (!sortStatus) {
  //   // 单独更新组
  //   if (items.value[groupIndex] && items.value[groupIndex].id)
  //     updateItemIconGroupByNet(groupIndex, items.value[groupIndex].id as number)
  // }
}

function handleSetHoverStatus(hoverStatus: boolean) {
  monitorGroup.value.hoverStatus = hoverStatus
}

const cardStyle: CardStyle = {
  background: '#2a2a2a6b',
}

const monitorDatas = ref<MonitorData[]>([])

function handleClick(index: number, item: MonitorData) {
  editShowStatus.value = true

  editData.value = item
  editIndex.value = index

  console.log(editData.value)
}

async function getData() {
  monitorDatas.value = await getAll()
}

onMounted(() => {
  getData()
})

function handleSaveDone() {
  getData()
}

function handleSaveSort() {

}
</script>

<template>
  <div class="w-full">
    <div
      class="mt-[50px]"
      :class="monitorGroup.sortStatus ? 'shadow-2xl border shadow-[0_0_30px_10px_rgba(0,0,0,0.3)]  p-[10px] rounded-2xl' : ''"
      @mouseenter="handleSetHoverStatus(true)"
      @mouseleave="handleSetHoverStatus(false)"
    >
      <!-- 分组标题 -->
      <div class="text-white text-xl font-extrabold mb-[20px] ml-[10px] flex items-center">
        <span class="text-shadow">
          系统状态
        </span>
        <div
          v-if="authStore.visitMode === VisitMode.VISIT_MODE_LOGIN"
          class="ml-2 delay-100 transition-opacity flex"
          :class="monitorGroup.hoverStatus ? 'opacity-100' : 'opacity-0'"
        >
          <span class="mr-2 cursor-pointer" title="添加快捷图标" @click="handleAddItem()">
            <SvgIcon class="text-white font-xl" icon="typcn:plus" />
          </span>
          <span class="mr-2 cursor-pointer " title="排序组快捷图标" @click="handleSetSortStatus(!monitorGroup.sortStatus)">
            <SvgIcon class="text-white font-xl" icon="ri:drag-drop-line" />
          </span>
        </div>
      </div>

      <!-- 详情图标 -->
      <div v-if="panelState.panelConfig.iconStyle === PanelPanelConfigStyleEnum.info">
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
      </div>

      <!-- APP图标宫型盒子 -->
      <div v-if="panelState.panelConfig.iconStyle === PanelPanelConfigStyleEnum.icon">
        <div v-if="monitorDatas">
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
        </div>
      </div>

      <!-- 编辑栏 -->
      <div v-if="monitorGroup.sortStatus" class="flex mt-[10px]">
        <div>
          <NButton color="#2a2a2a6b" @click="handleSaveSort()">
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

    <Edit v-model:visible="editShowStatus" :monitor-data="editData" :index="editIndex" @done="handleSaveDone" />
  </div>
</template>

<style>
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
