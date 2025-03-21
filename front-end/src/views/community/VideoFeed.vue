<!--
 * @Author: Created for IdeaCosmos
 * @Description: 抖音风格视频内容播放页面
-->

<template>
  <div class="video-feed-container min-h-screen bg-black">
    <!-- 顶部导航栏 -->
    <div class="fixed top-0 left-0 right-0 bg-black/30 backdrop-blur-md z-10 text-white p-4 flex items-center justify-between">
      <div class="text-lg font-medium">视频内容</div>
      <div class="flex items-center gap-4">
        <button class="p-2" @click="goBack">
          <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
          </svg>
        </button>
      </div>
    </div>

    <!-- 视频容器 -->
    <div 
      class="video-swiper-container h-full w-full rounded-lg" 
      ref="swiperContainer"
      @touchstart="touchStart"
      @touchend="touchEnd"
      @wheel="handleWheel">
      <!-- 视频项目列表 -->
      <template v-if="videos.length > 0">
        <div 
          v-for="(video, index) in videos" 
          :key="video.id"
          class="video-item relative h-screen w-full rounded-lg" 
          :class="{ 
            'active': currentVideoIndex === index,
            'prev': currentVideoIndex > index,
            'next': currentVideoIndex < index,
            'transition-active': isTransitioning
          }"
          :style="{
            transform: getVideoPosition(index)
          }"
          :id="`video-${index}`">
          <!-- 视频元素 -->
          <video 
            :ref="el => { if (el) videoRefs[index] = el }"
            class="h-full w-full object-cover" 
            :src="video.videoUrl" 
            loop
            preload="auto"
            webkit-playsinline
            playsinline
            x5-playsinline
            @click="togglePlay">
          </video>

          <!-- 视频信息层 -->
          <div class="video-info absolute bottom-0 left-0 right-0 p-4 text-white bg-black/50">
            <!-- 用户信息 -->
            <div class="flex items-center gap-3 mb-4">
              <img :src="logo" class="w-12 h-12 rounded-full border-2 border-white" />
              <div>
                <div class="font-medium">{{ video.author.username }}</div>
                <div class="text-sm opacity-80">{{ formatDate(video.createdAt) }}</div>
              </div>
              <button class="ml-auto px-4 py-1 bg-blue-500 rounded-full text-white text-sm">关注</button>
            </div>
            
            <!-- 视频标题和描述 -->
            <div class="mb-6">
              <div class="text-lg font-medium mb-1">{{ video.title }}</div>
              <div class="text-sm opacity-80">{{ video.description.slice(0, 100) }}</div>
            </div>
          </div>

          <!-- 交互按钮组 -->
          <div class="absolute right-4 bottom-24 flex flex-col items-center gap-6">
            <!-- 点赞按钮 -->
            <div class="flex flex-col items-center" @click="toggleLike(video)">
              <div class="bg-black/50 p-3 rounded-full mb-1" :class="{'text-red-500': video.isLiked}">
                <svg xmlns="http://www.w3.org/2000/svg" class="h-8 w-8" :fill="video.isLiked ? 'currentColor' : 'white'" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4.318 6.318a4.5 4.5 0 000 6.364L12 20.364l7.682-7.682a4.5 4.5 0 00-6.364-6.364L12 7.636l-1.318-1.318a4.5 4.5 0 00-6.364 0z" />
                </svg>
              </div>
              <span class="text-sm text-white font-bold">{{ video.likes }}</span>
            </div>
            
            <!-- 收藏按钮 -->
            <div class="flex flex-col items-center" @click="toggleFavorite(video)">
              <div class="bg-black/50 p-3 rounded-full mb-1" :class="{'text-yellow-500': video.isFavorited}">
                <svg xmlns="http://www.w3.org/2000/svg" class="h-8 w-8" :fill="video.isFavorited ? 'currentColor' : 'white'" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11.049 2.927c.3-.921 1.603-.921 1.902 0l1.519 4.674a1 1 0 00.95.69h4.915c.969 0 1.371 1.24.588 1.81l-3.976 2.888a1 1 0 00-.363 1.118l1.518 4.674c.3.922-.755 1.688-1.538 1.118l-3.976-2.888a1 1 0 00-1.176 0l-3.976 2.888c-.783.57-1.838-.197-1.538-1.118l1.518-4.674a1 1 0 00-.363-1.118l-3.976-2.888c-.784-.57-.38-1.81.588-1.81h4.914a1 1 0 00.951-.69l1.519-4.674z" />
                </svg>
              </div>
              <span class="text-sm text-white font-bold">{{ video.favorites }}</span>
            </div>
            
            <!-- 评论按钮 -->
            <div class="flex flex-col items-center" @click="openCommentsModal(video)">
              <div class="bg-black/50 p-3 rounded-full mb-1">
                <svg xmlns="http://www.w3.org/2000/svg" class="h-8 w-8" fill="white" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z" />
                </svg>
              </div>
              <span class="text-sm text-white font-bold">{{ video.comments.length }}</span>
            </div>
            
            <!-- 分享按钮 -->
            <div class="flex flex-col items-center" @click="shareVideo(video)">
              <div class="bg-black/30 p-3 rounded-full mb-1">
                <svg xmlns="http://www.w3.org/2000/svg" class="h-8 w-8" fill="white" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8.684 13.342C8.886 12.938 9 12.482 9 12c0-.482-.114-.938-.316-1.342m0 2.684a3 3 0 110-2.684m0 2.684l6.632 3.316m-6.632-6l6.632-3.316m0 0a3 3 0 105.367-2.684 3 3 0 00-5.367 2.684zm0 9.316a3 3 0 105.368 2.684 3 3 0 00-5.368-2.684z" />
                </svg>
              </div>
              <span class="text-sm text-white font-bold">分享</span>
            </div>
          </div>

          <!-- 加载指示器 -->
          <div v-if="isLoading" class="absolute inset-0 flex items-center justify-center bg-black/50">
            <div class="animate-spin rounded-full h-16 w-16 border-t-2 border-b-2 border-white"></div>
          </div>
        </div>
      </template>

      <!-- 加载更多提示 -->
      <div v-if="isLoadingMore" class="loading-more text-white text-center p-4">
        <div class="animate-spin rounded-full h-8 w-8 border-t-2 border-b-2 border-white mx-auto mb-2"></div>
        <div>加载更多视频...</div>
      </div>

      <!-- 没有更多视频提示 -->
      <div v-if="noMoreVideos" class="no-more-videos text-white text-center p-8">
        没有更多视频了
      </div>

      <!-- 没有视频提示 -->
      <div v-if="!isLoading && videos.length === 0" class="no-videos-message h-screen flex flex-col items-center justify-center text-white">
        <svg xmlns="http://www.w3.org/2000/svg" class="h-16 w-16 mb-4 text-gray-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 10l4.553-2.276A1 1 0 0121 8.618v6.764a1 1 0 01-1.447.894L15 14M5 18h8a2 2 0 002-2V8a2 2 0 00-2-2H5a2 2 0 00-2 2v8a2 2 0 002 2z" />
        </svg>
        <div class="text-xl font-medium mb-2">还没有视频内容</div>
        <div class="text-gray-400 text-center max-w-xs">
          暂时没有可播放的视频，请稍后再来查看
        </div>
      </div>
    </div>

    <!-- 评论模态框 -->
    <a-modal
      v-model:visible="commentsModalVisible"
      title="评论"
      :footer="null"
      class="comments-modal"
      :width="420"
      :mask-closable="true"
      :destroyOnClose="true"
    >
      <div class="comments-container">
        <!-- 评论列表 -->
        <div class="comments-list max-h-96 overflow-y-auto px-2">
          <template v-if="selectedVideo && selectedVideo.comments.length > 0">
            <div 
              v-for="comment in selectedVideo.comments" 
              :key="comment.id"
              class="comment-item p-3 border-b"
            >
              <div class="flex items-start">
                <img :src="BACKEND_DOMAIN + comment.author.avatar" class="w-8 h-8 rounded-full mr-3" />
                <div class="flex-1">
                  <div class="font-medium text-sm">{{ comment.author.username }}</div>
                  <div class="text-sm text-gray-600 mt-1">{{ comment.content }}</div>
                  <div class="flex items-center justify-between mt-2">
                    <div class="text-xs text-gray-500">{{ formatDate(comment.createdAt) }}</div>
                    <button 
                      class="text-sm flex items-center gap-1" 
                      :class="{'text-red-500': comment.isLiked}"
                      @click="toggleCommentLike(comment)"
                    >
                      <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" :fill="comment.isLiked ? 'currentColor' : 'none'" viewBox="0 0 24 24" stroke="currentColor">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4.318 6.318a4.5 4.5 0 000 6.364L12 20.364l7.682-7.682a4.5 4.5 0 00-6.364-6.364L12 7.636l-1.318-1.318a4.5 4.5 0 00-6.364 0z" />
                      </svg>
                      {{ comment.likes }}
                    </button>
                  </div>
                </div>
              </div>
            </div>
          </template>
          <div v-else class="flex flex-col items-center justify-center py-8 text-gray-500">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-10 w-10 mb-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 10h.01M12 10h.01M16 10h.01M9 16H5a2 2 0 01-2-2V6a2 2 0 012-2h14a2 2 0 012 2v8a2 2 0 01-2 2h-5l-5 5v-5z" />
            </svg>
            <div class="text-center">还没有评论，快来发表第一条评论吧</div>
          </div>
        </div>

        <!-- 评论输入框 -->
        <div class="comment-input-container mt-4 flex items-center gap-2">
          <a-input 
            v-model:value="newComment" 
            placeholder="写下你的评论..." 
            class="flex-1" 
            :bordered="true"
            @pressEnter="addComment"
          />
          <a-button type="primary" @click="addComment" :disabled="!newComment.trim()">
            发送
          </a-button>
        </div>
      </div>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, onMounted, watch, computed, nextTick, onUnmounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { get, post, postJSON } from '@/util/request'
