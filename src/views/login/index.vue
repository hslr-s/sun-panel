<script setup lang="ts">
import { NButton, NCard, NForm, NFormItem, NGradientText, NInput, useMessage } from 'naive-ui'
import { ref } from 'vue'
import { login } from '@/api'
import { useAuthStore, useUserStore } from '@/store'
import { router } from '@/router'
import { Captcha, SvgIcon } from '@/components/common'

const userStore = useUserStore()
const authStore = useAuthStore()
const ms = useMessage()
const isShowCaptcha = ref<boolean>(false)
const isShowRegister = ref<boolean>(false)

const captchaRef = ref()

const form = ref<Login.LoginReqest>({
  username: '',
  password: '',
})

const loginPost = async () => {
  const res = await login<Login.LoginResponse>(form.value)

  userStore.updateUserInfo(res.data)

  if (res.code === 0) {
    authStore.setToken(res.data.token)
    ms.success(`Hi ${res.data.name}，欢迎回来!`)
    router.push({ path: '/' })
  }
  else {
    captchaRef.value.refresh()
  }
}

function handleSubmit() {
  // 点击登录按钮触发
  loginPost()
}
</script>

<template>
  <div class="login-container">
    <NCard class="login-card">
      <div class="login-title">
        <NGradientText :size="30" type="success">
          {{ $t('common.appName') }}
        </NGradientText>
      </div>
      <NForm :model="form" label-width="100px">
        <NFormItem>
          <NInput v-model:value="form.username" placeholder="请输入邮箱地址作为账号">
            <template #prefix>
              <SvgIcon icon="ph:user-bold" />
            </template>
          </NInput>
        </NFormItem>

        <NFormItem>
          <NInput v-model:value="form.password" type="password" placeholder="请输入密码">
            <template #prefix>
              <SvgIcon icon="mdi:password-outline" />
            </template>
          </NInput>
        </NFormItem>

        <NFormItem v-if="isShowCaptcha">
          <div class="w-[120px] h-[34px] mr-[20px] rounded border flex cursor-pointer">
            <Captcha ref="captchaRef" src="/api/captcha/getImage" />
          </div>
          <NInput v-model:value="form.vcode" type="text" placeholder="请输入图像验证码" />
        </NFormItem>
        <NFormItem style="margin-top: 10px">
          <NButton type="primary" block @click="handleSubmit">
            登录
          </NButton>
        </NFormItem>

        <div class="flex justify-end">
          <NButton v-if="isShowRegister" quaternary type="info" class="flex" @click="$router.push({ path: '/register' })">
            注册
          </NButton>
          <!-- <NButton quaternary type="info" class="flex" @click="$router.push({ path: '/resetPassword' })">
            忘记密码?
          </NButton> -->
        </div>
      </NForm>
    </NCard>
  </div>
</template>

  <style>
    .login-container {
        padding: 20px;
        display: flex;
        justify-content: center;
        align-items: center;
        height: 100vh;
        background-color: #f2f6ff;
    }

    /* 夜间模式 */
    .dark .login-container{
      background-color: rgb(43, 43, 43);
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
        min-width:400px;
    }

  .login-title{
    text-align: center;
    margin: 20px;
  }
  </style>
