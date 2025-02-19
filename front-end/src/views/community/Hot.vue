<script setup>
import { get } from '@/util/request';
import { message } from 'ant-design-vue';
import { onMounted, reactive } from 'vue'
import { imagePrefix } from '@/util/VARRIBLES';

const options = reactive({
  projects: [],
  pageIndex: 0
})

const fetchData = () => {
  get('/api/public/getHotProjects', {
    pageIndex: options.pageIndex
  }, (messageer, data) => {
    options.pageIndex += 10
    options.projects = data
  }, (messageer) => {
    message.warning(messageer)
  }, (messageer) => {
    message.error(messageer)
  })
}

onMounted(() => {
  fetchData()
})
</script>

<template>
  <div class="w-full space-y-6">
    <!-- 内容列表 -->
    <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-6">
      <div 
        v-for="project in options.projects" 
        :key="project.ID" 
        class="group bg-white/90 dark:bg-zinc-900 border theme-border rounded-xl overflow-hidden hover:shadow-lg dark:hover:shadow-zinc-800 transition-all duration-300 hover:-translate-y-1"
      >
        <!-- 使用与 Index.vue 相同的项目卡片结构 -->
        <!-- ... 项目卡片内容 ... -->
      </div>
    </div>
  </div>
</template> 