import { message, Modal } from 'ant-design-vue'
import { BACKEND_DOMAIN } from '@/util/VARRIBLES'
import { useUserStore } from '@/stores/user'
import { parseDateTime } from '@/util/common'
import logo from '@/assets/img/logo.webp'
const router = useRouter()
const route = useRoute()
const userStore = useUserStore()

// 视频数据
const videos = ref([])
const currentVideoIndex = ref(0)
const videoRefs = ref([])
const isLoading = ref(true)
const isLoadingMore = ref(false)
const noMoreVideos = ref(false)
const page = ref(1)
const pageSize = ref(10)

// 评论相关
const commentsModalVisible = ref(false)
const selectedVideo = ref(null)
const newComment = ref('')

// 触摸处理
const touchStartY = ref(0)
const touchEndY = ref(0)
const swiperContainer = ref(null)
const isTransitioning = ref(false)
const wheelThrottleTimeout = ref(null)
const wheelThrottleDelay = 700 // 滚轮节流延迟时间(ms)
const wheelDebounceTimer = ref(null)
const wheelDebounceDelay = 200 // 滚轮防抖延迟时间(ms)
const wheelDeltaAccumulator = ref(0) // 累积滚动距离

// 获取视频数据
const fetchVideos = async (loadMore = false) => {
  if (loadMore) {
    if (noMoreVideos.value) return
    isLoadingMore.value = true
    page.value++
  } else {
    isLoading.value = true
    page.value = 1
  }

  try {
    // 调用API获取篇章版本视频数据
    get(`/api/video/getVideosChapterVersion`, {
      page: page.value,
      limit: pageSize.value
    }, (msg, data) => {
      if (loadMore) {
        if (data.length === 0) {
          noMoreVideos.value = true
        } else {
          videos.value = [...videos.value, ...processVideosData(data)]
        }
      } else {
        videos.value = processVideosData(data)
      }
      
      isLoading.value = false
      isLoadingMore.value = false
      
      // 加载完成后播放第一个视频
      if (!loadMore && videos.value.length > 0) {
        nextTick(() => {
          playCurrentVideo()
        })
      }
    }, (msg) => {
      message.error(msg || '获取视频失败')
      isLoading.value = false
      isLoadingMore.value = false
    }, (msg) => {
      message.error('网络错误，请稍后再试')
      isLoading.value = false
      isLoadingMore.value = false
    })
  } catch (error) {
    message.error('获取视频失败，请稍后再试')
    console.error('获取视频失败:', error)
    isLoading.value = false
    isLoadingMore.value = false
  }
}

