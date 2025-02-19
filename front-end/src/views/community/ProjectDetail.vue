<script setup>
import { ref, onMounted, watch, nextTick, onBeforeMount } from "vue";
import { useRoute } from "vue-router";
import { get, post , postJSON } from "@/util/request";
import { message } from "ant-design-vue";
import { BACKEND_DOMAIN, imagePrefix } from "@/util/VARRIBLES";
import Loader from "@/components/loader.vue";
import * as echarts from "echarts";
import { useUserStore } from "@/stores/user";

const route = useRoute();
const project = ref(null);
const loading = ref(true);
const chaptersLoading = ref(true);

// 角色关系图相关
const characterChart = ref(null);
const characters = ref([]);
const relationships = ref([]);

// 章节列表
const chapters = ref([]);

const userStore = useUserStore();
const commentContent = ref("");
const comments = ref([]);
const commentsLoading = ref(true);

onMounted(async () => {
  const projectId = route.params.id;
  try {
    await fetchProjectDetail(projectId);
    await Promise.all([
      fetchCharacters(projectId),
      fetchChapters(projectId)
    ]);
  } catch (error) {
    message.error('加载数据失败');
  }
});

const fetchProjectDetail = (projectId) => {
  loading.value = true;
  get(
    "/api/public/getProjectDetail",
    { project_id: projectId },
    (messager, data) => {
      project.value = data;
      loading.value = false;
    },
    (messager) => {
      message.warning(messager);
      loading.value = false;
    },
    (messager) => {
      message.error(messager);
      loading.value = false;
    }
  );
};

// 获取角色数据
const fetchCharacters = (projectId) => {
  get(
    "/api/public/getProjectCharacters",
    { project_id: projectId },
    (messager, data) => {
      characters.value = data.characters;
      relationships.value = data.relationships;
      nextTick(() => {
        initCharacterChart();
      });
    },
    (messager) => message.warning(messager),
    (messager) => message.error(messager)
  );
};

// 获取章节数据
const fetchChapters = async (projectId) => {
  chaptersLoading.value = true;
  try {
    const response = await new Promise((resolve, reject) => {
      get(
        "/api/public/getProjectChapters",
        { project_id: projectId },
        (messager, data) => resolve(data),
        (messager) => reject(messager),
        (messager) => reject(messager)
      );
    });
    
    // 确保 CurrentVersion 字段存在
    chapters.value = response.map(chapter => ({
      ...chapter,
      CurrentVersion: chapter.CurrentVersion || { ID: 0 }
    }));
  } catch (error) {
    message.error(error);
  } finally {
    chaptersLoading.value = false;
  }
};

// 获取评论列表
const fetchComments = async (projectId) => {
  commentsLoading.value = true;
  try {
    const response = await new Promise((resolve, reject) => {
      get(
        "/api/public/getProjectComments",
        { project_id: projectId },
        (messager, data) => resolve(data),
        (messager) => reject(messager),
        (messager) => reject(messager)
      );
    });
    comments.value = response;
  } catch (error) {
    message.error("获取评论失败");
  } finally {
    commentsLoading.value = false;
  }
};

// 提交评论
const submitComment = async () => {
  if (!userStore.isLogin) {
    message.warning("请先登录后再评论");
    return;
  }
  
  if (!commentContent.value.trim()) {
    message.warning("请输入评论内容");
    return;
  }

  try {
    await new Promise((resolve, reject) => {
      postJSON(
        "/api/user/addProjectComment",
        {
          content: commentContent.value,
          project_id: project.value.ID
        },
        (messager, data) => {
          comments.value.unshift(data);
          commentContent.value = "";
          message.success("评论成功");
          resolve();
        },
        (messager) => reject(messager),
        (messager) => reject(messager)
      );
    });
  } catch (error) {
    message.error("评论失败");
  }
};

// 生成随机颜色
const getRandomColor = () => {
  const colors = [
    "#FF6B6B", "#4ECDC4", "#45B7D1", "#96CEB4", "#FFEEAD",
    "#D4A5A5", "#9B59B6", "#3498DB", "#1ABC9C", "#F1C40F",
    "#E74C3C", "#2ECC71", "#34495E", "#16A085", "#27AE60",
    "#FFB347", "#A8E6CF", "#3D84A8", "#46B3E6", "#4A90E2",
    "#D35400", "#8E44AD", "#2980B9", "#F39C12", "#7F8C8D",
    "#E67E22", "#C0392B", "#6C5CE7", "#00B894", "#FDA7DF",
    "#FF9FF3", "#FFA07A", "#20B2AA", "#87CEEB", "#DDA0DD"
  ];
  return colors[Math.floor(Math.random() * colors.length)];
};

