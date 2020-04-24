package cmd

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"

	"github.com/cfi2017/bl3-save-core/pkg/assets"
	"github.com/cfi2017/bl3-save-core/pkg/item"
	"github.com/spf13/cobra"
)

var (
	literals []string
	files    []string
)

// deserializeCmd represents the deserialize command
var ConvertCmd = &cobra.Command{
	Use:   "convert",
	Short: "Convert an item from gibbed to digital_marine or vice versa",
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		queue := make(chan string)
		done := sync.WaitGroup{}

		done.Add(1)
		go func() {
			for literal := range queue {
				c, err := convert(literal)
				if err != nil {
					cmd.PrintErr(err)
					return
				}
				fmt.Println(c)
			}
			done.Done()
		}()

		for _, literal := range literals {
			queue <- literal
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
				queue <- scanner.Text()
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
		close(queue)
		done.Wait()
	},
}

type Anointments struct {
	CopyType       string
	ComponentNames []string
	Components     []int
}

func convert(arg string) (string, error) {
	anoints := make([]string, 0)
	if parts := strings.Split(arg, " "); len(parts) > 1 {
		var anointments Anointments
		bs, err := base64.StdEncoding.DecodeString(parts[1])
		if err != nil {
			return "", err
		}
		err = json.Unmarshal(bs, &anointments)
		if anointments.CopyType != "anointment" {
			return "", errors.New("not a valid anointment code")
		}
		for _, i := range anointments.Components {
			anoints = append(anoints, item.DmKeyToInvKey(anointments.ComponentNames[i],
				assets.GetDB().GetData("InventoryGenericPartData").Assets))
		}
		arg = parts[0]
	}
	arg = strings.TrimPrefix(arg, "bl3(")
	arg = strings.TrimPrefix(arg, "BL3(")
	arg = strings.TrimSuffix(arg, ")")
	bs, err := base64.StdEncoding.DecodeString(arg)
	if err != nil {
		panic(err)
	}
	var dmi item.DigitalMarineItem
	err = json.Unmarshal(bs, &dmi)
	if err != nil {
		// try deserializing item
		i, err := item.Deserialize(bs)
		if err != nil {
			return "", errors.New("couldn't deserialize item")
		}
		// convert to dm item
		bs, err = json.Marshal(item.GibbedToDm(i))
		if err != nil {
			return "", errors.New("couldn't convert item to dm format")
		}
		return base64.StdEncoding.EncodeToString(bs), nil
	}
	i := item.DmToGibbed(dmi)
	// we don't check for existing anoints at the moment,
	// nor anoint count (todo: add sanity checks)
	i.Generics = append(i.Generics, anoints...)
	bs, err = item.Serialize(i, 0) // encrypt with 0 seed
	if err != nil {
		return "", errors.New("error decoding item. invalid serial")
	}
	return fmt.Sprintf("bl3(%s)", base64.StdEncoding.EncodeToString(bs)), nil
}

func init() {
	rootCmd.AddCommand(ConvertCmd)

	ConvertCmd.PersistentFlags().StringSliceVar(&literals, "from-literal", []string{}, "literal code inputs")
	ConvertCmd.PersistentFlags().StringSliceVar(&files, "from-file", []string{}, "input files")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// deserializeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// deserializeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
