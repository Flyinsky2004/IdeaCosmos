<script setup>
import {ref} from "vue";
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
    path: '/admin/dashboard',
    name: '管理面板',
    icon: `<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="size-6">
  <path stroke-linecap="round" stroke-linejoin="round" d="M3.75 3v11.25A2.25 2.25 0 0 0 6 16.5h2.25M3.75 3h-1.5m1.5 0h16.5m0 0h1.5m-1.5 0v11.25A2.25 2.25 0 0 1 18 16.5h-2.25m-7.5 0h7.5m-7.5 0-1 3m8.5-3 1 3m0 0 .5 1.5m-.5-1.5h-9.5m0 0-.5 1.5M9 11.25v1.5M12 11.25v1.5m3-1.5v1.5" />
</svg>
`
  }, {
    path: '/admin/users',
    name: '用户管理',
    icon: `<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="size-6">
  <path stroke-linecap="round" stroke-linejoin="round" d="M15 19.128a9.38 9.38 0 0 0 2.625.372 9.337 9.337 0 0 0 4.121-.952 4.125 4.125 0 0 0-7.533-2.493M15 19.128v-.003c0-1.113-.285-2.16-.786-3.07M15 19.128v.106A12.318 12.318 0 0 1 8.624 21c-2.331 0-4.512-.645-6.374-1.766l-.001-.109a6.375 6.375 0 0 1 11.964-3.07M12 6.375a3.375 3.375 0 1 1-6.75 0 3.375 3.375 0 0 1 6.75 0Zm8.25 2.25a2.625 2.625 0 1 1-5.25 0 2.625 2.625 0 0 1 5.25 0Z" />
</svg>
`
  }, {
    path: '/admin/chapters',
    name: '篇章管理',
    icon: `<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="size-6">
  <path stroke-linecap="round" stroke-linejoin="round" d="M12 6.042A8.967 8.967 0 0 0 6 3.75c-1.052 0-2.062.18-3 .512v14.25A8.987 8.987 0 0 1 6 18c2.305 0 4.408.867 6 2.292m0-14.25a8.966 8.966 0 0 1 6-2.292c1.052 0 2.062.18 3 .512v14.25A8.987 8.987 0 0 0 18 18a8.967 8.967 0 0 0-6 2.292m0-14.25v14.25" />
</svg>
`
  }, {
    path: '/admin/projects',
    name: '项目管理',
    icon: `<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="size-6">
  <path stroke-linecap="round" stroke-linejoin="round" d="m21 7.5-9-5.25L3 7.5m18 0-9 5.25m9-5.25v9l-9 5.25M3 7.5l9 5.25M3 7.5v9l9 5.25m0-9v9" />
</svg>
`
  }, {
    path: '/admin/statistics',
    name: '数据统计',
    icon: `<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="size-6">
  <path stroke-linecap="round" stroke-linejoin="round" d="M3 13.125C3 12.504 3.504 12 4.125 12h2.25c.621 0 1.125.504 1.125 1.125v6.75C7.5 20.496 6.996 21 6.375 21h-2.25A1.125 1.125 0 0 1 3 19.875v-6.75ZM9.75 8.625c0-.621.504-1.125 1.125-1.125h2.25c.621 0 1.125.504 1.125 1.125v11.25c0 .621-.504 1.125-1.125 1.125h-2.25a1.125 1.125 0 0 1-1.125-1.125V8.625ZM16.5 4.125c0-.621.504-1.125 1.125-1.125h2.25C20.496 3 21 3.504 21 4.125v15.75c0 .621-.504 1.125-1.125 1.125h-2.25a1.125 1.125 0 0 1-1.125-1.125V4.125Z" />
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
            class="rounded-full border-2 border-red-500 dark:border-red-400 object-cover transition-all duration-300" 
            :class="sidebarExpanded ? 'w-16 h-16' : 'w-10 h-10'" 
            :src="BACKEND_DOMAIN + userStore.user.avatar" 
            alt="用户头像"
          >
          <div class="absolute bottom-0 right-0 w-3 h-3 rounded-full bg-red-500 border-2 border-white dark:border-gray-800"></div>
        </div>
        <div 
          v-if="sidebarExpanded" 
          class="mt-2 w-full text-center overflow-hidden transition-all duration-300"
          :class="textVisible ? 'opacity-100 max-h-20' : 'opacity-0 max-h-0'"
        >
          <span class="font-medium truncate">{{ userStore.user.username }}</span>
          <span class="text-red-500 text-xs block">(管理员)</span>
        </div>
      </div>

      <!-- 导航菜单 -->
      <div class="flex-1 overflow-y-auto scrollbar-hide">
        <nav class="flex flex-col gap-1">
          <Tooltip v-for="cat in category" :key="cat.path" :title="!sidebarExpanded ? cat.name : ''" placement="right">
            <div
                :class="[
                  route.path === cat.path 
                    ? 'bg-red-50 dark:bg-gray-800 text-red-600 dark:text-red-400 border-l-4 border-red-600 dark:border-red-400' 
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
