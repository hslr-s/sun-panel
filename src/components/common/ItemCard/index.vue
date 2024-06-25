<script setup lang="ts">
import { ref } from 'vue'
import { PanelPanelConfigStyleEnum } from '@/enums'

interface Prop {
  cardTypeStyle: PanelPanelConfigStyleEnum
  class?: string
  backgroundColor?: string
  iconTextIconHideTitle?: boolean // 隐藏小图标标题
  iconTextColor?: string // 小图标文字颜色
  iconText?: string // 小图标文字
}

const props = withDefaults(defineProps<Prop>(), {})

const defaultBackground = '#2a2a2a6b'
const propClass = ref(props.class)
</script>

<template>
  <div class="item-card w-full">
    <!-- 详情图标 -->
    <div
      v-if="cardTypeStyle === PanelPanelConfigStyleEnum.info"
      class="item-card-info w-full rounded-2xl transition-all duration-200 flex"
      :class="propClass"
      :style="{
        'backgroundColor': backgroundColor ?? defaultBackground,
        '--custom-box-shadow-color': backgroundColor ?? defaultBackground,
      }"
    >
      <slot name="info" />
    </div>

    <!-- 极简图标（APP） -->
    <div
      v-if="cardTypeStyle === PanelPanelConfigStyleEnum.icon"
      class="item-card-small"
    >
      <div
        class="item-card-small-icon overflow-hidden rounded-2xl sunpanel w-[70px] h-[70px] mx-auto transition-all duration-200"
        :class="propClass"
        :style="{
          'backgroundColor': backgroundColor ?? defaultBackground,
          '--custom-box-shadow-color': backgroundColor ?? defaultBackground,
        }"
      >
        <slot name="small" />
      </div>

      <div
        v-if="!iconTextIconHideTitle"
        class="item-card-small-title text-center app-icon-text-shadow cursor-pointer mt-[2px]"
        :style="{ color: iconTextColor }"
      >
        {{ iconText }}
      </div>
    </div>
  </div>
</template>

<style scoped>
.item-card-info:hover,.item-card-small-icon:hover {
  box-shadow: 0 0 20px var(--custom-box-shadow-color,#2a2a2a6b),
  0 0 1px 1px rgba(255, 255, 255,.1),
  0 0 20px 5px rgba(0,0,0,.1);
}
</style>
