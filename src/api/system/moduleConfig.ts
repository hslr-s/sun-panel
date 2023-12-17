import { post } from '@/utils/request'

export function getValueByName<T>(name: string) {
  return post<T>({
    url: '/system/moduleConfig/getByName',
    data: { name },
  })
}

export function save<T>(name: string, value: any) {
  return post<T>({
    url: '/system/moduleConfig/save',
    data: {
      name,
      value,
    },
  })
}
