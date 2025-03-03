const app = getApp()
const BASE_URL = 'https://idea.1024110.xyz' // 开发环境

Page({
  data: {
    loading: true,
    project: null,
    characters: [],
    chapters: [],
    comments: [],
    commentContent: '',
    isFavorited: false
  },

  onLoad(options) {
    this.projectId = options.id
    this.fetchProjectDetail()
  },

  onShow() {
    if (this.projectId) {
      this.fetchProjectDetail()
    }
  },

  // 封装请求函数
  request(url, method, data) {
    return new Promise((resolve, reject) => {
      wx.request({
        url: `${BASE_URL}${url}`,
        method: method,
        data: data,
        header: {
          'content-type': 'application/json',
          'Cache-Control': 'no-cache'  // 添加no-cache头
        },
        success: (res) => {
          if (res.statusCode === 200) {
            resolve(res.data)
          } else {
            reject(res)
          }
        },
        fail: reject
      })
    })
  },

  // 获取项目详情
  async fetchProjectDetail() {
    try {
      const params = {
        project_id: this.projectId
      }
      
      if (app.globalData.userInfo) {
        params.user_id = app.globalData.userInfo.ID
      }

      const res = await this.request('/api/public/getProjectDetail', 'GET', params)

      if (res.code === 200) {
        // 处理图片路径
        const project = {
          ...res.data,
          cover_image: this.getImageUrl(res.data.cover_image),
          team: {
            ...res.data.team,
            avatar: this.getImageUrl(res.data.team?.avatar)
          }
        }

        this.setData({
          project,
          loading: false
        })
        
        // 获取其他相关数据
        this.fetchCharacters()
        this.fetchChapters()
        this.fetchComments()
        this.checkFavorite()
      } else {
        wx.showToast({
          title: res.message || '获取项目详情失败',
          icon: 'none'
        })
      }
    } catch (error) {
      console.error('获取项目详情失败:', error)
      wx.showToast({
        title: '获取项目详情失败',
        icon: 'none'
      })
    } finally {
      this.setData({ loading: false })
    }
  },

  // 获取角色列表
  async fetchCharacters() {
    try {
      const res = await this.request('/api/public/getProjectCharacters', 'GET', {
        id: this.projectId
      })

      if (res.code === 200) {
        const characters = res.data.map(character => ({
          ...character,
          avatar: this.getImageUrl(character.avatar)
        }))
        this.setData({ characters })
      }
    } catch (error) {
      console.error('获取角色列表失败:', error)
    }
  },

  // 获取章节列表
  async fetchChapters() {
    try {
      const res = await this.request('/api/public/getProjectChapters', 'GET', {
        id: this.projectId
      })

      if (res.code === 200) {
        // 确保 CurrentVersion 字段存在
        const chapters = res.data.map(chapter => ({
          ...chapter,
          CurrentVersion: chapter.CurrentVersion || { ID: 0 }
        }))
        
        this.setData({ chapters })
      }
    } catch (error) {
      console.error('获取章节列表失败:', error)
    }
  },

  // 获取评论列表
  async fetchComments() {
    try {
      const res = await this.request('/api/public/getProjectComments', 'GET', {
        project_id: this.projectId
      })

      if (res.code === 200) {
        const comments = res.data.map(comment => ({
          ...comment,
          user: {
            ...comment.user,
            avatar: this.getImageUrl(comment.user?.avatar)
          }
        }))
        this.setData({ comments })
      }
    } catch (error) {
      console.error('获取评论列表失败:', error)
    }
  },

  // 检查收藏状态
  async checkFavorite() {
    if (!app.globalData.userInfo) return

    try {
      const res = await this.request('/api/user/favorite/check', 'GET', {
        project_id: this.projectId
      })

      if (res.code === 200) {
        this.setData({
          isFavorited: res.data
        })
      }
    } catch (error) {
      console.error('检查收藏状态失败:', error)
    }
  },

  // 切换收藏状态
  async toggleFavorite() {
    if (!app.globalData.userInfo) {
      wx.showToast({
        title: '请先登录',
        icon: 'none'
      })
      return
    }

    try {
      const api = this.data.isFavorited ? '/api/user/favorite/remove' : '/api/user/favorite/add'
      const res = await this.request(api, 'GET', {
        project_id: this.projectId
      })

      if (res.code === 200) {
        this.setData({
          isFavorited: !this.data.isFavorited
        })
        
        wx.showToast({
          title: this.data.isFavorited ? '收藏成功' : '已取消收藏',
          icon: 'success'
        })
        
        // 更新项目信息
        this.fetchProjectDetail()
      }
    } catch (error) {
      console.error('操作收藏失败:', error)
      wx.showToast({
        title: '操作失败',
        icon: 'none'
      })
    }
  },

  // 评论输入
  onCommentInput(e) {
    this.setData({
      commentContent: e.detail.value
    })
  },

  // 提交评论
  async submitComment() {
    if (!app.globalData.userInfo) {
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
      const res = await this.request('/api/user/addProjectComment', 'POST', {
        content: this.data.commentContent,
        project_id: this.projectId
      })

      if (res.code === 200) {
        wx.showToast({
          title: '评论成功',
          icon: 'success'
        })
        
        this.setData({
          commentContent: ''
        })
        
        // 重新获取评论列表
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

  // 跳转到章节详情
  navigateToChapter(e) {
    const { id } = e.currentTarget.dataset
    const chapter = this.data.chapters.find(c => c.ID === id)
    
    if (!chapter.CurrentVersion || chapter.CurrentVersion.ID === 0) {
      wx.showToast({
        title: '该章节尚未开始创作',
        icon: 'none'
      })
      return
    }
    
    wx.navigateTo({
      url: `/pages/chapter-detail/index?id=${id}`
    })
  },

  // 处理图片路径
  getImageUrl(path) {
    try {
      if (!path) {
        return `${BASE_URL}/api/uploads/default-avatar.png`
      }
      
      if (path.startsWith('http')) {
        return path
      }
      
      return `${BASE_URL}/api/uploads/${path}`
    } catch (error) {
      console.error('处理图片路径出错:', error)
      return `${BASE_URL}/api/uploads/default-avatar.png`
    }
  },

  onShareAppMessage() {
    return {
      title: this.data.project?.project_name || '精彩项目',
      path: `/pages/project-detail/index?id=${this.projectId}`
    }
  }
}) 