// 处理视频数据 - 根据后端实际返回的数据结构处理
const processVideosData = (data) => {
  return data.map(item => {
    // 确保用户对象存在基本属性
    const author = item.user || {
      ID: 0,
      username: '未知用户',
      avatar: '/default-avatar.png'
    }
    
    return {
      id: item.ID,
      title: `第${item.chapter_id}章视频`,
      description: item.content || '无描述',
      videoUrl: `${BACKEND_DOMAIN}videos/${item.video_path}`, // 视频路径
      createdAt: item.CreatedAt,
      updatedAt: item.UpdatedAt,
      chapterId: item.chapter_id,
      score: item.score,
      status: item.status,
      author: {
        id: author.ID,
        username: author.username,
        avatar: author.avatar || '/default-avatar.png'
      },
      // 本地状态 - 离线版点赞收藏
      isLiked: false,
      isFavorited: false,
      likes: Math.floor(Math.random() * 100) + 1, // 模拟点赞数
      favorites: Math.floor(Math.random() * 50) + 1, // 模拟收藏数
      comments: [] // 初始化空评论数组
    }
  })
}

// 播放/暂停当前视频
const togglePlay = () => {
  const currentVideo = videoRefs.value[currentVideoIndex.value]
  if (currentVideo) {
    if (currentVideo.paused) {
      currentVideo.play()
    } else {
      currentVideo.pause()
    }
  }
}

