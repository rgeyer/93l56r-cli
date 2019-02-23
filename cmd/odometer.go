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
	"fmt"

	"github.com/spf13/cobra"
)

// odometerCmd represents the odometer command
var odometerCmd = &cobra.Command{
	Use:   "odometer",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {

		target_mileage := 157889

		fmt.Printf("Encoding value: %d ...\n", target_mileage)

		repeat_count := target_mileage & 0x0F // Take lowest nibble and use it as the repeat count
		new_value := target_mileage >> 4      // The new count has that lowest nibble shaved off
		old_value := new_value - 1            // The old count will be new count - 1

		data_buffer := make([]byte, 0x20) // use this as a workspace

		encode_table := []byte{0x00, 0x07, 0x0C, 0x0B, 0x06, 0x01, 0x0A, 0x0D, 0x03, 0x04, 0x0F, 0x08, 0x05, 0x02, 0x09, 0x0E}
		// Encoding table for the least significant 4 bits..
		// as each value changes, 3 bits change and 1 stays the same. This is probably some wear leveling scheme for the EEPROM

		new_value = (new_value & 0xFFF0) + int(encode_table[new_value&0x0F]) // do encodings on both values
		old_value = (old_value & 0xFFF0) + int(encode_table[old_value&0x0F])

		data_buffer[0] = byte(new_value >> 8) // Store new_value in the first slot
		data_buffer[1] = byte(new_value & 0xFF)

		for count := 0; count < repeat_count; count++ { // Store new_value in more slots determined by repeat_count
			data_buffer[2+count*2] = byte(new_value >> 8)
			data_buffer[3+count*2] = byte(new_value & 0xFF)
		}

		for count := (1 + repeat_count); count < 16; count++ { // Fill in the remaining data with old_value
			data_buffer[0+count*2] = byte(old_value >> 8)
			data_buffer[1+count*2] = byte(old_value & 0xFF)
		}

		for count := 0; count < 8; count++ { // Invert bits in every other slot
			data_buffer[2+count*4] ^= 0xFF
			data_buffer[3+count*4] ^= 0xFF
		}

		for line_count := 0; line_count < 2; line_count++ { // print the results on two lines
			for byte_count := 0; byte_count < 0x10; byte_count++ {
				fmt.Printf("0x%x ", data_buffer[byte_count+line_count*0x10])
			}
			fmt.Print("\n")
		}
	},
}

func init() {
	cmCmd.AddCommand(odometerCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// odometerCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// odometerCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
