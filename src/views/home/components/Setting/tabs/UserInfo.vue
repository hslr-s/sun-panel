<script setup lang="ts">
import type { FormInst, FormRules } from 'naive-ui'
import { NButton, NCard, NDivider, NForm, NFormItem, NInput, useDialog, useMessage } from 'naive-ui'
import { ref } from 'vue'
import { useAuthStore, usePanelState, useUserStore } from '@/store'
import { logout } from '@/api'
import { router } from '@/router'
import { RoundCardModal, SvgIcon } from '@/components/common/'
import { updateInfo, updatePassword } from '@/api/system/user'
import { updateLocalUserInfo } from '@/utils/cmn'
import { t } from '@/locales'

const userStore = useUserStore()
const authStore = useAuthStore()
const panelState = usePanelState()
const ms = useMessage()
const dialog = useDialog()

const nickName = ref(authStore.userInfo?.name || '')
const isEditNickNameStatus = ref(false)
const formRef = ref<FormInst | null>(null)

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

async function logoutApi() {
  await logout()
  userStore.resetUserInfo()
  authStore.removeToken()
  panelState.removeState()
  ms.success(t('settingUserInfo.logoutSuccess'))
  router.push({ path: '/login' })
}

function handleSaveInfo() {
  updateInfo(nickName.value).then(({ code, msg }) => {
    if (code === 0) {
      updateLocalUserInfo()
      isEditNickNameStatus.value = false
    }
    else {
      ms.error(`${t('common.editFail')}:${msg}`)
    }
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
      if (code === 0) {
        // 成功
        updatePasswordModalState.value.show = false
        ms.success(t('common.success'))
      }
      else if (code === 0) {
        // 旧密码错误

      }
      else {
        // 其他错误
        ms.error(`${t('common.failed')}:${msg}`)
      }
    }).finally(() => {
      updatePasswordModalState.value.loading = false
    }).catch(() => {
      ms.error(t('common.serverError'))
    })
    console.log('aaa')
  })
}

function handleLogout() {
  dialog.warning({
    title: t('common.warning'),
    content: t('settingUserInfo.confirmLogoutText'),
    positiveText: t('common.confirm'),
    negativeText: t('common.cancel'),
    onPositiveClick: () => {
      logoutApi()
    },
  })
}
</script>

<template>
  <div class="bg-slate-200 rounded-[10px] p-[8px]">
    <NCard style="border-radius:10px" size="small">
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

          <NButton size="small" text type="info" @click="isEditNickNameStatus = !isEditNickNameStatus">
            {{ $t('common.edit') }}
          </NButton>
        </div>

        <div v-else class="flex items-center">
          <div class="max-w-[150px]">
            <NInput v-model:value="nickName" type="text" :placeholder="$t('common.inputPlaceholder')" />
          </div>
          <NButton size="small" quaternary type="info" @click="handleSaveInfo">
            {{ $t('common.save') }}
          </NButton>
        </div>
      </div>

      <NDivider style="margin: 10px 0;" dashed />
      <div>
        <NButton size="small" text type="info" @click="updatePasswordModalState.show = !updatePasswordModalState.show">
          {{ $t('settingUserInfo.updatePassword') }}
        </NButton>
      </div>
    </NCard>

    <NCard style="border-radius:10px" class="mt-[10px]" size="small">
      <NButton size="small" text type="error" @click="handleLogout">
        <template #icon>
          <SvgIcon icon="tabler:logout" />
        </template>
        {{ $t('settingUserInfo.logout') }}
      </NButton>
    </NCard>

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
