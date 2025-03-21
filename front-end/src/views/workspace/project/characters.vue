<script setup>
import { onMounted, reactive, computed, ref } from "vue";
import { get, post, postJSON } from "@/util/request.js";
import { message } from "ant-design-vue";
import router from "@/router/index.js";
import Loader from "@/components/loader.vue";
import { parseDateTime, washJSONStr } from "@/util/common.js";
import { imagePrefix } from "@/util/VARRIBLES.js";
import RelationShipGraph from "@/components/RelationShipGraph.vue";
import SpinLoader from "@/components/spinLoader.vue";
import SpinLoaderLarge from "@/components/spinLoaderLarge.vue";

const project = JSON.parse(localStorage.getItem("project"));
const loading = ref(true);

onMounted(() => {
  fetchCharacters();
  fetchCharacterRs();
});

const options = reactive({
  characters: [],
  isCharactersGenerating: false,
  isMainCharactersGenerating: false,
  generateResults: [],
  currentMoveId: -1,
  isCharacterAvatarGenerating: false,
  characterGeneratingId: -1,
  isAddRSShow: false,
  editRSMode: false,
  currentRS: {},
  isCharacterRSGenerating: false,
  characterRSs: [],
});
const hoverState = reactive({
  characterId: null,
  showActions: false,
});
const sortedCharacters = computed(() => {
  return [...options.characters].sort((a, b) => a.name.localeCompare(b.name));
});
const fetchCharacterRs = () => {
  get(
    "/api/project/characterRS/getAll",
    {
      project_id: project.ID,
    },
    (messager, data) => {
      options.characterRSs = data;
    },
    (messager, data) => {
      message.warning(messager);
    },
    (messager, data) => {
      message.error(messager);
    }
  );
};
const fetchCharacters = () => {
  loading.value = true;
  post(
    "/api/project/getCharacters",
    {
      project_id: project.ID,
    },
    (messager, data) => {
      options.characters = data;
      if (options.characters.length > 1) {
        editRSForm.first_character_id = options.characters[0].ID;
        editRSForm.second_character_id = options.characters[1].ID;
      }
      loading.value = false;
    },
    (messager, data) => {
      message.warning(messager);
      loading.value = false;
    },
    (messager, data) => {
      message.error(messager);
      loading.value = false;
    }
  );
};
const generateCharacter = () => {
  options.isCharactersGenerating = true;
  post(
    "/api/project/generateCharacter",
    {
      project_id: project.ID,
    },
    (messager, data) => {
      let  raw;
      const content = data.choices[0].message.content;
      if (content.includes("<think>") && content.includes("</think>")) {
        raw = content.replace(/<think>.*?<\/think>/gs, "");
      } else {
        raw = content;
      }
      options.generateResults = JSON.parse(washJSONStr(raw));
      for (let i = 0; i < options.generateResults.length; i++) {
        options.generateResults[i].project_id = project.ID;
      }
      options.isCharactersGenerating = false;
    }
  );
};
const applyResult = () => {
  options.generateResults[0].project_id = project.ID;
  postJSON(
    "/api/project/createCharacterArray",
    options.generateResults,
    (messageer, data) => {
      message.success(messageer);
      fetchCharacters();
      dropResult();
    },
    (messageer, data) => {
      message.warning(messageer);
    },
    (messageer, data) => {
      message.error(messageer);
    }
  );
};
const dropResult = () => {
  options.generateResults = [];
};
const moveIn = (id) => {
  options.currentMoveId = id;
};
const moveOut = () => {
  options.currentMoveId = -1;
};
const generateAvatar = async (character) => {
  options.characterGeneratingId = character.ID;
  options.isCharacterAvatarGenerating = true;
  await new Promise((resolve) => {
    post(
      "/api/project/generateCharacterAvatar",
      {
        character_id: character.ID,
      },
      (messageer, data) => {
        message.success(messageer);
        character.avatar = data;
        for (let i = 0; i < options.characters.length; i++) {
          if (options.characters[i].ID === character.ID) {
            options.characters[i].avatar = data;
            break;
          }
        }
        options.isCharacterAvatarGenerating = false;
        resolve();
      },
      (messageer, data) => {
        message.warning(messageer);
        options.isCharacterAvatarGenerating = false;
      },
      (messageer, data) => {
        message.error(messageer);
        options.isCharacterAvatarGenerating = false;
      }
    );
  }).then(() => {});
};
const editRSForm = reactive({
  first_character_id: 1,
  second_character_id: 2,
  name: "",
  content: "",
});
const generateCharacterRS = () => {
  if (editRSForm.first_character_id === editRSForm.second_character_id) {
    message.info("角色选择不合法");
    return;
  }
  options.isCharacterRSGenerating = true;
  post(
    "/api/project/generateCharacterRS",
    {
      firstCharacterId: editRSForm.first_character_id,
      secondCharacterId: editRSForm.second_character_id,
    },
    (messager, data) => {
      let  raw;
      const content = data;
      if (content.includes("<think>") && content.includes("</think>")) {
        raw = content.replace(/<think>.*?<\/think>/gs, "");
      } else {
        raw = content;
      }
      
      const washed = JSON.parse(washJSONStr(raw));
      console.log(washed);
      editRSForm.name = washed.name;
      editRSForm.content = washed.content;
      options.isCharacterRSGenerating = false;
      message.success("生成成功");
    },
    (messageer, data) => {
      message.warning(messageer);
    },
    (messageer, data) => {
      message.error(messageer);
    }
  );
};
const submitCharacterRS = () => {
  if (options.editRSMode) {
  } else {
    postJSON(
      "/api/project/characterRS/create",
      editRSForm,
      (messager, data) => {
        message.success(messager);
        options.isAddRSShow = false;
        fetchCharacterRs();
      },
      (messageer, data) => {
        message.warning(messageer);
      },
      (messageer, data) => {
        message.error(messageer);
      }
    );
  }
};
const generateMainCharacters = () => {
  options.isMainCharactersGenerating = true;
  post(
    "/api/project/generateCharacterFromDescription",
    {
      project_id: project.ID,
    },
    (messager, data) => {
      let  raw;
      const content = data.choices[0].message.content;
      if (content.includes("<think>") && content.includes("</think>")) {
        raw = content.replace(/<think>.*?<\/think>/gs, "");
      } else {
        raw = content;
      }
      options.generateResults = JSON.parse(washJSONStr(raw));
      for (let i = 0; i < options.generateResults.length; i++) {
        options.generateResults[i].project_id = project.ID;
      }
      options.isMainCharactersGenerating = false;
    }
  );
};

