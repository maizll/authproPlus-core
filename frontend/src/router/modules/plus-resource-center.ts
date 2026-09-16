import { AppRouteRecord } from '@/types/router'

export const plusResourceCenterRoutes: AppRouteRecord = {
  path: '/plus-resources',
  name: 'PlusResources',
  component: '/index/index',
  meta: { title: 'Plus 资源中心', icon: 'ri:command-line', roles: ['R_SUPER', 'R_ADMIN'] },
  children: [
    { path: 'overview', name: 'PlusResourceOverview', component: '/plus-resource-center/index', meta: { title: '资源总览', icon: 'ri:dashboard-line', keepAlive: true } },
    { path: 'plugins', name: 'PlusResourcePlugins', component: '/plus-resource-center/plugins', meta: { title: '插件管理', icon: 'ri:extension-line', keepAlive: true } },
    { path: 'templates', name: 'PlusResourceTemplates', component: '/plus-resource-center/templates', meta: { title: '首页模板', icon: 'ri:layout-4-line', keepAlive: true } },
    { path: 'ads', name: 'PlusResourceAds', component: '/plus-resource-center/ads', meta: { title: '广告策略', icon: 'ri:advertisement-line', keepAlive: true } },
    { path: 'sources', name: 'PlusResourceSources', component: '/plus-resource-center/sources', meta: { title: '插件源', icon: 'ri:database-2-line', keepAlive: true } }
  ]
}
