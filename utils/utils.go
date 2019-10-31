package utils

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// 对账号等进行掩码
func MaskContent(str string) string {
	content := []rune(str)
	if len(content) < 2 {
		return str
	}
	reserveNum := 1
	if len(content)/2 > 2 {
		reserveNum = 2
	}
	const MaskLen = 3
	var contentLen = reserveNum*2 + MaskLen
	data := make([]rune, contentLen, contentLen)
	for i := 0; i < reserveNum; i++ {
		data[i] = content[i]
		data[contentLen-i-1] = content[len(content)-i-1]
	}
	for i := 0; i < MaskLen; i++ {
		data[i+reserveNum] = '*'
	}
	return string(data)
}

func Min(a, b int) int {
	if a < b {
		return a
	} else {
		return b
	}
}

func Min64(a, b int64) int64 {
	if a < b {
		return a
	} else {
		return b
	}
}

func ToShowBalance(balance int64) string {
	str := fmt.Sprintf("%v", balance)
	for i := len(str); i <= 8; i++ {
		str = fmt.Sprintf("0%s", str)
	}
	nSplit := len(str) - 8
	str = str[:nSplit] + "." + str[nSplit:]
	return str
}

func Md5Byte(genSin []byte, body []byte) string {
	arrByte := make([]byte, len(genSin)+len(body))
	copy(arrByte, genSin)
	copy(arrByte[len(genSin):], body)

	h := md5.New()
	h.Write(arrByte)
	return hex.EncodeToString(h.Sum(nil))
}

func MaxInt(a, b int) int {
	if a > b {
		return a
	} else {
		return b
	}
}
