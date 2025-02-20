<template>
  <!-- 最外层容器使用 flex 布局 -->
  <div class="flex flex-col min-h-screen">
    <!-- 固定顶栏 -->
    <header class="sticky top-0 z-50 bg-white dark:bg-zinc-800 border-b theme-border">
      <div class="max-w-7xl mx-auto px-4">
        <div class="flex items-center justify-between h-16">
          <!-- 左侧 Logo 和导航 -->
          <div class="flex items-center gap-8">
            <router-link to="/" class="text-xl font-bold text-gray-900 dark:text-gray-100">
              创剧星球
            </router-link>
            <nav class="hidden md:flex items-center gap-4">
              <router-link 
                v-for="item in navItems" 
                :key="item.path"
                :to="item.path"
                class="px-3 py-2 rounded-lg text-gray-600 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-zinc-700/50 transition-colors"
                active-class="text-blue-500 bg-blue-50 dark:bg-blue-900/20 dark:text-blue-400"
              >
                {{ item.name }}
              </router-link>
            </nav>
          </div>

          <!-- 右侧用户信息 -->
          <div class="flex items-center gap-4">
            <!-- 现有的用户信息组件 -->
          </div>
        </div>
      </div>
    </header>

    <!-- 主要内容区域 -->
    <main class="flex-1 bg-gray-100 dark:bg-zinc-900">
      <router-view></router-view>
    </main>
  </div>
</template>

<script setup>
const navItems = [
  { name: '社区', path: '/community' },
  { name: '工作台', path: '/workspace' },
  // ... 其他导航项
]
</script>

<style>
/* 添加平滑滚动效果 */
html {
  scroll-behavior: smooth;
}

/* 确保固定顶栏有阴影效果 */
header {
  box-shadow: 0 1px 3px 0 rgb(0 0 0 / 0.1), 0 1px 2px -1px rgb(0 0 0 / 0.1);
  backdrop-filter: blur(8px);
  background-color: rgba(255, 255, 255, 0.9);
}

/* 暗色模式下的顶栏 */
.dark header {
  box-shadow: 0 1px 3px 0 rgb(0 0 0 / 0.3);
  background-color: rgba(39, 39, 42, 0.9);
}
</style> 