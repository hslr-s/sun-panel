<script setup lang="ts">
import { onMounted, ref } from 'vue'
import type { FormInst, FormRules } from 'naive-ui'
import { NButton, NCard, NForm, NFormItem, NInput, useDialog, useMessage } from 'naive-ui'
import { VueDraggable } from 'vue-draggable-plus'
import { deletes, edit, getList } from '@/api/panel/itemIconGroup'
import { RoundCardModal, SvgIcon } from '@/components/common'

interface EditModalArg {
  show: boolean
  editStatus: number // 1.添加 2.编辑
  model: Panel.ItemIconGroup
  rules: FormRules
}

const formRef = ref<FormInst | null>(null)
const ms = useMessage()
const dialog = useDialog()

const defaultMNodal = {
  title: '',
  icon: 'material-symbols:folder-outline',
  sort: 9999,
}

const editModalArg = ref<EditModalArg>({
  show: false,
  editStatus: 1,
  model: defaultMNodal,
  rules: {
    title: [
      {
        required: true,
        trigger: 'blur',
        message: '必填项',
      },
    ],
  },
})

const groups = ref<Panel.ItemIconGroup[]>([])

function handleAddGroup() {
  editModalArg.value.show = !editModalArg.value.show
}

function handleEditGroup(groupInfo: Panel.ItemIconGroup) {
  editModalArg.value.show = true
  editModalArg.value.model = groupInfo
}

function handleDelete(groupInfo: Panel.ItemIconGroup) {
  dialog.warning({
    title: '警告',
    content: `你确定删除此分组[ ${groupInfo.title} ]，删除后此分组应用图标将丢失？`,
    positiveText: '确定',
    negativeText: '取消',
    onPositiveClick: () => {
      if (groupInfo.id) {
        deletes([groupInfo.id]).then(({ code, msg }) => {
          if (code !== 0)
            ms.error(`删除失败:${msg}`)
          else
            refreshList()
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
          ms.error(msg)

        refreshList()
        editModalArg.value.show = false
        editModalArg.value.model = { ...defaultMNodal }
      })
    }
    else { console.log(errors) }
  })
}

function refreshList() {
  getList<Common.ListResponse<Panel.ItemIconGroup[]>>().then(({ code, data }) => {
    groups.value = data.list
  })
}

onMounted(() => {
  refreshList()
})
</script>

<template>
  <div class="h-[500px]">
    <NButton type="success" size="small" @click="handleAddGroup">
      新增分组
    </NButton>

    <VueDraggable
      v-model="groups"
      item-key="sort" :animation="300"
      class="mt-[20px] p-[20px] bg-slate-200 rounded-xl"
    >
      <div v-for="(item, index) in groups" :key="index" class="w-[60%]">
        <NCard size="small" style="border-radius:10px;margin-bottom: 10px;">
          <div class="flex">
            <div class="flex items-center">
              <span class="mr-[10px]">
                <SvgIcon class="text-[20px]" :icon="item.icon" />
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
    <RoundCardModal v-model:show="editModalArg.show" type="small" :title="editModalArg.editStatus === 1 ? '添加' : '编辑'" style="width: 400px;">
      <NForm ref="formRef" :model="editModalArg.model" :rules="editModalArg.rules">
        <NFormItem path="title" label="分组名称">
          <NInput v-model:value="editModalArg.model.title" type="text" :maxlength="20" show-count placeholder="请输入" />
        </NFormItem>

        <!-- <NFormItem path="name" label="昵称">
          <NInput v-model:value="editModalArg.model" type="text" placeholder="请输入昵称" />
        </NFormItem> -->
      </NForm>
      <template #footer>
        <NButton type="success" size="small" class="float-right" @click="handleSaveGroup">
          确定
        </NButton>
      </template>
    </RoundCardModal>
  </div>
</template>
