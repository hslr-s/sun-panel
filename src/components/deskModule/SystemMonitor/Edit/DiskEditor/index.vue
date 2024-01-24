<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import type { FormInst, FormRules } from 'naive-ui'
import { NColorPicker, NForm, NFormItem, NSelect } from 'naive-ui'
import type { DiskExtendParam } from '../../typings'
import GenericMonitorCard from '../../components/GenericMonitorCard/index.vue'
import GenericProgress from '../../components/GenericProgress/index.vue'
import { PanelPanelConfigStyleEnum } from '@/enums'
import { getDiskMountpoints } from '@/api/system/systemMonitor'
import { defautSwatchesBackground } from '@/utils/defaultData'
import { t } from '@/locales'

interface Emit {
  (e: 'update:diskExtendParam', visible: DiskExtendParam): void
}

const props = defineProps<{
  diskExtendParam: DiskExtendParam
}>()
const emit = defineEmits<Emit>()
const formRef = ref<FormInst | null>(null)
const rules: FormRules = {
  path: {
    required: true,
    message: t('form.required'),
    trigger: 'blur',
  },
}

const mountPointList = ref<{
  label: string
  value: string
}[]>([])

const extendParam = computed({
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
      if (i === 0 && !extendParam.value.path)
        extendParam.value.path = element.mountpoint

      mountPointList.value.push({
        label: element.mountpoint,
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

defineExpose({
  async verification(): Promise<boolean> {
    try {
      const errors = await formRef.value?.validate()

      if (!errors) {
        return Promise.resolve(true)
      }
      else {
        console.log(errors)
        return Promise.resolve(false)
      }
    }
    catch (error) {
      console.error('An error occurred during validation:', error)
      return Promise.resolve(false)
    }
  },
})
</script>

<template>
  <div>
    <!-- <div>{{ diskExtendParam }}</div> -->
    <div class="flex mb-5 justify-center transparent-grid p-2 rounded-xl border">
      <div class="w-[200px]">
        <GenericMonitorCard
          icon-text-icon-hide-title
          :card-type-style="PanelPanelConfigStyleEnum.info"
          icon="solar-cpu-bold"
          :background-color="extendParam.backgroundColor"
          :text-color="extendParam.color"
        >
          <template #info>
            <GenericProgress
              :progress-color="extendParam.progressColor"
              :progress-rail-color="extendParam.progressRailColor"
              :percentage="50"
              :progress-height="5"
              info-card-left-text="DEMO"
              info-card-right-text="TEXT"
              :text-color="extendParam.color"
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
          :background-color="extendParam.backgroundColor"
          :icon-text-color="extendParam.color"
        >
          <template #small>
            <GenericProgress
              :progress-color="extendParam.progressColor"
              :progress-rail-color="extendParam.progressRailColor"
              :percentage="50"
              :text-color="extendParam.color"
              :card-type-style="PanelPanelConfigStyleEnum.icon"
            />
          </template>
        </GenericMonitorCard>
      </div>
    </div>

    <NForm ref="formRef" v-model="extendParam" :model="extendParam" :rules="rules">
      <NFormItem :label="$t('deskModule.systemMonitor.diskMountPoint')" path="path">
        <NSelect v-model:value="extendParam.path" size="small" :options="mountPointList" />
      </NFormItem>
      <NFormItem :label="$t('deskModule.systemMonitor.progressColor')">
        <NColorPicker v-model:value="extendParam.progressColor" :swatches="defautSwatchesBackground" :modes="['hex']" size="small" />
      </NFormItem>
      <NFormItem :label="$t('deskModule.systemMonitor.progressRailColor')">
        <NColorPicker v-model:value="extendParam.progressRailColor" :swatches="defautSwatchesBackground" :modes="['hex']" size="small" />
      </NFormItem>
      <NFormItem :label="$t('common.textColor')">
        <NColorPicker v-model:value="extendParam.color" :swatches="defautSwatchesBackground" :modes="['hex']" size="small" />
      </NFormItem>
      <NFormItem :label="$t('common.backgroundColor')">
        <NColorPicker v-model:value="extendParam.backgroundColor" :swatches="defautSwatchesBackground" :modes="['hex']" size="small" />
      </NFormItem>
    </NForm>
  </div>
</template>

<style>
.transparent-grid {
    background-image: linear-gradient(45deg, #fff 25%, transparent 25%, transparent 75%, #fff 75%),
                      linear-gradient(45deg, #fff 25%, transparent 25%, transparent 75%, #fff 75%);
    background-size: 16px 16px;
    background-position: 0 0, 8px 8px;
    background-color: #e2e8f0;
}
</style>
