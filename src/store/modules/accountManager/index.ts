import { defineStore } from 'pinia'
import type { AccountInfo, AccountManagerStore } from './helper'
import { defaultStore, setLocalSetting } from './helper'
import { getInfo } from '@/api/system/user'

export const useAccountManagerStore = defineStore('account-manager-store', {
  state: (): AccountManagerStore => defaultStore(),
  actions: {

    // 添加账号
    addAccount(account: AccountInfo) {
      if (this.isExistAccount(account.username)) {
        return false
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
