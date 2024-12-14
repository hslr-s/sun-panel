<script setup lang="ts">
import { NButton, NInput, NRadioButton, NRadioGroup, NUpload } from 'naive-ui'
import type { UploadFileInfo } from 'naive-ui'
import { computed, defineProps } from 'vue'
// import { ItemIcon } from '@/components/common'
import { useAuthStore } from '@/store'
import { apiRespErrMsg } from '@/utils/cmn'
import { t } from '@/locales'
import { FileType } from '@/enums'

const props = defineProps<{
  itemIcon: Panel.ItemIcon | null
}>()
const emit = defineEmits<{
  (e: 'update:itemIcon', visible: Panel.ItemIcon): void // 定义修改父组件（prop内）的值的事件
}>()
const authStore = useAuthStore()

// 默认图标背景色
// const defautSwatchesBackground = [
//   '#00000000',
//   '#000000',
//   '#ffffff',
//   '#18A058',
//   '#2080F0',
//   '#F0A020',
//   'rgba(208, 48, 80, 1)',
//   '#C418D1FF',
// ]

const initData: Panel.ItemIcon = {
  itemType: 2,
  backgroundColor: '#2a2a2a6b',
}

const itemIconInfo = computed({
  get() {
    const v = {
      ...initData,
      ...props.itemIcon,
      backgroundColor: props.itemIcon?.backgroundColor || initData.backgroundColor,
    }
    return v
  },
  set() {
    handleChange()
  },
})

function handleIconTypeRadioChange(type: number) {
  // checkedValueRef.value = type
  itemIconInfo.value.itemType = type
  handleChange()
}

function handleChange() {
  emit('update:itemIcon', itemIconInfo.value || null)
}

// watch(() => itemIconInfo.value.itemType, (newVal: number) => {
//   console.log('itemIconInfo.itemType change', newVal)
//   handleChange()
// }, {
//   deep: true,
// })

const handleUploadFinish = ({
  file,
  event,
}: {
  file: UploadFileInfo
  event?: ProgressEvent
}) => {
  const res = JSON.parse((event?.target as XMLHttpRequest).response)
  if (res.code === 0) {
    const imageUrl = res.data.imageUrl
    itemIconInfo.value.src = imageUrl
    emit('update:itemIcon', itemIconInfo.value || null)
  }
  else {
    apiRespErrMsg(res)
    // ms.error(`${t('common.uploadFail')}:${res.msg}`)
  }

  return file
}
</script>

<template>
  <div class="w-full">
    <div class="mb-[10px] flex items-center">
      <NRadioGroup
        v-model:value="itemIconInfo.itemType" name="iconType"
        @update-value="handleIconTypeRadioChange"
      >
        <NRadioButton :value="1">
          {{ $t('common.text') }}
        </NRadioButton>

        <NRadioButton :value="2">
          {{ $t('common.image') }}
        </NRadioButton>

        <NRadioButton :value="3">
          {{ t('iconItem.onlineIcon') }}
        </NRadioButton>
      </NRadioGroup>
    </div>

    <div class="w-full flex">
      <!-- 文字 -->
      <div v-if="itemIconInfo.itemType === 1" class="items-center">
        <div class="w-auto mr-2 whitespace-nowrap">
          {{ t('iconItem.icon.textContent') }}
        </div>
        <NInput v-model:value="itemIconInfo.text" :maxlength="10" type="text" show-count clearable @input="handleChange" />
      </div>

      <!-- 在线图标 -->
      <div v-if="itemIconInfo.itemType === 3" class="items-center">
        <div class="w-auto mr-2">
          {{ t('iconItem.icon.iconName') }}
        </div>
        <div class="flex">
          <NInput v-model:value="itemIconInfo.text" type="text" class="mr-1" :placeholder="$t('iconItem.inputIconName')" clearable @input="handleChange" />
          <NButton quaternary type="info">
            <a target="_blank" href="https://icon-sets.iconify.design/">{{ $t('iconItem.onlineIconLibrary') }}</a>
          </NButton>
        </div>
      </div>

      <!-- 图片 -->
      <div v-if="itemIconInfo.itemType === 2" class="items-center w-full">
        <div class="w-auto mr-2 whitespace-nowrap">
          {{ t('iconItem.icon.imageUrl') }}
        </div>
        <div class="flex w-full">
          <NInput
            v-model:value="itemIconInfo.src"
            class="mr-1"
            type="text"
            style="width:300px"
            :placeholder="$t('iconItem.inputIconUrlOrUpload')"
            clearable
            @input="handleChange"
          />
          <NUpload
            action="/api/file/uploadImg"
            :show-file-list="false"
            name="imgfile"
            :headers="{
              token: authStore.token as string,
            }"
            :data="{ type: FileType.CardIcon }"
            @finish="handleUploadFinish"
          >
            <NButton size="small">
              {{ $t('iconItem.selectUpload') }}
            </NButton>
          </NUpload>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.transparent-grid {
    background-image: linear-gradient(45deg, #fff 25%, transparent 25%, transparent 75%, #fff 75%),
                      linear-gradient(45deg, #fff 25%, transparent 25%, transparent 75%, #fff 75%);
    background-size: 16px 16px;
    background-position: 0 0, 8px 8px;
}
</style>
