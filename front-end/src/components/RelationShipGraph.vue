<template>
  <div :class="['relative', themeStore.isDark ? 'dark' : '']">
    <v-chart :option="chartOptions" style="width: 100%; height: 600px;"/>
  </div>
</template>

<script>
import { defineComponent, watch } from "vue";
import { useThemeStore } from "@/stores/theme"; // 导入 themeStore
import { BACKEND_DOMAIN } from "@/util/VARRIBLES";
import VChart from "vue-echarts";
import { use } from "echarts/core";
import { GraphChart } from "echarts/charts";
import { TitleComponent, TooltipComponent } from "echarts/components";
import { CanvasRenderer } from "echarts/renderers";

use([GraphChart, TitleComponent, TooltipComponent, CanvasRenderer]);

// 默认用户头像 SVG
const defaultAvatar = `
<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor">
  <path stroke-linecap="round" stroke-linejoin="round" d="M15.75 6a3.75 3.75 0 11-7.5 0 3.75 3.75 0 017.5 0zM4.501 20.118a7.5 7.5 0 0114.998 0A17.933 17.933 0 0112 21.75c-2.676 0-5.216-.584-7.499-1.632z" />
</svg>
`;

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
            symbol: firstCharacter.avatar 
              ? `image://${BACKEND_DOMAIN}uploads/${firstCharacter.avatar}`
              : `path://${defaultAvatar}`,
            symbolSize: 80,
            itemStyle: {
              borderWidth: 4,
              borderColor: isDark ? 'rgba(255,255,255,0.3)' : 'rgba(0,0,0,0.2)'
            }
          });
        }

        if (!nodes.some((node) => node.name === secondCharacter.name)) {
          nodes.push({
            name: secondCharacter.name,
            info: secondCharacter.description,
            symbol: secondCharacter.avatar 
              ? `image://${BACKEND_DOMAIN}uploads/${secondCharacter.avatar}`
              : `path://${defaultAvatar}`,
            symbolSize: 80,
            itemStyle: {
              borderWidth: 4,
              borderColor: isDark ? 'rgba(255,255,255,0.3)' : 'rgba(0,0,0,0.2)'
            }
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
          show: true,
          confine: true,
          formatter: function (params) {
            if (params.dataType === "node") {
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
                    角色：${params.data.name}
                  </div>
                  <div style="
                    color: rgba(255,255,255,0.9);
                    white-space: normal;
                    word-break: break-all;
                  ">
                    描述：${params.data.info}
                  </div>
                </div>
              `;
            } else if (params.dataType === "edge") {
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
                    关系：${params.data.value}
                  </div>
                  <div style="
                    color: rgba(255,255,255,0.9);
                    white-space: normal;
                    word-break: break-all;
                  ">
                    描述：${params.data.story}
                  </div>
                </div>
              `;
            }
          },
          backgroundColor: 'rgba(0,0,0,0.75)',
          borderRadius: 8,
          padding: [10, 15],
          textStyle: {
            color: '#fff',
            fontSize: 14
          }
        },
        series: [{
          type: 'graph',
          layout: 'force',
          data: nodes,
          links: links,
          roam: true,
          draggable: true,
          circular: true,
          symbol: 'circle',
          symbolKeepAspect: true,
          label: {
            show: true,
            position: 'bottom', // 将名字显示在头像下方
            formatter: '{b}',
            backgroundColor: 'rgba(102, 102, 102, 0.75)',
            borderRadius: 4,
            padding: [4, 8],
            color: '#fff',
            fontSize: 14,
            fontWeight: 'bold',
            distance: 10 // 增加标签与节点的距离
          },
          edgeLabel: {
            show: true,
            formatter: '{c}',
            backgroundColor: 'rgba(0, 0, 0, 0.75)',
            borderRadius: 4,
            padding: [4, 8],
            color: '#fff',
            fontSize: 12,
            distance: 20,
            rotate: 0,
            offset: [0, 0]
          },
          force: {
            repulsion: 1000, // 增加斥力
            gravity: 0.1,
            edgeLength: 400, // 增加边长
            friction: 0.1,
            layoutAnimation: true,
            initLayout: 'circular'
          },
          cursor: 'move',
          zoom: 1.5,
          center: ['50%', '50%'],
          emphasis: {
            focus: 'adjacency',
            scale: 1.2,
            lineStyle: {
              width: 4
            },
            label: {
              fontSize: 16,
              show: true
            }
          },
          nodeScaleRatio: 0.6,
          itemStyle: {
            borderWidth: 2,
            borderColor: isDark ? 'rgba(255,255,255,0.3)' : 'rgba(0,0,0,0.2)',
            shadowColor: 'rgba(0, 0, 0, 0.2)',
            shadowBlur: 10
          },
          lineStyle: {
            width: 2,
            opacity: 0.7,
            curveness: 0.3,
            color: {
              type: 'linear',
              x: 0,
              y: 0,
              x2: 0,
              y2: 1,
              colorStops: [{
                offset: 0, color: 'rgb(77, 77, 245)'
              }, {
                offset: 1, color: 'rgb(77, 77, 245)'
              }]
            }
          },
          edgeSymbol: ['none', 'arrow'],
          edgeSymbolSize: [0, 8],
        }],
        backgroundColor: 'transparent'
      };
    },
  },
});
</script>

<style scoped>
/* 样式 */
</style>
