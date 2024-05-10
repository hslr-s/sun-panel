<script setup lang="ts">
import moment from 'moment'
import { onMounted, ref } from 'vue'
import { NAlert, NButton, NCard, NEllipsis, NForm, NFormItemRow, NInput, useMessage } from 'naive-ui'
import { getStatus as getStatusApi, login as loginApi, logout as logoutApi } from '@/api/panel/proAuth'
import { apiRespErrMsg, buildTimeString } from '@/utils/cmn'
import { ProCard, RoundCardModal, SvgIconOnline } from '@/components/common'
import { t } from '@/locales'
import { useProAuthorizeStore } from '@/store'

interface ProCardParam {
  nickName?: string
  proFlag?: string
  username?: string
  expired?: string
}

interface LoginModalForm {
  username?: string
  password?: string
}

const links: { [key: string]: string } = {
  proAuthSite: 'http://links.sun-panel.top/proAuthSite',
  proAuthSiteRegister: 'http://links.sun-panel.top/proAuthSiteRegister',
  proAuthSiteFAQ: 'http://links.sun-panel.top/proAuthSiteFAQ',
}

const loginModalShow = ref(false)
const loginModalForm = ref<LoginModalForm>({
  username: '',
  password: '',
})

const ms = useMessage()
const proAuthorizeStore = useProAuthorizeStore()

// const showGoToSiteButton = ref(false)
const loginLoading = ref(false)
const loginStatus = ref(false)
const proExpiration = ref('')
const proCardParam = ref<ProCardParam>()
const cardStyleSuccess = {
  backgroundColor: '#F0A020',
  color: 'white',
  textShadow: '1px 1px 1px #5d5c5c',
}
const cardStyleExpired = {
  backgroundColor: 'rgb(118, 124, 130)',
  color: 'white',
  textShadow: '1px 1px 1px #5d5c5c',
}
const cardStyle = ref(cardStyleExpired)

async function getStatus() {
  try {
    const { data } = await getStatusApi<ProAuth.GetStatusResp>()
    loginStatus.value = data.loginStatus
    proExpiration.value = data.proExpiration || ''
    if (isExpired(proExpiration.value))
      cardStyle.value = cardStyleExpired
    else
      cardStyle.value = cardStyleSuccess

    proCardParam.value = {
      username: data.username,
      nickName: data.nickName,
    }
  }
  catch (error) {
    console.error(error)
  }
}

async function login(username: string, password: string) {
  loginLoading.value = true

  await loginApi(username, password).then(() => {
    getStatus()
    loginModalShow.value = false
  }).catch((error) => {
    if (error.code === 1414) {
      ms.error(t('apps.proAuth.serverError'))
      return
    }
    apiRespErrMsg(error)
  })
  proAuthorizeStore.refreshProAuthStatus()
  loginLoading.value = false
}

async function logout() {
  try {
    await logoutApi()
    getStatus()
    loginModalShow.value = false
    proAuthorizeStore.refreshProAuthStatus()
  }
  catch (error) {
    apiRespErrMsg(error)
  }
}

function handleLogin() {
  loginModalShow.value = true
}

function handleSubmit() {
  if (loginModalForm.value.username === '' || loginModalForm.value.password === '') {
    ms.error(t('apps.proAuth.usernameWithPasswordIsNotEmpty'))
    return
  }

  login(loginModalForm.value.username || '', loginModalForm.value.password || '')
}

function openToLink(link: string) {
  if (links[link])
    window.open(links[link])
}

function isExpired(expired: string) {
  const momentDate1 = moment()
  const momentDate2 = moment(expired)

  if (expired === '')
    return true

  // 比较日期
  if (momentDate1.isAfter(momentDate2))
    return true
  else
    return false
}

onMounted(() => {
  getStatus()
})
</script>

