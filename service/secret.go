package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var Secret secret

type secret struct{}

// GetSecretss 列表
func (s *secret) GetSecretss(filterName, namespace string, limit, page int) ([]corev1.Secret, error) {
	// 获取
	list, err := K8s.ClientSet.CoreV1().Secrets(namespace).List(context.TODO(), v1.ListOptions{})
	if err != nil {
		fmt.Println("获取Secrets列表出现错误", err.Error())
		return nil, errors.New("获取Secrets列表出现错误" + err.Error())
	}
	// 组装数据进行数据处理
	dataselector := dataSelector{
		GenericDataList: s.toCells(list.Items),
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

	return s.fromCells(data.GenericDataList), nil
}

// GetSecretsDetail 获取Secrets详情
func (s *secret) GetSecretsDetail(SecretsName, namespace string) (*corev1.Secret, error) {
	sc, err := K8s.ClientSet.CoreV1().Secrets(namespace).Get(context.TODO(), SecretsName, v1.GetOptions{})
	if err != nil {
		fmt.Println("获取Secrets出现错误", err.Error())
		return nil, errors.New("获取Secrets出现错误" + err.Error())
	}
	return sc, nil
}

// DeleteSecrets 删除Secrets
func (s *secret) DeleteSecrets(SecretsName, namespace string) error {
	err := K8s.ClientSet.CoreV1().Secrets(namespace).Delete(context.TODO(), SecretsName, v1.DeleteOptions{})
	if err != nil {
		fmt.Println("删除Secrets出现错误", err.Error())
		return errors.New("删除Secrets出现错误" + err.Error())
	}
	return nil
}

// UpdateSecrets 更新Secrets
func (s *secret) UpdateSecrets(namespace, content string) error {
	var sc = corev1.Secret{}
	err := json.Unmarshal([]byte(content), &sc)
	if err != nil {
		fmt.Println("序列化出现错误", err.Error())
		return errors.New("序列化出现错误" + err.Error())
	}
	_, err = K8s.ClientSet.CoreV1().Secrets(namespace).Update(context.TODO(), &sc, v1.UpdateOptions{})
	if err != nil {
		fmt.Println("更新Secrets出现错误", err.Error())
		return errors.New("更新Secrets出现错误" + err.Error())
	}
	return nil
}

// toCells方法用于将secret类型数组，转换成DataCell类型数组
func (s *secret) toCells(std []corev1.Secret) []DataCell {
	cells := make([]DataCell, len(std))
	for i := range std {
		cells[i] = secretCell(std[i])
	}
	return cells
}

// fromCells方法用于将DataCell类型数组，转换成sc类型数组
func (s *secret) fromCells(cells []DataCell) []corev1.Secret {
	sc := make([]corev1.Secret, len(cells))
	for i := range cells {
		sc[i] = corev1.Secret(cells[i].(secretCell))
	}
	return sc
}
