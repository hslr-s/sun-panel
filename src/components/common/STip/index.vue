<script lang="ts" setup>
import { withDefaults } from 'vue'

import { NTooltip } from 'naive-ui'
import { SvgIcon } from '..'

withDefaults(defineProps<{
  icon?: string
  content?: string
}>(), {
  icon: 'pepicons-pop--exclamation-circle-filled',
  content: '',
})
</script>

<template>
  <span>
    <NTooltip>
      <template #trigger>
        <SvgIcon
          :icon="icon"
        />
      </template>
      {{ content }}
    </NTooltip>
  </span>
</template>