// 初始化角色关系图
const initCharacterChart = () => {
  if (!characterChart.value) return;

  const chartInstance = echarts.init(characterChart.value);

  // 构建图表数据
  const nodes = characters.value.map((char) => ({
    name: char.name,
    id: char.ID.toString(),
    symbolSize: 60,
    category: 0,
    itemStyle: {
      color: getRandomColor(),
    },
    label: {
      show: true,
      color: "#fff",
      fontSize: 14,
      fontWeight: "bold",
      backgroundColor: "#666",
      padding: [4, 8],
      borderRadius: 4,
      distance: 5,
    },
    value: char.description,
  }));

  const categories = [{ name: "角色" }];

  const links = relationships.value.map((rel) => ({
    source: rel.first_character_id.toString(),
    target: rel.second_character_id.toString(),
    name: rel.name,
    value: rel.content,
    label: {
      show: true,
      formatter: rel.name,
      fontSize: 12,
      color: "#fff",
      backgroundColor: "rgba(0,0,0,0.6)",
      padding: [4, 8],
      borderRadius: 4,
    },
    lineStyle: {
      width: 2,
      curveness: 0.2,
      color: "#666",
      opacity: 0.8,
    },
  }));

  const option = {
    title: {
      text: "角色关系图谱",
      subtext: "可拖拽调整位置",
      top: "bottom",
      left: "right"
    },
    tooltip: {
      show: true,
      confine: true,
      enterable: true,
      formatter: function (params) {
        const title = params.dataType === "edge" ? "关系" : "角色";
        const name = params.data.name;
        const desc = params.data.value || "暂无描述";
        
        return `
          <div style="
            width: 200px;
            max-height: 300px;
            overflow-y: auto;
            font-size: 14px;
            line-height: 1.5;
          ">
            <div style="
              font-weight: bold;
              margin-bottom: 8px;
              color: #fff;
              font-size: 15px;
              white-space: normal;
              word-break: break-all;
            ">
              ${title}：${name}
            </div>
            <div style="
              color: rgba(255,255,255,0.9);
              white-space: normal;
              word-break: break-all;
            ">
              描述：${desc}
            </div>
          </div>
        `;
      },
      backgroundColor: "rgba(0,0,0,0.75)",
      padding: [10, 15],
      borderRadius: 8,
      textStyle: {
        fontSize: 14,
        lineHeight: 20,
        rich: {
          title: {
            fontSize: 15,
            fontWeight: 'bold',
            lineHeight: 25
          },
          content: {
            lineHeight: 20,
            color: 'rgba(255,255,255,0.9)'
          }
        }
      }
    },
    legend: [
      {
        data: categories.map((a) => a.name),
      },
    ],
    animationDuration: 1500,
    animationEasingUpdate: "quinticInOut",
    series: [
      {
        name: "角色关系",
        type: "graph",
        layout: "force",
        data: nodes,
        links: links,
        categories: categories,
        roam: true,
        draggable: true,
        label: {
          position: "right",
          formatter: "{b}",
        },
        force: {
          repulsion: 200,
          gravity: 0.1,
          edgeLength: 120,
          layoutAnimation: true,
        },
        lineStyle: {
          color: "source",
          curveness: 0.3,
        },
        emphasis: {
          focus: "adjacency",
          lineStyle: {
            width: 4,
          },
        },
      },
    ],
  };

  chartInstance.setOption(option);

  // 监听窗口大小变化
  window.addEventListener("resize", () => {
    chartInstance.resize();
  });
};

// 在获取到项目信息后加载其他数据
watch(
  () => project.value,
  (newProject) => {
    if (newProject) {
      fetchCharacters(newProject.ID);
      fetchChapters(newProject.ID);
      fetchComments(newProject.ID);
    }
  }
);

// 添加角色列表展示
const getAvatarUrl = (avatar) => {
  return avatar ? imagePrefix + avatar : "/default-avatar.png";
};

// 格式化时间
const formatDate = (date) => {
  const d = new Date(date);
  return d.toLocaleDateString("zh-CN", {
    year: "numeric",
    month: "long",
    day: "numeric",
    hour: "2-digit",
    minute: "2-digit",
  });
};

