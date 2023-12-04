import { defineStore } from 'pinia'
import { ref, toRefs } from 'vue'
import { getStorage, removeToken as hRemoveToken, setStorage } from './helper'
import { VisitMode } from '@/enums/auth'

// interface SessionResponse {
//   auth: boolean
// }

export interface AuthState {
  token: string | undefined
  userInfo: User.Info | undefined
  // session: SessionResponse | null
  visitMode: VisitMode
}

const defaultState: AuthState = {
  token: undefined,
  userInfo: undefined,
  visitMode: VisitMode.VISIT_MODE_LOGIN,
}

export const useAuthStore = defineStore('auth-store', () => {
  const state = ref<AuthState>(getStorage() || defaultState)

  function setToken(token: string) {
    state.value.token = token
    saveStorage()
  }

  function setUserInfo(userInfo: User.Info) {
    state.value.userInfo = userInfo
    saveStorage()
  }

  function setVisitMode(visitMode: VisitMode) {
    state.value.visitMode = visitMode
    saveStorage()
  }

  function saveStorage() {
    setStorage(state.value)
  }

  function removeToken() {
    state.value = defaultState
    hRemoveToken()
  }
  const stateRefs = toRefs(state.value)
  return {
    state: state.value,
    ...stateRefs,
    setToken,
    setUserInfo,
    setVisitMode,
    removeToken,
  }
})
