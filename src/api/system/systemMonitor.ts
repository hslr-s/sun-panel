import { post } from '@/utils/request'

export function getAll<T>() {
  return post<T>({
    url: '/system/monitor/getAll',
  })
}
