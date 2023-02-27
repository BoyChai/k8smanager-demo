package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"strconv"
	"time"
)

var Deployment deployment

type deployment struct{}

// DeploymentsResp 定义列表的返回内容，Items是deployment元素列表，Total为deployment元素数量
type DeploymentsResp struct {
	Items []appsv1.Deployment `json:"items"`
	Total int                 `json:"total"`
}

// DeployCreate 自定创建deployment的结构体
type DeployCreate struct {
	Name          string            `json:"name"`
	Namespace     string            `json:"namespace"`
	Replicas      int32             `json:"replicas"`
	Image         string            `json:"image"`
	Label         map[string]string `json:"label"`
	Cpu           string            `json:"cpu"`
	Memory        string            `json:"memory"`
	ContainerPort int32             `json:"container_port"`
	HealthCheck   bool              `json:"health_check"`
	HealthPath    string            `json:"health_path"`
}

// GetDeployments 获取deployment列表，支持过滤、排序、分页
func (d *deployment) GetDeployments(filterName, namespace string, limit, page int) (*DeploymentsResp, error) {
	list, err := K8s.ClientSet.AppsV1().Deployments(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		fmt.Println("获取deployment列表失败", err.Error())
		return nil, errors.New("获取deployment列表失败" + err.Error())
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
	// 先过滤
	filtered := selectableData.Filter()
	total := len(filtered.GenericDataList)

	// 排序和分页
	data := filtered.Sort().Paginate()
	// 转换
	Deployments := d.fromCells(data.GenericDataList)
	return &DeploymentsResp{Items: Deployments, Total: total}, nil
}

// GetDeploymentDetail 获取deployment详情
func (d *deployment) GetDeploymentDetail(deploymentName, namespace string) (deployment *appsv1.Deployment, err error) {
	deploymentInfo, err := K8s.ClientSet.AppsV1().Deployments(namespace).Get(context.TODO(), deploymentName, metav1.GetOptions{})
	if err != nil {
		fmt.Println("获取deployment详情失败", err.Error())
		return nil, errors.New("获取deployment详情失败" + err.Error())
	}
	return deploymentInfo, nil
}

// CreateDeployment 创建Deployment
func (d *deployment) CreateDeployment(data *DeployCreate) (err error) {
	//将data中的数据组装成appsv1.Deployment对象
	deployment := &appsv1.Deployment{
		// ObjectMeta中定义资源名，命名空以及标签
		ObjectMeta: metav1.ObjectMeta{
			Name:      data.Name,
			Namespace: data.Namespace,
			Labels:    data.Label,
		},
		// Spec中定义副本数、选择器、以及pod属性
		Spec: appsv1.DeploymentSpec{
			Replicas: &data.Replicas,
			Selector: &metav1.LabelSelector{MatchLabels: data.Label},
			Template: corev1.PodTemplateSpec{
				// 定义pod名和标签
				ObjectMeta: metav1.ObjectMeta{
					Name:   data.Name,
					Labels: data.Label,
				},
				// 定义容器名、镜像和端口
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  data.Name,
							Image: data.Image,
							Ports: []corev1.ContainerPort{{
								Name:          "http",
								Protocol:      corev1.ProtocolTCP,
								ContainerPort: 80,
							}}},
					},
				},
			},
		},
		//Status定义资源的运行状态，这里由于是新建，传入空的appsv1.DeploymentStatus{}对象即可
		Status: appsv1.DeploymentStatus{},
	}
	//判断是否打开健康检查功能，若打开，则定义ReadinessProbe和LivenessProbe
	if data.HealthCheck {
		//设置第一个容器的ReadinessProbe，因为我们pod中只有一个容器，所以直接使用index 0即可
		//若pod中有多个容器，则这里需要使用for循环去定义了
		deployment.Spec.Template.Spec.Containers[0].ReadinessProbe = &corev1.Probe{
			ProbeHandler: corev1.ProbeHandler{
				HTTPGet: &corev1.HTTPGetAction{
					Path: data.HealthPath,
					Port: intstr.IntOrString{
						Type:   0,
						IntVal: data.ContainerPort,
					},
				},
			},
			//初始化等待事件
			InitialDelaySeconds: 5,
			//超时时间
			TimeoutSeconds: 5,
			//执行间隔
			PeriodSeconds: 5,
		}
		deployment.Spec.Template.Spec.Containers[0].LivenessProbe = &corev1.Probe{
			ProbeHandler: corev1.ProbeHandler{
				HTTPGet: &corev1.HTTPGetAction{
					Path: data.HealthPath,
					Port: intstr.IntOrString{
						Type:   0,
						IntVal: data.ContainerPort,
					},
				},
			},
			//初始化等待事件
			InitialDelaySeconds: 15,
			//超时时间
			TimeoutSeconds: 5,
			//执行间隔
			PeriodSeconds: 5,
		}
		// 定义容器的limit和request资源
		deployment.Spec.Template.Spec.Containers[0].Resources.Limits = map[corev1.ResourceName]resource.Quantity{
			corev1.ResourceCPU:    resource.MustParse(data.Cpu),
			corev1.ResourceMemory: resource.MustParse(data.Memory)}

		deployment.Spec.Template.Spec.Containers[0].Resources.Requests = map[corev1.ResourceName]resource.Quantity{
			corev1.ResourceCPU:    resource.MustParse(data.Cpu),
			corev1.ResourceMemory: resource.MustParse(data.Memory)}
	}
	// 调用sdk创建deployment
	_, err = K8s.ClientSet.AppsV1().Deployments(data.Namespace).Create(context.TODO(), deployment, metav1.CreateOptions{})
	// 错误处理
	if err != nil {
		fmt.Println("创建deployment出现错误:", err.Error())
		return errors.New("创建deployment出现错误:" + err.Error())
	}
	return nil
}

