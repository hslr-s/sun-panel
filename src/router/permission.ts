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
  })
}
