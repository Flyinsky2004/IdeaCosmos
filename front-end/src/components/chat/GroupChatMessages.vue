<!--
 * @Author: Flyinsky w2084151024@gmail.com
 * @Description: 群组聊天消息组件
-->
<script setup>
import { ref, onMounted, watch, onBeforeUnmount } from 'vue'
import { useUserStore } from '@/stores/user'
import { get, post, postJSON } from '@/util/request'
import { message } from 'ant-design-vue'
import SpinLoaderLarge from '@/components/spinLoaderLarge.vue'
import { BACKEND_DOMAIN } from '@/util/VARRIBLES'
import { parseDate } from 'element-plus'
import { parseDateTime } from '@/util/common'

const props = defineProps({
  groupId: {
    type: Number,
    required: true
  }
})

const userStore = useUserStore()
const messages = ref([])
const loading = ref(false)
const messageInput = ref('')
const mediaType = ref('text')
const mediaFile = ref(null)
const mediaPreview = ref('')
const isEmojiPickerVisible = ref(false)
const ws = ref(null) // WebSocket连接
const retryCount = ref(0)
const isClosing = ref(false)
const isConnected = ref(false)
const heartbeatInterval = ref(null) // 心跳定时器

// 处理WebSocket消息
const handleWebSocketMessage = (data) => {
  // 处理不同类型的消息
  switch (data.type) {
    case 'chat':
      // 添加新消息到消息列表
      messages.value.push(data)
      // 滚动到底部
      setTimeout(scrollToBottom, 100)
      break
      
    case 'join':
    case 'leave':
      // 处理系统消息（加入/离开）
      messages.value.push({
        ...data,
        MediaType: 'text',
        SenderID: data.senderId,
        Content: data.content,
        CreatedAt: data.timestamp
      })
      setTimeout(scrollToBottom, 100)
      break
      
    default:
      console.warn('未知的消息类型:', data.type)
  }
}

// 发送心跳消息，保持连接活跃
const startHeartbeat = () => {
  // 清除现有的心跳定时器
  stopHeartbeat()
  
  // 每30秒发送一次心跳
  heartbeatInterval.value = setInterval(() => {
    if (ws.value && ws.value.readyState === WebSocket.OPEN) {
      try {
        // 发送一个简单的ping消息
        ws.value.send(JSON.stringify({ type: 'ping' }))
        console.log('发送心跳消息')
      } catch (error) {
        console.error('发送心跳消息失败:', error)
        connectWebSocket() // 尝试重连
      }
    } else if (ws.value && ws.value.readyState !== WebSocket.CONNECTING) {
      connectWebSocket() // 如果连接已断开且不在连接中，尝试重连
    }
  }, 30000)
}

// 停止心跳
const stopHeartbeat = () => {
  if (heartbeatInterval.value) {
    clearInterval(heartbeatInterval.value)
    heartbeatInterval.value = null
  }
}

