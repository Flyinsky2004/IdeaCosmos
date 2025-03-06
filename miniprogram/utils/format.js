/**
 * 格式化日期
 * @param {string|number|Date} date - 日期
 * @param {string} format - 格式化模板，默认为 'YYYY-MM-DD HH:mm:ss'
 * @returns {string} - 格式化后的日期字符串
 */
const formatDate = (date, format = 'YYYY-MM-DD HH:mm:ss') => {
  if (!date) return '';
  
  const d = new Date(date);
  
  const year = d.getFullYear();
  const month = d.getMonth() + 1;
  const day = d.getDate();
  const hour = d.getHours();
  const minute = d.getMinutes();
  const second = d.getSeconds();
  
  const padZero = (num) => {
    return num < 10 ? `0${num}` : num;
  };
  
  return format
    .replace(/YYYY/g, year)
    .replace(/MM/g, padZero(month))
    .replace(/DD/g, padZero(day))
    .replace(/HH/g, padZero(hour))
    .replace(/mm/g, padZero(minute))
    .replace(/ss/g, padZero(second));
};

/**
 * 格式化相对时间
 * @param {string|number|Date} date - 日期
 * @returns {string} - 相对时间字符串
 */
const formatRelativeTime = (date) => {
  if (!date) return '';
  
  const d = new Date(date);
  const now = new Date();
  
  const diff = now - d;
  const seconds = Math.floor(diff / 1000);
  const minutes = Math.floor(seconds / 60);
  const hours = Math.floor(minutes / 60);
  const days = Math.floor(hours / 24);
  
  if (seconds < 60) {
    return '刚刚';
  } else if (minutes < 60) {
    return `${minutes}分钟前`;
  } else if (hours < 24) {
    return `${hours}小时前`;
  } else if (days < 30) {
    return `${days}天前`;
  } else {
    return formatDate(date, 'YYYY-MM-DD');
  }
};

/**
 * 格式化数字（添加千分位）
 * @param {number} num - 数字
 * @returns {string} - 格式化后的数字字符串
 */
const formatNumber = (num) => {
  if (num === undefined || num === null) return '';
  return num.toString().replace(/\B(?=(\d{3})+(?!\d))/g, ',');
};

/**
 * 格式化文件大小
 * @param {number} size - 文件大小（字节）
 * @returns {string} - 格式化后的文件大小
 */
const formatFileSize = (size) => {
  if (size === undefined || size === null) return '';
  
  const units = ['B', 'KB', 'MB', 'GB', 'TB'];
  let index = 0;
  let fileSize = size;
  
  while (fileSize >= 1024 && index < units.length - 1) {
    fileSize /= 1024;
    index++;
  }
  
  return `${fileSize.toFixed(2)} ${units[index]}`;
};

module.exports = {
  formatDate,
  formatRelativeTime,
  formatNumber,
  formatFileSize
}; 