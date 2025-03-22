<script lang="ts" setup>
import { h, nextTick, onMounted, onUnmounted, ref } from 'vue'
import { NAlert, NDataTable, NEllipsis, NTag } from 'naive-ui'
import type { DataTableColumns } from 'naive-ui'
import type { ComponentPublicInstance } from 'vue'
import { useDockerStore } from '../../store/docker'
import { RowAction, RowName } from './table'
import { DockerContainerInfoState } from '@/enums/docker'
import { t } from '@/locales'
import { ProTag, SvgIcon } from '@/components/common'
import { useDomainLinks } from '@/hooks/useDomainLinks'
import { PanelCardType, PanelItemIconTypeEnum, PanelItemOpenMethodEnum } from '@/enums'
import { useTriggerUpdateItemData } from '@/hooks/useTriggerUpdateItemData'
import { EditItem } from '@/views/home/components'

const loadingList = ref(false)
let intervalId: any
const dockerState = useDockerStore()
const domainLinks = useDomainLinks()
const nSpace = ref<HTMLElement | ComponentPublicInstance | null>(null)
const dataTableHeight = ref(0)

const alertBoxRef = ref<HTMLElement | ComponentPublicInstance | null>(null)
const editItemInfo = ref<Panel.Info | null>(null)
const editItemInfoShow = ref(false)

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
      sorter: 'default',
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
    {
      key: 'add',
      fixed: 'right',
      width: 50,
      align: 'center',
      render(row) {
        return h('div', {
          title: t('common.addCardToGroup'),
          style: {
            display: 'flex',
            justifyContent: 'center',
          },
          onClick: () => addToGroupCard(row),
        },
        h(
          SvgIcon,
          {
            class: 'add-card',
            style: {
              width: '20px',
              height: '20px',
              cursor: 'pointer',
            },
            icon: 'fluent--square-add-16-regular',

          },
        ))
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
      dataTableHeight.value = nSpaceHeight - 60 - alertboxHeight
    }
  })
}

// 添加卡片分组
function addToGroupCard(row: Docker.Container.Info) {
  editItemInfo.value = {
    title: row.name,
    icon: {
      text: row.name,
      itemType: PanelItemIconTypeEnum.url,
    },
    cardType: PanelCardType.Docker,
    openMethod: PanelItemOpenMethodEnum.newTab,
    url: '',
    expandParam: {
      containerId: row.id,
      containerName: row.name,
    },
  }
  editItemInfoShow.value = true
}

function handleEditSuccess() {
  editItemInfoShow.value = false
  useTriggerUpdateItemData().triggerUpdate()
}
onUnmounted(() => {
  clearInterval(intervalId)
  window.removeEventListener('resize', handleResize)
})
</script>

<template>
  <div class="h-full overflow-auto">
    <!-- 没有PRO授权仅能查看状态，不可以操作。 -->
    <div class="flex flex-col h-full gap-2">
      <div>
        <NAlert type="info" closable @after-leave="handleResize">
          {{ t('apps.dockerManage.alertText') }}
        </NAlert>

        <NAlert
          v-if="dockerState.refreshContainerListErrorMsg" class="mt-2" type="error" closable
          @after-leave="handleResize"
        >
          {{ dockerState.refreshContainerListErrorMsg }}
          <a target="_blank" class="text-lime-600" :href="domainLinks.faqLinkWithLang('docker_run_error')">
            [{{ t('common.viewSolution') }}]
          </a>
        </NAlert>
      </div>

      <div class="h-full flex">
        <NDataTable
          flex-height
          :bordered="false"
          size="small" :data="dockerState.containerList"
          :columns="columns" :loading="loadingList"
          :scroll-x="800"
        />
      </div>
    </div>

    <EditItem v-model:visible="editItemInfoShow" :item-info="editItemInfo" is-out-add @done="handleEditSuccess" />
  </div>
</template>

<style scoped>
:deep(.add-card:hover) {
  color: var(--n-color-target);
}
:deep(.add-card) {
  transition: color 0.3s ease;
}
</style>
