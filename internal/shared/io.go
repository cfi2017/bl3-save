package shared

import (
	"encoding/binary"
	"encoding/hex"
	"io"
)

func ReadInt(r io.Reader) int {
	bs := ReadNBytes(r, 4)
	return int(binary.LittleEndian.Uint32(bs))
}

func WriteInt(w io.Writer, i int) {
	bs := make([]byte, 4)
	binary.LittleEndian.PutUint32(bs, uint32(i))
	WriteBytes(w, bs)
}

func ReadShort(r io.Reader) int {
	bs := ReadNBytes(r, 2)
	return int(binary.LittleEndian.Uint16(bs))
}

func WriteShort(w io.Writer, i int) {
	bs := make([]byte, 2)
	binary.LittleEndian.PutUint16(bs, uint16(i))
	WriteBytes(w, bs)
}

func ReadGuid(r io.Reader) string {
	bs := ReadNBytes(r, 16)
	return hex.EncodeToString(bs)
}

func WriteGuid(w io.Writer, guid string) {
	bs, err := hex.DecodeString(guid)
	if err != nil {
		panic(err)
	}
	WriteBytes(w, bs)
}

func ReadString(r io.Reader) string {
	l := ReadInt(r)
	if l <= 1 {
		return ""
	}
	bs := ReadNBytes(r, l)
	// trim last byte (0-byte)
	return string(bs[:len(bs)-1])
}

func WriteString(w io.Writer, s string) {
	WriteInt(w, len(s)+1)
	bs := append([]byte(s), 0x00)
	WriteBytes(w, bs)
}

func ReadNBytes(r io.Reader, n int) []byte {
	bs := make([]byte, n)
	if _, err := io.ReadFull(r, bs); err != nil {
		panic("couldn't read specified byte count")
	}
	return bs
}

func WriteBytes(w io.Writer, bs []byte) {
	if _, err := w.Write(bs); err != nil {
		panic(err)
	}
}
