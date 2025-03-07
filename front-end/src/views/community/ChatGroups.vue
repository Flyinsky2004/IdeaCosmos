<!--
 * @Author: Flyinsky w2084151024@gmail.com
 * @Description: 群组聊天页面
-->
<script setup>
import { ref, computed, onMounted } from 'vue'
import { get, post, postJSON } from '@/util/request'
import { message, Modal, Input, Dropdown, Menu, Tooltip, Checkbox, Table, Tabs } from 'ant-design-vue'
import GroupChatList from '@/components/chat/GroupChatList.vue'
import GroupChatMessages from '@/components/chat/GroupChatMessages.vue'
import GroupMemberList from '@/components/chat/GroupMemberList.vue'
import { useUserStore } from '@/stores/user'
import { BACKEND_DOMAIN } from '@/util/VARRIBLES'

const userStore = useUserStore()
const activeGroup = ref(null)
const members = ref([])
const showMembers = ref(false)
const isCreator = computed(() => activeGroup.value?.CreatorID === userStore.user?.id)
const isAdmin = ref(false)
const settingsVisible = ref(false)
const inviteModalVisible = ref(false)
const inviteEmail = ref('')
const inviteActiveTab = ref('userList')
const joinModalVisible = ref(false)
const joinGroupCode = ref('')
const groupSettingsForm = ref({
  name: '',
  description: '',
  avatarUrl: ''
})
const allUsers = ref([])
const selectedUserIds = ref([])
const userTableLoading = ref(false)
const showBackToTop = ref(false)

// 回到顶部
const scrollToTop = (elementSelector) => {
  const element = document.querySelector(elementSelector)
  if (element) {
    element.scrollTop = 0
    showBackToTop.value = false
  }
}

// 监听滚动
const handleScroll = (event) => {
  showBackToTop.value = event.target.scrollTop > 200
}

// 获取群组成员
const fetchGroupMembers = () => {
  if (!activeGroup.value) return
  
  get(`/api/chat/groups/${activeGroup.value.ID}/members`, {}, (_, data) => {
    members.value = data.members
    // 检查当前用户是否为管理员
    const currentUserMember = members.value.find(m => m.UserID === userStore.user.id)
    isAdmin.value = currentUserMember?.IsAdmin || isCreator.value
  },(msg) => {
    message.error('获取成员失败' + msg)
  },(msg) => {
    message.error('获取成员失败' + msg)
  })
}

// 选择群组时，获取群组详情和成员
const handleSelectGroup = (group) => {
  // 获取群组详情
  get(`/api/chat/groups/${group.ID}`, {}, (_, data) => {
    activeGroup.value = data.group
    fetchGroupMembers()
  })
}

// 创建群组
const handleCreateGroup = (group) => {
  postJSON('/api/chat/groups', group, (_, data) => {
    message.success('群组创建成功')
    activeGroup.value = data
    fetchGroupMembers()
  })
}

// 打开群组设置
const openGroupSettings = () => {
  if (!activeGroup.value) return
  
  groupSettingsForm.value = {
    name: activeGroup.value.Name,
    description: activeGroup.value.Description || '',
    avatarUrl: activeGroup.value.AvatarURL || ''
  }
  
  settingsVisible.value = true
}

// 更新群组信息
const updateGroupSettings = () => {
  if (!activeGroup.value) return
  
  put(`/api/chat/groups/${activeGroup.value.ID}`, groupSettingsForm.value, () => {
    message.success('群组信息更新成功')
    // 更新本地群组数据
    activeGroup.value = {
      ...activeGroup.value,
      Name: groupSettingsForm.value.name,
      Description: groupSettingsForm.value.description,
      AvatarURL: groupSettingsForm.value.avatarUrl
    }
    settingsVisible.value = false
  })
}

// 打开邀请成员模态框
const openInviteModal = () => {
  inviteEmail.value = ''
  selectedUserIds.value = []
  inviteModalVisible.value = true
  fetchAllUsers()
}

