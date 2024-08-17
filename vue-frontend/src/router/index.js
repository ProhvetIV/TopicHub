import { createRouter, createWebHistory } from 'vue-router'
import LoginPage from '../views/LoginPage.vue'

// Initialize the WebSocket service with the router


const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'yes',
      component: () => import('@/views/LoginPage.vue') //LoginPage //() => import('@/views/LoginPage.vue') //
    },
    {
      path: '/login',
      name: 'Login',
      component: () => import('@/views/LoginPage.vue') //LoginPage //() => import('@/views/LoginPage.vue') //
    },
    {
      path: '/register',
      name: 'Register',
      component: () => import('@/views/RegisterPage.vue')
    },
    {
      path: '/home',
      name: 'Home',
      component: () => import('@/views/Home.vue')
    },
    {
      path: '/post',
      name: 'Post',
      component: () => import('@/views/PostingPage.vue')
    },
    {
      path: '/groups',
      name: 'Groups',
      component: () => import('@/views/Groups.vue')
    },
    {
      path: '/group',
      name: 'Group',
      component: () => import('@/views/Group.vue')
    },
    
    { 
      path: '/profile',
      name: 'Profile',
      component: () => import('@/views/ProfilePage.vue')
    },
    {
      path: '/accountprivate',
      name: 'AccountPrivate',
      component: () => import('@/views/AccountPrivatePage.vue')
    },

  ]
})

export default router
