<script setup>

import ThemeSwitcher from "@/components/button/ThemeSwitcher.vue";
import logo from '@/assets/img/logo.webp'
import router from "@/router/index.js";
import LoginForm from "@/components/auth/loginForm.vue";
import example from '@/assets/img/example.png'
import {useRoute} from "vue-router";
import RegisterForm from "@/components/auth/registerForm.vue";
const route = useRoute()
const topBarItems = [
  {
    name: '作品画廊',
    path: '/',
    icon: `<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="size-5">
  <path stroke-linecap="round" stroke-linejoin="round" d="M2.25 7.125C2.25 6.504 2.754 6 3.375 6h6c.621 0 1.125.504 1.125 1.125v3.75c0 .621-.504 1.125-1.125 1.125h-6a1.125 1.125 0 0 1-1.125-1.125v-3.75ZM14.25 8.625c0-.621.504-1.125 1.125-1.125h5.25c.621 0 1.125.504 1.125 1.125v8.25c0 .621-.504 1.125-1.125 1.125h-5.25a1.125 1.125 0 0 1-1.125-1.125v-8.25ZM3.75 16.125c0-.621.504-1.125 1.125-1.125h5.25c.621 0 1.125.504 1.125 1.125v2.25c0 .621-.504 1.125-1.125 1.125h-5.25a1.125 1.125 0 0 1-1.125-1.125v-2.25Z" />
</svg>
`
  }, {
    name: '订阅计划',
    path: '/',
    icon: `<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="size-5">
  <path stroke-linecap="round" stroke-linejoin="round" d="M15.75 5.25a3 3 0 0 1 3 3m3 0a6 6 0 0 1-7.029 5.912c-.563-.097-1.159.026-1.563.43L10.5 17.25H8.25v2.25H6v2.25H2.25v-2.818c0-.597.237-1.17.659-1.591l6.499-6.499c.404-.404.527-1 .43-1.563A6 6 0 1 1 21.75 8.25Z" />
</svg>
`
  }
]
</script>

<template>
  <div class="fixed top-0 left-0 w-full shadow-md z-50 flex flex-nowrap place-items-center text-theme-switch font-sans border-b-[1px] theme-border bkg-theme-switch">
    <div class="w-1/12"></div>
    <div class="flex flex-nowrap p-4 cursor-pointer" @click="router.push('/')">
      <img :src="logo" class="w-8 h-8 rounded-full"/>
      <span class="my-auto ml-4">创剧星球</span>
    </div>
    <div class="flex-grow"/>
    <div class="h-full flex gap-4">
      <div class="text-hover flex flex-nowrap cursor-pointer transition-all p-2" v-for="item in topBarItems" @click="router.push(item.path)">
        <div v-html="item.icon" class="h-6 w-6 rounded-full"/>
        <span>{{item.name}}</span>
      </div>
      <ThemeSwitcher class="mx-auto my-auto"/>
    </div>
    <div class="w-1/12"></div>
  </div>
  <div class="grid grid-cols-[8fr,7fr] relative w-full h-full overflow-hidden bg-white dark:bg-[#1a1a1a]">
    <kinesis-container class="h-screen animate__animated animate__fadeInDown">
      <kinesis-element :strength="40">
        <img :src="example" class="object-cover rotate-12 input1">
      </kinesis-element>
    </kinesis-container>
    <div class="place-content-center place-items-center animate__animated animate__fadeIn">
      <div class="box1 p-8 rounded-xl transition-all">
        <LoginForm v-if="route.path === '/auth/login'"/>
        <RegisterForm v-else/>
      </div>
    </div>
  </div>

</template>

<style scoped>
</style>