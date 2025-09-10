<template>
  <div class="flex flex-col min-h-screen bg-gray-50">
    <!-- 顶部导航栏 -->
    <nav class="bg-gradient-to-r from-indigo-600 via-purple-600 to-pink-500 text-white shadow-lg">
      <div class="container mx-auto px-6 py-4">
        <div class="flex items-center justify-between">
          <!-- Logo 区域 -->
          <div class="flex items-center space-x-3">
            <div class="w-10 h-10 bg-white/20 backdrop-blur-sm rounded-xl flex items-center justify-center">
              <i class="fa fa-server text-xl"></i>
            </div>
            <div>
              <h1 class="text-xl font-bold">DNSLog</h1>
              <p class="text-xs text-white/70">DNS 日志管理</p>
            </div>
          </div>

          <!-- 域名显示 -->
          <div class="flex-1 max-w-xl mx-8">
            <div class="flex items-center justify-center">
              <button @click="copyDomain"
                      class="group flex items-center space-x-2 bg-white/10 hover:bg-white/20 backdrop-blur-sm
                             px-4 py-2 rounded-xl transition-all duration-200">
                <span class="font-mono text-lg">{{ domain }}</span>
                <i class="fa fa-copy text-white/70 group-hover:text-white transition-colors"></i>
              </button>
            </div>
          </div>

          <!-- 用户区域 -->
          <div class="flex items-center space-x-4">
            <div class="flex items-center space-x-3">
              <div class="w-9 h-9 bg-white/20 rounded-full flex items-center justify-center">
                <i class="fa fa-user"></i>
              </div>
              <span class="font-medium hidden sm:block">{{ userName }}</span>
            </div>
            <button @click="handleLogout"
                    class="flex items-center space-x-2 bg-white/10 hover:bg-red-500/80 px-4 py-2 rounded-xl
                           transition-all duration-200">
              <i class="fa fa-sign-out"></i>
              <span class="hidden sm:block">退出</span>
            </button>
          </div>
        </div>
      </div>
    </nav>

    <!-- 主内容区 -->
    <div class="flex flex-1 overflow-hidden">
      <!-- 侧边栏 -->
      <aside class="w-64 flex-shrink-0 bg-white shadow-xl hidden md:block border-r border-gray-100">
        <div class="p-6">
          <div class="space-y-2">
            <router-link to="dns_logs"
                         class="group flex items-center px-4 py-3 rounded-xl transition-all duration-200"
                         :class="$route.path === '/dns_logs'
                                ? 'bg-gradient-to-r from-indigo-500 to-purple-500 text-white shadow-lg shadow-indigo-500/30'
                                : 'text-gray-600 hover:bg-gray-50 hover:text-gray-900'">
              <div class="w-10 h-10 rounded-lg flex items-center justify-center mr-3"
                   :class="$route.path === '/dns_logs' ? 'bg-white/20' : 'bg-gray-100 group-hover:bg-gray-200'">
                <i class="fa fa-server text-lg"
                   :class="$route.path === '/dns_logs' ? 'text-white' : 'text-gray-500'"></i>
              </div>
              <span class="font-medium">DNS 日志</span>
            </router-link>

            <router-link to="rebind"
                         class="group flex items-center px-4 py-3 rounded-xl transition-all duration-200"
                         :class="$route.path === '/rebind'
                                ? 'bg-gradient-to-r from-indigo-500 to-purple-500 text-white shadow-lg shadow-indigo-500/30'
                                : 'text-gray-600 hover:bg-gray-50 hover:text-gray-900'">
              <div class="w-10 h-10 rounded-lg flex items-center justify-center mr-3"
                   :class="$route.path === '/rebind' ? 'bg-white/20' : 'bg-gray-100 group-hover:bg-gray-200'">
                <i class="fa fa-exchange text-lg"
                   :class="$route.path === '/rebind' ? 'text-white' : 'text-gray-500'"></i>
              </div>
              <span class="font-medium">DNS Rebind</span>
            </router-link>
          </div>
        </div>

        <!-- 移动端菜单触发 -->
        <div class="md:hidden fixed bottom-6 right-6 z-50">
          <button class="w-14 h-14 bg-gradient-to-r from-indigo-500 to-purple-500 text-white rounded-full shadow-xl
                         flex items-center justify-center">
            <i class="fa fa-bars text-xl"></i>
          </button>
        </div>
      </aside>

      <!-- 主内容 -->
      <main class="flex-1 overflow-y-auto p-8 bg-gray-50">
        <div class="max-w-7xl mx-auto animate-fade-in">
          <router-view />
        </div>
      </main>
    </div>

    <!-- Toast 提示 -->
    <transition name="toast">
      <div v-if="showToast"
           class="fixed top-8 left-1/2 transform -translate-x-1/2 z-50
                  bg-gradient-to-r from-green-500 to-emerald-500 text-white
                  px-6 py-3 rounded-2xl shadow-2xl flex items-center space-x-3">
        <i class="fa fa-check-circle text-xl"></i>
        <span class="font-medium">{{ toastMessage }}</span>
      </div>
    </transition>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import { useRouter } from 'vue-router';

const userName = ref('');
const domain = ref('');
const router = useRouter();

const showToast = ref(false);
const toastMessage = ref('');

function showToastMsg(msg, duration = 2000) {
  toastMessage.value = msg;
  showToast.value = true;
  setTimeout(() => {
    showToast.value = false;
  }, duration);
}

onMounted(() => {
  const userInfo = localStorage.getItem('userInfo') || '';
  userName.value = userInfo || '';
  domain.value = localStorage.getItem('host') || '';
});

const handleLogout = () => {
  localStorage.removeItem('token');
  localStorage.removeItem('userInfo');
  localStorage.removeItem('domain');
  localStorage.removeItem('host');
  router.push('/login');
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

<style scoped>
.toast-enter-active,
.toast-leave-active {
  transition: all 0.3s ease;
}

.toast-enter-from,
.toast-leave-to {
  opacity: 0;
  transform: translate(-50%, -20px);
}
</style>