// 播放当前视频
const playCurrentVideo = () => {
  // 首先暂停所有视频
  videoRefs.value.forEach((video, index) => {
    if (video && index !== currentVideoIndex.value) {
      video.pause()
      video.currentTime = 0
    }
  })

  // 播放当前视频
  const currentVideo = videoRefs.value[currentVideoIndex.value]
  if (currentVideo) {
    currentVideo.play().catch(error => {
      console.error('视频播放失败:', error)
    })
  }
}

// 触摸开始事件
const touchStart = (event) => {
  touchStartY.value = event.touches[0].clientY
}

// 触摸结束事件
const touchEnd = (event) => {
  touchEndY.value = event.changedTouches[0].clientY
  const deltaY = touchEndY.value - touchStartY.value
  
  // 向上滑动，加载下一个视频
  if (deltaY < -100) {
    animateToNextVideo()
  }
  // 向下滑动，回到上一个视频
  else if (deltaY > 100) {
    animateToPrevVideo()
  }
}

// 滚轮事件处理
const handleWheel = (event) => {
  // 防止事件默认行为和冒泡
  event.preventDefault()
  
  // 如果正在过渡中，不要响应新的滚动
  if (isTransitioning.value) {
    return
  }
  
  // 累积滚动距离
  wheelDeltaAccumulator.value += event.deltaY
  
  // 防止快速连续滚动多次触发 - 节流
  if (wheelThrottleTimeout.value !== null) {
    return
  }
  
  // 清除之前的防抖定时器
  if (wheelDebounceTimer.value) {
    clearTimeout(wheelDebounceTimer.value)
  }
  
  // 设置防抖，等待用户停止滚动
  wheelDebounceTimer.value = setTimeout(() => {
    // 根据累积的滚动距离决定滚动方向
    if (Math.abs(wheelDeltaAccumulator.value) > 50) { // 设置一个阈值，避免轻微滚动触发
      if (wheelDeltaAccumulator.value > 0) {
        animateToNextVideo()
      } else {
        animateToPrevVideo()
      }
      
      // 设置节流，防止短时间内再次触发
      wheelThrottleTimeout.value = setTimeout(() => {
        wheelThrottleTimeout.value = null
      }, wheelThrottleDelay)
    }
    
    // 重置累积距离
    wheelDeltaAccumulator.value = 0
  }, wheelDebounceDelay)
}

// 滑动到下一个视频（带动画）
const animateToNextVideo = () => {
  if (currentVideoIndex.value < videos.value.length - 1) {
    isTransitioning.value = true
    currentVideoIndex.value++
    setTimeout(() => {
      isTransitioning.value = false
      playCurrentVideo()
    }, 300) // 过渡动画时间后播放
  } else {
    // 如果是最后一个视频，尝试加载更多
    if (!isLoadingMore.value && !noMoreVideos.value) {
      fetchVideos(true)
    }
  }
}

// 滑动到上一个视频（带动画）
const animateToPrevVideo = () => {
  if (currentVideoIndex.value > 0) {
    isTransitioning.value = true
    currentVideoIndex.value--
    setTimeout(() => {
      isTransitioning.value = false
      playCurrentVideo()
    }, 300) // 过渡动画时间后播放
  }
}

// 计算每个视频的位置
const getVideoPosition = (index) => {
  if (currentVideoIndex.value === index) {
    return 'translateY(0)'
  } else if (currentVideoIndex.value > index) {
    return 'translateY(-100%)'
  } else {
    return 'translateY(100%)'
  }
}

// 滑动到下一个视频
const goToNextVideo = () => {
  animateToNextVideo()
}

// 滑动到上一个视频
const goToPrevVideo = () => {
  animateToPrevVideo()
}

// 点赞操作 - 离线版
const toggleLike = (video) => {
  if (!userStore.isLogin) {
    message.warning('请先登录')
    return
  }

  // 本地状态更新
  video.isLiked = !video.isLiked
  video.likes += video.isLiked ? 1 : -1
  
  // 模拟成功提示
  message.success(video.isLiked ? '点赞成功' : '已取消点赞')
}

// 收藏操作 - 离线版
const toggleFavorite = (video) => {
  if (!userStore.isLogin) {
    message.warning('请先登录')
    return
  }

  // 本地状态更新
  video.isFavorited = !video.isFavorited
  video.favorites += video.isFavorited ? 1 : -1
  
  // 模拟成功提示
  message.success(video.isFavorited ? '收藏成功' : '已取消收藏')
}

