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
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

var (
	dockerHost string
	tls        bool
	tlscacert  string
	tlscert    string
	tlskey     string

	//sourceFiles []string
	source   string
	authFile string
	target   string
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "Docker Images Migration Tool",
	Short: "Docker Images Migration Tool Used for Tar Images and Push to a New Registry",
	Long: `Docker Images Migration Tool Used for Tar Images and Push to a New Registry

Registry Auth File Would Be Like:

{
  "username": "hello",
  "password": "world"
}

Images Source File Would Be Like:

{
  "version": "1.0.0",
  "images": [
    "nginx:latest"
  ],
  "targetRegistryDomain": "acs-reg.sqa.alipay.net"
}


	`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports Persistent Flags, which, if defined here,
	// will be global for your application.

	//RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.acs-images.yaml)")
	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	//RootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	RootCmd.PersistentFlags().StringVarP(&dockerHost, "host", "H", "", "Daemon socket(s) to connect to")
	RootCmd.PersistentFlags().BoolVar(&tls, "tls", false, "Use TLS; implied by --tlsverify")
	RootCmd.PersistentFlags().StringVar(&tlscacert, "tlscacert", "~/.docker/ca.pem", "Trust certs signed only by this CA")
	RootCmd.PersistentFlags().StringVar(&tlscert, "tlscert", "~/.docker/cert.pem", "Path to TLS certificate file")
	RootCmd.PersistentFlags().StringVar(&tlskey, "tlskey", "~/.docker/key.pem", "Path to TLS key file")

	RootCmd.PersistentFlags().StringVarP(&authFile, "auth", "a", "", "Registry Auth File")
	//sourceFiles = *RootCmd.PersistentFlags().StringArrayP("source", "s", []string{"./images.json"}, "Images Source File")
	RootCmd.PersistentFlags().StringVarP(&target, "target", "t", "./images-tar.tar", "Images Tar File")
	RootCmd.PersistentFlags().StringVarP(&source, "source", "s", "./images.json", "Images Source File")

}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" { // enable ability to specify config file via flag
		viper.SetConfigFile(cfgFile)
	}

	viper.SetConfigName(".docker-images-migration") // name of config file (without extension)
	viper.AddConfigPath("$HOME")                    // adding home directory as first search path
	viper.AutomaticEnv()                            // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
