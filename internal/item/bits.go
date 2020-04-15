package item

import (
	"errors"
	"fmt"
	"math"
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

type Writer string

func (w *Writer) WriteInt(v, n int) error {
	if float64(v) >= math.Pow(2, float64(n)) {
		return errors.New("invalid value exceeds requested length")
	}
	r := fmt.Sprintf("%0"+strconv.Itoa(n)+"b", v)
	*w = *w + Writer(r)
	return nil
}

func (w *Writer) GetBytes() []byte {
	bs := make([]byte, 0)
	for i := len(*w); i >= 8; i -= 8 {
		p := (*w)[i-8 : i]
		i, err := strconv.ParseUint(string(p), 2, 8)
		if err != nil {
			panic(err)
		}
		bs = append(bs, byte(i))
	}
	if len(*w)%8 > 0 {
		p := (*w)[:len(*w)%8]
		i, err := strconv.ParseInt(string(p), 2, 8)
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
