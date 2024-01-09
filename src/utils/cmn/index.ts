import moment from 'moment'
import { h } from 'vue'
import type { NotificationReactive } from 'naive-ui'
import { NButton, createDiscreteApi } from 'naive-ui'
import { useAuthStore, useNoticeStore, useUserStore } from '@/store'
import { getAuthInfo } from '@/api/system/user'
import type { VisitMode } from '@/enums/auth'
import { getListByDisplayType as getListByDisplayTypeApi } from '@/api/notice'

const noticeStore = useNoticeStore()
const userStore = useUserStore()
const authStore = useAuthStore()

const { notification } = createDiscreteApi(['notification'])

/**
 * 生成指定时间格式
 * @param format 时间格式 默认：'YYYY-MM-DD HH:mm:ss'
 * @returns string
 */
export function buildTimeString(format?: string): string {
  if (!format)
    format = 'YYYY-MM-DD HH:mm:ss'

  return moment().format(format)
}

export function timeFormat(timeString?: string) {
  return moment(timeString).format('YYYY-MM-DD HH:mm:ss')
}

/**
 * 创建新的公告
 * @param timeString
 */
export function noticeCreate(info: Notice.NoticeInfo) {
  const option: any = {
    title: info.title,
    content: info.content,
    meta: info.createTime ? timeFormat(info.createTime) : '',
  }

  const btns: any = []

  let n: NotificationReactive
  // 链接按钮
  if (info.url !== '') {
    btns.push(
      h(
        NButton,
        {
          text: true,
          type: 'info',
          onClick: () => {
            window.open(info.url, '_blank')
            n.destroy()
          },
        },
        {
          default: () => '打开链接',
        },
      ),
    )
  }
  if (info.oneRead === 1) {
    btns.push(
      h(
        NButton,
        {
          text: true,
          type: 'primary',
          style: { marginLeft: '20px' },
          onClick: () => {
            if (info.id) {
              if (info.isLogin === 1 && userStore.userInfo.username) {
                noticeStore.setReadByUsername(userStore.userInfo.username, info.id)
                console.log('设置用户已读', info.id)
              }
              else {
                noticeStore.setReadByGlobal(info.id)
                console.log('设置全局已读', info.id)
              }
            }
            n.destroy()
          },
        },
        {
          default: () => '不再提醒',
        },
      ),
    )
  }
  option.action = () => btns
  n = notification.create(option)
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

  const { data } = await getAuthInfo<Req>()
  userStore.updateUserInfo({ headImage: data.user.headImage, name: data.user.name })
  authStore.setUserInfo(data.user)
  authStore.setVisitMode(data.visitMode)
}

export async function getNotice(displayType: number | number[]) {
  let param: number[]
  if (typeof displayType === 'number')
    param = [displayType]
  else
    param = displayType

  const { data } = await getListByDisplayTypeApi<Common.ListResponse<Notice.NoticeInfo[]>>(param)

  for (let i = 0; i < data.list.length; i++) {
    const element = data.list[i]
    if (element.id && !noticeStore.getReadByNoticeId(element.id, userStore.userInfo.username))
      noticeCreate(element)
  }
}

// 权限受限暂时不用
// export async function getFaviconUrl(url: string, extName = 'ico'): Promise<string | null> {
//   try {
//     // 获取网址的域名
//     const { protocol, host } = new URL(url)
//     const domain = `${protocol}//${host}`

//     // 构建 favicon URL
//     const faviconUrl = `${domain}/favicon.${extName}`

//     // 检查 favicon 是否存在，包含 CORS 头部
//     const response = await fetch(faviconUrl, { method: 'HEAD', mode: 'cors' })

//     // 如果请求成功，返回 favicon URL
//     if (response.ok) {
//       return faviconUrl
//     }
//     else {
//       console.log('Favicon not found.')
//       return null
//     }
//   }
//   catch (error) {
//     // 如果出现错误，返回 null，表示找不到 favicon
//     console.error('Error:', error)
//     return null
//   }
// }

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
