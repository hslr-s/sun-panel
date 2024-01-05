<script setup lang="ts">
import { ref } from 'vue'
import { NColorPicker, NForm, NFormItem } from 'naive-ui'
import type { GenericProgressStyleExtendParam } from '../../typings'
import GenericMonitorCard from '../../components/GenericMonitorCard/index.vue'
import GenericProgress from '../../components/GenericProgress/index.vue'
import { PanelPanelConfigStyleEnum } from '@/enums'

const props = defineProps<{
  genericProgressStyleExtendParam: GenericProgressStyleExtendParam
}>()

// const emit = defineEmits<{
//   (e: 'update:genericProgressStyleExtendParam', visible: ProgressStyle): void
// }>()

const data = ref(props.genericProgressStyleExtendParam)
</script>

<template>
  <div>
    <div class="flex mb-5 justify-center">
      <div class="w-[200px]">
        <GenericMonitorCard
          icon-text-icon-hide-title
          :card-type-style="PanelPanelConfigStyleEnum.info"
          class="hover:shadow-[0_0_20px_10px_rgba(0,0,0,0.2)]"
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
          class="hover:shadow-[0_0_20px_10px_rgba(0,0,0,0.2)]"
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

    <NForm ref="formRef">
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
