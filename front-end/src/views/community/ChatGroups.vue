<!--
 * @Author: Flyinsky w2084151024@gmail.com
 * @Description: 群组聊天页面
-->
<script setup>
import { message, Modal, Input, Dropdown, Menu, Tooltip, Checkbox, Table, Tabs } from 'ant-design-vue'
import GroupChatList from '@/components/chat/GroupChatList.vue'
import GroupChatMessages from '@/components/chat/GroupChatMessages.vue'
import GroupMemberList from '@/components/chat/GroupMemberList.vue'
import { useChatGroups } from './chatGroupsLogic'

// 导入并使用封装的群组聊天逻辑
const {
  // 状态
  userStore,
  activeGroup,
  members,
  showMembers,
  isCreator,
  isAdmin,
  settingsVisible,
  inviteModalVisible,
  inviteEmail,
  inviteActiveTab,
  joinModalVisible,
  joinGroupCode,
  groupSettingsForm,
  allUsers,
  selectedUserIds,
  userTableLoading,
  showBackToTop,
  userColumns,
  BACKEND_DOMAIN,

  // 方法
  scrollToTop,
  handleScroll,
  fetchGroupMembers,
  handleSelectGroup,
  handleCreateGroup,
  openGroupSettings,
  updateGroupSettings,
  openInviteModal,
  fetchAllUsers,
  inviteMember,
  openJoinModal,
  joinGroup,
  leaveGroup,
  deleteGroup,
  toggleMuteMember,
  removeMember,
  toggleAdminStatus,
  updateMemberNickname
} = useChatGroups()
</script>

