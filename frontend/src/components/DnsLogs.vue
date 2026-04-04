<template>
  <div class="w-full">
    <!-- 页面标题和操作栏 -->
    <div class="mb-8">
      <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
        <div>
          <h1 class="text-3xl font-bold text-gray-800 flex items-center">
            <i class="fa fa-server text-indigo-500 mr-3"></i>
            DNS 日志
          </h1>
          <p class="text-gray-500 mt-1">查看和管理 DNS 解析记录</p>
        </div>
        <div class="flex items-center gap-3">
          <!-- 实时刷新按钮 -->
          <button @click="toggleAutoRefresh"
                  class="px-4 py-2.5 rounded-xl font-medium transition-all duration-200 flex items-center gap-2"
                  :class="autoRefresh
                    ? 'bg-gradient-to-r from-green-500 to-emerald-500 text-white shadow-lg shadow-green-500/25'
                    : 'bg-white border-2 border-gray-200 text-gray-600 hover:border-gray-300 hover:bg-gray-50'">
            <i :class="autoRefresh ? 'fa fa-pause' : 'fa fa-play'"></i>
            <span>{{ autoRefresh ? '停止刷新' : '实时刷新' }}</span>
            <span v-if="autoRefresh" class="text-xs opacity-75">({{ countdown }}s)</span>
          </button>
        </div>
      </div>
    </div>

    <!-- 监听域名管理卡片 -->
    <div class="bg-white rounded-2xl shadow-xl border border-gray-100 overflow-hidden mb-8">
      <div class="px-6 py-4 bg-gradient-to-r from-indigo-500 to-purple-500">
        <div class="flex items-center justify-between">
          <h2 class="text-lg font-semibold text-white flex items-center">
            <i class="fa fa-eye mr-2"></i>
            监听域名管理
          </h2>
          <span class="text-sm text-white/80">只有添加的域名才会记录 DNS 查询</span>
        </div>
      </div>

      <div class="p-6">
        <!-- 添加域名表单 -->
        <form @submit.prevent="handleAddWatchDomain" class="flex flex-col sm:flex-row gap-4 mb-6">
          <div class="flex-1">
            <div class="relative">
              <div class="absolute inset-y-0 left-0 pl-4 flex items-center pointer-events-none">
                <i class="fa fa-globe text-gray-400"></i>
              </div>
              <input type="text" v-model="newWatchDomain" required
                     :placeholder="'例如: test.' + baseDomain"
                     class="w-full pl-11 pr-4 py-3 bg-gray-50 border-2 border-gray-200 rounded-xl
                            focus:outline-none focus:border-indigo-500 focus:bg-white
                            transition-all duration-200 font-mono">
            </div>
          </div>
          <button type="submit" :disabled="addingDomain"
                  class="px-6 py-3 bg-gradient-to-r from-indigo-500 to-purple-500 text-white font-semibold rounded-xl
                         hover:from-indigo-600 hover:to-purple-600 focus:outline-none focus:ring-4 focus:ring-indigo-500/20
                         disabled:opacity-70 disabled:cursor-not-allowed transition-all duration-200
                         flex items-center space-x-2">
            <i v-if="addingDomain" class="fa fa-spinner fa-spin"></i>
            <i v-else class="fa fa-plus"></i>
            <span>{{ addingDomain ? '添加中...' : '添加监听' }}</span>
          </button>
          <button type="button" @click="handleRefreshWatchDomains"
                  class="px-4 py-3 bg-gray-100 text-gray-700 font-medium rounded-xl
                         hover:bg-gray-200 transition-all duration-200">
            <i class="fa fa-refresh"></i>
          </button>
        </form>

        <!-- 监听域名列表 -->
        <div class="space-y-3">
          <div v-if="loadingWatchDomains" class="text-center py-8">
            <div class="flex flex-col items-center space-y-3">
              <i class="fa fa-spinner fa-spin text-2xl text-indigo-500"></i>
              <span class="text-gray-500">加载中...</span>
            </div>
          </div>
          <div v-else-if="watchDomains.length === 0" class="text-center py-8">
            <div class="flex flex-col items-center space-y-3">
              <div class="w-16 h-16 bg-gray-100 rounded-full flex items-center justify-center">
                <i class="fa fa-inbox text-2xl text-gray-400"></i>
              </div>
              <span class="text-gray-500 font-medium">暂无监听域名</span>
              <span class="text-gray-400 text-sm">添加域名后才会记录 DNS 查询</span>
            </div>
          </div>
          <div v-else class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-3">
            <div v-for="domain in watchDomains" :key="domain.id"
                 class="group border-2 border-gray-100 rounded-xl p-4 hover:border-gray-200 hover:bg-gray-50
                        transition-all duration-200">
              <div class="flex items-start justify-between">
                <div class="flex-1 min-w-0">
                  <div class="flex items-center space-x-2">
                    <div class="w-8 h-8 bg-gradient-to-br from-green-500 to-emerald-500 rounded-lg flex items-center justify-center">
                      <i class="fa fa-link text-white text-sm"></i>
                    </div>
                    <p class="font-mono text-sm font-medium text-gray-800 break-all">{{ domain.domain }}</p>
                  </div>
                  <p class="text-xs text-gray-400 mt-2">
                    添加于 {{ formatTime(domain.created_at) }}
                  </p>
                </div>
                <button @click="handleDeleteWatchDomain(domain.id)"
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

    <!-- DNS 日志表格卡片 -->
    <div class="bg-white rounded-2xl shadow-xl border border-gray-100 overflow-hidden">
      <!-- 表格头部搜索和操作栏 -->
      <div class="px-6 py-4 bg-gradient-to-r from-gray-50 to-gray-100 border-b border-gray-100">
        <div class="flex flex-col lg:flex-row lg:items-center lg:justify-between gap-4">
          <div class="flex items-center space-x-4">
            <div class="flex items-center space-x-2">
              <div class="w-3 h-3 bg-green-500 rounded-full" :class="autoRefresh ? 'animate-pulse' : ''"></div>
              <span class="text-sm text-gray-600">
                {{ autoRefresh ? '自动刷新中' : '实时更新' }}
              </span>
            </div>
            <span class="text-sm text-gray-500">|</span>
            <span class="text-sm font-medium text-gray-700">
              共 <span class="text-indigo-600">{{ totalCount }}</span> 条记录
            </span>
          </div>
          <div class="flex items-center gap-3">
            <!-- 搜索框 -->
            <div class="relative flex-1 lg:w-72">
              <i class="fa fa-search absolute left-4 top-1/2 -translate-y-1/2 text-gray-400"></i>
              <input type="text" v-model="searchQuery" placeholder="搜索域名或IP..."
                     class="w-full pl-11 pr-12 py-2.5 bg-white border-2 border-gray-200 rounded-xl
                            focus:outline-none focus:border-indigo-500 focus:ring-4 focus:ring-indigo-500/10
                            transition-all duration-200"
                     @keyup.enter="handleSearch">
            </div>
            <!-- 搜索按钮 -->
            <button @click="handleSearch"
                    class="px-5 py-2.5 bg-gradient-to-r from-indigo-500 to-purple-500 text-white font-medium rounded-xl
                           hover:from-indigo-600 hover:to-purple-600 focus:outline-none focus:ring-4 focus:ring-indigo-500/20
                           transition-all duration-200 flex items-center gap-2">
              <i class="fa fa-search"></i>
              <span>搜索</span>
            </button>
            <!-- 清空日志按钮 -->
            <button @click="handleClearLogs"
                    class="px-4 py-2.5 bg-gradient-to-r from-red-500 to-rose-500 text-white font-medium rounded-xl
                           hover:from-red-600 hover:to-rose-600 focus:outline-none focus:ring-4 focus:ring-red-500/20
                           transition-all duration-200 flex items-center gap-2 disabled:opacity-50"
                    :disabled="loading">
              <i class="fa fa-trash"></i>
              <span class="hidden sm:inline">清空</span>
            </button>
          </div>
        </div>
      </div>

      <!-- 表格内容 -->
      <div class="overflow-x-auto">
        <table class="min-w-full">
          <thead>
            <tr class="bg-gray-50 border-b border-gray-100">
              <th class="px-6 py-4 text-left text-xs font-semibold text-gray-500 uppercase tracking-wider">
                ID
              </th>
              <th class="px-6 py-4 text-left text-xs font-semibold text-gray-500 uppercase tracking-wider">
                域名
              </th>
              <th class="px-6 py-4 text-left text-xs font-semibold text-gray-500 uppercase tracking-wider">
                源 IP
              </th>
              <th class="px-6 py-4 text-left text-xs font-semibold text-gray-500 uppercase tracking-wider">
                类型
              </th>
              <th class="px-6 py-4 text-left text-xs font-semibold text-gray-500 uppercase tracking-wider">
                时间
              </th>
              <th class="px-6 py-4 text-left text-xs font-semibold text-gray-500 uppercase tracking-wider">
                操作
              </th>
            </tr>
          </thead>
          <tbody class="divide-y divide-gray-100">
            <tr v-if="loading">
              <td colspan="6" class="px-6 py-12 text-center">
                <div class="flex flex-col items-center space-y-3">
                  <i class="fa fa-spinner fa-spin text-3xl text-indigo-500"></i>
                  <span class="text-gray-500">加载中...</span>
                </div>
              </td>
            </tr>
            <tr v-else-if="logs.length === 0">
              <td colspan="6" class="px-6 py-12 text-center">
                <div class="flex flex-col items-center space-y-3">
                  <div class="w-16 h-16 bg-gray-100 rounded-full flex items-center justify-center">
                    <i class="fa fa-inbox text-2xl text-gray-400"></i>
                  </div>
                  <span class="text-gray-500 font-medium">暂无日志记录</span>
                  <span class="text-gray-400 text-sm">添加监听域名后，DNS 查询记录将显示在这里</span>
                </div>
              </td>
            </tr>
            <tr v-for="log in logs" :key="log.id"
                class="hover:bg-gray-50 transition-colors duration-150">
              <td class="px-6 py-4 whitespace-nowrap">
                <span class="text-sm font-medium text-gray-400">#{{ log.id }}</span>
              </td>
              <td class="px-6 py-4 whitespace-nowrap">
                <div class="flex items-center space-x-3">
                  <div class="w-8 h-8 bg-indigo-100 rounded-lg flex items-center justify-center">
                    <i class="fa fa-globe text-indigo-600 text-sm"></i>
                  </div>
                  <span class="text-sm font-medium text-gray-800 font-mono">{{ log.domain }}</span>
                </div>
              </td>
              <td class="px-6 py-4 whitespace-nowrap">
                <div class="inline-flex items-center px-2.5 py-1 rounded-lg bg-gray-100">
                  <i class="fa fa-map-marker text-gray-400 mr-1.5 text-xs"></i>
                  <span class="text-sm font-mono text-gray-700">{{ log.client_ip }}</span>
                </div>
              </td>
              <td class="px-6 py-4 whitespace-nowrap">
                <span class="inline-flex items-center px-3 py-1 rounded-full text-xs font-semibold"
                      :class="getQueryTypeClass(log.query_type)">
                  {{ log.query_type }}
                </span>
              </td>
              <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                <div class="flex items-center space-x-2">
                  <i class="fa fa-clock-o text-gray-400"></i>
                  <span>{{ formatTime(log.timestamp) }}</span>
                </div>
              </td>
              <td class="px-6 py-4 whitespace-nowrap">
                <button @click="handleDelete(log.id)"
                        class="p-2 text-gray-400 hover:text-red-500 hover:bg-red-50 rounded-lg transition-all duration-200">
                  <i class="fa fa-trash"></i>
                </button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <!-- 分页 -->
      <div class="px-6 py-4 bg-gray-50 border-t border-gray-100 flex flex-col sm:flex-row items-center justify-between gap-4">
        <div class="text-sm text-gray-600">
          显示第 <span class="font-medium text-gray-800">{{ (page - 1) * pageSize + 1 }}</span> -
          <span class="font-medium text-gray-800">{{ Math.min(page * pageSize, totalCount) }}</span> 条
        </div>
        <div class="flex items-center space-x-2">
          <button @click="prevPage"
                  class="px-4 py-2 bg-white border-2 border-gray-200 rounded-lg text-sm font-medium text-gray-600
                         hover:bg-gray-50 hover:text-gray-900 hover:border-gray-300
                         disabled:opacity-50 disabled:cursor-not-allowed disabled:hover:bg-white
                         transition-all duration-200"
                  :disabled="page <= 1 || loading">
            <i class="fa fa-chevron-left mr-1"></i>
            上一页
          </button>
          <div class="flex items-center space-x-1">
            <span class="px-3 py-2 bg-indigo-50 text-indigo-600 rounded-lg text-sm font-medium">
              {{ page }}
            </span>
            <span class="text-gray-400 text-sm">/</span>
            <span class="text-gray-600 text-sm">{{ totalPages || 1 }}</span>
          </div>
          <button @click="nextPage"
                  class="px-4 py-2 bg-white border-2 border-gray-200 rounded-lg text-sm font-medium text-gray-600
                         hover:bg-gray-50 hover:text-gray-900 hover:border-gray-300
                         disabled:opacity-50 disabled:cursor-not-allowed disabled:hover:bg-white
                         transition-all duration-200"
                  :disabled="page >= totalPages || loading">
            下一页
            <i class="fa fa-chevron-right ml-1"></i>
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue';
import { useRouter } from 'vue-router';
import { getApiUrl } from '../config.js';

