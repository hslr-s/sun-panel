import { post } from '@/utils/request'

export function edit<T>(req: Panel.ItemInfo) {
  return post<T>({
    url: '/panel/itemIconGroup/edit',
    data: req,
  })
}

export function getListByGroupId<T>() {
  return post<T>({
    url: '/panel/itemIconGroup/getListByGroupId',
  })
}

export function deletes<T>(ids: number[]) {
  return post<T>({
    url: '/panel/itemIconGroup/deletes',
    data: { ids },
  })
}
