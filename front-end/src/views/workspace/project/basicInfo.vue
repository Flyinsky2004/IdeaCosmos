<script setup>
import {useProjectStore} from "@/stores/project.js";
import {reactive} from "vue";
import Loader from "@/components/loader.vue";
import {post, postJSON} from "@/util/request.js";
import {message} from "ant-design-vue";
import {imagePrefix} from "@/util/VARRIBLES.js";
import {washJSONStr} from "@/util/common.js";
const trueVarible = true
const trueVarible2 = true

const projectStore = useProjectStore()
const project = JSON.parse(localStorage.getItem("project"))
const options = reactive({
  isCoverGenerating: false,
  isTextGenerating: false,
  infoTemp: {
    social_story: '',
    start: '',
    high_point: '',
    resolved: ''
  },
  editMode: false,
})
const generateCover = () => {
  options.isCoverGenerating = true
  post('/api/project/generateCover', {
    project_id: project.ID
  }, (messageer, data) => {
    message.success(messageer)
    project.cover_image = data
    localStorage.setItem("project", JSON.stringify(project))
    options.isCoverGenerating = false
  }, (messageer, data) => {
    message.warning(messageer)
    options.isCoverGenerating = false
  }, (messageer, data) => {
    message.error(messageer)
    options.isCoverGenerating = false
  })
}

const generateInfo = () => {
  options.isTextGenerating = true
  post('/api/project/generateInfo', {
    project_id: project.ID
  }, (messageer, data) => {
    message.success(messageer)
    const raw = data.choices[0].message.content
    options.infoTemp = JSON.parse(washJSONStr(raw))
    options.isTextGenerating = false
    options.editMode = true
  }, (messageer, data) => {
    message.warning(messageer)
    options.isTextGenerating = false
  }, (messageer, data) => {
    message.error(messageer)
    options.isTextGenerating = false
  })
}
const applyResult = () => {
  project.social_story = options.infoTemp.social_story
  project.start = options.infoTemp.start
  project.high_point = options.infoTemp.start
  project.resolved = options.infoTemp.resolved
  postJSON('/api/project/updateProject', project,
  (messageer, data) => {
    message.success(messageer)
    localStorage.setItem("project", JSON.stringify(project))
    options.infoTemp = {
      social_story: '',
      start: '',
      high_point: '',
      resolved: ''
    }
    options.editMode = false
  },(messageer, data) => {
    message.warning(messageer)
      }, (messageer, data) => {
    message.error(messageer)
      })
}
const dropResult = () => {
  options.infoTemp = {
    social_story: '',
    start: '',
    high_point: '',
    resolved: ''
  }
  options.editMode = false
}

// 添加编辑状态控制
const editState = reactive({
  socialStory: false,
  start: false,
  highPoint: false,
  resolved: false
})

// 添加临时编辑内容
const tempContent = reactive({
  socialStory: '',
  start: '',
  highPoint: '',
  resolved: ''
})

// 开始编辑
const startEdit = (field) => {
  editState[field] = true
  tempContent[field] = project[field]
}

// 保存编辑
const saveEdit = async (field) => {
  project[field] = tempContent[field]
  editState[field] = false
  
  // 更新到服务器
  await postJSON('/api/project/updateProject', project,
    (messageer) => {
      message.success('保存成功')
      localStorage.setItem("project", JSON.stringify(project))
    },
    (messageer) => message.warning(messageer),
    (messageer) => message.error(messageer)
  )
}

// 取消编辑
const cancelEdit = (field) => {
  editState[field] = false
  tempContent[field] = ''
}
</script>

