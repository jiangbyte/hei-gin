/** Author: Charlie */

import type { RouteRecordRaw } from 'vue-router'

/**
 * 内置基础路由。
 *
 * 这些路由不是后端 SysResource 资源，不参与侧边菜单和标签页生成。
 * 授权业务路由会在 route store 初始化时动态挂载到 appRoot 下。
 */
export const routes: RouteRecordRaw[] = [
  {
    // 根路径只作为入口，由路由守卫重定向到 VITE_HOME_PATH。
    path: '/',
    name: 'root',
    children: [],
  },
  {
    // 显式 404 页面，用于守卫初始化失败或主动跳转。
    path: '/not-found',
    name: 'not-found',
    component: () => import('@/views/error/NotFound.vue'),
  },
  {
    // 认证页面是公开静态路由，不进入后台 Layout、菜单和标签栏。
    path: '/auth',
    name: 'auth',
    redirect: '/auth/login',
  },
  {
    path: '/auth/login',
    name: 'auth-login',
    component: () => import('@/views/auth/Login.vue'),
    meta: { name: '登录' },
  },
  {
    path: '/auth/forgot-password',
    name: 'auth-forgot-password',
    component: () => import('@/views/auth/ForgotPassword.vue'),
    meta: { name: '忘记密码' },
  },
  {
    // 兜底路由必须放在最后。动态路由注册完成后，守卫会重新匹配原始目标地址。
    path: '/:pathMatch(.*)*',
    name: 'not-found-catch',
    component: () => import('@/views/error/NotFound.vue'),
  },
]
