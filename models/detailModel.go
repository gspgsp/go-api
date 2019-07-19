package models

type Detail struct {
	Course//开始 我是用Course Course 对象的，但是数据库操作的find()赋值会有问题，所以还是直接继承Course struct的成员变量
	BuyId int `json:"buy_id"` //用来判断:用户是否观看过或者购买过
	Schedule float32 `json:"schedule"`//用户观看课程的进度，0~100.
}