// ScaleDeployment 设置deployment的副本数
func (d *deployment) ScaleDeployment(deploymentName, namespace string, scaleNum int) (replicas int32, err error) {
	scale, err := K8s.ClientSet.AppsV1().Deployments(namespace).GetScale(context.TODO(), deploymentName, metav1.GetOptions{})
	if err != nil {
		fmt.Println("获取Deployment pod数量错误", err.Error())
		return 0, errors.New("获取Deployment pod数量" + err.Error())
	}
	// 设置副本数
	scale.Spec.Replicas = int32(scaleNum)
	newScale, err := K8s.ClientSet.AppsV1().Deployments(namespace).UpdateScale(context.TODO(), deploymentName, scale, metav1.UpdateOptions{})
	if err != nil {
		fmt.Println("设置Deployment pod数量错误", err.Error())
	}
	return newScale.Spec.Replicas, nil
}

// DeleteDeployment 删除deployment
func (d *deployment) DeleteDeployment(deploymentName, namespace string) (err error) {
	err = K8s.ClientSet.AppsV1().Deployments(namespace).Delete(context.TODO(), deploymentName, metav1.DeleteOptions{})
	if err != nil {
		fmt.Println("删除deployment错误:", err.Error())
		return errors.New("删除deployment错误:" + err.Error())
	}
	return nil
}

// RestartDeployment 重启Deployment
// 重启原理是更新个env，只要更新yaml都会使deployment重启，所以不仅仅可以更新env，标签啥的也都行
func (d *deployment) RestartDeployment(deploymentName, namespace string) (err error) {
	patchData := map[string]interface{}{
		"spec": map[string]interface{}{
			"template": map[string]interface{}{
				"spec": map[string]interface{}{
					"containers": []map[string]interface{}{
						{
							"name": deploymentName,
							"env": []map[string]string{
								{
									"name":  "RESTART_",
									"value": strconv.FormatInt(time.Now().Unix(), 10)},
							},
						},
					},
				},
			},
		},
	}
	// 序列化为字节，因为patch方法只能接收字节类型参数
	patchByte, err := json.Marshal(patchData)
	if err != nil {
		fmt.Println("json序列化失败", err.Error())
		return errors.New("json序列化失败" + err.Error())
	}
	// 调用patch方法更新deployment
	_, err = K8s.ClientSet.AppsV1().Deployments(namespace).Patch(context.TODO(), deploymentName, "application/strategic-merge-patch+json", patchByte, metav1.PatchOptions{})
	if err != nil {
		fmt.Println("重启deployment失败", err.Error())
		return errors.New("重启deployment失败" + err.Error())
	}
	return nil
}

// UpdateDeployment 更新deployment
func (d *deployment) UpdateDeployment(namespace, content string) (err error) {
	var deploy = &appsv1.Deployment{}
	err = json.Unmarshal([]byte(content), deploy)
	if err != nil {
		fmt.Println("更新deployment反序列化失败", err.Error())
		return errors.New("更新deployment反序列化失败" + err.Error())
	}
	_, err = K8s.ClientSet.AppsV1().Deployments(namespace).Update(context.TODO(), deploy, metav1.UpdateOptions{})
	if err != nil {
		fmt.Println("更新deployment失败", err.Error())
		return errors.New("更新deployment失败" + err.Error())
	}
	return nil
}

// DeploysNp 定义DeploysNp类型，用于返回namespace中deployment的数量
type DeploysNp struct {
	Namespace string `json:"namespace"`
	DeployNum int    `json:"deployment_num"`
}

// GetDeployNumPerNp 获取每个namespace的deployment数量
func (d *deployment) GetDeployNumPerNp() (deploysNps []*DeploysNp, err error) {
	list, err := K8s.ClientSet.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		fmt.Println("获取命名空间错误", err.Error())
		return nil, errors.New("获取命名空间错误" + err.Error())
	}
	for _, item := range list.Items {
		deploymentList, err := K8s.ClientSet.AppsV1().Deployments(item.Name).List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			fmt.Println("获取deployment数量错误", err.Error())
			return nil, errors.New("获取deployment数量错误" + err.Error())
		}
		deploysNp := &DeploysNp{
			Namespace: item.Name,
			DeployNum: len(deploymentList.Items),
		}
		deploysNps = append(deploysNps, deploysNp)
	}
	return deploysNps, nil
}

// toCells方法用于将pod类型数组，转换成DataCell类型数组
func (p *deployment) toCells(std []appsv1.Deployment) []DataCell {
	cells := make([]DataCell, len(std))
	for i := range std {
		cells[i] = deploymentCell(std[i])
	}
	return cells
}

// fromCells方法用于将DataCell类型数组，转换成pod类型数组
func (p *deployment) fromCells(cells []DataCell) []appsv1.Deployment {
	pods := make([]appsv1.Deployment, len(cells))
	for i := range cells {
		pods[i] = appsv1.Deployment(cells[i].(deploymentCell))
	}
	return pods
}
