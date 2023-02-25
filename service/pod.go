package service

import (
	"context"
	"errors"
	"fmt"
	corev1 "k8s.io/api/core/v1"
	matev1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var Pod pod

type pod struct {
}

// PodsResp 定义列表的返回内容，items是pod元素列表，total是元素数量
type PodsResp struct {
	Total int          `json:"total"`
	Items []corev1.Pod `json:"items"`
}

// GetPods 获取pod列表
func (p *pod) GetPods(filterName, namespace string, limit, page int) (podsResp *PodsResp, err error) {
	// 获取podList类型的pod列表
	//context.TODD()用于声明一个空的context上下文，用于List方法内置设置这个请求的超时，这里的常用用法
	// matev1.ListOptions{}用来过滤List数据，如使用label，filed等
	//
	list, err := K8s.ClientSet.CoreV1().Pods(namespace).List(context.TODO(), matev1.ListOptions{})
	if err != nil {
		//写入日志，最后排错用
		fmt.Println("获取pod列表失败:" + err.Error())
		// 返回给浏览器 给前端
		return nil, errors.New(err.Error())
	}
	// 实例化dataselector结果体，组装数据
	selectableData := &dataSelector{
		GenericDataList: p.toCells(list.Items),
		dataSelectQuery: &DataSelectQuery{
			FilterQuery: &FilterQuery{Name: filterName},
			PaginateQuery: &PaginateQuery{
				Limit: limit,
				Page:  page,
			}},
	}
	// 先过滤
	filtered := selectableData.Filter()
	total := len(filtered.GenericDataList)

	// 排序和分页
	data := filtered.Sort().Paginate()
	// DataCell类型转换pod
	pods := p.fromCells(data.GenericDataList)

	// debug
	// 处理前后的比较
	// 处理后的数据
	//fmt.Println("处理后的数据")
	//for _, pod := range pods {
	//	fmt.Println(pod.Name, pod.CreationTimestamp.Time)
	//}
	// 原始数据
	//fmt.Println("原始数据")
	//for _, pod := range list.Items {
	//	fmt.Println(pod.Name, pod.CreationTimestamp.Time)
	//}
	return &PodsResp{Total: total, Items: pods}, nil
}

// toCells方法用于将pod类型数组，转换成DataCell类型数组
func (p *pod) toCells(std []corev1.Pod) []DataCell {
	cells := make([]DataCell, len(std))
	for i := range std {
		cells[i] = podCell(std[i])
	}
	return cells
}

// fromCells方法用于将DataCell类型数组，转换成pod类型数组
func (p *pod) fromCells(cells []DataCell) []corev1.Pod {
	pods := make([]corev1.Pod, len(cells))
	for i := range cells {
		//cells[i].(podCell)就使用到了断言,断言后转换成了podCell类型，然后又转换成了Pod类型
		pods[i] = corev1.Pod(cells[i].(podCell))
	}
	return pods
}
