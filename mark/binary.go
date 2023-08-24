package mark

import (
	"runtime"
)

const BitSize = 8

func NewBinaryMark(size uint64) *BinaryMark {
	f := &BinaryMark{}
	bitNum := size / BitSize
	if size%BitSize != 0 {
		bitNum += 1
	}
	f.memSize = size
	f.useSize = size
	f.bytes = make([]byte, bitNum)
	return f
}
func NewBinaryMarkBytes(data []byte) *BinaryMark {
	f := &BinaryMark{}
	size := uint64(len(data) * BitSize)
	bitNum := size / BitSize
	if size%BitSize != 0 {
		bitNum += 1
	}
	f.memSize = size
	f.useSize = size
	f.bytes = make([]byte, bitNum)
	copy(f.bytes, data)
	return f
}

type BinaryMark struct {
	bytes   []byte
	memSize uint64
	useSize uint64
}

func (f *BinaryMark) Bytes() []byte {
	return f.bytes
}

func (f *BinaryMark) Release() {
	f.bytes = nil
	f.memSize = 0
	f.useSize = 0
	runtime.GC() // free memory as soon as possible
}
func (f *BinaryMark) Set(index uint64) {
	if f.useSize < index {
		return
	}
	bitIndex := index / BitSize
	value := &f.bytes[bitIndex]
	*value = *value | (1 << (index % BitSize))
}

func (f *BinaryMark) Get(index uint64) bool {
	if f.useSize < index {
		return false
	}
	bitIndex := index / BitSize
	value := f.bytes[bitIndex]
	return value&(1<<(index%BitSize)) != 0
}

func (f *BinaryMark) Clean(index uint64) {
	if f.useSize < index {
		return
	}
	bitIndex := index / BitSize
	value := f.bytes[bitIndex]
	value = value ^ (value & (1 << index))
}

func (f *BinaryMark) Reset(size uint64) {
	f.useSize = size
	f.bytes = make([]byte, size)
}
