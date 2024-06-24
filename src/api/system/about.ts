import { post } from '@/utils/request'

export function get<T>() {
  return post<T>({
    url: '/about',
  })
}

export function checkNewVersion<T>(isCheckBeta: boolean) {
  return post<T>({
    url: '/about/version/checkNewVersion',
    data: { isCheckBeta },
  })
}

export function setAutoCheckNewVersion<T>(data: About.AutoCheckNewVersionConfig) {
  return post<T>({
    url: '/about/version/setAutoCheckNewVersion',
    data,
  })
}

export function getAutoCheckNewVersion<T>(isCheckBeta: boolean) {
  return post<T>({
    url: '/about/version/getAutoCheckNewVersion',
  })
}
