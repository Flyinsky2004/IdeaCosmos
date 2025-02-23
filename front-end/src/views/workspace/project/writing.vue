<script setup>
import { onMounted, reactive, ref, h } from "vue";
import { MdPreview,MdEditor } from "md-editor-v3";
import "md-editor-v3/lib/preview.css";
import 'md-editor-v3/lib/style.css';
import { useThemeStore } from "@/stores/theme";
import Loader from "@/components/loader.vue";
import { get, post , postJSON } from "@/util/request";
import { BACKEND_DOMAIN } from "@/util/VARRIBLES";
import { message, Modal, InputNumber } from "ant-design-vue";
const themeStore = useThemeStore();
const chapter = JSON.parse(localStorage.getItem("chapter"));

// 添加 loading 状态
const loading = ref(true);

onMounted(() => {
  // 修改现有的 onMounted
  Promise.all([
    fetchCurrentChapterVersion(),
    fetchVersionHistory()
  ]).finally(() => {
    loading.value = false;
  });
});

const fetchCurrentChapterVersion = () => {
  if (chapter.version_id === 0) {
    return Promise.resolve();
  }
  return new Promise((resolve, reject) => {
    get(
      "/api/project/getCurrentChapterVersion",
      {
        chapter_id: chapter.ID,
      },
      (messager, data) => {
        options.currentVersion = data;
        resolve();
      },
      (messager, data) => {
        message.warning(messager);
        reject();
      },
      (messager, data) => {
        message.error(messager);
        reject();
      }
    );
  });
};

const fetchVersionHistory = () => {
  return new Promise((resolve, reject) => {
    get(
      "/api/project/getChapterVersions",
      {
        chapter_id: chapter.ID,
      },
      (messager, data) => {
        options.historyVersions = data;
        resolve();
      },
      (messager) => {
        message.warning(messager);
        reject();
      },
      (messager) => {
        message.error(messager);
        reject();
      }
    );
  });
};

const options = reactive({
  historyVersions: [],
  currentVersion: {},
  isGenerating: false,
  generatedContent: "",
  isEditing: false,
  editingContent: "",
});

const processContent = (content) => {
  return content.replace(/^```markdown\s*/g, '').replace(/\s*```$/g, '');
};

const generateNewVersion = (wordsCount) => {
  options.isGenerating = true;
  options.generatedContent = "";

  // 创建 WebSocket 连接
  const token = localStorage.getItem('authToken');
  const ws = new WebSocket(
    `ws://${BACKEND_DOMAIN.replace('http://', '')}ws/generateNewChapterVersionStream`
);
ws.onopen = (event) => {
  ws.send(JSON.stringify({
    token: token,
    chapterId: String(chapter.ID),
    wordsCount: String(wordsCount)
  }))
};

  ws.onmessage = (event) => {
    const response = JSON.parse(event.data);
    if (response.code === 500) {
      message.error(response.message);
      options.isGenerating = false;
      ws.close();
      return;
    }

    if (response.done) {
      options.isGenerating = false;
      message.success("生成成功！");
      ws.close();
      return;
    }

    options.generatedContent += response.content;
  };

  ws.onerror = (error) => {
    console.error('WebSocket error:', error);
    message.error("连接发生错误，请重试");
    options.isGenerating = false;
  };

  ws.onclose = () => {
    if (options.isGenerating) {
      options.isGenerating = false;
      message.error("连接已关闭，请重试");
    }
  };
};

const optimizeVersion = (suggestion) => {
  if (!options.generatedContent) {
    message.warning("请先生成或选择一个版本");
    return;
  }

  options.isGenerating = true;
  const oldContent = options.generatedContent;
  options.generatedContent = "";

  post(
    "/api/project/optimizeChapterVersion",
    {
      chapter_id: chapter.ID,
      current_content: oldContent,
      suggestion: suggestion,
      words_count: "3000",
    },
    (messager, data) => {
      options.generatedContent = processContent(data.choices[0].message.content);
      options.isGenerating = false;
      message.success("优化成功！");
    },
    (messager, data) => {
      options.isGenerating = false;
      options.generatedContent = oldContent;
      message.warning(messager);
    },
    (messager, data) => {
      options.isGenerating = false;
      options.generatedContent = oldContent;
      message.error(messager);
    }
  );
};

