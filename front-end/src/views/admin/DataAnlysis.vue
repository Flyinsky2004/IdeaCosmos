<script setup>
import { ref, onMounted, nextTick, onUnmounted, computed, watch } from 'vue'
import { get } from '@/util/request'
import { message } from 'ant-design-vue'
import * as echarts from 'echarts'
import SpinLoaderLarge from '@/components/spinLoaderLarge.vue'
import { useThemeStore } from "@/stores/theme"

const themeStore = useThemeStore()
const loading = ref(true)
const statsData = ref(null)
const projectTypeData = ref([])

// 图表引用
const userTrendChart = ref(null)
const projectTrendChart = ref(null)
const reviewTrendChart = ref(null)
const scoreDistChart = ref(null)
const projectTypeChart = ref(null)

// 所有图表实例
const chartInstances = []

// 颜色主题
const colorPalette = [
  '#3B82F6', '#10B981', '#F59E0B', '#8B5CF6', '#EC4899', 
  '#EF4444', '#06B6D4', '#14B8A6', '#F97316', '#84CC16'
]

// 加载管理统计数据
const fetchStatistics = async () => {
  loading.value = true
  
  get('/api/admin/statistics/overview', {}, 
    (msg, data) => {
      statsData.value = data
      // 获取项目类型数据
      fetchProjectTypeStats()
    },
    (msg) => {
      message.warning(msg)
      loading.value = false
    },
    (msg) => {
      message.error(msg)
      loading.value = false
    }
  )
}

// 获取项目类型统计
const fetchProjectTypeStats = () => {
  get('/api/admin/statistics/project-types', {}, 
    (msg, data) => {
      projectTypeData.value = data
      loading.value = false
      
      nextTick(() => {
        initCharts()
      })
    },
    (msg) => {
      message.warning(msg)
      loading.value = false
    },
    (msg) => {
      message.error(msg)
      loading.value = false
    }
  )
}

const initCharts = () => {
  // 清除之前的图表实例
  chartInstances.forEach(chart => {
    chart.dispose()
  })
  chartInstances.length = 0
  
  // 初始化各图表
  initUserTrendChart()
  initProjectTrendChart()
  initReviewTrendChart()
  initScoreDistChart()
  initProjectTypeChart()
}

// 用户注册趋势
const initUserTrendChart = () => {
  if (!userTrendChart.value) return
  
  const chart = echarts.init(userTrendChart.value, null, {
    renderer: 'canvas'
  })
  chartInstances.push(chart)
  
  const data = statsData.value.trend_data.user_registration || []
  
  // 计算累计用户数
  const cumulativeUsers = []
  let cumulative = 0
  data.forEach(item => {
    cumulative += item.count
    cumulativeUsers.push(cumulative)
  })
  
  chart.setOption({
    title: {
      text: '用户注册趋势',
      left: 'center',
      textStyle: {
        color: themeStore.isDark ? '#e5e7eb' : '#333',
        fontWeight: 'normal'
      }
    },
    tooltip: {
      trigger: 'axis',
      axisPointer: {
        type: 'cross'
      },
      formatter: function(params) {
        const date = params[0].name
        let html = `<div>${date}</div>`
        
        params.forEach(param => {
          html += `<div style="display:flex;align-items:center;margin:5px 0;">
            <span style="display:inline-block;width:10px;height:10px;border-radius:50%;background:${param.color};margin-right:5px;"></span>
            <span>${param.seriesName}: ${param.value}</span>
          </div>`
        })
        
        return html
      }
    },
    legend: {
      data: ['每日新增', '累计用户'],
      textStyle: {
        color: themeStore.isDark ? '#e5e7eb' : '#333'
      },
      bottom: 0
    },
    grid: {
      left: '3%',
      right: '4%',
      bottom: '10%',
      containLabel: true
    },
    xAxis: {
      type: 'category',
      data: data.map(item => item.date),
      axisLabel: {
        color: themeStore.isDark ? '#9ca3af' : '#666'
      },
      axisLine: {
        lineStyle: {
          color: themeStore.isDark ? '#4b5563' : '#ddd'
        }
      }
    },
    yAxis: [
      {
        type: 'value',
        name: '每日新增',
        nameTextStyle: {
          color: themeStore.isDark ? '#9ca3af' : '#666'
        },
        axisLabel: {
          color: themeStore.isDark ? '#9ca3af' : '#666'
        },
        axisLine: {
          show: true,
          lineStyle: {
            color: themeStore.isDark ? '#4b5563' : '#ddd'
          }
        },
        splitLine: {
          lineStyle: {
            color: themeStore.isDark ? '#374151' : '#eee'
          }
        }
      },
      {
        type: 'value',
        name: '累计用户',
        nameTextStyle: {
          color: themeStore.isDark ? '#9ca3af' : '#666'
        },
        axisLabel: {
          color: themeStore.isDark ? '#9ca3af' : '#666'
        },
        axisLine: {
          show: true,
          lineStyle: {
            color: themeStore.isDark ? '#4b5563' : '#ddd'
          }
        },
        splitLine: {
          show: false
        }
      }
    ],
    series: [
      {
        name: '每日新增',
        type: 'bar',
        barWidth: '60%',
        data: data.map(item => item.count),
        itemStyle: {
          color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
            { offset: 0, color: '#3B82F6' },
            { offset: 1, color: '#60A5FA' }
          ])
        }
      },
      {
        name: '累计用户',
        type: 'line',
        yAxisIndex: 1,
        smooth: true,
        symbol: 'circle',
        symbolSize: 6,
        lineStyle: {
          width: 3,
          color: '#10B981'
        },
        itemStyle: {
          color: '#10B981',
          borderColor: themeStore.isDark ? '#1f2937' : '#fff',
          borderWidth: 2
        },
        data: cumulativeUsers
      }
    ]
  })
  
  window.addEventListener('resize', () => chart.resize())
}

