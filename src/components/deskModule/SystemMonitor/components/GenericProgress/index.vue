<script setup lang="ts">
import { NProgress } from 'naive-ui'
import { PanelPanelConfigStyleEnum } from '@/enums'

interface Prop {
  textColor: string
  progressColor: string
  progressRailColor: string
  progressHeight?: number
  percentage: number
  cardTypeStyle: PanelPanelConfigStyleEnum
  infoCardLeftText?: string
  infoCardRightText?: string
}

defineProps<Prop>()
</script>

<template>
  <template v-if="cardTypeStyle === PanelPanelConfigStyleEnum.info">
    <div class="w-full">
      <div class="mb-1 text-xs" :style="{ color: textColor }">
        <span>
          {{ infoCardLeftText }}
        </span>
        <span class="float-right">
          {{ infoCardRightText }}
        </span>
      </div>
      <NProgress
        type="line"
        :color="progressColor"
        :rail-color="progressRailColor"
        :height="progressHeight"
        :percentage="percentage"
        :show-indicator="false"
        :stroke-width="15"
        style="max-width: 135px;"
      />
    </div>
  </template>
  <template v-else>
    <div class="w-full flex justify-center h-full w-full mt-3">
      <NProgress
        :color="progressColor"
        :rail-color="progressRailColor"
        type="dashboard"
        :percentage="percentage" :stroke-width="15"
        style="width: 50px;"
      >
        <div class="text-white" style="font-size: 8px;" :style="{ color: textColor }">
          {{ percentage }}%
        </div>
      </NProgress>
    </div>
  </template>
</template>
