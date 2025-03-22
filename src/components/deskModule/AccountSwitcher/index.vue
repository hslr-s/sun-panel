<script setup lang="ts">
import { NPopover } from 'naive-ui'
import { onMounted, ref } from 'vue'
import { AccountSwitcherManager } from '@/components/common'
import { useAccountManagerStore, useAuthStore } from '@/store'

const settingShow = ref(false)
const accountManagerStore = useAccountManagerStore()
const authStore = useAuthStore()

function handleOpenSetting() {
  settingShow.value = true
}

onMounted(() => {
  if (authStore.userInfo && !accountManagerStore.isExistAccount(authStore.userInfo?.username || '')) {
    accountManagerStore.addAccount({
      username: authStore.userInfo?.username || '',
      name: authStore.userInfo?.name || '',
      id: authStore.userInfo?.id || 0,
      cToken: authStore.token || '',
      loginStatus: true,
    })
  }
  accountManagerStore.loginCheck()
})
</script>

<template>
  <div>
    <NPopover trigger="click" style="border-radius: 12px;">
      <template #trigger>
        <slot />
      </template>
      <div class="w-[300px]">
        <AccountSwitcherManager :is-popup="true" @open-setting="handleOpenSetting" />
      </div>
    </NPopover>

    <!-- <RoundCardModal
      v-model:show="settingShow"
      style="width:400px"
      :title="t('deskModule.accountSwitcher.loggedInAccount')"
    >
      <AccountSwitcherManager />
    </RoundCardModal> -->
  </div>
</template>