// 项目创建趋势
const initProjectTrendChart = () => {
  if (!projectTrendChart.value) return
  
  const chart = echarts.init(projectTrendChart.value, null, {
    renderer: 'canvas'
  })
  chartInstances.push(chart)
  
  const data = statsData.value.trend_data.project_creation || []
  
  chart.setOption({
    title: {
      text: '项目创建趋势',
      left: 'center',
      textStyle: {
        color: themeStore.isDark ? '#e5e7eb' : '#333',
        fontWeight: 'normal'
      }
    },
    tooltip: {
      trigger: 'axis',
      axisPointer: {
        type: 'line'
      }
    },
    grid: {
      left: '3%',
      right: '4%',
      bottom: '3%',
      containLabel: true
    },
    xAxis: {
      type: 'category',
      data: data.map(item => item.date),
      axisLabel: {
        color: themeStore.isDark ? '#9ca3af' : '#666'
      },
      axisLine: {
        lineStyle: {
          color: themeStore.isDark ? '#4b5563' : '#ddd'
        }
      }
    },
    yAxis: {
      type: 'value',
      axisLabel: {
        color: themeStore.isDark ? '#9ca3af' : '#666'
      },
      axisLine: {
        show: false
      },
      splitLine: {
        lineStyle: {
          color: themeStore.isDark ? '#374151' : '#eee'
        }
      }
    },
    series: [
      {
        name: '新建项目',
        type: 'line',
        smooth: true,
        symbol: 'circle',
        symbolSize: 8,
        data: data.map(item => item.count),
        lineStyle: {
          width: 4,
          color: new echarts.graphic.LinearGradient(0, 0, 1, 0, [
            { offset: 0, color: '#10B981' },
            { offset: 1, color: '#34D399' }
          ])
        },
        itemStyle: {
          color: '#10B981',
          borderColor: '#fff',
          borderWidth: 2
        },
        areaStyle: {
          color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
            { offset: 0, color: 'rgba(16, 185, 129, 0.4)' },
            { offset: 1, color: 'rgba(16, 185, 129, 0.1)' }
          ])
        }
      }
    ]
  })
  
  window.addEventListener('resize', () => chart.resize())
}

