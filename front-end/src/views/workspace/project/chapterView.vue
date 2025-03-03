<script setup>
import { ref, onMounted, computed, watch, onUnmounted, reactive } from "vue";
import { useRoute, useRouter } from "vue-router";
import { useThemeStore } from "@/stores/theme";
import { get, post, postJSON } from "@/util/request";
import { message } from "ant-design-vue";
import { BACKEND_DOMAIN, imagePrefix } from "@/util/VARRIBLES";
import { MdPreview } from "md-editor-v3";
import "md-editor-v3/lib/style.css";
import SpinLoaderLarge from "@/components/spinLoaderLarge.vue";
import { parseDateTime } from "@/util/common";
const route = useRoute();
const router = useRouter();
const themeStore = useThemeStore();
const loading = ref(true);
const chapter = ref(null);
const project = ref(null);

// 音频播放器相关状态
const audioPlayer = ref(null);

// 字体大小选项
const fontSizes = [
  { label: "小", value: "small" },
  { label: "中", value: "medium" },
  { label: "大", value: "large" },
];

const currentFontSize = ref("medium");

// 获取章节详情
const fetchChapterDetail = async () => {
  loading.value = true;
  await new Promise((resolve, reject) => {
    get(
      "/api/public/getChapterDetail",
      {
        id: route.params.id,
      },
      (messager, data) => {
        chapter.value = data.chapter;
        project.value = data.project;
        resolve();
      },
      (messager, data) => {
        message.warning(messager);
        reject();
      },
      (messager, data) => {
        message.warning(messager);
        reject();
      }
    );
  }).finally(() => {
    loading.value = false;
  });
};

// 音频错误处理
const onAudioError = () => {
  message.error("音频加载失败");
};

// 滚动到顶部函数
const scrollToTop = () => {
  window.scrollTo({
    top: 0,
    behavior: "smooth",
  });
};

// 控制回到顶部按钮显示
const showBackToTop = ref(false);

// 监听滚动事件
const handleScroll = () => {
  showBackToTop.value = window.scrollY > 300;
};

onMounted(() => {
  // 添加滚动监听
  window.addEventListener("scroll", handleScroll);
  // 立即执行滚动到顶部
  scrollToTop();
  // 然后获取章节详情
  fetchChapterDetail();
});

// 监听路由变化，每次进入页面都滚动到顶部
watch(
  () => route.params.id,
  () => {
    scrollToTop();
  }
);

// 组件卸载时移除监听器
onUnmounted(() => {
  window.removeEventListener("scroll", handleScroll);
});

// 情绪映射表
const emotionMap = {
  "喜悦": { 
    icon: "😊", 
    description: "充满快乐和满足" 
  },
  "感动": { 
    icon: "🥹", 
    description: "内心被深深触动" 
  },
  "惊喜": { 
    icon: "🤩", 
    description: "意外的惊喜" 
  },
  "期待": { 
    icon: "🤔", 
    description: "对后续充满期待" 
  },
  "伤感": { 
    icon: "😢", 
    description: "略带忧伤的感动" 
  },
  "愤怒": { 
    icon: "😠", 
    description: "对情节感到愤慨" 
  },
  "恐惧": { 
    icon: "😱", 
    description: "感到害怕或紧张" 
  },
  "平静": { 
    icon: "😐", 
    description: "内心平和安宁" 
  }
};

const selectedEmotion = ref("");
const userFeeling = ref(null);

// 获取用户对当前版本的情绪评价
const fetchUserFeeling = async () => {
  if (!chapter.value?.current_version?.ID) return;
  
  await get(
    "/api/user/feeling/get",
    { version_id: chapter.value.current_version.ID },
    (message, data) => {
      userFeeling.value = data;
    }
  );
};

// 提交情绪评价
const submitFeeling = async (emotion) => {
  if (!chapter.value?.current_version?.ID) return;
  
  selectedEmotion.value = emotion;
  await postJSON(
    "/api/user/feeling/add",
    {
      version_id: chapter.value.current_version.ID,
      feeling: emotion
    },
    (messager, data) => {
      message.success("评价成功");
      fetchUserFeeling(); // 重新获取评价状态
    },
    (messager) => {
      selectedEmotion.value = "";
    }
  );
};

// 在获取章节详情后获取情绪评价
watch(() => chapter.value?.current_version?.ID, () => {
  if (chapter.value?.current_version?.ID) {
    fetchUserFeeling();
  }
});

