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
		PUT("/api/k8s/statefulset/update", StatefulSet.UpdateStatefulSet).
		// node资源
		GET("/api/k8s/nodes", Node.GetNodes).
		GET("/api/k8s/node/detail", Node.GetNodeDetail).
		// namespace资源
		GET("/api/k8s/namespaces", Namespace.GetNamespaces).
		GET("/api/k8s/namespace/detail", Namespace.GetNamespaceDetail).
		DELETE("/api/k8s/namespace/del", Namespace.DeleteNamespace).
		// pv资源
		GET("/api/k8s/persistentvolumes", PersistentVolume.GetPersistentVolumes).
		GET("/api/k8s/persistentvolume/detail", PersistentVolume.GetPersistentVolumeDetail).
		DELETE("/api/k8s/persistentvolume/del", PersistentVolume.DeletePersistentVolume).
		// service
		GET("/api/k8s/services", Service.GetServices).
		GET("/api/k8s/service/detail", Service.GetServiceDetail).
		PUT("/api/k8s/service/create", Service.CreateService).
		DELETE("/api/k8s/service/delete", Service.DeleteService).
		PUT("/api/k8s/service/update", Service.UpdateService).
		// ingress
		GET("/api/k8s/ingresss", Ingress.GetIngress).
		GET("/api/k8s/ingress/detail", Ingress.GetIngressDetail).
		PUT("/api/k8s/ingress/create", Ingress.CreateIngress).
		DELETE("/api/k8s/ingress/delete", Ingress.DeleteIngress).
		PUT("/api/k8s/ingress/update", Ingress.UpdateIngress).
		// configmap
		GET("/api/k8s/configmaps", ConfigMap.GetConfigMaps).
		GET("/api/k8s/configmap/detail", ConfigMap.GetConfigMapDetail).
		DELETE("/api/k8s/configmap/del", ConfigMap.DeleteConfigMap).
		PUT("/api/k8s/configmap/update", ConfigMap.UpdateConfigMap).
		//secret
		GET("/api/k8s/secrets", Secret.GetSecrets).
		GET("/api/k8s/secret/detail", Secret.GetSecretDetail).
		DELETE("/api/k8s/secret/del", Secret.DeleteSecret).
		PUT("/api/k8s/secret/update", Secret.UpdateSecret).
		// pvc
		GET("/api/k8s/persistentvolumeclaims", PersistentVolumeClaim.GetPersistentVolumeClaims).
		GET("/api/k8s/persistentvolumeclaim/detail", PersistentVolumeClaim.GetPersistentVolumeClaimDetail).
		DELETE("/api/k8s/persistentvolumeclaim/del", PersistentVolumeClaim.DeletePersistentVolumeClaim).
		PUT("/api/k8s/persistentvolumeclaim/update", PersistentVolumeClaim.UpdatePersistentVolumeClaim)
}
