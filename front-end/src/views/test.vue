<template>
  <div>
    <button @click="downloadWebp('20250106_1ad83f.webp')">下载并转换</button>
    <img v-if="base64Image" :src="base64Image" alt="Converted WebP">
  </div>
</template>

<script>
import { ref } from 'vue';
import { imagePrefix } from "@/util/VARRIBLES.js";
import {get} from "@/util/request.js";
import {message} from "ant-design-vue";

export default {
  setup() {
    const base64Image = ref('');

    const downloadWebp = (path) => {
      get('/api/user/getWebpImageBase64',{
        filename: path,
      },(messager,data) => {
        base64Image.value = data;
      },(messager,data) => {
        message.warn(messager)
      },(messager,data) => {
        message.error(messager)
      })
    };

    return {
      downloadWebp,
      base64Image
    };
  }
};
</script>