// 建立WebSocket连接
const connectWebSocket = () => {
  if (!props.groupId) return
  
  // 如果已经连接或正在连接，不要重复连接
  if (ws.value && (ws.value.readyState === WebSocket.OPEN || ws.value.readyState === WebSocket.CONNECTING)) {
    console.log('WebSocket已连接或正在连接中')
    return
  }
  
  try {
    // 将 http:// 替换为 ws://，并移除末尾的斜杠
    const baseUrl = BACKEND_DOMAIN.replace(/^http/, 'ws').replace(/\/$/, '')
    const wsUrl = `${baseUrl}/ws/groupChat/${props.groupId}`
    
    console.log('正在连接WebSocket:', wsUrl)
    
    // 获取认证token
    const token = localStorage.getItem('authToken')
    if (!token) {
      message.error('未登录或登录已过期')
      return
    }

    // 如果已经有连接，先关闭
    if (ws.value) {
      stopHeartbeat()
      isClosing.value = true
      ws.value.close()
      isClosing.value = false
    }

    ws.value = new WebSocket(wsUrl)
    
    // 设置连接超时
    const connectionTimeout = setTimeout(() => {
      if (ws.value && ws.value.readyState === WebSocket.CONNECTING) {
        ws.value.close()
        message.error('WebSocket连接超时，请检查网络状态')
      }
    }, 5000)
    
    ws.value.onopen = () => {
      clearTimeout(connectionTimeout)
      console.log('WebSocket连接已建立，发送认证消息')
      
      // 发送认证消息
      const authMessage = {
        token: token
      }
      ws.value.send(JSON.stringify(authMessage))
    }
    
    ws.value.onclose = (event) => {
      clearTimeout(connectionTimeout)
      console.log('WebSocket连接已关闭:', event)
      isConnected.value = false
      stopHeartbeat()
      
      // 如果不是主动关闭，尝试重连
      if (!isClosing.value) {
        const delay = Math.min(1000 * Math.pow(2, retryCount.value), 30000)
        retryCount.value++
        
        setTimeout(() => {
          if (!isClosing.value) {
            console.log(`尝试第 ${retryCount.value} 次重连...`)
            connectWebSocket()
          }
        }, delay)
      }
    }
    
    ws.value.onerror = (error) => {
      clearTimeout(connectionTimeout)
      console.error('WebSocket错误:', error)
      isConnected.value = false
      message.error('WebSocket连接错误，正在尝试重连...')
    }
    
    ws.value.onmessage = (event) => {
      try {
        const data = JSON.parse(event.data)
        
        // 处理认证成功消息
        if (data.code === 200 && (data.message === '认证成功' || data.data === '认证成功' || data.data === '连接成功')) {
          console.log('WebSocket认证成功')
          isConnected.value = true
          retryCount.value = 0
          startHeartbeat() // 启动心跳
          return
        }
        
        // 处理pong响应，不做特殊处理
        if (data.type === 'pong') {
          console.log('收到pong响应')
          return
        }
        
        // 处理错误消息
        if (data.code !== undefined && data.code !== 200) {
          message.error(data.message || '连接错误')
          if (data.code === 401) {
            // 认证失败，重新连接
            ws.value.close()
          }
          return
        }
        
        // 处理正常的聊天消息
        if (data.type) {
          handleWebSocketMessage(data)
        }
      } catch (error) {
        console.error('解析WebSocket消息失败:', error)
      }
    }
  } catch (error) {
    console.error('创建WebSocket连接失败:', error)
    message.error('创建WebSocket连接失败，请稍后重试')
  }
}

// 关闭WebSocket连接
const closeWebSocket = () => {
  stopHeartbeat()
  if (ws.value) {
    isClosing.value = true
    if (ws.value.readyState === WebSocket.OPEN || ws.value.readyState === WebSocket.CONNECTING) {
      ws.value.close()
    }
    ws.value = null
    isConnected.value = false
  }
}

// 获取历史消息列表
const fetchMessages = () => {
  if (!props.groupId) return
  
  loading.value = true
  get(`/api/chat/groups/${props.groupId}/messages`, {}, (_, data) => {
    messages.value = data
    loading.value = false
    scrollToBottom()
  }, (msg) => {
    loading.value = false
    message.error('获取消息失败' + msg)
  })
}

