package pay

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/xml"
	"fmt"
	"github.com/astaxie/beego/httplib"
	"github.com/udoyu/utils"
	"io"
	rand "math/rand"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type WeiPay struct {
	appid      string //微信分配的公众账号ID（企业号corpid即为此appId）
	mch_id     string //微信支付分配的商户号
	notify_url string //异步回调url
	md5_key    string //md5签名用
}

func NewWeiPay(appid, mch_id, notify_url, md5_key string) WeiPay {
	return WeiPay{
		appid:      appid,
		mch_id:     mch_id,
		notify_url: notify_url,
		md5_key:    md5_key,
	}
}

type UnifiedOrderReq struct {
	XMLName          xml.Name `xml:"xml" json:"-"`
	Appid            string   `xml:"appid"`            //微信分配的公众账号ID（企业号corpid即为此appId）
	Mch_id           string   `xml:"mch_id"`           //微信支付分配的商户号
	Device_info      string   `xml:"device_info"`      //WEB
	Nonce_str        string   `xml:"nonce_str"`        //随机串
	Sign             string   `xml:"sign"`             //签名
	Body             string   `xml:"body"`             //商品或支付单简要描述
	Attach           string   `xml:"attach"`           //附加自定义数据
	Out_trade_no     string   `xml:"out_trade_no"`     //商户订单号
	Total_fee        int      `xml:"total_fee"`        //订单总金额，单位为分
	Spbill_create_ip string   `xml:"spbill_create_ip"` //APP和网页支付提交用户端ip，Native支付填调用微信支付API的机器IP。
	Notify_url       string   `xml:"notify_url"`       //异步回调url
	Trade_type       string   `xml:"trade_type"`       //JSAPI--公众号支付、NATIVE--原生扫码支付、APP--app支付
	Product_id       string   `xml:"product_id"`       //商品id
}

type UnifiedOrderRsp struct {
	XMLName     xml.Name `xml:"xml" json:"-"`
	Return_code string   `xml:"return_code" json:"return_code"`
	Return_msg  string  `xml:"return_msg,omitempty" json:"return_msg,omitempty"`

	Appid        string `xml:"appid,omitempty" json:"appid,omitempty"`
	Mch_id       string `xml:"mch_id,omitempty" json:"mch_id,omitempty"`
	Device_info  string `xml:"device_info,omitempty" json:"device_info,omitempty"`
	Nonce_str    string `xml:"nonce_str,omitempty" json:"nonce_str,omitempty"`
	Sign         string `xml:"sign,omitempty" json:"sign,omitempty"`
	Result_code  string `xml:"result_code,omitempty" json:"result_code,omitempty"`
	Err_code     string `xml:"err_code,omitempty" json:"err_code,omitempty"`
	Err_code_des string `xml:"err_code_des,omitempty" json:"return_msg,omitempty"`

	Trade_type string `xml:"trade_type,omitempty" json:"trade_type,omitempty"`
	Prepay_id  string `xml:"prepay_id,omitempty" json:"prepay_id,omitempty"`
	Code_url   string `xml:"code_url,omitempty" json:"code_url,omitempty"`
}

type WeiNotifyReq struct {
	XMLName     xml.Name `xml:"xml" json:"-"`
	Return_code string   `xml:"return_code" json:"return_code"`
	Return_msg  string   `xml:"return_msg,omitempty" json:"return_msg,omitempty"`

	Appid        string `xml:"appid,omitempty" json:"appid,omitempty"`
	Mch_id       string `xml:"mch_id,omitempty" json:"mch_id,omitempty"`
	Device_info  string `xml:"device_info,omitempty" json:"device_info,omitempty"`
	Nonce_str    string `xml:"nonce_str,omitempty" json:"nonce_str,omitempty"`
	Sign         string `xml:"sign,omitempty" json:"sign,omitempty"`
	Result_code  string `xml:"result_code,omitempty" json:"result_code,omitempty"`
	Err_code     string `xml:"err_code,omitempty" json:"err_code,omitempty"`
	Err_code_des string `xml:"err_code_des,omitempty" json:"return_msg,omitempty"`
	Openid       string `xml:"openid,omitempty" json:"openid,omitempty"`
	Is_subscribe string `xml:"is_subscribe,omitempty" json:"is_subscribe,omitempty"`
	Trade_type   string `xml:"trade_type,omitempty" json:"trade_type,omitempty"`

	Bank_type      string `xml:"bank_type,omitempty" json:"bank_type,omitempty"`
	Total_fee      int    `xml:"total_fee,omitempty" json:"total_fee,omitempty"`
	Fee_type       string `xml:"fee_type,omitempty" json:"fee_type,omitempty"`
	Cash_fee       string `xml:"cash_fee,omitempty" json:"cash_fee,omitempty"`
	Cash_fee_type  string `xml:"cash_fee_type,omitempty" json:"cash_fee_type,omitempty"`
	Transaction_id string `xml:"transaction_id,omitempty" json:"transaction_id,omitempty"`
	Out_trade_no   string `xml:"out_trade_no,omitempty" json:"out_trade_no,omitempty"`
	Prepay_id      string `xml:"prepay_id,omitempty" json:"prepay_id,omitempty"`
	Attach         string `xml:"attach,omitempty" json:"attach,omitempty"`
	Time_end       string `xml:"time_end,omitempty" json:"time_end,omitempty"`
	Code_url       string `xml:"code_url,omitempty" json:"code_url,omitempty"`
}

