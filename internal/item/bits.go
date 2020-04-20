package item

import (
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"strconv"
)

type Reader struct {
	stream string
}

var ErrOutOfRange = errors.New("error: out of range")

func (r *Reader) ReadInt(n int) (uint64, error) {
	if len(r.stream) < n {
		return 0, ErrOutOfRange
	}
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

type Writer string

func (w *Writer) WriteInt(v uint64, n int) error {
	if float64(v) >= math.Pow(2, float64(n)) {
		return errors.New("invalid value exceeds requested length")
	}
	bs := make([]byte, 8)
	binary.BigEndian.PutUint64(bs, v)
	var text string
	for _, b := range bs {
		text += fmt.Sprintf("%08b", b)
	}
	*w = *w + Writer(text[len(text)-n:])
	return nil
}

func (w *Writer) GetBytes() []byte {
	bs := make([]byte, 0)
	padding := 8 - len(*w)%8
	if padding != 8 {
		for i := 0; i < padding; i++ {
			*w = "0" + *w
		}
	}
	for i := len(*w)/8 - 1; i > -1; i-- {
		i, err := strconv.ParseUint(string((*w)[i*8:i*8+8]), 2, 8)
		if err != nil {
			panic(err)
		}
		bs = append(bs, byte(i))
	}
	return bs
}

func NewWriter(initial string) *Writer {
	w := Writer(initial)
	return &w
}
