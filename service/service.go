package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

var Service service

type service struct {
}

// ServiceCreate 定义Create结构体，用于创建service需要的参数属性的定义
type ServiceCreate struct {
	Name          string            `json:"name"`
	Namespace     string            `json:"namespace"`
	Type          string            `json:"type"`
	ContainerPort int32             `json:"container_port"`
	Port          int32             `json:"port"`
	NodePort      int32             `json:"node_port"`
	Label         map[string]string `json:"label"`
}

// CreateService 创建service,,接收ServiceCreate对象
func (s *service) CreateService(data *ServiceCreate) (err error) {
	//将data中的数据组装成corev1.Service对象
	svc := &corev1.Service{
		//ObjectMeta中定义资源名、命名空间以及标签
		ObjectMeta: metav1.ObjectMeta{
			Name:      data.Name,
			Namespace: data.Namespace,
			Labels:    data.Label,
		},
		//Spec中定义类型，端口，选择器
		Spec: corev1.ServiceSpec{
			Type: corev1.ServiceType(data.Type),
			Ports: []corev1.ServicePort{
				{
					Name:     "http",
					Port:     data.Port,
					Protocol: "TCP",
					TargetPort: intstr.IntOrString{
						Type:   0,
						IntVal: data.ContainerPort,
					},
				},
			},
			Selector: data.Label,
		},
	}
	//默认ClusterIP,这里是判断NodePort,添加配置
	if data.NodePort != 0 && data.Type == "NodePort" {
		svc.Spec.Ports[0].NodePort = data.NodePort
	}
	//创建Service
	_, err = K8s.ClientSet.CoreV1().Services(data.Namespace).Create(context.TODO(),
		svc, metav1.CreateOptions{})
	if err != nil {
		fmt.Println("创建Service失败, ", err.Error())
		return errors.New("创建Service失败, " + err.Error())
	}
	return nil
}

// GetServices 获取列表
func (s *service) GetServices(filterName, namespace string, limit, page int) ([]corev1.Service, error) {
	list, err := K8s.ClientSet.CoreV1().Services(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		fmt.Println("获取Service列表失败", err.Error())
		return nil, errors.New("获取Service列表失败" + err.Error())
	}
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

// GetServiceDetail 获取Service详情
func (s *service) GetServiceDetail(serviceName, namespace string) (*corev1.Service, error) {
	svc, err := K8s.ClientSet.CoreV1().Services(namespace).Get(context.TODO(), serviceName, metav1.GetOptions{})
	if err != nil {
		fmt.Println("获取svc信息失败", err.Error())
		return nil, errors.New("获取svc信息失败" + err.Error())
	}
	return svc, nil
}

// 删除Service

func (s *service) DeleteService(serviceName, namespace string) error {
	err := K8s.ClientSet.CoreV1().Services(namespace).Delete(context.TODO(), serviceName, metav1.DeleteOptions{})
	if err != nil {
		fmt.Println("删除svc失败", err.Error())
		return errors.New("删除svc失败" + err.Error())
	}
	return nil
}

// 更新Service

func (s *service) UpdateService(content, namespace string) error {
	var svc = &corev1.Service{}
	err := json.Unmarshal([]byte(content), &svc)
	if err != nil {
		fmt.Println("序列化失败", err.Error())
		return errors.New("序列化失败" + err.Error())
	}
	_, err = K8s.ClientSet.CoreV1().Services(namespace).Update(context.TODO(), svc, metav1.UpdateOptions{})
	if err != nil {
		fmt.Println("更新svc失败", err.Error())
		return errors.New("更新svc失败" + err.Error())
	}
	return nil
}

// toCells方法用于将service类型数组，转换成DataCell类型数组
func (s *service) toCells(std []corev1.Service) []DataCell {
	cells := make([]DataCell, len(std))
	for i := range std {
		cells[i] = serviceCell(std[i])
	}
	return cells
}

// fromCells方法用于将DataCell类型数组，转换成service类型数组
func (s *service) fromCells(cells []DataCell) []corev1.Service {
	services := make([]corev1.Service, len(cells))
	for i := range cells {
		services[i] = corev1.Service(cells[i].(serviceCell))
	}
	return services
}
