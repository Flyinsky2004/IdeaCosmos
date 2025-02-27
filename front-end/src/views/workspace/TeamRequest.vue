<script>
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { message } from 'ant-design-vue'
import { get, postJSON } from "@/util/request.js"

export default {
  name: 'TeamRequest',
  setup() {
    const route = useRoute()
    const router = useRouter()
    const inviteCode = route.params.inviteCode
    
    const team = ref({})
    const memberCount = ref(0)
    const loading = ref(true)
    const error = ref('')
    const submitting = ref(false)

    const getTeamInfo = () => {
      loading.value = true
      get(`/api/team/getTeamByInviteCode`, {
        invite_code: inviteCode
      }, (msg, data) => {
        team.value = data.team
        memberCount.value = data.memberCount
        loading.value = false
      }, null, (msg) => {
        error.value = '无效的邀请链接'
        loading.value = false
      }, () => {
        loading.value = false
      })
    }

    const handleJoinRequest = () => {
      submitting.value = true
      const hide = message.loading('正在提交加入请求...', 0)
      
      postJSON('/api/team/joinByInviteCode', {
        invite_code: inviteCode
      }, (msg) => {
        hide()
        message.success('加入请求已提交，请等待团队管理员审核')
        router.push('/workspace/teams')
      }, (msg) => {
        hide()
        message.warning(msg)
        submitting.value = false
      }, (msg) => {
        hide()
        message.error(msg || '提交请求时发生错误')
        submitting.value = false
      })
    }

    const goToTeams = () => {
      router.push('/workspace/teams')
    }

    onMounted(() => {
      if (!inviteCode) {
        error.value = '无效的邀请链接'
        loading.value = false
        return
      }
      getTeamInfo()
    })

    return {
      team,
      memberCount,
      loading,
      error,
      submitting,
      handleJoinRequest,
      goToTeams
    }
  }
}
</script>
<template>
  <div class="team-request-container">
    <div class="request-card">
      <div v-if="loading" class="loading">
        <a-spin />
        <p class="mt-2 text-gray-500 dark:text-gray-400">正在获取团队信息...</p>
      </div>
      
      <div v-else-if="error" class="error-message">
        <a-alert
          :message="error"
          type="error"
          show-icon
        />
      </div>

      <div v-else class="team-info animate__animated animate__fadeIn">
        <!-- 欢迎横幅 -->
        <div class="welcome-banner">
          <div class="banner-content animate__animated animate__fadeInDown">
            <i class="fas fa-handshake text-4xl text-blue-500 mb-4"></i>
            <h1 class="text-2xl font-bold text-gray-800 dark:text-gray-200">欢迎加入我们！</h1>
            <p class="text-gray-600 dark:text-gray-400 mt-2">很高兴您对我们团队感兴趣</p>
          </div>
        </div>

        <!-- 团队基本信息卡片 -->
        <div class="team-card">
          <div class="team-header">
            <div class="team-title">
              <h2 class="text-xl font-bold text-blue-500 dark:text-blue-400">{{ team.username }}</h2>
              <span class="team-id">#{{ team.ID }}</span>
            </div>
            
            <div class="team-description">
              <h3 class="section-title">团队愿景</h3>
              <p class="description-content">
                {{ team.teamDescription }}
              </p>
            </div>
          </div>

          <!-- 团队统计信息 -->
          <div class="stats-grid">
            <div class="stat-card">
              <i class="fas fa-users text-2xl text-blue-500 mb-2"></i>
              <div class="stat-label">当前成员数</div>
              <div class="stat-value">
                {{ memberCount }}<span class="text-sm text-gray-400">/50</span>
              </div>
            </div>
            <div class="stat-card">
              <i class="fas fa-calendar text-2xl text-green-500 mb-2"></i>
              <div class="stat-label">成立时间</div>
              <div class="stat-value text-sm">
                {{ new Date(team.CreatedAt).toLocaleDateString() }}
              </div>
            </div>
          </div>

          <!-- 团队优势 -->
          <div class="advantages-section">
            <h3 class="section-title">加入我们，您将获得</h3>
            <div class="advantages-grid">
              <div class="advantage-card">
                <i class="fas fa-users-cog text-blue-500"></i>
                <h4>多人协作</h4>
                <p>与优秀的团队成员共同成长</p>
              </div>
              <div class="advantage-card">
                <i class="fas fa-comments text-green-500"></i>
                <h4>实时反馈</h4>
                <p>高效的沟通与反馈机制</p>
              </div>
              <div class="advantage-card">
                <i class="fas fa-code-branch text-purple-500"></i>
                <h4>版本控制</h4>
                <p>专业的项目管理体系</p>
              </div>
            </div>
          </div>
        
          <!-- 操作按钮 -->
          <div class="action-section">
            <p class="text-gray-600 dark:text-gray-400 text-center mb-4">
              准备好开启新的旅程了吗？
            </p>
            <div class="action-buttons">
              <a-button 
                type="primary"
                :loading="submitting"
                @click="handleJoinRequest"
                class="action-button"
              >
                <i class="fas fa-user-plus mr-1"></i>申请加入团队
              </a-button>
              
              <a-button 
                @click="goToTeams" 
                class="action-button"
              >
                <i class="fas fa-arrow-left mr-1"></i>返回团队列表
              </a-button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped lang="postcss">
