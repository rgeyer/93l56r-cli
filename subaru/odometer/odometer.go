package odometer

import "bytes"

var encode_table = [...]byte{0x00, 0x07, 0x0C, 0x0B, 0x06, 0x01, 0x0A, 0x0D, 0x03, 0x04, 0x0F, 0x08, 0x05, 0x02, 0x09, 0x0E}

func Encode(mileage int) []byte {
	repeat_count := mileage & 0x0F // Take lowest nibble and use it as the repeat count
	new_value := mileage >> 4      // The new count has that lowest nibble shaved off
	old_value := new_value - 1     // The old count will be new count - 1

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

	return data_buffer
}

func Decode(encoded_buffer []byte) int {
	// Commence (de)coding
	for count := 0; count < 8; count++ {
		encoded_buffer[2+count*4] ^= 0xFF
		encoded_buffer[3+count*4] ^= 0xFF
	}

	// Count the number of times the two byte pair occurs
	repeat_count := 0
	new_value := int(encoded_buffer[0]) << 8
	new_value = new_value | int(encoded_buffer[1]&0xFF)
	old_value := new_value
	for count := 0; count < 16; count++ {
		new_value = int(encoded_buffer[count*2]) << 8
		new_value = new_value | int(encoded_buffer[count*2+1]&0xFF)

		if new_value != old_value {
			break
		}

		old_value = new_value
		repeat_count++
	}

	// (de)Code the last nibble
	last_nibble := old_value & 0x000F
	old_value = (old_value & 0xFFF0) + int(bytes.IndexByte(encode_table[:], byte(last_nibble)))
	final_value := (old_value << 4) + repeat_count - 1

	return final_value
}
