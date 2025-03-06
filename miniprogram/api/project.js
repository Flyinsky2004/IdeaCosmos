// api/project.js
const { get, post, put, delete: del } = require('../utils/request');

/**
 * 获取项目列表
 * @param {Object} params - 查询参数
 * @returns {Promise} - 返回Promise
 */
const getProjects = (params = {}) => {
  return get('/projects', params);
};

/**
 * 获取项目详情
 * @param {string} id - 项目ID
 * @returns {Promise} - 返回Promise
 */
const getProjectDetail = (id) => {
  return get(`/projects/${id}`);
};

/**
 * 创建项目
 * @param {Object} data - 项目数据
 * @returns {Promise} - 返回Promise
 */
const createProject = (data) => {
  return post('/projects', data);
};

/**
 * 更新项目
 * @param {string} id - 项目ID
 * @param {Object} data - 项目数据
 * @returns {Promise} - 返回Promise
 */
const updateProject = (id, data) => {
  return put(`/projects/${id}`, data);
};

/**
 * 删除项目
 * @param {string} id - 项目ID
 * @returns {Promise} - 返回Promise
 */
const deleteProject = (id) => {
  return del(`/projects/${id}`);
};

/**
 * 获取项目章节列表
 * @param {string} projectId - 项目ID
 * @returns {Promise} - 返回Promise
 */
const getChapters = (projectId) => {
  return get(`/projects/${projectId}/chapters`);
};

/**
 * 获取章节详情
 * @param {string} projectId - 项目ID
 * @param {string} chapterId - 章节ID
 * @returns {Promise} - 返回Promise
 */
const getChapterDetail = (projectId, chapterId) => {
  return get(`/projects/${projectId}/chapters/${chapterId}`);
};

module.exports = {
  getProjects,
  getProjectDetail,
  createProject,
  updateProject,
  deleteProject,
  getChapters,
  getChapterDetail
}; 