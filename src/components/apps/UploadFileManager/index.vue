<script setup lang="ts">
import { NButton, NButtonGroup, NCard, NEllipsis, NImage, NImageGroup, NSpace, useDialog, useMessage } from 'naive-ui'
import { onMounted, ref } from 'vue'
import { deletes, getList } from '@/api/system/file'
import { SvgIcon } from '@/components/common'
import { copyToClipboard, timeFormat } from '@/utils/cmn'
import { t } from '@/locales'

const imageList = ref<File.Info[]>([])
const ms = useMessage()
const dialog = useDialog()

async function getFileList() {
  const { data } = await getList<Common.ListResponse<File.Info[]>>()
  imageList.value = data.list
}

async function copyImageUrl(text: string) {
  const res = await copyToClipboard(text)
  if (res)
    ms.success('复制成功')

  else
    ms.error('复制失败')
}

function handleDelete(id: number) {
  dialog.warning({
    title: t('common.warning'),
    content: '你确定删除此图片吗？',
    positiveText: t('common.confirm'),
    negativeText: t('common.cancel'),
    onPositiveClick: () => {
      deletesImges(id)
    },
  })
}

async function deletesImges(id: number) {
  try {
    const { code, msg } = await deletes([id])
    if (code === 0) {
      getFileList()
      ms.success(t('common.success'))
    }
    else {
      ms.error(`${t('common.failed')}:${msg}`)
    }
  }
  catch (error) {
    ms.error(t('common.failed'))
  }
}

onMounted(() => {
  getFileList()
})
</script>

<template>
  <div>
    <div class="flex justify-center">
      <NImageGroup>
        <NSpace>
          <div v-for=" item, index in imageList" :key="index">
            <NCard style="width: 150px;" size="small">
              <template #cover>
                <div class="card transparent-grid">
                  <NImage :lazy="true" width="100%" style="width: auto;object-fit: contain;height: 100%;" :src="item.src" />
                </div>
              </template>
              <template #footer>
                <span class="text-xs">
                  <NEllipsis>
                    {{ item.fileName }}
                  </NEllipsis>
                </span>
                <div class="flex justify-center mt-[10px]">
                  <NButtonGroup>
                    <NButton size="tiny" tertiary style="cursor: pointer;" @click="copyImageUrl(item.src)">
                      <template #icon>
                        <SvgIcon icon="ion-copy" />
                      </template>
                    </NButton>
                    <NButton size="tiny" tertiary style="cursor: pointer;" :title="timeFormat(item.createTime)">
                      <template #icon>
                        <SvgIcon icon="mdi-information-box-outline" />
                      </template>
                    </NButton>
                    <NButton size="tiny" tertiary type="error" style="cursor: pointer;" @click="handleDelete(item.id as number)">
                      <template #icon>
                        <SvgIcon icon="material-symbols-delete" />
                      </template>
                    </NButton>
                  </NButtonGroup>
                </div>
              </template>
            </NCard>
          </div>
        </NSpace>
      </NImageGroup>
    </div>
  </div>
</template>

<style scoped>
.card {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100px;
}

.transparent-grid {
  background-image: linear-gradient(45deg, #e6e4e4 25%, transparent 25%, transparent 75%, #e6e4e4 75%),
    linear-gradient(45deg, #e6e4e4 25%, transparent 25%, transparent 75%, #e6e4e4 75%);
  background-size: 16px 16px;
  background-position: 0 0, 8px 8px;
}
</style>
