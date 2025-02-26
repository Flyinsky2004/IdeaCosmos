<script setup>
import {reactive} from "vue";
import ThemeSwitch from "@/components/button/ThemeSwitcher.vue";
import {useUserStore} from "@/stores/user.js";
import router from "@/router/index.js";
import {message} from "ant-design-vue";
import {useRoute} from "vue-router";
import {useThemeStore} from "@/stores/theme.js";
import { BACKEND_DOMAIN } from "@/util/VARRIBLES";

const [messageApi, contextHolder] = message.useMessage();
const route = useRoute()
const userStore = useUserStore()

const options = reactive({
  sideExpand: false,
  sideTextExpand: false
})
const category = [
  {
    path: '/admin/dashboard',
    name: '管理面板',
    icon: `<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="size-6">
  <path stroke-linecap="round" stroke-linejoin="round" d="M3.75 3v11.25A2.25 2.25 0 0 0 6 16.5h2.25M3.75 3h-1.5m1.5 0h16.5m0 0h1.5m-1.5 0v11.25A2.25 2.25 0 0 1 18 16.5h-2.25m-7.5 0h7.5m-7.5 0-1 3m8.5-3 1 3m0 0 .5 1.5m-.5-1.5h-9.5m0 0-.5 1.5M9 11.25v1.5M12 11.25v1.5m3-1.5v1.5" />
</svg>
`
  }, {
    path: '/admin/users',
    name: '用户管理',
    icon: `<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="size-6">
  <path stroke-linecap="round" stroke-linejoin="round" d="M15 19.128a9.38 9.38 0 0 0 2.625.372 9.337 9.337 0 0 0 4.121-.952 4.125 4.125 0 0 0-7.533-2.493M15 19.128v-.003c0-1.113-.285-2.16-.786-3.07M15 19.128v.106A12.318 12.318 0 0 1 8.624 21c-2.331 0-4.512-.645-6.374-1.766l-.001-.109a6.375 6.375 0 0 1 11.964-3.07M12 6.375a3.375 3.375 0 1 1-6.75 0 3.375 3.375 0 0 1 6.75 0Zm8.25 2.25a2.625 2.625 0 1 1-5.25 0 2.625 2.625 0 0 1 5.25 0Z" />
</svg>
`
  }, {
    path: '/admin/chapters',
    name: '篇章管理',
    icon: `<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="size-6">
  <path stroke-linecap="round" stroke-linejoin="round" d="M12 6.042A8.967 8.967 0 0 0 6 3.75c-1.052 0-2.062.18-3 .512v14.25A8.987 8.987 0 0 1 6 18c2.305 0 4.408.867 6 2.292m0-14.25a8.966 8.966 0 0 1 6-2.292c1.052 0 2.062.18 3 .512v14.25A8.987 8.987 0 0 0 18 18a8.967 8.967 0 0 0-6 2.292m0-14.25v14.25" />
</svg>
`
  }, {
    path: '/admin/projects',
    name: '项目管理',
    icon: `<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="size-6">
  <path stroke-linecap="round" stroke-linejoin="round" d="m21 7.5-9-5.25L3 7.5m18 0-9 5.25m9-5.25v9l-9 5.25M3 7.5l9 5.25M3 7.5v9l9 5.25m0-9v9" />
</svg>
`
  }, {
    path: '/admin/statistics',
    name: '数据统计',
    icon: `<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="size-6">
  <path stroke-linecap="round" stroke-linejoin="round" d="M3 13.125C3 12.504 3.504 12 4.125 12h2.25c.621 0 1.125.504 1.125 1.125v6.75C7.5 20.496 6.996 21 6.375 21h-2.25A1.125 1.125 0 0 1 3 19.875v-6.75ZM9.75 8.625c0-.621.504-1.125 1.125-1.125h2.25c.621 0 1.125.504 1.125 1.125v11.25c0 .621-.504 1.125-1.125 1.125h-2.25a1.125 1.125 0 0 1-1.125-1.125V8.625ZM16.5 4.125c0-.621.504-1.125 1.125-1.125h2.25C20.496 3 21 3.504 21 4.125v15.75c0 .621-.504 1.125-1.125 1.125h-2.25a1.125 1.125 0 0 1-1.125-1.125V4.125Z" />
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
      <img class="rounded-xl" :src="BACKEND_DOMAIN + userStore.user.avatar" alt="用户头像">
      <span v-if="options.sideExpand" class="text-sm font-bold">{{ userStore.user.username }} <span class="text-red-500">(管理员)</span></span>
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
        class="h-screen overflow-auto w-full container1 text-theme-switch p-0 relative bg-slate-50 dark:bg-slate-950 before:absolute before:inset-0 before:bg-glow-effect before:blur-3xl before:z-[-1]">
      <div class="min-h-screen w-full dark:bg-gray-800/30 bg-gray-100/50 backdrop-blur-xl m-0 p-2">
        <RouterView/>
      </div>
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
  background-size: 100% 100%;
  background-attachment: fixed;
  background-position: center;
  background-blend-mode: multiply;
}
.darkBKG{
  background-image: url('@/assets/img/darkbkg.png');
  background-color: rgba(80, 79, 79, 0.3);
}
.lightBKG{
  background-image: url('@/assets/img/lightbkg.png');
  background-color: rgba(221, 221, 221, 0.1);
}
</style>
