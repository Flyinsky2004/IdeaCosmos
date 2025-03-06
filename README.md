# 创剧星球小程序

## 项目结构

```
miniprogram/
├── pages/                # 页面文件夹
│   ├── index/            # 首页
│   ├── profile/          # 个人资料页
│   ├── notification/     # 通知页面
│   ├── message/          # 消息页面
│   ├── chapter-detail/   # 章节详情页
│   ├── project-detail/   # 项目详情页
│   └── logs/             # 日志页面
├── components/           # 公共组件
│   ├── button/           # 按钮组件
│   ├── card/             # 卡片组件
│   ├── avatar/           # 头像组件
│   ├── form-item/        # 表单项组件
│   ├── input/            # 输入框组件
│   └── divider/          # 分割线组件
├── utils/                # 工具函数
│   ├── request.js        # 网络请求封装
│   ├── auth.js           # 用户认证工具
│   ├── format.js         # 格式化工具
│   └── validator.js      # 验证工具
├── api/                  # API接口
│   ├── user.js           # 用户相关接口
│   ├── project.js        # 项目相关接口
│   └── message.js        # 消息相关接口
└── app.js, app.json, app.wxss  # 小程序入口文件
```

## 组件使用示例

### 按钮组件

```html
<!-- 主要按钮 -->
<custom-button type="primary" bindtap="handleTap">主要按钮</custom-button>

<!-- 成功按钮 -->
<custom-button type="success" bindtap="handleTap">成功按钮</custom-button>

<!-- 危险按钮 -->
<custom-button type="danger" bindtap="handleTap">危险按钮</custom-button>

<!-- 禁用按钮 -->
<custom-button disabled bindtap="handleTap">禁用按钮</custom-button>

<!-- 非块级按钮 -->
<custom-button block="{{false}}" bindtap="handleTap">非块级按钮</custom-button>
```

### 卡片组件

```html
<custom-card title="卡片标题" subtitle="卡片副标题">
  <view>卡片内容</view>
</custom-card>
```

### 头像组件

```html
<custom-avatar src="头像地址" editable bindtap="handleAvatarTap" />
```

### 表单项组件

```html
<custom-form-item label="用户名">
  <custom-input value="{{username}}" bindinput="handleInput" data-field="username" placeholder="请输入用户名" />
</custom-form-item>
```

### 输入框组件

```html
<custom-input value="{{value}}" bindinput="handleInput" placeholder="请输入" />
```

### 分割线组件

```html
<custom-divider />
```

## API使用示例

### 用户API

```javascript
const userApi = require('../../api/user');

// 登录
userApi.login({
  username: 'test',
  password: '123456'
}).then(res => {
  console.log(res);
}).catch(err => {
  console.error(err);
});

// 获取用户信息
userApi.getUserInfo().then(res => {
  console.log(res);
}).catch(err => {
  console.error(err);
});
```

### 项目API

```javascript
const projectApi = require('../../api/project');

// 获取项目列表
projectApi.getProjects().then(res => {
  console.log(res);
}).catch(err => {
  console.error(err);
});

// 获取项目详情
projectApi.getProjectDetail('项目ID').then(res => {
  console.log(res);
}).catch(err => {
  console.error(err);
});
```

### 消息API

```javascript
const messageApi = require('../../api/message');

// 获取消息列表
messageApi.getMessages().then(res => {
  console.log(res);
}).catch(err => {
  console.error(err);
});

// 获取未读消息数量
messageApi.getUnreadCount().then(res => {
  console.log(res);
}).catch(err => {
  console.error(err);
});
```

## 工具函数使用示例

### 请求工具

```javascript
const { get, post } = require('../../utils/request');

// GET请求
get('/api/path', { param: 'value' }).then(res => {
  console.log(res);
}).catch(err => {
  console.error(err);
});

// POST请求
post('/api/path', { data: 'value' }).then(res => {
  console.log(res);
}).catch(err => {
  console.error(err);
});
```

### 认证工具

```javascript
const auth = require('../../utils/auth');

// 检查是否登录
const isLoggedIn = auth.isLoggedIn();

// 获取用户信息
const userInfo = auth.getUserInfo();

// 登出
auth.logout();
```

### 格式化工具

```javascript
const format = require('../../utils/format');

// 格式化日期
const formattedDate = format.formatDate(new Date(), 'YYYY-MM-DD');

// 格式化相对时间
const relativeTime = format.formatRelativeTime(new Date() - 3600000);

// 格式化数字
const formattedNumber = format.formatNumber(1000000);

// 格式化文件大小
const fileSize = format.formatFileSize(1024 * 1024);
```

### 验证工具

```javascript
const validator = require('../../utils/validator');

// 验证邮箱
const isValidEmail = validator.isEmail('test@example.com');

// 验证表单
const { isValid, errors } = validator.validateForm({
  username: 'test',
  password: '123456'
}, {
  username: {
    required: true,
    message: '用户名不能为空'
  },
  password: {
    required: true,
    min: 6,
    message: '密码不能为空',
    minMessage: '密码长度不能小于6位'
  }
});
``` 