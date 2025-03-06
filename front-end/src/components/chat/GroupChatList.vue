<!--
 * @Author: Flyinsky w2084151024@gmail.com
 * @Description: 群组聊天列表组件
-->
<script setup>
import { ref, onMounted, defineEmits, watchEffect } from 'vue'
import { useUserStore } from '@/stores/user'
import { get, post, postJSON } from '@/util/request'
import { message } from 'ant-design-vue'
import SpinLoaderLarge from '@/components/spinLoaderLarge.vue'

const props = defineProps({
  activeGroupId: {
    type: [Number, null],
    default: null
  }
})

const emit = defineEmits(['select-group', 'create-group'])

const userStore = useUserStore()
const groups = ref([])
const loading = ref(true)
const pageNum = ref(1)
const pageSize = ref(10)
const total = ref(0)
const searchKeyword = ref('')
const showCreateModal = ref(false)
const newGroupForm = ref({
  name: '',
  description: '',
  avatarURL: ''
})
const hoveredGroupId = ref(null)

// 获取用户加入的群组列表
const fetchGroups = () => {
  loading.value = true
  get('/api/chat/groups', {
    pageNum: pageNum.value,
    pageSize: pageSize.value,
    keyword: searchKeyword.value
  }, (_, data) => {
    groups.value = data.groups
    total.value = data.total
    loading.value = false
  }, () => {
    loading.value = false
    message.error('获取群组列表失败')
  })
}

// 搜索群组
const searchGroups = () => {
  pageNum.value = 1
  fetchGroups()
}

// 切换页码
const changePage = (page) => {
  pageNum.value = page
  fetchGroups()
}

// 选择群组
const selectGroup = (group) => {
  emit('select-group', group)
}

// 创建新群组
const createGroup = () => {
  if (!newGroupForm.value.name) {
    message.error('群组名称不能为空')
    return
  }
  
  postJSON('/api/chat/groups', newGroupForm.value, (_, data) => {
    showCreateModal.value = false
    newGroupForm.value = {
      name: '',
      description: '',
      avatarURL: ''
    }
    fetchGroups()
    emit('create-group', data)
  }, () => {
    message.error('创建群组失败')
  })
}

// 格式化时间
const formatTime = (timestamp) => {
  const date = new Date(timestamp)
  const now = new Date()
  
  // 同一天只显示时间
  if (date.toDateString() === now.toDateString()) {
    return date.toLocaleTimeString('zh-CN', { hour: '2-digit', minute: '2-digit' })
  }
  
  // 一周内显示星期几
  const diffDays = Math.floor((now - date) / (24 * 60 * 60 * 1000))
  if (diffDays < 7) {
    const weekdays = ['周日', '周一', '周二', '周三', '周四', '周五', '周六']
    return weekdays[date.getDay()]
  }
  
  // 其他显示日期
  return date.toLocaleDateString('zh-CN', { month: '2-digit', day: '2-digit' })
}

// 在组件挂载时获取群组列表
onMounted(() => {
  fetchGroups()
})

// 监听搜索关键词的变化
watchEffect(() => {
  if (searchKeyword.value === '') {
    fetchGroups()
  }
})
</script>

