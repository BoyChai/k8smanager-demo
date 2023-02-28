package service

import (
	appv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"sort"
	"strings"
	"time"
)

// dataSelector 用于封装排序、过滤、分页的数据类型
type dataSelector struct {
	GenericDataList []DataCell
	dataSelectQuery *DataSelectQuery
}

// DataCell 接口,用于各种资源list的类型的转换，转换后可以使用dataSelector的自定义排序方法
type DataCell interface {
	GetCreation() time.Time
	GetName() string
}

// DataSelectQuery 定义过滤和分页的属性，过滤：Name，分页：Limit和Page
// limit是单页的数据条数
// page是第几页
type DataSelectQuery struct {
	FilterQuery   *FilterQuery
	PaginateQuery *PaginateQuery
}
type FilterQuery struct {
	Name string
}

type PaginateQuery struct {
	Limit int
	Page  int
}

// 排序

// Len 方法获取数组长度
func (d *dataSelector) Len() int {
	return len(d.GenericDataList)
}

func (d *dataSelector) Swap(i, j int) {
	d.GenericDataList[i], d.GenericDataList[j] = d.GenericDataList[j], d.GenericDataList[i]
}

func (d *dataSelector) Less(i, j int) bool {
	a := d.GenericDataList[i].GetCreation()
	b := d.GenericDataList[j].GetCreation()

	return b.Before(a)
}

// Sort 重写上面三个方法使用sort，Sort进行排序
func (d *dataSelector) Sort() *dataSelector {
	sort.Sort(d)
	return d
}

// 过滤

// Filter 方法用于过滤元素，比较元素的Name属性，若包含，再返回
func (d *dataSelector) Filter() *dataSelector {
	//若Name的传参为空，则返回所有元素
	if d.dataSelectQuery.FilterQuery.Name == "" {
		return d
	}
	//若Name的传参不为空，则返回元素名中包含Name的所有元素
	filteredList := []DataCell{}
	for _, value := range d.GenericDataList {
		matches := true
		objName := value.GetName()
		if !strings.Contains(objName, d.dataSelectQuery.FilterQuery.Name) {
			matches = false
			continue
		}
		if matches {
			filteredList = append(filteredList, value)
		}
	}
	d.GenericDataList = filteredList
	return d
}

// 分页

// Paginate 方法用于数组分页，根据Limit和Page的传参，返回数据
func (d *dataSelector) Paginate() *dataSelector {
	limit := d.dataSelectQuery.PaginateQuery.Limit
	page := d.dataSelectQuery.PaginateQuery.Page
	//验证参数合法，若参数不合法，则返回所有数据
	if limit <= 0 || page <= 0 {
		return d
	}
	//举例：25个元素的数组，limit是10，page是3，startIndex是20，endIndex是30（实际上endIndex是 25）
	startIndex := limit * (page - 1)
	endIndex := limit*page - 1
	//处理最后一页，这时候就把endIndex由30改为25了
	if len(d.GenericDataList) < endIndex {
		endIndex = len(d.GenericDataList)
	}
	d.GenericDataList = d.GenericDataList[startIndex:endIndex]
	return d
}

// 定义podCell类型，实现GetCreateion和GetName方法后，可进行类型转换
type podCell corev1.Pod

func (p podCell) GetCreation() time.Time {
	return p.CreationTimestamp.Time
}
func (p podCell) GetName() string {
	return p.Name
}

// 定义deploymentCell类型，实现GetCreateion和GetName方法后，可以进行类型转换
type deploymentCell appv1.Deployment

func (p deploymentCell) GetCreation() time.Time {
	return p.CreationTimestamp.Time
}
func (p deploymentCell) GetName() string {
	return p.Name
}

// 定义daemonSetCell类型，实现GetCreateion和GetName方法后，可以进行类型转换
type daemonSetCell appv1.DaemonSet

func (d daemonSetCell) GetCreation() time.Time {
	return d.CreationTimestamp.Time
}
func (d daemonSetCell) GetName() string {
	return d.Name
}

// 定义statefulSetCell类型，实现GetCreateion和GetName方法后，可以进行类型转换
type statefulSetCell appv1.StatefulSet

func (s statefulSetCell) GetCreation() time.Time {
	return s.CreationTimestamp.Time
}
func (s statefulSetCell) GetName() string {
	return s.Name
}

// 定义node类型，实现GetCreateion和GetName方法后，可以进行类型转换
type nodeCell corev1.Node

func (n nodeCell) GetCreation() time.Time {
	return n.CreationTimestamp.Time
}
func (n nodeCell) GetName() string {
	return n.Name
}
