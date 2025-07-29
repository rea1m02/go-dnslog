<template>
  <div class="w-full">
    <div class="flex justify-between items-center mb-6">
      <h1 class="text-2xl font-bold">DNS日志</h1>
      <div class="flex space-x-3">
        <div class="relative">
          <input type="text" v-model="searchQuery" placeholder="搜索域名或IP..."
            class="pl-10 pr-4 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-indigo-500"
            @input="handleSearch">
          <i class="fa fa-search absolute left-3 top-3 text-gray-400"></i>
        </div>
        <button @click="handleClearLogs" class="bg-red-600 text-white px-4 py-2 rounded-md hover:bg-red-700 transition-colors"
          :disabled="loading">
          <i class="fa fa-trash mr-1"></i>清空日志
        </button>
      </div>
    </div>

    <!-- 日志表格 -->
    <div class="bg-white rounded-lg shadow-md overflow-hidden">
      <div class="overflow-x-auto">
        <table class="min-w-full divide-y divide-gray-200">
          <thead class="bg-gray-50">
            <tr>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">ID</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">域名</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">源IP</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">查询类型</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">时间</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">操作</th>
            </tr>
          </thead>
          <tbody class="bg-white divide-y divide-gray-200">
            <tr v-if="loading">
              <td colspan="5" class="px-6 py-4 text-center text-gray-500">加载中...</td>
            </tr>
            <tr v-else-if="logs.length === 0">
              <td colspan="5" class="px-6 py-4 text-center text-gray-500">暂无日志记录</td>
            </tr>
            <tr v-for="log in logs" :key="log.id">
              <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">{{ log.id }}</td>
              <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">{{ log.domain }}</td>
              <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">{{ log.client_ip }}</td>
              <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">{{ log.query_type }}</td>
              <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">{{ formatTime(log.timestamp) }}</td>
              <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                <button @click="handleDelete(log.id)" class="text-write-600 hover:text-red-900">
                  <i class="fa fa-trash">删除</i>
                </button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
      <!-- 分页 -->
      <div class="px-6 py-4 border-t border-gray-200 flex justify-between items-center">
        <div class="text-sm text-gray-700">
          显示 <span>{{ logs.length }}</span> 条，共 <span>{{ totalCount }}</span> 条
        </div>
        <div class="flex space-x-2">
          <button @click="prevPage()" 
            class="px-3 py-1 border border-gray-300 rounded-md text-sm bg-white hover:bg-gray-50 disabled:opacity-50">
            上一页
          </button>
          <button @click="nextPage()" 
            class="px-3 py-1 border border-gray-300 rounded-md text-sm bg-white hover:bg-gray-50 disabled:opacity-50">
            下一页
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, watch } from 'vue';
import { useRouter } from 'vue-router';

const logs = ref([]);
const loading = ref(false);
const totalCount = ref(0);
const page = ref(1);
const pageSize = ref(10);
const totalPages = ref(0);
const searchQuery = ref('');
const router = useRouter();

// 格式化时间
const formatTime = (timeStr) => {
  const date = new Date(timeStr);
  return date.toLocaleString('zh-CN', { hour12: false });
};

// 获取DNS日志
const fetchDnsLogs = async () => {
  loading.value = true;
  try {
    const response = await fetch('https://dns.rea1m.top/api/dns/list', {
      method: 'POST',
      headers: { 'Authorization': `Bearer ${localStorage.getItem('token')}` },
      body: JSON.stringify({
        pageNumber: page.value,
        pageSize: pageSize.value,
        search: searchQuery.value
      })
    });

    if (!response.ok) {
      if (response.status === 401) {
        // 未授权，跳转到登录页
        localStorage.removeItem('token');
        router.push('/login');
        return;
      }
      throw new Error('获取DNS日志失败');
    }

    const data = await response.json();
    logs.value = (data.logs || []).map(log => ({
      id: log.id,
      timestamp: log.created_at,      // 映射
      domain: log.host,               // 映射
      query_type: log.type,           // 映射
      client_ip: log.ip,              // 映射
      response: log.response || '',   // 兼容后端无response字段
    }));
    totalCount.value = data.total || 0;
    totalPages.value = Math.ceil(totalCount.value / pageSize.value);
  } catch (err) {
    console.error('获取DNS日志错误:', err);
    console.log('获取日志失败: ' + err.message);
  } finally {
    loading.value = false;
  }
};

// 搜索日志
const handleSearch = () => {
  page.value = 1; // 重置到第一页
  fetchDnsLogs();
};


const handleDelete = async (id) => {
  if (!confirm('确定要删除该DNS日志吗？')) return;
  
  loading.value = true;
  try {
    const response = await fetch('https://dns.rea1m.top/api/dns/delete', {
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

    // 删除成功后重新获取数据
    await fetchDnsLogs();
    alert('删除成功');
  } catch (err) {
    console.error('删除DNS日志错误:', err);
    alert('删除DNS日志失败: ' + err.message);
  } finally {
    loading.value = false;
  }
};


// 清空日志
const handleClearLogs = async () => {
  if (!confirm('确定要清空所有DNS日志吗？')) return;

  loading.value = true;
  try {
    const response = await fetch('https://dns.rea1m.top/api/dns/deleteAll', {
      method: 'POST',
      headers: { 'Authorization': `Bearer ${localStorage.getItem('token')}` }
    });

    if (!response.ok) throw new Error('清空日志失败');
    await fetchDnsLogs(); // 重新加载日志
    alert('日志已清空');
  } catch (err) {
    console.error('清空日志错误:', err);
    alert('清空日志失败: ' + err.message);
  } finally {
    loading.value = false;
  }
};

// 上一页
const prevPage = async() => {
  // console.log("prev page: ", page.value, totalPages.value);  
  if (page.value > 1) {
    page.value--;

    fetchDnsLogs();
  }
};

// 下一页
const nextPage = async() => {
  // console.log("next page: ", page.value, totalPages.value);  
  if (page.value < totalPages.value) {
    page.value++;
    fetchDnsLogs();
  }
};

// 监听搜索查询变化
watch(searchQuery, (newVal, oldVal) => {
  if (newVal !== oldVal) {
    // 可以添加防抖处理
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