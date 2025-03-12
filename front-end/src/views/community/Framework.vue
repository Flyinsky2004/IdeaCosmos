<script setup>
import { ref, computed, watch, onMounted, nextTick } from 'vue'
import { RouterView, useRoute, useRouter } from 'vue-router'
import ThemeSwitcher from '@/components/button/ThemeSwitcher.vue'
import { useThemeStore } from '@/stores/theme'
import { useUserStore } from '@/stores/user'
import { BACKEND_DOMAIN } from '@/util/VARRIBLES'
import logo from '@/assets/img/logo.webp'
const route = useRoute()
const router = useRouter()
const searchKeyword = ref('')

const activeTab = computed(() => {
  const path = route.path
  if (path === '/community') return '推荐'
  if (path.includes('/hot')) return '热门'
  if (path.includes('/categories')) return '分类'
  if (path.includes('/search')) return '搜索'
  return ''
})

const tabs = [
  { name: '推荐', path: '/community' },
  { name: '热门', path: '/community/hot' },
  { name: '分类', path: '/community/categories' }
]

const handleSearch = () => {
  if (searchKeyword.value.trim()) {
    router.push(`/community/search/${encodeURIComponent(searchKeyword.value.trim())}`)
  }
}

// 添加聊天和通知相关的方法
const navigateToChat = () => {
  router.push('/community/chat')
}

const navigateToNotifications = () => {
  router.push('/community/notifications')
}

// 添加通知数量
const unreadNotifications = ref(3)

const userStore = useUserStore()

const tabRefs = ref([])
const activeTabPosition = ref(0)
const activeTabWidth = ref(0)

// 更新选中标签的位置和宽度
const updateActiveTabIndicator = () => {
  const activeIndex = tabs.findIndex(tab => tab.name === activeTab.value)
  if (activeIndex >= 0 && tabRefs.value[activeIndex]) {
    const el = tabRefs.value[activeIndex].$el
    activeTabPosition.value = el.offsetLeft
    activeTabWidth.value = el.offsetWidth
  }
}

// 监听路由变化更新指示器
watch(activeTab, () => {
  nextTick(() => {
    updateActiveTabIndicator()
  })
})

// 初始化时设置指示器位置
onMounted(() => {
  nextTick(() => {
    updateActiveTabIndicator()
  })
})

// 监听窗口大小变化重新计算位置
window.addEventListener('resize', () => {
  updateActiveTabIndicator()
})
</script>

