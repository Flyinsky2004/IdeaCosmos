<script setup>
import { onMounted, reactive, ref } from "vue";
import { get, post, postJSON } from "@/util/request";
import { message, Modal, Upload } from "ant-design-vue";
import { imagePrefix , BACKEND_DOMAIN } from "@/util/VARRIBLES";
import { LoadingOutlined, PlusOutlined } from "@ant-design/icons-vue";
import { useUserStore } from "@/stores/user";
const userStore = useUserStore()
const uploadUrl = BACKEND_DOMAIN + 'user/uploadImage'
const authToken = localStorage.getItem("authToken")
const userInfo = reactive({
  data: {},
  isEditing: false,
  imageUrl: '',
  loading: false,
  showAvatarModal: false,
});

const editForm = reactive({
  username: '',
});

onMounted(() => {
  fetchUserInfo();
});

const fetchUserInfo = () => {
  get(
    "/api/user/me",
    {},
    (messager, data) => {
      userInfo.data = data;
      userInfo.imageUrl = data.avatar ? data.avatar : '';
      editForm.username = data.username;
      userStore.login(data)
    },
    (messager) => {
      message.warning(messager);
    },
    (messager) => {
      message.error(messager);
    }
  );
};

const handleAvatarChange = (info) => {
  if (info.file.status === 'uploading') {
    userInfo.loading = true;
    return;
  }
  if (info.file.status === 'done') {
    userInfo.loading = false;
    if (info.file.response.code === 200) {
      const path = info.file.response.data.path;
      userInfo.imageUrl = path;
      // 更新用户信息
      updateUserInfo({ avatar: path });
    } else {
      message.error('上传失败');
    }
  }
  if (info.file.status === 'error') {
    userInfo.loading = false;
    message.error('上传失败');
  }
};

const beforeUpload = (file) => {
  const isJpgOrPng = file.type === 'image/jpeg' || file.type === 'image/png';
  if (!isJpgOrPng) {
    message.error('只能上传JPG/PNG格式的图片！');
  }
  const isLt2M = file.size / 1024 / 1024 < 2;
  if (!isLt2M) {
    message.error('图片大小不能超过2MB！');
  }
  return isJpgOrPng && isLt2M;
};

const showChangeAvatar = () => {
  userInfo.showAvatarModal = true;
};

const updateUserInfo = (data) => {
  postJSON(
    "/api/user/updateInfo",
    data,
    (messager) => {
      message.success("更新成功");
      fetchUserInfo();
      userInfo.isEditing = false;
    },
    (messager) => {
      message.warning(messager);
    },
    (messager) => {
      message.error(messager);
    }
  );
};

const handleEdit = () => {
  if (userInfo.isEditing) {
    // 保存修改
    updateUserInfo({
      username: editForm.username,
    });
  } else {
    userInfo.isEditing = true;
  }
};
</script>

<template>
  <div class="flex flex-col gap-2 animate__animated animate__fadeIn">
    <div class="border p-4 workspace-box w-fit">
      <h1 class="text-2xl font-serif">个人信息</h1>
      <span class="text-sm font-bold">在这里您可以查看和编辑您的个人信息</span>
    </div>

    <div class="border p-4 workspace-box">
      <div class="flex gap-8">
        <!-- 头像部分 -->
        <div class="flex flex-col items-center gap-2">
          <div 
            class="relative group"
            @click="showChangeAvatar"
          >
            <img
              :src="BACKEND_DOMAIN + userInfo.imageUrl || '/default-avatar.png'"
              alt="avatar"
              class="w-32 h-32 rounded-full object-cover border-2 theme-border cursor-pointer"
            />
            <div class="absolute inset-0 bg-black bg-opacity-50 rounded-full opacity-0 group-hover:opacity-100 transition-opacity duration-200 flex items-center justify-center cursor-pointer">
              <span class="text-white text-sm">更换头像</span>
            </div>
          </div>
        </div>

        <!-- 信息部分 -->
        <div class="flex flex-col gap-4 flex-grow">
          <div class="grid grid-cols-[100px,1fr] gap-2 items-center">
            <span class="text-gray-600 dark:text-gray-400">用户名：</span>
            <input
              v-if="userInfo.isEditing"
              v-model="editForm.username"
              class="input1"
              placeholder="请输入用户名"
            />
            <span v-else class="dark:text-gray-200">{{ userInfo.data.username }}</span>
          </div>


          <div class="grid grid-cols-[100px,1fr] gap-2 items-center">
            <span class="text-gray-600 dark:text-gray-400">邮箱：</span>
            <span class="dark:text-gray-200">{{ userInfo.data.email }}</span>
          </div>

          <div class="flex justify-end mt-4">
            <button
              :class="userInfo.isEditing ? 'basic-success-button' : 'basic-prinary-button'"
              @click="handleEdit"
            >
              {{ userInfo.isEditing ? '保存' : '编辑' }}
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>

  <!-- 头像上传模态框 -->
  <Modal
    v-model:open="userInfo.showAvatarModal"
    title="更换头像"
    okText="确定"
    cancelText="取消"
    @cancel="userInfo.showAvatarModal = false"
    @ok="userInfo.showAvatarModal = false"
  >
    <div class="flex flex-col items-center gap-4">
      <Upload
        name="image"
        :show-upload-list="false"
        :action= uploadUrl
        :headers="{
          Authorization: authToken,
        }"
        :before-upload="beforeUpload"
        @change="handleAvatarChange"

      >
        <div class="w-40 h-40 border-2 theme-border rounded-full flex items-center justify-center cursor-pointer hover:border-blue-500 transition-colors">
          <img 
            v-if="userInfo.imageUrl" 
            :src="BACKEND_DOMAIN + userInfo.imageUrl" 
            alt="avatar" 
            class="w-full h-full rounded-full object-cover" 
          />
          <div v-else class="text-center">
            <LoadingOutlined v-if="userInfo.loading" />
            <PlusOutlined v-else />
            <div class="mt-2">点击上传</div>
          </div>
        </div>
      </Upload>
      <div class="text-gray-500 dark:text-gray-400 text-sm">
        支持 JPG、PNG 格式，文件小于 2MB
      </div>
    </div>
  </Modal>
</template>

<style scoped>

</style>