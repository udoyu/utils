package alipay

import (
	"fmt"
	"errors"
)

type AliPay struct {
	MD5Key string
}

func (this AliPay) SignCheck (raw string) (bool, error) {
	kvs, sign, sign_type := NewKVS(raw)
	if sign == "" {
		return false, errors.New("Error 102 : Sign is Empty")
	}
	if sign_type == "MD5" {
		return MD5Sign(kvs.String(), this.MD5Key) == sign, nil
	}
	
	return false, errors.New(fmt.Sprint("Error 103 : Unkown SignType ", sign_type))
}