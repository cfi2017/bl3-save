/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/cfi2017/bl3-save/internal/pb"
	"github.com/cfi2017/bl3-save/internal/shared"
	character2 "github.com/cfi2017/bl3-save/pkg/character"
	profile2 "github.com/cfi2017/bl3-save/pkg/profile"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

var (
	format       string
	isProfile    bool
	outputFormat string
)

// deserializeCmd represents the deserialize command
var deserializeCmd = &cobra.Command{
	Use:   "deserialize",
	Short: "Deserialize a .sav file.",
	Long: `Deserialize a .sav file.

Tries to best-guess the sav file format (profile or character) based on the files name.
Override with --format <profile|character>
`,
	Args: cobra.ExactArgs(1),
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if format != "" {
			if format == "profile" {
				isProfile = true
			} else if format == "character" {
				isProfile = false
			} else {
				return errors.New("unknown format option")
			}
		} else {
			isProfile = shared.GuessIsProfileSav(args[0])
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		f, err := os.Open(args[0])
		if err != nil {
			panic(err)
		}

		if isProfile {
			s, p := profile2.Deserialize(f)
			r := struct {
				Sav     shared.SavFile
				Profile pb.Profile
			}{s, p}
			bs, err := yaml.Marshal(r)
			if err != nil {
				panic(err)
			}
			fmt.Print(string(bs))
		} else {
			s, c := character2.Deserialize(f)
			r := struct {
				Sav       shared.SavFile
				Character pb.Character
			}{s, c}
			bs, err := yaml.Marshal(r)
			if err != nil {
				panic(err)
			}
			fmt.Print(string(bs))
		}

	},
}

func init() {
	rootCmd.AddCommand(deserializeCmd)
	deserializeCmd.PersistentFlags().StringVarP(&format, "format", "f", "", "format <isProfile|character>")
	deserializeCmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", "json", "output <json>")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// deserializeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// deserializeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
