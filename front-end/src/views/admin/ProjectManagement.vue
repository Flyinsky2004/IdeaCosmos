<script setup>
import { ref, onMounted, reactive, computed } from 'vue'
import { get, post, postJSON } from '@/util/request'
import { message, Modal } from 'ant-design-vue'
import SpinLoaderLarge from '@/components/spinLoaderLarge.vue'
import { BACKEND_DOMAIN, imagePrefix } from "@/util/VARRIBLES"
import { useRouter } from 'vue-router'
import { SearchOutlined, DownloadOutlined, FileTextOutlined, UndoOutlined, StarOutlined, StarFilled, StopOutlined, DeleteOutlined, FileImageOutlined, TeamOutlined, ReadOutlined, EyeOutlined, HeartOutlined, FileSearchOutlined } from '@ant-design/icons-vue'
import router from '@/router'
const loading = ref(true)
const searchLoading = ref(false)
const projects = ref([])
const searchKeyword = ref('')
const selectedType = ref('全部')
const currentPage = ref(1)
const pageSize = ref(10)
const total = ref(0)
const selectedStatus = ref('全部')
const typesOptions = ref([
  { label: '全部', value: '全部' },
  { label: '小说', value: '小说' },
  { label: '剧本', value: '剧本' },
  { label: '剧情短片脚本', value: '剧情短片脚本' },
  { label: '故事大纲', value: '故事大纲' }
])

// 状态筛选选项
const statusOptions = ref([
  { label: '全部', value: '全部' },
  { label: '正常', value: 'normal' },
  { label: '推荐', value: 'featured' },
  { label: '下架', value: 'banned' }
])

// 项目详情模态框
const projectDetailModal = reactive({
  visible: false,
  loading: false,
  projectData: null
})

// 定义表格列
const columns = [
  { title: 'ID', dataIndex: 'ID', key: 'id', width: 70 },
  { title: '项目名称', dataIndex: 'project_name', key: 'project_name' },
  { title: '类型', dataIndex: 'types', key: 'types' },
  { title: '团队', key: 'team', customRender: ({record}) => record.team?.username || '-' },
  { title: '章节数', key: 'chapter_num', dataIndex: 'chapter_num', align: 'center' },
  { title: '内容评估', key: 'avg_score', align: 'center', width: 120 },
  { title: '状态', key: 'status' },
  { title: '查看/收藏', key: 'stats', customRender: ({record}) => `${record.watches || 0}/${record.favorites || 0}` },
  { title: '更新时间', key: 'updatedAt', customRender: ({record}) => formatDateTime(record.UpdatedAt) },
  { title: '操作', key: 'action', fixed: 'right', width: 160 }
]

// 获取项目列表
const fetchProjects = async () => {
  loading.value = true
  
  get('/api/admin/projects', {
    page: currentPage.value,
    pageSize: pageSize.value,
    keyword: searchKeyword.value,
    type: selectedType.value !== '全部' ? selectedType.value : '',
    status: selectedStatus.value !== '全部' ? selectedStatus.value : ''
  }, 
  (msg, data) => {
    projects.value = data.projects
    total.value = data.total
    loading.value = false
  },
  (msg) => {
    message.warning(msg)
    loading.value = false
  },
  (msg) => {
    message.error(msg)
    loading.value = false
  })
}

// 查看项目详情
const viewProjectDetail = (projectId) => {
  projectDetailModal.loading = true
  projectDetailModal.visible = true
  
  get(`/api/admin/projects/${projectId}`, {}, 
    (msg, data) => {
      projectDetailModal.projectData = data
      projectDetailModal.loading = false
    },
    (msg) => {
      message.warning(msg)
      projectDetailModal.loading = false
      projectDetailModal.visible = false
    },
    (msg) => {
      message.error(msg)
      projectDetailModal.loading = false
      projectDetailModal.visible = false
    }
  )
}