<template>
  <div class="flex flex-col h-[calc(100vh-5rem)] animate__animated animate__fadeIn bg-[#f8fafc] dark:bg-[#030616]">
    <!-- 聊天主区域 -->
    <div class="flex-1 mx-4 mb-4 rounded-2xl shadow-lg overflow-hidden animate__animated animate__fadeIn animate__delay-2s">
      <div class="h-full flex">
        <!-- 左侧群组列表 -->
        <div class="w-80 flex-shrink-0 border-r border-gray-200 dark:border-gray-700 bg-white dark:bg-zinc-900">
          <GroupChatList 
            :active-group-id="activeGroup?.ID"
            @select-group="handleSelectGroup"
            @create-group="handleCreateGroup"
          />
        </div>
        
        <!-- 右侧聊天区域 -->
        <div class="flex-1 flex flex-col bg-white dark:bg-zinc-900">
          <!-- 群组信息栏 -->
          <div v-if="activeGroup" class="flex items-center justify-between px-6 py-4 border-b border-gray-200 dark:border-gray-700">
            <div class="flex items-center gap-3">
              <div class="w-12 h-12 rounded-full overflow-hidden bg-gradient-to-br from-blue-100 to-indigo-100 dark:from-blue-900/30 dark:to-indigo-900/30 flex-shrink-0 border border-gray-200 dark:border-gray-700 shadow-sm">
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
                <h3 class="font-bold text-lg text-gray-800 dark:text-gray-200">{{ activeGroup.Name }}</h3>
                <p class="text-xs text-gray-500 dark:text-gray-400 flex items-center gap-1.5">
                  <span class="flex items-center gap-1">
                    <svg xmlns="http://www.w3.org/2000/svg" class="h-3.5 w-3.5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                      <path stroke-linecap="round" stroke-linejoin="round" d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z" />
                    </svg>
                    {{ activeGroup.MemberCount }}人
                  </span>
                  <span class="w-1 h-1 rounded-full bg-gray-300 dark:bg-gray-600"></span>
                  <span class="text-gray-400 dark:text-gray-500">{{ activeGroup.Description || '暂无群组描述' }}</span>
                </p>
              </div>
            </div>
            <div class="flex items-center gap-2">
              <!-- 邀请成员按钮 - 仅管理员可见 -->
              <button 
                @click="openInviteModal" 
                class="px-3 py-1.5 text-sm text-white bg-gradient-to-r from-blue-500 to-cyan-500 hover:opacity-90 hover:shadow-md hover:shadow-blue-500/20 rounded-lg flex items-center gap-1.5 transition-all duration-300 transform hover:scale-105"
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
                  class="p-2 text-gray-500 hover:text-blue-500 dark:text-gray-400 dark:hover:text-blue-400 transition-colors rounded-lg hover:bg-gray-100 dark:hover:bg-gray-800 transform hover:scale-105"
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
                  class="p-2 text-gray-500 hover:text-blue-500 dark:text-gray-400 dark:hover:text-blue-400 transition-colors rounded-lg hover:bg-gray-100 dark:hover:bg-gray-800 transform hover:scale-105"
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
                  class="p-2 text-gray-500 hover:text-blue-500 dark:text-gray-400 dark:hover:text-blue-400 transition-colors rounded-lg hover:bg-gray-100 dark:hover:bg-gray-800 transform hover:scale-105"
                >
                  <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" d="M12 5v.01M12 12v.01M12 19v.01M12 6a1 1 0 110-2 1 1 0 010 2zm0 7a1 1 0 110-2 1 1 0 010 2zm0 7a1 1 0 110-2 1 1 0 010 2z" />
                  </svg>
                </button>
                <template #overlay>
                  <Menu class="bg-white dark:bg-zinc-800 border border-gray-200 dark:border-gray-700 shadow-xl rounded-xl overflow-hidden">
                    <Menu.Item v-if="isAdmin" @click="openInviteModal">
                      <div class="flex items-center gap-2 px-4 py-2 hover:bg-gray-50 dark:hover:bg-zinc-700">
                        <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 text-blue-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                          <path stroke-linecap="round" stroke-linejoin="round" d="M18 9v3m0 0v3m0-3h3m-3 0h-3m-2-5a4 4 0 11-8 0 4 4 0 018 0zM3 20a6 6 0 0112 0v1H3v-1z" />
                        </svg>
                        <span>邀请成员</span>
                      </div>
                    </Menu.Item>
                    <Menu.Item @click="leaveGroup">
                      <div class="flex items-center gap-2 px-4 py-2 hover:bg-gray-50 dark:hover:bg-zinc-700 text-orange-500">
                        <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                          <path stroke-linecap="round" stroke-linejoin="round" d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1" />
                        </svg>
                        <span>离开群组</span>
                      </div>
                    </Menu.Item>
                    <Menu.Item v-if="isCreator" @click="deleteGroup">
                      <div class="flex items-center gap-2 px-4 py-2 hover:bg-gray-50 dark:hover:bg-zinc-700 text-red-500">
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
                  class="p-3 bg-gradient-to-r from-blue-50 to-indigo-50 dark:from-blue-900/20 dark:to-indigo-900/10 border-b border-gray-200 dark:border-gray-700 flex items-center justify-between"
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
                    class="px-3 py-1.5 text-xs bg-gradient-to-r from-blue-500 to-cyan-500 text-white rounded-lg transition-all hover:opacity-90 transform hover:scale-105"
                  >
                    邀请成员
                  </button>
                </div>
                
                <GroupChatMessages :group-id="activeGroup.ID" class="flex-1" />
              </div>
              <div v-else class="h-full flex flex-col items-center justify-center px-6 py-12 bg-gradient-to-b from-gray-50 to-white dark:from-zinc-900 dark:to-zinc-900/90">
                <div class="text-6xl mb-4 text-gray-300 dark:text-gray-600 animate__animated animate__fadeIn animate__delay-3s">
                  <svg xmlns="http://www.w3.org/2000/svg" class="h-24 w-24" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" d="M8.625 9.75a.375.375 0 11-.75 0 .375.375 0 01.75 0zm0 0H8.25m4.125 0a.375.375 0 11-.75 0 .375.375 0 01.75 0zm0 0H12m4.125 0a.375.375 0 11-.75 0 .375.375 0 01.75 0zm0 0h-.375m-13.5 3.01c0 1.6 1.123 2.994 2.707 3.227 1.087.16 2.185.283 3.293.369V21l4.184-4.183a1.14 1.14 0 01.778-.332 48.294 48.294 0 005.83-.498c1.585-.233 2.708-1.626 2.708-3.228V6.741c0-1.602-1.123-2.995-2.707-3.228A48.394 48.394 0 0012 3c-2.392 0-4.744.175-7.043.513C3.373 3.746 2.25 5.14 2.25 6.741v6.018z" />
                  </svg>
                </div>
                <div class="border p-6 w-fit max-w-md rounded-xl bg-gradient-to-r from-blue-50 to-indigo-50 dark:from-blue-900/30 dark:to-indigo-900/30 dark:border-gray-700 animate__animated animate__fadeIn animate__delay-4s">
                  <h2 class="text-2xl font-bold mb-2 text-gray-800 dark:text-gray-200">欢迎使用群组聊天</h2>
                  <p class="text-lg mb-6 text-gray-600 dark:text-gray-400">选择一个群组开始聊天，或创建新的群组</p>
                  <div class="flex flex-col sm:flex-row gap-4">
                    <button 
                      @click="openJoinModal"
                      class="px-6 py-3 bg-gradient-to-r from-blue-500 to-cyan-500 text-white rounded-lg hover:opacity-90 hover:shadow-md hover:shadow-blue-500/20 transition-all flex items-center gap-2 justify-center transform hover:scale-105"
                    >
                      <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                        <path stroke-linecap="round" stroke-linejoin="round" d="M18 9v3m0 0v3m0-3h3m-3 0h-3m-2-5a4 4 0 11-8 0 4 4 0 018 0zM3 20a6 6 0 0112 0v1H3v-1z" />
                      </svg>
                      加入群组
                    </button>
                    <button 
                      class="px-6 py-3 bg-gray-100 text-gray-700 rounded-lg hover:bg-gray-200 dark:bg-gray-700 dark:text-gray-300 dark:hover:bg-gray-600 transition-all flex items-center gap-2 justify-center transform hover:scale-105"
                    >
                      <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                        <path stroke-linecap="round" stroke-linejoin="round" d="M12 4.5v15m7.5-7.5h-15" />
                      </svg>
                      创建新群组
                    </button>
                  </div>
                </div>
              </div>
            </div>
            
            <!-- 成员列表侧边栏 -->
            <div v-if="showMembers && activeGroup" class="w-full md:w-80 border-l border-gray-200 dark:border-gray-700 bg-white dark:bg-zinc-900 animate__animated animate__fadeIn">
              <div class="p-4 border-b border-gray-200 dark:border-gray-700 flex items-center justify-between">
                <h3 class="font-bold text-gray-800 dark:text-gray-200">群组成员</h3>
                <button 
                  @click="showMembers = false"
                  class="text-gray-400 hover:text-gray-600 dark:hover:text-gray-200 p-1.5 rounded-full hover:bg-gray-100 dark:hover:bg-gray-800 transition-all"
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
                  class="fixed bottom-6 right-6 z-10 cursor-pointer bg-gradient-to-r from-blue-500 to-cyan-500 hover:opacity-90 text-white p-3 rounded-full shadow-lg hover:shadow-blue-500/20 transition-all duration-300 animate__animated animate__fadeIn transform hover:scale-110"
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
      class="custom-modal"
      :bodyStyle="{
        backgroundColor: 'var(--bg-color)',
        borderRadius: '0.75rem',
      }"
    >
      <div class="space-y-4 p-2">
        <div>
          <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">群组名称</label>
          <Input v-model:value="groupSettingsForm.name" placeholder="输入群组名称" class="hover:border-blue-500 focus:border-blue-500" />
        </div>
        <div>
          <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">群组描述</label>
          <Input.TextArea v-model:value="groupSettingsForm.description" placeholder="输入群组描述" :rows="3" class="hover:border-blue-500 focus:border-blue-500" />
        </div>
        <div>
          <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">群组头像链接</label>
          <Input v-model:value="groupSettingsForm.avatarUrl" placeholder="输入头像URL" class="hover:border-blue-500 focus:border-blue-500" />
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
      class="custom-modal"
      :bodyStyle="{
        backgroundColor: 'var(--bg-color)',
        borderRadius: '0.75rem',
      }"
    >
      <div class="space-y-4 p-2">
        <Tabs v-model:activeKey="inviteActiveTab" class="custom-tabs">
          <Tabs.TabPane key="userList" tab="选择用户">
            <div class="flex justify-between items-center mb-4">
              <h3 class="text-lg font-bold text-gray-800 dark:text-gray-200">选择要邀请的成员</h3>
              <span class="text-blue-500 dark:text-blue-400 px-3 py-1 bg-blue-50 dark:bg-blue-900/30 rounded-full text-sm">已选: {{ selectedUserIds.length }}人</span>
            </div>
            
            <Table 
              :columns="userColumns" 
              :data-source="allUsers" 
              :loading="userTableLoading"
              :pagination="{ pageSize: 5 }"
              :rowKey="record => record.id"
              class="custom-table"
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
                  class="custom-checkbox"
                />
              </template>
              <template #avatar="{ record }">
                <div class="w-10 h-10 rounded-full overflow-hidden bg-gradient-to-br from-blue-100 to-indigo-100 dark:from-blue-900/30 dark:to-indigo-900/30 flex-shrink-0 border theme-border"> 
                  <img 
                    v-if="record.avatar" 
                    :src="BACKEND_DOMAIN + record.avatar" 
                    class="w-full h-full object-cover" 
                    alt="用户头像" 
                  />
                  <div 
                    v-else 
                    class="w-full h-full flex items-center justify-center text-blue-500 dark:text-blue-400 text-sm font-bold"
                  >
                    {{ record.username.charAt(0) }}
                  </div>
                </div>
              </template>
            </Table>
          </Tabs.TabPane>
          
          <Tabs.TabPane key="email" tab="通过邮箱邀请">
            <div class="py-3">
              <h3 class="text-lg font-bold mb-3 text-gray-800 dark:text-gray-200">通过邮箱邀请用户</h3>
              <p class="text-gray-500 dark:text-gray-400 mb-4">输入用户邮箱地址发送邀请</p>
              <Input 
                v-model:value="inviteEmail" 
                placeholder="输入用户邮箱" 
                allowClear
                class="text-lg hover:border-blue-500 focus:border-blue-500"
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
      class="custom-modal"
      :bodyStyle="{
        backgroundColor: 'var(--bg-color)',
        borderRadius: '0.75rem',
      }"
    >
      <div class="space-y-4 p-2">
        <p class="text-gray-500 dark:text-gray-400">输入群组邀请码加入群组</p>
        <Input v-model:value="joinGroupCode" placeholder="输入群组邀请码" class="text-lg hover:border-blue-500 focus:border-blue-500" />
        <div class="bg-gray-50 dark:bg-zinc-800/50 p-3 rounded-lg border theme-border mt-4">
          <p class="text-xs text-gray-500 dark:text-gray-400 flex items-center gap-2">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 text-blue-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
            您可以通过群主或管理员分享的邀请码加入已有群组
          </p>
        </div>
      </div>
    </Modal>
  </div>
