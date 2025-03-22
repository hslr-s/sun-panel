<script setup lang="ts">
import { onMounted, ref, watch } from 'vue'
import type { FormInst, FormRules } from 'naive-ui'
import {
  NButton, NCard, NColorPicker,
  NDivider, NFlex, NForm, NFormItem, NInput,
  NSelect, NSpin, NSwitch, useDialog,
  useMessage,
} from 'naive-ui'
import { VueDraggable } from 'vue-draggable-plus'
import { deletes, edit, getList, saveSort } from '@/api/panel/itemIconGroup'
import { RoundCardModal, STip, SvgIcon } from '@/components/common'
import { t } from '@/locales'
import { PanelPanelConfigStyleEnum } from '@/enums'
import { defaultItemGroupInfoCardStyle } from '@/utils/defaultData'
import { apiRespErrMsg } from '@/utils/cmn/apiMessage'
import { useTriggerUpdateItemData } from '@/hooks/useTriggerUpdateItemData'
import { useDomainLinks } from '@/hooks/useDomainLinks'

interface EditModalArg {
  show: boolean
  editStatus: number // 1.添加 2.编辑
  model: Panel.ItemGroupInfo
  rules: FormRules
}

const domainLinks = useDomainLinks()

const formRef = ref<FormInst | null>(null)
const ms = useMessage()
const dialog = useDialog()
const sortStatus = ref(false)
const triggerUpdateItemData = useTriggerUpdateItemData()
const loading = ref(false)

const defaultItemGroupInfo: Panel.ItemGroupInfo = {
  title: '',
  cardStyle: defaultItemGroupInfoCardStyle,
}

const iconTypeOptions = [
  {
    label: t('apps.baseSettings.detailIcon'),
    value: PanelPanelConfigStyleEnum.info,
  },
  {
    label: t('apps.baseSettings.smallIcon'),
    value: PanelPanelConfigStyleEnum.icon,
  },
]

const editModalArg = ref<EditModalArg>({
  show: false,
  editStatus: 1,
  model: defaultItemGroupInfo,
  rules: {
    title: [
      {
        required: true,
        trigger: 'blur',
        message: t('form.required'),
      },
    ],
  },
})

const groups = ref<Panel.ItemGroupInfo[]>([])

function handleAddGroup() {
  editModalArg.value.show = !editModalArg.value.show
}

function handleEditGroup(groupInfo: Panel.ItemGroupInfo) {
  editModalArg.value.show = true
  editModalArg.value.model = groupInfo || defaultItemGroupInfo
  editModalArg.value.editStatus = 2

  if (editModalArg.value.model.cardStyle.textColor === '')
    editModalArg.value.model.cardStyle.textColor = defaultItemGroupInfoCardStyle.textColor
}

function handleDragSort() {
  sortStatus.value = true
}

function handleSaveSort() {
  const saveItems: Common.SortItemRequest[] = []
  for (let i = 0; i < groups.value.length; i++) {
    const element = groups.value[i]
    saveItems.push({
      id: element.id as number,
      sort: i + 1,
    })
  }
  saveSort(saveItems).then(({ code, msg }) => {
    ms.success(t('common.saveSuccess'))
    sortStatus.value = false
    triggerUpdateItemData.triggerUpdate()
  }).catch(() => {
    ms.error(t('common.saveFail'))
  })
}

function handleDelete(groupInfo: Panel.ItemGroupInfo) {
  dialog.warning({
    title: t('common.warning'),
    content: t('apps.itemGroupManage.deleteWarnText', { name: groupInfo.title }),
    positiveText: t('common.confirm'),
    negativeText: t('common.cancel'),
    onPositiveClick: () => {
      if (groupInfo.id) {
        deletes([groupInfo.id]).then(({ code, msg }) => {
          if (code !== 0)
            ms.error(t('common.deleteFail'))
          else
            refreshList()
        }).catch((res) => {
          // ms.error(t('common.deleteFail'))
          apiRespErrMsg(res)
        })
      }
    },

  })
}

function handleSaveGroup() {
  formRef.value?.validate((errors) => {
    if (!errors) {
      edit(editModalArg.value.model).then(({ code, msg }) => {
        if (code !== 0)
          ms.error(t('common.saveFail'))

        refreshList()
        triggerUpdateItemData.triggerUpdate()
        editModalArg.value.show = false
        editModalArg.value.model = { ...defaultItemGroupInfo }
      }).catch((res) => {
        if (res.code === -2)
          ms.error(`${t('common.saveFail')}: ${t('iconItem.onlyNameExisted')}`)
        else
          apiRespErrMsg(res)
      })
    }
    else { console.log(errors) }
  })
}

async function refreshList() {
  loading.value = true
  await getList<Common.ListResponse<Panel.ItemGroupInfo[]>>().then(({ data }) => {
    groups.value = data.list
  })
  loading.value = false
}

onMounted(() => {
  refreshList()
})

watch(() => editModalArg.value.show, (v) => {
  // console.log(v)
  if (!v)
    editModalArg.value.model = { ...defaultItemGroupInfo }
})
</script>

