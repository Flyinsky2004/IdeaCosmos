<script setup>
import router from "@/router/index.js";
import Loader from "@/components/loader.vue";
import {onMounted, reactive, ref} from "vue";
import {get, post, postJSON} from "@/util/request.js";
import {message} from "ant-design-vue";
import {washJSONStr} from "@/util/common.js";
import SpinLoaderLarge from "@/components/spinLoaderLarge.vue";

const project = JSON.parse(localStorage.getItem("project"));
const options = reactive({
  isChapterGenerating: false,
  generatedChapters: [],
  chapters: []
})

// 添加编辑状态控制
const editState = reactive({
  currentEditingChapter: null,
  editMode: false,
  deleteConfirmId: null,
  newChapter: {
    Title: "",
    Description: ""
  },
  isCreating: false
})

// 添加 loading 状态
const loading = ref(true)

onMounted(() => {
  fetchChapters()
})

const fetchChapters = () => {
  loading.value = true
  get('/api/project/getAllChapters', {
    project_id: project.ID
  }, (messageer, data) => {
    options.chapters = data
    loading.value = false
  }, (messageer, data) => {
    message.warning(messageer)
    loading.value = false
  }, (messageer, data) => {
    message.error(messageer)
    loading.value = false
  })
}

const generateChapters = () => {
  options.isChapterGenerating = true;
  post('/api/project/generateChapters', {
    project_id: project.ID
  }, (messageer, data) => {
    const raw = data.choices[0].message.content
    options.generatedChapters = JSON.parse(washJSONStr(raw))
    message.success("篇章生成成功！")
    options.isChapterGenerating = false
  }, (messageer, data) => {
    message.warning(messageer)
    options.isChapterGenerating = false
  }, (messageer, data) => {
    message.error(messageer)
    options.isChapterGenerating = false
  })
}

const acceptChapter = () => {
  for (let i = 0; i < options.generatedChapters.length; i++) options.generatedChapters[i].project_id = project.ID;
  postJSON('/api/project/createChapterMulti',
      options.generatedChapters
      , (messageer, data) => {
        fetchChapters()
        message.success(data)
        options.generatedChapters = []
      }, (messageer, data) => {
        message.warning(messageer)
      }, (messageer, data) => {
        message.error(messageer)
      })
}

const rejectChapter = () => {
  options.generatedChapters = []
  message.info("已清空生成章节目录。")
}

const goToWriting = (chapter) => {
  localStorage.setItem("chapter",JSON.stringify(chapter))
  router.push('/workspace/editProject/writing')
}

// 添加创建章节方法
const startCreate = () => {
  editState.isCreating = true
  editState.newChapter = {
    Title: "",
    Description: "",
    project_id: project.ID
  }
}

const submitCreate = () => {
  if (!editState.newChapter.Title) {
    message.warning("章节标题不能为空")
    return
  }
  
  postJSON('/api/project/chapter/create', editState.newChapter, 
    (msg, data) => {
      message.success(msg)
      editState.isCreating = false
      fetchChapters()
    },
    (msg) => {
      message.warning(msg)
    },
    (msg) => {
      message.error(msg)
    }
  )
}

const cancelCreate = () => {
  editState.isCreating = false
}

// 添加编辑相关方法
const startEdit = (chapter) => {
  editState.currentEditingChapter = { ...chapter }
  editState.editMode = true
}

const saveEdit = () => {
  if (!editState.currentEditingChapter.Title) {
    message.warning("章节标题不能为空")
    return
  }
  
  postJSON('/api/project/chapter/update', editState.currentEditingChapter,
    (msg, data) => {
      message.success(msg)
      fetchChapters()
      editState.editMode = false
      editState.currentEditingChapter = null
    },
    (msg) => {
      message.warning(msg)
    },
    (msg) => {
      message.error(msg)
    }
  )
}

const cancelEdit = () => {
  editState.editMode = false
  editState.currentEditingChapter = null
}

// 添加删除确认方法
const confirmDelete = (chapterId) => {
  editState.deleteConfirmId = chapterId
}

const cancelDelete = () => {
  editState.deleteConfirmId = null
}

const executeDelete = (chapterId) => {
  postJSON('/api/project/chapter/delete', 
    { chapter_id: chapterId },
    (msg) => {
      message.success(msg)
      fetchChapters()
      editState.deleteConfirmId = null
    },
    (msg) => {
      message.warning(msg)
    },
    (msg) => {
      message.error(msg)
    }
  )
}
</script>

