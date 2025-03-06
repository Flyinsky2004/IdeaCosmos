// utils/request.js

const BASE_URL = 'https://idea.1024110.xyz/api';

/**
 * 封装微信请求
 * @param {Object} options - 请求配置
 * @param {string} options.url - 请求地址
 * @param {string} options.method - 请求方法
 * @param {Object} options.data - 请求数据
 * @param {boolean} options.loading - 是否显示加载中
 * @returns {Promise} - 返回Promise
 */
const request = (options) => {
  const { url, method = 'GET', data = {}, loading = true } = options;

  // 显示加载中
  if (loading) {
    wx.showLoading({
      title: '加载中...',
      mask: true
    });
  }

  // 获取token
  const token = wx.getStorageSync('token') || '';

  return new Promise((resolve, reject) => {
    wx.request({
      url: `${BASE_URL}${url}`,
      method,
      data,
      header: {
        'Content-Type': 'application/json',
        'Authorization': token ? `Bearer ${token}` : ''
      },
      success: (res) => {
        if (loading) {
          wx.hideLoading();
        }

        // 请求成功
        if (res.statusCode >= 200 && res.statusCode < 300) {
          resolve(res.data);
        } 
        // 未授权
        else if (res.statusCode === 401) {
          // 清除token
          wx.removeStorageSync('token');
          wx.removeStorageSync('userInfo');
          
          // 跳转到登录页
          wx.navigateTo({
            url: '/pages/profile/index'
          });
          
          reject(new Error('未授权，请重新登录'));
        } 
        // 其他错误
        else {
          wx.showToast({
            title: res.data.message || '请求失败',
            icon: 'none',
            duration: 2000
          });
          reject(new Error(res.data.message || '请求失败'));
        }
      },
      fail: (err) => {
        if (loading) {
          wx.hideLoading();
        }
        
        wx.showToast({
          title: '网络错误，请稍后再试',
          icon: 'none',
          duration: 2000
        });
        
        reject(new Error('网络错误，请稍后再试'));
      }
    });
  });
};

// 导出请求方法
module.exports = {
  get: (url, data, loading) => request({ url, method: 'GET', data, loading }),
  post: (url, data, loading) => request({ url, method: 'POST', data, loading }),
  put: (url, data, loading) => request({ url, method: 'PUT', data, loading }),
  delete: (url, data, loading) => request({ url, method: 'DELETE', data, loading })
}; 