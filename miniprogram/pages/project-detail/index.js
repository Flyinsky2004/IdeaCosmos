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
    if (!options.id) {
      wx.showToast({
        title: '项目ID无效',
        icon: 'none'
      })
      setTimeout(() => {
        wx.navigateBack()
      }, 1500)
      return
    }
    this.projectId = options.id
    this.fetchProjectDetail()
  },

  onShow() {
    // 强制刷新数据，避免缓存问题
    if (this.projectId) {
      this.setData({
        project: null,
        characters: [],
        chapters: [],
        comments: [],
        loading: true
      })
      this.fetchProjectDetail()
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
          console.error(`请求失败 ${url}:`, error) // 添加错误日志
          reject(error)
        }
      })
    })
  },

  // 获取项目详情
  async fetchProjectDetail() {
    if (!this.projectId) return
    
    try {
      this.setData({ loading: true })
      
      const params = {
        id: this.projectId
      }
      
      if (app.globalData.userInfo) {
        params.user_id = app.globalData.userInfo.ID
      }

      const res = await this.request('/api/public/getProjectDetail', 'GET', params)

      if (res.code === 200 && res.data) {
        // 处理图片路径
        const project = {
          ...res.data,
          cover_image: this.getImageUrl(res.data.cover_image),
          team: {
            ...res.data.team,
            avatar: this.getImageUrl(res.data.team?.avatar)
          }
        }
        wx.setNavigationBarTitle({
          title: project.project_name
        });
        
        this.setData({
          project,
          loading: false
        })
        
        // 获取其他相关数据
        await Promise.all([
          this.fetchCharacters(),
          this.fetchChapters(),
          this.fetchComments(),
          this.checkFavorite()
        ])
      } else {
        throw new Error(res.message || '获取项目详情失败')
      }
    } catch (error) {
      console.error('获取项目详情失败:', error)
      wx.showToast({
        title: error.message || '获取项目详情失败',
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
        // 确保 current_version 字段存在
        const chapters = res.data.map(chapter => ({
          ...chapter,
          current_version: chapter.current_version || { ID: 0 }
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
        const comments = res.data.map(comment => {
          // 处理用户头像路径
          let avatarPath = comment.user?.avatar || '';
          
          // 如果头像路径包含uploads/前缀，则去除
          if (avatarPath.startsWith('uploads/')) {
            avatarPath = avatarPath.replace('uploads/', '');
          }
          
          return {
            ...comment,
            user: {
              ...comment.user,
              avatar: this.getImageUrl(avatarPath)
            }
          };
        });
        
        this.setData({ comments });
      }
    } catch (error) {
      console.error('获取评论列表失败:', error)
    }
  },

  // 检查收藏状态
  async checkFavorite() {
    // 检查是否有token
    const token = wx.getStorageSync('authToken')
    if (!token) {
      // 未登录状态下不检查收藏状态
      this.setData({ isFavorited: false })
      return
    }

    try {
      const res = await this.request('/api/user/favorite/check', 'GET', {
        project_id: this.projectId
      }, true) // 需要授权

      if (res.code === 200) {
        this.setData({
          isFavorited: res.data
        })
      }
    } catch (error) {
      console.error('检查收藏状态失败:', error)
      // 如果是授权问题，不显示错误提示
      if (error.message !== '请先登录') {
        wx.showToast({
          title: '检查收藏状态失败',
          icon: 'none'
        })
      }
    }
  },

  // 切换收藏状态
  async toggleFavorite() {
    // 检查是否有token
    const token = wx.getStorageSync('authToken')
    if (!token) {
      wx.showToast({
        title: '请先登录',
        icon: 'none'
      })
      
      // 可以选择跳转到登录页面
      setTimeout(() => {
        wx.navigateTo({
          url: '/pages/profile/index'
        })
      }, 1500)
      return
    }

    try {
      const api = this.data.isFavorited ? '/api/user/favorite/remove' : '/api/user/favorite/add'
      const res = await this.request(api, 'POST', {
        project_id: this.projectId
      }, true) // 需要授权

      if (res.code === 200) {
        // 更新收藏状态
        const newFavoritedState = !this.data.isFavorited
        
        // 更新收藏数量
        const newFavorites = this.data.project.favorites + (newFavoritedState ? 1 : -1)
        
        this.setData({
          isFavorited: newFavoritedState,
          'project.favorites': newFavorites
        })
        
        wx.showToast({
          title: newFavoritedState ? '收藏成功' : '已取消收藏',
          icon: 'success'
        })
      } else {
        throw new Error(res.message || '操作失败')
      }
    } catch (error) {
      console.error('操作收藏失败:', error)
      wx.showToast({
        title: error.message || '操作失败',
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
    // 检查是否有token
    const token = wx.getStorageSync('authToken')
    if (!token) {
      wx.showToast({
        title: '请先登录',
        icon: 'none'
      })
      
      // 可以选择跳转到登录页面
      setTimeout(() => {
        wx.navigateTo({
          url: '/pages/profile/index'
        })
      }, 1500)
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
      }, true) // 需要授权

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
      } else {
        throw new Error(res.message || '评论失败')
      }
    } catch (error) {
      console.error('提交评论失败:', error)
      wx.showToast({
        title: error.message || '评论失败',
        icon: 'none'
      })
    }
  },

  // 跳转到章节详情
  navigateToChapter(e) {
    const { id } = e.currentTarget.dataset
    const chapter = this.data.chapters.find(c => c.ID === id)
    
    if (!chapter.current_version || chapter.current_version.ID === 0) {
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
      
      return `${BASE_URL}/api/uploads/${path}`;
    } catch (error) {
      console.error('处理图片路径出错:', error);
      return `${BASE_URL}/api/uploads/default-avatar.png`;
    }
  },

  onShareAppMessage() {
    return {
      title: this.data.project?.project_name || '精彩项目',
      path: `/pages/project-detail/index?id=${this.projectId}`
    }
  }
}) 