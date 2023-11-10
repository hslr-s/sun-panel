<script setup lang="ts">
import { NButton, NColorPicker, NInput, NRadio, NUpload, useMessage } from 'naive-ui'
import type { UploadFileInfo } from 'naive-ui'
import { defineProps, ref } from 'vue'
import { ItemIcon } from '@/components/common'
import { useAuthStore } from '@/store'

const props = defineProps<{
  itemIcon: Panel.ItemIcon | null
}>()
const emit = defineEmits<{
  (e: 'update:itemIcon', visible: Panel.ItemIcon): void // 定义修改父组件（prop内）的值的事件
}>()
const authStore = useAuthStore()
const ms = useMessage()
const checkedValueRef = ref<number | null>(props.itemIcon?.itemType || 1)

// 默认图标背景色
const defautSwatchesBackground = [
  '#000',
  '#18A058',
  '#2080F0',
  '#F0A020',
  'rgba(208, 48, 80, 1)',
]

const initData: Panel.ItemIcon = {
  itemType: 1,
  bgColor: '#000',
}

const itemIconInfo = ref<Panel.ItemIcon>(props.itemIcon ? { ...props.itemIcon } : { ...initData })

function handleIconTypeRadioChange(type: number) {
  checkedValueRef.value = type
  itemIconInfo.value.itemType = type
  emit('update:itemIcon', itemIconInfo.value)
}

function handleChange() {
  // console.log('值', itemIconInfo.value)
  emit('update:itemIcon', itemIconInfo.value || null)
}

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
    ms.error(`上传错误：${res.msg}`)
  }

  return file
}
</script>

<template>
  <div>
    <div class="mb-[10px]">
      <NRadio
        :checked="checkedValueRef === 1 "
        :value="1"
        name="iconType"
        @change="handleIconTypeRadioChange(1)"
      >
        文字
      </NRadio>

      <NRadio
        :checked="checkedValueRef === 2"
        :value="2"
        name="iconType"
        @change="handleIconTypeRadioChange(2)"
      >
        图片/SVG
      </NRadio>

      <NRadio
        :checked="checkedValueRef === 3"
        :value="3"
        name="iconType"
        @change="handleIconTypeRadioChange(3)"
      >
        图标
      </NRadio>
    </div>

    <div class="flex h-[100px]">
      <div>
        <div class="border rounded-2xl bg-slate-200">
          <ItemIcon :item-icon="itemIconInfo" />
        </div>
      </div>
      <!-- 文字 -->
      <div class="ml-[20px]">
        <!-- <NImage :src="model.icon" preview-disabled /> -->
        <div v-if="checkedValueRef === 1">
          <NInput v-model:value="itemIconInfo.text" class="mb-[5px]" size="small" type="text" placeholder="请输入文字作为图标" @input="handleChange" />
          <NColorPicker
            v-model:value="itemIconInfo.bgColor"
            size="small"
            :modes="['hex']"
            :swatches="defautSwatchesBackground"
            @complete="handleChange"
            @update-value="handleChange"
          />
        </div>

        <div v-if="checkedValueRef === 3">
          <div>
            <NInput v-model:value="itemIconInfo.text" class="mb-[5px]" size="small" type="text" placeholder="请输入图标名字" @input="handleChange" />
            <a target="_blank" href="https://icon-sets.iconify.design/" class="text-[blue]">图标库</a>
          </div>
          <NColorPicker
            v-model:value="itemIconInfo.bgColor"
            size="small"
            :modes="['hex']"
            :swatches="defautSwatchesBackground"
            @complete="handleChange"
            @update-value="handleChange"
          />
        </div>

        <!-- 图片 -->
        <div v-if="checkedValueRef === 2">
          <NUpload
            action="/api/file/uploadImg"
            :show-file-list="false"
            name="imgfile"
            :headers="{
              token: authStore.token as string,
            }"
            @finish="handleUploadFinish"
          >
            <NButton size="small">
              点击上传
            </NButton>
          </NUpload>
        </div>
      </div>
    </div>
  </div>
</template>