// 添加手动创建角色的状态
const manualCharacterModal = reactive({
  visible: false,
  loading: false
});

// 添加新角色表单
const newCharacterForm = reactive({
  name: '',
  description: '',
  project_id: project.ID,
  avatar: 'default-avatar.png' // 默认头像
});

// 重置表单
const resetCharacterForm = () => {
  newCharacterForm.name = '';
  newCharacterForm.description = '';
};

// 手动创建角色方法
const createCharacter = () => {
  // 表单验证
  if (!newCharacterForm.name.trim()) {
    message.warning('请输入角色名称');
    return;
  }
  
  if (!newCharacterForm.description.trim()) {
    message.warning('请输入角色描述');
    return;
  }
  
  manualCharacterModal.loading = true;
  
  postJSON(
    '/api/project/createCharacter',
    newCharacterForm,
    (msg, data) => {
      message.success('角色创建成功');
      manualCharacterModal.visible = false;
      manualCharacterModal.loading = false;
      resetCharacterForm();
      fetchCharacters(); // 刷新角色列表
    },
    (msg) => {
      message.warning(msg);
      manualCharacterModal.loading = false;
    },
    (msg) => {
      message.error(msg);
      manualCharacterModal.loading = false;
    }
  );
};

// 添加角色编辑功能
const editCharacterModal = reactive({
  visible: false,
  loading: false,
  character: {
    ID: null,
    name: '',
    description: '',
    project_id: project.ID
  }
});

// 打开编辑角色模态框
const openEditCharacterModal = (character) => {
  editCharacterModal.character.ID = character.ID;
  editCharacterModal.character.name = character.name;
  editCharacterModal.character.description = character.description;
  editCharacterModal.character.project_id = character.project_id;
  editCharacterModal.visible = true;
};

