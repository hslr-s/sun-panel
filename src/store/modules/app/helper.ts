import { ss } from '@/utils/storage'
import loginDefaultBackgroundUrl from '@/assets/loginDefaultBackground.webp'

const LOCAL_NAME = 'appSetting'

export type Theme = 'light' | 'dark' | 'auto'

export type Language = 'zh-CN' | 'en-US'

export interface AppState {
  theme: Theme
  language: Language
  siteSetting: GlobalSetting.SiteSetting
}

export function defaultSetting(): AppState {
  const lan = (navigator.language).toLowerCase()
  let language: Language = 'en-US'
  if (lan.includes('zh'))
    language = 'zh-CN'

  const siteSetting: GlobalSetting.SiteSetting = {
    title: 'Sun-Panel',
    favicoUrl: '',
    loginBackgroundUrl: loginDefaultBackgroundUrl,
    isEnableLoginCaptcha: false,
    isDisableItemCardCache: false,
  }

  return { theme: 'auto', language, siteSetting }
}

export function getLocalSetting(): AppState {
  const localSetting: AppState | undefined = ss.get(LOCAL_NAME)
  return { ...defaultSetting(), ...localSetting }
}

export function setLocalSetting(setting: AppState): void {
  ss.set(LOCAL_NAME, setting)
}

export function removeLocalState() {
  ss.remove(LOCAL_NAME)
}
