// 小程序API接口地址常量
const BASE_URL = 'https://idea.1024110.xyz';
const API_BASE_URL = BASE_URL + '/api';

Page({
  data: {
    videos: [],
    currentVideoIndex: 0,
    currentVideo: null, // 当前播放的视频数据
    isLoading: true,
    isLoadingMore: false,
    noMoreVideos: false,
    page: 1,
    pageSize: 10,
    showCommentModal: false,
    selectedVideo: null,
    newComment: '',
    videoContext: null
  },

  onLoad: function () {
    this.fetchVideos();
  },

  onShow: function () {
    // 恢复当前页面视频播放
    if (this.data.videoContext) {
      this.playCurrentVideo();
    }
  },

  onHide: function () {
    // 离开页面时暂停视频
    this.pauseCurrentVideo();
  },

  onUnload: function () {
    // 离开页面时暂停视频
    this.pauseCurrentVideo();
  },

  // 上拉触底加载更多
  onReachBottom: function () {
    if (!this.data.isLoadingMore && !this.data.noMoreVideos) {
      this.fetchVideos(true);
    }
  },

  // 视频滑动切换事件
  onSwiperChange: function (e) {
    const currentIndex = e.detail.current;
    
    // 防止重复触发
    if (this.data.currentVideoIndex === currentIndex) {
      return;
    }
    
    // 获取当前视频对象
    const currentVideo = this.data.videos[currentIndex];
    
    if (!currentVideo) {
      console.error('当前索引没有对应的视频数据:', currentIndex);
      return;
    }
    
    // 更新当前索引和当前视频对象
    this.setData({
      currentVideoIndex: currentIndex,
      currentVideo: currentVideo
    });
    
    // 延迟播放当前视频，避免在过渡动画期间播放导致卡顿
    setTimeout(() => {
      this.playCurrentVideo();
    }, 100);

    // 检查是否需要加载更多视频
    if (currentIndex >= this.data.videos.length - 3 && !this.data.isLoadingMore && !this.data.noMoreVideos) {
      this.fetchVideos(true);
    }
  },

  // 获取视频数据
  fetchVideos: async function (loadMore = false) {
    if (loadMore) {
      if (this.data.noMoreVideos) return;
      this.setData({
        isLoadingMore: true,
        page: this.data.page + 1
      });
    } else {
      this.setData({
        isLoading: true,
        page: 1
      });
    }

    try {
      // 获取token
      const token = wx.getStorageSync('authToken') || '';
      
      // 准备请求参数
      const params = {
        page: this.data.page,
        limit: this.data.pageSize
      };
      
      // 发起请求
      const result = await this.request(`/video/getVideosChapterVersion`, 'GET', params);
      
      const data = result.data || [];
      
      if (loadMore) {
        if (data.length === 0) {
          this.setData({
            noMoreVideos: true,
            isLoadingMore: false
          });
        } else {
          const processedData = this.processVideosData(data);
          const updatedVideos = [...this.data.videos, ...processedData];
          this.setData({
            videos: updatedVideos,
            isLoadingMore: false
          });
        }
      } else {
        const processedData = this.processVideosData(data);
        this.setData({
          videos: processedData,
          isLoading: false,
          currentVideo: processedData.length > 0 ? processedData[0] : null // 设置第一个视频为当前视频
        }, () => {
          // 加载完成后初始化视频上下文并播放第一个视频
          if (processedData.length > 0) {
            this.initVideoContext();
            this.playCurrentVideo();
          }
        });
      }
    } catch (error) {
      console.error('获取视频失败:', error);
      wx.showToast({
        title: error.message || '获取视频失败，请稍后再试',
        icon: 'none',
        duration: 2000
      });
      
      this.handleFetchError(loadMore);
    }
  },
  
  // 处理获取视频数据失败的情况
  handleFetchError: function(loadMore) {
    this.setData({
      isLoading: false,
      isLoadingMore: false
    });
    
    if (!loadMore && !this.data.videos.length) {
      // 首次加载失败时使用模拟数据
      const mockData = this.generateMockVideos(10);
      this.setData({
        videos: mockData,
        currentVideo: mockData[0], // 设置第一个视频为当前视频
        isLoading: false
      }, () => {
        this.initVideoContext();
        this.playCurrentVideo();
      });
    }
  },

  // 处理视频数据 - 根据后端实际返回的数据结构处理
  processVideosData: function (data) {
    if (!Array.isArray(data)) {
      console.error('无效的视频数据格式:', data);
      return [];
    }
    
    return data.map((item, index) => {
      // 如果数据项为空，返回一个有效的默认值
      if (!item) {
        return this.createDefaultVideo(index);
      }
      
      // 确保用户对象存在基本属性
      const author = item.user || {
        ID: 0,
        username: '未知用户',
        avatar: '/assets/images/avatar.png'
      };
      
      // 确保ID唯一，使用时间戳和随机数
      const uniqueId = item.ID ? 
        `video-${item.ID}` : 
        `video-${Date.now()}-${Math.floor(Math.random() * 10000)}-${index}`;
      
      // 确保视频URL格式正确
      let videoUrl = '';
      if (item.video_path) {
        // 检查是否是完整URL
        if (item.video_path.startsWith('http')) {
          videoUrl = item.video_path;
        } else {
          videoUrl = `${API_BASE_URL}/videos/${item.video_path}`;
        }
      } else {
        // 使用默认视频
        videoUrl = 'https://developer.obs.cn-north-4.myhuaweicloud.com:443/mng/202205/3.mp4';
      }
      
      return {
        id: uniqueId,
        title: item.chapter_id ? `第${item.chapter_id}章视频` : `视频${index + 1}`,
        description: item.content.substring(0, 100) || '无描述信息',
        videoUrl: videoUrl,
        createdAt: item.CreatedAt || this.formatTime(new Date()),
        updatedAt: item.UpdatedAt || this.formatTime(new Date()),
        chapterId: item.chapter_id || index + 1,
        score: item.score || Math.floor(Math.random() * 50) + 50,
        status: item.status || 1,
        author: {
          id: author.ID || index + 100,
          username: author.username || `用户${index + 1}`,
          avatar: author.avatar || '/assets/images/avatar.png'
        },
        // 本地状态 - 离线版点赞收藏
        isLiked: false,
        isFavorited: false,
        likes: Math.floor(Math.random() * 100) + 1, // 模拟点赞数
        favorites: Math.floor(Math.random() * 50) + 1, // 模拟收藏数
        comments: [] // 初始化空评论数组
      };
    });
  },
  
  // 创建默认视频对象
  createDefaultVideo: function(index) {
    return {
      id: `default-video-${Date.now()}-${index}`,
      title: `视频${index + 1}`,
      description: '默认视频内容',
      videoUrl: 'https://developer.obs.cn-north-4.myhuaweicloud.com:443/mng/202205/3.mp4',
      createdAt: this.formatTime(new Date()),
      updatedAt: this.formatTime(new Date()),
      chapterId: index + 1,
      score: Math.floor(Math.random() * 50) + 50,
      status: 1,
      author: {
        id: index + 100,
        username: `创作者${index + 1}`,
        avatar: '/assets/images/avatar.png'
      },
      isLiked: false,
      isFavorited: false,
      likes: Math.floor(Math.random() * 100) + 1,
      favorites: Math.floor(Math.random() * 50) + 1,
      comments: []
    };
  },

  // 生成模拟视频数据
  generateMockVideos: function (count) {
    const videos = [];
    const baseIndex = this.data.videos.length;
    
    for (let i = 0; i < count; i++) {
      const index = baseIndex + i + 1;
      // 使用时间戳和随机数确保ID唯一
      const uniqueId = 'video-' + Date.now() + '-' + Math.floor(Math.random() * 10000) + '-' + index;
      videos.push({
        id: uniqueId,
        title: `第${index}章视频`,
        description: `这是第${index}章的视频内容描述，详细介绍了本章节的主要内容和知识点。`,
        videoUrl: 'https://developer.obs.cn-north-4.myhuaweicloud.com:443/mng/202205/3.mp4',
        createdAt: this.formatTime(new Date()),
        updatedAt: this.formatTime(new Date()),
        chapterId: index,
        score: Math.floor(Math.random() * 50) + 50,
        status: 1,
        author: {
          id: 1,
          username: '创作者' + index,
          avatar: '/assets/images/avatar.png'
        },
        isLiked: false,
        isFavorited: false,
        likes: Math.floor(Math.random() * 100) + 1,
        favorites: Math.floor(Math.random() * 50) + 1,
        comments: this.generateMockComments(uniqueId)
      });
    }
    
    return videos;
  },

  // 生成模拟评论数据
  generateMockComments: function (videoId) {
    const commentCount = Math.floor(Math.random() * 5) + 1;
    const comments = [];
    
    for (let i = 0; i < commentCount; i++) {
      // 使用时间戳和随机数确保ID唯一
      const uniqueId = Date.now() + '-' + Math.floor(Math.random() * 10000) + '-' + i;
      comments.push({
        id: uniqueId,
        content: `这是第${i+1}条评论，视频内容很精彩！学到了很多知识。`,
        createdAt: this.formatTime(new Date(Date.now() - Math.floor(Math.random() * 1000000000))),
        likes: Math.floor(Math.random() * 20),
        isLiked: false,
        author: {
          id: i + 100,
          username: `创剧星球`,
          avatar: API_BASE_URL + 'uploads/default-avatar.png'
        }
      });
    }
    
    return comments;
  },

  // 初始化视频上下文
  initVideoContext: function () {
    this.setData({
      videoContext: wx.createVideoContext(`video-${this.data.currentVideoIndex}`)
    });
  },

  // 播放当前视频
  playCurrentVideo: function () {
    // 确保videos数组有数据
    if (!this.data.videos || this.data.videos.length === 0) {
      return;
    }

    // 先暂停所有视频
    this.pauseAllVideos();
    
    // 检查当前索引是否有效
    const currentIndex = this.data.currentVideoIndex;
    if (currentIndex < 0 || currentIndex >= this.data.videos.length) {
      return;
    }

    try {
      // 创建当前视频的上下文并播放
      const videoContext = wx.createVideoContext(`video-${currentIndex}`);
      this.setData({ videoContext });
      
      // 重置视频状态以提高播放可靠性
      videoContext.seek(0);
      
      // 添加延时，确保视频上下文已经准备好
      setTimeout(() => {
        videoContext.play();
        // 设置播放状态标记
        this.videoPlaying = true;
        
        console.log(`正在播放视频：${currentIndex}`);
      }, 100);
    } catch (error) {
      console.error('视频播放失败:', error);
    }
  },

  // 暂停当前视频
  pauseCurrentVideo: function () {
    if (this.data.videoContext) {
      this.data.videoContext.pause();
    }
  },

  // 暂停所有视频
  pauseAllVideos: function () {
    const { videos } = this.data;
    if (!videos || videos.length === 0) {
      return;
    }
    
    videos.forEach((_, index) => {
      try {
        const ctx = wx.createVideoContext(`video-${index}`);
        ctx.pause();
      } catch (error) {
        console.error(`暂停视频${index}失败:`, error);
      }
    });
    
    // 重置播放状态标记
    this.videoPlaying = false;
  },

  // 切换播放/暂停
  togglePlay: function (e) {
    const { index } = e.currentTarget.dataset;
    
    // 暂停当前视频
    if (this.videoPlaying) {
      this.pauseCurrentVideo();
      this.videoPlaying = false;
    } else {
      // 播放当前视频
      const videoContext = wx.createVideoContext(`video-${index}`);
      videoContext.play();
      this.videoPlaying = true;
    }
  },

  // 点赞操作
  toggleLike: function (e) {
    const { index } = e.currentTarget.dataset;
    const videos = this.data.videos;
    const video = videos[index];
    
    // 判断是否已登录
    const userInfo = wx.getStorageSync('userInfo');
    if (!userInfo) {
      wx.showToast({
        title: '请先登录',
        icon: 'none'
      });
      return;
    }
    
    // 更新点赞状态
    video.isLiked = !video.isLiked;
    video.likes += video.isLiked ? 1 : -1;
    
    // 更新数据，同时更新currentVideo
    this.setData({ 
      videos,
      currentVideo: video
    });
    
    wx.showToast({
      title: video.isLiked ? '点赞成功' : '已取消点赞',
      icon: 'success'
    });
  },

  // 收藏操作
  toggleFavorite: function (e) {
    const { index } = e.currentTarget.dataset;
    const videos = this.data.videos;
    const video = videos[index];
    
    // 判断是否已登录
    const userInfo = wx.getStorageSync('userInfo');
    if (!userInfo) {
      wx.showToast({
        title: '请先登录',
        icon: 'none'
      });
      return;
    }
    
    // 更新收藏状态
    video.isFavorited = !video.isFavorited;
    video.favorites += video.isFavorited ? 1 : -1;
    
    // 更新数据，同时更新currentVideo
    this.setData({ 
      videos,
      currentVideo: video
    });
    
    wx.showToast({
      title: video.isFavorited ? '收藏成功' : '已取消收藏',
      icon: 'success'
    });
  },

  // 关注作者
  followAuthor: function (e) {
    const { id } = e.currentTarget.dataset;
    
    // 判断是否已登录
    const userInfo = wx.getStorageSync('userInfo');
    if (!userInfo) {
      wx.showToast({
        title: '请先登录',
        icon: 'none'
      });
      return;
    }
    
    // 这里可以添加实际的关注逻辑，例如调用API进行关注操作
    // 示例中仅显示成功提示
    
    wx.showToast({
      title: '关注成功',
      icon: 'success'
    });
    
    // 可以更新UI状态，如显示已关注标记等
    // 这里可添加后续的UI更新代码
  },

  // 打开评论模态框
  openCommentsModal: function (e) {
    const { index } = e.currentTarget.dataset;
    const selectedVideo = this.data.videos[index];
    
    // 暂停当前视频
    this.pauseCurrentVideo();
    
    // 如果还没有加载评论，从服务器获取或生成模拟评论
    if (!selectedVideo.comments || selectedVideo.comments.length === 0) {
      // 尝试从服务器获取评论
      this.fetchVideoComments(selectedVideo.id)
        .then(comments => {
          if (comments && comments.length > 0) {
            selectedVideo.comments = comments;
          } else {
            // 如果没有评论或获取失败，生成模拟评论
            selectedVideo.comments = this.generateMockComments(selectedVideo.id);
          }
          
          // 如果打开的是当前播放视频的评论，同时更新currentVideo
          const updatedState = {
            selectedVideo,
            showCommentModal: true,
            newComment: ''
          };
          
          if (index === this.data.currentVideoIndex) {
            updatedState.currentVideo = selectedVideo;
          }
          
          this.setData(updatedState);
        })
        .catch(() => {
          // 获取失败时使用模拟数据
          selectedVideo.comments = this.generateMockComments(selectedVideo.id);
          
          // 如果打开的是当前播放视频的评论，同时更新currentVideo
          const updatedState = {
            selectedVideo,
            showCommentModal: true,
            newComment: ''
          };
          
          if (index === this.data.currentVideoIndex) {
            updatedState.currentVideo = selectedVideo;
          }
          
          this.setData(updatedState);
        });
    } else {
      // 已有评论数据，直接显示
      // 如果打开的是当前播放视频的评论，同时更新currentVideo
      const updatedState = {
        selectedVideo,
        showCommentModal: true,
        newComment: ''
      };
      
      if (index === this.data.currentVideoIndex) {
        updatedState.currentVideo = selectedVideo;
      }
      
      this.setData(updatedState);
    }
  },

  // 获取视频评论
  fetchVideoComments: function(videoId) {
    return new Promise((resolve) => {
      // 这里模拟异步获取评论
      // 实际项目中应当调用后端API获取评论
      setTimeout(() => {
        // 返回模拟数据
        resolve(this.generateMockComments(videoId));
      }, 300);
    });
  },

  // 关闭评论模态框
  closeCommentsModal: function () {
    this.setData({
      showCommentModal: false
    });
    
    // 关闭评论后恢复视频播放
    this.playCurrentVideo();
  },

  // 评论输入变化
  onCommentInput: function (e) {
    this.setData({
      newComment: e.detail.value
    });
  },

  // 添加评论
  addComment: function () {
    // 判断是否已登录
    const userInfo = wx.getStorageSync('userInfo');
    if (!userInfo) {
      wx.showToast({
        title: '请先登录',
        icon: 'none'
      });
      return;
    }
    
    const { newComment, selectedVideo } = this.data;
    
    if (!newComment.trim()) {
      return;
    }
    
    // 创建新评论 - 使用时间戳确保ID唯一
    const commentId = `new-comment-${Date.now()}-${Math.floor(Math.random() * 10000)}`;
    const comment = {
      id: commentId,
      content: newComment.trim(),
      author: {
        id: userInfo.id || 999,
        username: userInfo.nickName || '用户',
        avatar: userInfo.avatarUrl || '/assets/images/avatar.png'
      },
      createdAt: this.formatTime(new Date()),
      likes: 0,
      isLiked: false
    };
    
    // 添加到评论列表
    selectedVideo.comments.unshift(comment);
    
    // 更新视频数据中的评论
    const videoIndex = this.data.videos.findIndex(v => v.id === selectedVideo.id);
    if (videoIndex !== -1) {
      const videos = this.data.videos;
      videos[videoIndex] = selectedVideo;
      
      // 更新数据，如果当前视频是评论的视频，则也更新currentVideo
      const updatedState = {
        videos,
        selectedVideo,
        newComment: ''
      };
      
      if (videoIndex === this.data.currentVideoIndex) {
        updatedState.currentVideo = selectedVideo;
      }
      
      this.setData(updatedState);
      
      wx.showToast({
        title: '评论已发布',
        icon: 'success'
      });
    }
  },

  // 给评论点赞
  toggleCommentLike: function (e) {
    const { id } = e.currentTarget.dataset;
    const { selectedVideo } = this.data;
    
    // 判断是否已登录
    const userInfo = wx.getStorageSync('userInfo');
    if (!userInfo) {
      wx.showToast({
        title: '请先登录',
        icon: 'none'
      });
      return;
    }
    
    // 找到并更新对应的评论
    const commentIndex = selectedVideo.comments.findIndex(c => c.id === id);
    if (commentIndex !== -1) {
      const comment = selectedVideo.comments[commentIndex];
      comment.isLiked = !comment.isLiked;
      comment.likes += comment.isLiked ? 1 : -1;
      
      // 更新视频数据中的评论
      const videoIndex = this.data.videos.findIndex(v => v.id === selectedVideo.id);
      if (videoIndex !== -1) {
        const videos = this.data.videos;
        videos[videoIndex] = selectedVideo;
        
        // 更新数据，如果当前视频是评论的视频，则也更新currentVideo
        const updatedState = {
          videos,
          selectedVideo
        };
        
        if (videoIndex === this.data.currentVideoIndex) {
          updatedState.currentVideo = selectedVideo;
        }
        
        this.setData(updatedState);
      }
    }
  },

  // 分享视频
  shareVideo: function (e) {
    const { index } = e.currentTarget.dataset;
    const video = this.data.videos[index];
    
    wx.showShareMenu({
      withShareTicket: true,
      menus: ['shareAppMessage', 'shareTimeline']
    });
    
    wx.showToast({
      title: '分享卡片已生成',
      icon: 'success'
    });
  },

  // 分享到朋友圈
  onShareTimeline: function () {
    const currentVideo = this.data.videos[this.data.currentVideoIndex];
    return {
      title: currentVideo.title,
      imageUrl: '/assets/images/share-cover.png'
    };
  },

  // 分享给好友
  onShareAppMessage: function () {
    const currentVideo = this.data.videos[this.data.currentVideoIndex];
    return {
      title: currentVideo.title,
      path: `/pages/video-feed/index?videoId=${currentVideo.id}`,
      imageUrl: '/assets/images/share-cover.png'
    };
  },

  // 返回上一页
  goBack: function () {
    wx.navigateBack();
  },

  // 防止评论模态框穿透滚动
  preventTouchMove: function () {
    return false;
  },

  // 格式化日期时间
  formatTime: function (date) {
    const year = date.getFullYear();
    const month = date.getMonth() + 1;
    const day = date.getDate();
    const hour = date.getHours();
    const minute = date.getMinutes();
    
    return `${year}-${this.formatNumber(month)}-${this.formatNumber(day)} ${this.formatNumber(hour)}:${this.formatNumber(minute)}`;
  },

  formatNumber: function (n) {
    n = n.toString();
    return n[1] ? n : `0${n}`;
  },
  
  // 封装请求函数
  request: function(url, method, data) {
    return new Promise((resolve, reject) => {
      // 获取token - 使用authToken而不是token，且不添加Bearer前缀
      const token = wx.getStorageSync('authToken') || '';
      
      wx.request({
        url: `${API_BASE_URL}${url}`,
        method,
        data,
        header: {
          'Content-Type': 'application/json',
          'Authorization': token, // 直接使用token，无前缀
          'Cache-Control': 'no-cache',
          'Pragma': 'no-cache',
          'If-Modified-Since': '0'
        },
        success: (res) => {
          // 请求成功
          if (res.statusCode >= 200 && res.statusCode < 300) {
            // 处理Web端的响应格式，Web端可能直接返回数据或包装在data字段中
            if (res.data.success !== undefined) {
              // 如果是Web端标准格式
              if (res.data.success) {
                resolve({
                  data: res.data.data,
                  message: res.data.message
                });
              } else {
                reject(new Error(res.data.message || '请求失败'));
              }
            } else {
              // 直接返回数据
              resolve({
                data: res.data,
                message: 'success'
              });
            }
          } 
          // 未授权
          else if (res.statusCode === 401) {
            // 清除token
            wx.removeStorageSync('authToken');
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
          wx.showToast({
            title: '网络错误，请稍后再试',
            icon: 'none',
            duration: 2000
          });
          
          reject(new Error('网络错误，请稍后再试'));
        }
      });
    });
  }
}); 