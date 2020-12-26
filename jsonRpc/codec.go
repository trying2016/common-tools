package jsonRpc

import (
	"bufio"
	"encoding/json"
	"io"
)

// VarintObjectCodec reads/writes JSON-RPC 2.0 objects with a varint
// header that encodes the byte length.
type VarintObjectCodec struct{}

// WriteObject implements ObjectCodec.
func (VarintObjectCodec) WriteObject(stream io.Writer, obj interface{}) error {
	data, err := json.Marshal(obj)
	if err != nil {
		return err
	}
	data = append(data, '\n')
	if _, err := stream.Write(data); err != nil {
		return err
	}
	return nil
}

// ReadObject implements ObjectCodec.
func (VarintObjectCodec) ReadObject(stream *bufio.Reader, v interface{}) error {
	body, err := bufio.NewReader(stream).ReadBytes('\n')
	//body, _, err := stream.ReadLine()
	if err != nil {
		return err
	}
	//fmt.Printf("read:%v", string(body))
	return json.Unmarshal(body, v)
	//return json.NewDecoder(io.LimitReader(stream, int64(b))).Decode(v)
}
