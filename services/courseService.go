package services

import (
	"edu_api/models"
	"github.com/ant0ine/go-json-rest/rest"
	"edu_api/utils"
	"log"
	"strconv"
	"net/http"
	jwt2 "github.com/dgrijalva/jwt-go"
)

type Result struct {
	UID      int    `json:"uid"`
	Id       int    `json:"id"`
	Title    string `json:"title"`
	CourseId int    `json:"course_id"`
}

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
		where         = make(map[string]interface{})
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

	var (
		where       = make(map[string]interface{})
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
			return detail, err
		}

		//这里必须要转一下类型,不然取不出来,自定义的id,这里一开始不知道token.Claims.(jwt2.MapClaims)["id"]返回值类型，所以token.Claims.(jwt2.MapClaims)["id"].(int),测了一下，结果报错了，说是不能将float64转化为int,就知道了是float64类型
		//也可以用断言，但是断言有个问题，就是case得到的结果，value赋值的问题，比如userId 我声明的是float64，要是有多个case的话，就需要多种类型的变量去接，
		//不能 case string:userId = value case int:userId = value
		switch value := token.Claims.(jwt2.MapClaims)["id"].(type) {
		case float64:
			userId = value
		}
	}

	params := r.URL.Query()
	courseType := params.Get("type")

	id, err := strconv.Atoi(r.PathParam("id"))

	if courseType == "" || err != nil {
		return detail, err
	}

	where["id"] = id
	where["type"] = courseType
	where["status"] = "published"

	//利用go 的ORM查找的话，限制条件太多了，相比Laravel的ORM操作，要麻烦很多，所以推荐用类似于Laravel的查询构造器进行数据的查询，但是有个问题就是怎么给表取别名类似as操作，go里面好像不行，还有就是设置默认表前缀是不能在这里用的
	//baseOrm.GetDB().Model(&userCourse).Where(tempWhere).Related(&userCourse.Course).Find(&userCourse)

	if userId > 0 {
		baseOrm.GetDB().
			Table("h_edu_courses").
			Joins("left join h_user_course on h_user_course.course_id = h_edu_courses.id").
			Where(where).
			Where("h_user_course.user_id = ?", userId).
			Select("h_edu_courses.*, h_user_course.id as buy_id, h_user_course.schedule").
			Find(&detail)
	} else {
		baseOrm.GetDB().
			Table("h_edu_courses").
			Where(where).
			Find(&detail)
	}

	return detail, nil
}

func (baseOrm *BaseOrm) GetCourseChapter(r *rest.Request) (chapters []models.Chapter, err error) {
	var (
		tmpChapter []models.Chapter
		where      = make(map[string]interface{})
	)

	id, err := strconv.Atoi(r.PathParam("id"))
	if err != nil {
		return tmpChapter, err
	}

	where["course_id"] = id

	if err := baseOrm.GetDB().Table("h_edu_chapters").Where(where).Find(&tmpChapter).Error; err != nil {
		return nil, err
	}

	//对当前分类进行无限极分类排序

	return tmpChapter, nil
}
