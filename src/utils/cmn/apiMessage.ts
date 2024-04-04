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

export function apiRespErrMsg(res: any): boolean {
  const appStore = useAppStore()
  const osTheme = useOsTheme()
  if (appStore.theme === 'auto')
    themeRef.value = osTheme.value as 'dark' | 'light'
  else
    themeRef.value = appStore.theme as 'dark' | 'light'

  const apiErrorCodeName = `apiErrorCode.${res.code}`
  const getI18nValue = t(apiErrorCodeName)
  if (apiErrorCodeName === getI18nValue) {
    message.error(t('common.unknownError'))
    console.error(res.msg)
    return false
  }
  else {
    message.error(t(`apiErrorCode.${res.code}`))
    return true
  }
}
