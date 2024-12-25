<script setup>

import router from "@/router/index.js";
</script>

<template>
  <div class="font-sans text-theme-switch">
    <div class="text-center text-4xl">注册</div>
    <div class="text-center text-xl mt-4">欢迎来到创剧星球</div>
    <div class="w-full grid grid-cols-[1fr,2fr] mt-2 place-items-center">
      <span class="text-xl">用户名:</span>
      <input class="input1"/>
    </div>
    <div class="w-full grid grid-cols-[1fr,2fr] mt-2 place-items-center">
      <span class="text-xl">密码:</span>
      <input class="input1"/>
    </div>
    <div class="w-full flex mt-4">
      <button class="mx-auto btn1">
        登录
      </button>
    </div>
    <div class="mt-2 flex"><span class="mx-auto">已有账号?<a class="text-hover cursor-pointer hover:underline" @click="router.push('/auth/login')">马上登录</a></span></div>
  </div>
</template>

<style scoped>

</style>