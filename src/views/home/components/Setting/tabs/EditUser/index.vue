<script setup lang="ts">
import { computed, defineEmits, defineProps, ref, watch } from 'vue'
import type { FormInst, FormRules } from 'naive-ui'
import { NButton, NForm, NFormItem, NInput, NSelect, useMessage } from 'naive-ui'
import { edit as userManageEdit } from '@/api/panel/users'
import { RoundCardModal } from '@/components/common'
import { t } from '@/locales'

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

const roleOtions = ref([
  {
    label: t('common.role.regularUser'),
    value: 2,
  },
  {
    label: t('common.role.admin'),
    value: 1,
  },
])

const rules: FormRules = {
  username: [
    {
      required: true,
      trigger: 'blur',
      message: t('adminSettingUsers.formRules.usernameRequired'),
      min: 5,
    },
  ],
  role: {
    required: true,
    trigger: 'blur',
    type: 'number',
    message: t('adminSettingUsers.formRules.roleRequired'),
  },
  // status: {
  //   required: true,
  //   trigger: 'blur',
  //   type: 'number',
  //   message: '请选择账号状态',
  // },
  password: {
    trigger: 'blur',
    min: 6,
    max: 20,
    message: t('adminSettingUsers.formRules.passwordLimit'),
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
    message.warning(t('common.failed'))
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
  <RoundCardModal v-model:show="show" size="small" preset="card" style="width: 400px" :title="`${userInfo?.id ? $t('common.edit') : $t('common.add')}`">
    <NForm ref="formRef" :model="model" :rules="rules">
      <NFormItem path="username" :label="$t('common.username')">
        <NInput v-model:value="model.username" type="text" :placeholder="$t('common.inputPlaceholder')" />
      </NFormItem>

      <NFormItem path="name" :label="$t('common.nikeName')">
        <NInput v-model:value="model.name" type="text" :placeholder="$t('common.inputPlaceholder')" />
      </NFormItem>

      <NFormItem path="role" :label="$t('adminSettingUsers.role')">
        <NSelect
          v-model:value="model.role"
          :options="roleOtions"
        />
      </NFormItem>

      <NFormItem path="password" :label="$t('common.password')">
        <NInput v-model:value="model.password" :maxlength="20" type="password" :placeholder="`${userInfo?.id ? $t('adminSettingUsers.EditpasswordPlaceholder') : $t('adminSettingUsers.passwordPlaceholder')}`" />
      </NFormItem>
    </NForm>

    <template #footer>
      <div class="float-right">
        <NButton type="success" size="small" @click="handleValidateButtonClick">
          {{ $t('common.save') }}
        </NButton>
      </div>
    </template>
  </RoundCardModal>
</template>
