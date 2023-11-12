import { post } from '@/utils/request'

export function getList<T>() {
  return post<T>({
    url: '/docker/container/getList',
  })
}
