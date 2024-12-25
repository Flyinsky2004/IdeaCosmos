import './assets/css/main.css'
import 'animate.css';
import "@flaticon/flaticon-uicons/css/all/all.css";

import {createApp} from 'vue'
import {createPinia} from 'pinia'
import ElementPlus from 'element-plus'
import VueKinesis from 'vue-kinesis'

import 'element-plus/dist/index.css'
import App from './App.vue'
import router from './router'

const app = createApp(App)
app.use(createPinia())
app.use(router).use(ElementPlus).use(VueKinesis)

app.mount('#app')
