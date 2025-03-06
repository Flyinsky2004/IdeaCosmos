<!--
 * @Author: Flyinsky w2084151024@gmail.com
 * @Description: 群组聊天页面
-->
<script setup>
import { ref } from 'vue'
import { get, post, postJSON } from '@/util/request'
import GroupChatList from '@/components/chat/GroupChatList.vue'
import GroupChatMessages from '@/components/chat/GroupChatMessages.vue'

const activeGroup = ref(null)

const handleSelectGroup = (group) => {
  // 获取群组详情
  get(`/api/chat/groups/${group.ID}`, {}, (_, data) => {
    activeGroup.value = data.group
  })
}

const handleCreateGroup = (group) => {
  // 创建群组
  postJSON('/api/chat/groups', group, (_, data) => {
    activeGroup.value = data
  })
}
</script>

<template>
  <div class="flex flex-col h-[calc(100vh-5rem)] animate__animated animate__fadeIn">

    <!-- 聊天主区域 -->
    <div class="flex-1 mx-4 mb-4 bg-white dark:bg-zinc-900 rounded-2xl border theme-border shadow-sm overflow-hidden animate__animated animate__fadeIn animate__delay-2s">
      <div class="h-full flex">
        <!-- 左侧群组列表 -->
        <div class="w-80 flex-shrink-0 border-r theme-border">
          <GroupChatList 
            :active-group-id="activeGroup?.ID"
            @select-group="handleSelectGroup"
            @create-group="handleCreateGroup"
          />
        </div>
        
        <!-- 右侧聊天区域 -->
        <div class="flex-1">
          <div v-if="activeGroup" class="h-full">
            <GroupChatMessages :group-id="activeGroup.ID" />
          </div>
          <div v-else class="h-full flex flex-col items-center justify-center text-gray-500 dark:text-gray-400 bg-gray-50 dark:bg-zinc-800/50">
            <div class="text-6xl mb-4 text-gray-300 dark:text-gray-600">
              <svg xmlns="http://www.w3.org/2000/svg" class="h-24 w-24" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" d="M8.625 9.75a.375.375 0 11-.75 0 .375.375 0 01.75 0zm0 0H8.25m4.125 0a.375.375 0 11-.75 0 .375.375 0 01.75 0zm0 0H12m4.125 0a.375.375 0 11-.75 0 .375.375 0 01.75 0zm0 0h-.375m-13.5 3.01c0 1.6 1.123 2.994 2.707 3.227 1.087.16 2.185.283 3.293.369V21l4.184-4.183a1.14 1.14 0 01.778-.332 48.294 48.294 0 005.83-.498c1.585-.233 2.708-1.626 2.708-3.228V6.741c0-1.602-1.123-2.995-2.707-3.228A48.394 48.394 0 0012 3c-2.392 0-4.744.175-7.043.513C3.373 3.746 2.25 5.14 2.25 6.741v6.018z" />
              </svg>
            </div>
            <h2 class="text-2xl font-bold mb-2 text-gray-700 dark:text-gray-300">欢迎使用群组聊天</h2>
            <p class="text-lg mb-6">选择一个群组开始聊天</p>
            <div class="flex gap-4">
              <button 
                @click="$emit('create-group')"
                class="px-6 py-3 bg-gradient-to-r from-blue-500 to-cyan-500 text-white rounded-lg hover:opacity-90 transition-all flex items-center gap-2"
              >
                <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" d="M12 4.5v15m7.5-7.5h-15" />
                </svg>
                创建新群组
              </button>
              <button 
                class="px-6 py-3 bg-gray-100 text-gray-700 rounded-lg hover:bg-gray-200 dark:bg-gray-700 dark:text-gray-300 dark:hover:bg-gray-600 transition-all flex items-center gap-2"
              >
                <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" d="M21 21l-5.197-5.197m0 0A7.5 7.5 0 105.196 5.196a7.5 7.5 0 0010.607 10.607z" />
                </svg>
                查找群组
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.theme-border {
  border-color: var(--border-color);
}

/* 动画相关样式 */
.animate__animated {
  animation-duration: 0.5s;
}

.animate__delay-1s {
  animation-delay: 0.1s;
}

.animate__delay-2s {
  animation-delay: 0.2s;
}

@keyframes fadeIn {
  from {
    opacity: 0;
    transform: translateY(10px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.animate__fadeIn {
  animation-name: fadeIn;
}
</style> 