// 更新角色
const updateCharacter = () => {
  // 表单验证
  if (!editCharacterModal.character.name.trim()) {
    message.warning('请输入角色名称');
    return;
  }
  
  if (!editCharacterModal.character.description.trim()) {
    message.warning('请输入角色描述');
    return;
  }
  
  editCharacterModal.loading = true;
  
  postJSON(
    '/api/project/updateCharacter',
    editCharacterModal.character,
    (msg, data) => {
      message.success('角色更新成功');
      editCharacterModal.visible = false;
      editCharacterModal.loading = false;
      fetchCharacters(); // 刷新角色列表
    },
    (msg) => {
      message.warning(msg);
      editCharacterModal.loading = false;
    },
    (msg) => {
      message.error(msg);
      editCharacterModal.loading = false;
    }
  );
};

// 删除角色相关功能
const deleteModal = reactive({
  visible: false,
  loading: false,
  characterId: null,
  characterName: ''
});

// 打开删除确认框
const openDeleteModal = (character) => {
  deleteModal.characterId = character.ID;
  deleteModal.characterName = character.name;
  deleteModal.visible = true;
};

// 执行删除角色
const deleteCharacter = () => {
  deleteModal.loading = true;
  
  postJSON(
    '/api/project/deleteCharacter',
    { character_id: deleteModal.characterId },
    (msg) => {
      message.success('角色删除成功');
      deleteModal.visible = false;
      deleteModal.loading = false;
      fetchCharacters(); // 刷新角色列表
      fetchCharacterRs(); // 刷新角色关系
    },
    (msg) => {
      message.warning(msg);
      deleteModal.loading = false;
    },
    (msg) => {
      message.error(msg);
      deleteModal.loading = false;
    }
  );
};
</script>

