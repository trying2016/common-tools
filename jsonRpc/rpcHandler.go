package jsonRpc

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/trying2016/common-tools/jsonRpc/jsonrpc2"

	"github.com/tidwall/gjson"
	"github.com/trying2016/common-tools/log"
	"github.com/trying2016/common-tools/utils"
)

var rpcHandlerManagerInstance *RpcHandlerManager

func init() {
	rpcHandlerManagerInstance = &RpcHandlerManager{
		handlers: make(map[string]*RpcHandler),
	}
}

type RpcHandlerManager struct {
	handlersLock sync.RWMutex
	handlers     map[string]*RpcHandler
}

func GetRpcHandlerManager() *RpcHandlerManager {
	return rpcHandlerManagerInstance
}

func (s *RpcHandlerManager) Add(key string, handler *RpcHandler) {
	s.handlersLock.Lock()
	defer s.handlersLock.Unlock()
	s.handlers[key] = handler
}

func (s *RpcHandlerManager) Remove(key string) {
	s.handlersLock.Lock()
	defer s.handlersLock.Unlock()
	if _, ok := s.handlers[key]; ok {
		delete(s.handlers, key)
	}
}

func (s *RpcHandlerManager) Broadcast(method string, params Params) {
	s.handlersLock.Lock()
	defer s.handlersLock.Unlock()
	for _, handler := range s.handlers {
		handler.Send(method, params)
	}
}

//
type RpcHandler struct {
	conn         *jsonrpc2.Conn
	updateTime   time.Time
	key          string
	agent        string
	miner        string
	minerName    string
	ip           string
	methodHandle func(path string, param Params) (result Params, err error)
}

func (handler *RpcHandler) Send(method string, params Params) {
	if err := handler.conn.Notify(context.Background(), method, params); err != nil {
		log.Warn("Notify fail, error:%v", err)
	}
}

func (handler *RpcHandler) Handle(ctx context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request) {
	if req.Notif {
		return // notification
	}
	var params = Params{}

	if req.Params != nil {
		jsonRet := gjson.ParseBytes(*req.Params)
		for k, v := range jsonRet.Map() {
			params[strings.ToLower(k)] = v
		}
	}

	// 登录
	if req.Method == "login" {
		// 拦截登录
		_, err := handler.methodHandle(req.Method, params)
		if err != nil && err == ErrorNotHandle {
			if err := conn.Reply(ctx, req.ID, map[string]interface{}{
				"status": "error",
				"error":  err.Error(),
			}); err != nil {
				log.Warn("send keep lived Reply fail, error:%v", err)
			}
			return
		}

		handler.conn = conn
		handler.key = fmt.Sprintf("%v-%v-%v", handler.ip, time.Now().UnixNano(), utils.CreateRandomAllString(16))
		// 添加到队列
		GetRpcHandlerManager().Add(handler.key, handler)
	} else if req.Method == "keepalived" {
		handler.updateTime = time.Now()
		if err := conn.Reply(ctx, req.ID, map[string]interface{}{
			"status": "KEEPALIVED",
		}); err != nil {
			log.Warn("send keep lived Reply fail, error:%v", err)
		}
	} else if req.Method == "close" {
		GetRpcHandlerManager().Remove(handler.key)
	}

	result, err := handler.methodHandle(req.Method, params)
	if err != nil {
		if err == ErrorNotHandle {
			return
		}
		if err := conn.Reply(ctx, req.ID, map[string]interface{}{
			"status": "error",
			"error":  err.Error(),
		}); err != nil {
			log.Warn("send keep lived Reply fail, error:%v", err)
		}
	} else {
		result["status"] = "OK"
		if err := conn.Reply(ctx, req.ID, result); err != nil {
			log.Warn("send keep lived Reply fail, error:%v", err)
		}
	}
}
