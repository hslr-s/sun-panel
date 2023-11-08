import type { AxiosProgressEvent, GenericAbortSignal } from 'axios'
import { post } from '@/utils/request'

export function fetchChatAPI<T = any>(
  prompt: string,
  options?: { conversationId?: string; parentMessageId?: string },
  signal?: GenericAbortSignal,
) {
  return post<T>({
    url: '/chat',
    data: { prompt, options },
    signal,
  })
}

export function fetchChatConfig<T = any>() {
  return post<T>({
    url: '/config',
  })
}

export function fetchChatAPIProcess<T = any>(
  params: {
    aiChatDialogId: number
    prompt: string
    options?: { conversationId?: string; parentMessageId?: string }
    signal?: GenericAbortSignal

    onDownloadProgress?: (progressEvent: AxiosProgressEvent) => void },
) {
  // const settingStore = useSettingStore()
  // const authStore = useAuthStore()
  const data: Record<string, any> = {
    prompt: params.prompt,
    options: params.options,
    aiChatDialogId: params.aiChatDialogId,
  }

  // if (authStore.isChatGPTAPI) {
  //   data = {
  //     ...data,
  //     systemMessage: settingStore.systemMessage,
  //     temperature: settingStore.temperature,
  //     top_p: settingStore.top_p,
  //   }
  // }
  return post<T>({
    url: '/chatGpt/chatCompletion',
    data,
    signal: params.signal,
    onDownloadProgress: params.onDownloadProgress,
  })
}

export function againFetchChatAPIProcess<T = any>(
  params: {
    aiChatDialogId: number
    prompt?: string
    options?: { conversationId?: string; parentMessageId?: string }
    signal?: GenericAbortSignal
    id?: number // 记录id

    onDownloadProgress?: (progressEvent: AxiosProgressEvent) => void },
) {
  const data: Record<string, any> = {
    prompt: params.prompt,
    options: params.options,
    aiChatDialogId: params.aiChatDialogId,
    id: params.id,
  }

  return post<T>({
    url: '/chatGpt/againChatCompletion',
    data,
    signal: params.signal,
    onDownloadProgress: params.onDownloadProgress,
  })
}

export function fetchSession<T>() {
  return post<T>({
    url: '/chatGpt/session',
  })
}

export function fetchVerify<T>(token: string) {
  return post<T>({
    url: '/verify',
    data: { token },
  })
}

// 获取对话列表
export function chatDialogGetList<T>(page: number, limit: number, keyword?: string) {
  return post<T>({
    url: '/aiChatDialog/getList',
    data: { page, limit, keyword },
  })
}

// 新建对话
export function chatDialogAdd<T>(title: string, aiRoleId: number) {
  return post<T>({
    url: '/aiChatDialog/add',
    data: { title, aiRoleId },
  })
}

// 修改
export function chatDialogUpdate<T>(aiChatDialogId: number, title: string) {
  return post<T>({
    url: '/aiChatDialog/update',
    data: { aiChatDialogId, title },
  })
}

// 删除对话
export function chatDialogDelete<T>(aiChatDialogId: number) {
  return post<T>({
    url: '/aiChatDialog/delete',
    data: { aiChatDialogId },
  })
}

export function chatDialogGetInfo<T>(aiChatDialogId: number) {
  return post<T>({
    url: '/aiChatDialog/getInfo',
    data: { aiChatDialogId },
  })
}

// 获取某对话框聊天记录
export function chatRecordGetList<T>(aiChatDialogId: number) {
  return post<T>({
    url: '/aiChatRecord/getList',
    data: { aiChatDialogId },
  })
}

// export function chatRecordAddOne<T>(data: ChatRecord.AddOneRequest) {
//   return post<T>({
//     url: '/aiChatRecord/addOne',
//     data,
//   })
// }

export function chatRecordDelete<T>(aiChatDialogId: number, recordId: number) {
  return post<T>({
    url: '/aiChatRecord/delete',
    data: { aiChatDialogId, id: recordId },
  })
}

export function chatRoleGetSystemList<T>(data: Common.ListRequest) {
  return post<T>({
    url: '/aiChatRole/getSystemList',
    data,
  })
}

export function chatRoleGetMyCreateList<T>(data: Common.ListRequest) {
  return post<T>({
    url: '/aiChatRole/getMyCreateList',
    data,
  })
}

export function chatRoleGetInfo<T>(aiRoleId: number) {
  return post<T>({
    url: '/aiChatRole/getInfo',
    data: { aiRoleId },
  })
}

// export function chatRoleEdit<T>(roleInfo: ChatRole.RoleInfo) {
//   return post<T>({
//     url: '/aiChatRole/edit',
//     data: roleInfo,
//   })
// }

export function chatRoleEditDeletes<T>(aiRoleIds: number[]) {
  return post<T>({
    url: '/aiChatRole/deletes',
    data: { aiRoleIds },
  })
}

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
