<script setup>
import { ref, onMounted, nextTick, onUnmounted, computed, reactive, watch } from 'vue'
import { get } from '@/util/request'
import { message } from 'ant-design-vue'
import * as echarts from 'echarts'
import SpinLoaderLarge from '@/components/spinLoaderLarge.vue'
import { BACKEND_DOMAIN } from "@/util/VARRIBLES";
import { MdPreview } from "md-editor-v3"
import "md-editor-v3/lib/preview.css"
import { useThemeStore } from "@/stores/theme"

const themeStore = useThemeStore()

const loading = ref(true)
const projectData = ref([])
const styleData = ref([])
const emotionData = ref([])

// 分析建议相关的状态
const analysisOptions = reactive({
  isGenerating: false,
  showAnalysis: false,
  content: "",
})

// 将渐变色数组提取为公共变量
const gradientColors = [
  ['#3B82F6', '#60A5FA'], // 蓝色
  ['#10B981', '#34D399'], // 绿色
  ['#F59E0B', '#FBBF24'], // 橙色
  ['#8B5CF6', '#A78BFA'], // 紫色
  ['#EC4899', '#F472B6'], // 粉色
  ['#EF4444', '#F87171'], // 红色
  ['#06B6D4', '#22D3EE'], // 青色
  ['#14B8A6', '#2DD4BF'], // 蓝绿色
  ['#F97316', '#FB923C'], // 深橙色
  ['#84CC16', '#A3E635'], // 青柠色
  ['#6366F1', '#818CF8'], // 靛蓝色
  ['#D946EF', '#E879F9'], // 洋红色
  ['#0EA5E9', '#38BDF8'], // 天蓝色
  ['#4F46E5', '#6366F1'], // 深蓝色
  ['#7C3AED', '#A78BFA'], // 深紫色
  ['#DB2777', '#EC4899'], // 玫红色
  ['#EA580C', '#FB923C'], // 赭石色
  ['#16A34A', '#4ADE80'], // 翠绿色
  ['#2563EB', '#60A5FA'], // 皇家蓝
  ['#9333EA', '#C084FC']  // 紫罗兰色
]

// 公共配置项
const commonOptions = {
  tooltip: {
    trigger: 'axis',
    backgroundColor: 'rgba(255, 255, 255, 0.9)',
    borderColor: '#eee',
    borderWidth: 1,
    textStyle: {
      color: '#666'
    },
    axisPointer: {
      type: 'line',
      lineStyle: {
        color: '#ccc',
        width: 1
      }
    }
  },
  legend: {
    bottom: 0,
    type: 'scroll',
    pageIconColor: '#666',
    pageTextStyle: {
      color: '#666'
    },
    textStyle: {
      color: '#666',
      fontSize: 12
    },
    itemWidth: 12,
    itemHeight: 12,
    itemGap: 20
  },
  grid: {
    left: '3%',
    right: '4%',
    bottom: '15%',
    top: '15%',
    containLabel: true
  },
  xAxis: {
    type: 'category',
    boundaryGap: false,
    axisLabel: {
      rotate: 45,
      formatter: (value) => value.split('T')[0],
      color: '#666',
      fontSize: 12
    },
    axisLine: {
      lineStyle: {
        color: '#eee'
      }
    },
    axisTick: {
      show: false
    }
  },
  yAxis: {
    type: 'value',
    axisLabel: {
      color: '#666',
      fontSize: 12,
      formatter: (value) => value.toLocaleString()
    },
    splitLine: {
      lineStyle: {
        type: 'dashed',
        color: '#eee'
      }
    },
    axisTick: {
      show: false
    },
    axisLine: {
      show: false
    }
  },
  animation: true,
  animationDuration: 1000,
  animationEasing: 'cubicOut'
}

// 在 styleData 和 emotionData 的 ref 声明后添加新的计算属性
const styleAnalytics = computed(() => {
  if (!styleData.value.length) return null
  
  return {
    // 按风格统计总观看和点赞数
    styleStats: Object.entries(
      styleData.value.reduce((acc, curr) => {
        if (!acc[curr.style]) {
          acc[curr.style] = { watches: 0, likes: 0, count: 0 }
        }
        acc[curr.style].watches += curr.watch_count
        acc[curr.style].likes += curr.like_count
        acc[curr.style].count++
        return acc
      }, {})
    ),
    
    // 按类型统计
    typeStats: Object.entries(
      styleData.value.reduce((acc, curr) => {
        if (!acc[curr.type]) {
          acc[curr.type] = { watches: 0, likes: 0, count: 0 }
        }
        acc[curr.type].watches += curr.watch_count
        acc[curr.type].likes += curr.like_count
        acc[curr.type].count++
        return acc
      }, {})
    ),
    
    // 计算互动率 (观看+点赞)/作品数
    engagementRate: Object.entries(
      styleData.value.reduce((acc, curr) => {
        if (!acc[curr.style]) {
          acc[curr.style] = { total: 0, count: 0 }
        }
        acc[curr.style].total += (curr.watch_count + curr.like_count)
        acc[curr.style].count++
        return acc
      }, {})
    ).map(([style, data]) => ({
      style,
      rate: data.total / data.count
    }))
  }
})

const emotionAnalytics = computed(() => {
  if (!emotionData.value.length) return null
  
  return {
    // 按项目统计情绪分布
    projectEmotions: Object.entries(
      emotionData.value.reduce((acc, curr) => {
        if (!acc[curr.project_name]) {
          acc[curr.project_name] = {}
        }
        if (!acc[curr.project_name][curr.emotion]) {
          acc[curr.project_name][curr.emotion] = 0
        }
        acc[curr.project_name][curr.emotion] += curr.count
        return acc
      }, {})
    ),
    
    // 情绪变化趋势
    emotionTrends: emotionData.value.reduce((acc, curr) => {
      const date = curr.date.split('T')[0]
      if (!acc[date]) {
        acc[date] = {}
      }
      if (!acc[date][curr.emotion]) {
        acc[date][curr.emotion] = 0
      }
      acc[date][curr.emotion] += curr.count
      return acc
    }, {})
  }
})

