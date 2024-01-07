<script setup lang="ts">
import { computed } from 'vue'
import { NColorPicker, NForm, NFormItem } from 'naive-ui'
import type { GenericProgressStyleExtendParam } from '../../typings'
import GenericMonitorCard from '../../components/GenericMonitorCard/index.vue'
import GenericProgress from '../../components/GenericProgress/index.vue'
import { PanelPanelConfigStyleEnum } from '@/enums'
import { defautSwatchesBackground } from '@/utils/defaultData'

interface Emit {
  (e: 'update:genericProgressStyleExtendParam', visible: GenericProgressStyleExtendParam): void
}

const props = defineProps<{
  genericProgressStyleExtendParam: GenericProgressStyleExtendParam
}>()
const emit = defineEmits<Emit>()

const data = computed({
  get: () => props.genericProgressStyleExtendParam,
  set: (visible) => {
    emit('update:genericProgressStyleExtendParam', visible)
  },
})
</script>

<template>
  <div>
    <!-- <div>{{ genericProgressStyleExtendParam }}</div>
    <div>{{ data }}</div> -->
    <div class="flex mb-5 justify-center transparent-grid p-2 rounded-xl border">
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
      <NFormItem :label="$t('deskModule.systemMonitor.progressColor')">
        <NColorPicker v-model:value="data.progressColor" :swatches="defautSwatchesBackground" :modes="['hex']" size="small" />
      </NFormItem>
      <NFormItem :label="$t('deskModule.systemMonitor.progressRailColor')">
        <NColorPicker v-model:value="data.progressRailColor" :swatches="defautSwatchesBackground" :modes="['hex']" size="small" />
      </NFormItem>
      <NFormItem :label="$t('common.textColor')">
        <NColorPicker v-model:value="data.color" :swatches="defautSwatchesBackground" :modes="['hex']" size="small" />
      </NFormItem>
      <NFormItem :label="$t('common.backgroundColor')">
        <NColorPicker v-model:value="data.backgroundColor" :swatches="defautSwatchesBackground" :modes="['hex']" size="small" />
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
