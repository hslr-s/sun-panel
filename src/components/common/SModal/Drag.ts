interface HTMLElementWithCurrentStyle extends HTMLElement {
  currentStyle?: {
    [key: string]: string
  }
}

// 获取相关CSS属性
const getCss = function (o: HTMLElementWithCurrentStyle, key: string): string {
  if (o.currentStyle) {
    return o.currentStyle[key]
  }
  else {
    const computedStyle = document.defaultView?.getComputedStyle(o, null)
    return computedStyle ? computedStyle.getPropertyValue(key) : ''
  }
}

const params = {
  left: 0,
  top: 0,
  currentX: 0,
  currentY: 0,
  flag: false,
}

const startDrag = function (
  bar: HTMLElement,
  target: HTMLElementWithCurrentStyle,
  callback?: (left: number, top: number) => void,
) {
  const screenWidth = document.body.clientWidth // body当前宽度
  const screenHeight = document.documentElement.clientHeight // 可见区域高度

  const dragDomW = target.offsetWidth // 对话框宽度
  const dragDomH = target.offsetHeight // 对话框高度

  const minDomLeft = target.offsetLeft
  const minDomTop = target.offsetTop

  const maxDragDomLeft = screenWidth - minDomLeft - dragDomW
  const maxDragDomTop = screenHeight - minDomTop - dragDomH

  if (getCss(target, 'left') !== 'auto')
    params.left = parseInt(getCss(target, 'left'))

  if (getCss(target, 'top') !== 'auto')
    params.top = parseInt(getCss(target, 'top'))

  // o是移动对象
  bar.onmousedown = function (event: MouseEvent) {
    params.flag = true
    const e = event || (window.event as MouseEvent)

    // 防止IE文字选中
    bar.onselectstart = function () {
      return false
    }

    params.currentX = e.clientX
    params.currentY = e.clientY
  }

  document.onmouseup = function () {
    params.flag = false
    if (getCss(target, 'left') !== 'auto')
      params.left = parseInt(getCss(target, 'left'))

    if (getCss(target, 'top') !== 'auto')
      params.top = parseInt(getCss(target, 'top'))
  }

  document.onmousemove = function (event: MouseEvent) {
    const e = event || (window.event as MouseEvent)
    if (params.flag) {
      const nowX = e.clientX
      const nowY = e.clientY
      const disX = nowX - params.currentX
      const disY = nowY - params.currentY

      let left = params.left + disX
      let top = params.top + disY

      // 拖出屏幕边缘
      if (-left > minDomLeft)
        left = -minDomLeft
      else if (left > maxDragDomLeft)
        left = maxDragDomLeft

      if (-top > minDomTop)
        top = -minDomTop
      else if (top > maxDragDomTop)
        top = maxDragDomTop

      target.style.left = `${left}px`
      target.style.top = `${top}px`

      if (typeof callback === 'function')
        callback(left, top)

      if (event.preventDefault)
        event.preventDefault()

      return false
    }
  }
}

export default startDrag
