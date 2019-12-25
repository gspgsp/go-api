package services

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"edu_api/middlewares"
	"edu_api/models"
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
		var notifyReq models.NotifyRequestModel
		if err := r.DecodeJsonPayload(&notifyReq); err != nil {
			return err.Error()
		}

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
			return ""
		}

		//进行rsa verify
		pass, _err := RSAVerify([]byte(strPreSign), []byte(sign))

		if pass {
			fmt.Println("verify sig pass.")
		} else {
			fmt.Println("verify sig not pass. error:", _err)
		}
	}

	if notify_type == "wechat_pay" {

	}

	return "success"
}

/**
本来准备用这个三方的拓展包实现支付以及回调验证，但是发现支付还可以用，但是回调验证就不好用了，
因为回调验证用到了echo框架(和我用的rest框架有太多不一样，主要是rest框架实现的方法数量有限)，所以很多结构体不能用,
所以网上找了下面的两个方法来验证
*/
func aliPay(notifyReq *models.NotifyRequestModel) string {
	//验签操作
	ok, err := gopay.VerifyAliPaySign(aliPayConfig.Alipay.PublicKey, notifyReq)
	if err != nil {
		log.Info("alipay verify sign error", err.Error())
		return "fail"
	}

	log.Info("alipay verify sign:", ok)

	//数据库操作

	return "success"

}

func weixPay() {

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
		fmt.Println("k=", k, "v =", mapBody[k])
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
		fmt.Printf("Failed to parse RSA public key: %s\n", err)
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
	fmt.Printf("base decoder: %v, %v\n", string(sign), hexSig)

	//步骤4，调用rsa包的VerifyPKCS1v15验证签名有效性
	err = rsa.VerifyPKCS1v15(rsaPub, crypto.SHA1, digest, data)
	if err != nil {
		fmt.Println("Verify sig error, reason: ", err)
		return false, err
	}

	return true, nil
}
