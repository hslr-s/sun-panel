import { post } from '@/utils/request'

// 登录相关

export function login<T>(data: Login.LoginReqest) {
  return post<T>({
    url: '/login',
    data,
  })
}

export function logout<T>() {
  return post<T>({
    url: '/logout',
  })
}

export function UserUpdateInfo<T>(headImage: string, name: string) {
  return post<T>({
    url: '/user/updateInfo',
    data: { headImage, name },
  })
}

export function AdminSystemSettingGetEmail<T>() {
  return post<T>({
    url: '/admin/systemSetting/getEmail',
  })
}

export function AdminSystemSettingGetWebsiteSetting<T>() {
  return post<T>({
    url: '/admin/systemSetting/getApplicationSetting',
  })
}

export function adminSystemSettingRoleManageGetSystemList<T>(data: Common.ListRequest) {
  return post<T>({
    url: '/admin/roleManage/getSystemList',
    data,
  })
}

export function adminSystemSettingRoleManageGetInfo<T>(aiRoleId: number) {
  return post<T>({
    url: '/admin/roleManage/getInfo',
    data: { aiRoleId },
  })
}

export function adminSystemSettingRoleManageDeletes<T>(aiRoleIds: number[]) {
  return post<T>({
    url: '/admin/roleManage/deletes',
    data: { aiRoleIds },
  })
}
