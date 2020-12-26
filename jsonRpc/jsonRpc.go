package jsonRpc

import (
	"context"
	"errors"
	"net"
	"regexp"
	"strings"

	"51shcp.com/trying/sync/definition/entity"

	"51shcp.com/trying/sync/server/jsonRpc/jsonrpc2"
)

var (
	ErrorNotHandle = errors.New("not handle")
)

type Server struct {
	mapMethods map[string]func(param entity.Params) (result entity.Params, err error)
}

func (server Server) handle(path string, param entity.Params) (result entity.Params, err error) {
	if fn, ok := server.mapMethods[path]; ok {
		return fn(param)
	} else {
		return nil, ErrorNotHandle
	}
}

func (server *Server) Method(path string, fn func(param entity.Params) (result entity.Params, err error)) {
	if server.mapMethods == nil {
		server.mapMethods = make(map[string]func(param entity.Params) (result entity.Params, err error))
	}
	server.mapMethods[path] = fn
}

func (server *Server) Start(host string) error {
	ctx := context.Background()
	lis, err := net.Listen("tcp", host) // any available address
	if err != nil {
		return err
	}
	serve := func(ctx context.Context, lis net.Listener, opts ...jsonrpc2.ConnOpt) error {
		for {
			conn, err := lis.Accept()
			if err != nil {
				return err
			}
			reg := regexp.MustCompile(`\d+\.\d+\.\d+\.\d+`)
			readIp := ""
			realIps := reg.FindAllString(conn.RemoteAddr().String(), -1)
			if len(realIps) > 0 {
				readIp = realIps[0]
			}
			jsonrpc2.NewConn(ctx, jsonrpc2.NewBufferedStream(conn, VarintObjectCodec{}), &RpcHandler{
				ip:           readIp,
				methodHandle: server.handle,
			}, opts...)
		}
	}
	go func() {
		if err = serve(ctx, lis); err != nil {
			if !strings.HasSuffix(err.Error(), "use of closed network connection") {
			}
		} else {
		}
	}()
	return nil
}

func (server *Server) BroadCast(method string, param entity.Params) {
	GetRpcHandlerManager().Broadcast(method, param)
}
