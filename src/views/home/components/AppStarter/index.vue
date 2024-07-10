<script setup lang="ts">
import { computed, onMounted, onUnmounted, ref } from 'vue'
import { NLayout, NLayoutContent, NLayoutSider, NSpace } from 'naive-ui'
import { useAuthStore } from '@/store'
import { AppLoader, SModal, SvgIcon } from '@/components/common'
import { t } from '@/locales'
import type { SystemApp } from '@/utils/cmn/builtInApps'
import { adminSystemApps, systemApps } from '@/utils/cmn/builtInApps'

const props = defineProps<{
  visible: boolean
}>()

const emit = defineEmits<{
  (e: 'update:visible', visible: boolean): void
}>()

const componentName = ref('UserInfo')
const collapsed = ref(false)
const screenWidth = ref(0)
const isSmallScreen = ref(false)
const defaultTitle = t('appLauncher.title')
const title = ref('')
const height = ref('600px')
const maxWidth = ref('1000px')

const apps = ref(systemApps)

const authStore = useAuthStore()

const show = computed({
  get: () => props.visible,
  set: (visible: boolean) => {
    emit('update:visible', visible)
  },
})

function handleClickApp(item: SystemApp) {
  componentName.value = item.componentName
  if (isSmallScreen.value)
    collapsed.value = true
}

function getScreenWidth() {
  return window.innerWidth
}

function handleResize() {
  screenWidth.value = getScreenWidth()
  if (screenWidth.value < 640) {
    collapsed.value = true
    isSmallScreen.value = true
  }
  else {
    collapsed.value = false
    isSmallScreen.value = false
  }
}

// function handleMaxMinButton() {
//   height.value = 'calc(100vh - 70px)'
//   maxWidth.value = '100vw'
// }

onMounted(() => {
  const adminApps = adminSystemApps
  // 初始化
  if (authStore.userInfo?.role === 1)
    apps.value = apps.value.concat(adminApps)

  window.addEventListener('resize', handleResize)
  handleResize()
})

onUnmounted(() => {
  window.removeEventListener('resize', handleResize)
})
</script>

<template>
  <div>
    <SModal
      v-model:show="show"
      style="max-width: 1000px;border-radius: 1rem;"
      :style="{ maxWidth }"
      size="small"
      move
    >
      <!-- <template #header-extra>
        <NButton size="small" @click="handleMaxMinButton">
          最大化
        </NButton>
      </template> -->
      <template #header>
        <div class="flex items-center select-none">
          <div
            class="text-3xl cursor-pointer"
            style="color:var(--n-color-target)"
            @click="collapsed = !collapsed"
          >
            <SvgIcon class=" transition-all duration-500" :icon="collapsed ? 'tabler-layout-sidebar-right-collapse-filled' : 'tabler-layout-sidebar-left-collapse-filled'" />
          </div>
          <div class="ml-1" @click="collapsed = !collapsed">
            {{ title === '' ? defaultTitle : title }}
          </div>
        </div>
      </template>
      <div class="w-full h-full app-starter-modal-content">
        <NSpace vertical size="large" style="height: 100%;width: 100%;">
          <NLayout has-sider style="border-radius:0.75rem;">
            <NLayoutSider
              v-model:collapsed="collapsed"
              collapse-mode="width"
              :collapsed-width="0"
              :width="isSmallScreen ? '100%' : 240"
              style="height: 100%;"
              content-style="overflow: hidden"
            >
              <div class="w-full h-full dark:bg-[#2c2c32]">
                <div
                  class="p-[5px] rounded-xl overflow-auto border dark:border-zinc-700 bg-slate-200 dark:bg-zinc-900 "
                  :style="{
                    width: isSmallScreen ? '100%' : '220px',
                    minWidth: '200px',
                    height,
                  }"
                >
                  <div
                    v-for=" (item, index) in apps"
                    :key="index"
                    :style="{ color: componentName === item.componentName ? 'var(--n-color-target)' : '' }"
                    @click="handleClickApp(item)"
                  >
                    <div
                      class="bg-white dark:bg-zinc-800 p-[10px] rounded-lg mb-[5px] font-bold cursor-pointer flex items-center hover:bg-slate-50 focus:bg-slate-50"
                    >
                      <div class="flex items-center justify-center">
                        <div class="text-lg">
                          <SvgIcon :icon="item.icon" />
                        </div>
                        <span class="ml-2">{{ item.name }}</span>
                      </div>
                    <!-- 更多按钮 -->
                    <!-- <div class="ml-auto">
                      <SvgIcon icon="mingcute-more-1-fill" />
                    </div> -->
                    </div>
                  </div>
                </div>
              </div>
            </NLayoutSider>
            <NLayoutContent :content-style="{ height }">
              <div class="rounded-2xl h-full overflow-auto transition-all duration-500 min-w-[300px] h-full border dark:border-zinc-700 bg-slate-200 dark:bg-zinc-900 " :class="(isSmallScreen && !collapsed) ? 'opacity-0' : 'opacity-100'">
                <AppLoader :component-name="componentName" class="h-full" />
              </div>
            </NLayoutContent>
          </NLayout>
        </NSpace>
      </div>
    </SModal>
  </div>
</template>

<style scoped>
.text-shadow {
  text-shadow: 0px 0px 5px gray;
}
</style>

<style>
.dark .app-starter-modal-content .n-layout{
    background-color: #2c2c32;
}
</style>
