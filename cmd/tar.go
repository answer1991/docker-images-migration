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

// tarCmd represents the tar command
var tarCmd = &cobra.Command{
	Use:   "tar",
	Short: "Package Images From Source File(s) to Target Tar File",
	Long:  `Package Images From Source File(s) to Target Tar File`,
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

		registryDomains, err := getTarRegistryDomains(sourceList)

		if nil != err {
			stderr(cmd, err)
			return
		}

		stdout(cmd, "Registry Domains: ", registryDomains)

		defer logoutRegistries(cmd, registryDomains)

		for _, registryDomain := range registryDomains {
			if "" == registryDomain {
				continue
			}

			stdout(cmd, fmt.Sprintf("Try to Login in %s", registryDomain))
			err := loginRegistry(cmd, authInfo, registryDomain)

			if nil != err {
				stderr(cmd, err)
			}
		}

		images, err := getImages(sourceList)

		if nil != err {
			stderr(cmd, err)
			return
		}

		err = pullImages(cmd, images)

		if nil != err {
			stderr(cmd, err)
			return
		}

		err = tarImages(cmd, images)

		if nil != err {
			stderr(cmd, err)
			return
		}

		stdout(cmd, fmt.Sprintf("Package Images Finished, Found file in %s", target))
	},
}

func init() {
	RootCmd.AddCommand(tarCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// tarCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// tarCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
