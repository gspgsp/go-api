package models

/**
成绩表模型
*/
type GradeModel struct {
	Id        int64  `json:"id"`
	Point     int64  `json:"point,omitempty"`
	Result    string `json:"result,omitempty"` //成绩结果
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	RollId    int64  `json:"'roll_id'"`
	CourseId  int64  `json:"course_id"`
	ChapterId int64  `json:"chapter_id,omitempty"`
	UserId    int    `json:"user_id"`
}

/**
成绩表结果
*/
type GradeResult struct {
	Point    int64 `json:"point"`
	Numbers  int   `json:"numbers"`
	Success  int   `json:"success"`
	UseTimes int64 `json:"use_times"`
	AllPoint int64 `json:"all_point"`
}

/**
成绩详情表模型
*/
type GradeLogModel struct {
	Id        int64  `json:"id"`
	IsCorrect int    `json:"is_correct"`
	Result    string `json:"result,omitempty"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	GradeId   int64  `json:"grade_id"`
	RollId    int64  `json:"roll_id"`
	CourseId  int64  `json:"course_id"`
	TopicId   int64  `json:"topic_id"`
	UserId    int64  `json:"user_id"`
}

/**
成绩详情结果
*/
type GradeLogResult struct {
	TopicId   int64  `json:"topic_id,omitempty"`
	Num       string `json:"num"`
	UserChose string `json:"user_chose"`
}

/**
答题返回结果
*/
type AnswerReturn struct {
	GradeResult
	Result     string `json:"result"`
	SubmitTime string `json:"submit_time"`
	UseTime    string `json:"use_time"`
}
