package services

import (
	"edu_api/models"
	_ "github.com/go-sql-driver/mysql"
)

/**
获取分类列表
 */
func (baseOrm *BaseOrm) CategoryList() (category []models.Category, err error) {

	var tmpCategory []models.Category

	if err := baseOrm.GetDB().Table("h_edu_categories").Find(&tmpCategory).Error; err != nil {
		return nil, err
	}

	//对当前分类进行无限极分类排序
	category = tree(tmpCategory)

	return category, nil
}

/**
预处理数据
 */
func tree(list []models.Category) []models.Category {
	data := formatData(list)

	result := formatCore(0, data)

	return result
}

/**
格式化数据
 */
func formatData(list []models.Category) map[int]map[int]models.Category {
	data := make(map[int]map[int]models.Category)

	for _, v := range list {
		id := int(v.Id)
		parentId := int(v.ParentId)

		if _, ok := data[parentId]; !ok {
			data[parentId] = make(map[int]models.Category)
		}

		data[parentId][id] = v
	}

	return data
}

/**
格式化数据核心方法
 */
func formatCore(index int, data map[int]map[int]models.Category) []models.Category {
	tmp := make([]models.Category, 0)

	for id, item := range data[index] {
		if data[id] != nil {
			item.Children = formatCore(id, data)
		}

		tmp = append(tmp, item)
	}

	return tmp
}
