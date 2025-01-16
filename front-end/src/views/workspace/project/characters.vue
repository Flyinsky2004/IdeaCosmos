<script setup>
import {onMounted, reactive} from "vue";
import {get, post, postJSON} from "@/util/request.js";
import {message} from "ant-design-vue";
import router from "@/router/index.js";
import Loader from "@/components/loader.vue";
import {parseDateTime, washJSONStr} from "@/util/common.js";
import {imagePrefix} from "@/util/VARRIBLES.js";
import RelationShipGraph from "@/components/RelationShipGraph.vue";

const project = JSON.parse(localStorage.getItem("project"));
onMounted(() => {
  fetchCharacters()
  fetchCharacterRs()
})
const options = reactive({
  characters: [],
  isCharactersGenerating: false,
  generateResults: [],
  currentMoveId: -1,
  isCharacterAvatarGenerating: false,
  characterGeneratingId: -1,
  isAddRSShow: false,
  editRSMode: false,
  currentRS: {},
  isCharacterRSGenerating: false,
  characterRSs:[]
})
const fetchCharacterRs = () => {
  get('/api/project/characterRS/getAll', {
    project_id: project.ID
  }, (messager, data) => {
    options.characterRSs = data
  }, (messager, data) => {
    message.warning(messager)
  }, (messager, data) => {
    message.error(messager)
  })
}
const fetchCharacters = () => {
  post('/api/project/getCharacters', {
    project_id: project.ID
  }, (messager, data) => {
    options.characters = data
  }, (messager, data) => {
    message.warning(messager)
  }, (messager, data) => {
    message.error(messager)
  })
}
const generateCharacter = () => {
  options.isCharactersGenerating = true;
  post('/api/project/generateCharacter', {
    project_id: project.ID
  }, (messager, data) => {
    const raw = data.choices[0].message.content
    options.generateResults = JSON.parse(washJSONStr(raw))
    for (let i = 0; i < options.generateResults.length; i++) {
      options.generateResults[i].project_id = project.ID
    }
    options.isCharactersGenerating = false
  })
}
const applyResult = () => {
  options.generateResults[0].project_id = project.ID
  postJSON('/api/project/createCharacterArray', options.generateResults,
      (messageer, data) => {
        message.success(messageer)
        fetchCharacters()
        dropResult()
      }, (messageer, data) => {
        message.warning(messageer)
      }, (messageer, data) => {
        message.error(messageer)
      })
}
const dropResult = () => {
  options.generateResults = []
}
const moveIn = (id) => {
  options.currentMoveId = id
}
const moveOut = () => {
  options.currentMoveId = -1
}
const generateAvatar = async (character) => {
  options.characterGeneratingId = character.ID
  options.isCharacterAvatarGenerating = true
  await new Promise(resolve => {
    post('/api/project/generateCharacterAvatar', {
      character_id: character.ID
    }, (messageer, data) => {
      message.success(messageer)
      character.avatar = data
      for (let i = 0; i < options.characters.length; i++) {
        if (options.characters[i].ID === character.ID) {
          options.characters[i].avatar = data
          break;
        }
      }
      options.isCharacterAvatarGenerating = false
      resolve()
    }, (messageer, data) => {
      message.warning(messageer)
      options.isCharacterAvatarGenerating = false
    }, (messageer, data) => {
      message.error(messageer)
      options.isCharacterAvatarGenerating = false
    })
  }).then(() => {

  })
}
const editRSForm = reactive({
  first_character_id: 1,
  second_character_id: 2,
  name: '',
  content: ''
})
const generateCharacterRS = () => {
  if (editRSForm.first_character_id === editRSForm.second_character_id) {
    message.info("角色选择不合法")
    return
  }
  options.isCharacterRSGenerating = true;
  post('/api/project/generateCharacterRS', {
    firstCharacterId: editRSForm.first_character_id,
    secondCharacterId: editRSForm.second_character_id,
  }, (messager, data) => {
    const raw = data.choices[0].message.content
    const washed = JSON.parse(washJSONStr(raw))
    editRSForm.name = washed.name
    editRSForm.content = washed.content
    options.isCharacterRSGenerating = false
    message.success("生成成功")
  }, (messageer, data) => {
    message.warning(messageer)
  }, (messageer, data) => {
    message.error(messageer)
  })
}
const submitCharacterRS = () => {
  if (options.editRSMode) {

  } else {
    postJSON('/api/project/characterRS/create', editRSForm,
        (messager, data) => {
          message.success(messager)
          options.isAddRSShow = false
          fetchCharacterRs()
        }, (messageer, data) => {
          message.warning(messageer)
        }, (messageer, data) => {
          message.error(messageer)
        })
  }
}

