<script setup lang="ts">
import { NEllipsis } from 'naive-ui'
defineProps<{
  value: Docker.Container.Info
}>()
</script>

<template>
  <div class="w-full">
    <div>
      <NEllipsis :line-clamp="2" style="max-width: 200px">
        {{ value.name }}
      </NEllipsis>
    </div>
    <div class="text-xs text-slate-400">
      <NEllipsis :line-clamp="2" style="max-width: 400px">
        {{ value.image }}
      </NEllipsis>
    </div>
  </div>
</template>
