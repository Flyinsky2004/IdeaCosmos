<script setup>
import { ref, onMounted, nextTick } from 'vue'
import { useThemeStore } from '@/stores/theme'
import { Transition } from 'vue'

const themeStore = useThemeStore()
const messageInput = ref('')
const sidebarOpen = ref(true)
const isTyping = ref(false)
const messagesContainer = ref(null)

const conversations = ref([
  {
    id: 1,
    title: '如何使用 Vue 3 组合式 API',
    lastMessage: '组合式 API 是 Vue 3 中...',
    timestamp: '3分钟前',
    active: true
  },
  {
    id: 2,
    title: '讨论项目架构设计',
    lastMessage: '关于项目架构，我们需要...',
    timestamp: '2小时前',
    active: false
  }
])

const currentChat = ref({
  messages: [
    {
      id: 1,
      role: 'assistant',
      content: '你好！我是 AI 助手，很高兴为你服务。我可以帮助你解决编程问题、讨论技术话题，或者回答你的其他问题。',
      timestamp: '14:00'
    },
    {
      id: 2,
      role: 'user',
      content: '你能帮我解释一下 Vue 3 的响应式原理吗？',
      timestamp: '14:01'
    },
    {
      id: 3,
      role: 'assistant',
      content: 'Vue 3 的响应式系统是基于 Proxy 实现的，这是一个重大改进。相比 Vue 2 使用的 Object.defineProperty，Proxy 可以监听数组变化、对象属性的添加和删除等操作，使得响应式系统更加强大和灵活。',
      timestamp: '14:02'
    }
  ]
})

const toggleSidebar = () => {
  sidebarOpen.value = !sidebarOpen.value
}

const scrollToBottom = async () => {
  await nextTick()
  if (messagesContainer.value) {
    messagesContainer.value.scrollTop = messagesContainer.value.scrollHeight
  }
}

const simulateTyping = async () => {
  isTyping.value = true
  await new Promise(resolve => setTimeout(resolve, 2000))
  isTyping.value = false
}

const sendMessage = async () => {
  if (!messageInput.value.trim()) return
  
  const userMessage = {
    id: currentChat.value.messages.length + 1,
    role: 'user',
    content: messageInput.value,
    timestamp: new Date().toLocaleTimeString('zh-CN', { hour: '2-digit', minute: '2-digit' })
  }
  
  currentChat.value.messages.push(userMessage)
  messageInput.value = ''
  await scrollToBottom()
  
  // 模拟 AI 响应
  await simulateTyping()
  const aiResponse = {
    id: currentChat.value.messages.length + 1,
    role: 'assistant',
    content: '这是一个模拟的 AI 响应消息...',
    timestamp: new Date().toLocaleTimeString('zh-CN', { hour: '2-digit', minute: '2-digit' })
  }
  currentChat.value.messages.push(aiResponse)
  await scrollToBottom()
}

const selectConversation = (id) => {
  conversations.value = conversations.value.map(conv => ({
    ...conv,
    active: conv.id === id
  }))
}

onMounted(() => {
  scrollToBottom()
})
</script>

