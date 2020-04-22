package cmd

import (
	"log"
	"os"

	"github.com/cfi2017/bl3-save-core/pkg/character"
	"github.com/cfi2017/bl3-save-core/pkg/profile"
	"github.com/spf13/cobra"
)

var (
	fileType string
)

var DecryptCommand = &cobra.Command{
	Use:   "decrypt",
	Short: "decrypt a file",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		f, err := os.Open(args[0])
		if err != nil {
			panic(err)
		}
		defer f.Close()
		var d []byte
		if fileType == "character" {
			_, d = character.Decrypt(f)
		} else if fileType == "profile" {
			_, d = profile.Decrypt(f)
		} else {
			log.Fatalln("invalid file type")
		}
		_, err = os.Stdout.Write(d)
	},
}

func init() {
	rootCmd.AddCommand(DecryptCommand)
	DecryptCommand.PersistentFlags().StringVarP(&fileType, "type", "t", "character", "character|profile")
}
