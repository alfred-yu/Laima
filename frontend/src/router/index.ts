import { createRouter, createWebHistory } from 'vue-router'
import { useUserStore } from '../stores'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/',
      name: 'dashboard',
      component: () => import('../views/Dashboard.vue'),
      meta: { requiresAuth: true }
    },
    {
      path: '/repos',
      name: 'repos',
      component: () => import('../views/RepoList.vue'),
      meta: { requiresAuth: true }
    },
    {
      path: '/repos/:owner/:repo',
      name: 'repo-detail',
      component: () => import('../views/RepoDetail.vue'),
      meta: { requiresAuth: true },
      children: [
        {
          path: 'code',
          name: 'repo-code',
          component: () => import('../views/RepoCode.vue'),
          meta: { requiresAuth: true }
        },
        {
          path: 'pulls',
          name: 'repo-pulls',
          component: () => import('../views/RepoPulls.vue'),
          meta: { requiresAuth: true }
        },
        {
          path: 'issues',
          name: 'repo-issues',
          component: () => import('../views/RepoIssues.vue'),
          meta: { requiresAuth: true }
        },
        {
          path: 'cicd',
          name: 'repo-cicd',
          component: () => import('../views/RepoCICD.vue'),
          meta: { requiresAuth: true }
        }
      ]
    },
    {
      path: '/pulls',
      name: 'pulls',
      component: () => import('../views/PullList.vue'),
      meta: { requiresAuth: true }
    },
    {
      path: '/pulls/:id',
      name: 'pull-detail',
      component: () => import('../views/PullDetail.vue'),
      meta: { requiresAuth: true }
    },
    {
      path: '/issues',
      name: 'issues',
      component: () => import('../views/IssueList.vue'),
      meta: { requiresAuth: true }
    },
    {
      path: '/issues/:id',
      name: 'issue-detail',
      component: () => import('../views/IssueDetail.vue'),
      meta: { requiresAuth: true }
    },
    {
      path: '/cicd',
      name: 'cicd',
      component: () => import('../views/CICDList.vue'),
      meta: { requiresAuth: true }
    },
    {
      path: '/cicd/:id',
      name: 'cicd-detail',
      component: () => import('../views/CICDDetail.vue'),
      meta: { requiresAuth: true }
    },
    {
      path: '/users',
      name: 'users',
      component: () => import('../views/UserList.vue'),
      meta: { requiresAuth: true }
    },
    {
      path: '/users/:id',
      name: 'user-detail',
      component: () => import('../views/UserDetail.vue'),
      meta: { requiresAuth: true }
    },
    {
      path: '/login',
      name: 'login',
      component: () => import('../views/Login.vue'),
      meta: { requiresAuth: false }
    },
    {
      path: '/register',
      name: 'register',
      component: () => import('../views/Register.vue'),
      meta: { requiresAuth: false }
    },
    {
      path: '/forgot-password',
      name: 'forgot-password',
      component: () => import('../views/ForgotPassword.vue'),
      meta: { requiresAuth: false }
    },
    {
      path: '/settings',
      name: 'settings',
      component: () => import('../views/Setting.vue'),
      meta: { requiresAuth: true }
    },
    {
      path: '/activity',
      name: 'activity',
      component: () => import('../views/UserActivity.vue'),
      meta: { requiresAuth: true }
    }
  ]
})

// 路由守卫
router.beforeEach((to, from, next) => {
  const userStore = useUserStore()
  
  // 检查是否需要认证
  const requiresAuth = to.matched.some(record => record.meta.requiresAuth !== false)
  
  if (requiresAuth) {
    if (!userStore.isLoggedIn) {
      // 未登录，跳转到登录页面
      next('/login')
    } else {
      // 已登录，检查是否是根路径
      if (to.path === '/') {
        // 跳转到用户页面（假设用户ID为1，实际应该从用户信息中获取）
        next('/dashboard')
      } else {
        // 其他受保护页面，正常访问
        next()
      }
    }
  } else {
    // 不需要认证的页面，正常访问
    next()
  }
})

export default router