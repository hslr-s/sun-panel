<script setup lang="ts">
import { nextTick, ref } from 'vue'
import { NDropdown, NModal } from 'naive-ui'
import StandardIcon from './common/StandardIcon.vue'
import type { PanelPanelConfigStyleEnum } from '@/enums'
import { t } from '@/locales'
import { useAuthStore } from '@/store'
import { AppLoader } from '@/components/common/'

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
const authStore = useAuthStore()
const item = ref(props.itemInfo)
const builtinAppModalShow = ref(false)

const dropdownMenuX = ref(0)
const dropdownMenuY = ref(0)
const dropdownShow = ref(false)

// 点击项目
function handleItemClick() {
  if (props.sortStatus) {
    emit('edit')
    return
  }

  console.log('点击了...')
  builtinAppModalShow.value = true
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

  switch (key) {
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
  const dropdownMenuOptions = []

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
    />

    <!-- 右键菜单 -->
    <NDropdown
      placement="bottom-start" trigger="manual" :x="dropdownMenuX" :y="dropdownMenuY"
      :options="getDropdownMenuOptions()" :show="dropdownShow" :on-clickoutside="onClickoutside" @select="handleRightMenuSelect"
    />

    <NModal
      v-model:show="builtinAppModalShow"
      :title="item.title"
      preset="card"
      style="max-width: 1000px;height: 600px; border-radius: 1rem;"
      :bordered="true"
      content-style="height: calc(100% - 50px);display:flex;"
      size="small"
    >
      <div class="w-full rounded-2xl border overflow-auto dark:border-zinc-700" style="">
        <AppLoader :component-name="item.expandParam.systemAppName" />
      </div>
    </NModal>
  </div>
</template>
