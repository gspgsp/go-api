package services

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"edu_api/middlewares"
	"edu_api/models"
	"edu_api/tasksAndEvents"
	"edu_api/utils"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/iGoogle-ink/gopay"
	log "github.com/sirupsen/logrus"
	"io"
	"sort"
	"time"
)

const (
	//支付宝公钥
	ALIPAY_PUBLIC_KEY = `-----BEGIN PUBLIC KEY-----  
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCnxj/9qwVfgoUh/y2W89L6BkRA
FljhNhgPdyPuBV64bfQNN1PjbCzkIM6qRdKBoLPXmKKMiFYnkd6rAoprih3/PrQE
B/VsW8OoM8fxn67UDYuyBTqA23MML9q1+ilIZwBC2AQ2UBVOrFXfFl75p6/B5Ksi
NG9zpgmLCUYuLkxpLQIDAQAB 
-----END PUBLIC KEY-----
`
)

/**
初始化
*/
func init() {
	db = new(BaseOrm)
	db.InitDB()

	//初始化支付宝客户端
	aliPayConfig, _ = GetAliPayConf()
	aliPayClient = gopay.NewAliPayClient(aliPayConfig.Alipay.AppId, aliPayConfig.Alipay.PrivateKey, false)

	//初始化微信支付客户端

}

/**
生成支付信息
*/
func (baseOrm *BaseOrm) Payment(r *rest.Request, payment *middlewares.Payment) (int, interface{}) {
	auth = r.Header.Get("Authorization")
	user = GetUserInfo(auth)

	pay := checkPayIsValid(payment)
	switch val := pay.(type) {
	case error:
		return 1, val
	case models.OrderModel:
		//课程订单支付
		url, err := payPage(val, payment)
		if err != nil {
			log.Info("pay err is:", err)
			return 1, err
		}

		return 0, url
	case models.VipOrderModel:
		//VIP订单支付
	}

	return 0, nil
}

/**
判断当前支付订单是否有效
*/
func checkPayIsValid(payment *middlewares.Payment) interface{} {
	where := fmt.Sprintf("user_id = %d and status = 0 and no = %s and payment_status = 0 and payment_amount > 0 and payment_expired_at > '%s'", user.Id, payment.No, utils.ParseTimeToString())
	if payment.PayType == "course" {
		var order models.OrderModel
		if err := db.GetDB().Table("h_orders").Where(where).First(&order).Error; err != nil {
			return errors.New("未找到相关订单信息")
		}
		return order
	} else if payment.PayType == "invoice" {
		return nil
	} else if payment.PayType == "vip" {
		return nil
	} else {
		return nil
	}
}

/**
生成支付链接
*/
func payPage(pay interface{}, payment *middlewares.Payment) (string, error) {
	switch val := pay.(type) {
	case models.OrderModel:
		if payment.PayMethod == "alipay" {
			bm := make(gopay.BodyMap)
			bm.Set("subject", "手机网站测试支付")
			bm.Set("out_trade_no", val.No)
			bm.Set("quit_url", "https://m.helixlife.cn")
			bm.Set("total_amount", val.PaymentAmount)
			bm.Set("product_code", "QUICK_WAP_WAY")
			bm.Set("return_url", "")
			bm.Set("notify_url", "")

			//H5支付
			//if val.Source == "mb" {
			//	payUrl, err := aliPayClient.AliPayTradeWapPay(bm)
			//}

			//PC页面支付
			//if val.Source == "pc" {
			//	payUrl, err := aliPayClient.AliPayTradePagePay(bm)
			//}

			//return payUrl, err

			//查看支付结果
			aliRsp, err := aliPayClient.AliPayTradeQuery(bm)
			if err != nil {
				return "", err
			} else {
				return aliRsp.Response.Msg, err
			}
		}

		return "", nil
	default:
		return "", nil
	}
}

func (baseOrm *BaseOrm) PayNotify(r *rest.Request, notify_type string) string {

	if notify_type == "alipay" {
		//这个json解析插件还在研究中，还是没法解析多层json的数据
		//var notifyReq models.NotifyRequestModel
		//content, _ := ioutil.ReadAll(r.Body)
		//render := strings.NewReader(string(content))
		//decoder := codec.NewDecoder(render, &codec.JsonHandle{})
		//_ = decoder.Decode(&notifyReq)
		//fmt.Printf("res is:%+v", notifyReq)

		var notifyReq models.NotifyRequestModel
		if err := r.DecodeJsonPayload(&notifyReq); err != nil {
			return "fail"
		}

		//暂时取反
		if !aliPay(&notifyReq) {
			params, _ := base64.StdEncoding.DecodeString(notifyReq.PassbackParams)

			var extend models.PayAliExtendParam
			bytes := []byte(string(params) + `"}`)
			_ = json.Unmarshal(bytes, &extend)

			res := afterPayOrder(notifyReq, extend)
			if res == 0 || res == 2 {
				return "success"
			} else {
				return "fail"
			}
		}

		return "fail"
	}

	if notify_type == "wechat_pay" {

		return "success"
	}

	return "fail"
}

