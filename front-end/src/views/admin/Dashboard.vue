<script setup>
import { ref, onMounted } from 'vue'
import { get } from '@/util/request'
import { message } from 'ant-design-vue'
import SpinLoaderLarge from '@/components/spinLoaderLarge.vue'
import { BACKEND_DOMAIN } from "@/util/VARRIBLES"

const loading = ref(true)
const dashboardData = ref({
  statistics: {
    userCount: 0,
    projectCount: 0,
    chapterCount: 0,
    commentCount: 0
  },
  recentUsers: [],
  hotProjects: [],
  activeUsers: []
})

onMounted(() => {
  fetchDashboardData()
})

const fetchDashboardData = () => {
  get('/api/admin/dashboard', {}, 
    (msg, data) => {
      dashboardData.value = data
      loading.value = false
    },
    (msg) => {
      message.warning(msg)
      loading.value = false
    },
    (msg) => {
      message.error(msg)
      loading.value = false
    }
  )
}

const formatDate = (dateString) => {
  const date = new Date(dateString)
  return date.toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit'
  })
}
</script>

<template>
  <div class="p-4">
    <h1 class="text-2xl font-bold mb-6 dark:text-gray-100">管理仪表盘</h1>

    <div v-if="loading" class="flex justify-center items-center h-64">
      <SpinLoaderLarge />
    </div>

    <div v-else class="animate__animated animate__fadeIn">
      <!-- 统计卡片 -->
      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4 mb-8">
        <a-card class="bg-white dark:bg-zinc-800 shadow-sm hover:shadow-md transition-shadow">
          <template #title>
            <div class="flex items-center text-blue-600 dark:text-blue-400">
              <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6 mr-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" d="M15 19.128a9.38 9.38 0 0 0 2.625.372 9.337 9.337 0 0 0 4.121-.952 4.125 4.125 0 0 0-7.533-2.493M15 19.128v-.003c0-1.113-.285-2.16-.786-3.07M15 19.128v.106A12.318 12.318 0 0 1 8.624 21c-2.331 0-4.512-.645-6.374-1.766l-.001-.109a6.375 6.375 0 0 1 11.964-3.07M12 6.375a3.375 3.375 0 1 1-6.75 0 3.375 3.375 0 0 1 6.75 0Zm8.25 2.25a2.625 2.625 0 1 1-5.25 0 2.625 2.625 0 0 1 5.25 0Z" />
              </svg>
              用户总数
            </div>
          </template>
          <h3 class="text-3xl font-bold text-gray-700 dark:text-gray-300">{{ dashboardData.statistics.userCount }}</h3>
          <p class="text-gray-500 dark:text-gray-400 text-sm">平台注册用户总数</p>
        </a-card>

        <a-card class="bg-white dark:bg-zinc-800 shadow-sm hover:shadow-md transition-shadow">
          <template #title>
            <div class="flex items-center text-green-600 dark:text-green-400">
              <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6 mr-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" d="m16.862 4.487 1.687-1.688a1.875 1.875 0 1 1 2.652 2.652L10.582 16.07a4.5 4.5 0 0 1-1.897 1.13L6 18l.8-2.685a4.5 4.5 0 0 1 1.13-1.897l8.932-8.931Zm0 0L19.5 7.125M18 14v4.75A2.25 2.25 0 0 1 15.75 21H5.25A2.25 2.25 0 0 1 3 18.75V8.25A2.25 2.25 0 0 1 5.25 6H10" />
              </svg>
              项目总数
            </div>
          </template>
          <h3 class="text-3xl font-bold text-gray-700 dark:text-gray-300">{{ dashboardData.statistics.projectCount }}</h3>
          <p class="text-gray-500 dark:text-gray-400 text-sm">平台创建的项目总数</p>
        </a-card>

        <a-card class="bg-white dark:bg-zinc-800 shadow-sm hover:shadow-md transition-shadow">
          <template #title>
            <div class="flex items-center text-purple-600 dark:text-purple-400">
              <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6 mr-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" d="M12 6.042A8.967 8.967 0 0 0 6 3.75c-1.052 0-2.062.18-3 .512v14.25A8.987 8.987 0 0 1 6 18c2.305 0 4.408.867 6 2.292m0-14.25a8.966 8.966 0 0 1 6-2.292c1.052 0 2.062.18 3 .512v14.25A8.987 8.987 0 0 0 18 18a8.967 8.967 0 0 0-6 2.292m0-14.25v14.25" />
              </svg>
              章节总数
            </div>
          </template>
          <h3 class="text-3xl font-bold text-gray-700 dark:text-gray-300">{{ dashboardData.statistics.chapterCount }}</h3>
          <p class="text-gray-500 dark:text-gray-400 text-sm">平台创建的章节总数</p>
        </a-card>

        <a-card class="bg-white dark:bg-zinc-800 shadow-sm hover:shadow-md transition-shadow">
          <template #title>
            <div class="flex items-center text-orange-600 dark:text-orange-400">
              <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6 mr-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" d="M7.5 8.25h9m-9 3H12m-9.75 1.51c0 1.6 1.123 2.994 2.707 3.227 1.129.166 2.27.293 3.423.379.35.026.67.21.865.501L12 21l2.755-4.133a1.14 1.14 0 0 1 .865-.501 48.172 48.172 0 0 0 3.423-.379c1.584-.233 2.707-1.626 2.707-3.228V6.741c0-1.602-1.123-2.995-2.707-3.228A48.394 48.394 0 0 0 12 3c-2.392 0-4.744.175-7.043.513C3.373 3.746 2.25 5.14 2.25 6.741v6.018Z" />
              </svg>
              评论总数
            </div>
          </template>
          <h3 class="text-3xl font-bold text-gray-700 dark:text-gray-300">{{ dashboardData.statistics.commentCount }}</h3>
          <p class="text-gray-500 dark:text-gray-400 text-sm">平台评论总数</p>
        </a-card>
      </div>

      <!-- 数据列表 -->
      <div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
        <!-- 最近注册用户 -->
        <a-card class="bg-white dark:bg-zinc-800 shadow-sm">
          <template #title>
            <div class="flex items-center font-bold text-gray-800 dark:text-gray-200">
              <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 mr-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" d="M19 7.5v3m0 0v3m0-3h3m-3 0h-3m-2.25-4.125a3.375 3.375 0 1 1-6.75 0 3.375 3.375 0 0 1 6.75 0ZM4 19.235v-.11a6.375 6.375 0 0 1 12.75 0v.109A12.318 12.318 0 0 1 10.374 21c-2.331 0-4.512-.645-6.374-1.766Z" />
              </svg>
              最近注册用户
            </div>
          </template>
          <div class="space-y-4">
            <div v-for="user in dashboardData.recentUsers" :key="user.ID" 
                class="flex items-center p-3 rounded-lg hover:bg-gray-50 dark:hover:bg-zinc-700 transition-colors">
              <a-avatar :src="user.avatar ? BACKEND_DOMAIN + user.avatar : ''" :size="40">
                {{ user.username ? user.username.charAt(0).toUpperCase() : 'U' }}
              </a-avatar>
              <div class="ml-3 flex-grow">
                <p class="font-medium text-gray-900 dark:text-gray-100">{{ user.username }}</p>
                <p class="text-sm text-gray-500 dark:text-gray-400">{{ formatDate(user.CreatedAt) }}</p>
              </div>
            </div>
          </div>
        </a-card>

        <!-- 热门项目 -->
        <a-card class="bg-white dark:bg-zinc-800 shadow-sm">
          <template #title>
            <div class="flex items-center font-bold text-gray-800 dark:text-gray-200">
              <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 mr-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" d="M15.362 5.214A8.252 8.252 0 0 1 12 21 8.25 8.25 0 0 1 6.038 7.047 8.287 8.287 0 0 0 9 9.601a8.983 8.983 0 0 1 3.361-6.867 8.21 8.21 0 0 0 3 2.48Z" />
                <path stroke-linecap="round" stroke-linejoin="round" d="M12 18a3.75 3.75 0 0 0 .495-7.468 5.99 5.99 0 0 0-1.925 3.547 5.975 5.975 0 0 1-2.133-1.001A3.75 3.75 0 0 0 12 18Z" />
              </svg>
              热门项目
            </div>
          </template>
          <div class="space-y-4">
            <div v-for="project in dashboardData.hotProjects" :key="project.ID" 
                class="p-3 rounded-lg hover:bg-gray-50 dark:hover:bg-zinc-700 transition-colors">
              <div class="flex justify-between items-center">
                <p class="font-medium text-gray-900 dark:text-gray-100">{{ project.project_name }}</p>
                <a-tag color="blue">{{ project.watches }} 观看</a-tag>
              </div>
              <p class="text-sm text-gray-500 dark:text-gray-400 line-clamp-2 mt-1">{{ project.social_story }}</p>
            </div>
          </div>
        </a-card>

        <!-- 活跃用户 -->
        <a-card class="bg-white dark:bg-zinc-800 shadow-sm">
          <template #title>
            <div class="flex items-center font-bold text-gray-800 dark:text-gray-200">
              <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 mr-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" d="M18 18.72a9.094 9.094 0 0 0 3.741-.479 3 3 0 0 0-4.682-2.72m.94 3.198.001.031c0 .225-.012.447-.037.666A11.944 11.944 0 0 1 12 21c-2.17 0-4.207-.576-5.963-1.584A6.062 6.062 0 0 1 6 18.719m12 0a5.971 5.971 0 0 0-.941-3.197m0 0A5.995 5.995 0 0 0 12 12.75a5.995 5.995 0 0 0-5.058 2.772m0 0a3 3 0 0 0-4.681 2.72 8.986 8.986 0 0 0 3.74.477m.94-3.197a5.971 5.971 0 0 0-.94 3.197M15 6.75a3 3 0 1 1-6 0 3 3 0 0 1 6 0Z" />
              </svg>
              活跃用户
            </div>
          </template>
          <div class="space-y-4">
            <div v-for="user in dashboardData.activeUsers" :key="user.userID" 
                class="flex items-center p-3 rounded-lg hover:bg-gray-50 dark:hover:bg-zinc-700 transition-colors">
              <a-avatar :src="user.avatar ? BACKEND_DOMAIN + user.avatar : ''" :size="40">
                {{ user.username ? user.username.charAt(0).toUpperCase() : 'U' }}
              </a-avatar>
              <div class="ml-3 flex-grow">
                <p class="font-medium text-gray-900 dark:text-gray-100">{{ user.username }}</p>
                <p class="text-sm text-gray-500 dark:text-gray-400">活跃度: {{ user.activityCount }}</p>
              </div>
            </div>
          </div>
        </a-card>
      </div>
    </div>
  </div>
</template>

<style scoped>
.animate__animated {
  animation-duration: 0.6s;
}
</style>
