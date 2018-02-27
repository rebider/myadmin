import Vue from 'vue'
import Router from 'vue-router'

// in development-env not use lazy-loading, because lazy-loading too many pages will cause webpack hot update too slow. so only in production use lazy-loading;
// detail: https://panjiachen.github.io/vue-element-admin-site/#/lazy-loading

Vue.use(Router)

/* Layout */
import Layout from '../views/layout/Layout'

/**
* hidden: true                   if `hidden:true` will not show in the sidebar(default is false)
* alwaysShow: true               if set true, will always show the root menu, whatever its child routes length
*                                if not set alwaysShow, only more than one route under the children
*                                it will becomes nested mode, otherwise not show the root menu
* redirect: noredirect           if `redirect:noredirect` will no redirct in the breadcrumb
* name:'router-name'             the name is used by <keep-alive> (must set!!!)
* meta : {
    title: 'title'               the name show in submenu and breadcrumb (recommend set)
    icon: 'svg-name'             the icon show in the sidebar,
  }
**/
export const asyncRouterMap = [
  {
    path: '/set',
    component: Layout,
    redirect: '/set/userList',
    name: 'Example',
    meta: { title: '设置', icon: 'example' },
    children: [
      {
        path: 'userList',
        name: 'Table',
        component: () => import('@/views/userList'),
        meta: { title: '用户管理', icon: 'table' }
      },
      {
        path: 'role',
        name: 'Role',
        component: () => import('@/views/tree/index'),
        meta: { title: '角色管理', icon: 'tree' }
      },
      {
        path: 'resource',
        name: 'Resource',
        component: () => import('@/views/tree/index'),
        meta: { title: '资源管理', icon: 'tree' }
      }
    ]
  },

  {
    path: '/server',
    component: Layout,
    redirect: '/server/gameserver',
    meta: { title: '服务器管理', icon: 'example' },
    children: [
      {
        path: 'gameserver',
        name: 'Gameserver',
        component: () => import('@/views/form/index'),
        meta: { title: '游戏服', icon: 'form' }
      }
      // {
      //   path: 'node',
      //   name: 'Node',
      //   component: () => import('@/views/form/index'),
      //   meta: { title: '节点', icon: 'form' }
      // }
    ]
  }
]
export const constantRouterMap = [
  { path: '/login', component: () => import('@/views/login/index'), hidden: true },
  { path: '/404', component: () => import('@/views/404'), hidden: true },
  {
    path: '',
    component: Layout,
    redirect: 'dashboard',
    children: [{
      path: 'dashboard',
      name: 'dashboard',
      component: () => import('@/views/dashboard/index'),
      meta: { title: '首页', icon: 'dashboard', noCache: true }
    }]
  },
  // {
  //   path: '/changePasswd',
  //   component: Layout,
  //   hidden: true,
  //   children: [{ path: '/changePasswd', component: () => import('@/views/changePasswd') }]
  // },
  { path: '*', redirect: '/404', hidden: true }
]

export default new Router({
  // mode: 'history', //后端支持可开
  scrollBehavior: () => ({ y: 0 }),
  routes: constantRouterMap
})

