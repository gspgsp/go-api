package services

import (
	"edu_api/models"
	"edu_api/utils"
	"errors"
	"github.com/ant0ine/go-json-rest/rest"
	jwt2 "github.com/dgrijalva/jwt-go"
	log "github.com/sirupsen/logrus"
	"net/http"
	"reflect"
	"strconv"
	"strings"
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

/**
获取课程章节
*/
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
	where["status"] = "published"

	if err := baseOrm.GetDB().Table("h_edu_chapters").Where(where).Find(&tmpChapter).Error; err != nil {
		return nil, err
	}

	//对当前分类进行无限极分类排序
	res := Trees(tmpChapter)

	return res.([]models.Chapter), nil
}

/**
评价列表
*/
func (baseOrm *BaseOrm) GetCourseReview(r *rest.Request) (reviews []models.Review, err error) {

	var (
		defaultLimit  = 20
		defaultOffset = 0
	)

	id, err := strconv.Atoi(r.PathParam("id"))
	if err != nil {
		return
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
		return reviews, errors.New("limit/page 参数必须")
	}

	if err := baseOrm.GetDB().
		Table("h_user_course").
		Where("course_id = ?", id).
		Select("id, anonymous, rating, practical_rating, popular_rating, logic_rating, status, review, reply, reviewed_at, reply_at, course_id").
		Limit(defaultLimit).
		Offset(defaultOffset).
		Find(&reviews).
		Error; err != nil {
		return reviews, err
	}

	return
}

/**
获取推荐课
*/
func (baseOrm *BaseOrm) GetRecommendCourse(r *rest.Request) (recommends []models.Recommend, err error) {
	var (
		defaultLimit = 4
	)

	id, err := strconv.Atoi(r.PathParam("id"))
	if err != nil {
		return
	}

	params := r.URL.Query()
	limit := params.Get("limit")
	intLimit, _ := strconv.Atoi(limit)

	//如果传了limit那么就限制取值数量
	if intLimit > 0 {
		defaultLimit = intLimit
	}

	//利用find方法，必须要有模型存储返回数据
	type tag struct {
		Id int `json:"id"`
	}

	var tags []tag

	if err := baseOrm.GetDB().Raw("select t.id from h_tags t where exists (select 1 from h_taggables t_g inner join h_edu_courses c on c.id = t_g.taggable_id where t_g.tag_id = t.id and c.id = ? and c.status = 'published')", id).Find(&tags).Error; err != nil {
		return nil, err
	}

	channel := make(chan []models.Recommend, 4)
	endChannel := make(chan bool)
	endNumber := 0
	number := 0

	if len(tags) == 0 {
		//当前课程如果没有标签，就从具有推荐属性的课程获取(不区分类型，包括免费和精品)
		if err := baseOrm.GetDB().Table("h_edu_courses").Where("type in ('free', 'boutique') and status = 'published' and is_recommended = 1").Limit(defaultLimit).Find(&recommends).Error; err != nil {
			return nil, err
		}

		return
	}

	number = len(tags)

	//通过goroutine获取数据
	for _, val := range tags {
		go getRecommend(channel, endChannel, val.Id, id, baseOrm)
	}

	//开始准备用Rows方法获取数据的，但是后来没法获取数据的长度，来给number赋值，所以还是用Find()方法
	//必须调用一次Next，不能直接调用Scan取值
	//for tagIds.Next() {
	//	tagIds.Scan(&temp)
	//
	//	//对当前temp进行处理，准备通过goroutine获取对应tag下的课程
	//	go test(channel, endChannel, temp.(int64))
	//
	//}

GetChannelData:
	for {
		select {
		case v, ok := <-channel:
			if ok {

				for _, val := range v {
					recommends = append(recommends, val)
				}
			} else {
				log.Printf("read channel error")
			}
		case <-endChannel:

			endNumber++

			if endNumber == number {
				close(channel)
				break GetChannelData
			}
		}
	}

	temp := RemoveDuplicateSlice(recommends)

	if len(temp) >= defaultLimit {
		return temp[0:defaultLimit], nil
	} else {
		//补充带标签的推荐课程
		external := make([]models.Recommend, defaultLimit)
		exceptIds := []int{id}

		for index, _ := range temp {
			exceptIds = append(exceptIds, temp[index].Id)
		}

		if err := baseOrm.GetDB().
			Table("h_edu_courses").
			Where("type in ('free', 'boutique') and status = 'published' and is_recommended = 1 and id not in (?) ", exceptIds).
			Limit(defaultLimit - len(temp)).Select("id, type, title, price, vip_price, discount, discount_end_at, cover_picture, learn_num, buy_num").
			Find(&external).
			Error; err != nil {
			//没有三目运算，只能if else
			if len(temp) > 0 {
				return temp, err
			} else {
				return nil, err
			}
		} else {
			for _, val := range external {
				temp = append(temp, val)
			}
			return temp, nil
		}
	}
}

