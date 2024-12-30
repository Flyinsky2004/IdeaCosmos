<script setup>
import {onMounted, reactive} from "vue";
import {get, postJSON} from "@/util/request.js";
import {message} from "ant-design-vue";
import {parseDateTime} from "@/util/common.js";

const options = reactive({
  myTeams: [],
  myJoinedTeams: [],
  isAddWindowOpen: false,
})
const fetchMyTeams = () => {
  get('/api/team/getMyTeams', {},
      (message, data) => {
        options.myTeams = data;
      })
}
const fetchMyJoinedTeams = () => {
  get('/api/team/getMyJoinedTeams', {},
      (message, data) => {
        options.myJoinedTeams = data;
      })
}
onMounted(() => {
  fetchMyTeams()
  fetchMyJoinedTeams()
})
const newTeamForm = reactive({
  team_name: '',
  team_description: '',
})
const submitNewTeam = () => {
  if (newTeamForm.team_name.length < 6 && newTeamForm.team_description.length < 20) {
    message.info("团队名称长度必须大于6个字符且团队描述必须大于20个字符。")
  } else {
    postJSON('/api/team/createTeam', newTeamForm,
        (messager, data) => {
          message.success(messager)
          fetchMyTeams()
          options.isAddWindowOpen = false
        },
        (messager, data) => {
          message.warning(messager)
        },
        (messager, data) => {
          message.error(messager)
        })
  }
}
</script>

<template>
  <a-modal v-model:open="options.isAddWindowOpen" title="创建新团队" @ok="submitNewTeam" ok-text="创建"
           cancel-text="取消">
    <div class="font-sans grid">
      <span>团队名称:</span><input class="input1" v-model="newTeamForm.team_name">
      <span>团队描述:</span><textarea class="input1" v-model="newTeamForm.team_description"/>
    </div>
  </a-modal>
  <div class="flex flex-col gap-2 animate__animated animate__fadeIn">
    <div class="border p-4 workspace-box w-fit">
      <h1 class="text-2xl font-serif">团队工作</h1>
      <span class="text-sm font-bold">我们的团队功能支持您将创剧空间中的项目与您的团队绑定，使您的团队成员一同加入您的精彩内容创作！<br/>我们支持:项目多人分工协作/实时评论反馈/项目版本控制</span>
    </div>
    <div class="border p-4 workspace-box">
      <h1 class="text-2xl">我管理的团队</h1>
      <div class="w-full grid grid-cols-5 gap-2">
        <div v-for="team in options.myTeams" class="border rounded-xl dark:border-[rgb(118,118,118)] outline-[1px] min-h-56 text-theme-switch p-4
hover:bg-gray-100/50 active:bg-gray-100/90 dark:hover:bg-gray-900/5 dark:active:bg-gray-900/10 cursor-pointer">
          <h1 class="text-2xl text-blue-500">#{{ team.ID }}{{ team.username }}</h1>
          <h1 class="text-sm font-bold text-blue-400">团队介绍：</h1>

          <h1 class="text-sm">{{ team.teamDescription }}</h1>
          <h1 class="text-sm font-bold text-blue-400">创建时间：</h1>

          <span class="text-sm">{{ parseDateTime(team.CreatedAt) }}</span>
        </div>
        <div class="border border-dashed rounded-xl dark:border-[rgb(118,118,118)] outline-[1px] min-h-56 place-items-center place-content-center text-gray-400
hover:bg-gray-100/50 active:bg-gray-100/90 dark:hover:bg-gray-900/5 dark:active:bg-gray-900/10 cursor-pointer"
             @click="options.isAddWindowOpen = true">
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