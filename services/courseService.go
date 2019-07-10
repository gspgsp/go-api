package services

import (
	"edu_api/models"
	"github.com/ant0ine/go-json-rest/rest"
	"edu_api/utils"
	"log"
	"strconv"
)

var (
	where = make(map[string]interface{})
)

/**
获取课程列表信息
 */
func (baseOrm *BaseOrm) CourseList(r *rest.Request) (course []models.Course, err error) {

	var (
		tmpCourse     []models.Course //这个变量可以不用声明,直接用course就可以了
		filter        utils.Filter
		order         = ""
		defaultLimit  = 20
		defaultOffset = 0
		not           = make(map[string]interface{})
	)

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
		log.Println("limit/page param require!")
		return
	}

	err = r.DecodeJsonPayload(&filter)

	if err != nil {
		log.Println("parse request param error!", err)
		return
	}

	if filter.Code != "" {
		var category models.Category

		if err = baseOrm.GetDB().Table("h_edu_categories").Where("code = ? ", filter.Code).First(&category).Error; err != nil {
			return nil, err
		}

		where["category_id"] = category.Id
	}

	if filter.DifficultyLevel != "" {
		where["difficulty_level"] = filter.DifficultyLevel
	}

	//排序
	if filter.Order != "" {
		switch filter.Order {
		case "new":
			order = "created_at desc"
		case "rating":
			order = "rating desc"
		case "sold_count":
			order = "(buy_num + learn_num) desc"
		}
	}

	if filter.VipPrice != "" {
		where["vip_level"] = filter.VipPrice
	}

	//当前课程的类别
	if filter.Type != "" {
		where["type"] = filter.Type
	} else {
		not["type"] = "class"
	}

	//是否推荐课程
	if filter.IsRecommended != "" {
		where["is_recommended"] = filter.IsRecommended
	}

	//必须是发布的课程
	where["status"] = "published"

	if err = baseOrm.GetDB().Table("h_edu_courses").Not(not).Where(where).Order(order).Limit(defaultLimit).Offset(defaultOffset).Find(&tmpCourse).Error; err != nil {
		return nil, err
	}

	return tmpCourse, nil
}

/**
获取课程详情信息
 */
func (baseOrm *BaseOrm) GetCourseDetail(r *rest.Request) (detail models.Detail, err error) {

	params := r.URL.Query()

	courseType := params.Get("type")

	id, err := strconv.Atoi(r.PathParam("id"))

	if courseType == "" || err != nil {
		return detail, err
	}

	where["id"] = id
	where["type"] = courseType

	var detailTemp models.Detail

	if err = baseOrm.GetDB().Table("h_edu_courses").Where(where).Find(&detailTemp.Course).Error; err != nil {
		return detail, err
	}

	log.Printf("the detail is:%v\n", detailTemp)

	return detailTemp, nil
}
