package jsonRpc

import (
	"bufio"
	"encoding/binary"
	"encoding/json"
	"io"

	"github.com/trying2016/common-tools/utils"
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
	data = utils.Compress(data)
	sendData := make([]byte, len(data)+4)
	var nLen uint32 = uint32(len(data))
	var buf4 [4]byte
	binary.LittleEndian.PutUint32(buf4[:], nLen)
	copy(sendData, buf4[:])
	copy(sendData[4:], data)
	//data = append(data, '\n')
	if _, err := stream.Write(sendData); err != nil {
		return err
	}
	return nil
}

// ReadObject implements ObjectCodec.
func (VarintObjectCodec) ReadObject(stream *bufio.Reader, v interface{}) error {
	var buf4 [4]byte
	reader := bufio.NewReader(stream)
	_, err := reader.Read(buf4[:])
	if err != nil {
		return err
	}
	nLen := binary.LittleEndian.Uint32(buf4[:])
	body := make([]byte, nLen)
	//body, err := bufio.NewReader(stream).ReadBytes('\n')
	//body, _, err := stream.ReadLine()
	_, err = reader.Read(body)
	if err != nil {
		return err
	}
	body = utils.UnCompress(body)
	//fmt.Printf("read:%v", string(body))
	return json.Unmarshal(body, v)
	//return json.NewDecoder(io.LimitReader(stream, int64(b))).Decode(v)
}
