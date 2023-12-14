<script setup lang="ts">
import { computed, onMounted, onUnmounted, ref } from 'vue'
import { NLayout, NLayoutContent, NLayoutSider, NSpace } from 'naive-ui'
import { useAuthStore } from '@/store'
import { AppLoader, RoundCardModal, SvgIcon } from '@/components/common'

interface App {
  name: string
  componentName: string
  icon: string
  auth?: number
}

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

const apps = ref<App[]>([
  {
    name: '用户信息',
    componentName: 'UserInfo',
    icon: 'material-symbols-person-edit-outline-rounded',
  },
  {
    name: '基本设置',
    componentName: 'Style',
    icon: 'ep-setting',
  },
  {
    name: '分组管理',
    componentName: 'ItemGroupManage',
    icon: 'material-symbols-ad-group-outline-rounded',
  },
  {
    name: '导入导出',
    componentName: 'ImportExport',
    icon: 'icon-park-outline-import-and-export',
  },
  {
    name: '关于',
    componentName: 'About',
    icon: 'lucide-info',
  },
])

const authStore = useAuthStore()

const show = computed({
  get: () => props.visible,
  set: (visible: boolean) => {
    emit('update:visible', visible)
  },
})

function handleClickApp(item: App) {
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
    isSmallScreen.value = false
  }
}

onMounted(() => {
  const adminApp: App = {
    name: '用户管理',
    componentName: 'Users',
    icon: 'lucide-users',
    auth: 1,
  }
  // 初始化
  if (authStore.userInfo?.role === 1)
    apps.value.push(adminApp)

  window.addEventListener('resize', handleResize)
  handleResize()
})

onUnmounted(() => {
  window.removeEventListener('resize', handleResize)
})
</script>

<template>
  <div>
    <RoundCardModal
      v-model:show="show"
      style="max-width: 900px;"
      title="应用列表"
    >
      <template #header>
        <div class="flex items-center select-none" @click="collapsed = !collapsed">
          <div class="text-3xl cursor-pointer" style="color:var(--n-color-target)">
            <SvgIcon class=" transition-all duration-500" :icon="collapsed ? 'tabler-layout-sidebar-right-collapse-filled' : 'tabler-layout-sidebar-left-collapse-filled'" />
          </div>
          <div class="ml-1">
            系统应用 & 设置
          </div>
        </div>
      </template>
      <div class="w-full h-full">
        <NSpace vertical size="large" style="height: 100%;width: 100%;">
          <NLayout has-sider>
            <NLayoutSider
              v-model:collapsed="collapsed"
              collapse-mode="width"
              :collapsed-width="0"
              :width="isSmallScreen ? '100%' : 240"
              style="height: 100%;"
              content-style="overflow: hidden"
            >
              <div class="h-[500px] p-[5px] bg-slate-200 rounded-xl overflow-auto" :style="{ width: isSmallScreen ? '100%' : '220px', minWidth: '200px' }">
                <div
                  v-for=" (item, index) in apps"
                  :key="index"
                  :style="{ color: componentName === item.componentName ? 'var(--n-color-target)' : '' }"
                  @click="handleClickApp(item)"
                >
                  <div
                    class="bg-white p-[10px] rounded-lg mb-[5px] font-bold cursor-pointer flex items-center"
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
            </NLayoutSider>
            <NLayoutContent content-style="height:500px;overflow:hidden">
              <div class="p-[5px] rounded-lg transition-all duration-500 min-w-[300px]" :class="(isSmallScreen && !collapsed) ? 'opacity-0' : 'opacity-100'">
                <AppLoader :component-name="componentName" />
              </div>
            </NLayoutContent>
          </NLayout>
        </NSpace>
      </div>
    </RoundCardModal>
  </div>
</template>

<style scoped>
.text-shadow {
  text-shadow: 0px 0px 5px gray;
}
</style>
