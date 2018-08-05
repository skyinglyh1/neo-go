package util

import (
	"encoding/binary"
	"io"
)

type BinReader struct {
	R   io.Reader
	Err error
}

func (r *BinReader) Read(v interface{}) {
	if r.Err != nil {
		return
	}
	binary.Read(r.R, binary.LittleEndian, v)
}
func (r *BinReader) ReadBigEnd(v interface{}) {
	if r.Err != nil {
		return
	}
	binary.Read(r.R, binary.BigEndian, v)
}

// func (br *BinReader) Bytes() {
// 	by, _ := util.ReaderToBuffer(br.R)
// 	return by
// }

func (r *BinReader) VarUint() uint64 {
	var b uint8
	binary.Read(r.R, binary.LittleEndian, &b)

	if b == 0xfd {
		var v uint16
		binary.Read(r.R, binary.LittleEndian, &v)
		return uint64(v)
	}
	if b == 0xfe {
		var v uint32
		binary.Read(r.R, binary.LittleEndian, &v)
		return uint64(v)
	}
	if b == 0xff {
		var v uint64
		binary.Read(r.R, binary.LittleEndian, &v)
		return v
	}

	return uint64(b)
}

func (r *BinReader) VarBytes() []byte {
	n := r.VarUint()
	b := make([]byte, n)
	r.Read(b)
	return b
}

func (r *BinReader) VarString() string {
	b := r.VarBytes()
	return string(b)
}
