import type { RouteRecordRaw } from 'vue-router'
export const dockerRouter: RouteRecordRaw = {
  path: '/docker',
  name: 'Docker',
  // component: () => import('@/views/home/Layout.vue'),
  component: () => import('@/views/docker/index.vue'),
}
