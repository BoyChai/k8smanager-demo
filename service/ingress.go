package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	networkv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var Ingress ingress

type ingress struct {
}

// IngressCreate 定义ServiceCreate结构体，用于创建service需要的参数属性的定义
type IngressCreate struct {
	Name      string                 `json:"name"`
	Namespace string                 `json:"namespace"`
	Label     map[string]string      `json:"label"`
	Hosts     map[string][]*HttpPath `json:"hosts"`
}

// HttpPath 定义ingress的path结构体
type HttpPath struct {
	Path        string             `json:"path"`
	PathType    networkv1.PathType `json:"path_type"`
	ServiceName string             `json:"service_name"`
	ServicePort int32              `json:"service_port"`
}

// Getingresss 获取ingress列表
func (i *ingress) Getingresss(filterName, namespace string, limit, page int) (*dataSelector, error) {
	ingressList, err := K8s.ClientSet.NetworkingV1().Ingresses(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		fmt.Println("获取ingress列表出现错误", err.Error())
		return nil, errors.New("获取ingress列表出现错误" + err.Error())
	}
	dataselector := dataSelector{
		GenericDataList: i.toCells(ingressList.Items),
		dataSelectQuery: &DataSelectQuery{
			FilterQuery: &FilterQuery{Name: filterName},
			PaginateQuery: &PaginateQuery{
				Limit: limit,
				Page:  page,
			},
		},
	}
	// 筛选过滤
	filter := dataselector.Filter()
	// 排序过滤
	data := filter.Sort().Paginate()
	return data, nil
}

// GetIngressDetail 获取Ingress详情
func (i *ingress) GetIngressDetail(ingressName, namespace string) (*networkv1.Ingress, error) {
	ingressInfo, err := K8s.ClientSet.NetworkingV1().Ingresses(namespace).Get(context.TODO(), ingressName, metav1.GetOptions{})
	if err != nil {
		fmt.Println("获取Ingress信息出现错误", err.Error())
		return nil, errors.New("获取Ingress信息出现错误" + err.Error())
	}
	return ingressInfo, nil
}

// CreateIngress 创建Ingress
func (i *ingress) CreateIngress(data *IngressCreate) (err error) {
	//声明nwv1.IngressRule和nwv1.HTTPIngressPath变量，后面组装数据于鏊用到
	var ingressRules []networkv1.IngressRule
	var httpIngressPATHs []networkv1.HTTPIngressPath
	//将data中的数据组装成nwv1.Ingress对象
	in := &networkv1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:      data.Name,
			Namespace: data.Namespace,
			Labels:    data.Label,
		},
		Status: networkv1.IngressStatus{},
	}
	//第一层for循环是将host组装成nwv1.IngressRule类型的对象
	// 一个host对应一个ingressrule，每个ingressrule中包含一个host和多个path
	for key, value := range data.Hosts {
		ir := networkv1.IngressRule{
			Host: key,
			//这里现将nwv1.HTTPIngressRuleValue类型中的Paths置为空，后面组装好数据再赋值
			IngressRuleValue: networkv1.IngressRuleValue{
				HTTP: &networkv1.HTTPIngressRuleValue{Paths: nil},
			},
		}
		//第二层for循环是将path组装成nwv1.HTTPIngressPath类型的对象
		for _, httpPath := range value {
			hip := networkv1.HTTPIngressPath{
				Path:     httpPath.Path,
				PathType: &httpPath.PathType,
				Backend: networkv1.IngressBackend{
					Service: &networkv1.IngressServiceBackend{
						Name: httpPath.ServiceName,
						Port: networkv1.ServiceBackendPort{
							Number: httpPath.ServicePort,
						},
					},
				},
			}
			//将每个hip对象组装成数组
			httpIngressPATHs = append(httpIngressPATHs, hip)
		}
		//给Paths赋值，前面置为空了
		ir.IngressRuleValue.HTTP.Paths = httpIngressPATHs
		//将每个ir对象组装成数组，这个ir对象就是IngressRule，每个元素是一个host和多个path
		ingressRules = append(ingressRules, ir)
	}
	//将ingressRules对象加入到ingress的规则中
	in.Spec.Rules = ingressRules
	//创建ingress
	_, err = K8s.ClientSet.NetworkingV1().Ingresses(data.Namespace).Create(context.TODO(), in, metav1.CreateOptions{})
	if err != nil {
		fmt.Println("创建Ingress失败, ", err.Error())
		return errors.New("创建Ingress失败, " + err.Error())
	}
	return nil
}

// DeleteIngress 删除Ingress
func (i *ingress) DeleteIngress(ingressName, namespace string) error {
	err := K8s.ClientSet.NetworkingV1().Ingresses(namespace).Delete(context.TODO(), ingressName, metav1.DeleteOptions{})
	if err != nil {
		fmt.Println("删除ingress出错", err.Error())
		return errors.New("删除ingress出错" + err.Error())
	}
	return nil
}

// UpdateIngress 更新Ingress
func (i *ingress) UpdateIngress(namespace, content string) error {
	var update = &networkv1.Ingress{}
	err := json.Unmarshal([]byte(content), &update)
	if err != nil {
		fmt.Println("序列化失败", err.Error())
		return errors.New("序列化失败" + err.Error())
	}
	_, err = K8s.ClientSet.NetworkingV1().Ingresses(namespace).Update(context.TODO(), update, metav1.UpdateOptions{})
	if err != nil {
		fmt.Println("更新ingress失败", err.Error())
		return errors.New("更新ingress失败" + err.Error())
	}
	return nil
}

// toCells方法用于将ingress类型数组，转换成DataCell类型数组
func (i *ingress) toCells(std []networkv1.Ingress) []DataCell {
	cells := make([]DataCell, len(std))
	for ix := range std {
		cells[ix] = ingressCell(std[ix])
	}
	return cells
}

// fromCells方法用于将DataCell类型数组，转换成ingress类型数组
func (i *ingress) fromCells(cells []DataCell) []networkv1.Ingress {
	ingresss := make([]networkv1.Ingress, len(cells))
	for ix := range cells {
		ingresss[ix] = networkv1.Ingress(cells[ix].(ingressCell))
	}
	return ingresss
}