// 发送消息 - 使用WebSocket
const sendMessage = async () => {
  if (!messageInput.value.trim() && !mediaFile.value) return
  
  // 如果连接正在建立中，等待一段时间
  if (ws.value && ws.value.readyState === WebSocket.CONNECTING) {
    try {
      await new Promise((resolve, reject) => {
        const timeout = setTimeout(() => {
          reject(new Error('连接超时'))
        }, 5000)

        const checkConnection = () => {
          if (isConnected.value) {
            clearTimeout(timeout)
            resolve()
          } else if (ws.value.readyState === WebSocket.CLOSED) {
            clearTimeout(timeout)
            reject(new Error('连接已关闭'))
          } else {
            setTimeout(checkConnection, 100)
          }
        }
        checkConnection()
      })
    } catch (error) {
      message.error('聊天连接失败，正在重新连接...')
      connectWebSocket()
      return
    }
  }

  // 检查连接状态
  if (!isConnected.value || !ws.value || ws.value.readyState !== WebSocket.OPEN) {
    message.error('聊天连接已断开，正在重新连接...')
    connectWebSocket()
    return
  }
  
  // 如果有媒体文件，先上传
  if (mediaFile.value) {
    uploadMedia(() => {
      sendTextMessage()
    })
    return
  }
  
  sendTextMessage()
}

// 发送文本消息
const sendTextMessage = () => {
  if (!isConnected.value) {
    message.error('等待连接建立...')
    return
  }

  const messageData = {
    type: 'chat',
    content: messageInput.value,
    mediaType: mediaType.value,
    mediaURL: mediaPreview.value
  }
  
  try {
    ws.value.send(JSON.stringify(messageData))
    messageInput.value = ''
    mediaFile.value = null
    mediaPreview.value = ''
    mediaType.value = 'text'
  } catch (error) {
    console.error('发送消息失败:', error)
    message.error('发送消息失败，正在重新连接...')
    connectWebSocket()
  }
}

// 上传媒体文件
const uploadMedia = (callback) => {
  if (!mediaFile.value) {
    callback()
    return
  }
  
  const formData = new FormData()
  formData.append('media', mediaFile.value)
  formData.append('type', mediaType.value)
  
  post('/api/upload/chat-media', formData, (_, data) => {
    mediaPreview.value = data.url
    callback()
  }, () => {
    message.error('媒体上传失败')
  })
}

// 处理媒体文件选择
const handleMediaSelect = (event) => {
  const file = event.target.files[0]
  if (!file) return
  
  mediaFile.value = file
  mediaType.value = file.type.startsWith('image/') ? 'image' : 
                    file.type.startsWith('video/') ? 'video' : 
                    file.type.startsWith('audio/') ? 'audio' : 'file'
  
  if (mediaType.value === 'image') {
    const reader = new FileReader()
    reader.onload = (e) => {
      mediaPreview.value = e.target.result
    }
    reader.readAsDataURL(file)
  }
}

// 取消媒体文件选择
const cancelMediaSelect = () => {
  mediaFile.value = null
  mediaPreview.value = ''
  mediaType.value = 'text'
}

// 滚动到底部
const scrollToBottom = () => {
  const container = document.querySelector('.messages-container')
  if (container) {
    container.scrollTop = container.scrollHeight
  }
}

// 格式化时间
const formatTime = (timestamp) => {
  const date = new Date(timestamp)
  const now = new Date()
  
  if (date.toDateString() === now.toDateString()) {
    return date.toLocaleTimeString('zh-CN', { hour: '2-digit', minute: '2-digit' })
  }
  
  return date.toLocaleString('zh-CN', {
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit'
  })
}

// 组件挂载时初始化
onMounted(() => {
  fetchMessages() // 首先获取历史消息
  connectWebSocket() // 然后建立WebSocket连接
})

// 监听groupId变化
watch(() => props.groupId, (newVal, oldVal) => {
  if (newVal !== oldVal) {
    closeWebSocket() // 关闭旧连接
    messages.value = [] // 清空消息列表
    fetchMessages() // 获取新群组的历史消息
    connectWebSocket() // 建立新连接
  }
})

// 组件卸载时清理
onBeforeUnmount(() => {
  stopHeartbeat()
  if (ws.value) {
    isClosing.value = true
    try {
      ws.value.close(1000, 'Component unmounted')
    } catch (error) {
      console.error('关闭WebSocket连接失败:', error)
    }
    ws.value = null
  }
})
</script>

