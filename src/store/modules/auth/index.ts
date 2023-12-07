import { defineStore } from 'pinia'
import { getStorage, removeToken as hRemoveToken, setStorage } from './helper'
import { VisitMode } from '@/enums/auth'

// interface SessionResponse {
//   auth: boolean
// }

export interface AuthState {
  token: string | null
  userInfo: User.Info | null
  // session: SessionResponse | null
  visitMode: VisitMode
}

const defaultState: AuthState = {
  token: null,
  userInfo: null,
  visitMode: VisitMode.VISIT_MODE_LOGIN,
}

export const useAuthStore = defineStore('auth-store', {
  state: (): AuthState => getStorage() || defaultState,

  actions: {
    setToken(token: string) {
      this.token = token
      this.saveStorage()
    },

    setUserInfo(userInfo: User.Info) {
      this.userInfo = userInfo
      this.saveStorage()
    },

    setVisitMode(visitMode: VisitMode) {
      this.visitMode = visitMode
      this.saveStorage()
    },

    saveStorage() {
      setStorage(this.$state)
    },

    removeToken() {
      this.$state = defaultState
      hRemoveToken()
    },
  },

})
