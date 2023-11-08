import { post } from '@/utils/request'

export function getInfo<T>() {
  return post<T>({
    url: '/user/getInfo',
  })
}

export function getReferralCode<T>() {
  return post<T>({
    url: '/user/getReferralCode',
  })
}
