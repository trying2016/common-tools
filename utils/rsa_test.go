package utils

import (
	"encoding/hex"
	"fmt"
	"testing"
)

var (
	publicKey   string
	privateKey  string
	encryptData string
)

func init() {
	var err error
	privateKey, publicKey, err = GenerateKey(1024)
	if err != nil {
		return
	}
	data, err := hex.DecodeString("a31bb3c0a683f42fad7c43ac50993dcdcb4b636195333f4bcdd34b5243881f0fd543bda03fc2d4ed7aa09351c1d8de1fdc4e19c1c1b918b7794d5d9de9a83af1d62b2feb27881c19e0fc482c82313b1ee77627c85b689f809aeb7efaf2c2cc7dcd2d6173783c684816715b62cf99bd0475cde801f596a463884eee5668e9d4e9b7fc416aac816e0100")
	if err == nil {
		encryptData, err := RsaEncryptRaw(publicKey, data)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(hex.EncodeToString(encryptData))
		}
	}
	encryptData, err = RsaEncrypt(publicKey, "a31bb3c0a683f42fad7c43ac50993dcdcb4b636195333f4bcdd34b5243881f0fd543bda03fc2d4ed7aa09351c1d8de1fdc4e19c1c1b918b7794d5d9de9a83af1d62b2feb27881c19e0fc482c82313b1ee77627c85b689f809aeb7efaf2c2cc7dcd2d6173783c684816715b62cf99bd0475cde801f596a463884eee5668e9d4e9b7fc416aac816e0100")
	if err != nil {
		fmt.Println(err)
	}
}

func TestGenerateKey(t *testing.T) {
	privateKey, publicKey, err := GenerateKey(256)
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
	privateKey, publicKey, err := GenerateKey(1024)
	if err != nil {
		t.Fatalf("GenerateKey fail, error: %v", err)
	}
	fmt.Printf("Private Key :%s \n", privateKey)
	fmt.Printf("Public Key :%s \n", publicKey)
	enRet, err := RsaEncrypt(publicKey, Map{
		"fpk":         "a76c35ca3c043c0a3a14c8aa4cc528e5010ac45df9845f69fdb1273e7895bdf624f4b78508ba0e293e997a75f3108454",
		"fingerprint": 1,
	}.ToJson())
	if err != nil {
		t.Fatalf("RsaEncrypt fail, error: %v", err)
	}
	deRet, err := RsaDecrypt(privateKey, enRet)
	if err != nil {
		t.Fatalf("RsaDecrypt fail, error: %v", err)
	}
	fmt.Printf("result %v", deRet)
}

