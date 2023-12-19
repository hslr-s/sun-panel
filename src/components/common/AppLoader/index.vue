<script setup lang="ts">
import { defineAsyncComponent, onMounted, shallowRef, watch } from 'vue'
import { NSpin } from 'naive-ui'

const props = defineProps<{
  componentName: string | null
}>()
const loading = shallowRef(false)
const dynamicComponent = shallowRef('')

function updateComponent() {
  loading.value = true
  dynamicComponent.value = defineAsyncComponent(() =>
    import(`../../apps/${props.componentName}/index.vue`)
      .finally(() => {
        loading.value = false
      }).catch(() => {
      // 组件不存在
        dynamicComponent.value = ''
        return null
      }),
  )
}

watch(() => props.componentName, () => {
  updateComponent()
})

onMounted(() => {
  updateComponent()
})
</script>

<template>
  <div class="h-full">
    <NSpin :show="loading" style="height: 100%;" content-style="height: 100%;" :delay="500" description="loading...">
      <component :is="dynamicComponent" v-if="dynamicComponent" />
      <!-- <component :is="getComponent(componentName || '')" v-if="dynamicComponent" /> -->
      <div
        v-else-if="!dynamicComponent"
      >
        Component not found!
      </div>
    </NSpin>
  </div>
</template>
