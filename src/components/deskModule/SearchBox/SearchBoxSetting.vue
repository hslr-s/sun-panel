<script setup lang="ts">
import { onMounted, ref } from 'vue'
import type {
  FormInst,
  UploadFileInfo,
} from 'naive-ui'
import {
  NAvatar, NButton, NCard, NFlex,
  NForm, NFormItem, NInput, NModal, NUpload,
  useDialog,
} from 'naive-ui'
import type { SearchSaveDataType } from './base'
import { defaultSearchEngineList, defaultSearchSaveData } from './base'
import { t } from '@/locales'
import { useAuthStore, useProAuthorizeStore } from '@/store'
import { useModuleConfig } from '@/store/modules'
import { apiRespErrMsg, message } from '@/utils/cmn/apiMessage'
import { randomCode } from '@/utils/cmn'
import { ProTag } from '@/components/common'

interface State extends SearchSaveDataType {
}

const dialog = useDialog()
const formRef = ref<FormInst | null>(null)
const moduleConfig = useModuleConfig()
const authStore = useAuthStore()
const proAuthStore = useProAuthorizeStore()
const editModalShow = ref(false)
const editModalIndex = ref(0)
const searchSaveData = ref<State>({ ...defaultSearchSaveData })

const modelDefault: DeskModule.SearchBox.SearchEngine = {
  iconSrc: '',
  title: '',
  url: '',
  key: '',
}
const editModalModel = ref<DeskModule.SearchBox.SearchEngine>({ ...modelDefault })
const editModalTitle = ref('')
const editModalIsEdit = ref(false)

const rules = {
  iconSrc: {
    required: true,
    message: t('form.required'),
    trigger: ['input', 'blur'],
  },
  title: {
    required: true,
    message: t('form.required'),
    trigger: ['input', 'blur'],
  },
  url: {
    required: true,
    message: t('form.required'),
    trigger: ['input', 'blur'],
  },
}
function handleAdd() {
  if (proAuthStore.isExpired && searchSaveData.value.searchEngineList.length >= 2) {
    message.warning(t('deskModule.searchBox.noProAuth'))
    return
  }

  editModalModel.value = { ...modelDefault }
  editModalTitle.value = t('common.add')
  editModalIsEdit.value = false
  editModalIndex.value = 0
  editModalShow.value = true
}

function handleEdit(item: DeskModule.SearchBox.SearchEngine, index: number) {
  editModalModel.value = { ...item }
  editModalTitle.value = t('common.edit')
  editModalIsEdit.value = true
  editModalIndex.value = index
  editModalShow.value = true
}

function handleUploadBackgroundFinish({
  file,
  event,
}: {
  file: UploadFileInfo
  event?: ProgressEvent
}) {
  const res = JSON.parse((event?.target as XMLHttpRequest).response)
  if (res.code === 0)
    editModalModel.value.iconSrc = res.data.imageUrl

  else
    apiRespErrMsg(res)

  return file
}

function getList() {
  moduleConfig.getValueByNameFromCloud<State>('deskModuleSearchBox').then(({ code, data }) => {
    if (code === 0)
      searchSaveData.value = data || defaultSearchSaveData
    else
      searchSaveData.value = defaultSearchSaveData

    if (!searchSaveData.value.searchEngineList)
      searchSaveData.value.searchEngineList = []
  })
}

function save(e: MouseEvent) {
  e.preventDefault()

  // 生成唯一的key值
  if (!editModalModel.value.key || editModalModel.value.key === '')
    editModalModel.value.key = randomCode(15)

  formRef.value?.validate((errors) => {
    if (!errors) {
      if (editModalIsEdit.value) {
        searchSaveData.value.searchEngineList[editModalIndex.value] = editModalModel.value

        if (searchSaveData.value.currentSearchEngine.key || searchSaveData.value.currentSearchEngine.key !== '') {
          // 如果当前使用的搜索引擎是当前的修改项，将对当前所使用的搜索引擎同步更新
          for (let i = 0; i < searchSaveData.value.searchEngineList.length; i++) {
            const element = searchSaveData.value.searchEngineList[i]
            if (searchSaveData.value.currentSearchEngine.key === element.key)
              searchSaveData.value.currentSearchEngine = element
          }
        }
      }

      else {
        searchSaveData.value.searchEngineList.push(editModalModel.value)
      }

      moduleConfig.saveToCloud('deskModuleSearchBox', searchSaveData.value)
      editModalShow.value = false
    }
    else {
      console.log(errors)
    }
  })
}

