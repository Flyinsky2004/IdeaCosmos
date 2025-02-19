import './assets/css/main.css'
import 'animate.css';
import "@flaticon/flaticon-uicons/css/all/all.css";
import Antd from 'ant-design-vue';
import '@/assets/css/remixicon.css';
import 'ant-design-vue/dist/reset.css';

import {createApp} from 'vue'
import {createPinia} from 'pinia'
import ElementPlus from 'element-plus'
import VueKinesis from 'vue-kinesis'

import 'element-plus/dist/index.css'
import App from './App.vue'
import router from './router'
import axios from "axios";

axios.defaults.baseURL="http://localhost:8080"
//axios.defaults.baseURL="https://ic.flyinsky.wiki"
const app = createApp(App)
app.use(createPinia())
app.use(router).use(ElementPlus).use(VueKinesis).use(Antd)

app.mount('#app')
