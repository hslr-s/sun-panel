<script setup lang="ts">
import { defineEmits, onBeforeUnmount, onMounted, ref } from 'vue'
import { NAvatar, NCheckbox, NModal } from 'naive-ui'
import SearchBoxSetting from './SearchBoxSetting.vue'
import type { SearchSaveDataType } from './base'
import { defaultSearchEngineList, defaultSearchSaveData } from './base'
import { SvgIcon } from '@/components/common'
import { useModuleConfig } from '@/store/modules'
import { useAuthStore } from '@/store'
import { VisitMode } from '@/enums/auth'
import { t } from '@/locales'
import { message } from '@/utils/cmn/apiMessage'

withDefaults(defineProps<{
  background?: string
  textColor?: string
  borderColor?: string
}>(), {
  background: '#2a2a2a6b',
  textColor: 'white',
})

const emits = defineEmits(['itemSearch'])

interface State extends SearchSaveDataType {
}

const moduleConfigName = 'deskModuleSearchBox'
const moduleConfig = useModuleConfig()
const authStore = useAuthStore()
const searchTerm = ref('')
const isFocused = ref(false)
const searchInputRef = ref<HTMLInputElement | null>(null) // 输入框的 ref
const searchSelectListShow = ref(false)
const searchBoxSettingShow = ref(false)

const defaultState: State = defaultSearchSaveData

const state = ref<State>({ ...defaultState })

const selectEngineList = ref<DeskModule.SearchBox.SearchEngine[]>([])

const onFocus = (): void => {
  isFocused.value = true
}

const onBlur = (): void => {
  isFocused.value = false
}

function handleEngineClick() {
  // 访客模式不允许修改
  // if (authStore.visitMode === VisitMode.VISIT_MODE_PUBLIC) {
  //   message.warning(t('common.publicModeRestricted'))
  //   return
  // }

  searchSelectListShow.value = !searchSelectListShow.value
}

function handleEngineUpdate(engine: DeskModule.SearchBox.SearchEngine) {
  state.value.currentSearchEngine = engine
  searchSelectListShow.value = false

  // 非登录模式不保存
  if (authStore.visitMode === VisitMode.VISIT_MODE_PUBLIC) {
    return
  }
  moduleConfig.saveToCloud(moduleConfigName, state.value)
}

function handleOnKeyDownEnter(e: any) {
  // 修复Mac safari浏览器中文输入法回车触发
  if (e.keyCode === 229)
    return
  if (e.key === 'Enter' && !e.isComposing) {
    handleSearchClick()
  }
}

function handleSearchClick() {
  const url = state.value.currentSearchEngine.url
  const keyword = searchTerm
  // 如果网址中存在 %s，则直接替换为关键字
  const fullUrl = replaceOrAppendKeywordToUrl(url, keyword.value)
  handleClearSearchTerm()
  if (state.value.newWindowOpen) {
    window.open(fullUrl)
  }
  else {
    if (window.self !== window.top)
      window.parent.location.href = fullUrl
    else
      window.location.href = fullUrl
  }
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

function handleOpenModalSetting() {
  if (authStore.visitMode === VisitMode.VISIT_MODE_PUBLIC) {
    message.warning(t('common.publicModeRestricted'))
    return
  }
  searchBoxSettingShow.value = !searchBoxSettingShow.value
}

function getData() {
  moduleConfig.getValueByNameFromCloud<State>('deskModuleSearchBox').then(({ code, data }) => {
    if (code === 0)
      state.value = data || defaultState
    else
      state.value = defaultState

    selectEngineList.value = defaultSearchEngineList.concat(state.value.searchEngineList)
  })
}

function handleUpdateOpenMethod() {
  if (authStore.visitMode === VisitMode.VISIT_MODE_PUBLIC) {
    return
  }
  moduleConfig.saveToCloud(moduleConfigName, state.value)
}

function handleKeyPress(event: KeyboardEvent) {
  // 检查当前是否是输入框
  const target = event.target as HTMLElement
  const isInputElement = target.tagName === 'INPUT' || target.tagName === 'TEXTAREA' || target.isContentEditable

  if (event.key === '/' && !isInputElement) {
    const input = searchInputRef.value

    if (input && !input.matches(':focus')) {
      // 如果输入框没有焦点，则设置焦点并阻止默认行为
      event.preventDefault()
      input.focus()
    }
  }
}

onMounted(() => {
  window.addEventListener('keydown', handleKeyPress)

  getData()
})

onBeforeUnmount(() => {
  window.removeEventListener('keydown', handleKeyPress)
})
</script>

<template>
  <div class="search-box w-full" @keydown.enter="handleOnKeyDownEnter" @keydown.esc="handleClearSearchTerm">
    <div
      class="search-container flex rounded-2xl items-center justify-center text-white w-full"
      :style="{ background, color: textColor, borderColor }"
      :class="{ focused: isFocused }"
    >
      <div class="search-box-btn-engine w-[40px] flex justify-center cursor-pointer" @click="handleEngineClick">
        <NAvatar :src="state.currentSearchEngine.iconSrc" style="background-color: transparent;" :size="20" />
      </div>

      <input ref="searchInputRef" v-model="searchTerm" :placeholder="$t('deskModule.searchBox.inputPlaceholder')" @focus="onFocus" @blur="onBlur" @input="handleItemSearch">

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
        <div class="flex items-center flex-wrap">
          <div
            v-for="item, index in selectEngineList"
            :key="index"
            :title="item.title"
            class="w-[40px] h-[40px] mr-[10px]  cursor-pointer bg-[#ffffff] flex items-center justify-center rounded-xl my-1"
            @click="handleEngineUpdate(item)"
          >
            <NAvatar :src="item.iconSrc" style="background-color: transparent;" :size="20" />
          </div>

          <!-- 搜索引擎管理 -->
          <div
            v-if="authStore.visitMode === VisitMode.VISIT_MODE_LOGIN"
            class="w-[40px] h-[40px] mr-[10px]  cursor-pointer bg-[#ffffff] flex items-center justify-center rounded-xl my-1"
            @click="handleOpenModalSetting"
          >
            <NAvatar style="background-color: transparent;" :size="30">
              <SvgIcon icon="lets-icons--setting-alt-fill" style="font-size: 20px;color:#475569" />
            </NAvatar>
          </div>
        </div>
      </div>

      <div class="mt-[10px]">
        <NCheckbox v-model:checked="state.newWindowOpen" @update-checked="handleUpdateOpenMethod">
          <span :style="{ color: textColor }">
            {{ $t('deskModule.searchBox.openWithNewOpen') }}
          </span>
        </NCheckbox>
      </div>
    </div>

    <!-- 搜索栏设置窗口 -->
    <NModal
      v-model:show="searchBoxSettingShow"
      preset="card"
      size="small"
      style="width: 600px;border-radius: 1rem;"
      :title="t('deskModule.searchBox.searchEngineManage')"
      @after-leave="getData"
    >
      <SearchBoxSetting />
    </NModal>
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
  backdrop-filter:blur(10px)
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