/**
获取相同标签下的推荐课程
*/
func getRecommend(channel chan []models.Recommend, endChannel chan bool, tagId int, id int, baseOrm *BaseOrm) {

	defer func() {
		endChannel <- true
	}()

	var recommend []models.Recommend
	//这个子查询in(少) 和 exists(多)效率差不多 还是join查询快一点(少了一层子结果集扫描)
	/*
		select * from h_edu_courses where id in (select taggable_id from h_taggables where tag_id = 7 and taggable_id != 106);
		select * from h_edu_courses where exists (select tag_id from h_taggables where tag_id = 7 and taggable_id != 106 and h_taggables.taggable_id = h_edu_courses.id);
		select * from h_edu_courses inner join h_taggables on h_edu_courses.id = h_taggables.taggable_id where h_taggables.tag_id = 7 and taggable_id != 106;
	*/
	baseOrm.GetDB().
		Raw("select id, type, title, price, vip_price, discount, discount_end_at, cover_picture, learn_num, buy_num from h_edu_courses left join h_taggables on h_edu_courses.id = h_taggables.taggable_id where h_taggables.tag_id = ? and taggable_id != ?;", tagId, id).
		Find(&recommend)

	//简单的读取操作
	channel <- recommend
}

/**
slice 去重操作(类似冒泡排序), 参数类型不能用interface{},否则返回值没法处理成指定类型
*/
func RemoveDuplicateSlice(a []models.Recommend) (ret []models.Recommend) {

	n := len(a)

	for i := 0; i < n; i++ {

		state := false

		for j := i + 1; j < n; j++ {
			if j > 0 && reflect.DeepEqual(a[i], a[j]) {
				state = true
				continue
			}
		}

		if !state {
			ret = append(ret, a[i])
		}
	}

	return
}

func (baseOrm *BaseOrm) GetTrySeeList(r *rest.Request) (ret []models.Chapter, err error) {
	//因为试看有两个地方，所以通过这种方式传递参数
	var (
		courseId  string //当为体系课的时候，这个id为ChapterId；当为精品课的时候，这个为CourseId
		isChapter string //是否为体系课 1：是 0：否
	)

	idStr := r.PathParam("id")
	if len(idStr) == 0 {
		log.Info("参数错误!")
		return nil, errors.New("参数错误!")
	}

	idArr := strings.Split(idStr, "_")
	if len(idArr) == 0 {
		log.Info("参数错误!")
		return nil, errors.New("参数错误!")
	}

	courseId = idArr[0]
	isChapter = idArr[1]

	if isChapter == "0" {
		if err = baseOrm.GetDB().Table("h_edu_chapters").Where("course_id = ? and is_free = 1 and status = 2", courseId).Find(&ret).Error; err != nil {
			return nil, err
		}

		return ret, nil
	} else if isChapter == "1" {

	}

	return
}

//递归获取所有子级id
func getRecursion(baseOrm *BaseOrm, courseId string) {
	type ParentId struct {
		Id string `json:"id"`
	}

	var (
		done = false
		ids  []ParentId
	)
	//db.Where("name IN (?)", []string{"jinzhu", "jinzhu 2"}).Find(&users)
	for {
		var temp = ""
		baseOrm.GetDB().Table("h_edu_chapters").Where("parent_id in (?)", []string{courseId}).Find(&ids)

		for _, value := range ids {
			temp += `"` + value.Id + `",`
		}

		if len(ids) == 0 {
			done = true
		} else {
			courseId = temp
		}

		if done {
			break
		}
	}
}
