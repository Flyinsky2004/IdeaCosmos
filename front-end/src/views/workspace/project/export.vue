<script setup>
import {onMounted, reactive} from "vue";
import {audioModels, BACKEND_DOMAIN, femaleAvatar, FRONTEND_DOMAIN, maleAvatar} from "@/util/VARRIBLES.js";
import {get} from "@/util/request.js";
import {message} from "ant-design-vue";
import {parseDateTime} from "@/util/common.js";
import {useUserStore} from "@/stores/user.js";
import axios from "axios";
import router from "@/router/index.js";
import {exportScript} from "@/util/exportTool.js";
const userStore = useUserStore()
const project = JSON.parse(localStorage.getItem("project"));
const options = reactive({
  currentChapter: {},
  currentAudioName: '',
  chapters: [],
  currentAudioFileName: '',
  isAudioPlaying: false,
  audioUrl: ''
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
const chapterClickHandler = (chapter, versionID) => {
  if (versionID !== 0) {
    options.currentChapterId = chapter.ID
    options.currentChapter = chapter
    if (chapter.current_version !== undefined) {
      fetchAudioResource()
    }
  } else {
    message.info("该篇章还没有创作版本，请前往创作板块创作！")
  }
}
const generateAudio = () => {
  get('/api/project/generateChapterAudio', {
    chapterId: options.currentChapterId,
    audioName: options.currentAudioName
  }, (messageer, data) => {
    message.success(messageer)
    options.currentAudioFileName = data
  }, (messageer, data) => {
    message.warning(messageer)
  }, (messageer, data) => {
    message.warning(messageer)
  })
}
const playAudio = () => {
  options.isAudioPlaying = true
}
const stopAudio = () => {
  options.isAudioPlaying = false
}
const fetchAudioResource = async () => {
  // 发送请求，获取MP3文件流
  const response = await axios.get(BACKEND_DOMAIN + 'audio/' + options.currentChapter.current_version.audio_path, {
    responseType: 'arraybuffer', // 确保获取的是二进制数据
  });

  // 将 ArrayBuffer 转化为 Blob
  const blob = new Blob([response.data], {type: 'audio/mpeg'});

  // 将 Blob 转化为临时URL
  options.audioUrl = URL.createObjectURL(blob);

}
</script>

<template>
  <div>
    <h1 class="text-3xl text-blue-500">选择篇章</h1>
    <div class="w-full grid grid-cols-5 gap-2">
      <div
          v-for="chapter in options.chapters"
          :key="chapter.ID"
          class="group aspect-square p-4 border-2 rounded-xl cursor-pointer transition-all duration-200"
          :class="{
      'bg-blue-50/70 dark:bg-blue-900/20 border-blue-400 dark:border-blue-800': chapter.ID === options.currentChapterId,
      'hover:bg-gray-50 dark:hover:bg-gray-800/50 border-gray-200 dark:border-gray-700': chapter.ID !== options.currentChapterId,
      'cursor-not-allowed' : chapter.version_id === 0
    }"
          @click="chapterClickHandler(chapter,chapter.version_id)"
      >
        <div class="flex flex-col h-full">
          <!-- 标题和元信息 -->
          <div class="flex items-center justify-between mb-3">
            <h3 class="text-lg font-semibold text-gray-900 dark:text-gray-100 truncate">
              {{ chapter.Title }}
            </h3>
            <span
                class="text-xs px-2 py-1 rounded-full bg-gray-100 text-gray-700 dark:bg-gray-700/30 dark:text-gray-300">
         版本 {{ chapter.current_version.ID }}
        </span>
          </div>

          <!-- 描述 -->
          <p class="text-sm text-gray-600 dark:text-gray-400 line-clamp-3 flex-1">
            {{ chapter.Description }}
          </p>

          <!-- 底部元信息 -->
          <div class="mt-4 pt-3 border-t border-gray-200 dark:border-gray-700">
            <div class="flex flex-col justify-between text-xs text-gray-500 dark:text-gray-400">
              <span>创建时间：{{ parseDateTime(chapter.CreatedAt) }}</span>
              <span>更新时间：{{ parseDateTime(chapter.UpdatedAt) }}</span>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
  <div class="mt-4">
    <h1 class="text-3xl text-blue-500">音频导出</h1>
    <div class="w-full grid grid-cols-4 gap-2">
      <div
          v-for="(model, index) in audioModels"
          :key="model.model"
          class="flex flex-col group aspect-square p-4 border-2 rounded-xl cursor-pointer transition-all duration-200"
          :class="{
      'bg-blue-50/70 dark:bg-blue-900/20 border-blue-400 dark:border-blue-800': model.model === options.currentAudioName,
      'hover:bg-gray-50 dark:hover:bg-gray-800/50 border-gray-200 dark:border-gray-700': model.model !== options.currentAudioName
    }"
          @click="options.currentAudioName = model.model"
      >
        <div class="flex items-start gap-4">
          <!-- 头像 -->
          <div class="shrink-0">
            <div class="w-fit h-fit rounded-full bg-gradient-to-br p-1"
                 :class="model.isMale
            ? 'from-blue-400 to-blue-600 dark:from-blue-600 dark:to-blue-800'
            : 'from-pink-400 to-rose-500 dark:from-rose-600 dark:to-rose-800'">
              <div class="w-full h-full rounded-full bg-white dark:bg-gray-900 p-1.5">
                <div v-html="model.isMale ? maleAvatar : femaleAvatar"
                     class="text-gray-800 dark:text-gray-200"></div>
              </div>
            </div>
          </div>
          <!-- 内容 -->
          <div class="flex-1 min-w-0">
            <h3 class="text-lg font-semibold text-gray-900 dark:text-gray-100">
              {{ model.name }}
              <span class="text-sm ml-2 px-2 py-0.5 rounded-full"
                    :class="model.isMale
              ? 'bg-blue-100 text-blue-800 dark:bg-blue-900/50 dark:text-blue-300'
              : 'bg-rose-100 text-rose-800 dark:bg-rose-900/50 dark:text-rose-300'">
            {{ model.isMale ? '男声' : '女声' }}
          </span>
            </h3>

            <p class="text-sm mt-1 text-gray-600 dark:text-gray-400 line-clamp-2">
              {{ model.description }}
            </p>

            <!-- 风格标签 -->
            <div class="mt-2 flex flex-wrap gap-2">
              <template v-if="Array.isArray(model.style)">
            <span v-for="(style, i) in model.style"
                  :key="i"
                  class="text-xs px-2 py-1 rounded-full border"
                  :class="index === options.currentAudioId
                    ? 'bg-blue-100/50 border-blue-400 text-blue-800 dark:bg-blue-900/30 dark:border-blue-800 dark:text-blue-300'
                    : 'bg-gray-200 border-gray-200 text-gray-700 dark:bg-gray-700/30 dark:border-gray-600 dark:text-gray-300'">
              {{ style }}
            </span>
              </template>
              <span v-else class="text-xs px-2 py-1 rounded-full border"
                    :class="index === options.currentAudioId
                  ? 'bg-blue-100/50 border-blue-200 text-blue-800 dark:bg-blue-900/30 dark:border-blue-800 dark:text-blue-300'
                  : 'bg-gray-100 border-gray-200 text-gray-700 dark:bg-gray-700/30 dark:border-gray-600 dark:text-gray-300'">
            {{ model.style }}
          </span>
            </div>
          </div>
        </div>
        <div class="flex flex-grow"></div>
        <audio controls>
          <source :src="FRONTEND_DOMAIN + '/audio/models/' + model.model + '.wav'" type="audio/wav">
        </audio>
      </div>
    </div>
    <div>
      <!--      <audio controls crossorigin="anonymous">-->
      <!--        <source src="http://localhost:8080/api/audio/20250129_vi5mms.mp3" type="audio/mpeg">-->
      <!--        Your browser does not support the audio element.-->
      <!--      </audio>-->
      <div class="w-fit p-4 bg-red-500 flex">
        <kinesis-container class="mx-auto my-auto" :playAudio="options.isAudioPlaying"
                           :audio="options.audioUrl">
          <kinesis-audio
              :audioIndex="17">
            <kinesis-element
                :strength="10"
                type="depth">
              <img class="rounded-xl w-32 aspect-square" :src="BACKEND_DOMAIN + userStore.user.avatar" alt="用户头像">
            </kinesis-element>
          </kinesis-audio>
        </kinesis-container>
      </div>
      <div>
        <button
            @click="playAudio"
            :disabled="options.isAudioPlaying"
        >
          Play
        </button>
        <button
            @click="stopAudio"
            :disabled="!options.isAudioPlaying"
        >
          Stop
        </button>
      </div>
      <button class="basic-prinary-button" @click="generateAudio">
        generateAudio
      </button>
    </div>
  </div>
  <div class="mt-4">
    <h1 class="text-3xl text-blue-500">文件导出</h1>
    <div class="grid grid-cols-5 gap-4">
      <div class="border border-dashed rounded-xl dark:border-[rgb(118,118,118)] outline-[1px] min-h-56 place-items-center place-content-center text-gray-400
