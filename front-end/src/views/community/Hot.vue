<script setup>
import { ref, onMounted } from 'vue'
import { get } from '@/util/request'
import { message } from 'ant-design-vue'
import { BACKEND_DOMAIN } from '@/util/VARRIBLES'
import router from '@/router'
import SpinLoaderLarge from '@/components/spinLoaderLarge.vue'

const hotProjects = ref([])
const loading = ref(true)

// 获取热门项目
const fetchHotProjects = async () => {
  loading.value = true
  get('/api/public/hot-projects', 
    {}, 
    (msg, data) => {
      hotProjects.value = data
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

// 查看项目详情
const viewProject = (projectId) => {
  router.push(`/community/project/${projectId}`)
}

onMounted(() => {
  fetchHotProjects()
})
</script>

<template>
  <div class="min-h-screen py-8 bg-gray-50 dark:bg-zinc-900/70 rounded-xl">
    <div class="max-w-7xl mx-auto px-4">
      <!-- 页面标题 -->
      <div class="flex items-center gap-3 mb-8 animate__animated animate__fadeIn">
        <div class="h-8 w-1 bg-gradient-to-b from-blue-500 to-cyan-500 rounded-full"></div>
        <h1 class="text-3xl font-bold text-gray-900 dark:text-gray-100">
          热门项目
        </h1>
      </div>

      <!-- 加载状态 -->
      <div v-if="loading" class="flex items-center justify-center py-12">
        <SpinLoaderLarge />
      </div>

      <!-- 项目列表 -->
      <div v-else class="space-y-6">
        <div v-for="(project, index) in hotProjects" 
          :key="project.ID"
          class="group bg-white dark:bg-zinc-800 rounded-xl border theme-border hover:shadow-xl transition-all duration-300 hover:-translate-y-1 animate__animated animate__fadeIn"
          :class="'animate__delay-' + (index + 1) + 's'"
        >
          <div class="flex gap-6">
            <!-- 项目封面和排名 -->
            <div class="w-72 shrink-0 relative">
              <!-- 排名标识 -->
              <div 
                class="absolute -top-4 -left-4 z-10 w-16 h-16 flex items-center justify-center"
                :class="{
                  'bg-gradient-to-br from-amber-400 to-yellow-500': index === 0,
                  'bg-gradient-to-br from-gray-300 to-gray-400': index === 1,
                  'bg-gradient-to-br from-amber-600 to-amber-700': index === 2,
                  'bg-gradient-to-br from-blue-400 to-blue-500': index > 2
                }"
                style="clip-path: polygon(50% 0%, 100% 25%, 100% 75%, 50% 100%, 0% 75%, 0% 25%)"
              >
                <div class="flex flex-col items-center text-white">
                  <!-- 前三名显示奖杯图标 -->
                  <svg 
                    v-if="index < 3" 
                    xmlns="http://www.w3.org/2000/svg" 
                    class="h-6 w-6 mb-0.5" 
                    fill="none" 
                    viewBox="0 0 24 24" 
                    stroke="currentColor"
                    :class="{
                      'text-yellow-200': index === 0,
                      'text-gray-200': index === 1,
                      'text-amber-200': index === 2
                    }"
                  >
                    <path 
                      stroke-linecap="round" 
                      stroke-linejoin="round" 
                      stroke-width="2"
                      d="M5 3v4M3 5h4M6 17v4m-2-2h4m5-16l2.286 6.857L21 12l-5.714 2.143L13 21l-2.286-6.857L5 12l5.714-2.143L13 3z" 
                    />
                  </svg>
                  <!-- 排名数字 -->
                  <span 
                    class="font-bold text-lg"
                    :class="{
                      'bg-gradient-to-r from-yellow-200 to-amber-100 bg-clip-text text-transparent': index === 0,
                      'bg-gradient-to-r from-gray-200 to-white bg-clip-text text-transparent': index === 1,
                      'bg-gradient-to-r from-amber-200 to-amber-100 bg-clip-text text-transparent': index === 2
                    }"
                  >
                    {{ index + 1 }}
                  </span>
                </div>
              </div>

              <div class="h-full rounded-l-xl overflow-hidden">
                <img 
                  :src="`${BACKEND_DOMAIN}uploads/${project.cover_image}`"
                  class="w-full h-full object-cover transform group-hover:scale-105 transition-transform duration-300"
                  alt="项目封面"
                />
              </div>
            </div>

            <!-- 项目内容 -->
            <div class="flex-1 min-w-0 py-6 pr-6">
              <!-- 头部信息 -->
              <div class="flex items-start justify-between gap-4 mb-4">
                <div>
                  <h2 class="text-2xl font-bold text-gray-900 dark:text-gray-100 group-hover:text-blue-500 transition-colors mb-2">
                    {{ project.project_name }}
                  </h2>
                  <div class="flex items-center gap-4 text-sm text-gray-600 dark:text-gray-400">
                    <div class="flex items-center gap-2">
                      <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                        <path stroke-linecap="round" stroke-linejoin="round" d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z" />
                      </svg>
                      <span>{{ project.Team?.username }}</span>
                    </div>
                    <div class="flex items-center gap-2">
                      <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                        <path stroke-linecap="round" stroke-linejoin="round" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
                        <path stroke-linecap="round" stroke-linejoin="round" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z" />
                      </svg>
                      <span>{{ project.watches }}</span>
                    </div>
                    <span class="px-2 py-1 bg-blue-100 dark:bg-blue-900/30 text-blue-700 dark:text-blue-300 rounded-md">
                      {{ project.types }}
                    </span>
                  </div>
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

              <!-- 项目描述 -->
              <p class="text-gray-600 dark:text-gray-400 mb-4 line-clamp-3">
                {{ project.social_story }}
              </p>

              <!-- 项目标签 -->
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
    </div>
  </div>
</template>

<style scoped>
.animate__animated {
  animation-duration: 0.5s;
}

/* 自定义延迟时间 */
.animate__delay-1s {
  animation-delay: 0.1s;
}
.animate__delay-2s {
  animation-delay: 0.2s;
}
.animate__delay-3s {
  animation-delay: 0.3s;
}
.animate__delay-4s {
  animation-delay: 0.4s;
}
.animate__delay-5s {
  animation-delay: 0.5s;
}
.animate__delay-6s {
  animation-delay: 0.6s;
}
.animate__delay-7s {
  animation-delay: 0.7s;
}
.animate__delay-8s {
  animation-delay: 0.8s;
}
.animate__delay-9s {
  animation-delay: 0.9s;
}
.animate__delay-10s {
  animation-delay: 1s;
}

/* 添加阴影效果 */
.absolute {
  filter: drop-shadow(0 4px 6px rgb(0 0 0 / 0.1));
}
</style>