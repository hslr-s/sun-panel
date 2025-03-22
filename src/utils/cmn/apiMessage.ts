import type { ConfigProviderProps } from 'naive-ui'
import { createDiscreteApi, darkTheme, lightTheme, useOsTheme } from 'naive-ui'
import { computed, ref } from 'vue'
import { t } from '@/locales'
import { useAppStore, useAuthStore } from '@/store'
import { router } from '@/router'

let loginMessageShow = false

const themeRef = ref<'light' | 'dark'>('light')
const configProviderPropsRef = computed<ConfigProviderProps>(() => ({
  theme: themeRef.value === 'light' ? lightTheme : darkTheme,
}))
export const { message } = createDiscreteApi(['message'], { configProviderProps: configProviderPropsRef })

export function apiRespErrMsg(res: any) {
  const appStore = useAppStore()
  const osTheme = useOsTheme()
  const authStore = useAuthStore()

  if (appStore.theme === 'auto')
    themeRef.value = osTheme.value as 'dark' | 'light'
  else
    themeRef.value = appStore.theme as 'dark' | 'light'

  if (res.code === 1001) {
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
    return
  }
  else if (res.code === 1000) {
    router.push({ path: '/login' })
    authStore.removeToken()
    return
  }

  if (res.code) {
    const apiErrorCodeName = `apiErrorCode.${res.code}`
    const getI18nValue = t(apiErrorCodeName)
    if (apiErrorCodeName === getI18nValue && apiErrorCodeName !== undefined) {
      message.error(`${t('common.unknownError')}: ${res.msg}`)
      console.error(res.msg)
    }
    else {
      message.error(getI18nValue)
    }
  }
  else {
    // 其他错误（可能非API返回的错误）
    console.error(res.msg)
  }
}
