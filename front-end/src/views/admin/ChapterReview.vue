<script setup>
import { ref, onMounted, reactive } from 'vue'
import { get, post, postJSON } from '@/util/request'
import { message } from 'ant-design-vue'
import SpinLoaderLarge from '@/components/spinLoaderLarge.vue'
import { BACKEND_DOMAIN } from "@/util/VARRIBLES"
import { MdPreview } from "md-editor-v3"
import "md-editor-v3/lib/style.css"
import { parseDateTime } from '@/util/common'
import { useThemeStore } from "@/stores/theme"
import { useRoute, useRouter } from 'vue-router'

const route = useRoute()
const router = useRouter()
const themeStore = useThemeStore()
const loading = ref(true)
const chapterInfo = ref(null)
// 获取状态文本
const getStatusText = (status) => {
  switch (status) {
    case 'pending': return '待审核'
    case 'approved': return '已通过'
    case 'rejected': return '已拒绝'
    default: return '未知状态'
  }
}

// AI评分功能
const aiScoreLoading = ref(false)
const requestAIScore = () => {
  const chapterId = route.params.id
  
  // 显示加载状态
  aiScoreLoading.value = true
  message.loading({ content: 'AI正在分析内容并评分...', key: 'aiScore', duration: 0 })
  
  get(`/api/admin/chapters/${chapterId}/ai-score`, {}, 
    (msg, data) => {
      // 更新评分和理由
      reviewState.score = data.score
      reviewState.reason = data.reason
      
      message.success({ content: 'AI评分完成！', key: 'aiScore' })
      aiScoreLoading.value = false
    },
    (msg) => {
      message.warning({ content: msg, key: 'aiScore' })
      aiScoreLoading.value = false
    },
    (msg) => {
      message.error({ content: '获取AI评分失败: ' + msg, key: 'aiScore' })
      aiScoreLoading.value = false
    }
  )
}

// 审核状态
const reviewState = reactive({
  loading: false,
  score: 70,
  reason: '',
})

// 获取章节详情
const fetchChapterDetail = async () => {
  loading.value = true
  const chapterId = route.params.id
  
  if (!chapterId) {
    message.error('章节ID无效')
    router.push('/admin/chapters')
    return
  }
  
  get(`/api/admin/chapters/${chapterId}`, {}, 
    (msg, data) => {
      // 验证是否有有效的当前版本
      if (!data || !data.current_version) {
        message.error('章节数据异常，可能没有关联的版本')
        setTimeout(() => router.push('/admin/chapters'), 1500)
        return
      }
      
      chapterInfo.value = data
      // 如果已有评分，使用现有评分
      if (data.current_version && data.current_version.score !== undefined) {
        reviewState.score = data.current_version.score
      }
      loading.value = false
    },
    (msg) => {
      message.warning(msg)
      loading.value = false
      router.push('/admin/chapters')
    },
    (msg) => {
      message.error(msg)
      loading.value = false
      router.push('/admin/chapters')
    }
  )
}

onMounted(() => {
  fetchChapterDetail()
})

// 返回列表页
const goBackToList = () => {
  router.push('/admin/chapters')
}

// 仅评分不审核
const updateScore = () => {
  reviewState.loading = true
  const chapterId = route.params.id
  
  postJSON(`/api/admin/chapters/${chapterId}/score`, {
    score: reviewState.score
  }, 
  (msg, data) => {
    message.success('章节评分更新成功')
    reviewState.loading = false
    setTimeout(() => {
      router.push('/admin/chapters')
    }, 1500)
  },
  (msg) => {
    message.warning(msg)
    reviewState.loading = false
  },
  (msg) => {
    message.error(msg)
    reviewState.loading = false
  })
}

// 审核章节
const reviewChapter = (status) => {
  reviewState.loading = true
  const chapterId = route.params.id
  
  postJSON(`/api/admin/chapters/${chapterId}/review`, {
    status: status,
    reason: reviewState.reason,
    score: reviewState.score
  }, 
  (msg, data) => {
    message.success(`章节${status === 'approved' ? '审核通过' : '审核拒绝'}成功`)
    reviewState.loading = false
    setTimeout(() => {
      router.push('/admin/chapters')
    }, 1500)
  },
  (msg) => {
    message.warning(msg)
    reviewState.loading = false
  },
  (msg) => {
    message.error(msg)
    reviewState.loading = false
  })
}

// 获取评分级别文本和颜色
const getScoreLevelText = (score) => {
  if (score >= 90) {
    return { text: '优秀', color: 'green' }
  } else if (score >= 70) {
    return { text: '良好', color: 'blue' }
  } else if (score >= 50) {
    return { text: '一般', color: 'orange' }
  } else {
    return { text: '敏感', color: 'red' }
  }
}
</script>

