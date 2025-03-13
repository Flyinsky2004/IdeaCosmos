<!--
 * @Author: Flyinsky w2084151024@gmail.com
 * @Description: None
-->

<template>
  <div class="flex h-[calc(100vh-5rem)] rounded-xl bg-white dark:bg-black font-song animate__animated animate__fadeInUp">
    <!-- 左侧历史记录栏 -->
    <div 
      class="transition-all rounded-xl duration-300 ease-in-out border-r theme-border bg-gray-50 dark:bg-zinc-900 flex flex-col"
      :class="[isExpanded ? 'w-64' : 'w-16']"
      @mouseenter="isExpanded = true"
      @mouseleave="isExpanded = false">
      <!-- 新建对话按钮 -->
      <div class="p-4">
        <button 
          class="w-full px-4 py-2 bg-blue-500 hover:bg-blue-600 text-white rounded-lg flex items-center justify-center gap-2 transition-colors"
          :class="{'px-2': !isExpanded}"
          @click="createNewChat"
        >
          <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
          </svg>
          <span v-if="isExpanded">新建对话</span>
        </button>
      </div>

      <!-- 历史记录列表 -->
      <div class="flex-1 overflow-y-auto overflow-x-hidden">
        <div v-for="(chat, index) in chatHistory" :key="chat.id" 
             class="px-4 py-3 hover:bg-gray-100 dark:hover:bg-zinc-800 cursor-pointer transition-colors relative group animate__animated animate__fadeInLeft"
             :class="{'bg-gray-100 dark:bg-zinc-800': currentChatId === chat.id}"
             :style="{'animation-delay': index * 0.1 + 's'}"
             @click="loadChatHistory(chat.ID)">
          <div class="text-sm text-gray-900 dark:text-gray-200 truncate">
            <div class="flex items-center">
              <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 mr-2 flex-shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 10h.01M12 10h.01M16 10h.01M9 16H5a2 2 0 01-2-2V6a2 2 0 012-2h14a2 2 0 012 2v8a2 2 0 01-2 2h-4l-4 4z" />
              </svg>
              <span v-if="isExpanded" class="flex-1">{{ chat.title }}</span>
              <!-- 删除按钮 -->
              <button 
                v-if="isExpanded"
                @click.stop="deleteChat(chat.ID)"
                class="opacity-0 group-hover:opacity-100 transition-all duration-300 p-1 hover:text-red-500 transform hover:scale-110">
                <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                </svg>
              </button>
            </div>
          </div>
          <div v-if="isExpanded" class="text-xs text-gray-500 dark:text-gray-400 mt-1">更新于{{ parseDateTime(chat.UpdatedAt) }}</div>
        </div>
      </div>
    </div>

    <!-- 主对话区域 -->
    <div class="flex-1 flex flex-col">
      <!-- 对话内容区 -->
      <div class="flex-1 overflow-y-auto p-4 space-y-6" ref="messageContainer">
        <div v-for="(message, index) in messages" :key="index" 
             class="flex gap-4 animate__animated animate__fadeInUp" 
             :class="{'justify-end': message.role === 'user'}"
             :style="{'animation-delay': index * 0.15 + 's'}">
          <!-- 头像 -->
          <div v-if="message.role !== 'user'" class="w-8 h-8 rounded-full bg-purple-100 flex items-center justify-center flex-shrink-0">
            <img :src="logo" class="w-full h-full rounded-full" />
          </div>
          
          <!-- 消息内容 -->
          <div class="max-w-[80%] rounded-lg" 
               :class="message.role === 'user' ? 'bg-blue-500 text-white p-4' : 'bg-gray-100 dark:bg-zinc-800'">
            <!-- 用户消息直接显示文本 -->
            <template v-if="message.role === 'user'">
              {{ message.content }}
            </template>
            <!-- AI消息使用MdPreview渲染 -->
            <template v-else>
              <MdPreview
                style="background: transparent"
                :theme="themeStore.currentTheme"
                :editorId="'msg-' + index"
                :modelValue="message.content"
                previewTheme="github"
                class="prose dark:prose-invert max-w-none p-4"

              />
            </template>
          </div>
          
          <!-- 用户头像 -->
          <div v-if="message.role === 'user'" class="w-8 h-8 rounded-full bg-blue-100 flex items-center justify-center flex-shrink-0">
            <img :src="BACKEND_DOMAIN + userStore.user.avatar" class="w-full h-full rounded-full" />
          </div>
        </div>
      </div>

      <!-- 输入区域 -->
      <div class="border-t theme-border p-4 bg-white dark:bg-black animate__animated animate__fadeInUp"
           style="animation-delay: 0.4s">
        <div class="flex gap-4 justify-center items-end">
          <a-textarea
            v-model:value="newMessage"
            :auto-size="{ minRows: 1, maxRows: 6 }"
            class="w-1/2 resize-none rounded-lg focus:ring-2 focus:ring-blue-500/20 dark:bg-zinc-900 dark:text-gray-200 hover:border-blue-400 focus:border-blue-500"
            placeholder="输入消息..."
            :bordered="true"
            @pressEnter="handleEnterPress"
          />
          <a-button 
            type="primary"
            :loading="isGenerating"
            @click="sendMessage"
            class="flex items-center gap-2 h-[40px]"
          >
            <template #icon>
              <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
                <path d="M10.894 2.553a1 1 0 00-1.788 0l-7 14a1 1 0 001.169 1.409l5-1.429A1 1 0 009 15.571V11a1 1 0 112 0v4.571a1 1 0 00.725.962l5 1.428a1 1 0 001.17-1.408l-7-14z" />
              </svg>
            </template>
            发送
          </a-button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, nextTick, onUnmounted, watch } from 'vue'