// 评论相关
const commentType = ref("all"); // all, reader, author
const comments = ref([]);
const commentForm = reactive({
  content: "",
  type: "reader", // 默认为读者评论
});

// 获取评论
const fetchComments = async () => {
  if (!chapter.value?.current_version?.ID) return;
  
  await get(
    "/api/user/getVersionComments",
    { 
      version_id: chapter.value.current_version.ID,
      type: commentType.value 
    },
    (msg, data) => {
      if (data !== null) {
        comments.value = data.reader_comments;
      }
    }
  );
};

// 提交评论
const submitComment = async () => {
  if (!chapter.value?.current_version?.ID || !commentForm.content.trim()) return;
  
  await postJSON(
    "/api/user/addVersionComment",
    {
      version_id: chapter.value.current_version.ID,
      content: commentForm.content,
      type: commentForm.type
    },
    (msg, data) => {
      message.success("评论发表成功");
      commentForm.content = "";
      fetchComments(); // 重新加载评论
    },
    (msg) => {
      message.error(msg || "评论发表失败");
    }
  );
};

// 获取用户头像
const getAvatarUrl = (user) => {
  return user.avatar ? `${BACKEND_DOMAIN}${user.avatar}` : 
    (user.gender === "male" ? maleAvatar : femaleAvatar);
};

// 监听评论类型变化
watch(commentType, () => {
  fetchComments();
});

// 监听章节变化时获取评论
watch(() => chapter.value?.current_version?.ID, () => {
  if (chapter.value?.current_version?.ID) {
    fetchComments();
  }
});

// 初始加载时获取评论
onMounted(() => {
  if (chapter.value?.current_version?.ID) {
    fetchComments();
  }
});
</script>

