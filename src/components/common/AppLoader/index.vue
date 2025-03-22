<script setup lang="ts">
import { defineAsyncComponent, onMounted, shallowRef, watch } from 'vue'
import { NButton, NResult, NSpin } from 'naive-ui'
import { t } from '@/locales'
import { router } from '@/router'

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
      <template v-if="dynamicComponent">
        <div v-show="!loading" class="fade-in-animation h-full">
          <component :is="dynamicComponent" />
        </div>
      </template>

      <!-- <component :is="getComponent(componentName || '')" v-if="dynamicComponent" /> -->
      <div v-else-if="!dynamicComponent">
        <!-- Component not found! -->
        <div class="mt-[10px]">
          <NResult status="404" :description="$t('apps.common.notFound')">
            <template #footer>
              <NButton type="success" @click="router.push({ name: 'login' })">
                {{ t('common.goToLogin') }}
              </NButton>
            </template>
          </NResult>
        </div>
      </div>
    </NSpin>
  </div>
</template>

<style scoped>
@keyframes fadeIn {
  from {
    opacity: 0;
  }

  to {
    opacity: 1;
  }
}

.fade-in-animation {
  animation: fadeIn 0.2s ease-in-out;
}
</style>
