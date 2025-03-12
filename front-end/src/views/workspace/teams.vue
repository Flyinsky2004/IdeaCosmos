<script setup>
import {onMounted, reactive, ref} from "vue";
import {get, postJSON} from "@/util/request.js";
import {message, Spin, Badge} from "ant-design-vue";
import {parseDateTime} from "@/util/common.js";
import { BACKEND_DOMAIN, FRONTEND_DOMAIN } from "@/util/VARRIBLES";

const options = reactive({
  myTeams: [],
  myJoinedTeams: [],
  isAddWindowOpen: false,
  isJoinWindowOpen: false,
  isTeamDetailOpen: false,
  currentTeam: null,
  teamMembers: [],
  pendingRequests: [],
  loading: {
    myTeams: true,
    joinedTeams: true,
    teamDetail: false,
    teamMembers: false,
    pendingRequests: false
  }
})

const inviteCode = ref('');
const selectedTeamId = ref(null);

const fetchMyTeams = () => {
  options.loading.myTeams = true;
  get('/api/team/getMyTeams', {},
      (message, data) => {
        if(data !== null) options.myTeams = data; else options.myTeams = [];
        options.loading.myTeams = false;
      }, null, null, () => {
        options.loading.myTeams = false;
      })
}

const fetchMyJoinedTeams = () => {
  options.loading.joinedTeams = true;
  get('/api/team/getMyJoinedTeams', {},
      (message, data) => {
        if(data !== null) options.myJoinedTeams = data;
        options.loading.joinedTeams = false;
      }, null, null, () => {
        options.loading.joinedTeams = false;
      })
}

onMounted(() => {
  fetchMyTeams()
  fetchMyJoinedTeams()
})

const newTeamForm = reactive({
  team_name: '',
  team_description: '',
})

const validateTeamForm = () => {
  if (!newTeamForm.team_name || newTeamForm.team_name.length < 6) {
    message.error("团队名称长度必须大于6个字符");
    return false;
  }
  
  if (!newTeamForm.team_description || newTeamForm.team_description.length < 20) {
    message.error("团队描述必须大于20个字符");
    return false;
  }
  
  return true;
}

const submitNewTeam = () => {
  if (!validateTeamForm()) {
    return;
  }
  
  const hide = message.loading("正在创建团队...", 0);
  
  postJSON('/api/team/createTeam', newTeamForm,
      (messager, data) => {
        hide();
        message.success(messager)
        fetchMyTeams()
        options.isAddWindowOpen = false
        newTeamForm.team_name = '';
        newTeamForm.team_description = '';
      },
      (messager, data) => {
        hide();
        message.warning(messager)
      },
      (messager, data) => {
        hide();
        message.error(messager)
      })
}

// 通过邀请码加入团队
const joinTeamByInviteCode = () => {
  if (!inviteCode.value || inviteCode.value.length !== 8) {
    message.error("请输入有效的8位邀请码");
    return;
  }
  
  const hide = message.loading("正在提交加入请求...", 0);
  
  postJSON('/api/team/joinByInviteCode', {
    invite_code: inviteCode.value
  }, (msg, data) => {
    hide();
    message.success(msg);
    fetchMyJoinedTeams();
    options.isJoinWindowOpen = false;
    inviteCode.value = '';
  }, (msg) => {
    hide();
    message.warning(msg);
  }, (msg) => {
    hide();
    message.error(msg);
  });
}

// 查看团队详情
const viewTeamDetail = (teamId) => {
  selectedTeamId.value = teamId;
  options.loading.teamDetail = true;
  get('/api/team/detail', {
    team_id: teamId
  }, (msg, data) => {
    options.currentTeam = data;
    options.isTeamDetailOpen = true;
    options.loading.teamDetail = false;
    
    // 获取团队成员
    fetchTeamMembers(teamId);
    
    // 如果是团队创建者，获取待处理的加入请求
    if (data.is_leader) {
      fetchPendingRequests(teamId);
    }
  }, null, null, () => {
    options.loading.teamDetail = false;
  });
}

