<script setup>
import { ref } from 'vue'
import { RouterView } from 'vue-router'
import ThemeSwitcher from '@/components/button/ThemeSwitcher.vue'
import { useThemeStore } from '@/stores/theme'
const activeTab = ref('推荐')
const tabs = [
  { name: '推荐', path: '/community' },
  { name: '最新', path: '/community/latest' },
  { name: '热门', path: '/community/hot' }
]
</script>

<template>
  <div :class="useThemeStore().currentTheme === 'light' ? 'lightBKG' : 'darkBKG'" class="min-h-screen container1 bg-gray-50 dark:bg-black">
    <!-- 顶部导航栏 -->
    <header class="fixed top-0 w-full bg-white/80 dark:bg-black border-b theme-border shadow-sm z-50">
      <div class="max-w-7xl mx-auto h-16 flex items-center justify-between px-4 sm:px-6">
        <div class="flex items-center gap-12">
          <div class="text-xl font-bold bg-gradient-to-r from-blue-500 to-cyan-500 bg-clip-text text-transparent">
            创意宇宙
          </div>
          <nav class="flex gap-8">
            <router-link
              v-for="tab in tabs"
              :key="tab.name"
              :to="tab.path"
              class="relative px-2 py-1 text-gray-700 dark:text-gray-300 hover:text-blue-500 dark:hover:text-blue-400 transition-colors"
              :class="{ 'text-transparent bg-clip-text bg-gradient-to-r from-blue-500 to-cyan-500': activeTab === tab.name }"
            >
              {{ tab.name }}
              <div
                v-if="activeTab === tab.name"
                class="absolute bottom-0 left-0 w-full h-0.5 bg-gradient-to-r from-blue-500 to-cyan-500"
              />
            </router-link>
          </nav>
        </div>
        <ThemeSwitcher />
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
</style>