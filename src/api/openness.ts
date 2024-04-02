import { get, post } from '@/utils/request'

export function getLoginConfig<T>() {
  return get<T>({
    url: '/openness/loginConfig',
  })
}

export function getSiteStting<T>() {
  return post<T>({
    url: '/panel/globalSetting/getSiteStting',
  })
}
