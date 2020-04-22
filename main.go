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
package main

import (
	"github.com/cfi2017/bl3-save-core/pkg/assets"
	"github.com/cfi2017/bl3-save/cmd"
	assets2 "github.com/cfi2017/bl3-save/internal/assets"
	"github.com/cfi2017/bl3-save/internal/server"
	"github.com/spf13/cobra"
)

var (
	version = "v100.0.0"
	commit  = ""
	date    = ""
	builtBy = "local"
)

func main() {
	assets.DefaultAssetLoader = assets2.HttpAssetsLoader{}
	server.BuildVersion = version
	server.BuildCommit = commit
	server.BuildDate = date
	server.BuiltBy = builtBy
	cobra.MousetrapHelpText = ""
	cmd.Execute()
}
