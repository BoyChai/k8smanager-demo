package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	corev1 "k8s.io/api/core/v1"
	matev1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8smanager-demo/config"
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

// GetPodDetail 获取pod详情
func (p *pod) GetPodDetail(podName, namespace string) (pod *corev1.Pod, err error) {
	pod, err = K8s.ClientSet.CoreV1().Pods(namespace).Get(context.TODO(), podName, matev1.GetOptions{})
	if err != nil {
		fmt.Println("获取pod详细信息失败", err)
		return nil, errors.New("获取pod详细信息失败" + err.Error())
	}
	return pod, nil
}

// DeletePod 删除pod
func (p *pod) DeletePod(podName, namespace string) (err error) {
	err = K8s.ClientSet.CoreV1().Pods(namespace).Delete(context.TODO(), podName, matev1.DeleteOptions{})
	if err != nil {
		fmt.Println("删除pod失败", err)
		return errors.New("删除pod失败" + err.Error())
	}
	return nil
}

// UpdatePod 更新pod
// func (p *pod) UpdatePod(podName, namespace, content string) (err error) {
func (p *pod) UpdatePod(namespace, content string) (err error) {
	var pod = &corev1.Pod{}
	// 反序列化为pod对象
	err = json.Unmarshal([]byte(content), &pod)
	if err != nil {
		fmt.Println("反序列化失败", err.Error())
		return errors.New("反序列化失败" + err.Error())
	}
	// 更新pod
	_, err = K8s.ClientSet.CoreV1().Pods(namespace).Update(context.TODO(), pod, matev1.UpdateOptions{})
	if err != nil {
		fmt.Println("更新pod失败", err.Error())
		return errors.New("更新pod失败" + err.Error())
	}

	return nil
}

// GetPodContainer 获取pod中的容器名称
func (p *pod) GetPodContainer(podName, namespace string) (containers []string, err error) {
	// 获取pod详情
	detail, err := p.GetPodDetail(podName, namespace)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	//从pod中拿容器名称
	for _, container := range detail.Spec.Containers {
		containers = append(containers, container.Name)
	}

	return containers, nil
}

// GetPodLog 获取pod内容器日志
func (p *pod) GetPodLog(containerName, podName, namespace string) (log string, err error) {
	// 设置容器配额，容器名称，tail行数
	lineLimit := int64(config.PodLogTailLine)
	options := &corev1.PodLogOptions{
		Container: containerName,
		TailLines: &lineLimit,
	}

	// 获取一个request实例
	req := K8s.ClientSet.CoreV1().Pods(namespace).GetLogs(podName, options)
	// 发起stream链接得到response.body
	podlogs, err := req.Stream(context.TODO())
	if err != nil {
		fmt.Println("获取日志请求失败", err.Error())
		return "", errors.New("获取日志请求失败" + err.Error())
	}
	defer podlogs.Close()
	// 把response.body写入到缓冲区，目的是为了转换string类型
	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, podlogs)
	if err != nil {
		fmt.Println("复制PodLog错误", err.Error())
		return "", errors.New("复制PodLog错误" + err.Error())
	}
	return buf.String(), nil
}

type PodsNp struct {
	Namespace string
	PodNum    int
}

// GetPodNumPerNp 获取每个namespace的pod数量
func (p *pod) GetPodNumPerNp() (podNps []*PodsNp, err error) {
	// 获取namespace列表
	nsList, err := K8s.ClientSet.CoreV1().Namespaces().List(context.TODO(), matev1.ListOptions{})
	if err != nil {
		fmt.Println("获取命名空间列表错误,", err.Error())
		return nil, errors.New("获取命名空间错误" + err.Error())
	}
	// 根据命名空间去获取pod
	for _, item := range nsList.Items {
		podList, err := K8s.ClientSet.CoreV1().Pods(item.Name).List(context.TODO(), matev1.ListOptions{})
		if err != nil {
			fmt.Println("获取pod列表失败,", err.Error())
			return nil, errors.New("获取pod列表失败" + err.Error())
		}
		// 组装数据计算数量
		podsNp := &PodsNp{
			Namespace: item.Name,
			PodNum:    len(podList.Items),
		}
		podNps = append(podNps, podsNp)
	}
	return podNps, nil
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
