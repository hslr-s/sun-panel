<script lang="ts" setup>
import { h, onMounted } from 'vue'
import { NButton, NDataTable, NSpace, NTag, useMessage } from 'naive-ui'
import type { DataTableColumns } from 'naive-ui'
import { useDockerStore } from '../../store/docker'
import { RowAction, RowName } from './table'
import { DockerContainerInfoState } from '@/enums/docker'

const ms = useMessage()
const dockerState = useDockerStore()
const pagination = {
  pageSize: 10,
}

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
            tagArg.name = '已暂停'
            tagArg.type = 'warning'
            break

          case DockerContainerInfoState.exited:
            tagArg.name = '已停止'
            // tagArg.type = 'error' // 默认样式
            break
          case DockerContainerInfoState.dead:
            tagArg.name = '已终止'
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
      title: 'Action',
      key: 'action',
      render(row) {
        return h(
          RowAction,
          {
            value: row,
          },
        )
      },
    },

    {
      title: '更多',
      key: 'option',
      render(row) {
        return h(
          NButton,
          {
            size: 'small',
            onClick: () => sendMail(row),
          },
          { default: () => '...' },
        )
      },
    },
  ]
}

const columns = createColumns({
  sendMail(rowData) {
    ms.info(`send mail to ${rowData.name}`)
  },
})

onMounted(() => {
  dockerState.refreshContainerList()
})
</script>

<template>
  <NSpace vertical :size="12">
    <NDataTable
      :bordered="false"
      size="large"
      :data="dockerState.containerList"
      :columns="columns"
      :pagination="pagination"
    />
  </NSpace>
</template>