// 获取团队成员
const fetchTeamMembers = (teamId) => {
  options.loading.teamMembers = true;
  get('/api/team/members', {
    team_id: teamId
  }, (msg, data) => {
    options.teamMembers = data;
    options.loading.teamMembers = false;
  }, null, null, () => {
    options.loading.teamMembers = false;
  });
}

// 获取待处理的加入请求
const fetchPendingRequests = (teamId) => {
  options.loading.pendingRequests = true;
  get('/api/team/getRequests', {
    status: 0 // 待处理状态
  }, (msg, data) => {
    options.pendingRequests = data.filter(req => req.teamId === teamId);
    options.loading.pendingRequests = false;
  }, null, null, () => {
    options.loading.pendingRequests = false;
  });
}

// 处理加入请求
const handleJoinRequest = (requestId, status) => {
  const statusText = status === 1 ? "批准" : "拒绝";
  const hide = message.loading(`正在${statusText}加入请求...`, 0);
  postJSON('/api/team/updateRequest', {
    requestId: requestId,
    status: status
  }, (msg) => {
    hide();
    message.success(status === 1 ? "已批准加入请求" : "已拒绝加入请求");
    fetchPendingRequests(selectedTeamId.value);
    fetchTeamMembers(selectedTeamId.value);
  }, null, (msg) => {
    hide();
    message.error(`${statusText}请求失败: ${msg}`);
  });
}

// 重新生成邀请码
const regenerateInviteCode = () => {
  const hide = message.loading("正在生成新的邀请码...", 0);
  get('/api/team/regenerateInviteCode', {
    team_id: selectedTeamId.value
  }, (msg, data) => {
    hide();
    message.success("邀请码已重新生成");
    options.currentTeam.invite_code = data;
  }, null, (msg) => {
    hide();
    message.error(`生成邀请码失败: ${msg}`);
  });
}

// 复制邀请码到剪贴板
const copyInviteCode = () => {
  navigator.clipboard.writeText(options.currentTeam.invite_code)
      .then(() => {
        message.success("邀请码已复制到剪贴板");
      })
      .catch(() => {
        message.error("复制失败，请手动复制");
      });
}
</script>

