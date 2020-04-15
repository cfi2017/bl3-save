package cmd

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/cfi2017/bl3-save/internal/item"
	"github.com/spf13/cobra"
)

type DigitalMarineItem struct {
	CopyType       string   `json:"copyType"`
	Level          int      `json:"level"`
	Blueprint      string   `json:"blueprint"`
	Balance        string   `json:"balance"`
	Manufacturer   string   `json:"manufacturer"`
	ComponentNames []string `json:"componentNames"`
	Components     []int    `json:"components"`
}

// deserializeCmd represents the deserialize command
var convertCmd = &cobra.Command{
	Use:   "convert",
	Short: "Convert an item from gibbed to digital_marine or vice versa",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		bs, err := base64.StdEncoding.DecodeString(args[0])
		if err != nil {
			panic(err)
		}
		var dmi DigitalMarineItem
		err = json.Unmarshal(bs, &dmi)
		if err != nil {
			// try deserializing item
			i, err := item.Deserialize(bs)
			if err != nil {
				panic(err)
			}
			// convert to dm item
			bs, err = json.Marshal(gibbedToDm(i))
			if err != nil {
				panic(err)
			}
			fmt.Print(base64.StdEncoding.EncodeToString(bs))
			return
		}
		i := dmToGibbed(dmi)
		bs, err = item.Serialize(i, 0) // encrypt with 0 seed
		if err != nil {
			panic(err)
		}
		fmt.Print(base64.StdEncoding.EncodeToString(bs))
	},
}

func dmToGibbed(dmi DigitalMarineItem) item.Item {
	i := item.Item{}
	db := item.GetDB()
	btik := item.GetBtik()
	i.Balance = dmKeyToInvKey(dmi.Balance, db.GetData("InventoryBalanceData").Assets)
	i.Manufacturer = dmKeyToInvKey(dmi.Manufacturer, db.GetData("ManufacturerData").Assets)
	i.Level = dmi.Level
	k := btik[strings.ToLower(i.Balance)]
	for _, i2 := range dmi.Components {
		i.Parts = append(i.Parts, dmKeyToInvKey(dmi.ComponentNames[i2], db.GetData(k).Assets))
	}
	i.Version = 55
	i.InvData = dmKeyToInvKey(strings.Split(dmi.Blueprint, " ")[1], db.GetData("InventoryData").Assets)
	return i
}

func getBlueprint(key, invdata string) string {
	key = strings.Replace(key, "Part", "", 1)
	parts := strings.Split(key, "_")
	parts[1], parts[2] = parts[2], parts[1]
	key = strings.Join(parts, "_")
	return key + " " + invdata
}

func gibbedToDm(i item.Item) DigitalMarineItem {
	m := DigitalMarineItem{}
	m.Manufacturer = getPartSuffix(i.Manufacturer)
	m.Level = i.Level
	m.CopyType = "item"
	m.Balance = getPartSuffix(i.Balance)
	btik := item.GetBtik()
	key := btik[strings.ToLower(i.Balance)]
	m.Blueprint = getBlueprint(key, getPartSuffix(i.InvData))
	for _, part := range i.Parts {
		p := getPartSuffix(part)
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

func dmKeyToInvKey(key string, assets []string) string {
	for _, a := range assets {
		if strings.HasSuffix(a, key) {
			return a
		}
	}
	return ""
}

func getPartSuffix(part string) string {
	p := strings.Split(part, "/")
	return p[len(p)-1]
}

func init() {
	rootCmd.AddCommand(convertCmd)
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// deserializeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// deserializeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