</template>

<style scoped>
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

.animate__delay-3s {
  animation-delay: 0.3s;
}

.animate__delay-4s {
  animation-delay: 0.4s;
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

/* 自定义模态框样式 */
:deep(.custom-modal .ant-modal-content) {
  @apply rounded-2xl overflow-hidden border border-gray-200 dark:border-gray-700;
}

:deep(.custom-modal .ant-modal-header) {
  @apply bg-white dark:bg-zinc-900 border-b border-gray-200 dark:border-gray-700;
}

:deep(.custom-modal .ant-modal-footer) {
  @apply bg-white dark:bg-zinc-900 border-t border-gray-200 dark:border-gray-700;
}

:deep(.custom-modal .ant-btn-primary) {
  @apply bg-gradient-to-r from-blue-500 to-cyan-500 border-0;
}

:deep(.custom-modal .ant-btn-primary:hover) {
  @apply opacity-90 shadow-md shadow-blue-500/20;
}

/* 自定义表格样式 */
:deep(.custom-table .ant-table) {
  @apply bg-white dark:bg-zinc-900 rounded-xl overflow-hidden;
}

:deep(.custom-table .ant-table-thead > tr > th) {
  @apply bg-blue-50/50 dark:bg-blue-900/30 text-gray-900 dark:text-gray-100;
}

:deep(.custom-tabs .ant-tabs-tab.ant-tabs-tab-active .ant-tabs-tab-btn) {
  @apply text-blue-500;
}

:deep(.custom-tabs .ant-tabs-ink-bar) {
  @apply bg-gradient-to-r from-blue-500 to-cyan-500;
}

/* 自定义复选框样式 */
:deep(.custom-checkbox .ant-checkbox-checked .ant-checkbox-inner) {
  @apply bg-blue-500 border-blue-500;
}
</style> 