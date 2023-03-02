package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var PersistentVolumeClaim pvc

type pvc struct{}

// GetPersistentVolumeClaims 列表
func (p *pvc) GetPersistentVolumeClaims(filterName, namespace string, limit, page int) ([]corev1.PersistentVolumeClaim, error) {
	// 获取
	list, err := K8s.ClientSet.CoreV1().PersistentVolumeClaims(namespace).List(context.TODO(), v1.ListOptions{})
	if err != nil {
		fmt.Println("获取Secrets列表出现错误", err.Error())
		return nil, errors.New("获取Secrets列表出现错误" + err.Error())
	}
	// 组装数据进行数据处理
	dataselector := dataSelector{
		GenericDataList: p.toCells(list.Items),
		dataSelectQuery: &DataSelectQuery{
			FilterQuery: &FilterQuery{Name: filterName},
			PaginateQuery: &PaginateQuery{
				Limit: limit,
				Page:  page,
			},
		},
	}
	// 过滤
	filter := dataselector.Filter()
	// 排序分页
	data := filter.Sort().Paginate()

	return p.fromCells(data.GenericDataList), nil
}

// GetPersistentVolumeClaimDetail 获取pvc详情
func (p *pvc) GetPersistentVolumeClaimDetail(pvcName, namespace string) (*corev1.PersistentVolumeClaim, error) {
	pc, err := K8s.ClientSet.CoreV1().PersistentVolumeClaims(namespace).Get(context.TODO(), pvcName, v1.GetOptions{})
	if err != nil {
		fmt.Println("获取pvc出现错误", err.Error())
		return nil, errors.New("获取pvc出现错误" + err.Error())
	}
	return pc, nil
}

// DeletePersistentVolumeClaim 删除pvc
func (p *pvc) DeletePersistentVolumeClaim(pvcName, namespace string) error {
	err := K8s.ClientSet.CoreV1().PersistentVolumeClaims(namespace).Delete(context.TODO(), pvcName, v1.DeleteOptions{})
	if err != nil {
		fmt.Println("删除pvc出现错误", err.Error())
		return errors.New("删除pvc出现错误" + err.Error())
	}
	return nil
}

// UpdatePersistentVolumeClaim 更新PersistentVolumeClaim
func (p *pvc) UpdatePersistentVolumeClaim(namespace, content string) error {
	var pc = corev1.PersistentVolumeClaim{}
	err := json.Unmarshal([]byte(content), &pc)
	if err != nil {
		fmt.Println("序列化出现错误", err.Error())
		return errors.New("序列化出现错误" + err.Error())
	}
	_, err = K8s.ClientSet.CoreV1().PersistentVolumeClaims(namespace).Update(context.TODO(), &pc, v1.UpdateOptions{})
	if err != nil {
		fmt.Println("更新pvc出现错误", err.Error())
		return errors.New("更新pvc出现错误" + err.Error())
	}
	return nil
}

// toCells方法用于将secret类型数组，转换成DataCell类型数组
func (p *pvc) toCells(std []corev1.PersistentVolumeClaim) []DataCell {
	cells := make([]DataCell, len(std))
	for i := range std {
		cells[i] = persistentVolumeClaimCell(std[i])
	}
	return cells
}

// fromCells方法用于将DataCell类型数组，转换成sc类型数组
func (p *pvc) fromCells(cells []DataCell) []corev1.PersistentVolumeClaim {
	pc := make([]corev1.PersistentVolumeClaim, len(cells))
	for i := range cells {
		pc[i] = corev1.PersistentVolumeClaim(cells[i].(persistentVolumeClaimCell))
	}
	return pc
}
