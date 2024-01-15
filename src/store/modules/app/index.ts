import { defineStore } from 'pinia'
import type { AppState, Language, Theme } from './helper'
import { defaultSetting, getLocalSetting, removeLocalState, setLocalSetting } from './helper'
import { store } from '@/store'
import { useTheme } from '@/hooks/useTheme'

export const useAppStore = defineStore('app-store', {
  state: (): AppState => getLocalSetting(),
  actions: {
    setSiderCollapsed(collapsed: boolean) {
      this.siderCollapsed = collapsed
      this.recordState()
    },

    setTheme(theme: Theme) {
      this.theme = theme
      this.recordState()
    },

    setLanguage(language: Language) {
      if (this.language !== language) {
        this.language = language
        this.recordState()
      }
    },

    getTheme() {
      const { theme } = useTheme()
      return theme

      // const appStore = useAppStore()
      // console.log('主题', appStore.theme)
      // if (appStore.theme === 'auto')
      //   return (window.matchMedia && window.matchMedia('(prefers-color-scheme: dark)').matches) ? darkTheme : lightTheme

      // else if (appStore.theme === 'light')
      //   return lightTheme

      // else
      //   return darkTheme
    },

    recordState() {
      setLocalSetting(this.$state)
    },

    removeToken() {
      this.$state = defaultSetting()
      removeLocalState()
    },
  },
})

export function useAppStoreWithOut() {
  return useAppStore(store)
}
