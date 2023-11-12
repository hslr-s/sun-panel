<script lang="ts" setup>
import { h, onMounted, ref } from 'vue'
import { NButton, NDataTable, NSpace, NTag, useMessage } from 'naive-ui'
import type { DataTableColumns } from 'naive-ui'
import { useDockerStore } from '../../store/docker'
import { getList } from '@/api/docker/container'
import { DockerContainerInfoState } from '@/enum'

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
      title: '容器名称',
      key: 'name',
    },
    {
      title: '状态',
      key: 'state',
      render(row) {
        if (row.state === DockerContainerInfoState.created)
          console.log('sss')

        return h(
          NTag,
          {
            style: {
              marginRight: '6px',
            },
            type: 'info',
            bordered: false,
          },
          {
            default: () => '未知',
          },
        )
        return tags
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
