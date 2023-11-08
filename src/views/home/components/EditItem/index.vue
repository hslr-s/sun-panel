<script setup lang="ts">
import { computed, defineEmits, defineProps, onMounted, ref, watch } from 'vue'
import type { FormInst, FormRules } from 'naive-ui'
import { NButton, NForm, NFormItem, NInput, NModal, NSelect, useMessage } from 'naive-ui'
import IconEditor from './IconEditor.vue'
import { edit } from '@/api/panel/itemIcon'

interface Props {
  visible: boolean
  itemInfo: Panel.Info | null
}

const props = defineProps<Props>()
const emit = defineEmits<Emit>()
const ms = useMessage()

const restoreDefault: Panel.Info = {
  icon: null,
  title: '',
  url: '',
  lanUrl: '',
  description: '',
  openMethod: 1,
}

interface Emit {
  (e: 'update:visible', visible: boolean): void
  (e: 'done', item: Panel.Info): void// 创建完成
}

const model = ref<Panel.Info>(props.itemInfo !== null ? { ...props.itemInfo } : { ...restoreDefault })
const formRef = ref<FormInst | null>(null)

const rules: FormRules = {
  title: {
    required: true,
    trigger: 'blur',
    message: '必填项',
  },
  url: {
    required: true,
    trigger: 'blur',
    type: 'string',
    message: '必填项',
  },
}

const options = [
  {
    default: true,
    label: '当前页面打开',
    value: 1,
  },
  {
    label: '新窗口打开',
    value: 2,
  },
  {
    label: '弹窗打开',
    value: 3,
  },
]

// 更新值父组件传来的值
const show = computed({
  get: () => props.visible,
  set: (visible: boolean) => {
    emit('update:visible', visible)
  },
})

async function editApi() {
  const { code, data } = await edit<Panel.ItemInfo>(model.value)
  if (code === 0) {
    show.value = false
    model.value = restoreDefault

    emit('done', data)
  }
  else {
    ms.error('保存失败')
  }
}

const handleValidateButtonClick = (e: MouseEvent) => {
  e.preventDefault()
  formRef.value?.validate((errors) => {
    if (!errors)
      editApi()
  })
}

watch(() => props.itemInfo, (newValue) => {
  model.value = newValue || restoreDefault
})

onMounted(() => {
  // rolesLoading.value = true
})
</script>

<template>
  <NModal v-model:show="show" preset="card" style="width: 600px;border-radius: 1rem;" :title="itemInfo ? '修改项目' : '添加项目'">
    <NForm ref="formRef" :model="model" :rules="rules">
      <NFormItem path="title" label="标题">
        <NInput v-model:value="model.title" type="text" show-count :maxlength="20" placeholder="请输入标题" />
      </NFormItem>
      <NFormItem path="icon" label="图标">
        <IconEditor v-model:item-icon="model.icon" />
      </NFormItem>
      <NFormItem path="url" label="跳转地址">
        <NInput v-model:value="model.url" type="text" :maxlength="1000" placeholder="请输入跳转地址" />
      </NFormItem>
      <NFormItem path="lanUrl" label="局域网跳转地址">
        <NInput v-model:value="model.lanUrl" type="text" :maxlength="1000" placeholder="（可以留空）切换到局域网模式，点击会使用该地址" />
      </NFormItem>
      <NFormItem path="description" label="描述信息">
        <NInput v-model:value="model.description" type="text" show-count :maxlength="100" placeholder="请填写描述信息" />
      </NFormItem>
      <NFormItem path="openMethod" label="打开方式">
        <NSelect v-model:value="model.openMethod" :options="options" />
      </NFormItem>
    </NForm>

    <template #footer>
      <NButton type="success" @click="handleValidateButtonClick">
        确定
      </NButton>
    </template>
  </NModal>
</template>
