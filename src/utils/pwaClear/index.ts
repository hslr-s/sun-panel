async function unregisterServiceWorkers() {
  const registrations = await navigator.serviceWorker.getRegistrations()
  for (const registration of registrations) {
    await registration.unregister()
    console.log('Service Worker unregistered:', registration)
  }
}

async function clearCacheStorage() {
  if ('caches' in window) {
    const cacheNames = await caches.keys()
    for (const name of cacheNames) {
      await caches.delete(name)
      console.log('Cache deleted:', name)
    }
  }
}

// function clearIndexedDB() {
//   return new Promise((resolve, reject) => {
//     if ('indexedDB' in window) {
//       const DBOpenRequest = indexedDB.deleteDatabase('your-database-name') // 替换为你的数据库名称
//       DBOpenRequest.onsuccess = () => {
//         console.log('IndexedDB deleted')
//         resolve()
//       }
//       DBOpenRequest.onerror = (error) => {
//         console.error('Error deleting IndexedDB:', error)
//         reject(error)
//       }
//     }
//     else {
//       resolve()
//     }
//   })
// }

function clearLocalStorage() {
  localStorage.clear()
  console.log('LocalStorage cleared')
}

function clearSessionStorage() {
  sessionStorage.clear()
  console.log('SessionStorage cleared')
}

function clearCookies() {
  document.cookie.split(';').forEach((cookie) => {
    const [name] = cookie.split('=')
    document.cookie = `${name}=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/;`
  })
  console.log('Cookies cleared')
}

export async function clearAllCache() {
  try {
    await unregisterServiceWorkers()
    await clearCacheStorage()
    // await clearIndexedDB()
    clearLocalStorage()
    clearSessionStorage()
    clearCookies()
    console.log('All cache cleared successfully!')
    alert('所有缓存已清除！')

    // await registerServiceWorker()
    window.location.reload() // 刷新页面
  }
  catch (error) {
    console.error('Error clearing cache:', error)
    alert('清除缓存时出错，请重试。')
  }
}
