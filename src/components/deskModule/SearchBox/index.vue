<script setup lang="ts">
import { NInput } from 'naive-ui'
import { SvgIcon } from '@/components/common'
</script>

<template>
  <div class="w-full">
    <NInput size="large" class="background" round placeholder="输入搜索内容">
      <template #prefix>
        百度
      </template>
      <template #suffix>
        <SvgIcon icon="iconamoon:search-fill" />
      </template>
    </NInput>
  </div>
</template>

<style scoped>
.background{
  background-color: #ffffff78;
}
</style>
