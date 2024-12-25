<script setup>

import ThemeSwitcher from "@/components/ThemeSwitcher.vue";
import logo from '@/assets/logo.webp'
</script>

<template>
  <div class="h-fit flex flex-nowrap text-theme-switch font-sans border-b-[1px] theme-border">
    <div class="flex flex-nowrap p-4">
      <img :src="logo" class="w-8 h-8 rounded-full" />
      <span class="my-auto ml-4">创剧星球</span>
    </div>
  </div>
  <div class="h-1/3 flex flex-col font-sans mt-32">
    <div class="mx-auto text-theme-switch">
      <span
          class="text-7xl font-bold bg-gradient-to-b from-gray-100 via-gray-300 to-gray-400 text-transparent bg-clip-text">
        内容创作者的
      </span>
    </div>
    <div class="mx-auto text-theme-switch">
      <span
          class="text-6xl font-bold bg-gradient-to-b from-gray-400 via-gray-500 to-gray-600 text-transparent bg-clip-text">
        得力助手
      </span>
    </div>
    <div class="mx-auto text-[#a9a9a9] mt-8 text-xl">
      一款集灵感激发、角色设计、剧本生成、团队协作于一体的智能创作平台，
    </div>
    <div class="mx-auto text-[#a9a9a9] mt-2">
      支持版本控制、版权管理及内容优化输出，助力创作者高效创作。
    </div>
    <div class="mx-auto mt-4 flex flex-wrap gap-2">
      <button class="color-mixed-button font-light hover:-translate-y-0.5 transition-all">立刻上手</button>
      <button class="transparent-button font-light hover:-translate-y-0.5 transition-all">用户评价</button>

    </div>

  </div>
  <div class="text-theme-switch">


    <ThemeSwitcher/>
  </div>
</template>

<style scoped>
</style>