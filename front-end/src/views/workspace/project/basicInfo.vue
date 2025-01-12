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
</script>

<template>
  <div class="animate__animated animate__fadeIn ">
    <h1 class="text-4xl font-bold text-indigo-500">{{ project.project_name }}</h1>
    <p class="text-blue-500 mt-2 font-bold">{{ project.types }}</p>
    <div class="flex flex-nowrap gap-4">
      <a-popover title="编辑作品封面">
        <template #content>
          <div class="flex flex-nowrap gap-2">
            <button class="color-mixed-button" @click="generateCover">AI生成</button>
            <button class="btn1">手动上传</button>
            <button class="transparent-button">清空封面</button>
          </div>
        </template>
        <div class="w-72 h-80 border theme-border rounded-xl">
          <div v-if="options.isCoverGenerating"
               class="w-full h-full flex bg-slate-50/20 dark:bg-[#242424]/50 rounded-xl">
            <div class="mx-auto my-auto place-items-center">
              <loader class="m-2"/>
              <span class="font-sans font-bold">吃奶的劲加载中...</span>
            </div>
          </div>
          <div v-else-if="project.cover_image === ''"
               class="w-full h-full flex bg-slate-50/20 dark:bg-[#242424]/50 rounded-xl">
            <div class="mx-auto my-auto place-items-center">
              <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5"
                   stroke="currentColor" class="size-6 ">
                <path stroke-linecap="round" stroke-linejoin="round"
                      d="M9.879 7.519c1.171-1.025 3.071-1.025 4.242 0 1.172 1.025 1.172 2.687 0 3.712-.203.179-.43.326-.67.442-.745.361-1.45.999-1.45 1.827v.75M21 12a9 9 0 1 1-18 0 9 9 0 0 1 18 0Zm-9 5.25h.008v.008H12v-.008Z"/>
              </svg>
              暂无封面
            </div>
          </div>
          <img :src="imagePrefix+project.cover_image" class="h-full w-auto border theme-border rounded-xl" v-else
               alt=""/>
        </div>
      </a-popover>
      <div class="flex flex-col w-full">
        <div class="p-2 border theme-border rounded-xl mb-4">
          <div v-if="options.isTextGenerating"
               class="w-full h-80 mb-4 flex bg-slate-50/20 dark:bg-[#242424]/50 rounded-xl">
            <div class="mx-auto my-auto place-items-center">
              <loader class="m-2"/>
              <span class="font-sans font-bold">AI使劲思考中...</span>
            </div>
          </div>
          <div v-else>
            <div class="mb-4">
              <h2 class="text-2xl font-semibold text-blue-500">剧情简介</h2>
              <p class="mt-2 ">{{ project.social_story }}</p>
            </div>
            <div v-if="options.editMode" class="mb-4 bg-green-400/20 dark:bg-green-400/10 rounded-xl border theme-border p-2">
              <h2 class="text-xl font-semibold text-blue-500">优化后:</h2>
              <p class="mt-2 ">{{ options.infoTemp.social_story }}</p>
            </div>

            <!-- 起点 -->
            <div class="mb-4">
              <h2 class="text-2xl font-semibold text-blue-500">起点</h2>
              <p class="mt-2 ">{{ project.start }}</p>
            </div>
            <div v-if="options.editMode" class="mb-4 bg-green-400/20 dark:bg-green-400/10 rounded-xl border theme-border p-2">
              <h2 class="text-xl font-semibold text-blue-500">优化后:</h2>
              <p class="mt-2 ">{{ options.infoTemp.start }}</p>
            </div>

            <!-- 高潮 -->
            <div class="mb-4">
              <h2 class="text-2xl font-semibold text-blue-500">高潮</h2>
              <p class="mt-2 ">{{ project.high_point }}</p>
            </div>
            <div v-if="options.editMode" class="mb-4 bg-green-400/20 dark:bg-green-400/10 rounded-xl border theme-border p-2">
              <h2 class="text-xl font-semibold text-blue-500">优化后:</h2>
              <p class="mt-2 ">{{ options.infoTemp.high_point }}</p>
            </div>

            <!-- 结局 -->
            <div class="mb-4">
              <h2 class="text-2xl font-semibold text-blue-500">结局</h2>
              <p class="mt-2 ">{{ project.resolved }}</p>
            </div>
            <div v-if="options.editMode" class="mb-4 bg-green-400/20 dark:bg-green-400/10 rounded-xl border theme-border p-2">
              <h2 class="text-xl font-semibold text-blue-500">优化后:</h2>
              <p class="mt-2 ">{{ options.infoTemp.resolved }}</p>
            </div>
          </div>
          <div v-if="options.editMode" class="flex flex-nowrap gap-2 mb-2">
            <span class="my-auto">是否采用优化结果?</span>
            <button class="basic-prinary-button my-auto" @click="applyResult">
              采用
            </button>
            <button class="basic-error-button" @click="dropResult">
              丢弃
            </button>
          </div>

          <div class="w-full flex flex-nowrap gap-4">
            <button class="color-mixed-button" @click="generateInfo">
              智能优化
            </button>
            <button class="transparent-button" @click="">
              手动编辑
            </button>
          </div>
        </div>

        <!-- 风格 -->
        <div class="mb-4 ml-2">
          <h2 class="text-2xl font-semibold text-blue-500">风格</h2>
          <ul class="mt-2 list-disc list-inside ">
            <a-checkable-tag
                class="text-md font-bold ml-[1px] mt-2"
                v-for="(tag, index) in project.style"
                :key="tag"
                v-model:checked="trueVarible"
            >
              {{ tag }}
            </a-checkable-tag>
          </ul>
        </div>

        <!-- 目标市场人群 -->
        <div class="mb-4 ml-2">
          <h2 class="text-2xl font-semibold text-blue-500">目标市场人群</h2>
          <a-checkable-tag
              class="text-md font-bold ml-[1px] mt-2"
              v-for="(tag) in project.market_people"
              :key="tag"
              v-model:checked="trueVarible2"
          >
            {{ tag }}
          </a-checkable-tag>
        </div>
      </div>
    </div>
    <div class="mt-6 bg-indigo-50 dark:bg-indigo-950/20 p-4 rounded-lg border theme-border">
      <h2 class="text-xl font-semibold text-indigo-600">团队信息</h2>
      <p class="mt-2 "><strong>团队名称：</strong>{{ project.team.username }}</p>
      <p class="mt-2 "><strong>团队简介：</strong>{{ project.team.teamDescription }}</p>
      <p class="mt-2 "><strong>邀请码：</strong>{{ project.team.invite_code }}</p>
    </div>
  </div>
</template>

<style scoped>

</style>