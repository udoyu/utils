package pay

import (
	"crypto/md5"
	"crypto/rsa"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/udoyu/utils/xrsa"
	"net/url"
	"sort"
	"strconv"
	"strings"
)

type KV struct {
	key   string
	value string
}

type KVS []KV

func (this KVS) Len() int           { return len(this) }
func (this KVS) Swap(i, j int)      { this[i], this[j] = this[j], this[i] }
func (this KVS) Less(i, j int) bool { return this[i].key < this[j].key }

//filter md5 and sort
func NewKVSFromString(str string) (kvs KVS, sign, sign_type string) {
	body, _ := url.QueryUnescape(str)
	vs := strings.Split(body, "&")
	for _, v := range vs {
		kv := strings.Split(v, "=")
		if kv[0] == "sign" {
			sign = kv[1]
		} else if kv[0] == "sign_type" {
			sign_type = kv[1]
		} else if kv[1] != "" {
			kvs = append(kvs, KV{key: kv[0], value: v})
		}
	}
	sort.Sort(kvs)
	return kvs, sign, sign_type
}

func NewKVSFromForm(form url.Values) (kvs KVS, sign, sign_type string) {
	for key, _ := range form {
		v := form.Get(key)
		if key == "sign" {
			sign = v
		} else if key == "sign_type" {
			sign_type = v
		} else {
			kvs = append(kvs, KV{key: key, value: key + "=" + v})
		}
	}
	sort.Sort(kvs)
	return kvs, sign, sign_type
}

func NewKVSFromMap(m map[string]interface{}) (kvs KVS, sign, sign_type string) {
	for key, value := range m {
		v := fmt.Sprint(value)
		if key == "sign" {
			sign = v
		} else if key == "sign_type" {
			sign_type = v
		} else {
			kvs = append(kvs, KV{key: key, value: key + "=" + v})
		}
	}
	sort.Sort(kvs)
	return kvs, sign, sign_type
}

func (this KVS) String() string {
	newBody := ""
	for _, v := range this {
		newBody = newBody + v.value + "&"
	}

	newBody = newBody[:len(newBody)-1]
	return newBody
}

func (this KVS) Form() url.Values {
	form := make(url.Values)
	for _, v := range this {
		form.Add(v.key, v.value)
	}
	return form
}

func MD5Sign(str, key string) string {

	m := md5.New()
	m.Write([]byte(str + key))
	sign := hex.EncodeToString(m.Sum(nil))
	return sign
}

func RSASignToBytes(str string, priv *rsa.PrivateKey) ([]byte, error) {
	return xrsa.RsaSignToBytes(priv, str)
}

func RsaVerifyBytes(srcStr string, signBytes []byte, pubc *rsa.PublicKey) error {
	return xrsa.RsaVerifyBytes(pubc, srcStr, signBytes)
}

// 按照支付宝规则生成sign
func MD5SignFromInterface(param interface{}, key string) string {
	//解析为字节数组
	paramBytes, err := json.Marshal(param)
	if err != nil {
		return ""
	}

	//重组字符串
	var sign string
	oldString := string(paramBytes)

	//为保证签名前特殊字符串没有被转码，这里解码一次
	oldString = strings.Replace(oldString, `\u003c`, "<", -1)
	oldString = strings.Replace(oldString, `\u003e`, ">", -1)

	//去除特殊标点
	oldString = strings.Replace(oldString, "\"", "", -1)
	oldString = strings.Replace(oldString, "{", "", -1)
	oldString = strings.Replace(oldString, "}", "", -1)
	paramArray := strings.Split(oldString, ",")

	for _, v := range paramArray {
		detail := strings.SplitN(v, ":", 2)
		//排除sign和sign_type
		if detail[0] != "sign" && detail[0] != "sign_type" {
			//total_fee转化为2位小数
			if detail[0] == "total_fee" {
				number, _ := strconv.ParseFloat(detail[1], 32)
				detail[1] = strconv.FormatFloat(number, 'f', 2, 64)
			}
			if sign == "" {
				sign = detail[0] + "=" + detail[1]
			} else {
				sign += "&" + detail[0] + "=" + detail[1]
			}
		}
	}

	//追加密钥
	sign += key

	//md5加密
	m := md5.New()
	m.Write([]byte(sign))
	sign = hex.EncodeToString(m.Sum(nil))
	return sign
}