// 添加新的计算属性用于数据概览
const dataOverview = computed(() => {
  if (!projectData.value.length) return null
  
  // 计算总观看和收藏
  const totalWatches = projectData.value.reduce((sum, item) => sum + item.watch_count, 0)
  const totalFavorites = projectData.value.reduce((sum, item) => sum + item.favorite_count, 0)
  
  // 计算互动率
  const interactionRate = totalWatches > 0 ? 
    ((totalFavorites / totalWatches) * 100).toFixed(1) : 0
  
  // 计算最受欢迎的项目
  const projectStats = projectData.value.reduce((acc, curr) => {
    if (!acc[curr.project_name]) {
      acc[curr.project_name] = { watches: 0, favorites: 0 }
    }
    acc[curr.project_name].watches += curr.watch_count
    acc[curr.project_name].favorites += curr.favorite_count
    return acc
  }, {})
  
  const popularProject = Object.entries(projectStats)
    .sort((a, b) => (b[1].watches + b[1].favorites) - (a[1].watches + a[1].favorites))
    .map(([name, stats]) => ({ name, ...stats }))[0] || { name: '暂无数据', watches: 0, favorites: 0 }
  
  // 计算最受欢迎的风格
  let popularStyle = { style: '暂无数据', count: 0 }
  if (styleData.value.length) {
    const styleStats = styleData.value.reduce((acc, curr) => {
      if (!acc[curr.style]) {
        acc[curr.style] = 0
      }
      acc[curr.style] += curr.watch_count + curr.like_count
      return acc
    }, {})
    
    const topStyle = Object.entries(styleStats)
      .sort((a, b) => b[1] - a[1])[0]
    
    if (topStyle) {
      popularStyle = { style: topStyle[0], count: topStyle[1] }
    }
  }
  
  // 计算主要情绪
  let mainEmotion = { emotion: '暂无数据', count: 0 }
  if (emotionData.value.length) {
    const emotionStats = emotionData.value.reduce((acc, curr) => {
      if (!acc[curr.emotion]) {
        acc[curr.emotion] = 0
      }
      acc[curr.emotion] += curr.count
      return acc
    }, {})
    
    const topEmotion = Object.entries(emotionStats)
      .sort((a, b) => b[1] - a[1])[0]
    
    if (topEmotion) {
      mainEmotion = { emotion: topEmotion[0], count: topEmotion[1] }
    }
  }
  
  // 计算近期趋势
  const recentDates = [...new Set(projectData.value.map(item => item.date))]
    .sort()
    .slice(-5)
  
  const recentTrend = recentDates.map(date => {
    const dayData = projectData.value.filter(item => item.date === date)
    return {
      date: date.split('T')[0],
      watches: dayData.reduce((sum, item) => sum + item.watch_count, 0),
      favorites: dayData.reduce((sum, item) => sum + item.favorite_count, 0)
    }
  })
  
  // 计算增长率
  let growthRate = 0
  if (recentTrend.length >= 2) {
    const oldValue = recentTrend[0].watches
    const newValue = recentTrend[recentTrend.length - 1].watches
    growthRate = oldValue > 0 ? 
      (((newValue - oldValue) / oldValue) * 100).toFixed(1) : 0
  }
  
  return {
    totalProjects: [...new Set(projectData.value.map(item => item.project_name))].length,
    totalWatches,
    totalFavorites,
    interactionRate,
    popularProject,
    popularStyle,
    mainEmotion,
    recentTrend,
    growthRate
  }
})

// 初始化图表
const initCharts = () => {
  const watchChart = echarts.init(document.getElementById('watchChart'))
  const favoriteChart = echarts.init(document.getElementById('favoriteChart'))

  // 处理数据
  const projects = [...new Set(projectData.value.map(item => item.project_name))]
  const dates = [...new Set(projectData.value.map(item => item.date))].sort()

  // 观看数据
  const watchSeries = projects.map((project, index) => ({
    name: project,
    type: 'line',
    smooth: true,
    symbolSize: 8,
    showSymbol: false, // 默认不显示标记点
    emphasis: {
      focus: 'series',
      showSymbol: true // 鼠标经过时显示标记点
    },
    lineStyle: {
      width: 4,
      shadowColor: 'rgba(0,0,0,0.2)',
      shadowBlur: 10
    },
    areaStyle: {
      opacity: 0.1,
      color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
        {
          offset: 0,
          color: gradientColors[index % gradientColors.length][0]
        },
        {
          offset: 1,
          color: gradientColors[index % gradientColors.length][1]
        }
      ])
    },
    itemStyle: {
      color: gradientColors[index % gradientColors.length][0]
    },
    data: dates.map(date => {
      const record = projectData.value.find(item => 
        item.project_name === project && item.date === date
      )
      return record ? record.watch_count : 0
    })
  }))

  // 收藏数据配置类似
  const favoriteSeries = projects.map((project, index) => ({
    name: project,
    type: 'line',
    smooth: true,
    symbolSize: 8,
    showSymbol: false,
    emphasis: {
      focus: 'series',
      showSymbol: true
    },
    lineStyle: {
      width: 4,
      shadowColor: 'rgba(0,0,0,0.2)',
      shadowBlur: 10
    },
    areaStyle: {
      opacity: 0.1,
      color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
        {
          offset: 0,
          color: gradientColors[index % gradientColors.length][0]
        },
        {
          offset: 1,
          color: gradientColors[index % gradientColors.length][1]
        }
      ])
    },
    itemStyle: {
      color: gradientColors[index % gradientColors.length][0]
    },
    data: dates.map(date => {
      const record = projectData.value.find(item => 
        item.project_name === project && item.date === date
      )
      return record ? record.favorite_count : 0
    })
  }))

  // 计算所有项目的总观看和总收藏数
  const totalStats = {
    totalProjects: projects.length,
    totalWatches: projectData.value.reduce((sum, item) => sum + item.watch_count, 0),
    totalFavorites: projectData.value.reduce((sum, item) => sum + item.favorite_count, 0)
  }

  // 配置项
  const commonOptions = {
    tooltip: {
      trigger: 'axis',
      backgroundColor: 'rgba(255, 255, 255, 0.9)',
      borderColor: '#eee',
      borderWidth: 1,
      textStyle: {
        color: '#666'
      },
      axisPointer: {
        type: 'line',
        lineStyle: {
          color: '#ccc',
          width: 1
        }
      },
      formatter: (params) => {
        let result = `<div class="font-medium text-gray-600">${params[0].axisValue}</div>`
        params.forEach(param => {
          result += `
            <div class="flex items-center justify-between gap-4 mt-2">
              <span style="color:${param.color}">${param.seriesName}</span>
              <span class="font-semibold">${param.value}</span>
            </div>`
        })
        return result
      }
    },
    legend: {
      data: projects,
      bottom: 0,
      type: 'scroll',
      pageIconColor: '#666',
      pageTextStyle: {
        color: '#666'
      },
      textStyle: {
        color: '#666',
        fontSize: 12
      },
      itemWidth: 12,
      itemHeight: 12,
      itemGap: 20
    },
    grid: {
      left: '3%',
      right: '4%',
      bottom: '15%',
      top: '15%',
      containLabel: true
    },
    xAxis: {
      type: 'category',
      data: dates,
      boundaryGap: false,
      axisLabel: {
        rotate: 45,
        formatter: (value) => value.split('T')[0],
        color: '#666',
        fontSize: 12
      },
      axisLine: {
        lineStyle: {
          color: '#eee'
        }
      },
      axisTick: {
        show: false
      }
    },
    yAxis: {
      type: 'value',
      axisLabel: {
        color: '#666',
        fontSize: 12,
        formatter: (value) => value.toLocaleString()
      },
      splitLine: {
        lineStyle: {
          type: 'dashed',
          color: '#eee'
        }
      },
      axisTick: {
        show: false
      },
      axisLine: {
        show: false
      }
    },
    animation: true,
    animationDuration: 1000,
    animationEasing: 'cubicOut'
  }

  // 观看图表配置
  watchChart.setOption({
    ...commonOptions,
    title: {
      text: '近10天项目观看趋势',
      left: 'center',
      top: 0,
      textStyle: {
        fontSize: 16,
        fontWeight: 'bold',
        color: document.documentElement.classList.contains('dark') ? '#e5e7eb' : '#333'
      },
      subtextStyle: {
        color: document.documentElement.classList.contains('dark') ? '#9ca3af' : '#666',
        fontSize: 12
      }
    },
    series: watchSeries
  })

  // 收藏图表配置
  favoriteChart.setOption({
    ...commonOptions,
    title: {
      text: '近10天项目收藏趋势',
      left: 'center',
      top: 0,
      textStyle: {
        fontSize: 16,
        fontWeight: 'bold',
        color: document.documentElement.classList.contains('dark') ? '#e5e7eb' : '#333'
      },
      subtextStyle: {
        color: document.documentElement.classList.contains('dark') ? '#9ca3af' : '#666',
        fontSize: 12
      }
    },
    series: favoriteSeries
  })

  // 添加暗色模式切换监听
  const observer = new MutationObserver(() => {
    const isDark = document.documentElement.classList.contains('dark')
    watchChart.setOption({
      title: {
        textStyle: {
          color: isDark ? '#e5e7eb' : '#333'
        },
        subtextStyle: {
          color: isDark ? '#9ca3af' : '#666'
        }
      }
    })
    favoriteChart.setOption({
      title: {
        textStyle: {
          color: isDark ? '#e5e7eb' : '#333'
        },
        subtextStyle: {
          color: isDark ? '#9ca3af' : '#666'
        }
      }
    })
  })

  observer.observe(document.documentElement, {
    attributes: true,
    attributeFilter: ['class']
  })

  // 在组件卸载时移除监听器
  onUnmounted(() => {
    observer.disconnect()
  })

  // 监听窗口大小变化
  window.addEventListener('resize', () => {
    watchChart.resize()
    favoriteChart.resize()
  })

  return { totalStats }
}

