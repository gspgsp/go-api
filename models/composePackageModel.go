package models

type ComposePackageModel struct {
	Package
	PackageCourse []Course `json:"package_course"`
}

type ComposeModel struct {
	ComposePackage []ComposePackageModel `json:"compose_package"`
}
