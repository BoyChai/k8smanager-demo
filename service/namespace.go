package service

import (
	"context"
	"errors"
	"fmt"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var Namespace namespace

type namespace struct {
}

// GetNamespaces 获取命名空间列表
func (n *namespace) GetNamespaces(filterName string, limit, page int) ([]corev1.Namespace, error) {
	list, err := K8s.ClientSet.CoreV1().Namespaces().List(context.TODO(), v1.ListOptions{})
	if err != nil {
		fmt.Println("获取namespace列表失败", err.Error())
		return nil, errors.New("获取namespace列表失败" + err.Error())
	}
	// 组装dataselector
	dataselector := dataSelector{
		GenericDataList: n.toCells(list.Items),
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

	return n.fromCells(data.GenericDataList), nil
}

// GetNamespaceDetail 获取namespace详情
func (n *namespace) GetNamespaceDetail(namespaceName string) (node *corev1.Namespace, err error) {
	namespaceInfo, err := K8s.ClientSet.CoreV1().Namespaces().Get(context.TODO(), namespaceName, v1.GetOptions{})
	if err != nil {
		fmt.Println("获取namespace详情错误", err.Error())
		return nil, errors.New("获取namespace详情错误" + err.Error())
	}
	return namespaceInfo, nil
}

// DeleteNamespace 删除Node
func (n *namespace) DeleteNamespace(namespaceName string) (err error) {
	err = K8s.ClientSet.CoreV1().Namespaces().Delete(context.TODO(), namespaceName, v1.DeleteOptions{})
	if err != nil {
		fmt.Println("获取namespace详情错误", err.Error())
		return errors.New("获取namespace详情错误" + err.Error())
	}
	return nil
}

// toCells方法用于将pod类型数组，转换成DataCell类型数组
func (n *namespace) toCells(std []corev1.Namespace) []DataCell {
	cells := make([]DataCell, len(std))
	for i := range std {
		cells[i] = namespaceCell(std[i])
	}
	return cells
}

// fromCells方法用于将DataCell类型数组，转换成pod类型数组
func (n *namespace) fromCells(cells []DataCell) []corev1.Namespace {
	namespaces := make([]corev1.Namespace, len(cells))
	for i := range cells {
		//cells[i].(podCell)就使用到了断言,断言后转换成了podCell类型，然后又转换成了Pod类型
		namespaces[i] = corev1.Namespace(cells[i].(namespaceCell))
	}
	return namespaces
}
