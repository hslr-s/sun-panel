<script setup lang="ts">
import { NBadge, NButton, NCheckbox, NDivider, NFlex, NGradientText, NTag } from 'naive-ui'
import { computed, onMounted, ref } from 'vue'
import { checkNewVersion as checkNewVersionApi, get } from '@/api/system/about'
import { useAppStore, useAuthStore } from '@/store'
import srcSvglogo from '@/assets/logo.svg'
import srcGitee from '@/assets/about_image/gitee.png'
import srcGithub from '@/assets/about_image/github.png'
import srcDocker from '@/assets/about_image/docker.png'
import srcBilibili from '@/assets/about_image/bilibili.png'
import srcYoutube from '@/assets/about_image/youtube.png'
import srcQQGroupQR from '@/assets/about_image/qq_group_qr2.png'
import { RoundCardModal, STip } from '@/components/common'
import { message } from '@/utils/cmn/apiMessage'
import { t } from '@/locales'
import { buildTimeString } from '@/utils/cmn'
import { useDomainLinks } from '@/hooks/useDomainLinks'

interface Version {
  versionName: string
  versionCode: number
}

const appStore = useAppStore()
const uauthStore = useAuthStore()
const versionName = ref('')
const qqGroupQRShow = ref(false)
const frontVersion = import.meta.env.VITE_APP_VERSION || 'unknown'
const isCheckBeta = ref(false)
const newVersionInfo = ref<About.NewVersionInfo | null>(null)
const checkLoading = ref(false)
const isCheckFailed = ref(false)
const domainLinks = useDomainLinks()
const isAdmin = computed(() => uauthStore.userInfo?.role === 1)

// 新增判断是否为beta版本的函数
function isBetaVersion(version: string): boolean {
  return version.includes('-beta')
}

function link(c: string): string {
  return domainLinks.linkWithLang(c)
}

function getUpdateLogLink(version: string): string {
  let url = isBetaVersion(version) ? domainLinks.linkWithLangParam('update_log_beta_lang') : domainLinks.linkWithLangParam('update_log_lang')
  url += `&version=${version}`
  return url
}

function getWebsiteUrl(): string {
  if (appStore.language === 'zh-CN')
    return 'https://doc.sun-panel.top/zh_cn/'

  return 'http://doc.sun-panel.top'
}

async function checkNewVersion() {
  checkLoading.value = true
  await checkNewVersionApi<About.NewVersionInfo | null>(isCheckBeta.value).then(({ data }) => {
    newVersionInfo.value = data
    isCheckFailed.value = false
  }).catch(() => {
    isCheckFailed.value = true
    message.error(t('common.failed'))
  })
  checkLoading.value = false
}

onMounted(() => {
  get<Version>().then((res) => {
    if (res.code === 0)
      versionName.value = res.data.versionName
  })
  if (isAdmin.value)
    checkNewVersion()
})
</script>

