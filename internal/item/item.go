package item

import (
	"encoding/binary"
	"errors"
	"fmt"
	"hash/crc32"
	"log"
	"strings"

	"github.com/cfi2017/bl3-save/internal/assets"
	"github.com/cfi2017/bl3-save/pkg/pb"
)

var (
	debug bool
)

type Item struct {
	Level        int                              `json:"level"`
	Balance      string                           `json:"balance"`
	Manufacturer string                           `json:"manufacturer"`
	InvData      string                           `json:"inv_data"`
	Parts        []string                         `json:"parts"`
	Generics     []string                         `json:"generics"`
	Overflow     string                           `json:"overflow"`
	Version      uint64                           `json:"version"`
	Wrapper      *pb.OakInventoryItemSaveGameData `json:"wrapper"`
}

func DecryptSerial(data []byte) ([]byte, error) {
	if len(data) < 5 {
		return nil, errors.New("invalid serial length")
	}
	if data[0] != 0x03 {
		return nil, errors.New("invalid serial")
	}
	seed := int32(binary.BigEndian.Uint32(data[1:])) // next four bytes of serial are bogo seed
	decrypted := bogoDecrypt(seed, data[5:])
	crc := binary.BigEndian.Uint16(decrypted)                          // first two bytes of decrypted data are crc checksum
	combined := append(append(data[:5], 0xFF, 0xFF), decrypted[2:]...) // combined data with checksum replaced with 0xFF to compute checksum
	computedChecksum := crc32.ChecksumIEEE(combined)
	check := uint16(((computedChecksum) >> 16) ^ ((computedChecksum & 0xFFFF) >> 0))

	if crc != check {
		return nil, errors.New("checksum failure in packed data")
	}

	return decrypted[2:], nil
}

func EncryptSerial(data []byte, seed int32) ([]byte, error) {
	prefix := []byte{0x03}
	seedBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(seedBytes, uint32(seed))
	prefix = append(prefix, seedBytes...)
	prefix = append(prefix, 0xFF, 0xFF)
	data = append(prefix, data...)
	crc := crc32.ChecksumIEEE(data)
	checksum := ((crc >> 16) ^ crc) & 0xFFFF
	sumBytes := make([]byte, 2)
	binary.BigEndian.PutUint16(sumBytes, uint16(checksum))
	data[5], data[6] = sumBytes[0], sumBytes[1] // set crc

	return append(append([]byte{0x03}, seedBytes...), bogoEncrypt(seed, data[5:])...), nil

}

func bogoEncrypt(seed int32, data []byte) []byte {
	if seed == 0 {
		return data
	}

	steps := int(seed&0x1F) % len(data)
	data = append(data[steps:], data[:steps]...)
	return xor(seed, data)
}

func GetSeedFromSerial(data []byte) (int32, error) {
	if len(data) < 5 {
		return 0, errors.New("invalid serial length")
	}
	return int32(binary.BigEndian.Uint32(data[1:])), nil
}

func bogoDecrypt(seed int32, data []byte) []byte {
	if seed == 0 {
		return data
	}

	data = xor(seed, data)
	steps := int(seed&0x1F) % len(data)
	return append(data[len(data)-steps:], data[:len(data)-steps]...)
}

func xor(seed int32, data []byte) []byte {
	x := uint64(seed>>5) & 0xFFFFFFFF
	// target 4248340707
	for i := range data {
		x = (x * 0x10A860C1) % 0xFFFFFFFB
		data[i] = byte((uint64(data[i]) ^ x) & 0xFF)
	}
	return data
}

