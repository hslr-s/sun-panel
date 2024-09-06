<script setup lang="ts">
import { nextTick, onMounted, ref, watch, withDefaults } from 'vue'
import { SvgIcon } from '..'
import startDrag from './Drag'

const props = withDefaults(defineProps<{
  show: boolean
  width?: number | string
  height?: number | string
  title?: string
  left?: number | string
  top?: number | string
  resize?: boolean
  move?: boolean
}>(), {
  width: 500,
  height: 500,
  resize: false,
  move: false,
})

const emit = defineEmits<{
  (e: 'update:show', show: boolean): void
  (e: 'onClose'): void
}>()

let isDragging = false
let startX: number, startY: number
let startWidth: number, startHeight: number
const elementId = getTimestampWithRandom()
const showModal = ref(props.show)

const containerId = `free-layer-container-${elementId}`
const moveId = `free-layer-move-${elementId}`
const resizerId = `free-layer-resizer-${elementId}`

const container = ref()
const resizer = ref()

// 值被动改变之后，需要手动更新
// 否则在iframe加载完成事件后会重新渲染，导致重置
const containerPosition: { [key: string]: string } = {
  left: numberToPx(props.left || 0),
  top: numberToPx(props.top || 0),
  width: numberToPx(props.width),
  height: numberToPx(props.height),
}

const freeLayerResizerStyle = ref<any>({
  width: '40px',
  height: '40px',
  // backgroundColor: '#4CAF50',
  position: 'absolute',
  bottom: 0,
  right: 0,
  cursor: 'se-resize',
})

function onMouseMove(e: MouseEvent) {
  if (isDragging) {
    const dx = e.clientX - startX
    const dy = e.clientY - startY
    let newWidth = startWidth + dx
    let newHeight = startHeight + dy

    // 限制新的宽度和高度不超过可视区域
    const maxWidth = window.innerWidth - container.value.offsetLeft
    const maxHeight = window.innerHeight - container.value.offsetTop

    // 确保新的宽度和高度不小于初始大小
    // newWidth = Math.max(startWidth, Math.min(newWidth, maxWidth))
    // newHeight = Math.max(startHeight, Math.min(newHeight, maxHeight))

    // 确保新的宽度和高度不小于200*200
    newWidth = Math.max(200, Math.min(newWidth, maxWidth))
    newHeight = Math.max(200, Math.min(newHeight, maxHeight))

    container.value.style.width = `${newWidth}px`
    container.value.style.height = `${newHeight}px`
  }
}

function numberToPx(value: string | number) {
  return typeof value === 'string' ? value : `${value}px`
}

function setZIndex() {
  // 动态设置 z-index 确保当前操作的弹窗在最前面
  const maxZIndex = Array.from(document.querySelectorAll('.free-layer-container')).reduce((max, el) => {
    const zIndex = parseInt(window.getComputedStyle(el).zIndex || '10', 10)
    return zIndex > max ? zIndex : max
  }, 10)

  const target = document.getElementById(containerId)

  if (target)
    target.style.zIndex = `${maxZIndex + 1}`
}

function onMouseUp(e: MouseEvent) {
  isDragging = false

  containerPosition.height = `${container.value.offsetHeight}px`
  containerPosition.width = `${container.value.offsetWidth}px`

  document.removeEventListener('mousemove', onMouseMove)
  document.removeEventListener('mouseup', onMouseUp)

  freeLayerResizerStyle.value.width = '40px'
  freeLayerResizerStyle.value.height = '40px'
}

function getTimestampWithRandom() {
  // 获取当前时间戳（毫秒级别）
  const timestamp = Date.now()
  // 生成一个0到9999之间的随机数
  const randomNum = Math.floor(Math.random() * 10000)
  // 将时间戳和随机数拼接成字符串
  return `${timestamp}${randomNum}`
}

function closeModal() {
  emit('update:show', false)
  emit('onClose')
}

function init() {
  containerPosition.left = props.left ? `${numberToPx(props.left)}` : `calc(100% / 2 - ${containerPosition.width} / 2)`
  containerPosition.top = props.top ? `${numberToPx(props.top)}` : `calc(100% / 2 - ${containerPosition.height} / 2)`
  //   containerPosition.value.left = `${window.innerWidth / 2 - props.width / 2}px`
  //   containerPosition.value.top = `${window.innerHeight / 2 - props.height / 2}px`
  nextTick(() => {
    resizer.value = document.getElementById(resizerId)
    container.value = document.getElementById(containerId)
    if (resizer.value) {
      // 监听尺寸改变
      resizer.value.addEventListener('mousedown', (e: MouseEvent) => {
        isDragging = true
        freeLayerResizerStyle.value.width = '100%'
        freeLayerResizerStyle.value.height = '100%'

        startX = e.clientX
        startY = e.clientY
        startWidth = container.value.offsetWidth
        startHeight = container.value.offsetHeight
        document.addEventListener('mousemove', onMouseMove)
        document.addEventListener('mouseup', onMouseUp)
        e.preventDefault() // 防止默认事件
      })
      const oBox = document.getElementById(containerId)
      const oBar = document.getElementById(moveId)
      if (props.move) {
        startDrag(oBar as HTMLElement, oBox as HTMLElement, (left: number, top: number) => {
          containerPosition.left = `${left}px`
          containerPosition.top = `${top}px`
        })
      }
    }
    setZIndex()
  })
}

watch(() => props.show, (v: boolean) => {
  showModal.value = v
  if (v) {
    init()
  }
})

onMounted(() => {
  init()
})
</script>

<template>
  <Teleport to="body">
    <transition name="free-layer-transition">
      <div
        v-if="showModal"
        v-show="showModal"
        :id="containerId" class="free-layer-container shadow-2xl bg-white dark:bg-[#2C2C32]"
        :style="containerPosition"
        @mousedown="setZIndex()"
      >
        <div
          :id="moveId"
          class="py-2 px-4 flex items-center h-[40px] dark:text-white select-none"
          :class="move ? 'cursor-move' : ''"
        >
          <div class="mr-auto">
            <slot name="header">
              {{ title }}
            </slot>
          </div>
          <div>
            <slot name="header-close">
              <div class="cursor-pointer" @click="closeModal">
                <SvgIcon icon="mingcute--close-line" style="width: 20px;height: 20px;" />
              </div>
            </slot>
          </div>
        </div>
        <slot>
          <div
            class="content m-2 p-2 rounded-lg bg-slate-500" :style="{
              height: 'calc(100% - 60px)',

            }"
          >
            <slot name="content" />
          </div>
        </slot>

        <div
          v-show="resize"
          :id="resizerId"
          :style="freeLayerResizerStyle"
          class="rounded-tl-2xl hover:border-spacing-4 hover:border-dashed hover:border-emerald-300 hover:border hover:bg-[#6ee7b761]"
        />
      </div>
    </transition>
  </Teleport>
</template>

<style scoped>
.free-layer-header{
    padding: 10px;
}

.free-layer-container {
    background-color: #fff;
    position: fixed ;
    overflow: hidden;
    border-radius: 1rem;
    z-index: 10;
}

.free-layer-transition-enter-active, .free-layer-transition-leave-active {
  transition: opacity 0.2s ease, transform 0.2s ease;
}
.free-layer-transition-enter-from , .free-layer-transition-leave-to {
  opacity: 0;
  transform: scale(0.9);
}

.free-layer-transition-enter-to, .free-layer-transition-leave-from {
  opacity: 1;
  transform: scale(1);
}
</style>
