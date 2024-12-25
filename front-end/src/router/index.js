import {createRouter, createWebHistory} from 'vue-router'
import {useUserStore} from "@/stores/user.js";
import {ElMessage} from "element-plus";
import {get} from "@/util/request.js";

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
        }

    ],
})

// BAD
router.beforeEach((to, from, next) => {
    const userStore = useUserStore()
    if (to.path === '/' || to.path.startsWith('/auth')) {
        next()
    } else {
        if (userStore.isLogin) {
            next()
        }
        get('/api/user/me', {},
            (message, data) => {
                userStore.login(data)
            }
        )
        if (!userStore.isLogin) {
            ElMessage.info('您还尚未登录，请先登录！');
            setTimeout(() => {
                router.push('/auth/login')
            }, 2000)
        }else{
            next()
        }
    }
})
export default router