<template>
  <div class="max-w-7xl mx-auto p-6 space-y-8">
    <!-- 页面标题 -->
    <div class="flex justify-between items-center">
      <h1
        class="text-3xl font-bold bg-gradient-to-r from-blue-500 to-cyan-500 bg-clip-text text-transparent"
      >
        角色管理
      </h1>
      <div class="flex items-center gap-4">
        <!-- 添加手动创建角色按钮 -->
        <button
          @click="manualCharacterModal.visible = true"
          class="px-4 py-2 bg-blue-500 text-white rounded-lg hover:bg-blue-600 transition-all flex items-center gap-2"
        >
          <svg
            xmlns="http://www.w3.org/2000/svg"
            class="h-5 w-5"
            fill="none"
            viewBox="0 0 24 24"
            stroke="currentColor"
          >
            <path
              stroke-linecap="round"
              stroke-linejoin="round"
              d="M12 4.5v15m7.5-7.5h-15"
            />
          </svg>
          创建角色
        </button>
        
        <button
          @click="generateMainCharacters"
          :disabled="options.isMainCharactersGenerating"
          class="px-4 py-2 bg-white dark:bg-zinc-800 border theme-border text-blue-500 rounded-lg hover:border-blue-500 transition-all flex items-center gap-2 disabled:opacity-50 disabled:cursor-not-allowed"
        >
          <svg
            xmlns="http://www.w3.org/2000/svg"
            class="h-5 w-5"
            fill="none"
            viewBox="0 0 24 24"
            stroke="currentColor"
          >
            <path
              stroke-linecap="round"
              stroke-linejoin="round"
              d="M15 19.128a9.38 9.38 0 002.625.372 9.337 9.337 0 004.121-.952 4.125 4.125 0 00-7.533-2.493M15 19.128v-.003c0-1.113-.285-2.16-.786-3.07M15 19.128v.106A12.318 12.318 0 018.624 21c-2.331 0-4.512-.645-6.374-1.766l-.001-.109a6.375 6.375 0 0111.964-3.07M12 6.375a3.375 3.375 0 11-6.75 0 3.375 3.375 0 016.75 0zm8.25 2.25a2.625 2.625 0 11-5.25 0 2.625 2.625 0 015.25 0z"
            />
          </svg>
          {{ options.isMainCharactersGenerating ? "补全中..." : "补全主角团" }}
          <SpinLoader
            v-if="options.isMainCharactersGenerating"
            class="ml-2 h-5 w-5"
          />
        </button>
        <button
          @click="generateCharacter"
          :disabled="options.isCharactersGenerating"
          class="px-4 py-2 bg-gradient-to-r from-blue-500 to-cyan-500 text-white rounded-lg hover:opacity-90 transition-all flex items-center gap-2"
        >
          <svg
            xmlns="http://www.w3.org/2000/svg"
            class="h-5 w-5"
            fill="none"
            viewBox="0 0 24 24"
            stroke="currentColor"
          >
            <path
              stroke-linecap="round"
              stroke-linejoin="round"
              d="M12 4.5v15m7.5-7.5h-15"
            />
          </svg>
          AI 生成角色
          <SpinLoader
            v-if="options.isCharactersGenerating"
            class="ml-2 h-5 w-5 mx-auto my-auto"
          />
        </button>
      </div>
    </div>

    <!-- 添加加载状态显示 -->
    <div v-if="loading" class="flex items-center justify-center py-12">
      <SpinLoaderLarge />
    </div>

    <template v-else>
      <!-- 角色列表 -->
      <div
        class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-6"
      >
        <!-- 现有角色卡片 -->
        <div
          v-for="character in sortedCharacters"
          :key="character.ID"
          class="group bg-white dark:bg-zinc-900 rounded-xl border theme-border overflow-hidden hover:shadow-lg dark:hover:shadow-zinc-800/20 transition-all duration-300"
          @mouseenter="hoverState.characterId = character.ID"
          @mouseleave="hoverState.characterId = null"
        >
          <!-- 角色头像区域 -->
          <div class="relative aspect-square">
            <div
              v-if="
                options.isCharacterAvatarGenerating &&
                options.characterGeneratingId === character.ID
              "
              class="absolute inset-0 bg-gray-900/60 flex items-center justify-center backdrop-blur-sm"
            >
              <div class="text-center">
                <loader class="mx-auto mb-2" />
                <span class="text-sm text-white">根据人物设定绘画中...</span>
              </div>
            </div>
            <div
              v-else-if="!character.avatar"
              class="h-full flex items-center justify-center bg-gray-100 dark:bg-zinc-800"
            >
              <div class="text-center text-gray-400">
                <svg
                  xmlns="http://www.w3.org/2000/svg"
                  class="h-12 w-12 mx-auto mb-2"
                  fill="none"
                  viewBox="0 0 24 24"
                  stroke="currentColor"
                >
                  <path
                    stroke-linecap="round"
                    stroke-linejoin="round"
                    d="M18.685 19.097A9.723 9.723 0 0021.75 12c0-5.385-4.365-9.75-9.75-9.75S2.25 6.615 2.25 12a9.723 9.723 0 003.065 7.097A9.716 9.716 0 0012 21.75a9.716 9.716 0 006.685-2.653zm-12.54-1.285A7.486 7.486 0 0112 15a7.486 7.486 0 015.855 2.812A8.224 8.224 0 0112 20.25a8.224 8.224 0 01-5.855-2.438zM15.75 9a3.75 3.75 0 11-7.5 0 3.75 3.75 0 017.5 0z"
                  />
                </svg>
                暂无角色画像
              </div>
            </div>
            <img
              v-else
              :src="imagePrefix + character.avatar"
              :alt="character.name"
              class="w-full h-full object-cover"
            />

            <!-- 悬浮操作按钮 -->
            <div
              v-show="hoverState.characterId === character.ID"
              class="absolute inset-0 bg-black/60 flex items-center justify-center gap-2 animate__animated animate__fadeIn animate__faster"
            >
              <button
                @click="generateAvatar(character)"
                class="px-3 py-1.5 bg-blue-500 text-white rounded-lg hover:bg-blue-600 transition-colors flex items-center gap-1"
              >
                <svg
                  xmlns="http://www.w3.org/2000/svg"
                  class="h-4 w-4"
                  fill="none"
                  viewBox="0 0 24 24"
                  stroke="currentColor"
                >
                  <path
                    stroke-linecap="round"
                    stroke-linejoin="round"
                    d="M9.813 15.904L9 18.75l-.813-2.846a4.5 4.5 0 00-3.09-3.09L2.25 12l2.846-.813a4.5 4.5 0 003.09-3.09L9 5.25l.813 2.846a4.5 4.5 0 003.09 3.09L15.75 12l-2.846.813a4.5 4.5 0 00-3.09 3.09z"
                  />
                </svg>
                AI 重绘
              </button>
              <button
                @click="openEditCharacterModal(character)"
                class="px-3 py-1.5 bg-amber-500 text-white rounded-lg hover:bg-amber-600 transition-colors flex items-center gap-1"
              >
                <svg
                  xmlns="http://www.w3.org/2000/svg"
                  class="h-4 w-4"
                  fill="none"
                  viewBox="0 0 24 24"
                  stroke="currentColor"
                >
                  <path
                    stroke-linecap="round"
                    stroke-linejoin="round"
                    d="M16.862 4.487l1.687-1.688a1.875 1.875 0 112.652 2.652L10.582 16.07a4.5 4.5 0 01-1.897 1.13L6 18l.8-2.685a4.5 4.5 0 011.13-1.897l8.932-8.931zm0 0L19.5
                    7.125"
                  />
                </svg>
                编辑
              </button>
              <button
                @click="openDeleteModal(character)"
                class="px-3 py-1.5 bg-red-500 text-white rounded-lg hover:bg-red-600 transition-colors flex items-center gap-1"
              >
                <svg
                  xmlns="http://www.w3.org/2000/svg"
                  class="h-4 w-4"
                  fill="none"
                  viewBox="0 0 24 24"
                  stroke="currentColor"
                >
                  <path
                    stroke-linecap="round"
                    stroke-linejoin="round"
                    d="M14.74 9l-.346 9m-4.788 0L9.26 9m9.968-3.21c.342.052.682.107 1.022.166m-1.022-.165L18.16 19.673a2.25 2.25 0 01-2.244 2.077H8.084a2.25 2.25 0 01-2.244-2.077L4.772 5.79m14.456 0a48.108 48.108 0 00-3.478-.397m-12 .562c.34-.059.68-.114 1.022-.165m0 0a48.11 48.11 0 013.478-.397m7.5 0v-.916c0-1.18-.91-2.164-2.09-2.201a51.964 51.964 0 00-3.32 0c-1.18.037-2.09 1.022-2.09 2.201v.916m7.5 0a48.667 48.667 0 00-7.5 0"
                  />
                </svg>
                删除
              </button>
            </div>
          </div>

          <!-- 角色信息 -->
          <div class="p-4">
            <h3
              class="text-lg font-semibold text-gray-900 dark:text-gray-100 mb-2"
            >
              {{ character.name }}
            </h3>
            <p class="text-sm text-gray-600 dark:text-gray-400 line-clamp-3">
              {{ character.description }}
            </p>
          </div>
        </div>

        <!-- AI 生成的角色建议 -->
        <div v-if="options.generateResults.length > 0" class="col-span-full">
          <div
            class="bg-green-50 dark:bg-green-900/20 rounded-xl border border-green-200 dark:border-green-800 p-6"
          >
            <div class="flex justify-between items-center mb-4">
              <h2
                class="text-lg font-medium text-green-800 dark:text-green-400"
              >
                AI 生成的角色建议
              </h2>
              <div class="flex items-center gap-2">
                <button
                  @click="applyResult"
                  class="px-4 py-2 bg-green-500 text-white rounded-lg hover:bg-green-600 transition-colors flex items-center gap-2"
                >
                  <svg
                    xmlns="http://www.w3.org/2000/svg"
                    class="h-5 w-5"
                    fill="none"
                    viewBox="0 0 24 24"
                    stroke="currentColor"
                  >
                    <path
                      stroke-linecap="round"
                      stroke-linejoin="round"
                      d="M5 13l4 4L19 7"
                    />
                  </svg>
                  应用建议
                </button>
                <button
                  @click="dropResult"
                  class="px-4 py-2 border border-red-200 dark:border-red-800 text-red-600 rounded-lg hover:bg-red-50 dark:hover:bg-red-900/20 transition-colors flex items-center gap-2"
                >
                  <svg
                    xmlns="http://www.w3.org/2000/svg"
                    class="h-5 w-5"
                    fill="none"
                    viewBox="0 0 24 24"
                    stroke="currentColor"
                  >
                    <path
                      stroke-linecap="round"
                      stroke-linejoin="round"
                      d="M6 18L18 6M6 6l12 12"
                    />
                  </svg>
                  放弃建议
                </button>
              </div>
            </div>

            <div
              class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4"
            >
              <div
                v-for="character in options.generateResults"
                :key="character.name"
                class="bg-white dark:bg-zinc-800 rounded-lg p-4"
              >
                <h3
                  class="text-lg font-medium text-gray-900 dark:text-gray-100 mb-2"
                >
                  {{ character.name }}
                </h3>
                <p class="text-sm text-gray-600 dark:text-gray-400">
                  {{ character.description }}
                </p>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- 角色关系图谱 -->
      <div class="space-y-4">
        <div class="flex justify-between items-center">
          <h2 class="text-2xl font-bold text-blue-500">角色关系</h2>
          <div class="flex items-center gap-2">
            <button
              @click="options.isAddRSShow = true"
              class="px-4 py-2 bg-blue-500 text-white rounded-lg hover:bg-blue-600 transition-colors flex items-center gap-2"
            >
              <svg
                xmlns="http://www.w3.org/2000/svg"
                class="h-5 w-5"
                fill="none"
                viewBox="0 0 24 24"
                stroke="currentColor"
              >
                <path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  d="M12 4.5v15m7.5-7.5h-15"
                />
              </svg>
              添加关系
            </button>
          </div>
        </div>

        <!-- 关系图谱组件 -->
        <div class="gradient-bkg rounded-xl border theme-border p-6">
          <RelationShipGraph
            :relationships="options.characterRSs"
            class="min-h-[400px]"
          />
        </div>
      </div>
    </template>

    <!-- 添加关系弹窗 -->
    <a-modal
      v-model:open="options.isAddRSShow"
      title="添加角色关系"
      @ok="submitCharacterRS"
      :okButtonProps="{
        class:
          'bg-blue-500 hover:bg-blue-600 border-blue-500 hover:border-blue-600',
      }"
    >
      <div class="space-y-4">
        <div class="grid grid-cols-[120px,1fr] items-center gap-4">
          <span class="text-gray-700 dark:text-gray-300">角色 1：</span>
          <a-select
            v-model:value="editRSForm.first_character_id"
            class="w-full"
            placeholder="选择第一个角色"
          >
            <a-select-option
              v-for="character in options.characters"
              :key="character.ID"
              :value="character.ID"
            >
              {{ character.name }}
            </a-select-option>
          </a-select>

          <span class="text-gray-700 dark:text-gray-300">角色 2：</span>
          <a-select
            v-model:value="editRSForm.second_character_id"
            class="w-full"
            placeholder="选择第二个角色"
          >
            <a-select-option
              v-for="character in options.characters"
              :key="character.ID"
              :value="character.ID"
            >
              {{ character.name }}
            </a-select-option>
          </a-select>

          <span class="text-gray-700 dark:text-gray-300">关系名称：</span>
          <input
            v-model="editRSForm.name"
            class="w-full px-4 py-2 border theme-border rounded-lg focus:ring-2 focus:ring-blue-500/20 focus:border-blue-500 dark:bg-zinc-800/50"
            placeholder="例如：师徒、恋人、敌对..."
          />
        </div>

        <div class="space-y-2">
          <span class="text-gray-700 dark:text-gray-300">关系描述：</span>
          <textarea
            v-model="editRSForm.content"
            class="w-full px-4 py-2 border theme-border rounded-lg focus:ring-2 focus:ring-blue-500/20 focus:border-blue-500 dark:bg-zinc-800/50 min-h-[100px]"
            placeholder="描述两个角色之间的关系..."
          />
        </div>

        <div class="flex justify-end">
          <button
            :disabled="options.isCharacterRSGenerating"
            @click="generateCharacterRS"
            class="px-4 py-2 bg-gradient-to-r from-blue-500 to-cyan-500 text-white rounded-lg hover:opacity-90 transition-all flex items-center gap-2"
          >
            <svg
              xmlns="http://www.w3.org/2000/svg"
              class="h-5 w-5"
              fill="none"
              viewBox="0 0 24 24"
              stroke="currentColor"
            >
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                d="M9.813 15.904L9 18.75l-.813-2.846a4.5 4.5 0 00-3.09-3.09L2.25 12l2.846-.813a4.5 4.5 0 003.09-3.09L9 5.25l.813 2.846a4.5 4.5 0 003.09 3.09L15.75 12l-2.846.813a4.5 4.5 0 00-3.09 3.09z"
              />
            </svg>
            {{ options.isCharacterRSGenerating ? "生成中..." : "AI 生成关系" }}
          </button>
        </div>
      </div>
    </a-modal>
    
    <!-- 手动创建角色模态框 -->
    <a-modal
      v-model:open="manualCharacterModal.visible"
      title="创建新角色"
      :confirmLoading="manualCharacterModal.loading"
      @ok="createCharacter"
      @cancel="resetCharacterForm"
      :okButtonProps="{
        class: 'bg-blue-500 hover:bg-blue-600 border-blue-500 hover:border-blue-600',
      }"
    >
      <div class="space-y-6">
        <div class="flex flex-col space-y-2">
          <label class="text-gray-700 dark:text-gray-300 font-medium">角色名称 <span class="text-red-500">*</span></label>
          <input
            v-model="newCharacterForm.name"
            class="input1"
            placeholder="例如：李明、张华、艾米莉..."
          />
        </div>
        
        <div class="flex flex-col space-y-2">
          <label class="text-gray-700 dark:text-gray-300 font-medium">角色描述 <span class="text-red-500">*</span></label>
          <textarea
            v-model="newCharacterForm.description"
            class="input1"
            placeholder="描述角色的性格、背景、特点、目标、动机等..."
          ></textarea>
        </div>
        
        <div class="flex items-center justify-center p-4 bg-gray-50 dark:bg-zinc-900 rounded-lg">
          <p class="text-gray-500 dark:text-gray-400 text-sm">
            角色头像将使用系统默认头像，创建后可通过"生成头像"功能为角色生成AI头像
          </p>
        </div>
      </div>
    </a-modal>

    <!-- 编辑角色模态框 -->
    <a-modal
      v-model:open="editCharacterModal.visible"
      title="编辑角色"
      :confirmLoading="editCharacterModal.loading"
      @ok="updateCharacter"
      :okButtonProps="{
        class: 'bg-blue-500 hover:bg-blue-600 border-blue-500 hover:border-blue-600',
      }"
    >
      <div class="space-y-6">
        <div class="flex flex-col space-y-2">
          <label class="text-gray-700 dark:text-gray-300 font-medium">角色名称 <span class="text-red-500">*</span></label>
          <input
            v-model="editCharacterModal.character.name"
            class="input1"
            placeholder="例如：李明、张华、艾米莉..."
          />
        </div>
        
        <div class="flex flex-col space-y-2">
          <label class="text-gray-700 dark:text-gray-300 font-medium">角色描述 <span class="text-red-500">*</span></label>
          <textarea
            v-model="editCharacterModal.character.description"
            class="input1"
            placeholder="描述角色的性格、背景、特点、目标、动机等..."
            rows="6"
          ></textarea>
        </div>
      </div>
    </a-modal>

    <!-- 删除角色确认模态框 -->
    <a-modal
      v-model:open="deleteModal.visible"
      title="删除角色"
      :confirmLoading="deleteModal.loading"
      @ok="deleteCharacter"
      okText="删除"
      cancelText="取消"
      :okButtonProps="{
        class: 'bg-red-500 hover:bg-red-600 border-red-500 hover:border-red-600',
      }"
    >
      <div class="p-4">
        <div class="flex items-center gap-4 mb-4">
          <div class="text-red-500">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-12 w-12" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
            </svg>
          </div>
          <div>
            <h3 class="text-lg font-medium text-gray-900 dark:text-gray-100">确认删除角色？</h3>
            <p class="text-gray-600 dark:text-gray-400 mt-1">您即将删除角色 <strong class="text-red-500">{{ deleteModal.characterName }}</strong>，此操作不可撤销。</p>
          </div>
        </div>
        
        <div class="bg-amber-50 dark:bg-amber-900/20 p-4 rounded-lg border border-amber-200 dark:border-amber-800">
          <p class="text-amber-800 dark:text-amber-400 text-sm">
            <strong>注意：</strong> 删除角色将同时删除与该角色相关的所有角色关系。
          </p>
        </div>
      </div>
    </a-modal>
  </div>
</template>

<style scoped>
/* 添加动画样式 */
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
.animate__delay-5s {
  animation-delay: 0.5s;
}
.animate__delay-6s {
  animation-delay: 0.6s;
}
.animate__delay-7s {
  animation-delay: 0.7s;
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
.grid > div:nth-child(1) {
  animation-delay: 0.1s;
}
.grid > div:nth-child(2) {
  animation-delay: 0.2s;
}
.grid > div:nth-child(3) {
  animation-delay: 0.3s;
}
.grid > div:nth-child(4) {
  animation-delay: 0.4s;
}
.grid > div:nth-child(5) {
  animation-delay: 0.5s;
}
.grid > div:nth-child(6) {
  animation-delay: 0.6s;
}
.grid > div:nth-child(7) {
  animation-delay: 0.7s;
}
.grid > div:nth-child(8) {
  animation-delay: 0.8s;
}
</style>
