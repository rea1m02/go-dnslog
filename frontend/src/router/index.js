import { createRouter, createWebHashHistory } from 'vue-router';
import Login from '../components/Login.vue';
import Register from '../components/Register.vue';
import DnsLogs from '../components/DnsLogs.vue';
import Rebind from '../components/Rebind.vue';
import Layout from '../components/Layout.vue';

// 路由守卫中间件
const requireAuth = (to, from, next) => {
  console.log('Route guard triggered:', to.path);
  const token = localStorage.getItem('token');
  if (!token && to.path !== '/login' && to.path !== '/register') {
    console.log('No token found, redirecting to login');
    next('/login');
  } else {
    console.log('Auth check passed, proceeding to:', to.path);
    next();
  }
};

const routes = [
  { path: '/login', component: Login },
  { path: '/register', component: Register },
  {
    path: '/',
    component: Layout,
    children: [
      { path: '', redirect: 'dns_logs' }, // 默认跳转到DNS日志
      { path: 'dns_logs', component: DnsLogs, beforeEnter: requireAuth },
      { path: 'rebind', component: Rebind, beforeEnter: requireAuth }
    ]
  },
  { path: '/:pathMatch(.*)*', redirect: '/dns_logs' }
];

const router = createRouter({
  history: createWebHashHistory(),
  routes
});

export default router;