<script setup lang="ts">
import { defineEmits, onMounted, ref } from 'vue'
import { NAvatar, NCheckbox } from 'naive-ui'
import { SvgIcon } from '@/components/common'
import { useModuleConfig } from '@/store/modules'
import { useAuthStore } from '@/store'
import { VisitMode } from '@/enums/auth'

import SvgSrcBaidu from '@/assets/search_engine_svg/baidu.svg'
import SvgSrcBing from '@/assets/search_engine_svg/bing.svg'
import SvgSrcGoogle from '@/assets/search_engine_svg/google.svg'

withDefaults(defineProps<{
  background?: string
  textColor?: string
}>(), {
  background: '#2a2a2a6b',
  textColor: 'white',
})

const emits = defineEmits(['itemSearch'])

interface State {
  currentSearchEngine: DeskModule.SearchBox.SearchEngine
  searchEngineList: DeskModule.SearchBox.SearchEngine[]
  newWindowOpen: boolean
}

const moduleConfigName = 'deskModuleSearchBox'
const moduleConfig = useModuleConfig()
const authStore = useAuthStore()
const searchTerm = ref('')
const isFocused = ref(false)
const searchSelectListShow = ref(false)
const defaultSearchEngineList = ref<DeskModule.SearchBox.SearchEngine[]>([
  {
    iconSrc: SvgSrcGoogle,
    title: 'Google',
    url: 'https://www.google.com/search?q=%s',
  },
  {
    iconSrc: SvgSrcBaidu,
    title: 'Baidu',
    url: 'https://www.baidu.com/s?wd=%s',
  },
  {
    iconSrc: SvgSrcBing,
    title: 'Bing',
    url: 'https://www.bing.com/search?q=%s',
  },
])

const defaultState: State = {
  currentSearchEngine: defaultSearchEngineList.value[0],
  searchEngineList: [] || defaultSearchEngineList,
  newWindowOpen: false,
}

const state = ref<State>({ ...defaultState })

const onFocus = (): void => {
  isFocused.value = true
}

const onBlur = (): void => {
  isFocused.value = false
}

function handleEngineClick() {
  // 访客模式不允许修改
  if (authStore.visitMode === VisitMode.VISIT_MODE_PUBLIC)
    return
  searchSelectListShow.value = !searchSelectListShow.value
}

function handleEngineUpdate(engine: DeskModule.SearchBox.SearchEngine) {
  state.value.currentSearchEngine = engine
  moduleConfig.saveToCloud(moduleConfigName, state.value)
  searchSelectListShow.value = false
}

function handleSearchClick() {
  const url = state.value.currentSearchEngine.url
  const keyword = searchTerm
  // 如果网址中存在 %s，则直接替换为关键字
  const fullUrl = replaceOrAppendKeywordToUrl(url, keyword.value)
  handleClearSearchTerm()
  if (state.value.newWindowOpen)
    window.open(fullUrl)
  else
    window.location.href = fullUrl
}

function replaceOrAppendKeywordToUrl(url: string, keyword: string) {
  // 如果网址中存在 %s，则直接替换为关键字
  if (url.includes('%s'))
    return url.replace('%s', encodeURIComponent(keyword))

  // 如果网址中不存在 %s，则将关键字追加到末尾
  return url + (keyword ? `${encodeURIComponent(keyword)}` : '')
}

const handleItemSearch = () => {
  emits('itemSearch', searchTerm.value)
}

function handleClearSearchTerm() {
  searchTerm.value = ''
  emits('itemSearch', searchTerm.value)
}

onMounted(() => {
  moduleConfig.getValueByNameFromCloud<State>('deskModuleSearchBox').then(({ code, data }) => {
    if (code === 0)
      state.value = data || defaultState
    else
      state.value = defaultState
  })
})
</script>

<template>
  <div class="search-box w-full" @keydown.enter="handleSearchClick" @keydown.esc="handleClearSearchTerm">
    <div class="search-container flex rounded-2xl items-center justify-center text-white w-full" :style="{ background, color: textColor }" :class="{ focused: isFocused }">
      <div class="search-box-btn-engine w-[40px] flex justify-center cursor-pointer" @click="handleEngineClick">
        <NAvatar :src="state.currentSearchEngine.iconSrc" style="background-color: transparent;" :size="20" />
      </div>

      <input v-model="searchTerm" :placeholder="$t('deskModule.searchBox.inputPlaceholder')" @focus="onFocus" @blur="onBlur" @input="handleItemSearch">

      <div v-if="searchTerm !== ''" class="search-box-btn-clear w-[25px] mr-[10px] flex justify-center cursor-pointer" @click="handleClearSearchTerm">
        <SvgIcon style="width: 20px;height: 20px;" icon="line-md:close-small" />
      </div>
      <div class="search-box-btn-search w-[25px] flex justify-center cursor-pointer" @click="handleSearchClick">
        <SvgIcon style="width: 20px;height: 20px;" icon="iconamoon:search-fill" />
      </div>
    </div>

    <!-- 搜索引擎选择 -->
    <div v-if="searchSelectListShow" class="w-full mt-[10px] rounded-xl p-[10px]" :style="{ background }">
      <div class="flex items-center">
        <div class="flex items-center">
          <div
            v-for="item, index in defaultSearchEngineList"
            :key="index"
            :title="item.title"
            class="w-[40px] h-[40px] mr-[10px]  cursor-pointer bg-[#ffffff] flex items-center justify-center rounded-xl"
            @click="handleEngineUpdate(item)"
          >
            <NAvatar :src="item.iconSrc" style="background-color: transparent;" :size="20" />
          </div>
        <!-- <div class="w-[40px] h-[40px] ml-[10px] flex justify-center items-center cursor-pointer" @click="handleEngineClick">
          <NAvatar style="background-color: transparent;" :size="30">
            <SvgIcon icon="lets-icons:setting-alt-fill" style="font-size: 20px;" />
          </NAvatar>
        </div> -->
        </div>
      </div>

      <div class="mt-[10px]">
        <NCheckbox v-model:checked="state.newWindowOpen" @update-checked="moduleConfig.saveToCloud(moduleConfigName, state)">
          <span :style="{ color: textColor }">
            {{ $t('deskModule.searchBox.openWithNewOpen') }}
          </span>
        </NCheckbox>
      </div>
    </div>
  </div>
</template>

<style scoped>
.search-container {
  border: 1px solid #ccc;
  transition: box-shadow 0.5s,backdrop-filter 0.5s;
  padding: 2px 10px;
  backdrop-filter:blur(2px)
}

.focused, .search-container:hover {
  box-shadow: 0px 0px 30px -5px rgba(41, 41, 41, 0.45);
  -webkit-box-shadow: 0px 0px 30px -5px rgba(0, 0, 0, 0.45);
  -moz-box-shadow: 0px 0px 30px -5px rgba(0, 0, 0, 0.45);
  backdrop-filter:blur(5px)
}

.before {
  left: 10px;
}

.after {
  right: 10px;
}

input {
  background-color: transparent;
  box-sizing: border-box;
  width: 100%;
  height: 40px;
  padding: 10px 5px;
  border: none;
  outline: none;
  font-size: 17px;
}
</style>
