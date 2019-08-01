package services

import (
	"edu_api/models"
	"github.com/ant0ine/go-json-rest/rest"
	"strconv"
	"edu_api/utils"
	"net/http"
	jwt2 "github.com/dgrijalva/jwt-go"
	"errors"
)

/**
资料列表
 */
func (baseOrm *BaseOrm) GetMaterialList(r *rest.Request) (material []models.Material, err error) {

	var (
		defaultLimit  = 20
		defaultOffset = 0
	)

	id, err := strconv.Atoi(r.PathParam("id"))
	if err != nil {
		return material, err
	}

	params := r.URL.Query()
	limit := params.Get("limit")
	intLimit, _ := strconv.Atoi(limit)

	page := params.Get("page")
	intPage, _ := strconv.Atoi(page)

	//如果传了limit那么就限制取值数量,如果传了page那么就分页查询,么次必须只能穿一个
	if intLimit > 0 {
		defaultLimit = intLimit
		defaultOffset = 0
	} else if intPage > 0 {
		if intPage > 1 {
			defaultOffset = (intPage - 1) * defaultLimit
		} else {
			defaultOffset = 0
		}
	} else {
		return material, errors.New("limit/page 参数必须")
	}

	//查看当前课程类型
	var types []string
	baseOrm.GetDB().Table("h_edu_courses").Select("type").Where("id = ?", id).Pluck("type", &types)

	if types == nil {
		return material, errors.New("当前课程类型不存在")
	}

	courseType := types[0]

	if courseType != "free" {
		var (
			header      http.Header
			accessToken string
			j           models.JwtClaim
			userId      float64
		)

		header = r.Header
		if _, ok := header["Authorization"]; ok {
			for _, v := range header["Authorization"] {
				accessToken = v
			}

			token, err := j.VerifyToken(accessToken)

			if err != nil {
				return material, err
			}

			switch value := token.Claims.(jwt2.MapClaims)["id"].(type) {
			case float64:
				userId = value
			}

			//当前课程是否会员免费
			var levels []string
			baseOrm.GetDB().Table("h_users").Where("id = ?", userId).Pluck("level", &levels)

			//是否会员免费
			var vip_levels []int
			baseOrm.GetDB().Table("h_edu_courses").Where("id = ?", id).Pluck("vip_level", &vip_levels)

			//查看用户是否购买过当前课程
			var ids []int
			baseOrm.GetDB().Table("h_user_course").Where("user_id = ? and course_id = ?", userId, id).Pluck("id", &ids)

			if ((levels[0] == "vip1" || levels[0] == "vip2") && vip_levels[0] == 1) || ids[0] > 0 {
				//直接查询
				if err := baseOrm.GetDB().Table("h_edu_materials").Where("course_id = ?", id).Select("id, title, size, type, download_num, course_id").Limit(defaultLimit).Offset(defaultOffset).Find(&material).Error; err != nil {
					return material, err
				}
			} else {
				return material, errors.New("未购买当前课程")
			}
		} else {
			return material, errors.New(utils.TOKEN_PARAM_REQUIRE)
		}
	} else {
		//直接查询
		if err := baseOrm.GetDB().Table("h_edu_materials").Where("course_id = ?", id).Select("id, title, size, type, download_num, course_id").Limit(defaultLimit).Offset(defaultOffset).Find(&material).Error; err != nil {
			return material, err
		}
	}

	if len(material) == 0 {
		return material, err
	}

	for index, value := range material {
		material[index].Type = utils.TransferMaterialType(value.Type)
		material[index].FormatSize = utils.TransferMaterialSize(value.Size)
		material[index].Size = 0 //直接置空
	}

	return material, nil
}
