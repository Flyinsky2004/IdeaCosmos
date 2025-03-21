const app = getApp();
const towxml = require("../../utils/towxml/index");

Page({
  data: {
    loading: true,
    chapter: null,
    project: null,
    parsedContent: "",
    audioUrl: "",
    comments: [],
    commentContent: "",
    isLoggedIn: false,

    // 字体大小选项
    fontSizes: [
      { label: "小", value: "small" },
      { label: "中", value: "medium" },
      { label: "大", value: "large" },
    ],
    currentFontSize: "medium",

    // 评论类型
    commentTypes: [
      { label: "全部评论", value: "all" },
      { label: "读者评论", value: "reader" },
      { label: "作者评论", value: "author" },
    ],
    commentType: "all",

    // 情绪映射
    emotionMap: {
      喜悦: {
        icon: "😊",
        description: "充满快乐和满足",
      },
      感动: {
        icon: "🥹",
        description: "内心被深深触动",
      },
      惊喜: {
        icon: "🤩",
        description: "意外的惊喜",
      },
      期待: {
        icon: "🤔",
        description: "对后续充满期待",
      },
      伤感: {
        icon: "😢",
        description: "略带忧伤的感动",
      },
      愤怒: {
        icon: "😠",
        description: "对情节感到愤慨",
      },
      恐惧: {
        icon: "😱",
        description: "感到害怕或紧张",
      },
      平静: {
        icon: "😐",
        description: "内心平和安宁",
      },
    },
    emotions: [], // 将在 onLoad 中初始化
    selectedEmotion: "",
    userFeeling: null,
    isPlaying: false,
    currentTime: 0,
    audioDuration: 0,
    currentTimeText: '00:00',
    durationText: '00:00',
    playbackRates: [
      { label: '0.8x', value: 0.8 },
      { label: '1.0x', value: 1.0 },
      { label: '1.2x', value: 1.2 },
      { label: '1.5x', value: 1.5 }
    ],
    currentRate: 1.0,
    innerAudioContext: null,
    
    // 插画模式相关
    isShowPic: false,
    chapterScenes: []
  },

  onLoad(options) {
    this.chapterId = options.id;

    // 初始化情绪数组
    const emotions = Object.entries(this.data.emotionMap).map(
      ([name, data]) => ({
        name,
        icon: data.icon,
        description: data.description,
      })
    );

    this.setData({
      emotions,
    });

    // 检查登录状态
    wx.getStorage({
      key: 'authToken',
      success: (res) => {
        this.setData({
          isLoggedIn: !!res.data
        });
      },
      fail: () => {
        this.setData({
          isLoggedIn: false
        });
      }
    });

    this.fetchChapterDetail();
  },

  // 初始化音频上下文
  initAudioContext() {
    try {
      // 如果已存在，先销毁
      if (this.innerAudioContext) {
        this.innerAudioContext.destroy();
      }

      // 创建新的音频上下文
      this.innerAudioContext = wx.createInnerAudioContext();

      // 绑定事件
      this.innerAudioContext.onPlay(() => {
        this.setData({ isPlaying: true });
      });

      this.innerAudioContext.onPause(() => {
        this.setData({ isPlaying: false });
      });

      this.innerAudioContext.onTimeUpdate(() => {
        this.setData({
          currentTime: this.innerAudioContext.currentTime,
          audioDuration: this.innerAudioContext.duration,
          currentTimeText: this.formatTime(this.innerAudioContext.currentTime),
          durationText: this.formatTime(this.innerAudioContext.duration)
        });
      });

      this.innerAudioContext.onEnded(() => {
        this.setData({
          isPlaying: false,
          currentTime: 0,
          currentTimeText: '00:00'
        });
      });

      this.innerAudioContext.onError((res) => {
        console.error('音频播放错误:', res);
        wx.showToast({
          title: '音频加载失败',
          icon: 'none'
        });
      });

      return true;
    } catch (error) {
      console.error('初始化音频上下文失败:', error);
      return false;
    }
  },

  // 获取章节详情
  async fetchChapterDetail() {
    try {
      this.setData({ loading: true });
      
      // 获取存储的认证令牌
      const { data: authToken } = await wx
        .getStorage({ key: "authToken" })
        .catch(() => ({ data: null }));

      // 将 wx.request 包装成 Promise
      const result = await new Promise((resolve, reject) => {
        wx.request({
          url: "https://idea.1024110.xyz/api/public/getChapterDetail",
          method: "GET",
          data: { id: this.chapterId },
          header: authToken ? { Authorization: authToken } : {},
          success: (res) => {
            resolve(res);
          },
          fail: (err) => {
            console.error('请求失败:', err);
            reject(err);
          }
        });
      });
      
      // 确保我们有数据
      if (!result.data) {
        throw new Error('响应数据为空');
      }
            
      // 根据实际API响应结构调整
      // 有些API会将数据包装在data字段中
      const responseData = result.data.data || result.data;
      
      const { chapter, project } = responseData;
      
      if (!chapter || !project) {
        throw new Error('章节或项目数据缺失');
      }

      // 设置页面标题为项目名称
      wx.setNavigationBarTitle({
        title: project.project_name + '-' + chapter.Title
      });

      // 解析 Markdown 内容
      const parsedContent = towxml(
        chapter.current_version.content,
        "markdown",
        {
          theme: "light",
          events: {
            tap: (e) => {
              console.log("tap", e);
            },
          },
        }
      );

      
      // 设置音频 URL
      const audioUrl = chapter.current_version.audio_path
        ? `https://idea.1024110.xyz/api/audio/${chapter.current_version.audio_path}`
        : "";

      if (audioUrl) {
        // 初始化音频上下文并设置URL
        if (this.initAudioContext()) {
          this.innerAudioContext.src = audioUrl;
        }
      }

      this.setData({
        chapter,
        project,
        audioUrl,
        loading: false,
        // 直接使用原始的markdown内容，不使用towxml解析结果
        // parsedContent: parsedContent
      });

      // 获取章节场景（插画模式需要）
      this.fetchChapterScenes();
      
      // 加载评论
      this.fetchComments();
      
      // 加载用户情绪评价
      this.fetchUserFeeling();
    } catch (error) {
      console.error("获取章节详情失败:", error);
      this.setData({ loading: false });
      wx.showToast({
        title: "获取章节详情失败",
        icon: "none",
      });
    }
  },

  // 获取章节场景（用于插画模式）
  async fetchChapterScenes() {
    if (!this.data.chapter || !this.data.chapter.current_version) return;
    
    try {
      // 获取存储的认证令牌
      const { data: authToken } = await wx
        .getStorage({ key: "authToken" })
        .catch(() => ({ data: null }));
      
      wx.request({
        url: "https://idea.1024110.xyz/api/video/getSceneByChapterVersionID",
        method: "GET",
        data: { 
          chapter_verison_id: this.data.chapter.current_version.ID 
        },
        header: authToken ? { Authorization: authToken } : {},
        success: (res) => {
          if (res.data && res.data.data) {
            this.setData({
              chapterScenes: res.data.data
            });
          }
        },
        fail: (err) => {
          console.error('获取章节场景失败:', err);
        }
      });
    } catch (error) {
      console.error("获取章节场景失败:", error);
    }
  },

  // 获取评论列表
  async fetchComments() {
    if (!this.data.chapter?.current_version?.ID) return;

    try {
      const authTokenRes = await wx.getStorage({ key: "authToken" });
      const authToken = authTokenRes.data;

      const res = await new Promise((resolve, reject) => {
        wx.request({
          url: "https://idea.1024110.xyz/api/user/getVersionComments",
          method: "GET",
          data: {
            version_id: this.data.chapter.current_version.ID,
            type: this.data.commentType,
          },
          header: authToken ? { Authorization: authToken } : {},
          success: resolve,
          fail: reject
        });
      });

      if (res.statusCode === 200 && res.data.data) {
        this.setData({
          comments: res.data.data.reader_comments || [],
        });
      }
    } catch (error) {
      console.error("获取评论失败:", error);
      if (error.errMsg && error.errMsg.includes('getStorage:fail')) {
        this.setData({ isLoggedIn: false });
      }
    }
  },

  // 获取用户情绪评价
  async fetchUserFeeling() {
    if (!this.data.chapter?.current_version?.ID || !this.data.isLoggedIn) return;

    try {
      const authTokenRes = await wx.getStorage({ key: "authToken" });
      const authToken = authTokenRes.data;

      if (!authToken) {
        this.setData({ isLoggedIn: false });
        return;
      }

      const res = await new Promise((resolve, reject) => {
        wx.request({
          url: "https://idea.1024110.xyz/api/user/feeling/get",
          method: "GET",
          data: {
            version_id: this.data.chapter.current_version.ID,
          },
          header: { Authorization: authToken },
          success: resolve,
          fail: reject
        });
      });

      if (res.statusCode === 200 && res.data.data != '获取失败') {
        this.setData({
          userFeeling: res.data.data,
        });
      }
    } catch (error) {
      console.error("获取情绪评价失败:", error);
      if (error.errMsg && error.errMsg.includes('getStorage:fail')) {
        this.setData({ isLoggedIn: false });
      }
    }
  },

  // 提交情绪评价
  async submitFeeling(e) {
    if (!this.data.isLoggedIn) {
      wx.showToast({
        title: "请先登录",
        icon: "none",
      });
      return;
    }

    const emotion = e.currentTarget.dataset.emotion;
    this.setData({ selectedEmotion: emotion });

    try {
      const authTokenRes = await wx.getStorage({ key: "authToken" });
      const authToken = authTokenRes.data;

      if (!authToken) {
        this.setData({ isLoggedIn: false });
        wx.showToast({
          title: "请先登录",
          icon: "none",
        });
        return;
      }

      const res = await new Promise((resolve, reject) => {
        wx.request({
          url: "https://idea.1024110.xyz/api/user/feeling/add",
          method: "POST",
          data: {
            version_id: this.data.chapter.current_version.ID,
            feeling: emotion,
          },
          header: {
            Authorization: authToken,
            "Content-Type": "application/json",
          },
          success: resolve,
          fail: reject
        });
      });

      if (res.statusCode === 200) {
        wx.showToast({
          title: "评价成功",
          icon: "success",
        });
        await this.fetchUserFeeling();
      }
    } catch (error) {
      console.error("提交情绪评价失败:", error);
      wx.showToast({
        title: "评价失败",
        icon: "none",
      });
      if (error.errMsg && error.errMsg.includes('getStorage:fail')) {
        this.setData({ isLoggedIn: false });
      }
    }
  },

  // 评论相关方法
  onCommentInput(e) {
    this.setData({
      commentContent: e.detail.value,
    });
  },

  async submitComment() {
    if (!this.data.isLoggedIn) {
      wx.showToast({
        title: "请先登录",
        icon: "none",
      });
      return;
    }

    if (!this.data.commentContent.trim()) {
      wx.showToast({
        title: "请输入评论内容",
        icon: "none",
      });
      return;
    }

    try {
      const authTokenRes = await wx.getStorage({ key: "authToken" });
      const authToken = authTokenRes.data;

      if (!authToken) {
        this.setData({ isLoggedIn: false });
        wx.showToast({
          title: "请先登录",
          icon: "none",
        });
        return;
      }

      const res = await new Promise((resolve, reject) => {
        wx.request({
          url: "https://idea.1024110.xyz/api/user/addVersionComment",
          method: "POST",
          data: {
            version_id: this.data.chapter.current_version.ID,
            content: this.data.commentContent,
            type: "reader",
          },
          header: {
            Authorization: authToken,
            "Content-Type": "application/json",
          },
          success: resolve,
          fail: reject
        });
      });

      if (res.statusCode === 200) {
        wx.showToast({
          title: "评论成功",
          icon: "success",
        });

        this.setData({
          commentContent: "",
        });

        await this.fetchComments();
      }
    } catch (error) {
      console.error("提交评论失败:", error);
      wx.showToast({
        title: "评论失败",
        icon: "none",
      });
      if (error.errMsg && error.errMsg.includes('getStorage:fail')) {
        this.setData({ isLoggedIn: false });
      }
    }
  },

  // 切换评论类型
  changeCommentType(e) {
    const type = e.currentTarget.dataset.type;
    this.setData({ commentType: type });
    this.fetchComments();
  },

  // 字体大小控制
  changeFontSize(e) {
    const size = e.currentTarget.dataset.size;
    this.setData({ currentFontSize: size });
  },

  // 音频控制相关方法
  togglePlayAudio() {
    if (this.data.isPlaying) {
      this.innerAudioContext.pause();
    } else {
      this.innerAudioContext.play();
    }
  },

  rewindAudio() {
    const newTime = Math.max(0, this.data.currentTime - 15);
    this.innerAudioContext.seek(newTime);
  },

  forwardAudio() {
    const newTime = Math.min(this.data.audioDuration, this.data.currentTime + 15);
    this.innerAudioContext.seek(newTime);
  },

  onSliderChange(e) {
    const position = e.detail.value;
    this.innerAudioContext.seek(position);
  },

  onSliderChanging(e) {
    const position = e.detail.value;
    this.setData({
      currentTime: position,
      currentTimeText: this.formatTime(position)
    });
  },

  changePlaybackRate(e) {
    const rate = e.currentTarget.dataset.rate;
    this.innerAudioContext.playbackRate = rate;
    this.setData({ currentRate: rate });
  },

  formatTime(seconds) {
    const minutes = Math.floor(seconds / 60);
    const remainingSeconds = Math.floor(seconds % 60);
    return `${minutes.toString().padStart(2, '0')}:${remainingSeconds.toString().padStart(2, '0')}`;
  },

  // 导航方法
  navigateBack() {
    wx.navigateBack();
  },

  navigateToLogin() {
    wx.navigateTo({
      url: "/pages/login/index",
    });
  },

  // 格式化日期
  formatDate(dateStr) {
    if (!dateStr) return '';
    
    try {
      // 移除时区信息后解析日期
      const cleanDate = dateStr.split('+')[0];
      const d = new Date(cleanDate);
      
      if (isNaN(d.getTime())) {
        return '';
      }

      const year = d.getFullYear();
      const month = String(d.getMonth() + 1).padStart(2, '0');
      const day = String(d.getDate()).padStart(2, '0');
      const hours = String(d.getHours()).padStart(2, '0');
      const minutes = String(d.getMinutes()).padStart(2, '0');
      
      return `${year}年${month}月${day}日 ${hours}:${minutes}`;
    } catch (error) {
      console.error('日期格式化错误:', error);
      return '';
    }
  },

  onShareAppMessage() {
    return {
      title: this.data.chapter?.Title || "精彩章节",
      path: `/pages/chapter-detail/index?id=${this.chapterId}`,
    };
  },

  // 页面卸载时清理音频上下文
  onUnload() {
    try {
      if (this.innerAudioContext) {
        // 停止播放
        this.innerAudioContext.stop();
        // 销毁实例
        this.innerAudioContext.destroy();
        this.innerAudioContext = null;
      }
    } catch (error) {
      console.error('销毁音频上下文失败:', error);
    }
  },

  // 页面隐藏时暂停播放
  onHide() {
    try {
      if (this.innerAudioContext && this.data.isPlaying) {
        this.innerAudioContext.pause();
      }
    } catch (error) {
      console.error('暂停音频播放失败:', error);
    }
  },

  // 切换插画模式/纯净模式
  togglePicMode() {
    this.setData({
      isShowPic: !this.data.isShowPic
    });
  },

  // 图片加载错误处理
  onImageError(e) {
    console.error('图片加载失败:', e);
    wx.showToast({
      title: '图片加载失败',
      icon: 'none'
    });
  },
});
