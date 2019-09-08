package models

import (
	"time"
	"fmt"
)

type JsonTime time.Time

const timeFormart = "2006-01-02 15:04:05"

/**
当需要转json的时候，将时间格式转换为本地时间格式
参考链接:https://www.jianshu.com/p/03003d5cbdbc
 */
func (this *JsonTime) MarshalJSON() ([]byte, error) {
	var stamp = fmt.Sprintf("\"%s\"", time.Time(*this).Format(timeFormart))
	return []byte(stamp), nil
}

/**
返回格式化时间字符串
 */
func (this *JsonTime) String() string {
	return fmt.Sprintf("%s", time.Time(*this).Format(timeFormart))
}

/**
当需要映射会数据库的时候就用下面的方法，实际上没啥用，因为数据库的时间类型是格式化以后的字符串，而不是[]byte
*/
func (this *JsonTime) UnmarshalJSON(data []byte) (err error) {
	now, err := time.ParseInLocation(`"`+timeFormart+`"`, string(data), time.Local)
	*this = JsonTime(now)
	return
}

/**
课程学习
 */
type CourseLearn struct {
	Id            int64    `json:"id"`
	Status        string   `json:"status"`
	StartAt       int64    `json:"start_at"`
	FinishAt      JsonTime `json:"finish_at,omitempty"`
	WatchDuration int64    `json:"watch_duration"`
	LessonLength  int64    `json:"lesson_length"`
	WatchNum      int64    `json:"watch_num"`
	CreatedAt     JsonTime `json:"created_at"`
	UpdatedAt     JsonTime `json:"updated_at"`
	UserId        int64    `json:"user_id"`
	CourseId      int64    `json:"course_id"`
	ChapterId     int64    `json:"chapter_id,omitempty"`
	UnitId        int64    `json:"unit_id,omitempty"`
	LessonId      int64    `json:"lesson_id"`
}
