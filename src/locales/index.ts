import type { App } from 'vue'
import { createI18n } from 'vue-i18n'
import enUS from './en-US.json'
// import koKR from './ko-KR'
import zhCN from './zh-CN.json'
// import ruRU from './ru-RU'

const defaultLocale = 'zh-CN'

const i18n = createI18n({
  locale: defaultLocale,
  fallbackLocale: defaultLocale,
  allowComposition: true,
  messages: {
    'en-US': enUS,
    // 'ko-KR': koKR,
    'zh-CN': zhCN,
    // 'zh-TW': zhTW,
    // 'ru-RU': ruRU,
  },
})

export const t = i18n.global.t

// 避免循环依赖appstore(authstore)language此处暂时先使用any
// 后面有时间调整
export function setLocale(locale: any) {
  i18n.global.locale = locale
}

export function setupI18n(app: App) {
  app.use(i18n)
}

export default i18n
