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
type CompressObjectCodec struct{}

// WriteObject implements ObjectCodec.
func (CompressObjectCodec) WriteObject(stream io.Writer, obj interface{}) error {
	data, err := json.Marshal(obj)
	if err != nil {
		return err
	}
	data = utils.Compress(data)
	var buf [binary.MaxVarintLen64]byte
	b := binary.PutUvarint(buf[:], uint64(len(data)))
	if _, err := stream.Write(buf[:b]); err != nil {
		return err
	}
	if _, err := stream.Write(data); err != nil {
		return err
	}
	return nil
}

// ReadObject implements ObjectCodec.
func (CompressObjectCodec) ReadObject(stream *bufio.Reader, v interface{}) error {
	b, err := binary.ReadUvarint(stream)
	if err != nil {
		return err
	}
	data := make([]byte, int64(b))
	_, err = stream.Read(data)
	if err != nil {
		return err
	}
	data = utils.UnCompress(data)
	return json.Unmarshal(data, v)
}
