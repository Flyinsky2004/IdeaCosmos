<script setup>

import {onMounted, reactive} from "vue";
import {get} from "@/util/request.js";
import {parseDateTime} from "@/util/common.js";
import router from "@/router/index.js";

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
            <h1 class="text-2xl text-blue-500">#{{ project.ID }}{{ project.project_name }}</h1>
            <h1 class="text-sm text-blue-600">项目团队:{{ project.team.username }}</h1>
            <h1 class="text-sm font-bold text-blue-400">项目简介：</h1>
            <h1 class="text-sm">{{ project.social_story }}...</h1>
            <h1 class="text-sm font-bold text-blue-400">创建时间：</h1>
            <span class="text-sm">{{ parseDateTime(project.CreatedAt) }}</span>
          </div>
          <div class="flex w-full h-full animate__animated animate__fadeIn animate__faster
          bg-gray-100/50 active:bg-gray-200/90 dark:bg-gray-950/10 dark:active:bg-gray-950/15 cursor-pointer" v-else>

            <span class="mx-auto my-auto flex flex-nowrap">
              <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5"
                   stroke="currentColor" class="size-5">
  <path stroke-linecap="round" stroke-linejoin="round"
        d="M15.59 14.37a6 6 0 0 1-5.84 7.38v-4.8m5.84-2.58a14.98 14.98 0 0 0 6.16-12.12A14.98 14.98 0 0 0 9.631 8.41m5.96 5.96a14.926 14.926 0 0 1-5.841 2.58m-.119-8.54a6 6 0 0 0-7.381 5.84h4.8m2.581-5.84a14.927 14.927 0 0 0-2.58 5.84m2.699 2.7c-.103.021-.207.041-.311.06a15.09 15.09 0 0 1-2.448-2.448 14.9 14.9 0 0 1 .06-.312m-2.24 2.39a4.493 4.493 0 0 0-1.757 4.306 4.493 4.493 0 0 0 4.306-1.758M16.5 9a1.5 1.5 0 1 1-3 0 1.5 1.5 0 0 1 3 0Z"/>
</svg>

              进入项目空间</span>
          </div>
        </div>
        <div class="border border-dashed rounded-xl dark:border-[rgb(118,118,118)] outline-[1px] min-h-56 place-items-center place-content-center text-gray-400
hover:bg-gray-100/50 active:bg-gray-100/90 dark:hover:bg-gray-900/5 dark:active:bg-gray-900/10 cursor-pointer"
             @click="router.push('/workspace/newProject')">
          <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5"
               stroke="currentColor" class="size-16">
            <path stroke-linecap="round" stroke-linejoin="round" d="M12 9v6m3-3H9m12 0a9 9 0 1 1-18 0 9 9 0 0 1 18 0Z"/>
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