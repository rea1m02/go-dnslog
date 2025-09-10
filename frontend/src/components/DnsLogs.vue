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
          <p class="text-gray-500 mt-1">查看和管理所有 DNS 解析记录</p>
        </div>
        <div class="flex flex-col sm:flex-row gap-3">
          <div class="relative">
            <i class="fa fa-search absolute left-4 top-1/2 -translate-y-1/2 text-gray-400"></i>
            <input type="text" v-model="searchQuery" placeholder="搜索域名或IP..."
                   class="pl-11 pr-4 py-2.5 bg-white border-2 border-gray-200 rounded-xl
                          focus:outline-none focus:border-indigo-500 focus:ring-4 focus:ring-indigo-500/10
                          transition-all duration-200 w-full sm:w-72">
          </div>
          <button @click="handleClearLogs"
                  class="px-5 py-2.5 bg-gradient-to-r from-red-500 to-rose-500 text-white font-medium rounded-xl
                         hover:from-red-600 hover:to-rose-600 focus:outline-none focus:ring-4 focus:ring-red-500/20
                         transition-all duration-200 flex items-center space-x-2 disabled:opacity-50"
                  :disabled="loading">
            <i class="fa fa-trash"></i>
            <span>清空日志</span>
          </button>
        </div>
      </div>
    </div>

    <!-- 日志表格卡片 -->
    <div class="bg-white rounded-2xl shadow-xl border border-gray-100 overflow-hidden">
      <!-- 表格头部统计 -->
      <div class="px-6 py-4 bg-gradient-to-r from-gray-50 to-gray-100 border-b border-gray-100">
        <div class="flex items-center justify-between">
          <div class="flex items-center space-x-4">
            <div class="flex items-center space-x-2">
              <div class="w-3 h-3 bg-green-500 rounded-full animate-pulse"></div>
              <span class="text-sm text-gray-600">实时更新</span>
            </div>
            <span class="text-sm text-gray-500">|</span>
            <span class="text-sm font-medium text-gray-700">
              共 <span class="text-indigo-600">{{ totalCount }}</span> 条记录
            </span>
          </div>
        </div>
      </div>

      <!-- 表格内容 -->
      <div class="overflow-x-auto">
        <table class="min-w-full">
          <thead>
            <tr class="bg-gray-50">
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
                  <span class="text-gray-400 text-sm">DNS 查询记录将显示在这里</span>
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
          <button @click="prevPage()"
                  class="px-4 py-2 bg-white border border-gray-200 rounded-lg text-sm font-medium text-gray-600
                         hover:bg-gray-50 hover:text-gray-900 disabled:opacity-50 disabled:cursor-not-allowed
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
          <button @click="nextPage()"
                  class="px-4 py-2 bg-white border border-gray-200 rounded-lg text-sm font-medium text-gray-600
                         hover:bg-gray-50 hover:text-gray-900 disabled:opacity-50 disabled:cursor-not-allowed
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
import { ref, onMounted, watch } from 'vue';
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

const prevPage = async() => {
  if (page.value > 1) {
    page.value--;
    fetchDnsLogs();
  }
};

const nextPage = async() => {
  if (page.value < totalPages.value) {
    page.value++;
    fetchDnsLogs();
  }
};

watch(searchQuery, (newVal, oldVal) => {
  if (newVal !== oldVal) {
    const timeout = setTimeout(() => {
      handleSearch();
    }, 500);
    return () => clearTimeout(timeout);
  }
});

onMounted(() => {
  fetchDnsLogs();
});
</script>
