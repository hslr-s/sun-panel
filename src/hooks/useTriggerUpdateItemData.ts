import { ref, watch } from 'vue'

// 定义全局状态
const updateCounter = ref(0) // 使用计数器来确保每次触发都有效

// 提供触发函数
export const useTriggerUpdateItemData = () => {
  const triggerUpdate = () => {
    // console.log('主动触发更新')
    updateCounter.value++ // 每次调用递增计数器
  }

  const onTriggered = (callback: () => void) => {
    // 监听计数器变化并执行回调
    watch(() => updateCounter.value, () => {
      // console.log('收到更新信号')
      callback()
    })
  }

  return { triggerUpdate, onTriggered }
}
