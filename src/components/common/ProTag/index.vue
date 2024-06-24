<script setup lang="ts">
import { NTag, NTooltip } from 'naive-ui'
import { SvgIcon } from '..'
defineProps<{
  showIcon?: boolean
  content?: string
}>()
</script>

<template>
  <NTooltip v-if="content">
    <template #trigger>
      <NTag size="small" style="border-radius: 10px;font-weight: 800;" type="warning">
        <template #icon>
          <slot v-if="content && showIcon" name="icon">
            <SvgIcon
              icon="pepicons-pop--exclamation-circle-filled"
            />
          </slot>
        </template>
        PRO
      </NTag>
    </template>
    {{ content }}
  </NTooltip>

  <NTag v-else size="small" style="border-radius: 10px;font-weight: 800;" type="warning">
    <template #icon>
      <slot v-if="content && showIcon" name="icon">
        <SvgIcon
          icon="pepicons-pop--exclamation-circle-filled"
        />
      </slot>
    </template>
    PRO
  </NTag>
</template>
