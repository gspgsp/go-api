package controllers

import "edu_api/services"

/**
公共控制器
 */
type Controller struct {
	Err error
	BaseOrm services.BaseOrm
}