// 获取数据
const fetchData = async () => {
  try {
    get('/api/project/analysis/watches-likes', {}, (messageer, data) => {
      projectData.value = data
      loading.value = false
      nextTick(() => {
        const { totalStats } = initCharts()
        projectStats.value = totalStats
      })
    })
  } catch (error) {
    message.error('获取数据失败')
    loading.value = false
  }
}

const projectStats = ref({})

// 初始化风格和类型图表
const initStyleAndTypeCharts = () => {
  const styleChart = echarts.init(document.getElementById('styleChart'))
  const typeChart = echarts.init(document.getElementById('typeChart'))
  
  const dates = [...new Set(styleData.value.map(item => item.date))].sort()
  const styles = [...new Set(styleData.value.map(item => item.style))]
  const types = [...new Set(styleData.value.map(item => item.type))]

  // 风格数据系列
  const styleSeries = styles.map((style, index) => {
    // 观看数据系列
    const watchSeries = {
      name: `${style}(观看)`,
      type: 'line',
      smooth: true,
      symbolSize: 8,
      showSymbol: false,
      emphasis: {
        focus: 'series',
        showSymbol: true
      },
      lineStyle: {
        width: 4,
        shadowColor: 'rgba(0,0,0,0.2)',
        shadowBlur: 10
      },
      areaStyle: {
        opacity: 0.1,
        color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
          {
            offset: 0,
            color: gradientColors[index % gradientColors.length][0]
          },
          {
            offset: 1,
            color: gradientColors[index % gradientColors.length][1]
          }
        ])
      },
      itemStyle: {
        color: gradientColors[index % gradientColors.length][0]
      },
      data: dates.map(date => {
        const records = styleData.value.filter(item => 
          item.style === style && 
          item.date === date
        )
        return records.reduce((sum, record) => sum + record.watch_count, 0)
      })
    }

    // 收藏数据系列
    const likeSeries = {
      name: `${style}(收藏)`,
      type: 'line',
      smooth: true,
      symbolSize: 8,
      showSymbol: false,
      emphasis: {
        focus: 'series',
        showSymbol: true
      },
      lineStyle: {
        width: 4,
        shadowColor: 'rgba(0,0,0,0.2)',
        shadowBlur: 10
      },
      areaStyle: {
        opacity: 0.1,
        color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
          {
            offset: 0,
            color: gradientColors[index % gradientColors.length][0]
          },
          {
            offset: 1,
            color: gradientColors[index % gradientColors.length][1]
          }
        ])
      },
      itemStyle: {
        color: gradientColors[index % gradientColors.length][1]
      },
      data: dates.map(date => {
        const records = styleData.value.filter(item => 
          item.style === style && 
          item.date === date
        )
        return records.reduce((sum, record) => sum + record.like_count, 0)
      })
    }

    return [watchSeries, likeSeries]
  }).flat()

  // 类型数据系列
  const typeSeries = types.map((type, index) => {
    // 观看数据系列
    const watchSeries = {
      name: `${type}(观看)`,
      type: 'line',
      smooth: true,
      symbolSize: 8,
      showSymbol: false,
      emphasis: {
        focus: 'series',
        showSymbol: true
      },
      lineStyle: {
        width: 4,
        shadowColor: 'rgba(0,0,0,0.2)',
        shadowBlur: 10
      },
      areaStyle: {
        opacity: 0.1,
        color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
          {
            offset: 0,
            color: gradientColors[index % gradientColors.length][0]
          },
          {
            offset: 1,
            color: gradientColors[index % gradientColors.length][1]
          }
        ])
      },
      itemStyle: {
        color: gradientColors[index % gradientColors.length][0]
      },
      data: dates.map(date => {
        const records = styleData.value.filter(item => 
          item.type === type && 
          item.date === date
        )
        return records.reduce((sum, record) => sum + record.watch_count, 0)
      })
    }

    // 收藏数据系列
    const likeSeries = {
      name: `${type}(收藏)`,
      type: 'line',
      smooth: true,
      symbolSize: 8,
      showSymbol: false,
      emphasis: {
        focus: 'series',
        showSymbol: true
      },
      lineStyle: {
        width: 4,
        shadowColor: 'rgba(0,0,0,0.2)',
        shadowBlur: 10
      },
      areaStyle: {
        opacity: 0.1,
        color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
          {
            offset: 0,
            color: gradientColors[index % gradientColors.length][0]
          },
          {
            offset: 1,
            color: gradientColors[index % gradientColors.length][1]
          }
        ])
      },
      itemStyle: {
        color: gradientColors[index % gradientColors.length][1]
      },
      data: dates.map(date => {
        const records = styleData.value.filter(item => 
          item.type === type && 
          item.date === date
        )
        return records.reduce((sum, record) => sum + record.like_count, 0)
      })
    }

    return [watchSeries, likeSeries]
  }).flat()

  // 设置风格图表配置
  styleChart.setOption({
    ...commonOptions,
    title: {
      text: '项目风格数据趋势',
      subtext: '按天统计观看和收藏数',
      left: 'center',
      top: 0,
      textStyle: {
        fontSize: 16,
        fontWeight: 'bold',
        color: document.documentElement.classList.contains('dark') ? '#e5e7eb' : '#333'
      },
      subtextStyle: {
        color: document.documentElement.classList.contains('dark') ? '#9ca3af' : '#666'
      }
    },
    xAxis: {
      type: 'category',
      data: dates,
      boundaryGap: false,
      axisLabel: {
        formatter: (value) => value.split('T')[0],
        rotate: 45,
        color: '#666',
        fontSize: 12
      },
      axisLine: {
        lineStyle: {
          color: '#eee'
        }
      },
      axisTick: {
        show: false
      }
    },
    yAxis: {
      type: 'value',
      axisLabel: {
        color: '#666',
        fontSize: 12,
        formatter: (value) => value.toLocaleString()
      },
      splitLine: {
        lineStyle: {
          type: 'dashed',
          color: '#eee'
        }
      },
      axisTick: {
        show: false
      },
      axisLine: {
        show: false
      }
    },
    tooltip: {
      trigger: 'axis',
      formatter: (params) => {
        let result = `<div class="font-medium text-gray-600">${params[0].axisValue.split('T')[0]}</div>`
        
        // 按风格分组显示数据
        const styleData = {}
        params.forEach(param => {
          const [style, metric] = param.seriesName.split('(')
          if (!styleData[style]) styleData[style] = {}
          styleData[style][metric.slice(0, -1)] = param.value
        })

        Object.entries(styleData).forEach(([style, data]) => {
          result += `
            <div class="flex justify-between gap-4 mt-1">
              <span>${style}</span>
              <span>观看: ${data.观看 || 0} / 收藏: ${data.收藏 || 0}</span>
            </div>`
        })

        return result
      }
    },
    series: styleSeries
  })

  // 设置类型图表配置
  typeChart.setOption({
    ...commonOptions,
    title: {
      text: '项目类型数据趋势',
      subtext: '按天统计观看和收藏数',
      left: 'center',
      top: 0,
      textStyle: {
        fontSize: 16,
        fontWeight: 'bold',
        color: document.documentElement.classList.contains('dark') ? '#e5e7eb' : '#333'
      },
      subtextStyle: {
        color: document.documentElement.classList.contains('dark') ? '#9ca3af' : '#666'
      }
    },
    xAxis: {
      type: 'category',
      data: dates,
      boundaryGap: false,
      axisLabel: {
        formatter: (value) => value.split('T')[0],
        rotate: 45,
        color: '#666',
        fontSize: 12
      },
      axisLine: {
        lineStyle: {
          color: '#eee'
        }
      },
      axisTick: {
        show: false
      }
    },
    yAxis: {
      type: 'value',
      axisLabel: {
        color: '#666',
        fontSize: 12,
        formatter: (value) => value.toLocaleString()
      },
      splitLine: {
        lineStyle: {
          type: 'dashed',
          color: '#eee'
        }
      },
      axisTick: {
        show: false
      },
      axisLine: {
        show: false
      }
    },
    tooltip: {
      trigger: 'axis',
      formatter: (params) => {
        let result = `<div class="font-medium text-gray-600">${params[0].axisValue.split('T')[0]}</div>`
        
        // 按类型分组显示数据
        const typeData = {}
        params.forEach(param => {
          const [type, metric] = param.seriesName.split('(')
          if (!typeData[type]) typeData[type] = {}
          typeData[type][metric.slice(0, -1)] = param.value
        })

        Object.entries(typeData).forEach(([type, data]) => {
          result += `
            <div class="flex justify-between gap-4 mt-1">
              <span>${type}</span>
              <span>观看: ${data.观看 || 0} / 收藏: ${data.收藏 || 0}</span>
            </div>`
        })

        return result
      }
    },
    series: typeSeries
  })

  // 监听窗口大小变化
  window.addEventListener('resize', () => {
    styleChart.resize()
    typeChart.resize()
  })
}

