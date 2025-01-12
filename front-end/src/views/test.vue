<template>
  <div :class="['relative', darkMode ? 'dark' : '']">
    <v-chart :option="chartOptions" style="width: 100%; height: 600px;" />
  </div>
</template>

<script>
import { defineComponent } from "vue";
import VChart from "vue-echarts"; // 引入 vue-echarts
import { use } from "echarts/core"; // 用于引入 ECharts 模块
import { GraphChart } from "echarts/charts"; // 力导向图
import { TitleComponent, TooltipComponent } from "echarts/components"; // 标题、提示组件
import { CanvasRenderer } from "echarts/renderers"; // 渲染器

use([GraphChart, TitleComponent, TooltipComponent, CanvasRenderer]); // 注册使用的模块

// 随机生成颜色
function randomColor() {
  const letters = '0123456789ABCDEF';
  let color = '#';
  for (let i = 0; i < 6; i++) {
    color += letters[Math.floor(Math.random() * 16)];
  }
  return color;
}

export default defineComponent({
  components: {
    VChart,
  },
  data() {
    return {
      darkMode: false, // 默认情况下不启用暗黑模式
      chartOptions: {
        tooltip: {
          trigger: "item", // 修改为 "item" 来触发边上的 tooltip
          formatter: function (params) {
            if (params.dataType === "node") {
              return `${params.data.name}<br />${params.data.info || "无详细信息"}`;
            } else if (params.dataType === "edge") {
              return `${params.data.source} 与 ${params.data.target} 的关系：${params.data.value}<br />故事：${params.data.story}`;
            }
            return "";
          },
        },
        series: [
          {
            type: "graph",
            layout: "force", // 力导向布局
            symbolSize: 50,
            roam: true, // 允许拖拽
            focusNodeAdjacency: true, // 鼠标悬浮时突出显示相关联节点
            label: {
              show: true,
              position: "right",
              formatter: "{b}", // 显示人物名字
            },
            edgeLabel: {
              show: true,
              formatter: "{c}",
            },
            force: {
              repulsion: 1000, // 排斥力
              edgeLength: [50, 150], // 边的长度范围
            },
            data: [
              { name: "简爱", info: "简爱是女主角，坚强勇敢，聪明睿智。", itemStyle: { color: randomColor() } },
              { name: "爱德华·罗切斯特", info: "简爱的爱人，庄园主人，复杂多面的角色。", itemStyle: { color: randomColor() } },
              { name: "海伦", info: "简爱在孤儿院的朋友，性格温柔，智慧超凡。", itemStyle: { color: randomColor() } },
              { name: "伯莎", info: "罗切斯特的妻子，神秘、悲剧的角色。", itemStyle: { color: randomColor() } },
              // 其他人物...
            ],
            links: [
              {
                source: "简爱",
                target: "爱德华·罗切斯特",
                value: "爱情",
                story: "简爱与罗切斯特之间的爱情历经了重重考验。两人跨越了阶级与困境，最终走到了一起。"
              },
              {
                source: "简爱",
                target: "海伦",
                value: "友情",
                story: "海伦是简爱在孤儿院的朋友，海伦的温柔与智慧深深影响了简爱，让她学会了宽容与坚强。"
              },
              {
                source: "爱德华·罗切斯特",
                target: "伯莎",
                value: "婚姻关系",
                story: "罗切斯特与伯莎的婚姻是悲剧的开始。伯莎的疯狂给两人的生活带来了毁灭性的影响。"
              },
              // 其他关系...
            ],
          },
        ],
      },
    };
  }
});
</script>

<style scoped>

/* 其他 TailwindCSS 样式 */
</style>
