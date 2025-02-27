<script setup>
import { ref, onMounted, reactive } from 'vue'
import { get, post, postJSON } from '@/util/request'
import { message,Modal } from 'ant-design-vue'
import SpinLoaderLarge from '@/components/spinLoaderLarge.vue'
import { BACKEND_DOMAIN } from "@/util/VARRIBLES"
import { useUserStore } from '@/stores/user'

const loading = ref(true)
const users = ref([])
const searchKeyword = ref('')
const selectedRole = ref('全部')
const currentPage = ref(1)
const pageSize = ref(10)
const total = ref(0)
const userStore = useUserStore()
// 查看/编辑用户的模态框
const editUserModal = reactive({
  visible: false,
  loading: false,
  type: 'view', // 'view' 或 'edit'
  currentUser: {}
})

// 获取用户列表
const fetchUsers = async () => {
  loading.value = true
  
  get('/api/admin/users', {
    page: currentPage.value,
    pageSize: pageSize.value,
    keyword: searchKeyword.value,
    role: selectedRole.value !== '全部' ? selectedRole.value : ''
  }, 
  (msg, data) => {
    users.value = data.users
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

// 删除用户 - 修改为JSON提交
const deleteUser = (userId) => {
  Modal.confirm({
    title: '确认删除',
    content: '您确定要删除此用户吗？此操作不可撤销。',
    okText: '确认',
    cancelText: '取消',
    onOk: () => {
      postJSON(`/api/admin/users/${userId}/delete`, {}, 
        (msg, data) => {
          message.success('用户删除成功')
          fetchUsers()
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

// 更新用户状态（封禁/解封） - 修改为JSON提交
const updateUserStatus = (userId, status) => {
  postJSON(`/api/admin/users/${userId}/status`, {
    status: status
  }, 
  (msg, data) => {
    message.success('用户状态更新成功')
    fetchUsers()
  },
  (msg) => {
    message.warning(msg)
  },
  (msg) => {
    message.error(msg)
  })
}

// 更新用户角色 - 修改为JSON提交
const updateUserRole = (userId, role) => {
  postJSON(`/api/admin/users/${userId}/role`, {
    role: role
  }, 
  (msg, data) => {
    message.success('用户角色更新成功')
    fetchUsers()
  },
  (msg) => {
    message.warning(msg)
  },
  (msg) => {
    message.error(msg)
  })
}

// 保存用户编辑 - 修改为JSON提交
const saveUserEdit = () => {
  editUserModal.loading = true
  
  postJSON(`/api/admin/users/${editUserModal.currentUser.ID}`, {
    username: editUserModal.currentUser.username,
    email: editUserModal.currentUser.email,
    permission: editUserModal.currentUser.permission,
    group: editUserModal.currentUser.group,
    tokens: editUserModal.currentUser.tokens
  }, 
  (msg, data) => {
    message.success('用户信息更新成功')
    editUserModal.loading = false
    editUserModal.visible = false
    fetchUsers() // 刷新列表
  },
  (msg) => {
    message.warning(msg)
    editUserModal.loading = false
  },
  (msg) => {
    message.error(msg)
    editUserModal.loading = false
  })
}

// 打开用户详情/编辑模态框
const openUserModal = (user, type = 'view') => {
  // 深拷贝用户信息，避免直接修改引用
  editUserModal.currentUser = JSON.parse(JSON.stringify(user))
  editUserModal.type = type
  editUserModal.visible = true
}

// 获取用户状态文本
const getUserStatusText = (user) => {
  return user.group === 0 ? '正常' : '已封禁'
}

// 获取用户角色文本
const getUserRoleText = (user) => {
  return user.permission >= 1 ? '管理员' : '普通用户'
}

// 页码变化
const onPageChange = (page) => {
  currentPage.value = page
  fetchUsers()
}

// 每页条数变化
const onPageSizeChange = (current, size) => {
  pageSize.value = size
  fetchUsers()
}

// 搜索
const onSearch = () => {
  currentPage.value = 1
  fetchUsers()
}

// 重置搜索
const resetSearch = () => {
  searchKeyword.value = ''
  selectedRole.value = '全部'
  currentPage.value = 1
  fetchUsers()
}

// 页面加载时获取用户列表
onMounted(() => {
  fetchUsers()
})
</script>

<template>
  <div class="p-4">
    <h1 class="text-2xl font-bold mb-6 dark:text-gray-100">用户管理</h1>

    <!-- 搜索和筛选 -->
    <div class="mb-6 bg-white dark:bg-zinc-800 rounded-lg p-4 shadow-sm">
      <div class="flex flex-wrap gap-4 items-end">
        <div>
          <p class="text-sm text-gray-500 dark:text-gray-400 mb-1">用户名/邮箱</p>
          <a-input v-model:value="searchKeyword" placeholder="搜索用户名或邮箱" allowClear @pressEnter="onSearch" />
        </div>
        <div>
          <p class="text-sm text-gray-500 dark:text-gray-400 mb-1">角色</p>
          <a-select v-model:value="selectedRole" style="width: 120px">
            <a-select-option value="全部">全部</a-select-option>
            <a-select-option value="admin">管理员</a-select-option>
            <a-select-option value="user">普通用户</a-select-option>
          </a-select>
        </div>
        <a-button type="primary" @click="onSearch">搜索</a-button>
        <a-button @click="resetSearch">重置</a-button>
      </div>
    </div>

    <!-- 用户列表 -->
    <div class="bg-white dark:bg-zinc-800 rounded-lg p-4 shadow-sm">
      <div v-if="loading" class="flex justify-center items-center py-12">
        <SpinLoaderLarge />
      </div>
      <div v-else class="animate__animated animate__fadeIn">
        <a-table 
          :dataSource="users" 
          :pagination="false"
          rowKey="ID"
        >
          <!-- 用户信息列 -->
          <a-table-column key="user" title="用户信息">
            <template #default="{ record }">
              <div class="flex items-center">
                <a-avatar :src="record.avatar ? BACKEND_DOMAIN + record.avatar : ''" :size="40">
                  {{ record.username ? record.username.charAt(0).toUpperCase() : 'U' }}
                </a-avatar>
                <div class="ml-3">
                  <p class="font-medium text-gray-900 dark:text-gray-100">{{ record.username }}</p>
                  <p class="text-sm text-gray-500 dark:text-gray-400">{{ record.email }}</p>
                </div>
              </div>
            </template>
          </a-table-column>

          <!-- 角色列 -->
          <a-table-column key="role" title="角色" width="120px">
            <template #default="{ record }">
              <a-tag :color="record.permission >= 1 ? 'red' : 'blue'">
                {{ getUserRoleText(record) }}
              </a-tag>
            </template>
          </a-table-column>

          <!-- 状态列 -->
          <a-table-column key="status" title="状态" width="120px">
            <template #default="{ record }">
              <a-tag :color="record.group === 0 ? 'green' : 'red'">
                {{ getUserStatusText(record) }}
              </a-tag>
            </template>
          </a-table-column>

          <!-- 积分列 -->
          <a-table-column key="tokens" title="积分" width="100px">
            <template #default="{ record }">
              {{ record.tokens }}
            </template>
          </a-table-column>

          <!-- 注册时间列 -->
          <a-table-column key="created_at" title="注册时间" width="160px">
            <template #default="{ record }">
              {{ new Date(record.CreatedAt).toLocaleString('zh-CN') }}
            </template>
          </a-table-column>

          <!-- 操作列 -->
          <a-table-column key="action" title="操作" width="280px">
            <template #default="{ record }">
              <div class="flex gap-2">
                <a-button size="small" @click="openUserModal(record, 'view')">查看</a-button>
                <a-button size="small" type="primary" @click="openUserModal(record, 'edit')">编辑</a-button>
                
                <!-- 封禁/解封按钮 -->
                <a-button 
                  size="small" 
                  :type="record.group === 0 ? 'danger' : 'default'" 
                  @click="updateUserStatus(record.ID, record.group === 0 ? 'banned' : 'active')"
                  :disabled="record.permission >= 1 && record.ID !== userStore?.id"
                >
                  {{ record.group === 0 ? '封禁' : '解封' }}
                </a-button>
                
                <!-- 角色切换按钮 -->
                <a-button 
                  size="small" 
                  :type="record.permission >= 1 ? 'default' : 'primary'" 
                  @click="updateUserRole(record.ID, record.permission >= 1 ? 'user' : 'admin')"
                >
                  {{ record.permission >= 1 ? '降为用户' : '设为管理员' }}
                </a-button>
                
                <!-- 删除按钮 -->
                <a-button 
                  size="small" 
                  danger 
                  @click="deleteUser(record.ID)"
                  :disabled="record.permission >= 1 && record.ID !== userStore?.id"
                >
                  删除
                </a-button>
              </div>
            </template>
          </a-table-column>
        </a-table>

        <!-- 分页 -->
        <div class="flex justify-end mt-4">
          <a-pagination
            v-model:current="currentPage"
            v-model:pageSize="pageSize"
            :total="total"
            show-size-changer
            show-total
            @change="onPageChange"
            @showSizeChange="onPageSizeChange"
          />
        </div>
      </div>
    </div>

    <!-- 用户详情/编辑模态框 -->
    <a-modal
      :title="editUserModal.type === 'view' ? '用户详情' : '编辑用户'"
      v-model:visible="editUserModal.visible"
      :footer="editUserModal.type === 'view' ? null : undefined"
      @ok="saveUserEdit"
      :confirmLoading="editUserModal.loading"
      width="700px"
    >
      <div v-if="Object.keys(editUserModal.currentUser).length > 0">
        <!-- 用户头像和基本信息 -->
        <div class="flex items-center mb-6">
          <a-avatar :src="editUserModal.currentUser.avatar ? BACKEND_DOMAIN + editUserModal.currentUser.avatar : ''" :size="80">
            {{ editUserModal.currentUser.username ? editUserModal.currentUser.username.charAt(0).toUpperCase() : 'U' }}
          </a-avatar>
          <div class="ml-6">
            <h2 class="text-xl font-bold dark:text-gray-200">{{ editUserModal.currentUser.username }}</h2>
            <p class="text-gray-500 dark:text-gray-400">
              <a-tag v-if="editUserModal.currentUser.permission >= 1" color="red">管理员</a-tag>
              <a-tag :color="editUserModal.currentUser.group === 0 ? 'green' : 'red'">
                {{ editUserModal.currentUser.group === 0 ? '正常' : '已封禁' }}
              </a-tag>
            </p>
          </div>
        </div>

        <!-- 详细信息表单 -->
        <a-divider />
        
        <div v-if="editUserModal.type === 'view'">
          <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div>
              <p class="text-gray-500 dark:text-gray-400 text-sm mb-1">用户ID</p>
              <p>{{ editUserModal.currentUser.ID }}</p>
            </div>
            <div>
              <p class="text-gray-500 dark:text-gray-400 text-sm mb-1">电子邮件</p>
              <p>{{ editUserModal.currentUser.email }}</p>
            </div>
            <div>
              <p class="text-gray-500 dark:text-gray-400 text-sm mb-1">注册时间</p>
              <p>{{ new Date(editUserModal.currentUser.CreatedAt).toLocaleString('zh-CN') }}</p>
            </div>
            <div>
              <p class="text-gray-500 dark:text-gray-400 text-sm mb-1">积分</p>
              <p>{{ editUserModal.currentUser.tokens }}</p>
            </div>
            <div>
              <p class="text-gray-500 dark:text-gray-400 text-sm mb-1">权限级别</p>
              <p>{{ editUserModal.currentUser.permission }}</p>
            </div>
            <div>
              <p class="text-gray-500 dark:text-gray-400 text-sm mb-1">状态值</p>
              <p>{{ editUserModal.currentUser.group }}</p>
            </div>
          </div>
        </div>

        <div v-else>
          <a-form layout="vertical">
            <a-form-item label="用户名">
              <a-input v-model:value="editUserModal.currentUser.username" />
            </a-form-item>
            <a-form-item label="电子邮件">
              <a-input v-model:value="editUserModal.currentUser.email" />
            </a-form-item>
            <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
              <a-form-item label="权限级别">
                <a-input-number v-model:value="editUserModal.currentUser.permission" :min="0" :max="9" />
              </a-form-item>
              <a-form-item label="状态值">
                <a-select v-model:value="editUserModal.currentUser.group">
                  <a-select-option :value="0">正常</a-select-option>
                  <a-select-option :value="1">封禁</a-select-option>
                </a-select>
              </a-form-item>
              <a-form-item label="积分">
                <a-input-number v-model:value="editUserModal.currentUser.tokens" :min="0" />
              </a-form-item>
            </div>
          </a-form>
        </div>
      </div>
    </a-modal>
  </div>
</template>

<style scoped>
.animate__animated {
  animation-duration: 0.5s;
}
</style> 