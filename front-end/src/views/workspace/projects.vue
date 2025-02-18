<script setup>

import {onMounted, reactive} from "vue";
import {get} from "@/util/request.js";
import {parseDateTime} from "@/util/common.js";
import router from "@/router/index.js";
import {useProjectStore} from "@/stores/project.js";
import {imagePrefix} from "@/util/VARRIBLES.js";

const options = reactive({
  projects: [],
  isAddWindowOpen: false,
  nowHoverId: -1,
  quickActions: [
    {
      name: '新建项目',
      icon: 'M12 9v6m3-3H9m12 0a9 9 0 1 1-18 0 9 9 0 0 1 18 0Z',
      route: '/workspace/newProject',
      description: '从零开始创建一个新的创意项目'
    },
    {
      name: '项目模板',
      icon: 'M7.5 8.25h9m-9 3H12m-9.75 1.51c0 1.6 1.123 2.994 2.707 3.227 1.129.166 2.27.293 3.423.379.35.026.67.21.865.501L12 21l2.755-4.133a1.14 1.14 0 0 1 .865-.501 48.172 48.172 0 0 0 3.423-.379c1.584-.233 2.707-1.626 2.707-3.228V6.741c0-1.602-1.123-2.995-2.707-3.228A48.394 48.394 0 0 0 12 3c-2.392 0-4.744.175-7.043.513C3.373 3.746 2.25 5.14 2.25 6.741v6.018Z',
      route: '/workspace/templates',
      description: '使用预设模板快速开始创作'
    },
    {
      name: '导入项目',
      icon: 'M7.5 7.5h-.75A2.25 2.25 0 0 0 4.5 9.75v7.5a2.25 2.25 0 0 0 2.25 2.25h7.5a2.25 2.25 0 0 0 2.25-2.25v-7.5a2.25 2.25 0 0 0-2.25-2.25h-.75m-6 3.75 3 3m0 0 3-3m-3 3V1.5m6 9h.75a2.25 2.25 0 0 1 2.25 2.25v7.5a2.25 2.25 0 0 1-2.25 2.25h-7.5a2.25 2.25 0 0 1-2.25-2.25v-.75',
      route: '/workspace/import',
      description: '导入已有项目继续创作'
    }
  ]
})
const fetchMyProjects = () => {
  get('/api/project/myProjects', {},
      (message, data) => {
        options.projects = data;
      })
}

onMounted(() => {
  fetchMyProjects()
})
const moveIn = (id) => {
  options.nowHoverId = id;
}
const moveOut = () => {
  options.nowHoverId = -1;
}
const projectStore = useProjectStore()
const enterProject = (project) => {
  projectStore.setProject(project)
  localStorage.setItem("project", JSON.stringify(project))
  router.push('/workspace/editProject/index')
}
</script>

<template>
  <div class="flex flex-col gap-4 animate__animated animate__fadeIn p-6">
    <!-- 顶部欢迎区域 -->
    <div class="bg-white dark:bg-zinc-900 rounded-2xl p-6 shadow-sm border theme-border">
      <h1 class="text-3xl font-bold bg-gradient-to-r from-blue-500 to-cyan-500 bg-clip-text text-transparent mb-2">
        创剧空间
      </h1>
      <p class="text-gray-600 dark:text-gray-400 leading-relaxed">
        我们的团队功能支持您将创剧空间中的项目与您的团队绑定，使您的团队成员一同加入您的精彩内容创作！
        <br/>
        <span class="text-sm font-medium mt-2 inline-block">
          支持功能：项目多人分工协作 / 实时评论反馈 / 项目版本控制
        </span>
      </p>
    </div>

    <!-- 项目列表区域 -->
    <div class="bg-white dark:bg-zinc-900 rounded-2xl p-6 shadow-sm border theme-border">
      <div class="flex justify-between items-center mb-6">
        <h2 class="text-2xl font-bold text-gray-900 dark:text-gray-100">我的项目</h2>
        <button 
          @click="router.push('/workspace/newProject')"
          class="flex items-center gap-2 px-4 py-2 bg-gradient-to-r from-blue-500 to-cyan-500 text-white rounded-lg hover:opacity-90 transition-all"
        >
          <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
          </svg>
          新建项目
        </button>
      </div>

      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-6">
        <!-- 项目卡片 -->
        <div 
          v-for="project in options.projects" 
          :key="project.ID"
          class="group relative bg-gray-50 dark:bg-zinc-800/50 rounded-xl overflow-hidden border theme-border transition-all duration-300 hover:-translate-y-1 hover:shadow-lg"
          @mouseover="moveIn(project.ID)" 
          @mouseleave="moveOut()"
        >
          <!-- 项目封面 -->
          <div class="aspect-video w-full overflow-hidden bg-gray-100 dark:bg-zinc-800">
            <img 
              v-if="project.cover_image" 
              :src="imagePrefix + project.cover_image" 
              :alt="project.project_name"
              class="w-full h-full object-cover"
            />
            <div v-else class="w-full h-full flex items-center justify-center text-gray-400">
              <svg xmlns="http://www.w3.org/2000/svg" class="h-12 w-12" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z" />
              </svg>
            </div>
          </div>

          <!-- 项目信息 -->
          <div class="p-4">
            <h3 class="text-lg font-semibold text-gray-900 dark:text-gray-100 mb-1">
              {{ project.project_name }}
            </h3>
            <p class="text-sm text-gray-500 dark:text-gray-400 mb-2">
              团队：{{ project.team.username }}
            </p>
            <div class="flex items-center gap-4 text-sm text-gray-500 dark:text-gray-400">
              <span class="flex items-center gap-1">
                <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
                </svg>
                {{ parseDateTime(project.CreatedAt) }}
              </span>
            </div>
          </div>

          <!-- 悬浮时显示的操作按钮 -->
          <div 
            v-if="options.nowHoverId === project.ID"
            class="absolute inset-0 bg-black/60 flex items-center justify-center animate__animated animate__fadeIn animate__faster"
          >
            <button 
              @click="enterProject(project)"
              class="px-6 py-3 bg-white text-gray-900 rounded-lg font-medium hover:bg-gray-100 transition-colors flex items-center gap-2"
            >
              <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" d="M13 7l5 5m0 0l-5 5m5-5H6" />
              </svg>
              进入项目
            </button>
          </div>
        </div>

        <!-- 快捷操作卡片 -->
        <div 
          v-for="action in options.quickActions"
          :key="action.name"
          @click="router.push(action.route)"
          class="group bg-gray-50 dark:bg-zinc-800/50 rounded-xl border theme-border p-6 flex flex-col items-center justify-center gap-4 cursor-pointer hover:bg-gray-100 dark:hover:bg-zinc-800 transition-all duration-300 hover:-translate-y-1"
        >
          <svg xmlns="http://www.w3.org/2000/svg" class="h-12 w-12 text-blue-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" :d="action.icon"/>
          </svg>
          <div class="text-center">
            <h3 class="text-lg font-medium text-gray-900 dark:text-gray-100 mb-1">{{ action.name }}</h3>
            <p class="text-sm text-gray-500 dark:text-gray-400">{{ action.description }}</p>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>

</style>