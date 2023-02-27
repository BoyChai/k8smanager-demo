package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"k8smanager-demo/service"

	"net/http"
)

var DaemonSet daemonSet

type daemonSet struct {
}

// GetDaemonSets 获取DaemonSet列表，支持过滤、排序、分页
func (d *daemonSet) GetDaemonSets(ctx *gin.Context) {
	params := new(struct {
		FilterName string `form:"filter_name"`
		Namespace  string `form:"namespace"`
		Page       int    `form:"page"`
		Limit      int    `form:"limit"`
	})
	if err := ctx.Bind(params); err != nil {
		fmt.Println("Bind请求参数失败, " + err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	data, err := service.DaemonSet.GetDaemonSets(params.FilterName, params.Namespace, params.Limit, params.Page)
	if err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "获取DaemonSet列表成功",
		"data": data,
	})
}

// GetDaemonSetDetail 获取DaemonSet详情
func (d *daemonSet) GetDaemonSetDetail(ctx *gin.Context) {
	params := new(struct {
		daemonSetName string `form:"daemonset_name"`
		Namespace     string `form:"namespace"`
	})
	if err := ctx.Bind(params); err != nil {
		fmt.Println("Bind请求参数失败, " + err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	data, err := service.DaemonSet.GetDaemonSetDetail(params.daemonSetName, params.Namespace)
	if err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "获取DaemonSet详细信息成功",
		"data": data,
	})
}

// DeleteDaemonSet 删除DaemonSet
func (d *daemonSet) DeleteDaemonSet(ctx *gin.Context) {
	params := new(struct {
		DaemonSetName string `json:"daemonset_name"`
		Namespace     string `json:"namespace"`
	})
	if err := ctx.ShouldBindJSON(params); err != nil {
		fmt.Println("Bind请求参数失败, " + err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	err := service.DaemonSet.DeleteDaemonSet(params.DaemonSetName, params.Namespace)
	if err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "获取DaemonSet详细信息成功",
		"data": nil,
	})
}

// UpdateDaemonSet 更新DaemonSet
func (d *daemonSet) UpdateDaemonSet(ctx *gin.Context) {
	params := new(struct {
		Namespace string `json:"namespace"`
		Content   string `json:"content"`
	})
	if err := ctx.ShouldBindJSON(params); err != nil {
		fmt.Println("Bind请求参数失败, " + err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	err := service.DaemonSet.UpdateDaemonSet(params.Namespace, params.Content)
	if err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "获取DaemonSet详细信息成功",
		"data": nil,
	})
}