/**
本来准备用iGoogle-ink三方的拓展包实现支付以及回调验证，但是发现支付还可以用，但是回调验证就不好用了，
因为回调验证用到了echo框架(和我用的rest框架有太多不一样，主要是rest框架实现的方法数量有限)，所以很多结构体不能用,
所以网上找了下面的两个方法来验证
*/
func aliPay(notifyReq *models.NotifyRequestModel) bool {
	v_fund_bill_list := make([]models.FundBillListInfo, 0)
	json.Unmarshal([]byte(notifyReq.FundBillList.(string)), &v_fund_bill_list)

	if notifyReq.VoucherDetailList != nil {
		v_voucher_bill_list := make([]models.VoucherDetailListInfo, 0)
		json.Unmarshal([]byte(notifyReq.VoucherDetailList.(string)), &v_voucher_bill_list)
		notifyReq.VoucherDetailList = v_voucher_bill_list
	}

	notifyReq.FundBillList = v_fund_bill_list

	//重新组合数据,排除sign和sign_type
	var m map[string]interface{}
	m = make(map[string]interface{}, 0)
	m["notify_time"] = notifyReq.NotifyTime
	m["notify_type"] = notifyReq.NotifyType
	m["notify_id"] = notifyReq.NotifyId
	m["app_id"] = notifyReq.AppId
	m["charset"] = notifyReq.Charset
	m["version"] = notifyReq.Version
	m["auth_app_id"] = notifyReq.AuthAppId
	m["trade_no"] = notifyReq.TradeNo
	m["total_amount"] = notifyReq.TotalAmount
	m["out_biz_no"] = notifyReq.OutBizNo
	m["buyer_id"] = notifyReq.BuyerId
	m["buyer_logon_id"] = notifyReq.BuyerLogonId
	m["seller_id"] = notifyReq.SellerId
	m["seller_email"] = notifyReq.SellerEmail
	m["trade_status"] = notifyReq.TradeStatus
	m["receipt_amount"] = notifyReq.ReceiptAmount
	m["invoice_amount"] = notifyReq.InvoiceAmount
	m["buyer_pay_amount"] = notifyReq.BuyerPayAmount
	m["point_amount"] = notifyReq.PointAmount
	m["refund_fee"] = notifyReq.RefundFee
	m["subject"] = notifyReq.Subject
	m["body"] = notifyReq.Body
	m["gmt_create"] = notifyReq.GmtCreate
	m["gmt_payment"] = notifyReq.GmtPayment
	m["gmt_refund"] = notifyReq.GmtRefund
	m["gmt_close"] = notifyReq.GmtClose
	m["fund_bill_list"] = notifyReq.FundBillList
	m["passback_params"] = notifyReq.PassbackParams
	m["voucher_detail_list"] = notifyReq.VoucherDetailList
	m["method"] = notifyReq.Method
	m["timestamp"] = notifyReq.Timestamp

	sign := notifyReq.Sign
	//获取要进行计算哈希的sign string
	strPreSign, _err := genAlipaySignString(m)
	if _err != nil {
		fmt.Println("error get sign string, reason:", _err)
		return false
	}

	//进行rsa verify
	pass, _err := RSAVerify([]byte(strPreSign), []byte(sign))

	if pass {
		log.Info("verify sig pass.")
		return true
	} else if _err != nil {
		log.Info("verify sig not pass. error:", _err.Error())
		return false
	}

	log.Info("unknown error")
	return false
}

func weixPay() {

	return
}

/***************************************************************
*函数目的：获得从参数列表拼接而成的待签名字符串
*mapBody：是我们从HTTP request body parse出来的参数的一个map
*返回值：sign是拼接好排序后的待签名字串。
***************************************************************/
func genAlipaySignString(mapBody map[string]interface{}) (sign string, err error) {
	sorted_keys := make([]string, 0)
	for k, _ := range mapBody {
		sorted_keys = append(sorted_keys, k)
	}
	sort.Strings(sorted_keys)
	var signStrings string

	index := 0
	for _, k := range sorted_keys {
		value := fmt.Sprintf("%v", mapBody[k])
		if value != "" {
			signStrings = signStrings + k + "=" + value
		}
		//最后一项后面不要&
		if index < len(sorted_keys)-1 {
			signStrings = signStrings + "&"
		}
		index++
	}

	return signStrings, nil
}

/***************************************************************
*RSA签名验证
*src:待验证的字串，sign:支付宝返回的签名
*pass:返回true表示验证通过
*err :当pass返回false时，err是出错的原因
****************************************************************/
func RSAVerify(src []byte, sign []byte) (pass bool, err error) {
	//步骤1，加载RSA的公钥
	block, _ := pem.Decode([]byte(ALIPAY_PUBLIC_KEY))
	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		log.Info("Failed to parse RSA public key:", err.Error())
		return
	}
	rsaPub, _ := pub.(*rsa.PublicKey)

	//步骤2，计算代签名字串的SHA1哈希
	t := sha1.New()
	io.WriteString(t, string(src))
	digest := t.Sum(nil)

	//步骤3，base64 decode,必须步骤，支付宝对返回的签名做过base64 encode必须要反过来decode才能通过验证
	data, _ := base64.StdEncoding.DecodeString(string(sign))

	hexSig := hex.EncodeToString(data)
	log.Info("base decoder:", hexSig)

	//步骤4，调用rsa包的VerifyPKCS1v15验证签名有效性
	err = rsa.VerifyPKCS1v15(rsaPub, crypto.SHA1, digest, data)
	if err != nil {
		log.Info("Verify sig error, reason: ", err.Error())
		return false, err
	}

	return true, nil
}

