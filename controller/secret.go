package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"k8smanager-demo/service"
	"net/http"
)

var Secret secret

type secret struct{}

// GetSecrets 获取Secret列表
func (s *secret) GetSecrets(ctx *gin.Context) {
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
	data, err := service.Secret.GetSecretss(params.FilterName, params.Namespace, params.Limit, params.Page)
	if err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "获取Secret列表成功",
		"data": data,
	})
}

// GetSecretDetail 获取Secret详情
func (s *secret) GetSecretDetail(ctx *gin.Context) {
	params := new(struct {
		SecretName string `form:"secret_name"`
		Namespace  string `form:"namespace"`
	})
	if err := ctx.Bind(params); err != nil {
		fmt.Println("Bind请求参数失败, " + err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	data, err := service.Secret.GetSecretsDetail(params.SecretName, params.Namespace)
	if err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "获取configmap详细信息成功",
		"data": data,
	})
}

// DeleteSecret 删除Secret
func (s *secret) DeleteSecret(ctx *gin.Context) {
	params := new(struct {
		SecretName string `json:"secret_name"`
		Namespace  string `json:"namespace"`
	})
	if err := ctx.ShouldBindJSON(params); err != nil {
		fmt.Println("Bind请求参数失败, " + err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	err := service.Secret.DeleteSecrets(params.SecretName, params.Namespace)
	if err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "删除Secret成功",
		"data": nil,
	})
}

// UpdateSecret 更新Secret
func (s *secret) UpdateSecret(ctx *gin.Context) {
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
	err := service.Secret.UpdateSecrets(params.Namespace, params.Content)
	if err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "更新Secret成功",
		"data": nil,
	})
}
