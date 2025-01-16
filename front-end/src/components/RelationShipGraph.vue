<template>
  <div :class="['relative', themeStore.isDark ? 'dark' : '']">
    <v-chart :option="chartOptions" style="width: 100%; height: 600px;"/>
  </div>
</template>

<script>
import { defineComponent, watch } from "vue";
import { useThemeStore } from "@/stores/theme"; // 导入 themeStore
import VChart from "vue-echarts";
import { use } from "echarts/core";
import { GraphChart } from "echarts/charts";
import { TitleComponent, TooltipComponent } from "echarts/components";
import { CanvasRenderer } from "echarts/renderers";

use([GraphChart, TitleComponent, TooltipComponent, CanvasRenderer]);

function randomColor() {
  const letters = "0123456789ABCDEF";
  let color = "#";
  for (let i = 0; i < 6; i++) {
    color += letters[Math.floor(Math.random() * 16)];
  }
  return color;
}

export default defineComponent({
  components: { VChart },
  props: {
    relationships: {
      type: Array,
      required: true,
    },
  },
  setup() {
    const themeStore = useThemeStore(); // 使用 themeStore
    return { themeStore };
  },
  data() {
    return {
      chartOptions: {},
    };
  },
  watch: {
    relationships: {
      immediate: true,
      handler(newRelationships) {
        this.updateChartOptions(newRelationships);
      },
    },
    "themeStore.isDark": {
      immediate: true,
      handler() {
        this.updateChartOptions(this.relationships); // 主题切换时更新图表
      },
    },
  },
  methods: {
    updateChartOptions(relationships) {
      const isDark = this.themeStore.isDark;

      const nodes = [];
      const links = [];

      relationships.forEach((rel) => {
        const firstCharacter = rel.first_character;
        const secondCharacter = rel.second_character;

        if (!nodes.some((node) => node.name === firstCharacter.name)) {
          nodes.push({
            name: firstCharacter.name,
            info: firstCharacter.description,
            itemStyle: { color: randomColor() },
          });
        }

        if (!nodes.some((node) => node.name === secondCharacter.name)) {
          nodes.push({
            name: secondCharacter.name,
            info: secondCharacter.description,
            itemStyle: { color: randomColor() },
          });
        }

        links.push({
          source: firstCharacter.name,
          target: secondCharacter.name,
          value: rel.name,
          story: rel.content,
        });
      });

      this.chartOptions = {
        tooltip: {
          trigger: "item",
          backgroundColor: isDark ? "#333" : "#FFF", // 提示框背景色
          textStyle: { color: isDark ? "#FFF" : "#000" }, // 提示框文字颜色
          formatter(params) {
            const style = `
              <div style="
                max-width: 400px;
                word-wrap: break-word;
                white-space: normal;
                line-height: 1.5;
              ">
            `;
            if (params.dataType === "node") {
              return `${style}${params.data.name}<br />${params.data.info || "无详细信息"}</div>`;
            } else if (params.dataType === "edge") {
              return `${style}${params.data.source} 与 ${params.data.target} 的关系：${params.data.value}<br />故事：${params.data.story}</div>`;
            }
            return "";
          },
        },
        series: [
          {
            type: "graph",
            layout: "force",
            symbolSize: 50,
            roam: true,
            focusNodeAdjacency: true,
            label: {
              show: true,
              position: "right",
              color: isDark ? "#FFF" : "#000", // 节点文字颜色
              formatter: "{b}",
            },
            edgeLabel: {
              show: true,
              formatter: "{c}",
              color: isDark ? "#AAA" : "#555", // 边文字颜色
            },
            force: {
              repulsion: 1000,
              edgeLength: [50, 150],
            },
            lineStyle: {
              color: isDark ? "#888" : "#CCC", // 连接线颜色
              opacity: isDark ? 0.7 : 0.9, // 透明度
            },
            data: nodes,
            links: links,
          },
        ],
      };
    },
  },
});
</script>

<style scoped>
/* 样式 */
</style>
