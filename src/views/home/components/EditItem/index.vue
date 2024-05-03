<script setup lang="ts">
import { computed, defineEmits, defineProps, ref, watch } from 'vue'
import type { FormInst, FormRules } from 'naive-ui'
import {
  NButton, NCheckbox, NColorPicker, NDivider, NForm, NFormItem,
  NGrid, NGridItem, NInput, NModal, NSelect,
  useMessage,
} from 'naive-ui'
import NormalCard from './itemCards/Normal.vue'
import SystemAppCard from './itemCards/BuiltInApp.vue'
import { ItemCards } from '@/views/home/components'
import { edit } from '@/api/panel/itemIcon'
import { getList as getGroupList } from '@/api/panel/itemIconGroup'
import { t } from '@/locales'
import { apiRespErrMsg } from '@/utils/cmn/apiMessage'
import { PanelCardType, PanelPanelConfigStyleEnum } from '@/enums'
import { defaultItemCardBackgroundColor } from '@/utils/defaultData'

interface Props {
  visible: boolean
  itemInfo: Panel.Info | null
  itemGroupId?: number
}

const props = defineProps<Props>()
const emit = defineEmits<Emit>()
const ms = useMessage()
const submitLoading = ref(false)
const expandAdvancedFeatures = ref(false)
const showPreview = ref(true)
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
  cardType: 1,
}

interface Emit {
  (e: 'update:visible', visible: boolean): void
  (e: 'done', item: Panel.Info): void// 创建完成
}

const model = ref<Panel.Info>(props.itemInfo ? { ...props.itemInfo } : { ...restoreDefault })
const formRef = ref<FormInst | null>(null)

const defaultRules: FormRules = {
  title: {
    required: true,
    trigger: 'blur',
    message: t('form.required'),
  },
  itemIconGroupId: {
    required: true,
    trigger: ['blur', 'change'],
    message: t('form.required'),
    type: 'number',
  },
}
const rules = ref<FormRules>(defaultRules)

// 卡片类型
const optionsCardType = [
  {
    default: true,
    label: t('iconItem.ordinaryIcon'),
    value: PanelCardType.Normal,
  },
  // {
  //   label: t('iconItem.dockerStatusCard'),
  //   value: 2,
  // },
  {
    label: t('iconItem.systemAppCard'),
    value: PanelCardType.BuiltInApp,
  },
  // {
  //   label: t('iconItem.iframeCard'),
  //   value: 4,
  // },
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

  await edit<Panel.ItemInfo>(model.value).then(({ data }) => {
    show.value = false
    model.value = { ...restoreDefault }
    emit('done', data)
  }).catch((error) => {
    if (error.code === -2)
      ms.error(`${t('common.saveFail')}: ${t('iconItem.onlyNameExisted')}`)
    else if (error.code === -1)
      ms.error(t('common.saveFail'))
    else
      apiRespErrMsg(error)
  })

  submitLoading.value = false
}

const handleValidateButtonClick = (e: MouseEvent) => {
  e.preventDefault()
  formRef.value?.validate((errors) => {
    if (!errors)
      editApi()
  })
}

watch(() => props.visible, (newValue) => {
  if (newValue === true) {
    model.value = props.itemInfo ? { ...props.itemInfo } : { ...restoreDefault }
    if (props.itemGroupId)
      model.value.itemIconGroupId = props.itemGroupId

    if (!model.value.backgroundColor || model.value.backgroundColor === '')
      model.value.backgroundColor = defaultItemCardBackgroundColor
    if (!model.value.cardType || model.value.cardType === 0)
      model.value.cardType = 1

    if (!model.value.expandParam)
      model.value.expandParam = {}
  }

  getGroupListOptions()
})

function getGroupListOptions() {
  getGroupList<Common.ListResponse<Panel.ItemGroupInfo[]>>().then(({ data, code, msg }) => {
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
  }).catch(() => {
    ms.error(t('iconItem.getGroupFail'))
  })
}

function handleUpdateCardType() {
  rules.value = { ...defaultRules }
  // model.value.expandParam = undefined
}
</script>

