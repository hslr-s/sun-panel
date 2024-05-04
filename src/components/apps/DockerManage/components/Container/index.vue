<script lang="ts" setup>
import { h, onMounted, ref } from 'vue'
import { NDataTable, NSpace, NTag, useMessage } from 'naive-ui'
import type { DataTableColumns } from 'naive-ui'
import { useDockerStore } from '../../store/docker'
import { RowAction, RowName } from './table'
import { DockerContainerInfoState } from '@/enums/docker'
import { t } from '@/locales'

const ms = useMessage()
const dockerState = useDockerStore()
const loadingList = ref(false)
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
      title: t('apps.dockerManage.containerName'),
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
      title: t('apps.dockerManage.state'),
      key: 'state',
      render(row) {
        const tagArg: {
          type: 'default' | 'success' | 'error' | 'warning' | 'primary' | 'info' | undefined
          name: string
        } = {
          type: undefined,
          name: t('apps.dockerManage.unknown'),
        }

        switch (row.state) {
          case DockerContainerInfoState.created:
            tagArg.name = t('apps.dockerManage.created')
            tagArg.type = 'info'
            break

          case DockerContainerInfoState.restarting:
            tagArg.name = t('apps.dockerManage.restarting')
            tagArg.type = 'warning'
            break
          case DockerContainerInfoState.removing:
            tagArg.name = t('apps.dockerManage.removing')
            tagArg.type = 'warning'
            break
          case DockerContainerInfoState.paused:
            tagArg.name = t('apps.dockerManage.paused')
            tagArg.type = 'warning'
            break

          case DockerContainerInfoState.exited:
            tagArg.name = t('apps.dockerManage.exited')
            // tagArg.type = 'error' // 默认样式
            break
          case DockerContainerInfoState.dead:
            tagArg.name = t('apps.dockerManage.dead')
            tagArg.type = 'error'
            break

          case DockerContainerInfoState.running:
            tagArg.name = t('apps.dockerManage.running')
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
    // {
    //   title: '资源概览',
    //   key: 'resource',
    // },
    {
      title: t('apps.dockerManage.action'),
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

    // {
    //   title: '更多',
    //   key: 'option',
    //   render(row) {
    //     return h(
    //       NButton,
    //       {
    //         size: 'small',
    //         onClick: () => sendMail(row),
    //       },
    //       { default: () => '...' },
    //     )
    //   },
    // },
  ]
}

const columns = createColumns({

  sendMail(rowData) {
    // ms.info(`send mail to ${rowData.name}`)
  },
})

onMounted(async () => {
  loadingList.value = true
  await dockerState.refreshContainerList()
  loadingList.value = false
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
      :loading="loadingList"
    />
  </NSpace>
</template>
