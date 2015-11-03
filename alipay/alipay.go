package alipay

import (
	"net/url"
	"fmt"
	"errors"
	"strconv"
)

type AliPay struct {
	Partner   string     // 合作者ID
	SellerEmail string	 // 合作者Email
	NotifyUrl string	 // 异步url
	ReturnUrl string     // 同步url
	MD5Key string        // md5key
}

func (this AliPay) SignCheckString (raw string) (bool, error) {
	kvs, sign, sign_type := NewKVSFromString(raw)
	if sign == "" {
		return false, errors.New("Error 102 : Sign is Empty")
	}
	if sign_type == "MD5" {
		ok := MD5Sign(kvs.String(), this.MD5Key) == sign
		if !ok {
			return false, errors.New("Error 103 : Sign not Equal")
		}
		return true, nil		
	}
	
	return false, errors.New(fmt.Sprint("Error 104 : Unkown SignType ", sign_type))
}

func (this AliPay) SignCheckForm(form url.Values) (bool, error) {
	kvs, sign, sign_type := NewKVSFromForm(form)
	if sign == "" {
		return false, errors.New("Error 102 : Sign is Empty")
	}
	if sign_type == "MD5" {
		ok := MD5Sign(kvs.String(), this.MD5Key) == sign
		if !ok {
			return false, errors.New("Error 103 : Sign not Equal")
		}
		return true, nil		
	}
	
	return false, errors.New(fmt.Sprint("Error 104 : Unkown SignType ", sign_type))
}

// 生成订单的参数
type Options struct {
	OrderId  string  // 订单唯一id
	Fee      float32 // 价格
	NickName string  // 充值账户名称
	Subject  string  // 充值描述
}

type Result struct {
	// 状态
	Status int
	// 本网站订单号
	OrderNo string
	// 支付宝交易号
	TradeNo string
	// 买家支付宝账号
	BuyerEmail string
	// 错误提示
	Message string
	// 消费
	Money float64
}

/* 生成支付宝即时到帐提交表单html代码 */
func (this AliPay) Form(opts Options) string {
	//实例化参数
	param := &AlipayParameters{}
	param.InputCharset = "utf-8"
	param.Body = "为" + opts.NickName + "充值" + strconv.FormatFloat(float64(opts.Fee), 'f', 2, 32) + "元"
	param.NotifyUrl = this.NotifyUrl
	param.OutTradeNo = opts.OrderId
	param.Partner = this.Partner
	param.PaymentType = 1
	param.ReturnUrl = this.ReturnUrl
	param.SellerEmail = this.SellerEmail
	param.Service = "create_direct_pay_by_user"
	param.Subject = opts.Subject
	param.TotalFee = opts.Fee

	//生成签名
	sign := MD5SignFromInterface(param, this.MD5Key)

	//追加参数
	param.Sign = sign
	param.SignType = "MD5"

	//生成自动提交form
	return `
		<form id="alipaysubmit" name="alipaysubmit" action="https://mapi.alipay.com/gateway.do?_input_charset=utf-8" method="get" style='display:none;'>
			<input type="hidden" name="_input_charset" value="` + param.InputCharset + `">
			<input type="hidden" name="body" value="` + param.Body + `">
			<input type="hidden" name="notify_url" value="` + param.NotifyUrl + `">
			<input type="hidden" name="out_trade_no" value="` + param.OutTradeNo + `">
			<input type="hidden" name="partner" value="` + param.Partner + `">
			<input type="hidden" name="payment_type" value="` + strconv.Itoa(int(param.PaymentType)) + `">
			<input type="hidden" name="return_url" value="` + param.ReturnUrl + `">
			<input type="hidden" name="seller_email" value="` + param.SellerEmail + `">
			<input type="hidden" name="service" value="` + param.Service + `">
			<input type="hidden" name="subject" value="` + param.Subject + `">
			<input type="hidden" name="total_fee" value="` + strconv.FormatFloat(float64(param.TotalFee), 'f', 2, 32) + `">
			<input type="hidden" name="sign" value="` + param.Sign + `">
			<input type="hidden" name="sign_type" value="` + param.SignType + `">
		</form>
		<script>
			document.forms['alipaysubmit'].submit();
		</script>
	`
}

func (this AliPay) Check(form url.Values) *Result {
	result := &Result{}
	
	result.OrderNo = form.Get("out_trade_no")
	result.TradeNo = form.Get("trade_no")
	result.BuyerEmail = form.Get("buyer_email")
	{
	money, err := strconv.ParseFloat(form.Get("total_fee"), 64)
	if err != nil {
		result.Status = -3
		result.Message = fmt.Sprint("消费额错误,total_fee=", form.Get("total_fee"))
		return result
	}
	result.Money = money
	}
	if result.OrderNo == "" {
		//不存在交易号
		result.Status = -1
		result.Message = "站交易号为空"
		return result
	}
	{
		ok, err := this.SignCheckForm(form)
		if !ok {
			result.Status = -2
			result.Message = err.Error()
			return result
		}
	}
	
	// 判断订单是否已完成
	tradeStatus := form.Get("trade_status")
	if tradeStatus == "TRADE_FINISHED" || tradeStatus == "TRADE_SUCCESS" { //交易成功
		result.Status = 1
		return result
	} else { // 交易未完成，返回错误代码-4
		result.Status = -4
		result.Message = "交易未完成"
		return result
	}
	
	return result
}

type AlipayParameters struct {
	InputCharset string  `json:"_input_charset"` //网站编码
	Body         string  `json:"body"`           //订单描述
	NotifyUrl    string  `json:"notify_url"`     //异步通知页面
	OutTradeNo   string  `json:"out_trade_no"`   //订单唯一id
	Partner      string  `json:"partner"`        //合作者身份ID
	PaymentType  uint8   `json:"payment_type"`   //支付类型 1：商品购买
	ReturnUrl    string  `json:"return_url"`     //回调url
	SellerEmail  string  `json:"seller_email"`   //卖家支付宝邮箱
	Service      string  `json:"service"`        //接口名称
	Subject      string  `json:"subject"`        //商品名称
	TotalFee     float32 `json:"total_fee"`      //总价
	Sign         string  `json:"sign"`           //签名，生成签名时忽略
	SignType     string  `json:"sign_type"`      //签名类型，生成签名时忽略
}