type WeiNotifyRsp struct {
	XMLName     xml.Name `xml:"xml" json:"-"`
	Return_code string   `xml:"return_code" json:"return_code"`
	Return_msg  string   `xml:"return_msg,omitempty" json:"return_msg,omitempty"`
}

func AnyToForm(v interface{}) url.Values {
	form := make(url.Values)
	values := reflect.ValueOf(v)
	vtypes := reflect.TypeOf(v)
	if values.Kind() == reflect.Ptr {
		values = values.Elem()
		vtypes = vtypes.Elem()
	}
	for i := 0; i < values.NumField(); i++ {
		tag := vtypes.Field(i).Tag.Get("xml")
		if strings.Contains(tag, ",") {
			tag = utils.StringCutRightExp(tag, ",", 1)
		}
		form.Add(tag, fmt.Sprint(values.Field(i).Interface()))
	}
	return form
}

func (this WeiPay) UnifiedOrder(wpOptions UnifiedOrderReq) (*UnifiedOrderRsp, error) {
	wpOptions.Appid = this.appid
	wpOptions.Mch_id = this.mch_id
	
	wpOptions.Nonce_str = NonceStr()[:16]
	wpOptions.Notify_url = this.notify_url
	
	wpOptions.Attach = "weipay"
	form := AnyToForm(wpOptions)
	form.Del("xml")
	wpOptions.Sign = this.MD5Sign(form)

	fmt.Println(wpOptions)
	xmlBuf, _ := xml.Marshal(&wpOptions)

	hp := httplib.Post("https://api.mch.weixin.qq.com/pay/unifiedorder")
	rsp := UnifiedOrderRsp{}
	fmt.Println(string(xmlBuf))
	if err := hp.Body(xmlBuf).ToXml(&rsp); err != nil {
		return &rsp, err
	}
	fmt.Println(rsp)
	return_code := rsp.Return_code
	if return_code != "SUCCESS" {
		if rsp.Return_msg != "" {
			return &rsp, fmt.Errorf("Error 302 : %s", rsp.Return_msg)
		} else {
			return &rsp, fmt.Errorf("Error 303 : unknown return_msg")
		}
	}
	if rsp.Result_code != "SUCCESS" {
		return &rsp, fmt.Errorf("Error 305 : %s %s ", rsp.Err_code, rsp.Err_code_des)
	}
	form = AnyToForm(rsp)
	fmt.Println(form)
	form.Del("xml")
	kvs, sign, _ := NewKVSFromForm(form)
	nsign := this.MD5SignFromKVS(kvs)
	if sign != nsign {
		return &rsp, fmt.Errorf("Error 304 : sign failed")
	}
	return &rsp, nil
}

func (this WeiPay) MD5SignFromKVS(kvs KVS) string {
	m := md5.New()
	fmt.Println(kvs.String())
	m.Write([]byte(kvs.String() + "&key=" + this.md5_key))
	sign := hex.EncodeToString(m.Sum(nil))
	return strings.ToUpper(sign)
}

func (this WeiPay) MD5Sign(form url.Values) string {
	//	form.Add("key", this.md5_key)
	form.Del("xml")
	kvs, _, _ := NewKVSFromForm(form)
	return this.MD5SignFromKVS(kvs)
}

func NonceStr() string {
	nano := time.Now().UnixNano()
	rand.Seed(nano)
	rndNum := rand.Int63()
	nonceStr := hash(hash(strconv.FormatInt(nano, 10)) + hash(strconv.FormatInt(rndNum, 10)))
	return nonceStr
}

func hash(text string) string {
	hashMd5 := md5.New()
	io.WriteString(hashMd5, text)
	return fmt.Sprintf("%X", hashMd5.Sum(nil))
}