<template>
  <div class="min-h-screen py-8">
    <!-- 加载状态 -->
    <div v-if="loading" class="flex items-center justify-center min-h-[400px]">
      <SpinLoaderLarge />
    </div>

    <div v-else>
      <div class="max-w-5xl mx-auto px-4">

        <!-- 项目信息卡片 -->
        <div
          class="bg-white dark:bg-zinc-800/80 rounded-xl p-6 border theme-border"
        >
        <!-- 导航栏 -->
        <nav class="flex items-center gap-2 mb-6 text-sm">
          <button
            @click="router.back()"
            class="flex items-center gap-1 text-gray-600 dark:text-gray-400 hover:text-gray-900 dark:hover:text-gray-200"
          >
            <svg
              xmlns="http://www.w3.org/2000/svg"
              class="h-4 w-4"
              fill="none"
              viewBox="0 0 24 24"
              stroke="currentColor"
            >
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                d="M15 19l-7-7 7-7"
              />
            </svg>
            返回项目
          </button>
          <span class="text-gray-400">/</span>
          <span class="text-gray-600 dark:text-gray-400">{{
            project?.project_name
          }}</span>
        </nav>
          <h3
            class="text-xl font-bold bg-gradient-to-r from-blue-500 to-cyan-500 bg-clip-text text-transparent mb-4"
          >
            {{ project.project_name }}
          </h3>
          <div class="space-y-4">
            <div>
              <h4
                class="text-sm font-medium text-gray-500 dark:text-gray-400 mb-2"
              >
                项目简介
              </h4>
              <p class="text-gray-700 dark:text-gray-300">
                {{ project?.social_story }}
              </p>
            </div>
            <div class="flex flex-wrap gap-2">
              <span
                v-for="style in project?.style"
                :key="style"
                class="px-3 py-1 text-sm bg-blue-100/80 dark:bg-blue-900/30 text-blue-700 dark:text-blue-300 rounded-full"
              >
                {{ style }}
              </span>
            </div>
          </div>
        </div>
        <!-- 主要内容 -->
        <div class="space-y-6 mt-4">
          <!-- 章节信息卡片 -->
          <div
            class="bg-white dark:bg-zinc-800/80 rounded-xl p-6 border theme-border"
          >
            <div class="space-y-4">
              <div class="flex items-center justify-between">
                <h1
                  class="text-3xl font-bold bg-gradient-to-r from-blue-500 to-cyan-500 bg-clip-text text-transparent"
                >
                  {{ chapter?.Title }}
                </h1>
              </div>

              <!-- 音频播放器 - 根据音频路径显示 -->
              <div
                v-if="chapter?.current_version?.audio_path"
                class="mt-4 p-4 bg-gray-50 dark:bg-zinc-900/50 rounded-xl border border-gray-200 dark:border-zinc-700/80"
              >
                <div class="flex items-center gap-3 mb-2">
                  <svg
                    xmlns="http://www.w3.org/2000/svg"
                    class="h-5 w-5 text-blue-500"
                    fill="none"
                    viewBox="0 0 24 24"
                    stroke="currentColor"
                  >
                    <path
                      stroke-linecap="round"
                      stroke-linejoin="round"
                      d="M9 19V6l12-3v13M9 19c0 1.105-1.343 2-3 2s-3-.895-3-2 1.343-2 3-2 3 .895 3 2zm12-3c0 1.105-1.343 2-3 2s-3-.895-3-2 1.343-2 3-2 3 .895 3 2zM9 10l12-3"
                    />
                  </svg>
                  <span
                    class="text-sm font-medium text-gray-700 dark:text-gray-300"
                    >章节音频</span
                  >
                </div>
                <audio
                  ref="audioPlayer"
                  :src="`${BACKEND_DOMAIN}audio/${chapter.current_version.audio_path}`"
                  class="w-full focus:outline-none"
                  controls
                  preload="metadata"
                  @error="onAudioError"
                ></audio>
              </div>

              <p class="text-gray-600 dark:text-gray-400 text-lg">
                {{ chapter?.Description }}
              </p>

              <div
                class="flex flex-wrap gap-4 text-sm text-gray-500 dark:text-gray-400 pt-4 border-t theme-border"
              >
                <div class="flex items-center gap-2">
                  <svg
                    xmlns="http://www.w3.org/2000/svg"
                    class="h-5 w-5"
                    fill="none"
                    viewBox="0 0 24 24"
                    stroke="currentColor"
                  >
                    <path
                      stroke-linecap="round"
                      stroke-linejoin="round"
                      d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"
                    />
                  </svg>
                  <span>创建于 {{ parseDateTime(chapter?.CreatedAt) }}</span>
                </div>
                <div class="flex items-center gap-2">
                  <svg
                    xmlns="http://www.w3.org/2000/svg"
                    class="h-5 w-5"
                    fill="none"
                    viewBox="0 0 24 24"
                    stroke="currentColor"
                  >
                    <path
                      stroke-linecap="round"
                      stroke-linejoin="round"
                      d="M15.232 5.232l3.536 3.536m-2.036-5.036a2.5 2.5 0 113.536 3.536L6.5 21.036H3v-3.572L16.732 3.732z"
                    />
                  </svg>
                  <span>最后更新 {{ parseDateTime(chapter?.UpdatedAt) }}</span>
                </div>
                <div class="flex items-center gap-2">
                  <svg
                    xmlns="http://www.w3.org/2000/svg"
                    class="h-5 w-5"
                    fill="none"
                    viewBox="0 0 24 24"
                    stroke="currentColor"
                  >
                    <path
                      stroke-linecap="round"
                      stroke-linejoin="round"
                      d="M3 7v10a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2h-6l-2-2H5a2 2 0 00-2 2z"
                    />
                  </svg>
                  <span
                    >字数统计：{{
                      chapter?.current_version?.content?.length || 0
                    }}</span
                  >
                </div>
              </div>
              <div class="space-y-4 select-text">
                <div class="flex items-center justify-end gap-2 mb-3">
                  <span class="text-sm text-gray-500 dark:text-gray-400">字体大小</span>
                  <div class="flex items-center gap-1 bg-gray-100 dark:bg-zinc-900 rounded-lg p-1">
                    <button
                      v-for="size in fontSizes"
                      :key="size.value"
                      @click="currentFontSize = size.value"
                      :class="[
                        'px-3 py-1 text-sm rounded-md transition-colors',
                        currentFontSize === size.value
                          ? 'bg-blue-500 text-white'
                          : 'text-gray-600 dark:text-gray-400 hover:bg-gray-200 dark:hover:bg-zinc-800',
                      ]"
                    >
                      {{ size.label }}
                    </button>
                  </div>
                </div>

                <div :class="['content-container', `font-${currentFontSize}`]">
                  <MdPreview
                    style="background: transparent"
                    :theme="themeStore.currentTheme"
                    editorId="chapter-preview"
                    :modelValue="chapter?.current_version?.content || ''"
                    previewTheme="github"
                    :previewOnly="true"
                  />
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- 在章节内容后添加情绪评价部分 -->
        <div v-if="chapter?.current_version" class="mt-8 p-6 bg-white dark:bg-zinc-800/80 rounded-xl border theme-border">
          <div class="text-center">
            <h3 class="text-lg font-semibold mb-2 text-gray-800 dark:text-gray-200">你对本篇剧情感受如何？</h3>
            <p class="text-sm text-gray-500 dark:text-gray-400 mb-6">选择一个最能代表你此刻感受的情绪</p>
          </div>
          
          <!-- 如果已经评价过，显示已选择的情绪 -->
          <div v-if="userFeeling !== '获取失败'" class="text-center">
            <p class="flex items-center justify-center gap-2 text-gray-600 dark:text-gray-400">
              <span>你的感受是:</span>
              <span class="inline-flex items-center gap-1 px-3 py-1 rounded-full bg-blue-50 dark:bg-blue-900/30 text-blue-600 dark:text-blue-300">
                {{ emotionMap[userFeeling.feeling]?.icon }}
                <span class="font-medium">{{ userFeeling.feeling }}</span>
              </span>
            </p>
          </div>
          
          <!-- 如果还没评价过，显示情绪选择按钮 -->
          <div v-else class="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 gap-4 max-w-2xl mx-auto">
            <button
              v-for="emotion in Object.entries(emotionMap)"
              :key="emotion"
              @click="submitFeeling(emotion[0])"
              class="flex flex-col items-center gap-2 p-4 rounded-xl border theme-border hover:bg-blue-50 dark:hover:bg-blue-900/30 transition-all transform hover:scale-105"
              :class="{
                'bg-blue-50 dark:bg-blue-900/30 border-blue-500': selectedEmotion === emotion[0]
              }"
            >
              <span class="text-2xl">{{ emotion[1].icon }}</span>
              <span class="text-sm font-medium text-gray-700 dark:text-gray-300">{{ emotion[0] }}</span>
              <span class="text-xs text-gray-500 dark:text-gray-400">{{ emotion[1].description }}</span>
            </button>
          </div>
        </div>

        <!-- 评论区域 -->
        <div class="bg-white dark:bg-zinc-800/80 rounded-xl p-6 border theme-border mt-6">
          <h2 class="text-xl font-bold text-gray-900 dark:text-gray-100 mb-4">读者评论</h2>
          
          <!-- 评论表单 -->
          <div class="mb-6 border-b theme-border pb-6">
            <textarea
              v-model="commentForm.content"
              rows="3"
              class="w-full rounded-lg border theme-border p-3 bg-white dark:bg-zinc-900 text-gray-800 dark:text-gray-200 focus:ring-2 focus:ring-blue-500 focus:border-transparent transition-all resize-none"
              placeholder="分享你对这个章节的看法..."
            ></textarea>
            <div class="flex justify-end mt-3">
              <button
                @click="submitComment"
                class="px-4 py-2 bg-blue-500 text-white rounded-lg hover:bg-blue-600 disabled:opacity-50 disabled:cursor-not-allowed transition-all"
              >
                发表评论
              </button>
            </div>
          </div>
          
          <!-- 切换评论类型 -->
          <div class="flex gap-4 mb-4">
            <button
              @click="commentType = 'all'"
              class="px-3 py-1 text-sm rounded-full transition-colors"
              :class="commentType === 'all' ? 'bg-blue-100 dark:bg-blue-900/30 text-blue-700 dark:text-blue-300' : 'text-gray-600 dark:text-gray-400'"
            >
              全部评论
            </button>
            <button
              @click="commentType = 'reader'"
              class="px-3 py-1 text-sm rounded-full transition-colors"
              :class="commentType === 'reader' ? 'bg-green-100 dark:bg-green-900/30 text-green-700 dark:text-green-300' : 'text-gray-600 dark:text-gray-400'"
            >
              读者评论
            </button>
            <button
              @click="commentType = 'author'"
              class="px-3 py-1 text-sm rounded-full transition-colors"
              :class="commentType === 'author' ? 'bg-purple-100 dark:bg-purple-900/30 text-purple-700 dark:text-purple-300' : 'text-gray-600 dark:text-gray-400'"
            >
              作者评论
            </button>
          </div>
          
          <!-- 评论列表 -->
          <div v-if="comments.length === 0" class="text-center py-10 text-gray-500">
            还没有评论，快来发表第一条评论吧！
          </div>
          
          <div v-else class="space-y-4">
            <div
              v-for="comment in comments"
              :key="comment.id"
              class="p-4 border theme-border rounded-lg"
            >
              <div class="flex items-start gap-3">
                <img
                  :src="getAvatarUrl(comment.user)"
                  class="w-10 h-10 rounded-full object-cover"
                  alt="用户头像"
                />
                <div class="flex-1">
                  <div class="flex items-center justify-between">
                    <div class="flex items-center gap-2">
                      <h4 class="font-medium text-gray-900 dark:text-gray-100">
                        {{ comment.user.username }}
                      </h4>
                      <span
                        v-if="comment.type === 'author'"
                        class="px-2 py-0.5 bg-purple-100 dark:bg-purple-900/30 text-purple-700 dark:text-purple-300 text-xs rounded-full"
                      >
                        作者
                      </span>
                    </div>
                    <span class="text-sm text-gray-500">
                      {{ parseDateTime(comment.CreatedAt) }}
                    </span>
                  </div>
                  <p class="mt-2 text-gray-700 dark:text-gray-300">
                    {{ comment.Content }}
                  </p>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
      <!-- 回到顶部按钮 -->
    </div>
  </div>
