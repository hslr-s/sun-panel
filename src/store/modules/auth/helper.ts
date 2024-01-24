import type { AuthState } from './index'
import { ss } from '@/utils/storage'

const LOCAL_NAME = 'AUTH_TOKEN'

// export function getToken() {
//   return ss.get(LOCAL_NAME)
// }

// export function setToken(token: string) {
//   return ss.set(LOCAL_NAME, token)
// }

// export function setUserInfo(userInfo: User.Info) {
//   return ss.set(LOCAL_NAME, userInfo)
// }

// export function getUserInfo() {
//   return ss.get(LOCAL_NAME)
// }

export function setStorage(state: AuthState) {
  return ss.set(LOCAL_NAME, state)
}

export function getStorage() {
  return ss.get(LOCAL_NAME)
}

export function removeToken() {
  // ss.clear()
  return ss.remove(LOCAL_NAME)
}
