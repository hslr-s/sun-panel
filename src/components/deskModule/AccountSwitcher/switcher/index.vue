<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { NButton, NDivider, NEllipsis, NFlex, NTag, NText } from 'naive-ui'
import { useAccountManagerStore, useAuthStore } from '@/store'
import type { AccountInfo } from '@/store/modules/accountManager/helper'
import { SvgIcon } from '@/components/common'
import { t } from '@/locales'

const accountManagerStore = useAccountManagerStore()
const authStore = useAuthStore()

const testData = ref<AccountInfo[]>([
  {
    name: '空间166666665555',
    id: 1,
    username: '95302870@qq.com',
    cToken: '123123213231',
    loginStatus: true,
  },
  {
    name: '空间2',
    id: 1,
    username: '123456',
    cToken: '123123213231',
    loginStatus: false,
  },
  {
    name: '空间3',
    id: 1,
    username: 'admin@sun.cc',
    cToken: '123123213231',
    loginStatus: true,
  },
])

onMounted(() => {
  accountManagerStore.accounts = testData.value
})
</script>

<template>
  <div class="w-[300px] p-1">
    <div class="font-bold text-base mb-3 flex">
      <div class="mr-auto">
        {{ t('deskModule.accountSwitcher.accountSwitch') }}
      </div>
      <NFlex justify="end" align="center">
        <NButton size="tiny" type="info">
          <template #icon>
            <SvgIcon icon="lets-icons--user-add" />
          </template>
        </NButton>

        <NButton size="tiny" type="error">
          <template #icon>
            <SvgIcon icon="tabler-logout" />
          </template>
        </NButton>
      </NFlex>
    </div>
    <NDivider style="margin: 10px 0;" />
    <div
      v-for="item, index in accountManagerStore.accounts"
      :key="index"
      class="my-1 cursor-pointer hover:bg-slate-100 dark:hover:bg-slate-500 rounded-lg p-2 transform-all duration-300"
    >
      <NFlex align="center">
        <div class="mr-auto">
          <div>
            <NEllipsis style="max-width: 250px;">
              <NTag
                v-if="authStore.userInfo?.username === item.username"
                type="info" size="tiny" class="mr-1"
              >
                {{ t('common.current') }}
              </NTag>
              <NTag v-if="!item.loginStatus" type="error" size="tiny" class="mr-1">
                {{ t('deskModule.accountSwitcher.failed') }}
              </NTag>
              <NText :depth="2">
                {{ item.name }}
              </NText>
            </NEllipsis>
          </div>
          <div>
            <NEllipsis style="max-width: 250px;">
              <NText :depth="3">
                {{ item.username }}
              </NText>
            </NEllipsis>
          </div>
        </div>

        <!-- 按钮 -->
        <NFlex v-if="!(authStore.userInfo?.username === item.username)">
          <NButton v-if="item.loginStatus" size="tiny" dashed>
            <template #icon>
              <SvgIcon icon="tabler-logout" />
            </template>
          </NButton>

          <NButton v-else size="tiny" dashed>
            <template #icon>
              <SvgIcon icon="material-symbols-delete" />
            </template>
          </NButton>
        </NFlex>
      </NFlex>
    </div>
  </div>
</template>
