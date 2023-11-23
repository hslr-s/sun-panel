<script setup lang="ts">
import { NDivider } from 'naive-ui'
import { onMounted, ref } from 'vue'
import { get } from '@/api/system/about'
import srcSvglogo from '@/assets/logo.svg'

interface Version {
  versionName: string
  versionCode: number
}

const versionName = ref('')

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
        <div class="text-3xl">
          {{ $t('common.appName') }}
        </div>
        <div class="text-xl">
          v{{ versionName }}
        </div>
      </div>
    </div>
    <NDivider />
    <div class="flex flex-col items-center justify-center">
      <div class="text-lg">
        开发者： <a href="https://blog.enianteam.com/u/sun/content/11" target="_blank" class="link">红烧猎人</a>
      </div>
      <div class="text-lg">
        QQ交流群：<a href="http://qm.qq.com/cgi-bin/qm/qr?_wv=1027&k=_I9WIoJn1roIdoaAqelSj9qClLKlXIa1&authKey=GfsQP2GagHnus0jMc7U8Sm6VhWjtsipXUzCHbFwQsGyHMgmYWx6ZbAP%2Bhut%2B4D6N&noverify=0&group_code=276594668" target="_blank" class="link">276594668</a>
      </div>
      <div class="text-lg">
        项目地址：
        <a href="https://github.com/hslr-s/sun-panel" target="_blank" class="link">Github</a> |
        <!-- <a href="https://gitee.com/hslr/sun-panel/stargazers"><img src="https://gitee.com/hslr/sun-panel/badge/star.svg?theme=dark" alt="star"></a> -->
        <a href="https://gitee.com/hslr/sun-panel" target="_blank" class="link">Gitee</a> |
        <a href="https://hub.docker.com/r/hslr/sun-panel" target="_blank" class="link">Docker</a>
      </div>
    </div>
  </div>
</template>

<style>
.link{
    color:blue
}
</style>
