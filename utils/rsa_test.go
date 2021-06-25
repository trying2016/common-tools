package utils

import (
	"fmt"
	"testing"
)

func TestGenerateKey(t *testing.T) {
	privateKey, publicKey, err := GenerateKey(1024)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Printf("Private Key :%s \n", privateKey)
	fmt.Printf("Public Key :%s \n", publicKey)
	str := "hello Sigin~"
	sig, err := SignData(privateKey, str)
	if err != nil {
		t.Error(err)
		return
	}
	// eyJhcHBpZCI6MTYsImV0IjoxNTY4ODc3Mjg5LCJnaWQiOjE2fQ==.gc1q6BfpXgEqgdOROLV36m8mCdSCeSZwCNJtcza8x8M=
	fmt.Printf("sig:%s", Base64Encoding(sig))
	if !VerifySign(publicKey, str, Base64Encoding(sig)) {
		t.Error("VerifySign Error")
	}
}

func TestGenerateSig(t *testing.T) {
	privateKey, publicKey, err := GenerateKey(256)
	publicKey = "LS0tLS1CRUdJTiBSU0EgUFVCTElDIEtFWS0tLS0tCk1DUXdEUVlKS29aSWh2Y05BUUVCQlFBREV3QXdFQUlKQU50ODhaeTBZTXgzQWdNQkFBRT0KLS0tLS1FTkQgUlNBIFBVQkxJQyBLRVktLS0tLQo="
	privateKey = "LS0tLS1CRUdJTiBSU0EgUFJJVkFURSBLRVktLS0tLQpNRDhDQVFBQ0NRRGJmUEdjdEdETWR3SURBUUFCQWdnVEhkSzlOZUV2b1FJRkFQODVFRDBDQlFEY0tBYkRBZ1VBCjg1QXloUUlFVmI3Tkt3SUZBSU81V3ZzPQotLS0tLUVORCBSU0EgUFJJVkFURSBLRVktLS0tLQo="
	fmt.Printf("Private Key :%s \n", privateKey)
	fmt.Printf("Public Key :%s \n", publicKey)
	str := `{"appid":16,"et":1568877289,"gid":16}`
	sig, err := SignData(privateKey, str)
	if err != nil {
		t.Error(err)
		return
	}
	// eyJhcHBpZCI6MTYsImV0IjoxNTY4ODc3Mjg5LCJnaWQiOjE2fQ==.gc1q6BfpXgEqgdOROLV36m8mCdSCeSZwCNJtcza8x8M=
	fmt.Printf("sig:%s", Base64Encoding(sig))
	if !VerifySign(publicKey, str, "gc1q6BfpXgEqgdOROLV36m8mCdSCeSZwCNJtcza8x8M=") {
		t.Error("VerifySign Error")
	}
}

func TestRsaEncrypt(t *testing.T) {
	privateKey, publicKey, err := GenerateKey(256)
	if err != nil {
		t.Fatalf("GenerateKey fail, error: %v", err)
	}
	fmt.Printf("Private Key :%s \n", privateKey)
	fmt.Printf("Public Key :%s \n", publicKey)
	enRet, err := RsaEncrypt(publicKey, "123123")
	if err != nil {
		t.Fatalf("RsaEncrypt fail, error: %v", err)
	}
	deRet, err := RsaDecrypt(privateKey, enRet)
	if err != nil {
		t.Fatalf("RsaDecrypt fail, error: %v", err)
	}
	fmt.Printf("result %v", deRet)
}
