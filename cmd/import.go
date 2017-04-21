// Copyright © 2017 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// importCmd represents the import command
var importCmd = &cobra.Command{
	Use:   "import",
	Short: "Load Images from Tar File, And Push to the New Registry",
	Long:  `Load Images from Tar File, And Push to the New Registry`,
	Run: func(cmd *cobra.Command, args []string) {
		authInfo, err := loadAuth()

		if nil != err {
			stderr(cmd, err)
			return
		}

		sourceList, err := loadSource()

		if nil != err {
			stderr(cmd, err)
			return
		}

		registryDomains, err := getImportRegistryDomains(sourceList)

		if nil != err {
			stderr(cmd, err)
			return
		}

		defer logoutRegistries(cmd, registryDomains)

		for _, registryDomain := range registryDomains {
			if "" != registryDomain {
				stdout(cmd, fmt.Sprintf("Try to Login in %s", registryDomain))
				err := loginRegistry(cmd, authInfo, registryDomain)

				if nil != err {
					stderr(cmd, err)
				}
			}
		}

		loadTar(cmd)

		images, err := translateTargetImages(cmd, sourceList)
		if nil != err {
			stderr(cmd, err)
			return
		}

		err = pushImages(cmd, images)
		if nil != err {
			stderr(cmd, err)
			return
		}
	},
}

func init() {
	RootCmd.AddCommand(importCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// importCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// importCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
