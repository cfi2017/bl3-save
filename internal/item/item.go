package item

import (
	"encoding/binary"
	"encoding/hex"
	"errors"
	"hash/crc32"
	"log"
)

type Item struct {
}

func Deserialize(data []byte) (i Item, err error) {
	if len(data) < 5 || len(data) > 40 {
		err = errors.New("invalid serial length")
		return
	}
	if data[0] != 0x03 {
		err = errors.New("invalid serial")
		return
	}
	seed := binary.LittleEndian.Uint32(data[1:]) // next four bytes of serial are bogo seed
	decrypted := BogoDecrypt(seed, data[5:])
	crc := binary.LittleEndian.Uint16(decrypted)                       // first two bytes of decrypted data are crc checksum
	combined := append(append(data[:5], 0xFF, 0xFF), decrypted[2:]...) // combined data with checksum replaced with 0xFF to compute checksum
	log.Println(hex.EncodeToString(combined))
	for len(combined) < 40 {
		combined = append(combined, 0xFF)
	}
	computedChecksum := crc32.ChecksumIEEE(combined)
	check := uint16(((computedChecksum) >> 16) ^ ((computedChecksum & 0xFFFF) >> 0))

	if crc != check {
		err = errors.New("checksum failure in packed data")
		return
	}

	return
}

func BogoDecrypt(seed uint32, data []byte) []byte {
	if seed == 0 {
		return data
	}

	tmp := make([]byte, len(data))
	copy(tmp, data)

	xor := seed >> 5
	for i := 0; i < len(data); i++ {
		xor = uint32((uint64(xor) * 0x10A860C1) % 0xFFFFFFFB)
		tmp[i] ^= byte(xor & 0xFF)
	}

	rightHalf := int(seed%32) % len(data)
	leftHalf := len(data) - rightHalf

	copy(data[:rightHalf], tmp[leftHalf:leftHalf+rightHalf])
	copy(data[rightHalf:rightHalf+leftHalf], tmp[:leftHalf])
	return data
}
