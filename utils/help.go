package utils

import (
	"math"
	"strconv"
)

/**
课件类型转换
 */
func TransferMaterialType(m string) (m_type string) {
	switch m {
	case "courseware":
		return "课件"
	case "notes":
		return "笔记"
	case "software":
		return "软件"
	case "literature":
		return "文献"
	case "other":
		return "其它"
	default:
		return "未知类型"
	}
}

/**
课件大小转换
 */
func TransferMaterialSize(s float64) (s_size string) {

	if s >= 1073741824 {
		s_size = strconv.FormatFloat(math.Round(s/1073741824*100) / 100, 'f', -1, 64) +"GB"
	} else if s >= 1048576 {
		s_size = strconv.FormatFloat(math.Round(s/1048576*100) / 100, 'f', -1, 64) +"MB"
	} else if s >= 1024 {
		s_size = strconv.FormatFloat(math.Round(s/1024*100) / 100, 'f', -1, 64) + "KB"
	} else {
		s_size = strconv.FormatFloat(s, 'f', -1, 64) + "字节"
	}

	return
}