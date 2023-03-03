package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"k8smanager-demo/config"
	"k8smanager-demo/controller"
	"k8smanager-demo/db"
	"k8smanager-demo/service"
)

func main() {
	// 初始化数据库
	db.Init()
	// 初始化gin
	r := gin.Default()
	// 初始化k8s client
	service.K8s.Init()
	// 跨包调用router的初始化方法
	controller.Router.InitApiRouter(r)

	// 启动gin server
	r.Run(config.ListenAddr)

	// 关闭数据库链接
	db.Close()
	fmt.Println("123")
}