<template>
  <div class="animate__animated animate__fadeIn max-w-6xl mx-auto">
    <!-- 项目标题区域 -->
    <div class="flex items-start justify-between mb-8">
      <div>
        <h1 class="text-4xl font-bold bg-gradient-to-r from-blue-500 to-cyan-500 bg-clip-text text-transparent">
          {{ project.project_name }}
        </h1>
        <div class="flex items-center gap-3 mt-2">
          <span class="px-3 py-1 bg-blue-100 dark:bg-blue-900 text-blue-600 dark:text-blue-400 rounded-full text-sm font-medium">
            {{ project.types }}
          </span>
          <span class="text-gray-500 dark:text-gray-400 text-sm">
            ID: #{{ project.ID }}
          </span>
        </div>
      </div>
    </div>

    <div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
      <!-- 左侧封面区域 -->
      <div class="lg:col-span-1">
        <a-popover title="编辑作品封面" trigger="hover">
          <template #content>
            <div class="flex flex-col gap-2">
              <button 
                class="px-4 py-2 bg-gradient-to-r from-blue-500 to-cyan-500 text-white rounded-lg hover:opacity-90 transition-all flex items-center gap-2"
                @click="generateCover"
                :disabled="options.isCoverGenerating"
              >
                <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" d="M13.5 16.875h3.375m0 0h3.375m-3.375 0V13.5m0 3.375v3.375M6 10.5h2.25a2.25 2.25 0 002.25-2.25V6a2.25 2.25 0 00-2.25-2.25H6A2.25 2.25 0 003.75 6v2.25A2.25 2.25 0 006 10.5zm0 9.75h2.25A2.25 2.25 0 0010.5 18v-2.25a2.25 2.25 0 00-2.25-2.25H6a2.25 2.25 0 00-2.25 2.25V18A2.25 2.25 0 006 20.25zm9.75-9.75H18a2.25 2.25 0 002.25-2.25V6A2.25 2.25 0 0018 3.75h-2.25A2.25 2.25 0 0013.5 6v2.25a2.25 2.25 0 002.25 2.25z" />
                </svg>
                AI 生成封面
              </button>
              <button class="px-4 py-2 border theme-border rounded-lg hover:bg-gray-50 dark:hover:bg-gray-800 transition-colors flex items-center gap-2">
                <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" d="M3 16.5v2.25A2.25 2.25 0 005.25 21h13.5A2.25 2.25 0 0021 18.75V16.5m-13.5-9L12 3m0 0l4.5 4.5M12 3v13.5" />
                </svg>
                上传封面
              </button>
              <button class="px-4 py-2 border border-red-200 dark:border-red-800 text-red-600 rounded-lg hover:bg-red-50 dark:hover:bg-red-900/20 transition-colors flex items-center gap-2">
                <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" d="M14.74 9l-.346 9m-4.788 0L9.26 9m9.968-3.21c.342.052.682.107 1.022.166m-1.022-.165L18.16 19.673a2.25 2.25 0 01-2.244 2.077H8.084a2.25 2.25 0 01-2.244-2.077L4.772 5.79m14.456 0a48.108 48.108 0 00-3.478-.397m-12 .562c.34-.059.68-.114 1.022-.165m0 0a48.11 48.11 0 013.478-.397m7.5 0v-.916c0-1.18-.91-2.164-2.09-2.201a51.964 51.964 0 00-3.32 0c-1.18.037-2.09 1.022-2.09 2.201v.916m7.5 0a48.667 48.667 0 00-7.5 0" />
                </svg>
                清除封面
              </button>
            </div>
          </template>
          <div class="aspect-[3/4] w-full bg-gray-100 dark:bg-zinc-800 rounded-xl overflow-hidden border theme-border">
            <div v-if="options.isCoverGenerating" class="w-full h-full flex items-center justify-center">
              <div class="text-center">
                <loader class="mx-auto mb-4"/>
                <span class="text-sm text-gray-500 dark:text-gray-400">AI 正在生成封面...</span>
              </div>
            </div>
            <div v-else-if="!project.cover_image" class="w-full h-full flex items-center justify-center">
              <div class="text-center text-gray-400">
                <svg xmlns="http://www.w3.org/2000/svg" class="h-12 w-12 mx-auto mb-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" d="M2.25 15.75l5.159-5.159a2.25 2.25 0 013.182 0l5.159 5.159m-1.5-1.5l1.409-1.409a2.25 2.25 0 013.182 0l2.909 2.909m-18 3.75h16.5a1.5 1.5 0 001.5-1.5V6a1.5 1.5 0 00-1.5-1.5H3.75A1.5 1.5 0 002.25 6v12a1.5 1.5 0 001.5 1.5zm10.5-11.25h.008v.008h-.008V8.25zm.375 0a.375.375 0 11-.75 0 .375.375 0 01.75 0z" />
                </svg>
                <span class="text-sm">暂无封面</span>
              </div>
            </div>
            <img 
              v-else 
              :src="imagePrefix + project.cover_image" 
              class="w-full h-full object-cover"
              :alt="project.project_name"
            />
          </div>
        </a-popover>
      </div>

      <!-- 右侧内容区域 -->
      <div class="lg:col-span-2 space-y-6">
        <!-- 故事内容卡片 -->
        <div class="bg-white dark:bg-zinc-900 rounded-xl border theme-border p-6">
          <div v-if="options.isTextGenerating" class="h-64 flex items-center justify-center">
            <div class="text-center">
              <loader class="mx-auto mb-4"/>
              <span class="text-sm text-gray-500 dark:text-gray-400">AI 正在优化内容...</span>
            </div>
          </div>
          <div v-else class="space-y-6">
            <!-- 剧情简介 -->
            <div class="space-y-4">
              <div class="flex items-center justify-between">
                <h2 class="text-xl font-semibold text-blue-500">剧情简介</h2>
                <div class="flex items-center gap-2">
                  <button 
                    v-if="!editState.socialStory"
                    @click="startEdit('socialStory')"
                    class="text-gray-500 hover:text-blue-500 transition-colors"
                  >
                    <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                      <path stroke-linecap="round" stroke-linejoin="round" d="M16.862 4.487l1.687-1.688a1.875 1.875 0 112.652 2.652L10.582 16.07a4.5 4.5 0 01-1.897 1.13L6 18l.8-2.685a4.5 4.5 0 011.13-1.897l8.932-8.931zm0 0L19.5 7.125M18 14v4.75A2.25 2.25 0 0115.75 21H5.25A2.25 2.25 0 013 18.75V8.25A2.25 2.25 0 015.25 6H10" />
                    </svg>
                  </button>
                </div>
              </div>
              <div v-if="editState.socialStory" class="space-y-2">
                <textarea
                  v-model="tempContent.socialStory"
                  class="w-full px-4 py-2 border theme-border rounded-lg focus:ring-2 focus:ring-blue-500/20 focus:border-blue-500 dark:bg-zinc-800/50 min-h-[100px]"
                ></textarea>
                <div class="flex justify-end gap-2">
                  <button 
                    @click="cancelEdit('socialStory')"
                    class="px-3 py-1 text-sm border theme-border rounded-lg hover:bg-gray-50 dark:hover:bg-gray-800"
                  >
                    取消
                  </button>
                  <button 
                    @click="saveEdit('socialStory')"
                    class="px-3 py-1 text-sm bg-blue-500 text-white rounded-lg hover:bg-blue-600"
                  >
                    保存
                  </button>
                </div>
              </div>
              <div v-if="options.editMode" class="p-4 bg-green-100 dark:bg-green-900/20 rounded-lg border border-green-200 dark:border-green-800">
                <div class="flex justify-between items-start mb-2">
                  <h3 class="text-sm font-medium text-green-800 dark:text-green-400">AI 优化建议</h3>
                  <div class="flex items-center gap-2">
                    <button 
                      @click="applyResult"
                      class="px-2 py-1 text-xs bg-green-500 text-white rounded hover:bg-green-600 transition-colors flex items-center gap-1"
                    >
                      <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                        <path stroke-linecap="round" stroke-linejoin="round" d="M5 13l4 4L19 7" />
                      </svg>
                      应用
                    </button>
                    <button 
                      @click="dropResult"
                      class="px-2 py-1 text-xs border border-red-200 dark:border-red-800 text-red-600 rounded hover:bg-red-50 dark:hover:bg-red-900/20 transition-colors flex items-center gap-1"
                    >
                      <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                        <path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12" />
                      </svg>
                      放弃
                    </button>
                  </div>
                </div>
                <p class="text-gray-600 dark:text-gray-300">{{ options.infoTemp.social_story }}</p>
              </div>
            </div>

            <!-- 起点 -->
            <div class="space-y-4">
              <div class="flex items-center justify-between">
                <h2 class="text-xl font-semibold text-blue-500">起点</h2>
                <div class="flex items-center gap-2">
                  <button 
                    v-if="!editState.start"
                    @click="startEdit('start')"
                    class="text-gray-500 hover:text-blue-500 transition-colors"
                  >
                    <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                      <path stroke-linecap="round" stroke-linejoin="round" d="M16.862 4.487l1.687-1.688a1.875 1.875 0 112.652 2.652L10.582 16.07a4.5 4.5 0 01-1.897 1.13L6 18l.8-2.685a4.5 4.5 0 011.13-1.897l8.932-8.931zm0 0L19.5 7.125M18 14v4.75A2.25 2.25 0 0115.75 21H5.25A2.25 2.25 0 013 18.75V8.25A2.25 2.25 0 015.25 6H10" />
                    </svg>
                  </button>
                </div>
              </div>
              <div v-if="editState.start" class="space-y-2">
                <textarea
                  v-model="tempContent.start"
                  class="w-full px-4 py-2 border theme-border rounded-lg focus:ring-2 focus:ring-blue-500/20 focus:border-blue-500 dark:bg-zinc-800/50 min-h-[100px]"
                ></textarea>
                <div class="flex justify-end gap-2">
                  <button 
                    @click="cancelEdit('start')"
                    class="px-3 py-1 text-sm border theme-border rounded-lg hover:bg-gray-50 dark:hover:bg-gray-800"
                  >
                    取消
                  </button>
                  <button 
                    @click="saveEdit('start')"
                    class="px-3 py-1 text-sm bg-blue-500 text-white rounded-lg hover:bg-blue-600"
                  >
                    保存
                  </button>
                </div>
              </div>
              <p v-else class="text-gray-600 dark:text-gray-300">{{ project.start }}</p>
              <div v-if="options.editMode" class="p-4 bg-green-100 dark:bg-green-900/20 rounded-lg border border-green-200 dark:border-green-800">
                <div class="flex justify-between items-start mb-2">
                  <h3 class="text-sm font-medium text-green-800 dark:text-green-400">AI 优化建议</h3>
                  <div class="flex items-center gap-2">
                    <button 
                      @click="applyResult"
                      class="px-2 py-1 text-xs bg-green-500 text-white rounded hover:bg-green-600 transition-colors flex items-center gap-1"
                    >
                      <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                        <path stroke-linecap="round" stroke-linejoin="round" d="M5 13l4 4L19 7" />
                      </svg>
                      应用
                    </button>
                    <button 
                      @click="dropResult"
                      class="px-2 py-1 text-xs border border-red-200 dark:border-red-800 text-red-600 rounded hover:bg-red-50 dark:hover:bg-red-900/20 transition-colors flex items-center gap-1"
                    >
                      <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                        <path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12" />
                      </svg>
                      放弃
                    </button>
                  </div>
                </div>
                <p class="text-gray-600 dark:text-gray-300">{{ options.infoTemp.start }}</p>
              </div>
            </div>

            <!-- 高潮 -->
            <div class="space-y-4">
              <div class="flex items-center justify-between">
                <h2 class="text-xl font-semibold text-blue-500">高潮</h2>
                <div class="flex items-center gap-2">
                  <button 
                    v-if="!editState.highPoint"
                    @click="startEdit('highPoint')"
                    class="text-gray-500 hover:text-blue-500 transition-colors"
                  >
                    <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                      <path stroke-linecap="round" stroke-linejoin="round" d="M16.862 4.487l1.687-1.688a1.875 1.875 0 112.652 2.652L10.582 16.07a4.5 4.5 0 01-1.897 1.13L6 18l.8-2.685a4.5 4.5 0 011.13-1.897l8.932-8.931zm0 0L19.5 7.125M18 14v4.75A2.25 2.25 0 0115.75 21H5.25A2.25 2.25 0 013 18.75V8.25A2.25 2.25 0 015.25 6H10" />
                    </svg>
                  </button>
                </div>
              </div>
              <div v-if="editState.highPoint" class="space-y-2">
                <textarea
                  v-model="tempContent.highPoint"
                  class="w-full px-4 py-2 border theme-border rounded-lg focus:ring-2 focus:ring-blue-500/20 focus:border-blue-500 dark:bg-zinc-800/50 min-h-[100px]"
                ></textarea>
                <div class="flex justify-end gap-2">
                  <button 
                    @click="cancelEdit('highPoint')"
                    class="px-3 py-1 text-sm border theme-border rounded-lg hover:bg-gray-50 dark:hover:bg-gray-800"
                  >
                    取消
                  </button>
                  <button 
                    @click="saveEdit('highPoint')"
                    class="px-3 py-1 text-sm bg-blue-500 text-white rounded-lg hover:bg-blue-600"
                  >
                    保存
                  </button>
                </div>
              </div>
              <p v-else class="text-gray-600 dark:text-gray-300">{{ project.high_point }}</p>
              <div v-if="options.editMode" class="p-4 bg-green-100 dark:bg-green-900/20 rounded-lg border border-green-200 dark:border-green-800">
                <div class="flex justify-between items-start mb-2">
                  <h3 class="text-sm font-medium text-green-800 dark:text-green-400">AI 优化建议</h3>
                  <div class="flex items-center gap-2">
                    <button 
                      @click="applyResult"
                      class="px-2 py-1 text-xs bg-green-500 text-white rounded hover:bg-green-600 transition-colors flex items-center gap-1"
                    >
                      <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                        <path stroke-linecap="round" stroke-linejoin="round" d="M5 13l4 4L19 7" />
                      </svg>
                      应用
                    </button>
                    <button 
                      @click="dropResult"
                      class="px-2 py-1 text-xs border border-red-200 dark:border-red-800 text-red-600 rounded hover:bg-red-50 dark:hover:bg-red-900/20 transition-colors flex items-center gap-1"
                    >
                      <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                        <path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12" />
                      </svg>
                      放弃
                    </button>
                  </div>
                </div>
                <p class="text-gray-600 dark:text-gray-300">{{ options.infoTemp.high_point }}</p>
              </div>
            </div>

            <!-- 结局 -->
            <div class="space-y-4">
              <div class="flex items-center justify-between">
                <h2 class="text-xl font-semibold text-blue-500">结局</h2>
                <div class="flex items-center gap-2">
                  <button 
                    v-if="!editState.resolved"
                    @click="startEdit('resolved')"
                    class="text-gray-500 hover:text-blue-500 transition-colors"
                  >
                    <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                      <path stroke-linecap="round" stroke-linejoin="round" d="M16.862 4.487l1.687-1.688a1.875 1.875 0 112.652 2.652L10.582 16.07a4.5 4.5 0 01-1.897 1.13L6 18l.8-2.685a4.5 4.5 0 011.13-1.897l8.932-8.931zm0 0L19.5 7.125M18 14v4.75A2.25 2.25 0 0115.75 21H5.25A2.25 2.25 0 013 18.75V8.25A2.25 2.25 0 015.25 6H10" />
                    </svg>
                  </button>
                </div>
              </div>
              <div v-if="editState.resolved" class="space-y-2">
                <textarea
                  v-model="tempContent.resolved"
                  class="w-full px-4 py-2 border theme-border rounded-lg focus:ring-2 focus:ring-blue-500/20 focus:border-blue-500 dark:bg-zinc-800/50 min-h-[100px]"
                ></textarea>
                <div class="flex justify-end gap-2">
                  <button 
                    @click="cancelEdit('resolved')"
                    class="px-3 py-1 text-sm border theme-border rounded-lg hover:bg-gray-50 dark:hover:bg-gray-800"
                  >
                    取消
                  </button>
                  <button 
                    @click="saveEdit('resolved')"
                    class="px-3 py-1 text-sm bg-blue-500 text-white rounded-lg hover:bg-blue-600"
                  >
                    保存
                  </button>
                </div>
              </div>
              <p v-else class="text-gray-600 dark:text-gray-300">{{ project.resolved }}</p>
              <div v-if="options.editMode" class="p-4 bg-green-100 dark:bg-green-900/20 rounded-lg border border-green-200 dark:border-green-800">
                <div class="flex justify-between items-start mb-2">
                  <h3 class="text-sm font-medium text-green-800 dark:text-green-400">AI 优化建议</h3>
                  <div class="flex items-center gap-2">
                    <button 
                      @click="applyResult"
                      class="px-2 py-1 text-xs bg-green-500 text-white rounded hover:bg-green-600 transition-colors flex items-center gap-1"
                    >
                      <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                        <path stroke-linecap="round" stroke-linejoin="round" d="M5 13l4 4L19 7" />
                      </svg>
                      应用
                    </button>
                    <button 
                      @click="dropResult"
                      class="px-2 py-1 text-xs border border-red-200 dark:border-red-800 text-red-600 rounded hover:bg-red-50 dark:hover:bg-red-900/20 transition-colors flex items-center gap-1"
                    >
                      <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                        <path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12" />
                      </svg>
                      放弃
                    </button>
                  </div>
                </div>
                <p class="text-gray-600 dark:text-gray-300">{{ options.infoTemp.resolved }}</p>
              </div>
            </div>

            <!-- AI优化按钮 -->
            <div class="flex justify-end gap-4 pt-4">
              <div v-if="options.editMode" class="flex items-center gap-2">
                <button 
                  @click="applyResult"
                  class="px-4 py-2 bg-green-500 text-white rounded-lg hover:bg-green-600 transition-colors flex items-center gap-2"
                >
                  <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" d="M5 13l4 4L19 7" />
                  </svg>
                  应用所有优化
                </button>
                <button 
                  @click="dropResult"
                  class="px-4 py-2 border border-red-200 dark:border-red-800 text-red-600 rounded-lg hover:bg-red-50 dark:hover:bg-red-900/20 transition-colors flex items-center gap-2"
                >
                  <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12" />
                  </svg>
                  放弃优化
                </button>
              </div>
              <button 
                @click="generateInfo"
                :disabled="options.isTextGenerating"
                class="px-4 py-2 bg-gradient-to-r from-blue-500 to-cyan-500 text-white rounded-lg hover:opacity-90 transition-all flex items-center gap-2"
              >
                <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" d="M9.813 15.904L9 18.75l-.813-2.846a4.5 4.5 0 00-3.09-3.09L2.25 12l2.846-.813a4.5 4.5 0 003.09-3.09L9 5.25l.813 2.846a4.5 4.5 0 003.09 3.09L15.75 12l-2.846.813a4.5 4.5 0 00-3.09 3.09zM18.259 8.715L18 9.75l-.259-1.035a3.375 3.375 0 00-2.455-2.456L14.25 6l1.036-.259a3.375 3.375 0 002.455-2.456L18 2.25l.259 1.035a3.375 3.375 0 002.456 2.456L21.75 6l-1.035.259a3.375 3.375 0 00-2.456 2.456zM16.894 20.567L16.5 21.75l-.394-1.183a2.25 2.25 0 00-1.423-1.423L13.5 18.75l1.183-.394a2.25 2.25 0 001.423-1.423l.394-1.183.394 1.183a2.25 2.25 0 001.423 1.423l1.183.394-1.183.394a2.25 2.25 0 00-1.423 1.423z" />
                </svg>
                AI 优化内容
                <loader v-if="options.isTextGenerating" class="ml-2 h-5 w-5"/>
              </button>
            </div>
          </div>
        </div>

        <!-- 标签区域 -->
        <div class="bg-white dark:bg-zinc-900 rounded-xl border theme-border p-6">
          <div class="space-y-4">
            <h2 class="text-xl font-semibold text-blue-500">风格标签</h2>
            <div class="flex flex-wrap gap-2">
              <span 
                v-for="tag in project.style" 
                :key="tag"
                class="px-3 py-1 bg-blue-100 dark:bg-blue-900 text-blue-600 dark:text-blue-400 rounded-full text-sm"
              >
                {{ tag }}
              </span>
            </div>
          </div>
          <div class="mt-6 space-y-4">
            <h2 class="text-xl font-semibold text-blue-500">目标受众</h2>
            <div class="flex flex-wrap gap-2">
              <span 
                v-for="tag in project.market_people" 
                :key="tag"
                class="px-3 py-1 bg-green-100 dark:bg-green-900 text-green-600 dark:text-green-400 rounded-full text-sm"
              >
                {{ tag }}
              </span>
            </div>
          </div>
        </div>

        <!-- 团队信息 -->
        <div class="bg-white dark:bg-zinc-900 rounded-xl border theme-border p-6">
          <h2 class="text-xl font-semibold text-blue-500 mb-4">团队信息</h2>
          <div class="space-y-3">
            <div class="flex items-center gap-2">
              <span class="text-gray-500 dark:text-gray-400">团队名称：</span>
              <span class="font-medium text-gray-900 dark:text-gray-100">{{ project.team.username }}</span>
            </div>
            <div class="flex items-center gap-2">
              <span class="text-gray-500 dark:text-gray-400">团队简介：</span>
              <span class="text-gray-600 dark:text-gray-300">{{ project.team.teamDescription || '暂无简介' }}</span>
            </div>
            <div class="flex items-center gap-2">
              <span class="text-gray-500 dark:text-gray-400">邀请码：</span>
              <code class="px-2 py-1 bg-gray-100 dark:bg-gray-800 rounded text-sm font-mono">
                {{ project.team.invite_code }}
              </code>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.loader {
  @apply animate-spin;
}
</style>