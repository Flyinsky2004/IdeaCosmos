const app = getApp()
const towxml = require('../../utils/towxml/index')

Page({
  data: {
    loading: true,
    chapter: null,
    project: null,
    parsedContent: '',
    audioUrl: '',
    comments: [],
    commentContent: '',
    isLoggedIn: false,

    // 字体大小选项
    fontSizes: [
      { label: '小', value: 'small' },
      { label: '中', value: 'medium' },
      { label: '大', value: 'large' }
    ],
    currentFontSize: 'medium',

    // 评论类型
    commentTypes: [
      { label: '全部评论', value: 'all' },
      { label: '读者评论', value: 'reader' },
      { label: '作者评论', value: 'author' }
    ],
    commentType: 'all',

    // 情绪映射
    emotionMap: {
      "喜悦": { 
        icon: "😊", 
        description: "充满快乐和满足" 
      },
      "感动": { 
        icon: "🥹", 
        description: "内心被深深触动" 
      },
      "惊喜": { 
        icon: "🤩", 
        description: "意外的惊喜" 
      },
      "期待": { 
        icon: "🤔", 
        description: "对后续充满期待" 
      },
      "伤感": { 
        icon: "😢", 
        description: "略带忧伤的感动" 
      },
      "愤怒": { 
        icon: "😠", 
        description: "对情节感到愤慨" 
      },
      "恐惧": { 
        icon: "😱", 
        description: "感到害怕或紧张" 
      },
      "平静": { 
        icon: "😐", 
        description: "内心平和安宁" 
      }
    },
    emotions: [], // 将在 onLoad 中初始化
    selectedEmotion: '',
    userFeeling: null
  },

  onLoad(options) {
    this.chapterId = options.id
    
    // 初始化情绪数组
    const emotions = Object.entries(this.data.emotionMap).map(([name, data]) => ({
      name,
      icon: data.icon,
      description: data.description
    }))
    
    this.setData({ 
      emotions,
      isLoggedIn: !!app.globalData.userInfo
    })
    
    this.fetchChapterDetail()
  },

  // 获取章节详情
  async fetchChapterDetail() {
    try {
      const res = await wx.cloud.callFunction({
        name: 'getChapterDetail',
        data: { id: this.chapterId }
      })

      if (res.result.code === 0) {
        const { chapter, project } = res.result.data
        
        // 解析 Markdown 内容
        const parsedContent = towxml(chapter.current_version.content, 'markdown', {
          theme: 'light',
          events: {
            tap: (e) => {
              // 处理链接点击等事件
              console.log('tap', e)
            }
          }
        })

        // 设置音频 URL
        const audioUrl = chapter.current_version.audio_path ? 
          `${app.globalData.BACKEND_DOMAIN}/audio/${chapter.current_version.audio_path}` : ''

        this.setData({
          chapter,
          project,
          parsedContent,
          audioUrl,
          loading: false
        })

        // 获取评论和情绪评价
        this.fetchComments()
        this.fetchUserFeeling()
      }
    } catch (error) {
      console.error('获取章节详情失败:', error)
      wx.showToast({
        title: '获取章节详情失败',
        icon: 'none'
      })
    }
  },

  // 获取评论列表
  async fetchComments() {
    if (!this.data.chapter?.current_version?.ID) return

    try {
      const res = await wx.cloud.callFunction({
        name: 'getVersionComments',
        data: {
          version_id: this.data.chapter.current_version.ID,
          type: this.data.commentType
        }
      })

      if (res.result.code === 0) {
        this.setData({
          comments: res.result.data
        })
      }
    } catch (error) {
      console.error('获取评论失败:', error)
    }
  },

  // 获取用户情绪评价
  async fetchUserFeeling() {
    if (!this.data.chapter?.current_version?.ID || !this.data.isLoggedIn) return

    try {
      const res = await wx.cloud.callFunction({
        name: 'getUserFeeling',
        data: {
          version_id: this.data.chapter.current_version.ID
        }
      })

      if (res.result.code === 0) {
        this.setData({
          userFeeling: res.result.data
        })
      }
    } catch (error) {
      console.error('获取情绪评价失败:', error)
    }
  },

  // 提交情绪评价
  async submitFeeling(e) {
    if (!this.data.isLoggedIn) {
      wx.showToast({
        title: '请先登录',
        icon: 'none'
      })
      return
    }

    const emotion = e.currentTarget.dataset.emotion
    this.setData({ selectedEmotion: emotion })

    try {
      const res = await wx.cloud.callFunction({
        name: 'addFeeling',
        data: {
          version_id: this.data.chapter.current_version.ID,
          feeling: emotion
        }
      })

      if (res.result.code === 0) {
        wx.showToast({
          title: '评价成功',
          icon: 'success'
        })
        this.fetchUserFeeling()
      }
    } catch (error) {
      console.error('提交情绪评价失败:', error)
      wx.showToast({
        title: '评价失败',
        icon: 'none'
      })
    }
  },

  // 评论相关方法
  onCommentInput(e) {
    this.setData({
      commentContent: e.detail.value
    })
  },

  async submitComment() {
    if (!this.data.isLoggedIn) {
      wx.showToast({
        title: '请先登录',
        icon: 'none'
      })
      return
    }

    if (!this.data.commentContent.trim()) {
      wx.showToast({
        title: '请输入评论内容',
        icon: 'none'
      })
      return
    }

    try {
      const res = await wx.cloud.callFunction({
        name: 'addVersionComment',
        data: {
          version_id: this.data.chapter.current_version.ID,
          content: this.data.commentContent,
          type: 'reader'
        }
      })

      if (res.result.code === 0) {
        wx.showToast({
          title: '评论成功',
          icon: 'success'
        })
        
        this.setData({
          commentContent: ''
        })
        
        this.fetchComments()
      }
    } catch (error) {
      console.error('提交评论失败:', error)
      wx.showToast({
        title: '评论失败',
        icon: 'none'
      })
    }
  },

  // 切换评论类型
  changeCommentType(e) {
    const type = e.currentTarget.dataset.type
    this.setData({ commentType: type })
    this.fetchComments()
  },

  // 字体大小控制
  changeFontSize(e) {
    const size = e.currentTarget.dataset.size
    this.setData({ currentFontSize: size })
  },

  // 音频相关方法
  onAudioPlay() {
    // 可以在这里添加音频播放的统计等逻辑
  },

  onAudioError() {
    wx.showToast({
      title: '音频加载失败',
      icon: 'none'
    })
  },

  // 导航方法
  navigateBack() {
    wx.navigateBack()
  },

  navigateToLogin() {
    wx.navigateTo({
      url: '/pages/login/index'
    })
  },

  // 格式化日期
  formatDate(date) {
    const d = new Date(date)
    return `${d.getFullYear()}年${d.getMonth() + 1}月${d.getDate()}日 ${d.getHours()}:${d.getMinutes()}`
  },

  onShareAppMessage() {
    return {
      title: this.data.chapter?.Title || '精彩章节',
      path: `/pages/chapter-detail/index?id=${this.chapterId}`
    }
  }
}) 