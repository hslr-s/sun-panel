import { defineStore } from 'pinia'
import type { AccountInfo, AccountManagerStore } from './helper'
import { defaultStore, setLocalSetting } from './helper'
import { getInfo } from '@/api/system/user'

export const useAccountManagerStore = defineStore('account-manager-store', {
  state: (): AccountManagerStore => defaultStore(),
  actions: {

    // 添加账号，已存在并且登录状态为失效状态 则先删除记录，再添加新的记录
    addAccount(account: AccountInfo) {
      const info = this.getAccountInfo(account.username)
      if (info && !info.loginStatus) {
        this.deleteAccount(account)
      }
      this.accounts.push(account)
      this.recordState()
      return true
    },

    isExistAccount(username: string) {
      return this.accounts.some(item => item.username === username)
    },

    // 删除账号
    deleteAccount(account: AccountInfo) {
      const index = this.accounts.findIndex(item => item.username === account.username || item.cToken === account.cToken)
      if (index !== -1)
        this.accounts.splice(index, 1)
      this.recordState()
    },

    getAccountInfo(username: string) {
      const index = this.accounts.findIndex(item => item.username === username)
      if (index !== -1)
        return this.accounts[index]
      return null
    },

    // 登录检测
    async loginCheck() {
      for (let i = 0; i < this.accounts.length; i++) {
        const element = this.accounts[i]
        if (element.loginStatus) {
          // 尝试连接检测账号的ctoken是否可用
          await getInfo<User.Info>(element.cToken).then(({ data }) => {
            this.accounts[i].name = data.name || ''
          }).catch(() => {
            this.accounts[i].loginStatus = false
          })
        }
      }

      this.recordState()
    },

    clearAccount() {
      this.accounts = []
      this.recordState()
    },

    recordState() {
      setLocalSetting(this.$state)
    },
  },
})
