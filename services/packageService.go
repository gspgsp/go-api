package services

import (
	"edu_api/models"
	"errors"
	"github.com/ant0ine/go-json-rest/rest"
	log "github.com/sirupsen/logrus"
	"strconv"
	"sync"
)

/**
全局变量
*/
var (
	packageCourse []models.PackageCourseModel
	compose       models.ComposeModel
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

/**
组合套餐
*/
func (baseOrm *BaseOrm) GetComposePackage(r *rest.Request) (models.ComposeModel, error) {
	//因为gorm好像还不支持模型里面套模型赋值，所以这里不用join了，直接分开获取，再组合成指定模型
	var (
		pcs []models.PackageCourseModel
	)

	id, err := strconv.Atoi(r.PathParam("id"))
	if err != nil {
		log.Info("路由参数错误!")
		return compose, errors.New("路由参数错误!")
	}

	if err := baseOrm.GetDB().Table("h_edu_package_course").Where("course_id = ?", id).Select("package_id").Find(&pcs).Error; err != nil {
		log.Info("资源获取错误!" + err.Error())
		return compose, errors.New("资源获取错误!" + err.Error())
	}

	//声明一个工作池
	var wg sync.WaitGroup
	for i := 0; i < len(pcs); i++ {
		wg.Add(1)
		go getPackageCourse(baseOrm, pcs[i].PackageId, &wg)
	}
	wg.Wait()

	return compose, nil
}

/**
通过当前package获取所有对应的课程
*/
func getPackageCourse(baseOrm *BaseOrm, package_id int64, wg *sync.WaitGroup) {
	var (
		pcs            []models.PackageCourseModel
		courseIds      []int64
		courses        []models.Course
		composePackage models.ComposePackageModel
	)

	defer func() {
		wg.Done()
	}()

	if err := baseOrm.GetDB().Table("h_edu_packages").Where("id = ?", package_id).First(&composePackage).Error; err != nil {
		log.Info("资源获取错误!" + err.Error())
		return
	}

	if err := baseOrm.GetDB().
		Table("h_edu_package_course").
		Select("course_id").
		Find(&pcs).Error; err != nil {
		log.Info("资源获取错误!" + err.Error())
		return
	}

	for _, value := range pcs {
		courseIds = append(courseIds, value.CourseId)
	}

	if err := baseOrm.GetDB().
		Table("h_edu_courses").
		Where("status = 2 and id in (?)", courseIds).
		Select("id, type, title, subtitle, price, vip_price, discount, discount_end_at, cover_picture, back_picture, learn_num, buy_num, video_url").
		Find(&courses).Error; err != nil {
		log.Info("资源获取错误!" + err.Error())
		return
	}

	composePackage.PackageCourse = courses
	compose.ComposePackage = append(compose.ComposePackage, composePackage)

	return
}
