import type { AxiosProgressEvent, AxiosResponse, GenericAbortSignal } from 'axios'
import { createDiscreteApi } from 'naive-ui'
import request from './axios'
import { t } from '@/locales'
import { useAppStore, useAuthStore } from '@/store'
import { router } from '@/router'

const { message } = createDiscreteApi(['message'])
let loginMessageShow = false
export interface HttpOption {
  url: string
  data?: any
  method?: string
  headers?: any
  onDownloadProgress?: (progressEvent: AxiosProgressEvent) => void
  signal?: GenericAbortSignal
  beforeRequest?: () => void
  afterRequest?: () => void
}

export interface Response<T = any> {
  data: T
  // message: string | null
  // status: string
  msg: string
  code: number
}

function http<T = any>(
  { url, data, method, headers, onDownloadProgress, signal, beforeRequest, afterRequest }: HttpOption,
) {
  const authStore = useAuthStore()
  const appStore = useAppStore()
  const successHandler = (res: AxiosResponse<Response<T>>) => {
    if (res.data.code === 0 || typeof res.data === 'string')
      return res.data

    if (res.data.code === 1001) {
      // 避免重复弹窗
      if (loginMessageShow === false) {
        loginMessageShow = true
        message.warning(t('api.loginExpires'), {
        // message.warning('登录过期', {
          onLeave() {
            loginMessageShow = false
          },
        })
      }

      router.push({ path: '/login' })
      authStore.removeToken()
      return res.data
    }

    if (res.data.code === 1000) {
      router.push({ path: '/login' })
      authStore.removeToken()
      return res.data
    }

    if (res.data.code === 1005) {
      message.warning(res.data.msg)
      return res.data
    }

    if (res.data.code === -1) {
      message.warning(res.data.msg)
      // router.push({ path: '/login' })
      // authStore.removeToken()
      return res.data
    }

    // 验证码相关错误
    if (res.data.code > 1100 && res.data.code < 1200)
      return res.data

    return Promise.reject(res.data)
  }

  const failHandler = (error: Response<Error>) => {
    afterRequest?.()
    // message.error('网络错误，请稍后重试', {
    //   duration: 50000,
    //   closable: true,
    // })
    throw new Error(error?.msg || 'Error')
  }

  beforeRequest?.()

  method = method || 'GET'

  const params = Object.assign(typeof data === 'function' ? data() : data ?? {}, {})
  if (!headers)
    headers = {}

  headers.token = authStore.token
  headers.lang = appStore.language
  return method === 'GET'
    ? request.get(url, { params, signal, onDownloadProgress }).then(successHandler, failHandler)
    : request.post(url, params, { headers, signal, onDownloadProgress }).then(successHandler, failHandler)
}

export function get<T = any>(
  { url, data, method = 'GET', onDownloadProgress, signal, beforeRequest, afterRequest }: HttpOption,
): Promise<Response<T>> {
  return http<T>({
    url,
    method,
    data,
    onDownloadProgress,
    signal,
    beforeRequest,
    afterRequest,
  })
}

export function post<T = any>(
  { url, data, method = 'POST', headers, onDownloadProgress, signal, beforeRequest, afterRequest }: HttpOption,
): Promise<Response<T>> {
  return http<T>({
    url,
    method,
    data,
    headers,
    onDownloadProgress,
    signal,
    beforeRequest,
    afterRequest,
  })
}

export default post
