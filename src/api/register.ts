import { post } from '@/utils/request'

export function sendRegisterVcode<T>(data: System.Register.SendRegisterVcodeRquest) {
  return post<T>({
    url: '/register/sendRegisterVcode',
    data,
  })
}

export function commit<T>(data: System.Register.SendRegisterVcodeRquest) {
  return post<T>({
    url: '/register/commit',
    data,
  })
}
