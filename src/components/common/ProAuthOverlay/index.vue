<script setup lang="ts">
import { Overlay } from '@/components/common'

withDefaults(defineProps<{
  show: boolean
  title?: string
}>(), {
  show: true,
  title: '',
})
</script>

<template>
  <Overlay :show="show" :title="title || $t('common.noPorAuth')">
    <slot />
  </Overlay>
</template>
