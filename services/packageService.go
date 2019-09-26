package services

import (
	"edu_api/models"
	"errors"
	"fmt"
	"github.com/ant0ine/go-json-rest/rest"
	log "github.com/sirupsen/logrus"
	"strconv"
)

/**
获取套餐列表
*/
func (baseOrm *BaseOrm) PackageList(r *rest.Request) (packages []models.Package, err error) {

	var (
		defaultLimit  = 20
		defaultOffset = 0
		where         = make(map[string]interface{})
		order         = "created_at desc"
	)

	params := r.URL.Query()
	packageType := params.Get("type")

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
		log.Println("limit/page param require!")
		return nil, errors.New("limit/page param require!")
	}

	//套餐类型
	if packageType != "" {
		where["type"] = packageType
	}

	//必须是发布的课程
	where["status"] = "published"

	if err = baseOrm.GetDB().Table("h_edu_packages").Where(where).Order(order).Limit(defaultLimit).Offset(defaultOffset).Find(&packages).Error; err != nil {
		return nil, err
	}

	return packages, nil
}

func (baseOrm BaseOrm) GetComposePackage(r *rest.Request) ([]models.ComposePackageModel, error) {
	//因为gorm好像还不支持模型里面套模型赋值，所以这里不用join了，直接分开获取，再组合成指定模型
	var (
		pcs []models.PackageCourseModel
	)

	id, err := strconv.Atoi(r.PathParam("id"))
	if err != nil {
		log.Info("路由参数错误!")
		return nil, errors.New("路由参数错误!")
	}

	if err := baseOrm.GetDB().Table("h_edu_package_course").Where("course_id = ?", id).Select("package_id").Find(&pcs).Error; err != nil {
		log.Info("资源获取错误!" + err.Error())
		return nil, errors.New("资源获取错误!" + err.Error())
	}

	fmt.Printf("the len is:%v\n", len(pcs))

	return nil, errors.New("未知")
}
