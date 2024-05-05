<script setup lang="ts">
import { onMounted, ref, watch } from 'vue'
import type { FormInst, FormRules } from 'naive-ui'
import {
  NButton, NCard, NColorPicker, NDivider,
  NForm, NFormItem, NInput, NSelect, NSwitch,
  useDialog,
  useMessage,
} from 'naive-ui'
import { VueDraggable } from 'vue-draggable-plus'
import { deletes, edit, getList, saveSort } from '@/api/panel/itemIconGroup'
import { RoundCardModal, SvgIcon } from '@/components/common'
import { t } from '@/locales'
import { PanelPanelConfigStyleEnum } from '@/enums'
import { defaultItemGroupInfoCardStyle } from '@/utils/defaultData'

interface EditModalArg {
  show: boolean
  editStatus: number // 1.添加 2.编辑
  model: Panel.ItemGroupInfo
  rules: FormRules
}

const formRef = ref<FormInst | null>(null)
const ms = useMessage()
const dialog = useDialog()
const sortStatus = ref(false)

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
        }).catch(() => {
          ms.error(t('common.deleteFail'))
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
        editModalArg.value.show = false
        editModalArg.value.model = { ...defaultItemGroupInfo }
      }).catch(() => {
        ms.error(t('common.saveFail'))
      })
    }
    else { console.log(errors) }
  })
}

function refreshList() {
  getList<Common.ListResponse<Panel.ItemGroupInfo[]>>().then(({ code, data }) => {
    groups.value = data.list
  })
}

onMounted(() => {
  refreshList()
})

watch(() => editModalArg.value.show, (v) => {
  console.log(v)
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
                <div class="flex items-center">
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
      <NForm ref="formRef" :model="editModalArg.model" :rules="editModalArg.rules">
        <NFormItem path="title" :label="$t('apps.itemGroupManage.groupName')">
          <NInput v-model:value="editModalArg.model.title" type="text" :maxlength="20" show-count />
        </NFormItem>

        <NDivider style="margin-top:5px" />

        <NFormItem :label="t('apps.itemGroupManage.cardStyle')">
          <NSelect v-model:value="editModalArg.model.cardStyle.style" :options="iconTypeOptions" />
        </NFormItem>

        <NFormItem v-if="editModalArg.model.cardStyle.style === PanelPanelConfigStyleEnum.info" :label="t('apps.baseSettings.hideDescription')">
          <NSwitch v-model:value="editModalArg.model.cardStyle.textInfoHideDescription" />
        </NFormItem>

        <NFormItem v-if="editModalArg.model.cardStyle.style === PanelPanelConfigStyleEnum.icon" :label="t('apps.baseSettings.hideTitle')">
          <NSwitch v-model:value="editModalArg.model.cardStyle.textIconHideTitle" />
        </NFormItem>

        <NFormItem :label="t('common.textColor')">
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
      </NForm>
      <template #footer>
        <NButton type="success" size="small" class="float-right" @click="handleSaveGroup">
          {{ $t('common.confirm') }}
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
