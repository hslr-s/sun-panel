// import { ss } from '@/utils/storage'

// const LOCAL_NAME = 'appSetting'

export type Theme = 'light' | 'dark' | 'auto'

export type Language = 'zh-CN' | 'zh-TW' | 'en-US' | 'ko-KR'

export interface AdminState {
  siderCollapsed: boolean
  theme: Theme
  language: Language
}

export function defaultSetting(): AdminState {
  return { siderCollapsed: false, theme: 'light', language: 'zh-CN' }
}

// export function getLocalSetting(): AdminState {
//   const localSetting: AdminState | undefined = ss.get(LOCAL_NAME)
//   return { ...defaultSetting(), ...localSetting }
// }

// export function setLocalSetting(setting: AdminState): void {
//   ss.set(LOCAL_NAME, setting)
// }
