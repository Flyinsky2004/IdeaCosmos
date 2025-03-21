<script setup>
import {onMounted, reactive, ref} from "vue";
import {audioModels, BACKEND_DOMAIN, femaleAvatar, FRONTEND_DOMAIN, maleAvatar} from "@/util/VARRIBLES.js";
import {get, post} from "@/util/request.js";
import {message} from "ant-design-vue";
import {parseDateTime} from "@/util/common.js";
import {useUserStore} from "@/stores/user.js";
import axios from "axios";
import router from "@/router/index.js";
import {exportScript} from "@/util/exportTool.js";
import SpinLoaderLarge from "@/components/spinLoaderLarge.vue";
const userStore = useUserStore()
const project = JSON.parse(localStorage.getItem("project"));
const options = reactive({
  currentChapter: {},
  currentAudioName: '',
  chapters: [],
  currentAudioFileName: '',
  isAudioPlaying: false,
  audioUrl: '',
  activeTab: 'chapters',
  audioGenerating: false,
  exportLoading: false,
  loading: true,
  videoGenerating: false,
  sceneGenerating: false,
  imagesGenerating: false,
  scenes: [],
  videoPath: '',
})
onMounted(() => {
  fetchChapters()
})
const fetchChapters = () => {
  get('/api/project/getAllChapters', {
    project_id: project.ID
  }, (messageer, data) => {
    options.chapters = data
    options.loading = false;
  }, (messageer, data) => {
    message.warning(messageer)
  }, (messageer, data) => {
    message.error(messageer)
  })
}
const chapterClickHandler = (chapter, versionID) => {
  if (versionID !== 0) {
    options.currentChapterId = chapter.ID
    options.currentChapter = chapter
    if (chapter.current_version !== undefined) {
      fetchAudioResource()
      // 获取场景数据
      getScenes()
      // 设置视频路径
      if (chapter.current_version.video_path) {
        options.videoPath = chapter.current_version.video_path
      } else {
        options.videoPath = ''
      }
    }
  } else {
    message.info("该篇章还没有创作版本，请前往创作板块创作！")
  }
}
const generateAudio = async () => {
  if (!options.currentChapterId || !options.currentAudioName) {
    message.warning('请先选择篇章和音频模型');
    return;
  }
  
  options.audioGenerating = true;
  try {
    await new Promise((resolve, reject) => {
      get('/api/project/generateChapterAudio', {
        chapterId: options.currentChapterId,
        audioName: options.currentAudioName
      }, (messageer, data) => {
        options.currentAudioFileName = data;
        message.success('音频生成成功');
        resolve(data);
      }, (messageer) => reject(messageer),
         (messageer) => reject(messageer)
      );
    });
  } catch (error) {
    message.error('音频生成失败');
  } finally {
    options.audioGenerating = false;
    chapterClickHandler(options.currentChapter, options.currentChapterId)
  }
};
const playAudio = () => {
  options.isAudioPlaying = true
}
const stopAudio = () => {
  options.isAudioPlaying = false
}
const fetchAudioResource = async () => {
  // 发送请求，获取MP3文件流
  const response = await axios.get(BACKEND_DOMAIN + 'audio/' + options.currentChapter.current_version.audio_path, {
    responseType: 'arraybuffer', // 确保获取的是二进制数据
  });

  // 将 ArrayBuffer 转化为 Blob
  const blob = new Blob([response.data], {type: 'audio/mpeg'});

  // 将 Blob 转化为临时URL
  options.audioUrl = URL.createObjectURL(blob);

}

// 生成视频场景
const generateScenes = () => {
  if (!options.currentChapterId) {
    message.warning('请先选择篇章');
    return;
  }
  
  options.sceneGenerating = true;
  get('/api/video/generateScene', {
    chapter_verison_id: options.currentChapter.current_version.ID
  }, (messageer, data) => {
    options.scenes = data;
    message.success('场景生成成功');
    options.sceneGenerating = false;
  }, (messageer) => {
    message.warning(messageer);
    options.sceneGenerating = false;
  }, (messageer) => {
    message.error(messageer);
    options.sceneGenerating = false;
  });
}

// 获取现有场景
const getScenes = () => {
  if (!options.currentChapterId || !options.currentChapter.current_version) {
    return;
  }
  
  get('/api/video/getSceneByChapterVersionID', {
    chapter_verison_id: options.currentChapter.current_version.ID
  }, (messageer, data) => {
    options.scenes = data;
  }, (messageer) => {
    message.warning(messageer);
  }, (messageer) => {
    message.error(messageer);
  });
}