func TestRsa(t *testing.T) {
	pubKey := "LS0tLS1CRUdJTiBQVUJMSUMgS0VZLS0tLS0KTUlJQklqQU5CZ2txaGtpRzl3MEJBUUVGQUFPQ0FROEFNSUlCQ2dLQ0FRRUF2NEZoM1dkTHdsMFdDcmxuRjhvTgpKcjNkNzkzMlNOUU1SQXlMRmVtVFJtdm5xQll5aFNsS3ZpNDhwNE84S24rK25TQ1RDNVBWQ0NuVDNIZXA2L2lKClFNZ2s0WlBWNGJrTURvb2EwVk1QaGN2QlpqR0V6VmNQRGJxY0xwTDJOTjFUcmE3aEw1bitkS2ZCRysyeVlmdWMKa2Z4dElTajdRYzFGa1Y0eGFCSkdZcGZYbWEwazJZWlJKOW9jYnVubGpySmdIeWhlaVQxSTdkajNvazQ4ZWU2dgp3QUdpK2xmdWRQZHpjOE9SL0hDUnZZT25vM2R6YjZJN3VJZHc1TGRUNTlrY3UwdVFKUmh1UkZLY2RLVDU4VEt4Cmtlc3p0bHpTbmgxSXJraHJrbzFUOWluSFdodkVUbVdUNHUrek5TaHVRb1hCZGQyWGovQkV4eXQ4ZlhTUFBMdzYKL3dJREFRQUIKLS0tLS1FTkQgUFVCTElDIEtFWS0tLS0tCg=="
	privateKey := "LS0tLS1CRUdJTiBSU0EgUFJJVkFURSBLRVktLS0tLQpNSUlFcEFJQkFBS0NBUUVBdjRGaDNXZEx3bDBXQ3JsbkY4b05KcjNkNzkzMlNOUU1SQXlMRmVtVFJtdm5xQll5CmhTbEt2aTQ4cDRPOEtuKytuU0NUQzVQVkNDblQzSGVwNi9pSlFNZ2s0WlBWNGJrTURvb2EwVk1QaGN2QlpqR0UKelZjUERicWNMcEwyTk4xVHJhN2hMNW4rZEtmQkcrMnlZZnVja2Z4dElTajdRYzFGa1Y0eGFCSkdZcGZYbWEwawoyWVpSSjlvY2J1bmxqckpnSHloZWlUMUk3ZGozb2s0OGVlNnZ3QUdpK2xmdWRQZHpjOE9SL0hDUnZZT25vM2R6CmI2STd1SWR3NUxkVDU5a2N1MHVRSlJodVJGS2NkS1Q1OFRLeGtlc3p0bHpTbmgxSXJraHJrbzFUOWluSFdodkUKVG1XVDR1K3pOU2h1UW9YQmRkMlhqL0JFeHl0OGZYU1BQTHc2L3dJREFRQUJBb0lCQVFDUE5BSHBuU2V5dFFjWgpoK0RHa2tuWlVadVhsZ1JvRzJEOHRlQi94MFZoSUtsL01QSWdUMnRiNFpscnJuL1R5K2pPK0ovY3hYUkZBWG95CjM2ektEdlViNDA1by9MS3djejdIMUpBUFBheGE0YTNDYkg4aFNkdXc2WDJHK2xCdjRaMkVRRVNWNHZLN2F3SmwKanc2WVpKMkZNUnl0OGtaSXcyWWxPU2w2NkVlSHl6T0tERTZkWjc0WTZjTG9DZTBuUzh2czE0d255MmRrb202WQpYeHJvdGJlUXk4SHZIYjhmeGwxdlJzcFlocnBYTU5ZWkRGN3k5QWtyM3BVOWRhTlU2MXNpbjFPZGM4clczdEpTCndOSldmMlJKN0xvcWdxTzVCR1dyZWxpSGlqK0tJb05nZTlhVElKdzdoSkgyVDBYblVvS3Fvekp2SmtRSDZRekMKaEhneUtRdXhBb0dCQVA2Q016ZUc2VUpyUTFzR1NxZi9ZdE9PV25QRU1BRUFPQTFac3RlMy91Wk5PVkNvTGZRNQpJMU9ZSzNSSWthSFRkaEJua0pyVnAybFNQRm5BV3VPZDFRem9YcTJLSlhpNHhudXk1RytXaHU3SW43ZWhaZUlDCll0S3JuaWNxWjBNOS9IR0MyQnFIY0I4WGF2TUZ6TndyRkhUbjRPRGd1dXNrK3drQWRDQnVnQzQzQW9HQkFNQ2cKcXhMalhlWmVqKzhxSEVFTVVLL0w1SlBpK3lOUFFZdEpJQzlWMXFZWXUrS2ZUMVRLTUZFMEFVQkZqakhINTVETQpoM0hqK3JvcGNzejhUUU9WMkdseitlVnl1Vjcvb3VzTEZtVzlRcE9pbU9mRFRZY0Z4N25IV0IwRk9QdzhNa3pvCmhiNUZKZGVGNkhHQUdEcnpnWnVOY3lnSnIybVRVOGxIcytBSmt6VjVBb0dBQTk4dE1rb09JR0dMVzhZanVweUsKLzFicUQxckx0Q2d4c2hwTU96WGtYZEtNN2FveFVNYlJ2OExQM213QU15c0pYOFNEa2Fkd2JZeS91RW5SMkNhZAppQjI4MnYwQUJ1OGdyZDhSMUpUQXByOU1scm1RMkRoYkVvTmoyNHFzbVh4RzY5OG10SGljL3d3WEoyMU9LWWRLClAyRUxyY0FkZDloUExWcmhhV0RrK0U4Q2dZQVRYcFNWTzZPdmpJYXdwKzFiWlIrZjdjSzRWRFNvb2ttVzllMTAKbFE4V2VKbzcrWVVDbzZva0lEU1gvK2FDZnZWOEMvVDZzTS8vZERlRkFVSEZRSVlZWkg4V1lXamVjcG94UkZZbQpPTmUrL0xTZmFYWWNRNnFIdGRIWWFUUUh1UTNkeWV2WTdCOFlBdkVTcUt6SEVNRlJvUEdaYncvaWI5d1crMHpOCkZHbG5FUUtCZ1FDMTZDN0poaTUyZW56VWdrMXkrSmsvRVR2SWFNalNVdlN2elVHV045VkNGcEhTcTVRbXBwTW4KVG5ZU2RZelY5Z0VUQmJwWjFXSlZmNFNIdERINmVzMEx5VVBTVjF0UnpEY3QxZzZtUTlBMzFheWw2anJTd2ZGSwpUa3ZtV0x3RFBLTVQycHFXNzUwZzhRWHMxRlkwNFArVXI0ZGlGQ05MQ2RhSEk5ZkQ2STBQWGc9PQotLS0tLUVORCBSU0EgUFJJVkFURSBLRVktLS0tLQo="
	enStr, err := RsaEncrypt(pubKey, "1234567")
	if err != nil {
		t.Fatalf("RsaEncrypt fail, error: %v", err)
	}
	deStr, err := RsaDecrypt(privateKey, enStr)
	if err != nil {
		t.Fatalf("RsaDecrypt fail, error: %v", err)
	}
	fmt.Printf("result %v", deStr)
}

func BenchmarkRsa(b *testing.B) {
	for i := 0; i < b.N; i++ {
		if _, err := RsaDecrypt(privateKey, encryptData); err != nil {
			b.Error(err)
		}
	}
}
