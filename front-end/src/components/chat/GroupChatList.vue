<!--
 * @Author: Flyinsky w2084151024@gmail.com
 * @Description: 群组聊天列表组件
-->
<script setup>
import { ref, onMounted, watch, computed } from 'vue'
import { get, post, postJSON } from '@/util/request'
import { message, Modal, Input } from 'ant-design-vue'
import { useUserStore } from '@/stores/user'

const props = defineProps({
  activeGroupId: {
    type: Number,
    default: null
  }
})

const emit = defineEmits(['select-group', 'create-group'])

const userStore = useUserStore()
const groups = ref([])
const loading = ref(false)
const searchText = ref('')
const createModalVisible = ref(false)
const createGroupForm = ref({
  name: '',
  description: '',
  avatarUrl: ''
})

// 获取群组列表
const fetchGroups = () => {
  loading.value = true
  get('/api/chat/groups', {}, (_, data) => {
    groups.value = data.groups
    loading.value = false
  }, () => {
    loading.value = false
    message.error('获取群组失败')
  })
}

// 过滤搜索结果
const filteredGroups = computed(() => {
  if (!searchText.value) return groups.value
  
  const query = searchText.value.toLowerCase()
  return groups.value.filter(group => 
    group.Name.toLowerCase().includes(query) || 
    (group.Description && group.Description.toLowerCase().includes(query))
  )
})

// 选择群组
const selectGroup = (group) => {
  emit('select-group', group)
}

// 打开创建群组模态框
const openCreateModal = () => {
  createGroupForm.value = {
    name: '',
    description: '',
    avatarUrl: ''
  }
  createModalVisible.value = true
}

// 创建群组
const handleCreateGroup = () => {
  if (!createGroupForm.value.name) {
    message.warning('请输入群组名称')
    return
  }
  
  // 提交创建请求
  postJSON('/api/chat/groups', {
    name: createGroupForm.value.name,
    description: createGroupForm.value.description,
    avatarUrl: createGroupForm.value.avatarUrl
  }, (_, data) => {
    message.success('群组创建成功')
    createModalVisible.value = false
    
    // 刷新群组列表或直接添加到列表
    if (data) {
      groups.value.unshift(data)
      emit('create-group', data)
    } else {
      fetchGroups()
    }
  })
}

// 组件挂载时获取群组列表
onMounted(() => {
  fetchGroups()
})
</script>

<template>
  <div class="group-chat-list h-full flex flex-col">
    <!-- 顶部搜索和创建 -->
    <div class="p-4 border-b theme-border">
      <div class="flex items-center gap-2 mb-4">
        <h2 class="text-lg font-bold text-gray-800 dark:text-gray-200">我的群组</h2>
        <button 
          @click="openCreateModal"
          class="ml-auto p-1.5 text-gray-500 hover:text-blue-500 dark:text-gray-400 dark:hover:text-blue-400 transition-colors rounded-lg hover:bg-gray-100 dark:hover:bg-gray-800"
          title="创建群组"
        >
          <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" d="M12 4.5v15m7.5-7.5h-15" />
          </svg>
        </button>
      </div>
      
      <Input 
        v-model:value="searchText" 
        placeholder="搜索群组" 
        allowClear
      >
        <template #prefix>
          <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
          </svg>
        </template>
      </Input>
    </div>
    
    <!-- 群组列表 -->
    <div class="flex-1 overflow-y-auto py-2">
      <!-- 加载中 -->
      <div v-if="loading" class="py-4 text-center text-gray-500 dark:text-gray-400">
        <div class="animate-spin h-6 w-6 border-2 border-blue-500 dark:border-blue-400 rounded-full border-t-transparent mx-auto mb-2"></div>
        <p>加载中...</p>
      </div>
      
      <!-- 群组列表 -->
      <div 
        v-else-if="filteredGroups.length > 0"
        class="space-y-1 px-2"
      >
        <div 
          v-for="group in filteredGroups" 
          :key="group.ID"
          class="flex items-center p-3 rounded-lg cursor-pointer transition-colors"
          :class="group.ID === activeGroupId ? 
            'bg-blue-100 text-blue-700 dark:bg-blue-900/40 dark:text-blue-300' : 
            'hover:bg-gray-100 dark:hover:bg-zinc-800/70'"
          @click="selectGroup(group)"
        >
          <!-- 群组头像 -->
          <div class="w-10 h-10 rounded-full overflow-hidden bg-gradient-to-br from-blue-100 to-indigo-100 dark:from-blue-900/30 dark:to-indigo-900/30 flex-shrink-0 border theme-border">
            <img 
              v-if="group.AvatarURL" 
              :src="group.AvatarURL" 
              class="w-full h-full object-cover"
              alt="群组头像"
            >
            <div v-else class="w-full h-full flex items-center justify-center text-blue-500 dark:text-blue-400 text-sm font-bold">
              {{ group.Name.charAt(0) }}
            </div>
          </div>
          
          <!-- 群组信息 -->
          <div class="ml-3 overflow-hidden">
            <h3 class="font-medium text-gray-800 dark:text-gray-200 truncate group-hover:text-blue-600 dark:group-hover:text-blue-400">
              {{ group.Name }}
            </h3>
            <p class="text-xs text-gray-500 dark:text-gray-400 truncate">
              {{ group.MemberCount }}人
              <span v-if="group.CreatorID === userStore.user?.id" class="ml-1 text-red-500">· 创建者</span>
            </p>
          </div>
        </div>
      </div>
      
      <!-- 空列表提示 -->
      <div v-else-if="!loading" class="py-12 text-center text-gray-500 dark:text-gray-400">
        <svg xmlns="http://www.w3.org/2000/svg" class="h-12 w-12 mx-auto text-gray-300 dark:text-gray-600 mb-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z" />
        </svg>
        <p v-if="searchText">没有找到匹配的群组</p>
        <p v-else>暂无群组</p>
        <button 
          @click="openCreateModal"
          class="mt-4 px-4 py-2 bg-gradient-to-r from-blue-500 to-cyan-500 text-white rounded-lg hover:opacity-90 transition-all flex items-center gap-2 mx-auto"
        >
          <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" d="M12 4.5v15m7.5-7.5h-15" />
          </svg>
          创建新群组
        </button>
      </div>
    </div>
    
    <!-- 创建群组模态框 -->
    <Modal
      v-model:visible="createModalVisible"
      title="创建新群组"
      @ok="handleCreateGroup"
      okText="创建"
      cancelText="取消"
    >
      <div class="space-y-4">
        <div>
          <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">群组名称<span class="text-red-500">*</span></label>
          <Input v-model:value="createGroupForm.name" placeholder="输入群组名称" />
        </div>
        <div>
          <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">群组描述</label>
          <Input.TextArea v-model:value="createGroupForm.description" placeholder="输入群组描述" :rows="3" />
        </div>
        <div>
          <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">群组头像链接</label>
          <Input v-model:value="createGroupForm.avatarUrl" placeholder="输入头像URL" />
        </div>
      </div>
    </Modal>
  </div>
</template>

<style scoped>
.theme-border {
  border-color: var(--border-color);
}
</style> 