// 生成场景图片
const generateImages = () => {
  if (!options.currentChapterId || !options.scenes || options.scenes.length === 0) {
    message.warning('请先生成场景');
    return;
  }
  
  options.imagesGenerating = true;
  get('/api/video/generateChapterImages', {
    chapter_verison_id: options.currentChapter.current_version.ID
  }, (messageer, data) => {
    options.scenes = data;
    message.success('图片生成成功');
    options.imagesGenerating = false;
  }, (messageer) => {
    message.warning(messageer);
    options.imagesGenerating = false;
  }, (messageer) => {
    message.error(messageer);
    options.imagesGenerating = false;
  });
}

// 生成视频
const generateVideo = () => {
  if (!options.currentChapterId || !options.scenes || options.scenes.length === 0) {
    message.warning('请先生成场景和图片');
    return;
  }
  
  options.videoGenerating = true;
  get('/api/video/generateChapterVideo', {
    chapter_verison_id: options.currentChapter.current_version.ID
  }, (messageer, data) => {
    options.videoPath = data;
    options.currentChapter.current_version.video_path = data;
    message.success('视频生成成功');
    options.videoGenerating = false;
  }, (messageer) => {
    message.warning(messageer);
    options.videoGenerating = false;
  }, (messageer) => {
    message.error(messageer);
    options.videoGenerating = false;
  });
}

// 社交媒体分享方法
const shareToDouyin = () => {
  if (!options.videoPath) {
    message.error('请先生成视频');
    return;
  }
  const videoUrl = `${BACKEND_DOMAIN}video/${options.videoPath}`;
  window.open(`https://creator.douyin.com/creator-micro/content/upload`);
}
import bilibili from '@/assets/img/bilibili.ico'
import douyin from '@/assets/img/douyin.ico'
import xiaohongshu from '@/assets/img/xiaohongshu.ico'
import wechatChannels from '@/assets/img/shipinhao.ico'
import qqXiaoshijie from '@/assets/img/qvideo_tob_icon_mcn.png'
const shareToBilibili = () => {
  if (!options.videoPath) {
    message.error('请先生成视频');
    return;
  }
  const videoUrl = `${BACKEND_DOMAIN}video/${options.videoPath}`;
  window.open(`https://member.bilibili.com/platform/upload/video/frame`);
}

const shareToXiaohongshu = () => {
  if (!options.videoPath) {
    message.error('请先生成视频');
    return;
  }
  const videoUrl = `${BACKEND_DOMAIN}video/${options.videoPath}`;
  window.open(`https://creator.xiaohongshu.com/publish/publish`);
}

const shareToWechatChannels = () => {
  if (!options.videoPath) {
    message.error('请先生成视频');
    return;
  }
  const videoUrl = `${BACKEND_DOMAIN}video/${options.videoPath}`;
  window.open(`https://channels.weixin.qq.com/platform/post/create`);
}

const shareToQQXiaoshijie = () => {
  if (!options.videoPath) {
    message.error('请先生成视频');
    return;
  }
  const videoUrl = `${BACKEND_DOMAIN}video/${options.videoPath}`;
  window.open(`https://qqzz.qq.com/publish`);
}
</script>

