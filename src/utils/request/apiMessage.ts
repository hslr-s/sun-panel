// 此文件待废弃 日期：2025-02-25
//
import type { ConfigProviderProps } from 'naive-ui'
import { createDiscreteApi, darkTheme, lightTheme } from 'naive-ui'
import { computed, ref } from 'vue'
import type { Response } from './index'
import { apiRespErrMsg as newApiRespErrMsg } from '@/utils/cmn/apiMessage'

const themeRef = ref<'light' | 'dark'>('light')
const configProviderPropsRef = computed<ConfigProviderProps>(() => ({
  theme: themeRef.value === 'light' ? lightTheme : darkTheme,
}))
export const { message } = createDiscreteApi(['message'], { configProviderProps: configProviderPropsRef })

/**
 * @deprecated 该函数即将弃用，请使用 utils/cmn/apiMessage.ts 中的 apiRespErrMsg
 */
export function apiRespErrMsg(res: Response): boolean {
  return newApiRespErrMsg(res)
}
