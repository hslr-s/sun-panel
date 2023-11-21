<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { NAvatar } from 'naive-ui'
import { SvgIcon } from '@/components/common'
import { useModuleConfig } from '@/store/modules'

import SvgSrcBaidu from '@/assets/search_engine_svg/baidu.svg'
import SvgSrcBing from '@/assets/search_engine_svg/bing.svg'
import SvgSrcGoogle from '@/assets/search_engine_svg/google.svg'

interface State {
  currentSearchEngine: DeskModule.SearchBox.SearchEngine
  searchEngineList: DeskModule.SearchBox.SearchEngine[]
}

const props = withDefaults(defineProps<{
  background?: string
  textColor?: string
}>(), {
  background: '#2a2a2a6b',
  textColor: 'white',
})

const moduleConfigName = 'deskModuleSearchBox'
const moduleConfig = useModuleConfig()
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
}

const state = ref<State>(defaultState)

const onFocus = (): void => {
  isFocused.value = true
}

const onBlur = (): void => {
  isFocused.value = false
}

function handleEngineClick() {
  searchSelectListShow.value = !searchSelectListShow.value
}

function handleEngineUpdate(engine: DeskModule.SearchBox.SearchEngine) {
  state.value.currentSearchEngine = engine
  moduleConfig.saveToCloud(moduleConfigName, state.value)
}

function handleSearchClick() {
  // 文本替换关键字 %s 同时 如果没有%s,将吧关键字放在网址最后
}

onMounted(() => {
  moduleConfig.getValueByNameFromCloud<State>('deskModuleSearchBox').then(({ code, data }) => {
    if (code === 0)
      state.value = data
    else
      state.value = defaultState
  })
})
</script>

<template>
  <div class="w-full">
    <div class="search-container flex rounded-2xl items-center justify-center text-white w-full" :style="{ background, color: textColor }" :class="{ focused: isFocused }">
      <div class="w-[40px] flex justify-center cursor-pointer" @click="handleEngineClick">
        <NAvatar :src="state.currentSearchEngine.iconSrc" style="background-color: transparent;" :size="20" />
      </div>

      <input v-model="searchTerm" placeholder="请输入搜索内容" @focus="onFocus" @blur="onBlur">
      <div class="w-[20px] flex justify-center cursor-pointer" @click="handleSearchClick">
        <SvgIcon icon="iconamoon:search-fill" />
      </div>
    </div>

    <!-- 搜索引擎选择 -->
    <div v-if="searchSelectListShow" class="w-full mt-[10px] bg-white h-[60px] flex items-center rounded-xl p-[10px]" :style="{ background }">
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
  </div>
</template>

<style scoped>
.search-container {
  border: 1px solid #ccc;
  transition: border-color 0.5s;
  padding: 2px 10px;
}

.focused {
  border-color: white;
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