<template>
  <NModal v-model:show="show" preset="card" size="small" style="width: 600px;border-radius: 1rem;" :title="itemInfo ? t('iconItem.edit') : t('iconItem.add')">
    <!-- 卡片预览 -->
    <div class="mb-2">
      <span class="flex justify-center mb-1">
        <NCheckbox v-model:checked="showPreview" size="small">
          {{ t('iconItem.preview') }}
        </NCheckbox>
      </span>
      <div v-if="showPreview" class="flex justify-center transparent-grid preview-box p-2 rounded-xl border">
        <div class="w-[210px] mr-4">
          <ItemCards
            :sort-status="false"
            class="cursor-pointer"
            :card-info="model"
            :icon-text-info-hide-description="false"
            :icon-text-icon-hide-title="false"
            :card-style="PanelPanelConfigStyleEnum.info"
            :card-type="PanelCardType.Normal"
          />
        </div>

        <div>
          <ItemCards
            :sort-status="false"
            class="cursor-pointer"
            :card-info="model"
            icon-text-color="#000"
            :icon-text-info-hide-description="false"
            :icon-text-icon-hide-title="false"
            :card-style="PanelPanelConfigStyleEnum.small"
            :card-type="PanelCardType.Normal"
          />
        </div>
      </div>
    </div>

    <NDivider class="!my-2" />
    <!-- {{ JSON.stringify(rules) }} -->

    <!-- 基础设置 -->
    <div class="h-[500px] overflow-auto p-[5px]">
      <NForm ref="formRef" :model="model" :rules="rules" size="small">
        <NGrid cols="3" :x-gap="10" item-responsive>
          <NGridItem span="3 500:1">
            <NFormItem path="itemIconGroupId" :label="t('iconItem.iconGroup')">
              <NSelect v-model:value="model.itemIconGroupId" :options="itemIconGroupOptions" />
            </NFormItem>
          </NGridItem>
          <NGridItem span="3 500:1">
            <NFormItem :label="$t('iconItem.cardType')">
              <NSelect v-model:value="model.cardType" :options="optionsCardType" @update:value="handleUpdateCardType" />
            </NFormItem>
          </NGridItem>
          <NGridItem span="3 500:1">
            <NFormItem :label="t('iconItem.cardBackground')">
              <div class="flex w-full mr-1">
                <NColorPicker
                  v-model:value="model.backgroundColor"
                  size="small"
                  :modes="['hex']"
                  :swatches="[
                    '#00000000',
                    '#000000',
                    '#ffffff',
                    '#18A058',
                    '#2080F0',
                    '#F0A020',
                    '#D03050',
                    '#C418D1',
                  ]"
                />
              </div>
              <div v-if="model.backgroundColor !== defaultItemCardBackgroundColor" class="ml-auto">
                <NButton @click="model.backgroundColor = defaultItemCardBackgroundColor">
                  {{ t('iconItem.resetBackground') }}
                </NButton>
              </div>
            </NFormItem>
          </NGridItem>
        </NGrid>

        <NDivider class="!my-2" />
        <!-- 个性化卡片的设置了 -->
        <template v-if="model.cardType === 1">
          <NormalCard v-model:itemInfo="model" v-model:rules="rules" />
        </template>
        <template v-if="model.cardType === 3">
          <SystemAppCard v-model:itemInfo="model" v-model:rules="rules" />
        </template>

        <NDivider class="!my-2" />

        <!-- 高级功能 -->
        <NCheckbox v-model:checked="expandAdvancedFeatures" class="my-2">
          {{ $t('iconItem.advancedFeatures') }}
        </NCheckbox>

        <div v-if="expandAdvancedFeatures">
          <NFormItem path="onlyName" :label="$t('iconItem.onlyName')">
            <NInput v-model:value="model.onlyName" type="text" show-count :maxlength="20" />
          </NFormItem>
        </div>
      </NForm>
    </div>

    <template #footer>
      <NButton type="success" :loading="submitLoading" style="float: right;" @click="handleValidateButtonClick">
        {{ $t('common.save') }}
      </NButton>
    </template>
  </NModal>
</template>

<style scoped>
.preview-box{
  pointer-events: none;
}

.dark .preview-box{
  filter: brightness(80%)
}
</style>