<template>
  <div class="p-2">
    <NCard style="border-radius: 10px;margin-bottom:10px ;" size="small">
      <div>
        <NButton type="success" size="small" style="margin-right: 10px;" @click="handleAddGroup">
          {{ $t('common.add') }}
        </NButton>

        <NButton v-if="!sortStatus" size="small" @click="handleDragSort">
          {{ $t('common.sort') }}
        </NButton>

        <NButton v-else type="warning" size="small" @click="handleSaveSort">
          {{ $t('common.saveSort') }}
        </NButton>
      </div>
    </NCard>

    <NCard style="border-radius: 10px;" size="small">
      <NSpin v-if="loading" />
      <div class="w-full mt-[20px] rounded-xl my-20 pt-[5px]" :class="sortStatus ? 'border border-dashed' : ''">
        <VueDraggable
          v-model="groups"
          item-key="sort" :animation="300"
          :style="{ padding: sortStatus ? '20px' : '5px' }"
          :disabled="!sortStatus"
        >
          <div v-for="(item, index) in groups" :key="index" class="w-full">
            <NCard size="small" style="border-radius:10px;margin-bottom: 10px;" class="card">
              <div class="flex" :class="sortStatus ? 'cursor-move' : ''">
                <div class="flex items-center cursor-pointer" @click="handleEditGroup(item)">
                  <span class="mr-[10px]">
                    <SvgIcon class="text-[20px]" icon="material-symbols:ad-group-outline-rounded" />
                  <!-- <SvgIcon class="text-[20px]" :icon="item.icon" /> -->
                  </span>
                  <span>
                    {{ item.title }}
                  </span>
                </div>
                <div class="ml-auto">
                  <span>
                    <NButton strong secondary type="success" size="small" @click="handleEditGroup(item)">
                      <template #icon>
                        <SvgIcon icon="basil:edit-solid" />
                      </template>
                    </NButton>
                  </span>
                  <span class="ml-[10px]">
                    <NButton strong secondary type="error" size="small" class="ml-[10px]" @click="handleDelete(item)">
                      <template #icon>
                        <SvgIcon icon="material-symbols:delete" />
                      </template>
                    </NButton>
                  </span>
                </div>
              </div>
            </NCard>
          </div>
        </VueDraggable>
      </div>
    </NCard>

    <RoundCardModal v-model:show="editModalArg.show" size="small" type="small" :title="editModalArg.editStatus === 1 ? t('common.add') : t('common.edit')" style="width: 400px;">
      <NForm ref="formRef" :model="editModalArg.model" size="small" :rules="editModalArg.rules">
        <NFormItem path="title" :label="$t('apps.itemGroupManage.groupName')">
          <NInput v-model:value="editModalArg.model.title" type="text" :maxlength="20" show-count />
        </NFormItem>

        <NFormItem :label="t('apps.itemGroupManage.publicVisitHide')">
          <NSwitch v-model:value="editModalArg.model.publicVisitHide" />
        </NFormItem>

        <NDivider class="!my-2" />

        <NFormItem :label="t('apps.itemGroupManage.cardStyle')">
          <NSelect v-model:value="editModalArg.model.cardStyle.style" :options="iconTypeOptions" />
        </NFormItem>

        <NFormItem v-if="editModalArg.model.cardStyle.style === PanelPanelConfigStyleEnum.info" :label="t('apps.baseSettings.hideDescription')">
          <NSwitch v-model:value="editModalArg.model.cardStyle.textInfoHideDescription" />
        </NFormItem>

        <NFormItem v-if="editModalArg.model.cardStyle.style === PanelPanelConfigStyleEnum.icon" :label="t('apps.baseSettings.hideTitle')">
          <NSwitch v-model:value="editModalArg.model.cardStyle.textIconHideTitle" />
        </NFormItem>

        <NFormItem>
          <template #label>
            <NFlex align="center">
              {{ t('common.textColor') }}
              <STip :content="t('apps.itemGroupManage.textAutoColor')" />
            </NFlex>
          </template>
          <NColorPicker
            v-model:value="editModalArg.model.cardStyle.textColor"
            :show-alpha="false"
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
        </NFormItem>

        <NDivider class="!my-2" />

        <NFormItem path="onlyName">
          <template #label>
            <NFlex align="center">
              {{ $t('common.onlyName') }}
              <STip>
                {{ t('apps.itemGroupManage.onlyNameDescribe') }}
                <NButton
                  text
                  tag="a"
                  :href="domainLinks.linkWithLang('only_name_beautify_doc')"
                  target="_blank"
                  type="primary"
                >
                  {{ t("common.lookDoc") }}
                </NButton>
              </STip>
            </NFlex>
          </template>
          <NInput v-model:value="editModalArg.model.onlyName" type="text" show-count :maxlength="20" />
        </NFormItem>
      </NForm>
      <template #footer>
        <NButton type="success" size="small" class="float-right" @click="handleSaveGroup">
          {{ $t('common.save') }}
        </NButton>
      </template>
    </RoundCardModal>
  </div>
</template>

<style scoped>
.card{
  background-color: rgb(241 241 241)
}

.dark .card{
  background-color: #24242a
}
</style>
