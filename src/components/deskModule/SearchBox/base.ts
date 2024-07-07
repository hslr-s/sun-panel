import SvgSrcBing from '@/assets/search_engine_svg/bing.svg'
import SvgSrcGoogle from '@/assets/search_engine_svg/google.svg'

export interface SearchSaveDataType {
  currentSearchEngine: DeskModule.SearchBox.SearchEngine
  searchEngineList: DeskModule.SearchBox.SearchEngine[]
  newWindowOpen: boolean
}

export const defaultSearchEngineList: DeskModule.SearchBox.SearchEngine[] = [
  {
    iconSrc: SvgSrcGoogle,
    title: 'Google',
    url: 'https://www.google.com/search?q=',
    key: 'google',
  },
  {
    iconSrc: SvgSrcBing,
    title: 'Bing',
    url: 'https://www.bing.com/search?q=',
    key: 'bing',
  },
  // {
  //   iconSrc: SvgSrcBing,
  //   title: 'Bing',
  //   url: 'https://duckduckgo.com/?q=',
  // },
  // {
  //   iconSrc: SvgSrcBaidu,
  //   title: 'Baidu',
  //   url: 'https://www.baidu.com/s?wd=',
  // },
]

export const defaultSearchSaveData: SearchSaveDataType = {
  currentSearchEngine: defaultSearchEngineList[0],
  searchEngineList: [] || defaultSearchEngineList,
  newWindowOpen: false,
}
