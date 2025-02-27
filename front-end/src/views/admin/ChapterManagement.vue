<script setup>
import { ref, onMounted } from 'vue'
import { get, postJSON } from '@/util/request'
import { message, Modal } from 'ant-design-vue'
import SpinLoaderLarge from '@/components/spinLoaderLarge.vue'
import { parseDateTime } from '@/util/common'
import { useRouter } from 'vue-router'

const router = useRouter()
const loading = ref(true)
const chapters = ref([])
const searchKeyword = ref('')
const selectedStatus = ref('全部')
const currentPage = ref(1)
const pageSize = ref(10)
const total = ref(0)

// 获取章节列表
const fetchChapters = async () => {
  loading.value = true
  
  get('/api/admin/chapters', {
    page: currentPage.value,
    pageSize: pageSize.value,
    keyword: searchKeyword.value,
    status: selectedStatus.value !== '全部' ? selectedStatus.value : ''
  }, 
  (msg, data) => {
    chapters.value = data.chapters
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

onMounted(() => {
  fetchChapters()
})

// 跳转到审核页面
const goToReview = (chapterId) => {
  router.push(`/admin/chapters/review/${chapterId}`)
}

// 搜索章节
const handleSearch = () => {
  currentPage.value = 1
  fetchChapters()
}

// 重置搜索
const handleReset = () => {
  searchKeyword.value = ''
  selectedStatus.value = '全部'
  currentPage.value = 1
  fetchChapters()
}

// 当页码变化时
const onPageChange = (page) => {
  currentPage.value = page
  fetchChapters()
}

// 当每页数量变化时
const onPageSizeChange = (current, size) => {
  currentPage.value = 1
  pageSize.value = size
  fetchChapters()
}

// 获取状态文本
const getStatusText = (status) => {
  switch (status) {
    case 'pending': return '待审核'
    case 'approved': return '已通过'
    case 'rejected': return '已拒绝'
    default: return '未知状态'
  }
}

// 获取状态颜色
const getStatusColor = (status) => {
  switch (status) {
    case 'pending': return 'processing'
    case 'approved': return 'success'
    case 'rejected': return 'error'
    default: return 'default'
  }
}

// 获取评分级别文本和颜色
const getScoreLevelText = (score) => {
  if (score >= 90) {
    return { text: '优秀', color: 'green' }
  } else if (score >= 70) {
    return { text: '良好', color: 'blue' }
  } else if (score >= 50) {
    return { text: '一般', color: 'orange' }
  } else if (score > 0) {
    return { text: '敏感', color: 'red' }
  } else {
    return { text: '暂无评分', color: 'default' }
  }
}

// 删除章节
const deleteChapter = (chapterId) => {
  Modal.confirm({
    title: '确认删除',
    content: '您确定要删除该章节吗？此操作不可撤销。',
    okText: '确认',
    cancelText: '取消',
    onOk: () => {
      postJSON(`/api/admin/chapters/${chapterId}/delete`, {}, 
        (msg, data) => {
          message.success('章节删除成功')
          fetchChapters()
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
</script>

<template>
  <div class="p-4">
    <h1 class="text-2xl font-bold mb-6 dark:text-gray-100">章节管理</h1>

    <div v-if="loading" class="flex justify-center items-center h-64">
      <SpinLoaderLarge />
    </div>

    <div v-else class="animate__animated animate__fadeIn">
      <!-- 搜索和筛选 -->
      <div class="bg-white dark:bg-zinc-800 p-4 rounded-lg shadow-sm mb-6">
        <div class="flex flex-wrap gap-4">
          <div class="flex-grow">
            <a-input
              v-model:value="searchKeyword"
              placeholder="搜索章节标题或描述"
              allow-clear
              @press-enter="handleSearch"
            >
              <template #prefix>
                <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
                </svg>
              </template>
            </a-input>
          </div>
          
          <a-select 
            v-model:value="selectedStatus" 
            style="min-width: 120px;"
            class="dark:bg-zinc-700 dark:text-white"
          >
            <a-select-option value="全部">全部状态</a-select-option>
            <a-select-option value="pending">待审核</a-select-option>
            <a-select-option value="approved">已通过</a-select-option>
            <a-select-option value="rejected">已拒绝</a-select-option>
          </a-select>
          
          <a-button type="primary" @click="handleSearch">搜索</a-button>
          <a-button @click="handleReset">重置</a-button>
        </div>
      </div>

      <!-- 章节列表 -->
      <a-table
        :loading="loading"
        :dataSource="chapters"
        :columns="[
          { title: '标题', dataIndex: 'Title', key: 'Title' },
          { title: '章节描述', dataIndex: 'Description', key: 'description' },
          { title: '项目ID', dataIndex: 'project_id', key: 'projectId' },
          { title: '评分', key: 'score' },
          { title: '状态', key: 'status' },
          { title: '更新时间', key: 'updatedAt' },
          { title: '操作', key: 'action' },
        ]"
        :pagination="false"
        rowKey="ID"
        class="bg-white dark:bg-zinc-800 rounded-lg shadow-sm"
      >
        <!-- 标题列 -->
        <template #bodyCell="{ column, record }">
          <template v-if="column.key === 'Title'">
            <div class="flex items-center">
              {{ record.Title }}
            </div>
          </template>
          
          <!-- 描述列 -->
          <template v-else-if="column.key === 'description'">
            <span class="line-clamp-1">{{ record.Description }}</span>
          </template>
          
          <!-- 评分列 -->
          <template v-else-if="column.key === 'score'">
            <div v-if="record.current_version && record.current_version.score !== undefined && record.current_version.score !== null">
              <a-tag :color="getScoreLevelText(record.current_version.score).color">
                {{ record.current_version.score }} - {{ getScoreLevelText(record.current_version.score).text }}
              </a-tag>
            </div>
            <span v-else class="text-gray-400 dark:text-gray-500">暂无评分</span>
          </template>
          
          <!-- 状态列 -->
          <template v-else-if="column.key === 'status'">
            <a-tag v-if="record.current_version && record.current_version.status" 
                   :color="getStatusColor(record.current_version.status)">
              {{ getStatusText(record.current_version.status) }}
            </a-tag>
            <a-tag v-else color="default">未审核</a-tag>
          </template>
          
          <!-- 更新时间列 -->
          <template v-else-if="column.key === 'updatedAt'">
            {{ parseDateTime(record.UpdatedAt) }}
          </template>
          
          <!-- 操作列 -->
          <template v-else-if="column.key === 'action'">
            <div class="flex space-x-2">
              <a-button 
                type="primary" 
                size="small" 
                @click="goToReview(record.ID)"
              >
                审核评分
              </a-button>
              
              <a-button 
                danger 
                size="small" 
                @click="deleteChapter(record.ID)"
              >
                删除
              </a-button>
            </div>
          </template>
        </template>
      </a-table>

      <!-- 分页 -->
      <div class="mt-4 flex justify-end">
        <a-pagination
          v-model:current="currentPage"
          :total="total"
          :pageSize="pageSize"
          :pageSizeOptions="['10', '20', '50']"
          showSizeChanger
          @change="onPageChange"
          @showSizeChange="onPageSizeChange"
        />
      </div>
    </div>
  </div>
</template>

<style scoped lang="postcss">
.animate__animated {
  animation-duration: 0.5s;
}

:deep(.ant-table-thead > tr > th) {
  @apply bg-gray-50 dark:bg-zinc-700 dark:text-gray-200 border-b dark:border-zinc-600;
}

:deep(.ant-table-tbody > tr > td) {
  @apply dark:border-zinc-700 dark:text-gray-300;
}

:deep(.ant-table-tbody > tr.ant-table-row:hover > td) {
  @apply bg-gray-50 dark:bg-zinc-700;
}

:deep(.ant-pagination-item) {
  @apply dark:bg-zinc-800 dark:border-zinc-700;
}

:deep(.ant-pagination-item a) {
  @apply dark:text-gray-300;
}

:deep(.ant-pagination-item-active) {
  @apply dark:bg-blue-900/30 dark:border-blue-700;
}
</style>
