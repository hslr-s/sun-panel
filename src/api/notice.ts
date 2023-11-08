import { post } from '@/utils/request'

// 下发重置密码的验证码到邮箱
export function getListByDisplayType<T>(displayType: number[]) {
  return post<T>({
    url: '/notice/getListByDisplayType',
    data: { displayType },
  })
}
