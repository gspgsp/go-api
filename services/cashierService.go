package services

import (
	"edu_api/middlewares"
	"edu_api/models"
	"edu_api/utils"
	"errors"
	"fmt"
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/iGoogle-ink/gopay"
	log "github.com/sirupsen/logrus"
)

var ()

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
