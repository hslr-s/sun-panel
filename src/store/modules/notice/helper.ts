import { ss } from '@/utils/storage'

const LOCAL_NAME = 'noticeStore'

export interface NoticeStore {
  global: number[]
  username: { [key: string]: number[] }
}

export function defaultSetting(): NoticeStore {
  return { global: [], username: {} }
}

export function getLocalSetting(): NoticeStore {
  const localSetting: NoticeStore | undefined = ss.get(LOCAL_NAME)
  return { ...defaultSetting(), ...localSetting }
}

export function setLocalSetting(setting: NoticeStore): void {
  ss.set(LOCAL_NAME, setting)
}
