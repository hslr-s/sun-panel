<script setup lang="ts">
import type { FormInst, FormRules } from 'naive-ui'
import { NButton, NCard, NDivider, NFlex, NForm, NFormItem, NInput, NSelect, useMessage } from 'naive-ui'
import { ref } from 'vue'
import { useAppStore, useAuthStore } from '@/store'
import { languageOptions } from '@/utils/defaultData'
import type { Language, Theme } from '@/store/modules/app/helper'
import { AccountSwitcherManager, RoundCardModal, STip } from '@/components/common/'
import { updateInfo, updatePassword } from '@/api/system/user'
import { updateLocalUserInfo } from '@/utils/cmn'
import { t } from '@/locales'

// const userStore = useUserStore()
const authStore = useAuthStore()
const appStore = useAppStore()
// const panelState = usePanelState()
// const accountManagerStore = useAccountManagerStore()

const ms = useMessage()
// const dialog = useDialog()

const languageValue = ref(appStore.language)
const themeValue = ref(appStore.theme)
const nickName = ref(authStore.userInfo?.name || '')
const isEditNickNameStatus = ref(false)
const formRef = ref<FormInst | null>(null)
const themeOptions: { label: string; key: string; value: Theme }[] = [
  { label: t('apps.userInfo.themeStyle.dark'), key: 'dark', value: 'dark' },
  { label: t('apps.userInfo.themeStyle.light'), key: 'light', value: 'light' },
  { label: t('apps.userInfo.themeStyle.auto'), key: 'Auto', value: 'auto' },
]
const updatePasswordModalState = ref({
  show: false,
  loading: false,
  form: {
    password: '',
    oldPassword: '',
    confirmPassword: '',
  },
})

const updatePasswordModalFormRules: FormRules = {
  oldPassword: {
    required: true,
    trigger: 'blur',
    min: 6,
    max: 20,
    message: t('adminSettingUsers.formRules.passwordLimit'),
  },
  password: {
    required: true,
    trigger: 'blur',
    min: 6,
    max: 20,
    message: t('adminSettingUsers.formRules.passwordLimit'),
  },
  confirmPassword: {
    required: true,
    trigger: 'blur',
    min: 6,
    max: 20,
    message: t('adminSettingUsers.formRules.passwordLimit'),
  },
}

// async function logoutApi() {
//   // await logout()
//   if (panelState.panelConfig.accountSwitcherShow) {
//     for (let i = 0; i < accountManagerStore.accounts.length; i++) {
//       const element = accountManagerStore.accounts[i]
//       await logout(element.cToken).then(() => {
//         accountManagerStore.deleteAccount(element)
//       }).catch((res) => {
//         apiRespErrMsg(res)
//       })
//     }
//     accountManagerStore.clearAccount()
//   }

//   userStore.resetUserInfo()
//   authStore.removeToken()
//   panelState.removeState()
//   // appStore.removeToken()
//   ms.success(t('settingUserInfo.logoutSuccess'))
//   // router.push({ path: '/login' })
//   location.href = '/'// 强制刷新一下页面
// }

function handleSaveInfo() {
  updateInfo(nickName.value).then(({ code, msg }) => {
    updateLocalUserInfo()
    isEditNickNameStatus.value = false
  }).catch((error) => {
    ms.error(`${t('common.editFail')}:${error.msg}`)
  })
}

function handleUpdatePassword(e: MouseEvent) {
  e.preventDefault()
  formRef.value?.validate((errors) => {
    if (errors) {
      console.log(errors)
      return
    }

    if (updatePasswordModalState.value.form.password !== updatePasswordModalState.value.form.confirmPassword) {
      ms.error(t('settingUserInfo.confirmPasswordInconsistentMsg'))
      return
    }
    updatePasswordModalState.value.loading = true
    updatePassword(updatePasswordModalState.value.form.oldPassword, updatePasswordModalState.value.form.password).then(({ code, msg }) => {
      // 成功
      updatePasswordModalState.value.show = false
      ms.success(t('common.success'))
    }).finally(() => {
      updatePasswordModalState.value.loading = false
    }).catch((error) => {
      ms.error(`${t('common.serverError')}:${error.msg}`)
    })
  })
}

// function handleLogout() {
//   dialog.warning({
//     title: t('common.warning'),
//     content: t('settingUserInfo.confirmLogoutText'),
//     positiveText: t('common.confirm'),
//     negativeText: t('common.cancel'),
//     onPositiveClick: () => {
//       logoutApi()
//     },
//   })
// }

function handleChangeLanuage(value: Language) {
  languageValue.value = value
  appStore.setLanguage(value)
  location.reload()
}

function handleChangeTheme(value: Theme) {
  themeValue.value = value
  appStore.setTheme(value)
  // location.reload()
}
</script>

