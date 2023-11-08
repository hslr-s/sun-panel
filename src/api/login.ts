import { post } from '@/utils/request'

// 下发重置密码的验证码到邮箱
export function sendResetPasswordVCode<T>(email: string, verification: Common.VerificationRequest) {
  return post<T>({
    url: '/login/sendResetPasswordVCode',
    data: { email, verification },
  })
}

// 下发重置密码的验证码到邮箱
export function resetPasswordByVCode<T>(data: Login.ResetPasswordByVCodeReqest) {
  return post<T>({
    url: '/login/resetPasswordByVCode',
    data,
  })
}