const logs = ref([]);
const loading = ref(false);
const totalCount = ref(0);
const page = ref(1);
const pageSize = ref(10);
const totalPages = ref(0);
const searchQuery = ref('');
const router = useRouter();

// 自动刷新相关
const autoRefresh = ref(false);
const countdown = ref(1);
let refreshTimer = null;

// 监听域名相关
const watchDomains = ref([]);
const loadingWatchDomains = ref(false);
const addingDomain = ref(false);
const newWatchDomain = ref('');
const baseDomain = ref('');

const formatTime = (timeStr) => {
  const date = new Date(timeStr);
  return date.toLocaleString('zh-CN', { hour12: false });
};

const getQueryTypeClass = (type) => {
  const classes = {
    'A': 'bg-green-100 text-green-700',
    'AAAA': 'bg-blue-100 text-blue-700',
    'CNAME': 'bg-purple-100 text-purple-700',
    'MX': 'bg-orange-100 text-orange-700',
    'NS': 'bg-pink-100 text-pink-700',
    'TXT': 'bg-yellow-100 text-yellow-700',
  };
  return classes[type] || 'bg-gray-100 text-gray-700';
};

// 获取DNS日志
const fetchDnsLogs = async () => {
  loading.value = true;
  try {
    const response = await fetch(getApiUrl('/api/dns/list'), {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${localStorage.getItem('token')}`
      },
      body: JSON.stringify({
        pageNumber: page.value,
        pageSize: pageSize.value,
        search: searchQuery.value
      })
    });

    if (!response.ok) {
      if (response.status === 401) {
        localStorage.removeItem('token');
        router.push('/login');
        return;
      }
      throw new Error('获取DNS日志失败');
    }

    const data = await response.json();
    logs.value = (data.logs || []).map(log => ({
      id: log.id,
      timestamp: log.created_at,
      domain: log.host,
      query_type: log.type,
      client_ip: log.ip,
      response: log.response || '',
    }));
    totalCount.value = data.total || 0;
    totalPages.value = Math.ceil(totalCount.value / pageSize.value);
  } catch (err) {
    console.error('获取DNS日志错误:', err);
  } finally {
    loading.value = false;
  }
};

// 搜索日志 - 重置到第一页并查询
const handleSearch = () => {
  page.value = 1;
  fetchDnsLogs();
};

const handleDelete = async (id) => {
  if (!confirm('确定要删除该DNS日志吗？')) return;

  loading.value = true;
  try {
    const response = await fetch(getApiUrl('/api/dns/delete'), {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${localStorage.getItem('token')}`
      },
      body: JSON.stringify({ id: id })
    });

    if (!response.ok) {
      const errorData = await response.json();
      throw new Error(errorData.message || '删除DNS日志失败');
    }

    await fetchDnsLogs();
  } catch (err) {
    console.error('删除DNS日志错误:', err);
    alert('删除DNS日志失败: ' + err.message);
  } finally {
    loading.value = false;
  }
};

