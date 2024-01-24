import { post } from '@/utils/request'

export function getAll<T>() {
  return post<T>({
    url: '/system/monitor/getAll',
  })
}

export function getCpuState<T>() {
  return post<T>({
    url: '/system/monitor/getCpuState',
  })
}

export function getDiskStateByPath<T>(path: string) {
  return post<T>({
    url: '/system/monitor/getDiskStateByPath',
    data: { path },
  })
}

export function getMemonyState<T>() {
  return post<T>({
    url: '/system/monitor/getMemonyState',
  })
}

export function getDiskMountpoints<T>() {
  return post<T>({
    url: '/system/monitor/getDiskMountpoints',
  })
}
