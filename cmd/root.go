// Copyright © 2019 NAME HERE <EMAIL ADDRESS>
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
	"errors"
	"fmt"
	"math"
	"os"

	homedir "github.com/mitchellh/go-homedir"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var logger *log.Logger
var cfgFile string
var serPort string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "93l56r-cli",
	Short: "Commandline utility for reading and writing 93L56R EEPROMs using an Arduino",
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		logger = log.New()
		logger.SetLevel(log.InfoLevel)

		if serPort == "" {
			errorMsg := "You must supply the --serial-port flag."
			return errors.New(errorMsg)
		}
		return nil
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.93l56r-cli.yaml)")
	rootCmd.PersistentFlags().StringVar(&serPort, "serial-port", "", "Device path or name for the serial port your arduino is connected to. I.E. COM1, /dev/cu.usbmodem*")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".93l56r-cli" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".93l56r-cli")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

func MicrosecondsOnTheWireByteCount(count int) int {
	// bit_time = 1 / baud * databytes * (1 start bit, 8 bit word, 1 stop bit = 10)
	return int(math.Round(1.0 / 9600.00 * 1000000.0 * float64(count) * 10.0))
}
