package item

import (
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"hash/crc32"
	"io/ioutil"
	"log"
	"strings"
	"sync"

	"github.com/cfi2017/bl3-save/pkg/pb"
)

var (
	db   PartsDatabase
	btik map[string]string
	once = sync.Once{}
)

type Item struct {
	Level        int
	Balance      string
	Manufacturer string
	InvData      string
	Parts        []string
	Generics     []string
	Overflow     string
	Version      uint64
	Wrapper      *pb.OakInventoryItemSaveGameData
}

func DecryptSerial(data []byte) ([]byte, error) {
	if len(data) < 5 || len(data) > 40 {
		return nil, errors.New("invalid serial length")
	}
	if data[0] != 0x03 {
		return nil, errors.New("invalid serial")
	}
	seed := int32(binary.BigEndian.Uint32(data[1:])) // next four bytes of serial are bogo seed
	decrypted := bogoDecrypt(seed, data[5:])
	crc := binary.BigEndian.Uint16(decrypted)                          // first two bytes of decrypted data are crc checksum
	combined := append(append(data[:5], 0xFF, 0xFF), decrypted[2:]...) // combined data with checksum replaced with 0xFF to compute checksum
	log.Println(hex.EncodeToString(combined))
	computedChecksum := crc32.ChecksumIEEE(combined)
	check := uint16(((computedChecksum) >> 16) ^ ((computedChecksum & 0xFFFF) >> 0))

	if crc != check {
		return nil, errors.New("checksum failure in packed data")
	}

	return decrypted[2:], nil
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

	once.Do(func() {
		btik, err = loadPartMap("balance_to_inv_key.json")
		if err != nil {
			return
		}
		db, err = loadPartsDatabase("inventory_raw.json")
	})
	if err != nil {
		return
	}

	item.Version = readNBits(r, 7)

	item.Balance = getPart("InventoryBalanceData", readNBits(r,
		getBits("InventoryBalanceData", item.Version))-1)
	item.InvData = getPart("InventoryData", readNBits(r,
		getBits("InventoryData", item.Version))-1)
	item.Manufacturer = getPart("ManufacturerData", readNBits(r,
		getBits("ManufacturerData", item.Version))-1)
	item.Level = int(readNBits(r, 7))

	if k, e := btik[strings.ToLower(item.Balance)]; e {
		bits := db.GetData(k).GetBits(item.Version)
		partCount := int(readNBits(r, 6))
		item.Parts = make([]string, partCount)
		for i := 0; i < partCount; i++ {
			item.Parts[i] = getPart(k, readNBits(r, bits)-1)
		}
		genericCount := int(readNBits(r, 4))
		item.Generics = make([]string, genericCount)
		for i := 0; i < genericCount; i++ {
			// looks like the bits are the same
			// for all the parts and generics
			item.Generics[i] = getPart(k, readNBits(r, bits)-1)
		}
		item.Overflow = r.Overflow()

	} else {
		err = errors.New(fmt.Sprintf("unknown category %s, skipping part introspection", item.Balance))
	}

	return
}

func getBits(k string, v uint64) int {
	return db.GetData(k).GetBits(v)
}

/*func Serialize(item Item) []byte {
	data := make([]byte, 0)
	buffer := item.Overflow

	if k, e := btik[strings.ToLower(item.Balance)]; e {
		for i := len(item.Generics); i <= 0; i-- {
			bits := getBits(k, item.Version)
			item.Generics[i]
		}
	}

}*/

func getPart(key string, index uint64) string {
	data := db.GetData(key)
	return data.GetPart(index)
}

func readNBits(r *Reader, n int) uint64 {
	i, err := r.ReadInt(n)
	if err != nil {
		panic(err)
	}
	return i
}

func loadPartMap(file string) (m map[string]string, err error) {
	bs, err := ioutil.ReadFile(file)
	if err != nil {
		return
	}
	err = json.Unmarshal(bs, &m)
	return
}

func loadPartsDatabase(file string) (db PartsDatabase, err error) {
	bs, err := ioutil.ReadFile(file)
	if err != nil {
		return
	}
	err = json.Unmarshal(bs, &db)
	return
}
