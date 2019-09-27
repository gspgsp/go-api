package models

/**
单个套餐
*/
type ComposePackageModel struct {
	Package
	PackageCourse []Course `json:"package_course"`
}

/**
多个套餐
*/
type ComposeModel struct {
	ComposePackage []ComposePackageModel `json:"compose_package"`
}
