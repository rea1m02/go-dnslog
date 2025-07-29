<template>
  <div class="min-h-screen flex items-center justify-center bg-gradient-to-br from-indigo-200 via-white to-indigo-100">
    <div class="bg-white p-8 rounded-lg shadow-md w-full max-w-xl">
      <h2 class="text-2xl font-bold mb-6 text-center">DNSLog 登录</h2>
      <form @submit.prevent="handleLogin" class="space-y-4">
        <div v-if="errorMessage" class="p-3 bg-red-50 border border-red-200 rounded-md text-red-700 text-sm">
          {{ errorMessage }}
        </div>
        <div>
          <label for="username" class="block text-sm font-medium text-gray-700">用户名</label>
          <div class="mt-1 relative rounded-md shadow-sm">
            <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
              <i class="fa fa-user text-gray-400"></i>
            </div>
            <input type="text" id="username" v-model="username" required
              class="focus:ring-indigo-500 focus:border-indigo-500 block w-full pl-10 pr-3 py-2 border border-gray-300 rounded-md">
          </div>
        </div>
        <div>
          <label for="password" class="block text-sm font-medium text-gray-700">密码</label>
          <div class="mt-1 relative rounded-md shadow-sm">
            <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
              <i class="fa fa-lock text-gray-400"></i>
            </div>
            <input type="password" id="password" v-model="password" required
              class="focus:ring-indigo-500 focus:border-indigo-500 block w-full pl-10 pr-3 py-2 border border-gray-300 rounded-md">
          </div>
        </div>
        <div>
          <button type="submit" :disabled="loading"
            class="w-full flex justify-center py-2 px-4 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 disabled:opacity-70">
            <i v-if="loading" class="fa fa-spinner fa-spin mr-2"></i>
            登录
          </button>
        </div>
        <div class="text-center">
          <router-link to="/register" class="text-sm text-indigo-600 hover:text-indigo-500">
            没有账号？注册
          </router-link>
          <span class="mx-2 text-gray-400">|</span>
          <a href="#" @click.prevent="handleRandomLogin" class="text-sm text-indigo-600 hover:text-indigo-500">
            随机ID登录
          </a>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue';
import { useRouter } from 'vue-router';

const username = ref('');
const password = ref('');
const errorMessage = ref('');
const loading = ref(false);
const router = useRouter();

const handleLogin = async () => {
  loading.value = true;
  errorMessage.value = '';
  try {
    const response = await fetch('https://dns.rea1m.top/api/login', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ username: username.value, password: password.value })
    });
    const data = await response.json();
    if (!response.ok) throw new Error(data.message || '登录失败');
    
    // 保存token和用户信息
    localStorage.setItem('token', data.token);
    localStorage.setItem('userInfo', data.username);
    localStorage.setItem('domain', data.user_domain);
    localStorage.setItem('host', data.host);
    
    // 跳转到DNS日志页面
    router.push('/dns_logs');
  } catch (err) {
    errorMessage.value = err.message;
  } finally {
    loading.value = false;
  }
};

const handleRandomLogin = async () => {
  loading.value = true;
  errorMessage.value = '';
  try {
    const response = await fetch('https://dns.rea1m.top/api/random_id_login');
    const data = await response.json();
    if (!response.ok) throw new Error(data.message || '随机登录失败');
    
    // 保存token和用户信息
    localStorage.setItem('token', data.token);
    localStorage.setItem('userInfo', data.username);
    localStorage.setItem('domain', data.user_domain);
    localStorage.setItem('host', data.host);
    // 跳转到DNS日志页面
    router.push('/dns_logs');
  } catch (err) {
    errorMessage.value = err.message;
  } finally {
    loading.value = false;
  }
};
</script>