// 修改 fetchStyleData 函数
const fetchStyleData = async () => {
  try {
    get('/api/project/analysis/style-type', {}, (messageer, data) => {
      styleData.value = data
      nextTick(() => {
        initStyleAndTypeCharts()
      })
    })
  } catch (error) {
    message.error('获取数据失败')
  }
}

// 获取情绪分析数据
const fetchEmotionData = async () => {
  try {
    get('/api/user/analysis/emotions', {}, (messageer, data) => {
      emotionData.value = data
      nextTick(() => {
        initEmotionChart()
      })
    })
  } catch (error) {
    message.error('获取情绪分析数据失败')
  }
}

// 初始化情绪分析图表
const initEmotionChart = () => {
  const emotionChart = echarts.init(document.getElementById('emotionChart'))
  
  const projects = [...new Set(emotionData.value.map(item => item.project_name))]
  const dates = [...new Set(emotionData.value.map(item => item.date))].sort()
  const emotions = ["喜悦", "感动", "惊喜", "期待", "伤感", "愤怒", "恐惧", "平静"]
  
  // 准备数据系列
  const series = projects.map((project, index) => {
    return emotions.map(emotion => ({
      name: `${project}-${emotion}`,
      type: 'bar',
      stack: project,
      emphasis: {
        focus: 'series'
      },
      itemStyle: {
        color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
          {
            offset: 0,
            color: gradientColors[index % gradientColors.length][0]
          },
          {
            offset: 1,
            color: gradientColors[index % gradientColors.length][1]
          }
        ])
      },
      data: dates.map(date => {
        const record = emotionData.value.find(item => 
          item.project_name === project && 
          item.date === date && 
          item.emotion === emotion
        )
        return record ? record.count : 0
      })
    }))
  }).flat()

  emotionChart.setOption({
    title: {
      text: '项目情绪分布趋势',
      subtext: '近10天用户情绪统计',
      left: 'center',
      top: 0,
      textStyle: {
        fontSize: 16,
        fontWeight: 'bold',
        color: document.documentElement.classList.contains('dark') ? '#e5e7eb' : '#333'
      },
      subtextStyle: {
        color: document.documentElement.classList.contains('dark') ? '#9ca3af' : '#666'
      }
    },
    tooltip: {
      trigger: 'axis',
      axisPointer: {
        type: 'shadow'
      },
      formatter: (params) => {
        let result = `<div class="font-medium text-gray-600">${params[0].axisValue}</div>`
        const projectGroups = {}
        
        params.forEach(param => {
          const [project, emotion] = param.seriesName.split('-')
          if (!projectGroups[project]) {
            projectGroups[project] = []
          }
          if (param.value > 0) {
            projectGroups[project].push({
              emotion,
              count: param.value,
              color: param.color
            })
          }
        })
        
        Object.entries(projectGroups).forEach(([project, emotions]) => {
          if (emotions.length > 0) {
            result += `<div class="mt-2 mb-1 font-medium">${project}</div>`
            emotions.forEach(({ emotion, count, color }) => {
              result += `
                <div class="flex items-center justify-between gap-4">
                  <span style="color:${color}">${emotion}</span>
                  <span class="font-semibold">${count}</span>
                </div>`
            })
          }
        })
        
        return result
      }
    },
    legend: {
      type: 'scroll',
      bottom: 0,
      height: 80,
      formatter: (name) => {
        const [project, emotion] = name.split('-')
        return `${project}\n${emotion}`
      }
    },
    grid: {
      left: '3%',
      right: '4%',
      bottom: 100,
      top: 60,
      containLabel: true
    },
    xAxis: {
      type: 'category',
      data: dates,
      axisLabel: {
        formatter: (value) => value.split('T')[0],
        rotate: 45
      }
    },
    yAxis: {
      type: 'value',
      name: '情绪计数'
    },
    series
  })

  // 监听窗口大小变化
  window.addEventListener('resize', () => {
    emotionChart.resize()
  })
}

