import { defineStore } from 'pinia'
import type { ModuleConfigState } from './helper'
import { getLocalState, setLocalState } from './helper'
import { getValueByName, save } from '@/api/system/moduleConfig'

export const useModuleConfig = defineStore('module-config-store', {
  state: (): ModuleConfigState => getLocalState(),
  actions: {

    // 保存
    // save(name: string, value: any) {
    //   const moduleName = `module-${name}`
    //   // 保存至网络
    //   console.log('保存模块配置', name, value)
    //   this.$state[moduleName] = value
    //   this.recordState()
    //   save(moduleName, value)
    // },

    // // 获取值
    // getValueByName<T>(name: string): T | null {
    //   const moduleName = `module-${name}`
    //   this.syncFromCloud(moduleName)
    //   if (this.$state[moduleName])
    //     return this.$state[moduleName]
    //   return null
    // },

    // 获取值
    async getValueByNameFromCloud<T>(name: string) {
      const moduleName = `module-${name}`
      return await getValueByName<T>(moduleName)
    },

    // 保存到网络
    saveToCloud(name: string, value: any) {
      const moduleName = `module-${name}`
      // 保存至网络
      save(moduleName, value)
    },

    // 从网络同步
    // syncFromCloud(moduleName: string) {
    //   getValueByName<any>(moduleName).then(({ code, data, msg }) => {
    //     if (code === 0)
    //       this.$state[moduleName] = data
    //   })
    // },

    recordState() {
      setLocalState(this.$state)
    },
  },
})
