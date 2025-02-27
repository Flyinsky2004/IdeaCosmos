<script setup>
import {reactive, ref} from "vue";
import ThemeSwitch from "@/components/button/ThemeSwitcher.vue";
import {useUserStore} from "@/stores/user.js";
import router from "@/router/index.js";
import {message, Tooltip} from "ant-design-vue";
import {useRoute} from "vue-router";
import {useThemeStore} from "@/stores/theme.js";
import { BACKEND_DOMAIN } from "@/util/VARRIBLES";
const [messageApi, contextHolder] = message.useMessage();
const route = useRoute()
const userStore = useUserStore()

// 使用ref替代reactive，更适合简单状态管理
const sidebarExpanded = ref(false);
const textVisible = ref(false);

// 鼠标悬浮展开侧边栏
const expandSidebar = () => {
  sidebarExpanded.value = true;
  // 延迟显示文本，创建平滑过渡效果
  setTimeout(() => {
    textVisible.value = true;
  }, 150);
}

// 鼠标离开收起侧边栏
const collapseSidebar = () => {
  textVisible.value = false;
  // 延迟收起侧边栏，让文本先消失
  setTimeout(() => {
    sidebarExpanded.value = false;
  }, 100);
}

const category = [
  {
    path: '/workspace/dataAnlysis',
    name: '热点看板',
    icon: `<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="size-6">
  <path stroke-linecap="round" stroke-linejoin="round" d="M10.5 6a7.5 7.5 0 1 0 7.5 7.5h-7.5V6Z" />
  <path stroke-linecap="round" stroke-linejoin="round" d="M13.5 10.5H21A7.5 7.5 0 0 0 13.5 3v7.5Z" />
</svg>
`
  }, {
    path: '/workspace/projects',
    name: '创剧空间',
    icon: `<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="size-6">
  <path stroke-linecap="round" stroke-linejoin="round" d="m21 7.5-2.25-1.313M21 7.5v2.25m0-2.25-2.25 1.313M3 7.5l2.25-1.313M3 7.5l2.25 1.313M3 7.5v2.25m9 3 2.25-1.313M12 12.75l-2.25-1.313M12 12.75V15m0 6.75 2.25-1.313M12 21.75V19.5m0 2.25-2.25-1.313m0-16.875L12 2.25l2.25 1.313M21 14.25v2.25l-2.25 1.313m-13.5 0L3 16.5v-2.25" />
</svg>

`
  }, {
    path: '/workspace/personalInfo',
    name: '个人信息',
    icon: `
    <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="size-6">
  <path stroke-linecap="round" stroke-linejoin="round" d="M15.75 6a3.75 3.75 0 1 1-7.5 0 3.75 3.75 0 0 1 7.5 0ZM4.501 20.118a7.5 7.5 0 0 1 14.998 0A17.933 17.933 0 0 1 12 21.75c-2.676 0-5.216-.584-7.499-1.632Z" />
</svg>

`
  }, {
    path: '/workspace/teams',
    name: '我的团队',
    icon: `<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="size-6">
  <path stroke-linecap="round" stroke-linejoin="round" d="M18 18.72a9.094 9.094 0 0 0 3.741-.479 3 3 0 0 0-4.682-2.72m.94 3.198.001.031c0 .225-.012.447-.037.666A11.944 11.944 0 0 1 12 21c-2.17 0-4.207-.576-5.963-1.584A6.062 6.062 0 0 1 6 18.719m12 0a5.971 5.971 0 0 0-.941-3.197m0 0A5.995 5.995 0 0 0 12 12.75a5.995 5.995 0 0 0-5.058 2.772m0 0a3 3 0 0 0-4.681 2.72 8.986 8.986 0 0 0 3.74.477m.94-3.197a5.971 5.971 0 0 0-.94 3.197M15 6.75a3 3 0 1 1-6 0 3 3 0 0 1 6 0Zm6 3a2.25 2.25 0 1 1-4.5 0 2.25 2.25 0 0 1 4.5 0Zm-13.5 0a2.25 2.25 0 1 1-4.5 0 2.25 2.25 0 0 1 4.5 0Z" />
</svg>

`
  }
]
</script>

