import type { Router } from 'vue-router'
import { useUserStore } from '@/store/modules/user'

export function setupPageGuard(router: Router) {
  router.beforeEach(async (to, from, next) => {
    // const authStore = useAuthStoreWithout()
    const userStore = useUserStore()
    // 非管理员路由拦截
    if (userStore.userInfo.role !== 1 && to.path.includes('admin'))
      next({ name: '404' })

    else
      next()

    // if (!authStore.session) {
    //   try {
    //     const data = await authStore.getSession()
    //     if (String(data.auth) === 'false' && authStore.token)
    //       authStore.removeToken()
    //     if (to.path === '/500')
    //       next({ name: 'Root' })
    //     else
    //       next()
    //   }
    //   catch (error) {
    //     if (to.path !== '/500')
    //       next({ name: '500' })
    //     else
    //       next()
    //   }
    // }
    // else {
    //   next()
    // }
  })
}
