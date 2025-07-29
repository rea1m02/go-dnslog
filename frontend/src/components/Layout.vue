<template>
  <div class="flex flex-col min-h-screen bg-gray-100">
    <!-- 顶部导航栏 -->
    <nav class="bg-indigo-600 text-white shadow-md">
      <div class="container mx-auto px-0 py-3 flex justify-between items-center">
        <div class="flex items-center space-x-2">
          <i class="fa fa-server text-2xl"></i>
          <span class="text-xl font-bold">DNSLog</span>
        </div>

        <div class="flex items-center space-x-3">
          <span class="text-xl font-bold" @click="copyDomain()">{{ domain }}</span>
        </div>

     
        <div class="flex items-center space-x-4">
          <span v-if="userName" class="md:inline-block">{{ userName }}</span>
          <button @click="handleLogout" class="hover:text-red-200 transition-colors">
            <i class="fa fa-sign-out mr-1"></i>退出
          </button>
        </div>
      </div>
    </nav>

    <!-- 主内容区 -->
    <div class="flex flex-1 overflow-hidden">
      <!-- 侧边栏 -->
      <aside class="bg-white shadow-md w-64 flex-shrink-0 hidden md:block">
        <div class="p-4">
          <div class="space-y-1">
            <!-- <router-link to="/" class="flex items-center px-4 py-2 rounded-md transition-colors" :class="$route.path === '/' ? 'text-indigo-600 bg-indigo-50' : 'text-gray-700 hover:bg-gray-100'">
              <i class="fa fa-tachometer w-5 h-5 mr-3"></i>
              <span>仪表盘</span>
            </router-link> -->
            <router-link
              to="dns_logs"
              class="flex items-center px-4 py-2 rounded-md transition-colors"
              :class="$route.path === '/dns_logs' ? 'text-indigo-600 bg-indigo-50' : 'text-gray-700 hover:bg-gray-100'"
            >
              <i class="fa fa-server w-5 h-5 mr-3"></i>
              <span>DNS日志</span>
            </router-link>
            <!-- <router-link to="web_logs" class="flex items-center px-4 py-2 rounded-md transition-colors" :class="$route.path === '/web_logs' ? 'text-indigo-600 bg-indigo-50' : 'text-gray-700 hover:bg-gray-100'" @click="handleNavClick('web_logs')">
              <i class="fa fa-globe w-5 h-5 mr-3"></i>
              <span>Web日志</span>
            </router-link> -->
            <router-link 
              to="rebind"
              class="flex items-center px-4 py-2 rounded-md transition-colors"
              :class="$route.path === '/rebind' ? 'text-indigo-600 bg-indigo-50' : 'text-gray-700 hover:bg-gray-100'"
            >
              <i class="fa fa-exchange w-5 h-5 mr-3"></i>
              <span>DNS Rebind</span>
            </router-link>
            <!-- <router-link 
              to="settings"
              class="flex items-center px-4 py-2 rounded-md transition-colors"
              :class="$route.path === '/settings' ? 'text-indigo-600 bg-indigo-50' : 'text-gray-700 hover:bg-gray-100'"
            >
              <i class="fa fa-cog w-5 h-5 mr-3"></i>
              <span>设置</span>
            </router-link> -->
          </div>
        </div>
      </aside>

      <!-- 主内容 -->
      <main class="flex-1 overflow-y-auto p-6">
        <router-view />
      </main>
    </div>
    <!-- Toast 提示 -->
    <transition name="fade">
      <div
        v-if="showToast"
        class="fixed top-6 left-1/2 transform -translate-x-1/2 bg-green-600 text-white px-4 py-2 rounded shadow-lg z-50"
      >
        {{ toastMessage }}
      </div>
    </transition>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import { useRouter, useRoute } from 'vue-router';

const userName = ref('');
const domain = ref('');
const router = useRouter();
const route = useRoute();

const showToast = ref(false);
const toastMessage = ref('');

function showToastMsg(msg, duration = 1500) {
  toastMessage.value = msg;
  showToast.value = true;
  setTimeout(() => {
    showToast.value = false;
  }, duration);
}

onMounted(() => {
  // 从localStorage获取用户名
  const userInfo = localStorage.getItem('userInfo') || '';
  userName.value = userInfo || '';
  domain.value = localStorage.getItem('domain') + '.' + localStorage.getItem('host') || '';
  // 调试信息
  console.log('顶栏用户名:', userName.value);
  console.log('Current route:', route.path);
  console.log('Router instance:', router);
});

const handleLogout = () => {
  // 清除本地存储
  localStorage.removeItem('token');
  localStorage.removeItem('userInfo');
  localStorage.removeItem('domain');
  localStorage.removeItem('host');
  // 重定向到登录页
  router.push('/login');
};

const handleNavClick = (path) => {
  // 实现路由跳转逻辑
  console.log(`Navigating to: ${path}`);
  router.push(path);
};


const copyDomain = async () => {
  try {
    await navigator.clipboard.writeText(domain.value);
    showToastMsg('域名已复制到剪贴板');
  } catch (err) {
    showToastMsg('复制失败，请手动复制', 2000);
  }
};



</script>

<style>
.fade-enter-active, .fade-leave-active {
  transition: opacity 0.3s;
}
.fade-enter-from, .fade-leave-to {
  opacity: 0;
}
</style>