// 章节审核趋势
const initReviewTrendChart = () => {
  if (!reviewTrendChart.value) return
  
  const chart = echarts.init(reviewTrendChart.value, null, {
    renderer: 'canvas'
  })
  chartInstances.push(chart)
  
  const data = statsData.value.trend_data.chapter_reviews || []
  
  // 计算总审核数
  const totalReviews = data.map(item => item.approved + item.rejected)
  
  chart.setOption({
    title: {
      text: '章节审核趋势',
      left: 'center',
      textStyle: {
        color: themeStore.isDark ? '#e5e7eb' : '#333',
        fontWeight: 'normal'
      }
    },
    tooltip: {
      trigger: 'axis',
      axisPointer: {
        type: 'shadow'
      }
    },
    legend: {
      data: ['通过', '拒绝', '总数'],
      bottom: 0,
      textStyle: {
        color: themeStore.isDark ? '#e5e7eb' : '#333'
      }
    },
    grid: {
      left: '3%',
      right: '4%',
      bottom: '10%',
      containLabel: true
    },
    xAxis: {
      type: 'category',
      data: data.map(item => item.date),
      axisLabel: {
        color: themeStore.isDark ? '#9ca3af' : '#666'
      },
      axisTick: {
        alignWithLabel: true
      },
      axisLine: {
        lineStyle: {
          color: themeStore.isDark ? '#4b5563' : '#ddd'
        }
      }
    },
    yAxis: [
      {
        type: 'value',
        name: '审核数量',
        nameTextStyle: {
          color: themeStore.isDark ? '#9ca3af' : '#666'
        },
        axisLabel: {
          color: themeStore.isDark ? '#9ca3af' : '#666'
        },
        axisLine: {
          show: true,
          lineStyle: {
            color: themeStore.isDark ? '#4b5563' : '#ddd'
          }
        },
        splitLine: {
          lineStyle: {
            color: themeStore.isDark ? '#374151' : '#eee'
          }
        }
      }
    ],
    series: [
      {
        name: '通过',
        type: 'bar',
        barWidth: '30%',
        barGap: '0%',
        data: data.map(item => item.approved),
        itemStyle: {
          color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
            { offset: 0, color: '#10B981' },
            { offset: 1, color: '#34D399' }
          ]),
          borderRadius: [4, 4, 0, 0]
        }
      },
      {
        name: '拒绝',
        type: 'bar',
        barWidth: '30%',
        data: data.map(item => item.rejected),
        itemStyle: {
          color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
            { offset: 0, color: '#EF4444' },
            { offset: 1, color: '#F87171' }
          ]),
          borderRadius: [4, 4, 0, 0]
        }
      },
      {
        name: '总数',
        type: 'line',
        smooth: true,
        symbol: 'diamond',
        symbolSize: 8,
        lineStyle: {
          width: 3,
          color: '#8B5CF6'
        },
        itemStyle: {
          color: '#8B5CF6',
          borderColor: themeStore.isDark ? '#1f2937' : '#fff',
          borderWidth: 2
        },
        data: totalReviews
      }
    ]
  })
  
  window.addEventListener('resize', () => chart.resize())
}

// 评分分布
const initScoreDistChart = () => {
  if (!scoreDistChart.value) return
  
  const chart = echarts.init(scoreDistChart.value, null, {
    renderer: 'canvas'
  })
  chartInstances.push(chart)
  
  const data = statsData.value.trend_data.score_distribution || []
  
  chart.setOption({
    title: {
      text: '评分分布',
      left: 'center',
      textStyle: {
        color: themeStore.isDark ? '#e5e7eb' : '#333',
        fontWeight: 'normal'
      }
    },
    tooltip: {
      trigger: 'item',
      formatter: '{b}: {c} ({d}%)'
    },
    legend: {
      bottom: 0,
      itemWidth: 12,
      itemHeight: 12,
      textStyle: {
        color: themeStore.isDark ? '#e5e7eb' : '#333'
      }
    },
    series: [
      {
        name: '评分分布',
        type: 'pie',
        radius: ['30%', '80%'],
        center: ['50%', '40%'],
        roseType: 'area',
        itemStyle: {
          borderRadius: 6,
          borderColor: themeStore.isDark ? '#1f2937' : '#fff',
          borderWidth: 2
        },
        label: {
          formatter: '{b}: {c}',
          color: themeStore.isDark ? '#e5e7eb' : '#333'
        },
        data: data.map((item, index) => ({
          value: item.count,
          name: item.score_range,
          itemStyle: {
            color: colorPalette[index % colorPalette.length]
          }
        }))
      }
    ]
  })
  
  window.addEventListener('resize', () => chart.resize())
}

