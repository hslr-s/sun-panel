<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { NColorPicker, NForm, NFormItem, NSelect } from 'naive-ui'
import type { DiskExtendParam } from '../../typings'
import GenericMonitorCard from '../../components/GenericMonitorCard/index.vue'
import GenericProgress from '../../components/GenericProgress/index.vue'
import { PanelPanelConfigStyleEnum } from '@/enums'
import { getDiskMountpoints } from '@/api/system/systemMonitor'

interface Emit {
  (e: 'update:diskExtendParam', visible: DiskExtendParam): void
}

const props = defineProps<{
  diskExtendParam: DiskExtendParam
}>()
const emit = defineEmits<Emit>()

const mountPointList = ref<{
  label: string
  value: string
}[]>([])

const data = computed({
  get: () => props.diskExtendParam,
  set: (visible) => {
    emit('update:diskExtendParam', visible)
  },
})

async function getMountPointList() {
  try {
    const { data } = await getDiskMountpoints<SystemMonitor.Mountpoint[]>()
    mountPointList.value = []
    for (let i = 0; i < data.length; i++) {
      const element = data[i]
      mountPointList.value.push({
        label: element.device,
        value: element.mountpoint,
      })
    }
  }
  catch (error) {

  }
}

onMounted(() => {
  getMountPointList()
})
</script>

<template>
  <div>
    <!-- <div>{{ genericProgressStyleExtendParam }}</div>
    <div>{{ data }}</div> -->
    <div class="flex mb-5 justify-center">
      <div class="w-[200px]">
        <GenericMonitorCard
          icon-text-icon-hide-title
          :card-type-style="PanelPanelConfigStyleEnum.info"
          icon="solar-cpu-bold"
          :background-color="data.backgroundColor"
          :text-color="data.color"
        >
          <template #info>
            <GenericProgress
              :progress-color="data.progressColor"
              :progress-rail-color="data.progressRailColor"
              :percentage="50"
              :progress-height="5"
              info-card-left-text="DEMO"
              info-card-right-text="TEXT"
              :text-color="data.color"
              :card-type-style="PanelPanelConfigStyleEnum.info"
            />
          </template>
        </GenericMonitorCard>
      </div>

      <div class="w-[80px] ml-2">
        <GenericMonitorCard
          icon-text-icon-hide-title
          :card-type-style="PanelPanelConfigStyleEnum.icon"
          icon="solar-cpu-bold"
          :background-color="data.backgroundColor"
          :icon-text-color="data.color"
        >
          <template #small>
            <GenericProgress
              :progress-color="data.progressColor"
              :progress-rail-color="data.progressRailColor"
              :percentage="50"
              :text-color="data.color"
              :card-type-style="PanelPanelConfigStyleEnum.icon"
            />
          </template>
        </GenericMonitorCard>
      </div>
    </div>

    <NForm ref="formRef" v-model="data">
      <NFormItem label="挂载点">
        <NSelect v-model:value="data.path" size="small" :options="mountPointList" />
      </NFormItem>
      <NFormItem label="主色">
        <NColorPicker v-model:value="data.progressColor" :modes="['hex']" size="small" />
      </NFormItem>
      <NFormItem label="副色">
        <NColorPicker v-model:value="data.progressRailColor" :modes="['hex']" size="small" />
      </NFormItem>
      <NFormItem label="文字图标颜色">
        <NColorPicker v-model:value="data.color" :modes="['hex']" size="small" />
      </NFormItem>
      <NFormItem label="卡片背景色">
        <NColorPicker v-model:value="data.backgroundColor" :modes="['hex']" size="small" />
      </NFormItem>
    </NForm>
  </div>
</template>
