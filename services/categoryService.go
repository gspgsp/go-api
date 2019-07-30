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
	res := Trees(tmpCategory)

	return res.([]models.Category), nil
}
