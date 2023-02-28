package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"k8smanager-demo/service"
	"net/http"
)

var PersistentVolume persistentvolume

type persistentvolume struct {
}

// GetPersistentVolumes 获取pv列表
func (p *persistentvolume) GetPersistentVolumes(ctx *gin.Context) {
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
	pvs, err := service.PersistentVolume.GetPersistentVolume(params.FilterName, params.Limit, params.Page)
	if err != nil {
		fmt.Println(err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "获取persistentvolume列表成功",
		"data": pvs,
	})
}

// GetPersistentVolumeDetail 获取pv详细信息
func (p *persistentvolume) GetPersistentVolumeDetail(ctx *gin.Context) {
	// 处理入参
	// 匿名结构体，用于定义入参，get请求为from格式，其他请求为json格式
	params := new(struct {
		PersistentVolumeName string `form:"persistentvolume_name"`
	})
	// form格式使用bind方法绑传入的参数
	if err := ctx.Bind(params); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  "Binding参数失败：" + err.Error(),
			"data": nil,
		})
	}
	persistentvolumeInfo, err := service.PersistentVolume.GetPersistentVolumeDetail(params.PersistentVolumeName)
	if err != nil {
		fmt.Println(err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "获取persistentvolume信息成功",
		"data": persistentvolumeInfo,
	})
}

// DeletePersistentVolume 删除pv
func (p *persistentvolume) DeletePersistentVolume(ctx *gin.Context) {
	// 处理入参
	// 匿名结构体，用于定义入参，get请求为from格式，其他请求为json格式
	params := new(struct {
		PersistentVolumeName string `form:"persistentvolume_name"`
	})
	// json格式使用ShouldBindJSON方法绑传入的参数
	if err := ctx.ShouldBindJSON(params); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  "Binding参数失败：" + err.Error(),
			"data": nil,
		})
	}
	err := service.PersistentVolume.DeletePersistentVolume(params.PersistentVolumeName)
	if err != nil {
		fmt.Println(err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "删除pv成功",
		"data": nil,
	})
}
