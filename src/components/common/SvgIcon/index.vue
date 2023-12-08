<script setup lang="ts">
import { computed, onMounted, ref, useAttrs } from 'vue'

const props = defineProps<{
  icon?: string
}>()

const attrs = useAttrs()
const iconName = ref(compatibleName(props.icon || ''))
const iconContent = ref<string | null>(null)
const bindAttrs = computed<{ class: string; style: string }>(() => ({
  class: (attrs.class as string) || '',
  style: (attrs.style as string) || '',
}))

const loadIcon = async () => {
  try {
    const iconPath = import.meta.glob('@/assets/svg-icons/*.svg')
    const iconModule = await iconPath[`/src/assets/svg-icons/${iconName.value}.svg`]()
    const moduleDefault = iconModule as { default: string }
    const response = await fetch(moduleDefault.default)
    const content = await response.text()

    if (isValidSvg(content))
      iconContent.value = content
    else
      console.error(`Invalid SVG format for icon ${iconName.value}`)
  }
  catch (error) {
    console.error(`Error loading icon ${iconName.value}:`, error)
  }
}

function isValidSvg(content: string): boolean {
  const parser = new DOMParser()
  const doc = parser.parseFromString(content, 'image/svg+xml')
  return doc.documentElement.tagName.toLowerCase() === 'svg'
}

function compatibleName(inputString: string): string {
  // 使用正则表达式替换所有的冒号
  const resultString = inputString.replace(/:/g, '-')
  return resultString
}

onMounted(() => {
  loadIcon()
})
</script>

<template>
  <div v-if="iconContent" v-bind="bindAttrs" v-html="iconContent" />
</template>
