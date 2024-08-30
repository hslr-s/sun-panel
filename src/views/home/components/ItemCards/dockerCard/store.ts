import { defineStore } from 'pinia'
import { onMounted, ref } from 'vue'
import { getContainersStats as getContainersStatsApi } from '@/api/docker/container'

export const useDockerCardStore = defineStore('docker-card', () => {
  interface ContainerQueryStatsNum {
    [key: string]: number
  }

  interface ContainerQueryStats {
    [key: string]: Docker.Container.Stats
  }

  // let intervalId: any
  const containerQueryStatsNum: ContainerQueryStatsNum = {}
  const containerRuningStats = ref<ContainerQueryStats>({})
  const containerQueryIds: string[] = []

  // const ms = useMessage()

  //   async function refreshContainerList() {
  //     containerRuningIds.value = []
  //     await getList<Common.ListResponse<Docker.Container.Info[]>>().then((res) => {
  //       containerList.value = res.data.list
  //       for (let i = 0; i < res.data.list.length; i++) {
  //         const element = res.data.list[i]
  //         containerRuningStats.value[element.id] = {
  //         } as Docker.Container.Stats
  //         if (element.state === DockerContainerInfoState.running) {
  //           containerRuningIds.value.push(element.id)
  //           console.log(element.id)
  //         }
  //       }
  //       refreshContainerListStats()
  //     }).catch((res) => {
  //       if (res.code === -1) {
  //         ms.error(`${t('apps.dockerManage.dockerError')}: ${res.msg}`, {
  //           closable: true,
  //           duration: 10000,
  //           keepAliveOnHover: true,
  //         })
  //       }
  //     })
  //   }

  async function refreshContainerListStats() {
    const ids = []
    for (let index = 0; index < containerQueryIds.length; index++) {
      const element = containerQueryIds[index]
      if (element)
        ids.push(element)
    }

    if (ids.length !== 0) {
      await getContainersStatsApi<ContainerQueryStats>(ids).then((res) => {
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
  }

  // 注册容器查询状态
  function registerContainerQueryStats(containerId: string) {
    if (!containerQueryStatsNum[containerId])
      containerQueryStatsNum[containerId] = 0

    containerQueryStatsNum[containerId]++

    if (!containerQueryIds.includes(containerId))
      containerQueryIds.push(containerId)
    if (!containerRuningStats.value[containerId])
      containerRuningStats.value[containerId] = {} as Docker.Container.Stats
  }

  // 注销容器查询状态
  function cancelContainerQueryStats(containerId: string) {
    if (containerQueryStatsNum[containerId]) {
      if (containerQueryStatsNum[containerId]) {
        containerQueryStatsNum[containerId]--

        if (containerQueryStatsNum[containerId] === 0) {
          // 解除查询状态
          delete containerQueryStatsNum[containerId]
          delete containerRuningStats.value[containerId]
          containerQueryIds.splice(containerQueryIds.indexOf(containerId), 1)
        //   console.log('解除查询状态', containerId)
        }
      }
    }
  }

  onMounted(() => {
    refreshContainerListStats()
    // console.log('docker状态定时器启动')
    setInterval(() => {
      if (containerQueryIds.length !== 0)
        refreshContainerListStats()
    //   console.log('计数', containerQueryStatsNum)
    //   console.log('ids', containerQueryIds)
    }, 5000) // 5000毫秒，即5秒
  })

  //   onUnmounted(() => {
  // console.log('docker状态定时器关闭')
  // clearInterval(intervalId)
  //   })

  return {
    containerRuningStats,
    registerContainerQueryStats,
    cancelContainerQueryStats,
  }
})
