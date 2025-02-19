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

// 添加表单步骤控制
const currentStep = ref(1)
const totalSteps = 4

const steps = [
  { title: '基本信息', description: '设置项目的基本信息' },
  { title: '故事架构', description: '构建故事的核心框架' },
  { title: '目标受众', description: '定义项目的目标群体' },
  { title: '创作设置', description: '设置创作风格和类型' }
]

const nextStep = () => {
  if (currentStep.value < totalSteps) {
    currentStep.value++
  }
}

const prevStep = () => {
  if (currentStep.value > 1) {
    currentStep.value--
  }
}

// 优化表单验证
const validateStep = (step) => {
  switch(step) {
    case 1:
      return project.project_name.length >= 6
    case 2:
      return project.social_story.length >= 6 &&
             project.start.length >= 6 &&
             project.high_point.length >= 6 &&
             project.resolved.length >= 6
    case 3:
      return project.market_people.length > 0
    case 4:
      return project.types && project.style.length > 0 && project.team_id
    default:
      return false
  }
}
</script>
<template>
  <div class="max-w-4xl mx-auto animate__animated animate__fadeIn p-8">
    <!-- 进度指示器 -->
    <div class="mb-8">
      <div class="flex justify-between items-center mb-4">
        <div v-for="(step, index) in steps" :key="index"
             class="flex-1 relative">
          <div class="flex items-center">
            <div class="mt-2" :class="[
              'w-8 h-8 rounded-full flex items-center justify-center text-sm font-medium',
              currentStep > index + 1 ? 'bg-blue-500 text-white' :
              currentStep === index + 1 ? 'bg-blue-100 text-blue-600 border-2 border-blue-500' :
              'bg-gray-100 text-gray-500 dark:bg-gray-800'
            ]">
              {{ index + 1 }}
            </div>
            <div class="ml-3">
              <p class="text-sm font-medium text-gray-900 dark:text-gray-100">{{ step.title }}</p>
              <p class="text-xs text-gray-500 dark:text-gray-400">{{ step.description }}</p>
            </div>
          </div>
          <div v-if="index < steps.length - 1" 
               :class="[
                 'absolute top-4 w-full h-0.5',
                 currentStep > index + 1 ? 'bg-blue-500' : 'bg-gray-200 dark:bg-gray-700'
               ]">
          </div>
        </div>
      </div>
    </div>

    <!-- 表单内容 -->
    <div class="bg-white dark:bg-zinc-900 rounded-xl p-6 shadow-sm border theme-border">
      <!-- 步骤 1: 基本信息 -->
      <div v-if="currentStep === 1" class="space-y-6 animate__animated animate__fadeIn">
        <div class="mb-6">
          <label for="projectName" class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
            项目名称
          </label>
          <input
            v-model="project.project_name"
            type="text"
            id="projectName"
            class="input1"
            placeholder="请输入项目名称"
          />
          <p class="mt-1 text-sm text-gray-500">
            至少输入6个字符
          </p>
        </div>
      </div>

      <!-- 步骤 2: 故事架构 -->
      <div v-if="currentStep === 2" class="space-y-6 animate__animated animate__fadeIn">
        <div v-for="field in fields.filter(f => f.model !== 'project_name')" 
             :key="field.id" 
             class="mb-6">
          <label :for="field.id" class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
            {{ field.label }}
          </label>
          <textarea
            v-model="project[field.model]"
            :id="field.id"
            :placeholder="field.placeholder"
            class="input1 min-h-[120px]"
          ></textarea>
        </div>
      </div>

      <!-- 步骤 3: 目标受众 -->
      <div v-if="currentStep === 3" class="animate__animated animate__fadeIn">
        <h3 class="text-lg font-medium text-gray-900 dark:text-gray-100 mb-4">选择目标受众</h3>
        <div class="space-y-2">
          <div class="flex flex-wrap gap-2">
            <a-checkable-tag
              v-for="(tag, index) in selectPeople.tagsData"
              :key="tag"
              v-model:checked="selectPeople.selectTags[index]"
              @change="checked => handlePeopleChange(tag, checked)"
              class="px-3 py-1.5 rounded-full text-sm"
            >
              {{ tag }}
            </a-checkable-tag>
          </div>
          {{ selectTags }}
        </div>
      </div>

      <!-- 步骤 4: 创作设置 -->
      <div v-if="currentStep === 4" class="space-y-6 animate__animated animate__fadeIn">
        <!-- 内容类型选择 -->
        <div class="mb-6">
          <h3 class="text-lg font-medium text-gray-900 dark:text-gray-100 mb-4">内容类型</h3>
          <a-radio-group v-model:value="project.types" class="flex flex-wrap gap-2">
            <a-radio-button 
              v-for="type in [
                '电影剧本', '短视频脚本', '剧情短片脚本', '连载小说',
                '喜剧小品剧本', 'Vlog脚本', '广告脚本', '商品测评视频脚本',
                '品牌故事视频', '创意动画脚本', '电竞或游戏视频脚本'
              ]" 
              :key="type" 
              :value="type"
              class="mb-2"
            >
              {{ type }}
            </a-radio-button>
          </a-radio-group>
        </div>

        <!-- 风格类型选择 -->
        <div class="mb-6">
          <h3 class="text-lg font-medium text-gray-900 dark:text-gray-100 mb-4">风格类型</h3>
          <div class="flex flex-wrap gap-2">
            <a-checkable-tag
              v-for="(tag, index) in selectTags.tagsData"
              :key="tag"
              v-model:checked="selectTags.selectTags[index]"
              @change="checked => handleChange(tag, checked)"
              class="px-3 py-1.5 rounded-full text-sm"
            >
              {{ tag }}
            </a-checkable-tag>
          </div>
        </div>

        <!-- 团队选择 -->
        <div class="mb-6">
          <h3 class="text-lg font-medium text-gray-900 dark:text-gray-100 mb-4">选择项目团队</h3>
          <a-select
            v-model:value="project.team_id"
            class="w-full"
            placeholder="请选择项目团队"
          >
            <a-select-option 
              v-for="team in options.myTeam"
              :key="team.ID"
              :value="team.ID"
            >
              {{ team.username }}
            </a-select-option>
          </a-select>
        </div>
      </div>

      <!-- 导航按钮 -->
      <div class="flex justify-between mt-8">
        <button
          v-if="currentStep > 1"
          @click="prevStep"
          class="px-4 py-2 border border-gray-300 dark:border-gray-600 rounded-lg text-gray-700 dark:text-gray-300 hover:bg-gray-50 dark:hover:bg-gray-800 transition-colors"
        >
          上一步
        </button>
        <div class="ml-auto">
          <button
            v-if="currentStep < totalSteps"
            @click="nextStep"
            :disabled="!validateStep(currentStep)"
            class="px-4 py-2 bg-blue-500 text-white rounded-lg hover:bg-blue-600 transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
          >
            下一步
          </button>
          <button
            v-else
            @click="submitProject"
            :disabled="!validateStep(currentStep)"
            class="px-4 py-2 bg-gradient-to-r from-blue-500 to-cyan-500 text-white rounded-lg hover:opacity-90 transition-all disabled:opacity-50 disabled:cursor-not-allowed"
          >
            创建项目
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>

</style>
