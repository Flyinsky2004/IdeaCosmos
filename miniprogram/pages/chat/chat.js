// 获取应用实例
const app = getApp()
import { formatTime } from '../../utils/util'

// 定义基础URL
const BASE_URL = 'https://idea.1024110.xyz'

Page({
  // 标记是否已经从 onLoad 中建立了连接
  hasInitializedConnection: false,
  
  data: {
    // 用户信息
    userInfo: {},
    hasUserInfo: false,
    
    // 聊天相关
    chatHistory: [],
    messages: [],
    newMessage: '',
    currentChatId: null,
    isGenerating: false,
    scrollToMessage: '',
    
    // UI 相关
    isExpanded: true,
    showDeleteModal: false,
    deleteTargetId: null,
    
    // WebSocket 相关
    isConnected: false
  },

  onLoad: function() {
    // 获取用户信息
    if (app.globalData.userInfo) {
      this.setData({
        userInfo: app.globalData.userInfo,
        hasUserInfo: true
      })
    } else {
      // 监听用户信息变化
      app.userInfoReadyCallback = res => {
        this.setData({
          userInfo: res.userInfo,
          hasUserInfo: true
        })
      }
    }
    
    // 获取聊天历史
    this.fetchChatHistory()
  },
  
  onShow: function() {
    // 页面显示时，不自动重连WebSocket
    // 只有在发送消息时才会检查连接状态并重连
  },
  
  onHide: function() {
    // 页面隐藏时，不关闭 WebSocket，保持连接
  },
  
  onUnload: function() {
    // 页面卸载时，关闭 WebSocket
    this.closeWebSocket()
  },
  
  // 封装请求函数
  request(url, method, data, needAuth = true) {
    return new Promise((resolve, reject) => {      
      // 构建请求头
      const header = {
        'content-type': 'application/json'
      }
      
      // 如果需要授权，添加token
      if (needAuth) {
        const token = wx.getStorageSync('authToken')
        if (token) {
          header['Authorization'] = token
        } else {
          // 如果需要授权但没有token，直接返回未登录错误
          reject(new Error('请先登录'))
          return
        }
      }
      
      wx.request({
        url: `${BASE_URL}${url}`,
        method: method,
        data: data,
        header: header,
        success: (res) => {
          if (res.statusCode === 200) {
            resolve(res.data)
          } else if (res.statusCode === 401) {
            // 处理授权失败
            wx.showToast({
              title: '登录已过期，请重新登录',
              icon: 'none'
            })
            reject(new Error('登录已过期'))
          } else {
            reject(res)
          }
        },
        fail: (error) => {
          console.error(`请求失败 ${url}:`, error)
          reject(error)
        }
      })
    })
  },
  
  // 初始化 WebSocket 连接
  initWebSocket() {
    // 如果已经有连接，先关闭
    if (this.socketTask) {
      this.closeWebSocket()
    }
    
    try {
      // 获取认证token
      const token = wx.getStorageSync('authToken')
      if (!token) {
        wx.showToast({
          title: '请先登录',
          icon: 'none'
        })
        return
      }
      
      // 将 http:// 替换为 ws://
      const wsUrl = `${BASE_URL.replace(/^http/, 'ws')}/ws/projectSuggest`
      
      console.log('正在连接WebSocket:', wsUrl)
      
      // 创建WebSocket连接
      this.socketTask = wx.connectSocket({
        url: wsUrl,
        success: () => {
          console.log('WebSocket连接创建成功')
        },
        fail: (error) => {
          console.error('WebSocket连接创建失败:', error)
          this.setData({ isConnected: false })
          
          wx.showToast({
            title: '连接服务器失败，请稍后再试',
            icon: 'none'
          })
        }
      })
      
      // 监听WebSocket连接打开
      this.socketTask.onOpen(() => {
        console.log('WebSocket连接已打开，发送认证消息')
        this.setData({ isConnected: true })
        
        // 发送认证消息
        const authMessage = {
          auth_token: token,
          chat_id: this.data.currentChatId || null,
          messages: this.data.messages.map(msg => ({
            role: msg.role === 'user' ? 'user' : 'assistant',
            content: msg.content
          }))
        }
        
        // 发送初始化数据
        this.socketTask.send({
          data: JSON.stringify(authMessage),
          success: () => {
            console.log('认证消息发送成功')
          },
          fail: (error) => {
            console.error('认证消息发送失败:', error)
          }
        })
      })
      
      // 监听WebSocket接收到服务器的消息
      this.socketTask.onMessage((res) => {
        try {
          const response = JSON.parse(res.data)
          
          if (response.code === 500) {
            wx.showToast({
              title: response.message || '请求失败',
              icon: 'none'
            })
            this.setData({
              isGenerating: false
            })
            this.closeWebSocket()
            return
          }
          
          if (response.done) {
            this.setData({
              isGenerating: false
            })
            
            // 如果是新对话，需要获取新的对话列表
            if (!this.data.currentChatId) {
              this.fetchChatHistory()
            }
            
            this.closeWebSocket()
            return
          }
          
          // 添加AI回复到消息列表
          if (response.content) {
            const messages = this.data.messages
            const lastMessage = messages.find(m => m.role === 'assistant' && m.isGenerating)
            
            if (!lastMessage) {
              // 创建新消息
              messages.push({
                role: 'assistant',
                content: response.content,
                formattedContent: this.formatMessage(response.content),
                isGenerating: true
              })
            } else {
              // 追加到现有消息
              const index = messages.findIndex(m => m.role === 'assistant' && m.isGenerating)
              if (index !== -1) {
                messages[index].content += response.content
                messages[index].formattedContent = this.formatMessage(messages[index].content)
              }
            }
            
            this.setData({
              messages: messages,
              scrollToMessage: `msg-${messages.length - 1}`
            })
          }
        } catch (error) {
          console.error('解析WebSocket消息失败:', error)
        }
      })
      
      // 监听WebSocket错误
      this.socketTask.onError((error) => {
        console.error('WebSocket错误:', error)
        this.setData({ 
          isConnected: false,
          isGenerating: false
        })
        
        wx.showToast({
          title: 'WebSocket连接错误',
          icon: 'none'
        })
      })
      
      // 监听WebSocket连接关闭
      this.socketTask.onClose((res) => {
        console.log('WebSocket连接已关闭:', res)
        this.setData({ 
          isConnected: false 
        })
        
        // 移除消息的生成状态
        const messages = this.data.messages
        const index = messages.findIndex(m => m.role === 'assistant' && m.isGenerating)
        if (index !== -1) {
          messages[index].isGenerating = false
          this.setData({
            messages: messages
          })
        }
        
        this.socketTask = null
      })
    } catch (error) {
      console.error('创建WebSocket连接失败:', error)
      this.setData({ 
        isConnected: false,
        isGenerating: false
      })
      
      wx.showToast({
        title: '创建WebSocket连接失败，请稍后重试',
        icon: 'none'
      })
    }
  },
  
  // 关闭WebSocket连接
  closeWebSocket() {
    if (this.socketTask) {
      try {
        this.socketTask.close({
          success: () => {
            console.log('WebSocket连接已关闭')
          },
          fail: (error) => {
            console.error('关闭WebSocket连接失败:', error)
          },
          complete: () => {
            this.socketTask = null
            this.setData({ isConnected: false })
          }
        })
      } catch (error) {
        console.error('关闭WebSocket连接出错:', error)
        this.socketTask = null
        this.setData({ isConnected: false })
      }
    }
  },
  
  // 获取聊天历史列表
  fetchChatHistory() {
    get('/api/agent/chats', null, 
      (messager, data) => {
        try {
          // 确保 data 是数组
          const chatList = Array.isArray(data) ? data : []
          if (!Array.isArray(data)) {
            console.warn('聊天历史数据格式不正确:', data)
          }
          
          // 格式化时间
          const chatHistory = chatList.map(chat => {
            if (!chat || typeof chat !== 'object') return null
            return {
              ...chat,
              formattedTime: chat.UpdatedAt ? formatTime(new Date(chat.UpdatedAt)) : '未知时间'
            }
          }).filter(Boolean) // 过滤掉无效的记录
          
          this.setData({
            chatHistory: chatHistory
          })
        } catch (error) {
          console.error('处理聊天历史数据失败:', error)
          wx.showToast({
            title: '处理聊天历史数据失败',
            icon: 'none'
          })
        }
      },
      (messager) => {
        wx.showToast({
          title: messager || '获取聊天历史失败',
          icon: 'none'
        })
      },
      (messager) => {
        wx.showToast({
          title: messager || '获取聊天历史失败',
          icon: 'none'
        })
      }
    )
  },
  
  // 加载特定聊天记录
  loadChatHistory(e) {
    const chatId = e.currentTarget.dataset.id
    
    get(`/api/agent/chats/${chatId}`, null,
      (messager, data) => {
        try {
          // 确保 data 是数组
          const messageList = Array.isArray(data) ? data : []
          if (!Array.isArray(data)) {
            console.warn('聊天记录数据格式不正确:', data)
          }
          
          // 格式化消息内容
          const messages = messageList.map(msg => {
            if (!msg || typeof msg !== 'object') return null
            return {
              ...msg,
              formattedContent: msg.content ? this.formatMessage(msg.content) : ''
            }
          }).filter(Boolean) // 过滤掉无效的消息
          
          this.setData({
            messages: messages,
            currentChatId: chatId,
            scrollToMessage: messages.length > 0 ? `msg-${messages.length - 1}` : ''
          })
        } catch (error) {
          console.error('处理聊天记录数据失败:', error)
          wx.showToast({
            title: '处理聊天记录数据失败',
            icon: 'none'
          })
        }
      },
      (messager) => {
        wx.showToast({
          title: messager || '获取聊天记录失败',
          icon: 'none'
        })
      },
      (messager) => {
        wx.showToast({
          title: messager || '获取聊天记录失败',
          icon: 'none'
        })
      }
    )
  },
  
  // 创建新对话
  createNewChat() {
    this.setData({
      messages: [],
      currentChatId: null
    })
  },
  
  // 处理输入变化
  onInputChange(e) {
    this.setData({
      newMessage: e.detail.value
    })
  },
  
  // 发送消息
  sendMessage() {
    if (!this.data.newMessage.trim() || this.data.isGenerating) {
      return
    }
    
    const messageContent = this.data.newMessage.trim()
    
    // 添加用户消息到列表
    const messages = this.data.messages
    messages.push({
      role: 'user',
      content: messageContent
    })
    
    this.setData({
      messages: messages,
      newMessage: '',
      isGenerating: true,
      scrollToMessage: `msg-${messages.length - 1}`
    })
    
    // 确保WebSocket连接存在
    if (!this.socketTask || !this.data.isConnected) {
      this.initWebSocket()
    }
  },
  
  // 切换侧边栏展开状态
toggleSidebar() {
  this.setData({
    isExpanded: !this.data.isExpanded
  })
  // 调整主体区域的宽度以适应侧边栏的变化
  const mainArea = wx.createSelectorQuery().select('.main-area')
  mainArea.boundingClientRect((rect) => {
    if (rect) {
      this.setData({
        mainAreaWidth: rect.width
      })
    }
  }).exec()
},
  
  // 显示删除确认弹窗
  showDeleteConfirm(e) {
    const chatId = e.currentTarget.dataset.id
    this.setData({
      showDeleteModal: true,
      deleteTargetId: chatId
    })
  },
  
  // 取消删除
  cancelDelete() {
    this.setData({
      showDeleteModal: false,
      deleteTargetId: null
    })
  },
  
  // 确认删除
  confirmDelete() {
    const chatId = this.data.deleteTargetId
    
    get(`/api/agent/chats/${chatId}/delete`, null,
      (messager, data) => {
        // 重新获取聊天列表
        this.fetchChatHistory()
        
        // 如果删除的是当前对话，清空消息
        if (this.data.currentChatId === chatId) {
          this.setData({
            messages: [],
            currentChatId: null
          })
        }
        
        wx.showToast({
          title: '删除成功',
          icon: 'success'
        })
        
        this.setData({
          showDeleteModal: false,
          deleteTargetId: null
        })
      },
      (messager) => {
        wx.showToast({
          title: messager || '删除失败',
          icon: 'none'
        })
        this.setData({
          showDeleteModal: false,
          deleteTargetId: null
        })
      },
      (messager) => {
        wx.showToast({
          title: messager || '删除失败',
          icon: 'none'
        })
        this.setData({
          showDeleteModal: false,
          deleteTargetId: null
        })
      }
    )
  },
  
  // 格式化消息内容（处理 Markdown、代码块等）
  formatMessage(content) {
    if (!content) return ''
    
    // 简单的 Markdown 转 HTML
    let html = content
      // 代码块
      .replace(/```([\s\S]*?)```/g, '<view class="code-block">$1</view>')
      // 粗体
      .replace(/\*\*(.*?)\*\*/g, '<strong>$1</strong>')
      // 斜体
      .replace(/\*(.*?)\*/g, '<em>$1</em>')
      // 链接
      .replace(/\[(.*?)\]\((.*?)\)/g, '<a href="$2">$1</a>')
      // 换行
      .replace(/\n/g, '<br>')
    
    return html
  }
})

// 封装 get 请求
function get(url, data, successCallback, warningCallback, errorCallback) {
  const token = wx.getStorageSync('authToken')
  
  wx.request({
    url: `${BASE_URL}${url}`,
    method: 'GET',
    data: data,
    header: {
      'Authorization': token
    },
    success(res) {
      if (res.statusCode === 200) {
        if (typeof successCallback === 'function') {
          successCallback('', res.data)
        }
      } else {
        if (typeof warningCallback === 'function') {
          warningCallback(res.data.message || '请求失败')
        }
      }
    },
    fail(err) {
      if (typeof errorCallback === 'function') {
        errorCallback(err.errMsg || '网络错误')
      }
    }
  })
}
