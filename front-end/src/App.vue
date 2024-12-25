<script setup>
import {RouterView} from "vue-router";
import {useThemeStore} from "@/stores/theme.js";
const themeStore = useThemeStore()
themeStore.initTheme()
</script>

<template>
  <div class="container1 select-none h-screen darkBKG">
<!--       :class="themeStore.currentTheme === 'light' ? 'lightBKG' : 'darkBKG'">-->
    <div class="h-full backdrop-blur-sm overflow-y-auto">
      <RouterView/>
    </div>
  </div>
</template>

<style scoped>
.container1{
  font-family: 'srchan', serif;
  background-size: 100% 100%; /* 确保背景图片覆盖整个容器 */
  background-attachment: fixed; /* 背景图片固定不动 */
  background-position: center; /* 背景图片居中 */
  background-blend-mode: multiply;
}
.darkBKG{
  background-image: url('@/assets/img/bkg.png'); /* 使用背景图片 */
  background-color: rgba(128,128,128,0.95);
  background-blend-mode: multiply;

}
.lightBKG{
  background-image: url('@/assets/img/light-bkg.webp'); /* 使用背景图片 */
}
</style>
