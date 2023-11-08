import { defineStore } from 'pinia'
import type { AdminState, Language, Theme } from './helper'
import { defaultSetting } from './helper'
import { store } from '@/store'

export const useAdminStore = defineStore('admin-store', {
  // state: (): AdminState => getLocalSetting(),
  state: (): AdminState => defaultSetting(),
  actions: {
    setSiderCollapsed(collapsed: boolean) {
      this.siderCollapsed = collapsed
      // this.recordState()
    },

    setTheme(theme: Theme) {
      this.theme = theme
      // this.recordState()
    },

    setLanguage(language: Language) {
      if (this.language !== language)
        this.language = language
        // this.recordState()
    },

    // recordState() {
    //   setLocalSetting(this.$state)
    // },
  },
})

export function useAdminStoreWithOut() {
  return useAdminStore(store)
}
