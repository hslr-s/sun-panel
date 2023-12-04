<script setup lang="ts">
import { NButton, NCard, useDialog, useMessage } from 'naive-ui'
import { useAuthStore, usePanelState, useUserStore } from '@/store'
import { logout } from '@/api'
import { router } from '@/router'
import { SvgIcon } from '@/components/common/'

const userStore = useUserStore()
const authStore = useAuthStore()
const panelState = usePanelState()

const ms = useMessage()
const dialog = useDialog()

async function logoutApi() {
  await logout()
  userStore.resetUserInfo()
  authStore.removeToken()
  panelState.removeState()
  ms.success('您已经安全退出，期待与你再次相见！')
  router.push({ path: '/login' })
}

function handleLogiut() {
  dialog.warning({
    title: '警告',
    content: '你确定要退出登录',
    positiveText: '确定',
    negativeText: '取消',
    onPositiveClick: () => {
      logoutApi()
    },
  })
}
</script>

<template>
  <div class="bg-slate-200 rounded-[10px] p-[8px]">
    <NCard style="border-radius:10px" size="small">
      <div class="text-slate-500 mb-[5px]">
        账号/邮箱
      </div>
      {{ authStore.userInfo?.username }}
    </NCard>

    <NCard style="border-radius:10px" class="mt-[10px]" size="small">
      <div class="text-slate-500 mb-[5px]">
        昵称
      </div>
      {{ authStore.userInfo?.name }}
    </NCard>

    <NCard style="border-radius:10px" class="mt-[10px]" size="small">
      <NButton size="small" quaternary type="error" @click="handleLogiut">
        <template #icon>
          <SvgIcon icon="tabler:logout" />
        </template>
        退出登录
      </NButton>
    </NCard>
  </div>
</template>
