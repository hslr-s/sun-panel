<script setup lang="ts">
import { ref } from 'vue'
import { VueDraggable } from 'vue-draggable-plus'
import AppIconSystemMonitor from './AppIconSystemMonitor/index.vue'
import { MonitorType } from './typings'
import type { CardStyle, MonitorData, ProgressStyle } from './typings'
import { useAuthStore, usePanelState } from '@/store'
import { PanelPanelConfigStyleEnum } from '@/enums'

const panelState = usePanelState()
const authStore = useAuthStore()
const systemMonitorData = ref<SystemMonitor.GetAllRes | null>(null)
const progressStyle = ref<ProgressStyle>({
  color: 'white',
  railColor: 'rgba(0, 0, 0, 0.2)',
  height: 5,
})

const cardStyle: CardStyle = {
  background: '#2a2a2a6b',
}

const monitorDatas = ref<MonitorData[]>([
  {
    monitorType: MonitorType.cpu,
    extendParam: {
      progressStyle,
    },
    description: 'CPU状态',
    cardStyle,
  },
  {
    monitorType: MonitorType.memory,
    cardStyle,
    extendParam: {
      progressStyle,
    },
  },
  {
    monitorType: MonitorType.disk,
    extendParam: {
      progressStyle,
      path: 'C:',
    },
    cardStyle,
  },
])
</script>

<template>
  <div class="w-full">
    <!-- 分组标题 -->
    <div class="text-white text-xl font-extrabold mb-[20px] ml-[10px] flex items-center">
      <span class="text-shadow">
        系统状态
      </span>
      <!-- <div
        v-if="authStore.visitMode === VisitMode.VISIT_MODE_LOGIN"
        class="ml-2 delay-100 transition-opacity flex"
        :class="itemGroup.hoverStatus ? 'opacity-100' : 'opacity-0'"
      >
        <span class="mr-2 cursor-pointer" title="添加快捷图标" @click="handleAddItem(itemGroup.id)">
          <SvgIcon class="text-white font-xl" icon="typcn:plus" />
        </span>
        <span class="mr-2 cursor-pointer " title="排序组快捷图标" @click="handleSetSortStatus(itemGroupIndex, !itemGroup.sortStatus)">
          <SvgIcon class="text-white font-xl" icon="ri:drag-drop-line" />
        </span>
      </div> -->
    </div>

    <!-- 详情图标 -->
    <!-- <div v-if="panelState.panelConfig.iconStyle === PanelPanelConfigStyleEnum.info"> -->
    <div v-if="panelState.panelConfig.iconStyle === PanelPanelConfigStyleEnum.info">
      <VueDraggable
        v-model="monitorDatas" item-key="sort" :animation="300"
        class="icon-info-box"
        filter=".not-drag"
      >
        <div v-for="item, index in monitorDatas" :key="index" :title="item.description" @contextmenu="(e) => handleContextMenu(e, itemGroupIndex, item)">
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
        >
          <div v-for="item, index in monitorDatas" :key="index" :title="item.description" @contextmenu="(e) => handleContextMenu(e, itemGroupIndex, item)">
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
    <!-- <div v-if="itemGroup.sortStatus" class="flex mt-[10px]">
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
      </div> -->
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