<template>
  <div :class="useThemeStore().currentTheme === 'light' ? 'lightBKG' : 'darkBKG'" class="min-h-screen container1 bg-gray-50 dark:bg-black">
    <!-- 顶部导航栏 -->
    <header class="fixed top-0 w-full bg-white/80 dark:bg-black border-b theme-border shadow-sm z-50 backdrop-blur-sm">
      <div class="max-w-7xl mx-auto h-16 flex items-center justify-between px-4 sm:px-6">
        <div class="flex items-center gap-12">
          <div class="flex flex-nowrap p-4 cursor-pointer" @click="router.push('/')">
      <img :src="logo" class="w-8 h-8 rounded-full"/>
      <span class="my-auto ml-4 dark:text-white">创剧星球</span>
    </div>
          <nav class="flex gap-8 relative">
            <!-- 背景动画条 -->
            <div 
              class="absolute bottom-0 h-0.5 bg-gradient-to-r from-blue-500 to-cyan-500 transition-all duration-300"
              :style="{
                left: activeTabPosition + 'px',
                width: activeTabWidth + 'px'
              }"
            />
            
            <router-link
              v-for="tab in tabs"
              :key="tab.name"
              :to="tab.path"
              class="relative px-2 py-1 text-gray-700 dark:text-gray-300 hover:text-blue-500 dark:hover:text-blue-400 transition-colors"
              :class="{ 
                'text-transparent bg-clip-text bg-gradient-to-r from-blue-500 to-cyan-500': activeTab === tab.name
              }"
              ref="tabRefs"
            >
              {{ tab.name }}
            </router-link>
          </nav>
        </div>
        
        <div class="flex items-center gap-4">
          <!-- 搜索框 -->
          <div class="relative">
            <a-input
              v-model="searchKeyword"
              @keyup.enter="handleSearch"
              type="text"
              placeholder="搜索项目..."
              class="w-64 px-4 py-2 pr-10 bg-gray-100 dark:bg-zinc-800 border-0 rounded-lg focus:ring-2 focus:ring-blue-500/20 dark:text-gray-300 dark:placeholder-gray-500"
            />
            <button 
              @click="handleSearch"
              class="absolute right-2 top-1/2 -translate-y-1/2 text-gray-400 hover:text-blue-500 dark:text-gray-500 dark:hover:text-blue-400"
            >
              <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" d="M21 21l-5.197-5.197m0 0A7.5 7.5 0 105.196 5.196a7.5 7.5 0 0010.607 10.607z" />
              </svg>
            </button>
          </div>
          
          <!-- 聊天按钮 -->
          <button 
            @click="navigateToChat"
            class="relative p-2 text-gray-600 hover:text-blue-500 dark:text-gray-400 dark:hover:text-blue-400 transition-colors rounded-full hover:bg-gray-100 dark:hover:bg-gray-800"
            title="聊天"
          >
            <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z" />
            </svg>
          </button>
          
          <!-- 通知按钮 -->
          <button 
            @click="navigateToNotifications"
            class="relative p-2 text-gray-600 hover:text-blue-500 dark:text-gray-400 dark:hover:text-blue-400 transition-colors rounded-full hover:bg-gray-100 dark:hover:bg-gray-800"
            title="通知"
          >
            <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 17h5l-1.405-1.405A2.032 2.032 0 0118 14.158V11a6.002 6.002 0 00-4-5.659V5a2 2 0 10-4 0v.341C7.67 6.165 6 8.388 6 11v3.159c0 .538-.214 1.055-.595 1.436L4 17h5m6 0v1a3 3 0 11-6 0v-1m6 0H9" />
            </svg>
            
            <!-- 未读通知数量 -->
            <span 
              v-if="unreadNotifications > 0" 
              class="absolute top-0 right-0 inline-flex items-center justify-center px-2 py-1 text-xs font-bold leading-none text-white transform translate-x-1/2 -translate-y-1/2 bg-red-500 rounded-full"
            >
              {{ unreadNotifications > 99 ? '99+' : unreadNotifications }}
            </span>
          </button>

          <!-- AI对话按钮 -->
          <button 
            @click="router.push('/community/chatS')"
            class="relative p-2 text-gray-600 hover:text-blue-500 dark:text-gray-400 dark:hover:text-blue-400 transition-colors rounded-full hover:bg-gray-100 dark:hover:bg-gray-800"
            title="AI对话"
          >
            <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9.75 3.104c-.767.217-1.5.5-2.17.87m0 0a10.07 10.07 0 0 1 2.17-.87M3 13h8m0 0v7m0-7c-3.866 0-7-3.134-7-7m7 7c3.866 0 7-3.134 7-7m-7 7v7m7-7h8m-8-7V3m0 0c3.866 0 7 3.134 7 7m-7-7V3M3 13c0 3.866 3.134 7 7 7m0 0c3.866 0 7-3.134 7-7" />
            </svg>
          </button>
          
          <ThemeSwitcher />
          <img 
            v-if="userStore.isLogin" 
            class="w-8 h-8 rounded-full" 
            :src="BACKEND_DOMAIN + userStore.user.avatar"
          >
        </div>
      </div>
    </header>

    <!-- 子路由渲染区域 -->
    <main class="dark:bg-gray-800/30 bg-gray-100/50 backdrop-blur-xl mx-auto pt-20 px-4 sm:px-6 pb-8">
      <RouterView />
    </main>
  </div>
</template>

<style scoped>
.community-framework {
  min-height: 100vh;
  background-color: #f1f2f3;
}

.header {
  position: fixed;
  top: 0;
  width: 100%;
  height: 64px;
  background-color: #fff;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  z-index: 100;
}

.header-inner {
  max-width: 1280px;
  margin: 0 auto;
  height: 100%;
  display: flex;
  align-items: center;
  padding: 0 24px;
}

.logo {
  font-size: 20px;
  font-weight: bold;
  color: #00aeec;
  margin-right: 48px;
}

.nav-tabs {
  display: flex;
  gap: 32px;
}

.tab-item {
  font-size: 16px;
  color: #18191c;
  text-decoration: none;
  padding: 8px 0;
  position: relative;
}

.tab-item.active::after {
  content: '';
  position: absolute;
  bottom: 0;
  left: 0;
  width: 100%;
  height: 2px;
  background-color: #00aeec;
}

.main-content {
  max-width: 1280px;
  margin: 84px auto 20px;
  padding: 0 24px;
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

/* 保留这些样式 */
nav a {
  position: relative;
  transition: all 0.3s ease;
}

nav a:hover {
  transform: translateY(-1px);
}

/* 优化动画性能 */
nav a {
  will-change: transform;
}
</style>