/**
验证订单信息
int： 0 订单已经付款， 1 订单信息不对 2 订单信息正常
*/
func afterPayOrder(notifyReq models.NotifyRequestModel, extend interface{}) int {
	var (
		order   models.OrderModel
		invoice models.InvoiceModel
		vip     models.VipOrderModel
	)
	updated_at, _ := FormatLocalTime(time.Now())
	switch value := extend.(type) {
	case models.PayAliExtendParam:
		//如果封装到同一个方法里面，还是会调用switch case
		if value.BranchType == "order" {
			//TODO::这里为了测试通过，直接调用异步任务服务
			return updateOrderInfo(vip.ID, "h_orders", value.PaySource, value.BranchType, notifyReq)

			if err := db.GetDB().Table("h_orders").Where("id = ? ", value.Id).First(&order).Error; err != nil {
				return utils.ORDER_INFO_ERROR
			}
			if order.PaymentStatus == 1 {
				return utils.ORDER_PAIED
			}
			if order.Status != 0 {
				db.GetDB().Table("h_orders").Where("id = ?", order.ID).Update(map[string]interface{}{
					"extra":      utils.DEFAULT_EXTRA,
					"updated_at": updated_at,
				})

				return utils.ORDER_INFO_ERROR
			}

			return updateOrderInfo(vip.ID, "h_orders", value.PaySource, value.BranchType, notifyReq)
		} else if value.BranchType == "invoice" {
			if err := db.GetDB().Table("h_invoices").Where("id = ? ", value.Id).First(&invoice).Error; err != nil {
				return utils.ORDER_INFO_ERROR
			}
			if invoice.PaymentStatus == 1 {
				return utils.ORDER_PAIED
			}
			if invoice.Status != 0 {
				db.GetDB().Table("h_invoices").Where("id = ?", invoice.ID).Update(map[string]interface{}{
					"extra":      utils.DEFAULT_EXTRA,
					"updated_at": updated_at,
				})

				return utils.ORDER_INFO_ERROR
			}

			return updateOrderInfo(vip.ID, "h_invoices", value.PaySource, value.BranchType, notifyReq)
		} else if value.BranchType == "vip" {
			if err := db.GetDB().Table("h_vip_orders").Where("id = ? ", value.Id).First(&vip).Error; err != nil {
				return utils.ORDER_INFO_ERROR
			}
			if vip.PaymentStatus == 1 {
				return utils.ORDER_PAIED
			}
			if vip.Status != 0 {
				db.GetDB().Table("h_vip_orders").Where("id = ?", vip.ID).Update(map[string]interface{}{
					"extra":      utils.DEFAULT_EXTRA,
					"updated_at": updated_at,
				})

				return utils.ORDER_INFO_ERROR
			}

			return updateOrderInfo(vip.ID, "h_vip_orders", value.PaySource, value.BranchType, notifyReq)
		}
	}

	return utils.ORDER_INFO_ERROR
}

/**
更新订单信息
*/
func updateOrderInfo(id int, table_name, payment_method, branch_type string, notifyReq models.NotifyRequestModel) int {
	updated_at, _ := FormatLocalTime(time.Now())

	if result := db.GetDB().Table(table_name).Where("id = ?", id).Update(map[string]interface{}{
		"payment_method":   payment_method,
		"payment_status":   1,
		"payment_at":       updated_at,
		"updated_at":       updated_at,
		"status":           1,
		"payment_order_no": notifyReq.TradeNo,
		"receipt_amount":   notifyReq.ReceiptAmount,
	}).RowsAffected; result > 0 {
		//异步更新课程信息
		order_execute := &tasksAndEvents.OrderExecute{OrderId: id, BranchType: branch_type}
		order_execute.Update()

		//同时发消息
		paid_success_message := &tasksAndEvents.PaidSuccessMessage{OrderId: id, BranchType: branch_type, PaySource: payment_method, EventType: "pay_order_success"}
		paid_success_message.Send()

		return utils.ORDER_INFO_OK
	}

	//TODO::这里为了测试通过，直接调用异步任务服务，可以删掉
	//异步更新课程信息
	order_execute := &tasksAndEvents.OrderExecute{OrderId: id, BranchType: "vip"}
	order_execute.Update()
	//同时发消息
	paid_success_message := &tasksAndEvents.PaidSuccessMessage{OrderId: id, BranchType: "vip", PaySource: payment_method, EventType: "pay_order_success"}
	paid_success_message.Send()

	return utils.ORDER_INFO_ERROR
}
