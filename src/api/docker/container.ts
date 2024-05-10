import { post } from '@/utils/request'

export function getList<T>() {
  return post<T>({
    url: '/docker/container/getList',
  })
}

export function execSwitchAction<T>(param: Docker.Container.ExecSwitchActionReq) {
  return post<T>({
    url: '/docker/container/execSwitchActionByContainerID',
    data: param,
  })
}

export function getContainerStatus<T>(containerId: string) {
  return post<T>({
    url: '/docker/container/getContainerStatus',
    data: { id: containerId },
  })
}

export function getContainersStats<T>(containerIds: string[]) {
  return post<T>({
    url: '/docker/container/getContainersStats',
    data: containerIds,
  })
}
