<script setup lang="ts">
import { ref } from 'vue'

const props = withDefaults(defineProps<{
  background?: string
  textColor?: string
  value: string
  placeholder?: string
}>(), {
  background: '#2a2a2a6b',
  textColor: 'white',
})

const isFocused = ref(false)

const onFocus = (): void => {
  isFocused.value = true
}

const onBlur = (): void => {
  isFocused.value = false
}

function handleEngineClick() {

}

function handleSearchClick() {

}
</script>

<template>
  <div class="search-container  flex rounded-2xl items-center justify-center text-white w-full" :style="{ background, color: textColor }" :class="{ focused: isFocused }">
    <slot name="prefix" />
    <input v-model="value" placeholder="请输入搜索内容" @focus="onFocus" @blur="onBlur">
    <slot name="suffix" />
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
