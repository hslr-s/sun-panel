import { defineStore } from 'pinia'
import type { NoticeStore } from './helper'
import { getLocalSetting, setLocalSetting } from './helper'
import { store } from '@/store'

export const useNoticeStore = defineStore('notice-store', {
  state: (): NoticeStore => getLocalSetting(),
  actions: {

    // 设置已读
    setReadByGlobal(noticeId: number) {
      console.log('设置全局已读', noticeId)
      this.global.push(noticeId)
      this.recordState()
    },

    // 设置用户已读
    setReadByUsername(username: string, noticeId: number) {
      if (!this.username[username])
        this.username[username] = []
      this.username[username].push(noticeId)
      this.recordState()
    },

    // 判断是否为已读通知
    getReadByNoticeId(noticeId: number, username?: string): boolean {
      if (!this.global.includes(noticeId) && (!username || (!this.username[username] || !this.username[username].includes(noticeId))))
        return false
      return true
    },

    recordState() {
      setLocalSetting(this.$state)
    },
  },
})

export function useNoticeStoreWithOut() {
  return useNoticeStore(store)
}