<template>
  <div class="chat-messages-container h-full flex flex-col animate__animated animate__fadeIn">
    <!-- 消息列表 -->
    <div class="messages-container flex-1 overflow-y-auto p-4 space-y-4">
      <div v-if="loading" class="flex justify-center py-8">
        <SpinLoaderLarge />
      </div>
      
      <div v-else-if="messages.length === 0" class="flex flex-col items-center justify-center py-12">
        <div class="text-6xl text-gray-300 dark:text-gray-600 mb-4">
          <svg xmlns="http://www.w3.org/2000/svg" class="h-16 w-16" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z" />
          </svg>
        </div>
        <p class="text-lg text-gray-500 dark:text-gray-400">开始聊天吧</p>
      </div>
      
      <template v-else>
        <div 
          v-for="message in messages" 
          :key="message.ID"
          class="flex items-start gap-3 group animate__animated animate__fadeInUp animate__faster"
          :class="message.SenderID === userStore.user.id ? 'flex-row-reverse' : ''"
        >
          <!-- 发送者头像 -->
          <div class="w-10 h-10 rounded-full overflow-hidden bg-gradient-to-br from-blue-100 to-indigo-100 dark:from-blue-900/30 dark:to-indigo-900/30 flex-shrink-0 border theme-border">
            <img 
              v-if="message.SenderAvatar" 
              :src="BACKEND_DOMAIN + message.senderAvatar" 
              class="w-full h-full object-cover"
              alt="用户头像"
            >
            <div v-else class="w-full h-full flex items-center justify-center text-blue-500 dark:text-blue-400 text-sm font-bold">
              {{ message.senderName?.charAt(0) }}
            </div>
          </div>
          
          <!-- 消息内容 -->
          <div 
            class="max-w-[70%] rounded-lg p-3 group-hover:shadow-md transition-all"
            :class="message.SenderID === userStore.user.id 
              ? 'bg-gradient-to-r from-blue-500 to-cyan-500 text-white' 
              : 'bg-gray-50 dark:bg-zinc-800/50 dark:text-gray-200 border theme-border'"
          >
            <!-- 发送者名称 -->
            <div 
              v-if="message.SenderID !== userStore.user.id"
              class="text-xs mb-1"
              :class="message.SenderID === userStore.user.id 
                ? 'text-blue-200' 
                : 'text-blue-500 dark:text-blue-400'"
            >
              {{ message.SenderName }}
            </div>
            
            <!-- 文本消息 -->
            <p v-if="message.MediaType === 'text'" class="break-words">
              {{ message.Content }}
            </p>
            
            <!-- 图片消息 -->
            <img 
              v-else-if="message.MediaType === 'image'" 
              :src="message.MediaURL"
              class="max-w-full rounded-lg shadow-sm hover:shadow-lg transition-shadow cursor-pointer"
              alt="图片消息"
            >
            
            <!-- 音频消息 -->
            <audio 
              v-else-if="message.MediaType === 'audio'" 
              :src="message.MediaURL"
              controls
              class="max-w-full"
            ></audio>
            
            <!-- 视频消息 -->
            <video 
              v-else-if="message.MediaType === 'video'" 
              :src="message.MediaURL"
              controls
              class="max-w-full rounded-lg shadow-sm"
            ></video>
            
            <!-- 发送时间 -->
            <div 
              class="text-xs mt-1 opacity-0 group-hover:opacity-100 transition-opacity"
              :class="message.SenderID === userStore.user.id 
                ? 'text-blue-200' 
                : 'text-gray-400 dark:text-gray-500'"
            >
              {{ parseDateTime(message.CreatedAt) }}
            </div>
          </div>
        </div>
      </template>
    </div>

    <!-- 输入区域 -->
    <div class="p-4 border-t theme-border bg-white dark:bg-zinc-900">
      <!-- 媒体预览 -->
      <div v-if="mediaPreview" class="mb-4 relative animate__animated animate__fadeIn">
        <img 
          v-if="mediaType === 'image'"
          :src="mediaPreview"
          class="max-h-32 rounded-lg shadow-sm"
          alt="媒体预览"
        >
        <div v-else class="p-3 bg-gray-50 dark:bg-zinc-800/50 rounded-lg border theme-border">
          <span class="text-gray-500 dark:text-gray-400">
            {{ mediaFile?.name }}
          </span>
        </div>
        <button 
          @click="cancelMediaSelect"
          class="absolute -top-2 -right-2 w-6 h-6 bg-red-500 text-white rounded-full flex items-center justify-center hover:bg-red-600 transition-colors shadow-lg"
        >
          <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12" />
          </svg>
        </button>
      </div>
      
      <div class="flex items-center gap-2">
        <!-- 媒体上传按钮 -->
        <div class="relative">
          <input 
            type="file" 
            accept="image/*,audio/*,video/*"
            class="hidden"
            @change="handleMediaSelect"
            ref="fileInput"
          >
          <button 
            @click="$refs.fileInput.click()"
            class="p-2 text-gray-500 hover:text-blue-500 dark:text-gray-400 dark:hover:text-blue-400 transition-colors rounded-lg hover:bg-gray-100 dark:hover:bg-gray-800"
            title="上传媒体文件"
          >
            <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" d="M12 9v6m3-3H9m12 0a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
          </button>
        </div>

        <!-- 表情按钮 -->
        <button 
          class="p-2 text-gray-500 hover:text-blue-500 dark:text-gray-400 dark:hover:text-blue-400 transition-colors rounded-lg hover:bg-gray-100 dark:hover:bg-gray-800"
          title="选择表情"
        >
          <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" d="M14.828 14.828a4 4 0 01-5.656 0M9 10h.01M15 10h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>
        </button>
        
        <!-- 消息输入框 -->
        <input 
          v-model="messageInput"
          @keyup.enter="sendMessage"
          type="text"
          placeholder="输入消息..."
          class="flex-1 px-4 py-2 bg-gray-50 dark:bg-zinc-800/50 border theme-border rounded-full focus:ring-2 focus:ring-blue-500/20 dark:text-gray-300 dark:placeholder-gray-500"
        >
        
        <!-- 发送按钮 -->
        <button 
          @click="sendMessage"
          class="px-4 py-2 bg-gradient-to-r from-blue-500 to-cyan-500 text-white rounded-full hover:opacity-90 transition-all disabled:opacity-50 disabled:cursor-not-allowed flex items-center gap-2"
          :disabled="!messageInput.trim() && !mediaFile"
        >
          <span class="hidden sm:inline">发送</span>
          <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" d="M6 12L3.269 3.126A59.768 59.768 0 0121.485 12 59.77 59.77 0 013.27 20.876L5.999 12zm0 0h7.5" />
          </svg>
        </button>
      </div>
    </div>
  </div>
</template>

<style scoped>
.messages-container {
  scroll-behavior: smooth;
}

.messages-container::-webkit-scrollbar {
  width: 6px;
}

.messages-container::-webkit-scrollbar-track {
  background: transparent;
}

.messages-container::-webkit-scrollbar-thumb {
  background-color: rgba(156, 163, 175, 0.5);
  border-radius: 3px;
}

.messages-container::-webkit-scrollbar-thumb:hover {
  background-color: rgba(156, 163, 175, 0.7);
}

.theme-border {
  border-color: var(--border-color);
}

/* 动画相关样式 */
.animate__animated {
  animation-duration: 0.5s;
}

.animate__faster {
  animation-duration: 0.3s;
}

@keyframes fadeInUp {
  from {
    opacity: 0;
    transform: translateY(10px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.animate__fadeInUp {
  animation-name: fadeInUp;
}
</style> 