import { post } from '@/utils/request'

export function login<T>(username: string, password: string) {
  return post<T>({
    url: '/panel/proAuth/login',
    data: {
      username,
      password,
    },
  })
}

export function logout<T>() {
  return post<T>({
    url: '/panel/proAuth/logout',
  })
}

export function getStatus<T>() {
  return post<T>({
    url: '/panel/proAuth/getStatus',
  })
}

export function refreshStatus<T>() {
  return post<T>({
    url: '/panel/proAuth/refreshStatus',
  })
}

export function setHideProBadgeStatus<T>(status: boolean) {
  return post<T>({
    url: '/panel/proAuth/setHideProBadgeStatus',
    data: { status },
  })
}
