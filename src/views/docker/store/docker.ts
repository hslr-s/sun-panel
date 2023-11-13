import { defineStore } from 'pinia'
import { ref } from 'vue'
import { useMessage } from 'naive-ui'
import { getList } from '@/api/docker/container'

export const useDockerStore = defineStore('docker-store', () => {
  const containerList = ref<Docker.Container.Info[]>([])
  const actionMenuBar = ref('container')

  const ms = useMessage()

  async function refreshContainerList() {
    await getList<Common.ListResponse<Docker.Container.Info[]>>().then((res) => {
      if (res.code === 0)
        containerList.value = res.data.list
      else
        ms.error(res.msg)
    })
  }

  return {
    actionMenuBar,
    refreshContainerList,
    containerList,
  }
})