const showOptimizeDialog = () => {
  const suggestion = ref("");

  Modal.confirm({
    title: "优化建议",
    content: h("div", {}, [
      h("p", "请输入您的优化建议："),
      h("textarea", {
        value: suggestion.value,
        onInput: (e) => (suggestion.value = e.target.value),
        style: {
          width: "100%",
          minHeight: "100px",
          marginTop: "10px",
          padding: "8px",
        },
      }),
    ]),
    onOk: () => {
      if (!suggestion.value) {
        message.warning("请输入优化建议");
        return;
      }
      optimizeVersion(suggestion.value);
    },
  });
};

const confirmAdoptContent = () => {
  Modal.confirm({
    title: "采纳确认",
    content: "确定要采纳这个版本吗？采纳后将创建新的章节版本。",
    okText: "采纳",
    cancelText: "不采纳",
    onOk() {
      // 调用创建新版本的接口
      postJSON(
        "/api/project/createNewChapterVersion",
        {
          chapter_id: chapter.ID,
          content: options.generatedContent,
        },
        (messager, data) => {
          message.success("版本创建成功！");
          options.currentVersion = data;
          chapter.version_id = data.ID;
          localStorage.setItem("chapter", JSON.stringify(chapter));
          fetchCurrentChapterVersion();
          fetchVersionHistory();
        },
        (messager) => {
          message.warning(messager);
        },
        (messager) => {
          message.error(messager);
        }
      );
    },
    onCancel() {
      // 清空生成的内容
      options.generatedContent = "";
      message.info("已清空生成的内容");
    },
  });
};

const showWordsCountDialog = () => {
  const wordsCount = ref("3000");
  
  Modal.confirm({
    title: "设置字数",
    content: h("div", {}, [
      h("p", "请输入期望生成的字数（1500-8000）："),
      h("div", { class: "flex items-center gap-2" }, [
        h(InputNumber, {
          value: wordsCount.value,
          min: 1500,
          max: 8000,
          step: 500,
          "onUpdate:value": (val) => (wordsCount.value = val),
          style: {
            width: "150px",
          },
          addonAfter: "字",
        }),
      ]),
      h("p", { 
        class: "text-gray-500 text-sm mt-2",
      }, "注：实际生成字数可能会略有偏差")
    ]),
    onOk: () => {
      if (!wordsCount.value) {
        message.warning("请输入字数");
        return;
      }
      generateNewVersion(wordsCount.value);
    },
  });
};

const switchToVersion = (version) => {
  options.currentVersion = version;
  chapter.version_id = version.ID;
  localStorage.setItem("chapter", JSON.stringify(chapter));
};

const startEditing = () => {
  options.isEditing = true;
  options.editingContent = options.generatedContent || options.currentVersion.content || "";
};

const cancelEditing = () => {
  options.isEditing = false;
  options.editingContent = "";
};

const submitEditing = () => {
  postJSON(
    "/api/project/createNewChapterVersion",
    {
      chapter_id: chapter.ID,
      content: options.editingContent,
    },
    (messager, data) => {
      message.success("保存成功！");
      options.isEditing = false;
      options.editingContent = "";
      fetchVersionHistory();
      fetchCurrentChapterVersion();
    },
    (messager) => {
      message.warning(messager);
    },
    (messager) => {
      message.error(messager);
    }
  );
};
</script>

