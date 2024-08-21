<script lang="ts" setup>
import { h, onMounted, onUnmounted, ref } from 'vue'
import { NAlert, NDataTable, NEllipsis, NSpace, NTag } from 'naive-ui'
import type { DataTableColumns } from 'naive-ui'
import { useDockerStore } from '../../store/docker'
import { RowAction, RowName } from './table'
import { DockerContainerInfoState } from '@/enums/docker'
import { t } from '@/locales'
import { ProTag } from '@/components/common'
import { useDomainLinks } from '@/hooks/useDomainLinks'

const loadingList = ref(false)
let intervalId: any
const dockerState = useDockerStore()
const domainLinks = useDomainLinks()

const pagination = {
  pageSize: 10,
}

// async function getContainerList() {
//   await getContainerListApi<Common.ListResponse<Docker.Container.Info[]>>().then((res) => {
//     containerList.value = res.data.list
//     for (let i = 0; i < res.data.list.length; i++) {
//       const element = res.data.list[i]
//       containerRuningState.value[element.id] = {
//       } as Docker.Container.Stats
//       if (element.state === DockerContainerInfoState.running) {
//         containerRuningIds.value.push(element.id)

//         getContainersStats()
//       }
//     }
//   }).catch((res) => {
//     if (res.code === -1) {
//       ms.error(`${t('apps.dockerManage.dockerError')}: ${res.msg}`, {
//         closable: true,
//         duration: 10000,
//         keepAliveOnHover: true,
//       })
//     }
//   })
// }

// async function getContainersStats() {
//   await getContainersStatsApi<{ [key: string]: Docker.Container.Stats }>(containerRuningIds.value).then((res) => {
//     containerRuningState.value = res.data
//   }).catch((res) => {
//     // if (res.code === -1) {
//     //   ms.error(`${t('apps.dockerManage.dockerError')}: ${res.msg}`, {
//     //     closable: true,
//     //     duration: 10000,
//     //     keepAliveOnHover: true,
//     //   })
//     // }
//   })
// }

const createColumns = ({
  sendMail,
}: {
  sendMail: (rowData: Docker.Container.Info) => void
}): DataTableColumns<Docker.Container.Info> => {
  return [
    {
      title: t('apps.dockerManage.containerName'),
      key: 'name',
      width: 400,
      render(row) {
        const rowInfo = h(
          RowName,
          {
            value: row,
          },
        )
        return h(
          'div',
          {
            style: { minWidth: 200 },
          },
          rowInfo,
        )
      },
    },

    {
      title: t('apps.dockerManage.state'),
      align: 'center',
      key: 'state',
      width: 100,
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
    {
      title: t('apps.dockerManage.ports'),
      key: 'name',
      width: 60,
      render(row) {
        const group = (g: Docker.Container.Port) => {
          return h(NTag,
            {
              size: 'tiny',
              style: { margin: '0 2px' },
              type: 'info',
              title: `(${g.type}) ${g.ip}:${g.publicPort} => ${g.privatePort}`,
            },
          `${g.type} ${g.ip}:${g.publicPort}`)
        }

        const c = []
        for (const key in row.ports) {
          if (Object.prototype.hasOwnProperty.call(row.ports, key)) {
            const element = row.ports[key]
            if (element.publicPort !== 0 && element.privatePort !== 0) {
              c.push(group(element))
            }
          }
        }

        return h(NEllipsis, { lineClamp: 2 }, c)
      },
    },
    {
      title: t('apps.dockerManage.resourceOverview'),
      key: 'resource',
      width: 200,
      render(row) {
        const status = dockerState.containerRuningStats[row.id]
        let cpu = ''
        let RAM = ''

        if (row.status === DockerContainerInfoState.running) {
          cpu = 'CPU: --'
          RAM = 'RAM: --'
        }

        if (status) {
          if (status.cpuPercent || status.cpuPercent === 0) {
            status.cpuPercent = roundNumberToTwoDecimalPlaces(status.cpuPercent || 0)
            cpu = `CPU: ${status.cpuPercent}%`
          }

          if (status.memory || status.memory === 0) {
            status.memory = roundNumberToTwoDecimalPlaces(status.memory || 0)
            RAM = `RAM: ${status.memory}MB`
          }
        }

        return h('div',
          [
            h(
              'div', { class: 'text-xs' },
              cpu,
            ),
            h(
              'div', { class: 'text-xs' },
              RAM,
            ),
          ],
        )
      },
    },
    {
      title() {
        return [
          h('span', t('apps.dockerManage.action')),
          h('span', {
            style: 'max-width:80px',
          }, h(
            ProTag,
            {
              style: 'margin-left:8px',
            },
          )),
        ]
      },
      key: 'action',
      fixed: 'right',
      width: 100,
      align: 'center',
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

function roundNumberToTwoDecimalPlaces(num: number): number {
  return Number(num.toFixed(2))
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

  intervalId = setInterval(() => {
    dockerState.refreshContainerListStats()
  }, 5000) // 5000毫秒，即5秒
})

onUnmounted(() => {
  clearInterval(intervalId)
})
</script>

<template>
  <NSpace vertical :size="12">
    <!-- 没有PRO授权仅能查看状态，不可以操作。 -->
    <NAlert type="info" closable>
      {{ t('apps.dockerManage.alertText') }}
    </NAlert>

    <NAlert v-if="dockerState.refreshContainerListErrorMsg" type="error" closable>
      {{ dockerState.refreshContainerListErrorMsg }}
      <a target="_blank" class="text-lime-600" :href="domainLinks.faqLinkWithLang('docker_run_error')">
        [{{ t('common.viewSolution') }}]
      </a>
    </NAlert>

    <NDataTable
      :bordered="false" size="small"
      :data="dockerState.containerList"
      :columns="columns"
      :pagination="pagination"
      :loading="loadingList"
    />
  </NSpace>
</template>