// 评论操作
const openCommentsModal = (video) => {
  selectedVideo.value = video
  
  // 如果还没有加载评论，模拟一些评论数据
  if (video.comments.length === 0) {
    video.comments = generateMockComments(video.id)
  }
  
  commentsModalVisible.value = true
}

// 生成模拟评论数据
const generateMockComments = (videoId) => {
  const commentCount = Math.floor(Math.random() * 5) + 1
  const comments = []
  
  for (let i = 0; i < commentCount; i++) {
    comments.push({
      id: `${videoId}-comment-${i}`,
      content: `这是第${i+1}条评论，视频内容很精彩！`,
      createdAt: new Date(Date.now() - Math.floor(Math.random() * 1000000000)).toISOString(),
      likes: Math.floor(Math.random() * 20),
      isLiked: false,
      author: {
        id: i + 100,
        username: `用户${i + 1}`,
        avatar: '/default-avatar.png'
      }
    })
  }
  
  return comments
}

// 添加评论 - 离线版
const addComment = () => {
  if (!userStore.isLogin) {
    message.warning('请先登录')
    return
  }

  if (!newComment.value.trim()) {
    return
  }

  // 添加新评论到列表
  const newCommentObj = {
    id: `new-comment-${Date.now()}`,
    content: newComment.value.trim(),
    author: {
      id: userStore.user.id,
      username: userStore.user.username,
      avatar: userStore.user.avatar
    },
    createdAt: new Date().toISOString(),
    likes: 0,
    isLiked: false
  }

  selectedVideo.value.comments.unshift(newCommentObj)
  newComment.value = ''
  message.success('评论已发布')
}

// 给评论点赞 - 离线版
const toggleCommentLike = (comment) => {
  if (!userStore.isLogin) {
    message.warning('请先登录')
    return
  }

  comment.isLiked = !comment.isLiked
  comment.likes += comment.isLiked ? 1 : -1
}

// 分享视频
const shareVideo = (video) => {
  // 复制分享链接
  const shareUrl = `${window.location.origin}/community/video/${video.id}`
  navigator.clipboard.writeText(shareUrl)
    .then(() => {
      message.success('分享链接已复制到剪贴板')
    })
    .catch(() => {
      message.error('复制链接失败')
    })
}

// 返回上一页
const goBack = () => {
  router.back()
}

// 监听当前视频索引变化
watch(currentVideoIndex, (newIndex) => {
  // 检查是否需要加载更多视频
  if (newIndex >= videos.value.length - 3 && !isLoadingMore.value && !noMoreVideos.value) {
    fetchVideos(true)
  }
  
  // 播放当前视频
  if (!isTransitioning.value) {
    nextTick(() => {
      playCurrentVideo()
    })
  }
})

// 过滤器：格式化日期
const formatDate = (timestamp) => {
  return parseDateTime(timestamp)
}

// 生命周期钩子
onMounted(() => {
  fetchVideos()
  
  // 监听滚动事件，用于检测是否需要加载更多内容
  window.addEventListener('scroll', handleScroll)
})

// 组件卸载时清理
onUnmounted(() => {
  window.removeEventListener('scroll', handleScroll)
})

// 滚动处理
const handleScroll = () => {
  if (window.innerHeight + window.scrollY >= document.body.offsetHeight - 100) {
    if (!isLoadingMore.value && !noMoreVideos.value) {
      fetchVideos(true)
    }
  }
}
</script>

<style scoped>
.video-feed-container {
  position: relative;
  overflow: hidden;
}

.video-swiper-container {
  position: relative;
  height: 100vh;
  overflow: hidden;
}

.video-item {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  overflow: hidden;
  will-change: transform;
}

.video-item.transition-active {
  transition: transform 0.4s cubic-bezier(0.22, 0.61, 0.36, 1);
}

.video-item video {
  object-fit: cover;
  width: 100%;
  height: 100%;
  display: block;
}

/* 评论模态框自定义样式 */
:deep(.comments-modal .ant-modal-content) {
  border-radius: 16px;
  overflow: hidden;
}

:deep(.comments-modal .ant-modal-header) {
  border-bottom: 1px solid #f0f0f0;
  padding: 16px;
}

:deep(.comments-modal .ant-modal-body) {
  padding: 16px;
}

:deep(.comments-modal .ant-modal-footer) {
  border-top: 1px solid #f0f0f0;
  padding: 12px 16px;
}

/* 占位提示 */
.loading-more,
.no-more-videos {
  padding: 20px;
  text-align: center;
  font-size: 14px;
  position: relative;
  z-index: 5;
}
</style> 