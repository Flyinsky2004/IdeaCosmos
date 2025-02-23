<script setup>
import { ref, onMounted, nextTick, onUnmounted } from 'vue'
import { get } from '@/util/request'
import { message } from 'ant-design-vue'
import * as echarts from 'echarts'
import SpinLoaderLarge from '@/components/spinLoaderLarge.vue'

const loading = ref(true)
const projectData = ref([])
const styleData = ref([])
const emotionData = ref([])

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

onMounted(() => {
  fetchData()
  fetchStyleData()
  fetchEmotionData()
})
</script>

<template>
  <div class="min-h-screen bg-gray-50 dark:bg-zinc-900/20 p-6">
    <div class="max-w-7xl mx-auto space-y-8">
      <!-- 页面标题 -->
      <div class="flex items-center justify-between animate__animated animate__fadeIn">
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

      <!-- 加载状态 -->
      <div v-if="loading" class="flex items-center justify-center py-12">
        <SpinLoaderLarge />
      </div>

      <template v-else>
        <!-- 数据概览卡片 -->
        <div class="grid grid-cols-1 md:grid-cols-3 gap-6">
          <!-- 项目数量卡片 -->
          <div class="bg-white dark:bg-zinc-800 rounded-xl p-6 border theme-border hover:shadow-lg transition-shadow animate__animated animate__fadeIn animate__delay-1s">
            <div class="flex items-center gap-3 mb-4">
              <div class="p-2 bg-indigo-100 dark:bg-indigo-900/30 rounded-lg">
                <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6 text-indigo-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10" />
                </svg>
              </div>
              <h3 class="text-lg font-semibold text-gray-900 dark:text-gray-100">
                项目总数
              </h3>
            </div>
            <div class="text-3xl font-bold text-indigo-500">
              {{ projectStats.totalProjects }}
            </div>
          </div>

          <!-- 总观看数卡片 -->
          <div class="bg-white dark:bg-zinc-800 rounded-xl p-6 border theme-border hover:shadow-lg transition-shadow animate__animated animate__fadeIn animate__delay-2s">
            <div class="flex items-center gap-3 mb-4">
              <div class="p-2 bg-blue-100 dark:bg-blue-900/30 rounded-lg">
                <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6 text-blue-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
                  <path stroke-linecap="round" stroke-linejoin="round" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z" />
                </svg>
              </div>
              <h3 class="text-lg font-semibold text-gray-900 dark:text-gray-100">
                总观看
              </h3>
            </div>
            <div class="text-3xl font-bold text-blue-500">
              {{ projectStats.totalWatches }}
            </div>
          </div>

          <!-- 总收藏数卡片 -->
          <div class="bg-white dark:bg-zinc-800 rounded-xl p-6 border theme-border hover:shadow-lg transition-shadow animate__animated animate__fadeIn animate__delay-3s">
            <div class="flex items-center gap-3 mb-4">
              <div class="p-2 bg-amber-100 dark:bg-amber-900/30 rounded-lg">
                <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6 text-amber-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" d="M11.48 3.499a.562.562 0 011.04 0l2.125 5.111a.563.563 0 00.475.345l5.518.442c.499.04.701.663.321.988l-4.204 3.602a.563.563 0 00-.182.557l1.285 5.385a.562.562 0 01-.84.61l-4.725-2.885a.563.563 0 00-.586 0L6.982 20.54a.562.562 0 01-.84-.61l1.285-5.386a.562.562 0 00-.182-.557l-4.204-3.602a.563.563 0 01.321-.988l5.518-.442a.563.563 0 00.475-.345L11.48 3.5z" />
                </svg>
              </div>
              <h3 class="text-lg font-semibold text-gray-900 dark:text-gray-100">
                总收藏
              </h3>
            </div>
            <div class="text-3xl font-bold text-amber-500">
              {{ projectStats.totalFavorites }}
            </div>
          </div>
        </div>

        <!-- 图表展示 -->
        <div class="grid grid-cols-1 lg:grid-cols-2 gap-8">
          <!-- 观看趋势图 -->
          <div class="bg-white dark:bg-zinc-800 rounded-xl p-6 border theme-border hover:shadow-lg transition-shadow animate__animated animate__fadeIn animate__delay-4s">
            <div id="watchChart" class="w-full h-[400px]"></div>
          </div>

          <!-- 收藏趋势图 -->
          <div class="bg-white dark:bg-zinc-800 rounded-xl p-6 border theme-border hover:shadow-lg transition-shadow animate__animated animate__fadeIn animate__delay-5s">
            <div id="favoriteChart" class="w-full h-[400px]"></div>
          </div>

          <!-- 风格分析图表 -->
          <div class="bg-white dark:bg-zinc-800 rounded-xl p-6 border theme-border hover:shadow-lg transition-shadow animate__animated animate__fadeIn animate__delay-6s">
            <div id="styleChart" class="w-full h-[400px]"></div>
          </div>

          <!-- 类型分析图表 -->
          <div class="bg-white dark:bg-zinc-800 rounded-xl p-6 border theme-border hover:shadow-lg transition-shadow animate__animated animate__fadeIn animate__delay-7s">
            <div id="typeChart" class="w-full h-[400px]"></div>
          </div>

          <!-- 情绪分析图表 -->
          <div class="bg-white dark:bg-zinc-800 rounded-xl p-6 border theme-border hover:shadow-lg transition-shadow animate__animated animate__fadeIn animate__delay-8s lg:col-span-2">
            <div id="emotionChart" class="w-full h-[500px]"></div>
          </div>
        </div>
      </template>
    </div>
  </div>
</template>

<style scoped>
.animate__animated {
  animation-duration: 0.5s; /* 调整动画持续时间 */
}

/* 自定义延迟时间 */
.animate__delay-1s {
  animation-delay: 0.1s;
}
.animate__delay-2s {
  animation-delay: 0.2s;
}
.animate__delay-3s {
  animation-delay: 0.3s;
}
.animate__delay-4s {
  animation-delay: 0.4s;
}
.animate__delay-5s {
  animation-delay: 0.5s;
}
.animate__delay-6s {
  animation-delay: 0.6s;
}
.animate__delay-7s {
  animation-delay: 0.7s;
}
</style>