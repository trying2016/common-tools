package jsonRpc

import (
	"testing"
	"time"

	"51shcp.com/trying/sync/server/jsonRpc/jsonrpc2"
)

func TestRpcClient(t *testing.T) {
	client, err := jsonrpc2.NewClient("127.0.0.1:9017", VarintObjectCodec{})
	if err != nil {
		return
	}
	client.Call("login", "")
	for {
		client.Call("keepalived", "")
		time.Sleep(time.Second * 60)
	}
}
