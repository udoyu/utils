package xrsa

import (
	"encoding/pem"
	"errors"
	"fmt"
	"crypto"
	"crypto/x509"
	"crypto/rsa"
	"crypto/rand"
	"strconv"
)

func hash(msg string) []byte {
	sh := crypto.SHA1.New()
	sh.Write([]byte(msg))
	hash := sh.Sum(nil)
	return hash
}

func HexDecode(value string) []byte {
	buf := make([]byte, len(value)/2)
	for i := 0; i < len(buf); i++ {
		high, _ := strconv.ParseInt(value[i*2:i*2+1], 16, 64)
		low, _ := strconv.ParseInt(value[i*2+1:i*2+2], 16, 64)
		buf[i] = byte(high*16 + low)
	}
	return buf
}
func HexEncode(value []byte) string {
	str := ""
	for _, b := range value {
		str += fmt.Sprintf("%02x", uint8(b)&0xFF)
	}
	return str
}
func HexStringToPublicKey(strPub string) (*rsa.PublicKey, error) {
	pk, err := x509.ParsePKIXPublicKey(HexDecode(strPub))
	if err != nil {
		return nil, err
	}
	return pk.(*rsa.PublicKey), nil
}

func RsaPriKeyFromPEM(priKeyPEM []byte) (*rsa.PrivateKey, error) {
	PEMBlock, _ := pem.Decode(priKeyPEM)
	if PEMBlock == nil {
		return nil, errors.New("Error 200 : pem.Decode failed")
	}
	if PEMBlock.Type != "RSA PRIVATE KEY" {
		return nil, errors.New(fmt.Sprint("Error 201 : Wrong key type, type=", PEMBlock.Type))
	}
	return x509.ParsePKCS1PrivateKey(PEMBlock.Bytes)
}

func RsaPubKeyFromPEM(pubKeyPEM []byte) (*rsa.PublicKey, error) {
	PEMBlock,_ := pem.Decode(pubKeyPEM)
	if PEMBlock == nil {
		return nil, errors.New("Error 200 : pem.Decode failed")
	}
	if PEMBlock.Type != "PUBLIC KEY" {
		return nil, errors.New(fmt.Sprint("Error 201 : Wrong key type, type=", PEMBlock.Type))
	}
	pub, err := x509.ParsePKIXPublicKey(PEMBlock.Bytes)
	if err != nil {
		return nil, err
	}
	return pub.(*rsa.PublicKey), nil
}

func RsaVerify(pub *rsa.PublicKey, srcStr, signStr string) error {
	return rsa.VerifyPKCS1v15(pub, crypto.SHA1, hash(srcStr), HexDecode(signStr))
}

func RsaVerifyBytes(pub *rsa.PublicKey, srcStr string, signBytes []byte) error {
	return rsa.VerifyPKCS1v15(pub, crypto.SHA1, hash(srcStr), signBytes)
}

func RsaSign(pri *rsa.PrivateKey, data string) (string, error) {
	signBuf, err := rsa.SignPKCS1v15(rand.Reader, pri, crypto.SHA1, hash(data))
	if err != nil {
		return "", err
	}
	return HexEncode(signBuf), err
}

func RsaSignToBytes(pri *rsa.PrivateKey, data string) ([]byte, error) {
	signBuf, err := rsa.SignPKCS1v15(rand.Reader, pri, crypto.SHA1, hash(data))
	if err != nil {
		return nil, err
	}
	return signBuf, err
}