// 项目类型分布
const initProjectTypeChart = () => {
  if (!projectTypeChart.value) return
  
  const chart = echarts.init(projectTypeChart.value, null, {
    renderer: 'canvas'
  })
  chartInstances.push(chart)
  
  // 对数据进行排序
  const sortedData = [...projectTypeData.value].sort((a, b) => b.count - a.count)
  
  chart.setOption({
    title: {
      text: '项目类型分布',
      left: 'center',
      textStyle: {
        color: themeStore.isDark ? '#e5e7eb' : '#333',
        fontWeight: 'normal'
      }
    },
    tooltip: {
      trigger: 'axis',
      axisPointer: {
        type: 'shadow'
      }
    },
    grid: {
      left: '3%',
      right: '4%',
      bottom: '3%',
      top: '15%',
      containLabel: true
    },
    xAxis: {
      type: 'value',
      axisLabel: {
        color: themeStore.isDark ? '#9ca3af' : '#666'
      },
      axisLine: {
        lineStyle: {
          color: themeStore.isDark ? '#4b5563' : '#ddd'
        }
      },
      splitLine: {
        lineStyle: {
          color: themeStore.isDark ? '#374151' : '#eee'
        }
      }
    },
    yAxis: {
      type: 'category',
      data: sortedData.map(item => item.type || '未分类'),
      axisLabel: {
        color: themeStore.isDark ? '#9ca3af' : '#666'
      },
      axisLine: {
        lineStyle: {
          color: themeStore.isDark ? '#4b5563' : '#ddd'
        }
      }
    },
    series: [
      {
        name: '项目数量',
        type: 'bar',
        barWidth: '60%',
        data: sortedData.map((item, index) => ({
          value: item.count,
          itemStyle: {
            color: new echarts.graphic.LinearGradient(1, 0, 0, 0, [
              { offset: 0, color: colorPalette[index % colorPalette.length] },
              { offset: 1, color: echarts.color.lift(colorPalette[index % colorPalette.length], 0.3) }
            ]),
            borderRadius: [0, 4, 4, 0]
          }
        })),
        label: {
          show: true,
          position: 'right',
          formatter: '{c}',
          color: themeStore.isDark ? '#e5e7eb' : '#333'
        }
      }
    ]
  })
  
  window.addEventListener('resize', () => chart.resize())
}

// 处理深色模式切换
const handleThemeChange = () => {
  nextTick(() => {
    initCharts()
  })
}

// 监听主题变化
const themeWatcher = computed(() => themeStore.isDark)
watch(themeWatcher, handleThemeChange)

onMounted(() => {
  fetchStatistics()
})

onUnmounted(() => {
  // 销毁所有图表实例
  chartInstances.forEach(chart => {
    chart.dispose()
  })
  
  // 移除窗口resize监听
  chartInstances.forEach(chart => {
    window.removeEventListener('resize', () => chart.resize())
  })
})
</script>