// 生成分析建议
const generateAnalysis = () => {
  if (!styleData.value || styleData.value.length === 0) {
    message.warning("暂无数据可供分析")
    return
  }

  analysisOptions.showAnalysis = true
  analysisOptions.isGenerating = true
  analysisOptions.content = ""

  // 创建 WebSocket 连接
  const token = localStorage.getItem("authToken")
  const baseUrl = BACKEND_DOMAIN.replace(/^http/, 'ws').replace(/\/$/, '')
  const ws = new WebSocket(
    `${baseUrl}/ws/newProjectAnalysis`
  )

  ws.onopen = () => {
    ws.send(
      JSON.stringify({
        data: styleData.value,
        auth_token: token
      })
    )
  }

  ws.onmessage = (event) => {
    const response = JSON.parse(event.data)
    if (response.code === 500) {
      message.error(response.message)
      analysisOptions.isGenerating = false
      ws.close()
      return
    }

    if (response.done) {
      ws.close()
      return
    }

    analysisOptions.content += response.content
    // 滚动到底部
    nextTick(() => {
      const element = document.querySelector('.prose')
      if (element) {
        element.scrollTop = element.scrollHeight
      }
    })
  }

  ws.onerror = (error) => {
    console.error("WebSocket error:", error)
    message.error("连接发生错误，请重试")
    analysisOptions.isGenerating = false
  }

  ws.onclose = () => {
    analysisOptions.isGenerating = false
  }
}

// 修改 initAdvancedCharts 函数
const initAdvancedCharts = () => {
  // 确保 DOM 元素存在
  const engagementChartEl = document.getElementById('engagementChart')
  const emotionHeatmapEl = document.getElementById('emotionHeatmap')
  
  if (!engagementChartEl || !emotionHeatmapEl || !styleAnalytics.value || !emotionAnalytics.value) {
    return
  }

  // 互动率雷达图
  const engagementChart = echarts.init(engagementChartEl)
  const { engagementRate } = styleAnalytics.value
  
  engagementChart.setOption({
    title: {
      text: '内容风格互动率分析',
      left: 'center',
      textStyle: {
        fontSize: 16,
        fontWeight: 'bold',
        color: document.documentElement.classList.contains('dark') ? '#e5e7eb' : '#333'
      }
    },
    radar: {
      indicator: engagementRate.map(item => ({
        name: item.style,
        max: Math.max(...engagementRate.map(i => i.rate)) * 1.2
      })),
      splitArea: {
        areaStyle: {
          color: document.documentElement.classList.contains('dark') 
            ? ['rgba(255,255,255,0.02)', 'rgba(255,255,255,0.05)']
            : ['rgba(200,200,200,0.1)', 'rgba(250,250,250,0.3)']
        }
      },
      axisLine: {
        lineStyle: {
          color: document.documentElement.classList.contains('dark') ? 'rgba(255,255,255,0.1)' : 'rgba(0,0,0,0.1)'
        }
      },
      splitLine: {
        lineStyle: {
          color: document.documentElement.classList.contains('dark') ? 'rgba(255,255,255,0.1)' : 'rgba(0,0,0,0.1)'
        }
      },
      name: {
        textStyle: {
          color: document.documentElement.classList.contains('dark') ? '#e5e7eb' : '#666'
        }
      }
    },
    series: [{
      type: 'radar',
      data: [{
        value: engagementRate.map(item => item.rate),
        name: '互动率',
        symbol: 'circle',
        symbolSize: 8,
        lineStyle: {
          width: 3,
          color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
            { offset: 0, color: '#3B82F6' },
            { offset: 1, color: '#60A5FA' }
          ])
        },
        areaStyle: {
          color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
            { offset: 0, color: 'rgba(59, 130, 246, 0.3)' },
            { offset: 1, color: 'rgba(96, 165, 250, 0.1)' }
          ])
        },
        itemStyle: {
          color: '#3B82F6'
        }
      }]
    }]
  })

  // 情绪分布热力图
  const emotionHeatmap = echarts.init(emotionHeatmapEl)
  const { projectEmotions } = emotionAnalytics.value
  
  const projects = projectEmotions.map(([name]) => name)
  const emotions = ["喜悦", "感动", "惊喜", "期待", "伤感", "愤怒", "恐惧", "平静"]
  const heatmapData = projects.flatMap((project, i) => 
    emotions.map((emotion, j) => [
      i,
      j,
      projectEmotions.find(([name]) => name === project)[1][emotion] || 0
    ])
  )

  emotionHeatmap.setOption({
    title: {
      text: '项目情绪分布热力图',
      left: 'center',
      textStyle: {
        fontSize: 16,
        fontWeight: 'bold',
        color: document.documentElement.classList.contains('dark') ? '#e5e7eb' : '#333'
      }
    },
    tooltip: {
      position: 'top',
      formatter: function(params) {
        return `${params.name}<br/>
                ${params.marker}${emotions[params.data[1]]}: ${params.data[2]}`
      },
      backgroundColor: 'rgba(255, 255, 255, 0.9)',
      borderColor: '#eee',
      borderWidth: 1,
      textStyle: {
        color: '#666'
      }
    },
    grid: {
      top: '15%',
      bottom: '15%',
      left: '10%',
      right: '10%'
    },
    xAxis: {
      type: 'category',
      data: projects,
      splitArea: {
        show: true
      },
      axisLabel: {
        rotate: 45,
        color: document.documentElement.classList.contains('dark') ? '#e5e7eb' : '#666',
        fontSize: 12
      },
      axisLine: {
        lineStyle: {
          color: document.documentElement.classList.contains('dark') ? 'rgba(255,255,255,0.1)' : '#eee'
        }
      }
    },
    yAxis: {
      type: 'category',
      data: emotions,
      splitArea: {
        show: true
      },
      axisLabel: {
        color: document.documentElement.classList.contains('dark') ? '#e5e7eb' : '#666',
        fontSize: 12
      },
      axisLine: {
        lineStyle: {
          color: document.documentElement.classList.contains('dark') ? 'rgba(255,255,255,0.1)' : '#eee'
        }
      }
    },
    visualMap: {
      min: 0,
      max: Math.max(...heatmapData.map(item => item[2])),
      calculable: true,
      orient: 'horizontal',
      left: 'center',
      bottom: '5%',
      textStyle: {
        color: document.documentElement.classList.contains('dark') ? '#e5e7eb' : '#666'
      },
      inRange: {
        color: document.documentElement.classList.contains('dark')
          ? [
              'rgba(255,255,255,0.05)',  // 最少 - 深色模式下接近透明
              '#6B7280',  // 中低
              '#4B5563',  // 中等
              '#3B82F6',  // 中高
              '#2563EB'   // 最多 - 深蓝色
            ]
          : [
              'rgba(0,0,0,0.05)',  // 最少 - 浅色模式下接近透明
              '#E5E7EB',  // 中低
              '#9CA3AF',  // 中等
              '#3B82F6',  // 中高
              '#1D4ED8'   // 最多 - 深蓝色
            ]
      }
    },
    series: [{
      name: '情绪分布',
      type: 'heatmap',
      data: heatmapData,
      label: {
        show: true,
        color: (params) => {
          // 根据数值动态调整文字颜色，确保在深色背景上清晰可见
          const value = params.data[2]
          const max = Math.max(...heatmapData.map(item => item[2]))
          const threshold = max * 0.5
          return value > threshold 
            ? '#ffffff'  // 深色背景上用白色文字
            : (document.documentElement.classList.contains('dark') ? '#e5e7eb' : '#666')
        },
        formatter: (params) => {
          return params.data[2] > 0 ? params.data[2] : ''
        }
      },
      emphasis: {
        itemStyle: {
          shadowBlur: 10,
          shadowColor: 'rgba(0, 0, 0, 0.5)'
        }
      }
    }]
  })

  // 监听窗口大小变化
  window.addEventListener('resize', () => {
    engagementChart.resize()
    emotionHeatmap.resize()
  })

  // 在组件卸载时移除监听器
  onUnmounted(() => {
    window.removeEventListener('resize', () => {
      engagementChart.resize()
      emotionHeatmap.resize()
    })
    engagementChart.dispose()
    emotionHeatmap.dispose()
  })
}

