<script setup>
import router from "@/router/index.js";
import Loader from "@/components/loader.vue";
import {onMounted, reactive} from "vue";
import {get, post, postJSON} from "@/util/request.js";
import {message} from "ant-design-vue";
import {washJSONStr} from "@/util/common.js";

const project = JSON.parse(localStorage.getItem("project"));
const options = reactive({
  isChapterGenerating: false,
  generatedChapters: [],
  chapters: []
})

onMounted(() => {
  fetchChapters()
})

const fetchChapters = () => {
  get('/api/project/getAllChapters', {
    project_id: project.ID
  }, (messageer, data) => {
    options.chapters = data
  }, (messageer, data) => {
    message.warning(messageer)
  }, (messageer, data) => {
    message.error(messageer)
  })
}

const generateChapters = () => {
  options.isChapterGenerating = true;
  post('/api/project/generateChapters', {
    project_id: project.ID
  }, (messageer, data) => {
    const raw = data.choices[0].message.content
    options.generatedChapters = JSON.parse(washJSONStr(raw))
    message.success("篇章生成成功！")
    options.isChapterGenerating = false
  }, (messageer, data) => {
    message.warning(messageer)
    options.isChapterGenerating = false
  }, (messageer, data) => {
    message.error(messageer)
    options.isChapterGenerating = false
  })
}

const acceptChapter = () => {
  for (let i = 0; i < options.generatedChapters.length; i++) options.generatedChapters[i].project_id = project.ID;
  postJSON('/api/project/createChapterMulti',
      options.generatedChapters
      , (messageer, data) => {
        fetchChapters()
        message.success(data)
        options.generatedChapters = []
      }, (messageer, data) => {
        message.warning(messageer)
      }, (messageer, data) => {
        message.error(messageer)
      })
}

const rejectChapter = () => {
  options.generatedChapters = []
  message.info("已清空生成章节目录。")
}

const goToWriting = (chapter) => {
  localStorage.setItem("chapter",JSON.stringify(chapter))
  router.push('/workspace/editProject/writing')
}
</script>

<template>
  <div class="w-full grid grid-cols-5 gap-2">
<!--    <div class="border border-dashed rounded-xl dark:border-[rgb(118,118,118)] outline-[1px] min-h-56 place-items-center place-content-center text-gray-400-->
<!--hover:bg-gray-100/50 active:bg-gray-100/90 dark:hover:dark:bg-gray-950/10 dark:active:bg-gray-900/10 cursor-pointer"-->
<!--         @click="router.push('/workspace/newProject')">-->
<!--      <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5"-->
<!--           stroke="currentColor" class="size-16">-->
<!--        <path stroke-linecap="round" stroke-linejoin="round" d="M12 9v6m3-3H9m12 0a9 9 0 1 1-18 0 9 9 0 0 1 18 0Z"/>-->
<!--      </svg>-->
<!--    </div>-->
    <div @click="generateChapters"
         class="border border-dashed rounded-xl dark:border-[rgb(118,118,118)] outline-[1px] min-h-56 place-items-center place-content-center
