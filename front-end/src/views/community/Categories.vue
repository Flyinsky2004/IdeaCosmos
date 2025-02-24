<script setup>
import { ref, reactive } from 'vue'
import { get } from '@/util/request'
import { message } from 'ant-design-vue'
import { imagePrefix } from '@/util/VARRIBLES'
import router from '@/router'
import SpinLoaderLarge from '@/components/spinLoaderLarge.vue'

// 使用 tagsData 作为分类列表
const categories = [
  "幽默", "悬疑", "黑暗", "科幻", "奇幻", "浪漫", "荒诞", "励志", "动作", "恐怖",
  "史诗", "温情", "讽刺", "文艺", "纪实", "冒险", "家庭剧", "超现实", "战争", "公路",
  "青春", "复仇", "政治", "犯罪", "悬幻", "心理", "怪诞", "温暖", "现实主义", "虚构",
  "哲学", "校园", "灾难", "武侠", "神秘", "励志成长", "古风", "穿越", "音乐", "童话"
].map(name => ({
  name,
  icon: 'M12 6.042A8.967 8.967 0 006 3.75c-1.052 0-2.062.18-3 .512v14.25A8.987 8.987 0 016 18c2.305 0 4.408.867 6 2.292m0-14.25a8.966 8.966 0 016-2.292c1.052 0 2.062.18 3 .512v14.25A8.987 8.987 0 0018 18a8.967 8.967 0 00-6 2.292m0-14.25v14.25'
}))

const selectedCategory = ref(null)
const options = reactive({
  projects: [],
  pageIndex: 0,
  hasMore: true
})

const loading = ref(false)

const fetchCategoryProjects = (category) => {
  selectedCategory.value = category
  options.pageIndex = 0
  options.hasMore = true
  loading.value = true
  
  get('/api/public/getCategoryProjects', {
    category: category.name,
    pageIndex: options.pageIndex
  }, (msg, data) => {
    options.projects = data
    options.hasMore = data.length === 12
    loading.value = false
  }, (msg) => {
    message.warning(msg)
    loading.value = false
  }, (msg) => {
    message.error(msg)
    loading.value = false
  })
}

// 加载更多项目
const loadMore = () => {
  if (!selectedCategory.value) return
  loading.value = true
  
  options.pageIndex++
  get('/api/public/getCategoryProjects', {
    category: selectedCategory.value.name,
    pageIndex: options.pageIndex
  }, (msg, data) => {
    options.projects.push(...data)
    options.hasMore = data.length === 12
    loading.value = false
  }, (msg) => {
    message.warning(msg)
    loading.value = false
  }, (msg) => {
    message.error(msg)
    loading.value = false
  })
}

// 查看项目详情
const viewProject = (projectId) => {
  router.push(`/community/project/${projectId}`)
}
</script>

<template>
  <div class="min-h-screen py-8 bg-gray-50 dark:bg-zinc-900/60 rounded-xl">
    <div class="max-w-7xl mx-auto px-4 space-y-8">
      <!-- 页面标题 -->
      <div class="flex items-center gap-3 animate__animated animate__fadeIn">
        <div class="h-8 w-1 bg-gradient-to-b from-blue-500 to-cyan-500 rounded-full"></div>
        <h1 class="text-3xl font-bold text-gray-900 dark:text-gray-100">
          分类浏览
        </h1>
      </div>

      <!-- 分类导航 -->
      <div class="grid grid-cols-2 sm:grid-cols-8 md:grid-cols-10 lg:grid-cols-12 gap-3 animate__animated animate__fadeIn animate__delay-1s">
        <button
          v-for="category in categories"
          :key="category.name"
          @click="fetchCategoryProjects(category)"
          :class="[
            'px-4 py-2 rounded-full border transition-all duration-300 text-sm font-medium',
            selectedCategory?.name === category.name
              ? 'border-blue-500 bg-blue-500 text-white'
              : 'theme-border bg-white dark:bg-zinc-800 text-gray-700 dark:text-gray-300 hover:border-blue-500 hover:text-blue-500 dark:hover:text-blue-400'
          ]"
        >
          {{ category.name }}
        </button>
      </div>

      <!-- 加载状态 -->
      <div v-if="loading" class="flex items-center justify-center py-12">
        <SpinLoaderLarge />
      </div>

      <!-- 空状态提示 -->
      <div 
        v-else-if="selectedCategory && options.projects.length === 0" 
        class="flex flex-col items-center justify-center py-12 text-gray-500 dark:text-gray-400"
      >
        <svg xmlns="http://www.w3.org/2000/svg" class="h-16 w-16 mb-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4" />
        </svg>
        <p class="text-lg">暂无{{ selectedCategory.name }}类作品</p>
      </div>

      <!-- 项目列表 -->
      <div v-else-if="selectedCategory" class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        <div 
          v-for="(project, index) in options.projects" 
          :key="project.ID"
          class="group bg-white dark:bg-zinc-800 rounded-xl border theme-border hover:shadow-xl transition-all duration-300 hover:-translate-y-1 overflow-hidden animate__animated animate__fadeIn"
          :style="`animation-delay: ${(index + 2) * 0.1}s`"
          @click="viewProject(project.ID)"
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
            </div>
          </div>
        </div>
      </div>

      <!-- 加载更多按钮 -->
      <div 
        v-if="options.projects.length > 0 && options.hasMore" 
        class="text-center animate__animated animate__fadeIn"
      >
        <button 
          @click="loadMore"
          class="px-6 py-2 bg-gradient-to-r from-blue-500 to-cyan-500 text-white rounded-lg hover:opacity-90 transition-opacity"
          :disabled="loading"
        >
          {{ loading ? '加载中...' : '加载更多' }}
        </button>
      </div>

      <!-- 没有更多数据的提示 -->
      <div
        v-if="options.projects.length > 0 && !options.hasMore"
        class="text-center text-gray-500 dark:text-gray-400 py-4"
      >
        没有更多作品了
      </div>
    </div>
  </div>
</template>

<style scoped>
/* 添加动画相关样式 */
.animate__animated {
  animation-duration: 0.5s;
}

.animate__delay-1s {
  animation-delay: 0.1s;
}

.animate__delay-2s {
  animation-delay: 0.2s;
}

.animate__delay-3s {
  animation-delay: 0.3s;
}
</style> 