<script setup lang="ts">
import { computed, defineEmits, defineProps, ref, watch } from 'vue'
import type { FormInst, FormRules } from 'naive-ui'
import { NButton, NForm, NFormItem, NGrid, NGridItem, NInput, NModal, NSelect, useMessage } from 'naive-ui'
import IconEditor from './IconEditor.vue'
import { edit } from '@/api/panel/itemIcon'
import { getList as getGroupList } from '@/api/panel/itemIconGroup'

interface Props {
  visible: boolean
  itemInfo: Panel.Info | null
}

const props = defineProps<Props>()
const emit = defineEmits<Emit>()
const ms = useMessage()
const itemIconGroupOptions = ref<{
  label: string
  value: number
}[]>([])

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
// const model = computed(()=>{
//   return props.itemInfo !== null ? { ...props.itemInfo } : { ...restoreDefault }
// })
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
  // itemIconGroupId: {
  //   required: true,
  //   trigger: ['blur', 'change'],
  //   message: '必填项',
  // },
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

// const urlProtocolOptions = [
//   {
//     default: true,
//     label: 'http://',
//     value: 'http://',
//   },
//   {
//     label: 'https://',
//     value: 'https://',
//   },
//   {
//     label: '不使用',
//     value: '',
//   },
// ]

// 更新值父组件传来的值
const show = computed({
  get: () => props.visible,
  set: (visible: boolean) => {
    emit('update:visible', visible)
  },
})

async function editApi() {
  const { code, data, msg } = await edit<Panel.ItemInfo>(model.value)
  if (code === 0) {
    show.value = false
    model.value = { ...restoreDefault }

    emit('done', data)
  }
  else {
    ms.error(`保存失败：${msg}`)
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
  model.value = newValue || { ...restoreDefault }
  getGroupListOptions()
})

watch(() => props.visible, (newValue) => {
  if (newValue)
    getGroupListOptions()
})

function getGroupListOptions() {
  getGroupList<Common.ListResponse<Panel.ItemIconGroup[]>>().then(({ data, code, msg }) => {
    if (code === 0) {
      itemIconGroupOptions.value = []

      for (let i = 0; i < data.list.length; i++) {
        const element = data.list[i]
        if (i === 0 && !model.value.itemIconGroupId) {
          restoreDefault.itemIconGroupId = element.id
          model.value.itemIconGroupId = element.id
        }

        itemIconGroupOptions.value.push({
          value: element.id as number,
          label: element.title as string,
        })
      }
    }
    else {
      ms.error(`分组信息获取失败:${msg}`)
    }
  })
}
</script>

<template>
  <NModal v-model:show="show" preset="card" style="width: 600px;border-radius: 1rem;" :title="itemInfo ? '修改项目' : '添加项目'">
    <div class="h-[600px] overflow-auto p-[5px]">
      <NForm ref="formRef" :model="model" :rules="rules">
        <NGrid cols="2" :x-gap="10" item-responsive>
          <NGridItem span="2 500:1">
            <NFormItem path="itemIconGroupId" label="分组">
              <NSelect v-model:value="model.itemIconGroupId" :options="itemIconGroupOptions" />
            </NFormItem>
          </NGridItem>
          <NGridItem span="2 500:1">
            <NFormItem path="title" label="标题">
              <NInput v-model:value="model.title" type="text" show-count :maxlength="20" placeholder="请输入标题" />
            </NFormItem>
          </NGridItem>
        </NGrid>

        <NFormItem path="icon" label="图标">
          <IconEditor v-model:item-icon="model.icon" />
        </NFormItem>
        <NFormItem path="url" label="跳转地址">
          <!-- <NSelect :style="{ width: '100px' }" :options="urlProtocolOptions" /> -->
          <NInput v-model:value="model.url" type="text" :maxlength="1000" placeholder="http(s)://" />
        </NFormItem>
        <NFormItem path="lanUrl" label="局域网跳转地址">
          <NInput v-model:value="model.lanUrl" type="text" :maxlength="1000" placeholder="http(s)://（可以留空，切换到局域网模式，点击会使用该地址）" />
        </NFormItem>
        <NFormItem path="description" label="描述信息">
          <NInput v-model:value="model.description" type="text" show-count :maxlength="100" placeholder="请填写描述信息" />
        </NFormItem>
        <NFormItem path="openMethod" label="打开方式">
          <NSelect v-model:value="model.openMethod" :options="options" />
        </NFormItem>
      </NForm>
    </div>

    <template #footer>
      <NButton type="success" style="float: right;" @click="handleValidateButtonClick">
        确定
      </NButton>
    </template>
  </NModal>
</template>
