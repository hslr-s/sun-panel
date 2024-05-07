import { defineStore } from 'pinia'
import { ref } from 'vue'
import { useMessage } from 'naive-ui'
import { getContainersStats as getContainersStatsApi, getList } from '@/api/docker/container'
import { t } from '@/locales'
import { DockerContainerInfoState } from '@/enums/docker'

export const useDockerStore = defineStore('docker-store', () => {
  const containerList = ref<Docker.Container.Info[]>([])
  const actionMenuBar = ref('container')
  const isLoading = ref(false)

  const containerRuningIds = ref<string[]>([])
  const containerRuningStats = ref<{ [key: string]: Docker.Container.Stats }>({})

  const ms = useMessage()

  async function refreshContainerList() {
    containerRuningIds.value = []
    await getList<Common.ListResponse<Docker.Container.Info[]>>().then((res) => {
      containerList.value = res.data.list
      for (let i = 0; i < res.data.list.length; i++) {
        const element = res.data.list[i]
        containerRuningStats.value[element.id] = {
        } as Docker.Container.Stats
        if (element.state === DockerContainerInfoState.running) {
          containerRuningIds.value.push(element.id)
          console.log(element.id)
        }
      }
      refreshContainerListStats()
    }).catch((res) => {
      if (res.code === -1) {
        ms.error(`${t('apps.dockerManage.dockerError')}: ${res.msg}`, {
          closable: true,
          duration: 10000,
          keepAliveOnHover: true,
        })
      }
    })
  }

  async function refreshContainerListStats() {
    await getContainersStatsApi<{ [key: string]: Docker.Container.Stats }>(containerRuningIds.value).then((res) => {
      containerRuningStats.value = res.data
    }).catch((res) => {
      // if (res.code === -1) {
      //   ms.error(`${t('apps.dockerManage.dockerError')}: ${res.msg}`, {
      //     closable: true,
      //     duration: 10000,
      //     keepAliveOnHover: true,
      //   })
      // }
    })
  }

  return {
    actionMenuBar,
    refreshContainerList,
    containerRuningStats,
    refreshContainerListStats,
    containerList,
    isLoading,
  }
})
