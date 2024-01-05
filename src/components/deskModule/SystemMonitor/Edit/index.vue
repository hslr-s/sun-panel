<script setup lang="ts">
import { computed, defineEmits, defineProps, ref } from 'vue'
import { NButton, NModal, NTabPane, NTabs, useMessage } from 'naive-ui'
import type { GenericProgressStyleExtendParam, MonitorData } from '../typings'
import GenericProgressStyleEditor from './GenericProgressStyleEditor/index.vue'

interface Props {
  visible: boolean
  monitorData: MonitorData | null
}

const props = defineProps<Props>()
const emit = defineEmits<Emit>()

// 默认通用的进度扩展参数
const defaultGenericProgressStyleExtendParam: GenericProgressStyleExtendParam = {
  progressColor: 'white',
  progressRailColor: 'white',
  color: 'white',
  backgroundColor: '#2a2a2a6b',
}

// const currentData = ref<MonitorData>(props.monitorData)

const currentDataExterdParam = ref<GenericProgressStyleExtendParam>(defaultGenericProgressStyleExtendParam)

const ms = useMessage()
const submitLoading = ref(false)

interface Emit {
  (e: 'update:visible', visible: boolean): void
  (e: 'done', item: MonitorData): void// 创建完成
}

// 更新值父组件传来的值
const show = computed({
  get: () => props.visible,
  set: (visible: boolean) => {
    emit('update:visible', visible)
  },
})
</script>

<template>
  <NModal v-model:show="show" preset="card" size="small" style="width: 600px;border-radius: 1rem;" :title="monitorData ? '修改项目' : '添加项目'">
    <!-- 选择监视器 -->
    <div>
      progressStyle值：{{ JSON.stringify(currentDataExterdParam) }}
    </div>
    <NTabs type="segment">
      <NTabPane name="chap1" tab="CPU">
        <GenericProgressStyleEditor v-model:genericProgressStyleExtendParam="currentDataExterdParam" />
      </NTabPane>
      <NTabPane name="chap2" tab="内存">
        <!--  -->
      </NTabPane>
      <NTabPane name="chap3" tab="磁盘">
        <!--  -->
      </NTabPane>
    </NTabs>

    <template #footer>
      <NButton type="success" :loading="submitLoading" style="float: right;" @click="handleValidateButtonClick">
        确定
      </NButton>
    </template>
  </NModal>
</template>
