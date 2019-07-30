package services

import (
	"edu_api/models"
	"github.com/ant0ine/go-json-rest/rest"
	"strconv"
)

/**
资料列表
 */
func (baseOrm *BaseOrm) GetMaterialList(r *rest.Request) (material []models.Material, err error) {

	id, err := strconv.Atoi(r.PathParam("id"))
	if err != nil {
		return material, err
	}

	if err := baseOrm.GetDB().Table("h_edu_materials").Where("course_id = ?", id).Select("id, title, size, type, download_num, course_id").Find(&material).Error; err != nil {
		return material, err
	}

	if len(material) == 0 {
		return material, err
	}

	for index, value := range material {
		material[index].Type = TransferMaterialType(value.Type)
	}

	return material, nil
}

/**
课件类型转换
 */
func TransferMaterialType(m string) (m_type string) {
	switch m {
	case "courseware":
		return "课件"
	case "notes":
		return "笔记"
	case "software":
		return "软件"
	case "literature":
		return "文献"
	case "other":
		return "其它"
	default:
		return "未知类型"
	}
}
