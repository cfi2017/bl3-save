package main

import (
	"fmt"
	"os"

	"github.com/cfi2017/bl3-save/internal/pb"
	profile2 "github.com/cfi2017/bl3-save/pkg/profile"
	"github.com/spf13/cobra"
)

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

var (
	trueAddr       = true
	customizations = []string{
		"/Game/PatchDLC/Customizations/PlayerCharacters/_Customizations/Operative/Heads/CustomHead_Operative_27.CustomHead_Operative_27",
		"/Game/PatchDLC/Customizations/PlayerCharacters/_Customizations/Gunner/Heads/CustomHead_Gunner_27.CustomHead_Gunner_27",
		"/Game/PatchDLC/Customizations/PlayerCharacters/_Customizations/Beastmaster/Heads/CustomHead_Beastmaster_27.CustomHead_Beastmaster_27",
		"/Game/PatchDLC/Customizations/PlayerCharacters/_Customizations/SirenBrawler/Heads/CustomHead_Siren_27.CustomHead_Siren_27",
		"/Game/PlayerCharacters/_Customizations/Operative/Skins/CustomSkin_Operative_42.CustomSkin_Operative_42",
		"/Game/PlayerCharacters/_Customizations/Gunner/Skins/CustomSkin_Gunner_42.CustomSkin_Gunner_42",
		"/Game/PlayerCharacters/_Customizations/Beastmaster/Skins/CustomSkin_Beastmaster_42.CustomSkin_Beastmaster_42",
		"/Game/PlayerCharacters/_Customizations/SirenBrawler/Skins/CustomSkin_Siren_42.CustomSkin_Siren_42",
	}
)

var rootCmd = &cobra.Command{
	Use:     "bl3-pax-skins",
	Short:   "Unlocks pax skins.",
	Long:    "Unlocks all the pax skins for the given profile. Needs an input (profile.sav) and an output. Provided with no guarantees. Make backups.",
	Args:    cobra.ExactArgs(2),
	Example: "bl3-pax-skins <infile> <outfile>",
	RunE: func(cmd *cobra.Command, args []string) error {
		in, err := os.Open(args[0])
		if err != nil {
			return err
		}
		s, p := profile2.Deserialize(in)
		err = in.Close()
		if err != nil {
			return err
		}

		for i := range customizations {
			found := false
			for _, c := range p.UnlockedCustomizations {
				if *c.CustomizationAssetPath == customizations[i] {
					found = true
				}
			}
			if !found {
				p.UnlockedCustomizations = append(p.UnlockedCustomizations, &pb.OakCustomizationSaveGameData{
					IsNew:                  &trueAddr,
					CustomizationAssetPath: &customizations[i],
				})
			}
		}

		out, err := os.Create(args[1])
		if err != nil {
			return err
		}
		profile2.Serialize(out, s, p)
		err = out.Close()
		if err != nil {
			return err
		}
		return nil
	},
}
