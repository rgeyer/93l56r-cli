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

	"github.com/spf13/cobra"
)

var serPort string
var eepromAddr int
var icType string

// eepromCmd represents the eeprom command
var eepromCmd = &cobra.Command{
	Use:   "eeprom",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if serPort == "" {
			errorMsg := "You must supply the --serial-port flag."
			return errors.New(errorMsg)
		}

		switch test := icType; test {
		case "microwire":
			break
		case "i2c":
			break
		default:
			return errors.New("You must supply the --type flag, and it must be one of: microwire, i2c")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("eeprom called")
	},
}

func init() {
	rootCmd.AddCommand(eepromCmd)

	eepromCmd.PersistentFlags().StringVar(&serPort, "serial-port", "", "Device path or name for the serial port your arduino is connected to. I.E. COM1, /dev/cu.usbmodem*")
	eepromCmd.PersistentFlags().IntVar(&eepromAddr, "start-address", 0, "The starting address of the EEPROM to begin the read or write operation. Default is 0")
	eepromCmd.PersistentFlags().StringVar(&icType, "type", "", "The type of EEPROM you're trying to read. One of: microwire, i2c")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// eepromCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// eepromCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
