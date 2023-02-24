package service

import (
	"fmt"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8smanager-demo/config"
)

// K8s 用于初始化k8s clientSet

var K8s k8s

type k8s struct {
	ClientSet *kubernetes.Clientset
}

func (k *k8s) Init() {
	// 获取k8s配置
	conf, err := clientcmd.BuildConfigFromFlags("", config.KubeConfig)
	if err != nil {
		panic("获取kubeconfig配置失败：" + err.Error())
	}
	// 根据rest.confi类型对象创建一个clientSet
	clientSet, err := kubernetes.NewForConfig(conf)
	if err != nil {
		panic("创建clientSet对象失败：" + err.Error())
	} else {
		fmt.Println("clientSet 初始化成功")
	}
	k.ClientSet = clientSet
}
