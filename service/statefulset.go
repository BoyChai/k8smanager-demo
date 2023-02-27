package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var StatefulSet statefulSet

type statefulSet struct{}

type StatefulSetsResp struct {
	Items []appsv1.StatefulSet `json:"items"`
	Total int                  `json:"total"`
}

// GetStatefulSets 获取StatefulSet列表
func (s *statefulSet) GetStatefulSets(filterName, namespace string, limit, page int) (*StatefulSetsResp, error) {
	list, err := K8s.ClientSet.AppsV1().StatefulSets(namespace).List(context.TODO(), v1.ListOptions{})
	if err != nil {
		fmt.Println("获取StatefulSet列表失败", err.Error())
		return nil, errors.New("获取StatefulSet列表失败" + err.Error())
	}
	// 实例化selectableData
	selectableData := dataSelector{
		GenericDataList: s.toCells(list.Items),
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
	StatefulSets := s.fromCells(data.GenericDataList)
	return &StatefulSetsResp{Total: total, Items: StatefulSets}, nil
}

// GetStatefulSetDetail 获取StatefulSet详情
func (s *statefulSet) GetStatefulSetDetail(statefulSetName, namespace string) (*appsv1.StatefulSet, error) {
	StatefulSetInfo, err := K8s.ClientSet.AppsV1().StatefulSets(namespace).Get(context.TODO(), statefulSetName, v1.GetOptions{})
	if err != nil {
		fmt.Println("获取StatefulSet详情失败", err.Error())
		return nil, errors.New("获取StatefulSet详情失败" + err.Error())
	}
	return StatefulSetInfo, nil
}

// DeleteStatefulSet 删除StatefulSet
func (s *statefulSet) DeleteStatefulSet(statefulSetName, namespace string) error {
	err := K8s.ClientSet.AppsV1().StatefulSets(namespace).Delete(context.TODO(), statefulSetName, v1.DeleteOptions{})
	if err != nil {
		fmt.Println("删除StatefulSet失败", err.Error())
		return errors.New("删除StatefulSet失败" + err.Error())
	}
	return nil
}

// UpdateStatefulSet 更新StatefulSet
func (s *statefulSet) UpdateStatefulSet(namespace, content string) error {
	var stateful = &appsv1.StatefulSet{}
	err := json.Unmarshal([]byte(content), &stateful)
	if err != nil {
		fmt.Println("序列化StatefulSet失败", err.Error())
		return errors.New("序列化StatefulSet失败" + err.Error())
	}
	_, err = K8s.ClientSet.AppsV1().StatefulSets(namespace).Update(context.TODO(), stateful, v1.UpdateOptions{})
	if err != nil {
		fmt.Println("更新StatefulSet失败", err.Error())
		return errors.New("更新StatefulSet失败" + err.Error())
	}
	return nil
}

// toCells方法用于将pod类型数组，转换成DataCell类型数组
func (s *statefulSet) toCells(std []appsv1.StatefulSet) []DataCell {
	cells := make([]DataCell, len(std))
	for i := range std {
		cells[i] = statefulSetCell(std[i])
	}
	return cells
}

// fromCells方法用于将DataCell类型数组，转换成pod类型数组
func (s *statefulSet) fromCells(cells []DataCell) []appsv1.StatefulSet {
	pods := make([]appsv1.StatefulSet, len(cells))
	for i := range cells {
		//cells[i].(podCell)就使用到了断言,断言后转换成了podCell类型，然后又转换成了Pod类型
		pods[i] = appsv1.StatefulSet(cells[i].(statefulSetCell))
	}
	return pods
}
