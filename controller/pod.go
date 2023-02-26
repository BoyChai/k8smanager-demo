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

// GetPodDetail 获取pod详情
func (p *pod) GetPodDetail(ctx *gin.Context) {
	// 处理入参
	// 匿名结构体，用于定义入参，get请求为from格式，其他请求为json格式
	params := new(struct {
		PodName   string `form:"pod_name"`
		Namespace string `form:"namespace"`
	})
	// form格式使用bind方法绑传入的参数
	if err := ctx.Bind(params); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  "Binding参数失败：" + err.Error(),
			"data": nil,
		})
	}
	podInfo, err := service.Pod.GetPodDetail(params.PodName, params.Namespace)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "获取pod详细信息成功",
		"data": podInfo,
	})
}

// DeletePod 删除某个pod
func (p *pod) DeletePod(ctx *gin.Context) {
	// 处理入参
	// 匿名结构体，用于定义入参，get请求为from格式，其他请求为json格式
	params := new(struct {
		PodName   string `json:"pod_name"`
		Namespace string `json:"namespace"`
	})
	// form格式使用bind方法绑传入的参数
	if err := ctx.Bind(params); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  "Binding参数失败：" + err.Error(),
			"data": nil,
		})
	}
	err := service.Pod.DeletePod(params.PodName, params.Namespace)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "删除pod成功",
		"data": nil,
	})
}

// UpdatePod 更新pod
func (p *pod) UpdatePod(ctx *gin.Context) {
	// 处理入参
	// 匿名结构体，用于定义入参，get请求为from格式，其他请求为json格式
	params := new(struct {
		//PodName   string `form:"pod_name"`
		Namespace string `json:"namespace"`
		content   string `json:"content"`
	})
	// form格式使用bind方法绑传入的参数
	if err := ctx.Bind(params); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  "Binding参数失败：" + err.Error(),
			"data": nil,
		})
	}
	//err := service.Pod.UpdatePod(params.PodName,params.Namespace, params.content)
	err := service.Pod.UpdatePod(params.Namespace, params.content)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "更新pod成功",
		"data": nil,
	})
}

// GetPodContainer 获取pod容器
func (p *pod) GetPodContainer(ctx *gin.Context) {
	// 处理入参
	// 匿名结构体，用于定义入参，get请求为from格式，其他请求为json格式
	params := new(struct {
		PodName   string `form:"pod_name"`
		Namespace string `form:"namespace"`
	})
	// form格式使用bind方法绑传入的参数
	if err := ctx.Bind(params); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  "Binding参数失败：" + err.Error(),
			"data": nil,
		})
	}
	containers, err := service.Pod.GetPodContainer(params.PodName, params.Namespace)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "获取pod容器成功",
		"data": containers,
	})
}

// GetPodLog 获取pod容器日志
func (p *pod) GetPodLog(ctx *gin.Context) {
	// 处理入参
	// 匿名结构体，用于定义入参，get请求为from格式，其他请求为json格式
	params := new(struct {
		ContainerName string `form:"container_name"`
		PodName       string `form:"pod_name"`
		Namespace     string `form:"namespace"`
	})
	// form格式使用bind方法绑传入的参数
	if err := ctx.Bind(params); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  "Binding参数失败：" + err.Error(),
			"data": nil,
		})
	}
	log, err := service.Pod.GetPodLog(params.ContainerName, params.PodName, params.Namespace)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "获取pod容器日志成功",
		"data": log,
	})
}

// GetPodNumPerNp 获取命名空间的pod数量
func (p *pod) GetPodNumPerNp(ctx *gin.Context) {
	nps, err := service.Pod.GetPodNumPerNp()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "获取pod容器日志成功",
		"data": nps,
	})
}
