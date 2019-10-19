package utils

import (
	"math"
	"strconv"
	"reflect"
	"errors"
	"bytes"
	"net"
	"strings"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"fmt"
	"github.com/jmcvetta/randutil"
	"time"
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
		s_size = strconv.FormatFloat(math.Round(s/1073741824*100)/100, 'f', -1, 64) + "GB"
	} else if s >= 1048576 {
		s_size = strconv.FormatFloat(math.Round(s/1048576*100)/100, 'f', -1, 64) + "MB"
	} else if s >= 1024 {
		s_size = strconv.FormatFloat(math.Round(s/1024*100)/100, 'f', -1, 64) + "KB"
	} else {
		s_size = strconv.FormatFloat(s, 'f', -1, 64) + "字节"
	}

	return
}

/**
查找数组、切片或者字典中是否存在某个值
 */
func Contain(obj interface{}, target interface{}) (bool, error) {

	targetValue := reflect.ValueOf(target)

	switch reflect.TypeOf(target).Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < targetValue.Len(); i++ {
			if targetValue.Index(i).Interface() == obj {
				return true, nil
			}
		}
	case reflect.Map:
		if targetValue.MapIndex(reflect.ValueOf(obj)).IsValid() {
			return true, nil
		}
	}

	return false, errors.New("不包含当前元素")
}

func ContactHashKey(args ...string) string {
	var buffer bytes.Buffer
	for _, val := range args {
		buffer.WriteString(val)
	}

	return buffer.String()
}

/**
获取内网ip
 */
func GetLocalIP() string {
	conn, _ := net.Dial("udp", "8.8.8.8:80")
	defer conn.Close()
	localAddr := conn.LocalAddr().String()
	idx := strings.LastIndex(localAddr, ":")
	return localAddr[0:idx]
}

/**
获取外网ip
 */
func GetPublicIp() string {
	resp, err := http.Get("http://myexternalip.com/raw")
	if err != nil {
		return ""
	}
	defer resp.Body.Close()
	content, _ := ioutil.ReadAll(resp.Body)
	return string(content)
}

func TaoBaoAPI(ip string) *IPInfo {
	url := "http://ip.taobao.com/service/getIpInfo.php?ip="
	url += ip

	resp, err := http.Get(url)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()

	out, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil
	}
	var result IPInfo
	if err := json.Unmarshal(out, &result); err != nil {
		return nil
	}

	return &result
}

/**
重新处理小数：保留n位以及是否四舍五入
 */
func RetainNumber(number float64) float64 {

	value := math.Trunc(number*1e1+0.5) * 1e-1

	pValue, _ := strconv.ParseFloat(fmt.Sprintf("%.1f", value), 64)

	return pValue
}

/**
生成订单号:YmdHis
 */
func GenerateOrderNo() string {
	tim := time.Now().Format("20060102150405")
	num, _ := randutil.IntRange(100000, 999999)
	return fmt.Sprintf("%s%s", tim, strconv.Itoa(num))
}
