import { defineStore } from 'pinia'
import { getToken, getUserInfo, removeToken, setToken } from './helper'
import { store } from '@/store'
import { VisitMode } from '@/enums/auth'

interface SessionResponse {
  auth: boolean
}

export interface AuthState {
  token: string | undefined
  userInfo: User.Info | undefined
  session: SessionResponse | null
  visitMode: VisitMode
}

export const useAuthStore = defineStore('auth-store', {
  state: (): AuthState => ({
    userInfo: getUserInfo(),
    token: getToken(),
    session: null,
    visitMode: VisitMode.VISIT_MODE_PUBLIC,
  }),

  getters: {
    // isChatGPTAPI(state): boolean {
    //   return state.session?.model === 'ChatGPTAPI'
    // },
  },

  actions: {
    // async getSession() {
    //   try {
    //     const { data } = await fetchSession<SessionResponse>()
    //     this.session = { ...data }
    //     return Promise.resolve(data)
    //   }
    //   catch (error) {
    //     return Promise.reject(error)
    //   }
    // },

    setToken(token: string) {
      this.token = token
      setToken(token)
    },

    setUserInfo(userInfo: User.Info) {
      this.userInfo = userInfo
      // this.setUserInfo(userInfo)
    },

    setVisitMode(mode: VisitMode) {
      this.visitMode = mode
    },

    // 清除所有的本地储存
    removeToken() {
      this.token = undefined
      removeToken()
    },
  },
})

export function useAuthStoreWithout() {
  return useAuthStore(store)
}
