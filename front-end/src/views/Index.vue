<script setup>
import ThemeSwitcher from "@/components/button/ThemeSwitcher.vue";
import logo from "@/assets/img/logo.webp";
import router from "@/router/index.js";
import { useUserStore } from "@/stores/user.js";
import { get } from "@/util/request.js";
import { useThemeStore } from "@/stores/theme.js";
import { message } from "ant-design-vue";
import { useFullPageScroll } from "@/util/useFullPageScroll.js";
import * as echarts from "echarts";
import { ref, onMounted, onUnmounted, watch } from "vue";
import titleDark from "@/assets/img/title-dark.png";
import { MdPreview } from "md-editor-v3";
import "md-editor-v3/lib/preview.css";
import webs from "@/assets/img/webs.webp";
import { useMotion } from "@vueuse/motion";
import ai from '@/assets/img/ai-active.png'
const { currentSection } = useFullPageScroll();
// 将渐变色数组提取为公共变量
const gradientColors = [
  ["#3B82F6", "#60A5FA"], // 蓝色
  ["#10B981", "#34D399"], // 绿色
  ["#F59E0B", "#FBBF24"], // 橙色
  ["#8B5CF6", "#A78BFA"], // 紫色
  ["#EC4899", "#F472B6"], // 粉色
  ["#EF4444", "#F87171"], // 红色
  ["#06B6D4", "#22D3EE"], // 青色
  ["#14B8A6", "#2DD4BF"], // 蓝绿色
  ["#F97316", "#FB923C"], // 深橙色
  ["#84CC16", "#A3E635"], // 青柠色
  ["#6366F1", "#818CF8"], // 靛蓝色
  ["#D946EF", "#E879F9"], // 洋红色
  ["#0EA5E9", "#38BDF8"], // 天蓝色
  ["#4F46E5", "#6366F1"], // 深蓝色
  ["#7C3AED", "#A78BFA"], // 深紫色
  ["#DB2777", "#EC4899"], // 玫红色
  ["#EA580C", "#FB923C"], // 赭石色
  ["#16A34A", "#4ADE80"], // 翠绿色
  ["#2563EB", "#60A5FA"], // 皇家蓝
  ["#9333EA", "#C084FC"], // 紫罗兰色
];

const themeStore = useThemeStore();
const getFrontColor = () => {
  return themeStore.isDark ? "#FFFFFF" : "#000000";
};
const topBarItems = [
  {
    name: "作品画廊",
    path: "/",
    icon: `<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="size-5">
  <path stroke-linecap="round" stroke-linejoin="round" d="M2.25 7.125C2.25 6.504 2.754 6 3.375 6h6c.621 0 1.125.504 1.125 1.125v3.75c0 .621-.504 1.125-1.125 1.125h-6a1.125 1.125 0 0 1-1.125-1.125v-3.75ZM14.25 8.625c0-.621.504-1.125 1.125-1.125h5.25c.621 0 1.125.504 1.125 1.125v8.25c0 .621-.504 1.125-1.125 1.125h-5.25a1.125 1.125 0 0 1-1.125-1.125v-8.25ZM3.75 16.125c0-.621.504-1.125 1.125-1.125h5.25c.621 0 1.125.504 1.125 1.125v2.25c0 .621-.504 1.125-1.125 1.125h-5.25a1.125 1.125 0 0 1-1.125-1.125v-2.25Z" />
</svg>
`,
  },
  {
    name: "订阅计划",
    path: "/",
    icon: `<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="size-5">
  <path stroke-linecap="round" stroke-linejoin="round" d="M15.75 5.25a3 3 0 0 1 3 3m3 0a6 6 0 0 1-7.029 5.912c-.563-.097-1.159.026-1.563.43L10.5 17.25H8.25v2.25H6v2.25H2.25v-2.818c0-.597.237-1.17.659-1.591l6.499-6.499c.404-.404.527-1 .43-1.563A6 6 0 1 1 21.75 8.25Z" />
</svg>
`,
  },
];
const quickStartClickHandler = () => {
  const userStore = useUserStore();
  if (userStore.isLogin) {
    console.log(userStore.user);
    message.success("欢迎回来," + userStore.user.username + "!正在为您跳转...");
    setTimeout(() => {
      router.push("/workspace/dataAnlysis");
    }, 1000);
  } else {
    get(
      "/api/user/me",
      {},
      (messager, data) => {
        message.success("欢迎回来," + data.username + "!正在为您跳转...");
        userStore.login(data);
        setTimeout(() => {
          router.push("/workspace/dataAnlysis");
        }, 1000);
      },
      (messager, data) => {
        message.info(messager);
        console.log(data);
        setTimeout(() => {
          router.push("/auth/login");
        }, 1000);
      },
      (messager, data) => {
        console.log(messager);
        message.info(messager);
        setTimeout(() => {
          router.push("/auth/login");
        }, 1000);
      }
    );
  }
};

const characterChart = ref(null);