<template>
  <!-- 创建新团队对话框 -->
  <a-modal v-model:open="options.isAddWindowOpen" title="创建新团队" @ok="submitNewTeam" ok-text="创建"
           cancel-text="取消">
    <div class="font-sans grid gap-4">
      <div>
        <div class="flex justify-between">
          <label class="text-sm font-medium mb-1 dark:text-gray-300">团队名称</label>
          <span class="text-xs text-gray-500 dark:text-gray-400">至少6个字符</span>
        </div>
        <input class="input1 " v-model="newTeamForm.team_name" placeholder="请输入团队名称">
      </div>
      <div>
        <div class="flex justify-between">
          <label class="text-sm font-medium mb-1 dark:text-gray-300">团队描述</label>
          <span class="text-xs text-gray-500 dark:text-gray-400">至少20个字符</span>
        </div>
        <textarea class="input1 min-h-24" v-model="newTeamForm.team_description" placeholder="请描述团队的目标和工作内容"/>
      </div>
    </div>
  </a-modal>
  
  <!-- 通过邀请码加入团队对话框 -->
  <a-modal v-model:open="options.isJoinWindowOpen" title="通过邀请码加入团队" @ok="joinTeamByInviteCode" ok-text="加入"
           cancel-text="取消">
    <div class="font-sans grid gap-4">
      <div>
        <label class="text-sm font-medium mb-1 block dark:text-gray-300">邀请码</label>
        <input class="input1 font-mono text-center tracking-wider text-lg dark:bg-gray-700 dark:text-gray-200 dark:border-gray-600" v-model="inviteCode" placeholder="请输入8位邀请码" maxlength="8">
        <p class="text-sm text-gray-500 mt-2 dark:text-gray-400">请向团队管理员获取8位邀请码</p>
      </div>
    </div>
  </a-modal>
  
  <!-- 团队详情对话框 -->
  <a-modal v-model:open="options.isTeamDetailOpen" title="团队详情" width="700px" footer="" cancel-text="关闭">
    <Spin :spinning="options.loading.teamDetail"/>
      <div v-if="options.currentTeam" class="font-sans">
      <div class="mb-4">
        <div class="flex items-center gap-2 mb-2">
          <h2 class="text-xl font-bold text-blue-500 dark:text-blue-400">{{ options.currentTeam.team.username }}</h2>
          <span class="text-xs bg-blue-100 text-blue-800 px-2 py-0.5 rounded-full dark:bg-blue-900 dark:text-blue-200">#{{ options.currentTeam.team.ID }}</span>
        </div>
        <p class="text-gray-600 mb-3 border-l-4 border-blue-300 pl-3 py-2 bg-blue-50 rounded-r dark:text-gray-300 dark:bg-blue-900/30 dark:border-blue-700">{{ options.currentTeam.team.teamDescription }}</p>
        <div class="flex justify-between text-sm text-gray-500 dark:text-gray-400">
          <p>创建时间: {{ parseDateTime(options.currentTeam.team.CreatedAt) }}</p>
          <p>成员数量: <span class="font-medium">{{ options.currentTeam.member_count }}/50</span></p>
        </div>
      </div>
      
      <!-- 邀请码区域 (仅团队创建者可见) -->
      <div v-if="options.currentTeam.is_leader" class="mb-4 p-3 bg-gray-50 rounded-lg dark:bg-gray-800">
        <div class="flex justify-between items-center">
          <div>
            <h3 class="font-bold dark:text-gray-200">团队邀请码</h3>
            <p class="text-lg font-mono tracking-wider bg-white p-2 rounded border mt-1 dark:bg-gray-700 dark:border-gray-600 dark:text-gray-200">{{ options.currentTeam.invite_code }}</p>
          </div>
          <div class="flex gap-2">
            <a-button type="primary" size="small" @click="copyInviteCode">
              <i class="fas fa-copy mr-1"></i> 复制
            </a-button>
            <a-button size="small" @click="regenerateInviteCode">
              <i class="fas fa-sync-alt mr-1"></i> 重新生成
            </a-button>
          </div>
        </div>
        <p class="text-xs text-gray-500 mt-2 dark:text-gray-400">分享此邀请码给您想邀请的成员</p>
      </div>
      
      <!-- 团队成员列表 -->
      <div class="mb-4">
        <div class="flex justify-between items-center mb-2">
          <h3 class="font-bold dark:text-gray-200">团队成员</h3>
          <Badge :count="options.teamMembers.length" :overflowCount="99" />
        </div>
        <Spin :spinning="options.loading.teamMembers">
          <div class="grid grid-cols-2 gap-2">
            <div v-for="member in options.teamMembers" :key="member.user.ID" 
                 class="flex items-center p-2 border rounded-lg hover:bg-gray-50 transition-colors dark:border-gray-700 dark:hover:bg-gray-800">
              <img :src="BACKEND_DOMAIN + member.user.avatar || '/default-avatar.png'" class="w-10 h-10 rounded-full mr-2 object-cover border dark:border-gray-600" alt="头像">
              <div>
                <p class="font-bold dark:text-gray-200">{{ member.user.username }}</p>
                <p v-if="member.is_leader" class="text-xs bg-yellow-100 text-yellow-800 px-2 py-0.5 rounded-full inline-block dark:bg-yellow-900 dark:text-yellow-200">团队创建者</p>
                <p v-else class="text-xs bg-green-100 text-green-800 px-2 py-0.5 rounded-full inline-block dark:bg-green-900 dark:text-green-200">团队成员</p>
              </div>
            </div>
          </div>
        </Spin>
      </div>
      
      <!-- 待处理的加入请求 (仅团队创建者可见) -->
      <div v-if="options.currentTeam.is_leader" class="mb-4">
        <div class="flex justify-between items-center mb-2">
          <h3 class="font-bold dark:text-gray-200">待处理的加入请求</h3>
          <Badge :count="options.pendingRequests.length" :overflowCount="99" />
        </div>
        <Spin :spinning="options.loading.pendingRequests">
          <div v-if="options.pendingRequests.length === 0" class="text-center py-4 text-gray-500 dark:text-gray-400">
            暂无待处理的请求
          </div>
          <div v-else class="space-y-2">
            <div v-for="request in options.pendingRequests" :key="request.ID" 
                 class="flex justify-between items-center p-3 border rounded-lg hover:bg-gray-50 transition-colors dark:border-gray-700 dark:hover:bg-gray-800">
              <div class="flex items-center">
                <img :src="BACKEND_DOMAIN + request.user.avatar || '/default-avatar.png'" class="w-8 h-8 rounded-full mr-2 object-cover border dark:border-gray-600" alt="头像">
                <div>
                  <p class="font-medium dark:text-gray-200">{{ request.user.username }}</p>
                  <p class="text-xs text-gray-500 dark:text-gray-400">申请时间: {{ parseDateTime(request.CreatedAt) }}</p>
                </div>
              </div>
              <div class="flex gap-2">
                <a-button type="primary" size="small" @click="handleJoinRequest(request.ID, 1)">
                  <i class="fas fa-check mr-1"></i> 批准
                </a-button>
                <a-button danger size="small" @click="handleJoinRequest(request.ID, 2)">
                  <i class="fas fa-times mr-1"></i> 拒绝
                </a-button>
              </div>
            </div>
          </div>
        </Spin>
      </div>
    </div>
  </a-modal>
  
  <div class="flex flex-col gap-2"
       v-motion
       :initial="{ opacity: 0, y: 20 }"
       :enter="{ opacity: 1, y: 0, transition: { duration: 500 } }">
    <div class="border p-4 workspace-box w-fit rounded-lg bg-gradient-to-r from-blue-50 to-indigo-50 dark:from-blue-900/30 dark:to-indigo-900/30 dark:border-gray-700"
         v-motion
         :initial="{ opacity: 0, scale: 0.9 }"
         :enter="{ 
           opacity: 1, 
           scale: 1,
           transition: {
             type: 'spring',
             damping: 12,
             stiffness: 100
           }
         }"
         :hover="{ scale: 1.02, transition: { duration: 200 } }">
      <h1 class="text-2xl font-serif mb-2 text-blue-700 dark:text-blue-400">团队工作</h1>
      <span class="text-sm dark:text-gray-300">我们的团队功能支持您将创剧空间中的项目与您的团队绑定，使您的团队成员一同加入您的精彩内容创作！</span>
      <div class="mt-2 flex flex-wrap gap-3">
        <div class="text-xs bg-blue-100 text-blue-800 px-2 py-1 rounded-full dark:bg-blue-900 dark:text-blue-200">项目多人分工协作</div>
        <div class="text-xs bg-green-100 text-green-800 px-2 py-1 rounded-full dark:bg-green-900 dark:text-green-200">实时评论反馈</div>
        <div class="text-xs bg-purple-100 text-purple-800 px-2 py-1 rounded-full dark:bg-purple-900 dark:text-purple-200">项目版本控制</div>
      </div>
    </div>
    
    <!-- 我管理的团队 -->
    <div class="border p-4 workspace-box dark:border-gray-700"
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
      <div class="flex justify-between items-center mb-4 border-b pb-3 dark:border-gray-700">
        <h1 class="text-2xl font-medium text-blue-700 dark:text-blue-400">我管理的团队</h1>
        <div class="flex gap-2">
          <a-button @click="options.isJoinWindowOpen = true">
            <i class="fas fa-ticket-alt mr-1"></i> 通过邀请码加入
          </a-button>
          <a-button type="primary" @click="options.isAddWindowOpen = true">
            <i class="fas fa-plus mr-1"></i> 创建团队
          </a-button>
        </div>
      </div>
      <Spin :spinning="options.loading.myTeams">
        <div class="w-full grid grid-cols-1 md:grid-cols-3 lg:grid-cols-5 gap-3">
          <div v-for="(team, index) in options.myTeams" :key="team.ID" 
               class="border rounded-lg shadow-sm outline-[1px] min-h-56 text-theme-switch p-4
               hover:shadow-md hover:-translate-y-1 transition-all cursor-pointer
               dark:border-gray-700 dark:bg-gray-800/50 dark:shadow-gray-900/30"
               v-motion
               :initial="{ opacity: 0, y: 20 }"
               :enter="{ 
                 opacity: 1, 
                 y: 0,
                 transition: { 
                   delay: 300 + (index * 100),
                   duration: 500
                 }
               }"
               :hover="{ scale: 1.02, transition: { duration: 200 } }"
               @click="viewTeamDetail(team.ID)">
            <div class="flex justify-between mb-1">
              <h1 class="text-lg font-bold text-blue-500 truncate dark:text-blue-400">{{ team.username }}</h1>
              <span class="text-xs px-1.5 py-0.5 rounded bg-blue-100 text-blue-800 dark:bg-blue-900 dark:text-blue-200">#{{ team.ID }}</span>
            </div>
            <div class="border-t border-b my-2 py-2 dark:border-gray-700">
              <h1 class="text-xs font-bold text-gray-500 mb-1 dark:text-gray-400">团队介绍：</h1>
              <div class="text-sm line-clamp-3 text-gray-600 dark:text-gray-300">{{ team.teamDescription }}</div>
            </div>
            <div class="text-xs text-gray-500 space-y-1 dark:text-gray-400">
              <div class="flex justify-between">
                <span class="font-medium">创建时间：</span>
                <span>{{ parseDateTime(team.CreatedAt).split(' ')[0] }}</span>
              </div>
              <div class="flex justify-between">
                <span class="font-medium">团队人数：</span>
                <span class="text-blue-600 font-bold dark:text-blue-400">{{team.member_count}}/50</span>
              </div>
            </div>
            <div class="mt-3 pt-1">
              <span class="text-xs bg-yellow-100 text-yellow-800 px-2 py-0.5 rounded-full inline-flex items-center dark:bg-yellow-900 dark:text-yellow-200">
                <i class="fas fa-crown mr-1"></i> 团队创建者
              </span>
            </div>
          </div>
          <div class="border border-dashed rounded-lg outline-[1px] min-h-56 flex flex-col items-center justify-center text-gray-400
               hover:bg-gray-50 active:bg-gray-100 cursor-pointer
               dark:border-gray-700 dark:hover:bg-gray-800 dark:active:bg-gray-700"
               v-motion
               :initial="{ opacity: 0, scale: 0.9 }"
               :enter="{ 
                 opacity: 1, 
                 scale: 1,
                 transition: { 
                   delay: 500,
                   type: 'spring',
                   damping: 12,
                   stiffness: 100
                 }
               }"
               :hover="{ scale: 1.05, transition: { duration: 200 } }"
               @click="options.isAddWindowOpen = true">
            <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5"
                 stroke="currentColor" class="size-12 mb-2">
              <path stroke-linecap="round" stroke-linejoin="round" d="M12 9v6m3-3H9m12 0a9 9 0 1 1-18 0 9 9 0 0 1 18 0Z"/>
            </svg>
            <span class="text-sm">创建新团队</span>
          </div>
        </div>
      </Spin>
    </div>
    
    <!-- 我加入的团队 -->
    <div class="border p-4 workspace-box dark:border-gray-700"
         v-motion
         :initial="{ opacity: 0, y: 20 }"
         :enter="{ 
           opacity: 1, 
           y: 0,
           transition: { 
             delay: 400,
             duration: 500
           }
         }">
      <h1 class="text-2xl font-medium text-blue-700 mb-4 border-b pb-3 dark:text-blue-400 dark:border-gray-700">我加入的团队</h1>
      <Spin :spinning="options.loading.joinedTeams">
        <div v-if="!options.myJoinedTeams || options.myJoinedTeams.length === 0" 
             class="text-center py-8 text-gray-500 border border-dashed rounded-lg dark:text-gray-400 dark:border-gray-700"
             v-motion
             :initial="{ opacity: 0, scale: 0.9 }"
             :enter="{ 
               opacity: 1, 
               scale: 1,
               transition: { 
                 delay: 600,
                 type: 'spring',
                 damping: 12,
                 stiffness: 100
               }
             }">
          <i class="fas fa-users text-3xl mb-2 opacity-30"></i>
          <p>您还没有加入任何团队</p>
          <a-button class="mt-4" @click="options.isJoinWindowOpen = true">
            <i class="fas fa-ticket-alt mr-1"></i> 通过邀请码加入团队
          </a-button>
        </div>
        <div v-else class="w-full grid grid-cols-1 md:grid-cols-3 lg:grid-cols-5 gap-3">
          <div v-for="(team, index) in options.myJoinedTeams" :key="team.ID" 
               class="border rounded-lg shadow-sm outline-[1px] min-h-56 text-theme-switch p-4
               hover:shadow-md hover:-translate-y-1 transition-all cursor-pointer
               dark:border-gray-700 dark:bg-gray-800/50 dark:shadow-gray-900/30"
               v-motion
               :initial="{ opacity: 0, y: 20 }"
               :enter="{ 
                 opacity: 1, 
                 y: 0,
                 transition: { 
                   delay: 600 + (index * 100),
                   duration: 500
                 }
               }"
               :hover="{ scale: 1.02, transition: { duration: 200 } }"
               @click="viewTeamDetail(team.ID)">
            <div class="flex justify-between mb-1">
              <h1 class="text-lg font-bold text-blue-500 truncate dark:text-blue-400">{{ team.username }}</h1>
              <span class="text-xs px-1.5 py-0.5 rounded bg-blue-100 text-blue-800 dark:bg-blue-900 dark:text-blue-200">#{{ team.ID }}</span>
            </div>
            <div class="border-t border-b my-2 py-2 dark:border-gray-700">
              <h1 class="text-xs font-bold text-gray-500 mb-1 dark:text-gray-400">团队介绍：</h1>
              <div class="text-sm line-clamp-3 text-gray-600 dark:text-gray-300">{{ team.teamDescription }}</div>
            </div>
            <div class="text-xs text-gray-500 space-y-1 dark:text-gray-400">
              <div class="flex justify-between">
                <span class="font-medium">创建时间：</span>
                <span>{{ parseDateTime(team.CreatedAt).split(' ')[0] }}</span>
              </div>
              <div class="flex justify-between">
                <span class="font-medium">团队人数：</span>
                <span class="text-blue-600 font-bold dark:text-blue-400">{{team.member_count}}/50</span>
              </div>
            </div>
            <div class="mt-3 pt-1">
              <span class="text-xs bg-green-100 text-green-800 px-2 py-0.5 rounded-full inline-flex items-center dark:bg-green-900 dark:text-green-200">
                <i class="fas fa-user-check mr-1"></i> 团队成员
              </span>
            </div>
          </div>
        </div>
      </Spin>
    </div>
  </div>
</template>

<style scoped>
.input1 {
  @apply border rounded p-2 w-full focus:outline-none focus:ring-2 focus:ring-blue-300 focus:border-blue-500 transition-all dark:focus:ring-blue-700 dark:focus:border-blue-600;
}

.line-clamp-3 {
  display: -webkit-box;
  -webkit-line-clamp: 3;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

/* 悬停效果 */
.hover\:shadow-md {
  transition: all 0.2s ease-in-out;
}

.hover\:-translate-y-1:hover {
  transform: translateY(-4px);
}

/* 响应式适配 */
@media (max-width: 768px) {
  .md\:grid-cols-3 {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }
}
</style>