import { post } from '@/utils/request'

export function edit<T>(param: User.Info) {
  let url = '/panel/users/create'
  if (param.id)
    url = '/panel/users/update'

  return post<T>({
    url,
    data: param,
  })
}

// 用户相关
export function getList<T>(param: AdminUserManage.GetListRequest) {
  return post<T>({
    url: '/panel/users/getList',
    data: param,
  })
}

export function deletes<T>(userIds: number[]) {
  return post<T>({
    url: '/panel/users/deletes',
    data: { userIds },
  })
}

export function getPublicVisitUser<T>() {
  return post<T>({
    url: '/panel/users/getPublicVisitUser',
  })
}

export function setPublicVisitUser<T>(userId: number | null) {
  return post<T>({
    url: '/panel/users/setPublicVisitUser',
    data: { userId },
  })
}