// 删除项目
const deleteProject = (projectId) => {
    Modal.confirm({
    title: '确认删除',
    content: '您确定要删除此项目吗？此操作将删除项目的所有数据，包括章节、角色等，且不可恢复！',
    okText: '确认删除',
    okType: 'danger',
    cancelText: '取消',
    onOk: () => {
      postJSON(`/api/admin/projects/${projectId}/delete`, {}, 
        (msg, data) => {
          message.success('项目删除成功')
          fetchProjects()
        },
        (msg) => {
          message.warning(msg)
        },
        (msg) => {
          message.error(msg)
        }
      )
    }
  })
}

// 更新项目状态
const updateProjectStatus = (projectId, status) => {
  postJSON(`/api/admin/projects/${projectId}/status`, {
    status: status
  }, 
  (msg, data) => {
    message.success('项目状态更新成功')
    fetchProjects()
  },
  (msg) => {
    message.warning(msg)
  },
  (msg) => {
    message.error(msg)
  })
}

// 下架项目
const banProject = (projectId) => {
    Modal.confirm({
    title: '确认下架',
    content: '您确定要下架此项目吗？下架后用户将无法访问此项目',
    okText: '确认下架',
    okType: 'danger',
    cancelText: '取消',
    onOk: () => {
      updateProjectStatus(projectId, 'banned')
    }
  })
}

// 恢复项目
const restoreProject = (projectId) => {
  updateProjectStatus(projectId, 'normal')
}

// 推荐项目
const featureProject = (projectId) => {
  updateProjectStatus(projectId, 'featured')
}

// 搜索
const handleSearch = () => {
  currentPage.value = 1
  searchLoading.value = true
  fetchProjects().finally(() => {
    searchLoading.value = false
  })
}

// 重置搜索
const resetSearch = () => {
  searchKeyword.value = ''
  selectedType.value = '全部'
  selectedStatus.value = '全部'
  currentPage.value = 1
  fetchProjects()
}

// 页码变化
const onPageChange = (page) => {
  currentPage.value = page
  fetchProjects()
}

// 每页条数变化
const onPageSizeChange = (current, size) => {
  currentPage.value = 1
  pageSize.value = size
  fetchProjects()
}

// 格式化日期
const formatDateTime = (dateString) => {
  if (!dateString) return '-'
  const date = new Date(dateString)
  return date.toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit'
  })
}

// 解析JSON字符串
const parseJsonArray = (jsonStr) => {
  if (!jsonStr) return []
  try {
    return JSON.parse(jsonStr)
  } catch (e) {
    return []
  }
}

// 跳转到章节审核页
const goToChapterReview = (chapterId) => {
  router.push(`/admin/chapters/${chapterId}`)
}

// 获取项目状态标签类型
const getStatusTagType = (status) => {
  switch(status) {
    case 'featured': return 'success'
    case 'banned': return 'error'
    default: return 'default'
  }
}

// 获取项目状态显示文本
const getStatusText = (status) => {
  switch(status) {
    case 'featured': return '推荐'
    case 'banned': return '下架'
    default: return '正常'
  }
}

// 导出项目数据
const exportProjects = () => {
  message.loading({ content: '正在准备导出数据...', key: 'export' })
  
  get('/api/admin/projects/export', {
    keyword: searchKeyword.value,
    type: selectedType.value !== '全部' ? selectedType.value : '',
    status: selectedStatus.value !== '全部' ? selectedStatus.value : ''
  }, 
  (msg, data) => {
    // 创建CSV内容
    let csvContent = "ID,项目名称,类型,团队,章节数,内容评估分数,内容评估等级,状态,查看数,收藏数,创建时间,更新时间\n"
    
    projects.value.forEach(project => {
      csvContent += `${project.ID},"${project.project_name}",${project.types},${project.team?.username || ''},`
      csvContent += `${project.chapter_num || 0},${project.avg_score ? project.avg_score.toFixed(1) : 0},`
      csvContent += `${getContentReviewLevel(project.avg_score)},${getStatusText(project.status || 'normal')},`
      csvContent += `${project.watches || 0},${project.favorites || 0},`
      csvContent += `${formatDateTime(project.CreatedAt)},${formatDateTime(project.UpdatedAt)}\n`
    })
    
    // 创建下载链接
    const blob = new Blob([csvContent], { type: 'text/csv;charset=utf-8;' })
    const url = window.URL.createObjectURL(blob)
    const link = document.createElement('a')
    link.href = url
    link.setAttribute('download', `项目数据_${new Date().toLocaleDateString()}.csv`)
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
    
    message.success({ content: '导出成功！', key: 'export' })
  },
  (msg) => {
    message.warning({ content: msg, key: 'export' })
  },
  (msg) => {
    message.error({ content: msg, key: 'export' })
  })
}

