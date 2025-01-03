<script setup>
import {onMounted, reactive, ref} from "vue";
import {get, post, postJSON} from "@/util/request.js";
import {message} from "ant-design-vue";
import router from "@/router/index.js";

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
const selectPeople = reactive({
  tagsData: [
    "儿童",
    "青少年",
    "大学生",
    "年轻人",
    "成年人",
    "中老年人",
    "家庭观众",
    "男性",
    "女性",
    "情侣",
    "科幻迷",
    "奇幻迷",
    "历史爱好者",
    "冒险爱好者",
    "恐怖爱好者",
    "喜剧迷",
    "文艺爱好者",
    "悬疑迷",
    "动作片爱好者",
    "家庭伦理关注者",
    "教育从业者",
    "学生",
    "职场新人",
    "企业高管",
    "创业者",
    "社会活动家",
    "心理学爱好者",
    "哲学爱好者",
    "音乐爱好者",
    "电影迷",
    "小说读者",
    "旅行者",
    "历史学者",
    "文化研究者",
    "游戏玩家",
    "运动爱好者",
    "环保主义者",
    "科技爱好者",
    "战争历史迷",
    "未来主义者"
  ],
  selectTags: [
    false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false]

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
  market_people: [],
  custom_prompt: "",
  team_id: null,
});
const checkProjectValidate = (form) => {
  // 遍历表单的每个属性进行检查
  for (const key in form) {
    if (Object.prototype.hasOwnProperty.call(form, key)) {
      const value = form[key];

      // 如果是数组类型，确保数组不为空
      if (Array.isArray(value) && value.length === 0) {
        message.info(value+"至少选择一个可选类型")
        return false;
      }

      // 如果是字符串类型或其他基础类型，确保非空
      if (  project.project_name.length > 6 &&
          project.social_story.length > 6 &&
          project.start.length > 6 &&
          project.high_point.length > 6 &&
          project.resolved.length > 6 &&
          project.types.length > 6 &&
          project.custom_prompt.length > 6) {
        message.info("至少大于6个字符")
        return false;
      }

      // 如果是 null 类型，也视为未填写
      if (value === null) {
        message.info("请选择项目团队")
        return false;
      }
    }
  }

  // 所有字段都通过了检查，返回 true
  return true;
}
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
  }
];

const typesInput = ref("");

const submitProject = async () => {
  if (checkProjectValidate(project)){
    postJSON('/api/project/createProject',
        project,
        (messager, data) => {
          message.success(messager)
          router.push("/workspace/projects")
        },
        (messager, data) => {
          message.warning(messager);
        }, (messager, data) => {
          message.error(messager, data);
        })
  }
}
const handleChange = (tag, checked) => {
  if (checked) {
    project.style.push(tag)
  } else {
    project.style = project.style.filter(t => t !== tag);
  }
};
const handlePeopleChange = (tag, checked) => {
  if (checked) {
    project.market_people.push(tag)
  } else {
    project.market_people = project.market_people.filter(t => t !== tag);
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
    <div>
      目标群众:
      <a-checkable-tag
          class="text-md font-bold ml-[1px] mt-2"
          v-for="(tag, index) in selectPeople.tagsData"
          :key="tag"
          v-model:checked="selectPeople[index]"
          @change="checked => handlePeopleChange(tag, checked)"
      >
        {{ tag }}
      </a-checkable-tag>
    </div>
    <div class="mb-4 mt-4">
      <span style="margin-right: 8px;">内容类型:</span>
      <a-radio-group v-model:value="project.types">
        <a-radio-button value="电影剧本">电影剧本</a-radio-button>
        <a-radio-button value="短视频脚本">短视频脚本</a-radio-button>
        <a-radio-button value="剧情短片脚本">剧情短片脚本</a-radio-button>
        <a-radio-button value="连载小说">连载小说</a-radio-button>
        <a-radio-button value="喜剧小品剧本">喜剧小品剧本</a-radio-button>
        <a-radio-button value="Vlog脚本">Vlog脚本</a-radio-button>
        <a-radio-button value="广告脚本">广告脚本</a-radio-button>
        <a-radio-button value="商品测评视频脚本">商品测评视频脚本</a-radio-button>
        <a-radio-button value="品牌故事视频">品牌故事视频</a-radio-button>
        <a-radio-button value="创意动画脚本">创意动画脚本</a-radio-button>
        <a-radio-button value="small">电竞或游戏视频脚本</a-radio-button>
      </a-radio-group>
    </div>
    <div class="mb-4">
      <span style="margin-right: 8px">风格类型: </span>
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
      >
        <a-select-option :value="team.ID" v-for="team in options.myTeam">{{ team.username }}</a-select-option>
        <
      </a-select>
    </div>
    <button
        type="submit"
        class="btn1 w-full"
        @click="submitProject"
    >
      创建新项目
    </button>
  </div>
</template>

<style>

</style>
