<script setup>
import {reactive} from "vue";
import ThemeSwitch from "@/components/button/ThemeSwitcher.vue";
import {useUserStore} from "@/stores/user.js";
import router from "@/router/index.js";
import {message} from "ant-design-vue";
import {useRoute} from "vue-router";
import {useThemeStore} from "@/stores/theme.js";

const [messageApi, contextHolder] = message.useMessage();
const route = useRoute()
const userStore = useUserStore()

const options = reactive({
  sideExpand: false,
  sideTextExpand: false
})
const category = [
  {
    path: '/workspace/dataAnlysis',
    name: '热点看板',
    icon: `<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="size-6">
  <path stroke-linecap="round" stroke-linejoin="round" d="M10.5 6a7.5 7.5 0 1 0 7.5 7.5h-7.5V6Z" />
  <path stroke-linecap="round" stroke-linejoin="round" d="M13.5 10.5H21A7.5 7.5 0 0 0 13.5 3v7.5Z" />
</svg>
`
  }, {
    path: '/workspace/projects',
    name: '创剧空间',
    icon: `<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="size-6">
  <path stroke-linecap="round" stroke-linejoin="round" d="m21 7.5-2.25-1.313M21 7.5v2.25m0-2.25-2.25 1.313M3 7.5l2.25-1.313M3 7.5l2.25 1.313M3 7.5v2.25m9 3 2.25-1.313M12 12.75l-2.25-1.313M12 12.75V15m0 6.75 2.25-1.313M12 21.75V19.5m0 2.25-2.25-1.313m0-16.875L12 2.25l2.25 1.313M21 14.25v2.25l-2.25 1.313m-13.5 0L3 16.5v-2.25" />
</svg>

`
  }, {
    path: '/workspace/personalInfo',
    name: '个人信息',
    icon: `
    <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="size-6">
  <path stroke-linecap="round" stroke-linejoin="round" d="M15.75 6a3.75 3.75 0 1 1-7.5 0 3.75 3.75 0 0 1 7.5 0ZM4.501 20.118a7.5 7.5 0 0 1 14.998 0A17.933 17.933 0 0 1 12 21.75c-2.676 0-5.216-.584-7.499-1.632Z" />
</svg>

`
  }, {
    path: '/workspace/teams',
    name: '我的团队',
    icon: `<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="size-6">
  <path stroke-linecap="round" stroke-linejoin="round" d="M18 18.72a9.094 9.094 0 0 0 3.741-.479 3 3 0 0 0-4.682-2.72m.94 3.198.001.031c0 .225-.012.447-.037.666A11.944 11.944 0 0 1 12 21c-2.17 0-4.207-.576-5.963-1.584A6.062 6.062 0 0 1 6 18.719m12 0a5.971 5.971 0 0 0-.941-3.197m0 0A5.995 5.995 0 0 0 12 12.75a5.995 5.995 0 0 0-5.058 2.772m0 0a3 3 0 0 0-4.681 2.72 8.986 8.986 0 0 0 3.74.477m.94-3.197a5.971 5.971 0 0 0-.94 3.197M15 6.75a3 3 0 1 1-6 0 3 3 0 0 1 6 0Zm6 3a2.25 2.25 0 1 1-4.5 0 2.25 2.25 0 0 1 4.5 0Zm-13.5 0a2.25 2.25 0 1 1-4.5 0 2.25 2.25 0 0 1 4.5 0Z" />
</svg>

`
  }
]
const sideBarMouseOver = () => {
  options.sideExpand = true
  if (!options.sideTextExpand) {
    setTimeout(() => {
      options.sideTextExpand = true
    }, 150)
  }

}
const sideBarMouseLeave = () => {
  options.sideExpand = false
  options.sideTextExpand = false
}
</script>

