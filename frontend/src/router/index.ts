import { createRouter, createWebHistory } from 'vue-router'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/',
      name: 'dashboard',
      component: () => import('../views/Dashboard.vue')
    },
    {
      path: '/repos',
      name: 'repos',
      component: () => import('../views/RepoList.vue')
    },
    {
      path: '/repos/:owner/:repo',
      name: 'repo-detail',
      component: () => import('../views/RepoDetail.vue'),
      children: [
        {
          path: 'code',
          name: 'repo-code',
          component: () => import('../views/RepoCode.vue')
        },
        {
          path: 'pulls',
          name: 'repo-pulls',
          component: () => import('../views/RepoPulls.vue')
        },
        {
          path: 'issues',
          name: 'repo-issues',
          component: () => import('../views/RepoIssues.vue')
        },
        {
          path: 'cicd',
          name: 'repo-cicd',
          component: () => import('../views/RepoCICD.vue')
        }
      ]
    },
    {
      path: '/pulls',
      name: 'pulls',
      component: () => import('../views/PullList.vue')
    },
    {
      path: '/pulls/:id',
      name: 'pull-detail',
      component: () => import('../views/PullDetail.vue')
    },
    {
      path: '/issues',
      name: 'issues',
      component: () => import('../views/IssueList.vue')
    },
    {
      path: '/issues/:id',
      name: 'issue-detail',
      component: () => import('../views/IssueDetail.vue')
    },
    {
      path: '/cicd',
      name: 'cicd',
      component: () => import('../views/CICDList.vue')
    },
    {
      path: '/cicd/:id',
      name: 'cicd-detail',
      component: () => import('../views/CICDDetail.vue')
    },
    {
      path: '/users',
      name: 'users',
      component: () => import('../views/UserList.vue')
    },
    {
      path: '/users/:id',
      name: 'user-detail',
      component: () => import('../views/UserDetail.vue')
    },
    {
      path: '/settings',
      name: 'settings',
      component: () => import('../views/Settings.vue')
    },
    {
      path: '/login',
      name: 'login',
      component: () => import('../views/Login.vue')
    },
    {
      path: '/register',
      name: 'register',
      component: () => import('../views/Register.vue')
    }
  ]
})

export default router