<template>
  <div class="min-h-screen bg-gray-50 dark:bg-zinc-900/20 p-6">
    <div class="max-w-7xl mx-auto space-y-8">
      <!-- 加载状态 -->
      <div v-if="options.loading" class="flex items-center justify-center py-12">
        <SpinLoaderLarge />
      </div>

      <template v-else>
        <!-- 标签页导航 -->
        <div class="flex gap-4 border-b border-gray-200 dark:border-gray-700">
          <button 
            v-for="tab in [
              { key: 'chapters', text: '选择篇章', icon: 'M12 6.042A8.967 8.967 0 006 3.75c-1.052 0-2.062.18-3 .512v14.25A8.987 8.987 0 016 18c2.305 0 4.408.867 6 2.292' },
              { key: 'audio', text: '音频导出', icon: 'M9 9l10.5-3m0 6.553v3.75a2.25 2.25 0 01-1.632 2.163l-1.32.377a1.803 1.803 0 11-.99-3.467l2.31-.66a2.25 2.25 0 001.632-2.163zm0 0V2.25L9 5.25v10.303m0 0v3.75a2.25 2.25 0 01-1.632 2.163l-1.32.377a1.803 1.803 0 01-.99-3.467l2.31-.66A2.25 2.25 0 009 15.553z' },
              { key: 'video', text: '视频导出', icon: 'M15.75 10.5l4.72-4.72a.75.75 0 011.28.53v11.38a.75.75 0 01-1.28.53l-4.72-4.72M4.5 18.75h9a2.25 2.25 0 002.25-2.25v-9a2.25 2.25 0 00-2.25-2.25h-9A2.25 2.25 0 002.25 7.5v9a2.25 2.25 0 002.25 2.25z' },
              { key: 'export', text: '文件导出', icon: 'M19.5 14.25v-2.625a3.375 3.375 0 00-3.375-3.375h-1.5A1.125 1.125 0 0113.5 7.125v-1.5a3.375 3.375 0 00-3.375-3.375H8.25m2.25 0H5.625c-.621 0-1.125.504-1.125 1.125v17.25c0 .621.504 1.125 1.125 1.125h12.75c.621 0 1.125-.504 1.125-1.125V11.25a9 9 0 00-9-9z' }
            ]"
            :key="tab.key"
            @click="options.activeTab = tab.key"
            class="px-6 py-3 font-medium transition-colors relative"
            :class="options.activeTab === tab.key ? 
              'text-blue-600 dark:text-blue-400' : 
              'text-gray-500 hover:text-gray-700 dark:text-gray-400 dark:hover:text-gray-300'"
          >
            <div class="flex items-center gap-2">
              <svg 
                xmlns="http://www.w3.org/2000/svg" 
                fill="none" 
                viewBox="0 0 24 24" 
                stroke-width="1.5" 
                stroke="currentColor" 
                class="w-5 h-5"
              >
                <path stroke-linecap="round" stroke-linejoin="round" :d="tab.icon" />
              </svg>
              {{ tab.text }}
            </div>
            <div 
              v-if="options.activeTab === tab.key"
              class="absolute bottom-0 left-0 w-full h-0.5 bg-blue-600 dark:bg-blue-400"
            />
          </button>
        </div>

        <!-- 篇章选择 -->
        <div v-show="options.activeTab === 'chapters'" class="space-y-6 animate__animated animate__fadeIn animate__delay-1s">
          <div class="flex items-center justify-between mb-4">
            <h2 class="text-2xl font-bold text-gray-900 dark:text-gray-100">
              选择要导出的篇章
            </h2>
            <div class="flex items-center gap-2 text-sm text-gray-600 dark:text-gray-400">
              <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
              <span>点击卡片选择要导出的篇章</span>
            </div>
          </div>

          <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-6">
            <div
              v-for="chapter in options.chapters"
              :key="chapter.ID"
              class="group relative bg-white dark:bg-zinc-800 rounded-xl overflow-hidden cursor-pointer transition-all duration-300 hover:-translate-y-1"
              :class="[
                chapter.ID === options.currentChapterId ? 
                  'ring-2 ring-blue-500 dark:ring-blue-400 shadow-lg shadow-blue-500/10' : 
                  'border theme-border hover:shadow-xl',
                !chapter.current_version?.ID && 'opacity-50'
              ]"
              @click="chapterClickHandler(chapter, chapter.version_id)"
            >
              <!-- 状态标签 -->
              <div 
                class="absolute top-3 right-3 px-2 py-1 rounded-full text-xs font-medium z-10"
                :class="chapter.current_version?.ID ? 
                  'bg-green-100 text-green-800 dark:bg-green-900/30 dark:text-green-400' : 
                  'bg-gray-100 text-gray-800 dark:bg-gray-700 dark:text-gray-300'"
              >
                {{ chapter.current_version?.ID ? '可导出' : '未创作' }}
              </div>

              <!-- 渐变背景 -->
              <div class="aspect-[4/3] relative">
                <div class="absolute inset-0 bg-gradient-to-b from-black/30 to-black/70"></div>
                
                <!-- 内容 -->
                <div class="absolute inset-0 p-4 flex flex-col justify-between">
                  <div class="space-y-2">
                    <h3 class="text-lg font-semibold text-white group-hover:text-blue-200 transition-colors">
                      {{ chapter.Title }}
                    </h3>
                    <p class="text-sm text-gray-200 line-clamp-2">
                      {{ chapter.Description }}
                    </p>
                  </div>
                  
                  <div class="flex items-center justify-between text-xs">
                    <div class="flex items-center gap-2 text-gray-200">
                      <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                        <path stroke-linecap="round" stroke-linejoin="round" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
                      </svg>
                      <span>版本 {{ chapter.current_version?.ID || '-' }}</span>
                    </div>
                    <div class="flex items-center gap-2 text-gray-200">
                      <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                        <path stroke-linecap="round" stroke-linejoin="round" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
                      </svg>
                      <span>{{ parseDateTime(chapter.UpdatedAt) }}</span>
                    </div>
                  </div>
                </div>
              </div>

              <!-- 选中指示器 -->
              <div 
                v-if="chapter.ID === options.currentChapterId"
                class="absolute inset-0 border-2 border-blue-500 dark:border-blue-400 rounded-xl pointer-events-none"
              ></div>
            </div>
          </div>
        </div>

        <!-- 音频导出 -->
        <div v-show="options.activeTab === 'audio'" class="space-y-6 animate__animated animate__fadeIn animate__delay-2s">
          <div class="flex items-center justify-between">
            <h2 class="text-2xl font-bold text-gray-900 dark:text-gray-100">
              选择音频模型
            </h2>
            <a-button 
              type="primary"
              :loading="options.audioGenerating"
              :disabled="!options.currentChapterId || !options.currentAudioName"
              @click="generateAudio"
            >
              生成音频
            </a-button>
          </div>

          <!-- 音频播放器 -->
          <div v-if="options.currentChapter?.current_version?.audio_path" 
            class="mt-4 p-4 bg-gray-50 dark:bg-zinc-900/50 rounded-xl border border-gray-200 dark:border-zinc-700/80"
          >
            <div class="flex items-center gap-3 mb-2">
              <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 text-blue-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" d="M9 19V6l12-3v13M9 19c0 1.105-1.343 2-3 2s-3-.895-3-2 1.343-2 3-2 3 .895 3 2zm12-3c0 1.105-1.343 2-3 2s-3-.895-3-2 1.343-2 3-2 3 .895 3 2zM9 10l12-3" />
              </svg>
              <span class="text-sm font-medium text-gray-700 dark:text-gray-300">章节音频</span>
            </div>
            <audio
              ref="audioPlayer"
              :src="`${BACKEND_DOMAIN}audio/${options.currentChapter.current_version.audio_path}`"
              class="w-full focus:outline-none"
              controls
              preload="metadata"
              @error="onAudioError"
            ></audio>
          </div>

          <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
            <div
              v-for="(model, index) in audioModels"
              :key="model.model"
              class="group relative p-4 border theme-border rounded-xl transition-all duration-200 cursor-pointer"
              :class="model.model === options.currentAudioName ? 
                'bg-blue-50/70 dark:bg-blue-900/20 border-blue-400 dark:border-blue-800' :
                'hover:bg-gray-50 dark:hover:bg-gray-800/50'"
              @click="options.currentAudioName = model.model"
            >
              <!-- 音频模型卡片内容 -->
              <div class="flex items-start gap-4">
                <div class="shrink-0">
                  <div class="w-12 h-12 rounded-full bg-gradient-to-br p-1"
                       :class="model.isMale ? 
                         'from-blue-400 to-blue-600' : 
                         'from-pink-400 to-rose-500'"
                  >
                    <div class="w-full h-full rounded-full bg-white dark:bg-gray-900">
                      <div v-html="model.isMale ? maleAvatar : femaleAvatar" />
                    </div>
                  </div>
                </div>
                
                <div class="flex-1 min-w-0">
                  <div class="flex items-center gap-2">
                    <h3 class="text-lg font-semibold text-gray-900 dark:text-gray-100">
                      {{ model.name }}
                    </h3>
                    <span class="text-sm px-2 py-0.5 rounded-full"
                          :class="model.isMale ? 
                            'bg-blue-100 text-blue-800' : 
                            'bg-rose-100 text-rose-800'"
                    >
                      {{ model.isMale ? '男声' : '女声' }}
                    </span>
                  </div>

                  <p class="mt-2 text-sm text-gray-600 dark:text-gray-400">
                    {{ model.description }}
                  </p>

                  <div class="mt-3 flex flex-wrap gap-2">
                    <template v-if="Array.isArray(model.style)">
                      <span 
                        v-for="style in model.style" 
                        :key="style"
                        class="text-xs px-2 py-1 rounded-full bg-gray-100 text-gray-700 dark:bg-gray-800 dark:text-gray-300"
                      >
                        {{ style }}
                      </span>
                    </template>
                  </div>
                </div>
              </div>

              <!-- 音频预览 -->
              <div class="mt-4 pt-4 border-t theme-border">
                <audio 
                  :id="`audio-${index}`"
                  class="w-full" 
                  controls
                  preload="none"
                >
                  <source :src="`${FRONTEND_DOMAIN}/audio/models/${model.model}.wav`" type="audio/wav">
                </audio>
              </div>
            </div>
          </div>

          <!-- 生成的音频预览 -->
          <div v-if="options.audioUrl" class="mt-8 p-6 border theme-border rounded-xl">
            <h3 class="text-lg font-semibold text-gray-900 dark:text-gray-100 mb-4">
              生成的音频
            </h3>
            <div class="flex items-center gap-6">
              <kinesis-container 
                class="w-32 h-32"
                :playAudio="options.isAudioPlaying"
                :audio="options.audioUrl"
              >
                <kinesis-audio :audioIndex="17">
                  <kinesis-element :strength="10" type="depth">
                    <img 
                      class="w-full h-full rounded-xl object-cover"
                      :src="BACKEND_DOMAIN + userStore.user.avatar" 
                      alt="用户头像"
                    >
                  </kinesis-element>
                </kinesis-audio>
              </kinesis-container>

              <div class="flex-1">
                <audio 
                  class="w-full" 
                  controls
                  :src="options.audioUrl"
                />
              </div>
            </div>
          </div>
        </div>

        <!-- 文件导出 -->
        <div v-show="options.activeTab === 'export'" class="space-y-6 animate__animated animate__fadeIn animate__delay-3s">
          <h2 class="text-2xl font-bold text-gray-900 dark:text-gray-100">
            选择导出格式
          </h2>
          
          <div class="grid grid-cols-1 md:grid-cols-3 gap-6">
            <div
              v-for="format in [
                { key: 'pdf', name: 'PDF', icon: 'PDF', desc: '导出为PDF文档' },
                { key: 'word', name: 'Word', icon: 'DOCX', desc: '导出为Word文档' },
                { key: 'markdown', name: 'Markdown', icon: 'MD', desc: '导出为富文本' }
              ]"
              :key="format.key"
              class="group relative p-6 border border-dashed theme-border rounded-xl hover:border-blue-500 transition-all cursor-pointer"
              @click="exportScript(project, project.project_name, options.currentChapter?.Title, options.currentChapter?.current_version?.content, format.key)"
            >
              <div class="text-center space-y-4">
                <span class="block text-4xl font-bold text-gray-400 group-hover:text-blue-500 transition-colors">
                  {{ format.icon }}
                </span>
                <div>
                  <h3 class="text-lg font-semibold text-gray-900 dark:text-gray-100">
                    {{ format.name }}
                  </h3>
                  <p class="mt-1 text-sm text-gray-500">
                    {{ format.desc }}
                  </p>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- 视频导出 -->
        <div v-show="options.activeTab === 'video'" class="space-y-6 animate__animated animate__fadeIn animate__delay-3s">
          <h2 class="text-2xl font-bold text-gray-900 dark:text-gray-100">
            视频导出
          </h2>
          
          <!-- 无选择章节提示 -->
          <div v-if="!options.currentChapterId" class="p-8 text-center bg-gray-50 dark:bg-zinc-900/50 rounded-xl border border-gray-200 dark:border-zinc-700/80">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-12 w-12 mx-auto text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" d="M15.75 10.5l4.72-4.72a.75.75 0 011.28.53v11.38a.75.75 0 01-1.28.53l-4.72-4.72M4.5 18.75h9a2.25 2.25 0 002.25-2.25v-9a2.25 2.25 0 00-2.25-2.25h-9A2.25 2.25 0 002.25 7.5v9a2.25 2.25 0 002.25 2.25z" />
            </svg>
            <p class="mt-4 text-gray-600 dark:text-gray-400">请先选择要生成视频的章节</p>
            <a-button 
              class="mt-4" 
              type="primary" 
              @click="options.activeTab = 'chapters'"
            >
              前往选择章节
            </a-button>
          </div>
          
          <!-- 视频生成步骤 -->
          <div v-else class="space-y-8">
            <!-- 步骤指示器 -->
            <div class="flex justify-between items-center pt-8 pb-12">
              <div class="w-full flex items-center">
                <div class="relative flex items-center justify-center">
                  <div class="h-12 w-12 rounded-full border-2 flex items-center justify-center"
                       :class="options.videoPath ? 'bg-blue-500 text-white border-blue-500' : 'bg-blue-500 text-white border-blue-500'">
                    1
                  </div>
                  <div class="absolute -bottom-6 w-max text-xs font-medium"
                       :class="options.videoPath ? 'text-blue-500' : 'text-blue-500'">
                    生成场景
                  </div>
                </div>
                <div class="flex-auto border-t-2"
                     :class="options.videoPath ? 'border-blue-500' : 'border-blue-500'"></div>
                <div class="relative flex items-center justify-center">
                  <div class="h-12 w-12 rounded-full border-2 flex items-center justify-center"
                       :class="options.videoPath || (options.scenes && options.scenes.length > 0) ? 'bg-blue-500 text-white border-blue-500' : 'bg-white text-gray-400 border-gray-300 dark:bg-zinc-900 dark:border-gray-700'">
                    2
                  </div>
                  <div class="absolute -bottom-6 w-max text-xs font-medium"
                       :class="options.videoPath || (options.scenes && options.scenes.length > 0) ? 'text-blue-500' : 'text-gray-400'">
                    生成图片
                  </div>
                </div>
                <div class="flex-auto border-t-2"
                     :class="options.videoPath ? 'border-blue-500' : options.scenes && options.scenes.length > 0 ? 'border-blue-500' : 'border-gray-300 dark:border-gray-700'">
                </div>
                <div class="relative flex items-center justify-center">
                  <div class="h-12 w-12 rounded-full border-2 flex items-center justify-center"
                       :class="options.videoPath ? 'bg-blue-500 text-white border-blue-500' : 'bg-white text-gray-400 border-gray-300 dark:bg-zinc-900 dark:border-gray-700'">
                    3
                  </div>
                  <div class="absolute -bottom-6 w-max text-xs font-medium"
                       :class="options.videoPath ? 'text-blue-500' : 'text-gray-400'">
                    生成视频
                  </div>
                </div>
              </div>
            </div>
            
            <!-- 所选章节信息 -->
            <div class="p-4 bg-gray-50 dark:bg-zinc-900/50 rounded-xl border border-gray-200 dark:border-zinc-700/80">
              <div class="flex items-center gap-4">
                <div class="w-12 h-12 flex items-center justify-center bg-blue-100 dark:bg-blue-900/30 rounded-lg">
                  <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6 text-blue-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
                  </svg>
                </div>
                <div>
                  <h3 class="text-lg font-semibold text-gray-900 dark:text-gray-100">
                    {{ options.currentChapter.Title }}
                  </h3>
                  <p class="text-sm text-gray-600 dark:text-gray-400">
                    版本: {{ options.currentChapter.current_version?.ID || '-' }}
                  </p>
                </div>
              </div>
            </div>
            
            <!-- 步骤1：生成场景 -->
            <div class="space-y-4">
              <div class="flex items-center justify-between">
                <h3 class="text-lg font-semibold text-gray-900 dark:text-gray-100">
                  第一步：生成场景
                </h3>
                <div v-if="options.videoPath" class="flex items-center gap-2 text-sm text-green-600 dark:text-green-400">
                  <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" d="M5 13l4 4L19 7" />
                  </svg>
                  <span>已完成</span>
                </div>
                <a-button 
                  v-else
                  type="primary"
                  :loading="options.sceneGenerating"
                  @click="generateScenes"
                >
                  生成场景
                </a-button>
              </div>
              
              <!-- 场景列表 -->
              <div v-if="options.scenes && options.scenes.length > 0" class="space-y-4">
                <div 
                  v-for="(scene, index) in options.scenes" 
                  :key="index"
                  class="p-4 bg-white dark:bg-zinc-800 rounded-lg border theme-border"
                >
                  <div class="flex items-start gap-4">
                    <div class="flex-shrink-0 text-sm font-medium bg-blue-100 dark:bg-blue-900/30 text-blue-800 dark:text-blue-400 rounded-full h-6 w-6 flex items-center justify-center">
                      {{ index + 1 }}
                    </div>
                    <div class="flex-1 space-y-2">
                      <div class="text-sm text-gray-700 dark:text-gray-300">{{ scene.text }}</div>
                      <div class="text-xs text-gray-500">
                        <span class="font-medium">时间: </span>
                        {{ scene.start_time }}s - {{ scene.end_time }}s
                      </div>
                    </div>
                    <div class="flex-shrink-0 w-16 h-16 rounded bg-gray-100 dark:bg-gray-700 overflow-hidden" v-if="scene.image_path">
                      <img 
                        :src="`${BACKEND_DOMAIN}uploads/${scene.image_path}`" 
                        class="w-full h-full object-cover"
                        alt="场景图片"
                      />
                    </div>
                  </div>
                </div>
              </div>
              
              <div v-else-if="!options.sceneGenerating" class="text-center py-8 bg-gray-50 dark:bg-zinc-900/50 rounded-lg border border-dashed theme-border">
                <p class="text-gray-500">请点击"生成场景"按钮开始生成</p>
              </div>
            </div>
            
            <!-- 步骤2：生成图片 -->
            <div class="space-y-4">
              <div class="flex items-center justify-between">
                <h3 class="text-lg font-semibold text-gray-900 dark:text-gray-100">
                  第二步：生成场景图片
                </h3>
                <div v-if="options.videoPath" class="flex items-center gap-2 text-sm text-green-600 dark:text-green-400">
                  <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" d="M5 13l4 4L19 7" />
                  </svg>
                  <span>已完成</span>
                </div>
                <a-button 
                  v-else
                  type="primary"
                  :disabled="!options.scenes || options.scenes.length === 0"
                  :loading="options.imagesGenerating"
                  @click="generateImages"
                >
                  生成图片
                </a-button>
              </div>
              
              <div v-if="options.scenes && options.scenes.length > 0 && options.scenes.some(s => s.image_path)" class="grid grid-cols-1 md:grid-cols-3 gap-4">
                <div 
                  v-for="(scene, index) in options.scenes.filter(s => s.image_path)" 
                  :key="`img-${index}`"
                  class="bg-white dark:bg-zinc-800 rounded-lg overflow-hidden border theme-border"
                >
                  <div class="aspect-w-16 aspect-h-9 bg-gray-100 dark:bg-gray-700">
                    <img 
                      :src="`${BACKEND_DOMAIN}uploads/${scene.image_path}`" 
                      class="w-full h-full object-cover"
                      alt="场景图片"
                    />
                  </div>
                  <div class="p-3">
                    <div class="text-xs text-gray-500 truncate">
                      {{ scene.text.substring(0, 50) }}{{ scene.text.length > 50 ? '...' : '' }}
                    </div>
                  </div>
                </div>
              </div>
              
              <div 
                v-else-if="options.scenes && options.scenes.length > 0" 
                class="text-center py-8 bg-gray-50 dark:bg-zinc-900/50 rounded-lg border border-dashed theme-border"
              >
                <p class="text-gray-500">请点击"生成图片"按钮开始生成场景图片</p>
              </div>
              
              <div 
                v-else 
                class="text-center py-8 bg-gray-50 dark:bg-zinc-900/50 rounded-lg border border-dashed theme-border"
              >
                <p class="text-gray-500">请先完成场景生成</p>
              </div>
            </div>
            
            <!-- 步骤3：生成视频 -->
            <div class="space-y-4">
              <div class="flex items-center justify-between">
                <h3 class="text-lg font-semibold text-gray-900 dark:text-gray-100">
                  第三步：生成视频
                </h3>
                <a-button 
                  v-if="!options.videoPath"
                  type="primary"
                  :disabled="!options.scenes || options.scenes.length === 0 || !options.scenes.every(s => s.image_path)"
                  :loading="options.videoGenerating"
                  @click="generateVideo"
                >
                  生成视频
                </a-button>
                <div v-else class="flex items-center gap-2 text-sm text-green-600 dark:text-green-400">
                  <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" d="M5 13l4 4L19 7" />
                  </svg>
                  <span>视频已生成</span>
                </div>
              </div>
              
              <div v-if="options.videoPath" class="rounded-lg overflow-hidden border theme-border">
                <div class="aspect-w-16 aspect-h-9 bg-black">
                  <video 
                    :src="`${BACKEND_DOMAIN}videos/${options.videoPath}`" 
                    class="w-full h-full" 
                    controls
                  ></video>
                </div>
                <div class="p-4 bg-white dark:bg-zinc-800">
                  <div class="flex items-center justify-between">
                    <div>
                      <h4 class="font-medium text-gray-900 dark:text-gray-100">
                        {{ options.currentChapter.Title }} - 视频
                      </h4>
                      <p class="text-sm text-gray-500 mt-1">
                        点击播放按钮即可观看
                      </p>
                    </div>
                    <div class="flex gap-2">
                      <a-button
                        danger
                        v-if="options.videoPath"
                        @click="generateVideo"
                        :loading="options.videoGenerating"
                      >
                        重新生成
                      </a-button>
                      <a 
                        :href="`${BACKEND_DOMAIN}videos/${options.videoPath}`" 
                        download
                        class="inline-flex items-center gap-2 px-4 py-2 bg-blue-500 hover:bg-blue-600 transition-colors rounded-lg text-white text-sm"
                      >
                        <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                          <path stroke-linecap="round" stroke-linejoin="round" d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4" />
                        </svg>
                        下载视频
                      </a>
                    </div>
                  </div>
                  
                  <!-- 社交媒体分享按钮 -->
                  <div class="mt-4 flex flex-wrap gap-2">
                    <a 
                      href="#" 
                      @click.prevent="shareToDouyin" 
                      class="inline-flex items-center gap-1 px-3 py-1.5 bg-black hover:bg-gray-800 transition-colors rounded-lg text-white text-xs"
                    >
                    <img :src="douyin" class="w-4 h-4">
                      <span>发布至抖音</span>
                    </a>
                    <a 
                      href="#" 
                      @click.prevent="shareToBilibili" 
                      class="inline-flex items-center gap-1 px-3 py-1.5 bg-pink-500 hover:bg-pink-600 transition-colors rounded-lg text-white text-xs"
                    >
                    <img :src="bilibili" class="w-4 h-4">
                      <span>发布至哔哩哔哩</span>
                    </a>
                    <a 
                      href="#" 
                      @click.prevent="shareToXiaohongshu" 
                      class="inline-flex items-center gap-1 px-3 py-1.5 bg-red-500 hover:bg-red-600 transition-colors rounded-lg text-white text-xs"
                    >
                    <img :src="xiaohongshu" class="w-4 h-4">
                      <span>发布至小红书</span>
                    </a>
                    <a 
                      href="#" 
                      @click.prevent="shareToWechatChannels" 
                      class="inline-flex items-center gap-1 px-3 py-1.5 bg-green-500 hover:bg-green-600 transition-colors rounded-lg text-white text-xs"
                    >
                    <img :src="wechatChannels" class="w-4 h-4">
                      <span>发布至微信视频号</span>
                    </a>
                    <a 
                      href="#" 
                      @click.prevent="shareToQQXiaoshijie" 
                      class="inline-flex items-center gap-1 px-3 py-1.5 bg-blue-400 hover:bg-blue-500 transition-colors rounded-lg text-white text-xs"
                    >
                    <img :src="qqXiaoshijie" class="w-4 h-4">
                      <span>发布至QQ小世界</span>
                    </a>
                  </div>
                </div>
              </div>
              
              <div 
                v-else-if="options.scenes && options.scenes.length > 0 && options.scenes.every(s => s.image_path)" 
                class="text-center py-8 bg-gray-50 dark:bg-zinc-900/50 rounded-lg border border-dashed theme-border"
              >
                <p class="text-gray-500">请点击"生成视频"按钮开始生成视频</p>
              </div>
              
              <div 
                v-else 
                class="text-center py-8 bg-gray-50 dark:bg-zinc-900/50 rounded-lg border border-dashed theme-border"
              >
                <p class="text-gray-500">请先完成场景图片生成</p>
              </div>
            </div>
          </div>
        </div>
      </template>
    </div>
  </div>
</template>

<style scoped>
/* 添加动画相关样式 */
.animate__animated {
  animation-duration: 0.5s;
}

.animate__delay-1s {
  animation-delay: 0.1s;
}

.animate__delay-2s {
  animation-delay: 0.2s;
}

.animate__delay-3s {
  animation-delay: 0.3s;
}

.animate__delay-4s {
  animation-delay: 0.4s;
}

/* 为卡片添加淡入动画 */
@keyframes fadeIn {
  from {
    opacity: 0;
    transform: translateY(10px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

/* 为章节卡片添加动画 */
.grid > div {
  animation: fadeIn 0.5s ease-out forwards;
  animation-delay: calc(var(--index) * 0.1s);
  opacity: 0;
}

/* 视频容器样式 */
.aspect-w-16 {
  position: relative;
  padding-bottom: 56.25%; /* 16:9 Aspect Ratio */
}

.aspect-h-9 {
  position: relative;
}

.aspect-w-16 > * {
  position: absolute;
  height: 100%;
  width: 100%;
  top: 0;
  right: 0;
  bottom: 0;
  left: 0;
}
</style>