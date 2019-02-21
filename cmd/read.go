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
	"encoding/hex"
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/spf13/cobra"
)

var outFile string
var binLen int
var icType string

// readCmd represents the read command
var readCmd = &cobra.Command{
	Use:   "read",
	Short: "Reads the EEPROM content and stores it in the specified --output-file",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if outFile == "" {
			errorMsg := "You must supply the --output-file flag."
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

		if binLen == 0 {
			binLen = 256
		}
		return nil
	},

	RunE: func(cmd *cobra.Command, args []string) error {
		buf := make([]byte, binLen)
		var err error
		arduino := NewArduino93L56R(serPort)
		if err := arduino.Connect(); err != nil {
			return err
		}
		defer arduino.Close()

		if icType == "microwire" {
			buf, err = arduino.Read(eepromAddr, binLen)
		}

		if icType == "i2c" {
			buf, err = arduino.I2CRead(eepromAddr, binLen)
		}

		if err != nil {
			return err
		}

		err = ioutil.WriteFile(outFile, buf, 0644)
		if err != nil {
			return fmt.Errorf("Unable to save EEPROM contents to file %s. Error: %s", outFile, err)
		}

		fmt.Println(hex.Dump(buf))

		return nil
	},
}

func init() {
	rootCmd.AddCommand(readCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// readCmd.PersistentFlags().String("foo", "", "A help for foo")
	readCmd.Flags().StringVar(&outFile, "output-file", "", "A file to store the contents read from the EEPROM")
	readCmd.Flags().IntVar(&binLen, "read-length", 256, "The number of bytes to read from the EEPROM. Default is 256")
	readCmd.Flags().StringVar(&icType, "type", "", "The type of EEPROM you're trying to read. One of: microwire, i2c")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// readCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
