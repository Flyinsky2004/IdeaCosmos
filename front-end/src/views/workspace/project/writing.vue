<script setup>
import { onMounted, reactive, ref, h, watch } from "vue";
import { MdPreview, MdEditor } from "md-editor-v3";
import "md-editor-v3/lib/preview.css";
import "md-editor-v3/lib/style.css";
import { useThemeStore } from "@/stores/theme";
import Loader from "@/components/loader.vue";
import { get, post, postJSON } from "@/util/request";
import { BACKEND_DOMAIN, FRONTEND_DOMAIN, imagePrefix } from "@/util/VARRIBLES";
import { message, Modal, InputNumber } from "ant-design-vue";
const themeStore = useThemeStore();
const chapter = JSON.parse(localStorage.getItem("chapter"));
import SpinLoaderLarge from "@/components/spinLoaderLarge.vue";
import { useUserStore } from "@/stores/user";
// 添加 loading 状态
const loading = ref(true);
const userId = localStorage.getItem("userId");
const userStore = useUserStore()

onMounted(() => {
  // 修改现有的 onMounted
  Promise.all([fetchCurrentChapterVersion(), fetchVersionHistory()]).finally(
    () => {
      loading.value = false;
    }
  );
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
        fetchVersionComments(data.ID);
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
  aiContentCollapsed: false,
  currentVersionCollapsed: false,
});

const processContent = (content) => {
  return content.replace(/^```markdown\s*/g, "").replace(/\s*```$/g, "");
};

const generateNewVersion = (wordsCount) => {
  options.isGenerating = true;
  options.generatedContent = "";

  // 创建 WebSocket 连接
  const token = localStorage.getItem("authToken");
  const baseUrl = BACKEND_DOMAIN.replace(/^http/, 'ws').replace(/\/$/, '')
  const ws = new WebSocket(
    `${baseUrl}/ws/generateNewChapterVersionStream`
  );
  ws.onopen = (event) => {
    ws.send(
      JSON.stringify({
        token: token,
        chapterId: String(chapter.ID),
        wordsCount: String(wordsCount),
      })
    );
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
    console.error("WebSocket error:", error);
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
  if (!options.generatedContent && !options.currentVersion.content) {
    message.warning("请先生成或选择一个版本");
    return;
  }

  options.isGenerating = true;
  const currentContent = options.generatedContent || options.currentVersion.content;
  options.generatedContent = "";

  // 创建 WebSocket 连接
  const token = localStorage.getItem("authToken");
  const ws = new WebSocket(
    `ws://${BACKEND_DOMAIN.replace("http://", "")}ws/modifyChapterVersionStream`
  );

  ws.onopen = (event) => {
    ws.send(
      JSON.stringify({
        token: token,
        chapterId: String(chapter.ID),
        currentContent: currentContent,
        modifyPreference: suggestion,
      })
    );
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
      message.success("优化成功！");
      ws.close();
      return;
    }

    options.generatedContent += response.content;
  };

  ws.onerror = (error) => {
    console.error("WebSocket error:", error);
    message.error("连接发生错误，请重试");
    options.isGenerating = false;
    options.generatedContent = currentContent;
  };

  ws.onclose = () => {
    if (options.isGenerating) {
      options.isGenerating = false;
      options.generatedContent = currentContent;
      message.error("连接已关闭，请重试");
    }
  };
};