<template>
  <!-- 添加加载状态显示 -->
  <div v-if="loading" class="h-full w-full flex items-center justify-center">
    <SpinLoaderLarge />
  </div>

  <!-- 将原有内容包装在 v-else 中，并添加动画类 -->
  <div v-else class="h-full w-full grid grid-cols-[4fr,1fr] gap-4 p-4">
    <!-- 主编辑区域 -->
    <div class="flex flex-col gap-4">
      <!-- 章节信息 -->
      <div class="bg-white dark:bg-zinc-900 border theme-border rounded-xl p-4 animate__animated animate__fadeIn animate__delay-1s">
        <h1 class="text-xl font-bold bg-gradient-to-r from-blue-500 to-cyan-500 bg-clip-text text-transparent mb-2">
          {{ chapter.Title }}
        </h1>
        <p class="text-sm text-gray-600 dark:text-gray-400">{{ chapter.Description }}</p>
      </div>

      <!-- 工具栏 -->
      <div class="bg-white dark:bg-zinc-900 border theme-border rounded-xl p-4 animate__animated animate__fadeIn animate__delay-2s">
        <div class="flex items-center justify-between mb-4">
          <h2 class="text-lg font-semibold text-blue-500">创作工具</h2>
          <div class="flex items-center gap-2">
            <button
              @click="showWordsCountDialog"
              :disabled="options.isGenerating"
              class="flex items-center gap-2 px-4 py-2 bg-gradient-to-r from-blue-500 to-cyan-500 text-white rounded-lg hover:opacity-90 transition-all disabled:opacity-50 disabled:cursor-not-allowed"
            >
              <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" d="M9.813 15.904L9 18.75l-.813-2.846a4.5 4.5 0 00-3.09-3.09L2.25 12l2.846-.813a4.5 4.5 0 003.09-3.09L9 5.25l.813 2.846a4.5 4.5 0 003.09 3.09L15.75 12l-2.846.813a4.5 4.5 0 00-3.09 3.09z" />
              </svg>
              AI 创作
            </button>
            <button
              @click="showOptimizeDialog"
              :disabled="options.isGenerating || !options.generatedContent"
              class="flex items-center gap-2 px-4 py-2 bg-white dark:bg-zinc-800 border theme-border rounded-lg hover:border-blue-500 transition-all disabled:opacity-50 disabled:cursor-not-allowed"
            >
              <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 text-blue-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" d="M16.862 4.487l1.687-1.688a1.875 1.875 0 112.652 2.652L10.582 16.07a4.5 4.5 0 01-1.897 1.13L6 18l.8-2.685a4.5 4.5 0 011.13-1.897l8.932-8.931zm0 0L19.5 7.125M18 14v4.75A2.25 2.25 0 0115.75 21H5.25A2.25 2.25 0 013 18.75V8.25A2.25 2.25 0 015.25 6H10" />
              </svg>
              优化内容
            </button>
            <button
              @click="startEditing"
              :disabled="options.isGenerating"
              class="flex items-center gap-2 px-4 py-2 bg-white dark:bg-zinc-800 border theme-border rounded-lg hover:border-blue-500 transition-all disabled:opacity-50 disabled:cursor-not-allowed"
            >
              <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 text-blue-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
              </svg>
              手动编辑
            </button>
          </div>
        </div>
      </div>

      <!-- 内容编辑/预览区域 -->
      <div class="flex-grow bg-white dark:bg-zinc-900 border theme-border rounded-xl overflow-hidden animate__animated animate__fadeIn animate__delay-3s">
        <div class="relative h-full">
          <!-- 编辑模式 -->
          <div v-if="options.isEditing" class="h-full p-4">
            <div class="flex items-center justify-between mb-4">
              <h3 class="text-lg font-semibold text-blue-500">编辑内容</h3>
              <div class="flex items-center gap-2">
                <button
                  @click="submitEditing"
                  class="flex items-center gap-2 px-4 py-2 bg-green-500 text-white rounded-lg hover:bg-green-600 transition-colors"
                >
                  <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" d="M5 13l4 4L19 7" />
                  </svg>
                  保存
                </button>
                <button
                  @click="cancelEditing"
                  class="flex items-center gap-2 px-4 py-2 border border-gray-200 dark:border-gray-700 text-gray-600 dark:text-gray-300 rounded-lg hover:bg-gray-50 dark:hover:bg-gray-800 transition-colors"
                >
                  <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12" />
                  </svg>
                  取消
                </button>
              </div>
            </div>
            <MdEditor
              v-model="options.editingContent"
              :theme="themeStore.currentTheme"
              previewTheme="github"
              :preview="false"
              style="height: calc(100% - 60px)"
              showCodeRowNumber
              :toolbars="[
                'bold', 'italic', 'strikethrough', '|',
                'title', 'quote', 'unorderedList', 'orderedList', '|',
                'codeRow', 'code', 'link', 'image', '|',
                'preview'
              ]"
            />
          </div>

          <!-- 预览模式 (当不在编辑模式时显示) -->
          <template v-else>
            <!-- AI生成的内容 -->
            <div v-if="options.generatedContent || options.isGenerating" class="border-b theme-border p-4">
              <div class="flex items-center justify-between mb-4">
                <h3 class="text-lg font-semibold bg-gradient-to-r from-blue-500 to-cyan-500 bg-clip-text text-transparent flex items-center gap-2">
                  AI 创作内容
                  <span v-if="options.isGenerating" 
                        class="text-sm font-normal text-blue-500 animate-pulse">
                    正在创作...
                  </span>
                </h3>
                <div class="flex items-center gap-2">
                  <button
                    @click="confirmAdoptContent"
                    :disabled="options.isGenerating"
                    class="flex items-center gap-2 px-4 py-2 bg-green-500 text-white rounded-lg hover:bg-green-600 transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
                  >
                    <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                      <path stroke-linecap="round" stroke-linejoin="round" d="M4.5 12.75l6 6 9-13.5" />
                    </svg>
                    采纳
                  </button>
                  <button
                    @click="options.generatedContent = ''"
                    :disabled="options.isGenerating"
                    class="flex items-center gap-2 px-4 py-2 border border-red-200 dark:border-red-800 text-red-600 rounded-lg hover:bg-red-50 dark:hover:bg-red-900/20 transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
                  >
                    <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                      <path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12" />
                    </svg>
                    放弃
                  </button>
                </div>
              </div>
              <MdPreview
                style="background: transparent"
                :theme="themeStore.currentTheme"
                editorId="ai-preview"
                :modelValue="options.generatedContent"
              />
            </div>

            <!-- 当前版本内容 -->
            <div class="p-4" :class="{ 'h-[calc(100%-200px)]': options.generatedContent || options.isGenerating }">
              <h3 class="text-lg font-semibold text-gray-900 dark:text-gray-100 mb-4">当前版本</h3>
              <div v-if="chapter.version_id === 0" class="flex items-center justify-center h-64 text-gray-500">
                <div class="text-center">
                  <svg xmlns="http://www.w3.org/2000/svg" class="h-12 w-12 mx-auto mb-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" d="M12 6.042A8.967 8.967 0 006 3.75c-1.052 0-2.062.18-3 .512v14.25A8.987 8.987 0 016 18c2.305 0 4.408.867 6 2.292m0-14.25a8.966 8.966 0 016-2.292c1.052 0 2.062.18 3 .512v14.25A8.987 8.987 0 0018 18a8.967 8.967 0 00-6 2.292m0-14.25v14.25" />
                  </svg>
                  <p>此章节暂无内容</p>
                </div>
              </div>
              <MdPreview
                v-else
                style="background: transparent"
                :theme="themeStore.currentTheme"
                editorId="current-preview"
                :modelValue="options.currentVersion.content"
              />
            </div>
          </template>
        </div>
      </div>
    </div>

    <!-- 历史版本侧边栏 -->
    <div class="bg-white dark:bg-zinc-900 border theme-border rounded-xl p-4 animate__animated animate__fadeIn animate__delay-4s">
      <h2 class="text-lg font-semibold text-blue-500 mb-4">历史版本</h2>
      <div class="space-y-4 max-h-[calc(100vh-10rem)] overflow-y-auto">
        <div v-if="options.historyVersions.length === 0" 
             class="flex items-center justify-center h-32 text-gray-500">
          <p>暂无历史版本</p>
        </div>
        
        <div v-for="version in options.historyVersions"
             :key="version.ID"
             @click="switchToVersion(version)"
             class="p-4 rounded-xl cursor-pointer transition-all duration-200"
             :class="[
               version.ID === chapter.version_id 
                 ? 'bg-blue-50 dark:bg-blue-900/20 border-blue-200 dark:border-blue-800' 
                 : 'hover:bg-gray-50 dark:hover:bg-gray-800/50 border-transparent',
               'border'
             ]"
        >
          <div class="flex items-center gap-3 mb-2">
            <img :src="BACKEND_DOMAIN + version.user.avatar" 
                 class="w-8 h-8 rounded-full object-cover"
                 :alt="version.user.nickname || version.user.username">
            <div>
              <p class="text-sm font-medium text-gray-900 dark:text-gray-100">
                {{ version.user.nickname || version.user.username }}
              </p>
              <p class="text-xs text-gray-500">
                {{ new Date(version.CreatedAt).toLocaleString() }}
              </p>
            </div>
          </div>
          <p class="text-sm text-gray-600 dark:text-gray-400 line-clamp-2">
            {{ version.content.substring(0, 100) + '...' }}
          </p>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
/* 添加动画相关样式 */
.animate__animated {
  animation-duration: 0.5s;
}

.animate__delay-1s {
  animation-delay: 0.1s;
}

.animate__delay-2s {
  animation-delay: 0.2s;
}

.animate__delay-3s {
  animation-delay: 0.3s;
}

.animate__delay-4s {
  animation-delay: 0.4s;
}

/* 为卡片添加淡入动画 */
@keyframes fadeIn {
  from {
    opacity: 0;
    transform: translateY(10px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}
</style>
