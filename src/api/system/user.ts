import { post } from '@/utils/request'

// export function getInfo<T>() {
//   return post<T>({
//     url: '/user/getInfo',
//   })
// }

export function getAuthInfo<T>() {
  return post<T>({
    url: '/user/getAuthInfo',
  })
}

export function getReferralCode<T>() {
  return post<T>({
    url: '/user/getReferralCode',
  })
}
