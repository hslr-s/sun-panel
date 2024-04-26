import type { ConfigProviderProps } from 'naive-ui'
import { createDiscreteApi, darkTheme, lightTheme, useOsTheme } from 'naive-ui'
import { computed, ref } from 'vue'
import { t } from '@/locales'
import { useAppStore } from '@/store'

const themeRef = ref<'light' | 'dark'>('light')
const configProviderPropsRef = computed<ConfigProviderProps>(() => ({
  theme: themeRef.value === 'light' ? lightTheme : darkTheme,
}))
export const { message } = createDiscreteApi(['message'], { configProviderProps: configProviderPropsRef })

export function apiRespErrMsg(res: any) {
  const appStore = useAppStore()
  const osTheme = useOsTheme()
  if (appStore.theme === 'auto')
    themeRef.value = osTheme.value as 'dark' | 'light'
  else
    themeRef.value = appStore.theme as 'dark' | 'light'

  if (res.code) {
    const apiErrorCodeName = `apiErrorCode.${res.code}`
    const getI18nValue = t(apiErrorCodeName)
    console.log(apiErrorCodeName)
    if (apiErrorCodeName === getI18nValue && apiErrorCodeName !== undefined) {
      message.error(t('common.unknownError'))
      console.error(res.msg)
    }
    else {
      message.error(t(`apiErrorCode.${res.code}`))
    }
  }
  else {
    // 其他错误（可能非API返回的错误）
    console.error(res.msg)
  }
}
