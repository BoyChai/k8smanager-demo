package dao

import (
	"errors"
	"fmt"
	"k8smanager-demo/db"
	"k8smanager-demo/model"
)

var Workflow workflow

type workflow struct {
}
type workflowResp struct {
	Item  []*model.Workflow
	Total int
}

// GetList 获取workflow列表
func (w *workflow) GetList(filterName string, limit, page int) (data *workflowResp, err error) {
	// 定义分页起始位置
	startSet := (page - 1) * limit
	// 定义数据库查询返回的内容
	var workflowList []*model.Workflow
	// 数据库查询，limit方法用于限制条数，Offset方法用于设置起始位置
	tx := db.GORM.Where("name like ?", "%"+filterName+"%").
		Limit(limit).
		Offset(startSet).
		Order("id desc").
		Find(&workflowList)
	if tx.Error != nil && tx.Error.Error() != "record not found" {
		fmt.Println("获取workflow列表失败", tx.Error.Error())
		return nil, errors.New("获取workflow列表失败" + tx.Error.Error())
	}

	return &workflowResp{
		Item:  workflowList,
		Total: len(workflowList),
	}, nil
}

// GetById 查询workflow单条数据
func (w *workflow) GetById(id int) (workflow *model.Workflow, err error) {
	workflow = &model.Workflow{}
	tx := db.GORM.Where("id = ?", id).First(&workflow)
	if tx.Error != nil && tx.Error.Error() != "record not found" {
		fmt.Println("获取Workflow单条数据失败, ", tx.Error.Error())
		return nil, errors.New("获取Workflow单条数据失败, " + tx.Error.Error())
	}
	return
}

// Add 新增workflow
func (w *workflow) Add(workflow *model.Workflow) (err error) {
	tx := db.GORM.Create(&workflow)
	if tx.Error != nil {
		fmt.Println("添加Workflow失败, ", tx.Error.Error())
		return errors.New("添加Workflow失败, " + tx.Error.Error())
	}
	return nil
}

// DelById 删除workflow
// 软删除 db.GORM.Delete("id = ?", id)
// 软删除执行的是UPDATE语句，将deleted_at字段设置为时间即可，gorm 默认就是软删。
// 实际执行语句 UPDATE `workflow` SET `deleted_at` = '2021-03-01 08:32:11' WHERE `id` IN ('1'
// 硬删除 db.GORM.Unscoped().Delete("id = ?", id)) 直接从表中删除这条数据
// 实际执行语句 DELETE FROM `workflow` WHERE `id` IN ('1');
func (w *workflow) DelById(id int) (err error) {
	tx := db.GORM.Where("id = ?", id).Delete(&model.Workflow{})
	if tx.Error != nil {
		fmt.Println("删除Workflow失败, ", tx.Error.Error())
		return errors.New("删除Workflow失败, " + tx.Error.Error())
	}
	return nil
}