.team-request-container {
  @apply min-h-[80vh] flex items-center justify-center p-4 bg-gray-50 dark:bg-gray-900;
}

.request-card {
  @apply bg-white dark:bg-gray-800 rounded-xl shadow-lg w-full max-w-3xl
         border dark:border-gray-700 overflow-hidden;
}

.welcome-banner {
  @apply bg-gradient-to-r from-blue-500 to-purple-500 p-8 text-center text-white;
}


.team-card {
  @apply p-6 space-y-8;
}

.team-header {
  @apply text-center space-y-4;
}

.team-title {
  @apply flex items-center justify-center gap-2;
}

.team-id {
  @apply text-xs bg-blue-100 text-blue-800 px-2 py-0.5 rounded-full 
         dark:bg-blue-900 dark:text-blue-200;
}

.section-title {
  @apply text-sm font-bold text-gray-500 mb-2 dark:text-gray-400;
}

.description-content {
  @apply text-gray-600 mb-3 border-l-4 border-blue-300 pl-3 py-2 
         bg-blue-50 rounded-r dark:text-gray-300 dark:bg-blue-900/30 
         dark:border-blue-700;
}

.stats-grid {
  @apply grid grid-cols-2 gap-4;
}

.stat-card {
  @apply text-center p-4 bg-gray-50 rounded-lg 
         dark:bg-gray-800/50 transition-transform hover:scale-105;
}

.stat-label {
  @apply text-sm text-gray-500 dark:text-gray-400;
}

.stat-value {
  @apply text-xl font-bold text-blue-600 dark:text-blue-400;
}

.advantages-section {
  @apply text-center;
}

.advantages-grid {
  @apply grid grid-cols-1 md:grid-cols-3 gap-4 mt-4;
}

.advantage-card {
  @apply p-4 rounded-lg bg-gray-50 dark:bg-gray-800/50
         transition-all hover:shadow-md hover:-translate-y-1;
}

.advantage-card i {
  @apply text-2xl mb-2;
}

.advantage-card h4 {
  @apply font-bold text-gray-700 dark:text-gray-300 mb-1;
}

.advantage-card p {
  @apply text-sm text-gray-600 dark:text-gray-400;
}

.action-section {
  @apply mt-8 border-t pt-6 dark:border-gray-700;
}

.action-buttons {
  @apply flex flex-col md:flex-row gap-3 justify-center;
}

.action-button {
  @apply w-full md:w-[180px] text-lg h-12 px-8 flex items-center justify-center;
}

.loading, .error-message {
  @apply text-center py-12;
}
</style>
