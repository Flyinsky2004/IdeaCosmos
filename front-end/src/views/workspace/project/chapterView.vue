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
                  <span>创建于 {{ formatDate(chapter?.CreatedAt) }}</span>
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
                  <span>最后更新 {{ formatDate(chapter?.UpdatedAt) }}</span>
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
              <div class="space-y-4">
                <div class="flex items-center justify-end gap-2">
                  <span class="text-sm text-gray-500 dark:text-gray-400"
                    >字体大小</span
                  >
                  <div
                    class="flex items-center gap-1 bg-gray-100 dark:bg-zinc-900 rounded-lg p-1"
                  >
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

                <MdPreview
                  style="background: transparent"
                  :theme="themeStore.currentTheme"
                  editorId="chapter-preview"
                  :modelValue="wrappedContent"
                  previewTheme="github"
                  :previewOnly="true"
                />
              </div>
            </div>
          </div>
        </div>
      </div>
      <!-- 回到顶部按钮 -->
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, computed, watch, onUnmounted } from "vue";
import { useRoute, useRouter } from "vue-router";
import { useThemeStore } from "@/stores/theme";
import { get } from "@/util/request";
import { message } from "ant-design-vue";
import { BACKEND_DOMAIN, imagePrefix } from "@/util/VARRIBLES";
import { MdPreview } from "md-editor-v3";
import "md-editor-v3/lib/style.css";
import SpinLoaderLarge from "@/components/spinLoaderLarge.vue";

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

// 包装 Markdown 内容
const wrappedContent = computed(() => {
  if (!chapter.value?.current_version?.content) return "";

  const sizeClasses = {
    small: "text-sm leading-relaxed",
    medium: "text-base leading-relaxed",
    large: "text-lg leading-relaxed",
  };

  return `
<div class="${sizeClasses[currentFontSize.value]}">

${chapter.value.current_version.content}

</div>
`;
});

// 格式化日期
const formatDate = (date) => {
  if (!date) return "未知时间";
  return new Date(date).toLocaleDateString("zh-CN", {
    year: "numeric",
    month: "long",
    day: "numeric",
    hour: "2-digit",
    minute: "2-digit",
  });
};

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
</script>

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
</style>
