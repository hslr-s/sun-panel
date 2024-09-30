<script lang="ts" setup>
import { h, nextTick, onMounted, onUnmounted, ref } from 'vue'
import { NAlert, NDataTable, NEllipsis, NSpace, NTag } from 'naive-ui'
import type { DataTableColumns } from 'naive-ui'
import type { ComponentPublicInstance } from 'vue'
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
const nSpace = ref<HTMLElement | ComponentPublicInstance | null>(null)
const dataTableHeight = ref(0)

const alertBoxRef = ref<HTMLElement | ComponentPublicInstance | null>(null)

const createColumns = ({
  sendMail,
}: {
  sendMail: (rowData: Docker.Container.Info) => void
}): DataTableColumns<Docker.Container.Info> => {
  return [
    {
      title: t('apps.dockerManage.containerName'),
      key: 'name',
      minWidth: 300,
      sorter: 'default',
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
      key: 'port',
      width: 200,
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
  handleResize()
  window.addEventListener('resize', handleResize)
  loadingList.value = true
  await dockerState.refreshContainerList()
  loadingList.value = false

  intervalId = setInterval(() => {
    dockerState.refreshContainerListStats()
  }, 5000) // 5000毫秒，即5秒
})

function handleResize() {
  nextTick(() => {
    if (nSpace.value && alertBoxRef.value) {
      const nspaceElement = (nSpace.value as ComponentPublicInstance)?.$el || nSpace.value
      const alertboxElement = (alertBoxRef.value as ComponentPublicInstance)?.$el || alertBoxRef.value
      const nSpaceHeight = nspaceElement ? nspaceElement.offsetHeight : 0
      const alertboxHeight = alertboxElement ? alertboxElement?.offsetHeight : 0
      dataTableHeight.value = nSpaceHeight - 50 - alertboxHeight
    }
  })
}

onUnmounted(() => {
  clearInterval(intervalId)
  window.removeEventListener('resize', handleResize)
})
</script>

<template>
  <NSpace ref="nSpace" vertical :size="12" style="height: 100%;">
    <!-- 没有PRO授权仅能查看状态，不可以操作。 -->

    <div ref="alertBoxRef">
      <NAlert type="info" closable @after-leave="handleResize">
        {{ t('apps.dockerManage.alertText') }}
      </NAlert>

      <NAlert v-if="dockerState.refreshContainerListErrorMsg" type="error" closable @after-leave="handleResize">
        {{ dockerState.refreshContainerListErrorMsg }}
        <a target="_blank" class="text-lime-600" :href="domainLinks.faqLinkWithLang('docker_run_error')">
          [{{ t('common.viewSolution') }}]
        </a>
      </NAlert>
    </div>

    <NDataTable
      :max-height="dataTableHeight"
      :bordered="false" size="small"
      :data="dockerState.containerList"
      :columns="columns"
      :loading="loadingList"
    />
  </NSpace>
</template>
