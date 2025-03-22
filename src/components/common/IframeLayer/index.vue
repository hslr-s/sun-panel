<script setup lang="ts">
import { ref, watch } from 'vue'
import { NSkeleton, NSpin } from 'naive-ui'
import { FreeLayer } from '..'

interface Prop {
  show: boolean
  title: string
  width?: number | string
  height?: number | string
  left?: number
  top?: number
  src: string
  resize?: boolean
  move?: boolean
}
const props = withDefaults(defineProps<Prop>(), {
  width: (window.innerWidth > 800) ? '80%' : '100%',
  height: (window.innerWidth > 800) ? '80%' : '100%',
  resize: false,
  move: false,
})

const emit = defineEmits<{
  (e: 'update:show', show: boolean): void
}>()

const freeLayerShow = ref(props.show)
const windowIframeIsLoad = ref<boolean>(true)

function handWindowIframeIdLoad(payload: Event) {
  windowIframeIsLoad.value = false
}

function handleOnClose() {
  emit('update:show', false)
  freeLayerShow.value = false
}

watch(() => props.show, (v: boolean) => {
  freeLayerShow.value = v
  if (v) {
    windowIframeIsLoad.value = true
  }
})
</script>

<template>
  <FreeLayer
    v-model:show="freeLayerShow"
    :width="width"
    :height="height"
    :left="props.left"
    :top="props.top"
    :move="move"
    :resize="resize"
    @on-close="handleOnClose"
  >
    <template #header>
      <div class="flex items-center">
        <span class="mr-[20px]">
          {{ title }}
        </span>

        <NSpin v-if="windowIframeIsLoad" size="small" />
      </div>
    </template>

    <div class="content m-[10px] rounded-lg border dark:border-zinc-700 overflow-auto" style="height: calc(100% - 60px);">
      <!-- loading层 -->
      <div v-if="windowIframeIsLoad" class="p-5 max-w-[500px] m-auto">
        <div class="flex justify-center">
          <NSkeleton height="35px" width="80%" class="rounded-lg" />
        </div>

        <div class="flex justify-center">
          <NSkeleton height="12px" width="30%" class="mt-[20px] rounded" />
        </div>

        <div class="flex justify-end">
          <NSkeleton height="20px" width="90%" class="mt-[50px] rounded-lg" />
        </div>
        <NSkeleton height="20px" width="100%" class="mt-[20px] rounded-lg" />
        <NSkeleton height="20px" width="100%" class="mt-[20px] rounded-lg" />
        <NSkeleton height="20px" width="80%" class="mt-[20px] rounded-lg" />

        <div class="flex justify-end">
          <NSkeleton height="20px" width="90%" class="mt-[40px] rounded-lg" />
        </div>
        <NSkeleton height="20px" width="80%" class="mt-[20px] rounded-lg" />

        <div class="flex justify-center">
          <NSkeleton height="15px" width="20%" class="mt-[20px] rounded-lg" />
        </div>
      </div>

      <iframe
        v-show="!windowIframeIsLoad"
        :src="src"
        frameborder="0"
        style="height: 100%;width: 100%;"
        @load="handWindowIframeIdLoad"
      />
    </div>
  </FreeLayer>
</template>