// 查看团队详情
const viewTeamDetail = (teamId) => {
  router.push(`/admin/teams/${teamId}`)
}

// 关闭模态框时清理数据
const handleModalClose = () => {
  projectDetailModal.projectData = null
}

// 获取内容评估颜色
const getContentReviewColor = (score) => {
  if (!score) return '#999999'
  if (score >= 90) return '#52c41a' // 绿色，安全
  if (score >= 70) return '#1890ff' // 蓝色，正常
  if (score >= 50) return '#faad14' // 黄色，警告
  return '#f5222d' // 红色，风险
}

// 获取内容评估等级
const getContentReviewLevel = (score) => {
  if (!score) return '未评估'
  if (score >= 90) return '安全'
  if (score >= 70) return '合规'
  if (score >= 50) return '警告'
  return '风险'
}

// 获取内容评估提示
const getContentReviewTip = (score) => {
  if (!score) return '内容尚未经过评估'
  if (score >= 90) return '内容安全，符合价值观和内容规范'
  if (score >= 70) return '内容基本合规，个别内容需谨慎'
  if (score >= 50) return '内容包含部分敏感话题，建议审核'
  return '内容存在明显敏感信息，触及政治、宗教等禁区，需严格审核'
}

onMounted(() => {
  fetchProjects()
})
</script>