hover:bg-gray-100/50 active:bg-gray-100/90 dark:hover:dark:bg-gray-950/10 dark:active:bg-gray-900/10 cursor-pointer">
      <div v-if="!options.isChapterGenerating" class="flex flex-col gap-4">
        <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="currentColor"
             class="size-7 dark:text-purple-400 mx-auto">
          <path fill-rule="evenodd"
                d="M9 4.5a.75.75 0 0 1 .721.544l.813 2.846a3.75 3.75 0 0 0 2.576 2.576l2.846.813a.75.75 0 0 1 0 1.442l-2.846.813a3.75 3.75 0 0 0-2.576 2.576l-.813 2.846a.75.75 0 0 1-1.442 0l-.813-2.846a3.75 3.75 0 0 0-2.576-2.576l-2.846-.813a.75.75 0 0 1 0-1.442l2.846-.813A3.75 3.75 0 0 0 7.466 7.89l.813-2.846A.75.75 0 0 1 9 4.5ZM18 1.5a.75.75 0 0 1 .728.568l.258 1.036c.236.94.97 1.674 1.91 1.91l1.036.258a.75.75 0 0 1 0 1.456l-1.036.258c-.94.236-1.674.97-1.91 1.91l-.258 1.036a.75.75 0 0 1-1.456 0l-.258-1.036a2.625 2.625 0 0 0-1.91-1.91l-1.036-.258a.75.75 0 0 1 0-1.456l1.036-.258a2.625 2.625 0 0 0 1.91-1.91l.258-1.036A.75.75 0 0 1 18 1.5ZM16.5 15a.75.75 0 0 1 .712.513l.394 1.183c.15.447.5.799.948.948l1.183.395a.75.75 0 0 1 0 1.422l-1.183.395c-.447.15-.799.5-.948.948l-.395 1.183a.75.75 0 0 1-1.422 0l-.395-1.183a1.5 1.5 0 0 0-.948-.948l-1.183-.395a.75.75 0 0 1 0-1.422l1.183-.395c.447-.15.799-.5.948-.948l.395-1.183A.75.75 0 0 1 16.5 15Z"
                clip-rule="evenodd"/>
        </svg>
        <span>AI剧匠智能生成</span>
      </div>
      <div v-else class="text-center place-items-center">
        <loader class="mb-4"/>
        <span>正在理解你的项目...</span>
      </div>
    </div>
  </div>

  <div class="w-full border theme-border min-h-48 rounded-xl bg-[#f1f4fa]/20 dark:bg-gray-950/10 p-4 mt-4 flex flex-col">
    <div v-if="options.chapters.length === 0 && options.generatedChapters.length === 0" class="mx-auto my-auto">
      <span class="mx-auto my-auto">暂无篇章</span>
    </div>

    <!-- Existing Chapters Workflow -->
    <div v-if="options.chapters.length !== 0" class="flex flex-wrap items-center gap-16">
      <div v-for="(gc, index) in options.chapters"
           class="relative flex items-center group">
        <!-- Node -->
        <div class="relative">
          <!-- Connection Point Left -->
          <div v-if="index > 0"
               class="absolute left-0 top-1/2 transform -translate-x-2 -translate-y-1/2 w-3 h-3 rounded-full bg-purple-400"></div>

          <!-- Connection Point Right -->
          <div v-if="index < options.chapters.length"
               class="absolute right-0 top-1/2 transform translate-x-2 -translate-y-1/2 w-3 h-3 rounded-full bg-purple-400"></div>

          <!-- Node Content -->
          <div class="w-48 p-4 rounded-lg shadow-lg border-2 bg-white dark:bg-indigo-950/20 border-purple-200 hover:border-purple-400 dark:border-purple-900 hover:dark:border-purple-600 transition-all transition-duration-300">
            <div class="text-lg font-medium mb-2">{{ gc.Title }}</div>
            <div class="text-sm text-gray-600 dark:text-gray-400">{{ gc.Description }}</div>
            <!-- Action Buttons -->
            <div class="flex gap-2 mt-4">
              <button @click="goToWriting(gc)"
                      class="px-2 py-1 text-xs bg-green-500 text-white rounded hover:bg-green-600">
                创作
              </button>
              <button class="px-2 py-1 text-xs bg-blue-500 text-white rounded hover:bg-blue-600">
                编辑
              </button>
              <button class="px-2 py-1 text-xs bg-red-500 text-white rounded hover:bg-red-600">
                删除
              </button>
            </div>
          </div>
        </div>

        <!-- Connector Line with Arrow -->
        <div v-if="index < options.chapters.length" class="absolute -right-16 top-1/2 w-16 flex items-center">
          <div class="h-[2px] flex-grow bg-purple-400"></div>
        </div>
      </div>
      <div class="relative">
        <!-- Connection Point Left -->
        <div v-if="index > 0"
             class="absolute left-0 top-1/2 transform -translate-x-2 -translate-y-1/2 w-3 h-3 rounded-full bg-purple-400"></div>

        <!-- Connection Point Right -->
        <div v-if="index < options.chapters.length"
             class="absolute right-0 top-1/2 transform translate-x-2 -translate-y-1/2 w-3 h-3 rounded-full bg-purple-400"></div>

        <!-- Node Content -->
        <div class="w-48 p-12 rounded-lg shadow-lg border-2 bg-white dark:bg-indigo-950/20 border-purple-200 hover:border-purple-400 dark:border-purple-900 hover:dark:border-purple-600 transition-all transition-duration-300 border-dashed cursor-pointer">
          <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5"
               stroke="currentColor" class="size-16 mx-auto my-auto">
            <path stroke-linecap="round" stroke-linejoin="round" d="M12 9v6m3-3H9m12 0a9 9 0 1 1-18 0 9 9 0 0 1 18 0Z"/>
          </svg>
        </div>
      </div>
    </div>

    <!-- Generated Chapters Workflow -->
    <div v-if="options.generatedChapters.length !== 0"
         class="bg-green-300/10 rounded-xl w-full p-4 mt-4">
      <div class="flex flex-wrap items-center gap-16">
        <div v-for="(gc, index) in options.generatedChapters"
             class="relative flex items-center group">
          <!-- Node -->
          <div class="relative">
            <!-- Connection Point Left -->
            <div v-if="index > 0"
                 class="absolute left-0 top-1/2 transform -translate-x-2 -translate-y-1/2 w-3 h-3 rounded-full bg-green-500"></div>

            <!-- Connection Point Right -->
            <div v-if="index < options.generatedChapters.length - 1"
                 class="absolute right-0 top-1/2 transform translate-x-2 -translate-y-1/2 w-3 h-3 rounded-full bg-green-500"></div>

            <!-- Node Content -->
            <div class="w-48 p-4 rounded-lg bg-white dark:bg-gray-800 shadow-lg border-2 border-green-400 dark:border-green-600">
              <div class="text-lg font-medium mb-2">{{ gc.Title }}</div>
              <div class="text-sm text-gray-600 dark:text-gray-400">{{ gc.Description }}</div>
            </div>
          </div>

          <!-- Connector Line with Arrow -->
          <div v-if="index < options.generatedChapters.length - 1" class="absolute -right-16 top-1/2 w-16 flex items-center">
            <div class="h-[2px] flex-grow bg-green-400"></div>
            <div class="absolute right-0 w-2 h-2 rotate-45 border-t-2 border-r-2 border-green-400 transform -translate-y-1/2"></div>
          </div>
        </div>
      </div>

      <!-- Accept/Reject Buttons -->
      <div class="mt-6 text-center">
        <p class="mb-4">您是否接受剧匠AI为您提供的建议？</p>
        <div class="flex justify-center gap-4">
          <button @click="acceptChapter"
                  class="px-4 py-2 bg-green-500 text-white rounded-lg hover:bg-green-600">
            采用
          </button>
          <button @click="rejectChapter"
                  class="px-4 py-2 bg-red-500 text-white rounded-lg hover:bg-red-600">
            丢弃
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
/* Add any additional styles here */
</style>