package controller

import (
	"github.com/gin-gonic/gin"
	"k8smanager-demo/service"
	"net/http"
)

var Service svc

type svc struct {
}

func (s *svc) GetServices(ctx *gin.Context) {
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
	svcs, err := service.Service.GetServices(params.FilterName, params.Namespace, params.Limit, params.Page)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "获取svc列表成功",
		"data": svcs,
	})
}

// GetServiceDetail 获取Service详情
func (s *svc) GetServiceDetail(ctx *gin.Context) {
	params := new(struct {
		ServiceName string `form:"service_name"`
		Namespace   string `form:"namespace"`
	})
	// form格式使用bind方法绑传入的参数
	if err := ctx.Bind(params); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  "Binding参数失败：" + err.Error(),
			"data": nil,
		})
	}
	svcInfo, err := service.Service.GetServiceDetail(params.ServiceName, params.Namespace)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "获取svc详细信息成功",
		"data": svcInfo,
	})
}

// CreateService 创建Service
func (s *svc) CreateService(ctx *gin.Context) {
	var (
		svcCreate = new(service.Create)
		err       error
	)
	// form格式使用bind方法绑传入的参数
	if err = ctx.ShouldBindJSON(svcCreate); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  "Binding参数失败：" + err.Error(),
			"data": nil,
		})
	}
	svcInfo, err := service.Service.CreateService(svcCreate)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "创建service                                                                                                                                                  成功",
		"data": svcInfo,
	})
}

// DeleteService 删除Service
func (s *svc) DeleteService(ctx *gin.Context) {
	params := new(struct {
		ServiceName string `form:"service_name"`
		Namespace   string `form:"namespace"`
	})
	// form格式使用bind方法绑传入的参数
	if err := ctx.ShouldBindJSON(params); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  "Binding参数失败：" + err.Error(),
			"data": nil,
		})
	}
	err := service.Service.DeleteService(params.ServiceName, params.Namespace)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "更新service成功",
		"data": nil,
	})
}

// UpdateService 更新Service
func (s *svc) UpdateService(ctx *gin.Context) {
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
	err := service.Service.UpdateService(params.Namespace, params.Content)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "更新svc成功",
		"data": nil,
	})
}