<template>
  <div class="bg-slate-200 dark:bg-zinc-900 p-2 h-full">
    <NCard style="border-radius:10px" size="small">
      <div class="flex justify-center">
        <div class="m-5 w-[300px]">
          <ProCard :card-style="cardStyle" class="dark:brightness-90">
            <template #nickName>
              <template v-if="loginStatus">
                {{ proCardParam?.nickName }}
              </template>
              <div v-else class="cursor-pointer">
                <span @click="handleLogin">{{ t("apps.proAuth.notLogin") }}</span>
              </div>
            </template>
            <template #proFlag>
              PRO
            </template>
            <template #username>
              <template v-if="loginStatus">
                <div class="flex items-center">
                  <span>
                    {{ proCardParam?.username }}
                  </span>
                  <span class="cursor-pointer ml-2" :title="t('apps.proAuth.logout')" @click="logout">
                    <SvgIconOnline icon="tabler:logout" class="w-[20px] h-[20px]" :title="t('apps.proAuth.logout')" />
                  </span>
                </div>
              </template>
              <template v-else>
                <NButton size="small" ghost type="success" @click="handleLogin">
                  <NEllipsis style="max-width: 200px">
                    <span class="text-white">
                      {{ t("apps.proAuth.loginSunPanel") }}
                    </span>
                  </NEllipsis>
                </NButton>
              </template>
            </template>
            <template #expired>
              <template v-if="loginStatus">
                <template v-if="!proExpiration || proExpiration === ''">
                  <NButton type="primary" dashed size="tiny" @click="openToLink('proAuthSite')">
                    {{ t("apps.proAuth.card.upgradePRO") }}
                  </NButton>
                </template>
                <template v-else-if="isExpired(proExpiration)">
                  {{ t("apps.proAuth.card.expiredText", { date: buildTimeString(proExpiration, 'YYYY-MM-DD') }) }}
                </template>
                <template v-else>
                  {{ t("apps.proAuth.card.notExpiredText", { date: buildTimeString(proExpiration, 'YYYY-MM-DD') }) }}
                </template>
              </template>
            </template>
          </ProCard>

          <div class="mt-2 flex justify-center">
            <NButton size="small" style="margin-right: 5px;" @click="openToLink('proAuthSiteFAQ')">
              {{ t("apps.proAuth.FAQButton") }}
            </NButton>

            <NButton size="small" @click="openToLink('proAuthSite')">
              {{ t("apps.proAuth.goToSunPanelAuthSite") }}
            </NButton>
          </div>
        </div>
      </div>

      <NAlert :bordered="false" :show-icon="false">
        {{ t("apps.proAuth.textDescription") }}
      </NAlert>

      <br>

      <!-- <NAlert type="info" title="常见问题" :bordered="false">
      <p>
        Q：为什么会频繁让我手动登录
      </p>
      <p>
        A：可能原因：授权账号在其他客户端登录、修改了账号密码、程序升级重装或者其他破解原因造成
      </p>
      <br>
      <p>
        Q：可以离线授权吗
      </p>
      <p>
        A：由于Docker容器可移植特性，暂没有好的方案离线授权，如果日后找到解决方案，会增加离线授权，目前必须在线授权
      </p>
    </NAlert> -->

      <RoundCardModal v-model:show="loginModalShow" style="width:400px" :title="t('apps.proAuth.loginSunPanel')">
        <div class="flex items-center">
          <div>
            <NButton size="tiny" type="info" @click="openToLink('proAuthSiteRegister')">
              {{ t("apps.proAuth.registerAccount") }}
            </NButton>
          </div>
        </div>

        <div @keydown.enter="handleSubmit">
          <NForm>
            <NFormItemRow>
              <NInput v-model:value="loginModalForm.username" :placeholder="t('apps.proAuth.username')" />
            </NFormItemRow>
            <NFormItemRow>
              <NInput v-model:value="loginModalForm.password" type="password" :placeholder="t('apps.proAuth.password')" />
            </NFormItemRow>
          </NForm>
          <NButton type="primary" :loading="loginLoading" block secondary strong @click="handleSubmit">
            {{ t("apps.proAuth.login") }}
          </NButton>
        </div>
      </RoundCardModal>
    </NCard>
  </div>
</template>
