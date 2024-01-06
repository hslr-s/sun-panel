<script setup lang="ts">
import { computed, defineEmits, defineProps, ref, watch } from 'vue'
import { NButton, NModal, NTabPane, NTabs, useMessage } from 'naive-ui'
import { MonitorType } from '../typings'
import type { DiskExtendParam, GenericProgressStyleExtendParam, MonitorData } from '../typings'
import { add, saveByIndex } from '../common'

import GenericProgressStyleEditor from './GenericProgressStyleEditor/index.vue'
import DiskEditor from './DiskEditor/index.vue'

interface Props {
  visible: boolean
  monitorData: MonitorData | null
  index: number | null
}

const props = defineProps<Props>()
const emit = defineEmits<Emit>()

// 默认通用的进度扩展参数
const defaultGenericProgressStyleExtendParam: GenericProgressStyleExtendParam = {
  progressColor: '#fff',
  progressRailColor: '#CFCFCFA8',
  color: '#fff',
  backgroundColor: '#2a2a2a6b',
}

const defaultDiskExtendParam: DiskExtendParam = {
  progressColor: '#fff',
  progressRailColor: '#CFCFCFA8',
  color: '#fff',
  backgroundColor: '#2a2a2a6b',
  path: '',
}

const defaultMonitorData: MonitorData = {
  extendParam: defaultGenericProgressStyleExtendParam,
  monitorType: MonitorType.cpu,
}

const active = ref<string>(MonitorType.cpu)
const currentMonitorData = ref<MonitorData>(props.monitorData || { ...defaultMonitorData })
const currentGenericProgressStyleExtendParam = ref<GenericProgressStyleExtendParam>({ ...defaultGenericProgressStyleExtendParam })
const currentDiskExtendParam = ref<DiskExtendParam>({ ...defaultDiskExtendParam })

const ms = useMessage()
const submitLoading = ref(false)

interface Emit {
  (e: 'update:visible', visible: boolean): void
  (e: 'done', item: boolean): void
}

// 更新值父组件传来的值
const show = computed({
  get: () => props.visible,
  set: (visible: boolean) => {
    emit('update:visible', visible)
  },
})

watch(() => props.visible, (value) => {
  active.value = props.monitorData?.monitorType || MonitorType.cpu
  if (props.monitorData?.monitorType === MonitorType.cpu || props.monitorData?.monitorType === MonitorType.memory)
    currentGenericProgressStyleExtendParam.value = { ...props.monitorData?.extendParam }
  else if (props.monitorData?.monitorType === MonitorType.disk)
    currentDiskExtendParam.value = { ...props.monitorData?.extendParam }

  if (!value)
    handleResetExtendParam()
})

function handleResetExtendParam() {
  currentGenericProgressStyleExtendParam.value = { ...defaultGenericProgressStyleExtendParam }
  currentDiskExtendParam.value = { ...defaultDiskExtendParam }
}

// 保存提交
async function handleSubmit() {
  currentMonitorData.value.monitorType = active.value as MonitorType
  if (currentMonitorData.value.monitorType === MonitorType.cpu || currentMonitorData.value.monitorType === MonitorType.memory)
    currentMonitorData.value.extendParam = currentGenericProgressStyleExtendParam

  else if (currentMonitorData.value.monitorType === MonitorType.disk)
    currentMonitorData.value.extendParam = currentDiskExtendParam

  console.log('保存', currentMonitorData.value.extendParam)

  if (props.index !== null) {
    const res = await saveByIndex(props.index, currentMonitorData.value)
    if (res) {
      show.value = false
      emit('done', true)
    }
    else {
      ms.error('保存失败')
    }
  }
  else {
    const res = await add(currentMonitorData.value)
    if (res) {
      show.value = false
      emit('done', true)
    }
    else {
      ms.error('保存失败')
    }
  }
}
</script>

<template>
  <NModal v-model:show="show" preset="card" size="small" style="width: 600px;border-radius: 1rem;" :title="monitorData ? '修改项目' : '添加项目'">
    <!-- 选择监视器 -->
    <!-- <div>
      {{ JSON.stringify(currentGenericProgressStyleExtendParam) }}
      {{ JSON.stringify(currentDiskExtendParam) }}
    </div> -->
    <NTabs v-model:value="active" type="segment">
      <NTabPane :name="MonitorType.cpu" tab="CPU状态">
        <GenericProgressStyleEditor v-model:genericProgressStyleExtendParam="currentGenericProgressStyleExtendParam" />
        <NButton @click="handleResetExtendParam">
          重置
        </NButton>
      </NTabPane>
      <NTabPane :name="MonitorType.memory" tab="内存状态">
        <GenericProgressStyleEditor v-model:genericProgressStyleExtendParam="currentGenericProgressStyleExtendParam" />
        <NButton @click="handleResetExtendParam">
          重置
        </NButton>
      </NTabPane>
      <NTabPane :name="MonitorType.disk" tab="磁盘状态">
        <DiskEditor v-model:disk-extend-param="currentDiskExtendParam" />
        <NButton @click="handleResetExtendParam">
          重置
        </NButton>
      </NTabPane>
    </NTabs>

    <template #footer>
      <NButton type="success" :loading="submitLoading" style="float: right;" @click="handleSubmit">
        确定
      </NButton>
    </template>
  </NModal>
</template>
