<script lang="ts" setup>
import { computed, ref } from 'vue'
import type { UploadFileInfo } from 'naive-ui'
import { NAvatar, NButton, NInput, NUpload, useDialog, useMessage } from 'naive-ui'
// import type { Language, Theme } from '@/store/modules/app/helper'
import type { Theme } from '@/store/modules/app/helper'
import { SvgIcon } from '@/components/common'
import { useAppStore, useAuthStore, useUserStore } from '@/store'
import type { UserInfo } from '@/store/modules/user/helper'
import { t } from '@/locales'
import { UserUpdateInfo, logout } from '@/api'
import { router } from '@/router'

// import defaultAvatar from '@/assets/userDefaultAvatar.png'

const appStore = useAppStore()
const userStore = useUserStore()
const authStore = useAuthStore()

const ms = useMessage()
const dialog = useDialog()

const theme = computed(() => appStore.theme)

const userInfo = computed(() => userStore.userInfo)

const avatar = ref(userInfo.value.headImage ?? '')

const name = ref(userInfo.value.name ?? '')

// const description = ref(userInfo.value.description ?? '')

// const language = computed({
//   get() {
//     return appStore.language
//   },
//   set(value: Language) {
//     appStore.setLanguage(value)
//   },
// })

const themeOptions: { label: string; key: Theme; icon: string }[] = [
  {
    label: 'Auto',
    key: 'auto',
    icon: 'ri:contrast-line',
  },
  {
    label: 'Light',
    key: 'light',
    icon: 'ri:sun-foggy-line',
  },
  {
    label: 'Dark',
    key: 'dark',
    icon: 'ri:moon-foggy-line',
  },
]

// const languageOptions: { label: string; key: Language; value: Language }[] = [
//   { label: '简体中文', key: 'zh-CN', value: 'zh-CN' },
//   { label: '繁體中文', key: 'zh-TW', value: 'zh-TW' },
//   { label: 'English', key: 'en-US', value: 'en-US' },
//   { label: '한국어', key: 'ko-KR', value: 'ko-KR' },
// ]

function updateUserInfo(options: Partial<UserInfo>) {
  userStore.updateUserInfo({
    headImage: userInfo.value.headImage,
    name: options.name as string,
  })
  UserUpdateInfo(userInfo.value.headImage as string, options.name as string)
  ms.success(t('common.success'))
}

// function uploadHeadImage(options: Partial<UserInfo>) {

// }

function onLogoutClick() {
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

async function logoutApi() {
  await logout()
  userStore.resetUserInfo()
  authStore.removeToken()
  ms.success('您已经安全退出，期待与你再次相见！')
  router.push({ path: '/login' })
}

const handleFinish = ({
  file,
  event,
}: {
  file: UploadFileInfo
  event?: ProgressEvent
}) => {
  const res = JSON.parse((event?.target as XMLHttpRequest).response)
  // {
  //     "code": 0,
  //     "data": {
  //         "imageUrl": "/uploads/2023/5/12/c94306cfdf37fe7844753cd98fd57aaf.jpg"
  //     },
  //     "msg": "OK"
  // }
  const imageUrl = res.data.imageUrl
  userStore.updateUserInfo({
    headImage: imageUrl,
    name: userInfo.value.name,
  })
  avatar.value = imageUrl
  UserUpdateInfo(imageUrl, userInfo.value.name || '')

  return file
}
</script>

<template>
  <div class="p-4 space-y-5 min-h-[200px]">
    <div class="space-y-6">
      <div class="flex items-center space-x-4">
        <span class="flex-shrink-0 w-[100px]">账号/邮箱</span>
        <div class="w-[200px]">
          <NInput v-model:value="userInfo.username" :disabled="true" readonly placeholder="" />
        </div>
      </div>
      <div class="flex items-center space-x-4">
        <span class="flex-shrink-0 w-[100px]">密码</span>
        <div class="w-[200px]">
          <NButton @click="router.push({ path: '/resetPassword', query: { u: userInfo.username } })">
            重置密码
          </NButton>
        </div>
      </div>
      <div class="flex items-center space-x-4">
        <span class="flex-shrink-0 w-[100px]">{{ $t('setting.avatar') }}</span>
        <div class="w-[50px]">
          <NAvatar
            round
            :size="48"
            :src="avatar ? avatar : ''"
          />
        </div>
        <NUpload
          action="/api/file/uploadImg"
          name="imgfile"
          :headers="{
            token: authStore.token as string,
          }"
          @finish="handleFinish"
        >
          <NButton size="tiny" text type="primary">
            上传文件
          </NButton>
        </NUpload>
        <!-- <NButton size="tiny" text type="primary" @click="uploadHeadImage({ avatar })">
          选择文件
        </NButton> -->
      </div>
      <!-- <div class="flex items-center space-x-4">
        <span class="flex-shrink-0 w-[100px]">{{ $t('setting.avatarLink') }}</span>
        <div class="flex-1">
          <NInput v-model:value="avatar" placeholder="" />
        </div>
        <NButton size="tiny" text type="primary" @click="updateUserInfo({ avatar })">
          {{ $t('common.save') }}
        </NButton>
      </div> -->

      <div class="flex items-center space-x-4">
        <span class="flex-shrink-0 w-[100px]">{{ $t('setting.name') }}</span>
        <div class="w-[200px]">
          <NInput v-model:value="name" placeholder="" />
        </div>
        <NButton size="tiny" text type="primary" @click="updateUserInfo({ name })">
          {{ $t('common.save') }}
        </NButton>
      </div>
      <div class="flex items-center space-x-4">
        <span class="flex-shrink-0 w-[100px]">{{ $t('setting.theme') }}</span>
        <div class="flex flex-wrap items-center gap-4">
          <template v-for="item of themeOptions" :key="item.key">
            <NButton
              size="small"
              :type="item.key === theme ? 'primary' : undefined"
              @click="appStore.setTheme(item.key)"
            >
              <template #icon>
                <SvgIcon :icon="item.icon" />
              </template>
            </NButton>
          </template>
        </div>
      </div>
      <!-- <div class="flex items-center space-x-4">
        <span class="flex-shrink-0 w-[100px]">{{ $t('setting.language') }}</span>
        <div class="flex flex-wrap items-center gap-4">
          <NSelect
            style="width: 140px"
            :value="language"
            :options="languageOptions"
            @update-value="value => appStore.setLanguage(value)"
          />
        </div>
      </div> -->

      <div class="flex items-center space-x-4">
        <span class="flex-shrink-0 w-[100px]" />
        <NButton size="small" type="error" strong secondary @click="onLogoutClick">
          <template #icon>
            <SvgIcon icon="oi:account-logout" />
          </template>
          安全退出账号
        </NButton>
      </div>
    </div>
  </div>
</template>
