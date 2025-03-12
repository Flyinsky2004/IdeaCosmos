/*
 * @Author: Flyinsky w2084151024@gmail.com
 * @Description: 群组聊天页面
 */
const app = getApp()
const BASE_URL = 'https://idea.1024110.xyz' // 开发环境

Page({
  // 标记是否已经从 onLoad 中建立了连接
  hasInitializedConnection: false,

  /**
   * 页面的初始数据
   */
  data: {
    // 群组信息
    groupId: null,
    groupInfo: {},
    
    // 用户信息
    userInfo: null,
    
    // 消息相关
    messages: [],
    messageInput: '',
    scrollToMessage: '',
    loading: false,
    page: 1,
    hasMore: true,
    
    // WebSocket相关
    isConnected: false,
    
    // 媒体相关
    mediaType: 'text',
    mediaFile: null,
    mediaPreview: ''
  },

  /**
   * 生命周期函数--监听页面加载
   */
  onLoad(options) {
    // 获取群组ID和名称
    const { id, name } = options
    
    if (!id) {
      wx.showToast({
        title: '群组ID不能为空',
        icon: 'none'
      })
      setTimeout(() => {
        wx.navigateBack()
      }, 1500)
      return
    }
    
    // 设置导航栏标题
    wx.setNavigationBarTitle({
      title: decodeURIComponent(name) || '群组聊天'
    })
    
    // 获取用户信息
    this.setData({
      groupId: id,
      userInfo: app.globalData.userInfo || {}
    })
    
    // 获取群组详情
    this.fetchGroupInfo()
    
    // 获取历史消息
    this.fetchMessages()
    
    // 建立WebSocket连接
    this.hasInitializedConnection = true
    this.connectWebSocket()
  },

  /**
   * 生命周期函数--监听页面初次渲染完成
   */
  onReady() {
    
  },

  /**
   * 生命周期函数--监听页面显示
   */
  onShow() {
    // 如果WebSocket未连接，尝试重新连接
    // 但只有在已经初始化过连接的情况下才重连
    if (this.hasInitializedConnection && !this.data.isConnected) {
      console.log('页面显示，检测到连接断开，尝试重新连接')
      this.connectWebSocket()
    }
  },

  /**
   * 生命周期函数--监听页面隐藏
   */
  onHide() {
    
  },

  /**
   * 生命周期函数--监听页面卸载
   */
  onUnload() {
    // 关闭WebSocket连接
    this.closeWebSocket()
  },

  /**
   * 页面相关事件处理函数--监听用户下拉动作
   */
  onPullDownRefresh() {
    // 下拉刷新，重新获取消息
    this.setData({
      page: 1,
      hasMore: true,
      messages: []
    })
    this.fetchMessages(() => {
      wx.stopPullDownRefresh()
    })
  },

  /**
   * 页面上拉触底事件的处理函数
   */
  onReachBottom() {
    
  },

  /**
   * 用户点击右上角分享
   */
  onShareAppMessage() {
    return {
      title: `加入群聊：${this.data.groupInfo.Name}`,
      path: `/pages/group-chat/index?id=${this.data.groupId}&name=${encodeURIComponent(this.data.groupInfo.Name)}`
    }
  },

  // 封装请求函数
  request(url, method, data, needAuth = false) {
    return new Promise((resolve, reject) => {      
      // 构建请求头
      const header = {
        'content-type': 'application/json',
        'Cache-Control': 'no-cache',
        'Pragma': 'no-cache',
        'If-Modified-Since': '0'
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
            // 处理Web端的响应格式，Web端可能直接返回数据或包装在data字段中
            if (res.data.code !== undefined) {
              // 如果响应包含code字段，按照标准格式处理
              resolve(res.data)
            } else {
              // 如果响应直接是数据，包装成标准格式
              resolve({
                code: 200,
                data: res.data,
                message: 'success'
              })
            }
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
          console.error(`请求失败 ${url}:`, error) // 添加错误日志
          reject(error)
        }
      })
    })
  },

  // 获取群组详情
  async fetchGroupInfo() {
    try {
      const res = await this.request(`/api/chat/groups/${this.data.groupId}`, 'GET', {}, true)
      
      if (res.code === 200 && res.data) {
        const groupInfo = res.data.group || res.data
        
        // 处理头像URL
        if (groupInfo.AvatarURL) {
          groupInfo.AvatarURL = this.getImageUrl(groupInfo.AvatarURL)
        }
        
        this.setData({
          groupInfo
        })
      } else {
        throw new Error(res.message || '获取群组详情失败')
      }
    } catch (error) {
      console.error('获取群组详情失败:', error)
      wx.showToast({
        title: error.message || '获取群组详情失败',
        icon: 'none'
      })
    }
  },

  // 获取历史消息
  async fetchMessages(callback) {
    if (!this.data.hasMore || this.data.loading) {
      if (typeof callback === 'function') callback()
      return
    }
    
    try {
      this.setData({ loading: true })
      
      const res = await this.request(`/api/chat/groups/${this.data.groupId}/messages`, 'GET', {
        page: this.data.page,
        limit: 20
      }, true)
      
      if (res.code === 200 && res.data) {
        let messages = res.data
        
        // 确保每条消息都有type属性，用于区分系统消息和普通消息
        messages = messages.map(msg => {
          // 处理发送者头像
          if (msg.senderAvatar) {
            msg.senderAvatar = this.getImageUrl(msg.senderAvatar)
          }
          
          // 假设消息的内容中包含"加入群聊"或"离开群聊"的是系统消息
          if (msg.Content && (msg.Content.includes('加入群聊') || msg.Content.includes('离开群聊'))) {
            return { ...msg, type: msg.Content.includes('加入群聊') ? 'join' : 'leave' }
          }
          return { ...msg, type: 'chat' }
        })
        
        // 更新消息列表
        this.setData({
          messages: [...this.data.messages, ...messages],
          loading: false,
          page: this.data.page + 1,
          hasMore: messages.length === 20 // 如果返回的消息数量小于请求的数量，说明没有更多消息了
        })
        
        // 如果是第一页，滚动到底部
        if (this.data.page === 2) {
          this.scrollToBottom()
        }
      } else {
        throw new Error(res.message || '获取消息失败')
      }
    } catch (error) {
      console.error('获取消息失败:', error)
      wx.showToast({
        title: error.message || '获取消息失败',
        icon: 'none'
      })
      this.setData({ loading: false })
    } finally {
      if (typeof callback === 'function') callback()
    }
  },

  // 加载更多消息
  loadMoreMessages() {
    if (this.data.hasMore && !this.data.loading) {
      this.fetchMessages()
    }
  },

  // 处理消息输入
  onMessageInput(e) {
    this.setData({
      messageInput: e.detail.value
    })
  },

  // 选择媒体文件
  chooseMedia() {
    wx.showActionSheet({
      itemList: ['图片', '视频', '文件'],
      success: (res) => {
        switch (res.tapIndex) {
          case 0: // 图片
            this.chooseImage()
            break
          case 1: // 视频
            this.chooseVideo()
            break
          case 2: // 文件
            wx.showToast({
              title: '小程序暂不支持发送文件',
              icon: 'none'
            })
            break
        }
      }
    })
  },

  // 选择图片
  chooseImage() {
    wx.chooseImage({
      count: 1,
      sizeType: ['compressed'],
      sourceType: ['album', 'camera'],
      success: (res) => {
        const tempFilePath = res.tempFilePaths[0]
        const file = res.tempFiles[0]
        
        this.setData({
          mediaType: 'image',
          mediaFile: file,
          mediaPreview: tempFilePath
        })
      }
    })
  },

  // 选择视频
  chooseVideo() {
    wx.chooseVideo({
      sourceType: ['album', 'camera'],
      maxDuration: 60,
      camera: 'back',
      success: (res) => {
        this.setData({
          mediaType: 'video',
          mediaFile: {
            path: res.tempFilePath,
            size: res.size,
            name: '视频消息'
          },
          mediaPreview: res.thumbTempFilePath || res.tempFilePath
        })
      }
    })
  },

  // 取消媒体选择
  cancelMediaSelect() {
    this.setData({
      mediaType: 'text',
      mediaFile: null,
      mediaPreview: ''
    })
  },

  // 发送消息
  async sendMessage() {
    // 如果没有输入内容且没有选择媒体文件，不发送消息
    if (!this.data.messageInput.trim() && !this.data.mediaFile) {
      return
    }
    
    // 如果WebSocket未连接，提示用户
    if (!this.data.isConnected) {
      wx.showToast({
        title: '正在连接服务器，请稍后再试',
        icon: 'none'
      })
      // 尝试重新连接
      this.connectWebSocket()
      return
    }
    
    try {
      // 如果有媒体文件，先上传媒体文件
      if (this.data.mediaFile) {
        await this.uploadAndSendMedia()
      } else {
        // 发送文本消息
        this.sendTextMessage()
      }
    } catch (error) {
      console.error('发送消息失败:', error)
      wx.showToast({
        title: '发送失败，请重试',
        icon: 'none'
      })
    }
  },

  // 发送文本消息
  sendTextMessage() {
    const message = {
      type: 'chat',
      content: this.data.messageInput,
      mediaType: 'text',
      groupId: this.data.groupId,
      timestamp: new Date().toISOString()
    }
    
    // 发送消息
    this.sendWebSocketMessage(message)
    
    // 清空输入框
    this.setData({
      messageInput: ''
    })
  },

  // 上传并发送媒体消息
  async uploadAndSendMedia() {
    try {
      wx.showLoading({
        title: '正在上传...',
        mask: true
      })
      
      // 上传媒体文件
      const mediaURL = await this.uploadMedia()
      
      // 发送媒体消息
      const message = {
        type: 'chat',
        content: this.data.messageInput || '发送了一个媒体文件',
        mediaType: this.data.mediaType,
        mediaUrl: mediaURL,
        groupId: this.data.groupId,
        timestamp: new Date().toISOString()
      }
      
      // 发送消息
      this.sendWebSocketMessage(message)
      
      // 清空输入框和媒体选择
      this.setData({
        messageInput: '',
        mediaType: 'text',
        mediaFile: null,
        mediaPreview: ''
      })
      
      wx.hideLoading()
    } catch (error) {
      wx.hideLoading()
      throw error
    }
  },

  // 上传媒体文件
  uploadMedia() {
    return new Promise((resolve, reject) => {
      const token = wx.getStorageSync('authToken')
      if (!token) {
        reject(new Error('请先登录'))
        return
      }
      
      wx.uploadFile({
        url: `${BASE_URL}/api/uploads`,
        filePath: this.data.mediaFile.path,
        name: 'file',
        header: {
          'Authorization': token
        },
        formData: {
          type: this.data.mediaType
        },
        success: (res) => {
          try {
            const data = JSON.parse(res.data)
            if (data.code === 200 && data.data) {
              // 返回完整的文件URL
              resolve(this.getImageUrl(data.data.path))
            } else {
              reject(new Error(data.message || '上传失败'))
            }
          } catch (error) {
            reject(error)
          }
        },
        fail: (error) => {
          reject(error)
        }
      })
    })
  },

  // 建立WebSocket连接
  connectWebSocket() {
    if (!this.data.groupId) return
    
    // 如果已经连接，不要重复连接
    if (this.socketTask && this.data.isConnected) {
      console.log('WebSocket已连接，跳过连接')
      return
    }
    
    // 如果正在连接中，不要重复连接
    if (this.isConnecting) {
      console.log('WebSocket正在连接中，跳过重复连接')
      return
    }
    
    this.isConnecting = true
    
    try {
      // 获取认证token
      const token = wx.getStorageSync('authToken')
      if (!token) {
        wx.showToast({
          title: '请先登录',
          icon: 'none'
        })
        this.isConnecting = false
        return
      }
      
      // 将 http:// 替换为 ws://
      const wsUrl = `${BASE_URL.replace(/^http/, 'ws')}/api/ws/groupChat/${this.data.groupId}`
      
      console.log('正在连接WebSocket:', wsUrl)
      
      // 关闭现有连接
      if (this.socketTask) {
        this.closeWebSocket()
      }
      
      // 创建WebSocket连接
      this.socketTask = wx.connectSocket({
        url: wsUrl,
        success: () => {
          console.log('WebSocket连接创建成功')
        },
        fail: (error) => {
          console.error('WebSocket连接创建失败:', error)
          this.setData({ isConnected: false })
          this.isConnecting = false
          
          wx.showToast({
            title: '连接服务器失败，请稍后再试',
            icon: 'none'
          })
        }
      })
      
      // 监听WebSocket连接打开
      this.socketTask.onOpen(() => {
        console.log('WebSocket连接已打开，发送认证消息')
        this.isConnecting = false
        
        // 发送认证消息
        const authMessage = {
          token: token
        }
        
        // 使用setTimeout确保连接完全建立后再发送消息
        setTimeout(() => {
          if (this.socketTask && !this.isClosing) {
            this.socketTask.send({
              data: JSON.stringify(authMessage),
              success: () => {
                console.log('认证消息发送成功')
              },
              fail: (error) => {
                console.error('认证消息发送失败:', error)
              }
            })
          }
        }, 500)
        
        // 启动心跳
        this.startHeartbeat()
      })
      
      // 监听WebSocket接收到服务器的消息
      this.socketTask.onMessage((res) => {
        try {
          const data = JSON.parse(res.data)
          
          // 处理认证成功消息
          if (data.code === 200 && (data.message === '认证成功' || data.data === '认证成功' || data.data === '连接成功')) {
            console.log('WebSocket认证成功')
            this.setData({ isConnected: true })
            this.retryCount = 0
            return
          }
          
          // 处理pong响应，不做特殊处理
          if (data.type === 'pong') {
            console.log('收到pong响应')
            return
          }
          
          // 处理错误消息
          if (data.code !== undefined && data.code !== 200) {
            wx.showToast({
              title: data.message || '连接错误',
              icon: 'none'
            })
            
            if (data.code === 401) {
              // 认证失败，关闭连接
              this.closeWebSocket()
            }
            return
          }
          
          // 处理正常的聊天消息
          if (data.type) {
            this.handleWebSocketMessage(data)
          }
        } catch (error) {
          console.error('解析WebSocket消息失败:', error)
        }
      })
      
      // 监听WebSocket连接关闭
      this.socketTask.onClose((res) => {
        console.log('WebSocket连接已关闭:', res)
        this.setData({ isConnected: false })
        this.isConnecting = false
        this.stopHeartbeat()
        
        // 如果不是主动关闭，尝试重连
        if (!this.isClosing) {
          const delay = Math.min(1000 * Math.pow(2, this.retryCount || 0), 30000)
          this.retryCount = (this.retryCount || 0) + 1
          
          setTimeout(() => {
            if (!this.isClosing && !this.isConnecting) {
              console.log(`尝试第 ${this.retryCount} 次重连...`)
              this.connectWebSocket()
            }
          }, delay)
        }
      })
      
      // 监听WebSocket错误
      this.socketTask.onError((error) => {
        console.error('WebSocket错误:', error)
        this.setData({ isConnected: false })
        this.isConnecting = false
        
        // 不要在这里立即重连，让onClose处理重连逻辑
        wx.showToast({
          title: 'WebSocket连接错误',
          icon: 'none'
        })
      })
    } catch (error) {
      console.error('创建WebSocket连接失败:', error)
      this.setData({ isConnected: false })
      this.isConnecting = false
      
      wx.showToast({
        title: '创建WebSocket连接失败，请稍后重试',
        icon: 'none'
      })
    }
  },

  // 关闭WebSocket连接
  closeWebSocket() {
    this.stopHeartbeat()
    
    if (this.socketTask) {
      this.isClosing = true
      
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
            this.isClosing = false
            this.isConnecting = false
          }
        })
      } catch (error) {
        console.error('关闭WebSocket连接出错:', error)
        this.socketTask = null
        this.setData({ isConnected: false })
        this.isClosing = false
        this.isConnecting = false
      }
    }
  },

  // 发送WebSocket消息
  sendWebSocketMessage(message) {
    if (!this.socketTask || !this.data.isConnected) {
      wx.showToast({
        title: '连接已断开，正在重连...',
        icon: 'none'
      })
      this.connectWebSocket()
      return
    }
    
    try {
      this.socketTask.send({
        data: JSON.stringify(message),
        success: () => {
          console.log('消息发送成功:', message)
        },
        fail: (error) => {
          console.error('消息发送失败:', error)
          wx.showToast({
            title: '消息发送失败，请重试',
            icon: 'none'
          })
        }
      })
    } catch (error) {
      console.error('发送消息时发生错误:', error)
      wx.showToast({
        title: '消息发送失败，请重试',
        icon: 'none'
      })
      
      // 尝试重新连接
      this.setData({ isConnected: false })
      this.connectWebSocket()
    }
  },

  // 处理WebSocket消息
  handleWebSocketMessage(data) {
    // 处理不同类型的消息
    switch (data.type) {
      case 'chat':
        // 添加新消息到消息列表，需要转换字段名
        const chatMessage = {
          ...data,
          ID: data.id || Date.now(), // 生成临时ID
          SenderID: data.senderId,
          senderName: data.nickname,
          senderAvatar: this.getImageUrl(data.avatarUrl),
          Content: data.content,
          MediaType: data.mediaType,
          MediaURL: data.mediaUrl,
          CreatedAt: data.timestamp,
          type: data.type
        }
        
        this.setData({
          messages: [...this.data.messages, chatMessage]
        })
        
        // 滚动到底部
        this.scrollToBottom()
        break
        
      case 'join':
      case 'leave':
        // 处理系统消息（加入/离开）
        const systemMessage = {
          ...data,
          ID: data.id || Date.now(), // 生成临时ID
          MediaType: 'text',
          SenderID: data.senderId || 0,
          Content: data.content,
          CreatedAt: data.timestamp,
          type: data.type // 保留原始类型以便于模板中区分
        }
        
        this.setData({
          messages: [...this.data.messages, systemMessage]
        })
        
        // 滚动到底部
        this.scrollToBottom()
        break
        
      default:
        console.warn('未知的消息类型:', data.type)
    }
  },

  // 发送心跳消息，保持连接活跃
  startHeartbeat() {
    // 清除现有的心跳定时器
    this.stopHeartbeat()
    
    // 每30秒发送一次心跳
    this.heartbeatInterval = setInterval(() => {
      if (this.socketTask && this.data.isConnected && !this.isClosing && !this.isConnecting) {
        try {
          // 发送一个简单的ping消息
          this.socketTask.send({
            data: JSON.stringify({ type: 'ping' }),
            success: () => {
              console.log('发送心跳消息')
            },
            fail: (error) => {
              console.error('发送心跳消息失败:', error)
              // 不要在心跳失败时立即重连，让onClose处理重连逻辑
              this.setData({ isConnected: false })
            }
          })
        } catch (error) {
          console.error('发送心跳消息失败:', error)
          this.setData({ isConnected: false })
        }
      } else if (!this.isClosing && !this.isConnecting) {
        // 如果连接已断开且不是正在关闭或连接中，尝试重连
        this.connectWebSocket()
      }
    }, 30000)
  },

  // 停止心跳
  stopHeartbeat() {
    if (this.heartbeatInterval) {
      clearInterval(this.heartbeatInterval)
      this.heartbeatInterval = null
    }
  },

  // 滚动到底部
  scrollToBottom() {
    // 使用延时确保消息列表已经渲染完成
    setTimeout(() => {
      if (this.data.messages.length > 0) {
        const lastMessage = this.data.messages[this.data.messages.length - 1]
        this.setData({
          scrollToMessage: `msg-${lastMessage.ID}`
        })
      }
    }, 100)
  },

  // 预览图片
  previewImage(e) {
    const url = e.currentTarget.dataset.url
    wx.previewImage({
      current: url,
      urls: [url]
    })
  },

  // 格式化时间
  formatTime(timestamp) {
    if (!timestamp) return ''
    
    const date = new Date(timestamp)
    const now = new Date()
    const diff = now - date
    
    // 如果是今天的消息，只显示时间
    if (diff < 24 * 60 * 60 * 1000 && date.getDate() === now.getDate()) {
      return this.formatTimeOnly(date)
    }
    
    // 如果是昨天的消息，显示"昨天 时间"
    const yesterday = new Date(now)
    yesterday.setDate(now.getDate() - 1)
    if (date.getDate() === yesterday.getDate() && date.getMonth() === yesterday.getMonth() && date.getFullYear() === yesterday.getFullYear()) {
      return `昨天 ${this.formatTimeOnly(date)}`
    }
    
    // 如果是今年的消息，显示"月-日 时间"
    if (date.getFullYear() === now.getFullYear()) {
      return `${date.getMonth() + 1}-${date.getDate()} ${this.formatTimeOnly(date)}`
    }
    
    // 其他情况，显示完整日期
    return `${date.getFullYear()}-${date.getMonth() + 1}-${date.getDate()} ${this.formatTimeOnly(date)}`
  },

  // 只格式化时间部分
  formatTimeOnly(date) {
    const hours = date.getHours()
    const minutes = date.getMinutes()
    return `${hours < 10 ? '0' + hours : hours}:${minutes < 10 ? '0' + minutes : minutes}`
  },

  // 处理图片路径
  getImageUrl(path) {
    try {
      if (!path) {
        return `${BASE_URL}/api/uploads/default-avatar.png`;
      }
      
      // 如果已经是完整URL，直接返回
      if (path.startsWith('http')) {
        return path;
      }
      
      // 确保路径不包含重复的uploads前缀
      if (path.startsWith('uploads/')) {
        path = path.replace('uploads/', '');
      }
      
      // 与Web端保持一致，使用 /api/uploads/ 路径
      return `${BASE_URL}/api/uploads/${path}`;
    } catch (error) {
      console.error('处理图片路径出错:', error);
      return `${BASE_URL}/api/uploads/default-avatar.png`;
    }
  }
}) 