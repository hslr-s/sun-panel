<script setup lang="ts">
import { computed } from 'vue'
import Standard from './Standard.vue'
import { PanelCardType } from '@/enums'
import type { PanelPanelConfigStyleEnum } from '@/enums'

interface Prop {
  cardInfo: Panel.ItemCard.CardInfo
  forceBackground?: string // 强制背景色
  iconTextColor?: string
  iconTextInfoHideDescription?: boolean
  iconTextIconHideTitle?: boolean
  cardType: PanelCardType
  cardStyle: PanelPanelConfigStyleEnum
  sortStatus: boolean
}

const props = defineProps<Prop>()
const emit = defineEmits<{
  (e: 'delete'): void
  (e: 'edit'): void
}>()
const backgroundColor = computed(() => {
  return (props.cardInfo?.backgroundColor ? props.cardInfo.backgroundColor : props.cardInfo?.icon?.backgroundColor) || '#2a2a2a6b'
})

function deleteItem() {
  emit('delete')
}

function editItem() {
  emit('edit')
}
</script>

<template>
  <div class="item-card-container w-full">
    <!-- 标准/普通图标 -->
    <template v-if="cardType === PanelCardType.Standard">
      <Standard
        :sort-status="sortStatus"
        :card-type-style="cardStyle"
        :background-color="backgroundColor"
        :icon-text-icon-hide-title="iconTextIconHideTitle || false"
        :icon-text-color="iconTextColor"
        :icon-text-info-hide-description="iconTextInfoHideDescription || false"
        :icon-text="cardInfo?.title"
        :item-info="cardInfo"
        :card-style="cardStyle"
        @delete="deleteItem"
        @edit="editItem"
      />
    </template>
  </div>
</template>
