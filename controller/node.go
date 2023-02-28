package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"k8smanager-demo/service"
	"net/http"
)

var Node node

type node struct {
}

// GetNodes 获取node列表
func (s *node) GetNodes(ctx *gin.Context) {
	// 处理入参
	// 匿名结构体，用于定义入参，get请求为from格式，其他请求为json格式
	params := new(struct {
		FilterName string `form:"filter_name"`
		Limit      int    `form:"limit"`
		Page       int    `form:"page"`
	})
	// form格式使用bind方法绑传入的参数
	if err := ctx.Bind(params); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  "Binding参数失败：" + err.Error(),
			"data": nil,
		})
	}
	nodes, err := service.Node.GetNodes(params.FilterName, params.Limit, params.Page)
	if err != nil {
		fmt.Println(err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "获取node列表成功",
		"data": nodes,
	})
}

// GetNodeDetail 获取node详细信息
func (s *node) GetNodeDetail(ctx *gin.Context) {
	// 处理入参
	// 匿名结构体，用于定义入参，get请求为from格式，其他请求为json格式
	params := new(struct {
		NodeName string `form:"node_name"`
	})
	// form格式使用bind方法绑传入的参数
	if err := ctx.Bind(params); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  "Binding参数失败：" + err.Error(),
			"data": nil,
		})
	}
	nodeInfo, err := service.Node.GetNodeDetail(params.NodeName)
	if err != nil {
		fmt.Println(err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "获取node信息成功",
		"data": nodeInfo,
	})
}
