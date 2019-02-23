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
	"reflect"
	"time"

	"github.com/spf13/cobra"
)

var inFile string

// writeCmd represents the write command
var writeCmd = &cobra.Command{
	Use:   "write",
	Short: "Writes the --input-file contents to the EEPROM",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if inFile == "" {
			errorMsg := "You must supply the --input-file flag."
			return errors.New(errorMsg)
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		var maxPacketDataLen int
		var wordSizeScale int

		if icType == "microwire" {
			maxPacketDataLen = 56
			wordSizeScale = 2
		}
		if icType == "i2c" {
			maxPacketDataLen = 55
			wordSizeScale = 1
		}
		// TODO: Require a --force option to write without reading
		buf, err := ioutil.ReadFile(inFile)
		if err != nil {
			return fmt.Errorf("Unable to read the input file %s. Error: %s", inFile, err)
		}

		arduino := NewArduino93L56R(serPort)
		if err := arduino.Connect(); err != nil {
			return err
		}
		defer arduino.Close()

		start := time.Now()

		wholePacketCount := len(buf) / maxPacketDataLen
		for l := 0; l < wholePacketCount; l++ {
			startAddr := maxPacketDataLen/wordSizeScale*l + eepromAddr
			bufslice := buf[startAddr*wordSizeScale : startAddr*wordSizeScale+maxPacketDataLen]
			if err = arduino.Write(startAddr, bufslice, icType); err != nil {
				return err
			}
		}

		remainder := buf[wholePacketCount*maxPacketDataLen:]
		if err = arduino.Write(maxPacketDataLen/wordSizeScale*wholePacketCount+eepromAddr, remainder, icType); err != nil {
			return err
		}

		duration := time.Since(start)

		ver, err := arduino.Read(eepromAddr, len(buf), icType)
		if err != nil {
			return err
		}

		if !reflect.DeepEqual(ver, buf) {
			return fmt.Errorf("The content written to the EEPROM was not the same as the content of the supplied file after writing.\n\nFile:\n%s\n\nEEPROM Content:\n%s", hex.Dump(buf), hex.Dump(ver))
		} else {
			fmt.Printf("Successfully wrote file to EEPROM in %s.\n\n%s", duration, hex.Dump(ver))
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(writeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// writeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// writeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	writeCmd.Flags().StringVar(&inFile, "input-file", "", "A file to write to the EEPROM")
}
