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

// 定义信息项配置
const infoItems = {
  social_story: { label: '剧情简介' },
  start: { label: '起点' },
  high_point: { label: '高潮' },
  resolved: { label: '结局' }
}
</script>

<template>
  <div class="animate__animated animate__fadeIn">
    <!-- 项目标题区域 -->
    <div class="mb-6">
      <h1 class="text-4xl font-bold bg-gradient-to-r from-indigo-500 to-blue-500 bg-clip-text text-transparent">
        {{ project.project_name }}
      </h1>
      <p class="text-blue-500 mt-2 font-bold flex items-center gap-2">
        <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" d="M3.75 13.5l10.5-11.25L12 10.5h8.25L9.75 21.75 12 13.5H3.75z" />
        </svg>
        {{ project.types }}
      </p>
    </div>

    <!-- 封面和信息区域 -->
    <div class="grid grid-cols-[400px,1fr] gap-8">
      <!-- 封面区域 -->
      <div>
        <a-popover title="编辑作品封面" trigger="hover">
          <template #content>
            <div class="flex flex-nowrap gap-2">
              <button class="flex items-center gap-1 px-3 py-1.5 bg-gradient-to-r from-blue-500 to-cyan-500 text-white rounded-lg hover:opacity-90 transition-all" @click="generateCover">
                <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" d="M9.813 15.904L9 18.75l-.813-2.846a4.5 4.5 0 00-3.09-3.09L2.25 12l2.846-.813a4.5 4.5 0 003.09-3.09L9 5.25l.813 2.846a4.5 4.5 0 003.09 3.09L15.75 12l-2.846.813a4.5 4.5 0 00-3.09 3.09z" />
                </svg>
                AI生成
              </button>
              <button class="flex items-center gap-1 px-3 py-1.5 border theme-border rounded-lg hover:border-blue-500 transition-all">
                <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" d="M3 16.5v2.25A2.25 2.25 0 005.25 21h13.5A2.25 2.25 0 0021 18.75V16.5m-13.5-9L12 3m0 0l4.5 4.5M12 3v13.5" />
                </svg>
                上传
              </button>
              <button class="flex items-center gap-1 px-3 py-1.5 text-red-500 border border-red-200 dark:border-red-800 rounded-lg hover:bg-red-50 dark:hover:bg-red-900/20 transition-all">
                <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" d="M14.74 9l-.346 9m-4.788 0L9.26 9m9.968-3.21c.342.052.682.107 1.022.166m-1.022-.165L18.16 19.673a2.25 2.25 0 01-2.244 2.077H8.084a2.25 2.25 0 01-2.244-2.077L4.772 5.79m14.456 0a48.108 48.108 0 00-3.478-.397m-12 .562c.34-.059.68-.114 1.022-.165m0 0a48.11 48.11 0 013.478-.397m7.5 0v-.916c0-1.18-.91-2.164-2.09-2.201a51.964 51.964 0 00-3.32 0c-1.18.037-2.09 1.022-2.09 2.201v.916m7.5 0a48.667 48.667 0 00-7.5 0" />
                </svg>
                清空
              </button>
            </div>
          </template>
          <div class="w-full h-[500px] border theme-border rounded-xl overflow-hidden group relative hover:border-blue-500 transition-all">
            <!-- 加载状态 -->
            <div v-if="options.isCoverGenerating"
                 class="absolute inset-0 bg-white/90 dark:bg-black/90 backdrop-blur-sm flex items-center justify-center">
              <div class="text-center">
                <loader class="mb-2"/>
                <span class="font-medium text-blue-500">AI 正在创作中...</span>
              </div>
            </div>
            
            <!-- 空状态 -->
            <div v-else-if="project.cover_image === ''"
                 class="w-full h-full flex items-center justify-center bg-slate-50/20 dark:bg-zinc-800/50">
              <div class="text-center text-gray-500">
                <svg xmlns="http://www.w3.org/2000/svg" class="h-12 w-12 mx-auto mb-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" d="M2.25 15.75l5.159-5.159a2.25 2.25 0 013.182 0l5.159 5.159m-1.5-1.5l1.409-1.409a2.25 2.25 0 013.182 0l2.909 2.909m-18 3.75h16.5a1.5 1.5 0 001.5-1.5V6a1.5 1.5 0 00-1.5-1.5H3.75A1.5 1.5 0 002.25 6v12a1.5 1.5 0 001.5 1.5zm10.5-11.25h.008v.008h-.008V8.25zm.375 0a.375.375 0 11-.75 0 .375.375 0 01.75 0z" />
                </svg>
                <p>点击添加封面</p>
              </div>
            </div>
            
            <!-- 封面图片 -->
            <img v-else
                 :src="imagePrefix+project.cover_image" 
                 class="h-full w-full object-cover"
                 alt="项目封面"/>
            
            <!-- 悬浮提示 -->
            <div class="absolute inset-0 bg-black/60 flex items-center justify-center opacity-0 group-hover:opacity-100 transition-opacity">
              <p class="text-white font-medium">点击编辑封面</p>
            </div>
          </div>
        </a-popover>

        <!-- 标签区域 -->
        <div class="space-y-4 mt-4">
          <!-- 风格标签 -->
          <div class="p-4 border theme-border rounded-xl bg-white/50 dark:bg-zinc-900/50 backdrop-blur-sm">
            <h2 class="text-lg font-semibold bg-gradient-to-r from-blue-500 to-cyan-500 bg-clip-text text-transparent mb-3">
              创作风格
            </h2>
            <div class="flex flex-wrap gap-2">
              <a-checkable-tag
                v-for="tag in project.style"
                :key="tag"
                :checked="true"
                class="!px-3 !py-1 !text-sm !font-medium"
              >
                {{ tag }}
              </a-checkable-tag>
            </div>
          </div>

          <!-- 目标人群标签 -->
          <div class="p-4 border theme-border rounded-xl bg-white/50 dark:bg-zinc-900/50 backdrop-blur-sm">
            <h2 class="text-lg font-semibold bg-gradient-to-r from-blue-500 to-cyan-500 bg-clip-text text-transparent mb-3">
              目标人群
            </h2>
            <div class="flex flex-wrap gap-2">
              <a-checkable-tag
                v-for="tag in project.market_people"
                :key="tag"
                :checked="true"
                class="!px-3 !py-1 !text-sm !font-medium"
              >
                {{ tag }}
              </a-checkable-tag>
            </div>
          </div>
        </div>
      </div>

      <!-- 项目信息区域 -->
      <div class="flex flex-col gap-2">
        <div class="p-6 border theme-border rounded-xl mb-6 gradient-bkg backdrop-blur-sm">
          <!-- AI 生成加载状态 -->
          <div v-if="options.isTextGenerating"
               class="h-80 flex items-center justify-center">
            <div class="text-center">
              <loader class="mb-2"/>
              <span class="font-medium text-blue-500">AI 正在思考中...</span>
            </div>
          </div>

          <template v-else>
            <!-- 信息展示区域 -->
            <div v-for="(item, key) in infoItems" :key="key" class="mb-6 last:mb-0">
              <div class="flex items-center justify-between mb-2">
                <h2 class="text-xl font-semibold bg-gradient-to-r from-blue-500 to-cyan-500 bg-clip-text text-transparent">
                  {{ item.label }}
                </h2>
                <div v-if="options.editMode" class="text-sm text-gray-500">
                  已优化 ✨
                </div>
              </div>
              
              <!-- 当前内容 -->
              <div class="mb-3">
                <p class="text-gray-700 dark:text-gray-300">{{ project[key] }}</p>
              </div>
              
              <!-- 优化后的内容 -->
              <div v-if="options.editMode" 
                   class="p-3 rounded-lg bg-green-50 dark:bg-green-900/20 border border-green-200 dark:border-green-800">
                <p class="text-gray-800 dark:text-gray-200">{{ options.infoTemp[key] }}</p>
              </div>
            </div>

            <!-- 操作按钮 -->
            <div class="flex items-center gap-4 mt-6">
              <template v-if="options.editMode">
                <button @click="applyResult"
                        class="flex items-center gap-2 px-4 py-2 bg-green-500 text-white rounded-lg hover:bg-green-600 transition-colors">
                  <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" d="M4.5 12.75l6 6 9-13.5" />
                  </svg>
                  采用优化结果
                </button>
                <button @click="dropResult"
                        class="flex items-center gap-2 px-4 py-2 border border-red-200 dark:border-red-800 text-red-600 rounded-lg hover:bg-red-50 dark:hover:bg-red-900/20 transition-colors">
                  <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12" />
                  </svg>
                  放弃优化
                </button>
              </template>
              <template v-else>
                <button @click="generateInfo"
                        :disabled="options.isTextGenerating"
                        class="flex items-center gap-2 px-4 py-2 bg-gradient-to-r from-blue-500 to-cyan-500 text-white rounded-lg hover:opacity-90 transition-all disabled:opacity-50 disabled:cursor-not-allowed">
                  <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" d="M9.813 15.904L9 18.75l-.813-2.846a4.5 4.5 0 00-3.09-3.09L2.25 12l2.846-.813a4.5 4.5 0 003.09-3.09L9 5.25l.813 2.846a4.5 4.5 0 003.09 3.09L15.75 12l-2.846.813a4.5 4.5 0 00-3.09 3.09z" />
                  </svg>
                  AI 优化
                </button>
                <button class="flex items-center gap-2 px-4 py-2 border theme-border rounded-lg hover:border-blue-500 transition-colors">
                  <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" d="M16.862 4.487l1.687-1.688a1.875 1.875 0 112.652 2.652L10.582 16.07a4.5 4.5 0 01-1.897 1.13L6 18l.8-2.685a4.5 4.5 0 011.13-1.897l8.932-8.931zm0 0L19.5 7.125M18 14v4.75A2.25 2.25 0 0115.75 21H5.25A2.25 2.25 0 013 18.75V8.25A2.25 2.25 0 015.25 6H10" />
                  </svg>
                  手动编辑
                </button>
              </template>
            </div>
          </template>
        </div>

        <!-- 团队信息卡片 -->
        <div class="mt-6 p-6 gradient-bkg rounded-xl border theme-border">
          <div class="flex items-center justify-between mb-4">
            <h2 class="text-xl font-semibold bg-gradient-to-r from-indigo-500 to-blue-500 bg-clip-text text-transparent">
              团队信息
            </h2>
            <button class="text-blue-500 hover:text-blue-600 transition-colors">
              管理团队 →
            </button>
          </div>
          <div class="space-y-3">
            <div class="flex items-center gap-2">
              <span class="font-medium text-gray-700 dark:text-gray-300">团队名称：</span>
              <span class="text-gray-600 dark:text-gray-400">{{ project.team.username }}</span>
            </div>
            <div class="flex items-center gap-2">
              <span class="font-medium text-gray-700 dark:text-gray-300">团队简介：</span>
              <span class="text-gray-600 dark:text-gray-400">{{ project.team.teamDescription }}</span>
            </div>
            <div class="flex items-center gap-2">
              <span class="font-medium text-gray-700 dark:text-gray-300">邀请码：</span>
              <span class="text-gray-600 dark:text-gray-400">{{ project.team.invite_code }}</span>
              <button class="text-blue-500 hover:text-blue-600 transition-colors text-sm">
                复制
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>

</style>