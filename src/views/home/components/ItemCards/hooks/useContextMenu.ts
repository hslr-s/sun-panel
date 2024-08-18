import { ref } from 'vue'

// 自定义钩子 useContextMenu
function useContextMenu() {
  const dropdownShow = ref(false)
  const dropdownMenuX = ref(0)
  const dropdownMenuY = ref(0)

  // 处理鼠标右键点击事件
  function handleContextMenu(event: MouseEvent) {
    event.preventDefault()
    if (props.sortStatus)
      return
    dropdownShow.value = true
    dropdownMenuX.value = event.clientX
    dropdownMenuY.value = event.clientY
  }

  // 返回钩子的结果
  return {
    dropdownShow,
    dropdownMenuX,
    dropdownMenuY,
    handleContextMenu,
  }
}