// 获取所有用户列表
const fetchAllUsers = () => {
  userTableLoading.value = true
  get('/api/user/all', {}, (_, data) => {
    allUsers.value = data.users.filter(user => {
      // 排除当前群组成员
      const memberIds = members.value.map(m => m.UserID)
      return !memberIds.includes(user.id)
    })
    userTableLoading.value = false
  }, (msg) => {
    message.error('获取用户列表失败' + msg)
    userTableLoading.value = false
  })
}

// 邀请成员加入群组
const inviteMember = () => {
  if (!activeGroup.value) return
  
  // 通过邮箱邀请的情况
  if (inviteActiveTab.value === 'email') {
    if (!inviteEmail.value) {
      message.warning('请输入邮箱地址')
      return
    }
    
    post(`/api/chat/groups/${activeGroup.value.ID}/members`, {
      email: inviteEmail.value
    }, () => {
      message.success('邀请已发送')
      inviteModalVisible.value = false
      fetchGroupMembers()
    })
    return
  }
  
  // 通过选择用户列表邀请的情况
  if (selectedUserIds.value.length === 0) {
    message.warning('请选择要邀请的用户')
    return
  }
  
  post(`/api/chat/groups/${activeGroup.value.ID}/members`, {
    userIds: selectedUserIds.value
  }, () => {
    message.success('邀请已发送')
    inviteModalVisible.value = false
    fetchGroupMembers()
  }, (msg) => {
    message.error('邀请失败: ' + msg)
  })
}

// 打开加入群组模态框
const openJoinModal = () => {
  joinGroupCode.value = ''
  joinModalVisible.value = true
}

// 加入群组
const joinGroup = () => {
  if (!joinGroupCode.value) return
  
  // 发送加入请求
  post(`/api/chat/groups/join`, {
    code: joinGroupCode.value
  }, (_, data) => {
    message.success('成功加入群组')
    joinModalVisible.value = false
    // 设置当前活动群组
    activeGroup.value = data
    fetchGroupMembers()
  })
}

// 离开群组
const leaveGroup = () => {
  if (!activeGroup.value) return
  
  Modal.confirm({
    title: '确认离开群组',
    content: '确定要离开该群组吗？',
    okText: '确认',
    cancelText: '取消',
    onOk: () => {
      post(`/api/chat/groups/${activeGroup.value.ID}/leave`, {}, () => {
        message.success('已离开群组')
        activeGroup.value = null
      })
    }
  })
}

// 删除群组
const deleteGroup = () => {
  if (!activeGroup.value || !isCreator.value) return
  
  Modal.confirm({
    title: '确认删除群组',
    content: '此操作不可恢复，确定要删除该群组吗？',
    okText: '删除',
    okType: 'danger',
    cancelText: '取消',
    onOk: () => {
      del(`/api/chat/groups/${activeGroup.value.ID}`, {}, () => {
        message.success('群组已删除')
        activeGroup.value = null
      })
    }
  })
}

// 禁言/解禁成员
const toggleMuteMember = (member, isMuted) => {
  if (!activeGroup.value || !isAdmin.value) return
  
  post(`/api/chat/groups/${activeGroup.value.ID}/members/${member.UserID}/mute`, {
    status: isMuted ? 2 : 1 // 状态：1-正常，2-禁言
  }, () => {
    message.success(isMuted ? '已禁言该成员' : '已解除禁言')
    fetchGroupMembers()
  })
}

// 移除成员
const removeMember = (member) => {
  if (!activeGroup.value || !isAdmin.value) return
  
  Modal.confirm({
    title: '确认移除成员',
    content: `确定要将 ${member.Nickname || member.Username} 移出群组吗？`,
    okText: '移除',
    okType: 'danger',
    cancelText: '取消',
    onOk: () => {
      del(`/api/chat/groups/${activeGroup.value.ID}/members/${member.UserID}`, {}, () => {
        message.success('成员已移除')
        fetchGroupMembers()
      })
    }
  })
}

