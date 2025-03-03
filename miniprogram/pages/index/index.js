// index.js
const defaultAvatarUrl = 'https://mmbiz.qpic.cn/mmbiz/icTdbqWNOwNRna42FI242Lcia07jQodd2FJGIYQfG0LAJGFxM4FbnQP6yfMxBgJ0F3YRqJCJ1aPAK2dQagdusBZg/0'

// 添加基础配置
const BASE_URL = 'https://idea.1024110.xyz' // 开发环境
// const BASE_URL = 'https://your-production-domain.com' // 生产环境

// 添加分类数据
const categories = [
  "幽默", "悬疑", "黑暗", "科幻", "奇幻", "浪漫", "荒诞", "励志", "动作", "恐怖",
  "史诗", "温情", "讽刺", "文艺", "纪实", "冒险", "家庭剧", "超现实", "战争", "公路",
  "青春", "复仇", "政治", "犯罪", "悬幻", "心理", "怪诞", "温暖", "现实主义", "虚构",
  "哲学", "校园", "灾难", "武侠", "神秘", "励志成长", "古风", "穿越", "音乐", "童话"
].map(name => ({
  name,
  icon: 'success' // 小程序内置图标
}))

Page({
  data: {
    motto: 'Hello World',
    userInfo: {
      avatarUrl: defaultAvatarUrl,
      nickName: '',
    },
    hasUserInfo: false,
    canIUseGetUserProfile: wx.canIUse('getUserProfile'),
    canIUseNicknameComp: wx.canIUse('input.type.nickname'),
    currentTab: 'latest',
    projects: [],
    loading: false,
    isRefreshing: false,
    page: 1,
    hasMore: true,
    categories: categories,
    selectedCategory: null,
    categoryProjects: [],
    categoryPage: 1,
    categoryHasMore: true
  },
  bindViewTap() {
    wx.navigateTo({
      url: '../logs/logs'
    })
  },
  onChooseAvatar(e) {
    const { avatarUrl } = e.detail
    const { nickName } = this.data.userInfo
    this.setData({
      "userInfo.avatarUrl": avatarUrl,
      hasUserInfo: nickName && avatarUrl && avatarUrl !== defaultAvatarUrl,
    })
  },
  onInputChange(e) {
    const nickName = e.detail.value
    const { avatarUrl } = this.data.userInfo
    this.setData({
      "userInfo.nickName": nickName,
      hasUserInfo: nickName && avatarUrl && avatarUrl !== defaultAvatarUrl,
    })
  },
  getUserProfile(e) {
    // 推荐使用wx.getUserProfile获取用户信息，开发者每次通过该接口获取用户个人信息均需用户确认，开发者妥善保管用户快速填写的头像昵称，避免重复弹窗
    wx.getUserProfile({
      desc: '展示用户信息', // 声明获取用户个人信息后的用途，后续会展示在弹窗中，请谨慎填写
      success: (res) => {
        console.log(res)
        this.setData({
          userInfo: res.userInfo,
          hasUserInfo: true
        })
      }
    })
  },
  onLoad() {
    this.loadProjects()
  },
  // 封装请求函数
  request(url, method, data) {
    return new Promise((resolve, reject) => {
      wx.request({
        url: `${BASE_URL}${url}`,
        method: method,
        data: data,
        header: {
          'content-type': 'application/json'
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
  // 切换标签
  async switchTab(e) {
    const tab = e.currentTarget.dataset.tab
    if (this.data.currentTab === tab) return
    
    this.setData({
      currentTab: tab,
      projects: [],
      page: 1,
      hasMore: true
    })
    this.loadProjects()
  },
  // 加载项目列表 - 对齐 Vue3 的 API 路径
  async loadProjects() {
    if (this.data.loading || !this.data.hasMore) return

    this.setData({ loading: true })
    try {
      let url = ''
      switch (this.data.currentTab) {
        case 'hot':
          url = '/api/public/hot-projects'
          break
        case 'latest':
          url = '/api/public/getIndexProject'
          break
        case 'following':
          url = '/api/public/getFollowingProjects'
          break
        default:
          url = '/api/public/getIndexProject'
      }

      const res = await this.request(url, 'GET', {
        pageIndex: this.data.page - 1 // 页码对齐
      })

      if (res.code === 200) { // 假设后端返回格式为 { code: 200, message: 'success', data: [] }
        const newProjects = res.data.map(project => ({
          ...project,
          coverImageUrl: this.getImageUrl(project.cover_image),
          teamAvatarUrl: this.getImageUrl(project.team?.avatar)
        }))
        console.log('处理后的项目数据：', newProjects)
        
        this.setData({
          projects: [...this.data.projects, ...newProjects],
          page: this.data.page + 1,
          hasMore: newProjects.length === 10, // 假设每页10条
          loading: false
        })
      } else {
        wx.showToast({
          title: res.message || '加载失败',
          icon: 'none'
        })
      }
    } catch (err) {
      console.error('加载项目失败：', err)
      wx.showToast({
        title: '网络错误',
        icon: 'none'
      })
    } finally {
      this.setData({ loading: false })
    }
  },
  // 下拉刷新
  async onRefresh() {
    this.setData({
      isRefreshing: true,
      page: 1,
      hasMore: true,
      projects: []
    })
    
    await this.loadProjects()
    this.setData({ isRefreshing: false })
  },
  // 上拉加载更多
  loadMore() {
    if (this.data.currentTab === 'categories') {
      this.loadCategoryProjects()
    } else {
      this.loadProjects()
    }
  },
  // 搜索输入 - 对齐搜索 API
  onSearchInput(e) {
    clearTimeout(this.searchTimer)
    this.searchTimer = setTimeout(() => {
      const keyword = e.detail.value
      if (keyword) {
        wx.navigateTo({
          url: `/pages/search/index?keyword=${encodeURIComponent(keyword)}`
        })
      }
    }, 500)
  },
  // 跳转到详情页
  goToDetail(e) {
    const id = e.currentTarget.dataset.id
    wx.navigateTo({
      url: `/pages/project-detail/index?id=${id}`
    })
  },
  // 处理项目图片路径
  getImageUrl: function(path) {
    try {
      // 如果路径为空，返回默认图片
      if (!path) {
        return 'https://idea.1024110.xyz/api/uploads/default-avatar.png';
      }
      
      // 如果路径已经包含完整URL，直接返回
      if (path.startsWith('http')) {
        return path;
      }
      
      // 否则拼接基础URL
      return `https://idea.1024110.xyz/api/uploads/${path}`;
    } catch (error) {
      console.error('处理图片路径出错:', error);
      return 'https://idea.1024110.xyz/api/uploads/default-avatar.png';
    }
  },
  // 选择分类
  selectCategory(e) {
    const category = e.currentTarget.dataset.category
    if (this.data.selectedCategory === category) return

    this.setData({
      selectedCategory: category,
      categoryProjects: [],
      categoryPage: 1,
      categoryHasMore: true
    })
    this.loadCategoryProjects()
  },
  // 加载分类项目
  async loadCategoryProjects() {
    if (this.data.loading || !this.data.categoryHasMore) return

    this.setData({ loading: true })
    try {
      const res = await this.request('/api/public/getCategoryProjects', 'GET', {
        category: this.data.selectedCategory,
        pageIndex: this.data.categoryPage - 1
      })

      if (res.code === 200) {
        const newProjects = res.data.map(project => ({
          ...project,
          coverImageUrl: this.getImageUrl(project.cover_image),
          teamAvatarUrl: this.getImageUrl(project.team?.avatar)
        }))

        this.setData({
          categoryProjects: [...this.data.categoryProjects, ...newProjects],
          categoryPage: this.data.categoryPage + 1,
          categoryHasMore: newProjects.length === 12,
          loading: false
        })
      } else {
        wx.showToast({
          title: res.message || '加载失败',
          icon: 'none'
        })
      }
    } catch (err) {
      console.error('加载分类项目失败：', err)
      wx.showToast({
        title: '网络错误',
        icon: 'none'
      })
    } finally {
      this.setData({ loading: false })
    }
  }
})
