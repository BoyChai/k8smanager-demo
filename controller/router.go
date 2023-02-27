package controller

import (
	"github.com/gin-gonic/gin"
)

// Router 实例化router类型对象，首字母大写用于跨包调用
var Router router

// 声明router结构体
type router struct{}

// InitApiRouter 初始化路由规则
func (r *router) InitApiRouter(router *gin.Engine) {
	router.
		// pod 资源
		GET("/api/k8s/pods", Pod.GetPods).
		GET("/api/k8s/pod/detail", Pod.GetPodDetail).
		DELETE("/api/k8s/pod/del", Pod.DeletePod).
		PUT("/api/k8s/pod/update", Pod.UpdatePod).
		GET("/api/k8s/pod/container", Pod.GetPodContainer).
		GET("/api/k8s/pod/log", Pod.GetPodLog).
		GET("/api/k8s/pod/numnp", Pod.GetPodNumPerNp).
		// deployment资源
		GET("/api/k8s/deployments", Deployment.GetDeployments).
		GET("/api/k8s/deployment/detail", Deployment.GetDeploymentDetail).
		PUT("/api/k8s/deployment/scale", Deployment.ScaleDeployment).
		DELETE("/api/k8s/deployment/del", Deployment.DeleteDeployment).
		PUT("/api/k8s/deployment/restart", Deployment.RestartDeployment).
		PUT("/api/k8s/deployment/update", Deployment.UpdateDeployment).
		GET("/api/k8s/deployment/numnp", Deployment.GetDeployNumPerNp).
		POST("/api/k8s/deployment/create", Deployment.CreateDeployment).
		// daemonset资源
		GET("/api/k8s/daemonsets", DaemonSet.GetDaemonSets).
		GET("/api/k8s/daemonset/detail", DaemonSet.GetDaemonSetDetail).
		DELETE("/api/k8s/daemonset/del", DaemonSet.DeleteDaemonSet).
		PUT("/api/k8s/daemonset/update", DaemonSet.UpdateDaemonSet).
		// statefulset资源
		GET("/api/k8s/statefulsets", StatefulSet.GetStatefulSets).
		GET("/api/k8s/statefulset/detail", StatefulSet.GetStatefulSetDetail).
		DELETE("/api/k8s/statefulset/del", StatefulSet.DeleteStatefulSet).
		PUT("/api/k8s/statefulset/update", StatefulSet.UpdateStatefulSet)

}