<template>
  <div class="p-4">
    <h1 class="text-2xl font-bold mb-6 dark:text-gray-100">项目管理</h1>

    <div v-if="loading" class="flex justify-center items-center h-64">
      <SpinLoaderLarge />
    </div>

    <div v-else class="animate__animated animate__fadeIn">
      <!-- 搜索工具栏 -->
      <div class="bg-white dark:bg-zinc-800 p-4 rounded-lg shadow-sm mb-6">
        <div class="flex flex-wrap gap-4 items-end">
          <div>
            <span class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">关键词搜索</span>
            <a-input 
              v-model:value="searchKeyword" 
              placeholder="项目名称/社会背景" 
              class="w-64" 
              @pressEnter="handleSearch"
              allowClear
            />
          </div>
          <div>
            <span class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">项目类型</span>
            <a-select 
              v-model:value="selectedType" 
              style="width: 150px"
              :options="typesOptions"
            />
          </div>
          <div>
            <span class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">项目状态</span>
            <a-select 
              v-model:value="selectedStatus" 
              style="width: 150px"
              :options="statusOptions"
            />
          </div>
          <div class="flex gap-2">
            <a-button type="primary" @click="handleSearch" :loading="searchLoading">
              <template #icon><SearchOutlined /></template>
              搜索
            </a-button>
            <a-button @click="resetSearch">重置</a-button>
          </div>
        </div>
      </div>

      <!-- 项目列表 -->
      <div class="bg-white dark:bg-zinc-800 rounded-lg shadow-sm p-6">
        <div class="flex justify-between mb-4">
          <h3 class="text-lg font-medium dark:text-gray-200">项目列表 (共 {{ total }} 个)</h3>
          <a-button type="primary" @click="exportProjects" :disabled="projects.length === 0">
            <template #icon><DownloadOutlined /></template>
            导出数据
          </a-button>
        </div>

        <a-table 
          :dataSource="projects" 
          :columns="columns" 
          :pagination="false"
          rowKey="ID"
          :rowClassName="(record) => record.status === 'banned' ? 'bg-gray-50 dark:bg-zinc-900/30' : ''"
          :loading="loading"
        >
          <!-- 名称列 -->
          <template #bodyCell="{ column, record }">
            <template v-if="column.key === 'project_name'">
              <div class="flex items-center">
                <a-avatar 
                  v-if="record.cover_image" 
                  :src="`${imagePrefix}${record.cover_image}`" 
                  :alt="record.project_name"
                  shape="square"
                  :size="40"
                  class="mr-2"
                />
                <a-avatar v-else :size="40" class="mr-2">
                  <template #icon><FileTextOutlined /></template>
                </a-avatar>
                <span class="font-medium dark:text-gray-200 hover:text-blue-500 cursor-pointer" @click="viewProjectDetail(record.ID)">
                  {{ record.project_name }}
                </span>
              </div>
            </template>
            
            <!-- 评分列 -->
            <template v-else-if="column.key === 'avg_score'">
              <div class="flex items-center">
                <a-progress
                  v-if="record.avg_score"
                  type="circle"
                  :percent="record.avg_score"
                  :width="30"
                  :size="50"
                  :strokeWidth="6"
                  :strokeColor="getContentReviewColor(record.avg_score)"
                  :format="percent => `${percent}`"
                />
                <span v-else>-</span>
                <a-tooltip :title="getContentReviewTip(record.avg_score)">
                  <a-tag class="ml-1" size="small" :color="getContentReviewColor(record.avg_score)">
                    {{ getContentReviewLevel(record.avg_score) }}
                  </a-tag>
                </a-tooltip>
              </div>
            </template>
            
            <!-- 状态列 -->
            <template v-else-if="column.key === 'status'">
              <a-tag 
                :color="getStatusTagType(record.status || 'normal')"
              >
                {{ getStatusText(record.status || 'normal') }}
              </a-tag>
            </template>
            
            <!-- 操作列 -->
            <template v-else-if="column.key === 'action'">
              <div class="flex gap-2">
                <a-tooltip title="查看详情">
                  <a-button size="small" class="flex justify-center items-center w-8 h-8 p-0" @click="viewProjectDetail(record.ID)">
                    <template #icon><EyeOutlined /></template>
                  </a-button>
                </a-tooltip>
                
                <!-- 根据状态显示不同操作按钮 -->
                <template v-if="record.status === 'banned'">
                  <a-tooltip title="恢复项目">
                    <a-button
                      size="small"
                      type="primary"
                      class="flex justify-center items-center w-8 h-8 p-0"
                      @click="restoreProject(record.ID)"
                    >
                      <template #icon><UndoOutlined /></template>
                    </a-button>
                  </a-tooltip>
                </template>
                <template v-else>
                  <!-- 推荐项目/取消推荐 -->
                  <a-tooltip :title="record.status === 'featured' ? '取消推荐' : '设为推荐'">
                    <a-button
                      size="small"
                      :type="record.status === 'featured' ? 'default' : 'primary'"
                      class="flex justify-center items-center w-8 h-8 p-0"
                      @click="record.status === 'featured' ? restoreProject(record.ID) : featureProject(record.ID)"
                    >
                      <template #icon>
                        <StarOutlined v-if="record.status !== 'featured'" />
                        <StarFilled v-else />
                      </template>
                    </a-button>
                  </a-tooltip>
                  
                  <!-- 下架按钮 -->
                  <a-tooltip title="下架项目">
                    <a-button 
                      size="small" 
                      danger 
                      class="flex justify-center items-center w-8 h-8 p-0"
                      @click="banProject(record.ID)"
                    >
                      <template #icon><StopOutlined /></template>
                    </a-button>
                  </a-tooltip>
                </template>
                
                <!-- 删除按钮 -->
                <a-tooltip title="删除项目">
                  <a-button 
                    size="small" 
                    danger 
                    class="flex justify-center items-center w-8 h-8 p-0"
                    @click="deleteProject(record.ID)"
                  >
                    <template #icon><DeleteOutlined /></template>
                  </a-button>
                </a-tooltip>
              </div>
            </template>
          </template>
        </a-table>

        <!-- 分页 -->
        <div class="flex justify-end mt-4">
          <a-pagination
            v-model:current="currentPage"
            v-model:pageSize="pageSize"
            :total="total"
            show-size-changer
            :pageSizeOptions="['10', '20', '50', '100']"
            :showTotal="total => `共 ${total} 条`"
            @change="onPageChange"
            @showSizeChange="onPageSizeChange"
          />
        </div>
      </div>
    </div>

    <!-- 项目详情模态框 -->
    <a-modal
      title="项目详情"
      v-model:visible="projectDetailModal.visible"
      width="800px"
      :footer="null"
      :afterClose="handleModalClose"
      :maskClosable="false"
      class="project-detail-modal"
    >
      <template #title>
        <div class="flex items-center">
          <FileTextOutlined class="mr-2" />
          <span>项目详情</span>
        </div>
      </template>

      <div v-if="projectDetailModal.loading" class="py-10 flex justify-center">
        <SpinLoaderLarge />
      </div>
      
      <div v-else-if="projectDetailModal.projectData" class="max-h-[70vh] overflow-y-auto pr-2">
        <!-- 顶部操作栏 -->
        <div class="sticky top-0 z-10 bg-white dark:bg-zinc-800 py-2 shadow-sm mb-4 flex justify-end gap-2">
          <a-button-group>
            <template v-if="projectDetailModal.projectData.project.status === 'banned'">
              <a-button type="primary" @click="restoreProject(projectDetailModal.projectData.project.ID)">
                <template #icon><UndoOutlined /></template>
                恢复项目
              </a-button>
            </template>
            <template v-else>
              <a-button 
                :type="projectDetailModal.projectData.project.status === 'featured' ? 'default' : 'primary'"
                @click="projectDetailModal.projectData.project.status === 'featured' ? 
                  restoreProject(projectDetailModal.projectData.project.ID) : 
                  featureProject(projectDetailModal.projectData.project.ID)"
              >
                <template #icon>
                  <StarOutlined v-if="projectDetailModal.projectData.project.status !== 'featured'" />
                  <StarFilled v-else />
                </template>
                {{ projectDetailModal.projectData.project.status === 'featured' ? '取消推荐' : '设为推荐' }}
              </a-button>
            
              <a-button danger @click="banProject(projectDetailModal.projectData.project.ID)">
                <template #icon><StopOutlined /></template>
                下架项目
              </a-button>
            </template>
          </a-button-group>
          
          <a-button danger @click="deleteProject(projectDetailModal.projectData.project.ID)">
            <template #icon><DeleteOutlined /></template>
            删除项目
          </a-button>
        </div>
        
        <!-- 项目基本信息 -->
        <div class="flex flex-wrap md:flex-nowrap gap-6 mb-6">
          <div class="w-full md:w-1/3">
            <div class="relative">
              <img 
                v-if="projectDetailModal.projectData.project.cover_image" 
                :src="`${imagePrefix}${projectDetailModal.projectData.project.cover_image}`" 
                :alt="projectDetailModal.projectData.project.project_name"
                class="w-full aspect-[3/4] object-cover rounded-lg shadow-md hover:shadow-lg transition-shadow"
              />
              <div v-else class="w-full aspect-[3/4] bg-gray-200 dark:bg-zinc-700 rounded-lg shadow-md flex items-center justify-center">
                <FileImageOutlined class="text-gray-400 text-6xl" />
              </div>
              
              <a-tag 
                :color="getStatusTagType(projectDetailModal.projectData.project.status || 'normal')" 
                class="absolute top-2 right-2 px-2 py-1"
              >
                {{ getStatusText(projectDetailModal.projectData.project.status || 'normal') }}
              </a-tag>
            </div>
          </div>
          
          <div class="w-full md:w-2/3">
            <h2 class="text-2xl font-bold dark:text-gray-200 mb-4 flex items-center">
              {{ projectDetailModal.projectData.project.project_name }}
              <a-button 
                v-if="projectDetailModal.projectData.project.team" 
                type="link" 
                class="ml-2"
                @click="viewTeamDetail(projectDetailModal.projectData.project.team.ID)"
              >
                <template #icon><TeamOutlined /></template>
                {{ projectDetailModal.projectData.project.team.username }}
              </a-button>
            </h2>
            
            <div class="mb-4 flex flex-wrap gap-2">
              <a-tag color="blue">{{ projectDetailModal.projectData.project.types }}</a-tag>
              <a-tag 
                v-for="style in parseJsonArray(projectDetailModal.projectData.project.style)" 
                :key="style"
                color="purple"
              >
                {{ style }}
              </a-tag>
              <a-tag 
                v-for="market in parseJsonArray(projectDetailModal.projectData.project.market_people)" 
                :key="market"
                color="cyan"
              >
                {{ market }}
              </a-tag>
            </div>
            
            <a-descriptions bordered size="small" class="mb-4">
              <a-descriptions-item label="创建时间" :span="2">
                {{ formatDateTime(projectDetailModal.projectData.project.CreatedAt) }}
              </a-descriptions-item>
              <a-descriptions-item label="更新时间" :span="2">
                {{ formatDateTime(projectDetailModal.projectData.project.UpdatedAt) }}
              </a-descriptions-item>
              <a-descriptions-item label="查看/收藏">
                <EyeOutlined /> {{ projectDetailModal.projectData.project.watches }}
                <HeartOutlined class="ml-3" /> {{ projectDetailModal.projectData.project.favorites }}
              </a-descriptions-item>
              <a-descriptions-item label="章节数量">
                <ReadOutlined /> {{ projectDetailModal.projectData.chapter_count || 0 }}
              </a-descriptions-item>
              <a-descriptions-item label="内容评估" :span="2">
                <div class="flex items-center">
                  <a-progress
                    v-if="projectDetailModal.projectData.avg_score"
                    type="circle"
                    :percent="projectDetailModal.projectData.avg_score"
                    :width="50"
                    :strokeColor="getContentReviewColor(projectDetailModal.projectData.avg_score)"
                    :format="percent => `${percent}分`"
                  />
                  <span v-else>-</span>
                  <div class="ml-3">
                    <div>
                      <a-tag :color="getContentReviewColor(projectDetailModal.projectData.avg_score)">
                        {{ getContentReviewLevel(projectDetailModal.projectData.avg_score) }}
                      </a-tag>
                    </div>
                    <div class="text-sm text-gray-500 dark:text-gray-400 mt-1">
                      {{ getContentReviewTip(projectDetailModal.projectData.avg_score) }}
                    </div>
                  </div>
                </div>
              </a-descriptions-item>
            </a-descriptions>
          </div>
        </div>
        
        <!-- 内容标签页 -->
        <a-tabs defaultActiveKey="1">
          <a-tab-pane key="1" tab="项目内容">
            <div class="grid grid-cols-1 gap-4 mb-6">
              <a-card title="社会背景" size="small" class="dark:bg-zinc-900">
                <p class="text-gray-700 dark:text-gray-400">
                  {{ projectDetailModal.projectData.project.social_story }}
                </p>
              </a-card>
              <a-card title="开始情景" size="small" class="dark:bg-zinc-900">
                <p class="text-gray-700 dark:text-gray-400">
                  {{ projectDetailModal.projectData.project.start }}
                </p>
              </a-card>
              <a-card title="高潮和冲突" size="small" class="dark:bg-zinc-900">
                <p class="text-gray-700 dark:text-gray-400">
                  {{ projectDetailModal.projectData.project.high_point }}
                </p>
              </a-card>
              <a-card title="解决结局" size="small" class="dark:bg-zinc-900">
                <p class="text-gray-700 dark:text-gray-400">
                  {{ projectDetailModal.projectData.project.resolved }}
                </p>
              </a-card>
            </div>
          </a-tab-pane>
          
          <a-tab-pane key="2" tab="章节列表">
            <div v-if="projectDetailModal.projectData.chapters.length > 0">
              <a-list 
                :dataSource="projectDetailModal.projectData.chapters"
                size="small"
                bordered
              >
                <template #renderItem="{ item, index }">
                  <a-list-item>
                    <template #actions>
                      <a-button type="link" @click="goToChapterReview(item.ID)">
                        <template #icon><FileSearchOutlined /></template>
                        审核
                      </a-button>
                    </template>
                    <a-list-item-meta>
                      <template #title>
                        <div class="flex items-center">
                          <span class="mr-2 w-6 h-6 bg-blue-500 text-white rounded-full flex items-center justify-center text-xs">
                            {{ index + 1 }}
                          </span>
                          {{ item.Tittle }}
                        </div>
                      </template>
                      <template #description>{{ item.Description }}</template>
                    </a-list-item-meta>
                  </a-list-item>
                </template>
              </a-list>
            </div>
            <a-empty v-else description="该项目暂无章节" />
          </a-tab-pane>
          
          <a-tab-pane key="3" tab="角色列表">
            <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
              <a-card 
                v-for="character in projectDetailModal.projectData.characters" 
                :key="character.ID"
                size="small"
                hoverable
                class="dark:bg-zinc-900"
              >
                <template #cover>
                  <div class="p-4 flex justify-center">
                    <a-avatar 
                      :size="64"
                      :src="character.avatar ? `${imagePrefix}${character.avatar}` : null"
                    >
                      <template #icon v-if="!character.avatar"><UserOutlined /></template>
                    </a-avatar>
                  </div>
                </template>
                <a-card-meta :title="character.name">
                  <template #description>
                    <div class="line-clamp-3 text-gray-700 dark:text-gray-400">
                      {{ character.description }}
                    </div>
                  </template>
                </a-card-meta>
              </a-card>
            </div>
            <a-empty description="该项目暂无角色" />
          </a-tab-pane>
        </a-tabs>
      </div>
    </a-modal>
  </div>
