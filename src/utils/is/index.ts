export function isNumber<T extends number>(value: T | unknown): value is number {
  return Object.prototype.toString.call(value) === '[object Number]'
}

export function isString<T extends string>(value: T | unknown): value is string {
  return Object.prototype.toString.call(value) === '[object String]'
}

export function isBoolean<T extends boolean>(value: T | unknown): value is boolean {
  return Object.prototype.toString.call(value) === '[object Boolean]'
}

export function isNull<T extends null>(value: T | unknown): value is null {
  return Object.prototype.toString.call(value) === '[object Null]'
}

export function isUndefined<T extends undefined>(value: T | unknown): value is undefined {
  return Object.prototype.toString.call(value) === '[object Undefined]'
}

export function isObject<T extends object>(value: T | unknown): value is object {
  return Object.prototype.toString.call(value) === '[object Object]'
}

export function isArray<T extends any[]>(value: T | unknown): value is T {
  return Object.prototype.toString.call(value) === '[object Array]'
}

export function isFunction<T extends (...args: any[]) => any | void | never>(value: T | unknown): value is T {
  return Object.prototype.toString.call(value) === '[object Function]'
}

export function isDate<T extends Date>(value: T | unknown): value is T {
  return Object.prototype.toString.call(value) === '[object Date]'
}

export function isRegExp<T extends RegExp>(value: T | unknown): value is T {
  return Object.prototype.toString.call(value) === '[object RegExp]'
}

export function isPromise<T extends Promise<any>>(value: T | unknown): value is T {
  return Object.prototype.toString.call(value) === '[object Promise]'
}

export function isSet<T extends Set<any>>(value: T | unknown): value is T {
  return Object.prototype.toString.call(value) === '[object Set]'
}

export function isMap<T extends Map<any, any>>(value: T | unknown): value is T {
  return Object.prototype.toString.call(value) === '[object Map]'
}

export function isFile<T extends File>(value: T | unknown): value is T {
  return Object.prototype.toString.call(value) === '[object File]'
}
/**
 * @Author: wintsa
 * @Date: 2024-03-18 16:32:53
 * @LastEditors: wintsa
 * @Description: 判断ip是否为局域网
 * @return {*}
 */
export function isLocalUrl(url:string) {
  // 创建一个 <a> 元素
  var anchor = document.createElement('a');
  // 设置其 href 属性为所检查的 URL
  anchor.href = url;
  // 获取主机名
  var hostname = anchor.hostname;
  
  // 检查主机名是否是局域网的
  // 这里假设局域网主机名的特征
  var ip = hostname.split('.');
  if (
      // 检查是否为localhost
      hostname === 'localhost' || 
      // 检查是否为私有IP地址
      (ip.length === 4 && (
          ip[0] === '10' ||
          (ip[0] === '172' && parseInt(ip[1]) >= 16 && parseInt(ip[1]) <= 31) ||
          (ip[0] === '192' && ip[1] === '168')
      ))
  ) {
      return true;
  } else {
      return false;
  }
}
