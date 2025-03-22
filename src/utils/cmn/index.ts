import moment from 'moment'
import { apiRespErrMsg as apiRespErrMsgApiMessage } from './apiMessage'
import { useAuthStore, useUserStore } from '@/store'
import { getAuthInfo } from '@/api/system/user'
import type { VisitMode } from '@/enums/auth'

const userStore = useUserStore()
const authStore = useAuthStore()

/**
 * 生成指定时间格式
 * @param format 时间格式 默认：'YYYY-MM-DD HH:mm:ss'
 * @returns string
 */
export function buildTimeString(timeString?: string, format?: string): string {
  if (!format)
    format = 'YYYY-MM-DD HH:mm:ss'

  return moment(timeString).format(format)
}

export function timeFormat(timeString?: string) {
  return moment(timeString).format('YYYY-MM-DD HH:mm:ss')
}
export function inArray<T>(array: T[], item: T): boolean {
  return array.includes(item)
}

export function setTitle(titile: string) {
  document.title = titile
}

export function getTitle(titile: string) {
  document.title = titile
}

//
export async function updateLocalUserInfo() {
  interface Req {
    user: User.Info
    visitMode: VisitMode
  }

  await getAuthInfo<Req>().then(({ data }) => {
    if (data.user) {
      userStore.updateUserInfo({ headImage: data?.user.headImage, name: data?.user.name })
      authStore.setUserInfo(data.user)
      authStore.setVisitMode(data.visitMode)
    }
  }).catch((res) => {
    apiRespErrMsgApiMessage(res)
  })
}

export function getFaviconUrl(url: string): string {
  // 获取网址的域名
  const { protocol, host } = new URL(url)
  const domain = `${protocol}//${host}`
  // 构建 favicon URL
  return `${domain}/favicon.ico`
}

/**
 * @description: 获取随机码
 * @param {number} size
 * @param {array} seed ["a","b"m"c]
 * @return {string}
 */
export function randomCode(size: number, seed?: Array<string>) {
  seed = seed || ['A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z',
    'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'm', 'n', 'p', 'Q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z',
    '2', '3', '4', '5', '6', '7', '8', '9',
  ]// 数组
  const seedlength = seed.length// 数组长度
  let createPassword = ''
  for (let i = 0; i < size; i++) {
    const j = Math.floor(Math.random() * seedlength)
    createPassword += seed[j]
  }
  return createPassword
}

// 复制文字到剪切板
export async function copyToClipboard(text: string): Promise<boolean> {
  if (navigator.clipboard) {
    // 使用 Clipboard API
    try {
      await navigator.clipboard.writeText(text)
      return true
    }
    catch (err) {
      console.error('copy fail', err)
      return false
    }
  }
  else {
    // 兼容旧版浏览器
    const textArea = document.createElement('textarea')
    textArea.value = text
    document.body.appendChild(textArea)
    textArea.select()

    try {
      document.execCommand('copy')
      return true
    }
    catch (err) {
      console.error('copy fail', err)
      return false
    }
    finally {
      document.body.removeChild(textArea)
    }
  }
}

export function bytesToSize(bytes: number) {
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB', 'PB', 'EB']
  if (bytes === 0)
    return '0B'
  const i = parseInt(String(Math.floor(Math.log(bytes) / Math.log(1024))))
  return `${(bytes / 1024 ** i).toFixed(1)} ${sizes[i]}`
}

export const apiRespErrMsg = apiRespErrMsgApiMessage
