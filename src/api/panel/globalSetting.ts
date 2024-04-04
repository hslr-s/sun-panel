import { post } from '@/utils/request'

export function saveSiteStting<T>(req: GlobalSetting.SiteSetting) {
  return post<T>({
    url: '/panel/globalSetting/saveSiteStting',
    data: req,
  })
}

export function getSiteStting<T>() {
  return post<T>({
    url: '/panel/globalSetting/getSiteStting',
  })
}

export function getCustomStylePath<T>() {
  return post<T>({
    url: '/panel/globalSetting/getCustomStylePath',
  })
}

export function getCustomJsAndCssCode<T>() {
  return post<T>({
    url: '/panel/globalSetting/getCustomJsAndCssCode',
  })
}

export function saveCustomJsAndCssCode<T>(data: GlobalSetting.GlobalSettingCustomJsAndCss) {
  return post<T>({
    url: '/panel/globalSetting/saveCustomJsAndCssCode',
    data,
  })
}
