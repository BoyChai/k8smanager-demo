package controller

import (
	"github.com/gin-gonic/gin"
	"k8smanager-demo/service"
	"net/http"
)

var Pod pod

type pod struct {
}

// GetPods 获取Pod列表，支持分页、过滤、排序
func (p *pod) GetPods(ctx *gin.Context) {
	// 处理入参
	// 匿名结构体，用于定义入参，get请求为from格式，其他请求为json格式
	params := new(struct {
		FilterName string `form:"filter_name"`
		Namespace  string `form:"namespace"`
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
	data, err := service.Pod.GetPods(params.FilterName, params.Namespace, params.Page, params.Limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "获取pod列表成功",
		"data": data,
	})
}
