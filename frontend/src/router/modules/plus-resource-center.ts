import { AppRouteRecord } from '@/types/router'

export const plusResourceCenterRoutes: AppRouteRecord = {
  path: '/plus-resources',
  name: 'PlusResources',
  component: '/index/index',
  meta: {
    title: 'Plus 资源中心',
    icon: 'ri:command-line',
    roles: ['R_SUPER', 'R_ADMIN']
  },
  children: [
    {
      path: 'control',
      name: 'PlusResourceControl',
      component: '/plus-resource-center/index',
      meta: {
        title: '插件 / 模板 / 广告',
        icon: 'ri:box-3-line',
        keepAlive: true
      }
    }
  ]
}
