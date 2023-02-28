package service

import (
	"context"
	"errors"
	"fmt"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var Node node

type node struct {
}

// GetNodes 获取node列表
func (n *node) GetNodes(filterName string, limit, page int) ([]corev1.Node, error) {
	list, err := K8s.ClientSet.CoreV1().Nodes().List(context.TODO(), v1.ListOptions{})
	if err != nil {
		fmt.Println("获取node列表错误", err.Error())
		return nil, errors.New("获取node列表错误" + err.Error())
	}
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
	// 排序分页
	paginate := filter.Sort().Paginate()
	data := n.fromCells(paginate.GenericDataList)
	return data, nil
}

// GetNodeDetail 获取Node详情
func (n *node) GetNodeDetail(nodeName string) (node *corev1.Node, err error) {
	nodeInfo, err := K8s.ClientSet.CoreV1().Nodes().Get(context.TODO(), nodeName, v1.GetOptions{})
	if err != nil {
		fmt.Println("获取Node详情错误", err.Error())
		return nil, errors.New("获取Node详情错误" + err.Error())
	}
	return nodeInfo, nil
}

// toCells方法用于将pod类型数组，转换成DataCell类型数组
func (n *node) toCells(std []corev1.Node) []DataCell {
	cells := make([]DataCell, len(std))
	for i := range std {
		cells[i] = nodeCell(std[i])
	}
	return cells
}

// fromCells方法用于将DataCell类型数组，转换成pod类型数组
func (n *node) fromCells(cells []DataCell) []corev1.Node {
	nodes := make([]corev1.Node, len(cells))
	for i := range cells {
		//cells[i].(podCell)就使用到了断言,断言后转换成了podCell类型，然后又转换成了Pod类型
		nodes[i] = corev1.Node(cells[i].(nodeCell))
	}
	return nodes
}
