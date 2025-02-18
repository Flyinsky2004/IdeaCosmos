<script setup>
import { onMounted, reactive, ref, h } from "vue";
import { MdPreview } from "md-editor-v3";
import "md-editor-v3/lib/preview.css";
import { useThemeStore } from "@/stores/theme";
import Loader from "@/components/loader.vue";
import { get, post , postJSON } from "@/util/request";
import { BACKEND_DOMAIN } from "@/util/VARRIBLES";
import { message, Modal, InputNumber } from "ant-design-vue";
const themeStore = useThemeStore();
const chapter = JSON.parse(localStorage.getItem("chapter"));
onMounted(() => {
  fetchCurrentChapterVersion();
  fetchVersionHistory();
});
const fetchCurrentChapterVersion = () => {
  if (chapter.version_id === 0) return;
  get(
    "/api/project/getCurrentChapterVersion",
    {
      chapter_id: chapter.ID,
    },
    (messager, data) => {
      options.currentVersion = data;
    },
    (messager, data) => {
      message.warning(messager);
    },
    (messager, data) => {
      message.error(messager);
    }
  );
};
const fetchVersionHistory = () => {
  get(
    "/api/project/getChapterVersions",
    {
      chapter_id: chapter.ID,
    },
    (messager, data) => {
      options.historyVersions = data;
    },
    (messager) => {
      message.warning(messager);
    },
    (messager) => {
      message.error(messager);
    }
  );
};
const options = reactive({
  historyVersions: [],
  currentVersion: {},
  isGenerating: false,
  generatedContent: "",
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
      generateNewVersion(wordsCount.value * 5);
    },
  });
};

const switchToVersion = (version) => {
  options.currentVersion = version;
  chapter.version_id = version.ID;
  localStorage.setItem("chapter", JSON.stringify(chapter));
};
</script>

