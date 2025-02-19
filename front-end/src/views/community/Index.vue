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
</script>

<template>
  <div class="w-full space-y-6">
    <!-- 轮播图区域 -->
    <div class="w-full h-64 bg-white dark:bg-zinc-900 border theme-border rounded-xl backdrop-blur-sm">
      <!-- 这里可以添加轮播图组件 -->
    </div>

    <!-- 内容列表 -->
    <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-6">
      <div 
        v-for="project in options.projects" 
        :key="project.ID" 
        @click="router.push('/community/project/' + project.ID)"
        class="group hover:cursor-pointer bg-white/90 dark:bg-zinc-900 border theme-border rounded-xl overflow-hidden hover:shadow-lg dark:hover:shadow-zinc-800 transition-all duration-300 hover:-translate-y-1"
      >
        <div class="aspect-video overflow-hidden bg-gray-100 dark:bg-zinc-800">
          <img 
            :src="imagePrefix + project.cover_image || '/default-cover.jpg'" 
            :alt="project.project_name"
            class="w-full h-full object-cover group-hover:scale-105 transition-transform duration-300"
          >
        </div>
        <div class="p-4 space-y-2">
          <h3 class="text-gray-900 dark:text-gray-100 text-base font-medium line-clamp-2 group-hover:text-transparent group-hover:bg-clip-text group-hover:bg-gradient-to-r group-hover:from-blue-500 group-hover:to-cyan-500 transition-colors">
            {{ project.project_name }}
          </h3>
          <div class="flex items-center gap-3 text-sm">
            <span class="bg-gradient-to-r from-blue-500 to-cyan-500 bg-clip-text text-transparent font-medium">
              {{ project.team?.username || '未知作者' }}
            </span>
            <div class="flex items-center gap-4 text-gray-500 dark:text-zinc-400">
              <span class="flex items-center gap-1">
                <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor" class="w-4 h-4">
                  <path d="M10 12.5a2.5 2.5 0 100-5 2.5 2.5 0 000 5z" />
                  <path fill-rule="evenodd" d="M.664 10.59a1.651 1.651 0 010-1.186A10.004 10.004 0 0110 3c4.257 0 7.893 2.66 9.336 6.41.147.381.146.804 0 1.186A10.004 10.004 0 0110 17c-4.257 0-7.893-2.66-9.336-6.41zM14 10a4 4 0 11-8 0 4 4 0 018 0z" clip-rule="evenodd" />
                </svg>
                {{ project.watches }}
              </span>
              <span class="flex items-center gap-1">
                <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor" class="w-4 h-4">
                  <path d="M9.653 16.915l-.005-.003-.019-.01a20.759 20.759 0 01-1.162-.682 22.045 22.045 0 01-2.582-1.9C4.045 12.733 2 10.352 2 7.5a4.5 4.5 0 018-2.828A4.5 4.5 0 0118 7.5c0 2.852-2.044 5.233-3.885 6.82a22.049 22.049 0 01-3.744 2.582l-.019.01-.005.003h-.002a.739.739 0 01-.69.001l-.002-.001z" />
                </svg>
                {{ project.favorites }}
              </span>
              <span class="flex items-center gap-1">
                <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor" class="w-4 h-4">
                  <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm.75-13a.75.75 0 00-1.5 0v5c0 .414.336.75.75.75h4a.75.75 0 000-1.5h-3.25V5z" clip-rule="evenodd" />
                </svg>
                {{ new Date(project.CreatedAt).toLocaleDateString() }}
              </span>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
/* 可以删除之前的样式，因为都使用了 Tailwind 类 */
</style>