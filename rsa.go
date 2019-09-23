package utils

import (
	"crypto"
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"math/big"
	"time"
)

// 生成私钥，公钥
func GenerateKey(bits int) (privateKey, publicKey string, err error) {
	rand.Int(rand.Reader, big.NewInt(time.Now().UnixNano()))
	privKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return
	}
	//生成私钥
	pkcs1PrivateKey := x509.MarshalPKCS1PrivateKey(privKey)
	block := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: pkcs1PrivateKey,
	}
	data := pem.EncodeToMemory(block)
	privateKey = Base64Encoding(data)

	pubKey := &privKey.PublicKey
	derPkix, err := x509.MarshalPKIXPublicKey(pubKey)
	if err != nil {
		return
	}
	block = &pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: derPkix,
	}
	data = pem.EncodeToMemory(block)
	publicKey = Base64Encoding(data)
	return
}

// 签名
func SignData(privateKey, msg string) (sig []byte, err error) {
	//准备签名的数据
	plaintxt := []byte(msg)
	h := md5.New()
	h.Write(plaintxt)
	hashed := h.Sum(nil)

	block, _ := pem.Decode(Base64Decoding(privateKey))
	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return
	}

	//签名
	opts := &rsa.PSSOptions{rsa.PSSSaltLengthAuto, crypto.MD5}
	sig, err = rsa.SignPSS(rand.Reader, priv, crypto.MD5, hashed, opts)
	return
}

// 验证签名
func VerifySign(publicKey, strData, sig string) bool {
	block, _ := pem.Decode(Base64Decoding(publicKey))
	if block == nil {
		return false
	}
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return false
	}
	pub := pubInterface.(*rsa.PublicKey)

	//验证发送方是否为zhaoyingkui
	h := md5.New()
	h.Write([]byte(strData))
	hashed := h.Sum(nil)

	e := rsa.VerifyPSS(pub, crypto.MD5, hashed, Base64Decoding(sig), nil)
	if e == nil {
		return true
	}
	return false
}
