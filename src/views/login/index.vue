<script setup lang="ts">
import { NButton, NCard, NForm, NFormItem, NGradientText, NInput, NSelect, useMessage } from 'naive-ui'
import { onMounted, ref } from 'vue'
import { login } from '@/api'
import { useAppStore, useAuthStore } from '@/store'
import { SvgIcon } from '@/components/common'
import { router } from '@/router'
import { t } from '@/locales'
import { languageOptions } from '@/utils/defaultData'
import type { Language } from '@/store/modules/app/helper'
import { getSiteStting } from '@/api/openness'
import { apiRespErrMsg, setTitle } from '@/utils/cmn'

// const userStore = useUserStore()
const authStore = useAuthStore()
const appStore = useAppStore()
const ms = useMessage()
const loading = ref(false)
const languageValue = ref<Language>(appStore.language)
const siteSetting = ref<GlobalSetting.SiteSetting>({
  favicoUrl: '',
  title: 'Sun-Panel',
  loginBackgroundUrl: '',
})

const form = ref<Login.LoginReqest>({
  username: '',
  password: '',
})

const loginPost = async () => {
  loading.value = true
  try {
    const res = await login<Login.LoginResponse>(form.value)
    if (res.code === 0) {
      authStore.setToken(res.data.token)
      authStore.setUserInfo(res.data)

      setTimeout(() => {
        ms.success(`Hi ${res.data.name},${t('login.welcomeMessage')}`)
        loading.value = false
        router.push({ path: '/' })
      }, 500)
    }
    else {
      loading.value = false
    }
  }
  catch (error) {
    apiRespErrMsg(error)
    loading.value = false
  }
}

function handleSubmit() {
  // 点击登录按钮触发
  loginPost()
}

function handleChangeLanuage(value: Language) {
  languageValue.value = value
  appStore.setLanguage(value)
}

async function getConfig() {
  try {
    const { data } = await getSiteStting<GlobalSetting.SiteSetting>()
    if (data.title === '')
      data.title = 'Sun-Panel'

    appStore.setSiteSetting(data)
    siteSetting.value = data
    setTitle(data.title)
  }
  catch (error) {
    apiRespErrMsg(error)
  }
}

onMounted(() => {
  setTitle(appStore.siteSetting.title)
  getConfig()
})
</script>

<template>
  <div>
    <div class="background" :style="{ backgroundImage: `url('${appStore.siteSetting.loginBackgroundUrl}')` }" />
    <div class="login-container">
      <NCard class="login-card" style="border-radius: 20px;">
        <div class="mb-5 flex items-center justify-end">
          <div class="mr-2">
            <SvgIcon icon="ion-language" style="width: 20;height: 20;" />
          </div>
          <div class="min-w-[100px]">
            <NSelect
              v-model:value="languageValue" size="small" :options="languageOptions"
              @update-value="handleChangeLanuage"
            />
          </div>
        </div>

        <div class="login-title  ">
          <NGradientText :size="30" type="success" class="!font-bold">
            {{ appStore.siteSetting.title }}
          </NGradientText>
        </div>
        <NForm :model="form" label-width="100px" @keydown.enter="handleSubmit">
          <NFormItem>
            <NInput v-model:value="form.username" :placeholder="$t('login.usernamePlaceholder')">
              <template #prefix>
                <SvgIcon icon="ph:user-bold" />
              </template>
            </NInput>
          </NFormItem>

          <NFormItem>
            <NInput v-model:value="form.password" type="password" :placeholder="$t('login.passwordPlaceholder')">
              <template #prefix>
                <SvgIcon icon="mdi:password-outline" />
              </template>
            </NInput>
          </NFormItem>

          <!-- <NFormItem v-if="isShowCaptcha">
          <div class="w-[120px] h-[34px] mr-[20px] rounded border flex cursor-pointer">
            <Captcha ref="captchaRef" src="/api/captcha/getImage" />
          </div>
          <NInput v-model:value="form.vcode" type="text" placeholder="请输入图像验证码" />
        </NFormItem> -->
          <NFormItem style="margin-top: 10px">
            <NButton type="primary" block :loading="loading" @click="handleSubmit">
              {{ $t('login.loginButton') }}
            </NButton>
          </NFormItem>

          <!-- <div class="flex justify-end">
          <NButton v-if="isShowRegister" quaternary type="info" class="flex" @click="$router.push({ path: '/register' })">
            注册
          </NButton>
          <NButton quaternary type="info" class="flex" @click="$router.push({ path: '/resetPassword' })">
            忘记密码?
          </NButton>
        </div> -->

          <div class="flex justify-center text-slate-300">
            Powered By <a
              href="https://github.com/hslr-s/sun-panel" target="_blank"
              class="ml-[5px] text-slate-500"
            >Sun-Panel</a>
          </div>
        </NForm>
      </NCard>
    </div>
  </div>
</template>

<style>
.background {
  position: absolute;
  width: 100%;
  height: 100%;
  background-size: cover;
  background-position: center;
  background-color: #f2f6ff;
}

.dark .background {
  filter: brightness(50%);
  background-color: rgb(82, 82, 82);
}

.login-container {
  padding: 20px;
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100vh;
}

@media (min-width: 600px) {
  .login-card {
    width: auto;
    margin: 0px 10px;
  }

  .login-button {
    width: 100%;
  }
}

.login-card {
  margin: 20px;
  min-width: 400px;
}

.login-title {
  text-align: center;
  margin: 20px;
}
</style>