</template>

<style lang="postcss">
/* 自定义音频播放器样式 */
audio {
  @apply h-12;
}

/* WebKit (Chrome, Safari) 样式 */
audio::-webkit-media-controls-panel {
  @apply bg-white dark:bg-zinc-800 rounded-lg;
}

audio::-webkit-media-controls-play-button {
  @apply bg-blue-500 hover:bg-blue-600 rounded-full transition-colors;
}

audio::-webkit-media-controls-timeline {
  @apply bg-gray-200 dark:bg-zinc-700 rounded-full overflow-hidden;
}

audio::-webkit-media-controls-current-time-display,
audio::-webkit-media-controls-time-remaining-display {
  @apply text-gray-700 dark:text-gray-300;
}

audio::-webkit-media-controls-volume-slider {
  @apply bg-gray-200 dark:bg-zinc-700 rounded-full;
}

/* Firefox 样式 */
audio::-moz-range-track {
  @apply bg-gray-200 dark:bg-zinc-700 rounded-full;
}

audio::-moz-range-thumb {
  @apply bg-blue-500;
}

/* 调整 Markdown 预览区域的样式 */
.md-editor-preview-wrapper {
  @apply px-0 !important;
}

/* 确保不同字体大小下的间距一致 */
.md-editor-preview {
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
  ul,
  ol {
    @apply mb-4;
  }
  li {
    @apply mb-2;
  }

  /* 根据字体大小调整代码块和引用的内边距 */
  pre,
  blockquote {
    @apply my-4;
  }
}

