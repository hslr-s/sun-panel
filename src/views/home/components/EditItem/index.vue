<script setup lang="ts">
import { computed, defineEmits, defineProps, ref, watch } from 'vue'
import type { FormInst, FormRules } from 'naive-ui'
import {
  NButton, NCheckbox, NColorPicker, NForm, NFormItem, NGrid,
  NGridItem, NInput, NModal, NSelect,
  useMessage,
} from 'naive-ui'
import { ItemCards } from '@/views/home/components'
import { edit, getSiteFavicon } from '@/api/panel/itemIcon'
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
const getIconLoading = ref([false, false])
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

const rules: FormRules = {
  title: {
    required: true,
    trigger: 'blur',
    message: t('form.required'),
  },
  url: [
    {
      required: true,
      trigger: 'blur',
      type: 'string',
      message: t('form.required'),
    },
  ],
  itemIconGroupId: {
    required: true,
    trigger: ['blur', 'change'],
    message: t('form.required'),
    type: 'number',
  },
}

// 卡片类型
const optionsCardType = [
  {
    default: true,
    label: t('iconItem.ordinaryIcon'),
    value: 1,
  },
  {
    label: t('iconItem.dockerStatusCard'),
    value: 2,
  },
  {
    label: t('iconItem.iframeCard'),
    value: 3,
  },
]

// 打开方式
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

function verifyUrlHttp(url: string): boolean {
  if (url === '')
    return true

  return url.startsWith('https://') || url.startsWith('http://')
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

    if (!model.value.backgroundColor || model.value.backgroundColor === '')
      model.value.backgroundColor = defaultItemCardBackgroundColor
    if (!model.value.cardType || model.value.cardType === 0)
      model.value.cardType = 1
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
</script>

<template>
  <NModal v-model:show="show" preset="card" size="small" style="width: 600px;border-radius: 1rem;" :title="itemInfo ? t('iconItem.edit') : t('iconItem.add')">
    <div class="mb-2">
      <span class="flex justify-center mb-1">
        <NCheckbox v-model:checked="showPreview" size="small">
          {{ t('iconItem.preview') }}
        </NCheckbox>
      </span>
      <div v-if="showPreview" class="flex justify-center transparent-grid preview-box p-2 rounded-xl border">
        <div class="mr-5">
          <ItemCards
            :sort-status="false"
            class="cursor-pointer"
            :card-info="model"
            icon-text-color="#000"
            :icon-text-info-hide-description="false"
            :icon-text-icon-hide-title="false"
            :card-style="PanelPanelConfigStyleEnum.small"
            :card-type="PanelCardType.Standard"
          />
        </div>
        <div class="w-[210px]">
          <ItemCards
            :sort-status="false"
            class="cursor-pointer"
            :card-info="model"
            icon-text-color="#fff"
            :icon-text-info-hide-description="false"
            :icon-text-icon-hide-title="false"
            :card-style="PanelPanelConfigStyleEnum.info"
            :card-type="PanelCardType.Standard"
          />
        </div>
      </div>
    </div>

    <div class="h-[500px] overflow-auto p-[5px]">
      <NForm ref="formRef" :model="model" :rules="rules" size="small">
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

        <NGrid cols="2" :x-gap="10" item-responsive>
          <NGridItem span="2 500:1">
            <NFormItem :label="t('iconItem.cardBackground')">
              <div class="flex w-full mr-1">
                <NColorPicker
                  v-model:value="model.backgroundColor"
                  size="small"
                  :modes="['hex']"
                  :swatches="[
                    '#000000',
                    '#ffffff',
                    '#18A058',
                    '#2080F0',
                    '#F0A020',
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
          <NGridItem span="2 500:1">
            <NFormItem :label="$t('iconItem.cardType')">
              <NSelect v-model:value="model.cardType" :options="optionsCardType" />
            </NFormItem>
          </NGridItem>
        </NGrid>

        <!-- 个性化卡片的设置了 -->

        <!-- <NFormItem v-show="false" path="icon" :label="$t('common.icon')">
          <IconEditor v-model:item-icon="model.icon" />
        </NFormItem>
        <NFormItem path="url" :label="$t('iconItem.url')">
          <div class="w-full">
            <NInputGroup>
              <NInput v-model:value="model.url" type="text" :maxlength="1000" placeholder="http(s)://" />
              <NButton :disabled="!model.url" :loading="getIconLoading[0]" @click="getIconByUrl(model.url, 0)">
                {{ $t('iconItem.getIcon') }}
              </NButton>
            </NInputGroup>
            <span v-if="!verifyUrlHttp(model.url)" class="text-xs text-orange-400">
              {{ t("iconItem.urlNoHttpStartWarn") }}
            </span>
          </div>
        </NFormItem>
        <NFormItem path="lanUrl" :label="$t('iconItem.lanUrl')">
          <div class="w-full">
            <NInputGroup>
              <NInput v-model:value="model.lanUrl" type="text" :maxlength="1000" :placeholder="$t('iconItem.lanUrlInputPlaceholder')" />
              <NButton :disabled="!model.lanUrl" :loading="getIconLoading[1]" @click="getIconByUrl(model.lanUrl || '', 1)">
                {{ $t('iconItem.getIcon') }}
              </NButton>
            </NInputGroup>
            <span v-if="!verifyUrlHttp(model.lanUrl || '')" class="text-xs text-orange-400">
              {{ t("iconItem.urlNoHttpStartWarn") }}
            </span>
          </div>
        </NFormItem>
        <NFormItem path="description" :label="$t('common.description')">
          <NInput v-model:value="model.description" type="text" show-count :maxlength="100" />
        </NFormItem>
        <NFormItem path="openMethod" :label="$t('iconItem.openMethod')">
          <NSelect v-model:value="model.openMethod" :options="options" />
        </NFormItem> -->
        <hr>
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
</style>
