import type { App } from 'vue'

import type { RouteRecordRaw } from 'vue-router'
import { createRouter, createWebHistory } from 'vue-router'
import { setupPageGuard } from './permission'
// import { dockerRouter } from './docker'

const routes: RouteRecordRaw[] = [
  {
    path: '/',
    name: 'Home',
    component: () => import('@/views/home/index.vue'),
  },

  {
    path: '/login',
    name: 'login',
    component: () => import('@/views/login/index.vue'),
  },

  {
    path: '/404',
    name: '404',
    component: () => import('@/views/exception/404/index.vue'),
  },

  {
    path: '/500',
    name: '500',
    component: () => import('@/views/exception/500/index.vue'),
  },

  {
    path: '/:pathMatch(.*)*',
    name: 'notFound',
    redirect: '/404',
  },

  {
    path: '/test',
    name: 'test',
    component: () => import('@/views/exception/test/index.vue'),
  },

  // dockerRouter,

  // {
  //   path: '/test/demo',
  //   name: 'test-demo',
  //   component: () => import('@/views/exception/test/demo.vue'),
  // },

  {
    path: '/clear',
    name: 'clearCache',
    component: () => import('@/views/clear/index.vue'),
  },

  // adminRouter,
]

export const router = createRouter({
  // history: createWebHashHistory(),// hash 模式
  history: createWebHistory(), // history 模式
  routes,
  scrollBehavior: () => ({ left: 0, top: 0 }),
})

setupPageGuard(router)

export async function setupRouter(app: App) {
  app.use(router)
  await router.isReady()
}