<template>
  <div class="bg-slate-200 dark:bg-zinc-900 p-2 h-full mb-10 overflow-auto">
    <!-- 账号信息 -->
    <NCard style="border-radius:10px" size="small">
      <template #header>
        <div class="text-slate-500 font-bold">
          {{ $t('apps.userInfo.currentAccount') }}
        </div>
      </template>
      <div>
        <div class="text-slate-500 font-bold">
          {{ $t('common.username') }}
        </div>
        {{ authStore.userInfo?.username }}
      </div>

      <div class="mt-[10px]">
        <div class="text-slate-500 font-bold">
          {{ $t('common.nikeName') }}
        </div>

        <div v-if="!isEditNickNameStatus">
          {{ authStore.userInfo?.name }}

          <NButton size="small" text type="success" @click="isEditNickNameStatus = !isEditNickNameStatus">
            {{ $t('common.edit') }}
          </NButton>
        </div>

        <div v-else class="flex items-center">
          <div class="max-w-[150px]">
            <NInput v-model:value="nickName" type="text" :placeholder="$t('common.inputPlaceholder')" />
          </div>
          <NButton size="small" quaternary type="success" @click="handleSaveInfo">
            {{ $t('common.save') }}
          </NButton>
        </div>
      </div>

      <NDivider style="margin: 10px 0;" dashed />
      <div>
        <NButton size="small" text type="success" @click="updatePasswordModalState.show = !updatePasswordModalState.show">
          {{ $t('settingUserInfo.updatePassword') }}
        </NButton>
      </div>
    </NCard>

    <!-- 基本设置 -->
    <NCard style="border-radius:10px" class="my-[10px]" size="small">
      <template #header>
        <div class="text-slate-500 font-bold">
          <NFlex align="center">
            {{ $t('common.setting') }}
            <STip :content="t('apps.userInfo.settingSaveLocal')" />
          </NFlex>
        </div>
      </template>
      <div>
        <div class="text-slate-500 font-bold">
          {{ $t('common.language') }}
        </div>
        <div class="max-w-[200px]">
          <NSelect v-model:value="languageValue" size="small" :options="languageOptions" @update-value="handleChangeLanuage" />
        </div>
      </div>

      <div class="mt-[10px]">
        <div class="text-slate-500 font-bold">
          {{ $t('apps.userInfo.theme') }}
        </div>
        <div class="max-w-[200px]">
          <NSelect v-model:value="themeValue" size="small" :options="themeOptions" @update-value="handleChangeTheme" />
        </div>
      </div>
    </NCard>

    <!-- 账号切换 -->
    <NCard style="border-radius:10px" class="mt-[10px]" size="small">
      <template #header>
        <div class="text-slate-500 font-bold">
          <NFlex align="center">
            {{ t('deskModule.accountSwitcher.loggedInAccount') }}
            <STip :content="t('apps.userInfo.settingSaveLocal')" />
          </NFlex>
        </div>
      </template>

      <div>
        <AccountSwitcherManager />
      </div>
    </NCard>
    <!--
    <NCard style="border-radius:10px;margin-bottom: 10px;" class="my-[10px]" size="small">
      <NButton size="small" text type="error" @click="handleLogout">
        <template #icon>
          <SvgIcon icon="tabler:logout" />
        </template>
        {{ $t('settingUserInfo.logout') }}
      </NButton>
    </NCard> -->

    <RoundCardModal v-model:show="updatePasswordModalState.show" size="small" preset="card" style="width: 400px" :title="$t('settingUserInfo.updatePassword')">
      <NForm ref="formRef" :model="updatePasswordModalState.form" :rules="updatePasswordModalFormRules">
        <NFormItem path="oldPassword" :label="$t('settingUserInfo.oldPassword')">
          <NInput v-model:value="updatePasswordModalState.form.oldPassword" :maxlength="20" type="password" :placeholder="$t('settingUserInfo.oldPassword')" />
        </NFormItem>

        <NFormItem path="password" :label="$t('settingUserInfo.newPassword')">
          <NInput v-model:value="updatePasswordModalState.form.password" :maxlength="20" type="password" :placeholder="$t('settingUserInfo.newPassword')" />
        </NFormItem>

        <NFormItem path="confirmPassword" :label="$t('settingUserInfo.confirmPassword')">
          <NInput v-model:value="updatePasswordModalState.form.confirmPassword" :maxlength="20" type="password" :placeholder="$t('settingUserInfo.confirmPassword')" />
        </NFormItem>
      </NForm>

      <template #footer>
        <div class="float-right">
          <NButton type="success" size="small" :loading="updatePasswordModalState.loading" @click="handleUpdatePassword">
            {{ $t('common.save') }}
          </NButton>
        </div>
      </template>
    </RoundCardModal>
  </div>
</template>