<template>
  <div class="group-list-container h-full flex flex-col animate__animated animate__fadeIn">
    <!-- 搜索和创建按钮 -->
    <div class="flex items-center gap-3 p-4 border-b theme-border bg-white dark:bg-zinc-900">
      <div class="relative flex-1">
        <input
          v-model="searchKeyword"
          @keyup.enter="searchGroups"
          type="text"
          placeholder="搜索群组..."
          class="w-full px-4 py-2 pr-10 bg-gray-50 dark:bg-zinc-800/50 border theme-border rounded-lg focus:ring-2 focus:ring-blue-500/20 dark:text-gray-300 dark:placeholder-gray-500"
        >
        <button 
          @click="searchGroups"
          class="absolute right-3 top-1/2 -translate-y-1/2 text-gray-400 hover:text-blue-500 dark:text-gray-500 dark:hover:text-blue-400 transition-colors"
        >
          <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" d="M21 21l-5.197-5.197m0 0A7.5 7.5 0 105.196 5.196a7.5 7.5 0 0010.607 10.607z" />
          </svg>
        </button>
      </div>
      <button 
        @click="showCreateModal = true"
        class="p-2 rounded-lg bg-gradient-to-r from-blue-500 to-cyan-500 text-white hover:opacity-90 transition-all flex items-center"
        title="创建新群组"
      >
        <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" d="M12 4.5v15m7.5-7.5h-15" />
        </svg>
      </button>
    </div>

    <!-- 群组列表 -->
    <div class="flex-1 overflow-y-auto bg-gray-50 dark:bg-zinc-800/50">
      <div v-if="loading" class="flex justify-center py-8">
        <SpinLoaderLarge />
      </div>
      
      <div v-else-if="groups.length === 0" class="flex flex-col items-center justify-center py-12">
        <div class="text-5xl text-gray-300 dark:text-gray-600 mb-4">
          <svg xmlns="http://www.w3.org/2000/svg" class="h-16 w-16" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" d="M18 18.72a9.094 9.094 0 003.741-.479 3 3 0 00-4.682-2.72m.94 3.198l.001.031c0 .225-.012.447-.037.666A11.944 11.944 0 0112 21c-2.17 0-4.207-.576-5.963-1.584A6.062 6.062 0 016 18.719m12 0a5.971 5.971 0 00-.941-3.197m0 0A5.995 5.995 0 0012 12.75a5.995 5.995 0 00-5.058 2.772m0 0a3 3 0 00-4.681 2.72 8.986 8.986 0 003.74.477m.94-3.197a5.971 5.971 0 00-.94 3.197M15 6.75a3 3 0 11-6 0 3 3 0 016 0zm6 3a2.25 2.25 0 11-4.5 0 2.25 2.25 0 014.5 0zm-13.5 0a2.25 2.25 0 11-4.5 0 2.25 2.25 0 014.5 0z" />
          </svg>
        </div>
        <p class="text-gray-500 dark:text-gray-400 mb-4">暂无群组</p>
        <button 
          @click="showCreateModal = true"
          class="px-4 py-2 bg-gradient-to-r from-blue-500 to-cyan-500 text-white rounded-lg hover:opacity-90 transition-all flex items-center gap-2"
        >
          <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" d="M12 4.5v15m7.5-7.5h-15" />
          </svg>
          创建新群组
        </button>
      </div>
      
      <div v-else class="p-2 space-y-2">
        <div 
          v-for="group in groups" 
          :key="group.ID"
          @click="selectGroup(group)"
          @mouseenter="hoveredGroupId = group.ID"
          @mouseleave="hoveredGroupId = null"
          class="flex items-center p-3 rounded-lg cursor-pointer transition-all animate__animated animate__fadeIn"
          :class="[
            group.ID === activeGroupId 
              ? 'bg-gradient-to-r from-blue-500/10 to-cyan-500/10 dark:from-blue-500/20 dark:to-cyan-500/20 border-blue-500/50' 
              : hoveredGroupId === group.ID
                ? 'bg-gray-100 dark:bg-gray-700/50'
                : 'hover:bg-gray-100 dark:hover:bg-gray-700/50',
            'border theme-border'
          ]"
        >
          <!-- 群组头像 -->
          <div class="relative w-12 h-12 rounded-lg overflow-hidden bg-gradient-to-br from-blue-100 to-indigo-100 dark:from-blue-900/30 dark:to-indigo-900/30 mr-3 flex-shrink-0 border theme-border">
            <img 
              v-if="group.AvatarURL" 
              :src="group.AvatarURL" 
              class="w-full h-full object-cover"
              alt="群组头像"
            >
            <div v-else class="w-full h-full flex items-center justify-center text-blue-500 dark:text-blue-400 text-xl font-bold">
              {{ group.Name.charAt(0) }}
            </div>
          </div>
          
          <!-- 群组信息 -->
          <div class="flex-1 min-w-0">
            <div class="flex justify-between items-center mb-1">
              <h3 class="font-medium truncate dark:text-white">{{ group.Name }}</h3>
              <span class="text-xs text-gray-500 dark:text-gray-400 ml-2 flex-shrink-0">
                {{ formatTime(group.UpdatedAt) }}
              </span>
            </div>
            <p class="text-sm text-gray-500 dark:text-gray-400 truncate">
              {{ group.Description || `${group.MemberCount}个成员` }}
            </p>
          </div>
        </div>
      </div>
    </div>

    <!-- 分页 -->
    <div v-if="total > pageSize" class="flex justify-center py-3 bg-white dark:bg-zinc-900 border-t theme-border">
      <div class="flex gap-1">
        <button 
          v-for="page in Math.min(5, Math.ceil(total / pageSize))" 
          :key="page"
          @click="changePage(page)"
          class="min-w-[2rem] h-8 flex items-center justify-center rounded transition-colors text-sm"
          :class="pageNum === page 
            ? 'bg-gradient-to-r from-blue-500 to-cyan-500 text-white' 
            : 'bg-gray-50 dark:bg-zinc-800/50 text-gray-700 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-700/50'"
        >
          {{ page }}
        </button>
        <span v-if="Math.ceil(total / pageSize) > 5" class="w-8 h-8 flex items-center justify-center text-gray-500 dark:text-gray-400">...</span>
      </div>
    </div>

    <!-- 创建群组弹窗 -->
    <div v-if="showCreateModal" class="fixed inset-0 bg-black/50 z-50 flex items-center justify-center p-4 animate__animated animate__fadeIn">
      <div class="bg-white dark:bg-zinc-900 rounded-xl p-6 w-full max-w-md border theme-border animate__animated animate__fadeInUp animate__faster">
        <div class="flex justify-between items-center mb-6">
          <h2 class="text-2xl font-bold dark:text-white">创建新群组</h2>
          <button 
            @click="showCreateModal = false"
            class="p-2 text-gray-500 hover:text-gray-700 dark:text-gray-400 dark:hover:text-gray-200 transition-colors"
          >
            <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>
        
        <div class="space-y-4">
          <div>
            <label class="block text-sm font-medium mb-1 dark:text-gray-300">群组名称</label>
            <input 
              v-model="newGroupForm.name" 
              type="text" 
              class="w-full px-4 py-2 bg-gray-50 dark:bg-zinc-800/50 border theme-border rounded-lg focus:ring-2 focus:ring-blue-500/20 dark:text-gray-300"
              placeholder="输入群组名称"
            >
          </div>
          
          <div>
            <label class="block text-sm font-medium mb-1 dark:text-gray-300">群组描述</label>
            <textarea 
              v-model="newGroupForm.description" 
              class="w-full px-4 py-2 bg-gray-50 dark:bg-zinc-800/50 border theme-border rounded-lg focus:ring-2 focus:ring-blue-500/20 dark:text-gray-300 resize-none"
              placeholder="输入群组描述（选填）"
              rows="3"
            ></textarea>
          </div>
          
          <div>
            <label class="block text-sm font-medium mb-1 dark:text-gray-300">群组头像</label>
            <input 
              v-model="newGroupForm.avatarURL" 
              type="text" 
              class="w-full px-4 py-2 bg-gray-50 dark:bg-zinc-800/50 border theme-border rounded-lg focus:ring-2 focus:ring-blue-500/20 dark:text-gray-300"
              placeholder="头像URL（选填）"
            >
          </div>
        </div>
        
        <div class="flex justify-end gap-3 mt-6">
          <button 
            @click="showCreateModal = false"
            class="px-4 py-2 bg-gray-100 text-gray-700 rounded-lg hover:bg-gray-200 dark:bg-gray-800 dark:text-gray-300 dark:hover:bg-gray-700 transition-colors"
          >
            取消
          </button>
          <button 
            @click="createGroup"
            class="px-4 py-2 bg-gradient-to-r from-blue-500 to-cyan-500 text-white rounded-lg hover:opacity-90 transition-all"
          >
            创建
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.group-list-container {
  background-color: var(--bg-content);
}

/* 动画相关样式 */
.animate__animated {
  animation-duration: 0.5s;
}

.animate__faster {
  animation-duration: 0.3s;
}

@keyframes fadeIn {
  from {
    opacity: 0;
  }
  to {
    opacity: 1;
  }
}

@keyframes fadeInUp {
  from {
    opacity: 0;
    transform: translateY(20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.animate__fadeIn {
  animation-name: fadeIn;
}

.animate__fadeInUp {
  animation-name: fadeInUp;
}
</style> 