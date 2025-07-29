<template>
  <div class="w-full">
    <h1 class="text-2xl font-bold mb-6">DNS Rebind</h1>

    <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
      <!-- 生成Rebind规则 -->
      <div class="bg-white p-6 rounded-lg shadow-md">
        <h2 class="text-lg font-semibold mb-4">生成Rebind规则</h2>
        <form @submit.prevent="handleGenerate" class="space-y-4">
          <div v-if="errorMessage" class="p-3 bg-red-50 border border-red-200 rounded-md text-red-700 text-sm">
            {{ errorMessage }}
          </div>
          <div v-if="successMessage" class="p-3 bg-green-50 border border-green-200 rounded-md text-green-700 text-sm">
            {{ successMessage }}
          </div>
          <div>
            <label for="firstIp" class="block text-sm font-medium text-gray-700">FirstIP</label>
            <input type="text" id="firstIp" v-model="firstIp" placeholder="例如: 192.168.1.1"
              class="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500">
          </div>
          <div>
            <label for="secondIp" class="block text-sm font-medium text-gray-700">SecondIP</label>
            <input type="text" id="secondIp" v-model="secondIp" placeholder="例如: 127.0.0.1"
              class="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500">
          </div>
          <div>
            <button type="submit" :disabled="loading"
              class="w-full flex justify-center py-2 px-4 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 disabled:opacity-70">
              <i v-if="loading" class="fa fa-spinner fa-spin mr-2"></i>
              生成Rebind域名
            </button>
          </div>
        </form>
        <div id="resultContainer" class="mt-4" v-if="rebindUrl">
          <div class="p-3 bg-green-50 border border-green-200 rounded-md">
            <h3 class="font-medium text-green-800">生成成功</h3>
            <p class="mt-1 text-sm text-green-700">Rebind链接: <span class="font-mono break-all">{{ rebindUrl }}</span></p>
          </div>
        </div>
      </div>

      <!-- 活跃Rebind规则 -->
      <div class="bg-white p-6 rounded-lg shadow-md">
        <h2 class="text-lg font-semibold mb-4">活跃Rebind规则</h2>
        <div class="space-y-3">
          <div v-if="loading" class="text-center text-gray-500 py-4">加载中...</div>
          <div v-else-if="rules.length === 0" class="text-center text-gray-500 py-4">暂无活跃规则</div>
          <div v-for="rule in rules" :key="rule.id" class="border border-gray-200 rounded-md p-4">
            <div class="flex justify-between items-start">
              <div>
                <p class="font-medium">{{ rule.domain }}</p>
                <p class="text-sm text-gray-500 mt-1">First_IP: {{ rule.first_ip }}</p>
                <p class="text-sm text-gray-500">Second_IP: {{ rule.second_ip }}</p>
              </div>
              <div class="text-right">
                <button @click="handleDeleteRule(rule.id)" class="text-red-500 hover:text-red-700 mt-1 text-sm">
                  <i class="fa fa-trash mr-1"></i>删除
                </button>
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

const firstIp = ref('');
const secondIp = ref('');
const rules = ref([]);
const rebindUrl = ref('');
const loading = ref(false);
const errorMessage = ref('');
const successMessage = ref('');
const router = useRouter();

// 格式化时间
const formatTime = (timeStr) => {
  const date = new Date(timeStr);
  return date.toLocaleString('zh-CN', { hour12: false });
};

// 获取活跃Rebind规则
const fetchRebindRules = async () => {
  loading.value = true;
  try {
    const response = await fetch('https://dns.rea1m.top/api/rebind/list', {
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
    errorMessage.value = '获取规则失败: ' + err.message;
  } finally {
    loading.value = false;
  }
};

// 生成Rebind规则
const handleGenerate = async () => {
  loading.value = true;
  errorMessage.value = '';
  successMessage.value = '';
  rebindUrl.value = '';

  try {
    const response = await fetch('https://dns.rea1m.top/api/rebind/gen', {
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
    if (!response.ok) throw new Error(data.message || '生成规则失败');

    rebindUrl.value = data.url;
    successMessage.value = 'Rebind规则生成成功';
    // 重新获取规则列表
    await fetchRebindRules();
  } catch (err) {
    console.error('生成Rebind规则错误:', err);
    errorMessage.value = err.message;
  } finally {
    loading.value = false;
  }
};

// 删除Rebind规则
const handleDeleteRule = async (id) => {
  if (!confirm('确定要删除此规则吗？')) return;

  try {
    const response = await fetch(`https://dns.rea1m.top/api/rebind/delete`, {
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
    // 重新获取规则列表
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