package cmd

import (
	"bufio"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"io"
	"os"
	"strings"

	"github.com/cfi2017/bl3-save-core/pkg/assets"
	"github.com/cfi2017/bl3-save-core/pkg/character"
	"github.com/cfi2017/bl3-save-core/pkg/item"
	"github.com/cfi2017/bl3-save-core/pkg/pb"
	"github.com/spf13/cobra"
)

var (
	items []string
)

var ItemsCommand = &cobra.Command{
	Use:  "items",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// first, grab the character
		f, err := os.Open(args[0])
		if err != nil {
			cmd.PrintErrf("couldn't open character: %v\n", err)
		}
		s, c, err := character.Deserialize(f, "pc")
		if err != nil {
			cmd.PrintErr(err)
			return
		}
		err = f.Close()
		if err != nil {
			cmd.PrintErr(err)
			return
		}

		for _, file := range files {
			var reader io.Reader
			var err error
			if file != "-" {
				reader, err = os.Open(file)
				if err != nil {
					cmd.PrintErr(err)
					return
				}
			} else {
				reader = os.Stdin
			}
			scanner := bufio.NewScanner(reader)
			for scanner.Scan() {
				items = append(items, scanner.Text())
			}
			if err := scanner.Err(); err != nil {
				cmd.PrintErr(err)
				return
			}
			if f, ok := reader.(*os.File); ok {
				err = f.Close()
				if err != nil {
					cmd.PrintErr(err)
					return
				}
			}
		}

		for i := range items {
			var anoints = make([]string, 0)
			if parts := strings.Split(items[i], " "); len(parts) > 1 {
				var anointments Anointments
				bs, err := base64.StdEncoding.DecodeString(parts[1])
				if err != nil {
					cmd.PrintErr(err)
					return
				}
				err = json.Unmarshal(bs, &anointments)
				if anointments.CopyType != "anointment" {
					cmd.PrintErrln("not a valid anointment code")
					return
				}
				for _, i := range anointments.Components {
					anoints = append(anoints, item.DmKeyToInvKey(anointments.ComponentNames[i],
						assets.GetDB().GetData("InventoryGenericPartData").Assets))
				}
				items[i] = parts[0]
			}
			if strings.HasPrefix(items[i], "bl3(") || strings.HasPrefix(items[i], "BL3(") {
				// assume bl3 format, verify and add
			} else {
				// try to convert base64
				bs, err := base64.StdEncoding.DecodeString(items[i])
				if err != nil {
					cmd.PrintErr(err)
					return
				}
				var dmi item.DigitalMarineItem
				err = json.Unmarshal(bs, &dmi)
				if err == nil {
					gi := item.DmToGibbed(dmi)
					bs, err = item.Serialize(gi, 0) // encrypt with 0 seed
					if err != nil {
						cmd.PrintErr(err)
						return
					}
					items[i] = hex.EncodeToString(bs)
				}
				// else assume bl3 format, do nothing for now
			}

			items[i] = strings.TrimPrefix(items[i], "bl3(")
			items[i] = strings.TrimPrefix(items[i], "BL3(")
			items[i] = strings.TrimSuffix(items[i], ")")

			bs, err := base64.StdEncoding.DecodeString(items[i])
			if err != nil {
				cmd.PrintErr(err)
				return
			}
			// we don't actually care about the item,
			// just that it deserializes correctly
			bsc := make([]byte, len(bs))
			copy(bsc, bs)
			current, err := item.Deserialize(bsc)
			if err != nil {
				cmd.PrintErr(err)
				return
			}
			current.Generics = append(current.Generics, anoints...)
			bs, err = item.Serialize(current, 0)
			if err != nil {
				cmd.PrintErr(err)
				return
			}
			c.InventoryItems = append(c.InventoryItems, &pb.OakInventoryItemSaveGameData{
				ItemSerialNumber: bs,
				PickupOrderIndex: 200, // set static, idk what this does
				Flags:            1,   // flag 1 should be "new"
				// WeaponSkinPath:      "",  // no skin applied
				DevelopmentSaveData: nil,
			})
		}

		// first, grab the character
		f, err = os.Create(args[0])
		if err != nil {
			cmd.PrintErrf("couldn't create character: %v\n", err)
		}
		character.Serialize(f, s, c, "pc")
		err = f.Close()
		if err != nil {
			cmd.PrintErr(err)
			return
		}

	},
}

func init() {
	rootCmd.AddCommand(ItemsCommand)
	ItemsCommand.PersistentFlags().StringSliceVar(&files, "from-file", []string{}, "import items from file (- for stdin)")
	ItemsCommand.PersistentFlags().StringSliceVar(&items, "from-literal", []string{}, "import item")
}
