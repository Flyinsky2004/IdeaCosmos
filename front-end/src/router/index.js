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
                    name: 'dataAnlysis',
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
                }
            ]
        },
    ]
})

// BAD
router.beforeEach((to, from, next) => {
    const userStore = useUserStore()
    if (to.path === '/' || to.path.startsWith('/auth')) {
        next()
    } else {
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

    }
})
export default router
