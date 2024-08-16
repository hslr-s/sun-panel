<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { NButton, NDivider, NEllipsis, NFlex, NForm, NFormItem, NInput, NTag, NText, useDialog } from 'naive-ui'
import { useAccountManagerStore, useAuthStore, usePanelState, useUserStore } from '@/store'
import type { AccountInfo } from '@/store/modules/accountManager/helper'
import { RoundCardModal, SvgIcon } from '@/components/common'
import { t } from '@/locales'
import { login, logout } from '@/api'
import { apiRespErrMsg, message } from '@/utils/request/apiMessage'

defineProps<{
  isPopup?: boolean
}>()

const emits = defineEmits<{
  (e: 'openSetting'): void
}>()

const accountManagerStore = useAccountManagerStore()
const authStore = useAuthStore()
const panelState = usePanelState()
const dialog = useDialog()
const userStore = useUserStore()

const loginModalShow = ref(false)
const loginModalLoading = ref(false)
const logoutAllLoading = ref(false)

const loginModalForm = ref({
  username: '',
  password: '',
})

function handleOpenSetting() {
  emits('openSetting')
}

function handleSubmit() {
  if (accountManagerStore.isExistAccount(loginModalForm.value.username)) {
    message.warning(t('deskModule.accountSwitcher.accountExist'))
    return
  }
  login<User.Info>({
    username: loginModalForm.value.username,
    password: loginModalForm.value.password,
  }).then(({ data }) => {
    console.log(data)
    const d: AccountInfo = {
      name: data.name as string,
      id: data.id as number,
      username: data.username as string,
      cToken: data.token as string,
      loginStatus: true,
    }
    accountManagerStore.addAccount(d)
    loginModalShow.value = false
  }).catch((res) => {
    apiRespErrMsg(res)
  })
}

// 退出
function handleLogoutItem(item: AccountInfo) {
  dialog.warning({
    title: t('common.warning'),
    content: t('deskModule.accountSwitcher.logoutConfirmText', { username: `${item.username} - ${item.name}` }),
    positiveText: t('common.confirm'),
    negativeText: t('common.cancel'),
    onPositiveClick: () => {
      logout(item.cToken).then(() => {
        handleDeleteItem(item)
      }).catch((res) => {
        apiRespErrMsg(res)
      })
    },
  })
}

function handleLogoutAll() {
  dialog.warning({
    title: t('common.warning'),
    content: t('deskModule.accountSwitcher.logoutAllConfirmText'),
    positiveText: t('common.confirm'),
    negativeText: t('common.cancel'),
    loading: logoutAllLoading.value,
    onPositiveClick: async () => {
      logoutAllLoading.value = true
      for (let i = 0; i < accountManagerStore.accounts.length; i++) {
        const element = accountManagerStore.accounts[i]
        await logout(element.cToken).then(() => {
          handleDeleteItem(element)
        }).catch((res) => {
          apiRespErrMsg(res)
        })
      }

      accountManagerStore.clearAccount()
      userStore.resetUserInfo()
      authStore.removeToken()
      panelState.removeState()
      message.success(t('settingUserInfo.logoutSuccess'))
      logoutAllLoading.value = false
      // 返回首页
      location.href = '/'
    },

  })
}

function handleSwitch(item: AccountInfo) {
  if (!item.loginStatus) {
    return
  }

  // 当前账号不切换
  if (item.username === authStore.userInfo?.username) {
    return
  }

  const userInfo: User.Info = {
    id: item.id,
    username: item.username,
    name: item.name,
  }

  panelState.setItemGroupCacheData([])
  authStore.setToken(item.cToken)
  authStore.setUserInfo(userInfo)
  location.reload()
}

function handleDeleteItem(item: AccountInfo) {
  accountManagerStore.deleteAccount(item)
}
onMounted(() => {
  // console.log(accountManagerStore.accounts)
  // accountManagerStore.accounts = testData.value
})
</script>

<template>
  <div class="p-1">
    <div v-if="isPopup" class="font-bold text-base mb-3 flex">
      <div class="mr-auto">
        {{ t('deskModule.accountSwitcher.accountSwitch') }}
      </div>
      <NFlex justify="end" align="center">
        <NButton v-if="isPopup" size="tiny" type="info" @click="handleOpenSetting">
          <template #icon>
            <SvgIcon icon="mingcute--add-fill" />
          </template>
        </NButton>

        <NButton size="tiny" type="error" @click="handleLogoutAll">
          <template #icon>
            <SvgIcon icon="tabler-logout" />
          </template>
        </NButton>
      </NFlex>
    </div>
    <NDivider v-if="isPopup" style="margin: 10px 0;" />
    <div
      v-for="item, index in accountManagerStore.accounts"
      :key="index"
      class="my-1 cursor-pointer hover:bg-slate-100 dark:hover:bg-slate-500 rounded-lg p-2 transform-all duration-300"
    >
      <NFlex align="center">
        <div class="mr-auto" style="width: calc(100% - 50px);" @click="handleSwitch(item)">
          <div>
            <NEllipsis>
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
            <NEllipsis>
              <NText :depth="3">
                {{ item.username }}
              </NText>
            </NEllipsis>
          </div>
        </div>

        <!-- 按钮 -->
        <NFlex v-if="!(authStore.userInfo?.username === item.username)">
          <NButton v-if="item.loginStatus" size="tiny" dashed :title="t('common.logout')" @click="handleLogoutItem(item)">
            <template #icon>
              <SvgIcon icon="tabler-logout" />
            </template>
          </NButton>

          <NButton v-else size="tiny" dashed :title="t('common.delete')" @click="handleDeleteItem(item)">
            <template #icon>
              <SvgIcon icon="material-symbols-delete" />
            </template>
          </NButton>
        </NFlex>
      </NFlex>
    </div>

    <NFlex v-if="!isPopup" justify="end" align="center">
      <!-- <NButton size="tiny" type="info" @click="loginModalShow = true">
        <template #icon>
          <SvgIcon icon="tabler--refresh" />
        </template>
      </NButton> -->

      <NButton size="tiny" type="info" @click="loginModalShow = true">
        <template #icon>
          <SvgIcon icon="lets-icons--user-add" />
        </template>
      </NButton>

      <NButton size="tiny" type="error" @click="handleLogoutAll">
        <template #icon>
          <SvgIcon icon="tabler-logout" />
        </template>
      </NButton>
    </NFlex>

    <div>
      <RoundCardModal v-model:show="loginModalShow" style="width:400px" :title="t('deskModule.accountSwitcher.addAccount')">
        <div class="mt-2">
          <NForm>
            <NFormItem :label="t('common.username')">
              <NInput v-model:value="loginModalForm.username" :placeholder="t('common.username')" />
            </NFormItem>
            <NFormItem :label="t('common.password')">
              <NInput v-model:value="loginModalForm.password" type="password" :placeholder="t('common.password')" />
            </NFormItem>
          </NForm>
          <NButton type="primary" :loading="loginModalLoading" block secondary strong @click="handleSubmit">
            {{ t("common.login") }}
          </NButton>
        </div>
      </RoundCardModal>
    </div>
  </div>
</template>
