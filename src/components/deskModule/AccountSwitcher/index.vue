<script setup lang="ts">
import { NPopover } from 'naive-ui'
import Switcher from './switcher/index.vue'
</script>

<template>
  <NPopover trigger="click" style="border-radius: 12px;">
    <template #trigger>
      <slot />
    </template>
    <Switcher />
  </NPopover>
</template>