const initCharacterChart = () => {
  if (!characterChart.value) return;

  const chartInstance = echarts.init(characterChart.value);

  // 扩展角色数据
  const characters = [
    {
      id: 1,
      name: "艾琳",
      value: "女主角，天才黑客，性格独立坚强",
      category: 0,
    },
    { id: 2, name: "马克", value: "男主角，前特工，充满正义感", category: 0 },
    {
      id: 3,
      name: "萨姆",
      value: "艾琳的导师，神秘的科技公司创始人",
      category: 2,
    },
    {
      id: 4,
      name: "莉莉",
      value: "艾琳的闺蜜，记者，擅长收集情报",
      category: 1,
    },
    { id: 5, name: "杰克", value: "马克的搭档，技术专家", category: 1 },
    { id: 6, name: "维克多", value: "主要反派，黑市军火商", category: 3 },
    { id: 7, name: "安娜", value: "萨姆的助手，AI专家", category: 1 },
    { id: 8, name: "克里斯", value: "警局内部线人", category: 2 },
    { id: 9, name: "艾米", value: "维克多的妹妹，双面间谍", category: 3 },
    { id: 10, name: "大卫", value: "政府特工，身份成谜", category: 2 },
    { id: 11, name: "索菲亚", value: "黑客组织领袖", category: 3 },
    { id: 12, name: "托马斯", value: "科技公司CEO，幕后推手", category: 2 },
  ];

  // 扩展关系数据
  const relationships = [
    {
      source: 0,
      target: 1,
      name: "命运相遇",
      value: "在一次网络安全事件中相识，共同对抗犯罪组织",
    },
    {
      source: 0,
      target: 2,
      name: "师徒关系",
      value: "萨姆发现艾琳的天赋，收她为徒",
    },
    {
      source: 0,
      target: 3,
      name: "挚友",
      value: "从大学时代就是形影不离的好友",
    },
    { source: 1, target: 4, name: "搭档", value: "多年的工作伙伴，配合默契" },
    {
      source: 2,
      target: 5,
      name: "宿敌",
      value: "曾是合作伙伴，后因理念不合成为对手",
    },
    { source: 3, target: 4, name: "信息合作", value: "莉莉为杰克提供重要情报" },
    {
      source: 1,
      target: 2,
      name: "互相怀疑",
      value: "对萨姆的真实身份产生怀疑",
    },
    { source: 2, target: 6, name: "秘密合作", value: "安娜协助萨姆进行AI研究" },
    {
      source: 5,
      target: 7,
      name: "非法交易",
      value: "克里斯暗中监视维克多的军火交易",
    },
    {
      source: 5,
      target: 8,
      name: "血亲",
      value: "艾米是维克多的妹妹，但暗中为政府工作",
    },
    { source: 8, target: 9, name: "上下级", value: "艾米向大卫汇报情报" },
    {
      source: 0,
      target: 10,
      name: "对手",
      value: "艾琳与索菲亚在多个案件中对抗",
    },
    {
      source: 2,
      target: 11,
      name: "商业合作",
      value: "萨姆与托马斯有密切的商业往来",
    },
    {
      source: 10,
      target: 11,
      name: "秘密联盟",
      value: "索菲亚与托马斯暗中结盟",
    },
    {
      source: 6,
      target: 11,
      name: "技术支持",
      value: "安娜为托马斯提供AI技术支持",
    },
    {
      source: 3,
      target: 7,
      name: "情报交换",
      value: "莉莉与克里斯交换警方情报",
    },
    { source: 4, target: 9, name: "暗中调查", value: "杰克调查大卫的真实身份" },
  ];

  const nodes = characters.map((character, index) => ({
    id: index,
    name: character.name,
    value: character.value,
    symbolSize: 60, // 增大节点大小
    category: character.category,
  }));

  const links = relationships.map((relation) => ({
    source: relation.source,
    target: relation.target,
    name: relation.name,
    value: relation.value,
  }));

  const option = {
    title: {
      text: "角色关系图谱",
      subtext: "示例项目：赛博朋克悬疑剧",
      top: "bottom",
      left: "right",
      textStyle: {
        color: themeStore.isDark ? "#e5e7eb" : "#111827",
        fontSize: 16,
        fontWeight: "bold",
      },
      subtextStyle: {
        color: themeStore.isDark ? "#9ca3af" : "#4b5563",
      },
    },
    tooltip: {
      show: true,
      confine: true,
      formatter: function (params) {
        const title = params.dataType === "edge" ? "关系" : "角色";
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
              ${title}：${params.data.name}
            </div>
            <div style="
              color: rgba(255,255,255,0.9);
              white-space: normal;
              word-break: break-all;
            ">
              描述：${params.data.value}
            </div>
          </div>
        `;
      },
      backgroundColor: "rgba(0,0,0,0.75)",
      borderRadius: 8,
      padding: [10, 15],
      textStyle: {
        color: "#fff",
        fontSize: 14,
      },
    },
    legend: [
      {
        data: ["主角", "重要角色", "神秘人物", "反派"],
        orient: "vertical",
        left: "left",
        textStyle: {
          color: themeStore.isDark ? "#e5e7eb" : "#111827",
        },
      },
    ],
    series: [
      {
        type: "graph",
        layout: "force",
        data: nodes,
        links: links,
        categories: [
          { name: "主角" },
          { name: "重要角色" },
          { name: "神秘人物" },
          { name: "反派" },
        ],
        roam: true,
        draggable: true,
        label: {
          show: true,
          position: "right",
          formatter: "{b}",
          backgroundColor: "rgba(102, 102, 102, 0.75)",
          borderRadius: 4,
          padding: [4, 8],
          color: "#fff",
          fontSize: 14,
          fontWeight: "bold",
        },
        edgeLabel: {
          show: true,
          formatter: "{c}",
          backgroundColor: "rgba(0, 0, 0, 0.75)",
          borderRadius: 4,
          padding: [4, 8],
          color: "#fff",
          fontSize: 12,
          distance: 20,
          rotate: 0,
          offset: [0, 0],
        },
        force: {
          repulsion: 800,
          gravity: 0.1,
          edgeLength: 300,
          friction: 0.1,
          layoutAnimation: true,
          initLayout: "circular",
        },
        cursor: "move",
        zoom: 1.5,
        center: ["50%", "50%"],
        emphasis: {
          focus: "adjacency",
          scale: 1.2,
          lineStyle: {
            width: 4,
          },
          label: {
            fontSize: 16,
            show: true,
          },
        },
        nodeScaleRatio: 0.6,
        itemStyle: {
          borderWidth: 2,
          borderColor: themeStore.isDark
            ? "rgba(255,255,255,0.3)"
            : "rgba(0,0,0,0.2)",
          shadowColor: "rgba(0, 0, 0, 0.2)",
          shadowBlur: 10,
        },
        animation: true,
        animationDuration: 1000,
        animationEasingUpdate: "quinticInOut",
        lineStyle: {
          width: 2,
          opacity: 0.7,
          curveness: 0.3,
          color: {
            type: "linear",
            x: 0,
            y: 0,
            x2: 0,
            y2: 1,
            colorStops: [
              {
                offset: 0,
                color: "rgb(77, 77, 245)",
              },
              {
                offset: 1,
                color: "rgb(77, 77, 245)",
              },
            ],
          },
        },
        emphasis: {
          lineStyle: {
            width: 4,
            curveness: 0.3,
          },
        },
        edgeSymbol: ["none", "arrow"],
        edgeSymbolSize: [0, 8],
      },
    ],
    backgroundColor: "transparent",
  };

  chartInstance.setOption(option);

  // 添加鼠标交互事件
  chartInstance.on("mousedown", "series", () => {
    chartInstance.getZr().setCursorStyle("move");
  });

  chartInstance.on("mouseup", "series", () => {
    chartInstance.getZr().setCursorStyle("default");
  });

  // 优化 resize 处理
  const resizeHandler = () => {
    chartInstance.resize();
  };

  window.addEventListener("resize", resizeHandler);

  // 确保组件卸载时移除事件监听
  onUnmounted(() => {
    window.removeEventListener("resize", resizeHandler);
    chartInstance.dispose();
  });
};

// 监听主题变化，更新图表样式
watch(
  () => themeStore.isDark,
  () => {
    initCharacterChart();
  }
);

onMounted(() => {
  initCharacterChart();
  //message.info("本网站处于开发状态,为学科竞赛参赛作品部署，请不要封锁");
});

// 示例数据
const demoData = {
  dates: ["2024-03-01", "2024-03-02", "2024-03-03", "2024-03-04", "2024-03-05"],
  projects: ["科幻冒险", "都市情感", "悬疑推理"],
  emotions: ["喜悦", "感动", "惊喜", "期待", "伤感"],
};

// 初始化示例图表
const initDemoCharts = () => {
  // 观看趋势图
  const watchChart = echarts.init(document.querySelector("#watchChart"));
  watchChart.setOption({
    tooltip: {
      trigger: "axis",
    },
    legend: {
      data: demoData.projects,
      bottom: 0,
    },
    xAxis: {
      type: "category",
      data: demoData.dates,
      axisLabel: { rotate: 45 },
    },
    yAxis: {
      type: "value",
      name: "观看次数",
    },
    series: demoData.projects.map((project, index) => ({
      name: project,
      type: "line",
      smooth: true,
      data: demoData.dates.map(() => Math.floor(Math.random() * 1000)),
      itemStyle: {
        color: gradientColors[index][0],
      },
    })),
  });

  // 情感分析图
  const emotionChart = echarts.init(document.querySelector("#emotionChart"));
  emotionChart.setOption({
    tooltip: {
      trigger: "item",
    },
    legend: {
      bottom: 0,
    },
    series: [
      {
        type: "pie",
        radius: ["40%", "70%"],
        data: demoData.emotions.map((emotion, index) => ({
          name: emotion,
          value: Math.floor(Math.random() * 100),
          itemStyle: {
            color: gradientColors[index][0],
          },
        })),
      },
    ],
  });

  // 监听窗口大小变化
  window.addEventListener("resize", () => {
    watchChart.resize();
    emotionChart.resize();
    styleChart.resize();
    aiChart.resize();
  });
};

onMounted(() => {
  initDemoCharts();
});

const demoText = ref("");
const fullDemoText = `正在分析已有内容...

故事发生在一个雨夜，城市的霓虹灯在雨中显得格外迷离。艾琳坐在她的工作室里，双手飞快地在键盘上敲击。显示器的蓝光映照在她专注的脸上，一行行代码如流水般在屏幕上滑过。

突然，她的系统发出了警告。有人正试图入侵她精心设计的防火墙。这不是普通的黑客，他们的手法太专业了。艾琳眯起眼睛，嘴角露出一丝微笑。"有意思，"她喃喃自语，"让我们看看你是谁。"

就在这时，她的通讯器响了。是一个加密频道。"艾琳，"一个熟悉的声音传来，"是我，马克。我们需要谈谈。"`;

let currentIndex = ref(0);
const isTyping = ref(true);

const typeText = () => {
  if (currentIndex.value < fullDemoText.length) {
    demoText.value = fullDemoText.slice(0, currentIndex.value + 1);
    currentIndex.value++;
  } else {
    // 重置动画
    setTimeout(() => {
      currentIndex.value = 0;
      demoText.value = "";
      isTyping.value = true;
    }, 2000);
  }
};

onMounted(() => {
  const interval = setInterval(() => {
    if (isTyping.value) {
      typeText();
    }
  }, 50);

  onUnmounted(() => {
    clearInterval(interval);
  });
});
</script>

<template>
  <div
    class="fixed top-0 left-0 w-full shadow-md z-50 flex flex-nowrap place-items-center text-theme-switch font-sans border-b-[1px] theme-border bkg-theme-switch"
  >
    <div class="w-1/12"></div>
    <div class="flex flex-nowrap p-4 cursor-pointer" @click="router.push('/')">
      <img :src="logo" class="w-8 h-8 rounded-full" />
      <span class="my-auto ml-4">创剧星球</span>
    </div>
    <div class="flex-grow" />
    <div class="h-full flex gap-4">
      <div
        class="text-hover flex flex-nowrap cursor-pointer transition-all my-auto"
        v-for="item in topBarItems"
        @click="router.push(item.path)"
      >
        <div v-html="item.icon" class="h-6 w-6 rounded-full" />
        <span>{{ item.name }}</span>
      </div>
      <ThemeSwitcher class="mx-auto my-auto" />
    </div>
    <div class="w-1/12"></div>
  </div>
  <div class="relative w-full h-screen overflow-hidden bg-black">
    <div class="absolute inset-0 overflow-hidden bg-slate-50 dark:bg-slate-950">
      <!-- 核心发光中心 -->
      <div
        class="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 w-32 h-32"
      >
        <div
          class="absolute inset-0 bg-gradient-to-r from-blue-400/50 via-violet-400/40 to-transparent dark:from-blue-400/30 dark:via-violet-400/25 blur-md"
        ></div>
      </div>

      <!-- 主渐变区域 -->
      <div
        class="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 w-[600px] h-[600px]"
      >
        <div
          class="absolute inset-0 bg-gradient-to-r from-blue-500/10 via-violet-500/40 to-transparent dark:from-blue-600/40 dark:via-violet-600/35 dark:to-slate-950/30 blur-3xl"
        ></div>
      </div>

      <!-- 线性光线 -->
      <div class="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2">
        <template v-for="i in 12" :key="i">
          <div
            class="absolute h-px w-[400px] bg-gradient-to-r from-blue-400/20 via-violet-400/15 to-transparent dark:from-blue-400/15 dark:via-violet-400/10 dark:to-slate-950/5"
            :style="{
              transform: `rotate(${(i - 1) * 30}deg)`,
              transformOrigin: '0 50%',
            }"
          ></div>
        </template>
      </div>

      <!-- 边缘渐隐效果 -->
      <div
        class="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 w-[800px] h-[800px]"
      >
        <div
          class="absolute inset-0 bg-gradient-to-r from-blue-400/20 via-violet-400/15 to-transparent dark:from-blue-500/15 dark:via-violet-500/10 dark:to-slate-950/5 blur-[150px]"
        ></div>
      </div>
    </div>

    <!-- 内容层 -->
    <div
      class="relative z-10 flex flex-col items-center justify-center h-full font-sans"
      v-motion
      :initial="{ opacity: 0, y: -50, scale: 0.9 }"
      :enter="{ 
        opacity: 1, 
        y: 0, 
        scale: 1,
        transition: {
          type: 'spring',
          damping: 12,
          stiffness: 100,
          duration: 800
        }
      }"
    >
      <div class="mx-auto font-serif"
           v-motion
           :initial="{ opacity: 0, y: 20 }"
           :enter="{ 
             opacity: 1, 
             y: 0,
             transition: { 
               delay: 200,
               duration: 500
             }
           }">
        <img :src="titleDark" class="invert dark:invert-0" />
      </div>

      <div class="mx-auto font-serif"></div>
      
      <div
        class="mx-auto font-serif text-gray-600 dark:text-[#8B8B8B] mt-8 text-2xl"
        v-motion
        :initial="{ opacity: 0, y: 20 }"
        :enter="{ 
          opacity: 1, 
          y: 0,
          transition: { 
            delay: 400,
            duration: 500
          }
        }"
      >
        让创意化为精彩剧集
      </div>
      
      <div class="mx-auto text-gray-600 dark:text-[#737373] mt-8 text-xl"
           v-motion
           :initial="{ opacity: 0, y: 20 }"
           :enter="{ 
             opacity: 1, 
             y: 0,
             transition: { 
               delay: 600,
               duration: 500
             }
           }">
        一款集AI智能创作、角色塑造、剧情构建、团队协同于一体的专业创作平台，
      </div>

      <div class="mx-auto text-gray-500 dark:text-[#666666] mt-2"
           v-motion
           :initial="{ opacity: 0, y: 20 }"
           :enter="{ 
             opacity: 1, 
             y: 0,
             transition: { 
               delay: 800,
               duration: 500
             }
           }">
        为您提供从灵感激发到作品完成的全流程支持，让创作更轻松，让故事更精彩。
      </div>

      <div class="mx-auto mt-4 flex flex-wrap gap-2"
           v-motion
           :initial="{ opacity: 0, y: 20 }"
           :enter="{ 
             opacity: 1, 
             y: 0,
             transition: { 
               delay: 1000,
               duration: 500
             }
           }">
        <button
          @click="quickStartClickHandler"
          class="px-6 py-2 rounded-2xl font-light bg-gradient-to-r from-blue-500 to-violet-500 text-white hover:-translate-y-0.5 transition-all"
          v-motion
          :hover="{ scale: 1.05, transition: { duration: 200 } }"
        >
          创剧工坊
        </button>
        <button
          @click="router.push('/community')"
          class="px-6 py-2 rounded-2xl font-light border border-gray-300 dark:border-gray-700 text-gray-700 dark:text-white hover:-translate-y-0.5 transition-all"
          v-motion
          :hover="{ scale: 1.05, transition: { duration: 200 } }"
        >
          星球社区
        </button>
      </div>
    </div>
  </div>
  <div class="h-screen bg-[#f8fafc] dark:bg-[#030616]">
    <div class="container mx-auto px-4 py-16 font-sans">
      <!-- AI扩写功能介绍 -->
      <div class="flex flex-col lg:flex-row gap-8 items-center">
        <!-- 左侧介绍文本 -->
        <div class="lg:w-1/2 space-y-6">
          <h2 class="text-4xl font-bold font-serif text-black dark:text-white">
            AI 智能创作
          </h2>
          <h3 class="text-2xl font-bold text-gray-700 dark:text-gray-300">
            让创作过程更轻松，让灵感源源不断
          </h3>
          <div class="space-y-4">
            <div class="flex items-start gap-4">
              <div class="p-2 bg-blue-100 dark:bg-blue-900/30 rounded-lg mt-1">
                <svg
                  xmlns="http://www.w3.org/2000/svg"
                  class="h-6 w-6 text-blue-500"
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
              </div>
              <div>
                <h4
                  class="text-lg font-semibold text-gray-900 dark:text-gray-100"
                >
                  智能续写
                </h4>
                <p class="text-gray-600 dark:text-gray-400">
                  根据已有内容智能分析，自动生成后续情节，让故事自然流畅地发展
                </p>
              </div>
            </div>

            <div class="flex items-start gap-4">
              <div
                class="p-2 bg-purple-100 dark:bg-purple-900/30 rounded-lg mt-1"
              >
                <svg
                  xmlns="http://www.w3.org/2000/svg"
                  class="h-6 w-6 text-purple-500"
                  fill="none"
                  viewBox="0 0 24 24"
                  stroke="currentColor"
                >
                  <path
                    stroke-linecap="round"
                    stroke-linejoin="round"
                    d="M15.042 21.672L13.684 16.6m0 0l-2.51 2.225.569-9.47 5.227 7.917-3.286-.672zM12 2.25V4.5m5.834.166l-1.591 1.591M20.25 10.5H18M7.757 14.743l-1.59 1.59M6 10.5H3.75m4.007-4.243l-1.59-1.59"
                  />
                </svg>
              </div>
              <div>
                <h4
                  class="text-lg font-semibold text-gray-900 dark:text-gray-100"
                >
                  情节优化
                </h4>
                <p class="text-gray-600 dark:text-gray-400">
                  提供情节建议和优化方案，让故事更加引人入胜
                </p>
              </div>
            </div>

            <div class="flex items-start gap-4">
              <div
                class="p-2 bg-green-100 dark:bg-green-900/30 rounded-lg mt-1"
              >
                <svg
                  xmlns="http://www.w3.org/2000/svg"
                  class="h-6 w-6 text-green-500"
                  fill="none"
                  viewBox="0 0 24 24"
                  stroke="currentColor"
                >
                  <path
                    stroke-linecap="round"
                    stroke-linejoin="round"
                    d="M15.75 6a3.75 3.75 0 11-7.5 0 3.75 3.75 0 017.5 0zM4.501 20.118a7.5 7.5 0 0114.998 0A17.933 17.933 0 0112 21.75c-2.676 0-5.216-.584-7.499-1.632z"
                  />
                </svg>
              </div>
              <div>
                <h4
                  class="text-lg font-semibold text-gray-900 dark:text-gray-100"
                >
                  角色塑造
                </h4>
                <p class="text-gray-600 dark:text-gray-400">
                  智能分析角色性格特征，提供人物对话和行为建议
                </p>
              </div>
            </div>

            <div class="flex items-start gap-4">
              <div class="p-2 bg-rose-100 dark:bg-rose-900/30 rounded-lg mt-1">
                <svg
                  xmlns="http://www.w3.org/2000/svg"
                  class="h-6 w-6 text-rose-500"
                  fill="none"
                  viewBox="0 0 24 24"
                  stroke="currentColor"
                >
                  <path
                    stroke-linecap="round"
                    stroke-linejoin="round"
                    d="M2.25 15.75l5.159-5.159a2.25 2.25 0 013.182 0l5.159 5.159m-1.5-1.5l1.409-1.409a2.25 2.25 0 013.182 0l2.909 2.909m-18 3.75h16.5a1.5 1.5 0 001.5-1.5V6a1.5 1.5 0 00-1.5-1.5H3.75A1.5 1.5 0 002.25 6v12a1.5 1.5 0 001.5 1.5zm10.5-11.25h.008v.008h-.008V8.25zm.375 0a.375.375 0 11-.75 0 .375.375 0 01.75 0z"
                  />
                </svg>
              </div>
              <div>
                <h4
                  class="text-lg font-semibold text-gray-900 dark:text-gray-100"
                >
                  角色画像生成
                </h4>
                <p class="text-gray-600 dark:text-gray-400">
                  基于角色设定智能生成个性化角色画像，让角色形象更加丰富立体
                </p>
              </div>
            </div>

            <div class="flex items-start gap-4">
              <div
                class="p-2 bg-amber-100 dark:bg-amber-900/30 rounded-lg mt-1"
              >
                <svg
                  xmlns="http://www.w3.org/2000/svg"
                  class="h-6 w-6 text-amber-500"
                  fill="none"
                  viewBox="0 0 24 24"
                  stroke="currentColor"
                >
                  <path
                    stroke-linecap="round"
                    stroke-linejoin="round"
                    d="M6.429 9.75L2.25 12l4.179 2.25m0-4.5l5.571 3 5.571-3m-11.142 0L2.25 7.5 12 2.25l9.75 5.25-4.179 2.25m0 0L21.75 12l-4.179 2.25m0 0l4.179 2.25L12 21.75 2.25 16.5l4.179-2.25m11.142 0l-5.571 3-5.571-3"
                  />
                </svg>
              </div>
              <div>
                <h4
                  class="text-lg font-semibold text-gray-900 dark:text-gray-100"
                >
                  作品海报设计
                </h4>
                <p class="text-gray-600 dark:text-gray-400">
                  根据作品主题和风格自动生成精美海报，提升作品视觉吸引力
                </p>
              </div>
            </div>
          </div>
        </div>

        <!-- 右侧演示窗口 -->
        <div class="lg:w-1/2">
          <div
            class="bg-white dark:bg-zinc-900 rounded-xl shadow-xl border theme-border p-6"
          >
            <div class="flex items-center gap-3 mb-4">
              <div class="flex space-x-2">
                <div class="w-3 h-3 rounded-full bg-red-500"></div>
                <div class="w-3 h-3 rounded-full bg-yellow-500"></div>
                <div class="w-3 h-3 rounded-full bg-green-500"></div>
              </div>
              <div
                class="flex-1 text-center text-sm text-gray-600 dark:text-gray-400"
              >
                AI 智能创作示例
              </div>
            </div>

            <div
              class="h-[400px] overflow-hidden rounded-lg bg-gray-50 dark:bg-zinc-800 p-4"
            >
              <MdPreview
                :modelValue="demoText"
                :theme="themeStore.currentTheme"
                previewTheme="github"
                style="background: transparent"
              />
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
  <div class="h-screen bg-[#f8fafc] dark:bg-[#030616]">
    <div class="container mx-auto px-4 py-16 h-full font-sans flex flex-col">
      <!-- 标题部分 -->
      <div class="flex-none">
        <h1
          class="text-5xl font-bold font-serif text-black dark:text-white text-center mb-4"
        >
          数据驱动创作
        </h1>
        <h1
          class="text-3xl font-bold font-serif text-black dark:text-white text-center mb-8"
        >
          让每个决策都有据可依
        </h1>
      </div>

      <!-- 图表展示区域 -->
      <div class="grid grid-cols-1 lg:grid-cols-2 gap-8 flex-1">
        <!-- 观看趋势示例图 -->
        <div
          class="bg-white dark:bg-[#141419] rounded-xl p-6 border theme-border hover:shadow-lg transition-shadow"
        >
          <div class="flex items-center gap-3 mb-4">
            <div class="p-2 bg-blue-100 dark:bg-blue-900/30 rounded-lg">
              <svg
                xmlns="http://www.w3.org/2000/svg"
                class="h-6 w-6 text-blue-500"
                fill="none"
                viewBox="0 0 24 24"
                stroke="currentColor"
              >
                <path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"
                />
                <path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"
                />
              </svg>
            </div>
            <h3 class="text-lg font-semibold text-gray-900 dark:text-gray-100">
              观看趋势分析
            </h3>
          </div>
          <div id="watchChart" ref="watchChart" class="w-full h-[300px]"></div>
        </div>

        <!-- 情感分析示例图 -->
        <div
          class="bg-white dark:bg-[#141419] rounded-xl p-6 border theme-border hover:shadow-lg transition-shadow"
        >
          <div class="flex items-center gap-3 mb-4">
            <div class="p-2 bg-purple-100 dark:bg-purple-900/30 rounded-lg">
              <svg
                xmlns="http://www.w3.org/2000/svg"
                class="h-6 w-6 text-purple-500"
                fill="none"
                viewBox="0 0 24 24"
                stroke="currentColor"
              >
                <path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  d="M14.828 14.828a4 4 0 01-5.656 0M9 10h.01M15 10h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
                />
              </svg>
            </div>
            <h3 class="text-lg font-semibold text-gray-900 dark:text-gray-100">
              情感分析洞察
            </h3>
          </div>
          <div
            id="emotionChart"
            ref="emotionChart"
            class="w-full h-[300px]"
          ></div>
        </div>
      </div>
    </div>
  </div>
  <div class="h-screen bg-[#f8fafc] dark:bg-[#030616]">
    <div class="container mx-auto px-4 py-16 h-full font-sans flex flex-col">
      <div class="flex-none">
        <h1
          class="text-5xl font-bold font-serif text-black dark:text-white text-center mb-4"
        >
          角色关系，一目了然
        </h1>
        <h1
          class="text-3xl font-bold font-serif text-black dark:text-white text-center mb-8"
        >
          我们提供最清晰的方式助你理解项目。
        </h1>
      </div>

      <div
        class="flex-1 min-h-0 bg-white dark:bg-zinc-900/80 rounded-xl shadow-xl p-6"
      >
        <div ref="characterChart" class="w-full h-full" />
      </div>
    </div>
  </div>
  <div class="h-screen bg-[#f8fafc] dark:bg-[#030616]">
    <div class="container mx-auto px-4 py-16 font-sans">
      <!-- 标题部分 -->
      <div class="text-center mb-16">
        <h1
          class="text-5xl font-bold font-serif text-black dark:text-white mb-4"
        >
          多样化导出，灵活应用
        </h1>
        <h2
          class="text-3xl font-bold font-serif text-black dark:text-white mb-8"
        >
          为不同场景提供专业的输出方案
        </h2>
      </div>

      <!-- 导出功能展示网格 -->
      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-8">
        <!-- PDF导出卡片 -->
        <div
          class="bg-white dark:bg-[#141419] rounded-xl p-8 border theme-border hover:shadow-xl transition-all group"
        >
          <div class="flex items-center gap-4 mb-6">
            <div
              class="p-3 rounded-xl bg-red-100 dark:bg-red-900/30 group-hover:scale-110 transition-transform"
            >
              <svg
                xmlns="http://www.w3.org/2000/svg"
                class="h-8 w-8 text-red-600"
                fill="none"
                viewBox="0 0 24 24"
                stroke="currentColor"
              >
                <path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  d="M19.5 14.25v-2.625a3.375 3.375 0 00-3.375-3.375h-1.5A1.125 1.125 0 0113.5 7.125v-1.5a3.375 3.375 0 00-3.375-3.375H8.25m2.25 0H5.625c-.621 0-1.125.504-1.125 1.125v17.25c0 .621.504 1.125 1.125 1.125h12.75c.621 0 1.125-.504 1.125-1.125V11.25a9 9 0 00-9-9z"
                />
              </svg>
            </div>
            <h3 class="text-xl font-bold text-gray-900 dark:text-gray-100">
              PDF格式
            </h3>
          </div>
          <p class="text-gray-600 dark:text-gray-400 mb-6">
            专业的排版布局，支持富文本样式，适合正式场合使用。完美支持跨平台阅读，保持内容的一致性展现。
          </p>
          <ul class="space-y-3 text-sm">
            <li
              class="flex items-center gap-2 text-gray-700 dark:text-gray-300"
            >
              <svg
                xmlns="http://www.w3.org/2000/svg"
                class="h-5 w-5 text-green-500"
                fill="none"
                viewBox="0 0 24 24"
                stroke="currentColor"
              >
                <path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  d="M9 12.75L11.25 15 15 9.75M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
                />
              </svg>
              专业排版与样式
            </li>
            <li
              class="flex items-center gap-2 text-gray-700 dark:text-gray-300"
            >
              <svg
                xmlns="http://www.w3.org/2000/svg"
                class="h-5 w-5 text-green-500"
                fill="none"
                viewBox="0 0 24 24"
                stroke="currentColor"
              >
                <path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  d="M9 12.75L11.25 15 15 9.75M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
                />
              </svg>
              跨平台阅读支持
            </li>
            <li
              class="flex items-center gap-2 text-gray-700 dark:text-gray-300"
            >
              <svg
                xmlns="http://www.w3.org/2000/svg"
                class="h-5 w-5 text-green-500"
                fill="none"
                viewBox="0 0 24 24"
                stroke="currentColor"
              >
                <path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  d="M9 12.75L11.25 15 15 9.75M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
                />
              </svg>
              内容版式锁定
            </li>
          </ul>
        </div>

        <!-- Word导出卡片 -->
        <div
          class="bg-white dark:bg-[#141419] rounded-xl p-8 border theme-border hover:shadow-xl transition-all group"
        >
          <div class="flex items-center gap-4 mb-6">
            <div
              class="p-3 rounded-xl bg-blue-100 dark:bg-blue-900/30 group-hover:scale-110 transition-transform"
            >
              <svg
                xmlns="http://www.w3.org/2000/svg"
                class="h-8 w-8 text-blue-600"
                fill="none"
                viewBox="0 0 24 24"
                stroke="currentColor"
              >
                <path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  d="M19.5 14.25v-2.625a3.375 3.375 0 00-3.375-3.375h-1.5A1.125 1.125 0 0113.5 7.125v-1.5a3.375 3.375 0 00-3.375-3.375H8.25m0 12.75h7.5m-7.5 3H12M10.5 2.25H5.625c-.621 0-1.125.504-1.125 1.125v17.25c0 .621.504 1.125 1.125 1.125h12.75c.621 0 1.125-.504 1.125-1.125V11.25a9 9 0 00-9-9z"
                />
              </svg>
            </div>
            <h3 class="text-xl font-bold text-gray-900 dark:text-gray-100">
              Word格式
            </h3>
          </div>
          <p class="text-gray-600 dark:text-gray-400 mb-6">
            支持二次编辑，方便修改和协作。兼容主流办公软件，适合需要进一步加工的场景。
          </p>
          <ul class="space-y-3 text-sm">
            <li
              class="flex items-center gap-2 text-gray-700 dark:text-gray-300"
            >
              <svg
                xmlns="http://www.w3.org/2000/svg"
                class="h-5 w-5 text-green-500"
                fill="none"
                viewBox="0 0 24 24"
                stroke="currentColor"
              >
                <path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  d="M9 12.75L11.25 15 15 9.75M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
                />
              </svg>
              可编辑与修改
            </li>
            <li
              class="flex items-center gap-2 text-gray-700 dark:text-gray-300"
            >
              <svg
                xmlns="http://www.w3.org/2000/svg"
                class="h-5 w-5 text-green-500"
                fill="none"
                viewBox="0 0 24 24"
                stroke="currentColor"
              >
                <path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  d="M9 12.75L11.25 15 15 9.75M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
                />
              </svg>
              广泛的软件兼容
            </li>
            <li
              class="flex items-center gap-2 text-gray-700 dark:text-gray-300"
            >
              <svg
                xmlns="http://www.w3.org/2000/svg"
                class="h-5 w-5 text-green-500"
                fill="none"
                viewBox="0 0 24 24"
                stroke="currentColor"
              >
                <path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  d="M9 12.75L11.25 15 15 9.75M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
                />
              </svg>
              团队协作友好
            </li>
          </ul>
        </div>

        <!-- Markdown导出卡片 -->
        <div
          class="bg-white dark:bg-[#141419] rounded-xl p-8 border theme-border hover:shadow-xl transition-all group"
        >
          <div class="flex items-center gap-4 mb-6">
            <div
              class="p-3 rounded-xl bg-purple-100 dark:bg-purple-900/30 group-hover:scale-110 transition-transform"
            >
              <svg
                xmlns="http://www.w3.org/2000/svg"
                class="h-8 w-8 text-purple-600"
                fill="none"
                viewBox="0 0 24 24"
                stroke="currentColor"
              >
                <path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  d="M14.25 9.75L16.5 12l-2.25 2.25m-4.5 0L7.5 12l2.25-2.25M6 20.25h12A2.25 2.25 0 0020.25 18V6A2.25 2.25 0 0018 3.75H6A2.25 2.25 0 003.75 6v12A2.25 2.25 0 006 20.25z"
                />
              </svg>
            </div>
            <h3 class="text-xl font-bold text-gray-900 dark:text-gray-100">
              Markdown格式
            </h3>
          </div>
          <p class="text-gray-600 dark:text-gray-400 mb-6">
            轻量级标记语言，适合技术写作和在线发布。支持快速转换为其他格式，保持良好的可读性。
          </p>
          <ul class="space-y-3 text-sm">
            <li
              class="flex items-center gap-2 text-gray-700 dark:text-gray-300"
            >
              <svg
                xmlns="http://www.w3.org/2000/svg"
                class="h-5 w-5 text-green-500"
                fill="none"
                viewBox="0 0 24 24"
                stroke="currentColor"
              >
                <path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  d="M9 12.75L11.25 15 15 9.75M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
                />
              </svg>
              纯文本易读性
            </li>
            <li
              class="flex items-center gap-2 text-gray-700 dark:text-gray-300"
            >
              <svg
                xmlns="http://www.w3.org/2000/svg"
                class="h-5 w-5 text-green-500"
                fill="none"
                viewBox="0 0 24 24"
                stroke="currentColor"
              >
                <path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  d="M9 12.75L11.25 15 15 9.75M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
                />
              </svg>
              多平台支持
            </li>
            <li
              class="flex items-center gap-2 text-gray-700 dark:text-gray-300"
            >
              <svg
                xmlns="http://www.w3.org/2000/svg"
                class="h-5 w-5 text-green-500"
                fill="none"
                viewBox="0 0 24 24"
                stroke="currentColor"
              >
                <path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  d="M9 12.75L11.25 15 15 9.75M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
                />
              </svg>
              版本控制友好
            </li>
          </ul>
        </div>

        <!-- 音频导出卡片 -->
        <div
          class="bg-white dark:bg-[#141419] rounded-xl p-8 border theme-border hover:shadow-xl transition-all group"
        >
          <div class="flex items-center gap-4 mb-6">
            <div
              class="p-3 rounded-xl bg-amber-100 dark:bg-amber-900/30 group-hover:scale-110 transition-transform"
            >
              <svg
                xmlns="http://www.w3.org/2000/svg"
                class="h-8 w-8 text-amber-600"
                fill="none"
                viewBox="0 0 24 24"
                stroke="currentColor"
              >
                <path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  d="M9 9l10.5-3m0 6.553v3.75a2.25 2.25 0 01-1.632 2.163l-1.32.377a1.803 1.803 0 11-.99-3.467l2.31-.66a2.25 2.25 0 001.632-2.163zm0 0V2.25L9 5.25v10.303m0 0v3.75a2.25 2.25 0 01-1.632 2.163l-1.32.377a1.803 1.803 0 01-.99-3.467l2.31-.66A2.25 2.25 0 009 15.553z"
                />
              </svg>
            </div>
            <h3 class="text-xl font-bold text-gray-900 dark:text-gray-100">
              音频格式
            </h3>
          </div>
          <p class="text-gray-600 dark:text-gray-400 mb-6">
            支持多种音色的语音合成，让作品以听书形式呈现。适合移动场景和视障人士使用。
          </p>
          <ul class="space-y-3 text-sm">
            <li
              class="flex items-center gap-2 text-gray-700 dark:text-gray-300"
            >
              <svg
                xmlns="http://www.w3.org/2000/svg"
                class="h-5 w-5 text-green-500"
                fill="none"
                viewBox="0 0 24 24"
                stroke="currentColor"
              >
                <path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  d="M9 12.75L11.25 15 15 9.75M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
                />
              </svg>
              多种音色选择
            </li>
            <li
              class="flex items-center gap-2 text-gray-700 dark:text-gray-300"
            >
              <svg
                xmlns="http://www.w3.org/2000/svg"
                class="h-5 w-5 text-green-500"
                fill="none"
                viewBox="0 0 24 24"
                stroke="currentColor"
              >
                <path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  d="M9 12.75L11.25 15 15 9.75M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
                />
              </svg>
              离线播放支持
            </li>
            <li
              class="flex items-center gap-2 text-gray-700 dark:text-gray-300"
            >
              <svg
                xmlns="http://www.w3.org/2000/svg"
                class="h-5 w-5 text-green-500"
                fill="none"
                viewBox="0 0 24 24"
                stroke="currentColor"
              >
                <path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  d="M9 12.75L11.25 15 15 9.75M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
                />
              </svg>
              自然语音合成
            </li>
            <li
              class="flex items-center gap-2 text-gray-700 dark:text-gray-300"
            >
              <svg
                xmlns="http://www.w3.org/2000/svg"
                class="h-5 w-5 text-green-500"
                fill="none"
                viewBox="0 0 24 24"
                stroke="currentColor"
              >
                <path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  d="M9 12.75L11.25 15 15 9.75M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
                />
              </svg>
              情感语气调节
            </li>
          </ul>
        </div>
      </div>
    </div>
  </div>
  <div class="h-screen bg-[#f8fafc] dark:bg-[#030616] flex flex-col">
    <!-- 介绍多端板块 -->
    <div class="flex-1 container mx-auto px-4 flex flex-col">
      <div class="text-center mt-12 mb-8">
        <h2
          class="text-5xl font-bold font-serif text-gray-900 dark:text-white mb-4"
        >
          跨平台无缝体验
        </h2>
        <h2
          class="text-3xl font-bold font-serif text-gray-900 dark:text-white mb-4"
        >
          一次创作，多端同步
        </h2>
      </div>

      <div class="flex-1 flex flex-row items-start justify-center gap-12 px-8">
        <!-- 左侧文字介绍 -->
        <div class="w-5/12 h-full">
          <div class="text-left space-y-8">
            <!-- 小程序部分 -->
            <div class="space-y-4">
              <div class="flex items-center gap-4">
                <div class="p-2 bg-green-100 dark:bg-green-900/30 rounded-lg">
                  <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6 text-green-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" d="M12 18.75a6 6 0 006-6v-1.5m-6 7.5a6 6 0 01-6-6v-1.5m6 7.5v3.75m-3.75 0h7.5M12 15.75a3 3 0 01-3-3V4.5a3 3 0 116 0v8.25a3 3 0 01-3 3z" />
                  </svg>
                </div>
                <h3 class="text-2xl font-semibold text-gray-900 dark:text-white">微信小程序</h3>
              </div>
              <p class="text-gray-600 dark:text-gray-400 mb-4">
                无需下载安装，扫码即用。随时浏览精彩剧集，接收创作动态，参与团队讨论。支持与"星知"AI助手实时对话，让灵感激发与交流不受场景限制。
              </p>
              <div class="grid grid-cols-2 gap-4">
                <div class="flex items-center gap-2">
                  <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 text-green-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" d="M2.036 12.322a1.012 1.012 0 010-.639C3.423 7.51 7.36 4.5 12 4.5c4.638 0 8.573 3.007 9.963 7.178.07.207.07.431 0 .639C20.577 16.49 16.64 19.5 12 19.5c-4.638 0-8.573-3.007-9.963-7.178z" />
                    <path stroke-linecap="round" stroke-linejoin="round" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
                  </svg>
                  <span class="text-gray-600 dark:text-gray-400">作品浏览</span>
                </div>
                <div class="flex items-center gap-2">
                  <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 text-green-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" d="M14.857 17.082a23.848 23.848 0 005.454-1.31A8.967 8.967 0 0118 9.75v-.7V9A6 6 0 006 9v.75a8.967 8.967 0 01-2.312 6.022c1.733.64 3.56 1.085 5.455 1.31m5.714 0a24.255 24.255 0 01-5.714 0m5.714 0a3 3 0 11-5.714 0" />
                  </svg>
                  <span class="text-gray-600 dark:text-gray-400">消息通知</span>
                </div>
                <div class="flex items-center gap-2">
                  <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 text-green-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" d="M18 18.72a9.094 9.094 0 003.741-.479 3 3 0 00-4.682-2.72m.94 3.198l.001.031c0 .225-.012.447-.037.666A11.944 11.944 0 0112 21c-2.17 0-4.207-.576-5.963-1.584A6.062 6.062 0 016 18.719m12 0a5.971 5.971 0 00-.941-3.197m0 0A5.995 5.995 0 0012 12.75a5.995 5.995 0 00-5.058 2.772m0 0a3 3 0 00-4.681 2.72 8.986 8.986 0 003.74.477m.94-3.197a5.971 5.971 0 00-.94 3.197M15 6.75a3 3 0 11-6 0 3 3 0 016 0zm6 3a2.25 2.25 0 11-4.5 0 2.25 2.25 0 014.5 0zm-13.5 0a2.25 2.25 0 11-4.5 0 2.25 2.25 0 014.5 0z" />
                  </svg>
                  <span class="text-gray-600 dark:text-gray-400">团队群聊</span>
                </div>
                <div class="flex items-center gap-2">
                  <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 text-green-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" d="M9.813 15.904L9 18.75l-.813-2.846a4.5 4.5 0 00-3.09-3.09L2.25 12l2.846-.813a4.5 4.5 0 003.09-3.09L9 5.25l.813 2.846a4.5 4.5 0 003.09 3.09L15.75 12l-2.846.813a4.5 4.5 0 00-3.09 3.09zM18.259 8.715L18 9.75l-.259-1.035a3.375 3.375 0 00-2.455-2.456L14.25 6l1.036-.259a3.375 3.375 0 002.455-2.456L18 2.25l.259 1.035a3.375 3.375 0 002.456 2.456L21.75 6l-1.035.259a3.375 3.375 0 00-2.456 2.456zM16.894 20.567L16.5 21.75l-.394-1.183a2.25 2.25 0 00-1.423-1.423L13.5 18.75l1.183-.394a2.25 2.25 0 001.423-1.423l.394-1.183.394 1.183a2.25 2.25 0 001.423 1.423l1.183.394-1.183.394a2.25 2.25 0 00-1.423 1.423z" />
                  </svg>
                  <span class="text-gray-600 dark:text-gray-400">星知助手</span>
                </div>
              </div>
            </div>

            <!-- Web端部分 -->
            <div class="space-y-4">
              <div class="flex items-center gap-4">
                <div class="p-2 bg-blue-100 dark:bg-blue-900/30 rounded-lg">
                  <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6 text-blue-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" d="M9 17.25v1.007a3 3 0 01-.879 2.122L7.5 21h9l-.621-.621A3 3 0 0115 18.257V17.25m6-12V15a2.25 2.25 0 01-2.25 2.25H5.25A2.25 2.25 0 013 15V5.25m18 0A2.25 2.25 0 0018.75 3H5.25A2.25 2.25 0 003 5.25m18 0V12a2.25 2.25 0 01-2.25 2.25H5.25A2.25 2.25 0 013 12V5.25" />
                  </svg>
                </div>
                <h3 class="text-2xl font-semibold text-gray-900 dark:text-white">Web 端</h3>
              </div>
              <p class="text-gray-600 dark:text-gray-400">
                跨平台网页应用，告别环境限制。支持所有主流浏览器，无论 Windows、Mac 还是 Linux，都能享受到完整的创作体验。专业的创作工具集成、云端数据实时同步，让您的创意无处不在。
              </p>
              <div class="grid grid-cols-2 gap-4">
                <div class="flex items-center gap-2">
                  <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 text-blue-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" d="M16.862 4.487l1.687-1.688a1.875 1.875 0 112.652 2.652L10.582 16.07a4.5 4.5 0 01-1.897 1.13L6 18l.8-2.685a4.5 4.5 0 011.13-1.897l8.932-8.931zm0 0L19.5 7.125M18 14v4.75A2.25 2.25 0 0115.75 21H5.25A2.25 2.25 0 013 18.75V8.25A2.25 2.25 0 015.25 6H10" />
                  </svg>
                  <span class="text-gray-600 dark:text-gray-400">专业创作</span>
                </div>
                <div class="flex items-center gap-2">
                  <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 text-blue-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" d="M3 13.125C3 12.504 3.504 12 4.125 12h2.25c.621 0 1.125.504 1.125 1.125v6.75C7.5 20.496 6.996 21 6.375 21h-2.25A1.125 1.125 0 013 19.875v-6.75zM9.75 8.625c0-.621.504-1.125 1.125-1.125h2.25c.621 0 1.125.504 1.125 1.125v11.25c0 .621-.504 1.125-1.125 1.125h-2.25a1.125 1.125 0 01-1.125-1.125V8.625zM16.5 4.125c0-.621.504-1.125 1.125-1.125h2.25C20.496 3 21 3.504 21 4.125v15.75c0 .621-.504 1.125-1.125 1.125h-2.25a1.125 1.125 0 01-1.125-1.125V4.125z" />
                  </svg>
                  <span class="text-gray-600 dark:text-gray-400">数据分析</span>
                </div>
                <div class="flex items-center gap-2">
                  <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 text-blue-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" d="M15 19.128a9.38 9.38 0 002.625.372 9.337 9.337 0 004.121-.952 4.125 4.125 0 00-7.533-2.493M15 19.128v-.003c0-1.113-.285-2.16-.786-3.07M15 19.128v.106A12.318 12.318 0 018.624 21c-2.331 0-4.512-.645-6.374-1.766l-.001-.109a6.375 6.375 0 0111.964-3.07M12 6.375a3.375 3.375 0 11-6.75 0 3.375 3.375 0 016.75 0zm8.25 2.25a2.625 2.625 0 11-5.25 0 2.625 2.625 0 015.25 0z" />
                  </svg>
                  <span class="text-gray-600 dark:text-gray-400">团队协作</span>
                </div>
                <div class="flex items-center gap-2">
                  <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 text-blue-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" d="M12 6v6h4.5m4.5 0a9 9 0 11-18 0 9 9 0 0118 0z" />
                  </svg>
                  <span class="text-gray-600 dark:text-gray-400">实时同步</span>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- 右侧展示图 -->
        <div class="w-7/12 flex flex-col justify-center items-center gap-4">
          <div class="flex-1 transition-all duration-500 hover:scale-105">
            <img
              :src="webs"
              alt="界面"
              class="h-[70vh] object-contain rounded-lg my-auto shadow-xl"
            />
          </div>
        </div>
      </div>
    </div>
  </div>
  <div class="h-screen bg-[#f8fafc] dark:bg-[#030616] flex flex-col">
    <!-- 介绍巨星全知板块 -->
    <div class="container mx-auto px-4 py-16">
      <div class="text-center mb-12">
        <h2 class="text-5xl font-bold font-serif text-gray-900 dark:text-white mb-4">
          遇见"星知"
        </h2>
        <h3 class="text-3xl font-bold font-serif text-gray-900 dark:text-white mb-4">
          创作者的智能伙伴，观众的贴心向导
        </h3>
        <p class="text-xl text-gray-600 dark:text-gray-400 max-w-3xl mx-auto">
          基于海量剧集数据与先进AI技术，为创作全程提供专业支持，让每个故事都绽放光彩
        </p>
      </div>

      <div class="flex gap-12">
        <!-- 左侧介绍 -->
        <div class="w-1/2 space-y-8">
          <div class="space-y-4">
            <!-- 创作者服务 -->
            <div class="flex items-start gap-4">
              <div class="p-2 bg-blue-100 dark:bg-blue-900/30 rounded-lg mt-1">
                <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6 text-blue-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" d="M12 18v-5.25m0 0a6.01 6.01 0 001.5-.189m-1.5.189a6.01 6.01 0 01-1.5-.189m3.75 7.478a12.06 12.06 0 01-4.5 0m3.75 2.383a14.406 14.406 0 01-3 0M14.25 18v-.192c0-.983.658-1.823 1.508-2.316a7.5 7.5 0 10-7.517 0c.85.493 1.509 1.333 1.509 2.316V18" />
                </svg>
              </div>
              <div>
                <h4 class="text-lg font-semibold text-gray-900 dark:text-gray-100">
                  创作者服务
                </h4>
                <p class="text-gray-600 dark:text-gray-400">
                  提供剧本创意构建、市场趋势分析、项目优化建议和创作技巧分享
                </p>
              </div>
            </div>

            <!-- 观众服务 -->
            <div class="flex items-start gap-4">
              <div class="p-2 bg-purple-100 dark:bg-purple-900/30 rounded-lg mt-1">
                <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6 text-purple-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" d="M15 19.128a9.38 9.38 0 002.625.372 9.337 9.337 0 004.121-.952 4.125 4.125 0 00-7.533-2.493M15 19.128v-.003c0-1.113-.285-2.16-.786-3.07M15 19.128v.106A12.318 12.318 0 018.624 21c-2.331 0-4.512-.645-6.374-1.766l-.001-.109a6.375 6.375 0 0111.964-3.07M12 6.375a3.375 3.375 0 11-6.75 0 3.375 3.375 0 016.75 0zm8.25 2.25a2.625 2.625 0 11-5.25 0 2.625 2.625 0 015.25 0z" />
                </svg>
              </div>
              <div>
                <h4 class="text-lg font-semibold text-gray-900 dark:text-gray-100">
                  观众服务
                </h4>
                <p class="text-gray-600 dark:text-gray-400">
                  提供个性化内容推荐、热门作品解析、类似作品推荐和观看建议
                </p>
              </div>
            </div>

            <!-- 平台服务 -->
            <div class="flex items-start gap-4">
              <div class="p-2 bg-green-100 dark:bg-green-900/30 rounded-lg mt-1">
                <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6 text-green-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" d="M9 12.75L11.25 15 15 9.75M21 12c0 1.268-.63 2.39-1.593 3.068a3.745 3.745 0 01-1.043 3.296 3.745 3.745 0 01-3.296 1.043A3.745 3.745 0 0112 21c-1.268 0-2.39-.63-3.068-1.593a3.746 3.746 0 01-3.296-1.043 3.745 3.745 0 01-1.043-3.296A3.745 3.745 0 013 12c0-1.268.63-2.39 1.593-3.068a3.745 3.745 0 011.043-3.296 3.746 3.746 0 013.296-1.043A3.746 3.746 0 0112 3c1.268 0 2.39.63 3.068 1.593a3.746 3.746 0 013.296 1.043 3.746 3.746 0 011.043 3.296A3.745 3.745 0 0121 12z" />
                </svg>
              </div>
              <div>
                <h4 class="text-lg font-semibold text-gray-900 dark:text-gray-100">
                  平台服务
                </h4>
                <p class="text-gray-600 dark:text-gray-400">
                  提供功能指南、创作工具使用建议、社区互动和平台资源推荐
                </p>
              </div>
            </div>

            <!-- 行业洞察 -->
            <div class="flex items-start gap-4">
              <div class="p-2 bg-amber-100 dark:bg-amber-900/30 rounded-lg mt-1">
                <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6 text-amber-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" d="M2.25 18L9 11.25l4.306 4.307a11.95 11.95 0 015.814-5.519l2.74-1.22m0 0l-5.94-2.28m5.94 2.28l-2.28 5.941" />
                </svg>
              </div>
              <div>
                <h4 class="text-lg font-semibold text-gray-900 dark:text-gray-100">
                  行业洞察
                </h4>
                <p class="text-gray-600 dark:text-gray-400">
                  提供剧集市场趋势、用户偏好研究、新兴题材机会和行业动态分析
                </p>
              </div>
            </div>

            <!-- 数据支持 -->
            <div class="flex items-start gap-4">
              <div class="p-2 bg-rose-100 dark:bg-rose-900/30 rounded-lg mt-1">
                <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6 text-rose-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" d="M3 13.125C3 12.504 3.504 12 4.125 12h2.25c.621 0 1.125.504 1.125 1.125v6.75C7.5 20.496 6.996 21 6.375 21h-2.25A1.125 1.125 0 013 19.875v-6.75zM9.75 8.625c0-.621.504-1.125 1.125-1.125h2.25c.621 0 1.125.504 1.125 1.125v11.25c0 .621-.504 1.125-1.125 1.125h-2.25a1.125 1.125 0 01-1.125-1.125V8.625zM16.5 4.125c0-.621.504-1.125 1.125-1.125h2.25C20.496 3 21 3.504 21 4.125v15.75c0 .621-.504 1.125-1.125 1.125h-2.25a1.125 1.125 0 01-1.125-1.125V4.125z" />
                </svg>
              </div>
              <div>
                <h4 class="text-lg font-semibold text-gray-900 dark:text-gray-100">
                  数据支持
                </h4>
                <p class="text-gray-600 dark:text-gray-400">
                  基于平台实时数据，提供作品表现分析、受众反馈和市场定位建议
                </p>
              </div>
            </div>
          </div>
        </div>

        <!-- 右侧对话示例 -->
        <div class="w-1/2">
          <div class="bg-white dark:bg-zinc-900 rounded-xl shadow-xl border theme-border p-6">
            <div class="flex items-center gap-3 mb-4">
              <div class="flex space-x-2">
                <div class="w-3 h-3 rounded-full bg-red-500"></div>
                <div class="w-3 h-3 rounded-full bg-yellow-500"></div>
                <div class="w-3 h-3 rounded-full bg-green-500"></div>
              </div>
              <div class="flex-1 text-center text-sm text-gray-600 dark:text-gray-400">
                与星知对话
              </div>
            </div>

            <div class="space-y-4 h-[400px] overflow-y-auto p-4 bg-gray-50 dark:bg-zinc-800 rounded-lg">
              <!-- 用户消息 -->
              <div class="flex gap-3 justify-end">
                <div class="bg-blue-500 text-white rounded-lg p-3 max-w-[80%]">
                  我想写一个科幻悬疑剧，但不知道如何构建吸引人的开场。
                </div>
                <div class="w-8 h-8 rounded-full bg-blue-100 flex items-center justify-center flex-shrink-0">
                  <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 text-blue-500" viewBox="0 0 20 20" fill="currentColor">
                    <path fill-rule="evenodd" d="M10 9a3 3 0 100-6 3 3 0 000 6zm-7 9a7 7 0 1114 0H3z" clip-rule="evenodd" />
                  </svg>
                </div>
              </div>

              <!-- AI回复 -->
              <div class="flex gap-3">
                <div class="w-8 h-8 rounded-full flex items-center justify-center flex-shrink-0">
                  <img :src="ai" class="p-1"/>
                </div>
                <div class="bg-gray-100 text-white dark:bg-zinc-700 rounded-lg p-3 max-w-[80%]">
                  根据平台数据分析，成功的科幻悬疑剧通常以一个令人费解的事件开场，比如：
                  
                  1. 一个看似普通的科技产品突然表现出超出设计的行为
                  2. 一位重要人物离奇失踪，留下一串神秘数字
                  3. 未来城市中出现不符合时代的古老物件
                  
                  建议你可以从这些方向入手，让观众在第一时间产生好奇心。要注意的是，开场悬念要与后续剧情有紧密联系，避免成为无意义的噱头。
                </div>
              </div>

              <!-- 用户消息 -->
              <div class="flex gap-3 justify-end">
                <div class="bg-blue-500 text-white rounded-lg p-3 max-w-[80%]">
                  这些建议很棒！能否帮我分析一下最近热门的科幻悬疑剧都有哪些共同特点？
                </div>
                <div class="w-8 h-8 rounded-full bg-blue-100 flex items-center justify-center flex-shrink-0">
                  <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 text-blue-500" viewBox="0 0 20 20" fill="currentColor">
                    <path fill-rule="evenodd" d="M10 9a3 3 0 100-6 3 3 0 000 6zm-7 9a7 7 0 1114 0H3z" clip-rule="evenodd" />
                  </svg>
                </div>
              </div>

              <!-- AI回复 -->
              <div class="flex gap-3">
                <div class="w-8 h-8 rounded-full flex items-center justify-center flex-shrink-0">
                  <img :src="ai" class="p-1"/>
                </div>
                <div class="bg-gray-100 text-white dark:bg-zinc-700 rounded-lg p-3 max-w-[80%]">
                  通过分析平台近期热门作品，我发现以下特点最受观众欢迎：

                  1. 多维叙事：在未来科技背景下展开多条时间线或平行宇宙的故事
                  2. 伦理思考：将人工智能、生物科技等前沿话题与道德困境相结合
                  3. 情感共鸣：即使在科幻背景下，也要注重人物的情感发展
                  4. 悬疑节奏：善用信息差和伏笔，让观众保持探索欲
                  
                  我建议你可以从这些角度出发，结合自己的创意，打造独特的作品。
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
  <div class="bg-white dark:bg-black flex font-sans min-h-32 text-theme-switch">
    <div class="mx-auto my-auto flex-col text-center">
      <span class="cursor-pointer hover:text-blue-600">
        创剧星球IdeaCosmos
      </span>
      <br />
      <span class="mx-auto my-auto cursor-pointer hover:text-blue-600">
        使用Vue.js&Golang构建
      </span>
    </div>
  </div>
</template>

<style scoped></style>
