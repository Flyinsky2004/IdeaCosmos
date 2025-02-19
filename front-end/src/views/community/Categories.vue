<script setup>
import { ref, reactive } from 'vue'
import { get } from '@/util/request'
import { message } from 'ant-design-vue'
import { imagePrefix } from '@/util/VARRIBLES'

const categories = [
  { name: '电影剧本', icon: 'M3.375 19.5h17.25m-17.25 0a1.125 1.125 0 01-1.125-1.125M3.375 19.5h1.5C5.496 19.5 6 18.996 6 18.375m-3.75 0V5.625m0 12.75v-1.5c0-.621.504-1.125 1.125-1.125m18.375 2.625V5.625m0 12.75c0 .621-.504 1.125-1.125 1.125m1.125-1.125v-1.5c0-.621-.504-1.125-1.125-1.125m0 3.75h-1.5A1.125 1.125 0 0118 18.375M20.625 4.5H3.375m17.25 0c.621 0 1.125.504 1.125 1.125M20.625 4.5h-1.5C18.504 4.5 18 5.004 18 5.625m3.75 0v1.5c0 .621-.504 1.125-1.125 1.125M3.375 4.5c-.621 0-1.125.504-1.125 1.125M3.375 4.5h1.5C5.496 4.5 6 5.004 6 5.625m-3.75 0v1.5c0 .621.504 1.125 1.125 1.125m0 0h1.5m-1.5 0c-.621 0-1.125.504-1.125 1.125v1.5c0 .621.504 1.125 1.125 1.125m1.5-3.75C5.496 8.25 6 7.746 6 7.125v-1.5M4.875 8.25C5.496 8.25 6 8.754 6 9.375v1.5m0-5.25v5.25m0-5.25C6 5.004 6.504 4.5 7.125 4.5h9.75c.621 0 1.125.504 1.125 1.125m1.125 2.625h1.5m-1.5 0A1.125 1.125 0 0118 7.125v-1.5m1.125 2.625c-.621 0-1.125.504-1.125 1.125v1.5m2.625-2.625c.621 0 1.125.504 1.125 1.125v1.5c0 .621-.504 1.125-1.125 1.125M18 5.625v5.25M7.125 12h9.75m-9.75 0A1.125 1.125 0 016 10.875M7.125 12C6.504 12 6 12.504 6 13.125m0-2.25C6 11.496 5.496 12 4.875 12M18 10.875c0 .621-.504 1.125-1.125 1.125M18 10.875c0 .621.504 1.125 1.125 1.125m-2.25 0c.621 0 1.125.504 1.125 1.125m-12 5.25v-5.25m0 5.25c0 .621.504 1.125 1.125 1.125h9.75c.621 0 1.125-.504 1.125-1.125m-12 0v-1.5c0-.621-.504-1.125-1.125-1.125M18 18.375v-5.25m0 5.25v-1.5c0-.621.504-1.125 1.125-1.125M18 13.125v1.5c0 .621.504 1.125 1.125 1.125M18 13.125c0-.621.504-1.125 1.125-1.125M6 13.125v1.5c0 .621-.504 1.125-1.125 1.125M6 13.125C6 12.504 5.496 12 4.875 12m-1.5 0h1.5m-1.5 0c-.621 0-1.125.504-1.125 1.125v1.5c0 .621.504 1.125 1.125 1.125M19.125 12h1.5m0 0c.621 0 1.125.504 1.125 1.125v1.5c0 .621-.504 1.125-1.125 1.125m-17.25 0h1.5M7.125 4.5v1.5m0-1.5h9.75m-9.75 0h-1.5M7.125 4.5C6.504 4.5 6 5.004 6 5.625M7.125 4.5h9.75C17.496 4.5 18 5.004 18 5.625' },
  { name: '短视频', icon: 'M15.75 10.5l4.72-4.72a.75.75 0 011.28.53v11.38a.75.75 0 01-1.28.53l-4.72-4.72M4.5 18.75h9a2.25 2.25 0 002.25-2.25v-9a2.25 2.25 0 00-2.25-2.25h-9A2.25 2.25 0 002.25 7.5v9a2.25 2.25 0 002.25 2.25z' },
  { name: '小说', icon: 'M12 6.042A8.967 8.967 0 006 3.75c-1.052 0-2.062.18-3 .512v14.25A8.987 8.987 0 016 18c2.305 0 4.408.867 6 2.292m0-14.25a8.966 8.966 0 016-2.292c1.052 0 2.062.18 3 .512v14.25A8.987 8.987 0 0018 18a8.967 8.967 0 00-6 2.292m0-14.25v14.25' },
  // ... 更多分类
]

const selectedCategory = ref(null)
const options = reactive({
  projects: [],
  pageIndex: 0
})

const fetchCategoryProjects = (category) => {
  selectedCategory.value = category
  get('/api/public/getCategoryProjects', {
    category: category.name,
    pageIndex: options.pageIndex
  }, (messageer, data) => {
    options.projects = data
  })
}
</script>

<template>
  <div class="w-full space-y-8">
    <!-- 分类导航 -->
    <div class="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-6 gap-4">
      <button
        v-for="category in categories"
        :key="category.name"
        @click="fetchCategoryProjects(category)"
        :class="[
          'p-4 rounded-xl border theme-border transition-all duration-300',
          selectedCategory?.name === category.name
            ? 'bg-blue-500 text-white'
            : 'bg-white dark:bg-zinc-900 hover:bg-blue-50 dark:hover:bg-blue-900/20'
        ]"
      >
        <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6 mx-auto mb-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" :d="category.icon" />
        </svg>
        <span class="text-sm font-medium">{{ category.name }}</span>
      </button>
    </div>

    <!-- 项目列表 -->
    <div v-if="selectedCategory" class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-6">
      <div 
        v-for="project in options.projects" 
        :key="project.ID" 
        class="group bg-white/90 dark:bg-zinc-900 border theme-border rounded-xl overflow-hidden hover:shadow-lg dark:hover:shadow-zinc-800 transition-all duration-300 hover:-translate-y-1"
      >
        <!-- 使用与 Index.vue 相同的项目卡片结构 -->
      </div>
    </div>
  </div>
</template> 