/* 添加过渡效果 */
.fixed {
  @apply transition-opacity duration-300 ease-in-out;
}

/* 添加情绪按钮的悬停动画 */
.transform {
  transition: all 0.2s ease-in-out;
}

/* 确保emoji在暗色模式下显示正常 */
.text-2xl {
  font-family: "Apple Color Emoji", "Segoe UI Emoji", "Segoe UI Symbol";
}

/* 字体大小控制 */
.content-container.font-small .md-editor-preview {
  font-size: 0.875rem !important;
  line-height: 1.6;
}

.content-container.font-medium .md-editor-preview {
  font-size: 1rem !important;
  line-height: 1.7;
}

.content-container.font-large .md-editor-preview {
  font-size: 1.125rem !important;
  line-height: 1.8;
}

/* 控制标题大小 */
.content-container.font-small .md-editor-preview h1 {
  font-size: 1.6rem !important;
}

.content-container.font-small .md-editor-preview h2 {
  font-size: 1.4rem !important;
}

.content-container.font-small .md-editor-preview h3 {
  font-size: 1.2rem !important;
}

.content-container.font-medium .md-editor-preview h1 {
  font-size: 1.8rem !important;
}

.content-container.font-medium .md-editor-preview h2 {
  font-size: 1.6rem !important;
}

.content-container.font-medium .md-editor-preview h3 {
  font-size: 1.4rem !important;
}

.content-container.font-large .md-editor-preview h1 {
  font-size: 2rem !important;
}

.content-container.font-large .md-editor-preview h2 {
  font-size: 1.8rem !important;
}

.content-container.font-large .md-editor-preview h3 {
  font-size: 1.6rem !important;
}
</style>