<template>
  <contextHolder/>
  <!-- 侧边栏 -->
  <div
      class="h-full fixed top-0 left-0 z-10 transition-all duration-300 text-theme-switch animate__animated animate__fadeIn shadow-lg"
      :class="[
        sidebarExpanded ? 'sidebar-expanded' : 'sidebar-collapsed',
        'bg-white/90 dark:bg-gray-900/90 backdrop-blur-md'
      ]"
      @mouseenter="expandSidebar"
      @mouseleave="collapseSidebar"
  >
    <div class="h-full flex flex-col p-3 gap-3">
      <!-- 用户信息区域 -->
      <div class="flex flex-col items-center bg-white dark:bg-gray-800 pt-3 pb-3 rounded-xl shadow-sm">
        <div class="relative">
          <img 
            class="rounded-full border-2 border-blue-500 dark:border-purple-500 object-cover transition-all duration-300" 
            :class="sidebarExpanded ? 'w-16 h-16' : 'w-10 h-10'" 
            :src="BACKEND_DOMAIN + userStore.user.avatar" 
            alt="用户头像"
          >
          <div class="absolute bottom-0 right-0 w-3 h-3 rounded-full bg-green-500 border-2 border-white dark:border-gray-800"></div>
        </div>
        <div 
          v-if="sidebarExpanded" 
          class="mt-2 w-full text-center overflow-hidden transition-all duration-300"
          :class="textVisible ? 'opacity-100 max-h-20' : 'opacity-0 max-h-0'"
        >
          <span class="font-medium truncate">{{ userStore.user.username }}</span>
        </div>
      </div>

      <!-- 导航菜单 -->
      <div class="flex-1 overflow-y-auto scrollbar-hide">
        <nav class="flex flex-col gap-1">
          <Tooltip v-for="cat in category" :key="cat.path" :title="!sidebarExpanded ? cat.name : ''" placement="right">
            <div
                :class="[
                  route.path === cat.path 
                    ? 'bg-blue-50 dark:bg-gray-800 text-blue-600 dark:text-purple-400 border-l-4 border-blue-600 dark:border-purple-400' 
                    : 'text-gray-700 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-800 border-l-4 border-transparent',
                ]"
                class="flex items-center p-2 rounded-lg cursor-pointer transition-all duration-200"
                @click="router.push(cat.path)"
            >
              <div class="flex-shrink-0" v-html="cat.icon"/>
              <div 
                v-if="sidebarExpanded" 
                class="ml-3 overflow-hidden transition-all duration-300"
                :class="textVisible ? 'opacity-100 max-w-40' : 'opacity-0 max-w-0'"
              >
                <span class="font-medium whitespace-nowrap">{{ cat.name }}</span>
              </div>
            </div>
          </Tooltip>
        </nav>
      </div>

      <!-- 底部操作区 -->
      <div class="border-t border-gray-200 dark:border-gray-700 pt-2">
        <div :class="sidebarExpanded ? 'flex items-center justify-between' : 'flex flex-col gap-2'"
             class="bg-white dark:bg-gray-800 p-2 rounded-lg shadow-sm transition-all duration-300">
          <ThemeSwitch class="flex-shrink-0" />
          <Tooltip :title="!sidebarExpanded ? '退出' : ''" placement="right">
            <div 
              class="flex items-center justify-center p-2 cursor-pointer bg-red-500 hover:bg-red-600 text-white rounded-lg transition-all duration-300"
              @click="router.push('/')"
            >
              <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="size-5">
                <path stroke-linecap="round" stroke-linejoin="round" d="M15.75 9V5.25A2.25 2.25 0 0 0 13.5 3h-6a2.25 2.25 0 0 0-2.25 2.25v13.5A2.25 2.25 0 0 0 7.5 21h6a2.25 2.25 0 0 0 2.25-2.25V15M12 9l-3 3m0 0 3 3m-3-3h12.75" />
              </svg>
              <div 
                v-if="sidebarExpanded" 
                class="ml-2 overflow-hidden transition-all duration-300"
                :class="textVisible ? 'opacity-100 max-w-20' : 'opacity-0 max-w-0'"
              >
                <span class="whitespace-nowrap">退出</span>
              </div>
            </div>
          </Tooltip>
        </div>
      </div>
    </div>
  </div>

  <!-- 主内容区域 -->
  <div class="w-full h-full transition-all duration-300"
       :class="sidebarExpanded ? 'pl-64' : 'pl-20'"
  >
    <div :class="useThemeStore().currentTheme === 'light' ? 'lightBKG' : 'darkBKG'"
        class="h-screen overflow-auto w-full container1 text-theme-switch p-0 relative bg-slate-50 dark:bg-slate-950 before:absolute before:inset-0 before:bg-glow-effect before:blur-3xl before:z-[-1]">
      <div class="min-h-screen w-full dark:bg-gray-800/30 bg-gray-100/50 backdrop-blur-xl m-0 p-2">
        <RouterView/>
      </div>
    </div>
  </div>
</template>

<style scoped>
.sidebar-expanded {
  width: 16rem; /* 64px */
}

.sidebar-collapsed {
  width: 5rem; /* 20px */
}

.scrollbar-hide::-webkit-scrollbar {
  display: none;
}

.scrollbar-hide {
  -ms-overflow-style: none;
  scrollbar-width: none;
}

.container1{
  background-size: 100% 100%; /* 确保背景图片覆盖整个容器 */
  background-attachment: fixed; /* 背景图片固定不动 */
  background-position: center; /* 背景图片居中 */
  background-blend-mode: multiply;
}
.darkBKG{
  background-image: url('@/assets/img/darkbkg.png'); /* 使用背景图片 */
  background-color: rgba(80, 79, 79, 0.3);
}
.lightBKG{
  background-image: url('@/assets/img/lightbkg.png'); /* 使用背景图片 */
  background-color: rgba(221, 221, 221, 0.1);
}
</style>