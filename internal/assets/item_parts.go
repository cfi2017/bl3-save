package assets

type PartsDatabase map[string]Data

func (db PartsDatabase) GetInventoryData() Data {
	return db["InventoryData"]
}

func (db PartsDatabase) GetInventoryBalanceData() Data {
	return db["InventoryBalanceData"]
}

func (db PartsDatabase) GetManufacturerData() Data {
	return db["ManufacturerData"]
}

func (db PartsDatabase) GetData(key string) Data {
	return db[key]
}

type Versions []Version

type Version struct {
	Version uint64 `json:"version"`
	Bits    int    `json:"bits"`
}

type Data struct {
	Versions Versions `json:"versions"`
	Assets   Assets   `json:"assets"`
}

func (d Data) GetBits(version uint64) int {
	curr := d.Versions[0].Bits
	for _, v := range d.Versions {
		if v.Version > version {
			return curr
		} else if version >= v.Version {
			curr = v.Bits
		}
	}
	return curr
}
func (d Data) GetPart(i uint64) string {
	return d.Assets[i]
}

type Assets []string
