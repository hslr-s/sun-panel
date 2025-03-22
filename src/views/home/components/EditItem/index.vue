<script setup lang="ts">
import { computed, defineEmits, defineProps, ref, watch } from 'vue'
import type { FormInst, FormRules } from 'naive-ui'
import {
  NButton,
  NCheckbox, NColorPicker, NFlex, NForm, NFormItem, NGrid, NGridItem,
  NInput, NModal, NSelect, useDialog, useMessage,
} from 'naive-ui'

import NormalCard from './itemCards/Normal.vue'
import SystemAppCard from './itemCards/BuiltInApp.vue'
import DockerCard from './itemCards/Docker.vue'

import { ItemCards } from '@/views/home/components'
import { deletes, edit } from '@/api/panel/itemIcon'
import { getList as getGroupList } from '@/api/panel/itemIconGroup'
import { t } from '@/locales'
import { apiRespErrMsg } from '@/utils/cmn/apiMessage'
import { PanelCardType, PanelPanelConfigStyleEnum } from '@/enums'
import { defaultItemCardBackgroundColor } from '@/utils/defaultData'
import { STip } from '@/components/common'
import { useDomainLinks } from '@/hooks/useDomainLinks'
import { useAuthStore } from '@/store'
import { AdminAuthRole } from '@/enums/admin'

interface Props {
  visible: boolean
  itemInfo: Panel.Info | null
  itemGroupId?: number
  isOutAdd?: boolean // 从外部添加,需要优先展示分组
}

const props = defineProps<Props>()
const emit = defineEmits<Emit>()
const ms = useMessage()
const dialog = useDialog()
const submitLoading = ref(false)
const deleteLoading = ref(false)
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

const domainLinks = useDomainLinks()
const authStore = useAuthStore()

const previewTransparentCanvas = ref(false)
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
const rules = ref<FormRules>({ ...defaultRules })

// 卡片类型
function getOptionsCardType() {
  const optionsCardType = [
    {
      default: true,
      label: t('iconItem.ordinaryIcon'),
      value: PanelCardType.Normal,
    },
    {
      label: t('iconItem.builtInApp'),
      value: PanelCardType.BuiltInApp,
    },
    // {
    //   label: t('iconItem.iframeCard'),
    //   value: 4,
    // },
  ]

  if (authStore.userInfo?.role === AdminAuthRole.admin) {
    optionsCardType.push({
      label: t('iconItem.dockerStatusCard'),
      value: PanelCardType.Docker,
    })
  }

  return optionsCardType
}

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
  handleResetRules()
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

function handleDelete() {
  if (props.itemInfo?.id) {
    dialog.warning({
      title: t('common.warning'),
      content: t('common.deleteConfirmByName', { name: props.itemInfo.title }),
      positiveText: t('common.confirm'),
      negativeText: t('common.cancel'),
      onPositiveClick: () => {
        deleteLoading.value = true
        deletes([props.itemInfo?.id as number]).then(({ code, msg }) => {
          if (code === 0) {
            ms.success(t('common.deleteSuccess'))
            show.value = false
            emit('done', props.itemInfo as Panel.Info)
          }
          else {
            ms.error(`${t('common.deleteFail')}:${msg}`)
          }
          deleteLoading.value = false
        }).catch((error) => {
          apiRespErrMsg(error)
          deleteLoading.value = false
        })
      },
    })
  }
}

// 重置表单规则
function handleResetRules() {
  rules.value = { ...defaultRules }
}
</script>

