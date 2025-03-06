/**
 * 检查用户是否已登录
 * @returns {boolean} - 是否已登录
 */
const isLoggedIn = () => {
  const token = wx.getStorageSync('token');
  return !!token;
};

/**
 * 获取用户信息
 * @returns {Object|null} - 用户信息
 */
const getUserInfo = () => {
  return wx.getStorageSync('userInfo') || null;
};

/**
 * 保存用户信息和token
 * @param {Object} data - 用户数据
 * @param {Object} data.userInfo - 用户信息
 * @param {string} data.token - 用户token
 */
const saveUserInfo = (data) => {
  const { userInfo, token } = data;
  wx.setStorageSync('userInfo', userInfo);
  wx.setStorageSync('token', token);
};

/**
 * 清除用户信息和token
 */
const clearUserInfo = () => {
  wx.removeStorageSync('userInfo');
  wx.removeStorageSync('token');
};

/**
 * 登出
 */
const logout = () => {
  clearUserInfo();
  // 跳转到登录页
  wx.reLaunch({
    url: '/pages/profile/index'
  });
};

module.exports = {
  isLoggedIn,
  getUserInfo,
  saveUserInfo,
  clearUserInfo,
  logout
}; 