</script>

<template>
  <div>
    <h1 class="text-3xl text-blue-500">角色管理</h1>
    <div class="w-full grid grid-cols-5 gap-2">
      <div v-for="character in options.characters" class="border rounded-xl dark:border-[rgb(118,118,118)] outline-[1px] min-h-56 text-theme-switch
hover:bg-gray-100/50 active:bg-gray-100/90 dark:hover:bg-gray-900/5 dark:active:bg-gray-900/10 cursor-pointer"
           @mouseover="moveIn(character.ID)" @mouseleave="moveOut()">
        <div class="p-2" v-if="options.nowHoverId !== project.ID">
          <a-popover title="编辑作品封面">
            <template #content>
              <div class="flex flex-nowrap gap-2">
                <button class="color-mixed-button" @click="generateAvatar(character)">AI生成</button>
                <button class="btn1">手动上传</button>
                <button class="transparent-button">清空封面</button>
              </div>
            </template>
            <div class="w-full aspect-square border theme-border rounded-xl">
              <div v-if="options.isCharacterAvatarGenerating && options.characterGeneratingId === character.ID"
                   class="w-full h-full flex bg-slate-50/20 dark:bg-[#242424]/50 rounded-xl">
                <div class="mx-auto my-auto place-items-center">
                  <loader class="m-2"/>
                  <span class="font-sans font-bold">根据人物设定绘画中...</span>
                </div>
              </div>
              <div v-else-if="character.avatar === ''"
                   class="w-full aspect-square flex bg-slate-50/20 dark:bg-[#242424]/50 rounded-xl">
                <div class="mx-auto my-auto place-items-center">
                  <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="currentColor" class="size-6">
                    <path fill-rule="evenodd"
                          d="M18.685 19.097A9.723 9.723 0 0 0 21.75 12c0-5.385-4.365-9.75-9.75-9.75S2.25 6.615 2.25 12a9.723 9.723 0 0 0 3.065 7.097A9.716 9.716 0 0 0 12 21.75a9.716 9.716 0 0 0 6.685-2.653Zm-12.54-1.285A7.486 7.486 0 0 1 12 15a7.486 7.486 0 0 1 5.855 2.812A8.224 8.224 0 0 1 12 20.25a8.224 8.224 0 0 1-5.855-2.438ZM15.75 9a3.75 3.75 0 1 1-7.5 0 3.75 3.75 0 0 1 7.5 0Z"
                          clip-rule="evenodd"/>
                  </svg>
                  暂无角色画像
                </div>
              </div>
              <img class="h-full w-auto border theme-border rounded-xl" :src="imagePrefix+character.avatar" v-else
                   alt="用户头像"/>
            </div>
          </a-popover>
          <h1 class="text-2xl text-blue-500">角色姓名:{{ character.name }}</h1>
          <h1 class="text-sm">角色介绍:{{ character.description }}</h1>
        </div>
        <div @click="enterProject(project)"
             class="flex flex-col w-full h-full animate__animated animate__fadeIn animate__faster p-2
          bg-gray-100/10 active:bg-gray-200/90 dark:bg-gray-950/10 dark:active:bg-gray-950/15 cursor-pointer" v-else>
          <h1 class="text-2xl text-blue-500">#{{ project.ID }}{{ project.project_name }}</h1>
          <h1 class="text-sm text-blue-600">项目团队:{{ project.team.username }}</h1>
          <h1 class="text-sm font-bold text-blue-400">项目简介：</h1>
          <h1 class="text-sm">{{ project.social_story }}...</h1>
          <h1 class="text-sm font-bold text-blue-400">创建时间：</h1>
          <span class="text-sm">{{ parseDateTime(project.CreatedAt) }}</span>
          <span class="mx-auto my-auto flex flex-nowrap">
              <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5"
                   stroke="currentColor" class="size-5">
  <path stroke-linecap="round" stroke-linejoin="round"
        d="M15.59 14.37a6 6 0 0 1-5.84 7.38v-4.8m5.84-2.58a14.98 14.98 0 0 0 6.16-12.12A14.98 14.98 0 0 0 9.631 8.41m5.96 5.96a14.926 14.926 0 0 1-5.841 2.58m-.119-8.54a6 6 0 0 0-7.381 5.84h4.8m2.581-5.84a14.927 14.927 0 0 0-2.58 5.84m2.699 2.7c-.103.021-.207.041-.311.06a15.09 15.09 0 0 1-2.448-2.448 14.9 14.9 0 0 1 .06-.312m-2.24 2.39a4.493 4.493 0 0 0-1.757 4.306 4.493 4.493 0 0 0 4.306-1.758M16.5 9a1.5 1.5 0 1 1-3 0 1.5 1.5 0 0 1 3 0Z"/>
