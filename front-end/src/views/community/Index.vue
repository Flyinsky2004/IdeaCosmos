<script setup>
import { get } from '@/util/request';
import { message } from 'ant-design-vue';
import { onMounted, reactive, ref } from 'vue'
import { imagePrefix } from '@/util/VARRIBLES';
import router from '@/router';
const options = reactive({
  projects: [],
  pageIndex: 0
})

const fetchData = () => {
  get('/api/public/getIndexProject',{
    pageIndex: options.pageIndex
  },(messageer,data) => {
    options.pageIndex += 10
    options.projects = data
  },(messageer,data) => {
    fetchData()
  },(messageer,data) => {
    fetchData()
  })
}

onMounted(() => {
  fetchData()
})

const posts = ref([
  {
    id: 1,
    title: '我的第一个AI创意项目',
    cover: 'https://example.com/image1.jpg',
    author: '创意达人',
    views: 1234,
    likes: 88,
    createTime: '2024-01-10'
  },
  // ... 更多帖子数据
])

// 查看项目详情
const viewProject = (projectId) => {
  router.push(`/community/project/${projectId}`)
}

// 加载更多
const loadMore = () => {
  fetchData()
}
</script>

<template>
  <div class="min-h-screen py-8 bg-gray-50 dark:bg-zinc-900/60 rounded-xl">
    <div class="max-w-7xl mx-auto px-4">
      <!-- 轮播图区域 -->
      <div class="mb-12 rounded-2xl overflow-hidden border theme-border">
        <a-carousel autoplay>
          <div 
            v-for="(project, index) in options.projects.slice(0, 4)" 
            :key="index"
            class="relative h-[400px] cursor-pointer"
            @click="viewProject(project.ID)"
          >
            <!-- 背景图 -->
            <div class="absolute inset-0">
              <img 
                :src="imagePrefix + project.cover_image"
                class="w-full h-full object-cover"
                alt="项目封面"
              />
              <!-- 渐变遮罩 -->
              <div class="absolute inset-0 bg-gradient-to-t from-black/80 via-black/30 to-transparent"></div>
            </div>
            
            <!-- 项目信息 -->
            <div class="absolute bottom-0 left-0 right-0 p-8 text-white">
              <div class="max-w-3xl">
                <h2 class="text-3xl font-bold mb-4 line-clamp-2">
                  {{ project.project_name }}
                </h2>
                <p class="text-gray-200 line-clamp-2 mb-4">
                  {{ project.social_story }}
                </p>
                <div class="flex items-center gap-6">
                  <div class="flex items-center gap-2">
                    <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                      <path stroke-linecap="round" stroke-linejoin="round" d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z" />
                    </svg>
                    <span>{{ project.team?.username }}</span>
                  </div>
                  <div class="flex items-center gap-2">
                    <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                      <path stroke-linecap="round" stroke-linejoin="round" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
                      <path stroke-linecap="round" stroke-linejoin="round" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z" />
                    </svg>
                    <span>{{ project.watches }}</span>
                  </div>
                  <div class="flex flex-wrap gap-2">
                    <span v-for="(style, idx) in project.style" 
                      :key="idx"
                      class="px-3 py-1 bg-white/20 backdrop-blur-sm rounded-full text-sm"
                    >
                      {{ style }}
                    </span>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </a-carousel>
      </div>

      <!-- 页面标题和筛选区 -->
      <div class="flex flex-col md:flex-row md:items-center justify-between gap-6 mb-8">
        <div class="flex items-center gap-3">
          <div class="h-8 w-1 bg-gradient-to-b from-blue-500 to-cyan-500 rounded-full"></div>
          <h1 class="text-3xl font-bold text-gray-900 dark:text-gray-100">
            社区作品
          </h1>
        </div>
        
        <!-- 筛选和排序 -->
        <div class="flex items-center gap-4">
          <div class="flex items-center gap-2 text-sm">
            <span class="text-gray-600 dark:text-gray-400">类型：</span>
            <select class="px-3 py-2 rounded-lg border theme-border bg-white dark:bg-zinc-800 text-gray-900 dark:text-gray-100">
              <option value="">全部</option>
              <option value="剧本">剧本</option>
              <option value="小说">小说</option>
            </select>
          </div>
          <div class="flex items-center gap-2 text-sm">
            <span class="text-gray-600 dark:text-gray-400">排序：</span>
            <select class="px-3 py-2 rounded-lg border theme-border bg-white dark:bg-zinc-800 text-gray-900 dark:text-gray-100">
              <option value="latest">最新发布</option>
              <option value="popular">最受欢迎</option>
            </select>
          </div>
        </div>
      </div>

      <!-- 项目网格 -->
      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        <div v-for="project in options.projects" 
          :key="project.ID"
          class="group bg-white dark:bg-zinc-800 rounded-xl border theme-border hover:shadow-xl transition-all duration-300 hover:-translate-y-1 overflow-hidden"
        >
          <!-- 项目封面 -->
          <div class="aspect-video relative overflow-hidden">
            <img 
              :src="imagePrefix + project.cover_image || '/default-cover.jpg'"
              class="w-full h-full object-cover transform group-hover:scale-105 transition-transform duration-300"
              alt="项目封面"
            />
            <!-- 项目类型标签 -->
            <div class="absolute top-3 left-3 px-3 py-1 bg-black/60 backdrop-blur-sm rounded-full text-white text-sm">
              {{ project.types }}
            </div>
          </div>

          <!-- 项目内容 -->
          <div class="p-6 space-y-4">
            <!-- 标题和团队信息 -->
            <div>
              <h3 class="text-xl font-bold text-gray-900 dark:text-gray-100 group-hover:text-blue-500 transition-colors mb-2">
                {{ project.project_name }}
              </h3>
              <div class="flex items-center gap-4 text-sm text-gray-600 dark:text-gray-400">
                <div class="flex items-center gap-2">
                  <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z" />
                  </svg>
                  <span>{{ project.team?.username || '未知作者' }}</span>
                </div>
                <div class="flex items-center gap-2">
                  <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
                    <path stroke-linecap="round" stroke-linejoin="round" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z" />
                  </svg>
                  <span>{{ project.watches }}</span>
                </div>
              </div>
            </div>

            <!-- 项目描述 -->
            <p class="text-gray-600 dark:text-gray-400 text-sm line-clamp-3">
              {{ project.social_story }}
            </p>

            <!-- 标签和操作按钮 -->
            <div class="flex items-center justify-between pt-4 border-t theme-border">
              <div class="flex flex-wrap gap-2">
                <span v-for="(style, index) in project.style" 
                  :key="index"
                  class="px-2 py-1 text-xs bg-gray-100 dark:bg-zinc-700/50 text-gray-600 dark:text-gray-400 rounded-full"
                >
                  {{ style }}
                </span>
              </div>
              <button 
                @click="viewProject(project.ID)"
                class="shrink-0 px-4 py-2 bg-gradient-to-r from-blue-500 to-cyan-500 text-white rounded-lg hover:opacity-90 transition-opacity flex items-center gap-2"
              >
                <span>查看详情</span>
                <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" d="M9 5l7 7-7 7" />
                </svg>
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
/* 自定义轮播图样式 */
:deep(.ant-carousel .slick-dots) {
  bottom: 20px;
}

:deep(.ant-carousel .slick-dots li button) {
  background: white;
  opacity: 0.5;
}

:deep(.ant-carousel .slick-dots li.slick-active button) {
  opacity: 1;
}
</style>