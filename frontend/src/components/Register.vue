<template>
  <div class="bg-gray-100 min-h-screen flex items-center justify-center">
    <div class="bg-white p-8 rounded-lg shadow-md w-full max-w-md">
      <h2 class="text-2xl font-bold mb-6 text-center">DNSLog 注册</h2>
      <form @submit.prevent="handleRegister" class="space-y-4">
        <div v-if="errorMessage" class="p-3 bg-red-50 border border-red-200 rounded-md text-red-700 text-sm">
          {{ errorMessage }}
        </div>
        <div v-if="successMessage" class="p-3 bg-green-50 border border-green-200 rounded-md text-green-700 text-sm">
          {{ successMessage }}
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
          <label for="email" class="block text-sm font-medium text-gray-700">邮箱</label>
          <div class="mt-1 relative rounded-md shadow-sm">
            <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
              <i class="fa fa-envelope text-gray-400"></i>
            </div>
            <input type="email" id="email" v-model="email" required
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
          <label for="confirmPassword" class="block text-sm font-medium text-gray-700">确认密码</label>
          <div class="mt-1 relative rounded-md shadow-sm">
            <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
              <i class="fa fa-lock text-gray-400"></i>
            </div>
            <input type="password" id="confirmPassword" v-model="confirmPassword" required
              class="focus:ring-indigo-500 focus:border-indigo-500 block w-full pl-10 pr-3 py-2 border border-gray-300 rounded-md">
            <div v-if="passwordMismatch" class="text-red-500 text-xs mt-1">
              密码不匹配
            </div>
          </div>
        </div>
        <div>
          <button type="submit" :disabled="loading"
            class="w-full flex justify-center py-2 px-4 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 disabled:opacity-70">
            <i v-if="loading" class="fa fa-spinner fa-spin mr-2"></i>
            注册
          </button>
        </div>
        <div class="text-center">
          <router-link to="/login" class="text-sm text-indigo-600 hover:text-indigo-500">
            已有账号？登录
          </router-link>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup>
import { ref, watch } from 'vue';
import { useRouter } from 'vue-router';

const username = ref('');
const email = ref('');
const password = ref('');
const confirmPassword = ref('');
const errorMessage = ref('');
const successMessage = ref('');
const loading = ref(false);
const passwordMismatch = ref(false);
const router = useRouter();

// 监听密码变化，检查密码是否匹配
watch([password, confirmPassword], () => {
  passwordMismatch.value = password.value !== confirmPassword.value && confirmPassword.value !== '';
});

const handleRegister = async () => {
  // 表单验证
  if (password.value !== confirmPassword.value) {
    passwordMismatch.value = true;
    return;
  }

  loading.value = true;
  errorMessage.value = '';
  successMessage.value = '';

  try {
    const response = await fetch('https://dns.rea1m.top/api/register', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        username: username.value,
        email: email.value,
        password: password.value
      })
    });

    const data = await response.json();
    if (!response.ok) throw new Error(data.message || '注册失败');

    // 注册成功，显示提示并跳转到登录页
    successMessage.value = '注册成功，请登录';
    setTimeout(() => {
      router.push('/login');
    }, 2000);
  } catch (err) {
    errorMessage.value = err.message;
  } finally {
    loading.value = false;
  }
};
</script>