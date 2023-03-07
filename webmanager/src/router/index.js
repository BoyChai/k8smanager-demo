import {createRouter, createWebHistory} from 'vue-router'
//导入进度条组件
import NProgress from 'nprogress'
import 'nprogress/nprogress.css'
//递增进度条，这将获取当前状态值并添加0.2直到状态为0.994
NProgress.inc(0.2)
//easing 动画字符串
//speed 动画速度
//showSpinner 进度环显示隐藏
NProgress.configure({ easing: 'ease', speed: 600, showSpinner: false })

const routes = [
    {
        path: '/',
        redirect: '/home',
    },
    {
        path: '/home', //视图
        component: import('@/layout/Layout.vue'),
        icon: "odometer", //图标
        meta: {title:"概要", requireAuth: false}, //定义meta元数据
        children: [
            {
                path: '/home', //视图
                name: '概要',
                component: () => import('@/views/home/Home.vue'), //视图组件
                icon: "odometer", //图标
                meta: {title:"概要", requireAuth: false}, //定义meta元数据
            }
        ]
    },
    {
        path: '/workload',
        name: '工作负载',
        component: import('@/layout/Layout.vue'),
        icon: 'menu',
        meta: {title: '工作负载', requireAuth: true},
        children: [
            {
                path: '/workload/deployment',
                name: 'Deployment',
                icon: 'el-icon-s-data',
                meta: {title: 'Deployment', requireAuth: true},
                // component: () => import('@/views/deployment/Deployment.vue')
            },
            {
                path: '/workload/pod',
                name: 'Pod',
                icon: 'el-icon-document-add',
                meta: {title: 'Pod', requireAuth: true},
                // component: () => import('@/views/pod/Pod.vue')
            },
            {
                path: '/workload/daemonset',
                name: 'DaemonSet',
                icon: 'el-icon-document-add',
                meta: {title: 'DaemonSet', requireAuth: true},
                // component: () => import('@/views/daemonset/DaemonSet.vue')
            }
        ]
    },
   {
        path: '/404',
        name: '404',
        meta: {title: "404"},
        component: () =>       import("@/common/404.vue"),
    },{
        path: '/403',
        name: '403',
        meta: {title: "403"},
        component: () =>       import("@/common/403.vue"),
    },
    {
        path: '/:pathMatch(.*)',
        redirect:'/404',
    },
]
//创建路由实例
const router = createRouter({
//hash模式：createWebHashHistory
//history模式：createWebHistory
    history: createWebHistory(),
    routes
})
// 路由守卫,路由拦截
router.beforeEach((to,from,next) => {
    // 启动进度条
    NProgress.start()
    // 设置title
    if (to.meta.title) {
        document.title = to.meta.title
    } else {
        document.title = "Kubernetes"
    }
    //放行
    next()
})

// 关闭进度条
router.afterEach((to, from, failure) => {
    NProgress.done()
})


//抛出路由实例，在main.js中引用
export default router