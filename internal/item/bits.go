package item

import (
	"fmt"
	"strconv"
)

type Reader struct {
	stream string
}

func (r *Reader) ReadInt(n int) (uint64, error) {
	val, err := strconv.ParseUint(r.stream[len(r.stream)-n:], 2, 64)
	r.stream = r.stream[:len(r.stream)-n]
	return val, err
}

func (r *Reader) Overflow() string {
	return r.stream
}

func (r *Reader) Remaining() int {
	return len(r.stream)
}

func NewReader(data []byte) *Reader {
	r := &Reader{}
	for i := len(data) - 1; i >= 0; i-- {
		r.stream += fmt.Sprintf("%08b", data[i])
	}
	return r
}
