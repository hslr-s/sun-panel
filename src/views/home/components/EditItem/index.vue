<script setup lang="ts">
import { computed, defineEmits, defineProps, ref, watch } from 'vue'
import type { FormInst, FormRules } from 'naive-ui'
import { NButton, NForm, NFormItem, NGrid, NGridItem, NInput, NInputGroup, NModal, NSelect, useMessage } from 'naive-ui'
import IconEditor from './IconEditor.vue'
import { edit, getSiteFavicon } from '@/api/panel/itemIcon'
import { getList as getGroupList } from '@/api/panel/itemIconGroup'
import { t } from '@/locales'

interface Props {
  visible: boolean
  itemInfo: Panel.Info | null
  itemGroupId?: number
}

const props = defineProps<Props>()
const emit = defineEmits<Emit>()
const ms = useMessage()
const submitLoading = ref(false)
const getIconLoading = ref([false, false])
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
  openMethod: 2,
}

interface Emit {
  (e: 'update:visible', visible: boolean): void
  (e: 'done', item: Panel.Info): void// 创建完成
}

const model = ref<Panel.Info>(props.itemInfo ? { ...props.itemInfo } : { ...restoreDefault })
const formRef = ref<FormInst | null>(null)

const rules: FormRules = {
  title: {
    required: true,
    trigger: 'blur',
    message: t('form.required'),
  },
  url: {
    required: true,
    trigger: 'blur',
    type: 'string',
    message: t('form.required'),
  },
  // itemIconGroupId: {
  //   required: true,
  //   trigger: ['blur', 'change'],
  //   message: t('form.required'),
  // },
}

const options = [
  {
    default: true,
    label: t('iconItem.currentPageOpen'),
    value: 1,
  },
  {
    label: t('iconItem.newWindowOpen'),
    value: 2,
  },
  {
    label: t('iconItem.currentPageLayerOpen'),
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
  submitLoading.value = true
  try {
    const { code, data, msg } = await edit<Panel.ItemInfo>(model.value)
    if (code === 0) {
      show.value = false
      model.value = { ...restoreDefault }

      emit('done', data)
    }
    else {
      ms.error(`${t('common.saveFail')}:${msg}`)
    }
  }
  catch (error) {
    ms.error(t('common.saveFail'))
  }
  submitLoading.value = false
}

const handleValidateButtonClick = (e: MouseEvent) => {
  e.preventDefault()
  formRef.value?.validate((errors) => {
    if (!errors)
      editApi()
  })
}

async function getIconByUrl(url: string, loadingIndex: number) {
  getIconLoading.value[loadingIndex] = true
  try {
    const { code, data } = await getSiteFavicon<{ iconUrl: string }>(url)
    if (code === 0) {
      model.value.icon = {
        itemType: 2,
        src: data.iconUrl,
      }
    }
    else {
      ms.error(t('iconItem.geticonFail'))
    }
  }
  catch (error) {
    ms.error(t('iconItem.geticonFail'))
  }
  getIconLoading.value[loadingIndex] = false
}

watch(() => props.visible, (newValue) => {
  if (newValue === true) {
    model.value = props.itemInfo ? { ...props.itemInfo } : { ...restoreDefault }
    if (props.itemGroupId)
      model.value.itemIconGroupId = props.itemGroupId
  }

  getGroupListOptions()
})

function getGroupListOptions() {
  getGroupList<Common.ListResponse<Panel.ItemIconGroup[]>>().then(({ data, code, msg }) => {
    if (code === 0) {
      itemIconGroupOptions.value = []

      for (let i = 0; i < data.list.length; i++) {
        const element = data.list[i]
        if (i === 0 && !model.value.itemIconGroupId) {
          model.value.itemIconGroupId = element.id
          restoreDefault.itemIconGroupId = element.id
        }

        itemIconGroupOptions.value.push({
          value: element.id as number,
          label: element.title as string,
        })
      }
    }
    else {
      ms.error(`${t('iconItem.getGroupFail')}:${msg}`)
    }
  })
}
</script>

<template>
  <NModal v-model:show="show" preset="card" size="small" style="width: 600px;border-radius: 1rem;" :title="itemInfo ? t('iconItem.edit') : t('iconItem.add')">
    <div class="h-[600px] overflow-auto p-[5px]">
      <NForm ref="formRef" :model="model" :rules="rules">
        <NGrid cols="2" :x-gap="10" item-responsive>
          <NGridItem span="2 500:1">
            <NFormItem path="itemIconGroupId" :label="t('iconItem.iconGroup')">
              <NSelect v-model:value="model.itemIconGroupId" :options="itemIconGroupOptions" />
            </NFormItem>
          </NGridItem>
          <NGridItem span="2 500:1">
            <NFormItem path="title" :label="$t('common.title')">
              <NInput v-model:value="model.title" type="text" show-count :maxlength="20" />
            </NFormItem>
          </NGridItem>
        </NGrid>

        <NFormItem path="icon" :label="$t('common.icon')">
          <IconEditor v-model:item-icon="model.icon" />
        </NFormItem>
        <NFormItem path="url" :label="$t('iconItem.url')">
          <!-- <NSelect :style="{ width: '100px' }" :options="urlProtocolOptions" /> -->
          <NInputGroup>
            <NInput v-model:value="model.url" type="text" :maxlength="1000" placeholder="http(s)://" />
            <NButton :disabled="!model.url" :loading="getIconLoading[0]" @click="getIconByUrl(model.url, 0)">
              {{ $t('iconItem.getIcon') }}
            </NButton>
          </NInputGroup>
        </NFormItem>
        <NFormItem path="lanUrl" :label="$t('iconItem.lanUrl')">
          <NInputGroup>
            <NInput v-model:value="model.lanUrl" type="text" :maxlength="1000" :placeholder="$t('iconItem.lanUrlInputPlaceholder')" />
            <NButton :disabled="!model.lanUrl" :loading="getIconLoading[1]" @click="getIconByUrl(model.lanUrl || '', 1)">
              {{ $t('iconItem.getIcon') }}
            </NButton>
          </NInputGroup>
        </NFormItem>
        <NFormItem path="description" :label="$t('common.description')">
          <NInput v-model:value="model.description" type="text" show-count :maxlength="100" />
        </NFormItem>
        <NFormItem path="openMethod" :label="$t('iconItem.openMethod')">
          <NSelect v-model:value="model.openMethod" :options="options" />
        </NFormItem>
      </NForm>
    </div>

    <template #footer>
      <NButton type="success" :loading="submitLoading" style="float: right;" @click="handleValidateButtonClick">
        {{ $t('common.save') }}
      </NButton>
    </template>
  </NModal>
</template>