import { get, post, postJSON } from '@/util/request'
import { Modal, message } from 'ant-design-vue'
import { BACKEND_DOMAIN } from '@/util/VARRIBLES'
import { useThemeStore } from '@/stores/theme'
import { MdPreview } from 'md-editor-v3'
import 'md-editor-v3/lib/preview.css'
import { parseDateTime } from '@/util/common'
import { useUserStore } from '@/stores/user'
import logo from '@/assets/img/logo.webp'
import 'animate.css'

// 主题store
const themeStore = useThemeStore()
const userStore = useUserStore()
// 控制侧边栏展开状态
const isExpanded = ref(false)
const currentChatId = ref(null)
const messageContainer = ref(null)
const isGenerating = ref(false)
let ws = null

// 聊天历史数据
const chatHistory = ref([])
const messages = ref([])
const newMessage = ref('')

// 初始化WebSocket连接
const initWebSocket = () => {
  const token = localStorage.getItem("authToken")
  const baseUrl = BACKEND_DOMAIN.replace(/^http/, 'ws').replace(/\/$/, '')
  ws = new WebSocket(
    `${baseUrl}/ws/projectSuggest`
  )

  ws.onopen = () => {
    console.log('WebSocket连接已建立')
    // 连接建立后立即发送初始化数据
    const initData = {
      auth_token: token,
      chat_id: currentChatId.value || null,
      messages: messages.value.map(msg => ({
        role: msg.role === 'user' ? 'user' : 'assistant',
        content: msg.content
      }))
    }
    ws.send(JSON.stringify(initData))
  }

  ws.onmessage = (event) => {
    const response = JSON.parse(event.data)
    if (response.code === 500) {
      message.error(response.message)
      isGenerating.value = false
      ws.close()
      return
    }

    if (response.done) {
      isGenerating.value = false
      // 如果是新对话，需要获取新的对话列表
      if (!currentChatId.value) {
        fetchChatHistory()
      }
      ws.close()
      return
    }

    // 添加AI回复到消息列表
    if (response.content) {
      // 如果是第一条回复，创建新消息
      if (!messages.value.find(m => m.role === 'assistant' && m.isGenerating)) {
        messages.value.push({
          role: 'assistant',
          content: response.content,
          isGenerating: true
        })
      } else {
        // 否则追加到现有消息
        const lastMessage = messages.value.find(m => m.role === 'assistant' && m.isGenerating)
        if (lastMessage) {
          lastMessage.content += response.content
        }
      }

      // 滚动到最新消息
      nextTick(() => {
        if (messageContainer.value) {
          messageContainer.value.scrollTop = messageContainer.value.scrollHeight
        }
      })
    }
  }

  ws.onerror = (error) => {
    console.error("WebSocket error:", error)
    message.error("连接发生错误，请重试")
    isGenerating.value = false
  }

  ws.onclose = () => {
    isGenerating.value = false
    // 移除消息的生成状态
    const generatingMessage = messages.value.find(m => m.type === 'assistant' && m.isGenerating)
    if (generatingMessage) {
      generatingMessage.isGenerating = false
    }
    ws = null
  }
}

// 发送消息
const sendMessage = () => {
  if (newMessage.value.trim() && !isGenerating.value) {
    const messageContent = newMessage.value.trim()
    
    // 先添加到本地显示
    messages.value.push({
      role: 'user',
      content: messageContent
    })
    
    // 清空输入框
    newMessage.value = ''
    
    // 确保WebSocket连接存在
    if (!ws || ws.readyState !== WebSocket.OPEN) {
      initWebSocket()
    }

    // 设置生成状态
    isGenerating.value = true
    
    // 滚动到底部
    nextTick(() => {
      if (messageContainer.value) {
        messageContainer.value.scrollTop = messageContainer.value.scrollHeight
      }
    })
  }
}

// 获取聊天历史列表
const fetchChatHistory = () => {
  get('/api/agent/chats', null, 
    (message, data) => {
      chatHistory.value = data
    }
  )
}

// 加载特定聊天记录
const loadChatHistory = (chatId) => {
  get(`/api/agent/chats/${chatId}`, null,
    async (message, data) => {
      messages.value = data
      currentChatId.value = chatId
      // 滚动到最新消息
      await nextTick()
      if (messageContainer.value) {
        messageContainer.value.scrollTop = messageContainer.value.scrollHeight
      }
    }
  )
}

