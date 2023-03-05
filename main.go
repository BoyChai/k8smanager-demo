package main

import (
	"github.com/gin-gonic/gin"
	"k8smanager-demo/config"
	"k8smanager-demo/controller"
	"k8smanager-demo/middle"
	"k8smanager-demo/service"
)

func main() {
	// 初始化gin
	r := gin.Default()
	// 初始化k8s client
	service.K8s.Init()
	// 加载跨域中间件
	// 加载中间件需要放在其他接口加载之前
	r.Use(middle.Cors())
	//jwt token验证
	r.Use(middle.JWTAuth())
	// 跨包调用router的初始化方法
	controller.Router.InitApiRouter(r)
	// 启动gin server
	r.Run(config.ListenAddr)
}
