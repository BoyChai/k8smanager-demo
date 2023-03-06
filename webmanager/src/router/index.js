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
        path: '/home',
        name: '概要',
        icon: 'odometer',
        meta: {title: "概要", requireAuth: true},
        component: () =>       import("@/views/home/Home.vue"),
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