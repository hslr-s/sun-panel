import { post } from '@/utils/request'

export function refreshToken<T>() {
  return post<T>({
    url: '/openApi/refreshToken',
  })
}

export function deleteToken<T>() {
  return post<T>({
    url: '/openApi/deleteToken',
  })
}

export function getTokenStatusInfo<T>() {
  return post<T>({
    url: '/openApi/getTokenStatusInfo',
  })
}
