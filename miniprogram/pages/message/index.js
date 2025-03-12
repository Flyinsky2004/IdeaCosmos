/*
 * @Author: Flyinsky w2084151024@gmail.com
 * @Description: None
 */
// pages/message/index.js
const app = getApp()
const BASE_URL = 'https://idea.1024110.xyz' // 开发环境

Page({

  /**
   * 页面的初始数据
   */
  data: {
    // 状态
    loading: true,
    groups: [],
    filteredGroups: [],
    searchText: '',
    activeGroupId: null,
    userInfo: null,
    
    // 创建群组相关
    createModalVisible: false,
    createGroupForm: {
      name: '',
      description: '',
      avatarUrl: ''
    }
  },

  /**
   * 生命周期函数--监听页面加载
   */
  onLoad(options) {
    // 获取用户信息
    this.setData({
      userInfo: app.globalData.userInfo || {}
    })
    
    // 获取群组列表
    this.fetchGroups()
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
    // 每次显示页面时刷新数据
    this.fetchGroups()
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

  },

  /**
   * 页面相关事件处理函数--监听用户下拉动作
   */
  onPullDownRefresh() {

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

  },

  // 处理搜索输入
  onSearchInput(e) {
    this.setData({
      searchText: e.detail.value
    })
    this.filterGroups()
  },
  
  // 处理群组名称输入
  onNameInput(e) {
    this.setData({
      'createGroupForm.name': e.detail.value
    })
  },
  
  // 处理群组描述输入
  onDescriptionInput(e) {
    this.setData({
      'createGroupForm.description': e.detail.value
    })
  },
  
  // 处理群组头像URL输入
  onAvatarUrlInput(e) {
    this.setData({
      'createGroupForm.avatarUrl': e.detail.value
    })
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
  
  // 获取群组列表
  async fetchGroups() {
    try {
      this.setData({ loading: true })
      
      // 检查是否有token
      const token = wx.getStorageSync('authToken')
      if (!token) {
        // 未登录状态下显示提示
        this.setData({ 
          loading: false,
          groups: [],
          filteredGroups: []
        })
        return
      }
      
      // 修改API路径，与Web端保持一致
      const res = await this.request('/api/chat/groups', 'GET', {}, true)
      
      if (res.code === 200 && res.data) {
        // 处理群组数据，注意Web端返回的是 data.groups
        const groups = (res.data.groups || res.data).map(group => ({
          ...group,
          AvatarURL: this.getImageUrl(group.AvatarURL)
        }))
        
        this.setData({
          groups,
          filteredGroups: groups,
          loading: false
        })
        
        this.filterGroups()
      } else {
        throw new Error(res.message || '获取群组列表失败')
      }
    } catch (error) {
      console.error('获取群组列表失败:', error)
      wx.showToast({
        title: error.message || '获取群组列表失败',
        icon: 'none'
      })
      
      this.setData({ 
        loading: false,
        groups: [],
        filteredGroups: []
      })
    }
  },
  
  // 清除搜索
  clearSearch() {
    this.setData({
      searchText: ''
    })
    this.filterGroups()
  },
  
  // 过滤群组
  filterGroups() {
    const { groups, searchText } = this.data
    
    if (!searchText) {
      this.setData({
        filteredGroups: groups
      })
      return
    }
    
    const filtered = groups.filter(group => 
      group.Name.toLowerCase().includes(searchText.toLowerCase())
    )
    
    this.setData({
      filteredGroups: filtered
    })
  },
  
  // 选择群组
  selectGroup(e) {
    const groupId = e.currentTarget.dataset.id
    const selectedGroup = this.data.groups.find(g => g.ID === groupId)
    
    if (!selectedGroup) return
    
    this.setData({
      activeGroupId: groupId
    })
    
    // 获取群组详情，与Web端保持一致
    this.request(`/api/chat/groups/${groupId}`, 'GET', {}, true)
      .then(res => {
        if (res.code === 200 && res.data) {
          const group = res.data.group || res.data
          
          console.log('跳转到群组聊天页面:', groupId, group.Name)
          
          // 跳转到群组聊天页面
          wx.navigateTo({
            url: `/pages/group-chat/index?id=${groupId}&name=${encodeURIComponent(group.Name)}`,
            fail: (err) => {
              console.error('跳转失败:', err)
              wx.showToast({
                title: '跳转失败，请重试',
                icon: 'none'
              })
            }
          })
        }
      })
      .catch(error => {
        console.error('获取群组详情失败:', error)
        wx.showToast({
          title: '获取群组详情失败',
          icon: 'none'
        })
      })
  },
  
  // 打开创建群组模态框
  openCreateModal() {
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
    
    this.setData({
      createModalVisible: true,
      createGroupForm: {
        name: '',
        description: '',
        avatarUrl: ''
      }
    })
  },
  
  // 关闭创建群组模态框
  closeCreateModal() {
    this.setData({
      createModalVisible: false
    })
  },
  
  // 创建群组
  async handleCreateGroup() {
    const { name, description, avatarUrl } = this.data.createGroupForm
    
    // 调试输出
    console.log('创建群组表单数据:', this.data.createGroupForm)
    
    if (!name.trim()) {
      wx.showToast({
        title: '请输入群组名称',
        icon: 'none'
      })
      return
    }
    
    try {
      // 修改API路径和参数，与Web端保持一致
      const res = await this.request('/api/chat/groups', 'POST', {
        name: name.trim(),
        description: description.trim(),
        avatarUrl: avatarUrl.trim()
      }, true)
      
      if (res.code === 200) {
        wx.showToast({
          title: '创建成功',
          icon: 'success'
        })
        
        // 关闭模态框
        this.setData({
          createModalVisible: false
        })
        
        // 重新获取群组列表
        this.fetchGroups()
      } else {
        throw new Error(res.message || '创建群组失败')
      }
    } catch (error) {
      console.error('创建群组失败:', error)
      wx.showToast({
        title: error.message || '创建群组失败',
        icon: 'none'
      })
    }
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