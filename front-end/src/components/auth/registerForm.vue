<script setup>
import { ref, reactive, onUnmounted } from 'vue'
import router from "@/router/index.js"
import { post, postJSON } from '@/util/request'
import { message } from 'ant-design-vue'

const formState = reactive({
  username: '',
  password: '',
  confirmPassword: '',
  email: '',
  verifyCode: ''
})

const currentStep = ref(0)
const loading = ref(false)
const sendingCode = ref(false)
const countdown = ref(0)
let timer = null

// 步骤配置
const steps = [
  { title: '设置账号', description: '创建您的用户名和密码' },
  { title: '验证邮箱', description: '验证您的邮箱地址' },
  { title: '完成注册', description: '确认信息并完成注册' }
]

// 验证当前步骤并前进
const nextStep = () => {
  if (currentStep.value === 0) {
    if (!formState.username || !formState.password || !formState.confirmPassword) {
      message.error('请填写完整信息')
      return
    }
    if (formState.password !== formState.confirmPassword) {
      message.error('两次输入的密码不一致')
      return
    }
    if (formState.password.length < 6) {
      message.error('密码长度不能少于6位')
      return
    }
    currentStep.value++
  } else if (currentStep.value === 1) {
    if (!formState.email) {
      message.error('请输入邮箱地址')
      return
    }
    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
    if (!emailRegex.test(formState.email)) {
      message.error('请输入正确的邮箱格式')
      return
    }
    currentStep.value++
  }
}

// 返回上一步
const prevStep = () => {
  if (currentStep.value > 0) {
    currentStep.value--
  }
}

// 发送验证码
const sendVerifyCode = async () => {
  try {
    sendingCode.value = true
    postJSON('/api/auth/sendCode', {
      email: formState.email
    }, (msg, data) => {
      message.success('验证码已发送至邮箱')
      countdown.value = 60
      timer = setInterval(() => {
        countdown.value--
        if (countdown.value <= 0) {
          clearInterval(timer)
        }
      }, 1000)
    },(msg,data) => {
      message.warning(msg)
    },(msg,data) => {
      message.warning(msg)
    })
  } catch (error) {
    message.error('发送验证码失败')
  } finally {
    sendingCode.value = false
  }
}

// 注册
const handleRegister = async () => {
  if (!formState.verifyCode) {
    message.error('请输入验证码')
    return
  }

  try {
    loading.value = true
    postJSON('/api/auth/register', {
      username: formState.username,
      password: formState.password,
      email: formState.email,
      code: formState.verifyCode
    }, (msg, data) => {
      message.success('注册成功')
      router.push('/auth/login')
    },(msg,data) => {
      message.warning(msg)
    },(msg,data) => {
      message.warning(msg)
    })
  } finally {
    loading.value = false
  }
}

onUnmounted(() => {
  if (timer) clearInterval(timer)
})
</script>

<template>
  <div class="font-sans text-theme-switch mt-8">
    <div class="text-center text-4xl font-bold">注册</div>
    <div class="text-center text-xl mt-4 text-gray-600 dark:text-gray-400">欢迎来到创剧星球</div>

    <!-- 步骤指示器 -->
    <div class="flex justify-between mt-8 mb-12">
      <div v-for="(step, index) in steps" :key="index" 
           class="flex-1 relative">
        <div class="flex flex-col items-center">
          <div :class="[
            'w-8 h-8 rounded-full flex items-center justify-center text-sm font-medium transition-all',
            index <= currentStep ? 'bg-blue-500 text-white' : 'bg-gray-200 dark:bg-gray-700 text-gray-500'
          ]">
            {{ index + 1 }}
          </div>
          <div class="mt-2 text-sm font-medium" :class="index <= currentStep ? 'text-blue-500' : 'text-gray-500'">
            {{ step.title }}
          </div>
          <div class="text-xs text-gray-400 mt-1 text-center">{{ step.description }}</div>
        </div>
        <div v-if="index < steps.length - 1" :class="[
          'absolute top-4 w-full h-[2px] left-1/2',
          index < currentStep ? 'bg-blue-500' : 'bg-gray-200 dark:bg-gray-700'
        ]"></div>
      </div>
    </div>

    <!-- 步骤内容 -->
    <div class="mt-8">
      <!-- 步骤1：账号信息 -->
      <div v-if="currentStep === 0" class="space-y-4">
        <div>
          <label class="block text-sm font-medium mb-2">用户名</label>
          <input 
            v-model="formState.username"
            type="text"
            class="input1 w-full"
            placeholder="请输入用户名"
          />
        </div>
        <div>
          <label class="block text-sm font-medium mb-2">密码</label>
          <input 
            v-model="formState.password"
            type="password"
            class="input1 w-full"
            placeholder="请输入密码"
          />
        </div>
        <div>
          <label class="block text-sm font-medium mb-2">确认密码</label>
          <input 
            v-model="formState.confirmPassword"
            type="password"
            class="input1 w-full"
            placeholder="请再次输入密码"
          />
        </div>
      </div>

      <!-- 步骤2：邮箱验证 -->
      <div v-if="currentStep === 1" class="space-y-4">
        <div>
          <label class="block text-sm font-medium mb-2">邮箱地址</label>
          <input 
            v-model="formState.email"
            type="email"
            class="input1 w-full"
            placeholder="请输入邮箱"
          />
        </div>
      </div>

      <!-- 步骤3：完成注册 -->
      <div v-if="currentStep === 2" class="space-y-4">
        <div class="bg-gray-50 dark:bg-zinc-800 p-4 rounded-lg">
          <h3 class="font-medium mb-2">确认信息</h3>
          <div class="space-y-2 text-sm">
            <div class="flex justify-between">
              <span class="text-gray-500">用户名</span>
              <span>{{ formState.username }}</span>
            </div>
            <div class="flex justify-between">
              <span class="text-gray-500">邮箱</span>
              <span>{{ formState.email }}</span>
            </div>
          </div>
        </div>
        <div>
          <label class="block text-sm font-medium mb-2">验证码</label>
          <div class="flex gap-4">
            <input 
              v-model="formState.verifyCode"
              type="text"
              class="input1 flex-1"
              placeholder="请输入验证码"
            />
            <button 
              @click="sendVerifyCode"
              :disabled="countdown > 0 || sendingCode"
              class="btn1 w-32"
            >
              {{ countdown > 0 ? `${countdown}s后重试` : '获取验证码' }}
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- 操作按钮 -->
    <div class="flex justify-between mt-8">
      <button 
        v-if="currentStep > 0"
        @click="prevStep"
        class="btn1 px-8"
      >
        上一步
      </button>
      <div v-else class="w-20"></div>

      <button 
        v-if="currentStep < 2"
        @click="nextStep"
        class="btn1 px-8 bg-blue-500 text-white hover:bg-blue-600"
      >
        下一步
      </button>
      <button 
        v-else
        @click="handleRegister"
        :disabled="loading"
        class="btn1 px-8 bg-gradient-to-r from-blue-500 to-violet-500 text-white hover:opacity-90"
      >
        {{ loading ? '注册中...' : '完成注册' }}
      </button>
    </div>

    <!-- 登录链接 -->
    <div class="mt-4 text-center text-sm">
      已有账号？
      <a 
        @click="router.push('/auth/login')"
        class="text-blue-500 hover:text-blue-600 cursor-pointer hover:underline"
      >
        马上登录
      </a>
    </div>
  </div>
</template>

<style scoped>

</style>