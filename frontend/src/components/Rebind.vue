<template>
  <div class="w-full">
    <!-- 页面标题 -->
    <div class="mb-8">
      <h1 class="text-3xl font-bold text-gray-800 flex items-center">
        <i class="fa fa-exchange text-purple-500 mr-3"></i>
        DNS Rebind
      </h1>
      <p class="text-gray-500 mt-1">管理 DNS 重绑定规则</p>
    </div>

    <div class="grid grid-cols-1 lg:grid-cols-2 gap-8">
      <!-- 生成 Rebind 规则 -->
      <div class="bg-white rounded-2xl shadow-xl border border-gray-100 overflow-hidden">
        <div class="px-6 py-4 bg-gradient-to-r from-purple-500 to-indigo-500">
          <h2 class="text-lg font-semibold text-white flex items-center">
            <i class="fa fa-plus-circle mr-2"></i>
            生成 Rebind 规则
          </h2>
        </div>

        <div class="p-6">
          <form @submit.prevent="handleGenerate" class="space-y-5">
            <!-- 提示信息 -->
            <transition name="fade">
              <div v-if="errorMessage"
                   class="p-4 bg-red-50 border border-red-200 rounded-xl flex items-start">
                <i class="fa fa-exclamation-circle text-red-500 mt-0.5 mr-3 text-lg"></i>
                <span class="text-red-700 text-sm">{{ errorMessage }}</span>
              </div>
            </transition>

            <transition name="fade">
              <div v-if="successMessage"
                   class="p-4 bg-green-50 border border-green-200 rounded-xl flex items-start">
                <i class="fa fa-check-circle text-green-500 mt-0.5 mr-3 text-lg"></i>
                <span class="text-green-700 text-sm">{{ successMessage }}</span>
              </div>
            </transition>

            <!-- First IP -->
            <div class="space-y-2">
              <label for="firstIp" class="block text-sm font-semibold text-gray-700 ml-1">
                <i class="fa fa-location-arrow mr-1 text-purple-500"></i>
                First IP
              </label>
              <div class="relative">
                <div class="absolute inset-y-0 left-0 pl-4 flex items-center pointer-events-none">
                  <span class="text-gray-400 font-mono text-sm">1st</span>
                </div>
                <input type="text" id="firstIp" v-model="firstIp" placeholder="例如: 192.168.1.1"
                       class="w-full pl-12 pr-4 py-3 bg-gray-50 border-2 border-gray-200 rounded-xl
                              focus:outline-none focus:border-purple-500 focus:bg-white
                              transition-all duration-200 font-mono text-gray-700">
              </div>
              <p class="text-xs text-gray-500">首次解析返回的 IP 地址</p>
            </div>

            <!-- Second IP -->
            <div class="space-y-2">
              <label for="secondIp" class="block text-sm font-semibold text-gray-700 ml-1">
                <i class="fa fa-location-arrow mr-1 text-indigo-500"></i>
                Second IP
              </label>
              <div class="relative">
                <div class="absolute inset-y-0 left-0 pl-4 flex items-center pointer-events-none">
                  <span class="text-gray-400 font-mono text-sm">2nd</span>
                </div>
                <input type="text" id="secondIp" v-model="secondIp" placeholder="例如: 127.0.0.1"
                       class="w-full pl-12 pr-4 py-3 bg-gray-50 border-2 border-gray-200 rounded-xl
                              focus:outline-none focus:border-indigo-500 focus:bg-white
                              transition-all duration-200 font-mono text-gray-700">
              </div>
              <p class="text-xs text-gray-500">后续解析返回的 IP 地址</p>
            </div>

            <!-- 生成按钮 -->
            <button type="submit" :disabled="loading"
                    class="w-full py-3 px-6 bg-gradient-to-r from-purple-500 to-indigo-500 text-white font-semibold rounded-xl
                           hover:from-purple-600 hover:to-indigo-600 focus:outline-none focus:ring-4 focus:ring-purple-500/20
                           disabled:opacity-70 disabled:cursor-not-allowed transition-all duration-200
                           flex items-center justify-center space-x-2">
              <i v-if="loading" class="fa fa-spinner fa-spin"></i>
              <span>{{ loading ? '生成中...' : '生成 Rebind 域名' }}</span>
            </button>
          </form>

          <!-- 生成结果 -->
          <transition name="fade">
            <div id="resultContainer" class="mt-6" v-if="rebindUrl">
              <div class="p-5 bg-gradient-to-r from-green-50 to-emerald-50 border border-green-200 rounded-xl">
                <div class="flex items-center mb-3">
                  <div class="w-10 h-10 bg-green-100 rounded-full flex items-center justify-center mr-3">
                    <i class="fa fa-check text-green-600"></i>
                  </div>
                  <h3 class="font-bold text-green-800">生成成功！</h3>
                </div>
                <div class="bg-white rounded-lg p-3 border border-green-200">
                  <p class="text-xs text-green-600 font-medium mb-1">Rebind 域名</p>
                  <p class="font-mono text-sm text-green-800 break-all">{{ rebindUrl }}</p>
                </div>
              </div>
            </div>
          </transition>
        </div>
      </div>

      <!-- 活跃 Rebind 规则列表 -->
      <div class="bg-white rounded-2xl shadow-xl border border-gray-100 overflow-hidden">
        <div class="px-6 py-4 bg-gradient-to-r from-gray-50 to-gray-100 border-b border-gray-100">
          <h2 class="text-lg font-semibold text-gray-800 flex items-center">
            <i class="fa fa-list-ul mr-2 text-gray-500"></i>
            活跃规则
            <span class="ml-2 px-2.5 py-0.5 bg-indigo-100 text-indigo-600 text-xs font-medium rounded-full">
              {{ rules.length }}
            </span>
          </h2>
        </div>

        <div class="p-6">
          <div class="space-y-4 max-h-[500px] overflow-y-auto pr-2">
            <!-- 加载状态 -->
            <div v-if="loading" class="text-center py-12">
              <div class="flex flex-col items-center space-y-3">
                <i class="fa fa-spinner fa-spin text-3xl text-indigo-500"></i>
                <span class="text-gray-500">加载中...</span>
              </div>
            </div>

            <!-- 空状态 -->
            <div v-else-if="rules.length === 0" class="text-center py-12">
              <div class="flex flex-col items-center space-y-3">
                <div class="w-16 h-16 bg-gray-100 rounded-full flex items-center justify-center">
                  <i class="fa fa-inbox text-2xl text-gray-400"></i>
                </div>
                <span class="text-gray-500 font-medium">暂无活跃规则</span>
                <span class="text-gray-400 text-sm">创建的 Rebind 规则将显示在这里</span>
              </div>
            </div>

            <!-- 规则列表 -->
            <div v-else>
              <div v-for="rule in rules" :key="rule.id"
                   class="group border-2 border-gray-100 rounded-xl p-4 hover:border-gray-200 hover:bg-gray-50
                          transition-all duration-200">
                <div class="flex items-start justify-between">
                  <div class="flex-1 min-w-0">
                    <div class="flex items-center space-x-2 mb-2">
                      <div class="w-8 h-8 bg-gradient-to-br from-purple-500 to-indigo-500 rounded-lg flex items-center justify-center">
                        <i class="fa fa-link text-white text-sm"></i>
                      </div>
                      <p class="font-mono text-sm font-medium text-gray-800 break-all">{{ rule.domain }}</p>
                    </div>
                    <div class="flex flex-wrap gap-2">
                      <div class="flex items-center space-x-1.5 bg-green-50 px-2.5 py-1 rounded-lg">
                        <span class="w-2 h-2 bg-green-500 rounded-full"></span>
                        <span class="text-xs font-mono text-green-700">{{ rule.first_ip }}</span>
                      </div>
                      <div class="flex items-center space-x-1">
                        <i class="fa fa-arrow-right text-gray-400 text-xs"></i>
                      </div>
                      <div class="flex items-center space-x-1.5 bg-blue-50 px-2.5 py-1 rounded-lg">
                        <span class="w-2 h-2 bg-blue-500 rounded-full"></span>
                        <span class="text-xs font-mono text-blue-700">{{ rule.second_ip }}</span>
                      </div>
                    </div>
                  </div>
                  <button @click="handleDeleteRule(rule.id)"
                          class="ml-4 p-2 text-gray-400 hover:text-red-500 hover:bg-red-50 rounded-lg
                                 transition-all duration-200 opacity-0 group-hover:opacity-100">
                    <i class="fa fa-trash"></i>
                  </button>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import { getApiUrl } from '../config.js';

