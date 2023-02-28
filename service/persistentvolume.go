package service

import (
	"context"
	"errors"
	"fmt"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var PersistentVolume persistentVolume

type persistentVolume struct {
}

// GetPersistentVolume 获取pv列表
func (p *persistentVolume) GetPersistentVolume(filterName string, limit, page int) ([]corev1.PersistentVolume, error) {
	list, err := K8s.ClientSet.CoreV1().PersistentVolumes().List(context.TODO(), v1.ListOptions{})
	if err != nil {
		fmt.Println("获取namespace列表失败", err.Error())
		return nil, errors.New("获取namespace列表失败" + err.Error())
	}
	// 组装dataselector
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
	//排序分页
	data := filter.Sort().Paginate()

	return p.fromCells(data.GenericDataList), nil
}

// GetPersistentVolumeDetail 获取pv详情
func (p *persistentVolume) GetPersistentVolumeDetail(pvName string) (node *corev1.PersistentVolume, err error) {
	pvInfo, err := K8s.ClientSet.CoreV1().PersistentVolumes().Get(context.TODO(), pvName, v1.GetOptions{})
	if err != nil {
		fmt.Println("获取pv详情错误", err.Error())
		return nil, errors.New("获取pv详情错误" + err.Error())
	}
	return pvInfo, nil
}

// DeletePersistentVolume 删除pv
func (p *persistentVolume) DeletePersistentVolume(pvName string) (err error) {
	err = K8s.ClientSet.CoreV1().PersistentVolumes().Delete(context.TODO(), pvName, v1.DeleteOptions{})
	if err != nil {
		fmt.Println("删除pv错误", err.Error())
		return errors.New("删除pv错误" + err.Error())
	}
	return nil
}

// toCells方法用于将pod类型数组，转换成DataCell类型数组
func (p *persistentVolume) toCells(std []corev1.PersistentVolume) []DataCell {
	cells := make([]DataCell, len(std))
	for i := range std {
		cells[i] = persistentVolumeCell(std[i])
	}
	return cells
}

// fromCells方法用于将DataCell类型数组，转换成pod类型数组
func (p *persistentVolume) fromCells(cells []DataCell) []corev1.PersistentVolume {
	pvs := make([]corev1.PersistentVolume, len(cells))
	for i := range cells {
		//cells[i].(podCell)就使用到了断言,断言后转换成了podCell类型，然后又转换成了Pod类型
		pvs[i] = corev1.PersistentVolume(cells[i].(persistentVolumeCell))
	}
	return pvs
}
