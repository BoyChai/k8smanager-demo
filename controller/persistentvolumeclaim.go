package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"k8smanager-demo/service"
	"net/http"
)

var PersistentVolumeClaim persistentVolumeClaim

type persistentVolumeClaim struct{}

// GetPersistentVolumeClaims 获取PersistentVolumeClaim列表
func (p *persistentVolumeClaim) GetPersistentVolumeClaims(ctx *gin.Context) {
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
	data, err := service.PersistentVolumeClaim.GetPersistentVolumeClaims(params.FilterName, params.Namespace, params.Limit, params.Page)
	if err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "获取PersistentVolumeClaim列表成功",
		"data": data,
	})
}

// GetPersistentVolumeClaimDetail 获取PersistentVolumeClaim详情
func (p *persistentVolumeClaim) GetPersistentVolumeClaimDetail(ctx *gin.Context) {
	params := new(struct {
		PersistentVolumeClaimName string `form:"persistentvolumeclaim_name"`
		Namespace                 string `form:"namespace"`
	})
	if err := ctx.Bind(params); err != nil {
		fmt.Println("Bind请求参数失败, " + err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	data, err := service.PersistentVolumeClaim.GetPersistentVolumeClaimDetail(params.PersistentVolumeClaimName, params.Namespace)
	if err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "获取PersistentVolumeClaim详细信息成功",
		"data": data,
	})
}

// DeletePersistentVolumeClaim 删除PersistentVolumeClaim
func (p *persistentVolumeClaim) DeletePersistentVolumeClaim(ctx *gin.Context) {
	params := new(struct {
		PersistentVolumeClaimName string `json:"persistentvolumeclaim_name"`
		Namespace                 string `json:"namespace"`
	})
	if err := ctx.ShouldBindJSON(params); err != nil {
		fmt.Println("Bind请求参数失败, " + err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	err := service.PersistentVolumeClaim.DeletePersistentVolumeClaim(params.PersistentVolumeClaimName, params.Namespace)
	if err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "删除pvc成功",
		"data": nil,
	})
}

// UpdatePersistentVolumeClaim 更新PersistentVolumeClaim
func (p *persistentVolumeClaim) UpdatePersistentVolumeClaim(ctx *gin.Context) {
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
	err := service.PersistentVolumeClaim.UpdatePersistentVolumeClaim(params.Namespace, params.Content)
	if err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "更新pvc成功",
		"data": nil,
	})
}
