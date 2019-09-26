package models

/**
套餐课程表
*/
type PackageCourseModel struct {
	PackageId int64 `json:"package_id,omitempty"`
	CourseId  int64 `json:"course_id,omitempty"`
	Sort      int64 `json:"sort,omitempty"`
}
