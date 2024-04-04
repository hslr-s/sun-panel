<script setup lang="ts">
import { Overlay } from '@/components/common'
import { t } from '@/locales'

withDefaults(defineProps<{
  show: boolean
  title?: string
}>(), {
  show: true,
  title: t('common.noPorAuth'),
})
</script>

<template>
  <Overlay :show="show" :title="title">
    <slot />
  </Overlay>
</template>
