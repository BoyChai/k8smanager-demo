package controller

import (
	"github.com/gin-gonic/gin"
)

// Router 实例化router类型对象，首字母大写用于跨包调用
var Router router

// 声明router结构体
type router struct{}

// InitApiRouter 初始化路由规则
func (r *router) InitApiRouter(router *gin.Engine) {
	router.
		GET("/api/k8s/pods", Pod.GetPods)
}
