import { post } from '@/utils/request'

export function get<T>() {
  return post<T>({
    url: '/about',
  })
}
