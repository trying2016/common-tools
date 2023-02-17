package utils

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestSubleString(t *testing.T) {
	str := "123abcd"
	fmt.Printf(SubleString(str, "3", "d"))
}

func TestSliceRemove(t *testing.T) {
	arr := [6]int{1, 2, 3, 4, 5, 6}
	SliceRemove(arr, 1)
}

func TestQueue(t *testing.T) {
	qu := NewQueue(100, func(data interface{}) {
		println(data)
	})
	qu.Push(1)
	qu.Push(1)
	count := qu.Count()
	println(count)
}

func TestToFloat(t *testing.T) {
	fValue := ToFloat("100.00")
	fmt.Printf("value %v", fValue)
}

func TestParseCookie(t *testing.T) {
	cookie := ParseCookie("pzapp_uid=%2FediapJ4DIhiwAu9aR8CaWPPqno%3D;pzapp_account=4Swh6TA7hw5n2gcZxn0E7aI33zZJGsF5;pzapp_type=xnZh4KLuphzLSduFcjQtjZzHmB3ZvIiW;pzapp_c=%2B62bbBMuh0puhMpJ%2BCpOSqgQLU8%3D;pzapp_uv=9fgG6kh3RRNRjNdvWdKZVA%3D%3D;pzapp_dp=W%2FIv7%2BtTvbvs67C%2FCl7tew%3D%3D;pzapp_num=9fgG6kh3RRNRjNdvWdKZVA%3D%3D;")
	account := cookie.Decode("pzapp_account", "bitnew_CxT6Egc3T3mBe")
	t.Log(account)
}

func TestHttpClient(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	client := NewHttpClient()
	client.SetContext(&ctx)
	go func() {
		time.Sleep(time.Second / 20)
		cancel()
	}()
	b, err := client.Get("https://www.hpool.com")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(b)
	}
}

func TestToString(t *testing.T) {
	t.Log(ToString([]string{"1", "2", "s 1"}))
}