</svg>

              进入项目空间</span>
        </div>
      </div>
      <div class="border border-dashed rounded-xl dark:border-[rgb(118,118,118)] outline-[1px] min-h-56 place-items-center place-content-center text-gray-400
hover:bg-gray-100/50 active:bg-gray-100/90 dark:hover:dark:bg-gray-950/10 dark:active:bg-gray-900/10 cursor-pointer"
           @click="router.push('/workspace/newProject')">
        <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5"
             stroke="currentColor" class="size-16">
          <path stroke-linecap="round" stroke-linejoin="round" d="M12 9v6m3-3H9m12 0a9 9 0 1 1-18 0 9 9 0 0 1 18 0Z"/>
        </svg>
      </div>
      <div v-if="options.generateResults.length === 0" @click="generateCharacter"
           class="border border-dashed rounded-xl dark:border-[rgb(118,118,118)] outline-[1px] min-h-56 place-items-center place-content-center
hover:bg-gray-100/50 active:bg-gray-100/90 dark:hover:dark:bg-gray-950/10 dark:active:bg-gray-900/10 cursor-pointer">
        <div v-if="!options.isCharactersGenerating">
          <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="currentColor"
               class="size-7 dark:text-purple-400">
            <path fill-rule="evenodd"
                  d="M9 4.5a.75.75 0 0 1 .721.544l.813 2.846a3.75 3.75 0 0 0 2.576 2.576l2.846.813a.75.75 0 0 1 0 1.442l-2.846.813a3.75 3.75 0 0 0-2.576 2.576l-.813 2.846a.75.75 0 0 1-1.442 0l-.813-2.846a3.75 3.75 0 0 0-2.576-2.576l-2.846-.813a.75.75 0 0 1 0-1.442l2.846-.813A3.75 3.75 0 0 0 7.466 7.89l.813-2.846A.75.75 0 0 1 9 4.5ZM18 1.5a.75.75 0 0 1 .728.568l.258 1.036c.236.94.97 1.674 1.91 1.91l1.036.258a.75.75 0 0 1 0 1.456l-1.036.258c-.94.236-1.674.97-1.91 1.91l-.258 1.036a.75.75 0 0 1-1.456 0l-.258-1.036a2.625 2.625 0 0 0-1.91-1.91l-1.036-.258a.75.75 0 0 1 0-1.456l1.036-.258a2.625 2.625 0 0 0 1.91-1.91l.258-1.036A.75.75 0 0 1 18 1.5ZM16.5 15a.75.75 0 0 1 .712.513l.394 1.183c.15.447.5.799.948.948l1.183.395a.75.75 0 0 1 0 1.422l-1.183.395c-.447.15-.799.5-.948.948l-.395 1.183a.75.75 0 0 1-1.422 0l-.395-1.183a1.5 1.5 0 0 0-.948-.948l-1.183-.395a.75.75 0 0 1 0-1.422l1.183-.395c.447-.15.799-.5.948-.948l.395-1.183A.75.75 0 0 1 16.5 15Z"
                  clip-rule="evenodd"/>
          </svg>
        </div>
        <div v-else class="text-center place-items-center">
          <loader class="mb-4"/>
          <span>正在理解你的项目...</span>
        </div>
      </div>
    </div>
    <div class="bg-green-300/10 rounded-xl">
      <div class="w-full grid grid-cols-5 gap-2 mt-2 p-2 rounded-xl"
           v-if="options.generateResults.length !== 0">
        <div class="border rounded-xl dark:border-[rgb(118,118,118)] outline-[1px] min-h-56 text-theme-switch