const showOptimizeDialog = () => {
  const suggestion = ref("");

  Modal.confirm({
    title: "优化建议",
    class: "dark:bg-zinc-800",
    content: h("div", {}, [
      h("p", { class: "text-gray-700 dark:text-gray-200" }, "请输入您的优化建议："),
      h("textarea", {
        value: suggestion.value,
        onInput: (e) => (suggestion.value = e.target.value),
        class: `
          w-full
          min-h-[100px]
          mt-3
          p-3
          rounded-lg
          border
          border-gray-200
          dark:border-gray-700
          bg-white
          dark:bg-zinc-900
          text-gray-900
          dark:text-gray-100
          placeholder-gray-400
          dark:placeholder-gray-500
          focus:border-blue-500
          dark:focus:border-blue-400
          focus:ring-1
          focus:ring-blue-500
          dark:focus:ring-blue-400
          focus:outline-none
          transition-colors
          duration-200
        `,
        placeholder: "请详细描述您希望如何优化当前内容...",
      }),
    ]),
    okText: "开始优化",
    cancelText: "取消",
    okButtonProps: {
      class: "bg-blue-500 hover:bg-blue-600 border-blue-500 hover:border-blue-600",
    },
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
      h(
        "p",
        {
          class: "text-gray-500 text-sm mt-2",
        },
        "注：实际生成字数可能会略有偏差"
      ),
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
  fetchVersionComments(version.ID);
};

const startEditing = () => {
  options.isEditing = true;
  options.editingContent =
    options.generatedContent || options.currentVersion.content || "";
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

const toggleContentCollapse = (type) => {
  if (type === 'ai') {
    options.aiContentCollapsed = !options.aiContentCollapsed;
  } else if (type === 'current') {
    options.currentVersionCollapsed = !options.currentVersionCollapsed;
  }
};

// 评论相关状态
const commentState = reactive({
  comments: [],
  loadingComments: false,
  commentContent: "",
  submittingComment: false,
});

// 获取版本评论
const fetchVersionComments = (versionId) => {
  if (!versionId) return;
  
  commentState.loadingComments = true;
  get(
    "/api/project/creator/comment/list",
    {
      version_id: versionId,
    },
    (msg, data) => {
      commentState.comments = data;
      commentState.loadingComments = false;
    },
    (msg) => {
      message.warning(msg);
      commentState.loadingComments = false;
    },
    (msg) => {
      message.error(msg);
      commentState.loadingComments = false;
    }
  );
};

// 添加评论
const submitComment = () => {
  if (!commentState.commentContent.trim()) {
    message.warning("评论内容不能为空");
    return;
  }
  
  if (!options.currentVersion.ID) {
    message.warning("请先选择一个版本");
    return;
  }
  
  commentState.submittingComment = true;
  
  postJSON(
    "/api/project/creator/comment/add",
    {
      version_id: options.currentVersion.ID,
      content: commentState.commentContent,
    },
    (msg, data) => {
      message.success(msg);
      commentState.commentContent = "";
      fetchVersionComments(options.currentVersion.ID);
      commentState.submittingComment = false;
    },
    (msg) => {
      message.warning(msg);
      commentState.submittingComment = false;
    },
    (msg) => {
      message.error(msg);
      commentState.submittingComment = false;
    }
  );
};

// 删除评论
const deleteComment = (commentId) => {
  Modal.confirm({
    title: "删除确认",
    content: "确定要删除这条评论吗？",
    okText: "确定",
    cancelText: "取消",
    onOk() {
      postJSON(
        "/api/project/creator/comment/delete",
        {
          comment_id: commentId,
        },
        (msg) => {
          message.success(msg);
          fetchVersionComments(options.currentVersion.ID);
        },
        (msg) => {
          message.warning(msg);
        },
        (msg) => {
          message.error(msg);
        }
      );
    },
  });
};

// 监听当前版本变化，加载评论
watch(
  () => options.currentVersion,
  (newVal) => {
    if (newVal && newVal.ID) {
      fetchVersionComments(newVal.ID);
    }
  },
  { deep: true }
);
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
      <div
        class="bg-white dark:bg-zinc-900 border theme-border rounded-xl p-4 animate__animated animate__fadeIn animate__delay-1s"
      >
        <h1
          class="text-xl font-bold bg-gradient-to-r from-blue-500 to-cyan-500 bg-clip-text text-transparent mb-2"
        >
          {{ chapter.Title }}
        </h1>
        <p class="text-sm text-gray-600 dark:text-gray-400">
          {{ chapter.Description }}
        </p>
      </div>

      <!-- 工具栏 -->
      <div
        class="bg-white dark:bg-zinc-900 border theme-border rounded-xl p-4 animate__animated animate__fadeIn animate__delay-2s"
      >
        <div class="flex items-center justify-between mb-4">
          <h2 class="text-lg font-semibold text-blue-500">创作工具</h2>
          <div class="flex items-center gap-2">
            <button
              @click="showWordsCountDialog"
              :disabled="options.isGenerating"
              class="flex items-center gap-2 px-4 py-2 bg-gradient-to-r from-blue-500 to-cyan-500 text-white rounded-lg hover:opacity-90 transition-all disabled:opacity-50 disabled:cursor-not-allowed"
            >
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
                  d="M9.813 15.904L9 18.75l-.813-2.846a4.5 4.5 0 00-3.09-3.09L2.25 12l2.846-.813a4.5 4.5 0 003.09-3.09L9 5.25l.813 2.846a4.5 4.5 0 003.09 3.09L15.75 12l-2.846.813a4.5 4.5 0 00-3.09 3.09z"
                />
              </svg>
              AI 创作
            </button>
            <button
              @click="showOptimizeDialog"
              :disabled="options.isGenerating"
              class="flex items-center gap-2 px-4 py-2 bg-white dark:bg-zinc-800 border theme-border rounded-lg hover:border-blue-500 transition-all disabled:opacity-50 disabled:cursor-not-allowed"
            >
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
                  d="M16.862 4.487l1.687-1.688a1.875 1.875 0 112.652 2.652L10.582 16.07a4.5 4.5 0 01-1.897 1.13L6 18l.8-2.685a4.5 4.5 0 011.13-1.897l8.932-8.931zm0 0L19.5 7.125M18 14v4.75A2.25 2.25 0 0115.75 21H5.25A2.25 2.25 0 013 18.75V8.25A2.25 2.25 0 015.25 6H10"
                />
              </svg>
              优化内容
            </button>
            <button
              @click="startEditing"
              :disabled="options.isGenerating"
              class="flex items-center gap-2 px-4 py-2 bg-white dark:bg-zinc-800 border theme-border rounded-lg hover:border-blue-500 transition-all disabled:opacity-50 disabled:cursor-not-allowed"
            >
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
                  d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"
                />
              </svg>
              手动编辑
            </button>
          </div>
        </div>
      </div>

      <!-- 内容编辑/预览区域 -->
      <div
        class="flex-grow bg-white dark:bg-zinc-900 border theme-border rounded-xl overflow-hidden animate__animated animate__fadeIn animate__delay-3s"
      >
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
                      d="M5 13l4 4L19 7"
                    />
                  </svg>
                  保存
                </button>
                <button
                  @click="cancelEditing"
                  class="flex items-center gap-2 px-4 py-2 border border-gray-200 dark:border-gray-700 text-gray-600 dark:text-gray-300 rounded-lg hover:bg-gray-50 dark:hover:bg-gray-800 transition-colors"
                >
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
                      d="M6 18L18 6M6 6l12 12"
                    />
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
                'bold',
                'italic',
                'strikethrough',
                '|',
                'title',
                'quote',
                'unorderedList',
                'orderedList',
                '|',
                'codeRow',
                'code',
                'link',
                'image',
                '|',
                'preview',
              ]"
            />
          </div>

          <!-- 预览模式 (当不在编辑模式时显示) -->
          <template v-else>
            <!-- AI生成的内容 -->
            <div
              v-if="options.generatedContent || options.isGenerating"
              class="border-b theme-border p-4 select-text"
            >
              <div class="flex items-center justify-between mb-4">
                <h3
                  class="text-lg font-semibold bg-gradient-to-r from-blue-500 to-cyan-500 bg-clip-text text-transparent flex items-center gap-2"
                >
                  AI 创作内容
                  <span
                    v-if="options.isGenerating"
                    class="text-sm font-normal text-blue-500 animate-pulse"
                  >
                    正在创作...
                  </span>
                </h3>
                <div class="flex items-center gap-2">
                  <button
                    @click="toggleContentCollapse('ai')"
                    class="flex items-center gap-1 px-2 py-1 text-sm text-gray-500 hover:text-blue-500 transition-colors"
                  >
                    <svg
                      xmlns="http://www.w3.org/2000/svg"
                      class="h-4 w-4"
                      fill="none"
                      viewBox="0 0 24 24"
                      stroke="currentColor"
                      :class="{ 'rotate-180': !options.aiContentCollapsed }"
                    >
                      <path
                        stroke-linecap="round"
                        stroke-linejoin="round"
                        stroke-width="2"
                        d="M19 9l-7 7-7-7"
                      />
                    </svg>
                    {{ options.aiContentCollapsed ? '展开' : '折叠' }}
                  </button>
                  <button
                    @click="confirmAdoptContent"
                    :disabled="options.isGenerating"
                    class="flex items-center gap-2 px-4 py-2 bg-green-500 text-white rounded-lg hover:bg-green-600 transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
                  >
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
                        d="M4.5 12.75l6 6 9-13.5"
                      />
                    </svg>
                    采纳
                  </button>
                  <button
                    @click="options.generatedContent = ''"
                    :disabled="options.isGenerating"
                    class="flex items-center gap-2 px-4 py-2 border border-red-200 dark:border-red-800 text-red-600 rounded-lg hover:bg-red-50 dark:hover:bg-red-900/20 transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
                  >
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
                        d="M6 18L18 6M6 6l12 12"
                      />
                    </svg>
                    放弃
                  </button>
                </div>
              </div>
              <div 
                class="content-wrapper overflow-hidden transition-all duration-300" 
                :class="{ 'collapsed': options.aiContentCollapsed }"
              >
                <MdPreview
                  style="background: transparent"
                  :theme="themeStore.currentTheme"
                  editorId="ai-preview"
                  :modelValue="options.generatedContent"
                />
              </div>
            </div>

            <!-- 当前版本内容 -->
            <div
              class="p-4 select-text"
              :class="{
                'h-[calc(100%-100px)]':
                  options.generatedContent || options.isGenerating,
              }"
            >
              <div class="flex items-center justify-between mb-4">
                <h3
                  class="text-lg font-semibold text-gray-900 dark:text-gray-100"
                >
                  当前版本
                </h3>
                <button
                  v-if="chapter.version_id !== 0"
                  @click="toggleContentCollapse('current')"
                  class="flex items-center gap-1 px-2 py-1 text-sm text-gray-500 hover:text-blue-500 transition-colors"
                >
                  <svg
                    xmlns="http://www.w3.org/2000/svg"
                    class="h-4 w-4"
                    fill="none"
                    viewBox="0 0 24 24"
                    stroke="currentColor"
                    :class="{ 'rotate-180': !options.currentVersionCollapsed }"
                  >
                    <path
                      stroke-linecap="round"
                      stroke-linejoin="round"
                      stroke-width="2"
                      d="M19 9l-7 7-7-7"
                    />
                  </svg>
                  {{ options.currentVersionCollapsed ? '展开' : '折叠' }}
                </button>
              </div>
              <div
                v-if="chapter.version_id === 0"
                class="flex items-center justify-center h-64 text-gray-500"
              >
                <div class="text-center">
                  <svg
                    xmlns="http://www.w3.org/2000/svg"
                    class="h-12 w-12 mx-auto mb-2"
                    fill="none"
                    viewBox="0 0 24 24"
                    stroke="currentColor"
                  >
                    <path
                      stroke-linecap="round"
                      stroke-linejoin="round"
                      d="M12 6.042A8.967 8.967 0 006 3.75c-1.052 0-2.062.18-3 .512v14.25A8.987 8.987 0 016 18c2.305 0 4.408.867 6 2.292m0-14.25a8.966 8.966 0 016-2.292c1.052 0 2.062.18 3 .512v14.25A8.987 8.987 0 0018 18a8.967 8.967 0 00-6 2.292m0-14.25v14.25"
                    />
                  </svg>
                  <p>此章节暂无内容</p>
                </div>
              </div>
              <div
                v-else
                class="content-wrapper overflow-hidden transition-all duration-300"
                :class="{ 'collapsed': options.currentVersionCollapsed }"
              >
                <MdPreview
                  style="background: transparent"
                  :theme="themeStore.currentTheme"
                  editorId="current-preview"
                  :modelValue="options.currentVersion.content"
                />
              </div>
            </div>

            <!-- 评论区域 -->
            <div v-if="chapter.version_id !== 0" class="bg-white dark:bg-zinc-900 border theme-border rounded-xl p-4 mt-4 animate__animated animate__fadeIn">
              <h3 class="text-lg font-semibold text-blue-500 mb-4">版本讨论</h3>
              
              <!-- 评论输入框 -->
              <div class="mb-6">
                <div class="flex items-start gap-3">
                  <div class="w-10 h-10 overflow-hidden rounded-full flex-shrink-0">
                    <img 
                      :src="BACKEND_DOMAIN + (userStore.user?.avatar || 'uploads/default-avatar.png')" 
                      class="w-full h-full object-cover" 
                      alt="用户头像"
                    />
                  </div>
                  <div class="flex-grow">
                    <textarea 
                      v-model="commentState.commentContent"
                      class="w-full min-h-[80px] p-3 rounded-lg border theme-border bg-white dark:bg-zinc-800 text-gray-900 dark:text-gray-100 focus:outline-none focus:ring-1 focus:ring-blue-500 transition-colors resize-none"
                      placeholder="在这里添加你对这个版本的想法和建议..."
                      :disabled="commentState.submittingComment"
                    ></textarea>
                    <div class="flex justify-end mt-2">
                      <button 
                        @click="submitComment"
                        :disabled="commentState.submittingComment"
                        class="flex items-center gap-2 px-4 py-2 bg-blue-500 text-white rounded-lg hover:bg-blue-600 transition-all disabled:opacity-50 disabled:cursor-not-allowed"
                      >
                        <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 19l9 2-9-18-9 18 9-2zm0 0v-8" />
                        </svg>
                        发送
                      </button>
                    </div>
                  </div>
                </div>
              </div>
              
              <!-- 评论列表 -->
              <div class="space-y-6">
                <!-- 加载状态 -->
                <div v-if="commentState.loadingComments" class="flex justify-center py-4">
                  <Loader />
                </div>
                
                <!-- 没有评论 -->
                <div v-else-if="commentState.comments.length === 0" class="text-center py-10 text-gray-500 dark:text-gray-400">
                  <svg xmlns="http://www.w3.org/2000/svg" class="h-12 w-12 mx-auto mb-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z" />
                  </svg>
                  <p>还没有评论，快来发表第一条评论吧！</p>
                </div>
                
                <!-- 评论列表 -->
                <div v-else class="space-y-4">
                  <div v-for="comment in commentState.comments" :key="comment.ID" class="border-b theme-border last:border-0 pb-4 last:pb-0">
                    <div class="flex items-start gap-3">
                      <div class="w-10 h-10 overflow-hidden rounded-full flex-shrink-0">
                        <img 
                          :src="BACKEND_DOMAIN + (comment.user?.avatar || '/avatar/default.jpg')" 
                          class="w-full h-full object-cover" 
                          alt="用户头像"
                        />
                      </div>
                      <div class="flex-grow">
                        <div class="flex items-center justify-between">
                          <div>
                            <span class="font-medium text-gray-900 dark:text-gray-100">
                              {{ comment.user?.nickname || comment.user?.username || '未知用户' }}
                            </span>
                            <span class="text-xs text-gray-500 ml-2">
                              {{ new Date(comment.CreatedAt).toLocaleString() }}
                            </span>
                          </div>
                          <button 
                            v-if="comment.UserId === parseInt(userId)"
                            @click="deleteComment(comment.ID)"
                            class="text-gray-500 hover:text-red-500 transition-colors"
                          >
                            <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                            </svg>
                          </button>
                        </div>
                        <p class="mt-2 text-gray-700 dark:text-gray-300">{{ comment.Content }}</p>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </template>
        </div>
      </div>
    </div>

    <!-- 历史版本侧边栏 -->
    <div
      class="bg-white dark:bg-zinc-900 border theme-border rounded-xl p-4 animate__animated animate__fadeIn animate__delay-4s"
    >
      <h2 class="text-lg font-semibold text-blue-500 mb-4">历史版本</h2>
      <div class="space-y-4 max-h-[calc(100vh-10rem)] overflow-y-auto">
        <div
          v-if="options.historyVersions.length === 0"
          class="flex items-center justify-center h-32 text-gray-500"
        >
          <p>暂无历史版本</p>
        </div>

        <div
          v-for="version in options.historyVersions"
          :key="version.ID"
          @click="switchToVersion(version)"
          class="p-4 rounded-xl cursor-pointer transition-all duration-200"
          :class="[
            version.ID === chapter.version_id
              ? 'bg-blue-50 dark:bg-blue-900/20 border-blue-200 dark:border-blue-800'
              : 'hover:bg-gray-50 dark:hover:bg-gray-800/50 border-transparent',
            'border',
          ]"
        >
          <div class="flex items-center gap-3 mb-2">
            <img
              :src="BACKEND_DOMAIN + version.user.avatar"
              class="w-8 h-8 rounded-full object-cover"
              :alt="version.user.nickname || version.user.username"
            />
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
            {{ version.content.substring(0, 100) + "..." }}
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

/* 折叠功能样式 */
.content-wrapper {
  max-height: 2000px;
}

.content-wrapper.collapsed {
  max-height: 300px;
  position: relative;
  mask-image: linear-gradient(to bottom, black 80%, transparent 100%);
  -webkit-mask-image: linear-gradient(to bottom, black 80%, transparent 100%);
}

.rotate-180 {
  transform: rotate(180deg);
}
</style>
