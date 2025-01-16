<script setup>
import {message} from 'ant-design-vue';
import router from "@/router/index.js";
import {reactive} from "vue";
import {postJSON} from "@/util/request.js";

const [messageApi, contextHolder] = message.useMessage();

const loginForm = reactive({
  username: "",
  password: ""
})
const loginHandler = () => {
  if (loginForm.username.length < 6 || loginForm.password.length < 6) {
    messageApi.warning('账号或密码不能少于6个字符哦~')
    return
  }
  postJSON('/api/auth/login',
      loginForm,
      (message,data) => {
        messageApi.success("登陆成功！欢迎回来！")
        localStorage.setItem("authToken", data)
      },(message,data) => {
        messageApi.warning(message)
      },(message,data) => {
        messageApi.warning(message)
      })
}
</script>

<template>
  <contextHolder/>
  <div class="font-sans text-theme-switch">
    <div class="text-center text-4xl">登录</div>
    <div class="text-center text-xl mt-4">欢迎来到创剧星球</div>
    <div class="w-full grid grid-cols-[1fr,2fr] mt-2 place-items-center">
      <span class="text-xl">用户名:</span>
      <input class="input1" v-model="loginForm.username"/>
    </div>
    <div class="w-full grid grid-cols-[1fr,2fr] mt-2 place-items-center">
      <span class="text-xl">密码:</span>
      <input type="password" class="input1" v-model="loginForm.password"/>
    </div>
    <div class="w-full flex mt-4">
      <button class="mx-auto btn1" @click="loginHandler()">
        登录
      </button>
    </div>
    <div class="mt-2 flex"><span class="mx-auto">还没有账号?<a class="text-hover cursor-pointer hover:underline"
                                                               @click="router.push('/auth/register')">注册一个</a></span>
    </div>
  </div>
</template>

<style scoped>

</style>