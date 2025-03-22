import type { MessageReactive } from 'naive-ui'
import { usePanelState } from '@/store'
import { PanelStateNetworkModeEnum } from '@/enums'
import { message } from '@/utils/cmn/apiMessage'

const panelState = usePanelState()

export function pingUrl(url: string, attempts: number, timeout: number): Promise<number> {
  let totalLatency = 0
  let successfulAttempts = 0

  const attemptPing = (): Promise<number> => {
    return new Promise((resolve, reject) => {
      const start = Date.now()

      const fetchPromise = fetch(url, { method: 'HEAD', mode: 'no-cors' })
        .then(() => {
          const end = Date.now()
          resolve(end - start)
        })
        .catch((err) => {
          reject(err)
        })

      const timeoutPromise = new Promise<number>((resolve, reject) =>
        setTimeout(() => reject(new Error('Timeout')), timeout),
      )

      Promise.race([fetchPromise, timeoutPromise])
        .then(latency => resolve(latency as number))
        .catch(err => reject(err))
    })
  }

  return new Promise((resolve, reject) => {
    const pingAttempt = (i: number) => {
      if (i < attempts) {
        attemptPing()
          .then((latency) => {
            totalLatency += latency
            successfulAttempts++
            pingAttempt(i + 1)
          })
          .catch(() => {
            // 失败的情况下继续尝试
            pingAttempt(i + 1)
          })
      }
      else {
        if (successfulAttempts > 0) {
          const averageLatency = totalLatency / successfulAttempts
          resolve(averageLatency)
        }
        else {
          reject(new Error('All ping attempts failed'))
        }
      }
    }

    pingAttempt(0)
  })
}

// 选择最优地址
async function selectBestUrl(lanUrl: string, url: string, timeout?: number) {
  // 开启消息类型的 loading
  let messageReactive: MessageReactive | null = null
  messageReactive = message.loading('loading...', {
    duration: 0,
  })

  await pingUrl(lanUrl, 2, timeout || 200).then(() => {
    url = lanUrl
  }).catch((err) => {
    // 无法ping通局域网,使用默认的url
    console.error(lanUrl, err)
  })

  // 关闭loading
  if (messageReactive) {
    messageReactive.destroy()
    messageReactive = null
  }
  return url
}

// 获取项目地址，如果没有局域网地址，将直接使用默认地址。
// 如果局域网地址可用并且开启了自动识别，优先尝试ping局域网地址，无法ping通将使用默认地址。
export async function getItemUrl(itemInfo: Panel.ItemInfo, timeout?: number) {
  let jumpUrl = ''
  if (itemInfo?.lanUrl === '') {
    jumpUrl = itemInfo.url
  }
  else {
    if (panelState.panelConfig.netModeChangeAuto) {
      await selectBestUrl(itemInfo.lanUrl || '', itemInfo.url, timeout).then((url) => {
        jumpUrl = url
      })
    }
    else {
      jumpUrl = (panelState.networkMode === PanelStateNetworkModeEnum.lan ? itemInfo.lanUrl : itemInfo.url) as string
    }
  }

  return jumpUrl
}
