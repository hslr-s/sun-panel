import { defineStore } from 'pinia'
import { defaultState, defaultStatePanelConfig, getLocalState, removeLocalState, setLocalState } from './helper'
import { router } from '@/router'
import type { PanelStateNetworkModeEnum } from '@/enums'
import { get as getUserConfig } from '@/api/panel/userConfig'
export const usePanelState = defineStore('panel', {
  state: (): Panel.State => getLocalState() || defaultState(),

  getters: {

  },

  actions: {
    setLeftSiderCollapsed(Collapsed: boolean) {
      this.leftSiderCollapsed = Collapsed
      // this.recordState()
    },

    setRightSiderCollapsed(Collapsed: boolean) {
      this.rightSiderCollapsed = Collapsed
      // this.recordState()
    },

    setNetworkMode(mode: PanelStateNetworkModeEnum) {
      this.networkMode = mode
      this.recordState()
    },

    // 获取云端（搭建的服务器）的面板配置
    updatePanelConfigByCloud() {
      getUserConfig<Panel.userConfig>().then((res) => {
        if (res.code === 0)
          this.panelConfig = { ...defaultStatePanelConfig(), ...res.data.panel }
        else
          this.resetPanelConfig() // 重置恢复默认
        this.recordState()
      })
    },

    resetPanelConfig() {
      this.panelConfig = defaultStatePanelConfig()
    },

    // 添加域名
    addDomainItem(item: Panel.DomainItem) {
      this.domainList.push(item)
      this.recordState()
    },

    // 删除域名
    removeDomainItem(id: string) {
      this.domainList = this.domainList.filter(item => item.id !== id)
      if (this.currentDomain === id)
        this.currentDomain = ''

      this.recordState()
    },

    // 设置当前域名
    setCurrentDomain(id: string) {
      this.currentDomain = id
      this.recordState()
    },

    // 更新域名列表
    updateDomainList(list: Panel.DomainItem[]) {
      this.domainList = list
      this.recordState()
    },

    // 更新域名
    updateDomainItem(id: string, updates: Partial<Panel.DomainItem>) {
      const index = this.domainList.findIndex(item => item.id === id)
      if (index !== -1) {
        this.domainList[index] = { ...this.domainList[index], ...updates }
        this.recordState()
      }
    },

    // 将临时域名转为正式域名
    convertTempDomain(id: string, domain: string, name?: string) {
      const index = this.domainList.findIndex(item => item.id === id)
      if (index !== -1) {
        this.domainList[index] = {
          ...this.domainList[index],
          domain,
          name,
          isTemp: false,
        }
        this.recordState()
      }
    },

    async reloadRoute(id?: number) {
      // this.recordState()
      await router.push({ name: 'AppletDialog', params: { aiAppletId: id } })
    },

    recordState() {
      setLocalState(this.$state)
    },

    removeState() {
      removeLocalState()
    },
  },
})