hover:bg-gray-100/50 active:bg-gray-100/90 dark:hover:dark:bg-gray-950/10 dark:active:bg-gray-900/10 cursor-pointer"
           @click="exportScript(project,project.project_name,options.currentChapter.Title,options.currentChapter.current_version.content,'pdf')">
        <h1 class="text-5xl">PDF</h1>
        <span class="font-bold text-theme-switch">导出为PDF文档</span>
      </div>
      <div class="border border-dashed rounded-xl dark:border-[rgb(118,118,118)] outline-[1px] min-h-56 place-items-center place-content-center text-gray-400
hover:bg-gray-100/50 active:bg-gray-100/90 dark:hover:dark:bg-gray-950/10 dark:active:bg-gray-900/10 cursor-pointer"
           @click="exportScript(project,project.project_name,options.currentChapter.Title,options.currentChapter.current_version.content,'word')">
        <h1 class="text-5xl">DOCX</h1>
        <span class="font-bold text-theme-switch">导出为Word文档</span>
      </div>
      <div class="border border-dashed rounded-xl dark:border-[rgb(118,118,118)] outline-[1px] min-h-56 place-items-center place-content-center text-gray-400
hover:bg-gray-100/50 active:bg-gray-100/90 dark:hover:dark:bg-gray-950/10 dark:active:bg-gray-900/10 cursor-pointer"
           @click="exportScript(project,project.project_name,options.currentChapter.Title,options.currentChapter.current_version.content,'markdown')">
        <h1 class="text-5xl">Markdown</h1>
        <span class="font-bold text-theme-switch">导出为富文本</span>
      </div>
    </div>
  </div>
</template>

<style scoped>

</style>