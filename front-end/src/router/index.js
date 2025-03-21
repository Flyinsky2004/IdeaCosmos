import {createRouter, createWebHistory} from 'vue-router'
import {useUserStore} from "@/stores/user.js";
import {get} from "@/util/request.js";
import {message} from "ant-design-vue";

const router = createRouter({
    history: createWebHistory(import.meta.env.BASE_URL),
    routes: [
        {
            path: '/',
            name: 'home',
            component: () => import('@/views/Index.vue')
        }, {
            path: '/auth',
            component: () => import('@/views/public/auth.vue'),
            children: [
                {
                    name: 'Login',
                    path: 'login',
                    component: () => import('@/views/public/auth.vue'),
                }, {
                    name: 'Register',
                    path: 'register',
                    component: () => import('@/views/public/auth.vue'),
                }
            ]
        }, {
            path: '/workspace',
            name: 'workspace',
            component: () => import('@/views/workspace/framework.vue'),
            children: [
                {
                    name: 'dataAnlysiss',
                    path: 'dataAnlysis',
                    component: () => import('@/views/workspace/dataAnlysis.vue')
                }, {
                    name: 'personalInfo',
                    path: 'personalInfo',
                    component: () => import('@/views/workspace/personalInfo.vue')
                },
                {
                    name: 'projects',
                    path: 'projects',
                    component: () => import('@/views/workspace/projects.vue')
                }, {
                    name: 'newProject',
                    path: 'newProject',
                    component: () => import('@/views/workspace/project/newProject.vue')
                }, {
                    name: 'teams',
                    path: 'teams',
                    component: () => import('@/views/workspace/teams.vue')
                },{
                    name: 'editProject',
                    path: 'editProject',
                    component: () => import('@/views/workspace/project/editProjectFramework.vue'),
                    children:[
                        {
                            name:'index',
                            path: 'index',
                            component: () => import('@/views/workspace/project/basicInfo.vue')
                        },{
                            name:'chapters',
                            path: 'chapters',
                            component: () => import('@/views/workspace/project/chapters.vue')
                        },{
                            name:'characters',
                            path: 'characters',
                            component: () => import('@/views/workspace/project/characters.vue')
                        },{
                            name:'writing',
                            path: 'writing',
                            component: () => import('@/views/workspace/project/writing.vue')
                        },{
                            name:'export',
                            path: 'export',
                            component: () => import('@/views/workspace/project/export.vue')
                        }
                    ]
                }
            ]
        },
        {
            name: 'test',
            path: '/test',
            component: () => import('@/views/test.vue')
        },{
            name: 'communityFramework',
            path: '/community',
            component: () => import('@/views/community/Framework.vue'),
            children:[
                {
                    name: 'sy',
                    path: '',
                    component: () => import('@/views/community/Index.vue')
                  },
                  {
                    path: 'hot',
                    component: () => import('@/views/community/Hot.vue')
                  },
                  {
                    path: 'categories',
                    component: () => import('@/views/community/Categories.vue')
                  },
                  {
                    path: 'search/:keyword',
                    component: () => import('@/views/community/Search.vue')
                  },{
                    path: 'project/:id',
                    name: 'ProjectDetail',
                    component: () => import('@/views/community/ProjectDetail.vue')
                },
                {
                    path: 'chapter/:id',
                    name: 'ChapterView',
                    component: () => import('@/views/workspace/project/chapterView.vue'),
                },{
                    path: 'chat',
                    component: () => import('@/views/community/ChatGroups.vue')

                },{
                    path: 'notifications',
                    component: () => import('@/views/community/Notifications.vue')
                },{
                    path: 'chatS',
                    component: () => import('@/views/community/Chat.vue')
                },{
                    path: 'videofeed',
                    component: () => import('@/views/community/VideoFeed.vue')
                },{
                    path: 'video/:id',
                    component: () => import('@/views/community/VideoFeed.vue')
                }
            ]
        },{
            path: '/admin',
            name: 'admin',
            component: () => import('@/views/admin/framework.vue'),
            
            children:[
                {
                    name: 'adminDashboard',
                    path: 'dashboard',
                    component: () => import('@/views/admin/Dashboard.vue')
                }, {
                    name: 'userManagement',
                    path: 'users',
                    component: () => import('@/views/admin/UserManagement.vue')
                }, {
                    name: 'chapterManagement',
                    path: 'chapters',
                    component: () => import('@/views/admin/ChapterManagement.vue')
                }, {
                    name: 'chapterReview',
                    path: 'chapters/review/:id',
                    component: () => import('@/views/admin/ChapterReview.vue')
                }, {
                    name: 'projectManagement',
                    path: 'projects',
                    component: () => import('@/views/admin/ProjectManagement.vue')
                },{
                    name: "dataAnlysis",
                    path: "statistics",
                    component: () => import('@/views/admin/DataAnlysis.vue')
                }
            ]
        },{
            path: '/team-request/:inviteCode',
            name: 'TeamRequest',
            component: () => import('@/views/workspace/TeamRequest.vue')
          }
    ]
})

// 路由守卫逻辑
router.beforeEach((to, from, next) => {
    const userStore = useUserStore()
    
    // 处理无需登录的路由
    if (to.path === '/' || to.path.startsWith('/auth')) {
        next()
        return
    }
    
    // 处理需要管理员权限的路由
    if (to.matched.some(record => record.meta.requiresAdmin)) {
        if (userStore.isLogin && userStore.user?.permission >= 1) {
            next()
        } else {
            message.info('此页面需要管理员权限，请使用管理员账号登录！');
            setTimeout(() => {
                next('/auth/login')
            }, 2000)
        }
        return
    }
    
    // 处理普通需登录路由
    if (userStore.isLogin) {
        next()
    } else {
        get('/api/user/me', {},
            (message, data) => {
                userStore.login(data)
                next()
            }, (messager, data) => {
                message.info('您还尚未登录，请先登录！');
                setTimeout(() => {
                    next('/auth/login')
                }, 2000)
            }, (messager, data) => {
                message.info('您还尚未登录，请先登录！');
                setTimeout(() => {
                    next('/auth/login')
                }, 2000)
            }
        )
    }
})
export default router
