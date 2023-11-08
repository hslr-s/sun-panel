import { post } from '@/utils/request'

// // 获取绘图的列表
// export function getMyDrawList<T>(req: Common.ListRequest) {
//   return post<T>({
//     url: '/aiDraw/getMyDrawList',
//     data: req,
//   })
// }

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

export function getSystemList<T>(data: Common.ListRequest) {
  return post<T>({
    url: '/aiApplet/getSystemList',
    data,
  })
}
