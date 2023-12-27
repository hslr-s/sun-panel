<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { getList } from '@/api/system/file'

import UploadFileManager from '@/components/apps/UploadFileManager/index.vue'

const imageList = ref<File.Info[]>([])

async function getFileList() {
  const { data } = await getList<Common.ListResponse<File.Info[]>>()
  console.log(data)
  imageList.value = data.list
}

onMounted(() => {
  getFileList()
})
</script>

<template>
  <div>
    <UploadFileManager />
  </div>
</template>

<style scoped>
.card {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100px;
  /* 或者其他容器高度 */
}

.transparent-grid {
  background-image: linear-gradient(45deg, #e6e4e4 25%, transparent 25%, transparent 75%, #e6e4e4 75%),
    linear-gradient(45deg, #e6e4e4 25%, transparent 25%, transparent 75%, #e6e4e4 75%);
  background-size: 16px 16px;
  background-position: 0 0, 8px 8px;
}
</style>
