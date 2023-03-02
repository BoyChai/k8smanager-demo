package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var ConfigMap configMap

type configMap struct{}

// GetConfigMaps 列表
func (c *configMap) GetConfigMaps(filterName, namespace string, limit, page int) ([]corev1.ConfigMap, error) {
	// 获取
	list, err := K8s.ClientSet.CoreV1().ConfigMaps(namespace).List(context.TODO(), v1.ListOptions{})
	if err != nil {
		fmt.Println("获取configmap列表出现错误", err.Error())
		return nil, errors.New("获取configmap列表出现错误" + err.Error())
	}
	// 组装数据进行数据处理
	dataselector := dataSelector{
		GenericDataList: c.toCells(list.Items),
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

	return c.fromCells(data.GenericDataList), nil
}

// GetConfigMapDetail 获取ConfigMap详情
func (c *configMap) GetConfigMapDetail(ConfigMapName, namespace string) (*corev1.ConfigMap, error) {
	cm, err := K8s.ClientSet.CoreV1().ConfigMaps(namespace).Get(context.TODO(), ConfigMapName, v1.GetOptions{})
	if err != nil {
		fmt.Println("获取configmap出现错误", err.Error())
		return nil, errors.New("获取configmap出现错误" + err.Error())
	}
	return cm, nil
}

// DeleteConfigMap 删除ConfigMap
func (c *configMap) DeleteConfigMap(ConfigMapName, namespace string) error {
	err := K8s.ClientSet.CoreV1().ConfigMaps(namespace).Delete(context.TODO(), ConfigMapName, v1.DeleteOptions{})
	if err != nil {
		fmt.Println("删除ConfigMap出现错误", err.Error())
		return errors.New("删除ConfigMap出现错误" + err.Error())
	}
	return nil
}

// UpdateConfigMap 更新Configmap
func (c *configMap) UpdateConfigMap(namespace, content string) error {
	var cm = corev1.ConfigMap{}
	err := json.Unmarshal([]byte(content), &cm)
	if err != nil {
		fmt.Println("序列化出现错误", err.Error())
		return errors.New("序列化出现错误" + err.Error())
	}
	_, err = K8s.ClientSet.CoreV1().ConfigMaps(namespace).Update(context.TODO(), &cm, v1.UpdateOptions{})
	if err != nil {
		fmt.Println("更新configmap出现错误", err.Error())
		return errors.New("更新configmap出现错误" + err.Error())
	}
	return nil
}

// toCells方法用于将pod类型数组，转换成DataCell类型数组
func (c *configMap) toCells(std []corev1.ConfigMap) []DataCell {
	cells := make([]DataCell, len(std))
	for i := range std {
		cells[i] = configMapCell(std[i])
	}
	return cells
}

// fromCells方法用于将DataCell类型数组，转换成pod类型数组
func (c *configMap) fromCells(cells []DataCell) []corev1.ConfigMap {
	cm := make([]corev1.ConfigMap, len(cells))
	for i := range cells {
		cm[i] = corev1.ConfigMap(cells[i].(configMapCell))
	}
	return cm
}