<template>
  <div class="space-y-8 p-6">
    <!-- 顶部操作区 -->
    <div class="flex items-center justify-between animate__animated animate__fadeIn">
      <h1 class="text-2xl font-bold bg-gradient-to-r from-blue-500 to-cyan-500 bg-clip-text text-transparent">
        章节管理
      </h1>
      <div class="flex gap-3">
        <button 
          @click="startCreate"
          class="px-4 py-2 bg-green-500 text-white rounded-lg hover:bg-green-600 transition-all flex items-center gap-2"
        >
          <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" d="M12 4v16m8-8H4" />
          </svg>
          添加章节
        </button>
        <button 
          @click="generateChapters"
          :disabled="options.isChapterGenerating"
          class="px-4 py-2 bg-gradient-to-r from-blue-500 to-cyan-500 text-white rounded-lg hover:opacity-90 transition-all flex items-center gap-2"
        >
          <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" d="M9.813 15.904L9 18.75l-.813-2.846a4.5 4.5 0 00-3.09-3.09L2.25 12l2.846-.813a4.5 4.5 0 003.09-3.09L9 5.25l.813 2.846a4.5 4.5 0 003.09 3.09L15.75 12l-2.846.813a4.5 4.5 0 00-3.09 3.09z" />
          </svg>
          {{ options.isChapterGenerating ? '生成中...' : 'AI 生成章节' }}
        </button>
      </div>
    </div>

    <!-- 加载状态显示 -->
    <div v-if="loading" class="flex items-center justify-center py-12">
      <SpinLoaderLarge />
    </div>

    <!-- 章节内容区域 -->
    <template v-else>
      <!-- 创建章节表单 -->
      <div v-if="editState.isCreating" class="bg-white dark:bg-zinc-900 rounded-xl border theme-border p-6 animate__animated animate__fadeIn">
        <h2 class="text-lg font-semibold text-blue-500 mb-4">添加新章节</h2>
        
        <div class="space-y-4">
          <div>
            <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">章节标题</label>
            <input 
              v-model="editState.newChapter.Title" 
              type="text" 
              class="w-full p-2 border theme-border rounded-lg bg-white dark:bg-zinc-800 text-gray-900 dark:text-gray-100 focus:outline-none focus:ring-1 focus:ring-blue-500 transition-colors"
              placeholder="请输入章节标题"
            />
          </div>
          
          <div>
            <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">章节简介</label>
            <textarea 
              v-model="editState.newChapter.Description" 
              class="w-full p-2 border theme-border rounded-lg bg-white dark:bg-zinc-800 text-gray-900 dark:text-gray-100 focus:outline-none focus:ring-1 focus:ring-blue-500 transition-colors resize-none min-h-[100px]"
              placeholder="请输入章节简介"
            ></textarea>
          </div>
          
          <div class="flex justify-end space-x-3 pt-3">
            <button 
              @click="cancelCreate" 
              class="px-4 py-2 border theme-border rounded-lg hover:bg-gray-50 dark:hover:bg-zinc-800 transition-colors"
            >
              取消
            </button>
            <button 
              @click="submitCreate" 
              class="px-4 py-2 bg-blue-500 text-white rounded-lg hover:bg-blue-600 transition-colors"
            >
              创建章节
            </button>
          </div>
        </div>
      </div>
      
      <!-- 编辑章节表单 -->
      <div v-if="editState.editMode" class="bg-white dark:bg-zinc-900 rounded-xl border theme-border p-6 animate__animated animate__fadeIn">
        <h2 class="text-lg font-semibold text-blue-500 mb-4">编辑章节</h2>
        
        <div class="space-y-4">
          <div>
            <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">章节标题</label>
            <input 
              v-model="editState.currentEditingChapter.Title" 
              type="text" 
              class="w-full p-2 border theme-border rounded-lg bg-white dark:bg-zinc-800 text-gray-900 dark:text-gray-100 focus:outline-none focus:ring-1 focus:ring-blue-500 transition-colors"
              placeholder="请输入章节标题"
            />
          </div>
          
          <div>
            <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">章节简介</label>
            <textarea 
              v-model="editState.currentEditingChapter.Description" 
              class="w-full p-2 border theme-border rounded-lg bg-white dark:bg-zinc-800 text-gray-900 dark:text-gray-100 focus:outline-none focus:ring-1 focus:ring-blue-500 transition-colors resize-none min-h-[100px]"
              placeholder="请输入章节简介"
            ></textarea>
          </div>
          
          <div class="flex justify-end space-x-3 pt-3">
            <button 
              @click="cancelEdit" 
              class="px-4 py-2 border theme-border rounded-lg hover:bg-gray-50 dark:hover:bg-zinc-800 transition-colors"
            >
              取消
            </button>
            <button 
              @click="saveEdit" 
              class="px-4 py-2 bg-blue-500 text-white rounded-lg hover:bg-blue-600 transition-colors"
            >
              保存修改
            </button>
          </div>
        </div>
      </div>

      <div class="relative min-h-[200px] bg-white dark:bg-zinc-900 rounded-xl border theme-border p-6">
        <!-- 空状态提示 -->
        <div v-if="options.chapters.length === 0 && options.generatedChapters.length === 0" 
             class="flex flex-col items-center justify-center h-48 text-gray-400">
          <svg xmlns="http://www.w3.org/2000/svg" class="h-12 w-12 mb-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" d="M12 6.042A8.967 8.967 0 006 3.75c-1.052 0-2.062.18-3 .512v14.25A8.987 8.987 0 016 18c2.305 0 4.408.867 6 2.292m0-14.25a8.966 8.966 0 016-2.292c1.052 0 2.062.18 3 .512v14.25A8.987 8.987 0 0018 18a8.967 8.967 0 00-6 2.292m0-14.25v14.25" />
          </svg>
          <span>暂无章节内容</span>
        </div>

        <!-- 现有章节展示 -->
        <div v-if="options.chapters.length > 0" class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 2xl:grid-cols-4 gap-8">
          <div v-for="(chapter, index) in options.chapters" 
               :key="chapter.ID"
               class="group relative animate__animated animate__fadeIn"
               :style="{ animationDelay: `${index * 0.1}s` }">
            
            <!-- 删除确认浮层 -->
            <div v-if="editState.deleteConfirmId === chapter.ID" 
                 class="absolute inset-0 bg-white dark:bg-zinc-800 bg-opacity-95 dark:bg-opacity-95 z-10 rounded-xl border border-red-500 dark:border-red-600 p-4 flex flex-col items-center justify-center animate__animated animate__fadeIn">
              <p class="text-center text-red-600 dark:text-red-500 mb-4">确定要删除此章节吗？</p>
              <div class="flex gap-4">
                <button @click="cancelDelete" class="px-4 py-2 border theme-border rounded-lg hover:bg-gray-50 dark:hover:bg-zinc-700 transition-colors">
                  取消
                </button>
                <button @click="executeDelete(chapter.ID)" class="px-4 py-2 bg-red-500 text-white rounded-lg hover:bg-red-600 transition-colors">
                  删除
                </button>
              </div>
            </div>
            
            <!-- 章节卡片 -->
            <div class="w-full p-4 rounded-xl bg-white dark:bg-zinc-800 border theme-border hover:border-blue-500 dark:hover:border-blue-400 transition-all duration-300 group-hover:shadow-lg">
              <div class="flex items-center justify-between mb-3">
                <span class="text-sm text-gray-500 dark:text-gray-400">第 {{ index + 1 }} 章</span>
                <div class="flex items-center gap-2 opacity-0 group-hover:opacity-100 transition-opacity">
                  <button @click="goToWriting(chapter)" 
                          class="p-1 text-green-500 hover:bg-green-50 dark:hover:bg-green-900/20 rounded-lg transition-colors">
                    <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                      <path stroke-linecap="round" stroke-linejoin="round" d="M15.232 5.232l3.536 3.536m-2.036-5.036a2.5 2.5 0 113.536 3.536L6.5 21.036H3v-3.572L16.732 3.732z" />
                    </svg>
                  </button>
                  <button @click="startEdit(chapter)"
                          class="p-1 text-blue-500 hover:bg-blue-50 dark:hover:bg-blue-900/20 rounded-lg transition-colors">
                    <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                      <path stroke-linecap="round" stroke-linejoin="round" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
                    </svg>
                  </button>
                  <button @click="confirmDelete(chapter.ID)"
                          class="p-1 text-red-500 hover:bg-red-50 dark:hover:bg-red-900/20 rounded-lg transition-colors">
                    <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                      <path stroke-linecap="round" stroke-linejoin="round" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                    </svg>
                  </button>
                </div>
              </div>
              <h3 class="text-lg font-medium mb-2 text-gray-900 dark:text-gray-100">{{ chapter.Title }}</h3>
              <p class="text-sm text-gray-600 dark:text-gray-400 line-clamp-3">{{ chapter.Description }}</p>
              
              <!-- 章节序号指示器 -->
              <div class="absolute -left-3 -top-3 w-6 h-6 rounded-full bg-gradient-to-r from-blue-500 to-cyan-500 flex items-center justify-center text-white text-sm font-medium">
                {{ index + 1 }}
              </div>
            </div>
          </div>

          <!-- 添加新章节按钮 -->
          <div 
               @click="startCreate"
               class="w-full p-4 rounded-xl border-2 border-dashed theme-border hover:border-blue-500 dark:hover:border-blue-400 transition-all duration-300 cursor-pointer group animate__animated animate__fadeIn"
               :style="{ animationDelay: `${options.chapters.length * 0.1}s` }">
            <div class="flex flex-col items-center justify-center h-40 text-gray-400 group-hover:text-blue-500">
              <svg xmlns="http://www.w3.org/2000/svg" class="h-8 w-8 mb-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" d="M12 4v16m8-8H4" />
              </svg>
              <span>添加新章节</span>
            </div>
          </div>
        </div>

        <!-- AI 生成的章节建议 -->
        <div v-if="options.generatedChapters.length > 0" 
             class="mt-8 p-6 bg-green-50 dark:bg-green-900/20 rounded-xl border border-green-200 dark:border-green-800 animate__animated animate__fadeIn">
          <div class="flex items-center justify-between mb-4">
            <h3 class="text-lg font-medium text-green-800 dark:text-green-400">AI 章节建议</h3>
            <div class="flex items-center gap-4">
              <button @click="acceptChapter"
                      class="px-4 py-2 bg-green-500 text-white rounded-lg hover:bg-green-600 transition-colors">
                采用建议
              </button>
              <button @click="rejectChapter"
                      class="px-4 py-2 border border-red-200 dark:border-red-800 text-red-600 rounded-lg hover:bg-red-50 dark:hover:bg-red-900/20 transition-colors">
                放弃建议
              </button>
            </div>
          </div>

          <!-- 修改 AI 建议章节的展示方式 -->
          <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 2xl:grid-cols-4 gap-8">
            <div v-for="(chapter, index) in options.generatedChapters" 
                 :key="index"
                 class="relative">
              <div class="w-full p-4 rounded-xl bg-white dark:bg-zinc-800 border border-green-200 dark:border-green-800">
                <span class="text-sm text-green-600 dark:text-green-400 mb-2 block">建议章节 {{ index + 1 }}</span>
                <h3 class="text-lg font-medium mb-2 text-gray-900 dark:text-gray-100">{{ chapter.Title }}</h3>
                <p class="text-sm text-gray-600 dark:text-gray-400 line-clamp-3">{{ chapter.Description }}</p>
                
                <!-- 章节序号指示器 -->
                <div class="absolute -left-3 -top-3 w-6 h-6 rounded-full bg-gradient-to-r from-green-500 to-emerald-500 flex items-center justify-center text-white text-sm font-medium">
                  {{ index + 1 }}
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </template>
  </div>
</template>

<style scoped>
.animate__animated {
  animation-duration: 0.5s;
}

/* 为卡片添加淡入动画 */
.grid > div {
  animation: fadeIn 0.5s ease-out forwards;
  opacity: 0;
}

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

/* 为每个卡片添加延迟 */
.grid > div:nth-child(1) { animation-delay: 0.1s; }
.grid > div:nth-child(2) { animation-delay: 0.2s; }
.grid > div:nth-child(3) { animation-delay: 0.3s; }
.grid > div:nth-child(4) { animation-delay: 0.4s; }
.grid > div:nth-child(5) { animation-delay: 0.5s; }
.grid > div:nth-child(6) { animation-delay: 0.6s; }
.grid > div:nth-child(7) { animation-delay: 0.7s; }
.grid > div:nth-child(8) { animation-delay: 0.8s; }
</style>