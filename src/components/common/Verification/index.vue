<script setup lang='ts'>
import { computed, ref, watch } from 'vue'
import { NButton, NInput, NModal, useMessage } from 'naive-ui'
import { Captcha } from '@/components/common'

const props = defineProps<{
  visible: boolean
  verificationId: string
  loading?: boolean
  title?: string
}>()
const emit = defineEmits<{
  (e: 'update:visible', visible: boolean): void
  (e: 'update:loading', loading: boolean): void
  (e: 'done', id: number): void// 创建完成
  (e: 'onSubmit', verificationId: string, vCode: string, obj: { refresh: () => void }): void // 提交
}>()
const ms = useMessage()
const vCode = ref<string>('')
const captchaRef = ref()

// 更新值父组件传来的值
const show = computed({
  get: () => props.visible,
  set: (visible: boolean) => {
    emit('update:visible', visible)
  },
})

watch(show, (newValue, oldValue) => {
  if (newValue === true) {
    captchaRef.value?.refresh()
    vCode.value = ''
    emit('update:loading', false)
  }
})

// 提交
function handleSubmit() {
  if (vCode.value === '') {
    ms.warning('验证码必须填写')
    return
  }

  if (!props.loading) {
    emit('onSubmit', props.verificationId, vCode.value, {
      refresh() {
        captchaRef.value?.refresh()
        vCode.value = ''
      },
    })
  }
}
</script>

<template>
  <div>
    <NModal v-model:show="show" preset="card" style="width: 400px" :title="title ?? '请输入验证码继续'" :mask-closable="false">
      <Captcha ref="captchaRef" class="rounded border" :src="`/api/captcha/getImageByCaptchaId/${verificationId}/200/60?0`" />
      <div class="flex">
        <div class="flex w-[80%]">
          <NInput v-model:value="vCode" placeholder="输入图中字母或数字后继续" @keydown.enter="handleSubmit" />
        </div>
        <div class="flex ml-[auto]">
          <NButton type="info" :loading="loading" :disabled="loading" @click="handleSubmit">
            继续
          </NButton>
        </div>
      </div>
    </NModal>
  </div>
</template>
