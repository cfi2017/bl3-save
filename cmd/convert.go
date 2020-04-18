package cmd

import (
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/cfi2017/bl3-save/internal/item"
	"github.com/spf13/cobra"
)

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
		var dmi item.DigitalMarineItem
		err = json.Unmarshal(bs, &dmi)
		if err != nil {
			// try deserializing item
			i, err := item.Deserialize(bs)
			if err != nil {
				panic(err)
			}
			// convert to dm item
			bs, err = json.Marshal(item.GibbedToDm(i))
			if err != nil {
				panic(err)
			}
			fmt.Print(base64.StdEncoding.EncodeToString(bs))
			return
		}
		i := item.DmToGibbed(dmi)
		bs, err = item.Serialize(i, 0) // encrypt with 0 seed
		if err != nil {
			panic(err)
		}
		fmt.Printf("bl3(%s)", base64.StdEncoding.EncodeToString(bs))
	},
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
