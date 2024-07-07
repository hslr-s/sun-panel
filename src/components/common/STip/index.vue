<script lang="ts" setup>
import { withDefaults } from 'vue'

import { NTooltip } from 'naive-ui'
import { SvgIcon } from '..'

withDefaults(defineProps<{
  icon?: string
  content?: string
}>(), {
  icon: 'fa-solid--info-circle',
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
      <slot>
        {{ content }}
      </slot>
    </NTooltip>
  </span>
</template>
