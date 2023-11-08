<script setup lang="ts">
import { computed, defineEmits, defineProps, ref, watch } from 'vue'
import type { FormInst, FormRules } from 'naive-ui'
import { NButton, NForm, NFormItem, NInput, useMessage } from 'naive-ui'
import { edit as userManageEdit } from '@/api/panel/users'
import { RoundCardModal } from '@/components/common'

interface Props {
  visible: boolean
  userId?: number
  userInfo?: User.Info
}

const props = defineProps<Props>()
const emit = defineEmits<Emit>()
const message = useMessage()

interface Emit {
  (e: 'update:visible', visible: boolean): void
  (e: 'done', id: number): void// 创建完成
}

const formInitValue = {
  name: '',
  username: '',
  role: 2,
  status: 3,
}

const model = ref<User.Info>(formInitValue)
const formRef = ref<FormInst | null>(null)

const rules: FormRules = {
  username: [
    {
      required: true,
      trigger: 'blur',
      message: '请输入账号且大于5个字符',
      min: 5,
    },
    {
      trigger: 'blur',
      message: '请输入邮箱作为账号',
      type: 'email',
    },
  ],
  role: {
    required: true,
    trigger: 'blur',
    type: 'number',
    message: '请选择角色',
  },
  status: {
    required: true,
    trigger: 'blur',
    type: 'number',
    message: '请选择账号状态',
  },
}

// 更新值父组件传来的值
const show = computed({
  get: () => props.visible,
  set: (visible: boolean) => {
    emit('update:visible', visible)
  },
})

watch(show, (newValue, oldValue) => {
  if (props.userInfo?.id)
    model.value = props.userInfo || {}

  else
    model.value = formInitValue
})

const add = async () => {
  const res = await userManageEdit<User.Info>(model.value)
  if (res.code === 0)
    emit('done', res.data.id as number)

  else if (res.code !== -1)
    message.warning('操作失败')
}

const handleValidateButtonClick = (e: MouseEvent) => {
  e.preventDefault()
  formRef.value?.validate((errors) => {
    if (!errors)
      add()
    else
      console.log(errors)
  })
}
</script>

<template>
  <RoundCardModal v-model:show="show" size="small" preset="card" style="width: 400px" :title="`${userInfo?.id ? '编辑' : '添加'}用户`">
    <NForm ref="formRef" :model="model" :rules="rules">
      <NFormItem path="username" label="账号">
        <NInput v-model:value="model.username" type="text" placeholder="邮箱地址作为账号" />
      </NFormItem>

      <NFormItem path="name" label="昵称">
        <NInput v-model:value="model.name" type="text" placeholder="请输入昵称" />
      </NFormItem>

      <NFormItem path="password" label="密码">
        <NInput v-model:value="model.password" type="text" :placeholder="`${userInfo?.id ? '请输入新密码，留空密码不变' : '请输入密码'}`" />
      </NFormItem>
    </NForm>

    <template #footer>
      <div class="float-right">
        <NButton type="success" size="small" @click="handleValidateButtonClick">
          保存
        </NButton>
      </div>
    </template>
  </RoundCardModal>
</template>
