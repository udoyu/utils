package alipay
import (
	"crypto/md5"
	"encoding/hex"
	"net/url"
	"sort"
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
func NewKVS (str string) (kvs KVS, sign, sign_type string) {
	body, _ := url.QueryUnescape(str)
	vs := strings.Split(body, "&")
	for _, v := range vs {
		kv := strings.Split(v, "=")
		if kv[0] == "sign" {
			sign = kv[1]
		} else if kv[0] == "sign_type" {
			sign_type = kv[1]
		} else {
			kvs = append(kvs, KV{key: kv[0], value: v})
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

func MD5Sign(str, key string) string {
	
	m := md5.New()
	m.Write([]byte(str+key))
	sign := hex.EncodeToString(m.Sum(nil))
	return sign
}