func Deserialize(data []byte) (item Item, err error) {
	data, err = DecryptSerial(data)
	if err != nil {
		return
	}

	r := NewReader(data)
	num := readNBits(r, 8)
	if num != 128 {
		err = errors.New("value should be 128")
		return
	}

	item.Version = readNBits(r, 7)

	balanceBits := getBits("InventoryBalanceData", item.Version)
	invDataBits := getBits("InventoryData", item.Version)
	manBits := getBits("ManufacturerData", item.Version)

	if debug {
		log.Printf("Got version: %v - balance bits: %v, invdata bits: %v, man bits: %v\n",
			item.Version, balanceBits, invDataBits, manBits,
		)
	}

	item.Balance = getPart("InventoryBalanceData", readNBits(r,
		balanceBits)-1)
	item.InvData = getPart("InventoryData", readNBits(r,
		invDataBits)-1)
	item.Manufacturer = getPart("ManufacturerData", readNBits(r,
		manBits)-1)
	item.Level = int(readNBits(r, 7))

	if k, e := assets.GetBtik()[strings.ToLower(item.Balance)]; e {
		bits := getBits(k, item.Version)
		partCount := int(readNBits(r, 6))
		item.Parts = make([]string, partCount)
		for i := 0; i < partCount; i++ {
			item.Parts[i] = getPart(k, readNBits(r, bits)-1)
		}
		genericCount := readNBits(r, 4)
		item.Generics = make([]string, genericCount)
		bits = getBits("InventoryGenericPartData", item.Version)
		for i := 0; i < int(genericCount); i++ {
			// looks like the bits are the same
			// for all the parts and generics
			item.Generics[i] = getPart("InventoryGenericPartData", readNBits(r, bits)-1)
		}
		item.Overflow = r.Overflow()

	} else {
		err = errors.New(fmt.Sprintf("unknown category %s, skipping part introspection", item.Balance))
	}

	return
}

func getBits(k string, v uint64) int {
	return assets.GetDB().GetData(k).GetBits(v)
}

func Serialize(item Item, seed int32) ([]byte, error) {
	w := NewWriter(item.Overflow)
	var err error

	// how many bits for each generic part?
	bits := getBits("InventoryGenericPartData", item.Version)

	// write each generic, bottom to top
	for i := len(item.Generics) - 1; i >= 0; i-- {
		index := getIndexFor("InventoryGenericPartData", item.Generics[i]) + 1
		err := w.WriteInt(uint64(index), bits)
		if err != nil {
			log.Printf("tried to fit index %v into %v bits for %s", index, bits, item.Generics[i])
			return nil, err
		}
	}
	// write generic count
	err = w.WriteInt(uint64(len(item.Generics)), 4)
	if err != nil {
		return nil, err
	}
	if k, e := assets.GetBtik()[strings.ToLower(item.Balance)]; e {
		// how many bits per part?
		bits = getBits(k, item.Version)
		// write each part, bottom to top
		for i := len(item.Parts) - 1; i >= 0; i-- {
			err := w.WriteInt(uint64(getIndexFor(k, item.Parts[i]))+1, bits)
			if err != nil {
				return nil, err
			}
		}
		// write part count
		err = w.WriteInt(uint64(len(item.Parts)), 6)
		if err != nil {
			return nil, err
		}
	}

	err = w.WriteInt(uint64(item.Level), 7)
	if err != nil {
		return nil, err
	}

	manIndex := getIndexFor("ManufacturerData", item.Manufacturer) + 1
	manBits := getBits("ManufacturerData", item.Version)
	err = w.WriteInt(uint64(manIndex), manBits)
	if err != nil {
		return nil, err
	}
	invIndex := getIndexFor("InventoryData", item.InvData) + 1
	invBits := getBits("InventoryData", item.Version)
	err = w.WriteInt(uint64(invIndex), invBits)
	if err != nil {
		return nil, err
	}
	balanceIndex := getIndexFor("InventoryBalanceData", item.Balance) + 1
	balanceBits := getBits("InventoryBalanceData", item.Version)
	err = w.WriteInt(uint64(balanceIndex), balanceBits)
	if err != nil {
		return nil, err
	}

	err = w.WriteInt(item.Version, 7)
	if err != nil {
		return nil, err
	}

	err = w.WriteInt(128, 8)
	if err != nil {
		return nil, err
	}

	return EncryptSerial(w.GetBytes(), seed)

}

func getIndexFor(k string, v string) int {
	for i, asset := range assets.GetDB().GetData(k).Assets {
		if asset == v {
			return i
		}
	}
	panic("no asset found while serializing")
}

func getPart(key string, index uint64) string {
	data := assets.GetDB().GetData(key)
	if int(index) >= len(data.Assets) {
		return ""
	}
	return data.GetPart(index)
}

func readNBits(r *Reader, n int) uint64 {
	i, err := r.ReadInt(n)
	if err != nil {
		panic(err)
	}
	return i
}
