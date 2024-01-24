<script lang="ts" setup>
import { h, onMounted, reactive, ref } from 'vue'
import { NAlert, NButton, NDataTable, NDropdown, NTag, useDialog, useMessage } from 'naive-ui'
import type { DataTableColumns, PaginationProps } from 'naive-ui'
import EditUser from './EditUser/index.vue'
import { getPublicVisitUser, setPublicVisitUser, deletes as usersDeletes, getList as usersGetList } from '@/api/panel/users'
import { SvgIcon } from '@/components/common'
import { useAuthStore } from '@/store'
import { t } from '@/locales'
import { AdminAuthRole } from '@/enums/admin'

const message = useMessage()
const authStore = useAuthStore()
const tableIsLoading = ref<boolean>(false)
const editUserDialogShow = ref<boolean>(false)
const keyWord = ref<string>()
const editUserUserInfo = ref<User.Info>()
const dialog = useDialog()
const publicVisitUserId = ref<number | null>(null)

const createColumns = ({
  update,
}: {
  update: (row: User.Info) => void
}): DataTableColumns<User.Info> => {
  return [
    {
      title: t('common.username'),
      key: 'username',
      render(row: User.Info) {
        let publicVisitHtml = ''
        if (publicVisitUserId.value && publicVisitUserId.value === row.id)
          publicVisitHtml = `[${t('adminSettingUsers.pblicText')}]-`

        if (row.username === authStore.userInfo?.username)
          return `${publicVisitHtml}${row.username} (${t('adminSettingUsers.currentUseUsername')})`
        return publicVisitHtml + row.username
      },
    },
    {
      title: t('common.nikeName'),
      key: 'name',
    },
    {
      title: t('adminSettingUsers.role'),
      key: 'role',
      render(row) {
        switch (row.role) {
          case AdminAuthRole.admin:
            return h(NTag, { type: 'info' }, t('common.role.admin'))
          case AdminAuthRole.regularUser:
            return h(NTag, t('common.role.regularUser'))
          default:
            return '-'
        }
      },
    },
    {
      title: t('common.action'),
      key: '',
      render(row) {
        const btn = h(
          NButton,
          {
            strong: true,
            tertiary: true,
            size: 'small',
          },
          {
            default() {
              return h(
                SvgIcon, {
                  icon: 'mingcute:more-1-fill',
                },
              )
            },
          },
        )

        return h(NDropdown, {
          trigger: 'click',
          onSelect(key: string | number) {
            switch (key) {
              case 'update':
                update(row)
                break
              case 'publicMode':
                // 取消
                if (publicVisitUserId.value && publicVisitUserId.value === row.id) {
                  setPublicVisitUser(null).then(({ code }) => {
                    if (code === 0)
                      publicVisitUserId.value = null
                  })
                }
                else {
                // 设置
                  setPublicVisitUser(row.id as number).then(({ code }) => {
                    if (code === 0)
                      publicVisitUserId.value = row.id as number
                  })
                }
                break
              case 'delete':
                dialog.warning({
                  title: t('common.warning'),
                  content: t('adminSettingUsers.deletePromptContent', { name: row.name, username: row.username }),
                  positiveText: t('common.confirm'),
                  negativeText: t('common.cancel'),
                  onPositiveClick: () => {
                    deletes([row.id as number])
                  },
                })
                break

              default:
                break
            }
          },
          options: [
            {
              label: t('common.edit'),
              key: 'update',
            },
            {
              label: t('adminSettingUsers.setOrUnsetPublicMode'),
              key: 'publicMode',
            },
            {
              label: t('common.delete'),
              key: 'delete',
            },
          ],
        }, { default: () => btn })
      },
    },
  ]
}

const userList = ref<User.Info[]>()

const columns = createColumns({
  update(row: User.Info) {
    editUserUserInfo.value = row
    editUserDialogShow.value = true
  },
})
const pagination = reactive({
  page: 1,
  showSizePicker: true,
  pageSizes: [10, 30, 50, 100, 200],
  pageSize: 10,
  itemCount: 0,
  onChange: (page: number) => {
    pagination.page = page
    getList(null)
  },
  onUpdatePageSize: (pageSize: number) => {
    pagination.pageSize = pageSize
    pagination.page = 1
    getList(null)
  },
  prefix(item: PaginationProps) {
    return t('adminSettingUsers.userCountText', { count: item.itemCount })
  },
})

function handlePageChange(page: number) {
  getList(page)
}

// 添加
function handleAdd() {
  editUserDialogShow.value = true
  editUserUserInfo.value = {}
}

function handelDone() {
  editUserDialogShow.value = false
  message.success(t('common.success'))
  getList(null)
}

async function getList(page: number | null) {
  tableIsLoading.value = true
  const req: AdminUserManage.GetListRequest = {
    page: page || pagination.page,
    limit: pagination.pageSize,
  }
  if (keyWord.value !== '')
    req.keyWord = keyWord.value

  const { data } = await usersGetList<Common.ListResponse<User.Info[]>>(req)
  pagination.itemCount = data.count
  if (data.list)
    userList.value = data.list
  tableIsLoading.value = false
}

async function deletes(ids: number[]) {
  const { code } = await usersDeletes(ids)
  if (code === 0) {
    message.success(t('common.deleteSuccess'))
    getList(null)
  }
}

onMounted(() => {
  getPublicVisitUser<User.Info>().then(({ data }) => {
    publicVisitUserId.value = data.id || null
  })
  getList(null)
})
</script>

<template>
  <div class="overflow-auto pt-2">
    <NAlert type="info" :bordered="false">
      {{ $t('adminSettingUsers.alertText') }}
    </NAlert>
    <div class="my-[10px]">
      <NButton type="primary" size="small" ghost @click="handleAdd">
        {{ $t('common.add') }}
      </NButton>
    </div>

    <NDataTable
      :columns="columns"
      :data="userList"
      :pagination="pagination"
      :bordered="false"
      :loading="tableIsLoading"
      :remote="true"

      @update:page="handlePageChange"
    />
    <EditUser v-model:visible="editUserDialogShow" :user-info="editUserUserInfo" @done="handelDone" />
  </div>
</template>