function deleteByIndex(index: number) {
  dialog.warning({
    title: t('common.warning'),
    content: t('common.deleteConfirm'),
    positiveText: t('common.confirm'),
    negativeText: t('common.cancel'),
    onPositiveClick: () => {
      searchSaveData.value.searchEngineList.splice(index, 1)
      moduleConfig.saveToCloud('deskModuleSearchBox', searchSaveData.value)
      editModalShow.value = false
    },

  })
}

onMounted(() => {
  getList()
})
</script>

<template>
  <div>
    <NFlex align="center">
      <NButton style="margin-right: 2px;" size="small" type="success" @click="handleAdd">
        {{ t("common.add") }}
      </NButton>

      <div class="flex items-center">
        <ProTag class="mx-2" show-icon :content="t('deskModule.searchBox.noProAuth')" />
      </div>
    </NFlex>

    <!-- list -->
    <div class="h-[500px] overflow-auto mt-5">
      <!-- 默认的搜索引擎，不可删除和修改 -->
      <NCard
        v-for="item, index in defaultSearchEngineList" :key="index" size="small"
        style="margin:4px 0" embedded
      >
        <div class="flex items-center">
          <NAvatar :src="item.iconSrc" style="background-color: transparent;" :size="20" />
          <div class="ml-2 w-[80px]">
            {{ item.title }}
          </div>
          <div class="ml-2 w-auto">
            {{ item.url }}
          </div>
        </div>
      </NCard>

      <!-- 用户自定义的搜索引擎 -->
      <NCard
        v-for="item, index in searchSaveData.searchEngineList" :key="index" size="small"
        style="margin:4px 0;" embedded
      >
        <div class="flex items-center">
          <NAvatar :src="item.iconSrc" style="background-color: transparent;" :size="20" />
          <div class="ml-2 w-[80px]">
            {{ item.title }}
          </div>
          <div class="ml-2 w-auto">
            {{ item.url }}
          </div>

          <div class="ml-auto">
            <div class="flex justify-end">
              <NButton size="tiny" style="margin-right: 2px;" type="success" quaternary @click="handleEdit(item, index)">
                {{ t('common.edit') }}
              </NButton>
              <NButton size="tiny" type="error" quaternary @click="deleteByIndex(index)">
                {{ t('common.delete') }}
              </NButton>
            </div>
          </div>
        </div>
      </NCard>
    </div>

    <!-- 添加/编辑窗口 -->
    <NModal
      v-model:show="editModalShow"
      preset="card"
      size="small"
      style="width: 500px;border-radius: 1rem;"
      :title="editModalTitle"
    >
      <NForm
        ref="formRef"
        :model="editModalModel"
        label-placement="left"
        label-width="auto"
        require-mark-placement="right-hanging"
        size="small"
        :rules="rules"
      >
        <NFormItem :label="t('common.icon')" path="iconSrc">
          <div class="w-full">
            <div class="mb-2 flex items-center">
              <NAvatar :src="editModalModel.iconSrc" :size="50" />
              <div class="ml-5">
                <NUpload
                  action="/api/file/uploadImg"
                  name="imgfile"
                  :headers="{
                    token: authStore.token as string,
                  }"
                  :show-file-list="false"
                  @finish="handleUploadBackgroundFinish"
                >
                  <NButton type="info" size="tiny">
                    {{ t('deskModule.searchBox.uploadImage') }}
                  </NButton>
                </NUpload>
              </div>
            </div>
            <div class="text-xs my-1">
              {{ t('deskModule.searchBox.iconUploadHint') }}
            </div>
            <NInput v-model:value="editModalModel.iconSrc" placeholder="Icon URL" />
          </div>
        </NFormItem>

        <NFormItem :label="t('common.title')" path="title">
          <NInput v-model:value="editModalModel.title" />
        </NFormItem>

        <NFormItem label="URL" path="url">
          <div class="w-full">
            <NInput v-model:value="editModalModel.url" />
            <div class="text-xs mt-1">
              {{ t('deskModule.searchBox.urlKeywordhint') }}
            </div>
          </div>
        </NFormItem>
      </NForm>

      <template #footer>
        <NFlex justify="end">
          <NButton type="success" @click="save">
            {{ t('common.save') }}
          </NButton>
        </NFlex>
      </template>
    </NModal>
  </div>
</template>
