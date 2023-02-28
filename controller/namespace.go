package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"k8smanager-demo/service"
	"net/http"
)

var Namespace namespace

type namespace struct {
}

// GetNamespaces 获取namespace列表
func (n *namespace) GetNamespaces(ctx *gin.Context) {
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
	namesapces, err := service.Namespace.GetNamespaces(params.FilterName, params.Limit, params.Page)
	if err != nil {
		fmt.Println(err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "获取namespace列表成功",
		"data": namesapces,
	})
}

// GetNamespaceDetail 获取namespace详细信息
func (n *namespace) GetNamespaceDetail(ctx *gin.Context) {
	// 处理入参
	// 匿名结构体，用于定义入参，get请求为from格式，其他请求为json格式
	params := new(struct {
		NamespaceName string `form:"namespace_name"`
	})
	// form格式使用bind方法绑传入的参数
	if err := ctx.Bind(params); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  "Binding参数失败：" + err.Error(),
			"data": nil,
		})
	}
	namespaceInfo, err := service.Namespace.GetNamespaceDetail(params.NamespaceName)
	if err != nil {
		fmt.Println(err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "获取namespace信息成功",
		"data": namespaceInfo,
	})
}

// DeleteNamespace 删除namespace
func (n *namespace) DeleteNamespace(ctx *gin.Context) {
	// 处理入参
	// 匿名结构体，用于定义入参，get请求为from格式，其他请求为json格式
	params := new(struct {
		NamespaceName string `form:"namespace_name"`
	})
	// json格式使用ShouldBindJSON方法绑传入的参数
	if err := ctx.ShouldBindJSON(params); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  "Binding参数失败：" + err.Error(),
			"data": nil,
		})
	}
	err := service.Namespace.DeleteNamespace(params.NamespaceName)
	if err != nil {
		fmt.Println(err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "删除namespace成功",
		"data": nil,
	})
}
