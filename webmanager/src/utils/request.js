import axios from "axios";

// 创建axios对象
const httpClient = axios.create({
    validateStatus(status){
        return status >= 200&& status<=504  // 设置请求的合法状态，若状态码不合法，则不会接收response
    },
    timeout: 10000
})
httpClient.defaults.retry = 3 // 请求重试次数
httpClient.defaults.retryDelay = 1000 // 请求重试时间间隔
httpClient.defaults.shouldRetry = true // 是否重试

//添加请求拦截器
httpClient.interceptors.request.use(
    config => {
//添加header
        config.headers['Content-Type'] = 'application/json'
        config.headers['Accept-Language'] = 'zh-CN'
        config.headers['Authorization'] = localStorage.getItem('token') // 可以全局设置接口请求header中带token
        if (config.method === 'post') {
            if (!config.data) { // 没有参数时，config.data为null，需要转下类型
                config.data = {}
            }
        }
        return config
    },
    err => {
        //Promise.reject()方法返回一个带有拒绝原因的Promise对象，在F12的console中显示报错
        return Promise.reject(err)
    }
);

httpClient.interceptors.response.use(
    response => {
        // 处理状态码
        if (response.status!=200) {
            return Promise.reject(response.data)
        } else {
            return response.data
        }
    }
)
export  default httpClient
