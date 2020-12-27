package jsonRpc

import (
	"strings"

	"github.com/trying2016/common-tools/utils"
)

var (
	BodyType = "_body"
)

type Params map[string]interface{}

func (params Params) GetString(key string) string {
	if v, ok := params[strings.ToLower(key)]; ok {
		return utils.ToString(v)
	} else {
		return ""
	}
}

func (params Params) GetInt(key string) int {
	if v, ok := params[strings.ToLower(key)]; ok {
		return utils.ToInt(v)
	} else {
		return 0
	}
}

func (params Params) GetInt64(key string) int64 {
	if v, ok := params[strings.ToLower(key)]; ok {
		return utils.ToInt64(v)
	} else {
		return 0
	}
}

func (params Params) GetUInt64(key string) uint64 {
	if v, ok := params[strings.ToLower(key)]; ok {
		return utils.ToUint64(v)
	} else {
		return 0
	}
}
func (params Params) GetUInt32(key string) uint32 {
	if v, ok := params[strings.ToLower(key)]; ok {
		return utils.ToUint32(v)
	} else {
		return 0
	}
}

func (params Params) GetBody() []byte {
	if v, ok := params[BodyType]; ok {
		return v.([]byte)
	} else {
		return nil
	}
}

func (params Params) GetRpcHandler() *RpcHandler {
	c, ok := params["_client"]
	if ok {
		return c.(*RpcHandler)
	} else {
		return nil
	}
}
