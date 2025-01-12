<script setup>

import {onMounted, reactive} from "vue";
import {get} from "@/util/request.js";
import {parseDateTime} from "@/util/common.js";
import router from "@/router/index.js";
import {useProjectStore} from "@/stores/project.js";
import {imagePrefix} from "@/util/VARRIBLES.js";

const options = reactive({
  projects: [],
  isAddWindowOpen: false,
  nowHoverId: -1,
})
const fetchMyProjects = () => {
  get('/api/project/myProjects', {},
      (message, data) => {
        options.projects = data;
      })
}

onMounted(() => {
  fetchMyProjects()
})
const moveIn = (id) => {
  options.nowHoverId = id;
}
const moveOut = () => {
  options.nowHoverId = -1;
}
const projectStore = useProjectStore()
const enterProject = (project) => {
  projectStore.setProject(project)
  localStorage.setItem("project", JSON.stringify(project))
  router.push('/workspace/editProject/index')
}
</script>

<template>
  <div class="flex flex-col gap-2 animate__animated animate__fadeIn">
    <div class="border p-4 workspace-box w-fit">
      <h1 class="text-2xl font-serif">创剧空间</h1>
      <span class="text-sm font-bold">我们的团队功能支持您将创剧空间中的项目与您的团队绑定，使您的团队成员一同加入您的精彩内容创作！<br/>我们支持:项目多人分工协作/实时评论反馈/项目版本控制</span>
    </div>
    <div class="border p-4 workspace-box">
      <h1 class="text-2xl">我的项目</h1>
      <div class="w-full grid grid-cols-5 gap-2">
        <div v-for="project in options.projects" class="border rounded-xl dark:border-[rgb(118,118,118)] outline-[1px] min-h-56 text-theme-switch
hover:bg-gray-100/50 active:bg-gray-100/90 dark:hover:bg-gray-900/5 dark:active:bg-gray-900/10 cursor-pointer"
             @mouseover="moveIn(project.ID)" @mouseleave="moveOut()">
          <div class="p-2" v-if="options.nowHoverId !== project.ID">
            <div class="w-full h-80 border theme-border rounded-xl">
              <div v-if="project.cover_image === ''" class="w-full h-full flex bg-slate-50/20 dark:bg-[#242424]/50 rounded-xl">
                <div class="mx-auto my-auto place-items-center">
                  <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="size-6 ">
                    <path stroke-linecap="round" stroke-linejoin="round" d="M9.879 7.519c1.171-1.025 3.071-1.025 4.242 0 1.172 1.025 1.172 2.687 0 3.712-.203.179-.43.326-.67.442-.745.361-1.45.999-1.45 1.827v.75M21 12a9 9 0 1 1-18 0 9 9 0 0 1 18 0Zm-9 5.25h.008v.008H12v-.008Z" />
                  </svg>
                  暂无封面
                </div>
              </div>
              <img class="h-full w-auto border theme-border rounded-xl" :src="imagePrefix+project.cover_image" v-else alt="用户头像"/>
            </div>
            <h1 class="text-2xl text-blue-500">#{{ project.ID }}{{ project.project_name }}</h1>
            <h1 class="text-sm text-blue-600">项目团队:{{ project.team.username }}</h1>
          </div>
          <div @click="enterProject(project)"
              class="flex flex-col w-full h-full animate__animated animate__fadeIn animate__faster p-2
          bg-gray-100/10 active:bg-gray-200/90 dark:bg-gray-950/10 dark:active:bg-gray-950/15 cursor-pointer" v-else>
            <h1 class="text-2xl text-blue-500">#{{ project.ID }}{{ project.project_name }}</h1>
            <h1 class="text-sm text-blue-600">项目团队:{{ project.team.username }}</h1>
            <h1 class="text-sm font-bold text-blue-400">项目简介：</h1>
            <h1 class="text-sm">{{ project.social_story }}...</h1>
            <h1 class="text-sm font-bold text-blue-400">创建时间：</h1>
            <span class="text-sm">{{ parseDateTime(project.CreatedAt) }}</span>
            <span class="mx-auto my-auto flex flex-nowrap">
              <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5"
                   stroke="currentColor" class="size-5">
  <path stroke-linecap="round" stroke-linejoin="round"
        d="M15.59 14.37a6 6 0 0 1-5.84 7.38v-4.8m5.84-2.58a14.98 14.98 0 0 0 6.16-12.12A14.98 14.98 0 0 0 9.631 8.41m5.96 5.96a14.926 14.926 0 0 1-5.841 2.58m-.119-8.54a6 6 0 0 0-7.381 5.84h4.8m2.581-5.84a14.927 14.927 0 0 0-2.58 5.84m2.699 2.7c-.103.021-.207.041-.311.06a15.09 15.09 0 0 1-2.448-2.448 14.9 14.9 0 0 1 .06-.312m-2.24 2.39a4.493 4.493 0 0 0-1.757 4.306 4.493 4.493 0 0 0 4.306-1.758M16.5 9a1.5 1.5 0 1 1-3 0 1.5 1.5 0 0 1 3 0Z"/>
