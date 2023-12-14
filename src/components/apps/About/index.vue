<script setup lang="ts">
import { NDivider, NGradientText } from 'naive-ui'
import { onMounted, ref } from 'vue'
import { get } from '@/api/system/about'
import srcSvglogo from '@/assets/logo.svg'
import srcGitee from '@/assets/about_image/gitee.png'
import srcGithub from '@/assets/about_image/github.png'
import srcDocker from '@/assets/about_image/docker.png'
import srcQQGroupQR from '@/assets/about_image/qq_group_qr.jpg'
import { RoundCardModal } from '@/components/common'

interface Version {
  versionName: string
  versionCode: number
}

const versionName = ref('')
const qqGroupQRShow = ref(false)

onMounted(() => {
  get<Version>().then((res) => {
    if (res.code === 0)
      versionName.value = res.data.versionName
  })
})
</script>

<template>
  <div>
    <div>
      <div class="flex flex-col items-center justify-center">
        <img :src="srcSvglogo" width="100" height="100" alt="">
        <div class="text-3xl font-semibold">
          {{ $t('common.appName') }}
        </div>
        <div class="text-xl">
          <NGradientText type="info">
            <a href="https://github.com/hslr-s/sun-panel/releases" class="font-semibold" title="点此查看更新说明" target="_blank">v{{ versionName }}</a>
          </NGradientText>
        </div>
      </div>
    </div>
    <NDivider />
    <div class="flex flex-col items-center justify-center text-base">
      <div>
        建议反馈：<a href="https://github.com/hslr-s/sun-panel/issues" target="_blank" class="link">Github Issues</a>
      </div>

      <div>
        QQ交流群：<a href="http://qm.qq.com/cgi-bin/qm/qr?_wv=1027&k=_I9WIoJn1roIdoaAqelSj9qClLKlXIa1&authKey=GfsQP2GagHnus0jMc7U8Sm6VhWjtsipXUzCHbFwQsGyHMgmYWx6ZbAP%2Bhut%2B4D6N&noverify=0&group_code=276594668" target="_blank" class="link">276594668</a>
        |
        <span class="link cursor-pointer" @click="qqGroupQRShow = !qqGroupQRShow">
          二维码(推荐)
        </span>
      </div>

      <div>
        开发者：<a href="https://blog.enianteam.com/u/sun/content/11" target="_blank" class="link">红烧猎人</a> | <a href="https://github.com/hslr-s/sun-panel/blob/master/doc/donate.md" target="_blank" class="text-red-600 hover:text-red-900">🧧打赏</a>
      </div>

      <div class="flex mt-[10px]">
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