onMounted(() => {
  fetchData()
  fetchStyleData()
  fetchEmotionData()
})

// 添加 watch 来监听数据变化
watch([styleData, emotionData], ([newStyleData, newEmotionData]) => {
  if (newStyleData.length && newEmotionData.length) {
    nextTick(() => {
      initAdvancedCharts()
    })
  }
}, { deep: true })
</script>

<template>
  <div class="min-h-screen bg-gray-50 dark:bg-zinc-900/20 p-6">
    <div class="max-w-7xl mx-auto space-y-8">
      <!-- 页面标题 -->
      <div class="flex items-center justify-between"
           v-motion
           :initial="{ opacity: 0, x: -20 }"
           :enter="{ opacity: 1, x: 0, transition: { duration: 500 } }">
        <div class="flex items-center gap-3">
          <div class="h-8 w-1 bg-gradient-to-b from-blue-500 to-cyan-500 rounded-full"></div>
          <h1 class="text-3xl font-bold text-gray-900 dark:text-gray-100">
            数据分析
          </h1>
        </div>
        <a-button 
          type="primary"
          :loading="loading"
          @click="fetchData"
          class="flex items-center gap-2"
        >
          <template #icon>
            <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
            </svg>
          </template>
          刷新数据
        </a-button>
      </div>

      <!-- 分析建议按钮 -->
      <div class="flex justify-end"
           v-motion
           :initial="{ opacity: 0, x: 20 }"
           :enter="{ opacity: 1, x: 0, transition: { duration: 500, delay: 200 } }">
        <a-button
          type="primary"
          :loading="analysisOptions.isGenerating"
          @click="generateAnalysis"
          class="flex items-center gap-2"
        >
          <template #icon>
            <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" d="M9.75 3.104v5.714a2.25 2.25 0 01-.659 1.591L5 14.5M9.75 3.104c-.251.023-.501.05-.75.082m.75-.082a24.301 24.301 0 014.5 0m0 0v5.714c0 .597.237 1.17.659 1.591L19.8 15.3M14.25 3.104c.251.023.501.05.75.082M19.8 15.3l-1.57.393A9.065 9.065 0 0112 15a9.065 9.065 0 00-6.23-.693L5 14.5m14.8.8l1.402 1.402c1.232 1.232.65 3.318-1.067 3.611A48.309 48.309 0 0112 21c-2.773 0-5.491-.235-8.135-.687-1.718-.293-2.3-2.379-1.067-3.61L5 14.5" />
            </svg>
          </template>
          获取创作建议
        </a-button>
      </div>

      <!-- AI分析建议区域 -->
      <div v-if="analysisOptions.showAnalysis" 
           class="bg-white dark:bg-zinc-800 select-text rounded-xl p-6 border theme-border hover:shadow-lg transition-shadow"
           v-motion
           :initial="{ opacity: 0, scale: 0.95 }"
           :enter="{ opacity: 1, scale: 1, transition: { duration: 500 } }">
        <div class="flex items-center justify-between mb-6">
          <h2 class="text-xl font-bold text-gray-900 dark:text-gray-100">创作建议分析</h2>
          <a-button
            type="text"
            @click="analysisOptions.showAnalysis = false"
            class="text-gray-500 hover:text-gray-700 dark:text-gray-400 dark:hover:text-gray-200"
          >
            <template #icon>
              <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12" />
              </svg>
            </template>
          </a-button>
        </div>
        
        <div ref="analysisOptions.content" class="prose dark:prose-invert max-w-none">
          <div v-if="analysisOptions.isGenerating" class="mb-4 flex items-center gap-2 text-blue-500">
            <div class="animate-spin h-4 w-4 border-2 border-blue-500 rounded-full border-t-transparent"></div>
            <span>AI 正在分析...</span>
          </div>
          <MdPreview
            style="background: transparent"
            :theme="themeStore.currentTheme"
            editorId="analysis-preview"
            :modelValue="analysisOptions.content"
            previewTheme="github"
          />
        </div>
      </div>

      <!-- 数据概览卡片 -->
      <div class="grid grid-cols-1 md:grid-cols-3 gap-6">
        <!-- 项目数据卡片 -->
        <div class="bg-white dark:bg-zinc-800 rounded-xl p-6 border theme-border hover:shadow-lg transition-shadow"
             v-motion
             :initial="{ opacity: 0, y: 20 }"
             :enter="{ opacity: 1, y: 0, transition: { duration: 500, delay: 100 } }"
             :hover="{ scale: 1.02, transition: { duration: 200 } }">
          <div class="flex items-center gap-3 mb-4">
            <div class="p-2 bg-indigo-100 dark:bg-indigo-900/30 rounded-lg">
              <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6 text-indigo-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10" />
              </svg>
            </div>
            <h3 class="text-lg font-semibold text-gray-900 dark:text-gray-100">
              项目概览
            </h3>
          </div>
          <div class="space-y-3">
            <div class="flex justify-between items-center">
              <span class="text-gray-600 dark:text-gray-400">项目总数</span>
              <span class="text-xl font-bold text-indigo-500">{{ dataOverview?.totalProjects || 0 }}</span>
            </div>
            <div class="flex justify-between items-center">
              <span class="text-gray-600 dark:text-gray-400">最受欢迎项目</span>
              <span class="text-md font-medium text-indigo-500">{{ dataOverview?.popularProject.name }}</span>
            </div>
            <div class="mt-2 pt-2 border-t dark:border-zinc-700">
              <div class="text-xs text-gray-500 dark:text-gray-400 mb-1">近期观看增长率</div>
              <div class="flex items-center gap-2">
                <div class="text-lg font-bold" :class="dataOverview?.growthRate > 0 ? 'text-green-500' : 'text-red-500'">
                  {{ dataOverview?.growthRate > 0 ? '+' : '' }}{{ dataOverview?.growthRate }}%
                </div>
                <svg v-if="dataOverview?.growthRate > 0" xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 text-green-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" d="M13 7h8m0 0v8m0-8l-8 8-4-4-6 6" />
                </svg>
                <svg v-else xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 text-red-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" d="M13 17h8m0 0V9m0 8l-8-8-4 4-6-6" />
                </svg>
              </div>
            </div>
          </div>
        </div>

        <!-- 观看与互动卡片 -->
        <div class="bg-white dark:bg-zinc-800 rounded-xl p-6 border theme-border hover:shadow-lg transition-shadow"
             v-motion
             :initial="{ opacity: 0, y: 20 }"
             :enter="{ opacity: 1, y: 0, transition: { duration: 500, delay: 200 } }"
             :hover="{ scale: 1.02, transition: { duration: 200 } }">
          <div class="flex items-center gap-3 mb-4">
            <div class="p-2 bg-blue-100 dark:bg-blue-900/30 rounded-lg">
              <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6 text-blue-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
                <path stroke-linecap="round" stroke-linejoin="round" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z" />
              </svg>
            </div>
            <h3 class="text-lg font-semibold text-gray-900 dark:text-gray-100">
              观看与互动
            </h3>
          </div>
          <div class="space-y-3">
            <div class="flex justify-between items-center">
              <span class="text-gray-600 dark:text-gray-400">总观看</span>
              <span class="text-xl font-bold text-blue-500">{{ dataOverview?.totalWatches.toLocaleString() || 0 }}</span>
            </div>
            <div class="flex justify-between items-center">
              <span class="text-gray-600 dark:text-gray-400">总收藏</span>
              <span class="text-xl font-bold text-blue-500">{{ dataOverview?.totalFavorites.toLocaleString() || 0 }}</span>
            </div>
            <div class="mt-2 pt-2 border-t dark:border-zinc-700">
              <div class="text-xs text-gray-500 dark:text-gray-400 mb-1">互动率 (收藏/观看)</div>
              <div class="flex items-center gap-2">
                <div class="w-full bg-gray-200 dark:bg-zinc-700 rounded-full h-2.5">
                  <div class="bg-blue-500 h-2.5 rounded-full" :style="`width: ${Math.min(dataOverview?.interactionRate || 0, 100)}%`"></div>
                </div>
                <span class="text-sm font-medium text-blue-500">{{ dataOverview?.interactionRate || 0 }}%</span>
              </div>
            </div>
          </div>
        </div>

        <!-- 内容偏好卡片 -->
        <div class="bg-white dark:bg-zinc-800 rounded-xl p-6 border theme-border hover:shadow-lg transition-shadow"
             v-motion
             :initial="{ opacity: 0, y: 20 }"
             :enter="{ opacity: 1, y: 0, transition: { duration: 500, delay: 300 } }"
             :hover="{ scale: 1.02, transition: { duration: 200 } }">
          <div class="flex items-center gap-3 mb-4">
            <div class="p-2 bg-amber-100 dark:bg-amber-900/30 rounded-lg">
              <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6 text-amber-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" d="M9.53 16.122a3 3 0 00-5.78 1.128 2.25 2.25 0 01-2.4 2.245 4.5 4.5 0 008.4-2.245c0-.399-.078-.78-.22-1.128zm0 0a15.998 15.998 0 003.388-1.62m-5.043-.025a15.994 15.994 0 011.622-3.395m3.42 3.42a15.995 15.995 0 004.764-4.648l3.876-5.814a1.151 1.151 0 00-1.597-1.597L14.146 6.32a15.996 15.996 0 00-4.649 4.763m3.42 3.42a6.776 6.776 0 00-3.42-3.42" />
              </svg>
            </div>
            <h3 class="text-lg font-semibold text-gray-900 dark:text-gray-100">
              内容偏好
            </h3>
          </div>
          <div class="space-y-3">
            <div class="flex justify-between items-center">
              <span class="text-gray-600 dark:text-gray-400">最受欢迎风格</span>
              <span class="text-md font-medium text-amber-500">{{ dataOverview?.popularStyle.style }}</span>
            </div>
            <div class="flex justify-between items-center">
              <span class="text-gray-600 dark:text-gray-400">主要情绪反应</span>
              <span class="text-md font-medium text-amber-500">{{ dataOverview?.mainEmotion.emotion }}</span>
            </div>
            <div class="mt-2 pt-2 border-t dark:border-zinc-700">
              <div class="text-xs text-gray-500 dark:text-gray-400 mb-1">近期观看趋势</div>
              <div class="flex items-end justify-between h-12 mt-1">
                <div v-for="(item, index) in dataOverview?.recentTrend" :key="index" class="flex flex-col items-center">
                  <div class="text-[10px] text-gray-400">{{ item.date.slice(-2) }}日</div>
                  <div class="w-8 bg-gradient-to-t from-amber-500 to-amber-300 dark:from-amber-600 dark:to-amber-400 rounded-sm" 
                       :style="`height: ${Math.max(item.watches / (Math.max(...(dataOverview?.recentTrend || []).map(i => i.watches)) || 1) * 100, 5)}%`"></div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- 图表展示 -->
      <div class="grid grid-cols-1 lg:grid-cols-2 gap-8">
        <!-- 观看趋势图 -->
        <div class="bg-white dark:bg-zinc-800 rounded-xl p-6 border theme-border hover:shadow-lg transition-shadow"
             v-motion
             :initial="{ opacity: 0, x: -20 }"
             :enter="{ opacity: 1, x: 0, transition: { duration: 500, delay: 400 } }"
             :hover="{ scale: 1.01, transition: { duration: 200 } }">
          <div class="flex items-center gap-3 mb-4">
            <div class="p-2 bg-blue-100 dark:bg-blue-900/30 rounded-lg">
              <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6 text-blue-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
                <path stroke-linecap="round" stroke-linejoin="round" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z" />
              </svg>
            </div>
            <h3 class="text-lg font-semibold text-gray-900 dark:text-gray-100">观看趋势分析</h3>
          </div>
          <p class="text-gray-600 dark:text-gray-400 mb-4">追踪不同项目的观看量变化趋势，帮助您了解内容受欢迎程度的时间分布。</p>
          <div id="watchChart" class="w-full h-[400px]"></div>
        </div>

        <!-- 收藏趋势图 -->
        <div class="bg-white dark:bg-zinc-800 rounded-xl p-6 border theme-border hover:shadow-lg transition-shadow"
             v-motion
             :initial="{ opacity: 0, x: 20 }"
             :enter="{ opacity: 1, x: 0, transition: { duration: 500, delay: 500 } }"
             :hover="{ scale: 1.01, transition: { duration: 200 } }">
          <div class="flex items-center gap-3 mb-4">
            <div class="p-2 bg-amber-100 dark:bg-amber-900/30 rounded-lg">
              <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6 text-amber-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" d="M11.48 3.499a.562.562 0 011.04 0l2.125 5.111a.563.563 0 00.475.345l5.518.442c.499.04.701.663.321.988l-4.204 3.602a.563.563 0 00-.182.557l1.285 5.385a.562.562 0 01-.84.61l-4.725-2.885a.563.563 0 00-.586 0L6.982 20.54a.562.562 0 01-.84-.61l1.285-5.386a.562.562 0 00-.182-.557l-4.204-3.602a.563.563 0 01.321-.988l5.518-.442a.563.563 0 00.475-.345L11.48 3.5z" />
              </svg>
            </div>
            <h3 class="text-lg font-semibold text-gray-900 dark:text-gray-100">收藏趋势分析</h3>
          </div>
          <p class="text-gray-600 dark:text-gray-400 mb-4">展示各项目收藏数量的变化情况，反映用户对内容的长期兴趣。</p>
          <div id="favoriteChart" class="w-full h-[400px]"></div>
        </div>

        <!-- 风格分析图表 -->
        <div class="bg-white dark:bg-zinc-800 rounded-xl p-6 border theme-border hover:shadow-lg transition-shadow"
             v-motion
             :initial="{ opacity: 0, x: -20 }"
             :enter="{ opacity: 1, x: 0, transition: { duration: 500, delay: 600 } }"
             :hover="{ scale: 1.01, transition: { duration: 200 } }">
          <div class="flex items-center gap-3 mb-4">
            <div class="p-2 bg-purple-100 dark:bg-purple-900/30 rounded-lg">
              <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6 text-purple-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" d="M9.53 16.122a3 3 0 00-5.78 1.128 2.25 2.25 0 01-2.4 2.245 4.5 4.5 0 008.4-2.245c0-.399-.078-.78-.22-1.128zm0 0a15.998 15.998 0 003.388-1.62m-5.043-.025a15.994 15.994 0 011.622-3.395m3.42 3.42a15.995 15.995 0 004.764-4.648l3.876-5.814a1.151 1.151 0 00-1.597-1.597L14.146 6.32a15.996 15.996 0 00-4.649 4.763m3.42 3.42a6.776 6.776 0 00-3.42-3.42" />
              </svg>
            </div>
            <h3 class="text-lg font-semibold text-gray-900 dark:text-gray-100">内容风格分布</h3>
          </div>
          <p class="text-gray-600 dark:text-gray-400 mb-4">分析不同创作风格的分布情况，了解受众对各类风格的偏好。</p>
          <div id="styleChart" class="w-full h-[400px]"></div>
        </div>

        <!-- 类型分析图表 -->
        <div class="bg-white dark:bg-zinc-800 rounded-xl p-6 border theme-border hover:shadow-lg transition-shadow"
             v-motion
             :initial="{ opacity: 0, x: 20 }"
             :enter="{ opacity: 1, x: 0, transition: { duration: 500, delay: 700 } }"
             :hover="{ scale: 1.01, transition: { duration: 200 } }">
          <div class="flex items-center gap-3 mb-4">
            <div class="p-2 bg-green-100 dark:bg-green-900/30 rounded-lg">
              <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6 text-green-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" d="M3.75 12h16.5m-16.5 3.75h16.5M3.75 19.5h16.5M5.625 4.5h12.75a1.875 1.875 0 010 3.75H5.625a1.875 1.875 0 010-3.75z" />
              </svg>
            </div>
            <h3 class="text-lg font-semibold text-gray-900 dark:text-gray-100">作品类型分析</h3>
          </div>
          <p class="text-gray-600 dark:text-gray-400 mb-4">展示不同类型作品的数量和表现，帮助优化内容类型策略。</p>
          <div id="typeChart" class="w-full h-[400px]"></div>
        </div>

        <!-- 情绪分析图表 -->
        <div class="bg-white dark:bg-zinc-800 rounded-xl p-6 border theme-border hover:shadow-lg transition-shadow"
             v-motion
             :initial="{ opacity: 0, x: -20 }"
             :enter="{ opacity: 1, x: 0, transition: { duration: 500, delay: 800 } }"
             :hover="{ scale: 1.01, transition: { duration: 200 } }">
          <div class="flex items-center gap-3 mb-4">
            <div class="p-2 bg-rose-100 dark:bg-rose-900/30 rounded-lg">
              <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6 text-rose-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" d="M15.182 15.182a4.5 4.5 0 01-6.364 0M21 12a9 9 0 11-18 0 9 9 0 0118 0zM9.75 9.75c0 .414-.168.75-.375.75S9 10.164 9 9.75 9.168 9 9.375 9s.375.336.375.75zm-.375 0h.008v.015h-.008V9.75zm5.625 0c0 .414-.168.75-.375.75s-.375-.336-.375-.75.168-.75.375-.75.375.336.375.75zm-.375 0h.008v.015h-.008V9.75z" />
              </svg>
            </div>
            <h3 class="text-lg font-semibold text-gray-900 dark:text-gray-100">情绪趋势分析</h3>
          </div>
          <p class="text-gray-600 dark:text-gray-400 mb-4">追踪观众对内容的情感反应变化，优化情节设计和叙事节奏。</p>
          <div id="emotionChart" class="w-full h-[500px]"></div>
        </div>

        <!-- 互动率雷达图 -->
        <div class="bg-white dark:bg-zinc-800 rounded-xl p-6 border theme-border hover:shadow-lg transition-shadow"
             v-motion
             :initial="{ opacity: 0, x: 20 }"
             :enter="{ opacity: 1, x: 0, transition: { duration: 500, delay: 900 } }"
             :hover="{ scale: 1.01, transition: { duration: 200 } }">
          <div class="flex items-center gap-3 mb-4">
            <div class="p-2 bg-indigo-100 dark:bg-indigo-900/30 rounded-lg">
              <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6 text-indigo-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" d="M7.5 14.25v2.25m3-4.5v4.5m3-6.75v6.75m3-9v9M6 20.25h12A2.25 2.25 0 0020.25 18V6A2.25 2.25 0 0018 3.75H6A2.25 2.25 0 003.75 6v12A2.25 2.25 0 006 20.25z" />
              </svg>
            </div>
            <h3 class="text-lg font-semibold text-gray-900 dark:text-gray-100">内容互动率分析</h3>
          </div>
          <p class="text-gray-600 dark:text-gray-400 mb-4">多维度展示不同风格内容的用户互动情况，发现最受欢迎的创作方向。</p>
          <div id="engagementChart" ref="engagementChart" class="w-full h-[400px]"></div>
        </div>

        <!-- 情绪分布热力图 -->
        <div class="bg-white dark:bg-zinc-800 rounded-xl p-6 border theme-border hover:shadow-lg transition-shadow"
             v-motion
             :initial="{ opacity: 0, x: -20 }"
             :enter="{ opacity: 1, x: 0, transition: { duration: 500, delay: 1000 } }"
             :hover="{ scale: 1.01, transition: { duration: 200 } }">
          <div class="flex items-center gap-3 mb-4">
            <div class="p-2 bg-cyan-100 dark:bg-cyan-900/30 rounded-lg">
              <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6 text-cyan-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" d="M3.75 3v11.25A2.25 2.25 0 006 16.5h2.25M3.75 3h-1.5m1.5 0h16.5m0 0h1.5m-1.5 0v11.25A2.25 2.25 0 0118 16.5h-2.25m-7.5 0h7.5m-7.5 0l-1 3m8.5-3l1 3m0 0l.5 1.5m-.5-1.5h-9.5m0 0l-.5 1.5m.75-9l3-3 2.148 2.148A12.061 12.061 0 0116.5 7.605" />
              </svg>
            </div>
            <h3 class="text-lg font-semibold text-gray-900 dark:text-gray-100">情绪分布热力分析</h3>
          </div>
          <p class="text-gray-600 dark:text-gray-400 mb-4">直观展示各项目引发的情绪分布强度，助您把握作品的情感基调。</p>
          <div id="emotionHeatmap" class="w-full h-[400px]"></div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
/* 移除 animate.css 相关样式 */
.prose {
  @apply text-gray-700 dark:text-gray-300;
  max-height: 60vh;
  overflow-y: auto;
  padding-right: 1rem;
  
  /* 自定义滚动条样式 */
  &::-webkit-scrollbar {
    width: 6px;
  }
  
  &::-webkit-scrollbar-track {
    @apply bg-gray-100 dark:bg-zinc-800 rounded-full;
  }
  
  &::-webkit-scrollbar-thumb {
    @apply bg-gray-300 dark:bg-zinc-600 rounded-full;
    &:hover {
      @apply bg-gray-400 dark:bg-zinc-500;
    }
  }
}

/* 调整 Markdown 预览区域的样式 */
:deep(.md-editor-preview-wrapper) {
  @apply px-0 !important;
}
</style>