const firstIp = ref('');
const secondIp = ref('');
const rules = ref([]);
const rebindUrl = ref('');
const loading = ref(false);
const errorMessage = ref('');
const successMessage = ref('');
const router = useRouter();

const fetchRebindRules = async () => {
  loading.value = true;
  try {
    const response = await fetch(getApiUrl('/api/rebind/list'), {
      method: 'GET',
      headers: { 'Authorization': `Bearer ${localStorage.getItem('token')}` }
    });

    if (!response.ok) {
      if (response.status === 401) {
        localStorage.removeItem('token');
        router.push('/login');
        return;
      }
      throw new Error('获取规则失败');
    }

    const data = await response.json();
    rules.value = data.rebind_list || [];
  } catch (err) {
    console.error('获取Rebind规则错误:', err);
  } finally {
    loading.value = false;
  }
};

const handleGenerate = async () => {
  loading.value = true;
  errorMessage.value = '';
  successMessage.value = '';
  rebindUrl.value = '';

  try {
    const response = await fetch(getApiUrl('/api/rebind/gen'), {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${localStorage.getItem('token')}`
      },
      body: JSON.stringify({
        first_ip: firstIp.value,
        second_ip: secondIp.value
      })
    });

    const data = await response.json();
    if (!response.ok) throw new Error(data.error || data.message || '生成规则失败');

    rebindUrl.value = data.url || data.rebind_domain;
    successMessage.value = 'Rebind规则生成成功';
    await fetchRebindRules();
  } catch (err) {
    console.error('生成Rebind规则错误:', err);
    errorMessage.value = err.message;
  } finally {
    loading.value = false;
  }
};

const handleDeleteRule = async (id) => {
  if (!confirm('确定要删除此规则吗？')) return;

  try {
    const response = await fetch(getApiUrl('/api/rebind/delete'), {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${localStorage.getItem('token')}`
      },
      body: JSON.stringify({
        id: id
      })
    });

    if (!response.ok) throw new Error('删除规则失败');
    await fetchRebindRules();
  } catch (err) {
    console.error('删除Rebind规则错误:', err);
    alert('删除规则失败: ' + err.message);
  }
};

onMounted(() => {
  fetchRebindRules();
});
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
