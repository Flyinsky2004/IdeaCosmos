<script setup>
import {onMounted, reactive, ref} from "vue";
import axios from "axios";
import {get} from "@/util/request.js";

const options = reactive({
  myTeam: []
})
const fetchMyTeam = () => {
  get('/api/team/myTeam', {},
      (message, data) => {
        options.myTeam = data
      })
}
onMounted(() => {
  fetchMyTeam()
})
const selectTags = reactive({
  tagsData: [
    "幽默",
    "悬疑",
    "黑暗",
    "科幻",
    "奇幻",
    "浪漫",
    "荒诞",
    "励志",
    "动作",
    "恐怖",
    "史诗",
    "温情",
    "讽刺",
    "文艺",
    "纪实",
    "冒险",
    "家庭剧",
    "超现实",
    "战争",
    "公路",
    "青春",
    "复仇",
    "政治",
    "犯罪",
    "悬幻",
    "心理",
    "怪诞",
    "温暖",
    "现实主义",
    "虚构",
    "哲学",
    "校园",
    "灾难",
    "武侠",
    "神秘",
    "励志成长",
    "古风",
    "穿越",
    "音乐",
    "童话"
  ]
  ,
  selectTags: [
    false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false]
})
const project = reactive({
  project_name: "",
  social_story: "",
  start: "",
  high_point: "",
  resolved: "",
  style: [],
  types: "",
  market_people: "",
  custom_prompt: "",
  team_id: null,
});

const fields = [
  {
    id: "projectName",
    label: "项目名称",
    model: "project_name",
    placeholder: "键入项目名称",
    type: "text",
    required: true
  },
  {
    id: "types",
    label: "项目类型",
    model: "types",
    placeholder: "项目类型...exp:",
    type: "text",
    required: true
  },
  {
    id: "socialStory",
    label: "时代背景",
    model: "social_story",
    placeholder: "故事发生年代，社会背景...",
    type: "textarea",
    required: true
  },
  {
    id: "start",
    label: "情节起点",
    model: "start",
    placeholder: "故事的开始情景...",
    type: "textarea",
    required: true
  },
  {
    id: "highPoint",
    label: "故事高潮/核心冲突",
    model: "high_point",
    placeholder: "角色目标与障碍之间的矛盾...",
    type: "textarea",
    required: true
  },
  {
    id: "resolved",
    label: "高潮与解决",
    model: "resolved",
    placeholder: "故事的转折与结局方向...",
    type: "textarea",
    required: true
  },
  {
    id: "目标群众",
    label: "Market People",
    model: "market_people",
    placeholder: "键入目标群众...exp:'老人，儿童，亲子，青春...'",
    type: "text",
    required: false
  }
];

const typesInput = ref("");

const submitProject = async () => {
  project.types = typesInput.value.split(",").map((type) => type.trim());
  try {
    const response = await axios.post("/api/project/create", project);
    alert(response.data.message || "Project created successfully!");
  } catch (error) {
    console.error(error);
    alert(error.response?.data?.message || "An error occurred!");
  }
}
const handleChange = (tag, checked) => {
  if (checked) {
    project.style.push(tag)
  } else {
    project.style = project.style.filter(t => t !== tag);
  }

};
</script>
<template>
  <div class="h-fit w-full workspace-box animate__animated animate__fadeIn p-8 font-serif">
    <div class="flex flex-nowrap mb-6 gap-2">
      <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor"
           class="size-6 my-auto">
        <path stroke-linecap="round" stroke-linejoin="round"
              d="M12 10.5v6m3-3H9m4.06-7.19-2.12-2.12a1.5 1.5 0 0 0-1.061-.44H4.5A2.25 2.25 0 0 0 2.25 6v12a2.25 2.25 0 0 0 2.25 2.25h15A2.25 2.25 0 0 0 21.75 18V9a2.25 2.25 0 0 0-2.25-2.25h-5.379a1.5 1.5 0 0 1-1.06-.44Z"/>
      </svg>
      <h2 class="text-3xl font-bold font-serif my-auto text-gray-800 dark:text-gray-100">
        创建新项目
      </h2>
    </div>
    <div v-for="field in fields" :key="field.id" class="mb-4">
      <label :for="field.id" class="block text-gray-700 dark:text-gray-300 font-medium mb-2">
        {{ field.label }}
      </label>
      <input
          v-if="field.type === 'text'"
          :type="field.inputType || 'text'"
          v-model="project[field.model]"
          :id="field.id"
          :placeholder="field.placeholder"
          class="input1"
          :required="field.required"
      />
      <textarea
          v-else
          v-model="project[field.model]"
          :id="field.id"
          :placeholder="field.placeholder"
          class="input1"
          :required="field.required"
      ></textarea>
    </div>
    <div class="mb-4">
      <span style="margin-right: 8px">风格类型:</span>
      <a-checkable-tag
          class="text-md font-bold ml-[1px] mt-2"
          v-for="(tag, index) in selectTags.tagsData"
          :key="tag"
          v-model:checked="selectTags[index]"
          @change="checked => handleChange(tag, checked)"
      >
        {{ tag }}
      </a-checkable-tag>
    </div>
    <div class="mb-4 flex flex-nowrap">
      <label for="teamID" class="block text-gray-700 dark:text-gray-300 font-medium my-auto">
        项目团队:
      </label>
      <a-select
          class="my-auto ml-2"
          ref="select"
          v-model:value="project.team_id"
          style="width: 120px"
          @focus="focus"
          @change="handleChange"
      >
        <a-select-option :value="team.ID" v-for="team in options.myTeam">{{ team.username }}</a-select-option>
        <
      </a-select>
    </div>
    <button
        type="submit"
        class="btn1 w-full"
    >
      创建新项目
    </button>
  </div>
</template>

<style>

</style>
