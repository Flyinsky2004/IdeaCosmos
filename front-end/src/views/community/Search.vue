<script setup>
import { ref, reactive, onMounted, watch } from 'vue'
import { useRoute } from 'vue-router'
import { get } from '@/util/request'
import { message } from 'ant-design-vue'
import { imagePrefix } from '@/util/VARRIBLES'

const route = useRoute()
const options = reactive({
  projects: [],
  pageIndex: 0,
  loading: false
})

const searchProjects = async () => {
  options.loading = true
  const keyword = route.params.keyword
  
  try {
    await get('/api/public/searchProjects', {
      keyword,
      pageIndex: options.pageIndex
    }, (messageer, data) => {
      options.projects = data
    })
  } finally {
    options.loading = false
  }
}

watch(() => route.params.keyword, () => {
  options.pageIndex = 0
  searchProjects()
})

onMounted(() => {
  searchProjects()
})
</script>

<template>
  <div class="w-full space-y-6">
    <!-- 搜索结果标题 -->
    <div class="flex items-center justify-between">
      <h2 class="text-xl font-semibold text-gray-900 dark:text-gray-100">
        搜索结果：{{ route.params.keyword }}
      </h2>
      <span class="text-sm text-gray-500 dark:text-gray-400">
        找到 {{ options.projects.length }} 个项目
      </span>
    </div>

    <!-- 加载中状态 -->
    <div v-if="options.loading" class="flex justify-center py-12">
      <div class="text-center">
        <svg class="animate-spin h-8 w-8 text-blue-500 mx-auto mb-4" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
          <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
          <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
        </svg>
        <span class="text-gray-500 dark:text-gray-400">搜索中...</span>
      </div>
    </div>

    <!-- 项目列表 -->
    <div v-else class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-6">
      <div 
        v-for="project in options.projects" 
        :key="project.ID" 
        class="group bg-white/90 dark:bg-zinc-900 border theme-border rounded-xl overflow-hidden hover:shadow-lg dark:hover:shadow-zinc-800 transition-all duration-300 hover:-translate-y-1"
      >
        <!-- 使用与 Index.vue 相同的项目卡片结构 -->
      </div>
    </div>

    <!-- 无结果提示 -->
    <div v-if="!options.loading && options.projects.length === 0" class="text-center py-12">
      <svg xmlns="http://www.w3.org/2000/svg" class="h-12 w-12 text-gray-400 mx-auto mb-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
        <path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12" />
      </svg>
      <p class="text-gray-500 dark:text-gray-400">未找到相关项目</p>
    </div>
  </div>
</template> 