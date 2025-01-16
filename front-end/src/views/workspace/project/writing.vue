<script setup>
import {reactive} from "vue";

const chapter = JSON.parse(localStorage.getItem('chapter'))
const options = reactive({
  historyVersions: [],
  currentVersion: {}
})
</script>

<template>
  <div class="h-full w-full grid grid-cols-[4fr,1fr] gap-2">
    <div class="border theme-border p-2 rounded-xl">
      <div>
        <span class="text-gray-600">此篇章还没有内容哦</span>
      </div>
      <div>
        <span class="text-gray-600">此篇章还没有内容哦</span>
      </div>
    </div>
    <div class="border theme-border p-2 rounded-xl">
      <h1 class="text-3xl text-blue-500">历史版本</h1>
      <div class="" v-if="options.historyVersions.length === 0">
        <span class="text-gray-600">还没有历史版本哦</span>
      </div>
    </div>
  </div>
{{chapter}}
</template>

<style scoped>

</style>