<template>
  <NModal
    v-model:show="show" preset="card" size="small" style="width: 600px;border-radius: 1rem;"
    :title="itemInfo?.id ? t('iconItem.edit') : t('iconItem.add')"
  >
    <!-- 卡片预览 -->
    <div class="mb-2">
      <span class="flex mb-1">
        <NCheckbox v-model:checked="showPreview" size="small">
          {{ t('iconItem.preview') }}
        </NCheckbox>
        <NCheckbox v-if="showPreview" v-model:checked="previewTransparentCanvas" size="small">
          {{ t('iconItem.previewTransparentCanvas') }}
        </NCheckbox>
      </span>
      <div
        v-if="showPreview" class="preview-box rounded-xl border"
        :class="previewTransparentCanvas ? 'transparent-grid' : 'bg-[#f1f8ff]'"
      >
        <div class="flex justify-center p-2">
          <div class="w-[210px] mr-4 z-[-1]">
            <ItemCards
              :sort-status="false" class="cursor-pointer" :card-info="model"
              :icon-text-info-hide-description="false" :icon-text-icon-hide-title="false"
              :card-style="PanelPanelConfigStyleEnum.info" :card-type="PanelCardType.Normal"
            />
          </div>

          <div class="z-[-1]">
            <ItemCards
              :sort-status="false" class="cursor-pointer" :card-info="model" icon-text-color="#000"
              :icon-text-info-hide-description="false" :icon-text-icon-hide-title="false"
              :card-style="PanelPanelConfigStyleEnum.small" :card-type="PanelCardType.Normal"
            />
          </div>
        </div>
      </div>
    </div>

    <!-- <NDivider class="!my-2" /> -->
    <!-- {{ JSON.stringify(rules) }} -->

    <!-- 基础设置 -->
    <div class="h-[500px] overflow-auto p-[5px]">
      <NForm ref="formRef" :model="model" :rules="rules" size="small">
        <!-- 其他窗口添加图标的时候所使用的分组 正常情况下隐藏 -->
        <NFormItem v-if="isOutAdd" path="itemIconGroupId">
          <template #label>
            <span class="font-bold">
              {{ $t('iconItem.iconGroup') }}
            </span>
          </template>
          <NSelect v-model:value="model.itemIconGroupId" :options="itemIconGroupOptions" />
        </NFormItem>

        <!-- <NGrid cols="3" :x-gap="10" item-responsive>
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
        </NGrid> -->

        <!-- <NDivider class="!my-2" /> -->
        <!-- 个性化卡片的设置了 -->
        <template v-if="model.cardType === PanelCardType.Normal">
          <NormalCard v-model:itemInfo="model" v-model:rules="rules" />
        </template>
        <template v-if="model.cardType === PanelCardType.Docker">
          <DockerCard v-model:itemInfo="model" v-model:rules="rules" />
        </template>
        <template v-if="model.cardType === PanelCardType.BuiltInApp">
          <SystemAppCard v-model:itemInfo="model" v-model:rules="rules" />
        </template>

        <!-- <NDivider class="!my-2" /> -->

        <div v-if="expandAdvancedFeatures">
          <NGrid cols="3" :x-gap="10" item-responsive>
            <NGridItem span="3 500:1">
              <NFormItem>
                <template #label>
                  <span class="font-bold">
                    {{ $t('iconItem.cardBackground') }}
                  </span>
                </template>
                <div class="flex w-full mr-1">
                  <NColorPicker
                    v-model:value="model.backgroundColor" size="small" :modes="['hex']" :swatches="[
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

            <NGridItem v-if="!isOutAdd" span="3 500:1">
              <NFormItem path="itemIconGroupId">
                <template #label>
                  <span class="font-bold">
                    {{ $t('iconItem.iconGroup') }}
                  </span>
                </template>
                <NSelect v-model:value="model.itemIconGroupId" :options="itemIconGroupOptions" />
              </NFormItem>
            </NGridItem>

            <NGridItem span="3 500:1">
              <NFormItem>
                <template #label>
                  <span class="font-bold">
                    {{ $t('iconItem.cardType') }}
                  </span>
                </template>
                <NSelect v-model:value="model.cardType" :options="getOptionsCardType()" @update:value="handleResetRules" />
              </NFormItem>
            </NGridItem>
          </NGrid>

          <NFormItem path="onlyName">
            <template #label>
              <span class="font-bold">
                <NFlex align="center">
                  {{ $t('common.onlyName') }}
                  <STip>
                    {{ $t('iconItem.onlyNameTip') }}
                    <NButton
                      text tag="a" :href="domainLinks.linkWithLang('only_name_beautify_doc')" target="_blank"
                      type="primary"
                    >
                      {{ t("common.lookDoc") }}
                    </NButton>
                  </STip>
                </NFlex>
              </span>
            </template>
            <NInput v-model:value="model.onlyName" type="text" show-count :maxlength="20" />
          </NFormItem>
        </div>

        <!-- 更多不常用的功能进行隐藏 -->
        <NCheckbox v-model:checked="expandAdvancedFeatures" class="mb-2" size="small">
          {{ $t('iconItem.moreOptions') }}
          <span class="text-xs text-gray-400">
            ({{ $t('iconItem.cardBackground') }} , {{ $t('iconItem.iconGroup') }} , {{ $t('iconItem.cardType') }} , {{ $t('common.onlyName') }})
          </span>
        </NCheckbox>
      </NForm>
    </div>

    <template #footer>
      <NFlex justify="end">
        <NButton v-if="itemInfo?.id" type="error" :loading="deleteLoading" style="float: right;" @click="handleDelete">
          {{ $t('common.delete') }}
        </NButton>

        <NButton type="success" :loading="submitLoading" style="float: right;" @click="handleValidateButtonClick">
          {{ $t('common.save') }}
        </NButton>
      </NFlex>
    </template>
  </NModal>
</template>

<style scoped>
.preview-box{
  position: relative;
  z-index: 1;
}

.dark .preview-box{
  filter: brightness(80%)
}
</style>