<template>
  <div class="flex h-[calc(100vh-5rem)] overflow-hidden font-song">
    <!-- 侧边栏切换按钮 -->
    <button 
      @click="toggleSidebar"
      class="fixed left-4 top-20 z-50 p-2 rounded-full bg-white/80 dark:bg-zinc-800/80 shadow-lg backdrop-blur-sm hover:bg-gray-100 dark:hover:bg-zinc-700 transition-colors"
    >
      <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 dark:text-gray-300" :class="{ 'rotate-180': !sidebarOpen }" fill="none" viewBox="0 0 24 24" stroke="currentColor">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
      </svg>
    </button>

    <!-- 左侧对话列表 -->
    <Transition name="slide">
      <div v-show="sidebarOpen" class="w-80 border-r theme-border bg-white/80 dark:bg-zinc-900/80 backdrop-blur-sm">
        <div class="p-4 border-b theme-border">
          <button class="w-full py-2 px-4 rounded-lg bg-gradient-to-r from-blue-500 to-blue-600 hover:from-blue-600 hover:to-blue-700 text-white transition-all shadow-md hover:shadow-lg">
            新建对话
          </button>
        </div>
        <div class="overflow-y-auto h-[calc(100%-4rem)]">
          <div
            v-for="conv in conversations"
            :key="conv.id"
            @click="selectConversation(conv.id)"
            class="p-4 cursor-pointer transition-all hover:bg-gray-100 dark:hover:bg-zinc-800 border-l-4 border-transparent"
            :class="{ 'border-l-blue-500 bg-gray-100 dark:bg-zinc-800': conv.active }"
          >
            <div class="flex items-center justify-between">
              <h3 class="font-medium text-gray-900 dark:text-gray-100 truncate">{{ conv.title }}</h3>
              <span class="text-xs text-gray-500">{{ conv.timestamp }}</span>
            </div>
            <p class="text-sm text-gray-500 dark:text-gray-400 mt-1 truncate">{{ conv.lastMessage }}</p>
          </div>
        </div>
      </div>
    </Transition>

    <!-- 主对话区域 -->
    <div class="flex-1 flex flex-col bg-gray-50 dark:bg-black/30">
      <!-- 消息列表 -->
      <div ref="messagesContainer" class="flex-1 overflow-y-auto p-6 space-y-6">
        <div
          v-for="message in currentChat.messages"
          :key="message.id"
          class="flex gap-4"
          :class="message.role === 'user' ? 'justify-end' : 'justify-start'"
        >
          <!-- AI 头像 -->
          <div v-if="message.role === 'assistant'" class="w-8 h-8 rounded-full bg-gradient-to-r from-blue-500 to-blue-600 flex items-center justify-center text-white text-sm">
            AI
          </div>
          
          <div
            class="max-w-[70%] rounded-2xl px-6 py-4 shadow-sm"
            :class="message.role === 'user' ? 'bg-gradient-to-r from-blue-500 to-blue-600 text-white' : 'bg-white dark:bg-zinc-800 dark:text-gray-100'"
          >
            <p class="text-[15px] leading-relaxed whitespace-pre-wrap">{{ message.content }}</p>
            <span class="text-xs opacity-70 mt-2 block">{{ message.timestamp }}</span>
          </div>

          <!-- 用户头像 -->
          <div v-if="message.role === 'user'" class="w-8 h-8 rounded-full bg-gray-200 dark:bg-zinc-700 flex items-center justify-center">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 text-gray-500 dark:text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" />
            </svg>
          </div>
        </div>

        <!-- 打字动画 -->
        <div v-if="isTyping" class="flex gap-4">
          <div class="w-8 h-8 rounded-full bg-gradient-to-r from-blue-500 to-blue-600 flex items-center justify-center text-white text-sm">
            AI
          </div>
          <div class="bg-white dark:bg-zinc-800 rounded-2xl px-6 py-4 shadow-sm">
            <div class="flex gap-2">
              <div class="w-2 h-2 rounded-full bg-blue-500 animate-bounce"></div>
              <div class="w-2 h-2 rounded-full bg-blue-500 animate-bounce" style="animation-delay: 0.2s"></div>
              <div class="w-2 h-2 rounded-full bg-blue-500 animate-bounce" style="animation-delay: 0.4s"></div>
            </div>
          </div>
        </div>
      </div>

      <!-- 输入区域 -->
      <div class="border-t theme-border p-4 bg-white/80 dark:bg-zinc-900/80 backdrop-blur-sm">
        <div class="max-w-4xl mx-auto flex gap-4">
          <textarea
            v-model="messageInput"
            @keydown.enter.prevent="sendMessage"
            rows="1"
            class="flex-1 resize-none rounded-xl border-0 bg-gray-100 dark:bg-zinc-800 p-4 focus:ring-2 focus:ring-blue-500/20 dark:text-gray-300 shadow-inner"
            placeholder="输入消息..."
          ></textarea>
          <button
            @click="sendMessage"
            class="px-6 py-2 bg-gradient-to-r from-blue-500 to-blue-600 hover:from-blue-600 hover:to-blue-700 text-white rounded-xl transition-all shadow-md hover:shadow-lg flex items-center gap-2"
          >
            <span>发送</span>
            <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 7l5 5m0 0l-5 5m5-5H6" />
            </svg>
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.theme-border {
  @apply border-gray-200 dark:border-zinc-800;
}

/* 自定义滚动条样式 */
::-webkit-scrollbar {
  width: 6px;
}

::-webkit-scrollbar-track {
  background: transparent;
}

::-webkit-scrollbar-thumb {
  background-color: rgba(156, 163, 175, 0.5);
  border-radius: 3px;
}

.dark ::-webkit-scrollbar-thumb {
  background-color: rgba(64, 64, 64, 0.5);
}

/* 侧边栏动画 */
.slide-enter-active,
.slide-leave-active {
  transition: transform 0.3s ease-in-out;
}

.slide-enter-from,
.slide-leave-to {
  transform: translateX(-100%);
}

/* 宋体字体 */
@font-face {
  font-family: 'Song';
  src: local('Songti SC'), local('SimSun');
}

.font-song {
  font-family: 'Song', serif;
}

/* 消息气泡动画 */
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

.flex-1 > div {
  animation: fadeIn 0.3s ease-out forwards;
}
</style>
