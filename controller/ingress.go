package controller

import (
	"github.com/gin-gonic/gin"
	"k8smanager-demo/service"
	"net/http"
)

var Ingress ingress

type ingress struct {
}

func (i *ingress) GetIngress(ctx *gin.Context) {
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
	ingresss, err := service.Ingress.Getingresss(params.FilterName, params.Namespace, params.Limit, params.Page)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "获取ingress列表成功",
		"data": ingresss,
	})
}

// GetIngressDetail 获取Ingress详情
func (i *ingress) GetIngressDetail(ctx *gin.Context) {
	params := new(struct {
		IngressName string `form:"ingress_name"`
		Namespace   string `form:"namespace"`
	})
	// form格式使用bind方法绑传入的参数
	if err := ctx.Bind(params); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  "Binding参数失败：" + err.Error(),
			"data": nil,
		})
	}
	ingressInfo, err := service.Ingress.GetIngressDetail(params.IngressName, params.Namespace)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "获取ingress详细信息成功",
		"data": ingressInfo,
	})
}

// CreateIngress 创建Ingress
func (i *ingress) CreateIngress(ctx *gin.Context) {
	var (
		ingressCreate = new(service.IngressCreate)
		err           error
	)
	// form格式使用bind方法绑传入的参数
	if err = ctx.ShouldBindJSON(ingressCreate); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  "Binding参数失败：" + err.Error(),
			"data": nil,
		})
	}
	err = service.Ingress.CreateIngress(ingressCreate)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "创建Ingress成功",
		"data": nil,
	})
}

// DeleteIngress 删除Ingress
func (i *ingress) DeleteIngress(ctx *gin.Context) {
	params := new(struct {
		IngressName string `form:"ingress_name"`
		Namespace   string `form:"namespace"`
	})
	// form格式使用bind方法绑传入的参数
	if err := ctx.ShouldBindJSON(params); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  "Binding参数失败：" + err.Error(),
			"data": nil,
		})
	}
	err := service.Ingress.DeleteIngress(params.IngressName, params.Namespace)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "更新ingress成功",
		"data": nil,
	})
}

// UpdateIngress 更新Ingress
func (i *ingress) UpdateIngress(ctx *gin.Context) {
	params := new(struct {
		Namespace string `json:"namespace"`
		Content   string `json:"content"`
	})
	// form格式使用bind方法绑传入的参数
	if err := ctx.Bind(params); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  "Binding参数失败：" + err.Error(),
			"data": nil,
		})
	}
	//err := service.Pod.UpdatePod(params.PodName,params.Namespace, params.content)
	err := service.Ingress.UpdateIngress(params.Namespace, params.Content)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "更新ingress成功",
		"data": nil,
	})
}
