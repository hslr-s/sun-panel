<script lang="ts" setup>
import { h, onMounted, ref } from 'vue'
import { NButton, NDataTable, NSpace, NTag, useMessage } from 'naive-ui'
import type { DataTableColumns } from 'naive-ui'
import { useDockerStore } from '../../store/docker'
import { RowName } from './table'
import { getList } from '@/api/docker/container'
import { DockerContainerInfoState } from '@/enums/docker'

const ms = useMessage()
const dockerState = useDockerStore()
const pagination = {
  pageSize: 10,
}
const data = ref<Docker.Container.Info[]>()

const createColumns = ({
  sendMail,
}: {
  sendMail: (rowData: Docker.Container.Info) => void
}): DataTableColumns<Docker.Container.Info> => {
  return [
    {
      title: '名称',
      key: 'name',
      render(row) {
        return h(
          RowName,
          {
            value: row,
          },
        )
      },
    },
    {
      title: '状态',
      key: 'state',
      render(row) {
        const tagArg: {
          type: 'default' | 'success' | 'error' | 'warning' | 'primary' | 'info' | undefined
          name: string
        } = {
          type: undefined,
          name: '未知',
        }

        switch (row.state) {
          case DockerContainerInfoState.created:
            tagArg.name = '已创建'
            tagArg.type = 'info'
            break

          case DockerContainerInfoState.restarting:
            tagArg.name = '重启中'
            tagArg.type = 'warning'
            break
          case DockerContainerInfoState.removing:
            tagArg.name = '迁移中'
            tagArg.type = 'warning'
            break
          case DockerContainerInfoState.paused:
            tagArg.type = 'warning'
            break

          case DockerContainerInfoState.exited:
            tagArg.name = '停止'
            break
          case DockerContainerInfoState.dead:
            tagArg.name = '终止'
            tagArg.type = 'error'
            break

          case DockerContainerInfoState.running:
            tagArg.name = '运行中'
            tagArg.type = 'success'
            break
        }

        return h(
          NTag,
          {
            type: tagArg.type,
            bordered: false,
            round: true,
            size: 'small',
            title: row.status,
          },
          {
            default: () => tagArg.name,
          },
        )
      },
    },
    {
      title: '资源概览',
      key: 'resource',
    },
    {
      title: 'Tags',
      key: 'tags',

    },
    {
      title: 'Action',
      key: 'actions',
      render(row) {
        return h(
          NButton,
          {
            size: 'small',
            onClick: () => sendMail(row),
          },
          { default: () => 'Send Email' },
        )
      },
    },
  ]
}

const columns = ref(createColumns({
  sendMail(rowData) {
    ms.info(`send mail to ${rowData.name}`)
  },
}))

// const data: RowData[] = [
//   {
//     key: 0,
//     name: 'John Brown',
//     age: 32,
//     address: 'New York No. 1 Lake Park',
//     tags: ['nice', 'developer'],
//   },
//   {
//     key: 1,
//     name: 'Jim Green',
//     age: 42,
//     address: 'London No. 1 Lake Park',
//     tags: ['wow'],
//   },
//   {
//     key: 2,
//     name: 'Joe Black',
//     age: 32,
//     address: 'Sidney No. 1 Lake Park',
//     tags: ['cool', 'teacher'],
//   },
// ]

onMounted(() => {
  getList<Common.ListResponse<Docker.Container.Info[]>>().then((res) => {
    if (res.code === 0)
      data.value = res.data.list
  })
})
</script>

<template>
  <NSpace vertical :size="12">
    <NDataTable
      :bordered="false"
      size="large"
      :data="data"
      :columns="columns"
      :pagination="pagination"
    />
  </NSpace>
</template>