<template>
  <div class="p-4 max-w-7xl mx-auto">
    <!-- 面包屑导航 -->
    <div class="flex items-center mb-6">
      <a-button 
        type="link" 
        @click="goBackToList" 
        class="flex items-center text-gray-600 dark:text-gray-300 hover:text-blue-500 dark:hover:text-blue-400"
      >
        <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 mr-1" viewBox="0 0 20 20" fill="currentColor">
          <path fill-rule="evenodd" d="M9.707 16.707a1 1 0 01-1.414 0l-6-6a1 1 0 010-1.414l6-6a1 1 0 011.414 1.414L5.414 9H17a1 1 0 110 2H5.414l4.293 4.293a1 1 0 010 1.414z" clip-rule="evenodd" />
        </svg>
        <span>返回章节列表</span>
      </a-button>
      <span class="mx-2 text-gray-400">/</span>
      <span class="text-gray-600 dark:text-gray-300">章节审核</span>
    </div>

    <div v-if="loading" class="flex justify-center items-center h-64">
      <SpinLoaderLarge />
    </div>

    <div v-else class="animate__animated animate__fadeIn">
      <!-- 章节标题和基本信息 -->
      <div class="bg-white dark:bg-zinc-800 rounded-lg shadow-sm p-6 mb-6">
        <h1 class="text-2xl font-bold mb-2 dark:text-gray-100">{{ chapterInfo.Title }}</h1>
        <div class="flex flex-wrap gap-2 mb-4">
          <a-tag>项目ID: {{ chapterInfo.project_id }}</a-tag>
          <a-tag color="blue">{{ chapterInfo.current_version?.status ? getStatusText(chapterInfo.current_version.status) : '待审核' }}</a-tag>
          <a-tag :color="getScoreLevelText(chapterInfo.current_version?.score || 0).color">
            评分: {{ chapterInfo.current_version?.score || '无' }}
          </a-tag>
          <a-tag color="purple">更新时间: {{ parseDateTime(chapterInfo.UpdatedAt) }}</a-tag>
        </div>
        <p class="text-gray-600 dark:text-gray-300">{{ chapterInfo.Description }}</p>
      </div>

      <!-- 评分和审核区域 -->
      <div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
        <!-- 左侧评分和审核 -->
        <div class="lg:col-span-1">
          <div class="bg-white dark:bg-zinc-800 rounded-lg shadow-sm p-6 sticky top-6">
            <h2 class="text-xl font-bold mb-4 dark:text-gray-100">内容评审</h2>
            
            <!-- 评分选择器 -->
            <div class="mb-6">
              <h3 class="text-lg font-medium mb-2 dark:text-gray-200">内容评分</h3>
              <p class="text-sm text-gray-500 dark:text-gray-400 mb-2">
                请为章节内容评分(0-100)，评分标准：
              </p>
              <ul class="list-disc list-inside text-sm text-gray-500 dark:text-gray-400 mb-4">
                <li>90-100分：优秀，内容健康积极，适合全年龄段</li>
                <li>70-89分：良好，内容适合青少年，无敏感内容</li>
                <li>50-69分：一般，包含少量敏感内容，需谨慎</li>
                <li>0-49分：敏感，包含较多政治、道德、宗教等敏感内容</li>
              </ul>
              
              <div class="flex items-center mb-2">
                <span class="mr-2 dark:text-gray-300">评分:</span>
                <a-slider 
                  v-model:value="reviewState.score" 
                  :min="0" 
                  :max="100" 
                  class="flex-1"
                />
                <a-input-number 
                  v-model:value="reviewState.score" 
                  :min="0" 
                  :max="100"
                  class="ml-4 w-16"
                />
              </div>
              
              <div 
                class="py-1 px-2 rounded-md text-sm mt-1"
                :class="{
                  'bg-green-100 text-green-800 dark:bg-green-900/30 dark:text-green-300': reviewState.score >= 90,
                  'bg-blue-100 text-blue-800 dark:bg-blue-900/30 dark:text-blue-300': reviewState.score >= 70 && reviewState.score < 90,
                  'bg-orange-100 text-orange-800 dark:bg-orange-900/30 dark:text-orange-300': reviewState.score >= 50 && reviewState.score < 70,
                  'bg-red-100 text-red-800 dark:bg-red-900/30 dark:text-red-300': reviewState.score < 50
                }"
              >
                评级: {{ getScoreLevelText(reviewState.score).text }}
              </div>
            </div>
            
            <!-- AI评分按钮 -->
            <div class="mb-4 mt-4">
              <a-button 
                @click="requestAIScore" 
                class="w-full" 
                :loading="aiScoreLoading"
                type="dashed"
              >
                <template v-if="!aiScoreLoading">
                  <div class="flex items-center justify-center w-full">
                    <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 mr-1 inline-block" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                      <path stroke-linecap="round" stroke-linejoin="round" d="M13 10V3L4 14h7v7l9-11h-7z" />
                    </svg>
                    <span>使用AI分析内容并评分</span>
                  </div>
                </template>
              </a-button>
              <div class="text-xs text-gray-500 dark:text-gray-400 mt-1 text-center">
                AI将分析章节内容并提供建议评分和理由
              </div>
            </div>
            
            <!-- 审核意见 -->
            <div class="mb-6">
              <h3 class="text-lg font-medium mb-2 dark:text-gray-200">审核意见</h3>
              <a-textarea 
                v-model:value="reviewState.reason" 
                placeholder="可选填写审核意见，特别是拒绝时的理由说明" 
                :rows="3"
                class="mb-4"
              />
            </div>
            
            <!-- 操作按钮 -->
            <div class="flex flex-col gap-3">
              <a-button 
                type="primary" 
                class="bg-green-500 hover:bg-green-600 focus:bg-green-600" 
                :loading="reviewState.loading" 
                @click="reviewChapter('approved')"
                block
              >
                <template v-if="!reviewState.loading">
                  <div class="flex items-center justify-center w-full">
                    <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 mr-1 inline-block" viewBox="0 0 20 20" fill="currentColor">
                      <path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clip-rule="evenodd" />
                    </svg>
                    <span>审核通过</span>
                  </div>
                </template>
              </a-button>
              
              <a-button 
                danger 
                :loading="reviewState.loading" 
                @click="reviewChapter('rejected')"
                block
              >
                <template v-if="!reviewState.loading">
                  <div class="flex items-center justify-center w-full">
                    <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 mr-1 inline-block" viewBox="0 0 20 20" fill="currentColor">
                      <path fill-rule="evenodd" d="M4.293 4.293a1 1 0 011.414 0L10 8.586l4.293-4.293a1 1 0 111.414 1.414L11.414 10l4.293 4.293a1 1 0 01-1.414 1.414L10 11.414l-4.293 4.293a1 1 0 01-1.414-1.414L8.586 10 4.293 5.707a1 1 0 010-1.414z" clip-rule="evenodd" />
                    </svg>
                    <span>审核拒绝</span>
                  </div>
                </template>
              </a-button>
              
              <a-button 
                :loading="reviewState.loading" 
                @click="updateScore"
                block
              >
                <template v-if="!reviewState.loading">
                  <div class="flex items-center justify-center w-full">
                    <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 mr-1 inline-block" viewBox="0 0 20 20" fill="currentColor">
                      <path d="M9.049 2.927c.3-.921 1.603-.921 1.902 0l1.07 3.292a1 1 0 00.95.69h3.462c.969 0 1.371 1.24.588 1.81l-2.8 2.034a1 1 0 00-.364 1.118l1.07 3.292c.3.921-.755 1.688-1.54 1.118l-2.8-2.034a1 1 0 00-1.175 0l-2.8 2.034c-.784.57-1.838-.197-1.539-1.118l1.07-3.292a1 1 0 00-.364-1.118L2.98 8.72c-.783-.57-.38-1.81.588-1.81h3.461a1 1 0 00.951-.69l1.07-3.292z" />
                    </svg>
                    <span>仅更新评分</span>
                  </div>
                </template>
              </a-button>
            </div>
          </div>
        </div>
        
        <!-- 右侧内容预览 -->
        <div class="lg:col-span-2">
          <div class="bg-white dark:bg-zinc-800 rounded-lg shadow-sm p-6">
            <h2 class="text-xl font-bold mb-4 dark:text-gray-100">章节内容</h2>
            
            <!-- 章节内容 -->
            <div class="prose dark:prose-invert max-w-none">
              <MdPreview 
                :modelValue="chapterInfo.current_version ? 
                          chapterInfo.current_version.content : 
                          '无内容'"
                style="background: transparent;"
                class="md-preview p-2 rounded-xl" 
                :theme="themeStore.isDark ? 'dark' : 'light'" 
              />
              
              <!-- 音频 -->
              <div v-if="chapterInfo.current_version && chapterInfo.current_version.audio_path" class="mt-6">
                <h3 class="text-lg font-medium mb-2 dark:text-gray-200">音频资源</h3>
                <audio controls class="w-full">
                  <source :src="`${BACKEND_DOMAIN}/uploads/${chapterInfo.current_version.audio_path}`" type="audio/mpeg">
                  您的浏览器不支持音频播放
                </audio>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped lang="postcss">
.animate__animated {
  animation-duration: 0.5s;
}

.prose {
  max-height: 80vh;
  overflow-y: auto;
  padding-right: 1rem;
}

/* 调整 Markdown 预览区域的样式 */
:deep(.md-editor-preview-wrapper) {
  @apply px-0 !important;
}

/* 确保不同字体大小下的间距一致 */
:deep(.md-editor-preview) {
  h1 {
    @apply mb-6;
  }
  h2 {
    @apply mb-5;
  }
  h3 {
    @apply mb-4;
  }
  p {
    @apply mb-4;
  }
  ul, ol {
    @apply mb-4;
  }
  li {
    @apply mb-2;
  }
  
  /* 根据字体大小调整代码块和引用的内边距 */
  pre, blockquote {
    @apply my-4;
  }
}
</style> 