// 设置/取消管理员
const toggleAdminStatus = (member, isAdminStatus) => {
  if (!activeGroup.value || !isCreator.value) return
  
  post(`/api/chat/groups/${activeGroup.value.ID}/members/${member.UserID}/admin`, {
    isAdmin: isAdminStatus
  }, () => {
    message.success(isAdminStatus ? '已设置为管理员' : '已取消管理员权限')
    fetchGroupMembers()
  })
}

// 更新成员昵称
const updateMemberNickname = (member, nickname) => {
  if (!activeGroup.value) return
  
  put(`/api/chat/groups/${activeGroup.value.ID}/members/${member.UserID}`, {
    nickname: nickname
  }, () => {
    message.success('昵称已更新')
    fetchGroupMembers()
  })
}

// 用户列表表格列定义
const userColumns = [
  {
    title: '选择',
    dataIndex: 'select',
    width: 60,
    slots: { 
      customRender: 'select' 
    }
  },
  {
    title: '用户名',
    dataIndex: 'username',
  },
  {
    title: '头像',
    dataIndex: 'avatar',
    width: 80,
    slots: { 
      customRender: 'avatar' 
    }
  }
]
</script>

<template>
  <div class="flex flex-col h-[calc(100vh-5rem)] animate__animated animate__fadeIn">

    <!-- 聊天主区域 -->
    <div class="flex-1 mx-4 mb-4 bg-white dark:bg-zinc-900 rounded-2xl border theme-border shadow-sm overflow-hidden animate__animated animate__fadeIn animate__delay-2s">
      <div class="h-full flex">
        <!-- 左侧群组列表 -->
        <div class="w-80 flex-shrink-0 border-r theme-border">
          <GroupChatList 
            :active-group-id="activeGroup?.ID"
            @select-group="handleSelectGroup"
            @create-group="handleCreateGroup"
          />
        </div>
        
        <!-- 右侧聊天区域 -->
        <div class="flex-1 flex flex-col">
          <!-- 群组信息栏 -->
          <div v-if="activeGroup" class="flex items-center justify-between px-4 py-3 border-b theme-border">
            <div class="flex items-center gap-3">
              <div class="w-10 h-10 rounded-full overflow-hidden bg-gradient-to-br from-blue-100 to-indigo-100 dark:from-blue-900/30 dark:to-indigo-900/30 flex-shrink-0 border theme-border">
                <img 
                  v-if="activeGroup.AvatarURL" 
                  :src="activeGroup.AvatarURL" 
                  class="w-full h-full object-cover"
                  alt="群组头像"
                >
                <div v-else class="w-full h-full flex items-center justify-center text-blue-500 dark:text-blue-400 text-sm font-bold">
                  {{ activeGroup.Name.charAt(0) }}
                </div>
              </div>
              <div>
                <h3 class="font-bold text-gray-800 dark:text-gray-200">{{ activeGroup.Name }}</h3>
                <p class="text-xs text-gray-500 dark:text-gray-400">{{ activeGroup.MemberCount }}人</p>
              </div>
            </div>
            <div class="flex items-center gap-2">
              <!-- 邀请成员按钮 - 仅管理员可见 -->
              <button 
                @click="openInviteModal" 
                class="px-3 py-1.5 text-sm text-white bg-blue-500 hover:bg-blue-600 rounded-lg flex items-center gap-1.5 transition-colors"
              >
                <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" d="M18 9v3m0 0v3m0-3h3m-3 0h-3m-2-5a4 4 0 11-8 0 4 4 0 018 0zM3 20a6 6 0 0112 0v1H3v-1z" />
                </svg>
                <span>邀请成员</span>
              </button>
              
              <!-- 成员列表按钮 -->
              <Tooltip title="成员管理">
                <button 
                  @click="showMembers = !showMembers" 
                  class="p-2 text-gray-500 hover:text-blue-500 dark:text-gray-400 dark:hover:text-blue-400 transition-colors rounded-lg hover:bg-gray-100 dark:hover:bg-gray-800"
                >
                  <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" d="M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197M13 7a4 4 0 11-8 0 4 4 0 018 0z" />
                  </svg>
                </button>
              </Tooltip>
              
              <!-- 群组设置按钮 - 仅管理员可见 -->
              <Tooltip v-if="isAdmin" title="群组设置">
                <button 
                  @click="openGroupSettings" 
                  class="p-2 text-gray-500 hover:text-blue-500 dark:text-gray-400 dark:hover:text-blue-400 transition-colors rounded-lg hover:bg-gray-100 dark:hover:bg-gray-800"
                >
                  <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z" />
                    <path stroke-linecap="round" stroke-linejoin="round" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
                  </svg>
                </button>
              </Tooltip>
              
              <!-- 更多操作按钮 -->
              <Dropdown>
                <button 
                  class="p-2 text-gray-500 hover:text-blue-500 dark:text-gray-400 dark:hover:text-blue-400 transition-colors rounded-lg hover:bg-gray-100 dark:hover:bg-gray-800"
                >
                  <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" d="M12 5v.01M12 12v.01M12 19v.01M12 6a1 1 0 110-2 1 1 0 010 2zm0 7a1 1 0 110-2 1 1 0 010 2zm0 7a1 1 0 110-2 1 1 0 010 2z" />
                  </svg>
                </button>
                <template #overlay>
                  <Menu>
                    <Menu.Item v-if="isAdmin" @click="openInviteModal">
                      <div class="flex items-center gap-2">
                        <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                          <path stroke-linecap="round" stroke-linejoin="round" d="M18 9v3m0 0v3m0-3h3m-3 0h-3m-2-5a4 4 0 11-8 0 4 4 0 018 0zM3 20a6 6 0 0112 0v1H3v-1z" />
                        </svg>
                        <span>邀请成员</span>
                      </div>
                    </Menu.Item>
                    <Menu.Item @click="leaveGroup">
                      <div class="flex items-center gap-2 text-orange-500">
                        <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                          <path stroke-linecap="round" stroke-linejoin="round" d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1" />
                        </svg>
                        <span>离开群组</span>
                      </div>
                    </Menu.Item>
                    <Menu.Item v-if="isCreator" @click="deleteGroup">
                      <div class="flex items-center gap-2 text-red-500">
                        <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                          <path stroke-linecap="round" stroke-linejoin="round" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                        </svg>
                        <span>删除群组</span>
                      </div>
                    </Menu.Item>
                  </Menu>
                </template>
              </Dropdown>
            </div>
          </div>
          
          <div class="flex-1 flex overflow-hidden">
            <!-- 聊天消息区域 -->
            <div class="flex-1 overflow-hidden" :class="{ 'hidden md:block': showMembers }">
              <div v-if="activeGroup" class="h-full flex flex-col">
                <!-- 成员较少时显示邀请提示 -->
                <div 
                  v-if="(isAdmin || isCreator) && activeGroup.MemberCount < 3"
                  class="p-3 bg-blue-50 dark:bg-blue-900/20 border-b theme-border flex items-center justify-between"
                >
                  <div class="flex items-center gap-2 text-blue-700 dark:text-blue-300">
                    <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                      <path stroke-linecap="round" stroke-linejoin="round" d="M11 14a2 2 0 100-4 2 2 0 000 4z" />
                      <path stroke-linecap="round" stroke-linejoin="round" d="M21 12c0 1.66-4 3-9 3s-9-1.34-9-3" />
                      <path stroke-linecap="round" stroke-linejoin="round" d="M3 12c0-1.66 4-3 9-3s9 1.34 9 3" />
                    </svg>
                    <span>群组成员较少，邀请更多朋友加入以增加互动</span>
                  </div>
                  <button 
                    @click="openInviteModal"
                    class="px-3 py-1.5 text-xs bg-blue-500 hover:bg-blue-600 text-white rounded-lg transition-colors"
                  >
                    邀请成员
                  </button>
                </div>
                
                <GroupChatMessages :group-id="activeGroup.ID" class="flex-1" />
              </div>
              <div v-else class="h-full flex flex-col items-center justify-center text-gray-500 dark:text-gray-400 bg-gray-50 dark:bg-zinc-800/50">
                <div class="text-6xl mb-4 text-gray-300 dark:text-gray-600">
                  <svg xmlns="http://www.w3.org/2000/svg" class="h-24 w-24" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" d="M8.625 9.75a.375.375 0 11-.75 0 .375.375 0 01.75 0zm0 0H8.25m4.125 0a.375.375 0 11-.75 0 .375.375 0 01.75 0zm0 0H12m4.125 0a.375.375 0 11-.75 0 .375.375 0 01.75 0zm0 0h-.375m-13.5 3.01c0 1.6 1.123 2.994 2.707 3.227 1.087.16 2.185.283 3.293.369V21l4.184-4.183a1.14 1.14 0 01.778-.332 48.294 48.294 0 005.83-.498c1.585-.233 2.708-1.626 2.708-3.228V6.741c0-1.602-1.123-2.995-2.707-3.228A48.394 48.394 0 0012 3c-2.392 0-4.744.175-7.043.513C3.373 3.746 2.25 5.14 2.25 6.741v6.018z" />
                  </svg>
                </div>
                <h2 class="text-2xl font-bold mb-2 text-gray-700 dark:text-gray-300">欢迎使用群组聊天</h2>
                <p class="text-lg mb-6">选择一个群组开始聊天</p>
                <div class="flex flex-col sm:flex-row gap-4">
                  <button 
                    @click="openJoinModal"
                    class="px-6 py-3 bg-gradient-to-r from-blue-500 to-cyan-500 text-white rounded-lg hover:opacity-90 transition-all flex items-center gap-2 justify-center"
                  >
                    <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                      <path stroke-linecap="round" stroke-linejoin="round" d="M18 9v3m0 0v3m0-3h3m-3 0h-3m-2-5a4 4 0 11-8 0 4 4 0 018 0zM3 20a6 6 0 0112 0v1H3v-1z" />
                    </svg>
                    加入群组
                  </button>
                  <button 
                    class="px-6 py-3 bg-gray-100 text-gray-700 rounded-lg hover:bg-gray-200 dark:bg-gray-700 dark:text-gray-300 dark:hover:bg-gray-600 transition-all flex items-center gap-2 justify-center"
                  >
                    <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                      <path stroke-linecap="round" stroke-linejoin="round" d="M12 4.5v15m7.5-7.5h-15" />
                    </svg>
                    创建新群组
                  </button>
                </div>
              </div>
            </div>
            
            <!-- 成员列表侧边栏 -->
            <div v-if="showMembers && activeGroup" class="w-full md:w-80 border-l theme-border bg-white dark:bg-zinc-900 animate__animated animate__fadeIn">
              <div class="p-4 border-b theme-border flex items-center justify-between">
                <h3 class="font-bold text-gray-800 dark:text-gray-200">群组成员</h3>
                <button 
                  @click="showMembers = false"
                  class="text-gray-400 hover:text-gray-600 dark:hover:text-gray-200"
                >
                  <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12" />
                  </svg>
                </button>
              </div>
              
              <!-- 成员列表 -->
              <div class="overflow-y-auto h-[calc(100%-57px)]" @scroll="handleScroll">
                <GroupMemberList 
                  :members="members" 
                  :is-admin="isAdmin"
                  :is-creator="isCreator"
                  :creator-id="activeGroup.CreatorID"
                  :current-user-id="userStore.user?.id"
                  @mute-member="toggleMuteMember"
                  @remove-member="removeMember"
                  @toggle-admin="toggleAdminStatus"
                  @update-nickname="updateMemberNickname"
                />
                
                <!-- 回到顶部按钮 -->
                <div 
                  v-show="showBackToTop" 
                  class="fixed bottom-6 right-6 z-10 cursor-pointer bg-blue-500 hover:bg-blue-600 text-white p-3 rounded-full shadow-lg transition-all duration-200 animate__animated animate__fadeIn"
                  @click="scrollToTop('.overflow-y-auto')"
                >
                  <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" d="M5 10l7-7m0 0l7 7m-7-7v18" />
                  </svg>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
    
    <!-- 群组设置模态框 -->
    <Modal
      v-model:visible="settingsVisible"
      title="群组设置"
      @ok="updateGroupSettings"
      okText="保存"
      cancelText="取消"
    >
      <div class="space-y-4">
        <div>
          <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">群组名称</label>
          <Input v-model:value="groupSettingsForm.name" placeholder="输入群组名称" />
        </div>
        <div>
          <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">群组描述</label>
          <Input.TextArea v-model:value="groupSettingsForm.description" placeholder="输入群组描述" :rows="3" />
        </div>
        <div>
          <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">群组头像链接</label>
          <Input v-model:value="groupSettingsForm.avatarUrl" placeholder="输入头像URL" />
        </div>
      </div>
    </Modal>
    
    <!-- 邀请成员模态框 -->
    <Modal
      v-model:visible="inviteModalVisible"
      title="邀请成员"
      @ok="inviteMember"
      okText="邀请"
      cancelText="取消"
      width="650px"
    >
      <div class="space-y-4">
        <Tabs v-model:activeKey="inviteActiveTab">
          <Tabs.TabPane key="userList" tab="选择用户">
            <div class="flex justify-between items-center mb-3">
              <h3 class="text-lg font-bold">选择要邀请的成员</h3>
              <span class="text-blue-500">已选: {{ selectedUserIds.length }}人</span>
            </div>
            
            <Table 
              :columns="userColumns" 
              :data-source="allUsers" 
              :loading="userTableLoading"
              :pagination="{ pageSize: 5 }"
              :rowKey="record => record.id"
              class="mt-4"
            >
              <template #select="{ record }">
                <Checkbox 
                  :checked="selectedUserIds.includes(record.id)" 
                  @change="e => {
                    if (e.target.checked) {
                      selectedUserIds.push(record.id)
                    } else {
                      selectedUserIds = selectedUserIds.filter(id => id !== record.id)
                    }
                  }"
                />
              </template>
              <template #avatar="{ record }">
                <img 
                  v-if="record.avatar" 
                  :src="BACKEND_DOMAIN + record.avatar" 
                  class="w-8 h-8 rounded-full object-cover" 
                  alt="用户头像" 
                />
                <div 
                  v-else 
                  class="w-8 h-8 rounded-full bg-blue-100 flex items-center justify-center text-blue-500"
                >
                  {{ record.username.charAt(0) }}
                </div>
              </template>
            </Table>
          </Tabs.TabPane>
          
          <Tabs.TabPane key="email" tab="通过邮箱邀请">
            <div class="py-3">
              <h3 class="text-lg font-bold mb-3">通过邮箱邀请用户</h3>
              <p class="text-gray-500 dark:text-gray-400 mb-4">输入用户邮箱地址发送邀请</p>
              <Input 
                v-model:value="inviteEmail" 
                placeholder="输入用户邮箱" 
                allowClear
                class="text-lg"
              />
            </div>
          </Tabs.TabPane>
        </Tabs>
      </div>
    </Modal>
    
    <!-- 加入群组模态框 -->
    <Modal
      v-model:visible="joinModalVisible"
      title="加入群组"
      @ok="joinGroup"
      okText="加入"
      cancelText="取消"
    >
      <div class="space-y-4">
        <p class="text-gray-500 dark:text-gray-400">输入群组邀请码加入群组</p>
        <Input v-model:value="joinGroupCode" placeholder="输入群组邀请码" />
      </div>
    </Modal>
  </div>
</template>

<style scoped>
.theme-border {
  border-color: var(--border-color);
}

/* 动画相关样式 */
.animate__animated {
  animation-duration: 0.5s;
}

.animate__delay-1s {
  animation-delay: 0.1s;
}

.animate__delay-2s {
  animation-delay: 0.2s;
}

@keyframes fadeIn {
  from {
    opacity: 0;
    transform: translateY(10px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.animate__fadeIn {
  animation-name: fadeIn;
}
</style> 