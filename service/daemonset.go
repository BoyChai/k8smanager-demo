package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var DaemonSet daemonSet

type daemonSet struct{}

type DaemonSetsResp struct {
	Items []appsv1.DaemonSet `json:"items"`
	Total int                `json:"total"`
}

// GetDaemonSets 获取DaemonSet列表
func (d *daemonSet) GetDaemonSets(filterName, namespace string, limit, page int) (*DaemonSetsResp, error) {
	list, err := K8s.ClientSet.AppsV1().DaemonSets(namespace).List(context.TODO(), v1.ListOptions{})
	if err != nil {
		fmt.Println("获取DaemonSet列表失败", err.Error())
		return nil, errors.New("获取DaemonSet列表失败" + err.Error())
	}
	// 实例化selectableData
	selectableData := dataSelector{
		GenericDataList: d.toCells(list.Items),
		dataSelectQuery: &DataSelectQuery{
			FilterQuery: &FilterQuery{
				Name: filterName,
			},
			PaginateQuery: &PaginateQuery{
				Limit: limit,
				Page:  page,
			},
		},
	}
	// 过滤
	filterd := selectableData.Filter()
	total := len(filterd.GenericDataList)

	// 排序分页
	data := filterd.Sort().Paginate()

	// 转换
	DaemonSets := d.fromCells(data.GenericDataList)
	return &DaemonSetsResp{Total: total, Items: DaemonSets}, nil
}

// GetDaemonSetDetail 获取DaemonSet详情
func (d *daemonSet) GetDaemonSetDetail(daemonSetName, namespace string) (*appsv1.DaemonSet, error) {
	daemonSetInfo, err := K8s.ClientSet.AppsV1().DaemonSets(namespace).Get(context.TODO(), daemonSetName, v1.GetOptions{})
	if err != nil {
		fmt.Println("获取DaemonSet详情失败", err.Error())
		return nil, errors.New("获取DaemonSet详情失败" + err.Error())
	}
	return daemonSetInfo, nil
}

// DeleteDaemonSet 删除DaemonSet
func (d *daemonSet) DeleteDaemonSet(daemonSetName, namespace string) error {
	err := K8s.ClientSet.AppsV1().DaemonSets(namespace).Delete(context.TODO(), daemonSetName, v1.DeleteOptions{})
	if err != nil {
		fmt.Println("删除DaemonSet失败", err.Error())
		return errors.New("删除DaemonSet失败" + err.Error())
	}
	return nil
}

// UpdateDaemonSet 更新DaemonSet
func (d daemonSet) UpdateDaemonSet(namespace, content string) error {
	var daemon = &appsv1.DaemonSet{}
	err := json.Unmarshal([]byte(content), &daemon)
	if err != nil {
		fmt.Println("序列化DaemonSet失败", err.Error())
		return errors.New("序列化DaemonSet失败" + err.Error())
	}
	_, err = K8s.ClientSet.AppsV1().DaemonSets(namespace).Update(context.TODO(), daemon, v1.UpdateOptions{})
	if err != nil {
		fmt.Println("更新DaemonSet失败", err.Error())
		return errors.New("更新DaemonSet失败" + err.Error())
	}
	return nil
}

// toCells方法用于将pod类型数组，转换成DataCell类型数组
func (d *daemonSet) toCells(std []appsv1.DaemonSet) []DataCell {
	cells := make([]DataCell, len(std))
	for i := range std {
		cells[i] = daemonSetCell(std[i])
	}
	return cells
}

// fromCells方法用于将DataCell类型数组，转换成pod类型数组
func (d *daemonSet) fromCells(cells []DataCell) []appsv1.DaemonSet {
	DaemonSets := make([]appsv1.DaemonSet, len(cells))
	for i := range cells {
		//cells[i].(podCell)就使用到了断言,断言后转换成了podCell类型，然后又转换成了Pod类型
		DaemonSets[i] = appsv1.DaemonSet(cells[i].(daemonSetCell))
	}
	return DaemonSets
}
