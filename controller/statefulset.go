package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"k8smanager-demo/service"
	"net/http"
)

var StatefulSet statefulSet

type statefulSet struct {
}

// GetStatefulSets 获取StatefulSet列表，支持过滤、排序、分页
func (s *statefulSet) GetStatefulSets(ctx *gin.Context) {
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
	data, err := service.StatefulSet.GetStatefulSets(params.FilterName, params.Namespace, params.Limit, params.Page)
	if err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "获取deployment列表成功",
		"data": data,
	})
}

// GetStatefulSetDetail 获取StatefulSet详情
func (s *statefulSet) GetStatefulSetDetail(ctx *gin.Context) {
	params := new(struct {
		StatefulSetName string `form:"statefulset_name"`
		Namespace       string `form:"namespace"`
	})
	if err := ctx.Bind(params); err != nil {
		fmt.Println("Bind请求参数失败, " + err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	data, err := service.StatefulSet.GetStatefulSetDetail(params.StatefulSetName, params.Namespace)
	if err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "获取statefulset详细信息成功",
		"data": data,
	})
}

// DeleteStatefulSet 删除StatefulSet
func (s *statefulSet) DeleteStatefulSet(ctx *gin.Context) {
	params := new(struct {
		StatefulSetName string `form:"statefulset_name"`
		Namespace       string `form:"namespace"`
	})
	if err := ctx.Bind(params); err != nil {
		fmt.Println("Bind请求参数失败, " + err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	err := service.StatefulSet.DeleteStatefulSet(params.StatefulSetName, params.Namespace)
	if err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "删除statefulset成功",
		"data": nil,
	})
}

// UpdateStatefulSet 更新StatefulSe
func (s *statefulSet) UpdateStatefulSet(ctx *gin.Context) {
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
	err := service.StatefulSet.UpdateStatefulSet(params.Namespace, params.Content)
	if err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "更新statefulset成功",
		"data": nil,
	})
}
