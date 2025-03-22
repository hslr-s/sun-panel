import { ss } from '@/utils/storage'

const LOCAL_NAME = 'accountManager'

export interface AccountInfo {
  name: string
  id: number
  username: string
  cToken: string
  loginStatus: boolean // 登录状态 false 为失效账号
}

export interface AccountManagerStore {
  accounts: AccountInfo[]
}

export function defaultStore(): AccountManagerStore {
  return ss.get(LOCAL_NAME) || { accounts: [] }
}

export function setLocalSetting(setting: AccountManagerStore): void {
  ss.set(LOCAL_NAME, setting)
}
