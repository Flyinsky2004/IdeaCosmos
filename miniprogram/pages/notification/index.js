// pages/notification/index.js
const app = getApp();

Page({

  /**
   * 页面的初始数据
   */
  data: {
    notifications: [],
    loading: true,
    pageNum: 1,
    pageSize: 10,
    total: 0,
    selectedTab: 'all',
    notificationSettings: null,
    settingsLoading: false,
    showSettings: false,
    showBackToTop: false,
    tabs: [
      { name: '全部', type: 'all' },
      { name: '系统通知', type: '1' },
      { name: '点赞通知', type: '2' },
      { name: '评论通知', type: '3' },
      { name: '关注通知', type: '4' },
      { name: '协作邀请', type: '5' },
      { name: '内容更新', type: '6' }
    ]
  },

  /**
   * 生命周期函数--监听页面加载
   */
  onLoad(options) {
    this.fetchNotifications();
    this.fetchNotificationSettings();
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
    this.setData({
      pageNum: 1
    }, () => {
      this.fetchNotifications(() => {
        wx.stopPullDownRefresh();
      });
    });
  },

  /**
   * 页面上拉触底事件的处理函数
   */
  onReachBottom() {
    if (this.data.notifications.length < this.data.total) {
      this.setData({
        pageNum: this.data.pageNum + 1
      }, () => {
        this.fetchNotifications(null, true);
      });
    }
  },

  /**
   * 用户点击右上角分享
   */
  onShareAppMessage() {

  },

  /**
   * 获取通知列表
   */
  fetchNotifications(callback, append = false) {
    this.setData({ loading: true });
    
    wx.request({
      url: `${app.globalData.baseUrl}/api/notifications`,
      method: 'GET',
      header: {
        'Authorization': `Bearer ${wx.getStorageSync('token')}`
      },
      data: {
        pageNum: this.data.pageNum,
        pageSize: this.data.pageSize,
        type: this.data.selectedTab === 'all' ? '' : this.data.selectedTab
      },
      success: (res) => {
        if (res.statusCode === 200) {
          const newNotifications = res.data.notifications || [];
          
          this.setData({
            notifications: append ? [...this.data.notifications, ...newNotifications] : newNotifications,
            total: res.data.total || 0,
            loading: false
          });
        } else {
          wx.showToast({
            title: '获取通知失败',
            icon: 'none'
          });
          this.setData({ loading: false });
        }
      },
      fail: () => {
        wx.showToast({
          title: '获取通知失败',
          icon: 'none'
        });
        this.setData({ loading: false });
      },
      complete: () => {
        if (callback) callback();
      }
    });
  },

  /**
   * 获取通知设置
   */
  fetchNotificationSettings() {
    this.setData({ settingsLoading: true });
    
    wx.request({
      url: `${app.globalData.baseUrl}/api/notifications/settings`,
      method: 'GET',
      header: {
        'Authorization': `Bearer ${wx.getStorageSync('token')}`
      },
      success: (res) => {
        if (res.statusCode === 200) {
          this.setData({
            notificationSettings: res.data,
            settingsLoading: false
          });
        } else {
          wx.showToast({
            title: '获取通知设置失败',
            icon: 'none'
          });
          this.setData({ settingsLoading: false });
        }
      },
      fail: () => {
        wx.showToast({
          title: '获取通知设置失败',
          icon: 'none'
        });
        this.setData({ settingsLoading: false });
      }
    });
  },

  /**
   * 更新通知设置
   */
  updateNotificationSettings() {
    wx.request({
      url: `${app.globalData.baseUrl}/api/notifications/settings`,
      method: 'POST',
      header: {
        'Authorization': `Bearer ${wx.getStorageSync('token')}`
      },
      data: this.data.notificationSettings,
      success: (res) => {
        if (res.statusCode === 200) {
          wx.showToast({
            title: '设置已更新',
            icon: 'success'
          });
        } else {
          wx.showToast({
            title: '更新设置失败',
            icon: 'none'
          });
        }
      },
      fail: () => {
        wx.showToast({
          title: '更新设置失败',
          icon: 'none'
        });
      }
    });
  },

  /**
   * 标记通知为已读
   */
  markAsRead(e) {
    const id = e.currentTarget.dataset.id;
    
    wx.request({
      url: `${app.globalData.baseUrl}/api/notifications/${id}/read`,
      method: 'POST',
      header: {
        'Authorization': `Bearer ${wx.getStorageSync('token')}`
      },
      success: (res) => {
        if (res.statusCode === 200) {
          const notifications = this.data.notifications.map(n => {
            if (n.ID === id) {
              n.IsRead = true;
            }
            return n;
          });
          
          this.setData({ notifications });
          
          wx.showToast({
            title: '已标记为已读',
            icon: 'success'
          });
        } else {
          wx.showToast({
            title: '标记已读失败',
            icon: 'none'
          });
        }
      },
      fail: () => {
        wx.showToast({
          title: '标记已读失败',
          icon: 'none'
        });
      }
    });
  },

  /**
   * 标记所有通知为已读
   */
  markAllAsRead() {
    wx.request({
      url: `${app.globalData.baseUrl}/api/notifications/read-all`,
      method: 'POST',
      header: {
        'Authorization': `Bearer ${wx.getStorageSync('token')}`
      },
      success: (res) => {
        if (res.statusCode === 200) {
          const notifications = this.data.notifications.map(n => {
            n.IsRead = true;
            return n;
          });
          
          this.setData({ notifications });
          
          wx.showToast({
            title: '已全部标记为已读',
            icon: 'success'
          });
        } else {
          wx.showToast({
            title: '标记全部已读失败',
            icon: 'none'
          });
        }
      },
      fail: () => {
        wx.showToast({
          title: '标记全部已读失败',
          icon: 'none'
        });
      }
    });
  },

  /**
   * 删除通知
   */
  deleteNotification(e) {
    const id = e.currentTarget.dataset.id;
    
    wx.request({
      url: `${app.globalData.baseUrl}/api/notifications/${id}`,
      method: 'DELETE',
      header: {
        'Authorization': `Bearer ${wx.getStorageSync('token')}`
      },
      success: (res) => {
        if (res.statusCode === 200) {
          const notifications = this.data.notifications.filter(n => n.ID !== id);
          
          this.setData({ notifications });
          
          wx.showToast({
            title: '删除成功',
            icon: 'success'
          });
        } else {
          wx.showToast({
            title: '删除失败',
            icon: 'none'
          });
        }
      },
      fail: () => {
        wx.showToast({
          title: '删除失败',
          icon: 'none'
        });
      }
    });
  },

  /**
   * 删除所有通知
   */
  deleteAllNotifications() {
    wx.showModal({
      title: '确认删除',
      content: '确定要删除所有通知吗？',
      success: (res) => {
        if (res.confirm) {
          wx.request({
            url: `${app.globalData.baseUrl}/api/notifications`,
            method: 'DELETE',
            header: {
              'Authorization': `Bearer ${wx.getStorageSync('token')}`
            },
            success: (res) => {
              if (res.statusCode === 200) {
                this.setData({
                  notifications: [],
                  total: 0
                });
                
                wx.showToast({
                  title: '已清空所有通知',
                  icon: 'success'
                });
              } else {
                wx.showToast({
                  title: '清空通知失败',
                  icon: 'none'
                });
              }
            },
            fail: () => {
              wx.showToast({
                title: '清空通知失败',
                icon: 'none'
              });
            }
          });
        }
      }
    });
  },

  /**
   * 切换标签
   */
  changeTab(e) {
    const type = e.currentTarget.dataset.type;
    
    this.setData({
      selectedTab: type,
      pageNum: 1
    }, () => {
      this.fetchNotifications();
    });
  },

  /**
   * 切换设置面板
   */
  toggleSettings() {
    this.setData({
      showSettings: !this.data.showSettings
    });
  },

  /**
   * 切换通知设置
   */
  toggleNotificationSetting(e) {
    const setting = e.currentTarget.dataset.setting;
    const value = e.detail.value;
    
    const notificationSettings = { ...this.data.notificationSettings };
    notificationSettings[setting] = value;
    
    this.setData({ notificationSettings });
  },

  /**
   * 滚动到顶部
   */
  scrollToTop() {
    wx.pageScrollTo({
      scrollTop: 0,
      duration: 300
    });
  },

  /**
   * 监听页面滚动
   */
  onPageScroll(e) {
    this.setData({
      showBackToTop: e.scrollTop > 300
    });
  },

  /**
   * 格式化时间
   */
  formatTime(timestamp) {
    const date = new Date(timestamp);
    return date.toLocaleString('zh-CN', {
      year: 'numeric',
      month: '2-digit',
      day: '2-digit',
      hour: '2-digit',
      minute: '2-digit'
    });
  },

  /**
   * 获取通知类型名称
   */
  getNotificationTypeName(type) {
    const typeMap = {
      1: '系统通知',
      2: '点赞通知',
      3: '评论通知',
      4: '关注通知',
      5: '协作邀请',
      6: '内容更新'
    };
    return typeMap[type] || '未知类型';
  },

  /**
   * 获取通知图标
   */
  getNotificationIcon(type) {
    const iconMap = {
      1: 'notification', // 系统通知
      2: 'like', // 点赞通知
      3: 'comment', // 评论通知
      4: 'user-follow', // 关注通知
      5: 'collaborate', // 协作邀请
      6: 'document-update' // 内容更新
    };
    return iconMap[type] || 'notification';
  }
})