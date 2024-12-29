<script setup>
import {onMounted, reactive} from "vue";
import {get} from "@/util/request.js";

const options = reactive({
  myTeams: [],
  myJoinedTeams: []
})
const fetchMyTeams = () => {
  get('/api/team/getMyTeams', {},
      (message,data) => {
        options.myTeams = data;
      })
}
const fetchMyJoinedTeams = () => {
  get('/api/team/getMyJoinedTeams', {},
      (message,data) => {
        options.myJoinedTeams = data;
      })
}
onMounted(() => {
  fetchMyTeams()
  fetchMyJoinedTeams()
})
</script>

<template>
  <div class="flex flex-col gap-2">
    <div class="border p-4 workspace-box w-fit">
      <h1 class="text-2xl">简述</h1>
      <span class="text-sm font-bold">我们的团队功能支持您将创剧空间中的项目与您的团队绑定，使您的团队成员一同加入您的精彩内容创作！<br/>我们支持:项目多人分工协作/实时评论反馈/项目版本控制</span>
    </div>
    <div class="border p-4 workspace-box">
      <h1 class="text-2xl">我管理的团队</h1>
      <div class="w-full grid grid-cols-5 gap-2">
        <div class="border border-dashed rounded-xl dark:border-[#cbcbc4] outline-[1px] min-h-56 place-items-center place-content-center text-gray-400
hover:bg-gray-100/50 active:bg-gray-100/90 dark:hover:bg-gray-900/5 dark:active:bg-gray-900/10 cursor-pointer">
            <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="size-16">
              <path stroke-linecap="round" stroke-linejoin="round" d="M12 9v6m3-3H9m12 0a9 9 0 1 1-18 0 9 9 0 0 1 18 0Z" />
            </svg>
        </div>
      </div>

    </div>
    <div class="border p-4 workspace-box">
      <h1 class="text-2xl">我加入的团队</h1>
      {{options}}
    </div>
  </div>

</template>

<style scoped>

</style>