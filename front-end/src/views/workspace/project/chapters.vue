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
    <div class="border border-dashed rounded-xl dark:border-[rgb(118,118,118)] outline-[1px] min-h-56 place-items-center place-content-center text-gray-400
hover:bg-gray-100/50 active:bg-gray-100/90 dark:hover:dark:bg-gray-950/10 dark:active:bg-gray-900/10 cursor-pointer"
         @click="router.push('/workspace/newProject')">
      <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5"
           stroke="currentColor" class="size-16">
        <path stroke-linecap="round" stroke-linejoin="round" d="M12 9v6m3-3H9m12 0a9 9 0 1 1-18 0 9 9 0 0 1 18 0Z"/>
      </svg>
    </div>
    <div @click="generateChapters"
         class="border border-dashed rounded-xl dark:border-[rgb(118,118,118)] outline-[1px] min-h-56 place-items-center place-content-center
hover:bg-gray-100/50 active:bg-gray-100/90 dark:hover:dark:bg-gray-950/10 dark:active:bg-gray-900/10 cursor-pointer">
      <div v-if="!options.isChapterGenerating">
        <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="currentColor"
             class="size-7 dark:text-purple-400">
          <path fill-rule="evenodd"
                d="M9 4.5a.75.75 0 0 1 .721.544l.813 2.846a3.75 3.75 0 0 0 2.576 2.576l2.846.813a.75.75 0 0 1 0 1.442l-2.846.813a3.75 3.75 0 0 0-2.576 2.576l-.813 2.846a.75.75 0 0 1-1.442 0l-.813-2.846a3.75 3.75 0 0 0-2.576-2.576l-2.846-.813a.75.75 0 0 1 0-1.442l2.846-.813A3.75 3.75 0 0 0 7.466 7.89l.813-2.846A.75.75 0 0 1 9 4.5ZM18 1.5a.75.75 0 0 1 .728.568l.258 1.036c.236.94.97 1.674 1.91 1.91l1.036.258a.75.75 0 0 1 0 1.456l-1.036.258c-.94.236-1.674.97-1.91 1.91l-.258 1.036a.75.75 0 0 1-1.456 0l-.258-1.036a2.625 2.625 0 0 0-1.91-1.91l-1.036-.258a.75.75 0 0 1 0-1.456l1.036-.258a2.625 2.625 0 0 0 1.91-1.91l.258-1.036A.75.75 0 0 1 18 1.5ZM16.5 15a.75.75 0 0 1 .712.513l.394 1.183c.15.447.5.799.948.948l1.183.395a.75.75 0 0 1 0 1.422l-1.183.395c-.447.15-.799.5-.948.948l-.395 1.183a.75.75 0 0 1-1.422 0l-.395-1.183a1.5 1.5 0 0 0-.948-.948l-1.183-.395a.75.75 0 0 1 0-1.422l1.183-.395c.447-.15.799-.5.948-.948l.395-1.183A.75.75 0 0 1 16.5 15Z"
                clip-rule="evenodd"/>
        </svg>
      </div>
      <div v-else class="text-center place-items-center">
        <loader class="mb-4"/>
        <span>正在理解你的项目...</span>
      </div>
    </div>
  </div>
  <div
      class="w-full border theme-border min-h-48 rounded-xl bg-[#f1f4fa]/20 dark:bg-gray-950/10 p-4 mt-4 flex flex-col">
    <div v-if="options.chapters.length === 0 && options.generatedChapters.length === 0" class="mx-auto my-auto">
      <span class="mx-auto my-auto">暂无篇章</span>
    </div>
    <div v-if="options.chapters.length !== 0">
      <a-timeline>
        <a-timeline-item v-for="gc in options.chapters">
          <a-popover title="操作" placement="topLeft">
            <template #content>
              <div class="flex flex-nowrap gap-2">
                <button class="basic-success-button" @click="goToWriting(gc)">
                  创作篇章
                </button>
                <button class="basic-prinary-button">编辑篇章</button>
                <button class="basic-error-button">删除篇章</button>
              </div>
            </template>
            {{ gc.Title }}:{{ gc.Description }}
          </a-popover>
        </a-timeline-item>
      </a-timeline>
    </div>
    <div v-if="options.generatedChapters.length !== 0" class="bg-green-300/10 rounded-xl w-full p-4">
      <a-timeline>
        <a-timeline-item v-for="gc in options.generatedChapters">{{ gc.Title }}:{{ gc.Description }}</a-timeline-item>
      </a-timeline>
      您是否接受剧匠AI为您提供的建议？
      <div class="flex flex-nowrap gap-2 mt-2">
        <button class="basic-prinary-button my-auto" @click="acceptChapter">
          采用
        </button>
        <button class="basic-error-button" @click="rejectChapter">
          丢弃
        </button>
      </div>
    </div>
  </div>

</template>

<style scoped>

</style>