<template>
  <div class="app-background h-full overflow-auto py-5">
    <div class="flex flex-col items-center justify-center">
      <img :src="srcSvglogo" width="80" height="80" alt="">
      <div class="text-3xl font-semibold">
        {{ $t('common.appName') }}
      </div>
      <div class="text-xl">
        <NGradientText type="info">
          <!-- 修改点击版本号跳转的逻辑 -->
          <a
            :href="getUpdateLogLink(versionName)"
            class="font-semibold"
            :title="$t('apps.about.viewUpdateLog')"
            target="_blank"
          >
            v{{ versionName }}
          </a>
        </NGradientText>
      </div>

      <div v-if="!checkLoading && isAdmin" class="mt-2">
        <NFlex v-if="newVersionInfo !== null" justify="center" align="center" style="margin:10px 0 ;" :title="t('apps.about.viewUpdateLog')">
          {{ t('apps.about.haveNewVersion') }}
          <NBadge value="New" processing :offset="[5, -5]">
            <NTag type="info">
              <!-- 修改检测新版本后跳转的逻辑 -->
              <a
                :href="getUpdateLogLink(newVersionInfo.version)"
                target="_blank"
                class="link"
              >
                v{{ newVersionInfo.version }}
              </a>
            </NTag>
          </NBadge>
          <template v-if="newVersionInfo.type === 'release'">
            ({{ buildTimeString(newVersionInfo.releaseTime, "YYYY-MM-DD") }})
          </template>
        </NFlex>

        <NFlex v-else-if="!isCheckFailed && newVersionInfo === null " justify="center" style="margin:10px 0 ;">
          <NTag type="success">
            {{ t('apps.about.alreadyLatestVersion') }}
          </NTag>
        </NFlex>

        <NFlex v-if="isCheckFailed" justify="center" style="margin:10px 0 ;">
          <NTag type="error">
            <a href="https://github.com/hslr-s/sun-panel/releases" target="_blank"> {{ t("apps.about.checkFailed") }}</a>
          </NTag>
        </NFlex>
      </div>

      <NFlex v-if="isAdmin" align="center" justify="center">
        <NCheckbox v-model:checked="isCheckBeta">
          Beta
        </NCheckbox>
        <NButton size="tiny" type="success" :loading="checkLoading" @click="checkNewVersion">
          {{ $t('apps.about.checkUpdate') }}
        </NButton>
        <STip>
          <div class="max-w-[300px]">
            <p>{{ t("apps.about.checkUploadTip1") }}</p>
            <br>
            <p>{{ t("apps.about.checkUploadTip2") }}</p>
            <br>
            <p>{{ t("apps.about.checkUploadTip3") }}</p>
          </div>
        </STip>
      </NFlex>
    </div>

    <NDivider style="margin:10px 0;padding: 0 10px">
      •
    </NDivider>
    <div class="flex flex-col items-center justify-center text-[13px]">
      <div>
        {{ $t('apps.about.website') }}
        <a :href="getWebsiteUrl()" target="_blank" class="link"> Sun-Panel.Top </a>
      </div>
      <div>
        {{ $t('apps.about.author') }}
        <a href="https://github.com/hslr-s" target="_blank" class="link"> 红烧猎人 </a>
        |
        <a :href="link('donate')" target="_blank" class="text-red-600 hover:text-red-900">
          {{ $t('apps.about.donate') }}
        </a>
      </div>
      <div>
        {{ $t('apps.about.issue') }}<a :href="link('faq')" target="_blank" class="link">{{ t('apps.about.faq') }}</a> | <a href="https://github.com/hslr-s/sun-panel/issues" target="_blank" class="link"> Github Issues</a>
      </div>
      <div>
        {{ $t('apps.about.discussions') }}<a href="https://github.com/hslr-s/sun-panel/discussions" target="_blank" class="link"> Github Discussions</a>
      </div>
      <div>
        {{ $t('apps.about.communicationGroup') }} <a :href="link('groupChat')" target="_blank" class="link"> {{ $t("apps.about.groupChat") }}</a>
      </div>

      <div class="flex mt-[10px] flex-wrap justify-center">
        <div class="flex items-center mx-[10px]">
          <img class="w-[20px] h-[20px] mr-[5px]" :src="srcGithub" alt="">
          <a href="https://github.com/hslr-s/sun-panel" target="_blank" class="link">Github</a>
        </div>
        <div class="flex items-center mx-[10px]">
          <img class="w-[20px] h-[20px] mr-[5px]" :src="srcGitee" alt="">
          <a href="https://gitee.com/hslr/sun-panel" target="_blank" class="link">Gitee</a>
        </div>
        <div class="flex items-center mx-[10px]">
          <img class="w-[20px] h-[20px] mr-[5px]" :src="srcDocker" alt="">
          <a href="https://hub.docker.com/r/hslr/sun-panel" target="_blank" class="link">Docker</a>
        </div>
        <div class="flex items-center mx-[10px]">
          <img class="w-[20px] h-[20px] mr-[5px]" :src="srcBilibili" alt="">
          <!-- <a href="https://space.bilibili.com/27407696/channel/collectiondetail?sid=2023810" target="_blank" class="link">Bilibili</a> -->
          <a href="https://space.bilibili.com/27407696/channel/collectiondetail?sid=2023810" target="_blank" class="link">Bilibili</a>
        </div>
        <div v-if="appStore.language !== 'zh-CN'" class="flex items-center mx-[10px]">
          <img class="w-[20px] h-[20px] mr-[5px]" :src="srcYoutube" alt="">
          <a href="https://www.youtube.com/channel/UCKwbFmKU25R602z6P2fgPYg" target="_blank" class="link">YouTube</a>
        </div>
      </div>

      <div class="mt-5">
        <NTag :bordered="false" size="small">
          {{ $t("apps.about.frontVersionText") }}: FV-{{ frontVersion }}
        </NTag>
      </div>

      <RoundCardModal v-model:show="qqGroupQRShow" title="交流群二维码" style="width: 300px;">
        <div class="text-center">
          - 如果失效请返回联系作者 -
        </div>
        <div class="flex justify-center">
          <img :src="srcQQGroupQR" class="h-[260px]">
        </div>
      </RoundCardModal>
    </div>
  </div>
</template>

<style>
.link{
    color:rgb(0, 89, 255)
}
</style>