<template>
  <div class="h-full w-full grid grid-cols-[4fr,1fr] gap-2">
    <div class="flex flex-col gap-2">
      <div class="h-fit border theme-border p-4 rounded-xl flex flex-col gap-2">
        <span class="text-blue-500 font-bold text-xl">工具栏</span>
        <div class="w-full flex flex-nowrap gap-2">
          <button
            class="color-mixed-button"
            @click="showWordsCountDialog"
            :disabled="options.isGenerating"
          >
            剧匠AI创作
          </button>
          <button
            class="basic-info-button cursor-pointer"
            @click="showOptimizeDialog"
            :disabled="options.isGenerating || !options.generatedContent"
          >
            优化内容
          </button>
          <button class="basic-prinary-button">编辑模式</button>
          <button class="basic-success-button">应用当前版本</button>
        </div>
      </div>
      <div
        class="flex flex-grow flex-col dark:bg-indigo-950/10 border theme-border rounded-xl"
      >
        <div class="flex flex-col flex-grow relative">
          <span
            v-if="options.generatedContent !== ''"
            class="text-blue-500 font-bold text-xl mt-4 ml-4"
            >剧匠AI版本</span
          >
          <MdPreview
            v-if="options.generatedContent !== ''"
            style="background: transparent"
            :theme="themeStore.currentTheme"
            editorId="preview-only"
            :modelValue="options.generatedContent"
          />
          <span class="text-blue-500 font-bold text-xl mt-4 ml-4"
            >当前版本</span
          >
          <div
            class="mx-auto my-auto"
            v-if="chapter.version_id === 0"
          >
            <span class="text-gray-600">此篇章还没有内容哦</span>
          </div>
          <MdPreview
            v-else-if="chapter.version_id !== 0"
            style="background: transparent"
            :theme="themeStore.currentTheme"
            editorId="preview-only"
            :modelValue="options.currentVersion.content"
          />

          <!-- 添加确认按钮区域 -->
          <div
            v-if="options.generatedContent"
            class="sticky bottom-0 mx-auto px-8 bg-white dark:bg-indigo-950/90 w-fit rounded-xl border theme-border p-4 border-t theme-border flex justify-center gap-4"
          >
            <button
              class="basic-success-button px-8 flex flex-nowrap"
              @click="confirmAdoptContent"
            >
              <svg
                xmlns="http://www.w3.org/2000/svg"
                viewBox="0 0 20 20"
                fill="currentColor"
                class="size-4"
              >
                <path
                  d="M1 8.25a1.25 1.25 0 1 1 2.5 0v7.5a1.25 1.25 0 1 1-2.5 0v-7.5ZM11 3V1.7c0-.268.14-.526.395-.607A2 2 0 0 1 14 3c0 .995-.182 1.948-.514 2.826-.204.54.166 1.174.744 1.174h2.52c1.243 0 2.261 1.01 2.146 2.247a23.864 23.864 0 0 1-1.341 5.974C17.153 16.323 16.072 17 14.9 17h-3.192a3 3 0 0 1-1.341-.317l-2.734-1.366A3 3 0 0 0 6.292 15H5V8h.963c.685 0 1.258-.483 1.612-1.068a4.011 4.011 0 0 1 2.166-1.73c.432-.143.853-.386 1.011-.814.16-.432.248-.9.248-1.388Z"
                />
              </svg>

              采纳这个版本
            </button>
            <button
              class="basic-error-button flex flex-nowrap"
              @click="options.generatedContent = ''"
            >
              <svg
                xmlns="http://www.w3.org/2000/svg"
                viewBox="0 0 20 20"
                fill="currentColor"
                class="size-4"
              >
                <path
                  d="M18.905 12.75a1.25 1.25 0 1 1-2.5 0v-7.5a1.25 1.25 0 0 1 2.5 0v7.5ZM8.905 17v1.3c0 .268-.14.526-.395.607A2 2 0 0 1 5.905 17c0-.995.182-1.948.514-2.826.204-.54-.166-1.174-.744-1.174h-2.52c-1.243 0-2.261-1.01-2.146-2.247.193-2.08.651-4.082 1.341-5.974C2.752 3.678 3.833 3 5.005 3h3.192a3 3 0 0 1 1.341.317l2.734 1.366A3 3 0 0 0 13.613 5h1.292v7h-.963c-.685 0-1.258.482-1.612 1.068a4.01 4.01 0 0 1-2.166 1.73c-.432.143-.853.386-1.011.814-.16.432-.248.9-.248 1.388Z"
                />
              </svg>
              放弃这个版本
            </button>
          </div>
        </div>
      </div>
    </div>
    <div class="border theme-border p-2 rounded-xl flex flex-col gap-2">
      <h1 class="text-xl font-bold text-blue-500">历史版本</h1>
      <div class="flex flex-col gap-2 overflow-y-auto">
        <div v-if="options.historyVersions.length === 0" class="text-center py-4">
          <span class="text-gray-600 dark:text-gray-400">还没有历史版本哦</span>
        </div>
        <div
          v-for="version in options.historyVersions"
          :key="version.ID"
          @click="switchToVersion(version)"
          class="p-4 border theme-border rounded-xl cursor-pointer hover:bg-gray-50 dark:hover:bg-gray-800 transition-colors duration-200"
          :class="{'bg-blue-50 dark:bg-blue-900/20': version.ID === chapter.version_id}"
        >
          <div class="flex items-center gap-2">
            <img
              :src="BACKEND_DOMAIN + version.user.avatar || ''"
              alt="avatar"
              class="w-8 h-8 rounded-full object-cover"
            />
            <div class="flex flex-col">
              <span class="text-sm font-medium dark:text-gray-200">
                {{ version.user.nickname || version.user.username }}
              </span>
              <span class="text-xs text-gray-500 dark:text-gray-400">
                {{ new Date(version.CreatedAt).toLocaleString() }}
              </span>
            </div>
          </div>
          <div class="mt-2">
            <p class="text-sm text-gray-600 dark:text-gray-300 line-clamp-2">
              {{ version.content.substring(0, 100) + '...' }}
            </p>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped></style>
