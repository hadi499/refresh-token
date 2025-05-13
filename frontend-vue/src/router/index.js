import { createRouter, createWebHistory } from 'vue-router';

import Home from '../views/Home.vue'
import Login from '../views/Login.vue';


const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/',
      name: 'Home', // Nama route tetap 'Home'
      component: Home,
      meta: { requiresAuth: true }
    },    
    {
      path: '/login',
      name: 'Login',
      component: Login,
      meta: { requiresAuth: false }
    },
  ]
});

// Navigation guard
router.beforeEach((to, from, next) => {
  const authTokens = localStorage.getItem('authTokens');

  if (to.meta.requiresAuth && !authTokens) {
    next('/login');
  } else {
    next();
  }
});

export default router;