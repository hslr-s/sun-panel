<script setup lang="ts">
import { NDivider, NGradientText } from 'naive-ui'
import { onMounted, ref } from 'vue'
import { get } from '@/api/system/about'

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
    <div class="text-3xl">
      <NGradientText type="danger">
        Sun-Panel
      </NGradientText>
      v{{ versionName }}
    </div>
    <NDivider />
    <div class="text-lg">
      开发者： <a href="https://blog.enianteam.com/u/sun/content/11" target="_blank" class="link">红烧猎人</a>
    </div>
    <div class="text-lg">
      项目开源地址：
      <a href="https://github.com/hslr-s/sun-panel" target="_blank" class="link">Github</a> |
      <a href="https://gitee.com/hslr/sun-panel" target="_blank" class="link">Gitee</a>
    </div>
  </div>
</template>

<style>
.link{
    color:blue
}
</style>
