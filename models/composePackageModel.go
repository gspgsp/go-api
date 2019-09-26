package models

type ComposePackageModel struct {
	Package
	PackageCourse []Course `json:"package_course"`
}