const handleClearLogs = async () => {
  if (!confirm('确定要清空所有DNS日志吗？此操作不可恢复！')) return;

  loading.value = true;
  try {
    const response = await fetch(getApiUrl('/api/dns/deleteAll'), {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${localStorage.getItem('token')}`
      }
    });

    if (!response.ok) throw new Error('清空日志失败');
    await fetchDnsLogs();
  } catch (err) {
    console.error('清空日志错误:', err);
    alert('清空日志失败: ' + err.message);
  } finally {
    loading.value = false;
  }
};

const prevPage = () => {
  if (page.value > 1) {
    page.value--;
    fetchDnsLogs();
  }
};

const nextPage = () => {
  if (page.value < totalPages.value) {
    page.value++;
    fetchDnsLogs();
  }
};

// 自动刷新功能
const toggleAutoRefresh = () => {
  autoRefresh.value = !autoRefresh.value;
  if (autoRefresh.value) {
    countdown.value = 1;
    startAutoRefresh();
  } else {
    stopAutoRefresh();
  }
};

const startAutoRefresh = () => {
  fetchDnsLogs();
  countdown.value = 1;

  refreshTimer = setInterval(() => {
    countdown.value--;
    if (countdown.value <= 0) {
      fetchDnsLogs();
      countdown.value = 1;
    }
  }, 1000);
};

const stopAutoRefresh = () => {
  if (refreshTimer) {
    clearInterval(refreshTimer);
    refreshTimer = null;
  }
};

// 监听域名管理
const fetchWatchDomains = async () => {
  loadingWatchDomains.value = true;
  try {
    const response = await fetch(getApiUrl('/api/watch_domains/list'), {
      method: 'GET',
      headers: { 'Authorization': `Bearer ${localStorage.getItem('token')}` }
    });

    if (!response.ok) {
      if (response.status === 401) {
        localStorage.removeItem('token');
        router.push('/login');
        return;
      }
      throw new Error('获取监听域名失败');
    }

    const data = await response.json();
    watchDomains.value = data.watch_domains || [];
  } catch (err) {
    console.error('获取监听域名错误:', err);
  } finally {
    loadingWatchDomains.value = false;
  }
};

const handleAddWatchDomain = async () => {
  if (!newWatchDomain.value.trim()) return;

  addingDomain.value = true;
  try {
    const response = await fetch(getApiUrl('/api/watch_domains/add'), {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${localStorage.getItem('token')}`
      },
      body: JSON.stringify({ domain: newWatchDomain.value.trim() })
    });

    if (!response.ok) {
      const errorData = await response.json();
      throw new Error(errorData.error || '添加监听域名失败');
    }

    newWatchDomain.value = '';
    await fetchWatchDomains();
    await handleRefreshWatchDomains();
  } catch (err) {
    console.error('添加监听域名错误:', err);
    alert('添加监听域名失败: ' + err.message);
  } finally {
    addingDomain.value = false;
  }
};

const handleDeleteWatchDomain = async (id) => {
  if (!confirm('确定要删除该监听域名吗？')) return;

  try {
    const response = await fetch(getApiUrl('/api/watch_domains/delete'), {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${localStorage.getItem('token')}`
      },
      body: JSON.stringify({ id: id })
    });

    if (!response.ok) {
      const errorData = await response.json();
      throw new Error(errorData.message || '删除监听域名失败');
    }

    await fetchWatchDomains();
    await handleRefreshWatchDomains();
  } catch (err) {
    console.error('删除监听域名错误:', err);
    alert('删除监听域名失败: ' + err.message);
  }
};

const handleRefreshWatchDomains = async () => {
  try {
    await fetch(getApiUrl('/api/watch_domains/refresh'), {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${localStorage.getItem('token')}`
      }
    });
    await fetchWatchDomains();
  } catch (err) {
    console.error('刷新监听域名错误:', err);
  }
};

// 组件生命周期
onMounted(() => {
  baseDomain.value = localStorage.getItem('host') || '';
  fetchDnsLogs();
  fetchWatchDomains();
});

onUnmounted(() => {
  stopAutoRefresh();
});
</script>
