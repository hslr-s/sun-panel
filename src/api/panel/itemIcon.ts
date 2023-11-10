import { post } from '@/utils/request'

export function edit<T>(req: Panel.ItemInfo) {
  return post<T>({
    url: '/panel/itemIcon/edit',
    data: req,
  })
}

// export function getInfo<T>(id: number) {
//   return post<T>({
//     url: '/aiApplet/getInfo',
//     data: { id },
//   })
// }

export function getListByGroupId<T>() {
  return post<T>({
    url: '/panel/itemIcon/getListByGroupId',
  })
}

export function deletes<T>(ids: number[]) {
  return post<T>({
    url: '/panel/itemIcon/deletes',
    data: { ids },
  })
}
