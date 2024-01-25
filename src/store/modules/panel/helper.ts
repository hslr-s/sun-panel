import { ss } from '@/utils/storage'
import { PanelPanelConfigStyleEnum, PanelStateNetworkModeEnum } from '@/enums'
import defaultBackground from '@/assets/defaultBackground.webp'
const LOCAL_NAME = 'panelStorage'

const defaultFooterHtml = '<div class="flex justify-center text-slate-300" style="margin-top:100px">Powered By <a href="https://github.com/hslr-s/sun-panel" target="_blank" class="ml-[5px]">Sun-Panel</a></div>'

export function defaultStatePanelConfig(): Panel.panelConfig {
  return {
    backgroundImageSrc: defaultBackground,
    backgroundBlur: 0,
    backgroundMaskNumber: 0,
    iconStyle: PanelPanelConfigStyleEnum.icon,
    iconTextColor: '#ffffff',
    iconTextInfoHideDescription: false,
    iconTextIconHideTitle: false,
    logoText: 'Sun-Panel',
    logoImageSrc: '',
    clockShowSecond: false,
    searchBoxShow: false,
    searchBoxSearchIcon: false,
    marginBottom: 10,
    marginTop: 10,
    maxWidth: 1200,
    maxWidthUnit: 'px',
    marginX: 5,
    footerHtml: defaultFooterHtml,
    systemMonitorShow: false,
    systemMonitorShowTitle: true,
    systemMonitorPublicVisitModeShow: false,
    netModeChangeButtonShow: true,

  }
}

export function defaultState(): Panel.State {
  return {
    rightSiderCollapsed: false,
    leftSiderCollapsed: false,
    networkMode: PanelStateNetworkModeEnum.wan,
    panelConfig: { ...defaultStatePanelConfig() },
  }
}

export function getLocalState(): Panel.State {
  const localState = ss.get(LOCAL_NAME)
  return { ...defaultState(), ...localState }
}

export function setLocalState(state: Panel.State) {
  ss.set(LOCAL_NAME, state)
}

export function removeLocalState() {
  ss.remove(LOCAL_NAME)
}
