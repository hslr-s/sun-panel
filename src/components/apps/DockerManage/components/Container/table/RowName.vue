<script setup lang="ts">
import { NEllipsis } from 'naive-ui'
defineProps<{
  value: Docker.Container.Info
}>()
</script>

<template>
  <NEllipsis :line-clamp="2">
    <div>
      {{ value.name }}
    </div>
    <div class="text-xs text-slate-400">
      {{ value.image }}
    </div>
  </NEllipsis>
</template>
