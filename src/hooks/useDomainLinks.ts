// import { enUS, koKR, zhCN, zhTW } from 'naive-ui'
import { useAppStore } from '@/store'

const linksDomain = 'http://links.sun-panel.top'

export function useDomainLinks() {
  const appStore = useAppStore()

  function getLangeuageSimple() {
    if (appStore.language === 'zh-CN')
      return 'zh_cn'

    return 'en'
  }

  // 获取链接带有语言后缀
  function linkWithLang(name: string): string {
    return `${linksDomain}/${name}_${getLangeuageSimple()}`
  }

  function linkWithLangParam(name: string): string {
    return `${linksDomain}/${name}?lang=${getLangeuageSimple()}`
  }

  function faqLinkWithLang(faqName: string): string {
    return `${linksDomain}/faq_${faqName}_${getLangeuageSimple()}`
  }

  return { linkWithLang, faqLinkWithLang, linkWithLangParam }
}