// 处理回车键按下事件
const handleEnterPress = (e) => {
  // 如果按下Shift+Enter，则插入换行
  if (e.shiftKey) {
    return
  }
  // 否则发送消息
  e.preventDefault()
  sendMessage()
}

// 删除聊天记录
const deleteChat = async (chatId) => {
  try {
    await get(`/api/agent/chats/${chatId}/delete`, null,
      (msg,data) => {
        // 重新获取聊天列表
        fetchChatHistory()
        
        // 如果删除的是当前对话，清空消息
        if (currentChatId.value === chatId) {
          messages.value = []
          currentChatId.value = null
        }
        message.success('删除成功')
      },(msg,data) => {
        console.error('删除失败:', msg)
      },(msg,data) => {
        console.error('删除失败:', msg)
      }
    )
  } catch (error) {
    console.error('删除失败:', error)
  }
}

// 创建新对话
const createNewChat = () => {
  messages.value = []
  currentChatId.value = null
}

// 格式化日期
const formatDate = (dateString) => {
  const date = new Date(dateString)
  return date.toLocaleDateString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit'
  })
}

// 组件挂载时获取聊天历史
onMounted(() => {
  fetchChatHistory()
})

// 组件卸载时关闭WebSocket连接
onUnmounted(() => {
  if (ws) {
    ws.close()
    ws = null
  }
})

// 监听消息变化自动滚动
watch(messages, () => {
  nextTick(() => {
    if (messageContainer.value) {
      messageContainer.value.scrollTop = messageContainer.value.scrollHeight
    }
  })
}, { deep: true })
</script>

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
  background-color: rgba(82, 82, 82, 0.5);
}

/* 添加宋体字体 */
@font-face {
  font-family: 'SongT';
  src: local('SimSun'), local('宋体');
}

.font-song {
  font-family: 'Song', serif;
}

/* Markdown 样式调整 */
:deep(.md-preview) {
  @apply text-gray-900 dark:text-gray-200;
}

:deep(.md-preview-wrapper) {
  background-color: transparent !important;
}

:deep(.md-preview h1),
:deep(.md-preview h2),
:deep(.md-preview h3),
:deep(.md-preview h4),
:deep(.md-preview h5),
:deep(.md-preview h6) {
  @apply border-b-0 pb-0 mt-2 mb-4;
}

:deep(.md-preview p) {
  @apply my-2;
}

:deep(.md-preview ul),
:deep(.md-preview ol) {
  @apply my-2 pl-6;
}

:deep(.md-preview li) {
  @apply my-1;
}

:deep(.md-preview code) {
  @apply bg-gray-100 dark:bg-zinc-700 px-1 rounded;
}

:deep(.md-preview pre) {
  @apply bg-gray-100 dark:bg-zinc-700 p-4 rounded-lg my-4;
}

:deep(.md-preview blockquote) {
  @apply border-l-4 border-gray-300 dark:border-gray-600 pl-4 my-4 italic;
}

:deep(.md-preview table) {
  @apply w-full border-collapse my-4;
}

:deep(.md-preview th),
:deep(.md-preview td) {
  @apply border border-gray-300 dark:border-gray-600 px-4 py-2;
}

:deep(.md-preview th) {
  @apply bg-gray-100 dark:bg-zinc-700;
}

/* 输入框样式优化 */
:deep(.ant-input) {
  @apply transition-all duration-300 ease-in-out;
}

:deep(.ant-input:hover) {
  @apply border-blue-400;
}

:deep(.ant-input:focus) {
  @apply border-blue-500 shadow-md;
}

/* 暗色模式下的输入框样式 */
.dark :deep(.ant-input) {
  @apply bg-zinc-900 text-gray-200 border-zinc-700;
}

.dark :deep(.ant-input:hover) {
  @apply border-blue-400;
}

.dark :deep(.ant-input:focus) {
  @apply border-blue-500;
}

/* 按钮样式优化 */
:deep(.ant-btn-primary) {
  @apply bg-blue-500 border-blue-500 hover:bg-blue-600 hover:border-blue-600;
}

.dark :deep(.ant-btn-primary) {
  @apply bg-blue-600 border-blue-600 hover:bg-blue-700 hover:border-blue-700;
}

/* Modal样式优化 */
:deep(.custom-modal-confirm) {
  .ant-modal-content {
    @apply bg-white dark:bg-zinc-900;
  }
  .ant-modal-header {
    @apply bg-white dark:bg-zinc-900 border-b border-gray-200 dark:border-zinc-800;
  }
  .ant-modal-title {
    @apply text-gray-900 dark:text-gray-200;
  }
  .ant-modal-body {
    @apply text-gray-700 dark:text-gray-300;
  }
  .ant-modal-footer {
    @apply border-t border-gray-200 dark:border-zinc-800;
  }
  .ant-btn-dangerous {
    @apply bg-red-500 border-red-500 text-white hover:bg-red-600 hover:border-red-600;
  }
}
</style>
