<template>
  <div class="min-h-screen flex items-center justify-center p-4"
       style="background: linear-gradient(135deg, #667eea 0%, #764ba2 50%, #f093fb 100%);">
    <div class="absolute inset-0 overflow-hidden">
      <div class="absolute -top-40 -right-32 w-80 h-80 bg-white/10 rounded-full blur-3xl"></div>
      <div class="absolute -bottom-40 -left-32 w-96 h-96 bg-purple-300/20 rounded-full blur-3xl"></div>
    </div>

    <div class="relative w-full max-w-md animate-fade-in">
      <!-- Logo 区域 -->
      <div class="text-center mb-8">
        <div class="inline-flex items-center justify-center w-20 h-20 bg-white/20 backdrop-blur-sm rounded-2xl mb-4 shadow-lg">
          <i class="fa fa-server text-4xl text-white"></i>
        </div>
        <h1 class="text-3xl font-bold text-white mb-2">DNSLog</h1>
        <p class="text-white/80">DNS 日志管理平台</p>
      </div>

      <!-- 登录卡片 -->
      <div class="glass rounded-3xl p-8 card-shadow">
        <h2 class="text-2xl font-bold text-gray-800 mb-6 text-center">欢迎回来</h2>

        <form @submit.prevent="handleLogin" class="space-y-5">
          <!-- 错误提示 -->
          <transition name="fade">
            <div v-if="errorMessage"
                 class="p-4 bg-red-50 border border-red-200 rounded-xl text-red-700 text-sm flex items-center">
              <i class="fa fa-exclamation-circle mr-3 text-lg"></i>
              <span>{{ errorMessage }}</span>
            </div>
          </transition>

          <!-- 用户名 -->
          <div class="space-y-2">
            <label for="username" class="block text-sm font-semibold text-gray-700 ml-1">
              用户名
            </label>
            <div class="relative">
              <div class="absolute inset-y-0 left-0 pl-4 flex items-center pointer-events-none">
                <i class="fa fa-user text-gray-400 text-lg"></i>
              </div>
              <input type="text" id="username" v-model="username" required
                     class="w-full pl-12 pr-4 py-4 bg-gray-50 border-2 border-gray-200 rounded-xl
                            focus:outline-none focus:border-indigo-500 focus:bg-white
                            transition-all duration-200 text-gray-800 placeholder-gray-400"
                     placeholder="请输入用户名">
            </div>
          </div>

          <!-- 密码 -->
          <div class="space-y-2">
            <label for="password" class="block text-sm font-semibold text-gray-700 ml-1">
              密码
            </label>
            <div class="relative">
              <div class="absolute inset-y-0 left-0 pl-4 flex items-center pointer-events-none">
                <i class="fa fa-lock text-gray-400 text-lg"></i>
              </div>
              <input :type="showPassword ? 'text' : 'password'" id="password" v-model="password" required
                     class="w-full pl-12 pr-12 py-4 bg-gray-50 border-2 border-gray-200 rounded-xl
                            focus:outline-none focus:border-indigo-500 focus:bg-white
                            transition-all duration-200 text-gray-800 placeholder-gray-400"
                     placeholder="请输入密码">
              <button type="button" @click="showPassword = !showPassword"
                      class="absolute inset-y-0 right-0 pr-4 flex items-center text-gray-400 hover:text-gray-600 transition-colors">
                <i :class="showPassword ? 'fa fa-eye-slash' : 'fa fa-eye'" class="text-lg"></i>
              </button>
            </div>
          </div>

          <!-- 登录按钮 -->
          <button type="submit" :disabled="loading"
                  class="w-full py-4 px-6 btn-gradient text-white font-semibold rounded-xl
                         disabled:opacity-70 disabled:cursor-not-allowed disabled:transform-none
                         flex items-center justify-center space-x-2">
            <i v-if="loading" class="fa fa-spinner fa-spin"></i>
            <span>{{ loading ? '登录中...' : '登 录' }}</span>
          </button>
        </form>
      </div>

      <!-- 底部版权 -->
      <p class="text-center text-white/60 text-sm mt-8">
        © 2024 DNSLog. All rights reserved.
      </p>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue';
import { useRouter } from 'vue-router';
import { getApiUrl } from '../config.js';

const username = ref('');
const password = ref('');
const showPassword = ref(false);
const errorMessage = ref('');
const loading = ref(false);
const router = useRouter();

const handleLogin = async () => {
  loading.value = true;
  errorMessage.value = '';
  try {
    const response = await fetch(getApiUrl('/api/login'), {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ username: username.value, password: password.value })
    });
    const data = await response.json();
    if (!response.ok) throw new Error(data.error || '登录失败');

    localStorage.setItem('token', data.token);
    localStorage.setItem('userInfo', data.username);
    localStorage.setItem('domain', data.user_domain);
    localStorage.setItem('host', data.host);

    router.push('/dns_logs');
  } catch (err) {
    errorMessage.value = err.message;
  } finally {
    loading.value = false;
  }
};
</script>

<style scoped>
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s ease, transform 0.3s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
  transform: translateY(-10px);
}
</style>
