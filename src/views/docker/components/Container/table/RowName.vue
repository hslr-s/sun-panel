<script setup lang="ts">
defineProps<{
  value: Docker.Container.Info
}>()
</script>

<template>
  <div class="w-full">
    <div>
      {{ value.name }}
    </div>
    <div class="text-xs text-slate-400">
      {{ value.image }}
    </div>
  </div>
</template>
