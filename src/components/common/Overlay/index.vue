<script lang='ts' setup>
import { defineProps } from 'vue'

withDefaults(defineProps<{
  show: boolean
  title?: string
}>(), {
  show: true,
})
</script>

<template>
  <div v-if="show" class="overlay-box" :title="title">
    <div class="cover" />
    <slot />
  </div>
  <template v-else>
    <slot />
  </template>
</template>

<style scoped>
.overlay-box{
    /* pointer-events: none; */
    cursor: not-allowed;
    position: relative;
    width:100%;
}

.cover {
    position: absolute;
    top: 0;
    left: 0;
    width:100%;
    height: 100%;
    background-color: black;
    opacity: 0;
    z-index: 10;
    transition: opacity 0.3s;
}

.overlay-box:hover .cover {
    opacity: 0.4;
}
</style>