<template>
  <div class="p-3 md:p-6">
    <h1 class="text-2xl font-bold mb-6 dark:text-gray-100">管理数据分析</h1>
    
    <SpinLoaderLarge v-if="loading" class="my-20" />
    
    <div v-else class="animate__animated animate__fadeIn">
      <!-- 统计卡片 -->
      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
        <!-- 项目统计 -->
        <div class="bg-white dark:bg-zinc-800 rounded-lg shadow-sm p-6 hover:shadow-md transition-shadow">
          <div class="flex justify-between items-start">
            <div>
              <p class="text-sm text-gray-500 dark:text-gray-400">项目总数</p>
              <p class="text-3xl font-bold text-blue-600 dark:text-blue-400">{{ statsData.project_stats.total_projects }}</p>
            </div>
            <div class="p-2 bg-blue-100 dark:bg-blue-900/30 rounded-lg">
              <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6 text-blue-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" d="M19.5 14.25v-2.625a3.375 3.375 0 00-3.375-3.375h-1.5A1.125 1.125 0 0113.5 7.125v-1.5a3.375 3.375 0 00-3.375-3.375H8.25m0 12.75h7.5m-7.5 3H12M10.5 2.25H5.625c-.621 0-1.125.504-1.125 1.125v17.25c0 .621.504 1.125 1.125 1.125h12.75c.621 0 1.125-.504 1.125-1.125V11.25a9 9 0 00-9-9z" />
              </svg>
            </div>
          </div>
          <div class="mt-4 grid grid-cols-3 gap-2">
            <div>
              <p class="text-xs text-gray-500 dark:text-gray-400">正常</p>
              <p class="text-sm font-medium text-green-600 dark:text-green-400">{{ statsData.project_stats.normal_projects }}</p>
            </div>
            <div>
              <p class="text-xs text-gray-500 dark:text-gray-400">下架</p>
              <p class="text-sm font-medium text-red-600 dark:text-red-400">{{ statsData.project_stats.banned_projects }}</p>
            </div>
            <div>
              <p class="text-xs text-gray-500 dark:text-gray-400">平均分</p>
              <p class="text-sm font-medium text-amber-600 dark:text-amber-400">{{ statsData.project_stats.avg_score.toFixed(1) }}</p>
            </div>
          </div>
        </div>
        
        <!-- 用户统计 -->
        <div class="bg-white dark:bg-zinc-800 rounded-lg shadow-sm p-6 hover:shadow-md transition-shadow">
          <div class="flex justify-between items-start">
            <div>
              <p class="text-sm text-gray-500 dark:text-gray-400">用户总数</p>
              <p class="text-3xl font-bold text-indigo-600 dark:text-indigo-400">{{ statsData.user_stats.total_users }}</p>
            </div>
            <div class="p-2 bg-indigo-100 dark:bg-indigo-900/30 rounded-lg">
              <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6 text-indigo-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" d="M15 19.128a9.38 9.38 0 002.625.372 9.337 9.337 0 004.121-.952 4.125 4.125 0 00-7.533-2.493M15 19.128v-.003c0-1.113-.285-2.16-.786-3.07M15 19.128v.106A12.318 12.318 0 018.624 21c-2.331 0-4.512-.645-6.374-1.766l-.001-.109a6.375 6.375 0 0111.964-3.07M12 6.375a3.375 3.375 0 11-6.75 0 3.375 3.375 0 016.75 0zm8.25 2.25a2.625 2.625 0 11-5.25 0 2.625 2.625 0 015.25 0z" />
              </svg>
            </div>
          </div>
          <div class="mt-4 grid grid-cols-3 gap-2">
            <div>
              <p class="text-xs text-gray-500 dark:text-gray-400">管理员</p>
              <p class="text-sm font-medium text-purple-600 dark:text-purple-400">{{ statsData.user_stats.admin_users }}</p>
            </div>
            <div>
              <p class="text-xs text-gray-500 dark:text-gray-400">正常</p>
              <p class="text-sm font-medium text-green-600 dark:text-green-400">{{ statsData.user_stats.normal_users }}</p>
            </div>
            <div>
              <p class="text-xs text-gray-500 dark:text-gray-400">封禁</p>
              <p class="text-sm font-medium text-red-600 dark:text-red-400">{{ statsData.user_stats.banned_users }}</p>
            </div>
          </div>
        </div>
        
        <!-- 审核统计 -->
        <div class="bg-white dark:bg-zinc-800 rounded-lg shadow-sm p-6 hover:shadow-md transition-shadow">
          <div class="flex justify-between items-start">
            <div>
              <p class="text-sm text-gray-500 dark:text-gray-400">审核总数</p>
              <p class="text-3xl font-bold text-green-600 dark:text-green-400">{{ statsData.review_stats.approved_reviews + statsData.review_stats.rejected_reviews }}</p>
            </div>
            <div class="p-2 bg-green-100 dark:bg-green-900/30 rounded-lg">
              <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6 text-green-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" d="M9 12.75L11.25 15 15 9.75M21 12c0 1.268-.63 2.39-1.593 3.068a3.745 3.745 0 01-1.043 3.296 3.745 3.745 0 01-3.296 1.043A3.745 3.745 0 0112 21c-1.268 0-2.39-.63-3.068-1.593a3.746 3.746 0 01-3.296-1.043 3.745 3.745 0 01-1.043-3.296A3.745 3.745 0 013 12c0-1.268.63-2.39 1.593-3.068a3.745 3.745 0 011.043-3.296 3.746 3.746 0 013.296-1.043A3.746 3.746 0 0112 3c1.268 0 2.39.63 3.068 1.593a3.746 3.746 0 013.296 1.043 3.746 3.746 0 011.043 3.296A3.745 3.745 0 0121 12z" />
              </svg>
            </div>
          </div>
          <div class="mt-4 grid grid-cols-3 gap-2">
            <div>
              <p class="text-xs text-gray-500 dark:text-gray-400">通过</p>
              <p class="text-sm font-medium text-green-600 dark:text-green-400">{{ statsData.review_stats.approved_reviews }}</p>
            </div>
            <div>
              <p class="text-xs text-gray-500 dark:text-gray-400">拒绝</p>
              <p class="text-sm font-medium text-red-600 dark:text-red-400">{{ statsData.review_stats.rejected_reviews }}</p>
            </div>
            <div>
              <p class="text-xs text-gray-500 dark:text-gray-400">待审核</p>
              <p class="text-sm font-medium text-amber-600 dark:text-amber-400">{{ statsData.review_stats.pending_reviews }}</p>
            </div>
          </div>
        </div>
        
        <!-- 交互统计 -->
        <div class="bg-white dark:bg-zinc-800 rounded-lg shadow-sm p-6 hover:shadow-md transition-shadow">
          <div class="flex justify-between items-start">
            <div>
              <p class="text-sm text-gray-500 dark:text-gray-400">总交互量</p>
              <p class="text-3xl font-bold text-pink-600 dark:text-pink-400">{{ statsData.project_stats.total_watches + statsData.project_stats.total_favorites }}</p>
            </div>
            <div class="p-2 bg-pink-100 dark:bg-pink-900/30 rounded-lg">
              <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6 text-pink-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" d="M15.042 21.672L13.684 16.6m0 0l-2.51 2.225.569-9.47 5.227 7.917-3.286-.672zM12 2.25V4.5m5.834.166l-1.591 1.591M20.25 10.5H18M7.757 14.743l-1.59 1.59M6 10.5H3.75m4.007-4.243l-1.59-1.59" />
              </svg>
            </div>
          </div>
          <div class="mt-4 grid grid-cols-2 gap-2">
            <div>
              <p class="text-xs text-gray-500 dark:text-gray-400">总浏览量</p>
              <p class="text-sm font-medium text-blue-600 dark:text-blue-400">{{ statsData.project_stats.total_watches }}</p>
            </div>
            <div>
              <p class="text-xs text-gray-500 dark:text-gray-400">总收藏量</p>
              <p class="text-sm font-medium text-amber-600 dark:text-amber-400">{{ statsData.project_stats.total_favorites }}</p>
            </div>
          </div>
        </div>
      </div>
      
      <!-- 热门项目信息 -->
      <div class="grid grid-cols-1 md:grid-cols-3 gap-6 mb-8">
        <div class="bg-white dark:bg-zinc-800 rounded-lg shadow-sm p-6 hover:shadow-md transition-shadow">
          <h3 class="text-lg font-semibold mb-4 text-gray-800 dark:text-gray-200">评分最高项目</h3>
          <div v-if="statsData.project_stats.highest_rated_id" class="flex flex-col">
            <p class="text-xl font-bold text-blue-600 dark:text-blue-400">{{ statsData.project_stats.highest_rated_name }}</p>
            <p class="mt-2 text-gray-600 dark:text-gray-400">项目ID: {{ statsData.project_stats.highest_rated_id }}</p>
            <div class="mt-2 bg-yellow-100 dark:bg-yellow-900/30 p-2 rounded-md">
              <p class="text-yellow-800 dark:text-yellow-200 font-semibold">
                平均评分: {{ statsData.project_stats.highest_rated_score.toFixed(1) }}
              </p>
            </div>
          </div>
          <p v-else class="text-gray-500 dark:text-gray-400">暂无数据</p>
        </div>
        
        <div class="bg-white dark:bg-zinc-800 rounded-lg shadow-sm p-6 hover:shadow-md transition-shadow">
          <h3 class="text-lg font-semibold mb-4 text-gray-800 dark:text-gray-200">浏览最多项目</h3>
          <div v-if="statsData.project_stats.most_watched_id" class="flex flex-col">
            <p class="text-xl font-bold text-blue-600 dark:text-blue-400">{{ statsData.project_stats.most_watched_name }}</p>
            <p class="mt-2 text-gray-600 dark:text-gray-400">项目ID: {{ statsData.project_stats.most_watched_id }}</p>
            <div class="mt-2 bg-blue-100 dark:bg-blue-900/30 p-2 rounded-md">
              <p class="text-blue-800 dark:text-blue-200 font-semibold">
                浏览量: {{ statsData.project_stats.most_watched_count }}
              </p>
            </div>
          </div>
          <p v-else class="text-gray-500 dark:text-gray-400">暂无数据</p>
        </div>
        
        <div class="bg-white dark:bg-zinc-800 rounded-lg shadow-sm p-6 hover:shadow-md transition-shadow">
          <h3 class="text-lg font-semibold mb-4 text-gray-800 dark:text-gray-200">收藏最多项目</h3>
          <div v-if="statsData.project_stats.most_favorited_id" class="flex flex-col">
            <p class="text-xl font-bold text-blue-600 dark:text-blue-400">{{ statsData.project_stats.most_favorited_name }}</p>
            <p class="mt-2 text-gray-600 dark:text-gray-400">项目ID: {{ statsData.project_stats.most_favorited_id }}</p>
            <div class="mt-2 bg-pink-100 dark:bg-pink-900/30 p-2 rounded-md">
              <p class="text-pink-800 dark:text-pink-200 font-semibold">
                收藏量: {{ statsData.project_stats.most_favorited_count }}
              </p>
            </div>
          </div>
          <p v-else class="text-gray-500 dark:text-gray-400">暂无数据</p>
        </div>
      </div>
      
      <!-- 图表区域 -->
      <div class="grid grid-cols-1 lg:grid-cols-2 gap-6 mb-8">
        <!-- 用户注册趋势 -->
        <div class="bg-white dark:bg-zinc-800 rounded-lg shadow-sm p-6 hover:shadow-md transition-shadow">
          <div ref="userTrendChart" style="width: 100%; height: 350px;"></div>
        </div>
        
        <!-- 项目创建趋势 -->
        <div class="bg-white dark:bg-zinc-800 rounded-lg shadow-sm p-6 hover:shadow-md transition-shadow">
          <div ref="projectTrendChart" style="width: 100%; height: 350px;"></div>
        </div>
        
        <!-- 章节审核趋势 -->
        <div class="bg-white dark:bg-zinc-800 rounded-lg shadow-sm p-6 hover:shadow-md transition-shadow">
          <div ref="reviewTrendChart" style="width: 100%; height: 350px;"></div>
        </div>
        
        <!-- 评分分布 -->
        <div class="bg-white dark:bg-zinc-800 rounded-lg shadow-sm p-6 hover:shadow-md transition-shadow">
          <div ref="scoreDistChart" style="width: 100%; height: 350px;"></div>
        </div>
      </div>
      
      <!-- 项目类型分布 -->
      <div class="bg-white dark:bg-zinc-800 rounded-lg shadow-sm p-6 mb-8 hover:shadow-md transition-shadow">
        <div ref="projectTypeChart" style="width: 100%; height: 400px;"></div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.animate__animated {
  animation-duration: 0.8s;
}

.animate__fadeIn {
  animation-name: fadeIn;
}

@keyframes fadeIn {
  from {
    opacity: 0;
  }
  to {
    opacity: 1;
  }
}

/* 自定义滚动条 */
::-webkit-scrollbar {
  width: 8px;
  height: 8px;
}

::-webkit-scrollbar-track {
  background: #f1f1f1;
  border-radius: 6px;
}

::-webkit-scrollbar-thumb {
  background: #c1c1c1;
  border-radius: 6px;
}

::-webkit-scrollbar-thumb:hover {
  background: #a8a8a8;
}

.dark ::-webkit-scrollbar-track {
  background: #374151;
}

.dark ::-webkit-scrollbar-thumb {
  background: #4b5563;
}

.dark ::-webkit-scrollbar-thumb:hover {
  background: #6b7280;
}
</style>