hover:bg-gray-100/50 active:bg-gray-100/90 dark:hover:bg-gray-900/5 dark:active:bg-gray-900/10 cursor-pointer p-2"
             v-for="character in options.generateResults">
          <div class="w-full h-80 border theme-border rounded-xl">
            <div class="w-full h-full flex bg-slate-50/20 dark:bg-[#242424]/50 rounded-xl">
              <div class="mx-auto my-auto place-items-center">
                <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="currentColor" class="size-6">
                  <path fill-rule="evenodd"
                        d="M18.685 19.097A9.723 9.723 0 0 0 21.75 12c0-5.385-4.365-9.75-9.75-9.75S2.25 6.615 2.25 12a9.723 9.723 0 0 0 3.065 7.097A9.716 9.716 0 0 0 12 21.75a9.716 9.716 0 0 0 6.685-2.653Zm-12.54-1.285A7.486 7.486 0 0 1 12 15a7.486 7.486 0 0 1 5.855 2.812A8.224 8.224 0 0 1 12 20.25a8.224 8.224 0 0 1-5.855-2.438ZM15.75 9a3.75 3.75 0 1 1-7.5 0 3.75 3.75 0 0 1 7.5 0Z"
                        clip-rule="evenodd"/>
                </svg>
                暂无角色画像
              </div>
            </div>
          </div>
          <h1 class="text-2xl text-blue-500">角色姓名:{{ character.name }}</h1>
          <h1 class="text-sm">角色介绍:{{ character.description }}</h1>
        </div>
      </div>
      <div class="flex flex-nowrap gap-2 mb-2 p-2 place-items-center" v-if="options.generateResults.length !== 0">
        <span>您是否要采用生成的角色?</span>
        <button class="basic-prinary-button my-auto" @click="applyResult">
          采用
        </button>
        <button class="basic-error-button" @click="dropResult">
          丢弃
        </button>
      </div>
    </div>


  </div>
  <div class="mt-3 mb-3">
    <h1 class="text-3xl text-blue-500">角色联系</h1>
    <div class="flex flex-nowrap gap-2">
      <a-modal v-model:open="options.isAddRSShow" title="角色关系维护" @ok="submitCharacterRS()">
        <div class="grid grid-cols-[1fr,2fr] w-fit gap-2 place-items-center">
          <span>角色1：</span>
          <a-select
              ref="select"
              v-model:value="editRSForm.first_character_id"
              style="width: 120px"
              @focus="focus"
          >
            <a-select-option v-for="character in options.characters" :value="character.ID">{{ character.name }}
            </a-select-option>
          </a-select>
          <span>角色2：</span>
          <a-select
              ref="select"
              v-model:value="editRSForm.second_character_id"
              style="width: 120px"
              @focus="focus"
          >
            <a-select-option v-for="character in options.characters" :value="character.ID">{{ character.name }}
            </a-select-option>
          </a-select>
          <span>关系名称:</span> <input v-model="editRSForm.name" class="input1">
          <span>人物故事:</span>
        </div>
        <textarea v-model="editRSForm.content" class="input1 mt-2"/>
        <div class="flex w-full">
          <div class="flex-grow"></div>
          <button :disabled="options.isCharacterRSGenerating" class="color-mixed-button" @click="generateCharacterRS">
            {{ options.isCharacterRSGenerating ? '正在努力思考中...' : 'AI生成' }}
          </button>
        </div>
        {{ editRSForm }}

      </a-modal>
      <button class="basic-success-button my-auto"
              @click="options.editRSMode = false;options.currentRS = {};options.isAddRSShow = true">
        添加关系
      </button>
      <button class="basic-prinary-button my-auto" @click="editRealationShip">
        编辑关系
      </button>
      <button class="basic-error-button my-auto" @click="manageRealationShip">
        删除关系
      </button>
    </div>
    <RelationShipGraph :relationships="options.characterRSs" class="border theme-border rounded-xl mt-6"/>
  </div>

</template>

<style scoped>

</style>