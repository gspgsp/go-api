package models

/**
分类
*/
type Category struct {
	Id          int64      `json:"id"`
	ParentId    int64      `json:"parent_id"`
	Type        string     `json:"type"`
	Code        string     `json:"code"`
	Name        string     `json:"name"`
	Icon        string     `json:"icon"`
	Keywords    string     `json:"keywords"`
	Description string     `json:"description"`
	Sort        int        `json:"sort"`
	Status      int        `json:"status"`
	IsDirectory int        `json:"is_directory"`
	Children    []Category `json:"children"` //子类
}
