import { defineStore } from 'pinia'
import type { AccountInfo, AccountManagerStore } from './helper'
import { defaultStore, setLocalSetting } from './helper'

export const useAccountManagerStore = defineStore('account-manager-store', {
  state: (): AccountManagerStore => defaultStore(),
  actions: {

    // 添加账号
    addAccount(account: AccountInfo) {
      this.accounts.push(account)
      this.recordState()
    },

    // 删除账号
    deleteAccount(account: AccountInfo) {
      const index = this.accounts.findIndex(item => item.cToken === account.cToken)
      if (index !== -1)
        this.accounts.splice(index, 1)
      this.recordState()
    },

    // 登录检测
    async loginCheck() {
      for (let i = 0; i < this.accounts.length; i++) {
        const element = this.accounts[i]
        if (element.loginStatus) {
          // 尝试连接检测账号的ctoken是否可用

        }
      }

      this.recordState()
    },

    // 更新账号

    // 设置已读
    // setReadByGlobal(noticeId: number) {
    //   console.log('设置全局已读', noticeId)
    //   this.global.push(noticeId)
    //   this.recordState()
    // },

    // // 设置用户已读
    // setReadByUsername(username: string, noticeId: number) {
    //   if (!this.username[username])
    //     this.username[username] = []
    //   this.username[username].push(noticeId)
    //   this.recordState()
    // },

    // // 判断是否为已读通知
    // getReadByNoticeId(noticeId: number, username?: string): boolean {
    //   if (!this.global.includes(noticeId) && (!username || (!this.username[username] || !this.username[username].includes(noticeId))))
    //     return false
    //   return true
    // },

    recordState() {
      setLocalSetting(this.$state)
    },
  },
})