</svg>

              进入项目空间</span>
          </div>
        </div>
        <div class="border border-dashed rounded-xl dark:border-[rgb(118,118,118)] outline-[1px] min-h-56 place-items-center place-content-center
hover:bg-gray-100/50 active:bg-gray-100/90 dark:hover:dark:bg-gray-950/10 dark:active:bg-gray-900/10 cursor-pointer"
             @click="router.push('/workspace/newProject')">
          <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5"
               stroke="currentColor" class="size-16">
            <path stroke-linecap="round" stroke-linejoin="round" d="M12 9v6m3-3H9m12 0a9 9 0 1 1-18 0 9 9 0 0 1 18 0Z"/>
          </svg>
        </div>
        <div class="border border-dashed rounded-xl dark:border-[rgb(118,118,118)] outline-[1px] min-h-56 place-items-center place-content-center
hover:bg-gray-100/50 active:bg-gray-100/90 dark:hover:dark:bg-gray-950/10 dark:active:bg-gray-900/10 cursor-pointer"
             @click="router.push('/workspace/newProject')">
          <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="size-14">
            <path stroke-linecap="round" stroke-linejoin="round" d="M7.5 8.25h9m-9 3H12m-9.75 1.51c0 1.6 1.123 2.994 2.707 3.227 1.129.166 2.27.293 3.423.379.35.026.67.21.865.501L12 21l2.755-4.133a1.14 1.14 0 0 1 .865-.501 48.172 48.172 0 0 0 3.423-.379c1.584-.233 2.707-1.626 2.707-3.228V6.741c0-1.602-1.123-2.995-2.707-3.228A48.394 48.394 0 0 0 12 3c-2.392 0-4.744.175-7.043.513C3.373 3.746 2.25 5.14 2.25 6.741v6.018Z" />
          </svg>
        </div>
        <div class="border border-dashed rounded-xl dark:border-[rgb(118,118,118)] outline-[1px] min-h-56 place-items-center place-content-center
hover:bg-gray-100/50 active:bg-gray-100/90 dark:hover:dark:bg-gray-950/10 dark:active:bg-gray-900/10 cursor-pointer"
             @click="router.push('/workspace/newProject')">
          <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="size-14">
            <path stroke-linecap="round" stroke-linejoin="round" d="M7.5 7.5h-.75A2.25 2.25 0 0 0 4.5 9.75v7.5a2.25 2.25 0 0 0 2.25 2.25h7.5a2.25 2.25 0 0 0 2.25-2.25v-7.5a2.25 2.25 0 0 0-2.25-2.25h-.75m-6 3.75 3 3m0 0 3-3m-3 3V1.5m6 9h.75a2.25 2.25 0 0 1 2.25 2.25v7.5a2.25 2.25 0 0 1-2.25 2.25h-7.5a2.25 2.25 0 0 1-2.25-2.25v-.75" />
          </svg>
        </div>


      </div>

    </div>
    <div class="border p-4 workspace-box">
      <h1 class="text-2xl">我加入的团队</h1>
      {{ options }}
    </div>
  </div>

</template>

<style scoped>

</style>