</template>

<style scoped>
.animate__animated {
  animation-duration: 0.5s;
}

/* 美化滚动条 */
::-webkit-scrollbar {
  width: 8px;
}

::-webkit-scrollbar-track {
  background: #f1f1f1;
  border-radius: 8px;
}

::-webkit-scrollbar-thumb {
  background: #888;
  border-radius: 8px;
}

::-webkit-scrollbar-thumb:hover {
  background: #555;
}

.dark ::-webkit-scrollbar-track {
  background: #333;
}

.dark ::-webkit-scrollbar-thumb {
  background: #555;
}

.dark ::-webkit-scrollbar-thumb:hover {
  background: #777;
}

/* 项目详情模态框样式 */
.project-detail-modal :deep(.ant-descriptions-item-label) {
  background-color: rgba(0, 0, 0, 0.02);
  width: 120px;
}

.dark .project-detail-modal :deep(.ant-descriptions-item-label) {
  background-color: rgba(255, 255, 255, 0.05);
  color: rgba(255, 255, 255, 0.85);
}

.dark .project-detail-modal :deep(.ant-descriptions-item-content) {
  color: rgba(255, 255, 255, 0.65);
}

.project-detail-modal :deep(.ant-tabs-tab) {
  transition: all 0.2s;
}

.project-detail-modal :deep(.ant-tabs-tab:hover) {
  color: #1890ff;
}

.dark .project-detail-modal :deep(.ant-tabs-tab) {
  color: rgba(255, 255, 255, 0.85);
}

.dark .project-detail-modal :deep(.ant-tabs-tab.ant-tabs-tab-active .ant-tabs-tab-btn) {
  color: #1890ff;
}

/* 列表项悬停效果 */
.ant-table-row:hover td {
  background-color: #f0f7ff !important;
}

.dark .ant-table-row:hover td {
  background-color: rgba(24, 144, 255, 0.1) !important;
}
</style>