// 获取章节状态标签
const getChapterStatus = (chapter) => {
  if (!chapter?.current_version) return { color: 'orange', text: '未开始创作' , isPublished: false }
  if (chapter.current_version?.ID === 0) return { color: 'orange', text: '未开始创作' , isPublished: false }
  return { color: 'green', text: '已发布' , isPublished: true }
};

// 获取章节卡片背景色
const getChapterCardStyle = (chapter) => {
  const status = getChapterStatus(chapter);
  const styles = {
    orange: "bg-orange-50 dark:bg-orange-950/80",
    green: "bg-indigo-50 dark:bg-indigo-900/10",
    blue: "bg-blue-50 dark:bg-blue-900/10",
  };
  return styles[status.color];
};

// 格式化评论时间
const formatCommentDate = (date) => {
  const d = new Date(date);
  return d.toLocaleString('zh-CN', {
    year: 'numeric',
    month: 'long',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  });
};
</script>

<template>
  <div class="min-h-screen bg-gray-50 dark:bg-zinc-900/20 p-6">
    <!-- 加载状态 -->
    <div v-if="loading" class="flex items-center justify-center min-h-[400px]">
      <loader />
    </div>

    <template v-else-if="project">
      <div class="max-w-7xl mx-auto">
        <!-- 项目头部信息 -->
        <div
          class="grid grid-cols-[400px,1fr] gap-8 bg-white dark:bg-zinc-800/80 rounded-xl p-6 border theme-border"
        >
          <!-- 封面 -->
          <div
            class="w-full h-[500px] rounded-xl overflow-hidden border theme-border"
          >
            <img
              :src="imagePrefix + project.cover_image"
              class="w-full h-full object-cover"
              alt="项目封面"
            />
          </div>

          <!-- 项目信息 -->
          <div class="space-y-6">
            <div>
              <h1
                class="text-4xl font-bold bg-gradient-to-r from-indigo-500 to-blue-500 bg-clip-text text-transparent"
              >
                {{ project.project_name }}
              </h1>
              <p class="text-blue-500 mt-2 font-bold flex items-center gap-2">
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
                    d="M3.75 13.5l10.5-11.25L12 10.5h8.25L9.75 21.75 12 13.5H3.75z"
                  />
                </svg>
                {{ project.types }}
              </p>
            </div>

            <!-- 数据统计 -->
            <div class="flex items-center gap-6 text-gray-500">
              <div class="flex items-center gap-2">
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
                    d="M2.036 12.322a1.012 1.012 0 010-.639C3.423 7.51 7.36 4.5 12 4.5c4.638 0 8.573 3.007 9.963 7.178.07.207.07.431 0 .639C20.577 16.49 16.64 19.5 12 19.5c-4.638 0-8.573-3.007-9.963-7.178z"
                  />
                  <path
                    stroke-linecap="round"
                    stroke-linejoin="round"
                    d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"
                  />
                </svg>
                {{ project.watches }} 浏览
              </div>
              <div class="flex items-center gap-2">
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
                    d="M21 8.25c0-2.485-2.099-4.5-4.688-4.5-1.935 0-3.597 1.126-4.312 2.733-.715-1.607-2.377-2.733-4.313-2.733C5.1 3.75 3 5.765 3 8.25c0 7.22 9 12 9 12s9-4.78 9-12z"
                  />
                </svg>
                {{ project.favorites }} 收藏
              </div>
            </div>

            <!-- 项目描述 -->
            <div class="space-y-4">
              <div>
                <h3
                  class="text-lg font-semibold text-gray-900 dark:text-gray-100 mb-2"
                >
                  剧情简介
                </h3>
                <p class="text-gray-600 dark:text-gray-400">
                  {{ project.social_story }}
                </p>
              </div>
              <div>
                <h3
                  class="text-lg font-semibold text-gray-900 dark:text-gray-100 mb-2"
                >
                  创作风格
                </h3>
                <div class="flex flex-wrap gap-2">
                  <a-tag v-for="tag in project.style" :key="tag" color="blue">{{
                    tag
                  }}</a-tag>
                </div>
              </div>
              <div>
                <h3
                  class="text-lg font-semibold text-gray-900 dark:text-gray-100 mb-2"
                >
                  适宜人群
                </h3>
                <div class="flex flex-wrap gap-2">
                  <a-tag
                    v-for="tag in project.market_people"
                    :key="tag"
                    color="green"
                    >{{ tag }}</a-tag
                  >
                </div>
              </div>
            </div>

            <!-- 团队信息 -->
            <div class="pt-4 border-t theme-border">
              <div class="flex flex-col gap-4">
                <h3
                  class="text-lg font-semibold text-gray-900 dark:text-gray-100"
                >
                  创作团队
                </h3>
                <span class="text-gray-600 dark:text-gray-400">{{
                  project.team.username
                }}
                </span>
              </div>
            </div>
          </div>
        </div>

        <!-- 在项目基本信息后添加新的内容区域 -->
        <div class="mt-8 space-y-8" v-if="project">
          <!-- 角色列表 -->
          <div
            class="bg-white dark:bg-zinc-800/80 rounded-xl p-6 border theme-border"
          >
            <h2
              class="text-2xl font-bold bg-gradient-to-r from-indigo-500 to-blue-500 bg-clip-text text-transparent mb-6"
            >
              角色列表
            </h2>
            <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
              <div
                v-for="character in characters"
                :key="character.ID"
                class="p-4 border theme-border rounded-lg hover:border-blue-500 transition-all"
              >
                <div class="flex items-center gap-4">
                  <img
                    :src="getAvatarUrl(character.avatar)"
                    class="w-12 h-12 rounded-full object-cover"
                    :alt="character.name"
                  />
                  <div>
                    <h3
                      class="text-lg font-semibold text-gray-900 dark:text-gray-100"
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

          <!-- 角色关系图 -->
          <div
            class="bg-white dark:bg-zinc-800/80 rounded-xl p-6 border theme-border"
          >
            <h2
              class="text-2xl font-bold bg-gradient-to-r from-indigo-500 to-blue-500 bg-clip-text text-transparent mb-6"
            >
              角色关系图谱
            </h2>
            <div ref="characterChart" class="w-full h-[600px]"></div>
          </div>

          <!-- 章节列表 -->
          <div class="bg-white dark:bg-zinc-800/80 rounded-xl p-6 border theme-border">
            <div class="flex items-center justify-between mb-6">
              <h2 class="text-2xl font-bold bg-gradient-to-r from-indigo-500 to-blue-500 bg-clip-text text-transparent">
                创作篇章
              </h2>
              <span class="text-sm text-gray-500" v-if="!chaptersLoading">
                共 {{ chapters.length }} 个章节
              </span>
            </div>

            <!-- 加载状态 -->
            <div v-if="chaptersLoading" class="flex items-center justify-center py-12">
              <loader />
            </div>

            <!-- 章节列表 -->
            <div v-else-if="chapters.length > 0" class="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-3 gap-4">
              <div v-for="chapter in chapters" 
                   :key="chapter.ID" 
                   :class="[
                     'group p-4 border theme-border rounded-lg hover:border-blue-500 transition-all',
                     getChapterCardStyle(chapter),
                     getChapterStatus(chapter).isPublished?'cursor-pointer':'cursor-not-allowed',
                   ]">
                <div class="flex flex-col h-full">
                  <div class="flex items-start justify-between mb-3">
                    <div class="space-y-2">
                      <div class="flex items-center gap-3">
                        <h3 class="text-lg font-semibold text-gray-900 dark:text-gray-100 group-hover:text-transparent group-hover:bg-clip-text group-hover:bg-gradient-to-r group-hover:from-blue-500 group-hover:to-cyan-500 transition-colors">
                          {{ chapter.Title }}
                        </h3>
                        <a-tag :color="getChapterStatus(chapter).color">
                          {{ getChapterStatus(chapter).text }}
                        </a-tag>
                      </div>
                      <p class="text-sm text-gray-600 dark:text-gray-400 line-clamp-2">
                        {{ getChapterStatus(chapter).isPublished?chapter.Description:'未完待续...' }}
                      </p>
                    </div>
                  </div>

                  <div class="flex-1 flex flex-col justify-end">
                    <div class="flex items-center justify-between text-sm text-gray-500 mb-3">
                      <div class="flex items-center gap-2">
                        <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                          <path stroke-linecap="round" stroke-linejoin="round" d="M12 6v6h4.5m4.5 0a9 9 0 11-18 0 9 9 0 0118 0z" />
                        </svg>
                        {{ formatDate(chapter.UpdatedAt) }}
                      </div>
                      <div v-if="chapter.CurrentVersion?.User" class="flex items-center gap-2">
                        <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                          <path stroke-linecap="round" stroke-linejoin="round" d="M15.75 6a3.75 3.75 0 11-7.5 0 3.75 3.75 0 017.5 0zM4.501 20.118a7.5 7.5 0 0114.998 0A17.933 17.933 0 0112 21.75c-2.676 0-5.216-.584-7.499-1.632z" />
                        </svg>
                        {{ chapter.CurrentVersion.User.nickname || chapter.CurrentVersion.User.username }}
                      </div>
                    </div>

                    <div v-if="chapter.CurrentVersion && chapter.CurrentVersion.ID !== 0" 
                         class="pt-3 border-t theme-border">
                      <div class="flex items-center justify-between text-sm">
                        <span class="text-gray-500">
                          当前版本字数：{{ chapter.CurrentVersion.content.length }}
                        </span>
                        <a-button type="link" size="small">
                          查看详情 →
                        </a-button>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            </div>

            <!-- 空状态 -->
            <div v-else class="flex flex-col items-center justify-center py-12 text-gray-500">
              <svg xmlns="http://www.w3.org/2000/svg" class="h-12 w-12 mb-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" d="M12 6.042A8.967 8.967 0 006 3.75c-1.052 0-2.062.18-3 .512v14.25A8.987 8.987 0 016 18c2.305 0 4.408.867 6 2.292m0-14.25a8.966 8.966 0 016-2.292c1.052 0 2.062.18 3 .512v14.25A8.987 8.987 0 0018 18a8.967 8.967 0 00-6 2.292m0-14.25v14.25" />
              </svg>
              <p>暂无篇章内容</p>
            </div>
          </div>

          <!-- 项目评论 -->
          <div class="bg-white dark:bg-zinc-800/80 rounded-xl p-6 border theme-border mt-8">
            <h2 class="text-2xl font-bold bg-gradient-to-r from-indigo-500 to-blue-500 bg-clip-text text-transparent mb-6">
              项目评论
            </h2>

            <!-- 评论输入框 -->
            <div class="mb-6">
              <div class="flex gap-4">
                <a-textarea
                  v-model:value="commentContent"
                  placeholder="分享你的想法..."
                  :rows="3"
                  :disabled="!userStore.isLogin"
                />
                <a-button
                  type="primary"
                  :disabled="!userStore.isLogin"
                  @click="submitComment"
                  class="whitespace-nowrap my-auto"
                >
                  发表评论
                </a-button>
              </div>
              <p v-if="!userStore.isLogin" class="mt-2 text-sm text-gray-500">
                请先 <a href="/auth/login" class="text-blue-500 hover:text-blue-600">登录</a> 后再评论
              </p>
            </div>

            <!-- 评论列表 -->
            <div v-if="commentsLoading" class="flex items-center justify-center py-12">
              <loader />
            </div>
            
            <div v-else-if="comments.length > 0" class="space-y-6">
              <div v-for="comment in comments" 
                   :key="comment.ID"
                   class="group p-4 border theme-border rounded-lg bg-gray-50/50 dark:bg-zinc-900/80">
                <div class="flex items-start gap-4">
                  <img 
                    :src="comment.user.avatar ? BACKEND_DOMAIN + comment.user.avatar : '/default-avatar.png'"
                    class="w-10 h-10 rounded-full object-cover"
                    :alt="comment.user.nickname || comment.user.username"
                  />
                  <div class="flex-1">
                    <div class="flex items-center justify-between mb-2">
                      <span class="font-medium text-gray-900 dark:text-gray-100">
                        {{ comment.user.nickname || comment.user.username }}
                      </span>
                      <span class="text-sm text-gray-500">
                        {{ formatCommentDate(comment.CreatedAt) }}
                      </span>
                    </div>
                    <p class="text-gray-700 dark:text-gray-300">
                      {{ comment.Content }}
                    </p>
                  </div>
                </div>
              </div>
            </div>

            <div v-else class="flex flex-col items-center justify-center py-12 text-gray-500">
              <svg xmlns="http://www.w3.org/2000/svg" class="h-12 w-12 mb-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z" />
              </svg>
              <p>暂无评论，来说两句吧~</p>
            </div>
          </div>
        </div>
      </div>
    </template>

    <!-- 错误状态 -->
    <div
      v-else
      class="flex items-center justify-center min-h-[400px] text-gray-500"
    >
      项目不存在或已被删除
    </div>
  </div>
</template>