<template>
  <contextHolder/>
  <div
      class="h-full absolute p-2 bg-[rgba(244,243,242,1)] dark:bg-[rgba(36,36,36,0.6)] transition-all duration-300 text-theme-switch animate__animated animate__fadeIn"
      :class="options.sideExpand ? 'sideAfter' : 'sideBefore'"
      @mouseover="sideBarMouseOver"
      @mouseleave="sideBarMouseLeave"
  >
    <div class="w-full bg-[#ffffff] dark:bg-[rgb(18,18,18)] p-2 rounded-xl text-center">
      <img :src="userStore.user.avator" alt="用户头像">
      <span v-if="options.sideExpand" class="text-sm font-bold">{{ userStore.user.username }}</span>
    </div>
    <div class="mt-2 theme-transition grid w-full my-auto  p-2 rounded-xl
         select-none cursor-pointer hover:bg-[rgb(250,250,250)] dark:hover:bg-[rgb(10,10,10)]"
         :class="[options.sideExpand? 'bg-[#ffffff] dark:bg-[rgb(18,18,18)]':'']">
      <div v-if="options.sideTextExpand"
           class="grid grid-cols-[1fr,1fr] animate__animated animate__fadeIn place-items-center">
        <ThemeSwitch/>
        <div class="bg-red-500 w-fit rounded-3xl font-bold text-white p-1 select-none cursor-pointer"
             @click="router.push('/')"
        >
          <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5"
               stroke="currentColor" class="size-6">
            <path stroke-linecap="round" stroke-linejoin="round"
                  d="M8.25 9V5.25A2.25 2.25 0 0 1 10.5 3h6a2.25 2.25 0 0 1 2.25 2.25v13.5A2.25 2.25 0 0 1 16.5 21h-6a2.25 2.25 0 0 1-2.25-2.25V15m-3 0-3-3m0 0 3-3m-3 3H15"/>
          </svg>
        </div>
      </div>
      <div v-else>
        <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor"
             class="size-6 mx-auto">
          <path stroke-linecap="round" stroke-linejoin="round"
                d="M8.625 12a.375.375 0 1 1-.75 0 .375.375 0 0 1 .75 0Zm0 0H8.25m4.125 0a.375.375 0 1 1-.75 0 .375.375 0 0 1 .75 0Zm0 0H12m4.125 0a.375.375 0 1 1-.75 0 .375.375 0 0 1 .75 0Zm0 0h-.375M21 12a9 9 0 1 1-18 0 9 9 0 0 1 18 0Z"/>
        </svg>
      </div>
    </div>
    <div
        :class="[route.path === cat.path ? 'bg-[#ffffff] dark:bg-[rgb(18,18,18)] text-blue-600 dark:text-purple-400': '',options.sideTextExpand?'grid-cols-[1fr,2fr]':'grid-cols-[1fr]']"
        class="mt-2 theme-transition grid w-full my-auto  p-2 rounded-xl text-hover
         select-none cursor-pointer hover:bg-[rgb(250,250,250)] dark:hover:bg-[rgb(10,10,10)]"
        v-for="cat in category" @click="router.push(cat.path)">
      <div class="mx-auto" v-html="cat.icon"/>
      <span class="animate__animated animate__fadeIn my-auto" v-if="options.sideTextExpand">{{ cat.name }}</span>
    </div>
  </div>

  <div class="w-full h-full basicBkgColorSwitch mx-auto rounded-2xl grid basicColorSwitch transition-all"
       :class="options.sideExpand ? 'grid-cols-[1fr,8fr]' : 'grid-cols-[1fr,19fr]' "
  >
    <div></div>
    <div :class="useThemeStore().currentTheme === 'light' ? 'lightBKG' : 'darkBKG'"
        class="h-screen overflow-auto w-full container1 text-theme-switch p-4 relative bg-slate-50 dark:bg-slate-950 before:absolute before:inset-0 before:bg-glow-effect before:blur-3xl before:z-[-1]">
      <RouterView/>
    </div>
  </div>
</template>

<style scoped>
.sideAfter {
  width: 11%;
}

.sideBefore {
  width: 5%;
}
.container1{
  background-size: 100% 100%; /* 确保背景图片覆盖整个容器 */
  background-attachment: fixed; /* 背景图片固定不动 */
  background-position: center; /* 背景图片居中 */
  background-blend-mode: multiply;
}
.darkBKG{
  background-image: url('@/assets/img/darkbkg.png'); /* 使用背景图片 */
  background-color: rgba(80, 79, 79, 0.3);

}
.lightBKG{
  background-image: url('@/assets/img/lightbkg.png'); /* 使用背景图片 */
  background-color: rgba(221, 221, 221, 0.1);
}
</style>