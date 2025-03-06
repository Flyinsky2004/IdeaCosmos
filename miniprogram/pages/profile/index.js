// pages/profile/index.js
const app = getApp()
const BACKEND_DOMAIN = 'https://idea.1024110.xyz/'

Page({

  /**
   * 页面的初始数据
   */
  data: {
    isLoggedIn: false,
    isEditing: false,
    loginForm: {
      username: '',
      password: ''
    },
    userInfo: {},
    editForm: {
      username: ''
    }
  },

  /**
   * 生命周期函数--监听页面加载
   */
  onLoad(options) {
    this.checkLoginStatus()
    this.setData({
      loginForm: {
        username: '',
        password: ''
      }
    })
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
    if (this.data.isLoggedIn) {
      this.fetchUserInfo()
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

  // 检查登录状态
  checkLoginStatus() {
    const token = wx.getStorageSync('authToken')
    this.setData({
      isLoggedIn: !!token
    })
    if (token) {
      this.fetchUserInfo()
    }
  },

  // 处理登录
  handleLogin() {
    const { username, password } = this.data.loginForm
    console.log(this.data.loginForm)
    if (!username || !password) {
      wx.showToast({
        title: '账号和密码不能为空',
        icon: 'none'
      })
      return
    }

    wx.request({
      url: BACKEND_DOMAIN + 'api/auth/login',
      method: 'POST',
      data: this.data.loginForm,
      success: (res) => {
        if (res.data.code === 200) {
          wx.setStorageSync('authToken', res.data.data)
          wx.showToast({
            title: '登录成功！欢迎回来！',
            icon: 'success'
          })
          this.setData({
            isLoggedIn: true
          })
          this.fetchUserInfo()
        } else {
          wx.showToast({
            title: res.data.message,
            icon: 'none'
          })
        }
      },
      fail: () => {
        wx.showToast({
          title: '网络错误，请稍后重试',
          icon: 'none'
        })
      }
    })
  },

  // 获取用户信息
  fetchUserInfo() {
    wx.request({
      url: BACKEND_DOMAIN + 'api/user/me',
      method: 'GET',
      header: {
        'Authorization': wx.getStorageSync('authToken')
      },
      success: (res) => {
        if (res.data.code === 200) {
          const userInfo = res.data.data
          this.setData({
            userInfo,
            'editForm.username': userInfo.username
          })
        } else if (res.data.code === 401) {
          this.handleLogout()
        } else {
          wx.showToast({
            title: res.data.message,
            icon: 'none'
          })
        }
      },
      fail: () => {
        wx.showToast({
          title: '网络错误，请稍后重试',
          icon: 'none'
        })
      }
    })
  },

  // 处理编辑/保存
  handleEdit() {
    if (this.data.isEditing) {
      // 保存修改
      wx.request({
        url: BACKEND_DOMAIN + 'api/user/updateInfo',
        method: 'POST',
        header: {
          'Authorization': wx.getStorageSync('authToken')
        },
        data: {
          username: this.data.editForm.username
        },
        success: (res) => {
          if (res.data.code === 200) {
            wx.showToast({
              title: '更新成功',
              icon: 'success'
            })
            this.fetchUserInfo()
            this.setData({
              isEditing: false
            })
          } else {
            wx.showToast({
              title: res.data.message,
              icon: 'none'
            })
          }
        },
        fail: () => {
          wx.showToast({
            title: '网络错误，请稍后重试',
            icon: 'none'
          })
        }
      })
    } else {
      this.setData({
        isEditing: true
      })
    }
  },

  // 显示头像选择
  showChangeAvatar() {
    wx.chooseMedia({
      count: 1,
      mediaType: ['image'],
      sourceType: ['album', 'camera'],
      success: (res) => {
        const tempFilePath = res.tempFiles[0].tempFilePath
        this.uploadAvatar(tempFilePath)
      }
    })
  },

  // 上传头像
  uploadAvatar(filePath) {
    wx.uploadFile({
      url: BACKEND_DOMAIN + 'api/user/uploadImage',
      filePath: filePath,
      name: 'image',
      header: {
        'Authorization': wx.getStorageSync('authToken')
      },
      success: (res) => {
        const data = JSON.parse(res.data)
        if (data.code === 200) {
          this.updateUserInfo({
            avatar: data.data.path
          })
        } else {
          wx.showToast({
            title: '上传失败',
            icon: 'none'
          })
        }
      },
      fail: () => {
        wx.showToast({
          title: '上传失败',
          icon: 'none'
        })
      }
    })
  },

  // 更新用户信息
  updateUserInfo(data) {
    wx.request({
      url: BACKEND_DOMAIN + 'api/user/updateInfo',
      method: 'POST',
      header: {
        'Authorization': wx.getStorageSync('authToken')
      },
      data: data,
      success: (res) => {
        if (res.data.code === 200) {
          wx.showToast({
            title: '更新成功',
            icon: 'success'
          })
          this.fetchUserInfo()
        } else {
          wx.showToast({
            title: res.data.message,
            icon: 'none'
          })
        }
      },
      fail: () => {
        wx.showToast({
          title: '网络错误，请稍后重试',
          icon: 'none'
        })
      }
    })
  },

  // 处理登出
  handleLogout() {
    wx.removeStorageSync('authToken')
    this.setData({
      isLoggedIn: false,
      userInfo: {},
      loginForm: {
        username: '',
        password: ''
      }
    })
  },

  // 跳转到注册页面
  goToRegister() {
    wx.showToast({
      title: '注册功能开发中',
      icon: 'none'
    })
  },

  // 确保有这样的输入事件处理函数
  handleInput(e) {
    const { field } = e.currentTarget.dataset
    const value = e.detail.value    
    this.setData({
      [`loginForm.${field}`]: value
    })
  },
  
  // 添加编辑表单的输入处理函数
  handleEditInput(e) {
    const { field } = e.currentTarget.dataset
    this.setData({
      [`editForm.${field}`]: e.detail.value
    })
  }
})