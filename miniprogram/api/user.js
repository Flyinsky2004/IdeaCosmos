// api/user.js
const { get, post, put } = require('../utils/request');

/**
 * 用户登录
 * @param {Object} data - 登录数据
 * @param {string} data.username - 用户名
 * @param {string} data.password - 密码
 * @returns {Promise} - 返回Promise
 */
const login = (data) => {
  return post('/user/login', data);
};

/**
 * 用户注册
 * @param {Object} data - 注册数据
 * @param {string} data.username - 用户名
 * @param {string} data.password - 密码
 * @param {string} data.email - 邮箱
 * @returns {Promise} - 返回Promise
 */
const register = (data) => {
  return post('/user/register', data);
};

/**
 * 获取用户信息
 * @returns {Promise} - 返回Promise
 */
const getUserInfo = () => {
  return get('/user/info');
};

/**
 * 更新用户信息
 * @param {Object} data - 用户信息
 * @returns {Promise} - 返回Promise
 */
const updateUserInfo = (data) => {
  return put('/user/info', data);
};

/**
 * 更新用户头像
 * @param {string} filePath - 头像文件路径
 * @returns {Promise} - 返回Promise
 */
const updateAvatar = (filePath) => {
  return new Promise((resolve, reject) => {
    wx.uploadFile({
      url: 'https://idea.1024110.xyz/api/user/avatar',
      filePath,
      name: 'avatar',
      header: {
        'Authorization': `Bearer ${wx.getStorageSync('token') || ''}`
      },
      success: (res) => {
        try {
          const data = JSON.parse(res.data);
          resolve(data);
        } catch (error) {
          reject(new Error('上传头像失败'));
        }
      },
      fail: () => {
        reject(new Error('上传头像失败'));
      }
    });
  });
};

module.exports = {
  login,
  register,
  getUserInfo,
  updateUserInfo,
  updateAvatar
}; 