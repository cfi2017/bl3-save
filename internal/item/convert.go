package item

import "strings"

type DigitalMarineItem struct {
	CopyType       string   `json:"copyType"`
	Level          int      `json:"level"`
	Blueprint      string   `json:"blueprint"`
	Balance        string   `json:"balance"`
	Manufacturer   string   `json:"manufacturer"`
	ComponentNames []string `json:"componentNames"`
	Components     []int    `json:"components"`
}

func DmToGibbed(dmi DigitalMarineItem) Item {
	i := Item{}
	db := GetDB()
	btik := GetBtik()
	i.Balance = DmKeyToInvKey(dmi.Balance, db.GetData("InventoryBalanceData").Assets)
	i.Manufacturer = DmKeyToInvKey(dmi.Manufacturer, db.GetData("ManufacturerData").Assets)
	i.Level = dmi.Level
	k := btik[strings.ToLower(i.Balance)]
	for _, i2 := range dmi.Components {
		i.Parts = append(i.Parts, DmKeyToInvKey(dmi.ComponentNames[i2], db.GetData(k).Assets))
	}
	i.Version = 55
	i.InvData = DmKeyToInvKey(strings.Split(dmi.Blueprint, " ")[1], db.GetData("InventoryData").Assets)
	return i
}

func GetBlueprint(key, invdata string) string {
	key = strings.Replace(key, "Part", "", 1)
	parts := strings.Split(key, "_")
	parts[1], parts[2] = parts[2], parts[1]
	key = strings.Join(parts, "_")
	return key + " " + invdata
}

func GibbedToDm(i Item) DigitalMarineItem {
	m := DigitalMarineItem{}
	m.Manufacturer = GetPartSuffix(i.Manufacturer)
	m.Level = i.Level
	m.CopyType = "item"
	m.Balance = GetPartSuffix(i.Balance)
	btik := GetBtik()
	key := btik[strings.ToLower(i.Balance)]
	m.Blueprint = GetBlueprint(key, GetPartSuffix(i.InvData))
	for _, part := range i.Parts {
		p := GetPartSuffix(part)
		found := false
		for i2, name := range m.ComponentNames {
			if name == p {
				m.Components = append(m.Components, i2)
				found = true
			}
		}
		if !found {
			m.ComponentNames = append(m.ComponentNames, p)
			m.Components = append(m.Components, len(m.ComponentNames)-1)
		}
	}
	return m
}

func DmKeyToInvKey(key string, assets []string) string {
	for _, a := range assets {
		if strings.HasSuffix(a, key) {
			return a
		}
	}
	return ""
}

func GetPartSuffix(part string) string {
	p := strings.Split(part, "/")
	return p[len(p)-1]
}
